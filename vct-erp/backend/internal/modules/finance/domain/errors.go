package domain

import "errors"

var (
	ErrCompanyRequired           = errors.New("company code is required")
	ErrCurrencyRequired          = errors.New("currency code is required")
	ErrAmountMustBePositive      = errors.New("amount must be greater than zero")
	ErrIdempotencyKeyRequired    = errors.New("idempotency key is required")
	ErrUnsupportedBusinessLine   = errors.New("unsupported business line")
	ErrUnsupportedOperation      = errors.New("unsupported finance operation")
	ErrIdempotencyConflict       = errors.New("idempotency key conflicts with a different payload")
	ErrIdempotencyInProgress     = errors.New("idempotency key is already being processed")
	ErrContractNumberRequired    = errors.New("contract number is required")
	ErrCustomerReferenceRequired = errors.New("customer reference is required")
	ErrAccountReferenceRequired  = errors.New("account reference is required")
	ErrTermMonthsRequired        = errors.New("term months must be greater than zero")
	ErrStartDateRequired         = errors.New("start date is required")
	ErrBillingMonthRequired      = errors.New("billing month is required")
	ErrRentalOrderRequired       = errors.New("rental order id is required")
	ErrDepositAlreadyReleased    = errors.New("rental deposit has already been released")
	ErrDepositNotFound           = errors.New("rental deposit not found")
	ErrReceivableNotFound        = errors.New("dojo receivable not found")
	ErrAmountExceedsBalance      = errors.New("amount exceeds remaining balance")
)
