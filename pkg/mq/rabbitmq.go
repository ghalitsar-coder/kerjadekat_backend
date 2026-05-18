package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
)

const (
	ExchangeDelayed = "kerjadekat.delayed"
	QueueMatchTimer = "kerjadekat.order.match.timer"
	RoutingMatch    = "order.match.timer"
)

// MatchTimerMessage is delivered after the offer window (60s) expires.
type MatchTimerMessage struct {
	OrderID uuid.UUID `json:"order_id"`
	Round   int       `json:"round"`
}

// Client wraps RabbitMQ delayed-message exchange publishing and consumption.
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func Connect(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("amqp dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("amqp channel: %w", err)
	}
	c := &Client{conn: conn, channel: ch}
	if err := c.declareTopology(); err != nil {
		_ = c.Close()
		return nil, err
	}
	return c, nil
}

func (c *Client) declareTopology() error {
	args := amqp.Table{
		"x-delayed-type": "direct",
	}
	if err := c.channel.ExchangeDeclare(
		ExchangeDelayed,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		return fmt.Errorf("declare delayed exchange: %w", err)
	}
	if _, err := c.channel.QueueDeclare(
		QueueMatchTimer,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}
	if err := c.channel.QueueBind(
		QueueMatchTimer,
		RoutingMatch,
		ExchangeDelayed,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}
	return nil
}

// PublishMatchTimer schedules a delayed match-expiry check (delayMs, typically 60000).
func (c *Client) PublishMatchTimer(ctx context.Context, orderID uuid.UUID, round int, delayMs int) error {
	body, err := json.Marshal(MatchTimerMessage{OrderID: orderID, Round: round})
	if err != nil {
		return err
	}
	return c.channel.PublishWithContext(ctx,
		ExchangeDelayed,
		RoutingMatch,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Headers: amqp.Table{
				"x-delay": int32(delayMs),
			},
			Body: body,
		},
	)
}

// ConsumeMatchTimer registers a consumer on the delayed match queue.
func (c *Client) ConsumeMatchTimer(ctx context.Context, handler func(ctx context.Context, msg MatchTimerMessage) error) error {
	deliveries, err := c.channel.Consume(
		QueueMatchTimer,
		"kerjadekat-api",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					return
				}
				var m MatchTimerMessage
				if err := json.Unmarshal(d.Body, &m); err != nil {
					_ = d.Nack(false, false)
					continue
				}
				if err := handler(ctx, m); err != nil {
					_ = d.Nack(false, true)
					continue
				}
				_ = d.Ack(false)
			}
		}
	}()
	return nil
}

func (c *Client) Close() error {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
