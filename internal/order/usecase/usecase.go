package usecase

import (
	"context"
	"fmt"
	"time"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
)

type Orders struct {
	orders       domain.OrderRepository
	skills       domain.SkillRepository
	pay          domain.PaymentGateway
	presence     domain.WorkerPresence
	matchPub     domain.OrderMatchPublisher
	offerTracker domain.OrderOfferTracker
	workers      domain.WorkerRepository
	wallets      domain.WalletRepository
	match        MatchSettings
}

// MatchSettings tunes Redis GEO radius and RabbitMQ re-broadcast rounds.
type MatchSettings struct {
	RadiusMeters float64
	TimerDelayMs int
	MaxRounds    int
}

func NewOrders(
	orders domain.OrderRepository,
	skills domain.SkillRepository,
	pay domain.PaymentGateway,
	presence domain.WorkerPresence,
	matchPub domain.OrderMatchPublisher,
	offerTracker domain.OrderOfferTracker,
	workers domain.WorkerRepository,
	wallets domain.WalletRepository,
	match MatchSettings,
) *Orders {
	if match.RadiusMeters <= 0 {
		match.RadiusMeters = defaultMatchRadiusM
	}
	if match.TimerDelayMs <= 0 {
		match.TimerDelayMs = matchTimerDelayMs
	}
	if match.MaxRounds <= 0 {
		match.MaxRounds = maxMatchRounds
	}
	return &Orders{
		orders:       orders,
		skills:       skills,
		pay:          pay,
		presence:     presence,
		matchPub:     matchPub,
		offerTracker: offerTracker,
		workers:      workers,
		wallets:      wallets,
		match:        match,
	}
}

func (o *Orders) appendLog(ctx context.Context, r domain.OrderRepository, orderID uuid.UUID, from *string, to string, actor *uuid.UUID, note *string) error {
	log := &domain.OrderStatusLog{
		ID:         uuid.New(),
		OrderID:    orderID,
		FromStatus: from,
		ToStatus:   to,
		ChangedBy:  actor,
		Note:       note,
	}
	return r.AppendStatusLog(ctx, log)
}

func (o *Orders) Get(ctx context.Context, id uuid.UUID, actorID uuid.UUID, actorRole string) (*domain.Order, error) {
	ord, err := o.orders.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canViewOrder(ord, actorID, actorRole) {
		return nil, domain.ErrForbidden
	}
	return ord, nil
}

func canViewOrder(o *domain.Order, actor uuid.UUID, role string) bool {
	if o.ConsumerID == actor {
		return true
	}
	if o.WorkerID != nil && *o.WorkerID == actor {
		return true
	}
	if role == domain.RoleAdmin {
		return true
	}
	return false
}

type XenditWebhookPayload struct {
	Event string `json:"event"`
	Data  struct {
		PaymentRequestID string `json:"payment_request_id"`
		ReferenceID      string `json:"reference_id"`
		Status           string `json:"status"`
	} `json:"data"`
}

func (o *Orders) HandleXenditWebhook(ctx context.Context, p XenditWebhookPayload) error {
	id := p.Data.PaymentRequestID
	if id == "" {
		return fmt.Errorf("missing payment_request_id in webhook")
	}

	ord, err := o.orders.FindByInvoiceID(ctx, id)
	if err != nil {
		return err
	}

	prev := ord.PaymentStatus
	note := fmt.Sprintf("Xendit v3 webhook: %s (event: %s)", p.Data.Status, p.Event)

	switch p.Data.Status {
	case "SUCCEEDED":
		ord.PaymentStatus = domain.PaymentCaptured
	case "AUTHORIZED":
		ord.PaymentStatus = domain.PaymentAuthorized
	case "EXPIRED", "FAILED", "CANCELED":
		ord.PaymentStatus = domain.PaymentFailed
	default:
		return nil
	}

	if prev == ord.PaymentStatus {
		return nil // no change needed
	}

	return o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		return o.appendLog(ctx, r, ord.ID, &prev, ord.Status, nil, &note)
	})
}

func (o *Orders) ListMine(ctx context.Context, userID uuid.UUID, role string, limit, offset int) ([]domain.Order, error) {
	switch role {
	case domain.RoleConsumer:
		return o.orders.ListByConsumer(ctx, userID, limit, offset)
	case domain.RoleWorker:
		return o.orders.ListByWorker(ctx, userID, limit, offset)
	default:
		return nil, domain.ErrForbidden
	}
}

type AcceptOrderInput struct {
	WorkerUserID uuid.UUID
	AgreedRate   *float64
}

func (o *Orders) publishStatus(ctx context.Context, ord *domain.Order, actorUserID *uuid.UUID) {
	o.matchPub.PublishOrderStatus(ctx, domain.OrderStatusEvent{
		OrderID:     ord.ID,
		NewStatus:   ord.Status,
		ActorUserID: actorUserID,
		ConsumerID:  ord.ConsumerID,
		WorkerID:    ord.WorkerID,
	})
}

func (o *Orders) Accept(ctx context.Context, orderID uuid.UUID, in AcceptOrderInput) (*domain.Order, error) {
	err := o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		ord, err := r.FindByID(ctx, orderID)
		if err != nil {
			return err
		}
		if ord.Status != domain.OrderStatusPendingMatch && ord.Status != domain.OrderStatusOffered {
			return domain.ErrInvalidTransition
		}
		if ord.WorkerID != nil {
			return domain.ErrConflict
		}
		if ord.ConsumerID == in.WorkerUserID {
			return domain.ErrForbidden
		}
		prev := ord.Status
		wid := in.WorkerUserID
		ord.WorkerID = &wid
		ord.Status = domain.OrderStatusAccepted
		ord.AgreedRate = in.AgreedRate
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		if err := o.appendLog(ctx, r, ord.ID, &prev, domain.OrderStatusAccepted, &in.WorkerUserID, nil); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	ord, err := o.orders.FindByID(ctx, orderID)
	if err == nil {
		uid := in.WorkerUserID
		o.publishStatus(ctx, ord, &uid)
	}
	return ord, err
}

func (o *Orders) Start(ctx context.Context, orderID uuid.UUID, workerUserID uuid.UUID) (*domain.Order, error) {
	err := o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		ord, err := r.FindByID(ctx, orderID)
		if err != nil {
			return err
		}
		if ord.WorkerID == nil || *ord.WorkerID != workerUserID {
			return domain.ErrForbidden
		}
		if ord.Status != domain.OrderStatusAccepted {
			return domain.ErrInvalidTransition
		}
		prev := ord.Status
		now := time.Now()
		ord.Status = domain.OrderStatusInProgress
		ord.StartedAt = &now
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		return o.appendLog(ctx, r, ord.ID, &prev, domain.OrderStatusInProgress, &workerUserID, nil)
	})
	if err != nil {
		return nil, err
	}
	ord, err := o.orders.FindByID(ctx, orderID)
	if err == nil {
		o.publishStatus(ctx, ord, &workerUserID)
	}
	return ord, err
}

func (o *Orders) Complete(ctx context.Context, orderID uuid.UUID, actor uuid.UUID) (*domain.Order, error) {
	ord, err := o.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if ord.ConsumerID != actor && (ord.WorkerID == nil || *ord.WorkerID != actor) {
		return nil, domain.ErrForbidden
	}
	if ord.Status != domain.OrderStatusInProgress {
		return nil, domain.ErrInvalidTransition
	}
	if ord.FeeAuthID == nil || *ord.FeeAuthID == "" {
		return nil, domain.ErrPaymentFailed
	}
	if err := o.pay.Capture(ctx, domain.CaptureRequest{AuthID: *ord.FeeAuthID}); err != nil {
		return nil, err
	}

	now := time.Now()
	prev := ord.Status
	ord.Status = domain.OrderStatusCompleted
	ord.PaymentStatus = domain.PaymentCaptured
	ord.CompletedAt = &now

	if err := o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		if err := o.appendLog(ctx, r, ord.ID, &prev, domain.OrderStatusCompleted, &actor, nil); err != nil {
			return err
		}
		if ord.WorkerID != nil && ord.AgreedRate != nil && *ord.AgreedRate > 0 {
			rec := &domain.IncomeRecord{
				ID:       uuid.New(),
				WorkerID: *ord.WorkerID,
				OrderID:  &ord.ID,
				Amount:   *ord.AgreedRate,
				Source:   domain.IncomeSourcePlatform,
				Verified: true,
			}
			if err := r.CreateIncomeRecord(ctx, rec); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if ord.WorkerID != nil && ord.AgreedRate != nil && *ord.AgreedRate > 0 {
		o.creditWallet(ctx, *ord.WorkerID, ord.ID, *ord.AgreedRate)
	}
	ord, err = o.orders.FindByID(ctx, orderID)
	if err == nil {
		o.publishStatus(ctx, ord, &actor)
	}
	return ord, err
}

func (o *Orders) Cancel(ctx context.Context, orderID uuid.UUID, actor uuid.UUID, reason *string) (*domain.Order, error) {
	ord, err := o.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var next string
	switch {
	case ord.ConsumerID == actor:
		next = domain.OrderStatusCancelledConsumer
	case ord.WorkerID != nil && *ord.WorkerID == actor:
		next = domain.OrderStatusCancelledWorker
	default:
		return nil, domain.ErrForbidden
	}

	if !isCancellable(ord.Status) {
		return nil, domain.ErrInvalidTransition
	}

	if ord.PaymentStatus == domain.PaymentAuthorized && ord.FeeAuthID != nil && *ord.FeeAuthID != "" {
		if err := o.pay.Void(ctx, domain.VoidRequest{AuthID: *ord.FeeAuthID}); err != nil {
			return nil, err
		}
		ord.PaymentStatus = domain.PaymentRefunded
	}

	prev := ord.Status
	ord.Status = next
	ord.CancelledReason = reason

	if err := o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		return o.appendLog(ctx, r, ord.ID, &prev, next, &actor, reason)
	}); err != nil {
		return nil, err
	}
	ord, err = o.orders.FindByID(ctx, orderID)
	if err == nil {
		o.publishStatus(ctx, ord, &actor)
	}
	return ord, err
}

func isCancellable(status string) bool {
	switch status {
	case domain.OrderStatusPendingMatch, domain.OrderStatusOffered, domain.OrderStatusAccepted, domain.OrderStatusInProgress:
		return true
	default:
		return false
	}
}

func (o *Orders) Reject(ctx context.Context, orderID uuid.UUID, workerUserID uuid.UUID) error {
	ord, err := o.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if ord.Status != domain.OrderStatusPendingMatch && ord.Status != domain.OrderStatusOffered {
		return domain.ErrInvalidTransition
	}
	if err := o.offerTracker.RecordOffered(ctx, orderID, []uuid.UUID{workerUserID}); err != nil {
		return err
	}
	return o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		return o.appendLog(ctx, r, ord.ID, &ord.Status, ord.Status, &workerUserID, nil)
	})
}

func (o *Orders) ConfirmPayment(ctx context.Context, orderID uuid.UUID, consumerUserID uuid.UUID) (*domain.Order, error) {
	ord, err := o.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if ord.ConsumerID != consumerUserID {
		return nil, domain.ErrForbidden
	}
	if ord.Status != domain.OrderStatusCompleted {
		return nil, domain.ErrInvalidTransition
	}
	prev := ord.PaymentStatus
	ord.PaymentStatus = domain.PaymentConfirmed

	if err := o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		return o.appendLog(ctx, r, ord.ID, &prev, ord.Status, &consumerUserID, nil)
	}); err != nil {
		return nil, err
	}
	return o.orders.FindByID(ctx, orderID)
}

type RateOrderInput struct {
	Score   int16
	Comment *string
}

func (o *Orders) creditWallet(ctx context.Context, workerID uuid.UUID, orderID uuid.UUID, amount float64) {
	wallet, err := o.wallets.FindByUserID(ctx, workerID)
	if err == domain.ErrNotFound {
		wallet = &domain.Wallet{ID: uuid.New(), UserID: workerID, Balance: 0}
		if err := o.wallets.CreateWallet(ctx, wallet); err != nil {
			return
		}
	} else if err != nil {
		return
	}
	before := wallet.Balance
	after := before + amount
	oid := orderID.String()
	_ = o.wallets.CreateTransaction(ctx, &domain.WalletTransaction{
		ID:            uuid.New(),
		WalletID:      wallet.ID,
		Type:          domain.WalletTxTypeCredit,
		Amount:        amount,
		BalanceBefore: before,
		BalanceAfter:  after,
		ReferenceType: "order_completion",
		ReferenceID:   &oid,
	})
	_ = o.wallets.UpdateBalance(ctx, wallet.ID, after)
}

func (o *Orders) GetWallet(ctx context.Context, userID uuid.UUID) (*domain.Wallet, error) {
	return o.wallets.FindByUserID(ctx, userID)
}

const defaultTxLimit = 20

func (o *Orders) ListWalletTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.WalletTransaction, error) {
	wallet, err := o.wallets.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return o.wallets.ListTransactions(ctx, wallet.ID, limit, offset)
}

func (o *Orders) Rate(ctx context.Context, orderID uuid.UUID, consumerUserID uuid.UUID, in RateOrderInput) (*domain.Order, error) {
	ord, err := o.orders.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if ord.ConsumerID != consumerUserID {
		return nil, domain.ErrForbidden
	}
	if ord.Status != domain.OrderStatusCompleted {
		return nil, domain.ErrInvalidTransition
	}
	if ord.Rating != nil {
		return nil, domain.ErrConflict
	}
	if in.Score < 1 || in.Score > 5 {
		return nil, domain.ErrInvalidInput
	}
	if ord.WorkerID == nil {
		return nil, domain.ErrInvalidInput
	}

	rating := &domain.OrderRating{
		ID:       uuid.New(),
		OrderID:  orderID,
		GivenBy:  consumerUserID,
		GivenTo:  *ord.WorkerID,
		Score:    in.Score,
		Comment:  in.Comment,
	}

	if err := o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.CreateRating(ctx, rating); err != nil {
			return err
		}
		prof, err := o.workers.FindProfileByUserID(ctx, *ord.WorkerID)
		if err != nil {
			return err
		}
		newCount := prof.RatingCount + 1
		newAvg := ((prof.RatingAvg * float64(prof.RatingCount)) + float64(in.Score)) / float64(newCount)
		prof.RatingAvg = newAvg
		prof.RatingCount = newCount
		return o.workers.UpdateProfile(ctx, prof)
	}); err != nil {
		return nil, err
	}
	return o.orders.FindByID(ctx, orderID)
}
