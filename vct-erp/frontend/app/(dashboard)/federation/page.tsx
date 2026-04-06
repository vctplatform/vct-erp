"use client";

import {
  Users, Building2, UserCog, Trophy, AlertTriangle, CheckCircle,
  Info, TrendingUp, TrendingDown, Minus,
  ArrowUpRight, Clock, GraduationCap, Award, Activity, Wallet,
} from "lucide-react";
import {
  AreaChart, Area, BarChart, Bar, PieChart, Pie, Cell,
  XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid,
} from "recharts";

import { StatCard } from "@/components/federation/ui/stat-card";
import { SectionPanel, StatusBadge, FederationDataTable, PageHeader } from "@/components/federation/ui/shared";
import { MOCK_DASHBOARD } from "@/lib/federation/mock-data";
import { BELT_CONFIG } from "@/lib/federation/constants";
import { cn } from "@/lib/utils";

const ALERT_ICONS = {
  urgent: AlertTriangle,
  warning: AlertTriangle,
  info: Info,
  success: CheckCircle,
};

const ALERT_COLORS = {
  urgent: "border-rose-500/20 bg-rose-500/5 dark:border-rose-400/12 dark:bg-rose-500/[0.06]",
  warning: "border-amber-500/20 bg-amber-500/5 dark:border-amber-400/12 dark:bg-amber-500/[0.06]",
  info: "border-sky-500/20 bg-sky-500/5 dark:border-sky-400/12 dark:bg-sky-500/[0.06]",
  success: "border-emerald-500/20 bg-emerald-500/5 dark:border-emerald-400/12 dark:bg-emerald-500/[0.06]",
};

const ALERT_ICON_COLORS = {
  urgent: "text-rose-500 dark:text-rose-400",
  warning: "text-amber-500 dark:text-amber-400",
  info: "text-sky-500 dark:text-sky-400",
  success: "text-emerald-500 dark:text-emerald-400",
};

const CHART_COLORS = ["#06B6D4", "#14B8A6", "#F59E0B", "#8B5CF6", "#EC4899", "#6366F1"];

const ACTIVITY_ICONS = {
  member_join: Users,
  exam_complete: GraduationCap,
  tournament: Trophy,
  club_approved: Building2,
  certificate_issued: Award,
  document: Activity,
};

export default function FederationDashboardPage() {
  const d = MOCK_DASHBOARD;
  const topBelts = d.beltDistribution.slice(0, 8);

  return (
    <section className="space-y-6">
      {/* Header */}
      <PageHeader
        kicker="Liên đoàn Võ Cổ Truyền Việt Nam"
        title="Trung tâm Chỉ huy"
        description="Tổng quan toàn quốc — Theo dõi tình hình hội viên, CLB, giải đấu và hoạt động của liên đoàn trên khắp 63 tỉnh thành."
      />

      {/* KPI Cards */}
      <div className="stagger-children grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          icon={Users}
          label="Tổng hội viên"
          value={d.totalMembers}
          trend={{ direction: "up", value: "+6.8%" }}
          color="cyan"
          sparkline={d.memberGrowth.map((m) => m.value)}
        />
        <StatCard
          icon={Building2}
          label="CLB / Võ đường"
          value={d.totalClubs}
          trend={{ direction: "up", value: "+3.2%" }}
          color="emerald"
          sparkline={[1100, 1120, 1140, 1160, 1180, 1200, 1215, 1225, 1235, 1240, 1245, 1247]}
        />
        <StatCard
          icon={UserCog}
          label="HLV & Trọng tài"
          value={d.totalCoaches + d.totalReferees}
          suffix="người"
          trend={{ direction: "stable", value: "+1.1%" }}
          color="indigo"
        />
        <StatCard
          icon={Trophy}
          label="Giải đấu năm nay"
          value={d.totalTournaments}
          trend={{ direction: "up", value: "+12%" }}
          color="amber"
        />
      </div>

      {/* ── Overview Ribbon ─────────────────────────── */}
      <div className="animate-slide-up grid grid-cols-2 gap-3 rounded-2xl border border-[var(--color-border)] bg-[var(--color-panel)] p-4 dark:border-white/6 dark:bg-white/[0.02] sm:grid-cols-4">
        {[
          { label: "Đang hoạt động", value: "4", icon: Activity, color: "text-emerald-500" },
          { label: "Chờ duyệt", value: "23", icon: Clock, color: "text-amber-500" },
          { label: "Thu Q4", value: "3.2 tỷ", icon: Wallet, color: "text-cyan-500" },
          { label: "Huy chương đã trao", value: "1,247", icon: Award, color: "text-indigo-500" },
        ].map((item) => (
          <div key={item.label} className="flex items-center gap-3 px-2">
            <item.icon className={cn("size-5", item.color)} />
            <div>
              <p className="text-lg font-bold text-[var(--color-ink)] dark:text-white">{item.value}</p>
              <p className="text-[0.65rem] text-[var(--color-ink-soft)] dark:text-white/40">{item.label}</p>
            </div>
          </div>
        ))}
      </div>

      {/* Charts Row */}
      <div className="grid gap-5 xl:grid-cols-[1.2fr_0.8fr]">
        {/* Member Growth Chart */}
        <SectionPanel
          title="Tăng trưởng Hội viên"
          subtitle="So sánh 12 tháng gần nhất với cùng kỳ năm trước"
          kicker="Biểu đồ tăng trưởng"
        >
          <div className="h-[280px]">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={d.memberGrowth} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="growthGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#06B6D4" stopOpacity={0.25} />
                    <stop offset="100%" stopColor="#06B6D4" stopOpacity={0} />
                  </linearGradient>
                  <linearGradient id="prevGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#94A3B8" stopOpacity={0.15} />
                    <stop offset="100%" stopColor="#94A3B8" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" opacity={0.5} />
                <XAxis dataKey="month" tick={{ fontSize: 12, fill: "var(--color-ink-soft)" }} axisLine={false} tickLine={false} />
                <YAxis tick={{ fontSize: 11, fill: "var(--color-ink-soft)" }} axisLine={false} tickLine={false} tickFormatter={(v: number) => `${(v / 1000).toFixed(0)}k`} />
                <Tooltip
                  contentStyle={{
                    background: "var(--color-panel)",
                    border: "1px solid var(--color-border)",
                    borderRadius: "12px",
                    fontSize: "13px",
                    boxShadow: "0 10px 30px rgba(0,0,0,0.15)",
                  }}
                  formatter={(value) => [Number(value).toLocaleString("vi-VN"), ""]}
                />
                <Area type="monotone" dataKey="previousValue" stroke="#94A3B8" strokeWidth={1.5} fill="url(#prevGrad)" dot={false} name="Năm trước" />
                <Area type="monotone" dataKey="value" stroke="#06B6D4" strokeWidth={2.5} fill="url(#growthGrad)" dot={false} name="Năm nay" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </SectionPanel>

        {/* Revenue Breakdown Pie */}
        <SectionPanel
          title="Cơ cấu Nguồn thu"
          subtitle="Phân bổ theo danh mục thu nhập"
          kicker="Tài chính"
        >
          <div className="flex flex-col items-center gap-4">
            <div className="h-[180px] w-[180px]">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={d.financeSummary.incomeByCategory}
                    cx="50%"
                    cy="50%"
                    innerRadius={55}
                    outerRadius={85}
                    paddingAngle={2}
                    dataKey="amount"
                  >
                    {d.financeSummary.incomeByCategory.map((_, i) => (
                      <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />
                    ))}
                  </Pie>
                </PieChart>
              </ResponsiveContainer>
            </div>
            <div className="w-full space-y-2">
              {d.financeSummary.incomeByCategory.map((cat, i) => (
                <div key={cat.category} className="flex items-center justify-between text-sm">
                  <div className="flex items-center gap-2">
                    <span
                      className="size-2.5 rounded-full"
                      style={{ backgroundColor: CHART_COLORS[i % CHART_COLORS.length] }}
                    />
                    <span className="text-[var(--color-ink)] dark:text-white/80">{cat.category}</span>
                  </div>
                  <span className="font-medium text-[var(--color-ink)] dark:text-white">
                    {cat.percentage}%
                  </span>
                </div>
              ))}
            </div>
          </div>
        </SectionPanel>
      </div>

      {/* Alerts + Activity Row */}
      <div className="grid gap-5 xl:grid-cols-2">
        {/* Alerts */}
        <SectionPanel
          title="Cảnh báo & Thông báo"
          kicker="Cần xử lý"
          actions={<StatusBadge status="pending" label={`${d.pendingApprovals} chờ duyệt`} size="md" />}
        >
          <div className="space-y-3">
            {d.alerts.map((alert) => {
              const Icon = ALERT_ICONS[alert.type];
              return (
                <div
                  key={alert.id}
                  className={cn(
                    "flex items-start gap-3 rounded-xl border p-4 transition-all duration-200 hover:-translate-y-0.5",
                    ALERT_COLORS[alert.type]
                  )}
                >
                  <Icon className={cn("mt-0.5 size-4 shrink-0", ALERT_ICON_COLORS[alert.type])} />
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium text-[var(--color-ink)] dark:text-white">
                      {alert.title}
                    </p>
                    <p className="mt-0.5 text-xs text-[var(--color-ink-soft)] dark:text-white/50">
                      {alert.description}
                    </p>
                  </div>
                  {alert.actionLabel && (
                    <button className="shrink-0 rounded-lg bg-[var(--color-canvas-soft)] px-3 py-1.5 text-xs font-medium text-[var(--color-ink)] transition-colors hover:bg-[var(--color-border)] dark:bg-white/6 dark:text-white dark:hover:bg-white/10">
                      {alert.actionLabel}
                    </button>
                  )}
                </div>
              );
            })}
          </div>
        </SectionPanel>

        {/* Recent Activity Timeline */}
        <SectionPanel
          title="Hoạt động Gần đây"
          kicker="Timeline"
        >
          <div className="space-y-4">
            {d.recentActivities.map((activity, i) => {
              const Icon = ACTIVITY_ICONS[activity.type] || Activity;
              return (
                <div key={activity.id} className="flex gap-3">
                  <div className="relative flex flex-col items-center">
                    <div className="rounded-lg bg-cyan-500/10 p-2 dark:bg-cyan-500/12">
                      <Icon className="size-3.5 text-cyan-600 dark:text-cyan-400" />
                    </div>
                    {i < d.recentActivities.length - 1 && (
                      <div className="mt-2 w-px flex-1 bg-[var(--color-border)] dark:bg-white/6" />
                    )}
                  </div>
                  <div className="pb-4">
                    <p className="text-sm font-medium text-[var(--color-ink)] dark:text-white">
                      {activity.title}
                    </p>
                    <p className="mt-0.5 text-xs text-[var(--color-ink-soft)] dark:text-white/50">
                      {activity.description}
                    </p>
                    <p className="mt-1 flex items-center gap-1 text-[0.65rem] text-[var(--color-ink-soft)] dark:text-white/35">
                      <Clock className="size-3" />
                      {new Date(activity.timestamp).toLocaleDateString("vi-VN")}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </SectionPanel>
      </div>

      {/* Top Clubs + Belt Distribution */}
      <div className="grid gap-5 xl:grid-cols-[1.2fr_0.8fr]">
        <SectionPanel
          title="Xếp hạng CLB Hoạt động"
          subtitle="Top CLB có điểm hoạt động cao nhất năm"
          kicker="Bảng xếp hạng"
        >
          <FederationDataTable
            columns={[
              { key: "rank", label: "#", width: "50px" },
              { key: "club", label: "Câu lạc bộ" },
              { key: "province", label: "Tỉnh/TP" },
              { key: "members", label: "Hội viên", align: "right" },
              { key: "score", label: "Điểm", align: "right" },
              { key: "trend", label: "Xu hướng", align: "center" },
            ]}
            rows={d.topClubs.map((club) => ({
              rank: (
                <span className={cn(
                  "inline-flex h-7 w-7 items-center justify-center rounded-full text-xs font-bold",
                  club.rank === 1 && "bg-amber-500/15 text-amber-600 dark:text-amber-400",
                  club.rank === 2 && "bg-slate-300/20 text-slate-600 dark:text-slate-300",
                  club.rank === 3 && "bg-orange-500/15 text-orange-600 dark:text-orange-400",
                  club.rank > 3 && "text-[var(--color-ink-soft)] dark:text-white/50"
                )}>
                  {club.rank}
                </span>
              ),
              club: (
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-500/15 to-teal-500/10 text-xs font-bold text-cyan-600 dark:text-cyan-400">
                    {club.clubName.charAt(0)}
                  </div>
                  <span className="font-medium">{club.clubName}</span>
                </div>
              ),
              province: club.province,
              members: club.memberCount.toLocaleString("vi-VN"),
              score: (
                <span className="font-mono font-semibold text-cyan-600 dark:text-cyan-400">
                  {club.activityScore}
                </span>
              ),
              trend: (
                <span className={cn(
                  "inline-flex items-center gap-1",
                  club.trend === "up" && "text-emerald-500",
                  club.trend === "down" && "text-rose-500",
                  club.trend === "stable" && "text-slate-400"
                )}>
                  {club.trend === "up" && <TrendingUp className="size-3.5" />}
                  {club.trend === "down" && <TrendingDown className="size-3.5" />}
                  {club.trend === "stable" && <Minus className="size-3.5" />}
                </span>
              ),
            }))}
          />
        </SectionPanel>

        <SectionPanel
          title="Phân bổ Đai cấp"
          subtitle="Tỷ lệ hội viên theo cấp đai hiện tại"
          kicker="Thống kê đai"
        >
          <div className="space-y-3">
            {topBelts.map((item) => {
              const belt = BELT_CONFIG[item.belt];
              return (
                <div key={item.belt} className="flex items-center gap-3">
                  <div
                    className="h-5 w-5 shrink-0 rounded-md border border-[var(--color-border)] dark:border-white/12"
                    style={{ backgroundColor: belt.color }}
                  />
                  <span className="w-20 shrink-0 text-sm text-[var(--color-ink)] dark:text-white/80">
                    {belt.label}
                  </span>
                  <div className="flex-1">
                    <div className="h-2 overflow-hidden rounded-full bg-[var(--color-canvas-soft)] dark:bg-white/6">
                      <div
                        className="h-full rounded-full bg-gradient-to-r from-cyan-500 to-teal-500 transition-all duration-1000"
                        style={{ width: `${Math.min(item.percentage * 3.5, 100)}%` }}
                      />
                    </div>
                  </div>
                  <span className="w-16 text-right text-xs font-medium text-[var(--color-ink-soft)] dark:text-white/50">
                    {item.count.toLocaleString("vi-VN")}
                  </span>
                  <span className="w-12 text-right text-xs text-[var(--color-ink-soft)] dark:text-white/40">
                    {item.percentage}%
                  </span>
                </div>
              );
            })}
          </div>
        </SectionPanel>
      </div>

      {/* Province Distribution (Bar Chart) */}
      <SectionPanel
        title="Phân bố CLB theo Tỉnh/Thành"
        subtitle="Top 10 tỉnh thành có nhiều CLB hoạt động nhất"
        kicker="Bản đồ phân bố"
      >
        <div className="h-[300px]">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={d.clubDistribution.slice(0, 10)} margin={{ top: 5, right: 5, left: -10, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" opacity={0.4} />
              <XAxis
                dataKey="province"
                tick={{ fontSize: 11, fill: "var(--color-ink-soft)" }}
                axisLine={false}
                tickLine={false}
                interval={0}
                angle={-20}
                textAnchor="end"
                height={60}
              />
              <YAxis tick={{ fontSize: 11, fill: "var(--color-ink-soft)" }} axisLine={false} tickLine={false} />
              <Tooltip
                contentStyle={{
                  background: "var(--color-panel)",
                  border: "1px solid var(--color-border)",
                  borderRadius: "12px",
                  fontSize: "13px",
                  boxShadow: "0 10px 30px rgba(0,0,0,0.15)",
                }}
              />
              <Bar dataKey="clubCount" name="Số CLB" fill="#06B6D4" radius={[6, 6, 0, 0]} />
              <Bar dataKey="memberCount" name="Hội viên" fill="#14B8A6" radius={[6, 6, 0, 0]} opacity={0.5} />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </SectionPanel>
    </section>
  );
}
