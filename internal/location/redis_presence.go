package location

import (
	"context"
	"fmt"
	"time"

	"kerjadekat/backend/internal/domain"
	"kerjadekat/backend/pkg/redisx"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	geoKeyPrefix         = "kerjadekat:geo:workers"
	presenceKeyPrefix    = "kerjadekat:presence:"
	skillsKeyPrefix      = "kerjadekat:skills:"
	availabilityKeyPrefix = "kerjadekat:availability:"
)

// RedisPresence tracks worker GPS in Redis (GEO + TTL), not PostgreSQL.
type RedisPresence struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewRedisPresence(rdb *redis.Client, ttl time.Duration) *RedisPresence {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &RedisPresence{rdb: rdb, ttl: ttl}
}

func presenceKey(userID uuid.UUID) string {
	return presenceKeyPrefix + userID.String()
}

func skillsKey(userID uuid.UUID) string {
	return skillsKeyPrefix + userID.String()
}

func availabilityKey(userID uuid.UUID) string {
	return availabilityKeyPrefix + userID.String()
}

// SetAvailability mirrors worker availability into Redis for GEO matching filters.
func (p *RedisPresence) SetAvailability(ctx context.Context, userID uuid.UUID, availability string) error {
	switch availability {
	case domain.WorkerAvailabilityOnline:
		return p.rdb.Set(ctx, availabilityKey(userID), availability, p.ttl).Err()
	case domain.WorkerAvailabilityOffline, domain.WorkerAvailabilityBusy:
		pipe := p.rdb.TxPipeline()
		pipe.Del(ctx, availabilityKey(userID))
		pipe.Del(ctx, presenceKey(userID))
		pipe.ZRem(ctx, geoKeyPrefix, userID.String())
		_, err := pipe.Exec(ctx)
		return err
	default:
		return fmt.Errorf("unknown availability: %s", availability)
	}
}

// UpdateWorkerLocation writes GEO position and refreshes online + skill set TTLs.
func (p *RedisPresence) UpdateWorkerLocation(ctx context.Context, userID uuid.UUID, lat, lng float64, skillIDs []int) error {
	member := userID.String()
	loc := &redis.GeoLocation{
		Longitude: lng,
		Latitude:  lat,
		Name:      member,
	}
	if err := p.rdb.GeoAdd(ctx, geoKeyPrefix, loc).Err(); err != nil {
		return fmt.Errorf("geoadd: %w", err)
	}
	if err := p.rdb.Set(ctx, presenceKey(userID), "1", p.ttl).Err(); err != nil {
		return fmt.Errorf("presence: %w", err)
	}
	skKey := skillsKey(userID)
	pipe := p.rdb.TxPipeline()
	pipe.Del(ctx, skKey)
	if len(skillIDs) > 0 {
		args := make([]interface{}, 0, len(skillIDs))
		for _, id := range skillIDs {
			args = append(args, redisx.IntToString(id))
		}
		pipe.SAdd(ctx, skKey, args...)
	}
	pipe.Expire(ctx, skKey, p.ttl)
	if n, _ := p.rdb.Exists(ctx, availabilityKey(userID)).Result(); n > 0 {
		pipe.Expire(ctx, availabilityKey(userID), p.ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("skills set: %w", err)
	}
	return nil
}

// NearbyWorkerUserIDs returns workers within radiusMeters of lat/lng who have a fresh presence TTL and the skill.
func (p *RedisPresence) NearbyWorkerUserIDs(ctx context.Context, skillID int, lat, lng float64, radiusMeters float64) ([]uuid.UUID, error) {
	res, err := p.rdb.GeoRadius(ctx, geoKeyPrefix, lng, lat, &redis.GeoRadiusQuery{
		Radius:    radiusMeters,
		Unit:      "m",
		WithCoord: false,
		WithDist:  false,
		Count:     200,
		Sort:      "ASC",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("georadius: %w", err)
	}
	skillMember := redisx.IntToString(skillID)
	out := make([]uuid.UUID, 0, len(res))
	for _, loc := range res {
		uid, err := uuid.Parse(loc.Name)
		if err != nil {
			continue
		}
		ok, err := p.rdb.Exists(ctx, presenceKey(uid)).Result()
		if err != nil || ok == 0 {
			continue
		}
		avail, err := p.rdb.Get(ctx, availabilityKey(uid)).Result()
		if err != nil || avail != domain.WorkerAvailabilityOnline {
			continue
		}
		if skillID > 0 {
			isMember, err := p.rdb.SIsMember(ctx, skillsKey(uid), skillMember).Result()
			if err != nil || !isMember {
				continue
			}
		}
		out = append(out, uid)
	}
	return out, nil
}

var _ domain.WorkerPresence = (*RedisPresence)(nil)
