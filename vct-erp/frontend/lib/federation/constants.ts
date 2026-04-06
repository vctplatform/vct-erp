/* ─────────────────────────────────────────────────────────
 * Federation Constants & Configuration
 * ───────────────────────────────────────────────────────── */

import type { BeltLevel } from "./types";

// ── Belt System ─────────────────────────────────────────

export const BELT_CONFIG: Record<
  BeltLevel,
  { label: string; labelEn: string; color: string; order: number }
> = {
  white: { label: "Trắng", labelEn: "White", color: "#F8FAFC", order: 1 },
  yellow: { label: "Vàng", labelEn: "Yellow", color: "#FACC15", order: 2 },
  orange: { label: "Cam", labelEn: "Orange", color: "#FB923C", order: 3 },
  green: { label: "Xanh lá", labelEn: "Green", color: "#22C55E", order: 4 },
  blue: { label: "Xanh dương", labelEn: "Blue", color: "#3B82F6", order: 5 },
  purple: { label: "Tím", labelEn: "Purple", color: "#A855F7", order: 6 },
  brown: { label: "Nâu", labelEn: "Brown", color: "#A16207", order: 7 },
  red: { label: "Đỏ", labelEn: "Red", color: "#EF4444", order: 8 },
  black_1dan: { label: "Đen - 1 Đẳng", labelEn: "Black 1st Dan", color: "#18181B", order: 9 },
  black_2dan: { label: "Đen - 2 Đẳng", labelEn: "Black 2nd Dan", color: "#18181B", order: 10 },
  black_3dan: { label: "Đen - 3 Đẳng", labelEn: "Black 3rd Dan", color: "#18181B", order: 11 },
  black_4dan: { label: "Đen - 4 Đẳng", labelEn: "Black 4th Dan", color: "#18181B", order: 12 },
  black_5dan: { label: "Đen - 5 Đẳng", labelEn: "Black 5th Dan", color: "#18181B", order: 13 },
  black_6dan: { label: "Đen - 6 Đẳng", labelEn: "Black 6th Dan", color: "#18181B", order: 14 },
  black_7dan: { label: "Đen - 7 Đẳng", labelEn: "Black 7th Dan", color: "#18181B", order: 15 },
  black_8dan: { label: "Đen - 8 Đẳng", labelEn: "Black 8th Dan", color: "#18181B", order: 16 },
  black_9dan: { label: "Đen - 9 Đẳng", labelEn: "Black 9th Dan", color: "#18181B", order: 17 },
};

// ── Vietnam Provinces (63 tỉnh/thành) ──────────────────

export const PROVINCES = [
  { code: "HCM", name: "TP. Hồ Chí Minh", region: "south" },
  { code: "HN", name: "Hà Nội", region: "north" },
  { code: "DN", name: "Đà Nẵng", region: "central" },
  { code: "BD", name: "Bình Dương", region: "south" },
  { code: "DNG", name: "Đồng Nai", region: "south" },
  { code: "KH", name: "Khánh Hòa", region: "central" },
  { code: "HP", name: "Hải Phòng", region: "north" },
  { code: "CT", name: "Cần Thơ", region: "south" },
  { code: "LA", name: "Long An", region: "south" },
  { code: "QN", name: "Quảng Nam", region: "central" },
  { code: "BR", name: "Bà Rịa - Vũng Tàu", region: "south" },
  { code: "LD", name: "Lâm Đồng", region: "central" },
  { code: "TH", name: "Thanh Hóa", region: "north" },
  { code: "NA", name: "Nghệ An", region: "north" },
  { code: "AG", name: "An Giang", region: "south" },
  { code: "BN", name: "Bắc Ninh", region: "north" },
  { code: "GL", name: "Gia Lai", region: "central" },
  { code: "TN", name: "Tây Ninh", region: "south" },
  { code: "BT", name: "Bình Thuận", region: "central" },
  { code: "TTH", name: "Thừa Thiên Huế", region: "central" },
  { code: "VL", name: "Vĩnh Long", region: "south" },
  { code: "TG", name: "Tiền Giang", region: "south" },
  { code: "BL", name: "Bạc Liêu", region: "south" },
  { code: "ST", name: "Sóc Trăng", region: "south" },
  { code: "TV", name: "Trà Vinh", region: "south" },
  { code: "VP", name: "Vĩnh Phúc", region: "north" },
  { code: "QNI", name: "Quảng Ninh", region: "north" },
  { code: "HD", name: "Hải Dương", region: "north" },
  { code: "HY", name: "Hưng Yên", region: "north" },
  { code: "PY", name: "Phú Yên", region: "central" },
  { code: "BDI", name: "Bình Định", region: "central" },
  { code: "QB", name: "Quảng Bình", region: "central" },
  { code: "DL", name: "Đắk Lắk", region: "central" },
  { code: "CM", name: "Cà Mau", region: "south" },
  { code: "KG", name: "Kiên Giang", region: "south" },
  { code: "HG", name: "Hậu Giang", region: "south" },
  { code: "DT", name: "Đồng Tháp", region: "south" },
  { code: "BG", name: "Bắc Giang", region: "north" },
  { code: "TB", name: "Thái Bình", region: "north" },
  { code: "NB", name: "Ninh Bình", region: "north" },
  { code: "PT", name: "Phú Thọ", region: "north" },
  { code: "LS", name: "Lạng Sơn", region: "north" },
  { code: "TQ", name: "Tuyên Quang", region: "north" },
  { code: "YB", name: "Yên Bái", region: "north" },
  { code: "LC", name: "Lào Cai", region: "north" },
  { code: "SL", name: "Sơn La", region: "north" },
  { code: "HB", name: "Hòa Bình", region: "north" },
  { code: "CB", name: "Cao Bằng", region: "north" },
  { code: "BK", name: "Bắc Kạn", region: "north" },
  { code: "DB", name: "Điện Biên", region: "north" },
  { code: "LCH", name: "Lai Châu", region: "north" },
  { code: "HGI", name: "Hà Giang", region: "north" },
  { code: "ND", name: "Nam Định", region: "north" },
  { code: "HNA", name: "Hà Nam", region: "north" },
  { code: "NTR", name: "Ninh Thuận", region: "central" },
  { code: "QT", name: "Quảng Trị", region: "central" },
  { code: "QNG", name: "Quảng Ngãi", region: "central" },
  { code: "KT", name: "Kon Tum", region: "central" },
  { code: "DNO", name: "Đắk Nông", region: "central" },
  { code: "TNN", name: "Thái Nguyên", region: "north" },
  { code: "BP", name: "Bình Phước", region: "south" },
  { code: "HT", name: "Hà Tĩnh", region: "north" },
  { code: "BNI", name: "Bến Tre", region: "south" },
] as const;

export type ProvinceCode = (typeof PROVINCES)[number]["code"];

// ── Martial Art Types ───────────────────────────────────

export const MARTIAL_ARTS = [
  "Võ Cổ Truyền",
  "Vovinam",
  "Karate",
  "Taekwondo",
  "Judo",
  "Wushu",
  "Pencak Silat",
  "Muay Thai",
  "Kickboxing",
  "Boxing",
] as const;

// ── Finance Categories ──────────────────────────────────

export const FINANCE_INCOME_CATEGORIES = [
  "Niên liễm hội viên",
  "Phí thi thăng đai",
  "Phí giải đấu",
  "Tài trợ",
  "Ngân sách Nhà nước",
  "Phí cấp chứng chỉ",
  "Khác",
] as const;

export const FINANCE_EXPENSE_CATEGORIES = [
  "Tổ chức giải đấu",
  "Tổ chức thi đai",
  "Lương nhân sự",
  "Thuê mặt bằng",
  "Đào tạo HLV",
  "IT & Công nghệ",
  "Văn phòng phẩm",
  "Đi lại công tác",
  "Truyền thông",
  "Khác",
] as const;

// ── Status Colors ───────────────────────────────────────

export const STATUS_COLORS = {
  active: { bg: "bg-emerald-500/15", text: "text-emerald-400", dot: "bg-emerald-400" },
  inactive: { bg: "bg-slate-500/15", text: "text-slate-400", dot: "bg-slate-400" },
  suspended: { bg: "bg-rose-500/15", text: "text-rose-400", dot: "bg-rose-400" },
  pending: { bg: "bg-amber-500/15", text: "text-amber-400", dot: "bg-amber-400" },
  expired: { bg: "bg-red-500/15", text: "text-red-400", dot: "bg-red-400" },
  pending_approval: { bg: "bg-amber-500/15", text: "text-amber-400", dot: "bg-amber-400" },
  dissolved: { bg: "bg-slate-500/15", text: "text-slate-400", dot: "bg-slate-400" },
  scheduled: { bg: "bg-sky-500/15", text: "text-sky-400", dot: "bg-sky-400" },
  in_progress: { bg: "bg-cyan-500/15", text: "text-cyan-400", dot: "bg-cyan-400" },
  completed: { bg: "bg-emerald-500/15", text: "text-emerald-400", dot: "bg-emerald-400" },
  cancelled: { bg: "bg-slate-500/15", text: "text-slate-400", dot: "bg-slate-400" },
  upcoming: { bg: "bg-indigo-500/15", text: "text-indigo-400", dot: "bg-indigo-400" },
  registration: { bg: "bg-violet-500/15", text: "text-violet-400", dot: "bg-violet-400" },
  live: { bg: "bg-rose-500/15", text: "text-rose-400", dot: "bg-rose-400 animate-pulse" },
  issued: { bg: "bg-emerald-500/15", text: "text-emerald-400", dot: "bg-emerald-400" },
  revoked: { bg: "bg-rose-500/15", text: "text-rose-400", dot: "bg-rose-400" },
  draft: { bg: "bg-slate-500/15", text: "text-slate-400", dot: "bg-slate-400" },
  pending_review: { bg: "bg-amber-500/15", text: "text-amber-400", dot: "bg-amber-400" },
  approved: { bg: "bg-emerald-500/15", text: "text-emerald-400", dot: "bg-emerald-400" },
  rejected: { bg: "bg-rose-500/15", text: "text-rose-400", dot: "bg-rose-400" },
  archived: { bg: "bg-slate-500/15", text: "text-slate-400", dot: "bg-slate-400" },
  published: { bg: "bg-emerald-500/15", text: "text-emerald-400", dot: "bg-emerald-400" },
} as const;

export type StatusKey = keyof typeof STATUS_COLORS;

// ── Navigation ──────────────────────────────────────────

export const FEDERATION_NAV = [
  {
    group: "Tổng quan",
    groupEn: "Overview",
    items: [
      { href: "/federation", label: "Dashboard", labelEn: "Dashboard", icon: "LayoutDashboard" },
    ],
  },
  {
    group: "Quản lý",
    groupEn: "Management",
    items: [
      { href: "/federation/members", label: "Hội viên", labelEn: "Members", icon: "Users" },
      { href: "/federation/clubs", label: "CLB / Võ đường", labelEn: "Clubs", icon: "Building2" },
      { href: "/federation/personnel", label: "Nhân sự & HLV", labelEn: "Personnel", icon: "UserCog" },
    ],
  },
  {
    group: "Nghiệp vụ",
    groupEn: "Operations",
    items: [
      { href: "/federation/examinations", label: "Thi & Thăng đai", labelEn: "Examinations", icon: "GraduationCap" },
      { href: "/federation/tournaments", label: "Giải đấu", labelEn: "Tournaments", icon: "Trophy" },
      { href: "/federation/certificates", label: "Chứng chỉ số", labelEn: "Certificates", icon: "Award" },
    ],
  },
  {
    group: "Vận hành",
    groupEn: "Administration",
    items: [
      { href: "/federation/finance", label: "Tài chính", labelEn: "Finance", icon: "Wallet" },
      { href: "/federation/documents", label: "Văn bản", labelEn: "Documents", icon: "FileText" },
      { href: "/federation/communications", label: "Truyền thông", labelEn: "Communications", icon: "Megaphone" },
    ],
  },
  {
    group: "Hệ thống",
    groupEn: "System",
    items: [
      { href: "/federation/reports", label: "Báo cáo", labelEn: "Reports", icon: "BarChart3" },
      { href: "/federation/settings", label: "Cài đặt", labelEn: "Settings", icon: "Settings" },
    ],
  },
] as const;
