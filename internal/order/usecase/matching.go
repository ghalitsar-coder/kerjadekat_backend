package usecase

import (
	"context"
	"fmt"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
)

const (
	defaultMatchRadiusM = 5000
	matchTimerDelayMs   = 60_000
	maxMatchRounds      = 5
)

type CreateOrderInput struct {
	ConsumerID       uuid.UUID
	SkillID          int
	Description      *string
	Latitude         float64
	Longitude        float64
	ConsumerAddress  *string
	PaymentMethodFee *string
}

func (o *Orders) Create(ctx context.Context, in CreateOrderInput) (*domain.Order, error) {
	if in.Latitude < -90 || in.Latitude > 90 || in.Longitude < -180 || in.Longitude > 180 {
		return nil, domain.ErrInvalidInput
	}
	if _, err := o.skills.FindCategoryByID(ctx, in.SkillID); err != nil {
		return nil, err
	}

	orderID := uuid.New()
	method := ""
	if in.PaymentMethodFee != nil {
		method = *in.PaymentMethodFee
	}
	authRes, err := o.pay.Authorize(ctx, domain.AuthorizeRequest{
		ReferenceID: orderID.String(),
		AmountIDR:   2000,
		Method:      method,
	})
	if err != nil {
		return nil, err
	}

	inv := authRes.InvoiceID
	authID := authRes.AuthID
	ord := &domain.Order{
		ID:               orderID,
		ConsumerID:       in.ConsumerID,
		SkillID:          in.SkillID,
		Status:           domain.OrderStatusPendingMatch,
		Description:      in.Description,
		ConsumerLocation: domain.NullPoint{Lat: in.Latitude, Lng: in.Longitude, Valid: true},
		ConsumerAddress:  in.ConsumerAddress,
		PlatformFee:      2000,
		PaymentMethodFee: in.PaymentMethodFee,
		XenditInvoiceID:  &inv,
		FeeAuthID:        &authID,
		PaymentStatus:    domain.PaymentAuthorized,
	}

	if err := o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.Create(ctx, ord); err != nil {
			return err
		}
		return o.appendLog(ctx, r, ord.ID, nil, domain.OrderStatusPendingMatch, &in.ConsumerID, nil)
	}); err != nil {
		return nil, err
	}

	if err := o.broadcastMatchRound(ctx, ord, 1); err != nil {
		return nil, fmt.Errorf("match broadcast: %w", err)
	}

	return o.orders.FindByID(ctx, ord.ID)
}

func (o *Orders) broadcastMatchRound(ctx context.Context, ord *domain.Order, round int) error {
	if o.presence == nil || o.matchPub == nil {
		return fmt.Errorf("matching engine not configured")
	}
	if !ord.ConsumerLocation.Valid {
		return domain.ErrInvalidInput
	}
	if !isMatchableStatus(ord.Status) {
		return nil
	}

	nearby, err := o.presence.NearbyWorkerUserIDs(
		ctx,
		ord.SkillID,
		ord.ConsumerLocation.Lat,
		ord.ConsumerLocation.Lng,
		o.match.RadiusMeters,
	)
	if err != nil {
		return err
	}

	eligible := nearby
	if o.offerTracker != nil && len(nearby) > 0 {
		eligible, err = o.offerTracker.FilterNotYetOffered(ctx, ord.ID, nearby)
		if err != nil {
			return fmt.Errorf("filter offered workers: %w", err)
		}
	}

	if len(eligible) > 0 {
		if err := o.matchPub.PublishNewOrder(ctx, domain.OrderMatchEvent{
			OrderID:         ord.ID,
			SkillID:         ord.SkillID,
			Latitude:        ord.ConsumerLocation.Lat,
			Longitude:       ord.ConsumerLocation.Lng,
			WorkerUserIDs:   eligible,
			PlatformFee:     ord.PlatformFee,
			ConsumerAddress: ord.ConsumerAddress,
		}); err != nil {
			return err
		}

		if o.offerTracker != nil {
			if err := o.offerTracker.RecordOffered(ctx, ord.ID, eligible); err != nil {
				return fmt.Errorf("record offered workers: %w", err)
			}
		}

		if ord.Status == domain.OrderStatusPendingMatch {
			if err := o.promoteToOffered(ctx, ord); err != nil {
				return err
			}
		}
	}

	return o.matchPub.ScheduleMatchTimer(ctx, ord.ID, round, o.match.TimerDelayMs)
}

func (o *Orders) promoteToOffered(ctx context.Context, ord *domain.Order) error {
	prev := ord.Status
	ord.Status = domain.OrderStatusOffered
	return o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		return o.appendLog(ctx, r, ord.ID, &prev, domain.OrderStatusOffered, nil, nil)
	})
}

// HandleMatchTimer is invoked by the RabbitMQ consumer when an offer window expires.
func (o *Orders) HandleMatchTimer(ctx context.Context, orderID uuid.UUID, round int) error {
	ord, err := o.orders.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if !isMatchableStatus(ord.Status) {
		return nil
	}

	if round >= o.match.MaxRounds {
		return o.expireOrder(ctx, ord)
	}

	return o.broadcastMatchRound(ctx, ord, round+1)
}

func isMatchableStatus(status string) bool {
	return status == domain.OrderStatusPendingMatch || status == domain.OrderStatusOffered
}

func (o *Orders) expireOrder(ctx context.Context, ord *domain.Order) error {
	if ord.PaymentStatus == domain.PaymentAuthorized && ord.FeeAuthID != nil && *ord.FeeAuthID != "" {
		if err := o.pay.Void(ctx, domain.VoidRequest{AuthID: *ord.FeeAuthID}); err != nil {
			return err
		}
		ord.PaymentStatus = domain.PaymentRefunded
	}
	prev := ord.Status
	reason := "no worker accepted within matching window"
	ord.Status = domain.OrderStatusExpired
	ord.CancelledReason = &reason

	return o.orders.WithTx(ctx, func(ctx context.Context, r domain.OrderRepository) error {
		if err := r.Update(ctx, ord); err != nil {
			return err
		}
		return o.appendLog(ctx, r, ord.ID, &prev, domain.OrderStatusExpired, nil, &reason)
	})
}
