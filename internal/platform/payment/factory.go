package payment

import (
	"kerjadekat/backend/config"
	"kerjadekat/backend/internal/domain"
)

// NewFromConfig returns XenditGateway if XENDIT_API_KEY is set; MockGateway otherwise.
func NewFromConfig(cfg *config.Config) domain.PaymentGateway {
	if cfg.XenditAPIKey != "" {
		return NewXenditGateway(cfg.XenditAPIKey, cfg.XenditCallbackToken)
	}
	return NewMockGateway()
}
