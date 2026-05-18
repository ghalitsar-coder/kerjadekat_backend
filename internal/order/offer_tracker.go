package order

import (
	"context"
	"fmt"
	"time"

	"kerjadekat/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const offeredKeyPrefix = "kerjadekat:order_offered:"

// RedisOfferTracker tracks workers already pinged for an order using a Redis SET.
type RedisOfferTracker struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisOfferTracker(rdb *redis.Client, ttl time.Duration) *RedisOfferTracker {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	return &RedisOfferTracker{rdb: rdb, ttl: ttl}
}

func offeredKey(orderID uuid.UUID) string {
	return offeredKeyPrefix + orderID.String()
}

func (t *RedisOfferTracker) FilterNotYetOffered(ctx context.Context, orderID uuid.UUID, candidates []uuid.UUID) ([]uuid.UUID, error) {
	if len(candidates) == 0 {
		return nil, nil
	}
	key := offeredKey(orderID)
	members, err := t.rdb.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("smembers offered: %w", err)
	}
	seen := make(map[string]struct{}, len(members))
	for _, m := range members {
		seen[m] = struct{}{}
	}
	out := make([]uuid.UUID, 0, len(candidates))
	for _, id := range candidates {
		if _, ok := seen[id.String()]; ok {
			continue
		}
		out = append(out, id)
	}
	return out, nil
}

func (t *RedisOfferTracker) RecordOffered(ctx context.Context, orderID uuid.UUID, workerIDs []uuid.UUID) error {
	if len(workerIDs) == 0 {
		return nil
	}
	key := offeredKey(orderID)
	pipe := t.rdb.TxPipeline()
	members := make([]interface{}, len(workerIDs))
	for i, id := range workerIDs {
		members[i] = id.String()
	}
	pipe.SAdd(ctx, key, members...)
	pipe.Expire(ctx, key, t.ttl)
	_, err := pipe.Exec(ctx)
	return err
}

var _ domain.OrderOfferTracker = (*RedisOfferTracker)(nil)
