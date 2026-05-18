package domain

import (
	"context"

	"github.com/google/uuid"
)

// WorkerPresence tracks live worker GPS and skills in Redis (not PostgreSQL).
type WorkerPresence interface {
	UpdateWorkerLocation(ctx context.Context, userID uuid.UUID, lat, lng float64, skillIDs []int) error
	SetAvailability(ctx context.Context, userID uuid.UUID, availability string) error
	NearbyWorkerUserIDs(ctx context.Context, skillID int, lat, lng float64, radiusMeters float64) ([]uuid.UUID, error)
}
