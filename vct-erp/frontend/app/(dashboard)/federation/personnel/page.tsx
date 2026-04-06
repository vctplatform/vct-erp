"use client";

import { UserCog, Shield, GraduationCap, Award, MapPin, Phone, Mail, Calendar, Plus } from "lucide-react";

import { PageHeader, SectionPanel, StatusBadge, GlassButton, SearchInput, Tabs, ProgressBar } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_PERSONNEL } from "@/lib/federation/mock-data";
import { BELT_CONFIG } from "@/lib/federation/constants";
import { cn } from "@/lib/utils";
import { useState } from "react";

const ROLE_MAP = {
  coach: { label: "Huấn luyện viên", color: "text-cyan-600 dark:text-cyan-400", bg: "bg-cyan-500/10" },
  referee: { label: "Trọng tài", color: "text-amber-600 dark:text-amber-400", bg: "bg-amber-500/10" },
  board_member: { label: "Ban Chấp hành", color: "text-indigo-600 dark:text-indigo-400", bg: "bg-indigo-500/10" },
  secretary: { label: "Thư ký", color: "text-emerald-600 dark:text-emerald-400", bg: "bg-emerald-500/10" },
  treasurer: { label: "Thủ quỹ", color: "text-rose-600 dark:text-rose-400", bg: "bg-rose-500/10" },
};

export default function PersonnelPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");

  const filtered = MOCK_PERSONNEL.filter((p) => {
    if (activeTab !== "all" && p.role !== activeTab) return false;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      return p.fullName.toLowerCase().includes(q) || p.province.toLowerCase().includes(q);
    }
    return true;
  });

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Tổ chức nhân sự"
        title="Nhân sự & HLV"
        description="Quản lý đội ngũ Huấn luyện viên, Trọng tài, Ban chấp hành Liên đoàn toàn quốc."
        actions={
          <GlassButton variant="primary" size="sm">
            <Plus className="size-4" />
            Thêm nhân sự
          </GlassButton>
        }
      />

      <div className="stagger-children grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={GraduationCap} label="Huấn luyện viên" value={3_456} color="cyan" />
        <StatCard icon={Shield} label="Trọng tài" value={892} color="amber" />
        <StatCard icon={UserCog} label="Ban Chấp hành" value={45} color="indigo" />
        <StatCard icon={Award} label="Chứng chỉ Quốc gia" value={2_840} color="emerald" />
      </div>

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          tabs={[
            { id: "all", label: "Tất cả", count: MOCK_PERSONNEL.length },
            { id: "coach", label: "HLV", count: MOCK_PERSONNEL.filter((p) => p.role === "coach").length },
            { id: "referee", label: "Trọng tài", count: MOCK_PERSONNEL.filter((p) => p.role === "referee").length },
            { id: "board_member", label: "Ban CH", count: MOCK_PERSONNEL.filter((p) => p.role === "board_member").length },
          ]}
          activeTab={activeTab}
          onChange={(id) => setActiveTab(id)}
        />
        <SearchInput value={searchQuery} onChange={setSearchQuery} placeholder="Tìm theo tên, tỉnh thành..." className="w-56" />
      </div>

      {/* Personnel Cards */}
      <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        {filtered.map((person) => {
          const role = ROLE_MAP[person.role];
          const belt = BELT_CONFIG[person.beltLevel];
          return (
            <div
              key={person.id}
              className="group relative overflow-hidden rounded-2xl border border-[var(--color-border)] bg-[var(--color-panel)] p-6 transition-all duration-300 hover:-translate-y-1 dark:border-white/6 dark:bg-white/[0.025] dark:hover:border-cyan-500/15 dark:hover:shadow-[0_0_30px_rgba(6,182,212,0.08)]"
            >
              <div className="flex items-start gap-4">
                {/* Avatar */}
                <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-500/20 to-indigo-500/10 text-lg font-bold text-cyan-700 dark:text-cyan-300">
                  {person.fullName.split(" ").pop()?.charAt(0)}
                </div>
                <div className="min-w-0 flex-1">
                  <h3 className="truncate font-semibold text-[var(--color-ink)] dark:text-white">
                    {person.fullName}
                  </h3>
                  <div className="mt-1 flex items-center gap-2">
                    <span className={cn("rounded-lg px-2 py-0.5 text-[0.65rem] font-medium", role.bg, role.color)}>
                      {role.label}
                    </span>
                    <StatusBadge status={person.status} />
                  </div>
                </div>
              </div>

              {/* Belt + Province */}
              <div className="mt-4 flex items-center gap-4 text-sm">
                <div className="flex items-center gap-2">
                  <span className="h-4 w-4 rounded-md border border-[var(--color-border)] dark:border-white/12" style={{ backgroundColor: belt.color }} />
                  <span className="text-xs text-[var(--color-ink-soft)] dark:text-white/60">{belt.label}</span>
                </div>
                <span className="flex items-center gap-1 text-xs text-[var(--color-ink-soft)] dark:text-white/45">
                  <MapPin className="size-3" />
                  {person.province}
                </span>
              </div>

              {/* Certifications */}
              <div className="mt-3 flex flex-wrap gap-1.5">
                {person.certifications.map((cert) => (
                  <span key={cert} className="inline-flex items-center gap-1 rounded-lg border border-emerald-500/15 bg-emerald-500/5 px-2 py-0.5 text-[0.65rem] font-medium text-emerald-700 dark:border-emerald-400/10 dark:text-emerald-400">
                    <Award className="size-3" />
                    {cert}
                  </span>
                ))}
              </div>

              {/* Specialization */}
              {person.specialization && (
                <p className="mt-3 text-xs text-[var(--color-ink-soft)] dark:text-white/40">
                  Chuyên môn: <span className="text-[var(--color-ink)] dark:text-white/70">{person.specialization}</span>
                </p>
              )}

              {/* Contact Footer */}
              <div className="mt-4 flex items-center gap-4 border-t border-[var(--color-border)] pt-3 text-xs text-[var(--color-ink-soft)] dark:border-white/5 dark:text-white/40">
                <span className="flex items-center gap-1">
                  <Phone className="size-3" />
                  {person.phone}
                </span>
                <span className="flex items-center gap-1">
                  <Calendar className="size-3" />
                  Từ {new Date(person.joinDate).getFullYear()}
                </span>
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}
