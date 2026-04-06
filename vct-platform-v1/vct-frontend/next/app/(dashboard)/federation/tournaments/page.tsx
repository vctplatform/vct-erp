"use client";

import { useState } from "react";
import { Trophy, Calendar, MapPin, Users, Medal, Swords, Plus, Radio, Filter } from "lucide-react";
import { PageHeader, SectionPanel, StatusBadge, FederationDataTable, GlassButton, Tabs, SearchInput } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_TOURNAMENTS } from "@/lib/federation/mock-data";
import { cn } from "@/lib/utils";

export default function TournamentsPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");

  const filtered = MOCK_TOURNAMENTS.filter((t) => {
    if (activeTab === "live" && t.status !== "live") return false;
    if (activeTab === "upcoming" && t.status !== "upcoming" && t.status !== "registration") return false;
    if (searchQuery) return t.name.toLowerCase().includes(searchQuery.toLowerCase());
    return true;
  });

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Nghiệp vụ chuyên môn"
        title="Giải đấu"
        description="Quản lý giải đấu Quốc gia, khu vực — Bảng thi đấu, kết quả, bảng tổng sắp huy chương."
        actions={
          <GlassButton variant="primary" size="sm">
            <Plus className="size-4" />
            Tạo giải đấu
          </GlassButton>
        }
      />

      <div className="stagger-fed grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={Trophy} label="Tổng giải đấu" value={156} color="amber" />
        <StatCard icon={Radio} label="Đang diễn ra" value={1} color="rose" />
        <StatCard icon={Users} label="Tổng VĐV" value={8920} color="cyan" />
        <StatCard icon={Medal} label="Huy chương đã trao" value={1247} color="emerald" />
      </div>

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          tabs={[
            { id: "all", label: "Tất cả", count: MOCK_TOURNAMENTS.length },
            { id: "live", label: "🔴 Đang diễn ra", count: MOCK_TOURNAMENTS.filter((t) => t.status === "live").length },
            { id: "upcoming", label: "Sắp tới", count: MOCK_TOURNAMENTS.filter((t) => t.status === "upcoming" || t.status === "registration").length },
          ]}
          activeTab={activeTab}
          onChange={setActiveTab}
        />
        <SearchInput value={searchQuery} onChange={setSearchQuery} placeholder="Tìm giải đấu..." className="w-56" />
      </div>

      {/* Tournament cards */}
      <div className="space-y-5">
        {filtered.map((t) => (
          <div
            key={t.id}
            className={cn(
              "group overflow-hidden rounded-2xl border bg-white/[0.025] transition-all duration-300 neon-hover",
              t.status === "live"
                ? "border-rose-500/30 shadow-[0_0_35px_rgba(225,29,72,0.12)]"
                : "border-white/6"
            )}
          >
            {/* Live indicator band */}
            {t.status === "live" && (
              <div className="flex items-center gap-2 border-b border-rose-400/10 bg-gradient-to-r from-rose-500/10 to-transparent px-6 py-2.5 text-sm font-medium text-rose-400">
                <Radio className="size-4 animate-pulse" />
                LIVE — Đang diễn ra
              </div>
            )}

            <div className="p-6">
              <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
                <div className="flex items-start gap-4">
                  <div className={cn(
                    "rounded-xl p-3",
                    t.status === "live"
                      ? "bg-rose-500/15"
                      : "bg-gradient-to-br from-amber-500/15 to-orange-500/10"
                  )}>
                    <Trophy className={cn(
                      "size-6",
                      t.status === "live" ? "text-rose-400" : "text-amber-400"
                    )} />
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-white">
                      {t.name}
                    </h3>
                    <div className="mt-2 flex flex-wrap items-center gap-4 text-sm text-white/50">
                      <span className="flex items-center gap-1.5">
                        <Calendar className="size-3.5" />
                        {new Date(t.startDate).toLocaleDateString("vi-VN")} — {new Date(t.endDate).toLocaleDateString("vi-VN")}
                      </span>
                      <span className="flex items-center gap-1.5">
                        <MapPin className="size-3.5" />
                        {t.location}
                      </span>
                      <StatusBadge status={t.status} />
                    </div>

                    {/* Categories */}
                    {t.categories.length > 0 && (
                      <div className="mt-3 flex flex-wrap gap-2">
                        {t.categories.map((cat) => (
                          <span
                            key={cat.id}
                            className="inline-flex items-center gap-1 rounded-lg border border-white/6 bg-white/[0.03] px-2 py-1 text-xs"
                          >
                            <Swords className="size-3 text-white/40" />
                            <span className="text-white/80">{cat.name}</span>
                            <span className="text-white/40">({cat.registeredCount})</span>
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                <div className="flex items-center gap-6 rounded-xl border border-white/5 bg-white/[0.02] px-5 py-3">
                  <div className="text-center">
                    <p className="text-2xl font-bold text-white">{t.teamCount}</p>
                    <p className="text-[0.65rem] uppercase tracking-wider text-white/35">Đoàn</p>
                  </div>
                  <div className="h-10 w-px bg-white/6" />
                  <div className="text-center">
                    <p className="text-2xl font-bold text-white">{t.athleteCount}</p>
                    <p className="text-[0.65rem] uppercase tracking-wider text-white/35">VĐV</p>
                  </div>
                </div>
              </div>

              {/* Medal Table */}
              {t.medalTable.length > 0 && (
                <div className="mt-5 rounded-xl border border-white/5 bg-white/[0.02] p-4 text-center">
                  <p className="mb-3 text-[0.62rem] uppercase tracking-[0.25em] text-white/35">
                    Bảng tổng sắp Huy chương
                  </p>
                  <FederationDataTable
                    compact
                    columns={[
                      { key: "rank", label: "#", width: "40px" },
                      { key: "team", label: "Đoàn" },
                      { key: "gold", label: "🥇", align: "center" },
                      { key: "silver", label: "🥈", align: "center" },
                      { key: "bronze", label: "🥉", align: "center" },
                      { key: "total", label: "Tổng", align: "right" },
                    ]}
                    rows={t.medalTable.map((m, i) => ({
                      rank: (
                        <span className={cn(
                          "inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold",
                          i === 0 && "bg-amber-500/15 text-amber-400",
                          i === 1 && "bg-slate-300/20 text-slate-300",
                          i === 2 && "bg-orange-500/15 text-orange-400",
                          i > 2 && "text-white/50"
                        )}>
                          {i + 1}
                        </span>
                      ),
                      team: <span className="font-medium text-white/80">{m.teamName} — {m.province}</span>,
                      gold: <span className="font-bold text-amber-400">{m.gold}</span>,
                      silver: <span className="font-medium text-slate-400">{m.silver}</span>,
                      bronze: <span className="font-medium text-orange-400">{m.bronze}</span>,
                      total: <span className="font-bold text-white">{m.total}</span>,
                    }))}
                  />
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
