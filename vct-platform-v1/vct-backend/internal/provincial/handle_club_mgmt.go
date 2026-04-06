package provincial

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/provincial"
	"vct-platform/backend/internal/shared/httputil"
)

// RegisterClubMgmtRoutes registers routes that were previously in club_handler.go
func (m *Module) RegisterClubMgmtRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/club/members", m.handleClubMembers)
	mux.HandleFunc("/api/v1/club/members/", m.handleClubMemberAction)
	mux.HandleFunc("/api/v1/club/classes", m.handleClubClasses)
	mux.HandleFunc("/api/v1/club/classes/", m.handleClubClassAction)
	mux.HandleFunc("/api/v1/club/finance", m.handleClubFinance)
	mux.HandleFunc("/api/v1/club/finance/summary", m.handleClubFinanceSummary)
}

func (m *Module) handleClubMembers(w http.ResponseWriter, r *http.Request) {
	if _, ok := httputil.GetPrincipal(r); !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Unauthorized")
		return
	}

	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}

	switch r.Method {
	case http.MethodGet:
		members, err := m.service.ListClubMembers(r.Context(), clubID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"members": members, "total": len(members)})

	case http.MethodPost:
		var member provincial.ClubMember
		if err := httputil.DecodeJSON(r, &member); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if member.ClubID == "" {
			member.ClubID = clubID
		}
		created, err := m.service.CreateClubMember(r.Context(), member)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleClubMemberAction(w http.ResponseWriter, r *http.Request) {
	if _, ok := httputil.GetPrincipal(r); !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Unauthorized")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/club/members/")
	parts := strings.Split(path, "/")
	id := parts[0]

	if len(parts) > 1 {
		action := parts[1]
		if r.Method != http.MethodPost {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		switch action {
		case "approve":
			if err := m.service.ApproveClubMember(r.Context(), id); err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, map[string]string{"message": "approved"})
		case "reject":
			if err := m.service.RejectClubMember(r.Context(), id); err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, map[string]string{"message": "rejected"})
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		member, err := m.service.GetClubMember(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "CLUB_404", "Member not found")
			return
		}
		httputil.Success(w, http.StatusOK, member)
	case http.MethodPatch, http.MethodPut:
		var patch map[string]any
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if err := m.service.UpdateClubMember(r.Context(), id, patch); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"message": "updated"})
	case http.MethodDelete:
		if err := m.service.DeleteClubMember(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusNoContent, nil)
	}
}

func (m *Module) handleClubClasses(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}

	switch r.Method {
	case http.MethodGet:
		classes, err := m.service.ListClubClasses(r.Context(), clubID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"classes": classes, "total": len(classes)})
	case http.MethodPost:
		var c provincial.ClubClass
		if err := httputil.DecodeJSON(r, &c); err != nil {
			httputil.Error(w, http.StatusBadRequest, "CLUB_400", err.Error())
			return
		}
		if c.ClubID == "" {
			c.ClubID = clubID
		}
		created, err := m.service.CreateClubClass(r.Context(), c)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusCreated, created)
	}
}

func (m *Module) handleClubClassAction(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/club/classes/")
	switch r.Method {
	case http.MethodGet:
		c, err := m.service.GetClubClass(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "CLUB_404", "Class not found")
			return
		}
		httputil.Success(w, http.StatusOK, c)
	case http.MethodDelete:
		if err := m.service.DeleteClubClass(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusNoContent, nil)
	}
}

func (m *Module) handleClubFinance(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}

	if r.Method == http.MethodGet {
		entries, err := m.service.ListClubFinance(r.Context(), clubID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, entries)
		return
	}
	httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
}

func (m *Module) handleClubFinanceSummary(w http.ResponseWriter, r *http.Request) {
	clubID := r.URL.Query().Get("club_id")
	if clubID == "" {
		clubID = "CLB-001"
	}

	summary, err := m.service.GetClubFinanceSummary(r.Context(), clubID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, summary)
}
