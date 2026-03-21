package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
	"vct-platform/backend/internal/shared/id"
	"vct-platform/backend/internal/shared/repository"
)

// DojoService automates tuition receivables for the martial-arts business.
type DojoService struct {
	txManager    financedomain.TxManager
	ledgerPoster LedgerPoster
	repo         financedomain.DojoReceivableRepository
	now          func() time.Time
}

// NewDojoService constructs the dojo accounting service.
func NewDojoService(
	txManager financedomain.TxManager,
	ledgerPoster LedgerPoster,
	repo financedomain.DojoReceivableRepository,
) *DojoService {
	return &DojoService{
		txManager:    txManager,
		ledgerPoster: ledgerPoster,
		repo:         repo,
		now:          time.Now,
	}
}

// AssessMonthlyTuition raises a receivable on the first day of the billing month.
func (s *DojoService) AssessMonthlyTuition(ctx context.Context, req AssessMonthlyTuitionRequest) (*AssessMonthlyTuitionResult, error) {
	if err := validateDojoAssessmentRequest(req); err != nil {
		return nil, err
	}

	now := s.now().UTC()
	billingMonth := monthStart(req.BillingMonth)
	receivable := financedomain.DojoReceivable{
		ID:                  id.NewUUID(),
		CompanyCode:         strings.TrimSpace(req.CompanyCode),
		StudentRef:          strings.TrimSpace(req.StudentRef),
		BillingMonth:        billingMonth,
		DueDate:             req.DueDate.UTC(),
		CurrencyCode:        normalizeCurrency(req.CurrencyCode),
		ReceivableAccountID: strings.TrimSpace(req.ReceivableAccountID),
		RevenueAccountID:    strings.TrimSpace(req.RevenueAccountID),
		AmountDue:           req.Amount,
		AmountPaid:          req.Amount.Sub(req.Amount),
		Status:              "open",
		SourceRef:           strings.TrimSpace(req.SourceRef),
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	result := &AssessMonthlyTuitionResult{
		ReceivableID:     receivable.ID,
		ReceivableStatus: receivable.Status,
	}

	err := s.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationSerializable,
	}, func(txCtx context.Context) error {
		postResult, err := s.ledgerPoster.PostEntry(txCtx, ledgerusecase.PostEntryRequest{
			CompanyCode:  receivable.CompanyCode,
			SourceModule: "dojo",
			ExternalRef:  firstNonEmpty(receivable.SourceRef, receivable.StudentRef),
			Description:  fmt.Sprintf("Tinh no hoc phi thang %s cho hoc vien %s", receivable.BillingMonth.Format("2006-01"), receivable.StudentRef),
			CurrencyCode: receivable.CurrencyCode,
			PostingDate:  receivable.BillingMonth,
			Metadata: map[string]any{
				"business_line": "dojo",
				"student_ref":   receivable.StudentRef,
				"billing_month": receivable.BillingMonth.Format("2006-01-02"),
				"receivable_id": receivable.ID,
			},
			Items: []ledgerusecase.PostEntryItemRequest{
				{
					AccountID: receivable.ReceivableAccountID,
					Side:      "debit",
					Amount:    receivable.AmountDue,
				},
				{
					AccountID: receivable.RevenueAccountID,
					Side:      "credit",
					Amount:    receivable.AmountDue,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("post dojo receivable: %w", err)
		}

		receivable.AssessmentEntryID = postResult.Entry.ID
		result.JournalEntryID = postResult.Entry.ID
		if err := s.repo.CreateReceivable(txCtx, receivable); err != nil {
			return fmt.Errorf("create dojo receivable: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CapturePayment clears a dojo receivable when cash is collected.
func (s *DojoService) CapturePayment(ctx context.Context, req CaptureDojoPaymentRequest) (*CaptureDojoPaymentResult, error) {
	if err := validateDojoPaymentRequest(req); err != nil {
		return nil, err
	}

	receivable, err := s.repo.GetReceivable(ctx, strings.TrimSpace(req.CompanyCode), strings.TrimSpace(req.StudentRef), monthStart(req.BillingMonth))
	if err != nil {
		return nil, fmt.Errorf("get dojo receivable: %w", err)
	}
	if receivable.ID == "" {
		return nil, financedomain.ErrReceivableNotFound
	}

	remaining := receivable.AmountDue.Sub(receivable.AmountPaid)
	if req.PaymentAmount.Cmp(remaining) > 0 {
		return nil, financedomain.ErrAmountExceedsBalance
	}

	paidAt := s.now().UTC()
	result := &CaptureDojoPaymentResult{
		ReceivableID: receivable.ID,
		Status:       "open",
	}

	err = s.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationSerializable,
	}, func(txCtx context.Context) error {
		postResult, err := s.ledgerPoster.PostEntry(txCtx, ledgerusecase.PostEntryRequest{
			CompanyCode:  receivable.CompanyCode,
			SourceModule: "dojo",
			ExternalRef:  firstNonEmpty(strings.TrimSpace(req.SourceRef), receivable.StudentRef),
			Description:  fmt.Sprintf("Thu hoc phi thang %s cho hoc vien %s", receivable.BillingMonth.Format("2006-01"), receivable.StudentRef),
			CurrencyCode: receivable.CurrencyCode,
			PostingDate:  paidAt,
			Metadata: map[string]any{
				"business_line": "dojo",
				"student_ref":   receivable.StudentRef,
				"billing_month": receivable.BillingMonth.Format("2006-01-02"),
				"receivable_id": receivable.ID,
			},
			Items: []ledgerusecase.PostEntryItemRequest{
				{
					AccountID: strings.TrimSpace(req.CashAccountID),
					Side:      "debit",
					Amount:    req.PaymentAmount,
				},
				{
					AccountID: receivable.ReceivableAccountID,
					Side:      "credit",
					Amount:    req.PaymentAmount,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("post dojo payment: %w", err)
		}

		result.JournalEntryID = postResult.Entry.ID
		if err := s.repo.MarkSettled(txCtx, receivable.ID, req.PaymentAmount, postResult.Entry.ID, paidAt); err != nil {
			return fmt.Errorf("mark dojo receivable settled: %w", err)
		}
		if req.PaymentAmount.Cmp(remaining) == 0 {
			result.Status = "settled"
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func validateDojoAssessmentRequest(req AssessMonthlyTuitionRequest) error {
	switch {
	case strings.TrimSpace(req.CompanyCode) == "":
		return financedomain.ErrCompanyRequired
	case strings.TrimSpace(req.StudentRef) == "":
		return financedomain.ErrCustomerReferenceRequired
	case req.BillingMonth.IsZero():
		return financedomain.ErrBillingMonthRequired
	case req.DueDate.IsZero():
		return financedomain.ErrStartDateRequired
	case strings.TrimSpace(req.ReceivableAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.RevenueAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.CurrencyCode) == "":
		return financedomain.ErrCurrencyRequired
	case !req.Amount.IsPositive():
		return financedomain.ErrAmountMustBePositive
	default:
		return nil
	}
}

func validateDojoPaymentRequest(req CaptureDojoPaymentRequest) error {
	switch {
	case strings.TrimSpace(req.CompanyCode) == "":
		return financedomain.ErrCompanyRequired
	case strings.TrimSpace(req.StudentRef) == "":
		return financedomain.ErrCustomerReferenceRequired
	case req.BillingMonth.IsZero():
		return financedomain.ErrBillingMonthRequired
	case strings.TrimSpace(req.CashAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.CurrencyCode) == "":
		return financedomain.ErrCurrencyRequired
	case !req.PaymentAmount.IsPositive():
		return financedomain.ErrAmountMustBePositive
	default:
		return nil
	}
}
