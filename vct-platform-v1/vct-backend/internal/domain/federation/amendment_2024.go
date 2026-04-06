package federation

import "time"

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — AMENDMENT 2024
// Luật số 128/2024/LĐVTCTVN
// Sửa đổi, bổ sung một số điều của Luật thi đấu Võ Cổ Truyền
//
// File này CHỈ chứa các điều khoản THAY ĐỔI so với 2021-QC01.
// Dữ liệu KHÔNG thay đổi vẫn giữ nguyên trong regulation_config.go.
// CẤM chỉnh sửa nếu không có văn bản pháp lý.
// ═══════════════════════════════════════════════════════════════

// ── Metadata Luật sửa đổi ────────────────────────────────────

const (
	AmendmentCode       = "128/2024/LĐVTCTVN"
	AmendmentName       = "Luật sửa đổi, bổ sung một số điều của Luật thi đấu Võ cổ truyền"
	AmendmentYear       = 2024
	AmendmentDate       = "2024-07-20"
	AmendmentSignedBy   = "Nguyễn Ngọc Anh"
	AmendmentSignedRole = "Chủ tịch Liên đoàn VTCT Việt Nam"
)

// ────────────────────────────────────────────────────────────────
// NHÓM TUỔI — SỬA ĐỔI ĐIỀU 4 & ĐIỀU 25
// Quan trọng: Nhóm tuổi đối kháng ≠ quyền thuật
// ────────────────────────────────────────────────────────────────

// CompetitionCategory phân biệt loại giải (đối kháng vs quyền thuật)
type CompetitionCategory string

const (
	CategoryDoiKhang   CompetitionCategory = "doi_khang"   // Đối kháng
	CategoryQuyenThuat CompetitionCategory = "quyen_thuat" // Quyền thuật
)

// AgeGroupAmended mở rộng MasterAgeGroup để hỗ trợ amendment 2024:
// phân biệt nhóm tuổi theo loại giải và cấp giải.
type AgeGroupAmended struct {
	MasterAgeGroup
	Category  CompetitionCategory `json:"category"`   // doi_khang hoặc quyen_thuat
	TierCode  string              `json:"tier_code"`  // "championship" hoặc "youth"
	AmendedBy string              `json:"amended_by"` // Mã văn bản sửa đổi
}

// Amendment2024AgeGroups trả về nhóm tuổi đã sửa đổi theo Luật 128/2024.
// THAY THẾ NationalAgeGroups() (2021) cho mục đích seed data.
func Amendment2024AgeGroups() []AgeGroupAmended {
	now := time.Now()
	scope := BeltScopeNational
	amd := AmendmentCode

	return []AgeGroupAmended{
		// ═══════════════════════════════════════════════════
		// ĐỐI KHÁNG — Điều 4 sửa đổi
		// ═══════════════════════════════════════════════════
		{
			MasterAgeGroup: MasterAgeGroup{ID: "dk-senior", Name: "Vô địch (17-40)", MinAge: 17, MaxAge: 40, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryDoiKhang, TierCode: "championship", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "dk-u13", Name: "Trẻ 12-13", MinAge: 12, MaxAge: 13, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryDoiKhang, TierCode: "youth", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "dk-u15", Name: "Trẻ 14-15", MinAge: 14, MaxAge: 15, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryDoiKhang, TierCode: "youth", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "dk-u17", Name: "Trẻ 16-17", MinAge: 16, MaxAge: 17, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryDoiKhang, TierCode: "youth", AmendedBy: amd,
		},

		// ═══════════════════════════════════════════════════
		// QUYỀN THUẬT — Điều 25 sửa đổi
		// ═══════════════════════════════════════════════════
		{
			MasterAgeGroup: MasterAgeGroup{ID: "qt-senior", Name: "VĐ 17-40", MinAge: 17, MaxAge: 40, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryQuyenThuat, TierCode: "championship", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "qt-masters-a", Name: "VĐ 41-50", MinAge: 41, MaxAge: 50, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryQuyenThuat, TierCode: "championship", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "qt-masters-b", Name: "VĐ 51-60", MinAge: 51, MaxAge: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryQuyenThuat, TierCode: "championship", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "qt-u10", Name: "Trẻ 6-10", MinAge: 6, MaxAge: 10, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryQuyenThuat, TierCode: "youth", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "qt-u14", Name: "Trẻ 11-14", MinAge: 11, MaxAge: 14, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryQuyenThuat, TierCode: "youth", AmendedBy: amd,
		},
		{
			MasterAgeGroup: MasterAgeGroup{ID: "qt-u17", Name: "Trẻ 15-17", MinAge: 15, MaxAge: 17, Scope: scope, CreatedAt: now, UpdatedAt: now},
			Category:       CategoryQuyenThuat, TierCode: "youth", AmendedBy: amd,
		},
	}
}

// ────────────────────────────────────────────────────────────────
// HẠNG CÂN — SỬA ĐỔI ĐIỀU 4
// Hoàn toàn thay thế hạng cân 2021
// ────────────────────────────────────────────────────────────────

// Amendment2024WeightClasses trả về hạng cân mới theo Luật 128/2024.
// THAY THẾ HOÀN TOÀN NationalWeightClasses() (2021).
func Amendment2024WeightClasses() []MasterWeightClass {
	now := time.Now()
	scope := BeltScopeNational

	return []MasterWeightClass{
		// ═════════════════════════════════════════════════════
		// NAM — VÔ ĐỊCH (17-40) — 14 hạng + 1 mở = 15
		// ═════════════════════════════════════════════════════
		{ID: "m-sr-u48", Gender: "MALE", Category: "dk-senior", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u51", Gender: "MALE", Category: "dk-senior", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u54", Gender: "MALE", Category: "dk-senior", MinWeight: 51.1, MaxWeight: 54, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u57", Gender: "MALE", Category: "dk-senior", MinWeight: 54.1, MaxWeight: 57, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u60", Gender: "MALE", Category: "dk-senior", MinWeight: 57.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u64", Gender: "MALE", Category: "dk-senior", MinWeight: 60.1, MaxWeight: 64, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u68", Gender: "MALE", Category: "dk-senior", MinWeight: 64.1, MaxWeight: 68, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u72", Gender: "MALE", Category: "dk-senior", MinWeight: 68.1, MaxWeight: 72, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u76", Gender: "MALE", Category: "dk-senior", MinWeight: 72.1, MaxWeight: 76, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u80", Gender: "MALE", Category: "dk-senior", MinWeight: 76.1, MaxWeight: 80, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u85", Gender: "MALE", Category: "dk-senior", MinWeight: 80.1, MaxWeight: 85, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u90", Gender: "MALE", Category: "dk-senior", MinWeight: 85.1, MaxWeight: 90, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u95", Gender: "MALE", Category: "dk-senior", MinWeight: 90.1, MaxWeight: 95, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u110", Gender: "MALE", Category: "dk-senior", MinWeight: 95.1, MaxWeight: 110, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-o110", Gender: "MALE", Category: "dk-senior", MinWeight: 110.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NỮ — VÔ ĐỊCH (17-40) — 11 hạng + 1 mở = 12
		// ═════════════════════════════════════════════════════
		{ID: "f-sr-u45", Gender: "FEMALE", Category: "dk-senior", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u48", Gender: "FEMALE", Category: "dk-senior", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u51", Gender: "FEMALE", Category: "dk-senior", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u54", Gender: "FEMALE", Category: "dk-senior", MinWeight: 51.1, MaxWeight: 54, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u57", Gender: "FEMALE", Category: "dk-senior", MinWeight: 54.1, MaxWeight: 57, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u60", Gender: "FEMALE", Category: "dk-senior", MinWeight: 57.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u64", Gender: "FEMALE", Category: "dk-senior", MinWeight: 60.1, MaxWeight: 64, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u68", Gender: "FEMALE", Category: "dk-senior", MinWeight: 64.1, MaxWeight: 68, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u72", Gender: "FEMALE", Category: "dk-senior", MinWeight: 68.1, MaxWeight: 72, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u76", Gender: "FEMALE", Category: "dk-senior", MinWeight: 72.1, MaxWeight: 76, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u80", Gender: "FEMALE", Category: "dk-senior", MinWeight: 76.1, MaxWeight: 80, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-o80", Gender: "FEMALE", Category: "dk-senior", MinWeight: 80.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NAM — TRẺ 12-13 (U13) — 9 hạng + 1 mở = 10
		// ═════════════════════════════════════════════════════
		{ID: "m-u13-u38", Gender: "MALE", Category: "dk-u13", MinWeight: 36.1, MaxWeight: 38, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u40", Gender: "MALE", Category: "dk-u13", MinWeight: 38.1, MaxWeight: 40, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u42", Gender: "MALE", Category: "dk-u13", MinWeight: 40.1, MaxWeight: 42, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u45", Gender: "MALE", Category: "dk-u13", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u48", Gender: "MALE", Category: "dk-u13", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u51", Gender: "MALE", Category: "dk-u13", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u54", Gender: "MALE", Category: "dk-u13", MinWeight: 51.1, MaxWeight: 54, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u57", Gender: "MALE", Category: "dk-u13", MinWeight: 54.1, MaxWeight: 57, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-u60", Gender: "MALE", Category: "dk-u13", MinWeight: 57.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u13-o60", Gender: "MALE", Category: "dk-u13", MinWeight: 60.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NỮ — TRẺ 12-13 (U13) — 8 hạng + 1 mở = 9
		// ═════════════════════════════════════════════════════
		{ID: "f-u13-u36", Gender: "FEMALE", Category: "dk-u13", MinWeight: 34.1, MaxWeight: 36, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-u38", Gender: "FEMALE", Category: "dk-u13", MinWeight: 36.1, MaxWeight: 38, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-u40", Gender: "FEMALE", Category: "dk-u13", MinWeight: 38.1, MaxWeight: 40, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-u42", Gender: "FEMALE", Category: "dk-u13", MinWeight: 40.1, MaxWeight: 42, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-u45", Gender: "FEMALE", Category: "dk-u13", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-u48", Gender: "FEMALE", Category: "dk-u13", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-u51", Gender: "FEMALE", Category: "dk-u13", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-u54", Gender: "FEMALE", Category: "dk-u13", MinWeight: 51.1, MaxWeight: 54, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u13-o54", Gender: "FEMALE", Category: "dk-u13", MinWeight: 54.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NAM — TRẺ 14-15 & 16-17 (cùng hạng cân) — 11 + 1 mở = 12
		// Dùng chung cho dk-u15 và dk-u17
		// ═════════════════════════════════════════════════════
		{ID: "m-u15-u45", Gender: "MALE", Category: "dk-u15", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u48", Gender: "MALE", Category: "dk-u15", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u51", Gender: "MALE", Category: "dk-u15", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u54", Gender: "MALE", Category: "dk-u15", MinWeight: 51.1, MaxWeight: 54, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u57", Gender: "MALE", Category: "dk-u15", MinWeight: 54.1, MaxWeight: 57, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u60", Gender: "MALE", Category: "dk-u15", MinWeight: 57.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u64", Gender: "MALE", Category: "dk-u15", MinWeight: 60.1, MaxWeight: 64, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u68", Gender: "MALE", Category: "dk-u15", MinWeight: 64.1, MaxWeight: 68, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u72", Gender: "MALE", Category: "dk-u15", MinWeight: 68.1, MaxWeight: 72, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u76", Gender: "MALE", Category: "dk-u15", MinWeight: 72.1, MaxWeight: 76, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u80", Gender: "MALE", Category: "dk-u15", MinWeight: 76.1, MaxWeight: 80, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-o80", Gender: "MALE", Category: "dk-u15", MinWeight: 80.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NỮ — TRẺ 14-15 & 16-17 (cùng hạng cân) — 8 + 1 mở = 9
		// ═════════════════════════════════════════════════════
		{ID: "f-u15-u45", Gender: "FEMALE", Category: "dk-u15", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u48", Gender: "FEMALE", Category: "dk-u15", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u51", Gender: "FEMALE", Category: "dk-u15", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u54", Gender: "FEMALE", Category: "dk-u15", MinWeight: 51.1, MaxWeight: 54, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u57", Gender: "FEMALE", Category: "dk-u15", MinWeight: 54.1, MaxWeight: 57, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u60", Gender: "FEMALE", Category: "dk-u15", MinWeight: 57.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u64", Gender: "FEMALE", Category: "dk-u15", MinWeight: 60.1, MaxWeight: 64, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u68", Gender: "FEMALE", Category: "dk-u15", MinWeight: 64.1, MaxWeight: 68, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-o68", Gender: "FEMALE", Category: "dk-u15", MinWeight: 68.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},
	}
}

// ────────────────────────────────────────────────────────────────
// HỆ THỐNG TÍNH ĐIỂM — SỬA ĐỔI ĐIỀU 14
// ────────────────────────────────────────────────────────────────

// ScoreType phân loại điểm số
type ScoreType string

const (
	ScoreHand      ScoreType = "hand"      // Đòn tay / chỏ
	ScoreFoot      ScoreType = "foot"      // Đòn chân / gối
	ScoreKnockdown ScoreType = "knockdown" // Đánh ngã / bắt chân đánh chân trụ
)

// ScoreRule quy định điểm cho từng loại đòn
type ScoreRule struct {
	Type        ScoreType `json:"type"`
	Points      int       `json:"points"`
	Description string    `json:"description"`
}

// Amendment2024ScoringRules trả về bảng tính điểm theo Điều 14 sửa đổi.
func Amendment2024ScoringRules() []ScoreRule {
	return []ScoreRule{
		{Type: ScoreHand, Points: 1, Description: "Đòn tay, đòn chỏ trúng vùng hợp lệ"},
		{Type: ScoreFoot, Points: 2, Description: "Đòn chân, đòn gối trúng vùng hợp lệ"},
		{Type: ScoreKnockdown, Points: 3, Description: "Đòn tay/chân đánh ngã, bắt chân đánh chân trụ hợp lệ"},
	}
}

// ────────────────────────────────────────────────────────────────
// PHÂN LOẠI LỖI — SỬA ĐỔI ĐIỀU 11
// ────────────────────────────────────────────────────────────────

// FoulSeverity mức độ vi phạm
type FoulSeverity string

const (
	FoulLight  FoulSeverity = "light"  // Lỗi nhẹ (K1 Đ11)
	FoulHeavy  FoulSeverity = "heavy"  // Lỗi nặng (K2 Đ11)
	FoulBanned FoulSeverity = "banned" // Đòn cấm (K3 Đ11)
)

// FoulRule mô tả loại vi phạm
type FoulRule struct {
	Code        string       `json:"code"`
	Severity    FoulSeverity `json:"severity"`
	Description string       `json:"description"`
}

// Amendment2024FoulRules trả về danh sách lỗi theo Điều 11 sửa đổi.
func Amendment2024FoulRules() []FoulRule {
	return []FoulRule{
		// ── Lỗi nhẹ (6 loại) ──
		{Code: "L01", Severity: FoulLight, Description: "Lôi kéo, xô đẩy, quăng quật, bốc, hất đối phương"},
		{Code: "L02", Severity: FoulLight, Description: "Kẹp găng"},
		{Code: "L03", Severity: FoulLight, Description: "Ôm ghì đối phương"},
		{Code: "L04", Severity: FoulLight, Description: "Cố tình di chuyển ra khỏi dây đài / vạch giới hạn"},
		{Code: "L05", Severity: FoulLight, Description: "Không tích cực thi đấu"},
		{Code: "L06", Severity: FoulLight, Description: "Không nghe khẩu lệnh trọng tài"},

		// ── Lỗi nặng (14 loại) ──
		{Code: "H01", Severity: FoulHeavy, Description: "Chẹn cổ đối phương"},
		{Code: "H02", Severity: FoulHeavy, Description: "Cố tình ôm vật đối phương"},
		{Code: "H03", Severity: FoulHeavy, Description: "Lợi dụng dây đài để ra đòn hoặc chống bị đánh ngã"},
		{Code: "H04", Severity: FoulHeavy, Description: "Cố tình tự ngã"},
		{Code: "H05", Severity: FoulHeavy, Description: "Ôm ghì tấn công đối phương"},
		{Code: "H06", Severity: FoulHeavy, Description: "Bắt chân tấn công (trừ bắt chân đánh chân trụ)"},
		{Code: "H07", Severity: FoulHeavy, Description: "Không tuân thủ theo lệnh trọng tài"},
		{Code: "H08", Severity: FoulHeavy, Description: "Bị đánh ngã cố tình lôi kéo đối phương ngã theo"},
		{Code: "H09", Severity: FoulHeavy, Description: "Cố tình nhả bảo vệ răng"},
		{Code: "H10", Severity: FoulHeavy, Description: "Không thực hiện / thực hiện sai phần trình diễn xe đài"},
		{Code: "H11", Severity: FoulHeavy, Description: "Lỗi trang phục thi đấu sai quy định"},
		{Code: "H12", Severity: FoulHeavy, Description: "Vi phạm đòn cấm chưa tới mức ảnh hưởng đối phương"},
		{Code: "H13", Severity: FoulHeavy, Description: "Lời nói/hành động xúc phạm đối phương, TT, BTC, khán giả"},

		// ── Đòn cấm (8 loại) ──
		{Code: "B01", Severity: FoulBanned, Description: "Húc đầu"},
		{Code: "B02", Severity: FoulBanned, Description: "Bẻ khớp"},
		{Code: "B03", Severity: FoulBanned, Description: "Cắn đối phương"},
		{Code: "B04", Severity: FoulBanned, Description: "Tấn công khớp gối, hạ bộ, gáy, cổ đối phương"},
		{Code: "B05", Severity: FoulBanned, Description: "Tấn công khi đối phương đã bị ngã"},
		{Code: "B06", Severity: FoulBanned, Description: "Ghì đầu/kẹp cổ/vít cổ để đánh chỏ, đánh gối"},
		{Code: "B07", Severity: FoulBanned, Description: "Đòn chỏ/gối (đối với giải cấm sử dụng)"},
		{Code: "B08", Severity: FoulBanned, Description: "Xoay lưng dùng chỏ đánh từ trên xuống đỉnh đầu (chỏ chồng)"},
	}
}

// ────────────────────────────────────────────────────────────────
// HỆ THỐNG PHẠT — SỬA ĐỔI ĐIỀU 12
// ────────────────────────────────────────────────────────────────

// PenaltyLevel cấp độ phạt
type PenaltyLevel string

const (
	PenaltyReminder   PenaltyLevel = "nhac_nho"      // Nhắc nhở
	PenaltyReprimand1 PenaltyLevel = "khien_trach_1" // Khiển trách lần 1
	PenaltyReprimand2 PenaltyLevel = "khien_trach_2" // Khiển trách lần 2
	PenaltyWarning1   PenaltyLevel = "canh_cao_1"    // Cảnh cáo lần 1
	PenaltyWarning2   PenaltyLevel = "canh_cao_2"    // Cảnh cáo lần 2
	PenaltyDQ         PenaltyLevel = "truat_quyen"   // Truất quyền thi đấu
)

// PenaltyRule mô tả hình thức phạt
type PenaltyRule struct {
	Level         PenaltyLevel `json:"level"`
	PointDeducted int          `json:"point_deducted"` // Số điểm bị trừ
	Scope         string       `json:"scope"`          // "round" hoặc "match"
	Description   string       `json:"description"`
}

// Amendment2024PenaltyRules trả về hệ thống phạt theo Điều 12 sửa đổi.
func Amendment2024PenaltyRules() []PenaltyRule {
	return []PenaltyRule{
		{Level: PenaltyReminder, PointDeducted: 0, Scope: "round", Description: "Nhắc nhở — không trừ điểm, bảo lưu trong hiệp"},
		{Level: PenaltyReprimand1, PointDeducted: 1, Scope: "round", Description: "Khiển trách lần 1 — trừ 1đ, bảo lưu trong hiệp"},
		{Level: PenaltyReprimand2, PointDeducted: 1, Scope: "round", Description: "Khiển trách lần 2 — trừ 1đ, bảo lưu trong hiệp"},
		{Level: PenaltyWarning1, PointDeducted: 2, Scope: "match", Description: "Cảnh cáo lần 1 — trừ 2đ, bảo lưu trong trận"},
		{Level: PenaltyWarning2, PointDeducted: 3, Scope: "match", Description: "Cảnh cáo lần 2 — trừ 3đ, bảo lưu trong trận"},
		{Level: PenaltyDQ, PointDeducted: 0, Scope: "match", Description: "Truất quyền thi đấu — loại khỏi trận"},
	}
}

// ────────────────────────────────────────────────────────────────
// THẢM ĐẤU — BỔ SUNG ĐIỀU 1 KHOẢN 3-4
// ────────────────────────────────────────────────────────────────

// MatSpec quy cách thảm đấu / sàn đấu
type MatSpec struct {
	TotalSizeM       float64 `json:"total_size_m"`       // 10 x 10
	CompetitionSizeM float64 `json:"competition_size_m"` // 8 x 8
	MaxThicknessCM   float64 `json:"max_thickness_cm"`   // 5 cm
	BoundaryWidthCM  float64 `json:"boundary_width_cm"`  // 5 cm
	CenterGapM       float64 `json:"center_gap_m"`       // 2m giữa 2 vạch
}

// Amendment2024MatSpec trả về quy cách thảm đấu theo Điều 1 sửa đổi.
func Amendment2024MatSpec() MatSpec {
	return MatSpec{
		TotalSizeM:       10.0,
		CompetitionSizeM: 8.0,
		MaxThicknessCM:   5.0,
		BoundaryWidthCM:  5.0,
		CenterGapM:       2.0,
	}
}

// ── Target Zone — Vùng đánh (Phụ lục 4) ─────────────────────

// TargetZone vùng đánh
type TargetZone string

const (
	ZoneScorable  TargetZone = "scorable"  // Được tính điểm: đầu/mặt, ngực/bụng/lườn
	ZoneLegal     TargetZone = "legal"     // Hợp lệ, không tính điểm: bắp tay, cẳng tay, đùi ngoài, bắp chân
	ZoneForbidden TargetZone = "forbidden" // Cấm đánh: đỉnh đầu, gáy, cổ, khớp gối, hạ bộ
)

// ────────────────────────────────────────────────────────────────
// EFFECTIVE (MERGED) FUNCTIONS
// Hàm trả về dữ liệu hiệu lực (gốc 2021 + sửa đổi 2024)
// Dùng cho store.seedData()
// ────────────────────────────────────────────────────────────────

// EffectiveAgeGroups trả về nhóm tuổi hiệu lực.
// Giữ nguyên MasterAgeGroup format cho backward compatibility,
// sử dụng category prefix trong ID để phân biệt.
func EffectiveAgeGroups() []MasterAgeGroup {
	amended := Amendment2024AgeGroups()
	result := make([]MasterAgeGroup, 0, len(amended))
	for _, ag := range amended {
		result = append(result, ag.MasterAgeGroup)
	}
	return result
}

// EffectiveWeightClasses trả về hạng cân hiệu lực (2024 thay thế hoàn toàn 2021).
func EffectiveWeightClasses() []MasterWeightClass {
	return Amendment2024WeightClasses()
}
