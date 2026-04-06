"use client";

import { useState } from "react";
import { Building2, Star, MapPin, Users, UserCog, Plus, Search, Calendar, Phone, Mail, ChevronRight } from "lucide-react";
import Link from "next/link";

import { PageHeader, SectionPanel, StatusBadge, GlassButton, Tabs, SearchInput, Pagination } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_CLUBS, MOCK_DASHBOARD } from "@/lib/federation/mock-data";
import { cn } from "@/lib/utils";

const PAGE_SIZE = 6;

function StarRating({ rating }: { rating: number }) {
  return (
    <div className="flex items-center gap-0.5">
      {[1, 2, 3, 4, 5].map((i) => (
        <Star
          key={i}
          className={cn(
            "size-3.5",
            i <= rating
              ? "fill-amber-400 text-amber-400"
              : "text-white/10"
          )}
        />
      ))}
    </div>
  );
}

export default function ClubsPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);

  const filteredClubs = MOCK_CLUBS.filter((c) => {
    if (activeTab === "active" && c.status !== "active") return false;
    if (activeTab === "pending" && c.status !== "pending_approval") return false;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      return c.name.toLowerCase().includes(q) || c.province.toLowerCase().includes(q) || c.code.toLowerCase().includes(q);
    }
    return true;
  });

  const totalPages = Math.ceil(filteredClubs.length / PAGE_SIZE);
  const paginatedClubs = filteredClubs.slice((currentPage - 1) * PAGE_SIZE, currentPage * PAGE_SIZE);

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Quản lý tổ chức"
        title="CLB / Võ đường"
        description="Quản lý toàn bộ các Câu lạc bộ, Võ đường trực thuộc Liên đoàn trên 63 tỉnh thành."
        actions={
          <GlassButton variant="primary" size="sm">
            <Plus className="size-4" />
            Đăng ký CLB mới
          </GlassButton>
        }
      />

      <div className="stagger-fed grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={Building2} label="Tổng CLB" value={MOCK_DASHBOARD.totalClubs} color="emerald" trend={{ direction: "up", value: "+42" }} />
        <StatCard icon={Building2} label="Đang hoạt động" value={1215} color="cyan" />
        <StatCard icon={Building2} label="Chờ phê duyệt" value={23} color="amber" />
        <StatCard icon={Users} label="Tổng hội viên CLB" value={MOCK_DASHBOARD.totalMembers} color="sky" />
      </div>

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          tabs={[
            { id: "all", label: "Tất cả", count: MOCK_CLUBS.length },
            { id: "active", label: "Hoạt động", count: MOCK_CLUBS.filter((c) => c.status === "active").length },
            { id: "pending", label: "Chờ duyệt", count: MOCK_CLUBS.filter((c) => c.status === "pending_approval").length },
          ]}
          activeTab={activeTab}
          onChange={(id) => { setActiveTab(id); setCurrentPage(1); }}
        />
        <SearchInput value={searchQuery} onChange={(v) => { setSearchQuery(v); setCurrentPage(1); }} placeholder="Tìm CLB, tỉnh thành..." className="w-56" />
      </div>

      {/* Clubs Card Grid */}
      <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        {paginatedClubs.map((club) => (
          <div
            key={club.id}
            className="group relative overflow-hidden rounded-2xl border border-white/6 bg-white/[0.025] p-6 transition-all duration-300 neon-hover"
          >
            <div className="relative">
              {/* Header */}
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-3">
                  <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-500/20 to-teal-500/10 text-lg font-bold text-cyan-300">
                    {club.name.charAt(0)}
                  </div>
                  <div className="min-w-0">
                    <h3 className="truncate font-semibold text-white">{club.name}</h3>
                    <p className="font-mono text-[0.65rem] text-white/40">{club.code}</p>
                  </div>
                </div>
                <StatusBadge status={club.status} />
              </div>

              {/* Rating + Location */}
              <div className="mt-4 flex items-center gap-3">
                <StarRating rating={club.rating} />
                <span className="flex items-center gap-1 text-xs text-white/45">
                  <MapPin className="size-3" />
                  {club.province}
                </span>
              </div>

              {/* Stats Grid */}
              <div className="mt-4 grid grid-cols-3 gap-3 rounded-xl border border-white/5 bg-white/[0.02] p-3">
                <div className="text-center">
                  <p className="text-lg font-bold text-white">{club.memberCount}</p>
                  <p className="text-[0.6rem] uppercase tracking-wider text-white/35">Hội viên</p>
                </div>
                <div className="text-center">
                  <p className="text-lg font-bold text-white">{club.coachCount}</p>
                  <p className="text-[0.6rem] uppercase tracking-wider text-white/35">HLV</p>
                </div>
                <div className="text-center">
                  <p className="text-lg font-bold text-cyan-400">{(club.monthlyFee / 1000).toFixed(0)}k</p>
                  <p className="text-[0.6rem] uppercase tracking-wider text-white/35">Phí/tháng</p>
                </div>
              </div>

              {/* Martial Arts Tags */}
              <div className="mt-3 flex flex-wrap gap-1.5">
                {club.martialArts.map((art) => (
                  <span key={art} className="rounded-lg border border-white/6 bg-white/[0.03] px-2 py-0.5 text-[0.65rem] text-white/50">
                    {art}
                  </span>
                ))}
              </div>

              {/* Footer */}
              <div className="mt-4 flex items-center justify-between border-t border-white/5 pt-3">
                <div className="flex items-center gap-3 text-xs text-white/40">
                  <span className="flex items-center gap-1">
                    <Calendar className="size-3" />
                    {new Date(club.foundedDate).getFullYear()}
                  </span>
                  <span className="flex items-center gap-1">
                    <UserCog className="size-3" />
                    {club.headCoach.split(" ").slice(-2).join(" ")}
                  </span>
                </div>
                <ChevronRight className="size-4 text-white/25 transition-transform group-hover:translate-x-0.5" />
              </div>
            </div>
          </div>
        ))}
      </div>

      <Pagination
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={setCurrentPage}
        totalItems={filteredClubs.length}
        pageSize={PAGE_SIZE}
      />
    </section>
  );
}
