package provincial

import (
	"context"
)

// ── Tournament ───────────────────────────────────────────────────────────────

type TournamentStore interface {
	ListTournaments(ctx context.Context, provinceID string) ([]ProvincialTournament, error)
	GetTournament(ctx context.Context, id string) (*ProvincialTournament, error)
	CreateTournament(ctx context.Context, t ProvincialTournament) (*ProvincialTournament, error)
	UpdateTournament(ctx context.Context, id string, patch map[string]interface{}) error

	ListRegistrations(ctx context.Context, tournamentID string) ([]TournamentRegistration, error)
	CreateRegistration(ctx context.Context, r TournamentRegistration) (*TournamentRegistration, error)

	ListResults(ctx context.Context, tournamentID string) ([]TournamentResult, error)
	CreateResult(ctx context.Context, r TournamentResult) (*TournamentResult, error)
}

// ── Finance ──────────────────────────────────────────────────────────────────

type FinanceStore interface {
	List(ctx context.Context, provinceID string) ([]FinanceEntry, error)
	Create(ctx context.Context, e FinanceEntry) (*FinanceEntry, error)
	Summary(ctx context.Context, provinceID string) (*FinanceSummary, error)
}

// ── Cert ─────────────────────────────────────────────────────────────────────

type CertStore interface {
	List(ctx context.Context, provinceID string) ([]ProvincialCert, error)
	Create(ctx context.Context, c ProvincialCert) (*ProvincialCert, error)
}

// ── Discipline ───────────────────────────────────────────────────────────────

type DisciplineStore interface {
	List(ctx context.Context, provinceID string) ([]DisciplineCase, error)
	Create(ctx context.Context, c DisciplineCase) (*DisciplineCase, error)
	Update(ctx context.Context, id string, patch map[string]interface{}) error
}

// ── Doc ──────────────────────────────────────────────────────────────────────

type DocStore interface {
	List(ctx context.Context, provinceID string) ([]ProvincialDoc, error)
	Create(ctx context.Context, d ProvincialDoc) (*ProvincialDoc, error)
	Update(ctx context.Context, id string, patch map[string]interface{}) error
}
