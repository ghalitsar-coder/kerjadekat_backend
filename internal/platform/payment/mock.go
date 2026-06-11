package payment

import (
	"context"
	"fmt"
	"sync/atomic"

	"kerjadekat/backend/internal/domain"
)

// MockGateway is a stand-in for Xendit authorize/capture/void flows.
type MockGateway struct {
	seq atomic.Uint64
}

func NewMockGateway() *MockGateway {
	return &MockGateway{}
}

func mockPaymentURL(method, refID string) string {
	switch method {
	case "qris":
		return fmt.Sprintf("https://api.qr-code-generator.com/v1/create?ref=%s&amount=2000", refID)
	case "gopay", "ovo", "dana":
		return fmt.Sprintf("https://mock.xendit.dev/pay/%s/%s", method, refID)
	default:
		return ""
	}
}

func (m *MockGateway) Authorize(ctx context.Context, in domain.AuthorizeRequest) (domain.AuthorizeResult, error) {
	_ = ctx
	n := m.seq.Add(1)
	return domain.AuthorizeResult{
		InvoiceID:  fmt.Sprintf("mock_inv_%d_%s", n, in.ReferenceID),
		AuthID:     fmt.Sprintf("mock_auth_%d_%s", n, in.ReferenceID),
		PaymentURL: mockPaymentURL(in.Method, in.ReferenceID),
	}, nil
}

func (m *MockGateway) Capture(ctx context.Context, in domain.CaptureRequest) error {
	_ = ctx
	if in.AuthID == "" {
		return domain.ErrPaymentFailed
	}
	return nil
}

func (m *MockGateway) Void(ctx context.Context, in domain.VoidRequest) error {
	_ = ctx
	if in.AuthID == "" {
		return domain.ErrPaymentFailed
	}
	return nil
}

var _ domain.PaymentGateway = (*MockGateway)(nil)
