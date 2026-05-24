package usecase

import (
	"context"

	"kerjadekat/backend/internal/domain"
)

type Kelurahans struct {
	repo domain.KelurahanRepository
}

func NewKelurahans(repo domain.KelurahanRepository) *Kelurahans {
	return &Kelurahans{repo: repo}
}

func (k *Kelurahans) ListAll(ctx context.Context) ([]domain.Kelurahan, error) {
	return k.repo.ListAll(ctx)
}
