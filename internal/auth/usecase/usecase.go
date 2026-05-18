package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/otp"
	"kerjadekat/backend/pkg/token"

	"github.com/google/uuid"
)

type Auth struct {
	users   domain.UserRepository
	workers domain.WorkerRepository
	otp     domain.OTPStore
	sms     domain.SMSNotifier
	tokens  *token.Issuer
}

func NewAuth(users domain.UserRepository, workers domain.WorkerRepository, otpStore domain.OTPStore, sms domain.SMSNotifier, tokens *token.Issuer) *Auth {
	return &Auth{users: users, workers: workers, otp: otpStore, sms: sms, tokens: tokens}
}

func normalizePhone(p string) string {
	return strings.TrimSpace(p)
}

func validateRole(role string) error {
	switch role {
	case domain.RoleWorker, domain.RoleConsumer, domain.RoleAgent, domain.RoleAdmin:
		return nil
	default:
		return domain.ErrInvalidInput
	}
}

func (a *Auth) RequestOTP(ctx context.Context, phone string) error {
	phone = normalizePhone(phone)
	if phone == "" || len(phone) < 8 {
		return domain.ErrInvalidInput
	}
	code, err := otp.Random6()
	if err != nil {
		return err
	}
	if err := a.otp.Store(ctx, phone, code); err != nil {
		return err
	}
	return a.sms.SendOTP(ctx, phone, code)
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (a *Auth) VerifyOTP(ctx context.Context, phone, code, role string) (*TokenPair, error) {
	phone = normalizePhone(phone)
	if phone == "" || code == "" {
		return nil, domain.ErrInvalidInput
	}
	if role == "" {
		role = domain.RoleConsumer
	}
	if err := validateRole(role); err != nil {
		return nil, err
	}
	ok, err := a.otp.VerifyAndConsume(ctx, phone, code)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrInvalidOTP
	}

	u, err := a.users.FindByPhone(ctx, phone)
	if err != nil {
		if err != domain.ErrNotFound {
			return nil, err
		}
		name := "Pengguna"
		u = &domain.User{
			ID:          uuid.New(),
			PhoneNumber: phone,
			FullName:    name,
			Role:        role,
			Status:      domain.UserStatusActive,
		}
		if err := a.users.Create(ctx, u); err != nil {
			return nil, err
		}
	} else {
		if u.Role != role {
			return nil, domain.ErrForbidden
		}
		if u.Status == domain.UserStatusSuspended {
			return nil, domain.ErrForbidden
		}
		if u.Status == domain.UserStatusPending {
			u.Status = domain.UserStatusActive
			if err := a.users.Update(ctx, u); err != nil {
				return nil, err
			}
		}
	}

	if role == domain.RoleWorker {
		if err := a.ensureWorkerProfile(ctx, u.ID); err != nil {
			return nil, err
		}
	}

	access, err := a.tokens.IssueAccess(token.Claims{UserID: u.ID, Role: u.Role})
	if err != nil {
		return nil, err
	}
	refresh, err := a.tokens.IssueRefresh(token.Claims{UserID: u.ID, Role: u.Role})
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresIn: a.tokens.AccessExpiresInSeconds()}, nil
}

func (a *Auth) ensureWorkerProfile(ctx context.Context, userID uuid.UUID) error {
	if _, err := a.workers.FindProfileByUserID(ctx, userID); err == nil {
		return nil
	} else if err != domain.ErrNotFound {
		return err
	}
	p := &domain.WorkerProfile{
		ID:           uuid.New(),
		UserID:       userID,
		Availability: domain.WorkerAvailabilityOffline,
	}
	return a.workers.CreateProfile(ctx, p)
}

type SocialLoginInput struct {
	Provider  string
	Subject   string
	Email     string
	Name      string
	Role      string
}

// SocialLogin issues JWT for OAuth users (Google/GitHub via Auth.js).
func (a *Auth) SocialLogin(ctx context.Context, in SocialLoginInput) (*TokenPair, error) {
	provider := strings.TrimSpace(strings.ToLower(in.Provider))
	subject := strings.TrimSpace(in.Subject)
	if provider == "" || subject == "" {
		return nil, domain.ErrInvalidInput
	}
	role := in.Role
	if role == "" {
		role = domain.RoleConsumer
	}
	if err := validateRole(role); err != nil {
		return nil, err
	}
	phone := oauthSyntheticPhone(provider, subject)
	name := strings.TrimSpace(in.Name)
	if name == "" {
		name = "Pengguna"
	}

	u, err := a.users.FindByPhone(ctx, phone)
	if err != nil {
		if err != domain.ErrNotFound {
			return nil, err
		}
		u = &domain.User{
			ID:          uuid.New(),
			PhoneNumber: phone,
			FullName:    name,
			Role:        role,
			Status:      domain.UserStatusActive,
		}
		if err := a.users.Create(ctx, u); err != nil {
			return nil, err
		}
	} else {
		if u.Role != role {
			return nil, domain.ErrForbidden
		}
		if u.Status == domain.UserStatusSuspended {
			return nil, domain.ErrForbidden
		}
		if name != "" && u.FullName != name {
			u.FullName = name
			_ = a.users.Update(ctx, u)
		}
	}

	if role == domain.RoleWorker {
		if err := a.ensureWorkerProfile(ctx, u.ID); err != nil {
			return nil, err
		}
	}

	access, err := a.tokens.IssueAccess(token.Claims{UserID: u.ID, Role: u.Role})
	if err != nil {
		return nil, err
	}
	refresh, err := a.tokens.IssueRefresh(token.Claims{UserID: u.ID, Role: u.Role})
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresIn: a.tokens.AccessExpiresInSeconds()}, nil
}

func oauthSyntheticPhone(provider, subject string) string {
	sum := sha256.Sum256([]byte(provider + ":" + subject))
	// Unique pseudo-phone (max 15 chars) for OAuth-only accounts.
	return "o" + hex.EncodeToString(sum[:6])
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	_ = ctx
	c, err := a.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	access, err := a.tokens.IssueAccess(token.Claims{UserID: c.UserID, Role: c.Role})
	if err != nil {
		return nil, err
	}
	newRefresh, err := a.tokens.IssueRefresh(token.Claims{UserID: c.UserID, Role: c.Role})
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: newRefresh, ExpiresIn: a.tokens.AccessExpiresInSeconds()}, nil
}
