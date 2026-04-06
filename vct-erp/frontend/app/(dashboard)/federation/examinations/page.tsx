"use client";

import { GraduationCap, Calendar, MapPin, Users, CheckCircle, XCircle, Clock, Plus, ArrowRight } from "lucide-react";

import { PageHeader, SectionPanel, StatusBadge, GlassButton, FederationDataTable, ProgressBar } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_EXAMINATIONS } from "@/lib/federation/mock-data";
import { BELT_CONFIG } from "@/lib/federation/constants";
import { cn } from "@/lib/utils";

export default function ExaminationsPage() {
  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Nghiệp vụ chuyên môn"
        title="Thi & Thăng đai"
        description="Quản lý kỳ thi thăng đai toàn quốc — Lịch thi, kết quả, thống kê tỷ lệ đạt."
        actions={
          <GlassButton variant="primary" size="sm">
            <Plus className="size-4" />
            Tạo kỳ thi mới
          </GlassButton>
        }
      />

      <div className="stagger-children grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={GraduationCap} label="Tổng kỳ thi" value={24} color="cyan" />
        <StatCard icon={Users} label="Thí sinh năm nay" value={4_526} color="emerald" />
        <StatCard icon={CheckCircle} label="Tỷ lệ đạt" value="82.4%" color="amber" />
        <StatCard icon={Calendar} label="Kỳ thi sắp tới" value={2} color="indigo" />
      </div>

      {/* Examination Timeline */}
      <div className="space-y-5">
        {MOCK_EXAMINATIONS.map((exam) => {
          const targetBelt = BELT_CONFIG[exam.targetBelt];
          const currentBelt = BELT_CONFIG[exam.beltLevel];
          const passRate = exam.candidateCount > 0 ? ((exam.passedCount / exam.candidateCount) * 100).toFixed(1) : "—";

          return (
            <SectionPanel
              key={exam.id}
              title={exam.title}
              kicker={exam.code}
              actions={<StatusBadge status={exam.status} size="md" />}
            >
              <div className="grid gap-5 lg:grid-cols-[1.5fr_1fr]">
                {/* Left — Details */}
                <div className="space-y-4">
                  <div className="flex flex-wrap items-center gap-4 text-sm text-[var(--color-ink-soft)] dark:text-white/50">
                    <span className="flex items-center gap-1.5">
                      <Calendar className="size-3.5" />
                      {new Date(exam.date).toLocaleDateString("vi-VN")}
                      {exam.endDate && ` — ${new Date(exam.endDate).toLocaleDateString("vi-VN")}`}
                    </span>
                    <span className="flex items-center gap-1.5">
                      <MapPin className="size-3.5" />
                      {exam.location}
                    </span>
                    <span className="flex items-center gap-1.5">
                      <Users className="size-3.5" />
                      {exam.candidateCount.toLocaleString("vi-VN")} thí sinh
                    </span>
                  </div>

                  {/* Belt Progression */}
                  <div className="flex items-center gap-3 rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-3 dark:border-white/5 dark:bg-white/[0.02]">
                    <div className="flex items-center gap-2">
                      <span className="h-6 w-6 rounded-lg border border-[var(--color-border)] dark:border-white/12" style={{ backgroundColor: currentBelt.color }} />
                      <span className="text-sm font-medium text-[var(--color-ink)] dark:text-white">{currentBelt.label}</span>
                    </div>
                    <ArrowRight className="size-4 text-cyan-500" />
                    <div className="flex items-center gap-2">
                      <span className="h-6 w-6 rounded-lg border border-[var(--color-border)] dark:border-white/12" style={{ backgroundColor: targetBelt.color }} />
                      <span className="text-sm font-medium text-[var(--color-ink)] dark:text-white">{targetBelt.label}</span>
                    </div>
                  </div>

                  {/* Judges */}
                  {exam.judges.length > 0 && (
                    <div>
                      <p className="mb-2 text-[0.62rem] uppercase tracking-[0.25em] text-[var(--color-ink-soft)] dark:text-white/35">
                        Hội đồng Giám khảo
                      </p>
                      <div className="flex flex-wrap gap-2">
                        {exam.judges.map((j) => (
                          <span key={j.id} className="inline-flex items-center gap-2 rounded-lg border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-3 py-1.5 text-xs dark:border-white/6 dark:bg-white/[0.03]">
                            <span className="font-medium text-[var(--color-ink)] dark:text-white">{j.name}</span>
                            <span className="text-[var(--color-ink-soft)] dark:text-white/40">({j.role === "chief" ? "Chủ khảo" : "Uỷ viên"})</span>
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </div>

                {/* Right — Results */}
                <div className="rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4 dark:border-white/5 dark:bg-white/[0.02]">
                  <p className="mb-3 text-[0.62rem] uppercase tracking-[0.25em] text-[var(--color-ink-soft)] dark:text-white/35">
                    Kết quả
                  </p>
                  {exam.status === "completed" ? (
                    <div className="space-y-4">
                      <div className="text-center">
                        <p className="text-4xl font-bold text-[var(--color-ink)] dark:text-white">
                          {passRate}<span className="text-lg text-[var(--color-ink-soft)] dark:text-white/50">%</span>
                        </p>
                        <p className="mt-1 text-xs text-[var(--color-ink-soft)] dark:text-white/40">Tỷ lệ đạt</p>
                      </div>
                      <ProgressBar value={exam.passedCount} max={exam.candidateCount} color="emerald" />
                      <div className="grid grid-cols-3 gap-3 text-center">
                        <div>
                          <p className="text-lg font-bold text-[var(--color-ink)] dark:text-white">{exam.candidateCount}</p>
                          <p className="text-[0.6rem] text-[var(--color-ink-soft)] dark:text-white/35">Dự thi</p>
                        </div>
                        <div>
                          <p className="text-lg font-bold text-emerald-600 dark:text-emerald-400">{exam.passedCount}</p>
                          <p className="text-[0.6rem] text-[var(--color-ink-soft)] dark:text-white/35">Đạt</p>
                        </div>
                        <div>
                          <p className="text-lg font-bold text-rose-600 dark:text-rose-400">{exam.failedCount}</p>
                          <p className="text-[0.6rem] text-[var(--color-ink-soft)] dark:text-white/35">Rớt</p>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="flex flex-col items-center justify-center py-6 text-center">
                      <Clock className="size-8 text-[var(--color-ink-soft)] dark:text-white/25" />
                      <p className="mt-3 text-sm font-medium text-[var(--color-ink)] dark:text-white">Chưa có kết quả</p>
                      <p className="mt-1 text-xs text-[var(--color-ink-soft)] dark:text-white/40">
                        Kỳ thi {exam.status === "scheduled" ? "sẽ diễn ra" : "đang diễn ra"}
                      </p>
                    </div>
                  )}
                </div>
              </div>
            </SectionPanel>
          );
        })}
      </div>
    </section>
  );
}
