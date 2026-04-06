package tournament

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — Tournament MgmtService Integration Tests
// Tests service-layer CRUD, batch operations, and stats
// using the in-memory store from adapter package.
// ═══════════════════════════════════════════════════════════════

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

var svcIDCounter atomic.Int64

func svcIDGen() string { return fmt.Sprintf("TSVC-%d", svcIDCounter.Add(1)) }

// newTestMgmtService creates a MgmtService backed by an in-memory repository for testing.
func newTestMgmtService(t *testing.T) *MgmtService {
	t.Helper()
	repo := NewInMemMgmtRepo()
	return NewMgmtService(repo, svcIDGen)
}

var svcCtx = context.Background()

// ── Category CRUD ────────────────────────────────────────────

func TestMgmtService_CategoryCRUD(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-TEST-001"

	// Create
	cat, err := svc.CreateCategory(svcCtx, &Category{
		TournamentID: tid,
		ContentType:  "doi_khang",
		AgeGroup:     "thanh_nien",
		WeightClass:  "60kg",
		Gender:       "nam",
		MaxAthletes:  32,
		MinAthletes:  4,
		Status:       "active",
	})
	if err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}
	if cat.ID == "" {
		t.Error("expected non-empty ID")
	}

	// List
	list, err := svc.ListCategories(svcCtx, tid)
	if err != nil {
		t.Fatalf("ListCategories: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1, got %d", len(list))
	}

	// Get
	got, err := svc.GetCategory(svcCtx, cat.ID)
	if err != nil {
		t.Fatalf("GetCategory: %v", err)
	}
	if got.ContentType != "doi_khang" {
		t.Errorf("expected doi_khang, got %s", got.ContentType)
	}

	// Update
	cat.MaxAthletes = 64
	updated, err := svc.UpdateCategory(svcCtx, cat)
	if err != nil {
		t.Fatalf("UpdateCategory: %v", err)
	}
	if updated.MaxAthletes != 64 {
		t.Errorf("expected 64, got %d", updated.MaxAthletes)
	}

	// Delete
	if err := svc.DeleteCategory(svcCtx, cat.ID); err != nil {
		t.Fatalf("DeleteCategory: %v", err)
	}
	list2, _ := svc.ListCategories(svcCtx, tid)
	if len(list2) != 0 {
		t.Errorf("expected 0 after delete, got %d", len(list2))
	}
}

// ── Registration Workflow ────────────────────────────────────

func TestMgmtService_RegistrationWorkflow(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-TEST-002"

	// Register team
	reg, err := svc.RegisterTeam(svcCtx, &Registration{
		TournamentID: tid,
		TeamName:     "Đoàn Hà Nội",
		TeamType:     "doan_tinh",
		Province:     "Hà Nội",
		HeadCoach:    "HLV Nguyễn Văn A",
		Status:       "nhap",
	})
	if err != nil {
		t.Fatalf("RegisterTeam: %v", err)
	}

	// Add athlete
	ath, err := svc.AddAthleteToRegistration(svcCtx, &RegistrationAthlete{
		RegistrationID: reg.ID,
		AthleteName:    "Trần Văn B",
		Gender:         "nam",
		Weight:         60,
		CategoryIDs:    []string{"cat-1"},
	})
	if err != nil {
		t.Fatalf("AddAthlete: %v", err)
	}
	if ath.ID == "" {
		t.Error("expected athlete ID")
	}

	// List athletes
	athletes, err := svc.ListRegistrationAthletes(svcCtx, reg.ID)
	if err != nil {
		t.Fatalf("ListRegistrationAthletes: %v", err)
	}
	if len(athletes) != 1 {
		t.Errorf("expected 1 athlete, got %d", len(athletes))
	}

	// Submit
	if err := svc.SubmitRegistration(svcCtx, reg.ID); err != nil {
		t.Fatalf("SubmitRegistration: %v", err)
	}
	got, _ := svc.GetRegistration(svcCtx, reg.ID)
	if got.Status != "cho_duyet" {
		t.Errorf("expected cho_duyet after submit, got %s", got.Status)
	}

	// Approve
	if err := svc.ApproveRegistration(svcCtx, reg.ID, "admin-001"); err != nil {
		t.Fatalf("ApproveRegistration: %v", err)
	}
	got2, _ := svc.GetRegistration(svcCtx, reg.ID)
	if got2.Status != "da_duyet" {
		t.Errorf("expected da_duyet after approve, got %s", got2.Status)
	}
	if got2.ApprovedBy != "admin-001" {
		t.Errorf("expected approvedBy admin-001, got %s", got2.ApprovedBy)
	}
}

func TestMgmtService_RejectRegistration(t *testing.T) {
	svc := newTestMgmtService(t)

	reg, _ := svc.RegisterTeam(svcCtx, &Registration{
		TournamentID: "TRN-REJ",
		TeamName:     "Đoàn Bị Từ Chối",
		TeamType:     "clb",
		HeadCoach:    "HLV X",
		Status:       "nhap",
	})
	_ = svc.SubmitRegistration(svcCtx, reg.ID)

	if err := svc.RejectRegistration(svcCtx, reg.ID, "admin-002", "Thiếu hồ sơ"); err != nil {
		t.Fatalf("RejectRegistration: %v", err)
	}
	got, _ := svc.GetRegistration(svcCtx, reg.ID)
	if got.Status != "tu_choi" {
		t.Errorf("expected tu_choi, got %s", got.Status)
	}
	if got.RejectReason != "Thiếu hồ sơ" {
		t.Errorf("expected reason 'Thiếu hồ sơ', got %s", got.RejectReason)
	}
}

// ── Schedule CRUD ────────────────────────────────────────────

func TestMgmtService_ScheduleCRUD(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-SCH-001"

	slot, err := svc.CreateScheduleSlot(svcCtx, &ScheduleSlot{
		TournamentID: tid,
		ArenaID:      "ARENA-001",
		ArenaName:    "Sân A1",
		Date:         "2026-06-15",
		Session:      "sang",
		StartTime:    "08:00",
		EndTime:      "12:00",
		Status:       "du_kien",
	})
	if err != nil {
		t.Fatalf("CreateScheduleSlot: %v", err)
	}

	list, _ := svc.ListScheduleSlots(svcCtx, tid)
	if len(list) != 1 {
		t.Errorf("expected 1 slot, got %d", len(list))
	}

	got, _ := svc.GetScheduleSlot(svcCtx, slot.ID)
	if got.Session != "sang" {
		t.Errorf("expected sang, got %s", got.Session)
	}

	slot.MatchCount = 10
	updated, err := svc.UpdateScheduleSlot(svcCtx, slot)
	if err != nil {
		t.Fatalf("UpdateScheduleSlot: %v", err)
	}
	if updated.MatchCount != 10 {
		t.Errorf("expected 10, got %d", updated.MatchCount)
	}

	if err := svc.DeleteScheduleSlot(svcCtx, slot.ID); err != nil {
		t.Fatalf("DeleteScheduleSlot: %v", err)
	}
}

// ── Arena Assignment ─────────────────────────────────────────

func TestMgmtService_ArenaAssignment(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-AR-001"

	assign, err := svc.AssignArena(svcCtx, &ArenaAssignment{
		TournamentID: tid,
		ArenaID:      "ARENA-001",
		ArenaName:    "Sân A1",
		Date:         "2026-06-15",
		Session:      "sang",
		ContentTypes: []string{"doi_khang"},
	})
	if err != nil {
		t.Fatalf("AssignArena: %v", err)
	}

	list, _ := svc.ListArenaAssignments(svcCtx, tid)
	if len(list) != 1 {
		t.Errorf("expected 1, got %d", len(list))
	}

	if err := svc.RemoveArenaAssignment(svcCtx, assign.ID); err != nil {
		t.Fatalf("RemoveArenaAssignment: %v", err)
	}
}

// ── Results & Finalize ───────────────────────────────────────

func TestMgmtService_ResultsAndFinalize(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-RES-001"

	result, err := svc.RecordResult(svcCtx, &TournamentResult{
		TournamentID: tid,
		CategoryID:   "cat-dk-60",
		CategoryName: "Đối kháng Nam 60kg",
		ContentType:  "doi_khang",
		GoldName:     "Nguyễn Văn A",
		GoldTeam:     "Hà Nội",
		SilverName:   "Trần Văn B",
		SilverTeam:   "TP.HCM",
	})
	if err != nil {
		t.Fatalf("RecordResult: %v", err)
	}
	if result.IsFinalized {
		t.Error("result should not be finalized initially")
	}

	// Finalize
	if err := svc.FinalizeResult(svcCtx, result.ID, "admin-001"); err != nil {
		t.Fatalf("FinalizeResult: %v", err)
	}

	list, _ := svc.ListResults(svcCtx, tid)
	if len(list) != 1 {
		t.Errorf("expected 1 result, got %d", len(list))
	}
	if !list[0].IsFinalized {
		t.Error("expected result to be finalized")
	}
}

// ── Team Standings ───────────────────────────────────────────

func TestMgmtService_TeamStandings(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-STD-001"

	// Record and finalize results for recalculation
	res1, _ := svc.RecordResult(svcCtx, &TournamentResult{TournamentID: tid, CategoryID: "c1", GoldName: "A", GoldTeam: "HN", SilverName: "B", SilverTeam: "SG", ContentType: "doi_khang"})
	res2, _ := svc.RecordResult(svcCtx, &TournamentResult{TournamentID: tid, CategoryID: "c2", GoldName: "C", GoldTeam: "HN", SilverName: "D", SilverTeam: "DN", ContentType: "quyen"})

	// Must finalize before recalculation (RecalculateTeamStandings only counts finalized)
	_ = svc.FinalizeResult(svcCtx, res1.ID, "admin")
	_ = svc.FinalizeResult(svcCtx, res2.ID, "admin")

	// Recalculate
	standings, err := svc.RecalculateTeamStandings(svcCtx, tid)
	if err != nil {
		t.Fatalf("RecalculateTeamStandings: %v", err)
	}
	if len(standings) == 0 {
		t.Error("expected standings after recalculation")
	}

	// HN should have 2 golds = 14 points, SG 1 silver = 5, DN 1 silver = 5
	// Also via GetTeamStandings
	list, err := svc.GetTeamStandings(svcCtx, tid)
	if err != nil {
		t.Fatalf("GetTeamStandings: %v", err)
	}
	if len(list) == 0 {
		t.Error("expected standings")
	}
}

// ── Stats ────────────────────────────────────────────────────

func TestMgmtService_Stats(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-STAT-001"

	svc.CreateCategory(svcCtx, &Category{TournamentID: tid, ContentType: "doi_khang", AgeGroup: "thanh_nien", Gender: "nam", Status: "active"})
	svc.CreateCategory(svcCtx, &Category{TournamentID: tid, ContentType: "quyen", AgeGroup: "thieu_nien", Gender: "nu", Status: "active"})
	svc.RegisterTeam(svcCtx, &Registration{TournamentID: tid, TeamName: "T1", TeamType: "doan_tinh", HeadCoach: "C1", Status: "nhap"})

	stats, err := svc.GetStats(svcCtx, tid)
	if err != nil {
		t.Fatalf("GetStats: %v", err)
	}
	if stats.TotalCategories != 2 {
		t.Errorf("expected 2 categories, got %d", stats.TotalCategories)
	}
	if stats.TotalRegistrations != 1 {
		t.Errorf("expected 1 registration, got %d", stats.TotalRegistrations)
	}
}

// ── Batch Operations ─────────────────────────────────────────

func TestMgmtService_BatchApproveRegistrations(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-BATCH-001"

	r1, _ := svc.RegisterTeam(svcCtx, &Registration{TournamentID: tid, TeamName: "T1", TeamType: "doan_tinh", HeadCoach: "C1", Status: "nhap"})
	r2, _ := svc.RegisterTeam(svcCtx, &Registration{TournamentID: tid, TeamName: "T2", TeamType: "clb", HeadCoach: "C2", Status: "nhap"})

	_ = svc.SubmitRegistration(svcCtx, r1.ID)
	_ = svc.SubmitRegistration(svcCtx, r2.ID)

	result, err := svc.BatchApproveRegistrations(svcCtx, tid, []string{r1.ID, r2.ID}, "admin-batch")
	if err != nil {
		t.Fatalf("BatchApprove: %v", err)
	}
	if result.Success != 2 {
		t.Errorf("expected 2 succeeded, got %d", result.Success)
	}
}

func TestMgmtService_BatchFinalizeResults(t *testing.T) {
	svc := newTestMgmtService(t)
	tid := "TRN-BATCH-002"

	res1, _ := svc.RecordResult(svcCtx, &TournamentResult{TournamentID: tid, CategoryID: "c1", GoldName: "A", SilverName: "B", ContentType: "doi_khang"})
	res2, _ := svc.RecordResult(svcCtx, &TournamentResult{TournamentID: tid, CategoryID: "c2", GoldName: "C", SilverName: "D", ContentType: "quyen"})

	result, err := svc.BatchFinalizeResults(svcCtx, tid, []string{res1.ID, res2.ID}, "admin-batch")
	if err != nil {
		t.Fatalf("BatchFinalize: %v", err)
	}
	if result.Success != 2 {
		t.Errorf("expected 2 succeeded, got %d", result.Success)
	}
}
