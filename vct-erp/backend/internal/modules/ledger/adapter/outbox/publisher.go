package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vct-platform/backend/internal/modules/ledger/domain"
)

// RedisStreamAppender captures the XADD capability needed for the outbox relay.
type RedisStreamAppender interface {
	XAdd(ctx context.Context, stream string, values map[string]string) (string, error)
}

// RedisStreamPublisher publishes outbox events to a Redis Stream after the DB commit succeeds.
type RedisStreamPublisher struct {
	client RedisStreamAppender
}

// NewRedisStreamPublisher constructs the Redis Stream publisher.
func NewRedisStreamPublisher(client RedisStreamAppender) *RedisStreamPublisher {
	return &RedisStreamPublisher{client: client}
}

// Publish appends a posted ledger event to Redis Stream for downstream modules.
func (p *RedisStreamPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("redis stream client is not configured")
	}
	if !json.Valid(event.Payload) {
		return fmt.Errorf("outbox payload for %s is not valid json", event.ID)
	}

	_, err := p.client.XAdd(ctx, event.StreamKey, map[string]string{
		"event_id":       event.ID,
		"aggregate_type": event.AggregateType,
		"aggregate_id":   event.AggregateID,
		"event_type":     event.EventType,
		"status":         string(event.Status),
		"payload":        string(event.Payload),
		"created_at":     event.CreatedAt.UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return fmt.Errorf("xadd redis stream %s: %w", event.StreamKey, err)
	}

	return nil
}
