package domain

import (
	"context"

	"github.com/google/uuid"
)

// OrderMatchPublisher emits realtime order events (Redis Pub/Sub) and schedules match timers.
type OrderMatchPublisher interface {
	PublishNewOrder(ctx context.Context, evt OrderMatchEvent) error
	PublishOrderStatus(ctx context.Context, evt OrderStatusEvent) error
	ScheduleMatchTimer(ctx context.Context, orderID uuid.UUID, round int, delayMs int) error
}

// OrderStatusEvent is published when an order transitions status.
type OrderStatusEvent struct {
	OrderID     uuid.UUID
	NewStatus   string
	ActorUserID *uuid.UUID
	ConsumerID  uuid.UUID
	WorkerID    *uuid.UUID
}

// OrderMatchEvent is fanned out to localized worker WebSocket connections.
type OrderMatchEvent struct {
	OrderID         uuid.UUID
	SkillID         int
	Latitude        float64
	Longitude       float64
	WorkerUserIDs   []uuid.UUID
	PlatformFee     float64
	ConsumerAddress *string
}
