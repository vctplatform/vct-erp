package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
	"vct-platform/backend/internal/shared/id"
	"vct-platform/backend/internal/shared/repository"
)

// RentalService automates deposit hold and release postings for the rental business.
type RentalService struct {
	txManager    financedomain.TxManager
	ledgerPoster LedgerPoster
	repo         financedomain.RentalDepositRepository
	now          func() time.Time
}

// NewRentalService constructs the rental deposit accounting service.
func NewRentalService(
	txManager financedomain.TxManager,
	ledgerPoster LedgerPoster,
	repo financedomain.RentalDepositRepository,
) *RentalService {
	return &RentalService{
		txManager:    txManager,
		ledgerPoster: ledgerPoster,
		repo:         repo,
		now:          time.Now,
	}
}

// CaptureDeposit locks a rental deposit into the intermediary holding account.
func (s *RentalService) CaptureDeposit(ctx context.Context, req CaptureRentalDepositRequest) (*CaptureRentalDepositResult, error) {
	if err := validateRentalCaptureRequest(req); err != nil {
		return nil, err
	}

	now := req.HeldAt.UTC()
	if now.IsZero() {
		now = s.now().UTC()
	}
	deposit := financedomain.RentalDeposit{
		ID:               id.NewUUID(),
		CompanyCode:      strings.TrimSpace(req.CompanyCode),
		RentalOrderID:    strings.TrimSpace(req.RentalOrderID),
		CustomerRef:      strings.TrimSpace(req.CustomerRef),
		CashAccountID:    strings.TrimSpace(req.CashAccountID),
		HoldingAccountID: strings.TrimSpace(req.HoldingAccountID),
		CurrencyCode:     normalizeCurrency(req.CurrencyCode),
		Amount:           req.Amount,
		Status:           "held",
		HeldAt:           now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	result := &CaptureRentalDepositResult{
		DepositID:     deposit.ID,
		DepositStatus: deposit.Status,
	}

	err := s.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationSerializable,
	}, func(txCtx context.Context) error {
		postResult, err := s.ledgerPoster.PostEntry(txCtx, ledgerusecase.PostEntryRequest{
			VoucherType:  "PT",
			CompanyCode:  deposit.CompanyCode,
			SourceModule: "rental",
			ExternalRef:  firstNonEmpty(strings.TrimSpace(req.SourceRef), deposit.RentalOrderID),
			Description:  fmt.Sprintf("Thu tien coc don thue %s", deposit.RentalOrderID),
			CurrencyCode: deposit.CurrencyCode,
			PostingDate:  now,
			Metadata: map[string]any{
				"business_line":   "rental",
				"cost_center":     "rental",
				"rental_order_id": deposit.RentalOrderID,
				"deposit_id":      deposit.ID,
				"customer_ref":    deposit.CustomerRef,
			},
			Items: []ledgerusecase.PostEntryItemRequest{
				{
					AccountID: deposit.CashAccountID,
					Side:      "debit",
					Amount:    deposit.Amount,
				},
				{
					AccountID: deposit.HoldingAccountID,
					Side:      "credit",
					Amount:    deposit.Amount,
				},
			},
		})
		if err != nil {
			return fmt.Errorf("post rental deposit hold: %w", err)
		}

		deposit.HeldEntryID = postResult.Entry.ID
		result.JournalEntryID = postResult.Entry.ID
		if err := s.repo.CreateDeposit(txCtx, deposit); err != nil {
			return fmt.Errorf("create rental deposit: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReleaseDeposit reverses the deposit hold once the equipment is returned.
func (s *RentalService) ReleaseDeposit(ctx context.Context, req ReleaseRentalDepositRequest) (*ReleaseRentalDepositResult, error) {
	if strings.TrimSpace(req.CompanyCode) == "" {
		return nil, financedomain.ErrCompanyRequired
	}
	if strings.TrimSpace(req.RentalOrderID) == "" {
		return nil, financedomain.ErrRentalOrderRequired
	}

	deposit, err := s.repo.GetByRentalOrder(ctx, strings.TrimSpace(req.CompanyCode), strings.TrimSpace(req.RentalOrderID))
	if err != nil {
		return nil, fmt.Errorf("get rental deposit: %w", err)
	}
	if deposit.ID == "" {
		return nil, financedomain.ErrDepositNotFound
	}
	if deposit.Status != "held" {
		return nil, financedomain.ErrDepositAlreadyReleased
	}

	releasedAt := req.ReleasedAt.UTC()
	if releasedAt.IsZero() {
		releasedAt = s.now().UTC()
	}
	damageAmount := req.DamageAmount
	if damageAmount.Sign() < 0 || damageAmount.Cmp(deposit.Amount) > 0 {
		return nil, financedomain.ErrDamageAmountInvalid
	}
	if damageAmount.IsPositive() && strings.TrimSpace(req.DamageRevenueAccountID) == "" {
		return nil, financedomain.ErrDamageAccountRequired
	}

	refundAmount := deposit.Amount.Sub(damageAmount)
	status := "released"
	if damageAmount.IsPositive() {
		status = "applied"
		if refundAmount.IsZero() {
			status = "forfeited"
		}
	}

	result := &ReleaseRentalDepositResult{
		DepositID:     deposit.ID,
		DepositStatus: status,
	}

	err = s.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationSerializable,
	}, func(txCtx context.Context) error {
		items := make([]ledgerusecase.PostEntryItemRequest, 0, 3)
		items = append(items, ledgerusecase.PostEntryItemRequest{
			AccountID: deposit.HoldingAccountID,
			Side:      "debit",
			Amount:    deposit.Amount,
		})
		if refundAmount.IsPositive() {
			cashAccountID := firstNonEmpty(strings.TrimSpace(req.CashAccountID), deposit.CashAccountID)
			items = append(items, ledgerusecase.PostEntryItemRequest{
				AccountID: cashAccountID,
				Side:      "credit",
				Amount:    refundAmount,
			})
		}
		if damageAmount.IsPositive() {
			items = append(items, ledgerusecase.PostEntryItemRequest{
				AccountID: strings.TrimSpace(req.DamageRevenueAccountID),
				Side:      "credit",
				Amount:    damageAmount,
			})
		}

		postResult, err := s.ledgerPoster.PostEntry(txCtx, ledgerusecase.PostEntryRequest{
			VoucherType:  "PC",
			CompanyCode:  deposit.CompanyCode,
			SourceModule: "rental",
			ExternalRef:  firstNonEmpty(strings.TrimSpace(req.SourceRef), deposit.RentalOrderID),
			Description:  s.releaseDescription(deposit.RentalOrderID, damageAmount),
			CurrencyCode: deposit.CurrencyCode,
			PostingDate:  releasedAt,
			Metadata: map[string]any{
				"business_line":   "rental",
				"cost_center":     "rental",
				"rental_order_id": deposit.RentalOrderID,
				"deposit_id":      deposit.ID,
				"customer_ref":    deposit.CustomerRef,
				"damage_amount":   damageAmount.String(),
				"refund_amount":   refundAmount.String(),
			},
			Items: items,
		})
		if err != nil {
			return fmt.Errorf("post rental deposit release: %w", err)
		}

		result.JournalEntryID = postResult.Entry.ID
		if err := s.repo.MarkDepositSettled(txCtx, deposit.ID, status, postResult.Entry.ID, releasedAt); err != nil {
			return fmt.Errorf("mark rental deposit settled: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *RentalService) releaseDescription(rentalOrderID string, damageAmount ledgerdomain.Money) string {
	if damageAmount.IsZero() {
		return fmt.Sprintf("Hoan tien coc don thue %s", rentalOrderID)
	}
	return fmt.Sprintf("Quyet toan tien coc don thue %s co khau tru hu hai", rentalOrderID)
}

func validateRentalCaptureRequest(req CaptureRentalDepositRequest) error {
	switch {
	case strings.TrimSpace(req.CompanyCode) == "":
		return financedomain.ErrCompanyRequired
	case strings.TrimSpace(req.RentalOrderID) == "":
		return financedomain.ErrRentalOrderRequired
	case strings.TrimSpace(req.CustomerRef) == "":
		return financedomain.ErrCustomerReferenceRequired
	case strings.TrimSpace(req.CashAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.HoldingAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.CurrencyCode) == "":
		return financedomain.ErrCurrencyRequired
	case !req.Amount.IsPositive():
		return financedomain.ErrAmountMustBePositive
	default:
		return nil
	}
}
