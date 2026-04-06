package adapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"vct-platform/backend/internal/domain/btc"
)

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — BTC POSTGRESQL ADAPTER
// Implements btc.Store interface using PostgreSQL.
// ═══════════════════════════════════════════════════════════════

type PgBTCStore struct {
	db *sql.DB
}

func NewPgBTCStore(db *sql.DB) *PgBTCStore {
	return &PgBTCStore{db: db}
}

// ── BTC Members ─────────────────────────────────────────────

func (s *PgBTCStore) ListMembers(ctx context.Context, giaiID string) ([]btc.BTCMember, error) {
	query := `SELECT id, ten, chuc_vu, ban, cap, sdt, email, don_vi, giai_id, is_active FROM btc_members`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY cap ASC, ten ASC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list members: %w", err)
	}
	defer rows.Close()

	var out []btc.BTCMember
	for rows.Next() {
		var m btc.BTCMember
		if err := rows.Scan(&m.ID, &m.Ten, &m.ChucVu, &m.Ban, &m.Cap, &m.Sdt, &m.Email, &m.DonVi, &m.GiaiID, &m.IsActive); err != nil {
			return nil, fmt.Errorf("btc scan member: %w", err)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) GetMember(ctx context.Context, id string) (*btc.BTCMember, error) {
	var m btc.BTCMember
	err := s.db.QueryRowContext(ctx,
		`SELECT id, ten, chuc_vu, ban, cap, sdt, email, don_vi, giai_id, is_active FROM btc_members WHERE id = $1`, id,
	).Scan(&m.ID, &m.Ten, &m.ChucVu, &m.Ban, &m.Cap, &m.Sdt, &m.Email, &m.DonVi, &m.GiaiID, &m.IsActive)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("member not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("btc get member: %w", err)
	}
	return &m, nil
}

func (s *PgBTCStore) CreateMember(ctx context.Context, m *btc.BTCMember) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_members (id, ten, chuc_vu, ban, cap, sdt, email, don_vi, giai_id, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		m.ID, m.Ten, m.ChucVu, m.Ban, m.Cap, m.Sdt, m.Email, m.DonVi, m.GiaiID, m.IsActive,
	)
	return err
}

func (s *PgBTCStore) UpdateMember(ctx context.Context, m *btc.BTCMember) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE btc_members SET ten=$2, chuc_vu=$3, ban=$4, cap=$5, sdt=$6, email=$7, don_vi=$8, giai_id=$9, is_active=$10 WHERE id=$1`,
		m.ID, m.Ten, m.ChucVu, m.Ban, m.Cap, m.Sdt, m.Email, m.DonVi, m.GiaiID, m.IsActive,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("member not found: %s", m.ID)
	}
	return nil
}

func (s *PgBTCStore) DeleteMember(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM btc_members WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("member not found: %s", id)
	}
	return nil
}

// ── Weigh-In ────────────────────────────────────────────────

func (s *PgBTCStore) ListWeighIns(ctx context.Context, giaiID string) ([]btc.WeighInRecord, error) {
	query := `SELECT id, giai_id, vdv_id, vdv_ten, doan_id, doan_ten, hang_can, can_nang, gioi_han, sai_so, ket_qua, lan_can, ghi_chu, nguoi_can, thoi_gian, created_at FROM btc_weigh_ins`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list weigh-ins: %w", err)
	}
	defer rows.Close()

	var out []btc.WeighInRecord
	for rows.Next() {
		var w btc.WeighInRecord
		if err := rows.Scan(&w.ID, &w.GiaiID, &w.VdvID, &w.VdvTen, &w.DoanID, &w.DoanTen, &w.HangCan, &w.CanNang, &w.GioiHan, &w.SaiSo, &w.KetQua, &w.LanCan, &w.GhiChu, &w.NguoiCan, &w.ThoiGian, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("btc scan weigh-in: %w", err)
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) CreateWeighIn(ctx context.Context, w *btc.WeighInRecord) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_weigh_ins (id, giai_id, vdv_id, vdv_ten, doan_id, doan_ten, hang_can, can_nang, gioi_han, sai_so, ket_qua, lan_can, ghi_chu, nguoi_can, thoi_gian, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		w.ID, w.GiaiID, w.VdvID, w.VdvTen, w.DoanID, w.DoanTen, w.HangCan, w.CanNang, w.GioiHan, w.SaiSo, w.KetQua, w.LanCan, w.GhiChu, w.NguoiCan, w.ThoiGian, w.CreatedAt,
	)
	return err
}

// ── Draw ────────────────────────────────────────────────────

func (s *PgBTCStore) ListDraws(ctx context.Context, giaiID string) ([]btc.DrawResult, error) {
	query := `SELECT id, giai_id, noi_dung_id, noi_dung_ten, loai_nd, hang_can, lua_tuoi, so_vdv, nhanh, thu_tu, created_at, created_by FROM btc_draws`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list draws: %w", err)
	}
	defer rows.Close()

	var out []btc.DrawResult
	for rows.Next() {
		var d btc.DrawResult
		var nhanhJSON, thuTuJSON []byte
		if err := rows.Scan(&d.ID, &d.GiaiID, &d.NoiDungID, &d.NoiDungTen, &d.LoaiND, &d.HangCan, &d.LuaTuoi, &d.SoVDV, &nhanhJSON, &thuTuJSON, &d.CreatedAt, &d.CreatedBy); err != nil {
			return nil, fmt.Errorf("btc scan draw: %w", err)
		}
		_ = json.Unmarshal(nhanhJSON, &d.Nhanh)
		_ = json.Unmarshal(thuTuJSON, &d.ThuTu)
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) CreateDraw(ctx context.Context, d *btc.DrawResult) error {
	nhanhJSON, _ := json.Marshal(d.Nhanh)
	thuTuJSON, _ := json.Marshal(d.ThuTu)
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_draws (id, giai_id, noi_dung_id, noi_dung_ten, loai_nd, hang_can, lua_tuoi, so_vdv, nhanh, thu_tu, created_at, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		d.ID, d.GiaiID, d.NoiDungID, d.NoiDungTen, d.LoaiND, d.HangCan, d.LuaTuoi, d.SoVDV, nhanhJSON, thuTuJSON, d.CreatedAt, d.CreatedBy,
	)
	return err
}

// ── Referee Assignment ──────────────────────────────────────

func (s *PgBTCStore) ListAssignments(ctx context.Context, giaiID string) ([]btc.RefereeAssignment, error) {
	query := `SELECT id, giai_id, trong_tai_id, trong_tai_ten, cap_bac, chuyen_mon, san_id, san_ten, ngay, phien, vai_tro, trang_thai, ghi_chu, created_at FROM btc_assignments`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY ngay, phien`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list assignments: %w", err)
	}
	defer rows.Close()

	var out []btc.RefereeAssignment
	for rows.Next() {
		var a btc.RefereeAssignment
		if err := rows.Scan(&a.ID, &a.GiaiID, &a.TrongTaiID, &a.TrongTaiTen, &a.CapBac, &a.ChuyenMon, &a.SanID, &a.SanTen, &a.Ngay, &a.Phien, &a.VaiTro, &a.TrangThai, &a.GhiChu, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("btc scan assignment: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) CreateAssignment(ctx context.Context, a *btc.RefereeAssignment) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_assignments (id, giai_id, trong_tai_id, trong_tai_ten, cap_bac, chuyen_mon, san_id, san_ten, ngay, phien, vai_tro, trang_thai, ghi_chu, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		a.ID, a.GiaiID, a.TrongTaiID, a.TrongTaiTen, a.CapBac, a.ChuyenMon, a.SanID, a.SanTen, a.Ngay, a.Phien, a.VaiTro, a.TrangThai, a.GhiChu, a.CreatedAt,
	)
	return err
}

// ── Results ─────────────────────────────────────────────────

func (s *PgBTCStore) ListTeamResults(ctx context.Context, giaiID string) ([]btc.TeamResult, error) {
	query := `SELECT id, giai_id, doan_id, doan_ten, tinh, hcv, hcb, hcd, tong_hc, diem, xep_hang FROM btc_team_results`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY xep_hang ASC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list team results: %w", err)
	}
	defer rows.Close()

	var out []btc.TeamResult
	for rows.Next() {
		var r btc.TeamResult
		if err := rows.Scan(&r.ID, &r.GiaiID, &r.DoanID, &r.DoanTen, &r.Tinh, &r.HCV, &r.HCB, &r.HCD, &r.TongHC, &r.Diem, &r.XepHang); err != nil {
			return nil, fmt.Errorf("btc scan team result: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) CreateTeamResult(ctx context.Context, r *btc.TeamResult) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_team_results (id, giai_id, doan_id, doan_ten, tinh, hcv, hcb, hcd, tong_hc, diem, xep_hang) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		r.ID, r.GiaiID, r.DoanID, r.DoanTen, r.Tinh, r.HCV, r.HCB, r.HCD, r.TongHC, r.Diem, r.XepHang,
	)
	return err
}

func (s *PgBTCStore) ListContentResults(ctx context.Context, giaiID string) ([]btc.ContentResult, error) {
	query := `SELECT id, giai_id, noi_dung_id, noi_dung_ten, hang_can, lua_tuoi, vdv_id_nhat, vdv_ten_nhat, doan_nhat, vdv_id_nhi, vdv_ten_nhi, doan_nhi, vdv_id_ba_1, vdv_ten_ba_1, doan_ba_1, vdv_id_ba_2, vdv_ten_ba_2, doan_ba_2 FROM btc_content_results`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list content results: %w", err)
	}
	defer rows.Close()

	var out []btc.ContentResult
	for rows.Next() {
		var r btc.ContentResult
		if err := rows.Scan(&r.ID, &r.GiaiID, &r.NoiDungID, &r.NoiDungTen, &r.HangCan, &r.LuaTuoi, &r.VdvIDNhat, &r.VdvTenNhat, &r.DoanNhat, &r.VdvIDNhi, &r.VdvTenNhi, &r.DoanNhi, &r.VdvIDBa1, &r.VdvTenBa1, &r.DoanBa1, &r.VdvIDBa2, &r.VdvTenBa2, &r.DoanBa2); err != nil {
			return nil, fmt.Errorf("btc scan content result: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ── Finance ─────────────────────────────────────────────────

func (s *PgBTCStore) ListFinance(ctx context.Context, giaiID string) ([]btc.FinanceEntry, error) {
	query := `SELECT id, giai_id, loai, danh_muc, mo_ta, so_tien, doan_id, doan_ten, trang_thai, ngay_gd, ghi_chu, created_by, created_at FROM btc_finance`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list finance: %w", err)
	}
	defer rows.Close()

	var out []btc.FinanceEntry
	for rows.Next() {
		var f btc.FinanceEntry
		if err := rows.Scan(&f.ID, &f.GiaiID, &f.Loai, &f.DanhMuc, &f.MoTa, &f.SoTien, &f.DoanID, &f.DoanTen, &f.TrangThai, &f.NgayGD, &f.GhiChu, &f.CreatedBy, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("btc scan finance: %w", err)
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) CreateFinance(ctx context.Context, f *btc.FinanceEntry) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_finance (id, giai_id, loai, danh_muc, mo_ta, so_tien, doan_id, doan_ten, trang_thai, ngay_gd, ghi_chu, created_by, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		f.ID, f.GiaiID, f.Loai, f.DanhMuc, f.MoTa, f.SoTien, f.DoanID, f.DoanTen, f.TrangThai, f.NgayGD, f.GhiChu, f.CreatedBy, f.CreatedAt,
	)
	return err
}

func (s *PgBTCStore) UpdateFinance(ctx context.Context, f *btc.FinanceEntry) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE btc_finance SET loai=$2, danh_muc=$3, mo_ta=$4, so_tien=$5, doan_id=$6, doan_ten=$7, trang_thai=$8, ngay_gd=$9, ghi_chu=$10, created_by=$11 WHERE id=$1`,
		f.ID, f.Loai, f.DanhMuc, f.MoTa, f.SoTien, f.DoanID, f.DoanTen, f.TrangThai, f.NgayGD, f.GhiChu, f.CreatedBy,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("finance entry not found: %s", f.ID)
	}
	return nil
}

// ── Technical Meeting ───────────────────────────────────────

func (s *PgBTCStore) ListMeetings(ctx context.Context, giaiID string) ([]btc.TechnicalMeeting, error) {
	query := `SELECT id, giai_id, tieu_de, ngay, dia_diem, chu_tri, tham_du, noi_dung, quyet_dinh, bien_ban_file, trang_thai, created_at FROM btc_meetings`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY ngay DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list meetings: %w", err)
	}
	defer rows.Close()

	var out []btc.TechnicalMeeting
	for rows.Next() {
		var m btc.TechnicalMeeting
		var thamDuJSON, quyetDinhJSON []byte
		if err := rows.Scan(&m.ID, &m.GiaiID, &m.TieuDe, &m.Ngay, &m.DiaDiem, &m.ChuTri, &thamDuJSON, &m.NoiDung, &quyetDinhJSON, &m.BienBanFile, &m.TrangThai, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("btc scan meeting: %w", err)
		}
		_ = json.Unmarshal(thamDuJSON, &m.ThamDu)
		_ = json.Unmarshal(quyetDinhJSON, &m.QuyetDinh)
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) CreateMeeting(ctx context.Context, m *btc.TechnicalMeeting) error {
	thamDuJSON, _ := json.Marshal(m.ThamDu)
	quyetDinhJSON, _ := json.Marshal(m.QuyetDinh)
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_meetings (id, giai_id, tieu_de, ngay, dia_diem, chu_tri, tham_du, noi_dung, quyet_dinh, bien_ban_file, trang_thai, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		m.ID, m.GiaiID, m.TieuDe, m.Ngay, m.DiaDiem, m.ChuTri, thamDuJSON, m.NoiDung, quyetDinhJSON, m.BienBanFile, m.TrangThai, m.CreatedAt,
	)
	return err
}

// ── Protests ────────────────────────────────────────────────

func (s *PgBTCStore) ListProtests(ctx context.Context, giaiID string) ([]btc.Protest, error) {
	query := `SELECT id, giai_id, tran_id, tran_mo_ta, nguoi_nop, doan_ten, loai_kn, ly_do, trang_thai, has_video, quyet_dinh, nguoi_xl, ngay_nop, ngay_xl, created_at FROM btc_protests`
	args := []any{}
	if giaiID != "" {
		query += ` WHERE giai_id = $1`
		args = append(args, giaiID)
	}
	query += ` ORDER BY ngay_nop DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("btc list protests: %w", err)
	}
	defer rows.Close()

	var out []btc.Protest
	for rows.Next() {
		var p btc.Protest
		var ngayXL sql.NullTime
		if err := rows.Scan(&p.ID, &p.GiaiID, &p.TranID, &p.TranMoTa, &p.NguoiNop, &p.DoanTen, &p.LoaiKN, &p.LyDo, &p.TrangThai, &p.HasVideo, &p.QuyetDinh, &p.NguoiXL, &p.NgayNop, &ngayXL, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("btc scan protest: %w", err)
		}
		if ngayXL.Valid {
			p.NgayXL = &ngayXL.Time
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *PgBTCStore) GetProtest(ctx context.Context, id string) (*btc.Protest, error) {
	var p btc.Protest
	var ngayXL sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, giai_id, tran_id, tran_mo_ta, nguoi_nop, doan_ten, loai_kn, ly_do, trang_thai, has_video, quyet_dinh, nguoi_xl, ngay_nop, ngay_xl, created_at FROM btc_protests WHERE id = $1`, id,
	).Scan(&p.ID, &p.GiaiID, &p.TranID, &p.TranMoTa, &p.NguoiNop, &p.DoanTen, &p.LoaiKN, &p.LyDo, &p.TrangThai, &p.HasVideo, &p.QuyetDinh, &p.NguoiXL, &p.NgayNop, &ngayXL, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("protest not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("btc get protest: %w", err)
	}
	if ngayXL.Valid {
		p.NgayXL = &ngayXL.Time
	}
	return &p, nil
}

func (s *PgBTCStore) CreateProtest(ctx context.Context, p *btc.Protest) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO btc_protests (id, giai_id, tran_id, tran_mo_ta, nguoi_nop, doan_ten, loai_kn, ly_do, trang_thai, has_video, quyet_dinh, nguoi_xl, ngay_nop, ngay_xl, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		p.ID, p.GiaiID, p.TranID, p.TranMoTa, p.NguoiNop, p.DoanTen, p.LoaiKN, p.LyDo, p.TrangThai, p.HasVideo, p.QuyetDinh, p.NguoiXL, p.NgayNop, p.NgayXL, p.CreatedAt,
	)
	return err
}

func (s *PgBTCStore) UpdateProtest(ctx context.Context, p *btc.Protest) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE btc_protests SET trang_thai=$2, quyet_dinh=$3, nguoi_xl=$4, ngay_xl=$5 WHERE id=$1`,
		p.ID, p.TrangThai, p.QuyetDinh, p.NguoiXL, p.NgayXL,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("protest not found: %s", p.ID)
	}
	return nil
}

// Compile-time interface check
var _ btc.Store = (*PgBTCStore)(nil)

// unused import guard
var _ = time.Now
