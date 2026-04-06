package btc

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"vct-platform/backend/internal/domain"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — BTC SERVICE (Business Logic)
// ═══════════════════════════════════════════════════════════════

type Store interface {
	// BTC Members
	ListMembers(ctx context.Context, giaiID string) ([]BTCMember, error)
	GetMember(ctx context.Context, id string) (*BTCMember, error)
	CreateMember(ctx context.Context, m *BTCMember) error
	UpdateMember(ctx context.Context, m *BTCMember) error
	DeleteMember(ctx context.Context, id string) error

	// Weigh-In
	ListWeighIns(ctx context.Context, giaiID string) ([]WeighInRecord, error)
	CreateWeighIn(ctx context.Context, w *WeighInRecord) error

	// Draw
	ListDraws(ctx context.Context, giaiID string) ([]DrawResult, error)
	CreateDraw(ctx context.Context, d *DrawResult) error

	// Referee Assignment
	ListAssignments(ctx context.Context, giaiID string) ([]RefereeAssignment, error)
	CreateAssignment(ctx context.Context, a *RefereeAssignment) error

	// Results
	ListTeamResults(ctx context.Context, giaiID string) ([]TeamResult, error)
	CreateTeamResult(ctx context.Context, r *TeamResult) error
	ListContentResults(ctx context.Context, giaiID string) ([]ContentResult, error)

	// Finance
	ListFinance(ctx context.Context, giaiID string) ([]FinanceEntry, error)
	CreateFinance(ctx context.Context, f *FinanceEntry) error
	UpdateFinance(ctx context.Context, f *FinanceEntry) error

	// Technical Meeting
	ListMeetings(ctx context.Context, giaiID string) ([]TechnicalMeeting, error)
	CreateMeeting(ctx context.Context, m *TechnicalMeeting) error

	// Protests
	ListProtests(ctx context.Context, giaiID string) ([]Protest, error)
	GetProtest(ctx context.Context, id string) (*Protest, error)
	CreateProtest(ctx context.Context, p *Protest) error
	UpdateProtest(ctx context.Context, p *Protest) error
}

type Service struct {
	store Store
	idGen func() string
}

func NewService(store Store, idGen func() string) *Service {
	return &Service{store: store, idGen: idGen}
}

// ── BTC Members ─────────────────────────────────────────────

func (s *Service) ListMembers(ctx context.Context, giaiID string) ([]BTCMember, error) {
	return s.store.ListMembers(ctx, giaiID)
}

func (s *Service) GetMember(ctx context.Context, id string) (*BTCMember, error) {
	return s.store.GetMember(ctx, id)
}

func (s *Service) CreateMember(ctx context.Context, m *BTCMember) error {
	if m.Ten == "" {
		return fmt.Errorf("tên thành viên BTC là bắt buộc")
	}
	if m.Ban == "" {
		return fmt.Errorf("ban là bắt buộc")
	}
	m.ID = s.idGen()
	m.IsActive = true
	return s.store.CreateMember(ctx, m)
}

func (s *Service) UpdateMember(ctx context.Context, m *BTCMember) error {
	if m.ID == "" {
		return fmt.Errorf("ID thành viên là bắt buộc")
	}
	return s.store.UpdateMember(ctx, m)
}

func (s *Service) DeleteMember(ctx context.Context, id string) error {
	return s.store.DeleteMember(ctx, id)
}

// ── Weigh-In ────────────────────────────────────────────────

func (s *Service) ListWeighIns(ctx context.Context, giaiID string) ([]WeighInRecord, error) {
	return s.store.ListWeighIns(ctx, giaiID)
}

func (s *Service) CreateWeighIn(ctx context.Context, w *WeighInRecord) error {
	if w.VdvID == "" || w.HangCan == "" {
		return fmt.Errorf("VĐV và hạng cân là bắt buộc")
	}
	if w.CanNang <= 0 {
		return fmt.Errorf("cân nặng phải lớn hơn 0")
	}
	w.ID = s.idGen()
	w.CreatedAt = time.Now().UTC()
	w.ThoiGian = time.Now().UTC()

	// BUG FIX: Use math.Abs to catch both overweight AND underweight
	diff := math.Abs(w.CanNang - w.GioiHan)
	if diff <= w.SaiSo {
		w.KetQua = "dat"
	} else {
		w.KetQua = "khong_dat"
		w.GhiChu = fmt.Sprintf("Lệch %.1fkg so với giới hạn %.1fkg", diff, w.GioiHan)
	}
	return s.store.CreateWeighIn(ctx, w)
}

// ── Draw ────────────────────────────────────────────────────

func (s *Service) ListDraws(ctx context.Context, giaiID string) ([]DrawResult, error) {
	return s.store.ListDraws(ctx, giaiID)
}

type DrawInput struct {
	GiaiID     string       `json:"giai_id"`
	NoiDungID  string       `json:"noi_dung_id"`
	NoiDungTen string       `json:"noi_dung_ten"`
	LoaiND     string       `json:"loai_nd"`
	HangCan    string       `json:"hang_can"`
	LuaTuoi    string       `json:"lua_tuoi"`
	Athletes   []DrawBranch `json:"athletes"`
	CreatedBy  string       `json:"created_by"`
}

func (s *Service) GenerateDraw(ctx context.Context, input DrawInput) (*DrawResult, error) {
	if len(input.Athletes) < 2 {
		return nil, fmt.Errorf("cần ít nhất 2 VĐV để bốc thăm")
	}

	d := &DrawResult{
		ID:         s.idGen(),
		GiaiID:     input.GiaiID,
		NoiDungID:  input.NoiDungID,
		NoiDungTen: input.NoiDungTen,
		LoaiND:     input.LoaiND,
		HangCan:    input.HangCan,
		LuaTuoi:    input.LuaTuoi,
		SoVDV:      len(input.Athletes),
		CreatedAt:  time.Now().UTC(),
		CreatedBy:  input.CreatedBy,
	}

	// Shuffle athletes randomly
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]DrawBranch, len(input.Athletes))
	copy(shuffled, input.Athletes)
	rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	if input.LoaiND == "doi_khang" {
		// Assign bracket positions
		for i := range shuffled {
			shuffled[i].Position = i + 1
		}
		d.Nhanh = shuffled
	} else {
		// Quyền: assign performance order
		orders := make([]DrawOrder, len(shuffled))
		for i, a := range shuffled {
			orders[i] = DrawOrder{
				ThuTu:   i + 1,
				VdvID:   a.VdvID,
				VdvTen:  a.VdvTen,
				DoanTen: a.DoanTen,
			}
		}
		d.ThuTu = orders
	}

	if err := s.store.CreateDraw(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

// ── Referee Assignment ──────────────────────────────────────

func (s *Service) ListAssignments(ctx context.Context, giaiID string) ([]RefereeAssignment, error) {
	return s.store.ListAssignments(ctx, giaiID)
}

func (s *Service) CreateAssignment(ctx context.Context, a *RefereeAssignment) error {
	if a.TrongTaiID == "" || a.SanID == "" {
		return fmt.Errorf("trọng tài và sàn là bắt buộc")
	}

	// Check for conflicts (same referee, same time slot)
	existing, err := s.store.ListAssignments(ctx, a.GiaiID)
	if err != nil {
		return err
	}
	for _, ex := range existing {
		if ex.TrongTaiID == a.TrongTaiID && ex.Ngay == a.Ngay && ex.Phien == a.Phien {
			return fmt.Errorf("trọng tài %s đã được phân công trong phiên %s ngày %s", a.TrongTaiTen, a.Phien, a.Ngay)
		}
	}

	a.ID = s.idGen()
	a.TrangThai = "phan_cong"
	a.CreatedAt = time.Now().UTC()
	return s.store.CreateAssignment(ctx, a)
}

// ── Results ─────────────────────────────────────────────────

func (s *Service) ListTeamResults(ctx context.Context, giaiID string) ([]TeamResult, error) {
	results, err := s.store.ListTeamResults(ctx, giaiID)
	if err != nil {
		return nil, err
	}
	// Sort by ranking
	sort.Slice(results, func(i, j int) bool {
		return results[i].XepHang < results[j].XepHang
	})
	return results, nil
}

func (s *Service) ListContentResults(ctx context.Context, giaiID string) ([]ContentResult, error) {
	return s.store.ListContentResults(ctx, giaiID)
}

// ── Finance ─────────────────────────────────────────────────

func (s *Service) ListFinance(ctx context.Context, giaiID string) ([]FinanceEntry, error) {
	return s.store.ListFinance(ctx, giaiID)
}

func (s *Service) CreateFinance(ctx context.Context, f *FinanceEntry) error {
	if f.SoTien <= 0 {
		return fmt.Errorf("số tiền phải lớn hơn 0")
	}
	if f.Loai != "thu" && f.Loai != "chi" {
		return fmt.Errorf("loại phải là 'thu' hoặc 'chi'")
	}
	f.ID = s.idGen()
	f.CreatedAt = time.Now().UTC()
	return s.store.CreateFinance(ctx, f)
}

func (s *Service) UpdateFinance(ctx context.Context, f *FinanceEntry) error {
	if f.ID == "" {
		return fmt.Errorf("ID bút toán là bắt buộc")
	}
	return s.store.UpdateFinance(ctx, f)
}

// FinanceSummary returns aggregated finance info for a tournament.
type FinanceSummaryResult struct {
	TongThu float64 `json:"tong_thu"`
	TongChi float64 `json:"tong_chi"`
	SoDu    float64 `json:"so_du"`
	SoBut   int     `json:"so_but"`
}

func (s *Service) FinanceSummary(ctx context.Context, giaiID string) (*FinanceSummaryResult, error) {
	entries, err := s.store.ListFinance(ctx, giaiID)
	if err != nil {
		return nil, err
	}
	result := &FinanceSummaryResult{SoBut: len(entries)}
	for _, f := range entries {
		if f.Loai == "thu" {
			result.TongThu += f.SoTien
		} else {
			result.TongChi += f.SoTien
		}
	}
	result.SoDu = result.TongThu - result.TongChi
	return result, nil
}

// ── Technical Meeting ───────────────────────────────────────

func (s *Service) ListMeetings(ctx context.Context, giaiID string) ([]TechnicalMeeting, error) {
	return s.store.ListMeetings(ctx, giaiID)
}

func (s *Service) CreateMeeting(ctx context.Context, m *TechnicalMeeting) error {
	if m.TieuDe == "" {
		return fmt.Errorf("tiêu đề cuộc họp là bắt buộc")
	}
	m.ID = s.idGen()
	m.TrangThai = "du_kien"
	m.CreatedAt = time.Now().UTC()
	return s.store.CreateMeeting(ctx, m)
}

// ── Protests ────────────────────────────────────────────────

func (s *Service) ListProtests(ctx context.Context, giaiID string) ([]Protest, error) {
	return s.store.ListProtests(ctx, giaiID)
}

func (s *Service) CreateProtest(ctx context.Context, p *Protest) error {
	if p.LyDo == "" {
		return fmt.Errorf("lý do khiếu nại là bắt buộc")
	}
	p.ID = s.idGen()
	p.TrangThai = "moi"
	p.NgayNop = time.Now().UTC()
	p.CreatedAt = time.Now().UTC()
	return s.store.CreateProtest(ctx, p)
}

func (s *Service) UpdateProtestStatus(ctx context.Context, id, newStatus, nguoiXL, quyetDinh string) error {
	p, err := s.store.GetProtest(ctx, id)
	if err != nil {
		return fmt.Errorf("không tìm thấy khiếu nại: %w", err)
	}

	// FIX: Use global ProtestTransitions from state_machine.go instead of duplicate local map
	if !domain.ProtestTransitions.CanTransition(p.TrangThai, newStatus) {
		return fmt.Errorf("không thể chuyển từ '%s' sang '%s'", p.TrangThai, newStatus)
	}

	p.TrangThai = newStatus
	p.NguoiXL = nguoiXL
	p.QuyetDinh = quyetDinh
	now := time.Now().UTC()
	p.NgayXL = &now

	return s.store.UpdateProtest(ctx, p)
}

// ── Stats ───────────────────────────────────────────────────

func (s *Service) GetStats(ctx context.Context, giaiID string) (*BTCStats, error) {
	stats := &BTCStats{}

	// FIX: Propagate errors instead of silently swallowing them
	weighIns, err := s.store.ListWeighIns(ctx, giaiID)
	if err != nil {
		return nil, fmt.Errorf("lỗi tải weigh-in: %w", err)
	}
	for _, w := range weighIns {
		if w.KetQua == "dat" {
			stats.DaCanKy++
		} else if w.KetQua == "cho_can" {
			stats.ChuaCanKy++
		}
	}
	totalCan := stats.DaCanKy + stats.ChuaCanKy
	if totalCan > 0 {
		stats.TyLeDatCan = float64(stats.DaCanKy) / float64(totalCan) * 100
	}
	stats.TongVDV = len(weighIns)

	// Referee stats
	assignments, err := s.store.ListAssignments(ctx, giaiID)
	if err != nil {
		return nil, fmt.Errorf("lỗi tải assignments: %w", err)
	}
	stats.DaPhanCong = len(assignments)
	// Count unique referees
	uniqTT := map[string]bool{}
	for _, a := range assignments {
		uniqTT[a.TrongTaiID] = true
	}
	stats.TongTrongTai = len(uniqTT)

	// Protest stats
	protests, err := s.store.ListProtests(ctx, giaiID)
	if err != nil {
		return nil, fmt.Errorf("lỗi tải protests: %w", err)
	}
	stats.TongKhieuNai = len(protests)
	for _, p := range protests {
		if p.TrangThai == "moi" || p.TrangThai == "tiep_nhan" || p.TrangThai == "xem_xet" {
			stats.KNChoXuLy++
		}
	}

	// Finance stats
	finances, err := s.store.ListFinance(ctx, giaiID)
	if err != nil {
		return nil, fmt.Errorf("lỗi tải finance: %w", err)
	}
	for _, f := range finances {
		if f.Loai == "thu" {
			stats.TongThu += f.SoTien
		} else {
			stats.TongChi += f.SoTien
		}
	}

	// Results stats
	results, err := s.store.ListTeamResults(ctx, giaiID)
	if err != nil {
		return nil, fmt.Errorf("lỗi tải results: %w", err)
	}
	stats.TongDoan = len(results)
	for _, r := range results {
		stats.TongHuyChuong += r.TongHC
		stats.TongTran += r.HCV + r.HCB + r.HCD // approximate match count
	}

	return stats, nil
}
