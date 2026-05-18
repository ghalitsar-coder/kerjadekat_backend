package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	claimTypeAccess  = "access"
	claimTypeRefresh = "refresh"
)

// Claims carried in access tokens.
type Claims struct {
	UserID uuid.UUID
	Role   string
}

type Issuer struct {
	secret       []byte
	accessTTL    time.Duration
	refreshTTL   time.Duration
	issuer       string
	audience     string
}

func NewIssuer(secret string, accessTTL, refreshTTL time.Duration) *Issuer {
	return &Issuer{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     "kerjadekat",
		audience:   "kerjadekat-api",
	}
}

func (i *Issuer) AccessExpiresInSeconds() int64 {
	if i.accessTTL <= 0 {
		return 0
	}
	return int64(i.accessTTL / time.Second)
}

func (i *Issuer) IssueAccess(c Claims) (string, error) {
	return i.sign(c, claimTypeAccess, i.accessTTL)
}

func (i *Issuer) IssueRefresh(c Claims) (string, error) {
	return i.sign(c, claimTypeRefresh, i.refreshTTL)
}

func (i *Issuer) sign(c Claims, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  c.UserID.String(),
		"role": c.Role,
		"typ":  typ,
		"iss":  i.issuer,
		"aud":  i.audience,
		"iat":  now.Unix(),
		"exp":  now.Add(ttl).Unix(),
		"nbf":  now.Unix(),
	})
	return t.SignedString(i.secret)
}

func (i *Issuer) ParseAccess(token string) (Claims, error) {
	return i.parse(token, claimTypeAccess)
}

func (i *Issuer) ParseRefresh(token string) (Claims, error) {
	return i.parse(token, claimTypeRefresh)
}

func (i *Issuer) parse(tokenString, wantTyp string) (Claims, error) {
	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return i.secret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(i.issuer),
		jwt.WithAudience(i.audience),
		jwt.WithLeeway(5*time.Second),
	)
	if err != nil || !tok.Valid {
		return Claims{}, ErrInvalidToken
	}

	mc, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, ErrInvalidToken
	}
	if typ, _ := mc["typ"].(string); typ != wantTyp {
		return Claims{}, ErrInvalidToken
	}
	sub, _ := mc["sub"].(string)
	uid, err := uuid.Parse(sub)
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	role, _ := mc["role"].(string)
	if role == "" {
		return Claims{}, ErrInvalidToken
	}
	return Claims{UserID: uid, Role: role}, nil
}

var ErrInvalidToken = errors.New("invalid token")
