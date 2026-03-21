package usecase

import (
	"context"
	"testing"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
)

func TestRentalServiceCaptureAndReleaseDeposit(t *testing.T) {
	now := time.Date(2026, 3, 21, 11, 0, 0, 0, time.UTC)
	repo := &fakeRentalRepo{}
	ledger := &fakeLedgerPoster{}
	service := NewRentalService(&fakeFinanceTxManager{}, ledger, repo)
	service.now = func() time.Time { return now }

	captureResult, err := service.CaptureDeposit(context.Background(), CaptureRentalDepositRequest{
		CompanyCode:      "VCT_GROUP",
		RentalOrderID:    "rent-001",
		CustomerRef:      "customer-001",
		CashAccountID:    "111",
		HoldingAccountID: "3388",
		CurrencyCode:     "VND",
		Amount:           ledgerdomain.MustParseMoney("500000.0000"),
	})
	if err != nil {
		t.Fatalf("CaptureDeposit returned error: %v", err)
	}

	if got, want := len(repo.created), 1; got != want {
		t.Fatalf("expected %d created deposit, got %d", want, got)
	}
	if captureResult.DepositStatus != "held" {
		t.Fatalf("expected held status, got %s", captureResult.DepositStatus)
	}

	repo.current = repo.created[0]
	releaseResult, err := service.ReleaseDeposit(context.Background(), ReleaseRentalDepositRequest{
		CompanyCode:   "VCT_GROUP",
		RentalOrderID: "rent-001",
	})
	if err != nil {
		t.Fatalf("ReleaseDeposit returned error: %v", err)
	}

	if releaseResult.DepositStatus != "released" {
		t.Fatalf("expected released status, got %s", releaseResult.DepositStatus)
	}
	if got, want := len(repo.released), 1; got != want {
		t.Fatalf("expected %d released deposit, got %d", want, got)
	}
	if got, want := len(ledger.requests), 2; got != want {
		t.Fatalf("expected %d ledger requests, got %d", want, got)
	}
}

type fakeRentalRepo struct {
	created  []financedomain.RentalDeposit
	current  financedomain.RentalDeposit
	released []string
}

func (f *fakeRentalRepo) CreateDeposit(_ context.Context, deposit financedomain.RentalDeposit) error {
	f.created = append(f.created, deposit)
	return nil
}

func (f *fakeRentalRepo) GetByRentalOrder(_ context.Context, _ string, _ string) (financedomain.RentalDeposit, error) {
	return f.current, nil
}

func (f *fakeRentalRepo) MarkDepositSettled(_ context.Context, depositID string, _ string, _ string, _ time.Time) error {
	f.released = append(f.released, depositID)
	return nil
}
