package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgercache "vct-platform/backend/internal/modules/ledger/adapter/cache"
)

// RedisValueStore captures the Redis GET/SET/DEL subset required by dashboard caching.
type RedisValueStore interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// DashboardRedisCache stores dashboard JSON snapshots in Redis for 60-second executive reads.
type DashboardRedisCache struct {
	redis     RedisValueStore
	keyPrefix string
}

// NewDashboardRedisCache constructs a Redis-backed analytics dashboard cache.
func NewDashboardRedisCache(redis RedisValueStore, keyPrefix string) *DashboardRedisCache {
	return &DashboardRedisCache{
		redis:     redis,
		keyPrefix: keyPrefix,
	}
}

// Get returns the cached dashboard snapshot for the requested company.
func (c *DashboardRedisCache) Get(ctx context.Context, companyCode string) (analyticsdomain.CommandCenterDashboardData, error) {
	if c == nil || c.redis == nil {
		return analyticsdomain.CommandCenterDashboardData{}, analyticsdomain.ErrDashboardCacheMiss
	}

	raw, err := c.redis.Get(ctx, c.keyFor(companyCode))
	if err != nil {
		if errors.Is(err, ledgercache.ErrCacheMiss) {
			return analyticsdomain.CommandCenterDashboardData{}, analyticsdomain.ErrDashboardCacheMiss
		}
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("redis get dashboard snapshot: %w", err)
	}

	var snapshot analyticsdomain.CommandCenterDashboardData
	if err := json.Unmarshal(raw, &snapshot); err != nil {
		return analyticsdomain.CommandCenterDashboardData{}, fmt.Errorf("decode dashboard snapshot: %w", err)
	}

	return snapshot, nil
}

// Set stores a dashboard snapshot with a short TTL.
func (c *DashboardRedisCache) Set(ctx context.Context, snapshot analyticsdomain.CommandCenterDashboardData, ttl time.Duration) error {
	if c == nil || c.redis == nil {
		return nil
	}

	raw, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("encode dashboard snapshot: %w", err)
	}

	if err := c.redis.Set(ctx, c.keyFor(snapshot.CompanyCode), raw, ttl); err != nil {
		return fmt.Errorf("redis set dashboard snapshot: %w", err)
	}

	return nil
}

// Delete invalidates the cached dashboard snapshot for the given company.
func (c *DashboardRedisCache) Delete(ctx context.Context, companyCode string) error {
	if c == nil || c.redis == nil {
		return nil
	}
	return c.redis.Delete(ctx, c.keyFor(companyCode))
}

func (c *DashboardRedisCache) keyFor(companyCode string) string {
	if c.keyPrefix == "" {
		return "finance:dashboard:" + companyCode
	}
	return c.keyPrefix + ":" + companyCode
}
