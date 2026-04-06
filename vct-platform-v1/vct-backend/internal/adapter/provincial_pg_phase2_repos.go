package adapter

import (
	"context"
	"time"

	"vct-platform/backend/internal/domain/provincial"
	"vct-platform/backend/internal/store"
)

// ── Tournament ───────────────────────────────────────────────────────────────

type pgTournamentStore struct {
	*StoreAdapter[provincial.ProvincialTournament]
	regStore *StoreAdapter[provincial.TournamentRegistration]
	resStore *StoreAdapter[provincial.TournamentResult]
}

func NewPgTournamentStore(ds store.DataStore) provincial.TournamentStore {
	return &pgTournamentStore{
		StoreAdapter: NewStoreAdapter[provincial.ProvincialTournament](ds, "provincial_tournaments"),
		regStore:     NewStoreAdapter[provincial.TournamentRegistration](ds, "provincial_tournament_regs"),
		resStore:     NewStoreAdapter[provincial.TournamentResult](ds, "provincial_tournament_results"),
	}
}

func (r *pgTournamentStore) ListTournaments(ctx context.Context, provinceID string) ([]provincial.ProvincialTournament, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialTournament
	for _, t := range items {
		if provinceID == "" || t.ProvinceID == provinceID {
			res = append(res, t)
		}
	}
	return res, nil
}

func (r *pgTournamentStore) GetTournament(ctx context.Context, id string) (*provincial.ProvincialTournament, error) {
	return r.StoreAdapter.GetByID(id)
}

func (r *pgTournamentStore) CreateTournament(ctx context.Context, t provincial.ProvincialTournament) (*provincial.ProvincialTournament, error) {
	return r.StoreAdapter.Create(t)
}

func (r *pgTournamentStore) UpdateTournament(ctx context.Context, id string, patch map[string]interface{}) error {
	patch["updated_at"] = time.Now().UTC()
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

func (r *pgTournamentStore) ListRegistrations(ctx context.Context, tournamentID string) ([]provincial.TournamentRegistration, error) {
	items, err := r.regStore.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.TournamentRegistration
	for _, reg := range items {
		if reg.TournamentID == tournamentID {
			res = append(res, reg)
		}
	}
	return res, nil
}

func (r *pgTournamentStore) CreateRegistration(ctx context.Context, reg provincial.TournamentRegistration) (*provincial.TournamentRegistration, error) {
	return r.regStore.Create(reg)
}

func (r *pgTournamentStore) ListResults(ctx context.Context, tournamentID string) ([]provincial.TournamentResult, error) {
	items, err := r.resStore.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.TournamentResult
	for _, reg := range items {
		if reg.TournamentID == tournamentID {
			res = append(res, reg)
		}
	}
	return res, nil
}

func (r *pgTournamentStore) CreateResult(ctx context.Context, res provincial.TournamentResult) (*provincial.TournamentResult, error) {
	return r.resStore.Create(res)
}

// ── Finance ──────────────────────────────────────────────────────────────────

type pgFinanceStore struct {
	*StoreAdapter[provincial.FinanceEntry]
}

func NewPgFinanceStore(ds store.DataStore) provincial.FinanceStore {
	return &pgFinanceStore{
		StoreAdapter: NewStoreAdapter[provincial.FinanceEntry](ds, "provincial_finances"),
	}
}

func (r *pgFinanceStore) List(ctx context.Context, provinceID string) ([]provincial.FinanceEntry, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.FinanceEntry
	for _, e := range items {
		if provinceID == "" || e.ProvinceID == provinceID {
			res = append(res, e)
		}
	}
	return res, nil
}

func (r *pgFinanceStore) Create(ctx context.Context, e provincial.FinanceEntry) (*provincial.FinanceEntry, error) {
	return r.StoreAdapter.Create(e)
}

func (r *pgFinanceStore) Summary(ctx context.Context, provinceID string) (*provincial.FinanceSummary, error) {
	entries, err := r.List(ctx, provinceID)
	if err != nil {
		return nil, err
	}
	sum := &provincial.FinanceSummary{ProvinceID: provinceID, EntryCount: len(entries)}
	for _, e := range entries {
		if e.Type == provincial.FinanceEntryIncome {
			sum.TotalIncome += e.Amount
		}
		if e.Type == provincial.FinanceEntryExpense {
			sum.TotalExpense += e.Amount
		}
	}
	sum.Balance = sum.TotalIncome - sum.TotalExpense
	return sum, nil
}

// ── Cert ─────────────────────────────────────────────────────────────────────

type pgCertStore struct {
	*StoreAdapter[provincial.ProvincialCert]
}

func NewPgCertStore(ds store.DataStore) provincial.CertStore {
	return &pgCertStore{
		StoreAdapter: NewStoreAdapter[provincial.ProvincialCert](ds, "provincial_certs"),
	}
}

func (r *pgCertStore) List(ctx context.Context, provinceID string) ([]provincial.ProvincialCert, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialCert
	for _, c := range items {
		if provinceID == "" || c.ProvinceID == provinceID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *pgCertStore) Create(ctx context.Context, c provincial.ProvincialCert) (*provincial.ProvincialCert, error) {
	return r.StoreAdapter.Create(c)
}

// ── Discipline ───────────────────────────────────────────────────────────────

type pgDisciplineStore struct {
	*StoreAdapter[provincial.DisciplineCase]
}

func NewPgDisciplineStore(ds store.DataStore) provincial.DisciplineStore {
	return &pgDisciplineStore{
		StoreAdapter: NewStoreAdapter[provincial.DisciplineCase](ds, "provincial_disciplines"),
	}
}

func (r *pgDisciplineStore) List(ctx context.Context, provinceID string) ([]provincial.DisciplineCase, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.DisciplineCase
	for _, c := range items {
		if provinceID == "" || c.ProvinceID == provinceID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *pgDisciplineStore) Create(ctx context.Context, c provincial.DisciplineCase) (*provincial.DisciplineCase, error) {
	return r.StoreAdapter.Create(c)
}

func (r *pgDisciplineStore) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}

// ── Doc ──────────────────────────────────────────────────────────────────────

type pgDocStore struct {
	*StoreAdapter[provincial.ProvincialDoc]
}

func NewPgDocStore(ds store.DataStore) provincial.DocStore {
	return &pgDocStore{
		StoreAdapter: NewStoreAdapter[provincial.ProvincialDoc](ds, "provincial_docs"),
	}
}

func (r *pgDocStore) List(ctx context.Context, provinceID string) ([]provincial.ProvincialDoc, error) {
	items, err := r.StoreAdapter.List()
	if err != nil {
		return nil, err
	}
	var res []provincial.ProvincialDoc
	for _, d := range items {
		if provinceID == "" || d.ProvinceID == provinceID {
			res = append(res, d)
		}
	}
	return res, nil
}

func (r *pgDocStore) Create(ctx context.Context, d provincial.ProvincialDoc) (*provincial.ProvincialDoc, error) {
	return r.StoreAdapter.Create(d)
}

func (r *pgDocStore) Update(ctx context.Context, id string, patch map[string]interface{}) error {
	_, err := r.StoreAdapter.Update(id, patch)
	return err
}
