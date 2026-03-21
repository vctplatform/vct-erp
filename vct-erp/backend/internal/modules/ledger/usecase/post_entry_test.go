package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"vct-platform/backend/internal/modules/ledger/domain"
	"vct-platform/backend/internal/shared/repository"
	"vct-platform/backend/internal/shared/sqltx"
)

func TestPostEntrySuccess(t *testing.T) {
	now := time.Date(2026, 3, 21, 9, 30, 0, 0, time.UTC)
	accountRepo := &fakeAccountRepo{
		accounts: map[string]domain.Account{
			"111": {ID: "111", CompanyCode: "VCT_GROUP", Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, NormalSide: domain.SideDebit, IsActive: true, IsPostable: true},
			"511": {ID: "511", CompanyCode: "VCT_GROUP", Code: "5113", Name: "Revenue", Type: domain.AccountTypeRevenue, NormalSide: domain.SideCredit, IsActive: true, IsPostable: true},
		},
	}
	cache := &fakeCache{}
	journalRepo := &fakeJournalRepo{}
	voucherRepo := &fakeVoucherRepo{}
	balanceRepo := &fakeBalanceRepo{}
	outboxRepo := &fakeOutboxRepo{}
	publisher := &fakePublisher{}
	txManager := &fakeTxManager{}

	uc := NewPostEntryUseCase(txManager, accountRepo, journalRepo, voucherRepo, balanceRepo, outboxRepo, cache, 15*time.Minute, publisher, "ledger.events")
	uc.now = func() time.Time { return now }

	result, err := uc.PostEntry(context.Background(), PostEntryRequest{
		CompanyCode:  "VCT_GROUP",
		SourceModule: "dojo",
		Description:  "Thu hoc phi",
		CurrencyCode: "VND",
		PostingDate:  now,
		Items: []PostEntryItemRequest{
			{AccountID: "111", Side: "debit", Amount: domain.MustParseMoney("1500000.0000")},
			{AccountID: "511", Side: "credit", Amount: domain.MustParseMoney("1500000.0000")},
		},
	})
	if err != nil {
		t.Fatalf("PostEntry returned error: %v", err)
	}

	if !txManager.called {
		t.Fatal("expected transaction manager to be used")
	}

	if got, want := len(journalRepo.entries), 1; got != want {
		t.Fatalf("expected %d journal entry, got %d", want, got)
	}

	if got, want := len(journalRepo.items), 2; got != want {
		t.Fatalf("expected %d journal items, got %d", want, got)
	}

	if got, want := len(balanceRepo.deltas), 2; got != want {
		t.Fatalf("expected %d balance deltas, got %d", want, got)
	}

	if got, want := len(outboxRepo.enqueued), 1; got != want {
		t.Fatalf("expected %d outbox event, got %d", want, got)
	}

	if got, want := len(outboxRepo.published), 1; got != want {
		t.Fatalf("expected %d published outbox mark, got %d", want, got)
	}

	if result.OutboxDeferred {
		t.Fatal("expected immediate publish success")
	}

	if !result.Totals.Debit.Equal(domain.MustParseMoney("1500000.0000")) {
		t.Fatalf("unexpected debit total: %s", result.Totals.Debit.String())
	}

	if cache.setCalls != 1 {
		t.Fatalf("expected cache warming to happen once, got %d", cache.setCalls)
	}
}

func TestPostEntryRejectsUnbalancedEntry(t *testing.T) {
	uc := NewPostEntryUseCase(
		&fakeTxManager{},
		&fakeAccountRepo{},
		&fakeJournalRepo{},
		nil,
		&fakeBalanceRepo{},
		&fakeOutboxRepo{},
		nil,
		15*time.Minute,
		nil,
		"ledger.events",
	)

	_, err := uc.PostEntry(context.Background(), PostEntryRequest{
		CompanyCode:  "VCT_GROUP",
		SourceModule: "dojo",
		Description:  "Invalid post",
		CurrencyCode: "VND",
		Items: []PostEntryItemRequest{
			{AccountID: "111", Side: "debit", Amount: domain.MustParseMoney("100.0000")},
			{AccountID: "511", Side: "credit", Amount: domain.MustParseMoney("90.0000")},
		},
	})
	if !errors.Is(err, domain.ErrEntryNotBalanced) {
		t.Fatalf("expected ErrEntryNotBalanced, got %v", err)
	}
}

func TestPostEntryDefersOutboxWhenPublishFails(t *testing.T) {
	now := time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC)
	outboxRepo := &fakeOutboxRepo{}
	publisher := &fakePublisher{err: errors.New("redis unavailable")}
	uc := NewPostEntryUseCase(
		&fakeTxManager{},
		&fakeAccountRepo{
			accounts: map[string]domain.Account{
				"111": {ID: "111", CompanyCode: "VCT_GROUP", Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, NormalSide: domain.SideDebit, IsActive: true, IsPostable: true},
				"511": {ID: "511", CompanyCode: "VCT_GROUP", Code: "5113", Name: "Revenue", Type: domain.AccountTypeRevenue, NormalSide: domain.SideCredit, IsActive: true, IsPostable: true},
			},
		},
		&fakeJournalRepo{},
		&fakeVoucherRepo{},
		&fakeBalanceRepo{},
		outboxRepo,
		nil,
		15*time.Minute,
		publisher,
		"ledger.events",
	)
	uc.now = func() time.Time { return now }

	result, err := uc.PostEntry(context.Background(), PostEntryRequest{
		CompanyCode:  "VCT_GROUP",
		SourceModule: "saas",
		Description:  "Thu subscription",
		CurrencyCode: "USD",
		PostingDate:  now,
		Items: []PostEntryItemRequest{
			{AccountID: "111", Side: "debit", Amount: domain.MustParseMoney("99.9900")},
			{AccountID: "511", Side: "credit", Amount: domain.MustParseMoney("99.9900")},
		},
	})
	if err != nil {
		t.Fatalf("PostEntry returned error: %v", err)
	}

	if !result.OutboxDeferred {
		t.Fatal("expected outbox publish to be deferred")
	}

	if got := len(outboxRepo.published); got != 0 {
		t.Fatalf("expected published marker to stay empty, got %d", got)
	}
}

func TestPostEntryDefersOutboxInsideOuterTransaction(t *testing.T) {
	now := time.Date(2026, 3, 21, 10, 30, 0, 0, time.UTC)
	outboxRepo := &fakeOutboxRepo{}
	publisher := &fakePublisher{}
	uc := NewPostEntryUseCase(
		&fakeTxManager{},
		&fakeAccountRepo{
			accounts: map[string]domain.Account{
				"111": {ID: "111", CompanyCode: "VCT_GROUP", Code: "1111", Name: "Cash", Type: domain.AccountTypeAsset, NormalSide: domain.SideDebit, IsActive: true, IsPostable: true},
				"511": {ID: "511", CompanyCode: "VCT_GROUP", Code: "5113", Name: "Revenue", Type: domain.AccountTypeRevenue, NormalSide: domain.SideCredit, IsActive: true, IsPostable: true},
			},
		},
		&fakeJournalRepo{},
		&fakeVoucherRepo{},
		&fakeBalanceRepo{},
		outboxRepo,
		nil,
		15*time.Minute,
		publisher,
		"ledger.events",
	)
	uc.now = func() time.Time { return now }

	nestedCtx := sqltx.WithTx(context.Background(), &sql.Tx{})
	result, err := uc.PostEntry(nestedCtx, PostEntryRequest{
		CompanyCode:  "VCT_GROUP",
		SourceModule: "saas",
		Description:  "Nested transaction post",
		CurrencyCode: "VND",
		PostingDate:  now,
		Items: []PostEntryItemRequest{
			{AccountID: "111", Side: "debit", Amount: domain.MustParseMoney("1000.0000")},
			{AccountID: "511", Side: "credit", Amount: domain.MustParseMoney("1000.0000")},
		},
	})
	if err != nil {
		t.Fatalf("PostEntry returned error: %v", err)
	}

	if !result.OutboxDeferred {
		t.Fatal("expected nested transaction to defer outbox publication")
	}
	if got := len(publisher.events); got != 0 {
		t.Fatalf("expected publisher to stay idle, got %d events", got)
	}
	if got := len(outboxRepo.published); got != 0 {
		t.Fatalf("expected no published marker inside nested transaction, got %d", got)
	}
}

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTransaction(ctx context.Context, _ repository.TxOptions, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

type fakeAccountRepo struct {
	accounts map[string]domain.Account
}

func (f *fakeAccountRepo) GetByIDs(_ context.Context, ids []string) (map[string]domain.Account, error) {
	result := make(map[string]domain.Account, len(ids))
	for _, id := range ids {
		if account, ok := f.accounts[id]; ok {
			result[id] = account
		}
	}
	return result, nil
}

type fakeJournalRepo struct {
	entries []domain.JournalEntry
	items   []domain.JournalItem
}

func (f *fakeJournalRepo) CreateEntry(_ context.Context, entry *domain.JournalEntry) error {
	f.entries = append(f.entries, *entry)
	return nil
}

func (f *fakeJournalRepo) CreateItems(_ context.Context, _ string, items []domain.JournalItem, _ time.Time, _ string, _ string) error {
	f.items = append(f.items, items...)
	return nil
}

type fakeVoucherRepo struct{}

func (f *fakeVoucherRepo) NextVoucherNumber(_ context.Context, _ string, voucherType domain.VoucherType, postingDate time.Time) (string, error) {
	return string(voucherType) + "-0001/" + postingDate.Format("01-06"), nil
}

type fakeBalanceRepo struct {
	deltas []domain.AccountBalanceDelta
}

func (f *fakeBalanceRepo) ApplyDeltas(_ context.Context, deltas []domain.AccountBalanceDelta, _ string, _ time.Time) error {
	f.deltas = append(f.deltas, deltas...)
	return nil
}

type fakeOutboxRepo struct {
	enqueued  []domain.OutboxEvent
	published []string
}

func (f *fakeOutboxRepo) Enqueue(_ context.Context, event domain.OutboxEvent) error {
	f.enqueued = append(f.enqueued, event)
	return nil
}

func (f *fakeOutboxRepo) MarkPublished(_ context.Context, eventID string, _ time.Time) error {
	f.published = append(f.published, eventID)
	return nil
}

type fakeCache struct {
	store    map[string]domain.Account
	setCalls int
}

func (f *fakeCache) GetMany(_ context.Context, ids []string) (map[string]domain.Account, error) {
	result := make(map[string]domain.Account)
	for _, id := range ids {
		if account, ok := f.store[id]; ok {
			result[id] = account
		}
	}
	return result, nil
}

func (f *fakeCache) SetMany(_ context.Context, accounts []domain.Account, _ time.Duration) error {
	if f.store == nil {
		f.store = make(map[string]domain.Account, len(accounts))
	}
	for _, account := range accounts {
		f.store[account.ID] = account
	}
	f.setCalls++
	return nil
}

type fakePublisher struct {
	events []domain.OutboxEvent
	err    error
}

func (f *fakePublisher) Publish(_ context.Context, event domain.OutboxEvent) error {
	f.events = append(f.events, event)
	return f.err
}
