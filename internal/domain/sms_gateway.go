package domain

import "context"

// SMSNotifier abstracts transactional SMS/WhatsApp OTP delivery.
type SMSNotifier interface {
	SendOTP(ctx context.Context, phoneNumber, code string) error
}
