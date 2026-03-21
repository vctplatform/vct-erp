package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"vct-platform/backend/internal/config"
	infraPostgres "vct-platform/backend/internal/infra/postgres"
	infraRedis "vct-platform/backend/internal/infra/redisclient"
	ledgeroutbox "vct-platform/backend/internal/modules/ledger/adapter/outbox"
	ledgerpg "vct-platform/backend/internal/modules/ledger/adapter/postgres"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := infraPostgres.Open(signalCtx, cfg.DatabaseDSN, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	redisClient := infraRedis.New(cfg.RedisAddr, cfg.RedisUsername, cfg.RedisPassword, cfg.RedisDB)
	if err := redisClient.Ping(signalCtx); err != nil {
		log.Fatalf("ping redis: %v", err)
	}
	defer redisClient.Close()

	ledgerStore := ledgerpg.NewStore(db)
	publisher := ledgeroutbox.NewRedisStreamPublisher(redisClient)
	worker := ledgerusecase.NewOutboxRelayWorker(
		ledgerStore,
		publisher,
		workerID(),
		cfg.OutboxRelayBatchSize,
		cfg.OutboxRelayInterval,
		cfg.OutboxRelayRetryDelay,
		cfg.OutboxRelayClaimTTL,
		cfg.OutboxRelayRunTimeout,
	)

	log.Printf("outbox relay worker %s started", workerID())
	if err := worker.Start(signalCtx); err != nil {
		log.Fatalf("outbox relay: %v", err)
	}
}

func workerID() string {
	if podName := os.Getenv("POD_NAME"); podName != "" {
		return podName
	}
	if host := os.Getenv("HOSTNAME"); host != "" {
		return host
	}
	return "outbox-relay"
}
