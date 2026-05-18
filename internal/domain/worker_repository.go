package domain

import (
	"context"

	"github.com/google/uuid"
)

// WorkerRepository persists worker profiles and related worker data.
type WorkerRepository interface {
	CreateProfile(ctx context.Context, p *WorkerProfile) error
	FindProfileByUserID(ctx context.Context, userID uuid.UUID) (*WorkerProfile, error)
	UpdateProfile(ctx context.Context, p *WorkerProfile) error
}
