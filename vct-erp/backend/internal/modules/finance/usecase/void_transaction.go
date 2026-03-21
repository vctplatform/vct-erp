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

// VoidTransactionUseCase creates a reversal entry and preserves an audit snapshot in MongoDB.
type VoidTransactionUseCase struct {
	txManager    ledgerdomain.TransactionManager
	journalRepo  ledgerdomain.JournalEntryQueryRepository
	ledgerPoster LedgerPoster
	auditRepo    financedomain.VoidAuditRepository
	now          func() time.Time
}

// NewVoidTransactionUseCase wires the privileged VAS reversal flow.
func NewVoidTransactionUseCase(
	txManager ledgerdomain.TransactionManager,
	journalRepo ledgerdomain.JournalEntryQueryRepository,
	ledgerPoster LedgerPoster,
	auditRepo financedomain.VoidAuditRepository,
) *VoidTransactionUseCase {
	return &VoidTransactionUseCase{
		txManager:    txManager,
		journalRepo:  journalRepo,
		ledgerPoster: ledgerPoster,
		auditRepo:    auditRepo,
		now:          time.Now,
	}
}

// VoidEntry reverses a posted journal entry using an opposite-side entry, then stores before/after snapshots for inspections.
func (uc *VoidTransactionUseCase) VoidEntry(ctx context.Context, entryID string, reason string, actorID string) (*financedomain.VoidEntryResult, error) {
	if strings.TrimSpace(entryID) == "" {
		return nil, ledgerdomain.ErrJournalEntryNotFound
	}

	var (
		originalEntry ledgerdomain.JournalEntry
		afterEntry    ledgerdomain.JournalEntry
		reversal      *ledgerusecase.PostEntryResult
	)

	reason = strings.TrimSpace(reason)
	reversedAt := uc.now().UTC()
	if err := uc.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationSerializable,
	}, func(txCtx context.Context) error {
		entry, err := uc.journalRepo.GetEntryForUpdate(txCtx, entryID)
		if err != nil {
			return err
		}
		if entry.Status == ledgerdomain.EntryStatusReversed || entry.ReversalEntryID != "" {
			return ledgerdomain.ErrEntryAlreadyReversed
		}
		if entry.ReversalOfEntryID != "" {
			return ledgerdomain.ErrEntryAlreadyReversed
		}
		if entry.Status != ledgerdomain.EntryStatusPosted {
			return ledgerdomain.ErrEntryNotPosted
		}

		originalEntry = entry
		reversal, err = uc.ledgerPoster.PostEntry(txCtx, ledgerusecase.PostEntryRequest{
			VoucherType:  string(ledgerdomain.ReversalVoucherType(entry.VoucherType)),
			CompanyCode:  entry.CompanyCode,
			SourceModule: entry.SourceModule,
			ExternalRef:  firstNonEmpty(entry.ExternalRef, entry.ReferenceNo),
			Description:  uc.voidDescription(entry),
			CurrencyCode: entry.CurrencyCode,
			PostingDate:  reversedAt,
			ReversalOfID: entry.ID,
			VoidReason:   reason,
			Metadata:     uc.voidMetadata(entry, reason, actorID, reversedAt),
			Items:        uc.reverseItems(entry.Items),
		})
		if err != nil {
			return fmt.Errorf("post reversal entry: %w", err)
		}

		if err := uc.journalRepo.MarkReversed(txCtx, entry.ID, reversal.Entry.ID, reversedAt, reason); err != nil {
			return err
		}

		afterEntry = entry
		afterEntry.Status = ledgerdomain.EntryStatusReversed
		afterEntry.ReversalEntryID = reversal.Entry.ID
		afterEntry.ReversedAt = &reversedAt
		afterEntry.VoidReason = reason
		return nil
	}); err != nil {
		return nil, err
	}

	if uc.auditRepo != nil {
		if err := uc.auditRepo.RecordVoid(ctx, financedomain.VoidAuditLog{
			ID:              id.NewUUID(),
			CompanyCode:     originalEntry.CompanyCode,
			OriginalEntryID: originalEntry.ID,
			ReversalEntryID: reversal.Entry.ID,
			ActorID:         strings.TrimSpace(actorID),
			Reason:          reason,
			VoidedAt:        reversedAt,
			Before:          originalEntry,
			After: map[string]any{
				"original_entry": afterEntry,
				"reversal_entry": reversal.Entry,
			},
		}); err != nil {
			return nil, fmt.Errorf("record void audit log: %w", err)
		}
	}

	return &financedomain.VoidEntryResult{
		OriginalEntryID: originalEntry.ID,
		ReversalEntryID: reversal.Entry.ID,
		VoucherNo:       reversal.Entry.ReferenceNo,
		Status:          "reversed",
	}, nil
}

func (uc *VoidTransactionUseCase) reverseItems(items []ledgerdomain.JournalItem) []ledgerusecase.PostEntryItemRequest {
	reversed := make([]ledgerusecase.PostEntryItemRequest, 0, len(items))
	for _, item := range items {
		side := ledgerdomain.SideCredit
		if item.Side == ledgerdomain.SideCredit {
			side = ledgerdomain.SideDebit
		}

		reversed = append(reversed, ledgerusecase.PostEntryItemRequest{
			AccountID:   item.AccountID,
			Side:        string(side),
			Amount:      item.Amount,
			Description: item.Description,
		})
	}
	return reversed
}

func (uc *VoidTransactionUseCase) voidDescription(entry ledgerdomain.JournalEntry) string {
	if entry.ReferenceNo != "" {
		return fmt.Sprintf("Dao but toan cho chung tu %s", entry.ReferenceNo)
	}
	return fmt.Sprintf("Dao but toan cho giao dich %s", entry.ID)
}

func (uc *VoidTransactionUseCase) voidMetadata(entry ledgerdomain.JournalEntry, reason string, actorID string, reversedAt time.Time) map[string]any {
	metadata := make(map[string]any, len(entry.Metadata)+4)
	for key, value := range entry.Metadata {
		metadata[key] = value
	}

	metadata["void_of_entry_id"] = entry.ID
	metadata["void_at"] = reversedAt.Format(time.RFC3339Nano)
	if reason != "" {
		metadata["void_reason"] = reason
	}
	if trimmedActor := strings.TrimSpace(actorID); trimmedActor != "" {
		metadata["void_actor_id"] = trimmedActor
	}
	return metadata
}
