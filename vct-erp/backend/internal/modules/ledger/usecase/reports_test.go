package usecase

import (
	"testing"

	"vct-platform/backend/internal/modules/ledger/domain"
)

func TestSummarizeTrialBalance(t *testing.T) {
	rows := []TrialBalanceRow{
		{
			AccountCode:    "1121",
			AccountName:    "Tien gui ngan hang VND",
			NormalSide:     domain.SideDebit,
			OpeningBalance: domain.MustParseMoney("10000000.0000"),
			PeriodDebit:    domain.MustParseMoney("3000000.0000"),
			PeriodCredit:   domain.MustParseMoney("500000.0000"),
			ClosingBalance: domain.MustParseMoney("12500000.0000"),
		},
		{
			AccountCode:    "5113",
			AccountName:    "Doanh thu cung cap dich vu",
			NormalSide:     domain.SideCredit,
			OpeningBalance: domain.MustParseMoney("0.0000"),
			PeriodDebit:    domain.MustParseMoney("0.0000"),
			PeriodCredit:   domain.MustParseMoney("2500000.0000"),
			ClosingBalance: domain.MustParseMoney("2500000.0000"),
		},
	}

	summary, err := SummarizeTrialBalance(rows)
	if err != nil {
		t.Fatalf("SummarizeTrialBalance returned error: %v", err)
	}

	if !summary.OpeningDebit.Equal(domain.MustParseMoney("10000000.0000")) {
		t.Fatalf("unexpected opening debit total: %s", summary.OpeningDebit.String())
	}
	if !summary.OpeningCredit.Equal(domain.MustParseMoney("0.0000")) {
		t.Fatalf("unexpected opening credit total: %s", summary.OpeningCredit.String())
	}
	if !summary.PeriodDebit.Equal(domain.MustParseMoney("3000000.0000")) {
		t.Fatalf("unexpected period debit total: %s", summary.PeriodDebit.String())
	}
	if !summary.PeriodCredit.Equal(domain.MustParseMoney("3000000.0000")) {
		t.Fatalf("unexpected period credit total: %s", summary.PeriodCredit.String())
	}
	if !summary.ClosingDebit.Equal(domain.MustParseMoney("12500000.0000")) {
		t.Fatalf("unexpected closing debit total: %s", summary.ClosingDebit.String())
	}
	if !summary.ClosingCredit.Equal(domain.MustParseMoney("2500000.0000")) {
		t.Fatalf("unexpected closing credit total: %s", summary.ClosingCredit.String())
	}
}

func TestSummarizeTrialBalanceRejectsInvalidClosingBalance(t *testing.T) {
	_, err := SummarizeTrialBalance([]TrialBalanceRow{
		{
			AccountCode:    "3387",
			AccountName:    "Doanh thu chua thuc hien",
			NormalSide:     domain.SideCredit,
			OpeningBalance: domain.MustParseMoney("12000000.0000"),
			PeriodDebit:    domain.MustParseMoney("1000000.0000"),
			PeriodCredit:   domain.MustParseMoney("0.0000"),
			ClosingBalance: domain.MustParseMoney("12000000.0000"),
		},
	})
	if err == nil {
		t.Fatal("expected invalid closing balance error")
	}
}
