package domain

import "context"

// AuthorizeRequest holds a platform-fee authorization attempt (Xendit-shaped stub).
type AuthorizeRequest struct {
	ReferenceID string
	AmountIDR   float64
	Method      string // e.g. qris, va_bca, ewallet
}

// AuthorizeResult returns gateway identifiers to persist on orders.
type AuthorizeResult struct {
	InvoiceID  string
	AuthID     string
	PaymentURL string // QR code image URL or deeplink for the customer
}

// CaptureRequest finalizes an authorized hold.
type CaptureRequest struct {
	AuthID string
}

// VoidRequest cancels an authorization without capture.
type VoidRequest struct {
	AuthID string
}

// PaymentGateway abstracts Xendit (authorize/capture/void).
type PaymentGateway interface {
	Authorize(ctx context.Context, in AuthorizeRequest) (AuthorizeResult, error)
	Capture(ctx context.Context, in CaptureRequest) error
	Void(ctx context.Context, in VoidRequest) error
}
