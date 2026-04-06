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
];

export const MOCK_CLUBS: Club[] = [
  { id: "CLB-HCM-001", code: "CLB-HCM-001", name: "Võ đường Bình Định Gia", province: "TP. Hồ Chí Minh", district: "Quận 1", address: "123 Nguyễn Huệ, Q.1", headCoach: "Đại Võ sư Nguyễn Văn Hùng", headCoachId: "P-001", phone: "028-38221234", email: "binhdinhgia@vct.vn", foundedDate: "2005-03-15", memberCount: 342, coachCount: 8, status: "active", rating: 5, facilities: ["Sân tập 200m²", "Phòng tập gym", "Phòng y tế"], martialArts: ["Võ Cổ Truyền"], monthlyFee: 500_000 },
];

export const MOCK_EXAMINATIONS: Examination[] = [
  { id: "EX-2025-004", code: "EX-2025-004", title: "Kỳ thi Thăng đai Quý IV/2025", date: "2025-12-20", location: "Nhà thi đấu Phú Thọ, TP.HCM", province: "TP. Hồ Chí Minh", beltLevel: "blue", targetBelt: "purple", status: "completed", candidateCount: 1247, passedCount: 1028, failedCount: 219, judges: [], candidates: [] },
];

export const MOCK_TOURNAMENTS: Tournament[] = [
  { id: "T-2025-001", code: "T-2025-001", name: "Giải Vô địch Võ Cổ Truyền Quốc gia 2025", startDate: "2025-12-15", endDate: "2025-12-22", location: "Nhà thi đấu Phú Thọ, TP.HCM", province: "TP. Hồ Chí Minh", status: "live", categories: [], teamCount: 45, athleteCount: 892, medalTable: [], organizer: "Liên đoàn Võ Cổ Truyền Việt Nam" },
];

export const MOCK_CERTIFICATES: Certification[] = [
  { id: "CERT-001", code: "CERT-2025-00001", type: "belt", title: "Chứng nhận Huyền đai Nhị đẳng", recipientId: "M-00001", recipientName: "Nguyễn Văn Minh", issuedDate: "2025-12-22", status: "issued", issuedBy: "ĐVS Nguyễn Văn Hùng", qrCode: "https://verify.vctplatform.vn/CERT-2025-00001" },
];

export const MOCK_FINANCE: FinanceRecord[] = [
  { id: "FIN-001", date: "2025-12-28", type: "income", category: "Niên liễm hội viên", description: "Thu niên liễm Q4/2025", amount: 1_228_000_000, status: "completed" },
];

export const MOCK_PERSONNEL: Personnel[] = [
  { id: "P-001", fullName: "Đại Võ sư Nguyễn Văn Hùng", role: "coach", beltLevel: "black_7dan", province: "TP. Hồ Chí Minh", phone: "0912000001", email: "hung.nvs@vct.vn", certifications: ["HLV Quốc gia"], joinDate: "2005-03-15", status: "active", specialization: "Võ Cổ Truyền Bình Định" },
];

export const MOCK_DOCUMENTS: Document[] = [
  { id: "DOC-001", code: "QC-2025-089", title: "Quy chế Thi đấu Võ Cổ Truyền 2026", type: "outgoing", category: "Quy chế", createdDate: "2025-12-01", createdBy: "Ban Chuyên môn", status: "approved", priority: "urgent", attachments: [], approvalChain: [] },
];

export const MOCK_COMMUNICATIONS: Communication[] = [
  { id: "COM-001", title: "Khai mạc Giải Vô địch Quốc gia 2025", type: "news", content: "Giải Vô địch Võ Cổ Truyền Quốc gia 2025 chính thức khai mạc...", publishDate: "2025-12-15", author: "Ban Truyền thông", status: "published", targetAudience: "all" },
];
