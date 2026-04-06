package tournament

import (
	"net/http"

	"vct-platform/backend/internal/shared/httputil"
)

// RegisterBracketRoutes registers all bracket-related and orchestrated tournament routes.
func (m *Module) RegisterBracketRoutes(mux *http.ServeMux) {
	// Bracket
	mux.HandleFunc("/api/v1/tournaments-action/generate-bracket", m.handleBracketGenerate)
	mux.HandleFunc("/api/v1/tournaments-action/brackets", m.handleBracketGet)

	// Tournament Orchestrated Actions
	mux.HandleFunc("/api/v1/tournaments-action/open-registration", m.handleTournamentOpenRegistration)
	mux.HandleFunc("/api/v1/tournaments-action/lock-registration", m.handleTournamentLockRegistration)
	mux.HandleFunc("/api/v1/tournaments-action/start", m.handleTournamentStart)
	mux.HandleFunc("/api/v1/tournaments-action/end", m.handleTournamentEnd)

	// Registration Validation
	mux.HandleFunc("/api/v1/registrations/validate", m.handleRegistrationValidate)

	// Team Actions
	mux.HandleFunc("/api/v1/teams-action/approve", m.handleTeamApprove)
	mux.HandleFunc("/api/v1/teams-action/reject", m.handleTeamReject)
	mux.HandleFunc("/api/v1/teams-action/checkin", m.handleTeamCheckin)

	// Results & Medals
	mux.HandleFunc("/api/v1/brackets/", m.handleAssignMedals)
}

// ── Bracket ──────────────────────────────────────────────────

func (m *Module) handleBracketGenerate(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "bracket_generate handler registered in tournament module",
	})
}

func (m *Module) handleBracketGet(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "bracket_get handler registered in tournament module",
	})
}

// ── Tournament Orchestrated Actions ──────────────────────────

func (m *Module) handleTournamentOpenRegistration(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "open_registration handler registered in tournament module",
	})
}

func (m *Module) handleTournamentLockRegistration(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "lock_registration handler registered in tournament module",
	})
}

func (m *Module) handleTournamentStart(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "tournament_start handler registered in tournament module",
	})
}

func (m *Module) handleTournamentEnd(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "tournament_end handler registered in tournament module",
	})
}

// ── Registration Validation ──────────────────────────────────

func (m *Module) handleRegistrationValidate(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "registration_validate handler registered in tournament module",
	})
}

// ── Team Approval ────────────────────────────────────────────

func (m *Module) handleTeamApprove(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "team_approve handler registered in tournament module",
	})
}

func (m *Module) handleTeamReject(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "team_reject handler registered in tournament module",
	})
}

func (m *Module) handleTeamCheckin(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "team_checkin handler registered in tournament module",
	})
}

// ── Results & Medals ─────────────────────────────────────────

func (m *Module) handleAssignMedals(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, http.StatusOK, map[string]string{
		"status": "assign_medals handler registered in tournament module",
	})
}
