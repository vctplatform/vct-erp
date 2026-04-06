"use client";

import { Settings, Shield, Globe, Bell, Palette, Lock, Database, Users, ChevronRight, Save, Monitor, Moon, Sun } from "lucide-react";

import { PageHeader, SectionPanel, GlassButton } from "@/components/federation/ui/shared";
import { cn } from "@/lib/utils";
import { useTheme } from "next-themes";

const PROFILE_FIELDS = [
  { label: "Tên liên đoàn", value: "Liên đoàn Võ Cổ Truyền Việt Nam" },
  { label: "Viết tắt", value: "VCT VN" },
  { label: "Ngày thành lập", value: "15/03/1992" },
  { label: "Website", value: "vctplatform.vn" },
  { label: "Email liên hệ", value: "contact@vct.org.vn" },
  { label: "Hotline", value: "1900 1234 56" },
];

const PERMISSION_ROLES = [
  { name: "Chủ tịch", permissions: ["Toàn quyền"], count: 1 },
  { name: "Phó Chủ tịch", permissions: ["Quản lý nghiệp vụ", "Phê duyệt"], count: 3 },
  { name: "Trưởng Ban", permissions: ["Quản lý phân hệ", "Tạo báo cáo"], count: 5 },
  { name: "Thư ký", permissions: ["Nhập liệu", "Xem báo cáo"], count: 8 },
  { name: "Nhân viên", permissions: ["Xem thông tin cơ bản"], count: 15 },
];

const COLOR_MAP = {
  cyan: { iconBg: "bg-cyan-500/10", iconText: "text-cyan-400" },
  indigo: { iconBg: "bg-indigo-500/10", iconText: "text-indigo-400" },
};

export default function SettingsPage() {
  const { theme, setTheme } = useTheme();

  return (
    <section className="space-y-6">
      <PageHeader
        kicker="Cấu hình hệ thống"
        title="Cài đặt"
        description="Quản lý hồ sơ liên đoàn, phân quyền, giao diện và tùy chỉnh hệ thống."
      />

      {/* Theme Settings */}
      <SectionPanel title="Giao diện" kicker="Hiển thị" subtitle="Chọn chế độ sáng/tối cho giao diện hệ thống">
        <div className="grid gap-3 sm:grid-cols-3">
          {[
            { id: "light", label: "Sáng", icon: Sun, desc: "Giao diện nền sáng, phù hợp ban ngày" },
            { id: "dark", label: "Tối", icon: Moon, desc: "Neo-Glassmorphism Cyber-Teal" },
            { id: "system", label: "Hệ thống", icon: Monitor, desc: "Tự động theo cài đặt thiết bị" },
          ].map((opt) => (
            <button
              key={opt.id}
              onClick={() => setTheme(opt.id)}
              className={cn(
                "flex items-start gap-3 rounded-xl border p-4 text-left transition-all duration-200",
                theme === opt.id
                  ? "border-cyan-500/25 bg-cyan-500/5 shadow-[0_0_20px_rgba(6,182,212,0.08)]"
                  : "border-white/5 bg-white/[0.02] hover:border-cyan-500/15"
              )}
            >
              <div className={cn(
                "rounded-lg p-2",
                theme === opt.id ? "bg-cyan-500/15 text-cyan-400" : "bg-white/5 text-white/35"
              )}>
                <opt.icon className="size-4" />
              </div>
              <div>
                <p className={cn(
                  "text-sm font-medium",
                  theme === opt.id ? "text-white" : "text-white/60"
                )}>
                  {opt.label}
                </p>
                <p className="mt-0.5 text-xs text-white/35">
                  {opt.desc}
                </p>
              </div>
            </button>
          ))}
        </div>
      </SectionPanel>

      {/* Federation Profile */}
      <SectionPanel
        title="Hồ sơ Liên đoàn"
        subtitle="Thông tin cơ bản của Liên đoàn Võ Cổ Truyền Việt Nam"
        kicker="Hồ sơ"
        actions={
          <GlassButton variant="primary" size="sm">
            <Save className="size-4" />
            Lưu thay đổi
          </GlassButton>
        }
      >
        <div className="grid gap-4 md:grid-cols-2">
          {PROFILE_FIELDS.map((field) => (
            <div key={field.label}>
              <label className="mb-1.5 block text-xs font-medium text-white/45">
                {field.label}
              </label>
              <input
                type="text"
                defaultValue={field.value}
                className="w-full rounded-xl border border-white/6 bg-white/[0.03] px-4 py-2.5 text-sm text-white outline-none transition-colors focus:border-cyan-400/30"
              />
            </div>
          ))}
        </div>
      </SectionPanel>

      {/* Permissions */}
      <SectionPanel
        title="Phân quyền Quản trị"
        subtitle="Quản lý vai trò và quyền truy cập trong hệ thống"
        kicker="Bảo mật"
      >
        <div className="space-y-3">
          {PERMISSION_ROLES.map((role) => (
            <div
              key={role.name}
              className="flex items-center justify-between rounded-xl border border-white/5 bg-white/[0.02] p-4 transition-colors hover:bg-white/[0.03]"
            >
              <div className="flex items-center gap-3">
                <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-indigo-500/10">
                  <Users className="size-4 text-indigo-400" />
                </div>
                <div>
                  <p className="text-sm font-medium text-white">{role.name}</p>
                  <div className="mt-0.5 flex flex-wrap gap-1">
                    {role.permissions.map((perm) => (
                      <span key={perm} className="text-[0.6rem] text-white/35">
                        {perm}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <span className="rounded-full bg-white/5 px-2 py-0.5 text-[0.65rem] font-medium text-white/40">
                  {role.count} người
                </span>
                <ChevronRight className="size-4 text-white/25" />
              </div>
            </div>
          ))}
        </div>
      </SectionPanel>

      {/* System Info */}
      <SectionPanel title="Thông tin Hệ thống" kicker="Platform">
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {[
            { label: "Phiên bản", value: "VCT Platform v2.0" },
            { label: "Module", value: "Federation Management" },
            { label: "Build", value: "2026.04.02" },
            { label: "Database", value: "PostgreSQL 18" },
          ].map((item) => (
            <div key={item.label} className="rounded-xl border border-white/5 bg-white/[0.02] p-3 text-center">
              <p className="text-[0.62rem] uppercase tracking-[0.25em] text-white/30">
                {item.label}
              </p>
              <p className="mt-1 text-sm font-medium text-white">
                {item.value}
              </p>
            </div>
          ))}
        </div>
      </SectionPanel>
    </section>
  );
}
