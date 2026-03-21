package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"vct-platform/backend/internal/config"
	"vct-platform/backend/internal/httpapi"
	infraMongo "vct-platform/backend/internal/infra/mongoclient"
	infraPostgres "vct-platform/backend/internal/infra/postgres"
	infraRedis "vct-platform/backend/internal/infra/redisclient"
	analyticscache "vct-platform/backend/internal/modules/analytics/adapter/cache"
	analyticspg "vct-platform/backend/internal/modules/analytics/adapter/postgres"
	analyticsrealtime "vct-platform/backend/internal/modules/analytics/adapter/realtime"
	analyticsusecase "vct-platform/backend/internal/modules/analytics/usecase"
	financemongo "vct-platform/backend/internal/modules/finance/adapter/mongoaudit"
	financepg "vct-platform/backend/internal/modules/finance/adapter/postgres"
	financeusecase "vct-platform/backend/internal/modules/finance/usecase"
	ledgercache "vct-platform/backend/internal/modules/ledger/adapter/cache"
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

	mongoClient, err := infraMongo.Open(signalCtx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("open mongo: %v", err)
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTPShutdownTimeout)
		defer cancel()
		_ = mongoClient.Disconnect(disconnectCtx)
	}()

	ledgerStore := ledgerpg.NewStore(db)
	financeStore := financepg.NewStore(db)
	analyticsRepo := analyticspg.NewRepository(db)
	chartCache := ledgercache.NewChartOfAccountsRedisCache(redisClient, "coa")
	streamPublisher := ledgeroutbox.NewRedisStreamPublisher(redisClient)
	dashboardCache := analyticscache.NewDashboardRedisCache(redisClient, "finance:dashboard")
	dashboardPublisher := analyticsrealtime.NewRedisDashboardEventPublisher(redisClient, cfg.FinanceEventsChannel)
	eventPublisher := ledgeroutbox.NewMultiPublisher(streamPublisher, dashboardPublisher)
	financeHub := analyticsrealtime.NewHub("cfo", "ceo", "system_admin")
	financeSocket := analyticsrealtime.NewWebsocketHandler(
		signalCtx,
		financeHub,
		cfg.CORSAllowedOrigins,
		cfg.AppRoleHeader,
		cfg.AppActorHeader,
		"cfo",
		"ceo",
		"system_admin",
	)
	bridge := analyticsrealtime.NewRedisBridge(
		analyticsrealtime.NewRedisClientSubscriber(redisClient),
		dashboardCache,
		financeHub,
		cfg.FinanceEventsChannel,
	)
	auditRepo := financemongo.NewRepository(
		mongoClient.Database(cfg.MongoDatabase).Collection(cfg.MongoAuditCollection),
	)

	postEntryUC := ledgerusecase.NewPostEntryUseCase(
		ledgerStore,
		ledgerStore,
		ledgerStore,
		ledgerStore,
		ledgerStore,
		ledgerStore,
		chartCache,
		cfg.AccountCacheTTL,
		eventPublisher,
		cfg.RedisStreamKey,
	)
	saasService := financeusecase.NewSaaSService(ledgerStore, postEntryUC, financeStore)
	dojoService := financeusecase.NewDojoService(ledgerStore, postEntryUC, financeStore)
	retailService := financeusecase.NewRetailService(postEntryUC)
	rentalService := financeusecase.NewRentalService(ledgerStore, postEntryUC, financeStore)
	captureUC := financeusecase.NewCaptureUseCase(financeStore, saasService, dojoService, retailService, rentalService)
	voidUC := financeusecase.NewVoidTransactionUseCase(ledgerStore, ledgerStore, postEntryUC, auditRepo)
	analyticsService := analyticsusecase.NewService(
		analyticsRepo,
		analyticsusecase.WithDashboardCache(dashboardCache, cfg.DashboardCacheTTL),
	)

	server := &http.Server{
		Addr: cfg.HTTPAddr,
		Handler: httpapi.New(httpapi.Dependencies{
			PostEntryUC:        postEntryUC,
			FinanceCaptureUC:   captureUC,
			FinanceVoidUC:      voidUC,
			AnalyticsRevenue:   analyticsService,
			AnalyticsRunway:    analyticsService,
			AnalyticsSummary:   analyticsService,
			AnalyticsSegments:  analyticsService,
			AnalyticsCashflow:  analyticsService,
			AnalyticsDashboard: analyticsService,
			AnalyticsCards:     analyticsService,
			AnalyticsMix:       analyticsService,
			AnalyticsChart:     analyticsService,
			FinanceRealtime:    financeSocket,
			CORSAllowedOrigins: cfg.CORSAllowedOrigins,
			AppRoleHeader:      cfg.AppRoleHeader,
			AppActorHeader:     cfg.AppActorHeader,
			IdempotencyHeader:  cfg.IdempotencyHeader,
		}),
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	errCh := make(chan error, 1)
	bridgeErrCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()
	go func() {
		bridgeErrCh <- bridge.Run(signalCtx)
	}()

	log.Printf("ledger api listening on %s", cfg.HTTPAddr)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http serve: %v", err)
		}
	case err := <-bridgeErrCh:
		if err != nil && !errors.Is(err, context.Canceled) && signalCtx.Err() == nil {
			log.Fatalf("finance realtime bridge: %v", err)
		}
	case <-signalCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTPShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("shutdown: %v", err)
		}
	}
}
