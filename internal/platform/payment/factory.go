package payment

import (
	"kerjadekat/backend/config"
	"kerjadekat/backend/internal/domain"
)

// NewFromConfig returns MockGateway in non-production; XenditGateway when ENV=production.
func NewFromConfig(cfg *config.Config) domain.PaymentGateway {
	if cfg.Env == "production" && cfg.XenditAPIKey != "" {
		return NewXenditGateway(cfg.XenditAPIKey, cfg.XenditCallbackToken)
	}
	return NewMockGateway()
}
