package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
)

// RetailService automates POS sales and refunds for the retail business line.
type RetailService struct {
	ledgerPoster LedgerPoster
}

// NewRetailService constructs the retail accounting service.
func NewRetailService(ledgerPoster LedgerPoster) *RetailService {
	return &RetailService{ledgerPoster: ledgerPoster}
}

// CaptureSale records a POS cash sale.
func (s *RetailService) CaptureSale(ctx context.Context, req CaptureRetailSaleRequest) (*CaptureRetailSaleResult, error) {
	if err := validateRetailSaleRequest(req); err != nil {
		return nil, err
	}

	postingDate := req.OccurredAt.UTC()
	if postingDate.IsZero() {
		postingDate = time.Now().UTC()
	}

	result, err := s.ledgerPoster.PostEntry(ctx, ledgerusecase.PostEntryRequest{
		VoucherType:  "PT",
		CompanyCode:  strings.TrimSpace(req.CompanyCode),
		SourceModule: "retail",
		ExternalRef:  firstNonEmpty(strings.TrimSpace(req.SourceRef), strings.TrimSpace(req.OrderNo)),
		Description:  fmt.Sprintf("Ban le POS don %s", req.OrderNo),
		CurrencyCode: normalizeCurrency(req.CurrencyCode),
		PostingDate:  postingDate,
		Metadata: map[string]any{
			"business_line": "retail",
			"cost_center":   "retail",
			"order_no":      strings.TrimSpace(req.OrderNo),
			"customer_ref":  strings.TrimSpace(req.CustomerRef),
			"flow_type":     "sale",
		},
		Items: []ledgerusecase.PostEntryItemRequest{
			{
				AccountID: strings.TrimSpace(req.CashAccountID),
				Side:      "debit",
				Amount:    req.Amount,
			},
			{
				AccountID: strings.TrimSpace(req.RevenueAccountID),
				Side:      "credit",
				Amount:    req.Amount,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("post retail sale: %w", err)
	}

	return &CaptureRetailSaleResult{
		OrderNo:        strings.TrimSpace(req.OrderNo),
		JournalEntryID: result.Entry.ID,
	}, nil
}

// CaptureRefund records a retail refund using contra revenue.
func (s *RetailService) CaptureRefund(ctx context.Context, req CaptureRetailRefundRequest) (*CaptureRetailRefundResult, error) {
	if err := validateRetailRefundRequest(req); err != nil {
		return nil, err
	}

	postingDate := req.RefundedAt.UTC()
	if postingDate.IsZero() {
		postingDate = time.Now().UTC()
	}

	result, err := s.ledgerPoster.PostEntry(ctx, ledgerusecase.PostEntryRequest{
		VoucherType:  "PC",
		CompanyCode:  strings.TrimSpace(req.CompanyCode),
		SourceModule: "retail",
		ExternalRef:  firstNonEmpty(strings.TrimSpace(req.SourceRef), strings.TrimSpace(req.OrderNo)),
		Description:  fmt.Sprintf("Hoan tien POS don %s", req.OrderNo),
		CurrencyCode: normalizeCurrency(req.CurrencyCode),
		PostingDate:  postingDate,
		Metadata: map[string]any{
			"business_line": "retail",
			"cost_center":   "retail",
			"order_no":      strings.TrimSpace(req.OrderNo),
			"customer_ref":  strings.TrimSpace(req.CustomerRef),
			"flow_type":     "refund",
		},
		Items: []ledgerusecase.PostEntryItemRequest{
			{
				AccountID: strings.TrimSpace(req.ContraRevenueAccountID),
				Side:      "debit",
				Amount:    req.Amount,
			},
			{
				AccountID: strings.TrimSpace(req.CashAccountID),
				Side:      "credit",
				Amount:    req.Amount,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("post retail refund: %w", err)
	}

	return &CaptureRetailRefundResult{
		OrderNo:        strings.TrimSpace(req.OrderNo),
		JournalEntryID: result.Entry.ID,
	}, nil
}

func validateRetailSaleRequest(req CaptureRetailSaleRequest) error {
	switch {
	case strings.TrimSpace(req.CompanyCode) == "":
		return financedomain.ErrCompanyRequired
	case strings.TrimSpace(req.OrderNo) == "":
		return financedomain.ErrContractNumberRequired
	case strings.TrimSpace(req.CashAccountID) == "":
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

func validateRetailRefundRequest(req CaptureRetailRefundRequest) error {
	switch {
	case strings.TrimSpace(req.CompanyCode) == "":
		return financedomain.ErrCompanyRequired
	case strings.TrimSpace(req.OrderNo) == "":
		return financedomain.ErrContractNumberRequired
	case strings.TrimSpace(req.CashAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.ContraRevenueAccountID) == "":
		return financedomain.ErrAccountReferenceRequired
	case strings.TrimSpace(req.CurrencyCode) == "":
		return financedomain.ErrCurrencyRequired
	case !req.Amount.IsPositive():
		return financedomain.ErrAmountMustBePositive
	default:
		return nil
	}
}
