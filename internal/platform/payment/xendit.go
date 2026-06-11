package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"kerjadekat/backend/internal/domain"
)

const xenditAPIBase = "https://api.xendit.co"

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
	ReferenceID   string  `json:"reference_id"`
	RequestAmount float64 `json:"request_amount"`
	Currency      string  `json:"currency"`
	Country       string  `json:"country"`
	ChannelCode   string  `json:"channel_code"`
	Type          string  `json:"type"`
}

type paymentRequestRes struct {
	PaymentRequestID string `json:"payment_request_id"`
	Status           string `json:"status"`
	Actions          []struct {
		Type       string `json:"type"`
		Descriptor string `json:"descriptor"`
		Value      string `json:"value"`
	} `json:"actions"`
}

var channelCodeMap = map[string]string{
	"qris":  "QRIS",
	"gopay": "ID_GOPAY",
	"ovo":   "ID_OVO",
	"dana":  "ID_DANA",
}

func xenditPaymentURL(actions []struct {
	Type       string `json:"type"`
	Descriptor string `json:"descriptor"`
	Value      string `json:"value"`
}) string {
	for _, a := range actions {
		switch a.Descriptor {
		case "QR_STRING":
			return "https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=" + url.QueryEscape(a.Value)
		case "WEB_URL":
			return a.Value
		default:
			if isURL(a.Value) {
				return a.Value
			}
		}
	}
	return ""
}

func isURL(s string) bool {
	return len(s) > 4 && (s[:4] == "http" || s[:3] == "goj" || s[:3] == "ovo")
}

func (x *XenditGateway) Authorize(ctx context.Context, in domain.AuthorizeRequest) (domain.AuthorizeResult, error) {
	channel, ok := channelCodeMap[in.Method]
	if !ok {
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
		log.Printf("[xendit] HTTP error: %v", err)
		return domain.AuthorizeResult{}, fmt.Errorf("xendit authorize: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Printf("[xendit] HTTP %d: %s", res.StatusCode, string(raw))
		return domain.AuthorizeResult{}, fmt.Errorf("xendit authorize http %d: %s", res.StatusCode, string(raw))
	}
	var pr paymentRequestRes
	if err := json.Unmarshal(raw, &pr); err != nil {
		log.Printf("[xendit] unmarshal error: %v body=%s", err, string(raw))
		return domain.AuthorizeResult{}, err
	}
	payURL := xenditPaymentURL(pr.Actions)
	log.Printf("[xendit] payment_request_id=%s status=%s actions=%d url=%q",
		pr.PaymentRequestID, pr.Status, len(pr.Actions), payURL)
	return domain.AuthorizeResult{
		InvoiceID:  pr.PaymentRequestID,
		AuthID:     pr.PaymentRequestID,
		PaymentURL: payURL,
	}, nil
}

func (x *XenditGateway) Capture(ctx context.Context, in domain.CaptureRequest) error {
	_ = ctx
	_ = in
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
