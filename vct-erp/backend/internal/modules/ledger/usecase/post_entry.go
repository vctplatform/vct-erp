package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"vct-platform/backend/internal/modules/ledger/domain"
	"vct-platform/backend/internal/shared/id"
	"vct-platform/backend/internal/shared/repository"
)

// PostEntryItemRequest models a single debit or credit line in the incoming request.
type PostEntryItemRequest struct {
	AccountID   string       `json:"account_id"`
	Side        string       `json:"side"`
	Amount      domain.Money `json:"amount"`
	Description string       `json:"description,omitempty"`
}

// PostEntryRequest carries the input needed to post a double-entry journal.
type PostEntryRequest struct {
	ReferenceNo  string                 `json:"reference_no,omitempty"`
	CompanyCode  string                 `json:"company_code"`
	SourceModule string                 `json:"source_module"`
	ExternalRef  string                 `json:"external_ref,omitempty"`
	Description  string                 `json:"description"`
	CurrencyCode string                 `json:"currency_code"`
	PostingDate  time.Time              `json:"posting_date"`
	Metadata     map[string]any         `json:"metadata,omitempty"`
	Items        []PostEntryItemRequest `json:"items"`
}

// Totals summarizes the debit and credit totals that were posted.
type Totals struct {
	Debit  domain.Money `json:"debit"`
	Credit domain.Money `json:"credit"`
}

// PostEntryResult captures the committed journal entry and downstream publish state.
type PostEntryResult struct {
	Entry          domain.JournalEntry `json:"entry"`
	Totals         Totals              `json:"totals"`
	OutboxDeferred bool                `json:"outbox_deferred"`
}

// PostEntryUseCase orchestrates the ACID posting flow for the general ledger.
type PostEntryUseCase struct {
	txManager       domain.TransactionManager
	accountRepo     domain.AccountCatalogRepository
	journalRepo     domain.JournalEntryRepository
	balanceRepo     domain.AccountBalanceRepository
	outboxRepo      domain.OutboxRepository
	accountCache    domain.ChartOfAccountsCache
	accountCacheTTL time.Duration
	eventPublisher  domain.EventPublisher
	streamKey       string
	now             func() time.Time
}

// NewPostEntryUseCase builds the posting use case with explicit dependencies.
func NewPostEntryUseCase(
	txManager domain.TransactionManager,
	accountRepo domain.AccountCatalogRepository,
	journalRepo domain.JournalEntryRepository,
	balanceRepo domain.AccountBalanceRepository,
	outboxRepo domain.OutboxRepository,
	accountCache domain.ChartOfAccountsCache,
	accountCacheTTL time.Duration,
	eventPublisher domain.EventPublisher,
	streamKey string,
) *PostEntryUseCase {
	return &PostEntryUseCase{
		txManager:       txManager,
		accountRepo:     accountRepo,
		journalRepo:     journalRepo,
		balanceRepo:     balanceRepo,
		outboxRepo:      outboxRepo,
		accountCache:    accountCache,
		accountCacheTTL: accountCacheTTL,
		eventPublisher:  eventPublisher,
		streamKey:       streamKey,
		now:             time.Now,
	}
}

// PostEntry posts a balanced journal entry, updates balances, stores the outbox row, and publishes a downstream event after commit.
func (uc *PostEntryUseCase) PostEntry(ctx context.Context, req PostEntryRequest) (*PostEntryResult, error) {
	if err := uc.validateRequest(req); err != nil {
		return nil, err
	}

	preparedItems, totals, deltas, accountIDs, err := uc.prepareItems(req)
	if err != nil {
		return nil, err
	}

	accounts, err := uc.loadAccounts(ctx, accountIDs)
	if err != nil {
		return nil, fmt.Errorf("load accounts: %w", err)
	}

	for _, item := range preparedItems {
		account, ok := accounts[item.AccountID]
		if !ok {
			return nil, fmt.Errorf("%w: %s", domain.ErrAccountNotFound, item.AccountID)
		}
		if !account.IsActive {
			return nil, fmt.Errorf("%w: %s", domain.ErrAccountInactive, item.AccountID)
		}
		if !account.IsPostable {
			return nil, fmt.Errorf("%w: %s", domain.ErrAccountNotPostable, item.AccountID)
		}
	}

	now := uc.now().UTC()
	postingDate := req.PostingDate.UTC()
	if postingDate.IsZero() {
		postingDate = now
	}

	entry := domain.JournalEntry{
		ID:           id.NewUUID(),
		ReferenceNo:  uc.referenceNo(req.ReferenceNo, now),
		CompanyCode:  strings.TrimSpace(req.CompanyCode),
		SourceModule: strings.TrimSpace(req.SourceModule),
		ExternalRef:  strings.TrimSpace(req.ExternalRef),
		Description:  strings.TrimSpace(req.Description),
		CurrencyCode: strings.ToUpper(strings.TrimSpace(req.CurrencyCode)),
		PostingDate:  postingDate,
		Status:       domain.EntryStatusPosted,
		PostedAt:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
		Items:        preparedItems,
	}

	event, err := uc.buildOutboxEvent(entry, totals, req.Metadata, now)
	if err != nil {
		return nil, fmt.Errorf("build outbox event: %w", err)
	}

	if err := uc.txManager.WithinTransaction(ctx, repository.TxOptions{
		Isolation: repository.IsolationSerializable,
	}, func(txCtx context.Context) error {
		if err := uc.journalRepo.CreateEntry(txCtx, &entry); err != nil {
			return fmt.Errorf("create journal entry: %w", err)
		}

		if err := uc.journalRepo.CreateItems(txCtx, entry.ID, entry.Items, entry.CreatedAt, entry.CompanyCode, entry.CurrencyCode); err != nil {
			return fmt.Errorf("create journal items: %w", err)
		}

		if err := uc.balanceRepo.ApplyDeltas(txCtx, deltas, entry.ID, now); err != nil {
			return fmt.Errorf("apply account balances: %w", err)
		}

		if err := uc.outboxRepo.Enqueue(txCtx, event); err != nil {
			return fmt.Errorf("enqueue outbox event: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	result := &PostEntryResult{
		Entry:  entry,
		Totals: totals,
	}

	if uc.eventPublisher == nil {
		result.OutboxDeferred = true
		return result, nil
	}

	if err := uc.eventPublisher.Publish(ctx, event); err != nil {
		result.OutboxDeferred = true
		return result, nil
	}

	publishedAt := uc.now().UTC()
	if err := uc.outboxRepo.MarkPublished(ctx, event.ID, publishedAt); err != nil {
		result.OutboxDeferred = true
		return result, nil
	}

	return result, nil
}

func (uc *PostEntryUseCase) validateRequest(req PostEntryRequest) error {
	if strings.TrimSpace(req.CompanyCode) == "" {
		return domain.ErrCompanyRequired
	}
	if strings.TrimSpace(req.CurrencyCode) == "" {
		return domain.ErrCurrencyRequired
	}
	if strings.TrimSpace(req.SourceModule) == "" {
		return domain.ErrSourceModuleRequired
	}
	if strings.TrimSpace(req.Description) == "" {
		return domain.ErrDescriptionRequired
	}
	if len(req.Items) < 2 {
		return domain.ErrEntryHasNoItems
	}

	return nil
}

func (uc *PostEntryUseCase) prepareItems(req PostEntryRequest) ([]domain.JournalItem, Totals, []domain.AccountBalanceDelta, []string, error) {
	var (
		items      = make([]domain.JournalItem, 0, len(req.Items))
		totals     Totals
		accountIDs = make([]string, 0, len(req.Items))
		deltaMap   = make(map[string]domain.AccountBalanceDelta, len(req.Items))
	)

	for index, line := range req.Items {
		if !line.Amount.IsPositive() {
			return nil, Totals{}, nil, nil, fmt.Errorf("%w at line %d", domain.ErrAmountMustBePositive, index+1)
		}

		side, err := normalizeSide(line.Side)
		if err != nil {
			return nil, Totals{}, nil, nil, fmt.Errorf("%w at line %d", err, index+1)
		}

		accountID := strings.TrimSpace(line.AccountID)
		accountIDs = append(accountIDs, accountID)

		item := domain.JournalItem{
			LineNo:      index + 1,
			AccountID:   accountID,
			Side:        side,
			Amount:      line.Amount,
			Description: strings.TrimSpace(line.Description),
		}
		items = append(items, item)

		delta := deltaMap[accountID]
		delta.CompanyCode = strings.TrimSpace(req.CompanyCode)
		delta.AccountID = accountID
		delta.CurrencyCode = strings.ToUpper(strings.TrimSpace(req.CurrencyCode))
		if side == domain.SideDebit {
			totals.Debit = totals.Debit.Add(line.Amount)
			delta.DebitDelta = delta.DebitDelta.Add(line.Amount)
			delta.NetDelta = delta.NetDelta.Add(line.Amount)
		} else {
			totals.Credit = totals.Credit.Add(line.Amount)
			delta.CreditDelta = delta.CreditDelta.Add(line.Amount)
			delta.NetDelta = delta.NetDelta.Sub(line.Amount)
		}
		deltaMap[accountID] = delta
	}

	if !totals.Debit.Equal(totals.Credit) {
		return nil, Totals{}, nil, nil, domain.ErrEntryNotBalanced
	}

	deltas := make([]domain.AccountBalanceDelta, 0, len(deltaMap))
	for _, accountID := range uniqueStrings(accountIDs) {
		deltas = append(deltas, deltaMap[accountID])
	}

	return items, totals, deltas, uniqueStrings(accountIDs), nil
}

func (uc *PostEntryUseCase) loadAccounts(ctx context.Context, ids []string) (map[string]domain.Account, error) {
	if len(ids) == 0 {
		return map[string]domain.Account{}, nil
	}

	accounts := make(map[string]domain.Account, len(ids))
	missingIDs := ids

	if uc.accountCache != nil {
		cached, err := uc.accountCache.GetMany(ctx, ids)
		if err != nil {
			return nil, err
		}

		for accountID, account := range cached {
			accounts[accountID] = account
		}
		missingIDs = subtractStrings(ids, mapKeys(cached))
	}

	if len(missingIDs) == 0 {
		return accounts, nil
	}

	repoAccounts, err := uc.accountRepo.GetByIDs(ctx, missingIDs)
	if err != nil {
		return nil, err
	}

	for accountID, account := range repoAccounts {
		accounts[accountID] = account
	}

	if uc.accountCache != nil && len(repoAccounts) > 0 {
		cachePayload := make([]domain.Account, 0, len(repoAccounts))
		for _, accountID := range missingIDs {
			if account, ok := repoAccounts[accountID]; ok {
				cachePayload = append(cachePayload, account)
			}
		}
		if len(cachePayload) > 0 {
			if err := uc.accountCache.SetMany(ctx, cachePayload, uc.cacheTTLOrDefault()); err != nil {
				return nil, err
			}
		}
	}

	return accounts, nil
}

func (uc *PostEntryUseCase) buildOutboxEvent(entry domain.JournalEntry, totals Totals, metadata map[string]any, now time.Time) (domain.OutboxEvent, error) {
	payload := map[string]any{
		"entry_id":         entry.ID,
		"reference_no":     entry.ReferenceNo,
		"company_code":     entry.CompanyCode,
		"source_module":    entry.SourceModule,
		"external_ref":     entry.ExternalRef,
		"currency_code":    entry.CurrencyCode,
		"description":      entry.Description,
		"posting_date":     entry.PostingDate.Format(time.RFC3339),
		"posted_at":        entry.PostedAt.Format(time.RFC3339Nano),
		"total_debit":      totals.Debit.String(),
		"total_credit":     totals.Credit.String(),
		"line_count":       len(entry.Items),
		"journal_accounts": entryAccountIDs(entry.Items),
		"metadata":         metadata,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return domain.OutboxEvent{}, err
	}

	return domain.OutboxEvent{
		ID:            id.NewUUID(),
		AggregateType: "journal_entry",
		AggregateID:   entry.ID,
		EventType:     "ledger.entry.posted",
		StreamKey:     uc.streamKeyOrDefault(),
		Payload:       raw,
		Status:        domain.OutboxStatusPending,
		CreatedAt:     now,
	}, nil
}

func (uc *PostEntryUseCase) streamKeyOrDefault() string {
	if strings.TrimSpace(uc.streamKey) == "" {
		return "ledger.events"
	}
	return uc.streamKey
}

func (uc *PostEntryUseCase) referenceNo(referenceNo string, now time.Time) string {
	if trimmed := strings.TrimSpace(referenceNo); trimmed != "" {
		return trimmed
	}

	return fmt.Sprintf("JE-%s-%s", now.Format("20060102-150405"), id.NewUUID()[:8])
}

func normalizeSide(value string) (domain.Side, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(domain.SideDebit):
		return domain.SideDebit, nil
	case string(domain.SideCredit):
		return domain.SideCredit, nil
	default:
		return "", domain.ErrUnsupportedSide
	}
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		unique = append(unique, trimmed)
	}
	return unique
}

func subtractStrings(left []string, right []string) []string {
	rightSet := make(map[string]struct{}, len(right))
	for _, item := range right {
		rightSet[item] = struct{}{}
	}

	result := make([]string, 0, len(left))
	for _, item := range left {
		if _, ok := rightSet[item]; ok {
			continue
		}
		result = append(result, item)
	}

	return uniqueStrings(result)
}

func mapKeys(values map[string]domain.Account) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func entryAccountIDs(items []domain.JournalItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.AccountID)
	}
	return ids
}

func (uc *PostEntryUseCase) cacheTTLOrDefault() time.Duration {
	if uc.accountCacheTTL <= 0 {
		return 15 * time.Minute
	}
	return uc.accountCacheTTL
}

// IsValidationError helps the HTTP layer map business-validation errors to 4xx responses.
func IsValidationError(err error) bool {
	return errors.Is(err, domain.ErrEntryHasNoItems) ||
		errors.Is(err, domain.ErrEntryNotBalanced) ||
		errors.Is(err, domain.ErrAmountMustBePositive) ||
		errors.Is(err, domain.ErrUnsupportedSide) ||
		errors.Is(err, domain.ErrCompanyRequired) ||
		errors.Is(err, domain.ErrCurrencyRequired) ||
		errors.Is(err, domain.ErrSourceModuleRequired) ||
		errors.Is(err, domain.ErrDescriptionRequired)
}
