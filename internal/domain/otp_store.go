package domain

import "context"

// OTPStore persists one-time passwords in a shared store (Redis) for multi-pod auth.
type OTPStore interface {
	Store(ctx context.Context, phone, code string) error
	VerifyAndConsume(ctx context.Context, phone, code string) (ok bool, err error)
}
