package usecase

import (
	"context"
	// "time"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
)

type Users struct {
	repo domain.UserRepository
}

func NewUsers(repo domain.UserRepository) *Users {
	return &Users{repo: repo}
}

func (u *Users) Me(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.repo.FindByID(ctx, id)
}
