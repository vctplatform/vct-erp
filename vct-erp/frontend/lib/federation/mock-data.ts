/* ─────────────────────────────────────────────────────────
 * Federation Mock Data
 * Realistic data for all 12 modules
 * ───────────────────────────────────────────────────────── */

import type {
  FederationDashboard,
  Member,
  Club,
  Examination,
  Tournament,
  Certificate,
  FinanceRecord,
  Personnel,
  Document,
  Communication,
  MedalEntry,
} from "./types";

// ── Helper ──────────────────────────────────────────────

const id = (prefix: string, n: number) => `${prefix}-${String(n).padStart(5, "0")}`;

// ── Dashboard Summary ───────────────────────────────────

export const MOCK_DASHBOARD: FederationDashboard = {
  totalMembers: 47_832,
  totalClubs: 1_247,
  totalCoaches: 3_456,
  totalReferees: 892,
  totalTournaments: 156,
  activeTournaments: 4,
  pendingApprovals: 23,
  memberGrowth: [
    { month: "T1", value: 41200, previousValue: 38500 },
    { month: "T2", value: 41800, previousValue: 39100 },
    { month: "T3", value: 42400, previousValue: 39800 },
    { month: "T4", value: 43100, previousValue: 40200 },
    { month: "T5", value: 43900, previousValue: 40900 },
    { month: "T6", value: 44500, previousValue: 41500 },
    { month: "T7", value: 45200, previousValue: 42200 },
    { month: "T8", value: 45800, previousValue: 42800 },
    { month: "T9", value: 46300, previousValue: 43300 },
    { month: "T10", value: 46900, previousValue: 43900 },
    { month: "T11", value: 47400, previousValue: 44400 },
    { month: "T12", value: 47832, previousValue: 44800 },
  ],
  clubDistribution: [
    { province: "TP. Hồ Chí Minh", provinceCode: "HCM", memberCount: 8420, clubCount: 187, percentage: 17.6 },
    { province: "Hà Nội", provinceCode: "HN", memberCount: 6350, clubCount: 142, percentage: 13.3 },
    { province: "Đà Nẵng", provinceCode: "DN", memberCount: 3120, clubCount: 78, percentage: 6.5 },
    { province: "Bình Dương", provinceCode: "BD", memberCount: 2840, clubCount: 65, percentage: 5.9 },
    { province: "Khánh Hòa", provinceCode: "KH", memberCount: 2100, clubCount: 52, percentage: 4.4 },
    { province: "Bình Định", provinceCode: "BDI", memberCount: 3450, clubCount: 89, percentage: 7.2 },
    { province: "Đồng Nai", provinceCode: "DNG", memberCount: 2200, clubCount: 48, percentage: 4.6 },
    { province: "Cần Thơ", provinceCode: "CT", memberCount: 1800, clubCount: 42, percentage: 3.8 },
    { province: "Lâm Đồng", provinceCode: "LD", memberCount: 1650, clubCount: 38, percentage: 3.4 },
    { province: "Thanh Hóa", provinceCode: "TH", memberCount: 1520, clubCount: 35, percentage: 3.2 },
  ],
  recentActivities: [
    { id: "act-1", type: "member_join", title: "478 hội viên mới đăng ký", description: "Tháng 12/2025 — tăng 12% so với tháng trước", timestamp: "2025-12-28T09:30:00Z" },
    { id: "act-2", type: "exam_complete", title: "Kỳ thi thăng đai Quý IV hoàn thành", description: "1.247 thí sinh tham dự, tỷ lệ đạt 82.4%", timestamp: "2025-12-22T16:00:00Z" },
    { id: "act-3", type: "tournament", title: "Giải Vô địch Quốc gia 2025 khởi tranh", description: "892 VĐV từ 45 tỉnh/thành tham dự", timestamp: "2025-12-15T08:00:00Z" },
    { id: "act-4", type: "club_approved", title: "Phê duyệt 12 CLB mới", description: "Các tỉnh: HCM (4), Hà Nội (3), Đà Nẵng (2), Bình Dương (3)", timestamp: "2025-12-10T14:30:00Z" },
    { id: "act-5", type: "certificate_issued", title: "Cấp 356 chứng chỉ HLV", description: "Đợt cấp chứng chỉ HLV cấp quốc gia lần 2/2025", timestamp: "2025-12-05T10:00:00Z" },
    { id: "act-6", type: "document", title: "Ban hành Quy chế thi đấu 2026", description: "Văn bản số QC-2025-089 đã được ký duyệt", timestamp: "2025-12-01T11:15:00Z" },
  ],
  alerts: [
    { id: "alert-1", type: "urgent", title: "23 đơn phê duyệt CLB mới đang chờ", description: "Vượt quá SLA 5 ngày xử lý", actionLabel: "Xử lý ngay", actionHref: "/federation/clubs?status=pending", timestamp: "2025-12-29T08:00:00Z" },
    { id: "alert-2", type: "warning", title: "Kỳ thi Q1/2026 chưa lên lịch", description: "Hạn chót đăng ký: 15/01/2026", actionLabel: "Tạo kỳ thi", actionHref: "/federation/examinations/new", timestamp: "2025-12-28T10:00:00Z" },
    { id: "alert-3", type: "info", title: "Giải Vô địch Quốc gia đang diễn ra", description: "Ngày thi đấu 3/7 — 245 trận đã hoàn thành", actionLabel: "Xem trực tiếp", actionHref: "/federation/tournaments/T-2025-001", timestamp: "2025-12-28T08:00:00Z" },
    { id: "alert-4", type: "success", title: "Ngân sách 2025 đã khóa sổ", description: "Tổng thu: 12.8 tỷ VND | Tổng chi: 10.2 tỷ VND", timestamp: "2025-12-27T17:00:00Z" },
  ],
  topClubs: [
    { rank: 1, clubId: "CLB-HCM-001", clubName: "Võ đường Bình Định Gia", province: "Bình Định", memberCount: 342, activityScore: 98.5, trend: "up" },
    { rank: 2, clubId: "CLB-HCM-012", clubName: "CLB Hùng Vương", province: "TP. Hồ Chí Minh", memberCount: 287, activityScore: 96.2, trend: "up" },
    { rank: 3, clubId: "CLB-HN-005", clubName: "Thiếu Lâm Bắc Phái", province: "Hà Nội", memberCount: 265, activityScore: 94.8, trend: "stable" },
    { rank: 4, clubId: "CLB-DN-003", clubName: "CLB Sơn Trà", province: "Đà Nẵng", memberCount: 198, activityScore: 93.1, trend: "up" },
    { rank: 5, clubId: "CLB-BD-007", clubName: "Võ đường Tân Sơn", province: "Bình Dương", memberCount: 176, activityScore: 91.5, trend: "down" },
  ],
  financeSummary: {
    totalIncome: 12_800_000_000,
    totalExpense: 10_200_000_000,
    balance: 2_600_000_000,
    monthlyIncome: [
      { month: "T1", value: 980_000_000 },
      { month: "T2", value: 1_050_000_000 },
      { month: "T3", value: 1_120_000_000 },
      { month: "T4", value: 890_000_000 },
      { month: "T5", value: 1_200_000_000 },
      { month: "T6", value: 1_350_000_000 },
      { month: "T7", value: 1_100_000_000 },
      { month: "T8", value: 950_000_000 },
      { month: "T9", value: 1_080_000_000 },
      { month: "T10", value: 1_150_000_000 },
      { month: "T11", value: 980_000_000 },
      { month: "T12", value: 950_000_000 },
    ],
    incomeByCategory: [
      { category: "Niên liễm hội viên", amount: 5_740_000_000, percentage: 44.8 },
      { category: "Phí thi thăng đai", amount: 2_560_000_000, percentage: 20.0 },
      { category: "Phí giải đấu", amount: 1_920_000_000, percentage: 15.0 },
      { category: "Tài trợ", amount: 1_280_000_000, percentage: 10.0 },
      { category: "Ngân sách Nhà nước", amount: 900_000_000, percentage: 7.0 },
      { category: "Khác", amount: 400_000_000, percentage: 3.2 },
    ],
  },
  beltDistribution: [
    { belt: "white", count: 12450, percentage: 26.0 },
    { belt: "yellow", count: 8960, percentage: 18.7 },
    { belt: "orange", count: 6720, percentage: 14.1 },
    { belt: "green", count: 5380, percentage: 11.2 },
    { belt: "blue", count: 4310, percentage: 9.0 },
    { belt: "purple", count: 3250, percentage: 6.8 },
    { belt: "brown", count: 2680, percentage: 5.6 },
    { belt: "red", count: 1920, percentage: 4.0 },
    { belt: "black_1dan", count: 1150, percentage: 2.4 },
    { belt: "black_2dan", count: 560, percentage: 1.2 },
    { belt: "black_3dan", count: 280, percentage: 0.6 },
    { belt: "black_4dan", count: 98, percentage: 0.2 },
    { belt: "black_5dan", count: 42, percentage: 0.1 },
    { belt: "black_6dan", count: 18, percentage: 0.04 },
    { belt: "black_7dan", count: 8, percentage: 0.02 },
    { belt: "black_8dan", count: 4, percentage: 0.01 },
    { belt: "black_9dan", count: 2, percentage: 0.004 },
  ],
};

// ── Members ─────────────────────────────────────────────

export const MOCK_MEMBERS: Member[] = [
  { id: id("M", 1), memberId: "VCT-2024-00001", fullName: "Nguyễn Văn Minh", dateOfBirth: "1998-03-15", gender: "male", phone: "0912345678", email: "minh.nv@email.com", province: "TP. Hồ Chí Minh", district: "Quận 1", address: "123 Nguyễn Huệ", currentBelt: "black_2dan", martialArt: "Võ Cổ Truyền", clubId: "CLB-HCM-001", clubName: "Võ đường Bình Định Gia", joinDate: "2018-06-01", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 2), memberId: "VCT-2024-00002", fullName: "Trần Thị Hương", dateOfBirth: "2001-07-22", gender: "female", phone: "0923456789", email: "huong.tt@email.com", province: "Hà Nội", district: "Cầu Giấy", address: "45 Xuân Thủy", currentBelt: "brown", martialArt: "Võ Cổ Truyền", clubId: "CLB-HN-005", clubName: "Thiếu Lâm Bắc Phái", joinDate: "2019-09-15", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 3), memberId: "VCT-2024-00003", fullName: "Lê Hoàng Phúc", dateOfBirth: "1995-11-08", gender: "male", phone: "0934567890", province: "Đà Nẵng", district: "Hải Châu", address: "78 Bạch Đằng", currentBelt: "black_3dan", martialArt: "Vovinam", clubId: "CLB-DN-003", clubName: "CLB Sơn Trà", joinDate: "2015-01-20", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 4), memberId: "VCT-2024-00004", fullName: "Phạm Ngọc Anh", dateOfBirth: "2003-05-30", gender: "female", phone: "0945678901", province: "Bình Dương", district: "Thủ Dầu Một", address: "12 Lê Lợi", currentBelt: "blue", martialArt: "Võ Cổ Truyền", clubId: "CLB-BD-007", clubName: "Võ đường Tân Sơn", joinDate: "2020-03-10", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 5), memberId: "VCT-2024-00005", fullName: "Võ Đình Khoa", dateOfBirth: "1992-09-18", gender: "male", phone: "0956789012", province: "Bình Định", district: "Quy Nhơn", address: "56 Trần Hưng Đạo", currentBelt: "black_4dan", martialArt: "Võ Cổ Truyền", clubId: "CLB-BDI-001", clubName: "Võ đường Tây Sơn", joinDate: "2010-08-05", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 6), memberId: "VCT-2024-00006", fullName: "Huỳnh Thị Mai", dateOfBirth: "2005-12-03", gender: "female", phone: "0967890123", province: "Cần Thơ", district: "Ninh Kiều", address: "89 Nguyễn Trãi", currentBelt: "green", martialArt: "Karate", clubId: "CLB-CT-002", clubName: "CLB Phong Cầm", joinDate: "2021-07-20", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 7), memberId: "VCT-2024-00007", fullName: "Đặng Quốc Bảo", dateOfBirth: "2000-01-25", gender: "male", phone: "0978901234", province: "Khánh Hòa", district: "Nha Trang", address: "34 Trần Phú", currentBelt: "red", martialArt: "Võ Cổ Truyền", clubId: "CLB-KH-004", clubName: "CLB Hải Phong", joinDate: "2017-04-12", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 8), memberId: "VCT-2024-00008", fullName: "Ngô Thanh Tùng", dateOfBirth: "1997-08-14", gender: "male", phone: "0989012345", province: "TP. Hồ Chí Minh", district: "Quận 7", address: "67 Nguyễn Thị Thập", currentBelt: "black_1dan", martialArt: "Võ Cổ Truyền", clubId: "CLB-HCM-012", clubName: "CLB Hùng Vương", joinDate: "2016-02-28", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 9), memberId: "VCT-2024-00009", fullName: "Bùi Thị Lan", dateOfBirth: "2004-04-19", gender: "female", phone: "0990123456", province: "Lâm Đồng", district: "Đà Lạt", address: "23 Phan Đình Phùng", currentBelt: "purple", martialArt: "Vovinam", clubId: "CLB-LD-001", clubName: "CLB Đà Lạt", joinDate: "2020-11-05", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 10), memberId: "VCT-2024-00010", fullName: "Trương Minh Đức", dateOfBirth: "1999-06-07", gender: "male", phone: "0901234567", province: "Thanh Hóa", district: "TP. Thanh Hóa", address: "45 Lê Hoàn", currentBelt: "brown", martialArt: "Võ Cổ Truyền", clubId: "CLB-TH-003", clubName: "Võ đường Lam Sơn", joinDate: "2018-09-22", status: "pending", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 11), memberId: "VCT-2024-00011", fullName: "Lý Thái Hòa", dateOfBirth: "1996-02-11", gender: "male", phone: "0911223344", province: "Hà Nội", district: "Đống Đa", address: "12 Tôn Đức Thắng", currentBelt: "black_1dan", martialArt: "Wushu", clubId: "CLB-HN-009", clubName: "CLB Thăng Long", joinDate: "2014-05-18", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
  { id: id("M", 12), memberId: "VCT-2024-00012", fullName: "Phan Thị Ngọc", dateOfBirth: "2006-10-28", gender: "female", phone: "0922334455", province: "Đồng Nai", district: "Biên Hòa", address: "78 Phạm Văn Thuận", currentBelt: "orange", martialArt: "Võ Cổ Truyền", clubId: "CLB-DNG-002", clubName: "CLB Biên Hòa", joinDate: "2022-01-10", status: "active", beltHistory: [], tournamentHistory: [], certificates: [] },
];

// ── Clubs ────────────────────────────────────────────────

export const MOCK_CLUBS: Club[] = [
  { id: "CLB-HCM-001", code: "CLB-HCM-001", name: "Võ đường Bình Định Gia", province: "TP. Hồ Chí Minh", district: "Quận 1", address: "123 Nguyễn Huệ, Q.1", headCoach: "Đại Võ sư Nguyễn Văn Hùng", headCoachId: "P-001", phone: "028-38221234", email: "binhdinhgia@vct.vn", foundedDate: "2005-03-15", memberCount: 342, coachCount: 8, status: "active", rating: 5, facilities: ["Sân tập 200m²", "Phòng tập gym", "Phòng y tế"], martialArts: ["Võ Cổ Truyền"], monthlyFee: 500_000 },
  { id: "CLB-HCM-012", code: "CLB-HCM-012", name: "CLB Hùng Vương", province: "TP. Hồ Chí Minh", district: "Quận 7", address: "67 Nguyễn Thị Thập, Q.7", headCoach: "VS Trần Minh Tuấn", headCoachId: "P-002", phone: "028-37751234", foundedDate: "2010-08-20", memberCount: 287, coachCount: 6, status: "active", rating: 4, facilities: ["Sân tập 150m²", "Kho vũ khí"], martialArts: ["Võ Cổ Truyền", "Vovinam"], monthlyFee: 450_000 },
  { id: "CLB-HN-005", code: "CLB-HN-005", name: "Thiếu Lâm Bắc Phái", province: "Hà Nội", district: "Cầu Giấy", address: "45 Xuân Thủy, Cầu Giấy", headCoach: "ĐVS Lê Quang Vinh", headCoachId: "P-003", phone: "024-37951234", foundedDate: "2002-11-10", memberCount: 265, coachCount: 7, status: "active", rating: 5, facilities: ["Sân tập 300m²", "Phòng lý thuyết", "Phòng y tế"], martialArts: ["Thiếu Lâm", "Võ Cổ Truyền"], monthlyFee: 600_000 },
  { id: "CLB-DN-003", code: "CLB-DN-003", name: "CLB Sơn Trà", province: "Đà Nẵng", district: "Hải Châu", address: "78 Bạch Đằng, Hải Châu", headCoach: "VS Phạm Đức Long", headCoachId: "P-004", phone: "0236-38221234", foundedDate: "2008-06-15", memberCount: 198, coachCount: 5, status: "active", rating: 4, facilities: ["Sân tập 180m²"], martialArts: ["Vovinam"], monthlyFee: 400_000 },
  { id: "CLB-BD-007", code: "CLB-BD-007", name: "Võ đường Tân Sơn", province: "Bình Dương", district: "Thủ Dầu Một", address: "12 Lê Lợi, TDM", headCoach: "VS Ngô Minh Phát", headCoachId: "P-005", phone: "0274-38221234", foundedDate: "2012-04-22", memberCount: 176, coachCount: 4, status: "active", rating: 4, facilities: ["Sân tập 120m²"], martialArts: ["Võ Cổ Truyền"], monthlyFee: 350_000 },
  { id: "CLB-BDI-001", code: "CLB-BDI-001", name: "Võ đường Tây Sơn", province: "Bình Định", district: "Quy Nhơn", address: "Tây Sơn, Quy Nhơn", headCoach: "ĐVS Võ Đình Hào", headCoachId: "P-006", phone: "0256-38221234", foundedDate: "1998-01-01", memberCount: 420, coachCount: 12, status: "active", rating: 5, facilities: ["Sân tập 500m²", "Phòng truyền thống", "Ký túc xá"], martialArts: ["Võ Cổ Truyền", "Bình Định Gia"], monthlyFee: 300_000 },
  { id: "CLB-CT-002", code: "CLB-CT-002", name: "CLB Phong Cầm", province: "Cần Thơ", district: "Ninh Kiều", address: "89 Nguyễn Trãi, Ninh Kiều", headCoach: "VS Huỳnh Minh Trí", headCoachId: "P-007", phone: "0292-38221234", foundedDate: "2015-09-30", memberCount: 134, coachCount: 3, status: "active", rating: 3, facilities: ["Sân tập 100m²"], martialArts: ["Karate", "Võ Cổ Truyền"], monthlyFee: 350_000 },
  { id: "CLB-NEW-001", code: "CLB-NEW-001", name: "Võ đường Long Hổ", province: "Đồng Nai", district: "Biên Hòa", address: "45 Phạm Văn Thuận", headCoach: "VS Trần Văn Lâm", headCoachId: "P-008", phone: "0251-38221234", foundedDate: "2025-11-01", memberCount: 0, coachCount: 1, status: "pending_approval", rating: 0, facilities: ["Sân tập 80m²"], martialArts: ["Võ Cổ Truyền"], monthlyFee: 300_000 },
];

// ── Examinations ────────────────────────────────────────

export const MOCK_EXAMINATIONS: Examination[] = [
  { id: "EX-2025-004", code: "EX-2025-004", title: "Kỳ thi Thăng đai Quý IV/2025 — Toàn quốc", date: "2025-12-20", endDate: "2025-12-22", location: "Nhà thi đấu Phú Thọ, TP.HCM", province: "TP. Hồ Chí Minh", beltLevel: "blue", targetBelt: "purple", status: "completed", candidateCount: 1247, passedCount: 1028, failedCount: 219, judges: [{ id: "J-1", name: "ĐVS Nguyễn Văn Hùng", role: "chief", beltLevel: "black_7dan", province: "TP. Hồ Chí Minh" }, { id: "J-2", name: "ĐVS Lê Quang Vinh", role: "member", beltLevel: "black_6dan", province: "Hà Nội" }], candidates: [] },
  { id: "EX-2026-001", code: "EX-2026-001", title: "Kỳ thi Thăng đai Quý I/2026 — Miền Nam", date: "2026-03-15", endDate: "2026-03-17", location: "Nhà thi đấu Tân Bình, TP.HCM", province: "TP. Hồ Chí Minh", beltLevel: "green", targetBelt: "blue", status: "scheduled", candidateCount: 856, passedCount: 0, failedCount: 0, judges: [], candidates: [] },
  { id: "EX-2026-002", code: "EX-2026-002", title: "Kỳ thi Thăng đai Quý I/2026 — Miền Bắc", date: "2026-03-22", endDate: "2026-03-23", location: "Cung thể thao Quần Ngựa, Hà Nội", province: "Hà Nội", beltLevel: "brown", targetBelt: "red", status: "scheduled", candidateCount: 423, passedCount: 0, failedCount: 0, judges: [], candidates: [] },
];

// ── Tournaments ─────────────────────────────────────────

export const MOCK_TOURNAMENTS: Tournament[] = [
  { id: "T-2025-001", code: "T-2025-001", name: "Giải Vô địch Võ Cổ Truyền Quốc gia 2025", startDate: "2025-12-15", endDate: "2025-12-22", location: "Nhà thi đấu Phú Thọ, TP.HCM", province: "TP. Hồ Chí Minh", status: "live", categories: [{ id: "C-1", name: "Đối kháng Nam 60kg", gender: "male", weightClass: "60kg", registeredCount: 32 }, { id: "C-2", name: "Đối kháng Nữ 52kg", gender: "female", weightClass: "52kg", registeredCount: 24 }, { id: "C-3", name: "Quyền thuật Nam", gender: "male", registeredCount: 48 }, { id: "C-4", name: "Biểu diễn Binh khí", gender: "mixed", registeredCount: 36 }], teamCount: 45, athleteCount: 892, medalTable: [{ teamName: "Bình Định", province: "Bình Định", gold: 8, silver: 5, bronze: 7, total: 20 }, { teamName: "TP. Hồ Chí Minh", province: "TP. Hồ Chí Minh", gold: 6, silver: 8, bronze: 9, total: 23 }, { teamName: "Hà Nội", province: "Hà Nội", gold: 5, silver: 6, bronze: 4, total: 15 }, { teamName: "Đà Nẵng", province: "Đà Nẵng", gold: 4, silver: 3, bronze: 5, total: 12 }, { teamName: "Bình Dương", province: "Bình Dương", gold: 3, silver: 4, bronze: 6, total: 13 }], organizer: "Liên đoàn Võ Cổ Truyền Việt Nam" },
  { id: "T-2026-001", code: "T-2026-001", name: "Giải Trẻ Võ Cổ Truyền Toàn quốc 2026", startDate: "2026-04-10", endDate: "2026-04-15", location: "Nhà thi đấu Đà Nẵng", province: "Đà Nẵng", status: "upcoming", categories: [{ id: "C-5", name: "U18 Nam 55kg", gender: "male", weightClass: "55kg", ageGroup: "U18", registeredCount: 28 }, { id: "C-6", name: "U18 Nữ 48kg", gender: "female", weightClass: "48kg", ageGroup: "U18", registeredCount: 20 }], teamCount: 38, athleteCount: 520, medalTable: [], organizer: "Liên đoàn Võ Cổ Truyền Việt Nam" },
  { id: "T-2026-002", code: "T-2026-002", name: "Giải Vô địch các CLB Mạnh 2026", startDate: "2026-07-20", endDate: "2026-07-25", location: "Cung thể thao Quần Ngựa, Hà Nội", province: "Hà Nội", status: "registration", categories: [], teamCount: 0, athleteCount: 0, medalTable: [], organizer: "Liên đoàn Võ Cổ Truyền Việt Nam" },
];

// ── Certificates ────────────────────────────────────────

export const MOCK_CERTIFICATES: Certificate[] = [
  { id: "CERT-001", code: "CERT-2025-00001", type: "belt", title: "Chứng nhận Huyền đai Nhị đẳng", recipientId: "M-00001", recipientName: "Nguyễn Văn Minh", issuedDate: "2025-12-22", status: "issued", issuedBy: "ĐVS Nguyễn Văn Hùng", qrCode: "https://verify.vctplatform.vn/CERT-2025-00001" },
  { id: "CERT-002", code: "CERT-2025-00002", type: "coach", title: "Chứng chỉ HLV Cấp Quốc gia", recipientId: "P-001", recipientName: "Nguyễn Văn Hùng", issuedDate: "2025-06-15", expiryDate: "2028-06-15", status: "issued", issuedBy: "Liên đoàn VCT Việt Nam", qrCode: "https://verify.vctplatform.vn/CERT-2025-00002" },
  { id: "CERT-003", code: "CERT-2025-00003", type: "referee", title: "Chứng chỉ Trọng tài Quốc gia", recipientId: "P-010", recipientName: "Trần Minh Tuấn", issuedDate: "2025-08-10", expiryDate: "2027-08-10", status: "issued", issuedBy: "Liên đoàn VCT Việt Nam", qrCode: "https://verify.vctplatform.vn/CERT-2025-00003" },
  { id: "CERT-004", code: "CERT-2025-00004", type: "achievement", title: "Huy chương Vàng — VĐQG 2025", recipientId: "M-00005", recipientName: "Võ Đình Khoa", issuedDate: "2025-12-22", status: "issued", issuedBy: "Liên đoàn VCT Việt Nam", qrCode: "https://verify.vctplatform.vn/CERT-2025-00004" },
  { id: "CERT-005", code: "CERT-2025-00005", type: "belt", title: "Chứng nhận Hồng đai", recipientId: "M-00007", recipientName: "Đặng Quốc Bảo", issuedDate: "2025-09-30", status: "pending", issuedBy: "", qrCode: "" },
];

// ── Finance Records ─────────────────────────────────────

export const MOCK_FINANCE: FinanceRecord[] = [
  { id: "FIN-001", date: "2025-12-28", type: "income", category: "Niên liễm hội viên", description: "Thu niên liễm Q4/2025 — 2.456 hội viên", amount: 1_228_000_000, status: "completed" },
  { id: "FIN-002", date: "2025-12-22", type: "income", category: "Phí thi thăng đai", description: "Phí thi đai Quý IV/2025 — 1.247 thí sinh", amount: 623_500_000, status: "completed" },
  { id: "FIN-003", date: "2025-12-15", type: "expense", category: "Tổ chức giải đấu", description: "Chi phí tổ chức Giải VĐQG 2025", amount: 850_000_000, reference: "T-2025-001", status: "completed" },
  { id: "FIN-004", date: "2025-12-10", type: "income", category: "Tài trợ", description: "Tài trợ từ Tập đoàn Viettel", amount: 500_000_000, status: "completed" },
  { id: "FIN-005", date: "2025-12-05", type: "expense", category: "Lương nhân sự", description: "Lương tháng 12/2025", amount: 245_000_000, status: "completed" },
  { id: "FIN-006", date: "2025-12-01", type: "expense", category: "IT & Công nghệ", description: "Chi phí hạ tầng VCT Platform — T12", amount: 85_000_000, status: "completed" },
  { id: "FIN-007", date: "2025-11-28", type: "income", category: "Phí giải đấu", description: "Phí đăng ký Giải VĐQG 2025", amount: 446_000_000, status: "completed" },
  { id: "FIN-008", date: "2025-11-20", type: "expense", category: "Đào tạo HLV", description: "Khóa đào tạo HLV đợt 2/2025", amount: 120_000_000, status: "completed" },
];

// ── Personnel ───────────────────────────────────────────

export const MOCK_PERSONNEL: Personnel[] = [
  { id: "P-001", fullName: "Đại Võ sư Nguyễn Văn Hùng", role: "coach", beltLevel: "black_7dan", province: "TP. Hồ Chí Minh", phone: "0912000001", email: "hung.nvs@vct.vn", certifications: ["HLV Quốc gia", "Trọng tài Quốc tế"], joinDate: "2005-03-15", status: "active", specialization: "Võ Cổ Truyền Bình Định" },
  { id: "P-002", fullName: "VS Trần Minh Tuấn", role: "coach", beltLevel: "black_5dan", province: "TP. Hồ Chí Minh", phone: "0912000002", certifications: ["HLV Quốc gia"], joinDate: "2010-08-20", status: "active", specialization: "Vovinam" },
  { id: "P-003", fullName: "ĐVS Lê Quang Vinh", role: "coach", beltLevel: "black_6dan", province: "Hà Nội", phone: "0912000003", certifications: ["HLV Quốc gia", "Giám khảo Quốc gia"], joinDate: "2002-11-10", status: "active", specialization: "Thiếu Lâm" },
  { id: "P-010", fullName: "Trần Đức Trọng", role: "referee", beltLevel: "black_4dan", province: "Đà Nẵng", phone: "0912000010", certifications: ["Trọng tài Quốc gia"], joinDate: "2012-05-20", status: "active", specialization: "Trọng tài Đối kháng" },
  { id: "P-020", fullName: "GS.TS Nguyễn Minh Châu", role: "board_member", beltLevel: "black_8dan", province: "Hà Nội", phone: "0912000020", certifications: ["Chủ tịch Liên đoàn"], joinDate: "2000-01-01", status: "active", specialization: "Quản lý Thể thao" },
  { id: "P-021", fullName: "PGS.TS Lê Thanh Hà", role: "board_member", beltLevel: "black_6dan", province: "TP. Hồ Chí Minh", phone: "0912000021", certifications: ["Phó Chủ tịch"], joinDate: "2005-06-15", status: "active", specialization: "Đào tạo VĐV" },
];

// ── Documents ───────────────────────────────────────────

export const MOCK_DOCUMENTS: Document[] = [
  { id: "DOC-001", code: "QC-2025-089", title: "Quy chế Thi đấu Võ Cổ Truyền 2026", type: "outgoing", category: "Quy chế", createdDate: "2025-12-01", createdBy: "Ban Chuyên môn", status: "approved", priority: "urgent", attachments: ["qc-2026.pdf"], approvalChain: [{ order: 1, approver: "Trưởng Ban Chuyên môn", role: "Phê duyệt nội dung", status: "approved", date: "2025-12-03" }, { order: 2, approver: "Chủ tịch Liên đoàn", role: "Phê duyệt ban hành", status: "approved", date: "2025-12-05" }] },
  { id: "DOC-002", code: "CV-2025-456", title: "Công văn triệu tập VĐV tập huấn SEA Games", type: "outgoing", category: "Công văn", createdDate: "2025-11-20", createdBy: "Ban Thư ký", status: "approved", priority: "critical", attachments: ["cv-456.pdf", "ds-vdv.xlsx"], approvalChain: [{ order: 1, approver: "Tổng Thư ký", role: "Soạn thảo", status: "approved", date: "2025-11-21" }, { order: 2, approver: "Chủ tịch", role: "Ký duyệt", status: "approved", date: "2025-11-22" }] },
  { id: "DOC-003", code: "TB-2025-123", title: "Thông báo lịch thi đai Quý I/2026", type: "outgoing", category: "Thông báo", createdDate: "2025-12-28", createdBy: "Ban Kỹ thuật", status: "pending_review", priority: "normal", attachments: [], approvalChain: [{ order: 1, approver: "Trưởng Ban Kỹ thuật", role: "Soạn thảo", status: "approved", date: "2025-12-28" }, { order: 2, approver: "Tổng Thư ký", role: "Phê duyệt", status: "pending" }] },
  { id: "DOC-004", code: "BC-2025-078", title: "Báo cáo Tổng kết hoạt động năm 2025", type: "internal", category: "Báo cáo", createdDate: "2025-12-25", createdBy: "Ban Thư ký", status: "draft", priority: "normal", attachments: ["bc-2025.docx"], approvalChain: [] },
];

// ── Communications ──────────────────────────────────────

export const MOCK_COMMUNICATIONS: Communication[] = [
  { id: "COM-001", title: "Khai mạc Giải Vô địch Quốc gia 2025", type: "news", content: "Giải Vô địch Võ Cổ Truyền Quốc gia 2025 chính thức khai mạc tại Nhà thi đấu Phú Thọ...", publishDate: "2025-12-15", author: "Ban Truyền thông", status: "published", targetAudience: "all", readCount: 15420 },
  { id: "COM-002", title: "Thông báo lịch thi đai Quý I/2026", type: "announcement", content: "Liên đoàn thông báo lịch thi thăng đai Quý I năm 2026 cho toàn bộ hội viên...", publishDate: "2025-12-28", author: "Ban Kỹ thuật", status: "published", targetAudience: "all", readCount: 8930 },
  { id: "COM-003", title: "Hội thảo Đào tạo HLV Quốc gia 2026", type: "event", content: "Chương trình đào tạo và cấp chứng chỉ HLV cấp Quốc gia năm 2026...", publishDate: "2026-01-05", author: "Ban Đào tạo", status: "scheduled", targetAudience: "coaches", eventDate: "2026-02-15", eventLocation: "Hà Nội" },
  { id: "COM-004", title: "Cập nhật Quy chế thi đấu mới 2026", type: "announcement", content: "Ban hành Quy chế thi đấu Võ Cổ Truyền mới áp dụng từ 01/01/2026...", publishDate: "2025-12-05", author: "Ban Chuyên môn", status: "published", targetAudience: "clubs", readCount: 12100 },
];
