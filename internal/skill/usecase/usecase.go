package usecase

import (
	"context"

	"kerjadekat/backend/internal/domain"
)

type Skills struct {
	repo domain.SkillRepository
}

func NewSkills(repo domain.SkillRepository) *Skills {
	return &Skills{repo: repo}
}

func (s *Skills) ListCategories(ctx context.Context) ([]domain.SkillCategory, error) {
	return s.repo.ListCategories(ctx)
}
