package provincial

import (
	"encoding/json"
	"net/http"
	"strings"
	"vct-platform/backend/internal/domain/provincial"
	"vct-platform/backend/internal/shared/httputil"
)

// resolveProvinceID extracts province_id from query or context.
func (m *Module) resolveProvinceID(r *http.Request) string {
	return r.URL.Query().Get("province_id")
}

// ── Tournaments ────────────────────────────────────────────

func (m *Module) handleProvincialTournaments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	provID := m.resolveProvinceID(r)

	switch r.Method {
	case http.MethodGet:
		list, err := m.tournamentStore.ListTournaments(ctx, provID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"tournaments": list, "total": len(list)})

	case http.MethodPost:
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		var t provincial.ProvincialTournament
		if err := httputil.DecodeJSON(r, &t); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		t.ProvinceID = provID
		created, err := m.tournamentStore.CreateTournament(ctx, t)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleProvincialTournamentDetail(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/tournaments/")
	parts := strings.SplitN(path, "/", 2)
	tournamentID := parts[0]

	if len(parts) == 2 && r.Method == http.MethodPost {
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		action := parts[1]
		switch action {
		case "open":
			m.tournamentStore.UpdateTournament(r.Context(), tournamentID, map[string]interface{}{"status": "open"})
			httputil.Success(w, http.StatusOK, map[string]string{"message": "Giải đấu đã mở đăng ký"})
		case "start":
			m.tournamentStore.UpdateTournament(r.Context(), tournamentID, map[string]interface{}{"status": "in_progress"})
			httputil.Success(w, http.StatusOK, map[string]string{"message": "Giải đấu đã bắt đầu"})
		case "complete":
			m.tournamentStore.UpdateTournament(r.Context(), tournamentID, map[string]interface{}{"status": "completed"})
			httputil.Success(w, http.StatusOK, map[string]string{"message": "Giải đấu đã kết thúc"})
		case "cancel":
			m.tournamentStore.UpdateTournament(r.Context(), tournamentID, map[string]interface{}{"status": "cancelled"})
			httputil.Success(w, http.StatusOK, map[string]string{"message": "Giải đấu đã bị hủy"})
		case "registrations":
			m.handleTournamentRegistrations(w, r, tournamentID)
		case "results":
			m.handleTournamentResults(w, r, tournamentID)
		default:
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", "unknown action: "+action)
		}
		return
	}

	t, err := m.tournamentStore.GetTournament(r.Context(), tournamentID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "PROVINCIAL_404", "tournament not found")
		return
	}
	httputil.Success(w, http.StatusOK, t)
}

func (m *Module) handleTournamentRegistrations(w http.ResponseWriter, r *http.Request, tournamentID string) {
	if r.Method == http.MethodGet {
		list, err := m.tournamentStore.ListRegistrations(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"registrations": list})
	} else {
		var reg provincial.TournamentRegistration
		if err := httputil.DecodeJSON(r, &reg); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		reg.TournamentID = tournamentID
		created, _ := m.tournamentStore.CreateRegistration(r.Context(), reg)
		httputil.Success(w, http.StatusCreated, created)
	}
}

func (m *Module) handleTournamentResults(w http.ResponseWriter, r *http.Request, tournamentID string) {
	if r.Method == http.MethodGet {
		list, err := m.tournamentStore.ListResults(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"results": list})
	} else {
		var res provincial.TournamentResult
		if err := httputil.DecodeJSON(r, &res); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		res.TournamentID = tournamentID
		created, _ := m.tournamentStore.CreateResult(r.Context(), res)
		httputil.Success(w, http.StatusCreated, created)
	}
}

// ── Finance ────────────────────────────────────────────────

func (m *Module) handleProvincialFinance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	provID := m.resolveProvinceID(r)

	switch r.Method {
	case http.MethodGet:
		list, err := m.financeStore.List(ctx, provID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"entries": list, "total": len(list)})

	case http.MethodPost:
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		var e provincial.FinanceEntry
		if err := httputil.DecodeJSON(r, &e); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		e.ProvinceID = provID
		created, err := m.financeStore.Create(ctx, e)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleProvincialFinanceSummary(w http.ResponseWriter, r *http.Request) {
	provID := m.resolveProvinceID(r)
	sum, err := m.financeStore.Summary(r.Context(), provID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, sum)
}

// ── Certifications ─────────────────────────────────────────

func (m *Module) handleProvincialCertifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	provID := m.resolveProvinceID(r)

	switch r.Method {
	case http.MethodGet:
		list, err := m.certStore.List(ctx, provID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"certifications": list, "total": len(list)})

	case http.MethodPost:
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		var c provincial.ProvincialCert
		if err := httputil.DecodeJSON(r, &c); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		c.ProvinceID = provID
		created, err := m.certStore.Create(ctx, c)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Discipline ─────────────────────────────────────────────

func (m *Module) handleProvincialDiscipline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	provID := m.resolveProvinceID(r)

	switch r.Method {
	case http.MethodGet:
		list, err := m.disciplineStore.List(ctx, provID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"cases": list, "total": len(list)})

	case http.MethodPost:
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		var c provincial.DisciplineCase
		if err := httputil.DecodeJSON(r, &c); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		c.ProvinceID = provID
		created, err := m.disciplineStore.Create(ctx, c)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleProvincialDisciplineAction(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/discipline/")
	parts := strings.SplitN(path, "/", 2)
	caseID := parts[0]
	if r.Method != http.MethodPost || len(parts) < 2 {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", "invalid request")
		return
	}
	if _, ok := m.authenticate(w, r); !ok {
		return
	}
	action := parts[1]
	switch action {
	case "resolve":
		var body struct {
			Penalty string `json:"penalty"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		m.disciplineStore.Update(r.Context(), caseID, map[string]interface{}{"status": "resolved", "penalty": body.Penalty})
		httputil.Success(w, http.StatusOK, map[string]string{"status": "resolved"})
	case "close":
		m.disciplineStore.Update(r.Context(), caseID, map[string]interface{}{"status": "closed"})
		httputil.Success(w, http.StatusOK, map[string]string{"status": "closed"})
	default:
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", "unknown action: "+action)
	}
}

// ── Documents ──────────────────────────────────────────────

func (m *Module) handleProvincialDocuments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	provID := m.resolveProvinceID(r)

	switch r.Method {
	case http.MethodGet:
		list, err := m.docStore.List(ctx, provID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"documents": list, "total": len(list)})

	case http.MethodPost:
		if _, ok := m.authenticate(w, r); !ok {
			return
		}
		var d provincial.ProvincialDoc
		if err := httputil.DecodeJSON(r, &d); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", err.Error())
			return
		}
		d.ProvinceID = provID
		created, err := m.docStore.Create(ctx, d)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleProvincialDocumentAction(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/provincial/documents/")
	parts := strings.SplitN(path, "/", 2)
	docID := parts[0]
	if r.Method != http.MethodPost || len(parts) < 2 {
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", "invalid request")
		return
	}
	if _, ok := m.authenticate(w, r); !ok {
		return
	}
	action := parts[1]
	switch action {
	case "publish":
		m.docStore.Update(r.Context(), docID, map[string]interface{}{"status": "published"})
		httputil.Success(w, http.StatusOK, map[string]string{"status": "published"})
	case "archive":
		m.docStore.Update(r.Context(), docID, map[string]interface{}{"status": "archived"})
		httputil.Success(w, http.StatusOK, map[string]string{"status": "archived"})
	default:
		httputil.Error(w, http.StatusBadRequest, "PROVINCIAL_400", "unknown action: "+action)
	}
}
