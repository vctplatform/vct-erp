package usecase

import (
	"fmt"

	"vct-platform/backend/internal/modules/ledger/domain"
)

// TrialBalanceRow mirrors the core columns of the VAS trial balance report.
type TrialBalanceRow struct {
	AccountCode    string
	AccountName    string
	NormalSide     domain.Side
	OpeningBalance domain.Money
	PeriodDebit    domain.Money
	PeriodCredit   domain.Money
	ClosingBalance domain.Money
}

// TrialBalanceSummary aggregates the report totals used for validation and exports.
type TrialBalanceSummary struct {
	OpeningDebit  domain.Money
	OpeningCredit domain.Money
	PeriodDebit   domain.Money
	PeriodCredit  domain.Money
	ClosingDebit  domain.Money
	ClosingCredit domain.Money
}

// ClosingBalance computes the natural ending balance for one account row.
func ClosingBalance(normalSide domain.Side, opening domain.Money, periodDebit domain.Money, periodCredit domain.Money) domain.Money {
	if normalSide == domain.SideCredit {
		return opening.Sub(periodDebit).Add(periodCredit)
	}
	return opening.Add(periodDebit).Sub(periodCredit)
}

// SummarizeTrialBalance validates each row and produces debit/credit presentation totals.
func SummarizeTrialBalance(rows []TrialBalanceRow) (TrialBalanceSummary, error) {
	var summary TrialBalanceSummary

	for _, row := range rows {
		expectedClosing := ClosingBalance(row.NormalSide, row.OpeningBalance, row.PeriodDebit, row.PeriodCredit)
		if !expectedClosing.Equal(row.ClosingBalance) {
			return TrialBalanceSummary{}, fmt.Errorf(
				"trial balance row %s has invalid closing balance: got %s want %s",
				row.AccountCode,
				row.ClosingBalance.String(),
				expectedClosing.String(),
			)
		}

		if row.NormalSide == domain.SideCredit {
			summary.OpeningCredit = summary.OpeningCredit.Add(row.OpeningBalance)
			summary.ClosingCredit = summary.ClosingCredit.Add(row.ClosingBalance)
		} else {
			summary.OpeningDebit = summary.OpeningDebit.Add(row.OpeningBalance)
			summary.ClosingDebit = summary.ClosingDebit.Add(row.ClosingBalance)
		}

		summary.PeriodDebit = summary.PeriodDebit.Add(row.PeriodDebit)
		summary.PeriodCredit = summary.PeriodCredit.Add(row.PeriodCredit)
	}

	return summary, nil
}
