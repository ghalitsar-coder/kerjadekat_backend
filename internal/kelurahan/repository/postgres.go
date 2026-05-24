package repository

import (
	"context"

	"kerjadekat/backend/internal/domain"

	"gorm.io/gorm"
)

type KelurahanPostgres struct {
	db *gorm.DB
}

func NewKelurahanPostgres(db *gorm.DB) *KelurahanPostgres {
	return &KelurahanPostgres{db: db}
}

func (r *KelurahanPostgres) ListAll(ctx context.Context) ([]domain.Kelurahan, error) {
	var rows []domain.Kelurahan
	err := r.db.WithContext(ctx).
		Order("name asc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

var _ domain.KelurahanRepository = (*KelurahanPostgres)(nil)
