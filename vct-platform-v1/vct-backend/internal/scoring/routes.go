package scoring

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/shared/httputil"
)

// RegisterRoutes registers all scoring HTTP routes on the given mux.
// This implements the httputil.Module interface for self-registering routes.
//
// Routes:
//
//	POST /api/v1/scoring/combat/{matchID}/start
//	POST /api/v1/scoring/combat/{matchID}/score
//	POST /api/v1/scoring/combat/{matchID}/penalty
//	POST /api/v1/scoring/combat/{matchID}/end
//	GET  /api/v1/scoring/combat/{matchID}/state
//	POST /api/v1/scoring/forms/{perfID}/score
//	POST /api/v1/scoring/forms/{perfID}/finalize
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/scoring/", m.handleScoringRoutes)
	m.RegisterRegistrationRoutes(mux)
	m.logger.Info("scoring module routes registered")
}

// handleScoringRoutes is the main router for the scoring module.
func (m *Module) handleScoringRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/scoring/")
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	if len(segments) < 2 {
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", "invalid scoring path, expected /scoring/{type}/{id}/{action}")
		return
	}

	// Auth check
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}

	matchType := segments[0] // "combat" or "forms"
	matchID := segments[1]
	action := ""
	if len(segments) >= 3 {
		action = segments[2]
	}

	switch matchType {
	case "combat":
		m.handleCombatRouting(w, r, matchID, action, userID)
	case "forms":
		m.handleFormsRouting(w, r, matchID, action, userID)
	default:
		httputil.Error(w, http.StatusBadRequest, "SCORING_400", "match type phải là 'combat' hoặc 'forms'")
	}
}

// ── Combat Routing ───────────────────────────────────────────

func (m *Module) handleCombatRouting(w http.ResponseWriter, r *http.Request, matchID, action, userID string) {
	switch action {
	case "start":
		if r.Method != http.MethodPost {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		m.handleStartMatch(w, r, matchID, userID)

	case "score":
		if r.Method != http.MethodPost {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		m.handleSubmitScore(w, r, matchID, userID)

	case "penalty":
		if r.Method != http.MethodPost {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		m.handleRecordPenalty(w, r, matchID, userID)

	case "end":
		if r.Method != http.MethodPost {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		m.handleEndMatch(w, r, matchID, userID)

	case "state", "":
		if r.Method != http.MethodGet {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		m.handleGetState(w, r, matchID)

	default:
		httputil.Error(w, http.StatusNotFound, "SCORING_404", "Không tìm thấy tài nguyên")
	}
}

// ── Forms Routing ────────────────────────────────────────────

func (m *Module) handleFormsRouting(w http.ResponseWriter, r *http.Request, perfID, action, userID string) {
	switch action {
	case "score":
		if r.Method != http.MethodPost {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		m.handleSubmitFormsScore(w, r, perfID, userID)

	case "finalize":
		if r.Method != http.MethodPost {
			httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
			return
		}
		m.handleFinalizeForms(w, r, perfID, userID)

	default:
		httputil.Error(w, http.StatusNotFound, "SCORING_404", "Không tìm thấy tài nguyên")
	}
}
