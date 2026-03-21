package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"vct-platform/backend/internal/modules/ledger/domain"
)

// OutboxRelayWorker polls pending outbox rows and publishes them with at-least-once delivery semantics.
type OutboxRelayWorker struct {
	repo                domain.OutboxRelayRepository
	publisher           domain.EventPublisher
	workerID            string
	batchSize           int
	interval            time.Duration
	retryDelay          time.Duration
	staleClaimThreshold time.Duration
	runTimeout          time.Duration
	now                 func() time.Time
}

// NewOutboxRelayWorker constructs the background relay worker.
func NewOutboxRelayWorker(
	repo domain.OutboxRelayRepository,
	publisher domain.EventPublisher,
	workerID string,
	batchSize int,
	interval time.Duration,
	retryDelay time.Duration,
	staleClaimThreshold time.Duration,
	runTimeout time.Duration,
) *OutboxRelayWorker {
	if batchSize <= 0 {
		batchSize = 100
	}
	if interval <= 0 {
		interval = 5 * time.Second
	}
	if retryDelay <= 0 {
		retryDelay = 30 * time.Second
	}
	if staleClaimThreshold <= 0 {
		staleClaimThreshold = 2 * time.Minute
	}
	if runTimeout <= 0 {
		runTimeout = 20 * time.Second
	}

	return &OutboxRelayWorker{
		repo:                repo,
		publisher:           publisher,
		workerID:            workerID,
		batchSize:           batchSize,
		interval:            interval,
		retryDelay:          retryDelay,
		staleClaimThreshold: staleClaimThreshold,
		runTimeout:          runTimeout,
		now:                 time.Now,
	}
}

// Start runs the relay loop until the context is cancelled.
func (w *OutboxRelayWorker) Start(ctx context.Context) error {
	if w.repo == nil {
		return fmt.Errorf("outbox relay repository is required")
	}
	if w.publisher == nil {
		return fmt.Errorf("outbox event publisher is required")
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	if err := w.runBatch(); err != nil && ctx.Err() == nil {
		log.Printf("outbox relay run failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := w.runBatch(); err != nil {
				log.Printf("outbox relay run failed: %v", err)
			}
		}
	}
}

// RunOnce claims a batch of outbox rows and publishes them.
func (w *OutboxRelayWorker) RunOnce(ctx context.Context) error {
	if w.repo == nil {
		return fmt.Errorf("outbox relay repository is required")
	}
	if w.publisher == nil {
		return fmt.Errorf("outbox event publisher is required")
	}

	now := w.now().UTC()
	events, err := w.repo.ClaimPending(ctx, w.workerIDOrDefault(), w.batchSize, now, w.staleClaimThreshold)
	if err != nil {
		return fmt.Errorf("claim pending outbox events: %w", err)
	}

	for _, event := range events {
		if err := w.publisher.Publish(ctx, event); err != nil {
			retryAt := now.Add(w.retryDelay)
			if markErr := w.repo.MarkFailed(ctx, event.ID, now, err.Error(), retryAt); markErr != nil {
				return fmt.Errorf("mark failed outbox event %s after publish error %v: %w", event.ID, err, markErr)
			}
			continue
		}

		if err := w.repo.MarkPublished(ctx, event.ID, now); err != nil {
			return fmt.Errorf("mark published outbox event %s: %w", event.ID, err)
		}
	}

	return nil
}

func (w *OutboxRelayWorker) workerIDOrDefault() string {
	if w.workerID == "" {
		return "outbox-relay"
	}
	return w.workerID
}

func (w *OutboxRelayWorker) runBatch() error {
	if w.runTimeout <= 0 {
		return w.RunOnce(context.Background())
	}

	runCtx, cancel := context.WithTimeout(context.Background(), w.runTimeout)
	defer cancel()
	return w.RunOnce(runCtx)
}
