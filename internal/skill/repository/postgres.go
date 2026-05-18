package repository

import (
	"context"
	"errors"

	"kerjadekat/backend/internal/domain"

	"gorm.io/gorm"
)

type SkillPostgres struct {
	db *gorm.DB
}

func NewSkillPostgres(db *gorm.DB) *SkillPostgres {
	return &SkillPostgres{db: db}
}

func (r *SkillPostgres) ListCategories(ctx context.Context) ([]domain.SkillCategory, error) {
	var rows []domain.SkillCategory
	if err := r.db.WithContext(ctx).Order("name asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *SkillPostgres) FindCategoryByID(ctx context.Context, id int) (*domain.SkillCategory, error) {
	var s domain.SkillCategory
	if err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

var _ domain.SkillRepository = (*SkillPostgres)(nil)
