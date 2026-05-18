package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"kerjadekat/backend/internal/domain"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "kerjadekat:otp:"

// Redis stores OTP codes in Redis with TTL (multi-pod safe).
type Redis struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedis(rdb *redis.Client, ttl time.Duration) *Redis {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &Redis{rdb: rdb, ttl: ttl}
}

func (s *Redis) key(phone string) string {
	return keyPrefix + phone
}

func (s *Redis) Store(ctx context.Context, phone, code string) error {
	return s.rdb.Set(ctx, s.key(phone), code, s.ttl).Err()
}

func (s *Redis) VerifyAndConsume(ctx context.Context, phone, code string) (bool, error) {
	key := s.key(phone)
	stored, err := s.rdb.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == code, nil
}

var _ domain.OTPStore = (*Redis)(nil)

// Random6 returns a 6-digit numeric OTP.
func Random6() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
