package realtime

import (
	"context"
	"fmt"

	infraRedis "vct-platform/backend/internal/infra/redisclient"
)

// RedisClientSubscriber adapts the shared Redis client to the bridge subscription interface.
type RedisClientSubscriber struct {
	client interface {
		Subscribe(ctx context.Context, channel string) (*infraRedis.Subscription, error)
	}
}

// NewRedisClientSubscriber constructs a realtime bridge subscriber backed by the shared Redis client.
func NewRedisClientSubscriber(client interface {
	Subscribe(ctx context.Context, channel string) (*infraRedis.Subscription, error)
}) *RedisClientSubscriber {
	return &RedisClientSubscriber{client: client}
}

// Subscribe opens the Redis subscription and wraps its messages for the bridge.
func (s *RedisClientSubscriber) Subscribe(ctx context.Context, channel string) (Subscription, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("redis realtime subscriber is not configured")
	}

	subscription, err := s.client.Subscribe(ctx, channel)
	if err != nil {
		return nil, err
	}

	return &redisSubscription{inner: subscription}, nil
}

type redisSubscription struct {
	inner interface {
		ReceiveMessage(ctx context.Context) (infraRedis.PubSubMessage, error)
		Close() error
	}
}

func (s *redisSubscription) ReceiveMessage(ctx context.Context) (Message, error) {
	message, err := s.inner.ReceiveMessage(ctx)
	if err != nil {
		return Message{}, err
	}
	return Message{
		Channel: message.Channel,
		Payload: message.Payload,
	}, nil
}

func (s *redisSubscription) Close() error {
	if s == nil || s.inner == nil {
		return nil
	}
	return s.inner.Close()
}
