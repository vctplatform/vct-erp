package federation

import (
	"log/slog"
	"net/http"

	"vct-platform/backend/internal/domain/approval"
	"vct-platform/backend/internal/domain/certification"
	"vct-platform/backend/internal/domain/discipline"
	"vct-platform/backend/internal/domain/document"
	"vct-platform/backend/internal/domain/federation"
	"vct-platform/backend/internal/domain/international"
	"vct-platform/backend/internal/shared/httputil"
)

// Module is the self-contained Federation module.
type Module struct {
	main          *federation.Service
	approval      *approval.Service
	certification *certification.Service
	discipline    *discipline.Service
	document      *document.Service
	international *international.Service
	broadcaster   httputil.EventBroadcaster
	authFn        func(r *http.Request) (string, error)
	logger        *slog.Logger
}

// Deps holds the dependencies for the Federation module.
type Deps struct {
	Main          *federation.Service
	Approval      *approval.Service
	Certification *certification.Service
	Discipline    *discipline.Service
	Document      *document.Service
	International *international.Service
	Broadcaster   httputil.EventBroadcaster
	AuthFn        func(r *http.Request) (string, error)
	Logger        *slog.Logger
}

// New creates a new Federation module.
func New(deps Deps) *Module {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	return &Module{
		main:          deps.Main,
		approval:      deps.Approval,
		certification: deps.Certification,
		discipline:    deps.Discipline,
		document:      deps.Document,
		international: deps.International,
		broadcaster:   deps.Broadcaster,
		authFn:        deps.AuthFn,
		logger:        deps.Logger.With(slog.String("module", "federation")),
	}
}

var _ httputil.Module = (*Module)(nil)

// RegisterRoutes registers all federation-related routes.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	// Infrastructure
	mux.HandleFunc("/api/v1/federation/provinces", m.handleProvinceRoutes)
	mux.HandleFunc("/api/v1/federation/provinces/", m.handleProvinceRoutes)
	mux.HandleFunc("/api/v1/federation/units", m.handleUnitRoutes)
	mux.HandleFunc("/api/v1/federation/units/", m.handleUnitRoutes)
	mux.HandleFunc("/api/v1/federation/org-chart", m.handleOrgChart)
	mux.HandleFunc("/api/v1/federation/statistics", m.handleStatistics)
	mux.HandleFunc("/api/v1/federation/stats", m.handleStatistics)
	mux.HandleFunc("/api/v1/federation/personnel", m.handlePersonnelRoutes)
	mux.HandleFunc("/api/v1/federation/personnel/", m.handlePersonnelRoutes)

	// Master Data
	mux.HandleFunc("/api/v1/federation/master/belts", m.handleMasterBelts)
	mux.HandleFunc("/api/v1/federation/master/belts/", m.handleMasterBelts)
	mux.HandleFunc("/api/v1/federation/master/weights", m.handleMasterWeights)
	mux.HandleFunc("/api/v1/federation/master/weights/", m.handleMasterWeights)

	// Approval Center (Consolidated from legacy)
	mux.HandleFunc("/api/v1/approvals/my-pending", m.handleApprovalMyPending)
	mux.HandleFunc("/api/v1/approvals/my-requests", m.handleApprovalMyRequests)
	mux.HandleFunc("/api/v1/approvals/workflows", m.handleWorkflowDefinitions)
	mux.HandleFunc("/api/v1/approvals/", m.handleApprovalCRUD)

	// Domain Modules
	mux.HandleFunc("/api/v1/documents", m.handleDocumentRoutes)
	mux.HandleFunc("/api/v1/documents/", m.handleDocumentRoutes)
	mux.HandleFunc("/api/v1/discipline/cases", m.handleDisciplineRoutes)
	mux.HandleFunc("/api/v1/discipline/cases/", m.handleDisciplineRoutes)
	mux.HandleFunc("/api/v1/certifications", m.handleCertificationRoutes)
	mux.HandleFunc("/api/v1/certifications/", m.handleCertificationRoutes)
	mux.HandleFunc("/api/v1/certifications/verify/", m.handleCertVerifyPublic)

	// Extended Features (PR, International, Workflow)
	mux.HandleFunc("/api/v1/federation/articles", m.handleArticleRoutes)
	mux.HandleFunc("/api/v1/federation/articles/", m.handleArticleRoutes)
	mux.HandleFunc("/api/v1/federation/partners", m.handlePartnerRoutes)
	mux.HandleFunc("/api/v1/federation/partners/", m.handlePartnerRoutes)
	mux.HandleFunc("/api/v1/federation/intl-events", m.handleIntlEventRoutes)
	mux.HandleFunc("/api/v1/federation/intl-events/", m.handleIntlEventRoutes)
	mux.HandleFunc("/api/v1/federation/workflows", m.handleWorkflowRoutes)
	mux.HandleFunc("/api/v1/federation/workflows/", m.handleWorkflowRoutes)

	// International Domain (Ported from legacy)
	mux.HandleFunc("/api/v1/international/partners", m.handlePartnerList)
	mux.HandleFunc("/api/v1/international/partners/", m.handlePartnerCRUD)
	mux.HandleFunc("/api/v1/international/events", m.handleIntlEventList)
	mux.HandleFunc("/api/v1/international/events/", m.handleIntlEventCRUD)
	mux.HandleFunc("/api/v1/international/delegations", m.handleDelegationList)
	mux.HandleFunc("/api/v1/international/delegations/", m.handleDelegationCRUD)

	// Portal Activity
	mux.HandleFunc("/api/v1/portal/activities", m.handlePortalActivities)

	m.logger.Info("federation module routes registered")
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
