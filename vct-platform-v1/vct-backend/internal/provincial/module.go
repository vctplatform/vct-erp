package provincial

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/federation"
	"vct-platform/backend/internal/domain/provincial"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Provincial module.
type Module struct {
	service         *provincial.Service
	fedService      *federation.Service
	tournamentStore provincial.TournamentStore
	financeStore    provincial.FinanceStore
	certStore       provincial.CertStore
	disciplineStore provincial.DisciplineStore
	docStore        provincial.DocStore
	authFn          func(r *http.Request) (string, error)
	logger          *slog.Logger
}

// Deps holds the dependencies for the Provincial module.
type Deps struct {
	Service         *provincial.Service
	FedService      *federation.Service
	TournamentStore provincial.TournamentStore
	FinanceStore    provincial.FinanceStore
	CertStore       provincial.CertStore
	DisciplineStore provincial.DisciplineStore
	DocStore        provincial.DocStore
	AuthFn          func(r *http.Request) (string, error)
	Logger          *slog.Logger
}

// New creates a new Provincial module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		service:         deps.Service,
		fedService:      deps.FedService,
		tournamentStore: deps.TournamentStore,
		financeStore:    deps.FinanceStore,
		certStore:       deps.CertStore,
		disciplineStore: deps.DisciplineStore,
		docStore:        deps.DocStore,
		authFn:          deps.AuthFn,
		logger:          deps.Logger.With(slog.String("module", "provincial")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers provincial-level routes.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	// Main Operations (Existing)
	mux.HandleFunc("/api/v1/provincial/dashboard", m.handleProvincialDashboard)
	mux.HandleFunc("/api/v1/provincial/clubs", m.handleProvincialClubs)
	mux.HandleFunc("/api/v1/provincial/clubs/", m.handleProvincialClubs)
	mux.HandleFunc("/api/v1/provincial/athletes", m.handleProvincialAthletes)
	mux.HandleFunc("/api/v1/provincial/athletes/", m.handleProvincialAthletes)
	mux.HandleFunc("/api/v1/provincial/vo-sinh", m.handleProvincialVoSinh)
	mux.HandleFunc("/api/v1/provincial/vo-sinh/", m.handleProvincialVoSinh)
	mux.HandleFunc("/api/v1/provincial/coaches", m.handleProvincialCoaches)
	mux.HandleFunc("/api/v1/provincial/coaches/", m.handleProvincialCoaches)
	mux.HandleFunc("/api/v1/provincial/referees", m.handleProvincialReferees)
	mux.HandleFunc("/api/v1/provincial/referees/", m.handleProvincialReferees)
	mux.HandleFunc("/api/v1/provincial/committee", m.handleProvincialCommittee)
	mux.HandleFunc("/api/v1/provincial/committee/", m.handleProvincialCommittee)
	mux.HandleFunc("/api/v1/provincial/transfers", m.handleProvincialTransfers)
	mux.HandleFunc("/api/v1/provincial/transfers/", m.handleProvincialTransfers)

	// Phase 2 Operations (Ported from legacy)
	mux.HandleFunc("/api/v1/provincial/tournaments", m.handleProvincialTournaments)
	mux.HandleFunc("/api/v1/provincial/tournaments/", m.handleProvincialTournamentDetail)
	mux.HandleFunc("/api/v1/provincial/finance", m.handleProvincialFinance)
	mux.HandleFunc("/api/v1/provincial/finance/summary", m.handleProvincialFinanceSummary)
	mux.HandleFunc("/api/v1/provincial/certifications", m.handleProvincialCertifications)
	mux.HandleFunc("/api/v1/provincial/discipline", m.handleProvincialDiscipline)
	mux.HandleFunc("/api/v1/provincial/discipline/", m.handleProvincialDisciplineAction)
	mux.HandleFunc("/api/v1/provincial/documents", m.handleProvincialDocuments)
	mux.HandleFunc("/api/v1/provincial/documents/", m.handleProvincialDocumentAction)

	// Federation-level Provincial Reporting (Ported)
	mux.HandleFunc("/api/v1/provincial/reports", m.handleProvReportRoutes)
	mux.HandleFunc("/api/v1/provincial/reports/", m.handleProvReportRoutes)
	mux.HandleFunc("/api/v1/provincial/stats", m.handleProvStats)
	
	m.RegisterClubMgmtRoutes(mux)
	m.logger.Info("provincial module routes registered")
}

// authenticate extracts userID from request using the authFn.
func (m *Module) authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	if m.authFn == nil {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "authentication required")
		return "", false
	}
	userID, err := m.authFn(r)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", err.Error())
		return "", false
	}
	return userID, true
}
