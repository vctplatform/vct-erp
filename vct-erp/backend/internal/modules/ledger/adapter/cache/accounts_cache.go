package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"vct-platform/backend/internal/modules/ledger/domain"
)

// ErrCacheMiss is returned when a key does not exist in Redis.
var ErrCacheMiss = errors.New("cache miss")

// RedisValueStore captures the Redis GET/SET subset needed by the chart-of-accounts cache.
type RedisValueStore interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// ChartOfAccountsRedisCache keeps account metadata in Redis to reduce database pressure.
type ChartOfAccountsRedisCache struct {
	redis     RedisValueStore
	keyPrefix string
}

// NewChartOfAccountsRedisCache constructs the Redis-backed chart cache adapter.
func NewChartOfAccountsRedisCache(redis RedisValueStore, keyPrefix string) *ChartOfAccountsRedisCache {
	return &ChartOfAccountsRedisCache{
		redis:     redis,
		keyPrefix: keyPrefix,
	}
}

// GetMany loads cached accounts by identifier and skips cache misses transparently.
func (c *ChartOfAccountsRedisCache) GetMany(ctx context.Context, ids []string) (map[string]domain.Account, error) {
	result := make(map[string]domain.Account, len(ids))
	if c == nil || c.redis == nil {
		return result, nil
	}

	for _, id := range ids {
		raw, err := c.redis.Get(ctx, c.keyFor(id))
		if err != nil {
			if errors.Is(err, ErrCacheMiss) {
				continue
			}
			return nil, fmt.Errorf("redis get account %s: %w", id, err)
		}

		var account domain.Account
		if err := json.Unmarshal(raw, &account); err != nil {
			return nil, fmt.Errorf("decode cached account %s: %w", id, err)
		}
		result[id] = account
	}

	return result, nil
}

// SetMany stores accounts in Redis as JSON documents.
func (c *ChartOfAccountsRedisCache) SetMany(ctx context.Context, accounts []domain.Account, ttl time.Duration) error {
	if c == nil || c.redis == nil {
		return nil
	}

	for _, account := range accounts {
		raw, err := json.Marshal(account)
		if err != nil {
			return fmt.Errorf("encode account %s: %w", account.ID, err)
		}
		if err := c.redis.Set(ctx, c.keyFor(account.ID), raw, ttl); err != nil {
			return fmt.Errorf("redis set account %s: %w", account.ID, err)
		}
	}

	return nil
}

func (c *ChartOfAccountsRedisCache) keyFor(accountID string) string {
	if c.keyPrefix == "" {
		return "coa:" + accountID
	}
	return c.keyPrefix + ":" + accountID
}
