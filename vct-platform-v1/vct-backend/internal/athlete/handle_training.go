package athlete

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/athlete"
	"vct-platform/backend/internal/shared/httputil"
)

// handleTrainingSessionRoutes registers training session related routes.
func (m *Module) handleTrainingSessionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/training-sessions/stats", m.handleTrainingSessionStats)
	mux.HandleFunc("/api/v1/training-sessions/", m.handleTrainingSessionByID)
	mux.HandleFunc("/api/v1/training-sessions", m.handleTrainingSessionList)
}

// ── Training Session Handlers ───────────────────────────────────

func (m *Module) handleTrainingSessionList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		athleteID := r.URL.Query().Get("athleteId")
		var list []athlete.TrainingSession
		var err error
		if athleteID != "" {
			list, err = m.training.ListByAthlete(r.Context(), athleteID)
		} else {
			list = []athlete.TrainingSession{}
		}
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		if list == nil {
			list = []athlete.TrainingSession{}
		}
		httputil.Success(w, http.StatusOK, list)

	case http.MethodPost:
		var payload athlete.TrainingSession
		if err := httputil.DecodeJSON(r, &payload); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TRAINING_400", err.Error())
			return
		}
		created, err := m.training.CreateSession(r.Context(), payload)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TRAINING_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleTrainingSessionByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/training-sessions/")
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "TRAINING_400", "missing session id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		sess, err := m.training.GetSession(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "TRAINING_404", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, sess)

	case http.MethodPatch:
		var patch map[string]interface{}
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TRAINING_400", err.Error())
			return
		}
		if err := m.training.UpdateSession(r.Context(), id, patch); err != nil {
			httputil.InternalError(w, err)
			return
		}
		sess, _ := m.training.GetSession(r.Context(), id)
		httputil.Success(w, http.StatusOK, sess)

	case http.MethodDelete:
		if err := m.training.DeleteSession(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusNoContent, nil)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleTrainingSessionStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	athleteID := r.URL.Query().Get("athleteId")
	if athleteID == "" {
		httputil.Error(w, http.StatusBadRequest, "TRAINING_400", "athleteId query parameter is required")
		return
	}
	stats, err := m.training.GetAttendanceStats(r.Context(), athleteID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}
