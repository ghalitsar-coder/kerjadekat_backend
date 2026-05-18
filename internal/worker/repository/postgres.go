package repository

import (
	"context"
	"errors"

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

var _ domain.WorkerRepository = (*WorkerPostgres)(nil)
