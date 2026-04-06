package parent

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/parent"
	"vct-platform/backend/internal/shared/httputil"
)

// GET /api/v1/parent/dashboard
func (m *Module) handleParentDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	dash, err := m.service.GetDashboard(r.Context(), userID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, dash)
}

// GET /api/v1/parent/children
func (m *Module) handleParentChildren(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		links, err := m.service.ListAllLinks(r.Context(), userID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		if links == nil {
			links = []parent.ParentLink{}
		}
		httputil.Success(w, http.StatusOK, links)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// POST /api/v1/parent/children/link
func (m *Module) handleParentLinkChild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	var req struct {
		AthleteID   string `json:"athlete_id"`
		AthleteName string `json:"athlete_name"`
		Relation    string `json:"relation"`
		ParentName  string `json:"parent_name"` // Ported from Principal.DisplayName in legacy
	}
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "PARENT_400", "invalid payload")
		return
	}

	link := parent.ParentLink{
		ParentID:    userID,
		ParentName:  req.ParentName,
		AthleteID:   req.AthleteID,
		AthleteName: req.AthleteName,
		Relation:    req.Relation,
	}
	created, err := m.service.RequestLink(r.Context(), link)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "PARENT_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, created)
}

// handleParentChildDetail handles sub-resources for a child.
func (m *Module) handleParentChildDetail(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/parent/children/")
	parts := strings.SplitN(path, "/", 2)

	if len(parts) == 0 || parts[0] == "" {
		httputil.Error(w, http.StatusBadRequest, "PARENT_400", "missing link ID")
		return
	}

	id := parts[0]

	// DELETE /api/v1/parent/children/{linkID} — unlink child
	if r.Method == http.MethodDelete {
		// Verify ownership (or delegat to service if it does ownership checks)
		link, err := m.service.GetLinkByID(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "PARENT_404", "link not found")
			return
		}
		if link.ParentID != userID {
			httputil.Error(w, http.StatusForbidden, "PARENT_403", "forbidden")
			return
		}
		if err := m.service.DeleteLink(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
		return
	}

	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	if len(parts) < 2 {
		httputil.Error(w, http.StatusBadRequest, "PARENT_400", "missing sub-resource")
		return
	}

	athleteID := parts[0]
	sub := parts[1]

	// Ownership check
	if !m.service.IsChildOfParent(r.Context(), userID, athleteID) {
		httputil.Error(w, http.StatusForbidden, "PARENT_403", "not your child")
		return
	}

	switch sub {
	case "attendance":
		records, err := m.service.GetChildAttendance(r.Context(), athleteID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, records)
	case "results":
		results, err := m.service.GetChildResults(r.Context(), athleteID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, results)
	default:
		httputil.Error(w, http.StatusBadRequest, "PARENT_400", "invalid sub-resource")
	}
}

// handleParentConsents handles consent lists and creation.
func (m *Module) handleParentConsents(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		consents, err := m.service.ListConsents(r.Context(), userID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		if consents == nil {
			consents = []parent.ConsentRecord{}
		}
		httputil.Success(w, http.StatusOK, consents)

	case http.MethodPost:
		var req struct {
			AthleteID   string `json:"athlete_id"`
			AthleteName string `json:"athlete_name"`
			Type        string `json:"type"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := httputil.DecodeJSON(r, &req); err != nil {
			httputil.Error(w, http.StatusBadRequest, "PARENT_400", "invalid payload")
			return
		}

		// Verify the athlete is linked to this parent
		if !m.service.IsChildOfParent(r.Context(), userID, req.AthleteID) {
			httputil.Error(w, http.StatusForbidden, "PARENT_403", "forbidden")
			return
		}

		c := parent.ConsentRecord{
			ParentID:    userID,
			AthleteID:   req.AthleteID,
			AthleteName: req.AthleteName,
			Type:        parent.ConsentType(req.Type),
			Title:       req.Title,
			Description: req.Description,
		}
		created, err := m.service.CreateConsent(r.Context(), c)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "PARENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// handleParentConsentAction handles consent revocation.
func (m *Module) handleParentConsentAction(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	if r.Method != http.MethodDelete {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/parent/consents/")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "PARENT_400", "missing consent ID")
		return
	}
	if err := m.service.RevokeConsent(r.Context(), id, userID); err != nil {
		httputil.Error(w, http.StatusForbidden, "PARENT_403", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": "revoked", "id": id})
}
