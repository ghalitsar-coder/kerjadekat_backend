package usecase

import (
	"context"
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
	return o.orders.FindByID(ctx, orderID)
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
	return o.orders.FindByID(ctx, orderID)
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
	return o.orders.FindByID(ctx, orderID)
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
	return o.orders.FindByID(ctx, orderID)
}

func isCancellable(status string) bool {
	switch status {
	case domain.OrderStatusPendingMatch, domain.OrderStatusOffered, domain.OrderStatusAccepted, domain.OrderStatusInProgress:
		return true
	default:
		return false
	}
}
