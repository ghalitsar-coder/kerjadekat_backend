package usecase

import (
	"context"
	"fmt"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
)

type Workers struct {
	repo     domain.WorkerRepository
	presence domain.WorkerPresence
}

func NewWorkers(repo domain.WorkerRepository, presence domain.WorkerPresence) *Workers {
	return &Workers{repo: repo, presence: presence}
}

func (w *Workers) Me(ctx context.Context, userID uuid.UUID) (*domain.WorkerProfile, error) {
	return w.repo.FindProfileByUserID(ctx, userID)
}

func (w *Workers) SetAvailability(ctx context.Context, userID uuid.UUID, availability string) error {
	switch availability {
	case domain.WorkerAvailabilityOnline, domain.WorkerAvailabilityOffline, domain.WorkerAvailabilityBusy:
	default:
		return domain.ErrInvalidInput
	}
	p, err := w.repo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	p.Availability = availability
	if err := w.repo.UpdateProfile(ctx, p); err != nil {
		return err
	}
	if w.presence != nil {
		return w.presence.SetAvailability(ctx, userID, availability)
	}
	return nil
}

// SetLocation stores live GPS in Redis GEO (5m TTL), not PostgreSQL.
func (w *Workers) SetLocation(ctx context.Context, userID uuid.UUID, lat, lng float64) error {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return domain.ErrInvalidInput
	}
	if w.presence == nil {
		return fmt.Errorf("worker presence store not configured")
	}
	p, err := w.repo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}
	skillIDs := make([]int, 0, len(p.Skills))
	for _, s := range p.Skills {
		skillIDs = append(skillIDs, s.SkillID)
	}
	return w.presence.UpdateWorkerLocation(ctx, userID, lat, lng, skillIDs)
}
