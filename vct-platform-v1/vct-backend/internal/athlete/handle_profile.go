package athlete

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/athlete"
	"vct-platform/backend/internal/shared/httputil"
)

// handleAthleteProfileRoutes registers specific athlete profile, club membership,
// and tournament entry routes.
func (m *Module) handleAthleteProfileRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/athlete-profiles/me", m.handleAthleteProfileMe)
	mux.HandleFunc("/api/v1/athlete-profiles/stats", m.handleAthleteProfileStats)
	mux.HandleFunc("/api/v1/athlete-profiles/search", m.handleAthleteProfileSearch)
	mux.HandleFunc("/api/v1/athlete-profiles/", m.handleAthleteProfileByID)
	mux.HandleFunc("/api/v1/athlete-profiles", m.handleAthleteProfileList)
	mux.HandleFunc("/api/v1/club-memberships/", m.handleClubMembershipByID)
	mux.HandleFunc("/api/v1/club-memberships", m.handleClubMembershipList)
	mux.HandleFunc("/api/v1/tournament-entries/", m.handleTournamentEntryByID)
	mux.HandleFunc("/api/v1/tournament-entries", m.handleTournamentEntryList)
}

// ── Profile List/Create ──────────────────────────────────────

func (m *Module) handleAthleteProfileList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		clubID := r.URL.Query().Get("clubId")
		var list []athlete.AthleteProfile
		var err error
		if clubID != "" {
			list, err = m.profile.ListByClub(r.Context(), clubID)
		} else {
			list, err = m.profile.ListProfiles(r.Context())
		}
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, list)

	case http.MethodPost:
		var payload athlete.AthleteProfile
		if err := httputil.DecodeJSON(r, &payload); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		created, err := m.profile.CreateProfile(r.Context(), payload)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Profile Me (current user) ────────────────────────────────

func (m *Module) handleAthleteProfileMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	userID, ok := m.authenticate(w, r)
	if !ok {
		return
	}
	profile, err := m.profile.GetByUserID(r.Context(), userID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "ATHLETE_404", "athlete profile not found for current user")
		return
	}
	httputil.Success(w, http.StatusOK, profile)
}

// ── Profile By ID + Sub-resources ────────────────────────────

func (m *Module) handleAthleteProfileByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/athlete-profiles/")
	parts := strings.SplitN(path, "/", 2)
	id := parts[0]

	// Sub-resource routing: /athlete-profiles/:id/clubs or /athlete-profiles/:id/tournaments
	if len(parts) == 2 {
		subResource := parts[1]
		switch {
		case subResource == "clubs" || strings.HasPrefix(subResource, "clubs"):
			m.handleProfileClubs(w, r, id)
			return
		case subResource == "tournaments" || strings.HasPrefix(subResource, "tournaments"):
			m.handleProfileTournaments(w, r, id)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		profile, err := m.profile.GetProfile(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "ATHLETE_404", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, profile)

	case http.MethodPatch:
		var patch map[string]interface{}
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		if err := m.profile.UpdateProfile(r.Context(), id, patch); err != nil {
			httputil.InternalError(w, err)
			return
		}
		profile, _ := m.profile.GetProfile(r.Context(), id)
		httputil.Success(w, http.StatusOK, profile)

	case http.MethodDelete:
		if err := m.profile.DeleteProfile(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"deleted": id})

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Profile → Clubs Sub-resource ─────────────────────────────

func (m *Module) handleProfileClubs(w http.ResponseWriter, r *http.Request, athleteID string) {
	switch r.Method {
	case http.MethodGet:
		list, err := m.profile.ListMyClubs(r.Context(), athleteID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, list)

	case http.MethodPost:
		var payload athlete.ClubMembership
		if err := httputil.DecodeJSON(r, &payload); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		payload.AthleteID = athleteID
		created, err := m.profile.JoinClub(r.Context(), payload)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Profile → Tournaments Sub-resource ───────────────────────

func (m *Module) handleProfileTournaments(w http.ResponseWriter, r *http.Request, athleteID string) {
	switch r.Method {
	case http.MethodGet:
		list, err := m.profile.ListMyTournaments(r.Context(), athleteID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, list)

	case http.MethodPost:
		var payload athlete.TournamentEntry
		if err := httputil.DecodeJSON(r, &payload); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		payload.AthleteID = athleteID
		created, err := m.profile.EnterTournament(r.Context(), payload)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Club Membership CRUD ─────────────────────────────────────

func (m *Module) handleClubMembershipList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		athleteID := r.URL.Query().Get("athleteId")
		clubID := r.URL.Query().Get("clubId")
		var list []athlete.ClubMembership
		var err error
		if athleteID != "" {
			list, err = m.profile.ListMyClubs(r.Context(), athleteID)
		} else if clubID != "" {
			list, err = m.profile.ListClubMembers(r.Context(), clubID)
		} else {
			// Return all — admin view
			list, err = m.profile.ListMyClubs(r.Context(), "")
		}
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		if list == nil {
			list = []athlete.ClubMembership{}
		}
		httputil.Success(w, http.StatusOK, list)

	case http.MethodPost:
		var payload athlete.ClubMembership
		if err := httputil.DecodeJSON(r, &payload); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		created, err := m.profile.JoinClub(r.Context(), payload)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleClubMembershipByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/club-memberships/")

	switch r.Method {
	case http.MethodGet:
		membership, err := m.profile.ListMyClubs(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "ATHLETE_404", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, membership)

	case http.MethodPatch:
		var patch map[string]interface{}
		if err := httputil.DecodeJSON(r, &patch); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		if err := m.profile.UpdateMembership(r.Context(), id, patch); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"updated": id})

	case http.MethodDelete:
		if err := m.profile.LeaveClub(r.Context(), id); err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"deleted": id})

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Tournament Entry CRUD ────────────────────────────────────

func (m *Module) handleTournamentEntryList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		athleteID := r.URL.Query().Get("athleteId")
		tournamentID := r.URL.Query().Get("tournamentId")
		var list []athlete.TournamentEntry
		var err error
		if athleteID != "" {
			list, err = m.profile.ListMyTournaments(r.Context(), athleteID)
		} else if tournamentID != "" {
			list, err = m.profile.ListByTournament(r.Context(), tournamentID)
		} else {
			list, err = m.profile.ListMyTournaments(r.Context(), "")
		}
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		if list == nil {
			list = []athlete.TournamentEntry{}
		}
		httputil.Success(w, http.StatusOK, list)

	case http.MethodPost:
		var payload athlete.TournamentEntry
		if err := httputil.DecodeJSON(r, &payload); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		created, err := m.profile.EnterTournament(r.Context(), payload)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

func (m *Module) handleTournamentEntryByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/tournament-entries/")

	switch r.Method {
	case http.MethodGet:
		entry, err := m.profile.GetEntry(r.Context(), id)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "ATHLETE_404", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, entry)

	case http.MethodPatch:
		var body struct {
			Status string `json:"status"`
			Notes  string `json:"notes,omitempty"`
		}
		if err := httputil.DecodeJSON(r, &body); err != nil {
			httputil.Error(w, http.StatusBadRequest, "ATHLETE_400", err.Error())
			return
		}
		if body.Status != "" {
			switch body.Status {
			case "approve":
				if err := m.profile.ApproveEntry(r.Context(), id); err != nil {
					httputil.InternalError(w, err)
					return
				}
			case "reject":
				if err := m.profile.RejectEntry(r.Context(), id); err != nil {
					httputil.InternalError(w, err)
					return
				}
			default:
				if err := m.profile.UpdateEntryStatus(r.Context(), id, athlete.EntryStatus(body.Status)); err != nil {
					httputil.InternalError(w, err)
					return
				}
			}
		}
		entry, _ := m.profile.GetEntry(r.Context(), id)
		httputil.Success(w, http.StatusOK, entry)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Stats Endpoint ───────────────────────────────────────────

func (m *Module) handleAthleteProfileStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	stats, err := m.profile.GetStats(r.Context())
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}

// ── Search Endpoint ──────────────────────────────────────────

func (m *Module) handleAthleteProfileSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	query := r.URL.Query().Get("q")
	list, err := m.profile.SearchProfiles(r.Context(), query)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	if list == nil {
		list = []athlete.AthleteProfile{}
	}
	httputil.Success(w, http.StatusOK, list)
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
