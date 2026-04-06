package btc

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — BTC IN-MEMORY STORE
// Thread-safe in-memory store with seed data for BTC domain.
// ═══════════════════════════════════════════════════════════════

type InMemStore struct {
	mu          sync.RWMutex
	members     []BTCMember
	weighIns    []WeighInRecord
	draws       []DrawResult
	assignments []RefereeAssignment
	teamResults []TeamResult
	contResults []ContentResult
	finance     []FinanceEntry
	meetings    []TechnicalMeeting
	protests    []Protest
}

func NewInMemStore() *InMemStore {
	s := &InMemStore{}
	s.seed()
	return s
}

// ── BTC Members ─────────────────────────────────────────────

func (s *InMemStore) ListMembers(_ context.Context, giaiID string) ([]BTCMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []BTCMember
	for _, m := range s.members {
		if m.GiaiID == giaiID || giaiID == "" {
			out = append(out, m)
		}
	}
	return out, nil
}

func (s *InMemStore) CreateMember(_ context.Context, m *BTCMember) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.members = append(s.members, *m)
	return nil
}

func (s *InMemStore) GetMember(_ context.Context, id string) (*BTCMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.members {
		if s.members[i].ID == id {
			m := s.members[i]
			return &m, nil
		}
	}
	return nil, fmt.Errorf("member not found: %s", id)
}

func (s *InMemStore) UpdateMember(_ context.Context, m *BTCMember) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.members {
		if s.members[i].ID == m.ID {
			s.members[i] = *m
			return nil
		}
	}
	return fmt.Errorf("member not found: %s", m.ID)
}

func (s *InMemStore) DeleteMember(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.members {
		if s.members[i].ID == id {
			s.members = append(s.members[:i], s.members[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("member not found: %s", id)
}

// ── Weigh-In ────────────────────────────────────────────────

func (s *InMemStore) ListWeighIns(_ context.Context, giaiID string) ([]WeighInRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []WeighInRecord
	for _, w := range s.weighIns {
		if w.GiaiID == giaiID || giaiID == "" {
			out = append(out, w)
		}
	}
	return out, nil
}

func (s *InMemStore) CreateWeighIn(_ context.Context, w *WeighInRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weighIns = append(s.weighIns, *w)
	return nil
}

// ── Draw ────────────────────────────────────────────────────

func (s *InMemStore) ListDraws(_ context.Context, giaiID string) ([]DrawResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []DrawResult
	for _, d := range s.draws {
		if d.GiaiID == giaiID || giaiID == "" {
			out = append(out, d)
		}
	}
	return out, nil
}

func (s *InMemStore) CreateDraw(_ context.Context, d *DrawResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.draws = append(s.draws, *d)
	return nil
}

// ── Referee Assignment ──────────────────────────────────────

func (s *InMemStore) ListAssignments(_ context.Context, giaiID string) ([]RefereeAssignment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []RefereeAssignment
	for _, a := range s.assignments {
		if a.GiaiID == giaiID || giaiID == "" {
			out = append(out, a)
		}
	}
	return out, nil
}

func (s *InMemStore) CreateAssignment(_ context.Context, a *RefereeAssignment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.assignments = append(s.assignments, *a)
	return nil
}

// ── Results ─────────────────────────────────────────────────

func (s *InMemStore) ListTeamResults(_ context.Context, giaiID string) ([]TeamResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []TeamResult
	for _, r := range s.teamResults {
		if r.GiaiID == giaiID || giaiID == "" {
			out = append(out, r)
		}
	}
	return out, nil
}

func (s *InMemStore) CreateTeamResult(_ context.Context, r *TeamResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.teamResults = append(s.teamResults, *r)
	return nil
}

func (s *InMemStore) ListContentResults(_ context.Context, giaiID string) ([]ContentResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []ContentResult
	for _, r := range s.contResults {
		if r.GiaiID == giaiID || giaiID == "" {
			out = append(out, r)
		}
	}
	return out, nil
}

// ── Finance ─────────────────────────────────────────────────

func (s *InMemStore) ListFinance(_ context.Context, giaiID string) ([]FinanceEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []FinanceEntry
	for _, f := range s.finance {
		if f.GiaiID == giaiID || giaiID == "" {
			out = append(out, f)
		}
	}
	return out, nil
}

func (s *InMemStore) CreateFinance(_ context.Context, f *FinanceEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.finance = append(s.finance, *f)
	return nil
}

func (s *InMemStore) UpdateFinance(_ context.Context, f *FinanceEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.finance {
		if s.finance[i].ID == f.ID {
			s.finance[i] = *f
			return nil
		}
	}
	return fmt.Errorf("finance entry not found: %s", f.ID)
}

// ── Technical Meeting ───────────────────────────────────────

func (s *InMemStore) ListMeetings(_ context.Context, giaiID string) ([]TechnicalMeeting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []TechnicalMeeting
	for _, m := range s.meetings {
		if m.GiaiID == giaiID || giaiID == "" {
			out = append(out, m)
		}
	}
	return out, nil
}

func (s *InMemStore) CreateMeeting(_ context.Context, m *TechnicalMeeting) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.meetings = append(s.meetings, *m)
	return nil
}

// ── Protests ────────────────────────────────────────────────

func (s *InMemStore) ListProtests(_ context.Context, giaiID string) ([]Protest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Protest
	for _, p := range s.protests {
		if p.GiaiID == giaiID || giaiID == "" {
			out = append(out, p)
		}
	}
	return out, nil
}

func (s *InMemStore) GetProtest(_ context.Context, id string) (*Protest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.protests {
		if s.protests[i].ID == id {
			p := s.protests[i]
			return &p, nil
		}
	}
	return nil, fmt.Errorf("protest not found: %s", id)
}

func (s *InMemStore) CreateProtest(_ context.Context, p *Protest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.protests = append(s.protests, *p)
	return nil
}

func (s *InMemStore) UpdateProtest(_ context.Context, p *Protest) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.protests {
		if s.protests[i].ID == p.ID {
			s.protests[i] = *p
			return nil
		}
	}
	return fmt.Errorf("protest not found: %s", p.ID)
}

// ═══════════════════════════════════════════════════════════════
// SEED DATA
// ═══════════════════════════════════════════════════════════════

func (s *InMemStore) seed() {
	gID := "giai-vct-2025"
	now := time.Now().UTC()

	// ── BTC Members ─────────────────────────────────────────
	s.members = []BTCMember{
		{ID: "btc-001", Ten: "Nguyễn Văn Hùng", ChucVu: "Trưởng ban", Ban: "ban_to_chuc", Cap: 1, Sdt: "0901234567", Email: "hung.nv@vct.vn", DonVi: "Liên đoàn VCT Việt Nam", GiaiID: gID, IsActive: true},
		{ID: "btc-002", Ten: "Trần Minh Đức", ChucVu: "Phó ban", Ban: "ban_to_chuc", Cap: 2, Sdt: "0901234568", Email: "duc.tm@vct.vn", DonVi: "Sở VH-TT TP.HCM", GiaiID: gID, IsActive: true},
		{ID: "btc-003", Ten: "Lê Hoàng Nam", ChucVu: "Trưởng ban", Ban: "ban_chuyen_mon", Cap: 1, Sdt: "0901234569", Email: "nam.lh@vct.vn", DonVi: "Liên đoàn VCT Việt Nam", GiaiID: gID, IsActive: true},
		{ID: "btc-004", Ten: "Phạm Thị Hồng", ChucVu: "Phó ban", Ban: "ban_chuyen_mon", Cap: 2, Sdt: "0901234570", Email: "hong.pt@vct.vn", DonVi: "ĐH TDTT TP.HCM", GiaiID: gID, IsActive: true},
		{ID: "btc-005", Ten: "Võ Thanh Sơn", ChucVu: "Trưởng ban", Ban: "ban_trong_tai", Cap: 1, Sdt: "0901234571", Email: "son.vt@vct.vn", DonVi: "Liên đoàn VCT Việt Nam", GiaiID: gID, IsActive: true},
		{ID: "btc-006", Ten: "Đặng Quốc Việt", ChucVu: "Phó ban", Ban: "ban_trong_tai", Cap: 2, Sdt: "0901234572", Email: "viet.dq@vct.vn", DonVi: "HLV Quốc gia", GiaiID: gID, IsActive: true},
		{ID: "btc-007", Ten: "BS. Nguyễn Thị Mai", ChucVu: "Trưởng ban", Ban: "ban_y_te", Cap: 1, Sdt: "0901234573", Email: "mai.nt@vct.vn", DonVi: "BV Thể thao", GiaiID: gID, IsActive: true},
		{ID: "btc-008", Ten: "Hoàng Minh Tuấn", ChucVu: "Trưởng ban", Ban: "ban_khang_nghi", Cap: 1, Sdt: "0901234574", Email: "tuan.hm@vct.vn", DonVi: "Liên đoàn VCT Việt Nam", GiaiID: gID, IsActive: true},
		{ID: "btc-009", Ten: "Lý Văn Phong", ChucVu: "Ủy viên", Ban: "ban_to_chuc", Cap: 3, Sdt: "0901234575", Email: "phong.lv@vct.vn", DonVi: "Sở VH-TT Đà Nẵng", GiaiID: gID, IsActive: true},
		{ID: "btc-010", Ten: "Trương Quang Hải", ChucVu: "Ủy viên", Ban: "ban_chuyen_mon", Cap: 3, Sdt: "0901234576", Email: "hai.tq@vct.vn", DonVi: "CLB VCT Bình Định", GiaiID: gID, IsActive: true},
	}

	// ── Weigh-In Records ────────────────────────────────────
	s.weighIns = []WeighInRecord{
		{ID: "wi-001", GiaiID: gID, VdvID: "vdv-001", VdvTen: "Nguyễn Anh Tuấn", DoanID: "doan-hcm", DoanTen: "TP.HCM", HangCan: "54kg Nam", CanNang: 53.8, GioiHan: 54, SaiSo: 0.5, KetQua: "dat", LanCan: 1, NguoiCan: "BS. Mai", ThoiGian: now, CreatedAt: now},
		{ID: "wi-002", GiaiID: gID, VdvID: "vdv-002", VdvTen: "Trần Minh Phúc", DoanID: "doan-hn", DoanTen: "Hà Nội", HangCan: "60kg Nam", CanNang: 60.2, GioiHan: 60, SaiSo: 0.5, KetQua: "dat", LanCan: 1, NguoiCan: "BS. Mai", ThoiGian: now, CreatedAt: now},
		{ID: "wi-003", GiaiID: gID, VdvID: "vdv-003", VdvTen: "Lê Thị Phương", DoanID: "doan-bd", DoanTen: "Bình Định", HangCan: "48kg Nữ", CanNang: 49.1, GioiHan: 48, SaiSo: 0.5, KetQua: "khong_dat", LanCan: 1, NguoiCan: "BS. Hùng", ThoiGian: now, CreatedAt: now},
		{ID: "wi-004", GiaiID: gID, VdvID: "vdv-004", VdvTen: "Phạm Văn Đạt", DoanID: "doan-dn", DoanTen: "Đà Nẵng", HangCan: "68kg Nam", CanNang: 67.5, GioiHan: 68, SaiSo: 0.5, KetQua: "dat", LanCan: 1, NguoiCan: "BS. Mai", ThoiGian: now, CreatedAt: now},
		{ID: "wi-005", GiaiID: gID, VdvID: "vdv-005", VdvTen: "Huỳnh Thanh Tâm", DoanID: "doan-hcm", DoanTen: "TP.HCM", HangCan: "78kg Nam", CanNang: 77.0, GioiHan: 78, SaiSo: 0.5, KetQua: "dat", LanCan: 1, NguoiCan: "BS. Hùng", ThoiGian: now, CreatedAt: now},
		{ID: "wi-006", GiaiID: gID, VdvID: "vdv-006", VdvTen: "Ngô Minh Hải", DoanID: "doan-hn", DoanTen: "Hà Nội", HangCan: "54kg Nam", CanNang: 54.8, GioiHan: 54, SaiSo: 0.5, KetQua: "khong_dat", LanCan: 1, NguoiCan: "BS. Mai", ThoiGian: now, CreatedAt: now},
	}

	// ── Referee Assignments ─────────────────────────────────
	s.assignments = []RefereeAssignment{
		{ID: "ra-001", GiaiID: gID, TrongTaiID: "tt-001", TrongTaiTen: "Võ Thanh Sơn", CapBac: "quoc_gia", ChuyenMon: "doi_khang", SanID: "san-01", SanTen: "Sàn 1", Ngay: "2025-11-15", Phien: "sang", VaiTro: "chu_toa", TrangThai: "xac_nhan", CreatedAt: now},
		{ID: "ra-002", GiaiID: gID, TrongTaiID: "tt-002", TrongTaiTen: "Đặng Quốc Việt", CapBac: "quoc_gia", ChuyenMon: "doi_khang", SanID: "san-01", SanTen: "Sàn 1", Ngay: "2025-11-15", Phien: "sang", VaiTro: "giam_dinh", TrangThai: "xac_nhan", CreatedAt: now},
		{ID: "ra-003", GiaiID: gID, TrongTaiID: "tt-003", TrongTaiTen: "Trần Đức Lộc", CapBac: "cap_1", ChuyenMon: "doi_khang", SanID: "san-02", SanTen: "Sàn 2", Ngay: "2025-11-15", Phien: "sang", VaiTro: "chu_toa", TrangThai: "phan_cong", CreatedAt: now},
		{ID: "ra-004", GiaiID: gID, TrongTaiID: "tt-004", TrongTaiTen: "Lê Minh Quang", CapBac: "cap_1", ChuyenMon: "quyen", SanID: "san-03", SanTen: "Sàn 3", Ngay: "2025-11-15", Phien: "chieu", VaiTro: "diem", TrangThai: "phan_cong", CreatedAt: now},
		{ID: "ra-005", GiaiID: gID, TrongTaiID: "tt-005", TrongTaiTen: "Nguyễn Thị Lan", CapBac: "cap_2", ChuyenMon: "quyen", SanID: "san-03", SanTen: "Sàn 3", Ngay: "2025-11-15", Phien: "chieu", VaiTro: "diem", TrangThai: "phan_cong", CreatedAt: now},
	}

	// ── Team Results ────────────────────────────────────────
	s.teamResults = []TeamResult{
		{ID: "tr-001", GiaiID: gID, DoanID: "doan-hcm", DoanTen: "TP. Hồ Chí Minh", Tinh: "TP.HCM", HCV: 5, HCB: 3, HCD: 4, TongHC: 12, Diem: 25, XepHang: 1},
		{ID: "tr-002", GiaiID: gID, DoanID: "doan-hn", DoanTen: "Hà Nội", Tinh: "Hà Nội", HCV: 4, HCB: 4, HCD: 3, TongHC: 11, Diem: 23, XepHang: 2},
		{ID: "tr-003", GiaiID: gID, DoanID: "doan-bd", DoanTen: "Bình Định", Tinh: "Bình Định", HCV: 3, HCB: 5, HCD: 6, TongHC: 14, Diem: 22, XepHang: 3},
		{ID: "tr-004", GiaiID: gID, DoanID: "doan-dn", DoanTen: "Đà Nẵng", Tinh: "Đà Nẵng", HCV: 3, HCB: 2, HCD: 3, TongHC: 8, Diem: 16, XepHang: 4},
		{ID: "tr-005", GiaiID: gID, DoanID: "doan-tth", DoanTen: "Thừa Thiên Huế", Tinh: "TT-Huế", HCV: 2, HCB: 3, HCD: 2, TongHC: 7, Diem: 14, XepHang: 5},
		{ID: "tr-006", GiaiID: gID, DoanID: "doan-ag", DoanTen: "An Giang", Tinh: "An Giang", HCV: 1, HCB: 2, HCD: 5, TongHC: 8, Diem: 11, XepHang: 6},
	}

	// ── Finance ─────────────────────────────────────────────
	s.finance = []FinanceEntry{
		{ID: "fi-001", GiaiID: gID, Loai: "thu", DanhMuc: "le_phi_doan", MoTa: "Lệ phí đoàn TP.HCM", SoTien: 5000000, DoanID: "doan-hcm", DoanTen: "TP.HCM", TrangThai: "da_thu", NgayGD: "2025-11-10", CreatedBy: "btc-002", CreatedAt: now},
		{ID: "fi-002", GiaiID: gID, Loai: "thu", DanhMuc: "le_phi_doan", MoTa: "Lệ phí đoàn Hà Nội", SoTien: 5000000, DoanID: "doan-hn", DoanTen: "Hà Nội", TrangThai: "da_thu", NgayGD: "2025-11-10", CreatedBy: "btc-002", CreatedAt: now},
		{ID: "fi-003", GiaiID: gID, Loai: "thu", DanhMuc: "le_phi_vdv", MoTa: "Lệ phí 120 VĐV x 200.000đ", SoTien: 24000000, TrangThai: "da_thu", NgayGD: "2025-11-12", CreatedBy: "btc-002", CreatedAt: now},
		{ID: "fi-004", GiaiID: gID, Loai: "thu", DanhMuc: "tai_tro", MoTa: "Tài trợ Công ty ABC", SoTien: 50000000, TrangThai: "da_thu", NgayGD: "2025-11-01", CreatedBy: "btc-001", CreatedAt: now},
		{ID: "fi-005", GiaiID: gID, Loai: "chi", DanhMuc: "thue_san", MoTa: "Thuê nhà thi đấu 3 ngày", SoTien: 30000000, TrangThai: "da_chi", NgayGD: "2025-11-14", CreatedBy: "btc-002", CreatedAt: now},
		{ID: "fi-006", GiaiID: gID, Loai: "chi", DanhMuc: "phu_cap_tt", MoTa: "Phụ cấp 15 trọng tài x 500.000đ", SoTien: 7500000, TrangThai: "da_chi", NgayGD: "2025-11-17", CreatedBy: "btc-002", CreatedAt: now},
		{ID: "fi-007", GiaiID: gID, Loai: "chi", DanhMuc: "huy_chuong", MoTa: "Huy chương + cúp", SoTien: 8000000, TrangThai: "da_chi", NgayGD: "2025-11-13", CreatedBy: "btc-002", CreatedAt: now},
	}

	// ── Content Results (Kết quả từng nội dung) ────────────
	s.contResults = []ContentResult{
		{ID: "cr-001", GiaiID: gID, NoiDungID: "nd-dk-54", NoiDungTen: "ĐK Nam 54kg", HangCan: "54kg", LuaTuoi: "Tuyển", VdvIDNhat: "vdv-001", VdvTenNhat: "Nguyễn Anh Tuấn", DoanNhat: "TP.HCM", VdvIDNhi: "vdv-006", VdvTenNhi: "Ngô Minh Hải", DoanNhi: "Hà Nội", VdvIDBa1: "vdv-010", VdvTenBa1: "Trần Văn Sơn", DoanBa1: "Bình Định", VdvIDBa2: "vdv-011", VdvTenBa2: "Lê Hoàng Phúc", DoanBa2: "Đà Nẵng"},
		{ID: "cr-002", GiaiID: gID, NoiDungID: "nd-dk-60", NoiDungTen: "ĐK Nam 60kg", HangCan: "60kg", LuaTuoi: "Tuyển", VdvIDNhat: "vdv-002", VdvTenNhat: "Trần Minh Phúc", DoanNhat: "Hà Nội", VdvIDNhi: "vdv-012", VdvTenNhi: "Phạm Đức Huy", DoanNhi: "TP.HCM", VdvIDBa1: "vdv-013", VdvTenBa1: "Võ Thanh Tùng", DoanBa1: "Bình Định", VdvIDBa2: "vdv-014", VdvTenBa2: "Huỳnh Minh Đạt", DoanBa2: "An Giang"},
		{ID: "cr-003", GiaiID: gID, NoiDungID: "nd-dk-48-nu", NoiDungTen: "ĐK Nữ 48kg", HangCan: "48kg", LuaTuoi: "Tuyển", VdvIDNhat: "vdv-020", VdvTenNhat: "Nguyễn Thị Hương", DoanNhat: "Bình Định", VdvIDNhi: "vdv-003", VdvTenNhi: "Lê Thị Phương", DoanNhi: "Bình Định", VdvIDBa1: "vdv-021", VdvTenBa1: "Trần Minh Châu", DoanBa1: "TP.HCM", VdvIDBa2: "vdv-022", VdvTenBa2: "Phan Thị Lan", DoanBa2: "Đà Nẵng"},
		{ID: "cr-004", GiaiID: gID, NoiDungID: "nd-quyen-01", NoiDungTen: "Quyền thuật Nam", LuaTuoi: "Tuyển", VdvIDNhat: "vdv-030", VdvTenNhat: "Lê Văn Tài", DoanNhat: "TP.HCM", VdvIDNhi: "vdv-031", VdvTenNhi: "Đặng Quốc Bảo", DoanNhi: "TT-Huế", VdvIDBa1: "vdv-032", VdvTenBa1: "Hoàng Minh Trí", DoanBa1: "Hà Nội", VdvIDBa2: "vdv-033", VdvTenBa2: "Nguyễn Đức Thịnh", DoanBa2: "An Giang"},
	}

	// ── Technical Meetings ──────────────────────────────────
	s.meetings = []TechnicalMeeting{
		{ID: "tm-001", GiaiID: gID, TieuDe: "Họp chuyên môn lần 1 — Quy chế giải", Ngay: "2025-11-13", DiaDiem: "Phòng họp A — Nhà thi đấu Phan Đình Phùng", ChuTri: "Lê Hoàng Nam", ThamDu: []string{"Ban CM", "Ban TT", "Trưởng đoàn"}, NoiDung: "Thống nhất quy chế thi đấu, hạng cân, nội dung thi đấu", QuyetDinh: []string{"Áp dụng luật VCTQG 2024", "Mỗi VĐV tối đa 2 nội dung đối kháng"}, TrangThai: "hoan_thanh", CreatedAt: now},
		{ID: "tm-002", GiaiID: gID, TieuDe: "Họp chuyên môn lần 2 — Bốc thăm & Lịch đấu", Ngay: "2025-11-14", DiaDiem: "Phòng họp A — Nhà thi đấu Phan Đình Phùng", ChuTri: "Lê Hoàng Nam", ThamDu: []string{"Ban CM", "Ban TT", "Trưởng đoàn"}, NoiDung: "Bốc thăm xếp nhánh đối kháng, lịch thi quyền", QuyetDinh: []string{"Hoàn thành bốc thăm 12 nội dung", "Lịch đấu 3 ngày"}, TrangThai: "hoan_thanh", CreatedAt: now},
	}

	// ── Protests ────────────────────────────────────────────
	s.protests = []Protest{
		{ID: "pr-001", GiaiID: gID, TranID: "tran-015", TranMoTa: "ĐK Nam 60kg — Trận tứ kết 3", NguoiNop: "HLV Trần Đức", DoanTen: "Hà Nội", LoaiKN: "cham_diem", LyDo: "Trọng tài không ghi nhận đòn đá vòng cầu hợp lệ ở hiệp 2", TrangThai: "xem_xet", HasVideo: true, NgayNop: now, CreatedAt: now},
		{ID: "pr-002", GiaiID: gID, TranID: "tran-022", TranMoTa: "ĐK Nữ 48kg — Trận bán kết 1", NguoiNop: "HLV Võ Minh", DoanTen: "Bình Định", LoaiKN: "pham_luat", LyDo: "Đối thủ sử dụng kỹ thuật không hợp lệ (chỏ)", TrangThai: "chap_nhan", QuyetDinh: "Xử phạt VĐV vi phạm, trừ 2 điểm", NguoiXL: "Hoàng Minh Tuấn", NgayNop: now, CreatedAt: now},
		{ID: "pr-003", GiaiID: gID, TranID: "tran-031", TranMoTa: "ĐK Nam 68kg — Trận bán kết 2", NguoiNop: "Trưởng đoàn Lê Hải", DoanTen: "Đà Nẵng", LoaiKN: "can_luong", LyDo: "Yêu cầu cân lại đối thủ — nghi vượt cân", TrangThai: "bac_bo", QuyetDinh: "VĐV đã cân đạt trước giờ thi đấu quy định", NguoiXL: "Hoàng Minh Tuấn", NgayNop: now, CreatedAt: now},
	}
}
