package redisclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	ledgercache "vct-platform/backend/internal/modules/ledger/adapter/cache"
)

// Client adapts go-redis to the cache and stream publisher interfaces used by the ledger module.
type Client struct {
	inner *redis.Client
}

// New constructs a Redis client for both cache and stream publishing.
func New(addr string, username string, password string, db int) *Client {
	return &Client{
		inner: redis.NewClient(&redis.Options{
			Addr:     addr,
			Username: username,
			Password: password,
			DB:       db,
		}),
	}
}

// Ping verifies the Redis connection.
func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.inner == nil {
		return fmt.Errorf("redis client is not configured")
	}
	return c.inner.Ping(ctx).Err()
}

// Close releases the underlying Redis client.
func (c *Client) Close() error {
	if c == nil || c.inner == nil {
		return nil
	}
	return c.inner.Close()
}

// Get loads a raw value from Redis.
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	if c == nil || c.inner == nil {
		return nil, fmt.Errorf("redis client is not configured")
	}

	value, err := c.inner.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ledgercache.ErrCacheMiss
		}
		return nil, err
	}
	return value, nil
}

// Set stores a raw value in Redis with TTL.
func (c *Client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if c == nil || c.inner == nil {
		return fmt.Errorf("redis client is not configured")
	}
	return c.inner.Set(ctx, key, value, ttl).Err()
}

// Delete removes a value from Redis.
func (c *Client) Delete(ctx context.Context, key string) error {
	if c == nil || c.inner == nil {
		return fmt.Errorf("redis client is not configured")
	}
	return c.inner.Del(ctx, key).Err()
}

// XAdd appends an event into a Redis stream.
func (c *Client) XAdd(ctx context.Context, stream string, values map[string]string) (string, error) {
	if c == nil || c.inner == nil {
		return "", fmt.Errorf("redis client is not configured")
	}

	payload := make(map[string]any, len(values))
	for key, value := range values {
		payload[key] = value
	}

	return c.inner.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: payload,
	}).Result()
}

// Publish sends a payload to a Redis pub/sub channel.
func (c *Client) Publish(ctx context.Context, channel string, payload []byte) error {
	if c == nil || c.inner == nil {
		return fmt.Errorf("redis client is not configured")
	}
	return c.inner.Publish(ctx, channel, payload).Err()
}

// PubSubMessage is a transport-agnostic Redis pub/sub message.
type PubSubMessage struct {
	Channel string
	Payload string
}

// Subscription wraps a Redis pub/sub subscription.
type Subscription struct {
	inner *redis.PubSub
}

// Subscribe creates a pub/sub subscription on the given channel.
func (c *Client) Subscribe(ctx context.Context, channel string) (*Subscription, error) {
	if c == nil || c.inner == nil {
		return nil, fmt.Errorf("redis client is not configured")
	}

	pubsub := c.inner.Subscribe(ctx, channel)
	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()
		return nil, err
	}

	return &Subscription{inner: pubsub}, nil
}

// ReceiveMessage blocks until the next pub/sub message arrives.
func (s *Subscription) ReceiveMessage(ctx context.Context) (PubSubMessage, error) {
	if s == nil || s.inner == nil {
		return PubSubMessage{}, fmt.Errorf("redis subscription is not configured")
	}

	message, err := s.inner.ReceiveMessage(ctx)
	if err != nil {
		return PubSubMessage{}, err
	}

	return PubSubMessage{
		Channel: message.Channel,
		Payload: message.Payload,
	}, nil
}

// Close terminates the underlying pub/sub subscription.
func (s *Subscription) Close() error {
	if s == nil || s.inner == nil {
		return nil
	}
	return s.inner.Close()
}
