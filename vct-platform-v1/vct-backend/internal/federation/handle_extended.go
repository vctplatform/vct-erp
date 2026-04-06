package federation

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/federation"
	"vct-platform/backend/internal/domain/international"
	"vct-platform/backend/internal/shared/httputil"
)

// ── PR Handlers ──────────────────────────────────────────────

func (m *Module) handleArticleRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/federation/articles")
	id := strings.Trim(path, "/")

	if id == "" {
		switch r.Method {
		case http.MethodGet:
			m.handleListArticles(w, r)
		case http.MethodPost:
			m.handleCreateArticle(w, r)
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		m.handleGetArticle(w, r, id)
	case http.MethodPut:
		m.handleUpdateArticle(w, r, id)
	case http.MethodDelete:
		m.handleDeleteArticle(w, r, id)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleListArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := m.main.ListArticles(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	status := r.URL.Query().Get("status")
	if status != "" {
		var filtered []federation.NewsArticle
		for _, a := range articles {
			if string(a.Status) == status {
				filtered = append(filtered, a)
			}
		}
		articles = filtered
	}
	httputil.Success(w, http.StatusOK, articles)
}

func (m *Module) handleGetArticle(w http.ResponseWriter, r *http.Request, id string) {
	a, err := m.main.GetArticle(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "Không tìm thấy bài viết")
		return
	}
	httputil.Success(w, http.StatusOK, a)
}

func (m *Module) handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	var article federation.NewsArticle
	if err := httputil.DecodeJSON(r, &article); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	article.AuthorID = userID
	if err := m.main.CreateArticle(r.Context(), article); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, map[string]string{"status": "article_created"})
}

func (m *Module) handleUpdateArticle(w http.ResponseWriter, r *http.Request, id string) {
	_, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	var article federation.NewsArticle
	if err := httputil.DecodeJSON(r, &article); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	article.ID = id
	if err := m.main.UpdateArticle(r.Context(), article); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": "article_updated"})
}

func (m *Module) handleDeleteArticle(w http.ResponseWriter, r *http.Request, id string) {
	_, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	if err := m.main.DeleteArticle(r.Context(), id); err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": "article_deleted"})
}

// ── International Domain Handlers (Partners, Events, Delegations) ─────────────────

func (m *Module) handlePartnerList(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		m.handlePartnerCreate(w, r)
		return
	}
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	country := r.URL.Query().Get("country")
	var partners []international.PartnerOrganization
	var err error
	if country != "" {
		partners, err = m.international.ListPartnersByCountry(r.Context(), country)
	} else {
		partners, err = m.international.ListPartners(r.Context())
	}
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"partners": partners, "total": len(partners)})
}

func (m *Module) handlePartnerCreate(w http.ResponseWriter, r *http.Request) {
	_, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	var partner international.PartnerOrganization
	if err := httputil.DecodeJSON(r, &partner); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	created, err := m.international.CreatePartner(r.Context(), partner)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handlePartnerCRUD(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/international/partners/")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", "partner ID required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		partner, err := m.international.GetPartner(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "partner not found")
			return
		}
		httputil.Success(w, http.StatusOK, partner)
	case http.MethodPut:
		_, ok := m.authenticate(w, r)
		if !ok {
			return
		}
		var partner international.PartnerOrganization
		if err := httputil.DecodeJSON(r, &partner); err != nil {
			httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
			return
		}
		updated, err := m.international.UpdatePartner(r.Context(), id, partner)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, updated)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleIntlEventList(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		m.handleIntlEventCreate(w, r)
		return
	}
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	events, err := m.international.ListEvents(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"events": events, "total": len(events)})
}

func (m *Module) handleIntlEventCreate(w http.ResponseWriter, r *http.Request) {
	_, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	var event international.InternationalEvent
	if err := httputil.DecodeJSON(r, &event); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	created, err := m.international.CreateEvent(r.Context(), event)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleIntlEventCRUD(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/international/events/")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", "event ID required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		event, err := m.international.GetEvent(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "event not found")
			return
		}
		httputil.Success(w, http.StatusOK, event)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleDelegationList(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		m.handleDelegationCreate(w, r)
		return
	}
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	eventID := r.URL.Query().Get("event_id")
	var delegations []international.Delegation
	var err error
	if eventID != "" {
		delegations, err = m.international.ListDelegationsByEvent(r.Context(), eventID)
	} else {
		delegations, err = m.international.ListDelegationsByEvent(r.Context(), "")
	}
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, map[string]any{"delegations": delegations, "total": len(delegations)})
}

func (m *Module) handleDelegationCreate(w http.ResponseWriter, r *http.Request) {
	_, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	var deleg international.Delegation
	if err := httputil.DecodeJSON(r, &deleg); err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	created, err := m.international.CreateDelegation(r.Context(), deleg)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

func (m *Module) handleDelegationCRUD(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/international/delegations/")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", "delegation ID required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		deleg, err := m.international.GetDelegation(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "FEDERATION_404", "delegation not found")
			return
		}
		httputil.Success(w, http.StatusOK, deleg)
	case http.MethodPut:
		_, ok := m.authenticate(w, r)
		if !ok {
			return
		}
		var deleg international.Delegation
		if err := httputil.DecodeJSON(r, &deleg); err != nil {
			httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
			return
		}
		updated, err := m.international.UpdateDelegation(r.Context(), id, deleg)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "FEDERATION_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, updated)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Legacy Compatibility (Old Federation Partners/Events) ──

func (m *Module) handlePartnerRoutes(w http.ResponseWriter, r *http.Request) {
	m.handlePartnerList(w, r)
}

func (m *Module) handleIntlEventRoutes(w http.ResponseWriter, r *http.Request) {
	m.handleIntlEventList(w, r)
}

// m.handleWorkflowRoutes ported from Federation Module.
func (m *Module) handleWorkflowRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/federation/workflows")
	id := strings.Trim(path, "/")
	if id == "" {
		m.handleListWorkflows(w, r)
		return
	}
	m.handleGetWorkflow(w, r, id)
}

func (m *Module) handleListWorkflows(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, []string{"Porting in progress"})
}

func (m *Module) handleGetWorkflow(w http.ResponseWriter, r *http.Request, id string) {
	httputil.Success(w, http.StatusOK, map[string]string{"id": id, "status": "Porting in progress"})
}
