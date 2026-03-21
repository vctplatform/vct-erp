package usecase

import (
	"context"
	"testing"
	"time"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

type fakeRepository struct {
	summary      analyticsdomain.FinanceSummary
	segments     []analyticsdomain.SegmentGrossProfit
	runway       analyticsdomain.CashRunway
	reconcile    []analyticsdomain.ReconciliationBalance
	seriesBySize map[int][]analyticsdomain.FinanceSeriesPoint
	accesses     []analyticsdomain.ReportAccessLog
	summaryCode  string
	segmentCode  string
	runwayCode   string
	runwayAsOf   time.Time
	runwayMonths int
	accessErr    error
}

func (f *fakeRepository) FinanceSummary(_ context.Context, companyCode string) (analyticsdomain.FinanceSummary, error) {
	f.summaryCode = companyCode
	return f.summary, nil
}

func (f *fakeRepository) Segments(_ context.Context, companyCode string) ([]analyticsdomain.SegmentGrossProfit, error) {
	f.segmentCode = companyCode
	return f.segments, nil
}

func (f *fakeRepository) RevenueStream(_ context.Context, _ string, _ time.Time, _ time.Time) ([]analyticsdomain.RevenueStreamPoint, error) {
	return nil, nil
}

func (f *fakeRepository) CashRunway(_ context.Context, companyCode string, asOf time.Time, months int) (analyticsdomain.CashRunway, error) {
	f.runwayCode = companyCode
	f.runwayAsOf = asOf
	f.runwayMonths = months
	return f.runway, nil
}

func (f *fakeRepository) ReconciliationBalances(_ context.Context, _ string) ([]analyticsdomain.ReconciliationBalance, error) {
	return append([]analyticsdomain.ReconciliationBalance(nil), f.reconcile...), nil
}

func (f *fakeRepository) FinanceSeries(_ context.Context, _ string, _ time.Time, months int) ([]analyticsdomain.FinanceSeriesPoint, error) {
	if f.seriesBySize == nil {
		return nil, nil
	}
	return append([]analyticsdomain.FinanceSeriesPoint(nil), f.seriesBySize[months]...), nil
}

func (f *fakeRepository) RecordReportAccess(_ context.Context, access analyticsdomain.ReportAccessLog) error {
	if f.accessErr != nil {
		return f.accessErr
	}
	f.accesses = append(f.accesses, access)
	return nil
}

func TestFinanceSummaryDefaultsCompanyAndRecordsAudit(t *testing.T) {
	repo := &fakeRepository{
		summary: analyticsdomain.FinanceSummary{
			CompanyCode:       "VCT_GROUP",
			TotalRevenue:      ledgerdomain.MustParseMoney("1000.0000"),
			GrossProfit:       ledgerdomain.MustParseMoney("450.0000"),
			GrossProfitMargin: 45,
			NetCash:           ledgerdomain.MustParseMoney("800.0000"),
			CurrencyCode:      "VND",
		},
	}
	service := NewService(repo)
	service.now = func() time.Time {
		return time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	}

	result, err := service.FinanceSummary(context.Background(), AccessMetadata{
		ActorID:   "ceo-001",
		ActorRole: "ceo",
		IPAddress: "127.0.0.1",
		Filters:   map[string]string{"company_code": "VCT_GROUP"},
	})
	if err != nil {
		t.Fatalf("FinanceSummary returned error: %v", err)
	}
	if repo.summaryCode != "VCT_GROUP" {
		t.Fatalf("expected default company code VCT_GROUP, got %s", repo.summaryCode)
	}
	if len(repo.accesses) != 1 {
		t.Fatalf("expected one audit access log, got %d", len(repo.accesses))
	}
	if repo.accesses[0].ReportCode != "finance.summary" {
		t.Fatalf("unexpected report code %s", repo.accesses[0].ReportCode)
	}
	if !result.NetCash.Equal(ledgerdomain.MustParseMoney("800.0000")) {
		t.Fatalf("unexpected net cash %s", result.NetCash.String())
	}
}

func TestDashboardCashRunwayDefaultsMonthsAndAudits(t *testing.T) {
	repo := &fakeRepository{
		runway: analyticsdomain.CashRunway{
			CurrentCash:        ledgerdomain.MustParseMoney("1000.0000"),
			AverageMonthlyBurn: ledgerdomain.MustParseMoney("100.0000"),
		},
	}
	service := NewService(repo)
	service.now = func() time.Time {
		return time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	}

	_, err := service.DashboardCashRunway(context.Background(), CashRunwayInput{
		Access: AccessMetadata{
			CompanyCode: "VCT_SIM",
			ActorID:     "cfo-001",
			ActorRole:   "cfo",
		},
	})
	if err != nil {
		t.Fatalf("DashboardCashRunway returned error: %v", err)
	}
	if repo.runwayCode != "VCT_SIM" {
		t.Fatalf("expected company code VCT_SIM, got %s", repo.runwayCode)
	}
	if repo.runwayMonths != 6 {
		t.Fatalf("expected default 6 months, got %d", repo.runwayMonths)
	}
	if len(repo.accesses) != 1 || repo.accesses[0].ReportCode != "finance.cash_runway" {
		t.Fatalf("expected finance.cash_runway audit log, got %+v", repo.accesses)
	}
}
