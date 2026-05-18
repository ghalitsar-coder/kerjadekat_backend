package sms

import (
	"context"
	"log"

	"kerjadekat/backend/internal/domain"
)

// MockNotifier logs OTP delivery; swap for a real SMS/WhatsApp provider later.
type MockNotifier struct{}

func NewMockNotifier() *MockNotifier {
	return &MockNotifier{}
}

func (m *MockNotifier) SendOTP(ctx context.Context, phoneNumber, code string) error {
	_ = ctx
	log.Printf("sms_mock: OTP for %s = %s", phoneNumber, code)
	return nil
}

var _ domain.SMSNotifier = (*MockNotifier)(nil)
