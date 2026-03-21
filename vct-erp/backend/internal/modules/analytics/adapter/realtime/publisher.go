package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

// PubSubPublisher captures the Redis publish capability required by the dashboard event publisher.
type PubSubPublisher interface {
	Publish(ctx context.Context, channel string, payload []byte) error
}

// RedisDashboardEventPublisher pushes transaction notifications into Redis Pub/Sub.
type RedisDashboardEventPublisher struct {
	client  PubSubPublisher
	channel string
	now     func() time.Time
}

// NewRedisDashboardEventPublisher constructs a Pub/Sub-backed finance dashboard publisher.
func NewRedisDashboardEventPublisher(client PubSubPublisher, channel string) *RedisDashboardEventPublisher {
	return &RedisDashboardEventPublisher{
		client:  client,
		channel: channel,
		now:     time.Now,
	}
}

// Publish transforms a committed ledger outbox event into a finance websocket payload.
func (p *RedisDashboardEventPublisher) Publish(ctx context.Context, event ledgerdomain.OutboxEvent) error {
	if p == nil || p.client == nil {
		return nil
	}
	if event.EventType != "ledger.entry.posted" {
		return nil
	}

	dashboardEvent, err := p.dashboardEvent(event)
	if err != nil {
		return fmt.Errorf("build dashboard event: %w", err)
	}

	raw, err := json.Marshal(dashboardEvent)
	if err != nil {
		return fmt.Errorf("marshal dashboard event: %w", err)
	}

	if err := p.client.Publish(ctx, p.channelOrDefault(), raw); err != nil {
		return fmt.Errorf("publish dashboard event: %w", err)
	}

	return nil
}

func (p *RedisDashboardEventPublisher) channelOrDefault() string {
	if strings.TrimSpace(p.channel) == "" {
		return "finance.dashboard.events"
	}
	return p.channel
}

func (p *RedisDashboardEventPublisher) dashboardEvent(event ledgerdomain.OutboxEvent) (analyticsdomain.DashboardEvent, error) {
	type payload struct {
		EntryID      string         `json:"entry_id"`
		ReferenceNo  string         `json:"reference_no"`
		CompanyCode  string         `json:"company_code"`
		SourceModule string         `json:"source_module"`
		PostedAt     string         `json:"posted_at"`
		TotalDebit   string         `json:"total_debit"`
		Metadata     map[string]any `json:"metadata"`
	}

	var source payload
	if err := json.Unmarshal(event.Payload, &source); err != nil {
		return analyticsdomain.DashboardEvent{}, err
	}

	amount, err := ledgerdomain.ParseMoney(source.TotalDebit)
	if err != nil {
		return analyticsdomain.DashboardEvent{}, err
	}

	timestamp := p.now().UTC()
	if strings.TrimSpace(source.PostedAt) != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, source.PostedAt); err == nil {
			timestamp = parsed.UTC()
		}
	}

	return analyticsdomain.DashboardEvent{
		Event:        "NEW_TRANSACTION",
		CompanyCode:  strings.TrimSpace(source.CompanyCode),
		EntryID:      firstNonEmpty(source.EntryID, event.AggregateID),
		ReferenceNo:  strings.TrimSpace(source.ReferenceNo),
		Amount:       moneyToFloat64(amount),
		Segment:      strings.ToUpper(resolveSegment(source.Metadata, source.SourceModule)),
		SourceModule: strings.ToUpper(strings.TrimSpace(source.SourceModule)),
		Timestamp:    timestamp,
	}, nil
}

func resolveSegment(metadata map[string]any, sourceModule string) string {
	for _, key := range []string{"cost_center", "business_line", "segment"} {
		if raw, ok := metadata[key]; ok {
			if value := strings.TrimSpace(fmt.Sprint(raw)); value != "" {
				return value
			}
		}
	}
	return strings.TrimSpace(sourceModule)
}

func moneyToFloat64(amount ledgerdomain.Money) float64 {
	parsed, err := strconv.ParseFloat(amount.String(), 64)
	if err != nil {
		return 0
	}
	return parsed
}
