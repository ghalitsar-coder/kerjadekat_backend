package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidOTP        = errors.New("invalid or expired otp")
	ErrInvalidTransition = errors.New("invalid order status transition")
	ErrPaymentFailed     = errors.New("payment operation failed")
)
