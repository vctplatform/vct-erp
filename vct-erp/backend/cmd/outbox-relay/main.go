package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"vct-platform/backend/internal/config"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	worker := ledgerusecase.NewOutboxRelayWorker(
		nil,
		nil,
		"outbox-relay",
		cfg.OutboxRelayBatchSize,
		cfg.OutboxRelayInterval,
		cfg.OutboxRelayRetryDelay,
		cfg.OutboxRelayClaimTTL,
	)

	if err := worker.Start(signalCtx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("outbox relay: %v", err)
	}
}
