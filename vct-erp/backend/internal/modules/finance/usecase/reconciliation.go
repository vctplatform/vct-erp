package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
	"vct-platform/backend/internal/shared/repository"
)

// ReconciliationService matches imported bank statement lines with ledger account 112 movements.
type ReconciliationService struct {
	txManager financedomain.TxManager
	repo      financedomain.BankReconciliationRepository
	now       func() time.Time
}

// NewReconciliationService constructs the bank reconciliation application service.
func NewReconciliationService(txManager financedomain.TxManager, repo financedomain.BankReconciliationRepository) *ReconciliationService {
	return &ReconciliationService{
		txManager: txManager,
		repo:      repo,
		now:       time.Now,
	}
}

// ReconcileBankAccount performs exact amount/date matching with optional reference boosts.
func (s *ReconciliationService) ReconcileBankAccount(ctx context.Context, req financedomain.ReconcileBankRequest) (*financedomain.ReconcileBankResult, error) {
	if strings.TrimSpace(req.CompanyCode) == "" {
		return nil, financedomain.ErrCompanyRequired
	}
	if strings.TrimSpace(req.LedgerAccountID) == "" {
		return nil, financedomain.ErrLedgerAccountRequired
	}
	if strings.TrimSpace(req.BankAccountNo) == "" {
		return nil, financedomain.ErrBankAccountRequired
	}

	dateFrom := req.DateFrom.UTC()
	dateTo := req.DateTo.UTC()
	if dateFrom.IsZero() {
		dateFrom = time.Now().UTC().AddDate(0, 0, -31)
	}
	if dateTo.IsZero() {
		dateTo = time.Now().UTC()
	}
	if dateTo.Before(dateFrom) {
		dateFrom, dateTo = dateTo, dateFrom
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 500
	}

	driftDays := req.MaxDateDriftDays
	if driftDays <= 0 {
		driftDays = 3
	}

	statementLines, err := s.repo.ListOpenStatementLines(ctx, strings.TrimSpace(req.CompanyCode), strings.TrimSpace(req.BankAccountNo), dateFrom, dateTo, limit)
	if err != nil {
		return nil, fmt.Errorf("list open statement lines: %w", err)
	}

	ledgerLines, err := s.repo.ListOpenLedgerMovements(ctx, strings.TrimSpace(req.CompanyCode), strings.TrimSpace(req.LedgerAccountID), dateFrom.AddDate(0, 0, -driftDays), dateTo.AddDate(0, 0, driftDays), limit)
	if err != nil {
		return nil, fmt.Errorf("list open ledger movements: %w", err)
	}

	matches := make([]financedomain.BankMatchResult, 0, len(statementLines))
	usedLedger := make(map[string]struct{}, len(ledgerLines))
	now := s.now().UTC()

	if err := s.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationRepeatableRead,
	}, func(txCtx context.Context) error {
		for _, statement := range statementLines {
			index, rule := findBestLedgerMatch(statement, ledgerLines, usedLedger, driftDays)
			if index < 0 {
				continue
			}

			movement := ledgerLines[index]
			if err := s.repo.MarkStatementMatched(txCtx, statement.ID, movement.JournalEntryID, rule, now); err != nil {
				return fmt.Errorf("mark statement line %s matched: %w", statement.ID, err)
			}

			usedLedger[movement.JournalEntryID] = struct{}{}
			matches = append(matches, financedomain.BankMatchResult{
				StatementLineID: statement.ID,
				JournalEntryID:  movement.JournalEntryID,
				EntryNo:         movement.EntryNo,
				MatchingRule:    rule,
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &financedomain.ReconcileBankResult{
		MatchedCount:            len(matches),
		UnmatchedStatementCount: len(statementLines) - len(matches),
		UnmatchedLedgerCount:    len(ledgerLines) - len(matches),
		Matches:                 matches,
	}, nil
}

func findBestLedgerMatch(
	statement financedomain.BankStatementLine,
	ledgerLines []financedomain.LedgerBankMovement,
	usedLedger map[string]struct{},
	maxDateDriftDays int,
) (int, string) {
	statementAmount := statement.Amount.Abs()
	expectedSide := ledgerdomain.SideDebit
	if statement.Amount.Sign() < 0 {
		expectedSide = ledgerdomain.SideCredit
	}

	bestIndex := -1
	bestScore := -1
	for index, movement := range ledgerLines {
		if _, alreadyUsed := usedLedger[movement.JournalEntryID]; alreadyUsed {
			continue
		}
		if movement.Side != expectedSide {
			continue
		}
		if !movement.Amount.Equal(statementAmount) {
			continue
		}
		if dateDriftInDays(statement.BookingDate, movement.PostingDate) > maxDateDriftDays {
			continue
		}

		score := 1
		rule := "exact_amount_side_date"
		if referenceMatches(statement.ReferenceNo, movement.ExternalRef) {
			score = 3
			rule = "exact_amount_ref_date"
		} else if referenceMatches(statement.Description, movement.Description) {
			score = 2
			rule = "exact_amount_desc_date"
		}

		if score > bestScore {
			bestScore = score
			bestIndex = index
			if rule == "exact_amount_ref_date" {
				return bestIndex, rule
			}
		}
	}

	if bestIndex < 0 {
		return -1, ""
	}

	switch bestScore {
	case 2:
		return bestIndex, "exact_amount_desc_date"
	default:
		return bestIndex, "exact_amount_side_date"
	}
}

func referenceMatches(left string, right string) bool {
	left = strings.ToLower(strings.TrimSpace(left))
	right = strings.ToLower(strings.TrimSpace(right))
	if left == "" || right == "" {
		return false
	}
	return left == right || strings.Contains(left, right) || strings.Contains(right, left)
}

func dateDriftInDays(left time.Time, right time.Time) int {
	left = time.Date(left.Year(), left.Month(), left.Day(), 0, 0, 0, 0, time.UTC)
	right = time.Date(right.Year(), right.Month(), right.Day(), 0, 0, 0, 0, time.UTC)
	if left.After(right) {
		left, right = right, left
	}
	return int(right.Sub(left).Hours() / 24)
}
