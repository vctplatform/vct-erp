package usecase

import (
	"context"
	"testing"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

type fakeDashboardCache struct {
	snapshot analyticsdomain.CommandCenterDashboardData
	sets     int
	miss     bool
}

func (f *fakeDashboardCache) Get(_ context.Context, _ string) (analyticsdomain.CommandCenterDashboardData, error) {
	if f.miss {
		return analyticsdomain.CommandCenterDashboardData{}, analyticsdomain.ErrDashboardCacheMiss
	}
	return f.snapshot, nil
}

func (f *fakeDashboardCache) Set(_ context.Context, snapshot analyticsdomain.CommandCenterDashboardData, _ time.Duration) error {
	f.snapshot = snapshot
	f.sets++
	return nil
}

func (f *fakeDashboardCache) Delete(_ context.Context, _ string) error {
	f.snapshot = analyticsdomain.CommandCenterDashboardData{}
	return nil
}

func TestDashboardBuildsLiveContractAndCachesSnapshot(t *testing.T) {
	now := time.Date(2026, 3, 22, 9, 0, 0, 0, time.UTC)
	repo := &fakeRepository{
		reconcile: []analyticsdomain.ReconciliationBalance{
			{AccountCode: "1111", EffectiveNetBalance: ledgerdomain.MustParseMoney("500.0000")},
			{AccountCode: "1121", EffectiveNetBalance: ledgerdomain.MustParseMoney("300.0000")},
		},
		segments: []analyticsdomain.SegmentGrossProfit{
			{SegmentKey: "saas", SegmentLabel: "SaaS", NetRevenue: ledgerdomain.MustParseMoney("400.0000")},
			{SegmentKey: "dojo", SegmentLabel: "Dojo", NetRevenue: ledgerdomain.MustParseMoney("250.0000")},
			{SegmentKey: "retail", SegmentLabel: "Retail", NetRevenue: ledgerdomain.MustParseMoney("175.0000")},
			{SegmentKey: "rental", SegmentLabel: "Rental", NetRevenue: ledgerdomain.MustParseMoney("50.0000")},
		},
		runway: analyticsdomain.CashRunway{
			CurrentCash:        ledgerdomain.MustParseMoney("800.0000"),
			AverageMonthlyBurn: ledgerdomain.MustParseMoney("100.0000"),
			Months: []analyticsdomain.CashRunwayMonth{
				{MonthLabel: "2026-04", OpeningCash: ledgerdomain.MustParseMoney("800.0000"), ContractedInflow: ledgerdomain.MustParseMoney("50.0000"), ProjectedBurn: ledgerdomain.MustParseMoney("100.0000"), ProjectedEnding: ledgerdomain.MustParseMoney("750.0000")},
				{MonthLabel: "2026-05", OpeningCash: ledgerdomain.MustParseMoney("750.0000"), ContractedInflow: ledgerdomain.MustParseMoney("50.0000"), ProjectedBurn: ledgerdomain.MustParseMoney("100.0000"), ProjectedEnding: ledgerdomain.MustParseMoney("700.0000")},
			},
		},
		seriesBySize: map[int][]analyticsdomain.FinanceSeriesPoint{
			7: buildSeries(7),
			6: buildSeries(6),
		},
	}
	cache := &fakeDashboardCache{miss: true}
	service := NewService(repo, WithDashboardCache(cache, time.Minute))
	service.now = func() time.Time { return now }

	dashboard, err := service.Dashboard(context.Background(), DashboardInput{
		Access: AccessMetadata{
			CompanyCode: "VCT_GROUP",
			ActorID:     "ceo-001",
			ActorRole:   "ceo",
		},
		AsOf: now,
	})
	if err != nil {
		t.Fatalf("Dashboard returned error: %v", err)
	}

	if dashboard.DataMode != "live" {
		t.Fatalf("expected data mode live, got %s", dashboard.DataMode)
	}
	if got, want := len(dashboard.Cards), 3; got != want {
		t.Fatalf("expected %d cards, got %d", want, got)
	}
	if got, want := len(dashboard.RevenueMix), 4; got != want {
		t.Fatalf("expected %d revenue slices, got %d", want, got)
	}
	if got, want := len(dashboard.CashflowChart.XAxis), 6; got != want {
		t.Fatalf("expected %d cashflow points, got %d", want, got)
	}
	if got, want := len(dashboard.RunwayProjection), 2; got != want {
		t.Fatalf("expected %d runway projections, got %d", want, got)
	}
	if got, want := cache.sets, 1; got != want {
		t.Fatalf("expected cache set count %d, got %d", want, got)
	}
	if len(repo.accesses) != 1 || repo.accesses[0].ReportCode != "finance.dashboard.live" {
		t.Fatalf("expected finance.dashboard.live audit access, got %+v", repo.accesses)
	}
	if dashboard.Cards[0].Value != 800 {
		t.Fatalf("expected cash card value 800, got %f", dashboard.Cards[0].Value)
	}
}

func TestDashboardUsesCachedSnapshotBeforeHittingRepository(t *testing.T) {
	cached := analyticsdomain.CommandCenterDashboardData{
		CompanyCode:        "VCT_GROUP",
		DataMode:           "live",
		RecommendedRefresh: "websocket",
		Cards: []analyticsdomain.DashboardCard{
			{Key: "cash_assets", Value: 42},
		},
	}
	cache := &fakeDashboardCache{snapshot: cached}
	repo := &fakeRepository{}
	service := NewService(repo, WithDashboardCache(cache, time.Minute))
	service.now = func() time.Time {
		return time.Date(2026, 3, 22, 9, 0, 0, 0, time.UTC)
	}

	result, err := service.Dashboard(context.Background(), DashboardInput{
		Access: AccessMetadata{
			CompanyCode: "VCT_GROUP",
			ActorID:     "cfo-001",
			ActorRole:   "cfo",
		},
	})
	if err != nil {
		t.Fatalf("Dashboard returned error: %v", err)
	}

	if got, want := result.Cards[0].Value, cached.Cards[0].Value; got != want {
		t.Fatalf("expected cached card value %f, got %f", want, got)
	}
	if got, want := len(repo.accesses), 1; got != want {
		t.Fatalf("expected %d audit entry, got %d", want, got)
	}
}

func buildSeries(points int) []analyticsdomain.FinanceSeriesPoint {
	series := make([]analyticsdomain.FinanceSeriesPoint, 0, points)
	start := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)
	for index := 0; index < points; index++ {
		series = append(series, analyticsdomain.FinanceSeriesPoint{
			PeriodStart: start.AddDate(0, index, 0),
			PeriodLabel: start.AddDate(0, index, 0).Format("2006-01"),
			CashEnding:  ledgerdomain.MustParseMoney("1000.0000").Add(ledgerdomain.MustParseMoney("100.0000").MulInt(index)),
			Revenue:     ledgerdomain.MustParseMoney("300.0000").Add(ledgerdomain.MustParseMoney("10.0000").MulInt(index)),
			Expense:     ledgerdomain.MustParseMoney("120.0000").Add(ledgerdomain.MustParseMoney("5.0000").MulInt(index)),
			Profit:      ledgerdomain.MustParseMoney("180.0000").Add(ledgerdomain.MustParseMoney("5.0000").MulInt(index)),
		})
	}
	return series
}
