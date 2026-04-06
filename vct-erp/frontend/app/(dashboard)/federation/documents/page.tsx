"use client";

import { FileText, AlertTriangle, CheckCircle, Clock, Eye, Download, Plus, Paperclip, ArrowRight } from "lucide-react";

import { PageHeader, SectionPanel, StatusBadge, GlassButton, Tabs, SearchInput, ProgressBar } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_DOCUMENTS } from "@/lib/federation/mock-data";
import { cn } from "@/lib/utils";
import { useState } from "react";

const PRIORITY_MAP = {
  normal: { label: "Bình thường", color: "text-slate-500", bg: "bg-slate-500/10" },
  urgent: { label: "Khẩn", color: "text-amber-600 dark:text-amber-400", bg: "bg-amber-500/10" },
  critical: { label: "Hỏa tốc", color: "text-rose-600 dark:text-rose-400", bg: "bg-rose-500/10" },
};

const TYPE_MAP = {
  incoming: "Đến",
  outgoing: "Đi",
  internal: "Nội bộ",
};

export default function DocumentsPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");

  const filtered = MOCK_DOCUMENTS.filter((d) => {
    if (activeTab === "approved" && d.status !== "approved") return false;
    if (activeTab === "pending" && d.status !== "pending_review" && d.status !== "draft") return false;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      return d.title.toLowerCase().includes(q) || d.code.toLowerCase().includes(q);
    }
    return true;
  });

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Quản lý công văn"
        title="Văn bản & Công văn"
        description="Quản lý toàn bộ hệ thống văn bản — Công văn đến/đi, quy chế, thông báo, quy trình phê duyệt."
        actions={
          <GlassButton variant="primary" size="sm">
            <Plus className="size-4" />
            Tạo văn bản
          </GlassButton>
        }
      />

      <div className="stagger-children grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={FileText} label="Tổng văn bản" value={MOCK_DOCUMENTS.length} color="cyan" />
        <StatCard icon={CheckCircle} label="Đã phê duyệt" value={2} color="emerald" />
        <StatCard icon={Clock} label="Chờ duyệt" value={1} color="amber" />
        <StatCard icon={FileText} label="Nháp" value={1} color="sky" />
      </div>

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          tabs={[
            { id: "all", label: "Tất cả", count: MOCK_DOCUMENTS.length },
            { id: "approved", label: "Đã duyệt", count: 2 },
            { id: "pending", label: "Chờ xử lý", count: 2 },
          ]}
          activeTab={activeTab}
          onChange={setActiveTab}
        />
        <SearchInput value={searchQuery} onChange={setSearchQuery} placeholder="Tìm theo mã, tiêu đề..." className="w-56" />
      </div>

      {/* Document Cards */}
      <div className="space-y-4">
        {filtered.map((doc) => {
          const priority = PRIORITY_MAP[doc.priority];
          const completedSteps = doc.approvalChain.filter((s) => s.status === "approved").length;
          const totalSteps = doc.approvalChain.length;

          return (
            <div
              key={doc.id}
              className="group rounded-2xl border border-[var(--color-border)] bg-[var(--color-panel)] p-5 transition-all duration-200 hover:-translate-y-0.5 dark:border-white/6 dark:bg-white/[0.025] dark:hover:border-cyan-500/12"
            >
              <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div className="flex items-start gap-4">
                  <div className="rounded-xl bg-gradient-to-br from-cyan-500/10 to-teal-500/5 p-3">
                    <FileText className="size-5 text-cyan-600 dark:text-cyan-400" />
                  </div>
                  <div className="min-w-0">
                    <div className="flex flex-wrap items-center gap-2">
                      <h3 className="font-semibold text-[var(--color-ink)] dark:text-white">{doc.title}</h3>
                      <span className={cn("rounded-lg px-2 py-0.5 text-[0.65rem] font-medium", priority.bg, priority.color)}>
                        {priority.label}
                      </span>
                    </div>
                    <div className="mt-1.5 flex flex-wrap items-center gap-3 text-xs text-[var(--color-ink-soft)] dark:text-white/45">
                      <span className="font-mono">{doc.code}</span>
                      <span>•</span>
                      <span>{TYPE_MAP[doc.type]}</span>
                      <span>•</span>
                      <span>{doc.category}</span>
                      <span>•</span>
                      <span>{new Date(doc.createdDate).toLocaleDateString("vi-VN")}</span>
                    </div>

                    {/* Attachments */}
                    {doc.attachments.length > 0 && (
                      <div className="mt-2 flex flex-wrap gap-2">
                        {doc.attachments.map((file) => (
                          <span key={file} className="inline-flex items-center gap-1 rounded-lg border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-2 py-1 text-[0.65rem] dark:border-white/5 dark:bg-white/[0.02]">
                            <Paperclip className="size-3 text-[var(--color-ink-soft)] dark:text-white/35" />
                            <span className="text-[var(--color-ink-soft)] dark:text-white/50">{file}</span>
                          </span>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  <StatusBadge status={doc.status} size="md" />
                </div>
              </div>

              {/* Approval Pipeline */}
              {doc.approvalChain.length > 0 && (
                <div className="mt-4 rounded-xl border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4 dark:border-white/5 dark:bg-white/[0.02]">
                  <div className="mb-3 flex items-center justify-between">
                    <p className="text-[0.62rem] uppercase tracking-[0.25em] text-[var(--color-ink-soft)] dark:text-white/35">
                      Quy trình phê duyệt
                    </p>
                    <span className="text-xs text-[var(--color-ink-soft)] dark:text-white/40">
                      {completedSteps}/{totalSteps}
                    </span>
                  </div>
                  <ProgressBar value={completedSteps} max={totalSteps} color="emerald" className="mb-3" />
                  <div className="flex items-center gap-2">
                    {doc.approvalChain.map((step, i) => (
                      <div key={step.order} className="flex items-center gap-2">
                        <div className={cn(
                          "flex items-center gap-2 rounded-lg px-3 py-2 text-xs",
                          step.status === "approved"
                            ? "border border-emerald-500/15 bg-emerald-500/5 text-emerald-700 dark:text-emerald-400"
                            : step.status === "rejected"
                            ? "border border-rose-500/15 bg-rose-500/5 text-rose-700 dark:text-rose-400"
                            : "border border-[var(--color-border)] bg-[var(--color-panel)] text-[var(--color-ink-soft)] dark:border-white/5 dark:text-white/40"
                        )}>
                          {step.status === "approved" && <CheckCircle className="size-3" />}
                          {step.status === "pending" && <Clock className="size-3" />}
                          <span className="font-medium">{step.approver}</span>
                        </div>
                        {i < doc.approvalChain.length - 1 && (
                          <ArrowRight className="size-3.5 text-[var(--color-ink-soft)] dark:text-white/25" />
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </section>
  );
}
