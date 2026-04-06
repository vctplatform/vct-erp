package tournament

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/tournament"
	"vct-platform/backend/internal/shared/httputil"
)

func (m *Module) handleTournamentMgmt(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/tournament-mgmt/")
	parts := strings.Split(strings.TrimRight(path, "/"), "/")

	if len(parts) < 1 || parts[0] == "" {
		httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", "tournament ID required")
		return
	}

	tournamentID := parts[0]
	resource := ""
	subID := ""
	action := ""

	if len(parts) > 1 {
		resource = parts[1]
	}
	if len(parts) > 2 {
		subID = parts[2]
	}
	if len(parts) > 3 {
		action = parts[3]
	}

	switch resource {
	case "categories":
		m.handleTournamentCategories(w, r, tournamentID, subID)
	case "registrations":
		m.handleTournamentRegistrations(w, r, tournamentID, subID, action)
	case "schedule":
		m.handleTournamentSchedule(w, r, tournamentID, subID)
	case "arenas":
		m.handleTournamentArenas(w, r, tournamentID, subID)
	case "results":
		m.handleTournamentResults(w, r, tournamentID, subID, action)
	case "standings":
		m.handleTournamentStandings(w, r, tournamentID)
	case "stats":
		m.handleTournamentStats(w, r, tournamentID)
	case "export":
		m.handleTournamentExport(w, r, tournamentID, subID)
	case "batch":
		m.handleTournamentBatch(w, r, tournamentID, subID)
	default:
		httputil.Error(w, http.StatusNotFound, "TOURNAMENT_404", "unknown resource: "+resource)
	}
}

// ── Categories ──────────────────────────────────────────────

func (m *Module) handleTournamentCategories(w http.ResponseWriter, r *http.Request, tournamentID, catID string) {
	switch {
	case r.Method == "GET" && catID == "":
		cats, err := m.mgmt.ListCategories(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"categories": cats, "total": len(cats)})

	case r.Method == "POST" && catID == "":
		var cat tournament.Category
		if err := httputil.DecodeJSON(r, &cat); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		cat.TournamentID = tournamentID
		created, err := m.mgmt.CreateCategory(r.Context(), &cat)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	case r.Method == "GET" && catID != "":
		cat, err := m.mgmt.GetCategory(r.Context(), catID)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "TOURNAMENT_404", "category not found")
			return
		}
		httputil.Success(w, http.StatusOK, cat)

	case r.Method == "PUT" && catID != "":
		var cat tournament.Category
		if err := httputil.DecodeJSON(r, &cat); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		cat.ID = catID
		updated, err := m.mgmt.UpdateCategory(r.Context(), &cat)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, updated)

	case r.Method == "DELETE" && catID != "":
		if err := m.mgmt.DeleteCategory(r.Context(), catID); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Registrations ───────────────────────────────────────────

func (m *Module) handleTournamentRegistrations(w http.ResponseWriter, r *http.Request, tournamentID, regID, action string) {
	switch {
	case r.Method == "GET" && regID == "":
		regs, err := m.mgmt.ListRegistrations(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"registrations": regs, "total": len(regs)})

	case r.Method == "POST" && regID == "" && action == "":
		var reg tournament.Registration
		if err := httputil.DecodeJSON(r, &reg); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		reg.TournamentID = tournamentID
		created, err := m.mgmt.RegisterTeam(r.Context(), &reg)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	case r.Method == "GET" && regID != "" && action == "":
		reg, err := m.mgmt.GetRegistration(r.Context(), regID)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "TOURNAMENT_404", "registration not found")
			return
		}
		httputil.Success(w, http.StatusOK, reg)

	case r.Method == "GET" && regID != "" && action == "athletes":
		athletes, err := m.mgmt.ListRegistrationAthletes(r.Context(), regID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"athletes": athletes, "total": len(athletes)})

	case r.Method == "POST" && regID != "" && action == "athletes":
		var athlete tournament.RegistrationAthlete
		if err := httputil.DecodeJSON(r, &athlete); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		athlete.RegistrationID = regID
		created, err := m.mgmt.AddAthleteToRegistration(r.Context(), &athlete)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	case r.Method == "POST" && regID != "" && action == "submit":
		if err := m.mgmt.SubmitRegistration(r.Context(), regID); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"status": "submitted"})

	case r.Method == "POST" && regID != "" && action == "approve":
		p, _ := httputil.GetPrincipal(r)
		if err := m.mgmt.ApproveRegistration(r.Context(), regID, p.User.ID); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"status": "approved"})

	case r.Method == "POST" && regID != "" && action == "reject":
		p, _ := httputil.GetPrincipal(r)
		var body struct {
			Reason string `json:"reason"`
		}
		_ = httputil.DecodeJSON(r, &body)
		if err := m.mgmt.RejectRegistration(r.Context(), regID, p.User.ID, body.Reason); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"status": "rejected"})

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Schedule ────────────────────────────────────────────────

func (m *Module) handleTournamentSchedule(w http.ResponseWriter, r *http.Request, tournamentID, slotID string) {
	switch {
	case r.Method == "GET" && slotID == "":
		slots, err := m.mgmt.ListScheduleSlots(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"schedule": slots, "total": len(slots)})

	case r.Method == "POST" && slotID == "":
		var slot tournament.ScheduleSlot
		if err := httputil.DecodeJSON(r, &slot); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		slot.TournamentID = tournamentID
		created, err := m.mgmt.CreateScheduleSlot(r.Context(), &slot)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	case r.Method == "GET" && slotID != "":
		slot, err := m.mgmt.GetScheduleSlot(r.Context(), slotID)
		if err != nil {
			httputil.Error(w, http.StatusNotFound, "TOURNAMENT_404", "schedule slot not found")
			return
		}
		httputil.Success(w, http.StatusOK, slot)

	case r.Method == "PUT" && slotID != "":
		var slot tournament.ScheduleSlot
		if err := httputil.DecodeJSON(r, &slot); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		slot.ID = slotID
		updated, err := m.mgmt.UpdateScheduleSlot(r.Context(), &slot)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, updated)

	case r.Method == "DELETE" && slotID != "":
		if err := m.mgmt.DeleteScheduleSlot(r.Context(), slotID); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Arenas ──────────────────────────────────────────────────

func (m *Module) handleTournamentArenas(w http.ResponseWriter, r *http.Request, tournamentID, assignID string) {
	switch {
	case r.Method == "GET" && assignID == "":
		assigns, err := m.mgmt.ListArenaAssignments(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"arena_assignments": assigns, "total": len(assigns)})

	case r.Method == "POST" && assignID == "":
		var assign tournament.ArenaAssignment
		if err := httputil.DecodeJSON(r, &assign); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		assign.TournamentID = tournamentID
		created, err := m.mgmt.AssignArena(r.Context(), &assign)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	case r.Method == "DELETE" && assignID != "":
		if err := m.mgmt.RemoveArenaAssignment(r.Context(), assignID); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Results ─────────────────────────────────────────────────

func (m *Module) handleTournamentResults(w http.ResponseWriter, r *http.Request, tournamentID, resultID, action string) {
	switch {
	case r.Method == "GET" && resultID == "":
		results, err := m.mgmt.ListResults(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"results": results, "total": len(results)})

	case r.Method == "POST" && resultID == "":
		var result tournament.TournamentResult
		if err := httputil.DecodeJSON(r, &result); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		result.TournamentID = tournamentID
		created, err := m.mgmt.RecordResult(r.Context(), &result)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusCreated, created)

	case r.Method == "POST" && resultID != "" && action == "finalize":
		p, _ := httputil.GetPrincipal(r)
		if err := m.mgmt.FinalizeResult(r.Context(), resultID, p.User.ID); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"status": "finalized"})

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Standings ───────────────────────────────────────────────

func (m *Module) handleTournamentStandings(w http.ResponseWriter, r *http.Request, tournamentID string) {
	switch r.Method {
	case "GET":
		standings, err := m.mgmt.GetTeamStandings(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"standings": standings, "total": len(standings)})

	case "POST":
		var ts tournament.TeamStanding
		if err := httputil.DecodeJSON(r, &ts); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		ts.TournamentID = tournamentID
		updated, err := m.mgmt.UpdateTeamStanding(r.Context(), &ts)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, updated)

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Stats ───────────────────────────────────────────────────

func (m *Module) handleTournamentStats(w http.ResponseWriter, r *http.Request, tournamentID string) {
	if r.Method != "GET" {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	stats, err := m.mgmt.GetStats(r.Context(), tournamentID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}

// ── Export CSV ──────────────────────────────────────────────

func (m *Module) handleTournamentExport(w http.ResponseWriter, r *http.Request, tournamentID, entity string) {
	if r.Method != "GET" {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	var csvContent string
	var filename string

	switch entity {
	case "categories":
		cats, err := m.mgmt.ListCategories(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		csvContent = tournament.ExportCategoriesToCSV(cats)
		filename = "noi_dung_" + tournamentID + ".csv"

	case "registrations":
		regs, err := m.mgmt.ListRegistrations(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		csvContent = tournament.ExportRegistrationsToCSV(regs)
		filename = "dang_ky_" + tournamentID + ".csv"

	case "schedule":
		slots, err := m.mgmt.ListScheduleSlots(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		csvContent = tournament.ExportScheduleToCSV(slots)
		filename = "lich_thi_" + tournamentID + ".csv"

	case "results":
		results, err := m.mgmt.ListResults(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		csvContent = tournament.ExportResultsToCSV(results)
		filename = "ket_qua_" + tournamentID + ".csv"

	case "standings":
		standings, err := m.mgmt.GetTeamStandings(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		csvContent = tournament.ExportStandingsToCSV(standings)
		filename = "toan_doan_" + tournamentID + ".csv"

	default:
		httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", "unknown export entity: "+entity)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("\xEF\xBB\xBF")) // UTF-8 BOM
	_, _ = w.Write([]byte(csvContent))
}

// ── Batch Operations ───────────────────────────────────────

func (m *Module) handleTournamentBatch(w http.ResponseWriter, r *http.Request, tournamentID, action string) {
	if r.Method != "POST" {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}

	switch action {
	case "approve-registrations":
		var body struct {
			IDs []string `json:"ids"`
		}
		if err := httputil.DecodeJSON(r, &body); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		p, _ := httputil.GetPrincipal(r)
		result, err := m.mgmt.BatchApproveRegistrations(r.Context(), tournamentID, body.IDs, p.User.ID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, result)

	case "reject-registrations":
		var body struct {
			IDs    []string `json:"ids"`
			Reason string   `json:"reason"`
		}
		if err := httputil.DecodeJSON(r, &body); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		p, _ := httputil.GetPrincipal(r)
		result, err := m.mgmt.BatchRejectRegistrations(r.Context(), tournamentID, body.IDs, p.User.ID, body.Reason)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, result)

	case "finalize-results":
		var body struct {
			IDs []string `json:"ids"`
		}
		if err := httputil.DecodeJSON(r, &body); err != nil {
			httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", err.Error())
			return
		}
		p, _ := httputil.GetPrincipal(r)
		result, err := m.mgmt.BatchFinalizeResults(r.Context(), tournamentID, body.IDs, p.User.ID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, result)

	case "recalculate-standings":
		standings, err := m.mgmt.RecalculateTeamStandings(r.Context(), tournamentID)
		if err != nil {
			httputil.InternalError(w, err)
			return
		}
		httputil.Success(w, http.StatusOK, map[string]any{"standings": standings, "total": len(standings)})

	default:
		httputil.Error(w, http.StatusBadRequest, "TOURNAMENT_400", "unknown batch action: "+action)
	}
}
