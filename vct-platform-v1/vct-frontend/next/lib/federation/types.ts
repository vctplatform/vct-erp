/* ─────────────────────────────────────────────────────────
 * Federation Domain Types
 * Phân hệ Quản lý Liên đoàn Cấp Quốc gia — VCT Platform
 * ───────────────────────────────────────────────────────── */

// ── Core Enums ──────────────────────────────────────────

export type BeltLevel =
  | "white"
  | "yellow"
  | "orange"
  | "green"
  | "blue"
  | "purple"
  | "brown"
  | "red"
  | "black_1dan"
  | "black_2dan"
  | "black_3dan"
  | "black_4dan"
  | "black_5dan"
  | "black_6dan"
  | "black_7dan"
  | "black_8dan"
  | "black_9dan";

export type MemberStatus = "active" | "inactive" | "suspended" | "pending" | "expired";
export type ClubStatus = "active" | "pending_approval" | "suspended" | "dissolved";
export type ExamStatus = "scheduled" | "in_progress" | "completed" | "cancelled";
export type TournamentStatus = "upcoming" | "registration" | "live" | "completed" | "cancelled";
export type CertificateStatus = "issued" | "pending" | "revoked" | "expired";
export type DocumentStatus = "draft" | "pending_review" | "approved" | "rejected" | "archived";
export type PersonnelRole = "coach" | "referee" | "board_member" | "secretary" | "treasurer";
export type Gender = "male" | "female" | "other";

// ── Core Entities ───────────────────────────────────────

export interface Member {
  id: string;
  memberId: string; // VCT-2024-XXXXX
  fullName: string;
  dateOfBirth: string;
  gender: Gender;
  phone: string;
  email?: string;
  province: string;
  district: string;
  address: string;
  avatar?: string;
  currentBelt: BeltLevel;
  martialArt: string; // Võ cổ truyền, Vovinam, etc.
  clubId: string;
  clubName: string;
  joinDate: string;
  status: MemberStatus;
  nationalCardNumber?: string; // CMND/CCCD
  beltHistory: BeltRecord[];
  tournamentHistory: TournamentParticipation[];
  certificates: CertificateRef[];
}

export interface BeltRecord {
  id: string;
  belt: BeltLevel;
  examId: string;
  examDate: string;
  examLocation: string;
  score: number;
  passed: boolean;
  issuedBy: string;
}

export interface Club {
  id: string;
  code: string; // CLB-HCM-001
  name: string;
  province: string;
  district: string;
  address: string;
  headCoach: string;
  headCoachId: string;
  phone: string;
  email?: string;
  foundedDate: string;
  memberCount: number;
  coachCount: number;
  status: ClubStatus;
  rating: number; // 1-5 stars
  logo?: string;
  facilities: string[];
  martialArts: string[];
  monthlyFee: number;
  lastInspectionDate?: string;
}

export interface Examination {
  id: string;
  code: string; // EX-2024-001
  title: string;
  date: string;
  endDate?: string;
  location: string;
  province: string;
  beltLevel: BeltLevel;
  targetBelt: BeltLevel; // Belt candidates are testing for
  status: ExamStatus;
  candidateCount: number;
  passedCount: number;
  failedCount: number;
  judges: Judge[];
  candidates: ExamCandidate[];
}

export interface Judge {
  id: string;
  name: string;
  role: "chief" | "member";
  beltLevel: BeltLevel;
  province: string;
}

export interface ExamCandidate {
  id: string;
  memberId: string;
  memberName: string;
  clubName: string;
  currentBelt: BeltLevel;
  targetBelt: BeltLevel;
  score?: number;
  result?: "pass" | "fail" | "pending";
  notes?: string;
}

export interface Tournament {
  id: string;
  code: string; // T-2024-001
  name: string;
  startDate: string;
  endDate: string;
  location: string;
  province: string;
  status: TournamentStatus;
  categories: TournamentCategory[];
  teamCount: number;
  athleteCount: number;
  medalTable: MedalEntry[];
  organizer: string;
  sponsoredBy?: string[];
}

export interface TournamentCategory {
  id: string;
  name: string;
  gender: Gender | "mixed";
  weightClass?: string;
  ageGroup?: string;
  registeredCount: number;
}

export interface TournamentParticipation {
  tournamentId: string;
  tournamentName: string;
  date: string;
  category: string;
  result: string;
  medal?: "gold" | "silver" | "bronze";
}

export interface MedalEntry {
  teamName: string;
  province: string;
  gold: number;
  silver: number;
  bronze: number;
  total: number;
}

export interface Certificate {
  id: string;
  code: string; // CERT-2024-XXXXX
  type: "belt" | "coach" | "referee" | "achievement" | "participation";
  title: string;
  recipientId: string;
  recipientName: string;
  issuedDate: string;
  expiryDate?: string;
  status: CertificateStatus;
  issuedBy: string;
  qrCode: string;
  digitalSignature?: string;
}

export type CertificateRef = Pick<Certificate, "id" | "code" | "type" | "title" | "issuedDate" | "status">;

export interface FinanceRecord {
  id: string;
  date: string;
  type: "income" | "expense";
  category: string;
  description: string;
  amount: number;
  reference?: string;
  approvedBy?: string;
  status: "completed" | "pending" | "cancelled";
}

export interface Personnel {
  id: string;
  fullName: string;
  role: PersonnelRole;
  beltLevel: BeltLevel;
  province: string;
  phone: string;
  email?: string;
  avatar?: string;
  certifications: string[];
  joinDate: string;
  status: MemberStatus;
  specialization?: string;
}

export interface Document {
  id: string;
  code: string; // CV-2024-001
  title: string;
  type: "incoming" | "outgoing" | "internal";
  category: string;
  content?: string;
  createdDate: string;
  createdBy: string;
  status: DocumentStatus;
  priority: "normal" | "urgent" | "critical";
  attachments: string[];
  approvalChain: ApprovalStep[];
}

export interface ApprovalStep {
  order: number;
  approver: string;
  role: string;
  status: "pending" | "approved" | "rejected";
  date?: string;
  comments?: string;
}

export interface Communication {
  id: string;
  title: string;
  type: "announcement" | "news" | "event";
  content: string;
  publishDate: string;
  author: string;
  status: "published" | "draft" | "scheduled";
  targetAudience: "all" | "clubs" | "coaches" | "members";
  readCount?: number;
  eventDate?: string;
  eventLocation?: string;
}

// ── Dashboard Aggregates ────────────────────────────────

export interface FederationDashboard {
  totalMembers: number;
  totalClubs: number;
  totalCoaches: number;
  totalReferees: number;
  totalTournaments: number;
  activeTournaments: number;
  pendingApprovals: number;
  memberGrowth: MonthlyDataPoint[];
  clubDistribution: ProvinceData[];
  recentActivities: ActivityItem[];
  alerts: AlertItem[];
  topClubs: ClubRanking[];
  financeSummary: FinanceSummary;
  beltDistribution: BeltDistributionItem[];
}

export interface MonthlyDataPoint {
  month: string;
  value: number;
  previousValue?: number;
}

export interface ProvinceData {
  province: string;
  provinceCode: string;
  memberCount: number;
  clubCount: number;
  percentage: number;
}

export interface ActivityItem {
  id: string;
  type: "member_join" | "exam_complete" | "tournament" | "club_approved" | "certificate_issued" | "document";
  title: string;
  description: string;
  timestamp: string;
  actor?: string;
}

export interface AlertItem {
  id: string;
  type: "warning" | "info" | "urgent" | "success";
  title: string;
  description: string;
  actionLabel?: string;
  actionHref?: string;
  timestamp: string;
}

export interface ClubRanking {
  rank: number;
  clubId: string;
  clubName: string;
  province: string;
  memberCount: number;
  activityScore: number;
  trend: "up" | "down" | "stable";
}

export interface FinanceSummary {
  totalIncome: number;
  totalExpense: number;
  balance: number;
  monthlyIncome: MonthlyDataPoint[];
  incomeByCategory: { category: string; amount: number; percentage: number }[];
}

export interface BeltDistributionItem {
  belt: BeltLevel;
  count: number;
  percentage: number;
}
