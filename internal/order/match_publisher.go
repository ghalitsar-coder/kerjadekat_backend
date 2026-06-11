package order

import (
	"context"
	"fmt"

	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/mq"

	"github.com/google/uuid"
)

// MatchPublisher bridges domain ports to Redis Pub/Sub and RabbitMQ delayed messages.
type MatchPublisher struct {
	redis *Publisher
	mq    *mq.Client
}

func NewMatchPublisher(redisPub *Publisher, mqClient *mq.Client) *MatchPublisher {
	return &MatchPublisher{redis: redisPub, mq: mqClient}
}

func (p *MatchPublisher) PublishNewOrder(ctx context.Context, evt domain.OrderMatchEvent) error {
	if p.redis == nil {
		return fmt.Errorf("redis publisher not configured")
	}
	return p.redis.PublishNewOrder(ctx, NewOrderAvailable{
		OrderID:         evt.OrderID,
		SkillID:         evt.SkillID,
		Latitude:        evt.Latitude,
		Longitude:       evt.Longitude,
		WorkerUserIDs:   evt.WorkerUserIDs,
		PlatformFee:     evt.PlatformFee,
		ConsumerAddress: evt.ConsumerAddress,
	})
}

func (p *MatchPublisher) PublishOrderStatus(ctx context.Context, evt domain.OrderStatusEvent) error {
	if p.redis == nil {
		return fmt.Errorf("redis publisher not configured")
	}
	return p.redis.PublishOrderStatus(ctx, OrderStatusChanged{
		OrderID:     evt.OrderID,
		NewStatus:   evt.NewStatus,
		ActorUserID: evt.ActorUserID,
		ConsumerID:  evt.ConsumerID,
		WorkerID:    evt.WorkerID,
	})
}

func (p *MatchPublisher) ScheduleMatchTimer(ctx context.Context, orderID uuid.UUID, round int, delayMs int) error {
	if p.mq == nil {
		return fmt.Errorf("rabbitmq not configured")
	}
	return p.mq.PublishMatchTimer(ctx, orderID, round, delayMs)
}

var _ domain.OrderMatchPublisher = (*MatchPublisher)(nil)
