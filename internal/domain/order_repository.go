package domain

import (
	"context"

	"github.com/google/uuid"
)

// OrderRepository persists orders, status logs, ratings, and income side-effects.
type OrderRepository interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, r OrderRepository) error) error

	Create(ctx context.Context, o *Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	FindByInvoiceID(ctx context.Context, invoiceID string) (*Order, error)
	ListByConsumer(ctx context.Context, consumerID uuid.UUID, limit, offset int) ([]Order, error)
	ListByWorker(ctx context.Context, workerUserID uuid.UUID, limit, offset int) ([]Order, error)
	Update(ctx context.Context, o *Order) error

	AppendStatusLog(ctx context.Context, log *OrderStatusLog) error
	CreateIncomeRecord(ctx context.Context, rec *IncomeRecord) error
}
