"use client";

import { BarChart3, FileDown, Printer, Calendar, TrendingUp, Users, Building2, Trophy, Wallet, GraduationCap, FileText, PieChart } from "lucide-react";

import { PageHeader, SectionPanel, GlassButton } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { cn } from "@/lib/utils";

const REPORT_TEMPLATES = [
  { id: "R-001", title: "Báo cáo Tổng hợp Hoạt động", icon: BarChart3, description: "Báo cáo tổng quan toàn bộ hoạt động liên đoàn theo kỳ (tháng/quý/năm).", color: "cyan", formats: ["PDF", "Excel"] },
  { id: "R-002", title: "Thống kê Hội viên Toàn quốc", icon: Users, description: "Chi tiết thành viên theo tỉnh, đai cấp, độ tuổi, giới tính.", color: "emerald", formats: ["PDF", "Excel", "CSV"] },
  { id: "R-003", title: "Báo cáo CLB & Võ đường", icon: Building2, description: "Danh sách CLB, xếp hạng hoạt động, số lượng hội viên theo vùng.", color: "indigo", formats: ["PDF", "Excel"] },
  { id: "R-004", title: "Kết quả Giải đấu", icon: Trophy, description: "Bảng tổng sắp huy chương, thống kê VĐV theo tỉnh và nội dung.", color: "amber", formats: ["PDF"] },
  { id: "R-005", title: "Báo cáo Tài chính", icon: Wallet, description: "Thu chi ngân sách, cơ cấu nguồn thu, so sánh cùng kỳ.", color: "rose", formats: ["PDF", "Excel"] },
  { id: "R-006", title: "Thống kê Thi thăng đai", icon: GraduationCap, description: "Kết quả thi đai theo kỳ, tỷ lệ đạt theo vùng và cấp đai.", color: "sky", formats: ["PDF", "Excel"] },
  { id: "R-007", title: "Báo cáo Nhân sự HLV", icon: Users, description: "Danh sách HLV, Trọng tài theo cấp chứng chỉ và vùng miền.", color: "indigo", formats: ["PDF", "CSV"] },
  { id: "R-008", title: "Phân tích Xu hướng", icon: TrendingUp, description: "Xu hướng tăng trưởng hội viên, hoạt động CLB, doanh thu qua các năm.", color: "cyan", formats: ["PDF"] },
];

const COLOR_MAP = {
  cyan: { iconBg: "bg-cyan-500/10 dark:bg-cyan-500/12", iconText: "text-cyan-600 dark:text-cyan-400" },
  emerald: { iconBg: "bg-emerald-500/10 dark:bg-emerald-500/12", iconText: "text-emerald-600 dark:text-emerald-400" },
  indigo: { iconBg: "bg-indigo-500/10 dark:bg-indigo-500/12", iconText: "text-indigo-600 dark:text-indigo-400" },
  amber: { iconBg: "bg-amber-500/10 dark:bg-amber-500/12", iconText: "text-amber-600 dark:text-amber-400" },
  rose: { iconBg: "bg-rose-500/10 dark:bg-rose-500/12", iconText: "text-rose-600 dark:text-rose-400" },
  sky: { iconBg: "bg-sky-500/10 dark:bg-sky-500/12", iconText: "text-sky-600 dark:text-sky-400" },
};

export default function ReportsPage() {
  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Phân tích & Báo cáo"
        title="Báo cáo"
        description="Tạo và xuất báo cáo tổng hợp — Phân tích dữ liệu hoạt động liên đoàn theo nhiều chiều."
      />

      <div className="stagger-children grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={BarChart3} label="Mẫu báo cáo" value={REPORT_TEMPLATES.length} color="cyan" />
        <StatCard icon={FileDown} label="Đã xuất tháng này" value={34} color="emerald" />
        <StatCard icon={Calendar} label="Báo cáo định kỳ" value={4} color="amber" />
        <StatCard icon={PieChart} label="Dashboard tùy chỉnh" value={2} color="indigo" />
      </div>

      {/* Report Templates Grid */}
      <SectionPanel title="Mẫu Báo cáo" subtitle="Chọn mẫu để tạo báo cáo mới" kicker="Thư viện">
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {REPORT_TEMPLATES.map((report) => {
            const ReportIcon = report.icon;
            const colorCfg = COLOR_MAP[report.color as keyof typeof COLOR_MAP];
            return (
              <div
                key={report.id}
                className="group cursor-pointer rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-5 transition-all duration-200 hover:-translate-y-0.5 hover:border-cyan-500/20 dark:border-white/5 dark:bg-white/[0.02] dark:hover:border-cyan-500/15 dark:hover:shadow-[0_0_20px_rgba(6,182,212,0.06)]"
              >
                <div className="flex items-start gap-3">
                  <div className={cn("rounded-xl p-2.5", colorCfg.iconBg)}>
                    <ReportIcon className={cn("size-5", colorCfg.iconText)} />
                  </div>
                  <div className="min-w-0 flex-1">
                    <h4 className="font-semibold text-[var(--color-ink)] dark:text-white">{report.title}</h4>
                    <p className="mt-1.5 text-xs leading-relaxed text-[var(--color-ink-soft)] dark:text-white/45">
                      {report.description}
                    </p>
                  </div>
                </div>

                <div className="mt-4 flex items-center justify-between">
                  <div className="flex gap-1.5">
                    {report.formats.map((fmt) => (
                      <span key={fmt} className="rounded border border-[var(--color-border)] px-1.5 py-0.5 text-[0.6rem] font-medium text-[var(--color-ink-soft)] dark:border-white/6 dark:text-white/40">
                        {fmt}
                      </span>
                    ))}
                  </div>
                  <div className="flex items-center gap-1.5">
                    <button className="rounded-lg p-1.5 text-[var(--color-ink-soft)] transition-colors hover:bg-[var(--color-border)] dark:text-white/35 dark:hover:bg-white/5">
                      <Printer className="size-3.5" />
                    </button>
                    <button className="rounded-lg p-1.5 text-[var(--color-ink-soft)] transition-colors hover:bg-[var(--color-border)] dark:text-white/35 dark:hover:bg-white/5">
                      <FileDown className="size-3.5" />
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </SectionPanel>
    </section>
  );
}
