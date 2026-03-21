package realtime

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

type capturePublisher struct {
	channel string
	payload []byte
}

func (c *capturePublisher) Publish(_ context.Context, channel string, payload []byte) error {
	c.channel = channel
	c.payload = append([]byte(nil), payload...)
	return nil
}

func TestRedisDashboardEventPublisherTransformsLedgerPayload(t *testing.T) {
	publisher := &capturePublisher{}
	realtimePublisher := NewRedisDashboardEventPublisher(publisher, "finance.dashboard.events")
	realtimePublisher.now = func() time.Time {
		return time.Date(2026, 3, 22, 9, 0, 0, 0, time.UTC)
	}

	event := ledgerdomain.OutboxEvent{
		EventType: "ledger.entry.posted",
		Payload: []byte(`{
			"entry_id":"entry-001",
			"reference_no":"PT-0001/03-26",
			"company_code":"VCT_GROUP",
			"source_module":"dojo",
			"posted_at":"2026-03-22T09:15:00Z",
			"total_debit":"5000000.0000",
			"metadata":{"cost_center":"dojo"}
		}`),
	}

	if err := realtimePublisher.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(publisher.payload, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if got, want := publisher.channel, "finance.dashboard.events"; got != want {
		t.Fatalf("expected channel %s, got %s", want, got)
	}
	if got, want := payload["event"], "NEW_TRANSACTION"; got != want {
		t.Fatalf("expected event %s, got %v", want, got)
	}
	if got, want := payload["segment"], "DOJO"; got != want {
		t.Fatalf("expected segment %s, got %v", want, got)
	}
	if got, want := payload["amount"], 5000000.0; got != want {
		t.Fatalf("expected amount %v, got %v", want, got)
	}
}
