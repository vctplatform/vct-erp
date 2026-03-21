package usecase

import (
	"context"
	"testing"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
	"vct-platform/backend/internal/shared/repository"
)

func TestSaaSServiceCaptureAnnualContract(t *testing.T) {
	now := time.Date(2026, 3, 21, 9, 0, 0, 0, time.UTC)
	repo := &fakeSaaSRepo{}
	ledger := &fakeLedgerPoster{}
	service := NewSaaSService(&fakeFinanceTxManager{}, ledger, repo)
	service.now = func() time.Time { return now }

	result, err := service.CaptureAnnualContract(context.Background(), CaptureAnnualContractRequest{
		CompanyCode:                "VCT_GROUP",
		ContractNo:                 "S-2026-001",
		CustomerRef:                "cust-001",
		CashAccountID:              "111",
		DeferredRevenueAccountID:   "3387",
		RecognizedRevenueAccountID: "5113",
		CurrencyCode:               "VND",
		ServiceStartDate:           now,
		TermMonths:                 12,
		TotalAmount:                ledgerdomain.MustParseMoney("12000000.0000"),
	})
	if err != nil {
		t.Fatalf("CaptureAnnualContract returned error: %v", err)
	}

	if got, want := len(repo.contracts), 1; got != want {
		t.Fatalf("expected %d contract, got %d", want, got)
	}
	if got, want := len(repo.schedules), 12; got != want {
		t.Fatalf("expected %d schedules, got %d", want, got)
	}
	if got, want := len(ledger.requests), 1; got != want {
		t.Fatalf("expected %d ledger post, got %d", want, got)
	}
	if result.ScheduleCount != 12 {
		t.Fatalf("expected schedule count 12, got %d", result.ScheduleCount)
	}
	if result.MonthlyAmounts[0] != "1000000.0000" {
		t.Fatalf("unexpected monthly amount: %s", result.MonthlyAmounts[0])
	}
}

func TestSaaSServiceRecognizeDueRevenue(t *testing.T) {
	now := time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC)
	repo := &fakeSaaSRepo{
		dueSchedules: []financedomain.DueRevenueSchedule{
			{
				ScheduleID:                 "sch-1",
				ContractID:                 "contract-1",
				ContractNo:                 "S-2026-001",
				CompanyCode:                "VCT_GROUP",
				CustomerRef:                "cust-001",
				CurrencyCode:               "VND",
				ServiceMonth:               time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
				Amount:                     ledgerdomain.MustParseMoney("1000000.0000"),
				DeferredRevenueAccountID:   "3387",
				RecognizedRevenueAccountID: "5113",
			},
		},
	}
	ledger := &fakeLedgerPoster{}
	service := NewSaaSService(&fakeFinanceTxManager{}, ledger, repo)
	service.now = func() time.Time { return now }

	result, err := service.RecognizeDueRevenue(context.Background(), RecognizeDueRevenueRequest{
		CompanyCode: "VCT_GROUP",
		UpTo:        now,
	})
	if err != nil {
		t.Fatalf("RecognizeDueRevenue returned error: %v", err)
	}

	if result.RecognizedCount != 1 {
		t.Fatalf("expected recognized count 1, got %d", result.RecognizedCount)
	}
	if got, want := len(repo.markRecognizedCalls), 1; got != want {
		t.Fatalf("expected %d mark recognized call, got %d", want, got)
	}
}

type fakeFinanceTxManager struct{}

func (f *fakeFinanceTxManager) WithinTransaction(ctx context.Context, _ repository.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type fakeLedgerPoster struct {
	requests []ledgerusecase.PostEntryRequest
}

func (f *fakeLedgerPoster) PostEntry(_ context.Context, req ledgerusecase.PostEntryRequest) (*ledgerusecase.PostEntryResult, error) {
	f.requests = append(f.requests, req)
	return &ledgerusecase.PostEntryResult{
		Entry: ledgerdomain.JournalEntry{ID: "je-" + req.SourceModule},
	}, nil
}

type fakeSaaSRepo struct {
	contracts           []financedomain.SaaSContract
	schedules           []financedomain.RevenueSchedule
	dueSchedules        []financedomain.DueRevenueSchedule
	markRecognizedCalls []string
}

func (f *fakeSaaSRepo) CreateContract(_ context.Context, contract financedomain.SaaSContract) error {
	f.contracts = append(f.contracts, contract)
	return nil
}

func (f *fakeSaaSRepo) CreateSchedules(_ context.Context, schedules []financedomain.RevenueSchedule) error {
	f.schedules = append(f.schedules, schedules...)
	return nil
}

func (f *fakeSaaSRepo) ListDueSchedules(_ context.Context, _ string, _ time.Time, _ int) ([]financedomain.DueRevenueSchedule, error) {
	return append([]financedomain.DueRevenueSchedule(nil), f.dueSchedules...), nil
}

func (f *fakeSaaSRepo) MarkScheduleRecognized(_ context.Context, scheduleID string, _ string, _ time.Time) error {
	f.markRecognizedCalls = append(f.markRecognizedCalls, scheduleID)
	return nil
}
