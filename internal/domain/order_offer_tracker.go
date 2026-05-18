package domain

import (
	"context"

	"github.com/google/uuid"
)

// OrderOfferTracker records which workers were already pinged for an order (anti-spam).
type OrderOfferTracker interface {
	FilterNotYetOffered(ctx context.Context, orderID uuid.UUID, candidates []uuid.UUID) ([]uuid.UUID, error)
	RecordOffered(ctx context.Context, orderID uuid.UUID, workerIDs []uuid.UUID) error
}
