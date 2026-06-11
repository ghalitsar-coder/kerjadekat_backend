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
	ListOnlineWithLocation(ctx context.Context) ([]WorkerProfile, error)
	FindNearbyOnline(ctx context.Context, lat, lng, radiusMeters float64, skillID *int) ([]WorkerProfile, error)
	FindBySkillIDs(ctx context.Context, skillIDs []int) ([]WorkerProfile, error)
}
