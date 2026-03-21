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

// SaaSService automates deferred-revenue accounting for prepaid software contracts.
type SaaSService struct {
	txManager    financedomain.TxManager
	ledgerPoster LedgerPoster
	repo         financedomain.SaaSContractRepository
	now          func() time.Time
}

// NewSaaSService constructs the SaaS accounting service.
func NewSaaSService(
	txManager financedomain.TxManager,
	ledgerPoster LedgerPoster,
	repo financedomain.SaaSContractRepository,
) *SaaSService {
	return &SaaSService{
		txManager:    txManager,
		ledgerPoster: ledgerPoster,
		repo:         repo,
		now:          time.Now,
	}
}

// CaptureAnnualContract records the initial cash receipt and seeds the monthly revenue schedule.
func (s *SaaSService) CaptureAnnualContract(ctx context.Context, req CaptureAnnualContractRequest) (*CaptureAnnualContractResult, error) {
	if err := validateSaaSContractRequest(req); err != nil {
		return nil, err
	}

	startDate := req.ServiceStartDate.UTC()
	monthlyAmounts, err := req.TotalAmount.Split(req.TermMonths)
	if err != nil {
		return nil, fmt.Errorf("split contract amount: %w", err)
	}

	now := s.now().UTC()
	capturedAt := req.CapturedAt.UTC()
	if capturedAt.IsZero() {
		capturedAt = now
	}
	contractID := id.NewUUID()
	contract := financedomain.SaaSContract{
		ID:                         contractID,
		CompanyCode:                strings.TrimSpace(req.CompanyCode),
		ContractNo:                 strings.TrimSpace(req.ContractNo),
		CustomerRef:                strings.TrimSpace(req.CustomerRef),
		CashAccountID:              strings.TrimSpace(req.CashAccountID),
		DeferredRevenueAccountID:   strings.TrimSpace(req.DeferredRevenueAccountID),
		RecognizedRevenueAccountID: strings.TrimSpace(req.RecognizedRevenueAccountID),
		CurrencyCode:               normalizeCurrency(req.CurrencyCode),
		StartDate:                  startDate,
		EndDate:                    startDate.AddDate(0, req.TermMonths, -1),
		TermMonths:                 req.TermMonths,
		TotalAmount:                req.TotalAmount,
		SourceRef:                  strings.TrimSpace(req.SourceRef),
		CreatedAt:                  capturedAt,
		UpdatedAt:                  capturedAt,
	}

	schedules := make([]financedomain.RevenueSchedule, 0, req.TermMonths)
	monthCursor := monthStart(startDate)
	monthlyAmountStrings := make([]string, 0, req.TermMonths)
	for index, amount := range monthlyAmounts {
		serviceMonth := monthCursor.AddDate(0, index, 0)
		schedules = append(schedules, financedomain.RevenueSchedule{
			ID:           id.NewUUID(),
			ContractID:   contractID,
			SequenceNo:   index + 1,
			ServiceMonth: serviceMonth,
			Amount:       amount,
			Status:       "scheduled",
			CreatedAt:    capturedAt,
			UpdatedAt:    capturedAt,
		})
		monthlyAmountStrings = append(monthlyAmountStrings, amount.String())
	}

	result := &CaptureAnnualContractResult{
		ContractID:     contractID,
		ScheduleCount:  len(schedules),
		MonthlyAmounts: monthlyAmountStrings,
	}

	err = s.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationSerializable,
	}, func(txCtx context.Context) error {
		postResult, err := s.ledgerPoster.PostEntry(txCtx, ledgerusecase.PostEntryRequest{
			VoucherType:  "PT",
			CompanyCode:  contract.CompanyCode,
			SourceModule: "saas",
			ExternalRef:  firstNonEmpty(contract.SourceRef, contract.ContractNo),
			Description:  fmt.Sprintf("Nhan tien tra truoc hop dong SaaS %s", contract.ContractNo),
			CurrencyCode: contract.CurrencyCode,
			PostingDate:  capturedAt,
			Metadata: map[string]any{
				"business_line": "saas",
				"cost_center":   "saas",
				"contract_id":   contract.ID,
				"contract_no":   contract.ContractNo,
				"customer_ref":  contract.CustomerRef,
				"term_months":   contract.TermMonths,
			},
			Items: []ledgerusecase.PostEntryItemRequest{
				{
					AccountID: contract.CashAccountID,
					Side:      "debit",
					Amount:    contract.TotalAmount,
				},
				{
					AccountID: contract.DeferredRevenueAccountID,
					Side:      "credit",
					Amount:    contract.TotalAmount,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("post initial saas cash receipt: %w", err)
		}

		contract.InitialJournalEntryID = postResult.Entry.ID
		result.InitialJournalEntryID = postResult.Entry.ID

		if err := s.repo.CreateContract(txCtx, contract); err != nil {
			return fmt.Errorf("create saas contract: %w", err)
		}
		if err := s.repo.CreateSchedules(txCtx, schedules); err != nil {
			return fmt.Errorf("create saas revenue schedules: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RecognizeDueRevenue recognizes all scheduled SaaS revenue slices due up to the provided month.
func (s *SaaSService) RecognizeDueRevenue(ctx context.Context, req RecognizeDueRevenueRequest) (*RecognizeDueRevenueResult, error) {
	if strings.TrimSpace(req.CompanyCode) == "" {
		return nil, financedomain.ErrCompanyRequired
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 120
	}

	upTo := monthStart(req.UpTo.UTC())
	dueSchedules, err := s.repo.ListDueSchedules(ctx, strings.TrimSpace(req.CompanyCode), upTo, limit)
	if err != nil {
		return nil, fmt.Errorf("list due saas schedules: %w", err)
	}

	result := &RecognizeDueRevenueResult{
		JournalEntryIDs: make([]string, 0, len(dueSchedules)),
	}

	for _, schedule := range dueSchedules {
		recognizedAt := s.now().UTC()
		var journalEntryID string
		err := s.txManager.WithinTransaction(ctx, repository.TxOptions{
			Isolation: repository.IsolationSerializable,
		}, func(txCtx context.Context) error {
			postResult, err := s.ledgerPoster.PostEntry(txCtx, ledgerusecase.PostEntryRequest{
				VoucherType:  "PK",
				CompanyCode:  schedule.CompanyCode,
				SourceModule: "saas",
				ExternalRef:  firstNonEmpty(schedule.ContractNo, schedule.ContractID),
				Description:  fmt.Sprintf("Phan bo doanh thu SaaS hop dong %s ky %s", schedule.ContractNo, schedule.ServiceMonth.Format("2006-01")),
				CurrencyCode: schedule.CurrencyCode,
				PostingDate:  schedule.ServiceMonth,
				Metadata: map[string]any{
					"business_line": "saas",
					"cost_center":   "saas",
					"contract_id":   schedule.ContractID,
					"schedule_id":   schedule.ScheduleID,
					"service_month": schedule.ServiceMonth.Format("2006-01-02"),
				},
				Items: []ledgerusecase.PostEntryItemRequest{
					{
						AccountID: schedule.DeferredRevenueAccountID,
						Side:      "debit",
						Amount:    schedule.Amount,
					},
					{
						AccountID: schedule.RecognizedRevenueAccountID,
						Side:      "credit",
						Amount:    schedule.Amount,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("post saas recognition for schedule %s: %w", schedule.ScheduleID, err)
			}

			journalEntryID = postResult.Entry.ID
			if err := s.repo.MarkScheduleRecognized(txCtx, schedule.ScheduleID, journalEntryID, recognizedAt); err != nil {
				return fmt.Errorf("mark schedule %s recognized: %w", schedule.ScheduleID, err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		result.RecognizedCount++
		result.JournalEntryIDs = append(result.JournalEntryIDs, journalEntryID)
	}

	return result, nil
}

func validateSaaSContractRequest(req CaptureAnnualContractRequest) error {
	switch {
	case strings.TrimSpace(req.CompanyCode) == "":
		return financedomain.ErrCompanyRequired
	case strings.TrimSpace(req.ContractNo) == "":
		return financedomain.ErrContractNumberRequired
	case strings.TrimSpace(req.CustomerRef) == "":
		return financedomain.ErrCustomerReferenceRequired
	case strings.TrimSpace(req.CashAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.DeferredRevenueAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.RecognizedRevenueAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case req.TermMonths <= 0:
		return financedomain.ErrTermMonthsRequired
	case req.ServiceStartDate.IsZero():
		return financedomain.ErrStartDateRequired
	case strings.TrimSpace(req.CurrencyCode) == "":
		return financedomain.ErrCurrencyRequired
	case !req.TotalAmount.IsPositive():
		return financedomain.ErrAmountMustBePositive
	default:
		return nil
	}
}
