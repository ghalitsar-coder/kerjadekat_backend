package redisx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"kerjadekat/backend/config"

	"github.com/redis/go-redis/v9"
)

// NewClient builds a standalone Redis client from application config.
func NewClient(cfg *config.Config) *redis.Client {
	addr := cfg.RedisHost + ":" + cfg.RedisPort
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.RedisPassword,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
}

// Ping verifies connectivity (optional startup check).
func Ping(ctx context.Context, c *redis.Client) error {
	if c == nil {
		return fmt.Errorf("nil redis client")
	}
	return c.Ping(ctx).Err()
}

// IntToString avoids strconv import at call sites for Redis set members.
func IntToString(v int) string {
	return strconv.Itoa(v)
}
