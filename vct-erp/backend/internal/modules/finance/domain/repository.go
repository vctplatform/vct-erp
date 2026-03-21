package domain

import (
	"context"
	"time"

	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
	"vct-platform/backend/internal/shared/repository"
)

// IdempotencyRepository persists and replays capture responses.
type IdempotencyRepository interface {
	Reserve(ctx context.Context, scope string, key string, requestHash string, now time.Time) (IdempotencyReservation, error)
	Complete(ctx context.Context, scope string, key string, responsePayload []byte, resourceID string, completedAt time.Time) error
	Fail(ctx context.Context, scope string, key string, lastError string, failedAt time.Time) error
}

// SaaSContractRepository stores SaaS contracts and revenue recognition schedules.
type SaaSContractRepository interface {
	CreateContract(ctx context.Context, contract SaaSContract) error
	CreateSchedules(ctx context.Context, schedules []RevenueSchedule) error
	ListDueSchedules(ctx context.Context, companyCode string, upTo time.Time, limit int) ([]DueRevenueSchedule, error)
	MarkScheduleRecognized(ctx context.Context, scheduleID string, journalEntryID string, recognizedAt time.Time) error
}

// DojoReceivableRepository stores dojo tuition receivables.
type DojoReceivableRepository interface {
	CreateReceivable(ctx context.Context, receivable DojoReceivable) error
	GetReceivable(ctx context.Context, companyCode string, studentRef string, billingMonth time.Time) (DojoReceivable, error)
	MarkSettled(ctx context.Context, receivableID string, amountPaid ledgerdomain.Money, journalEntryID string, settledAt time.Time) error
}

// RentalDepositRepository stores rental deposit lifecycles.
type RentalDepositRepository interface {
	CreateDeposit(ctx context.Context, deposit RentalDeposit) error
	GetByRentalOrder(ctx context.Context, companyCode string, rentalOrderID string) (RentalDeposit, error)
	MarkReleased(ctx context.Context, depositID string, releasedJournalEntryID string, releasedAt time.Time) error
}

// TxManager coordinates finance business writes with the ledger transaction.
type TxManager interface {
	repository.TxManager
}
