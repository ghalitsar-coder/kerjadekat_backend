package sms

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"

	"kerjadekat/backend/internal/domain"
)

// MockNotifier logs OTP delivery; swap for a real SMS/WhatsApp provider later.
type MockNotifier struct {
	logFile string
	mu      sync.Mutex
}

func NewMockNotifier(logFile string) *MockNotifier {
	return &MockNotifier{logFile: logFile}
}

func (m *MockNotifier) SendOTP(ctx context.Context, phoneNumber, code string) error {
	_ = ctx
	line := "sms_mock: OTP for " + phoneNumber + " = " + code
	log.Print(line)
	if m.logFile != "" {
		m.mu.Lock()
		defer m.mu.Unlock()
		_ = os.MkdirAll(filepath.Dir(m.logFile), 0o755)
		f, err := os.OpenFile(m.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err == nil {
			_, _ = f.WriteString(line + "\n")
			_ = f.Close()
		}
	}
	return nil
}

var _ domain.SMSNotifier = (*MockNotifier)(nil)
