package usecase

import (
	"context"
	"errors"
	"strings"

	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/otp"
	"kerjadekat/backend/pkg/token"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
		if err == domain.ErrNotFound {
			return nil, errors.New("user not found. please register first")
		}
		return nil, err
	}

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
	email := strings.TrimSpace(strings.ToLower(in.Email))
	name := strings.TrimSpace(in.Name)
	if name == "" {
		name = "Pengguna"
	}

	u, err := a.users.FindByProviderID(ctx, provider, subject)
	if err != nil {
		if err != domain.ErrNotFound {
			return nil, err
		}
		if email != "" {
			u, err = a.users.FindByEmail(ctx, email)
			if err != nil && err != domain.ErrNotFound {
				return nil, err
			}
		}
		if u == nil {
			u = &domain.User{
				ID:         uuid.New(),
				FullName:   name,
				Role:       role,
				Provider:   &provider,
				ProviderID: &subject,
				Status:     domain.UserStatusActive,
			}
			if email != "" {
				u.Email = &email
			}
			if err := a.users.Create(ctx, u); err != nil {
				return nil, err
			}
		} else {
			u.Provider = &provider
			u.ProviderID = &subject
			if err := a.users.Update(ctx, u); err != nil {
				return nil, err
			}
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

func (a *Auth) RegisterEmail(ctx context.Context, email, password, name, phone, role string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return domain.ErrInvalidInput
	}
	if role == "" {
		role = domain.RoleConsumer
	}
	if err := validateRole(role); err != nil {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashStr := string(hashed)

	phone = normalizePhone(phone)
	if phone == "" {
		return errors.New("phone number is required for 2FA")
	}

	_, err = a.users.FindByEmail(ctx, email)
	if err == nil {
		return domain.ErrConflict
	}
	_, err = a.users.FindByPhone(ctx, phone)
	if err == nil {
		return domain.ErrConflict
	}

	u := &domain.User{
		ID:          uuid.New(),
		Email:       &email,
		Password:    &hashStr,
		PhoneNumber: &phone,
		FullName:    name,
		Role:        role,
		Status:      domain.UserStatusActive,
	}

	if err := a.users.Create(ctx, u); err != nil {
		return err
	}

	if role == domain.RoleWorker {
		return a.ensureWorkerProfile(ctx, u.ID)
	}
	return nil
}

func (a *Auth) LoginEmail(ctx context.Context, email, password string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, err := a.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if u.Password == nil {
		return nil, domain.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password)); err != nil {
		return nil, domain.ErrUnauthorized
	}
	return u, nil
}
