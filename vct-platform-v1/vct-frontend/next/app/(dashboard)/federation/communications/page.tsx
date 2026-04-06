"use client";

import { Megaphone, Newspaper, Calendar, Eye, Plus, Users, Globe, Clock, Tag, Radio } from "lucide-react";

import { PageHeader, SectionPanel, StatusBadge, GlassButton, Tabs, SearchInput } from "@/components/federation/ui/shared";
import { StatCard } from "@/components/federation/ui/stat-card";
import { MOCK_COMMUNICATIONS } from "@/lib/federation/mock-data";
import { cn } from "@/lib/utils";
import { useState } from "react";

const TYPE_MAP = {
  announcement: { label: "Thông báo", icon: Megaphone, color: "text-amber-400", bg: "bg-amber-500/10" },
  news: { label: "Tin tức", icon: Newspaper, color: "text-cyan-400", bg: "bg-cyan-500/10" },
  event: { label: "Sự kiện", icon: Calendar, color: "text-indigo-400", bg: "bg-indigo-500/10" },
};

const AUDIENCE_LABELS = {
  all: "Tất cả",
  clubs: "CLB",
  coaches: "HLV",
  members: "Hội viên",
};

export default function CommunicationsPage() {
  const [activeTab, setActiveTab] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");

  const filtered = MOCK_COMMUNICATIONS.filter((c) => {
    if (activeTab !== "all" && c.type !== activeTab) return false;
    if (searchQuery) {
      return c.title.toLowerCase().includes(searchQuery.toLowerCase());
    }
    return true;
  });

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Kênh truyền thông"
        title="Truyền thông"
        description="Quản lý tin tức, thông báo, sự kiện — Phát hành nội dung đến toàn bộ hội viên và CLB."
        actions={
          <GlassButton variant="primary" size="sm">
            <Plus className="size-4" />
            Tạo bài viết
          </GlassButton>
        }
      />

      <div className="stagger-fed grid grid-cols-2 gap-4 lg:grid-cols-4">
        <StatCard icon={Newspaper} label="Tin đã đăng" value={156} color="cyan" />
        <StatCard icon={Eye} label="Lượt đọc tháng" value={45200} color="emerald" />
        <StatCard icon={Calendar} label="Sự kiện sắp tới" value={3} color="indigo" />
        <StatCard icon={Megaphone} label="Thông báo chờ" value={2} color="amber" />
      </div>

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <Tabs
          tabs={[
            { id: "all", label: "Tất cả", count: MOCK_COMMUNICATIONS.length },
            { id: "news", label: "Tin tức" },
            { id: "announcement", label: "Thông báo" },
            { id: "event", label: "Sự kiện" },
          ]}
          activeTab={activeTab}
          onChange={setActiveTab}
        />
        <SearchInput value={searchQuery} onChange={setSearchQuery} placeholder="Tìm bài viết..." className="w-56" />
      </div>

      {/* Communication Cards */}
      <div className="grid gap-5 md:grid-cols-2">
        {filtered.map((item) => {
          const type = TYPE_MAP[item.type as keyof typeof TYPE_MAP];
          const TypeIcon = type.icon;
          return (
            <div
              key={item.id}
              className="group overflow-hidden rounded-2xl border border-white/6 bg-white/[0.025] transition-all duration-300 neon-hover"
            >
              {/* Color band top */}
              <div className={cn("h-1", type.bg)} />

              <div className="p-5">
                <div className="flex items-start justify-between gap-3 text-center">
                  <div className="flex items-center gap-3">
                    <div className={cn("rounded-xl p-2.5", type.bg)}>
                      <TypeIcon className={cn("size-5", type.color)} />
                    </div>
                    <div>
                      <span className={cn("text-[0.65rem] font-medium", type.color)}>{type.label}</span>
                      <h3 className="mt-0.5 font-semibold text-white">{item.title}</h3>
                    </div>
                  </div>
                  <StatusBadge status={item.status} />
                </div>

                <p className="mt-3 line-clamp-2 text-sm text-white/50">
                  {item.content}
                </p>

                <div className="mt-4 flex flex-wrap items-center gap-3 text-xs text-white/40">
                  <span className="flex items-center gap-1">
                    <Clock className="size-3" />
                    {new Date(item.publishDate).toLocaleDateString("vi-VN")}
                  </span>
                  <span className="flex items-center gap-1">
                    <Tag className="size-3" />
                    {item.author}
                  </span>
                  <span className="flex items-center gap-1">
                    <Users className="size-3" />
                    {AUDIENCE_LABELS[item.targetAudience as keyof typeof AUDIENCE_LABELS]}
                  </span>
                  {item.readCount && (
                    <span className="flex items-center gap-1">
                      <Eye className="size-3" />
                      {item.readCount.toLocaleString("vi-VN")} lượt đọc
                    </span>
                  )}
                </div>

                {/* Event details */}
                {item.type === "event" && item.eventDate && (
                  <div className="mt-3 flex items-center gap-3 rounded-xl border border-indigo-400/10 bg-indigo-500/5 px-3 py-2 text-xs">
                    <Calendar className="size-3.5 text-indigo-400" />
                    <span className="font-medium text-indigo-300">
                      {new Date(item.eventDate).toLocaleDateString("vi-VN")}
                    </span>
                    {item.eventLocation && (
                      <>
                        <span className="text-white/30">•</span>
                        <span className="text-indigo-400">{item.eventLocation}</span>
                      </>
                    )}
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
