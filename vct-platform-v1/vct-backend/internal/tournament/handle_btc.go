package tournament

import (
	"net/http"
	"strings"

	"vct-platform/backend/internal/domain/btc"
	"vct-platform/backend/internal/shared/httputil"
)

// RegisterBTCRoutes registers all BTC related routes.
func (m *Module) RegisterBTCRoutes(mux *http.ServeMux) {
	// BTC Members
	mux.HandleFunc("/api/v1/btc/members", m.handleBTCMemberList)
	mux.HandleFunc("/api/v1/btc/members/create", m.handleBTCMemberCreate)
	mux.HandleFunc("/api/v1/btc/members/", m.handleBTCMemberByID) // GET, PATCH, DELETE by ID

	// Weigh-In
	mux.HandleFunc("/api/v1/btc/weigh-in", m.handleBTCWeighInList)
	mux.HandleFunc("/api/v1/btc/weigh-in/create", m.handleBTCWeighInCreate)

	// Draw
	mux.HandleFunc("/api/v1/btc/draws", m.handleBTCDrawList)
	mux.HandleFunc("/api/v1/btc/draws/generate", m.handleBTCDrawGenerate)

	// Referee Assignments
	mux.HandleFunc("/api/v1/btc/referee-assignments", m.handleBTCAssignmentList)
	mux.HandleFunc("/api/v1/btc/referee-assignments/create", m.handleBTCAssignmentCreate)

	// Results
	mux.HandleFunc("/api/v1/btc/results", m.handleBTCTeamResults)
	mux.HandleFunc("/api/v1/btc/results/content", m.handleBTCContentResults)

	// Finance
	mux.HandleFunc("/api/v1/btc/finance", m.handleBTCFinanceList)
	mux.HandleFunc("/api/v1/btc/finance/create", m.handleBTCFinanceCreate)
	mux.HandleFunc("/api/v1/btc/finance/summary", m.handleBTCFinanceSummary)

	// Technical Meetings
	mux.HandleFunc("/api/v1/btc/meetings", m.handleBTCMeetingList)
	mux.HandleFunc("/api/v1/btc/meetings/create", m.handleBTCMeetingCreate)

	// Protests
	mux.HandleFunc("/api/v1/btc/protests", m.handleBTCProtestList)
	mux.HandleFunc("/api/v1/btc/protests/create", m.handleBTCProtestCreate)
	mux.HandleFunc("/api/v1/btc/protests/", m.handleBTCProtestUpdate)

	// Stats
	mux.HandleFunc("/api/v1/btc/stats", m.handleBTCStats)
}

// ── BTC Members ─────────────────────────────────────────────

func (m *Module) handleBTCMemberList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	members, err := m.btc.ListMembers(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, members)
}

func (m *Module) handleBTCMemberCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	var member btc.BTCMember
	if err := httputil.DecodeJSON(r, &member); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	if err := m.btc.CreateMember(r.Context(), &member); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, member)
}

func (m *Module) handleBTCMemberByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	memberID := parts[len(parts)-1]
	if memberID == "" || memberID == "members" {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", "missing member ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		member, err := m.btc.GetMember(r.Context(), memberID)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, member)

	case http.MethodPatch:
		var member btc.BTCMember
		if err := httputil.DecodeJSON(r, &member); err != nil {
			httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
			return
		}
		member.ID = memberID
		if err := m.btc.UpdateMember(r.Context(), &member); err != nil {
			httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, member)

	case http.MethodDelete:
		if err := m.btc.DeleteMember(r.Context(), memberID); err != nil {
			httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
			return
		}
		httputil.Success(w, http.StatusOK, map[string]string{"status": "deleted"})

	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
	}
}

// ── Weigh-In ────────────────────────────────────────────────

func (m *Module) handleBTCWeighInList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	records, err := m.btc.ListWeighIns(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, records)
}

func (m *Module) handleBTCWeighInCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	var record btc.WeighInRecord
	if err := httputil.DecodeJSON(r, &record); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	if err := m.btc.CreateWeighIn(r.Context(), &record); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, record)
}

// ── Draw ────────────────────────────────────────────────────

func (m *Module) handleBTCDrawList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	draws, err := m.btc.ListDraws(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, draws)
}

func (m *Module) handleBTCDrawGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	var input btc.DrawInput
	if err := httputil.DecodeJSON(r, &input); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	result, err := m.btc.GenerateDraw(r.Context(), input)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, result)
}

// ── Referee Assignments ─────────────────────────────────────

func (m *Module) handleBTCAssignmentList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	assignments, err := m.btc.ListAssignments(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, assignments)
}

func (m *Module) handleBTCAssignmentCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	var assignment btc.RefereeAssignment
	if err := httputil.DecodeJSON(r, &assignment); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	if err := m.btc.CreateAssignment(r.Context(), &assignment); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, assignment)
}

// ── Results ─────────────────────────────────────────────────

func (m *Module) handleBTCTeamResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	results, err := m.btc.ListTeamResults(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, results)
}

func (m *Module) handleBTCContentResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	results, err := m.btc.ListContentResults(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, results)
}

// ── Finance ─────────────────────────────────────────────────

func (m *Module) handleBTCFinanceList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	entries, err := m.btc.ListFinance(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, entries)
}

func (m *Module) handleBTCFinanceCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	var entry btc.FinanceEntry
	if err := httputil.DecodeJSON(r, &entry); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	if err := m.btc.CreateFinance(r.Context(), &entry); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, entry)
}

func (m *Module) handleBTCFinanceSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	summary, err := m.btc.FinanceSummary(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, summary)
}

// ── Technical Meetings ──────────────────────────────────────

func (m *Module) handleBTCMeetingList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	meetings, err := m.btc.ListMeetings(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, meetings)
}

func (m *Module) handleBTCMeetingCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	var meeting btc.TechnicalMeeting
	if err := httputil.DecodeJSON(r, &meeting); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	if err := m.btc.CreateMeeting(r.Context(), &meeting); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, meeting)
}

// ── Protests ────────────────────────────────────────────────

func (m *Module) handleBTCProtestList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	protests, err := m.btc.ListProtests(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, protests)
}

func (m *Module) handleBTCProtestCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	var protest btc.Protest
	if err := httputil.DecodeJSON(r, &protest); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	if err := m.btc.CreateProtest(r.Context(), &protest); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusCreated, protest)
}

func (m *Module) handleBTCProtestUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	protestID := parts[len(parts)-1]
	if protestID == "" || protestID == "protests" {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", "missing protest ID")
		return
	}

	var body struct {
		TrangThai string `json:"trang_thai"`
		NguoiXL   string `json:"nguoi_xl"`
		QuyetDinh string `json:"quyet_dinh"`
	}
	if err := httputil.DecodeJSON(r, &body); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}

	if err := m.btc.UpdateProtestStatus(r.Context(), protestID, body.TrangThai, body.NguoiXL, body.QuyetDinh); err != nil {
		httputil.Error(w, http.StatusBadRequest, "BTC_400", err.Error())
		return
	}
	httputil.Success(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ── Stats ───────────────────────────────────────────────────

func (m *Module) handleBTCStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "METHOD_405", "Method not allowed")
		return
	}
	giaiID := r.URL.Query().Get("giai_id")
	stats, err := m.btc.GetStats(r.Context(), giaiID)
	if err != nil {
		httputil.InternalError(w, err)
		return
	}
	httputil.Success(w, http.StatusOK, stats)
}
