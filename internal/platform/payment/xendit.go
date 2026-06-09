package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"kerjadekat/backend/internal/domain"
)

const xenditAPIBase = "https://api.xendit.co"

// XenditGateway uses Xendit Payment Requests API (v3).
type XenditGateway struct {
	apiKey        string
	callbackToken string
	http          *http.Client
}

func NewXenditGateway(apiKey, callbackToken string) *XenditGateway {
	return &XenditGateway{
		apiKey:        apiKey,
		callbackToken: callbackToken,
		http:          &http.Client{Timeout: 30 * time.Second},
	}
}

type paymentRequestReq struct {
	ReferenceID   string `json:"reference_id"`
	RequestAmount float64 `json:"request_amount"`
	Currency      string `json:"currency"`
	Country       string `json:"country"`
	ChannelCode   string `json:"channel_code"`
	Type          string `json:"type"`
}

type paymentRequestRes struct {
	PaymentRequestID string `json:"payment_request_id"`
	Status           string `json:"status"`
}

func (x *XenditGateway) Authorize(ctx context.Context, in domain.AuthorizeRequest) (domain.AuthorizeResult, error) {
	channel := in.Method
	if channel == "" || channel == "qris" {
		channel = "QRIS"
	}

	body, _ := json.Marshal(paymentRequestReq{
		ReferenceID:   in.ReferenceID,
		RequestAmount: in.AmountIDR,
		Currency:      "IDR",
		Country:       "ID",
		ChannelCode:   channel,
		Type:          "PAY",
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, xenditAPIBase+"/v3/payment_requests", bytes.NewReader(body))
	if err != nil {
		return domain.AuthorizeResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-version", "2024-11-11")
	req.SetBasicAuth(x.apiKey, "")

	res, err := x.http.Do(req)
	if err != nil {
		return domain.AuthorizeResult{}, fmt.Errorf("xendit authorize: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return domain.AuthorizeResult{}, fmt.Errorf("xendit authorize http %d: %s", res.StatusCode, string(raw))
	}
	var pr paymentRequestRes
	if err := json.Unmarshal(raw, &pr); err != nil {
		return domain.AuthorizeResult{}, err
	}
	return domain.AuthorizeResult{
		InvoiceID: pr.PaymentRequestID,
		AuthID:    pr.PaymentRequestID,
	}, nil
}

func (x *XenditGateway) Capture(ctx context.Context, in domain.CaptureRequest) error {
	_ = ctx
	_ = in
	// Payment Request flow: capture is implicit when paid (webhook).
	return nil
}

func (x *XenditGateway) Void(ctx context.Context, in domain.VoidRequest) error {
	if in.AuthID == "" {
		return domain.ErrPaymentFailed
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/v3/payment_requests/%s/cancel", xenditAPIBase, in.AuthID), bytes.NewReader([]byte("{}")))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-version", "2024-11-11")
	req.SetBasicAuth(x.apiKey, "")
	res, err := x.http.Do(req)
	if err != nil {
		return fmt.Errorf("xendit void: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		raw, _ := io.ReadAll(res.Body)
		return fmt.Errorf("xendit void http %d: %s", res.StatusCode, string(raw))
	}
	return nil
}

var _ domain.PaymentGateway = (*XenditGateway)(nil)
