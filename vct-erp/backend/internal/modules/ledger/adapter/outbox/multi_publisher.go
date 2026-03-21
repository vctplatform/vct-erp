package outbox

import (
	"context"

	"vct-platform/backend/internal/modules/ledger/domain"
)

// MultiPublisher fans a committed outbox event to one required publisher and optional side publishers.
type MultiPublisher struct {
	primary   domain.EventPublisher
	secondary []domain.EventPublisher
}

// NewMultiPublisher constructs a fan-out publisher where only the primary publisher is critical.
func NewMultiPublisher(primary domain.EventPublisher, secondary ...domain.EventPublisher) *MultiPublisher {
	return &MultiPublisher{
		primary:   primary,
		secondary: secondary,
	}
}

// Publish sends the event to the primary publisher, then best-effort to secondary publishers.
func (p *MultiPublisher) Publish(ctx context.Context, event domain.OutboxEvent) error {
	if p == nil || p.primary == nil {
		return nil
	}
	if err := p.primary.Publish(ctx, event); err != nil {
		return err
	}

	for _, publisher := range p.secondary {
		if publisher == nil {
			continue
		}
		_ = publisher.Publish(ctx, event)
	}

	return nil
}
