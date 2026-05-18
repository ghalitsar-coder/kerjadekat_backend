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

// XenditGateway uses Xendit Invoice API for platform-fee holds (production only).
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

type invoiceCreateReq struct {
	ExternalID  string  `json:"external_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Currency    string  `json:"currency"`
}

type invoiceCreateRes struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (x *XenditGateway) Authorize(ctx context.Context, in domain.AuthorizeRequest) (domain.AuthorizeResult, error) {
	body, _ := json.Marshal(invoiceCreateReq{
		ExternalID:  in.ReferenceID,
		Amount:      in.AmountIDR,
		Description: fmt.Sprintf("KerjaDekat admin fee (%s)", in.Method),
		Currency:    "IDR",
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, xenditAPIBase+"/v2/invoices", bytes.NewReader(body))
	if err != nil {
		return domain.AuthorizeResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
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
	var inv invoiceCreateRes
	if err := json.Unmarshal(raw, &inv); err != nil {
		return domain.AuthorizeResult{}, err
	}
	return domain.AuthorizeResult{
		InvoiceID: inv.ID,
		AuthID:    inv.ID,
	}, nil
}

func (x *XenditGateway) Capture(ctx context.Context, in domain.CaptureRequest) error {
	_ = ctx
	_ = in
	// Invoice flow: capture is implicit when invoice is paid (webhook). No-op for MVP.
	return nil
}

func (x *XenditGateway) Void(ctx context.Context, in domain.VoidRequest) error {
	if in.AuthID == "" {
		return domain.ErrPaymentFailed
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/v2/invoices/%s/expire!", xenditAPIBase, in.AuthID), nil)
	if err != nil {
		return err
	}
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
