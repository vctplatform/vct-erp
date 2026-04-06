package scoring

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/authz"
	"vct-platform/backend/internal/domain"
	"vct-platform/backend/internal/shared/httputil"
)

// List of allowed route registration
func (m *Module) RegisterRegistrationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/registration", m.handleRegistrationRoutes)
	mux.HandleFunc("/api/v1/registration/", m.handleRegistrationRoutes)
}

func (m *Module) handleRegistrationRoutes(w http.ResponseWriter, r *http.Request) {
	p, ok := httputil.GetPrincipal(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "AUTH_401", "Yêu cầu đăng nhập")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/registration")
	id := strings.Trim(path, "/")

	if id == "" {
		switch r.Method {
		case http.MethodGet:
			if !authz.CanEntityAction(p.User.Role, "registration", authz.ActionView) {
				httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem đăng ký")
				return
			}

			athleteID := r.URL.Query().Get("athleteId")
			tournamentID := r.URL.Query().Get("tournamentId")
			var list []domain.Registration
			var err error

			if athleteID != "" {
				list, err = m.registration.ListByAthlete(r.Context(), athleteID)
			} else if tournamentID != "" {
				list, err = m.registration.ListByTournament(r.Context(), tournamentID)
			} else {
				list, err = m.registration.ListRegistrations(r.Context())
			}

			if err != nil {
				httputil.InternalError(w, err)
				return
			}
			httputil.Success(w, http.StatusOK, list)
			return

		case http.MethodPost:
			if !authz.CanEntityAction(p.User.Role, "registration", authz.ActionCreate) {
				httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền tạo đăng ký")
				return
			}
			var payload domain.Registration
			if err := httputil.DecodeJSON(r, &payload); err != nil {
				httputil.Error(w, http.StatusBadRequest, "REG_400", err.Error())
				return
			}
			created, err := m.registration.CreateRegistration(r.Context(), payload)
			if err != nil {
				httputil.Error(w, http.StatusBadRequest, "REG_400", err.Error())
				return
			}
			m.broadcast("registration", "created", created.ID, httputil.ToMap(created), nil)
			httputil.Success(w, http.StatusCreated, created)
			return
		default:
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
	}

	// Detail routes
	switch r.Method {
	case http.MethodGet:
		if !authz.CanEntityAction(p.User.Role, "registration", authz.ActionView) {
			httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xem đăng ký")
			return
		}
		reg, err := m.registration.GetRegistration(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "REG_404", "Không tìm thấy đăng ký")
			return
		}
		httputil.Success(w, http.StatusOK, reg)

	case http.MethodPatch:
		if !authz.CanEntityAction(p.User.Role, "registration", authz.ActionUpdate) {
			httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền cập nhật đăng ký")
			return
		}
		var patch map[string]interface{}
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "REG_400", err.Error())
			return
		}
		updated, err := m.registration.UpdateRegistration(r.Context(), id, patch)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "REG_400", err.Error())
			return
		}
		m.broadcast("registration", "updated", id, httputil.ToMap(updated), nil)
		httputil.Success(w, http.StatusOK, updated)

	case http.MethodDelete:
		if !authz.CanEntityAction(p.User.Role, "registration", authz.ActionDelete) {
			httputil.Error(w, http.StatusForbidden, "AUTH_403", "Không có quyền xóa đăng ký")
			return
		}
		if err := m.registration.DeleteRegistration(r.Context(), id); err != nil {
			httputil.Error(w, http.StatusBadRequest, "REG_400", err.Error())
			return
		}
		m.broadcast("registration", "deleted", id, nil, nil)
		w.WriteHeader(http.StatusNoContent)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}
