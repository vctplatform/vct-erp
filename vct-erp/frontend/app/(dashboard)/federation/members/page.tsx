"use client";

import { useState, useMemo } from "react";
import { Users, Search, Download, Plus, Filter, ChevronRight, MapPin } from "lucide-react";
import Link from "next/link";

import { PageHeader, SectionPanel, StatusBadge, FederationDataTable, GlassButton, Tabs, SearchInput, Pagination, ViewToggle, FilterChip, FilterBar } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_MEMBERS, MOCK_DASHBOARD } from "@/lib/federation/mock-data";
import { BELT_CONFIG } from "@/lib/federation/constants";
import { cn } from "@/lib/utils";

const PAGE_SIZE = 8;

export default function MembersPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [view, setView] = useState<"table" | "grid">("table");
  const [selectedProvinces, setSelectedProvinces] = useState<string[]>([]);

  const provinces = useMemo(() => [...new Set(MOCK_MEMBERS.map((m) => m.province))], []);

  const filteredMembers = useMemo(() => {
    return MOCK_MEMBERS.filter((m) => {
      if (activeTab === "active" && m.status !== "active") return false;
      if (activeTab === "pending" && m.status !== "pending") return false;
      if (activeTab === "inactive" && m.status !== "inactive") return false;
      if (selectedProvinces.length > 0 && !selectedProvinces.includes(m.province)) return false;
      if (searchQuery) {
        const q = searchQuery.toLowerCase();
        return m.fullName.toLowerCase().includes(q) || m.memberId.toLowerCase().includes(q) || m.province.toLowerCase().includes(q);
      }
      return true;
    });
  }, [activeTab, searchQuery, selectedProvinces]);

  const totalPages = Math.ceil(filteredMembers.length / PAGE_SIZE);
  const paginatedMembers = filteredMembers.slice((currentPage - 1) * PAGE_SIZE, currentPage * PAGE_SIZE);

  const toggleProvince = (p: string) => {
    setSelectedProvinces((prev) => prev.includes(p) ? prev.filter((x) => x !== p) : [...prev, p]);
    setCurrentPage(1);
  };

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Quản lý hội viên"
        title="Hội viên Toàn quốc"
        description="Quản lý hồ sơ, thẻ hội viên số và toàn bộ thông tin của võ sinh trên khắp 63 tỉnh thành."
        actions={
          <div className="flex items-center gap-2">
            <GlassButton variant="secondary" size="sm">
              <Download className="size-4" />
              Xuất Excel
            </GlassButton>
            <GlassButton variant="primary" size="sm">
              <Plus className="size-4" />
              Thêm mới
            </GlassButton>
          </div>
        }
      />

      {/* Quick Stats */}
      <div className="stagger-children grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={Users} label="Tổng hội viên" value={MOCK_DASHBOARD.totalMembers} color="cyan" trend={{ direction: "up", value: "+478" }} />
        <StatCard icon={Users} label="Đang hoạt động" value={45_120} color="emerald" />
        <StatCard icon={Users} label="Chờ duyệt" value={234} color="amber" />
        <StatCard icon={Users} label="Tạm ngưng" value={128} color="rose" />
      </div>

      {/* Filter Bar */}
      <div className="space-y-3">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <Tabs
            tabs={[
              { id: "all", label: "Tất cả", count: MOCK_MEMBERS.length },
              { id: "active", label: "Hoạt động", count: MOCK_MEMBERS.filter((m) => m.status === "active").length },
              { id: "pending", label: "Chờ duyệt", count: MOCK_MEMBERS.filter((m) => m.status === "pending").length },
              { id: "inactive", label: "Không hoạt động", count: 0 },
            ]}
            activeTab={activeTab}
            onChange={(id) => { setActiveTab(id); setCurrentPage(1); }}
          />
          <div className="flex items-center gap-2">
            <SearchInput value={searchQuery} onChange={(v) => { setSearchQuery(v); setCurrentPage(1); }} placeholder="Tìm theo tên, mã HV, tỉnh..." className="w-56" />
            <ViewToggle view={view} onChange={setView} />
          </div>
        </div>

        {/* Province Filters */}
        <FilterBar activeCount={selectedProvinces.length} onReset={() => setSelectedProvinces([])}>
          {provinces.slice(0, 8).map((p) => (
            <FilterChip key={p} label={p} selected={selectedProvinces.includes(p)} onClick={() => toggleProvince(p)} />
          ))}
        </FilterBar>
      </div>

      {/* Members Display */}
      {view === "table" ? (
        <SectionPanel title="Danh sách Hội viên">
          <FederationDataTable
            columns={[
              { key: "memberId", label: "Mã HV", width: "140px" },
              { key: "name", label: "Họ và tên" },
              { key: "belt", label: "Đai" },
              { key: "club", label: "CLB" },
              { key: "province", label: "Tỉnh/TP" },
              { key: "joinDate", label: "Ngày gia nhập" },
              { key: "status", label: "Trạng thái", align: "center" },
              { key: "action", label: "", width: "40px" },
            ]}
            rows={paginatedMembers.map((m) => {
              const belt = BELT_CONFIG[m.currentBelt];
              return {
                memberId: (
                  <span className="font-mono text-xs text-cyan-600 dark:text-cyan-400">
                    {m.memberId}
                  </span>
                ),
                name: (
                  <div className="flex items-center gap-3">
                    <div className="flex h-9 w-9 items-center justify-center rounded-full bg-gradient-to-br from-cyan-500/20 to-teal-500/10 text-xs font-bold text-cyan-700 dark:text-cyan-300">
                      {m.fullName.split(" ").pop()?.charAt(0)}
                    </div>
                    <div>
                      <p className="font-medium text-[var(--color-ink)] dark:text-white">{m.fullName}</p>
                      <p className="text-xs text-[var(--color-ink-soft)] dark:text-white/40">{m.martialArt}</p>
                    </div>
                  </div>
                ),
                belt: (
                  <div className="flex items-center gap-2">
                    <span className="h-4 w-4 rounded-md border border-[var(--color-border)] dark:border-white/12" style={{ backgroundColor: belt.color }} />
                    <span className="text-xs">{belt.label}</span>
                  </div>
                ),
                club: <span className="text-sm">{m.clubName}</span>,
                province: (
                  <span className="flex items-center gap-1 text-sm">
                    <MapPin className="size-3 text-[var(--color-ink-soft)] dark:text-white/30" />
                    {m.province}
                  </span>
                ),
                joinDate: (
                  <span className="text-xs text-[var(--color-ink-soft)] dark:text-white/50">
                    {new Date(m.joinDate).toLocaleDateString("vi-VN")}
                  </span>
                ),
                status: <StatusBadge status={m.status} />,
                action: (
                  <Link href={`/federation/members/${m.id}`}>
                    <ChevronRight className="size-4 text-[var(--color-ink-soft)] dark:text-white/30" />
                  </Link>
                ),
              };
            })}
          />
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={setCurrentPage}
            totalItems={filteredMembers.length}
            pageSize={PAGE_SIZE}
          />
        </SectionPanel>
      ) : (
        /* Grid View */
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {paginatedMembers.map((m) => {
            const belt = BELT_CONFIG[m.currentBelt];
            return (
              <Link
                key={m.id}
                href={`/federation/members/${m.id}`}
                className="group rounded-2xl border border-[var(--color-border)] bg-[var(--color-panel)] p-5 transition-all duration-200 hover:-translate-y-1 hover:shadow-lg dark:border-white/6 dark:bg-white/[0.025] dark:hover:border-cyan-500/15 dark:hover:shadow-[0_0_30px_rgba(6,182,212,0.1)]"
              >
                <div className="flex items-center gap-3">
                  <div className="flex h-12 w-12 items-center justify-center rounded-full bg-gradient-to-br from-cyan-500/20 to-teal-500/10 text-sm font-bold text-cyan-700 dark:text-cyan-300">
                    {m.fullName.split(" ").pop()?.charAt(0)}
                  </div>
                  <div className="min-w-0">
                    <p className="truncate font-medium text-[var(--color-ink)] dark:text-white">{m.fullName}</p>
                    <p className="font-mono text-[0.65rem] text-cyan-600 dark:text-cyan-400">{m.memberId}</p>
                  </div>
                </div>
                <div className="mt-4 flex items-center gap-2">
                  <span className="h-4 w-4 rounded-md border border-[var(--color-border)] dark:border-white/12" style={{ backgroundColor: belt.color }} />
                  <span className="text-xs text-[var(--color-ink-soft)] dark:text-white/60">{belt.label}</span>
                  <span className="mx-1 text-[var(--color-border)]">•</span>
                  <span className="text-xs text-[var(--color-ink-soft)] dark:text-white/50">{m.martialArt}</span>
                </div>
                <div className="mt-3 flex items-center justify-between">
                  <span className="flex items-center gap-1 text-xs text-[var(--color-ink-soft)] dark:text-white/40">
                    <MapPin className="size-3" />
                    {m.province}
                  </span>
                  <StatusBadge status={m.status} />
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </section>
  );
}
