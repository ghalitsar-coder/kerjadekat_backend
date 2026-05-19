package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository persists users.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByProviderID(ctx context.Context, provider, providerID string) (*User, error)
}
