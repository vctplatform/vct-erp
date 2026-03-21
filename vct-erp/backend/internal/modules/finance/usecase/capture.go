package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
)

const financeCaptureScope = "finance_capture"

// CaptureUseCase dispatches idempotent finance capture requests to the proper business adapter.
type CaptureUseCase struct {
	idempotencyRepo financedomain.IdempotencyRepository
	saasService     SaaSAccountingService
	dojoService     DojoAccountingService
	retailService   RetailAccountingService
	rentalService   RentalAccountingService
	now             func() time.Time
}

// NewCaptureUseCase constructs the finance capture entrypoint.
func NewCaptureUseCase(
	idempotencyRepo financedomain.IdempotencyRepository,
	saasService SaaSAccountingService,
	dojoService DojoAccountingService,
	retailService RetailAccountingService,
	rentalService RentalAccountingService,
) *CaptureUseCase {
	return &CaptureUseCase{
		idempotencyRepo: idempotencyRepo,
		saasService:     saasService,
		dojoService:     dojoService,
		retailService:   retailService,
		rentalService:   rentalService,
		now:             time.Now,
	}
}

// Capture transforms a business-specific payload into ledger-safe postings while enforcing idempotency.
func (uc *CaptureUseCase) Capture(ctx context.Context, req financedomain.CaptureRequest) (*financedomain.CaptureResult, error) {
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return nil, financedomain.ErrIdempotencyKeyRequired
	}
	if !operationBelongsToLine(req.BusinessLine, req.Operation) {
		if req.BusinessLine == "" {
			return nil, financedomain.ErrUnsupportedBusinessLine
		}
		return nil, financedomain.ErrUnsupportedOperation
	}

	requestHash, err := hashCaptureRequest(req)
	if err != nil {
		return nil, fmt.Errorf("hash finance capture request: %w", err)
	}

	now := uc.now().UTC()
	reservation, err := uc.idempotencyRepo.Reserve(ctx, financeCaptureScope, req.IdempotencyKey, requestHash, now)
	if err != nil {
		return nil, fmt.Errorf("reserve idempotency key: %w", err)
	}

	switch reservation.Status {
	case financedomain.IdempotencyStatusReplay:
		var replay financedomain.CaptureResult
		if err := json.Unmarshal(reservation.ResponsePayload, &replay); err != nil {
			return nil, fmt.Errorf("decode replay payload: %w", err)
		}
		replay.Replay = true
		return &replay, nil
	case financedomain.IdempotencyStatusConflict:
		return nil, financedomain.ErrIdempotencyConflict
	case financedomain.IdempotencyStatusInProgress:
		return nil, financedomain.ErrIdempotencyInProgress
	}

	result, err := uc.dispatch(ctx, req)
	if err != nil {
		_ = uc.idempotencyRepo.Fail(ctx, financeCaptureScope, req.IdempotencyKey, err.Error(), now)
		return nil, err
	}

	rawResult, err := json.Marshal(result)
	if err != nil {
		_ = uc.idempotencyRepo.Fail(ctx, financeCaptureScope, req.IdempotencyKey, err.Error(), now)
		return nil, fmt.Errorf("encode finance capture result: %w", err)
	}

	if err := uc.idempotencyRepo.Complete(ctx, financeCaptureScope, req.IdempotencyKey, rawResult, result.ResourceID, now); err != nil {
		return nil, fmt.Errorf("complete idempotency key: %w", err)
	}

	return result, nil
}

func (uc *CaptureUseCase) dispatch(ctx context.Context, req financedomain.CaptureRequest) (*financedomain.CaptureResult, error) {
	switch req.Operation {
	case financedomain.OperationSaaSCaptureAnnualContract:
		if uc.saasService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload CaptureAnnualContractRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.saasService.CaptureAnnualContract(ctx, payload)
		if err != nil {
			return nil, err
		}
		return capturePayload(req.BusinessLine, req.Operation, result.ContractID, result)
	case financedomain.OperationSaaSRecognizeDueRevenue:
		if uc.saasService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload RecognizeDueRevenueRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.saasService.RecognizeDueRevenue(ctx, payload)
		if err != nil {
			return nil, err
		}
		resourceID := ""
		if len(result.JournalEntryIDs) > 0 {
			resourceID = result.JournalEntryIDs[0]
		}
		return capturePayload(req.BusinessLine, req.Operation, resourceID, result)
	case financedomain.OperationDojoAssessMonthlyTuition:
		if uc.dojoService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload AssessMonthlyTuitionRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.dojoService.AssessMonthlyTuition(ctx, payload)
		if err != nil {
			return nil, err
		}
		return capturePayload(req.BusinessLine, req.Operation, result.ReceivableID, result)
	case financedomain.OperationDojoCapturePayment:
		if uc.dojoService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload CaptureDojoPaymentRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.dojoService.CapturePayment(ctx, payload)
		if err != nil {
			return nil, err
		}
		return capturePayload(req.BusinessLine, req.Operation, result.ReceivableID, result)
	case financedomain.OperationRetailCaptureSale:
		if uc.retailService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload CaptureRetailSaleRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.retailService.CaptureSale(ctx, payload)
		if err != nil {
			return nil, err
		}
		return capturePayload(req.BusinessLine, req.Operation, result.OrderNo, result)
	case financedomain.OperationRetailCaptureRefund:
		if uc.retailService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload CaptureRetailRefundRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.retailService.CaptureRefund(ctx, payload)
		if err != nil {
			return nil, err
		}
		return capturePayload(req.BusinessLine, req.Operation, result.OrderNo, result)
	case financedomain.OperationRentalCaptureDeposit:
		if uc.rentalService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload CaptureRentalDepositRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.rentalService.CaptureDeposit(ctx, payload)
		if err != nil {
			return nil, err
		}
		return capturePayload(req.BusinessLine, req.Operation, result.DepositID, result)
	case financedomain.OperationRentalReleaseDeposit:
		if uc.rentalService == nil {
			return nil, financedomain.ErrUnsupportedOperation
		}
		var payload ReleaseRentalDepositRequest
		if err := decodePayload(req.Payload, &payload); err != nil {
			return nil, err
		}
		result, err := uc.rentalService.ReleaseDeposit(ctx, payload)
		if err != nil {
			return nil, err
		}
		return capturePayload(req.BusinessLine, req.Operation, result.DepositID, result)
	default:
		return nil, financedomain.ErrUnsupportedOperation
	}
}

func hashCaptureRequest(req financedomain.CaptureRequest) (string, error) {
	canonical := struct {
		BusinessLine financedomain.BusinessLine     `json:"business_line"`
		Operation    financedomain.CaptureOperation `json:"operation"`
		Payload      json.RawMessage                `json:"payload"`
	}{
		BusinessLine: req.BusinessLine,
		Operation:    req.Operation,
		Payload:      req.Payload,
	}

	raw, err := json.Marshal(canonical)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

func decodePayload(raw json.RawMessage, target any) error {
	if len(raw) == 0 {
		return financedomain.ErrUnsupportedOperation
	}

	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("decode finance payload: %w", err)
	}
	return nil
}

func capturePayload(line financedomain.BusinessLine, operation financedomain.CaptureOperation, resourceID string, payload any) (*financedomain.CaptureResult, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &financedomain.CaptureResult{
		BusinessLine: line,
		Operation:    operation,
		ResourceID:   resourceID,
		Payload:      raw,
	}, nil
}

func operationBelongsToLine(line financedomain.BusinessLine, operation financedomain.CaptureOperation) bool {
	switch line {
	case financedomain.BusinessLineSaaS:
		return operation == financedomain.OperationSaaSCaptureAnnualContract || operation == financedomain.OperationSaaSRecognizeDueRevenue
	case financedomain.BusinessLineDojo:
		return operation == financedomain.OperationDojoAssessMonthlyTuition || operation == financedomain.OperationDojoCapturePayment
	case financedomain.BusinessLineRetail:
		return operation == financedomain.OperationRetailCaptureSale || operation == financedomain.OperationRetailCaptureRefund
	case financedomain.BusinessLineRental:
		return operation == financedomain.OperationRentalCaptureDeposit || operation == financedomain.OperationRentalReleaseDeposit
	default:
		return false
	}
}

// IsCaptureConflict reports whether the capture error should be surfaced as a client conflict.
func IsCaptureConflict(err error) bool {
	return errors.Is(err, financedomain.ErrIdempotencyConflict) ||
		errors.Is(err, financedomain.ErrIdempotencyInProgress)
}

// IsCaptureValidationError reports whether the capture error should map to 4xx validation errors.
func IsCaptureValidationError(err error) bool {
	return errors.Is(err, financedomain.ErrIdempotencyKeyRequired) ||
		errors.Is(err, financedomain.ErrUnsupportedBusinessLine) ||
		errors.Is(err, financedomain.ErrUnsupportedOperation) ||
		errors.Is(err, financedomain.ErrCompanyRequired) ||
		errors.Is(err, financedomain.ErrCurrencyRequired) ||
		errors.Is(err, financedomain.ErrAmountMustBePositive) ||
		errors.Is(err, financedomain.ErrContractNumberRequired) ||
		errors.Is(err, financedomain.ErrCustomerReferenceRequired) ||
		errors.Is(err, financedomain.ErrAccountReferenceRequired) ||
		errors.Is(err, financedomain.ErrTermMonthsRequired) ||
		errors.Is(err, financedomain.ErrStartDateRequired) ||
		errors.Is(err, financedomain.ErrBillingMonthRequired) ||
		errors.Is(err, financedomain.ErrRentalOrderRequired) ||
		errors.Is(err, financedomain.ErrAmountExceedsBalance)
}
