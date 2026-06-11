package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const ChannelOrdersNew = "orders:new"
const ChannelOrdersStatus = "orders:status"

// NewOrderAvailable is published to Redis for the WebSocket hub.
type NewOrderAvailable struct {
	Type            string      `json:"type"`
	OrderID         uuid.UUID   `json:"order_id"`
	SkillID         int         `json:"skill_id"`
	Latitude        float64     `json:"latitude"`
	Longitude       float64     `json:"longitude"`
	WorkerUserIDs   []uuid.UUID `json:"worker_user_ids"`
	PlatformFee     float64     `json:"platform_fee"`
	ConsumerAddress *string     `json:"consumer_address,omitempty"`
}

// OrderStatusChanged is published when an order's status changes.
type OrderStatusChanged struct {
	Type        string     `json:"type"`
	OrderID     uuid.UUID  `json:"order_id"`
	NewStatus   string     `json:"new_status"`
	ActorUserID *uuid.UUID `json:"actor_user_id,omitempty"`
	ConsumerID  uuid.UUID  `json:"consumer_id"`
	WorkerID    *uuid.UUID `json:"worker_id,omitempty"`
}

// Publisher emits order realtime events on Redis Pub/Sub.
type Publisher struct {
	rdb *redis.Client
}

func NewPublisher(rdb *redis.Client) *Publisher {
	return &Publisher{rdb: rdb}
}

func (p *Publisher) PublishNewOrder(ctx context.Context, evt NewOrderAvailable) error {
	evt.Type = "new_order_available"
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	if err := p.rdb.Publish(ctx, ChannelOrdersNew, b).Err(); err != nil {
		return fmt.Errorf("redis publish: %w", err)
	}
	return nil
}

func (p *Publisher) PublishOrderStatus(ctx context.Context, evt OrderStatusChanged) error {
	evt.Type = "order_status_changed"
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	if err := p.rdb.Publish(ctx, ChannelOrdersStatus, b).Err(); err != nil {
		return fmt.Errorf("redis publish status: %w", err)
	}
	return nil
}
