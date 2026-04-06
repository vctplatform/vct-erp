package federation

import "time"

// ═══════════════════════════════════════════════════════════════
// VCT PLATFORM — REGULATION CONFIG 2021-QC01
// Quy chế Quản lý Chuyên môn Võ Cổ Truyền Việt Nam
//
// Đây là file cấu hình CHUẨN QUỐC GIA, chứa toàn bộ
// hằng số và seed data theo quy chế 2021.
// CẤM chỉnh sửa dữ liệu nếu không có QĐ thay đổi quy chế.
// ═══════════════════════════════════════════════════════════════

// ── Metadata Quy chế ─────────────────────────────────────────

const (
	RegulationCode      = "2021-QC01"
	RegulationName      = "Quy chế Quản lý Chuyên môn Võ Cổ Truyền Việt Nam"
	RegulationYear      = 2021
	IssuingAuthority    = "Liên đoàn Võ thuật Cổ truyền Việt Nam"
	RegulationShortName = "QC QLCM VCT 2021"
)

// ── Hệ thống Đai — 18 bậc theo Quy chế (Chương II) ─────────

// NationalBelts trả về hệ thống 18 bậc đai chuẩn quốc gia.
// 12 bậc cấp (đai màu) + 6 bậc đẳng (đai đen).
func NationalBelts() []MasterBelt {
	now := time.Now()
	scope := BeltScopeNational

	return []MasterBelt{
		// ── 12 Bậc Cấp (Đai Màu) ────────────────────────────
		{Level: 1, Name: "Đai Trắng (Cấp 18)", ColorHex: "#FFFFFF", RequiredTimeMin: 0, IsDanLevel: false, Description: "Nhập môn — bắt đầu tập luyện", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 2, Name: "Đai Trắng Vạch Vàng (Cấp 17)", ColorHex: "#FFFFFF", RequiredTimeMin: 3, IsDanLevel: false, Description: "Nắm tấn pháp cơ bản", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 3, Name: "Đai Vàng (Cấp 16)", ColorHex: "#FBB724", RequiredTimeMin: 3, IsDanLevel: false, Description: "Quyền căn bản — tay chân phối hợp", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 4, Name: "Đai Vàng Vạch Xanh Lá (Cấp 15)", ColorHex: "#FBB724", RequiredTimeMin: 3, IsDanLevel: false, Description: "Kỹ thuật đá cơ bản", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 5, Name: "Đai Xanh Lá (Cấp 14)", ColorHex: "#4ADE80", RequiredTimeMin: 3, IsDanLevel: false, Description: "Quyền cước kết hợp — bài quyền đầu tiên", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 6, Name: "Đai Xanh Lá Vạch Xanh Dương (Cấp 13)", ColorHex: "#4ADE80", RequiredTimeMin: 4, IsDanLevel: false, Description: "Song luyện cơ bản", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 7, Name: "Đai Xanh Dương (Cấp 12)", ColorHex: "#3B82F6", RequiredTimeMin: 4, IsDanLevel: false, Description: "Bài quyền hoàn chỉnh — cước pháp nâng cao", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 8, Name: "Đai Xanh Dương Vạch Đỏ (Cấp 11)", ColorHex: "#3B82F6", RequiredTimeMin: 4, IsDanLevel: false, Description: "Đối luyện — kỹ thuật tự vệ", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 9, Name: "Đai Đỏ (Cấp 10)", ColorHex: "#EF4444", RequiredTimeMin: 6, IsDanLevel: false, Description: "Thực chiến cơ bản — bài binh khí đầu tiên", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 10, Name: "Đai Đỏ Vạch Nâu (Cấp 9)", ColorHex: "#EF4444", RequiredTimeMin: 6, IsDanLevel: false, Description: "Binh khí nâng cao — song luyện vũ khí", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 11, Name: "Đai Nâu (Cấp 8)", ColorHex: "#92400E", RequiredTimeMin: 6, IsDanLevel: false, Description: "Chuẩn bị đẳng cấp — toàn diện kỹ thuật", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 12, Name: "Đai Nâu Vạch Đen (Cấp 7)", ColorHex: "#92400E", RequiredTimeMin: 6, IsDanLevel: false, Description: "Hoàn thiện kỹ thuật — sẵn sàng thi đẳng", Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ── 6+ Bậc Đẳng (Đai Đen) ───────────────────────────
		{Level: 13, Name: "Đai Đen — Nhất Đẳng (1 Dan)", ColorHex: "#1E293B", RequiredTimeMin: 12, IsDanLevel: true, Description: "Sơ đẳng — đủ điều kiện HLV cấp cơ sở", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 14, Name: "Đai Đen — Nhị Đẳng (2 Dan)", ColorHex: "#1E293B", RequiredTimeMin: 24, IsDanLevel: true, Description: "Trung đẳng — đủ điều kiện HLV cấp cơ sở", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 15, Name: "Đai Đen — Tam Đẳng (3 Dan)", ColorHex: "#1E293B", RequiredTimeMin: 36, IsDanLevel: true, Description: "Trung đẳng — đủ điều kiện HLV cấp tỉnh", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 16, Name: "Đai Đen — Tứ Đẳng (4 Dan)", ColorHex: "#0F172A", RequiredTimeMin: 48, IsDanLevel: true, Description: "Cao đẳng — đủ điều kiện HLV cấp quốc gia", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 17, Name: "Đai Đen — Ngũ Đẳng (5 Dan)", ColorHex: "#0F172A", RequiredTimeMin: 60, IsDanLevel: true, Description: "Cao đẳng — Võ sư", Scope: scope, CreatedAt: now, UpdatedAt: now},
		{Level: 18, Name: "Đai Đen — Lục Đẳng trở lên (6+ Dan)", ColorHex: "#0F172A", RequiredTimeMin: 0, IsDanLevel: true, Description: "Đại sư — HĐ Thẩm định quốc gia xét duyệt", Scope: scope, CreatedAt: now, UpdatedAt: now},
	}
}

// ── Nhóm Tuổi — 7 nhóm theo Quy chế (Chương IV) ────────────

// NationalAgeGroups trả về 7 nhóm tuổi chuẩn quốc gia.
func NationalAgeGroups() []MasterAgeGroup {
	now := time.Now()
	scope := BeltScopeNational

	return []MasterAgeGroup{
		{ID: "u7", Name: "Mầm non", MinAge: 4, MaxAge: 6, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "u10", Name: "Nhi đồng nhỏ", MinAge: 7, MaxAge: 9, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "u12", Name: "Nhi đồng", MinAge: 10, MaxAge: 11, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "u15", Name: "Thiếu niên", MinAge: 12, MaxAge: 14, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "u18", Name: "Trẻ", MinAge: 15, MaxAge: 17, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "senior", Name: "Vô địch (18-35)", MinAge: 18, MaxAge: 35, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "masters", Name: "Cao niên (36+)", MinAge: 36, MaxAge: 99, Scope: scope, CreatedAt: now, UpdatedAt: now},
	}
}

// ── Hạng Cân Đối Kháng — theo Quy chế (Chương IV) ──────────

// NationalWeightClasses trả về hạng cân chuẩn quốc gia cho các nhóm tuổi.
func NationalWeightClasses() []MasterWeightClass {
	now := time.Now()
	scope := BeltScopeNational

	return []MasterWeightClass{
		// ═════════════════════════════════════════════════════
		// NAM — SENIOR (18-35)
		// ═════════════════════════════════════════════════════
		{ID: "m-sr-u48", Gender: "MALE", Category: "senior", MinWeight: 0, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u52", Gender: "MALE", Category: "senior", MinWeight: 48.1, MaxWeight: 52, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u56", Gender: "MALE", Category: "senior", MinWeight: 52.1, MaxWeight: 56, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u60", Gender: "MALE", Category: "senior", MinWeight: 56.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u65", Gender: "MALE", Category: "senior", MinWeight: 60.1, MaxWeight: 65, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u70", Gender: "MALE", Category: "senior", MinWeight: 65.1, MaxWeight: 70, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u75", Gender: "MALE", Category: "senior", MinWeight: 70.1, MaxWeight: 75, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u80", Gender: "MALE", Category: "senior", MinWeight: 75.1, MaxWeight: 80, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-u85", Gender: "MALE", Category: "senior", MinWeight: 80.1, MaxWeight: 85, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-sr-o85", Gender: "MALE", Category: "senior", MinWeight: 85.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NAM — U18 (15-17)
		// ═════════════════════════════════════════════════════
		{ID: "m-u18-u42", Gender: "MALE", Category: "u18", MinWeight: 0, MaxWeight: 42, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u45", Gender: "MALE", Category: "u18", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u48", Gender: "MALE", Category: "u18", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u51", Gender: "MALE", Category: "u18", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u54", Gender: "MALE", Category: "u18", MinWeight: 51.1, MaxWeight: 54, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u57", Gender: "MALE", Category: "u18", MinWeight: 54.1, MaxWeight: 57, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u60", Gender: "MALE", Category: "u18", MinWeight: 57.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u63", Gender: "MALE", Category: "u18", MinWeight: 60.1, MaxWeight: 63, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-u66", Gender: "MALE", Category: "u18", MinWeight: 63.1, MaxWeight: 66, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u18-o66", Gender: "MALE", Category: "u18", MinWeight: 66.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NAM — U15 (12-14)
		// ═════════════════════════════════════════════════════
		{ID: "m-u15-u30", Gender: "MALE", Category: "u15", MinWeight: 0, MaxWeight: 30, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u33", Gender: "MALE", Category: "u15", MinWeight: 30.1, MaxWeight: 33, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u36", Gender: "MALE", Category: "u15", MinWeight: 33.1, MaxWeight: 36, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u39", Gender: "MALE", Category: "u15", MinWeight: 36.1, MaxWeight: 39, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u42", Gender: "MALE", Category: "u15", MinWeight: 39.1, MaxWeight: 42, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u45", Gender: "MALE", Category: "u15", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u48", Gender: "MALE", Category: "u15", MinWeight: 45.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-u51", Gender: "MALE", Category: "u15", MinWeight: 48.1, MaxWeight: 51, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "m-u15-o51", Gender: "MALE", Category: "u15", MinWeight: 51.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NỮ — SENIOR (18-35)
		// ═════════════════════════════════════════════════════
		{ID: "f-sr-u44", Gender: "FEMALE", Category: "senior", MinWeight: 0, MaxWeight: 44, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u48", Gender: "FEMALE", Category: "senior", MinWeight: 44.1, MaxWeight: 48, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u52", Gender: "FEMALE", Category: "senior", MinWeight: 48.1, MaxWeight: 52, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u56", Gender: "FEMALE", Category: "senior", MinWeight: 52.1, MaxWeight: 56, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u60", Gender: "FEMALE", Category: "senior", MinWeight: 56.1, MaxWeight: 60, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u65", Gender: "FEMALE", Category: "senior", MinWeight: 60.1, MaxWeight: 65, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-u70", Gender: "FEMALE", Category: "senior", MinWeight: 65.1, MaxWeight: 70, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-sr-o70", Gender: "FEMALE", Category: "senior", MinWeight: 70.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NỮ — U18 (15-17)
		// ═════════════════════════════════════════════════════
		{ID: "f-u18-u40", Gender: "FEMALE", Category: "u18", MinWeight: 0, MaxWeight: 40, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u18-u43", Gender: "FEMALE", Category: "u18", MinWeight: 40.1, MaxWeight: 43, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u18-u46", Gender: "FEMALE", Category: "u18", MinWeight: 43.1, MaxWeight: 46, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u18-u49", Gender: "FEMALE", Category: "u18", MinWeight: 46.1, MaxWeight: 49, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u18-u52", Gender: "FEMALE", Category: "u18", MinWeight: 49.1, MaxWeight: 52, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u18-u55", Gender: "FEMALE", Category: "u18", MinWeight: 52.1, MaxWeight: 55, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u18-u58", Gender: "FEMALE", Category: "u18", MinWeight: 55.1, MaxWeight: 58, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u18-o58", Gender: "FEMALE", Category: "u18", MinWeight: 58.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},

		// ═════════════════════════════════════════════════════
		// NỮ — U15 (12-14)
		// ═════════════════════════════════════════════════════
		{ID: "f-u15-u28", Gender: "FEMALE", Category: "u15", MinWeight: 0, MaxWeight: 28, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u30", Gender: "FEMALE", Category: "u15", MinWeight: 28.1, MaxWeight: 30, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u33", Gender: "FEMALE", Category: "u15", MinWeight: 30.1, MaxWeight: 33, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u36", Gender: "FEMALE", Category: "u15", MinWeight: 33.1, MaxWeight: 36, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u39", Gender: "FEMALE", Category: "u15", MinWeight: 36.1, MaxWeight: 39, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u42", Gender: "FEMALE", Category: "u15", MinWeight: 39.1, MaxWeight: 42, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-u45", Gender: "FEMALE", Category: "u15", MinWeight: 42.1, MaxWeight: 45, Scope: scope, CreatedAt: now, UpdatedAt: now},
		{ID: "f-u15-o45", Gender: "FEMALE", Category: "u15", MinWeight: 45.1, MaxWeight: 0, IsHeavy: true, Scope: scope, CreatedAt: now, UpdatedAt: now},
	}
}

// ── Nội dung Thi đấu — 9 loại theo Quy chế (Chương III) ────

// NationalCompetitionContents trả về 9 nội dung thi đấu chuẩn quốc gia.
func NationalCompetitionContents() []MasterCompetitionContent {
	now := time.Now()
	scope := BeltScopeNational

	return []MasterCompetitionContent{
		{
			ID: "nd-doi-khang", Code: ContentDoiKhang,
			Name: "Đối kháng", Description: "Thi đấu 1 vs 1, có hạng cân, theo luật đối kháng",
			RequiresWeight: true, IsTeamEvent: false, MinAthletes: 1, MaxAthletes: 1, HasWeapon: false,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-quyen", Code: ContentQuyen,
			Name: "Quyền (cá nhân)", Description: "Biểu diễn bài quyền cá nhân — đánh giá kỹ thuật, sức mạnh, tốc độ",
			RequiresWeight: false, IsTeamEvent: false, MinAthletes: 1, MaxAthletes: 1, HasWeapon: false,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-quyen-dong-doi", Code: ContentQuyenDongDoi,
			Name: "Quyền đồng đội", Description: "Nhóm 3-5 người biểu diễn cùng bài quyền — đánh giá sự đồng đều",
			RequiresWeight: false, IsTeamEvent: true, MinAthletes: 3, MaxAthletes: 5, HasWeapon: false,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-song-luyen", Code: ContentSongLuyen,
			Name: "Song luyện", Description: "Hai người đấu mẫu — kỹ thuật công phòng phối hợp",
			RequiresWeight: false, IsTeamEvent: false, MinAthletes: 2, MaxAthletes: 2, HasWeapon: false,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-nhieu-luyen", Code: ContentNhieuLuyen,
			Name: "Nhiều luyện", Description: "Ba người trở lên đấu mẫu — kịch bản chiến đấu",
			RequiresWeight: false, IsTeamEvent: true, MinAthletes: 3, MaxAthletes: 10, HasWeapon: false,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-binh-khi", Code: ContentBinhKhi,
			Name: "Binh khí (quyền)", Description: "Cá nhân biểu diễn quyền binh khí: kiếm, đao, côn, thương...",
			RequiresWeight: false, IsTeamEvent: false, MinAthletes: 1, MaxAthletes: 1, HasWeapon: true,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-vu-khi-doi-luyen", Code: ContentVuKhiDoiLuyen,
			Name: "Vũ khí đối luyện", Description: "Hai người sử dụng binh khí đấu mẫu",
			RequiresWeight: false, IsTeamEvent: false, MinAthletes: 2, MaxAthletes: 2, HasWeapon: true,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-bieu-dien", Code: ContentBieuDien,
			Name: "Biểu diễn chiến lược", Description: "Kịch chiến đấu có kịch bản — đánh giá tính sáng tạo và kỹ thuật",
			RequiresWeight: false, IsTeamEvent: true, MinAthletes: 3, MaxAthletes: 15, HasWeapon: false,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
		{
			ID: "nd-tu-ve", Code: ContentTuVe,
			Name: "Tự vệ", Description: "Biểu diễn kỹ thuật tự vệ — ứng dụng chiến đấu thực tế",
			RequiresWeight: false, IsTeamEvent: false, MinAthletes: 1, MaxAthletes: 2, HasWeapon: false,
			Scope: scope, CreatedAt: now, UpdatedAt: now,
		},
	}
}

// ── Binh khí chính thức — theo Quy chế (Chương III) ─────────

// WeaponType represents types of traditional weapons.
type WeaponType string

const (
	WeaponKiem     WeaponType = "kiem"      // Kiếm (Sword)
	WeaponDao      WeaponType = "dao"       // Đao (Sabre/Broadsword)
	WeaponCon      WeaponType = "con"       // Côn/Bổng (Staff)
	WeaponThuong   WeaponType = "thuong"    // Thương/Giáo (Spear)
	WeaponSongDao  WeaponType = "song_dao"  // Song đao (Dual Sabres)
	WeaponSongKiem WeaponType = "song_kiem" // Song kiếm (Dual Swords)
	WeaponDaiDao   WeaponType = "dai_dao"   // Đại đao (Long-handled Sabre)
	WeaponRoi      WeaponType = "roi"       // Roi (Whip/Flexible weapon)
)

// WeaponInfo describes a regulation-approved weapon.
type WeaponInfo struct {
	Code        WeaponType `json:"code"`
	Name        string     `json:"name"`
	NameEN      string     `json:"name_en"`
	Description string     `json:"description"`
}

// NationalWeapons returns the list of regulation-approved weapons.
func NationalWeapons() []WeaponInfo {
	return []WeaponInfo{
		{Code: WeaponKiem, Name: "Kiếm", NameEN: "Sword", Description: "Kiếm thẳng — vũ khí truyền thống"},
		{Code: WeaponDao, Name: "Đao", NameEN: "Sabre/Broadsword", Description: "Đao cong — vũ khí chém"},
		{Code: WeaponCon, Name: "Côn/Bổng", NameEN: "Staff", Description: "Côn dài — vũ khí đánh tầm xa"},
		{Code: WeaponThuong, Name: "Thương/Giáo", NameEN: "Spear", Description: "Thương dài có mũi nhọn"},
		{Code: WeaponSongDao, Name: "Song đao", NameEN: "Dual Sabres", Description: "Hai đao sử dụng cùng lúc"},
		{Code: WeaponSongKiem, Name: "Song kiếm", NameEN: "Dual Swords", Description: "Hai kiếm sử dụng cùng lúc"},
		{Code: WeaponDaiDao, Name: "Đại đao", NameEN: "Long-handled Sabre", Description: "Đao cán dài — vũ khí nặng"},
		{Code: WeaponRoi, Name: "Roi", NameEN: "Whip/Flexible weapon", Description: "Roi mềm — vũ khí linh hoạt"},
	}
}

// ── Hằng số Chứng nhận — theo Quy chế (Chương V) ────────────

// CoachGrade represents coach certification levels.
type CoachGrade string

const (
	CoachGradeBasic    CoachGrade = "co_so"    // HLV Cấp cơ sở (Nhị Đẳng+)
	CoachGradeProvince CoachGrade = "cap_tinh" // HLV Cấp tỉnh (Tam Đẳng+)
	CoachGradeNational CoachGrade = "quoc_gia" // HLV Cấp quốc gia (Tứ Đẳng+)
)

// RefereeGrade represents referee certification levels.
type RefereeGrade string

const (
	RefereeGradeBasic    RefereeGrade = "co_so"    // TT Cấp cơ sở
	RefereeGradeProvince RefereeGrade = "cap_tinh" // TT Cấp tỉnh
	RefereeGradeNational RefereeGrade = "quoc_gia" // TT Cấp quốc gia
	RefereeGradeIntl     RefereeGrade = "quoc_te"  // TT Quốc tế
)

// CertValidity — Thời hạn chứng nhận (tháng)
const (
	CertValidityCoach   = 36 // 3 năm
	CertValidityReferee = 24 // 2 năm
	CertValidityClub    = 60 // 5 năm
	CertValidityMedical = 6  // 6 tháng
	CertValidityInsure  = 12 // 1 năm
)

// MinBeltForCoach — Đẳng cấp tối thiểu để đủ điều kiện HLV
var MinBeltForCoach = map[CoachGrade]int{
	CoachGradeBasic:    14, // Nhị Đẳng (level 14)
	CoachGradeProvince: 15, // Tam Đẳng (level 15)
	CoachGradeNational: 16, // Tứ Đẳng (level 16)
}
