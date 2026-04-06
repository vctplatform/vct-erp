"use client";

import { Award, QrCode, Shield, Calendar, User, Download, Plus, ExternalLink } from "lucide-react";

import { PageHeader, SectionPanel, StatusBadge, GlassButton, Tabs, SearchInput, FederationDataTable } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_CERTIFICATES } from "@/lib/federation/mock-data";
import { cn } from "@/lib/utils";
import { useState } from "react";

const TYPE_MAP = {
  belt: { label: "Đai cấp", icon: Shield, color: "text-cyan-600 dark:text-cyan-400", bg: "bg-cyan-500/10" },
  coach: { label: "HLV", icon: User, color: "text-indigo-600 dark:text-indigo-400", bg: "bg-indigo-500/10" },
  referee: { label: "Trọng tài", icon: Shield, color: "text-amber-600 dark:text-amber-400", bg: "bg-amber-500/10" },
  achievement: { label: "Thành tích", icon: Award, color: "text-emerald-600 dark:text-emerald-400", bg: "bg-emerald-500/10" },
  participation: { label: "Tham dự", icon: Calendar, color: "text-sky-600 dark:text-sky-400", bg: "bg-sky-500/10" },
};

export default function CertificatesPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");

  const filtered = MOCK_CERTIFICATES.filter((c) => {
    if (activeTab !== "all" && c.type !== activeTab) return false;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      return c.recipientName.toLowerCase().includes(q) || c.code.toLowerCase().includes(q) || c.title.toLowerCase().includes(q);
    }
    return true;
  });

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Chứng nhận số"
        title="Chứng chỉ Số"
        description="Quản lý cấp phát, xác thực chứng chỉ số — Đai cấp, HLV, Trọng tài, Thành tích thi đấu."
        actions={
          <GlassButton variant="primary" size="sm">
            <Plus className="size-4" />
            Cấp chứng chỉ
          </GlassButton>
        }
      />

      <div className="stagger-children grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={Award} label="Tổng chứng chỉ" value={12_847} color="cyan" trend={{ direction: "up", value: "+356" }} />
        <StatCard icon={Shield} label="Đai cấp" value={9_240} color="emerald" />
        <StatCard icon={User} label="HLV & Trọng tài" value={2_840} color="indigo" />
        <StatCard icon={Award} label="Thành tích" value={767} color="amber" />
      </div>

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          tabs={[
            { id: "all", label: "Tất cả", count: MOCK_CERTIFICATES.length },
            { id: "belt", label: "Đai cấp" },
            { id: "coach", label: "HLV" },
            { id: "referee", label: "Trọng tài" },
            { id: "achievement", label: "Thành tích" },
          ]}
          activeTab={activeTab}
          onChange={setActiveTab}
        />
        <SearchInput value={searchQuery} onChange={setSearchQuery} placeholder="Tìm theo mã, tên, loại..." className="w-56" />
      </div>

      {/* Certificate Cards */}
      <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        {filtered.map((cert) => {
          const type = TYPE_MAP[cert.type];
          const TypeIcon = type.icon;
          return (
            <div
              key={cert.id}
              className="group relative overflow-hidden rounded-2xl border border-[var(--color-border)] bg-[var(--color-panel)] transition-all duration-300 hover:-translate-y-1 dark:border-white/6 dark:bg-white/[0.025] dark:hover:border-cyan-500/15 dark:hover:shadow-[0_0_30px_rgba(6,182,212,0.08)]"
            >
              {/* Certificate Header Band */}
              <div className="relative border-b border-[var(--color-border)] bg-gradient-to-r from-[var(--color-canvas-soft)] to-[var(--color-panel)] px-6 py-4 dark:border-white/5 dark:from-cyan-500/[0.04] dark:to-transparent">
                <div className="pointer-events-none absolute inset-x-0 bottom-0 h-px bg-gradient-to-r from-transparent via-cyan-500/15 to-transparent dark:block" />
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className={cn("rounded-xl p-2", type.bg)}>
                      <TypeIcon className={cn("size-5", type.color)} />
                    </div>
                    <div>
                      <p className="font-mono text-[0.65rem] text-[var(--color-ink-soft)] dark:text-white/40">{cert.code}</p>
                      <p className={cn("text-xs font-medium", type.color)}>{type.label}</p>
                    </div>
                  </div>
                  <StatusBadge status={cert.status} size="md" />
                </div>
              </div>

              {/* Certificate Body */}
              <div className="p-6">
                <h3 className="text-sm font-semibold text-[var(--color-ink)] dark:text-white">
                  {cert.title}
                </h3>
                <p className="mt-2 flex items-center gap-1.5 text-xs text-[var(--color-ink-soft)] dark:text-white/50">
                  <User className="size-3" />
                  {cert.recipientName}
                </p>

                <div className="mt-4 flex items-center gap-4 text-xs text-[var(--color-ink-soft)] dark:text-white/40">
                  <span className="flex items-center gap-1">
                    <Calendar className="size-3" />
                    {new Date(cert.issuedDate).toLocaleDateString("vi-VN")}
                  </span>
                  {cert.expiryDate && (
                    <span>→ {new Date(cert.expiryDate).toLocaleDateString("vi-VN")}</span>
                  )}
                </div>

                {/* QR Code Section */}
                {cert.qrCode && cert.status === "issued" && (
                  <div className="mt-4 flex items-center justify-between rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-3 py-2 dark:border-white/5 dark:bg-white/[0.02]">
                    <div className="flex items-center gap-2 text-xs text-[var(--color-ink-soft)] dark:text-white/45">
                      <QrCode className="size-4 text-cyan-600 dark:text-cyan-400" />
                      <span>Xác thực online</span>
                    </div>
                    <a
                      href={cert.qrCode}
                      target="_blank"
                      rel="noreferrer"
                      className="flex items-center gap-1 text-xs font-medium text-cyan-600 hover:text-cyan-700 dark:text-cyan-400"
                    >
                      Mở link
                      <ExternalLink className="size-3" />
                    </a>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}
