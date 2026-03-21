package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
)

// Subscription represents the blocking receive loop needed for dashboard event fan-out.
type Subscription interface {
	ReceiveMessage(ctx context.Context) (Message, error)
	Close() error
}

// Message is the transport-agnostic Redis Pub/Sub message used by the bridge.
type Message struct {
	Channel string
	Payload string
}

// Subscriber creates realtime subscriptions for finance dashboard events.
type Subscriber interface {
	Subscribe(ctx context.Context, channel string) (Subscription, error)
}

// CacheInvalidator removes stale dashboard snapshots after a new transaction arrives.
type CacheInvalidator interface {
	Delete(ctx context.Context, companyCode string) error
}

// RedisBridge fans dashboard events from Redis Pub/Sub to the in-process websocket hub.
type RedisBridge struct {
	subscriber Subscriber
	cache      CacheInvalidator
	hub        *Hub
	channel    string
}

// NewRedisBridge constructs the realtime dashboard bridge.
func NewRedisBridge(subscriber Subscriber, cache CacheInvalidator, hub *Hub, channel string) *RedisBridge {
	return &RedisBridge{
		subscriber: subscriber,
		cache:      cache,
		hub:        hub,
		channel:    channel,
	}
}

// Run listens until the context is cancelled.
func (b *RedisBridge) Run(ctx context.Context) error {
	if b == nil || b.subscriber == nil {
		return nil
	}

	subscription, err := b.subscriber.Subscribe(ctx, b.channelOrDefault())
	if err != nil {
		return fmt.Errorf("subscribe dashboard redis channel: %w", err)
	}
	defer subscription.Close()

	for {
		message, err := subscription.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, context.Canceled) {
				return nil
			}
			return fmt.Errorf("receive dashboard event: %w", err)
		}

		var event analyticsdomain.DashboardEvent
		if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
			continue
		}
		if b.cache != nil && event.CompanyCode != "" {
			_ = b.cache.Delete(ctx, event.CompanyCode)
		}
		if b.hub != nil {
			b.hub.Broadcast(event)
		}
	}
}

func (b *RedisBridge) channelOrDefault() string {
	if b.channel == "" {
		return "finance.dashboard.events"
	}
	return b.channel
}
