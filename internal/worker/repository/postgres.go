package repository

import (
	"context"
	"errors"
	"fmt"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkerPostgres struct {
	db *gorm.DB
}

func NewWorkerPostgres(db *gorm.DB) *WorkerPostgres {
	return &WorkerPostgres{db: db}
}

func (r *WorkerPostgres) CreateProfile(ctx context.Context, p *domain.WorkerProfile) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *WorkerPostgres) FindProfileByUserID(ctx context.Context, userID uuid.UUID) (*domain.WorkerProfile, error) {
	var p domain.WorkerProfile
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Skills.Skill").
		First(&p, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *WorkerPostgres) UpdateProfile(ctx context.Context, p *domain.WorkerProfile) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: false}).Save(p).Error
}

func (r *WorkerPostgres) ListOnlineWithLocation(ctx context.Context) ([]domain.WorkerProfile, error) {
	var rows []domain.WorkerProfile
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Skills.Skill").
		Where("availability = ?", domain.WorkerAvailabilityOnline).
		Where("last_location IS NOT NULL").
		Find(&rows).Error
	return rows, err
}

func (r *WorkerPostgres) FindNearbyOnline(ctx context.Context, lat, lng, radiusMeters float64, skillID *int) ([]domain.WorkerProfile, error) {
	pointWKT := fmt.Sprintf("SRID=4326;POINT(%f %f)", lng, lat)
	q := r.db.WithContext(ctx).
		Preload("User").
		Preload("Skills.Skill").
		Joins("JOIN worker_skills ON worker_skills.worker_id = worker_profiles.id").
		Where("worker_profiles.availability = ?", domain.WorkerAvailabilityOnline).
		Where("worker_profiles.last_location IS NOT NULL").
		Where("ST_DWithin(worker_profiles.last_location, ST_GeographyFromText(?), ?)", pointWKT, radiusMeters)
	if skillID != nil && *skillID > 0 {
		q = q.Where("worker_skills.skill_id = ?", *skillID)
	}
	var rows []domain.WorkerProfile
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	seen := make(map[uuid.UUID]struct{}, len(rows))
	out := make([]domain.WorkerProfile, 0, len(rows))
	for _, row := range rows {
		if _, ok := seen[row.ID]; ok {
			continue
		}
		seen[row.ID] = struct{}{}
		out = append(out, row)
	}
	return out, nil
}

var _ domain.WorkerRepository = (*WorkerPostgres)(nil)
