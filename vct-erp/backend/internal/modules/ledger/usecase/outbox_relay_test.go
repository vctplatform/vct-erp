package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"vct-platform/backend/internal/modules/ledger/domain"
)

func TestOutboxRelayWorkerRunOnce(t *testing.T) {
	now := time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC)
	repo := &fakeOutboxRelayRepo{
		claimed: []domain.OutboxEvent{
			{ID: "event-1", StreamKey: "ledger.events", Payload: []byte(`{"ok":true}`), Status: domain.OutboxStatusProcessing, CreatedAt: now},
			{ID: "event-2", StreamKey: "ledger.events", Payload: []byte(`{"ok":true}`), Status: domain.OutboxStatusProcessing, CreatedAt: now},
		},
	}
	publisher := &fakeOutboxPublisher{failEventID: "event-1"}
	worker := NewOutboxRelayWorker(repo, publisher, "worker-1", 10, time.Second, time.Minute, time.Minute, 5*time.Second)
	worker.now = func() time.Time { return now }

	if err := worker.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce returned error: %v", err)
	}

	if got, want := len(repo.failed), 1; got != want {
		t.Fatalf("expected %d failed event, got %d", want, got)
	}
	if got, want := len(repo.published), 1; got != want {
		t.Fatalf("expected %d published event, got %d", want, got)
	}
}

type fakeOutboxRelayRepo struct {
	claimed   []domain.OutboxEvent
	published []string
	failed    []string
}

func (f *fakeOutboxRelayRepo) ClaimPending(_ context.Context, _ string, _ int, _ time.Time, _ time.Duration) ([]domain.OutboxEvent, error) {
	return append([]domain.OutboxEvent(nil), f.claimed...), nil
}

func (f *fakeOutboxRelayRepo) MarkPublished(_ context.Context, eventID string, _ time.Time) error {
	f.published = append(f.published, eventID)
	return nil
}

func (f *fakeOutboxRelayRepo) MarkFailed(_ context.Context, eventID string, _ time.Time, _ string, _ time.Time) error {
	f.failed = append(f.failed, eventID)
	return nil
}

type fakeOutboxPublisher struct {
	failEventID string
}

func (f *fakeOutboxPublisher) Publish(_ context.Context, event domain.OutboxEvent) error {
	if event.ID == f.failEventID {
		return errors.New("redis unavailable")
	}
	return nil
}
