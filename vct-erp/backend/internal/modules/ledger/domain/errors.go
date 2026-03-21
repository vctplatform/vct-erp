package domain

import "errors"

var (
	ErrEntryHasNoItems      = errors.New("journal entry must contain at least two lines")
	ErrEntryNotBalanced     = errors.New("total debit must equal total credit")
	ErrAmountMustBePositive = errors.New("journal item amount must be greater than zero")
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountNotPostable   = errors.New("account is not postable")
	ErrAccountInactive      = errors.New("account is inactive")
	ErrInvalidAccountNature = errors.New("account normal side does not match account type")
	ErrUnsupportedSide      = errors.New("unsupported journal side")
	ErrCompanyRequired      = errors.New("company code is required")
	ErrCurrencyRequired     = errors.New("currency code is required")
	ErrSourceModuleRequired = errors.New("source module is required")
	ErrDescriptionRequired  = errors.New("description is required")
	ErrVoucherTypeRequired  = errors.New("voucher type is required")
	ErrUnsupportedVoucher   = errors.New("unsupported voucher type")
	ErrJournalEntryNotFound = errors.New("journal entry not found")
	ErrEntryAlreadyReversed = errors.New("journal entry is already reversed")
	ErrEntryNotPosted       = errors.New("journal entry is not posted")
)
