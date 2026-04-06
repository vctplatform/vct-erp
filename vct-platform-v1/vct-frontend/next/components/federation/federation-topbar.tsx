"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { Bell, Search, Moon, Sun, Menu, X, ChevronRight, Shield, Check, Clock, AlertTriangle, Info } from "lucide-react";
import { useTheme } from "next-themes";
import { useState, useRef, useEffect } from "react";
import { cn } from "@/lib/utils";
import { FEDERATION_NAV } from "@/lib/federation/constants";
import {
  LayoutDashboard, Users, Building2, UserCog, GraduationCap,
  Trophy, Award, Wallet, FileText, Megaphone, BarChart3, Settings,
} from "lucide-react";

const ICON_MAP: Record<string, React.ElementType> = {
  LayoutDashboard, Users, Building2, UserCog, GraduationCap,
  Trophy, Award, Wallet, FileText, Megaphone, BarChart3, Settings,
};

/* ── Notification Data ── */
const MOCK_NOTIFICATIONS = [
  { id: 1, type: "urgent" as const, title: "CLB Bình Dương chờ duyệt", desc: "Hồ sơ đăng ký đã quá hạn 3 ngày", time: "5 phút trước", read: false },
  { id: 2, type: "info" as const, title: "Giải vô địch Quốc gia 2026", desc: "Đã đóng đăng ký - 128 VĐV tham gia", time: "1 giờ trước", read: false },
  { id: 3, type: "success" as const, title: "Thi thăng đai đợt 3 hoàn tất", desc: "42/45 võ sinh đạt yêu cầu", time: "2 giờ trước", read: false },
  { id: 4, type: "warning" as const, title: "Niên liễm Q2 sắp hết hạn", desc: "128 hội viên chưa đóng phí", time: "Hôm qua", read: true },
];

const NOTIF_ICONS = {
  urgent: AlertTriangle,
  warning: AlertTriangle,
  info: Info,
  success: Check,
};

const NOTIF_COLORS = {
  urgent: "text-rose-400 bg-rose-500/10",
  warning: "text-amber-400 bg-amber-500/10",
  info: "text-sky-400 bg-sky-500/10",
  success: "text-emerald-400 bg-emerald-500/10",
};

/* ── Breadcrumb ── */
const ROUTE_LABELS: Record<string, string> = {
  federation: "Liên đoàn",
  members: "Hội viên",
  clubs: "CLB / Võ đường",
  personnel: "Nhân sự & HLV",
  examinations: "Thi & Thăng đai",
  tournaments: "Giải đấu",
  certificates: "Chứng chỉ số",
  finance: "Tài chính",
  documents: "Văn bản",
  communications: "Truyền thông",
  reports: "Báo cáo",
  settings: "Cài đặt",
};

function Breadcrumb() {
  const pathname = usePathname();
  const segments = pathname.split("/").filter(Boolean);
  if (segments.length <= 1) return null;

  return (
    <nav className="flex items-center gap-1.5 text-sm">
      {segments.map((seg, i) => {
        const href = "/" + segments.slice(0, i + 1).join("/");
        const label = ROUTE_LABELS[seg] || seg;
        const isLast = i === segments.length - 1;
        return (
          <span key={href} className="flex items-center gap-1.5">
            {i > 0 && <ChevronRight className="size-3 text-white/25" />}
            {isLast ? (
              <span className="font-medium text-white">{label}</span>
            ) : (
              <Link href={href} className="text-white/45 transition-colors hover:text-white/75">
                {label}
              </Link>
            )}
          </span>
        );
      })}
    </nav>
  );
}

/* ── Notification Dropdown ── */
function NotificationDropdown() {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const unreadCount = MOCK_NOTIFICATIONS.filter((n) => !n.read).length;

  useEffect(() => {
    function handle(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, []);

  return (
    <div ref={ref} className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="relative rounded-xl border border-white/10 p-2.5 transition-colors hover:bg-white/5"
      >
        <Bell className="size-4 text-white/55" />
        {unreadCount > 0 && (
          <span className="absolute -right-0.5 -top-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-gradient-to-r from-rose-500 to-pink-500 text-[0.55rem] font-bold text-white shadow-[0_2px_8px_rgba(225,29,72,0.4)]">
            {unreadCount}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 top-full z-50 mt-2 w-[360px] animate-fed-slide-up overflow-hidden rounded-2xl border border-white/10 bg-[#081020]/95 backdrop-blur-2xl shadow-[0_20px_60px_rgba(0,0,0,0.5)]">
          <div className="flex items-center justify-between border-b border-white/10 px-4 py-3">
            <h3 className="text-sm font-semibold text-white">Thông báo</h3>
            <span className="rounded-full bg-cyan-500/10 px-2 py-0.5 text-[0.6rem] font-bold text-cyan-400">
              {unreadCount} mới
            </span>
          </div>

          <div className="max-h-[360px] overflow-y-auto vct-custom-scrollbar">
            {MOCK_NOTIFICATIONS.map((notif) => {
              const Icon = NOTIF_ICONS[notif.type];
              return (
                <div
                  key={notif.id}
                  className={cn(
                    "flex gap-3 border-b border-white/5 px-4 py-3 transition-colors hover:bg-white/[0.03]",
                    !notif.read && "bg-cyan-500/[0.04]"
                  )}
                >
                  <div className={cn("mt-0.5 rounded-lg p-1.5", NOTIF_COLORS[notif.type])}>
                    <Icon className="size-3.5" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className={cn("text-sm text-white/85", !notif.read && "font-medium")}>
                      {notif.title}
                    </p>
                    <p className="mt-0.5 text-xs text-white/40">
                      {notif.desc}
                    </p>
                    <p className="mt-1 flex items-center gap-1 text-[0.62rem] text-white/30">
                      <Clock className="size-2.5" />
                      {notif.time}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}

/* ── User Menu ── */
function UserMenu() {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handle(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, []);

  return (
    <div ref={ref} className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-2 transition-colors hover:border-white/20"
      >
        <div className="flex h-7 w-7 items-center justify-center rounded-full bg-gradient-to-br from-cyan-500 to-teal-500 text-xs font-bold text-white shadow-[0_4px_12px_rgba(6,182,212,0.3)]">
          CT
        </div>
        <div className="hidden sm:block text-left">
          <p className="text-xs font-medium text-white">Chủ tịch</p>
          <p className="text-[0.6rem] text-white/40">Admin</p>
        </div>
      </button>

      {open && (
        <div className="absolute right-0 top-full z-50 mt-2 w-56 animate-fed-slide-up overflow-hidden rounded-xl border border-white/10 bg-[#081020]/95 backdrop-blur-2xl shadow-[0_16px_48px_rgba(0,0,0,0.5)]">
          <div className="border-b border-white/10 px-4 py-3">
            <p className="text-sm font-semibold text-white">Võ sư Nguyễn Văn A</p>
            <p className="text-xs text-white/40">admin@vct.org.vn</p>
          </div>
          <div className="py-1">
            <Link href="#" className="block px-4 py-2 text-sm text-white/60 hover:bg-white/5 hover:text-white">Hồ sơ cá nhân</Link>
            <Link href="#" className="block px-4 py-2 text-sm text-white/60 hover:bg-white/5 hover:text-white">Cài đặt</Link>
          </div>
          <div className="border-t border-white/10 py-1">
            <button className="w-full px-4 py-2 text-left text-sm text-rose-400 hover:bg-rose-500/5">
              Đăng xuất
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

/* ── Mobile Drawer ── */
function MobileDrawer({ open, onClose }: { open: boolean; onClose: () => void }) {
  const pathname = usePathname();
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 lg:hidden">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="absolute left-0 top-0 h-full w-[280px] animate-fed-slide-up overflow-y-auto border-r border-white/10 bg-[#060e1e] p-4">
        <div className="mb-6 flex items-center justify-between">
          <div className="flex items-center gap-2 text-cyan-400">
            <Shield className="size-5" />
            <span className="text-sm font-bold tracking-tight text-white">VCT PORTAL</span>
          </div>
          <button onClick={onClose} className="rounded-lg p-1.5 hover:bg-white/5">
            <X className="size-5 text-white/55" />
          </button>
        </div>
        <nav className="space-y-6">
          {FEDERATION_NAV.map((group) => (
            <div key={group.group}>
              <p className="mb-2 px-2 text-[0.62rem] uppercase tracking-[0.28em] text-white/30">
                {group.group}
              </p>
              <div className="space-y-1">
                {group.items.map((item) => {
                  const Icon = ICON_MAP[item.icon];
                  const isActive = item.href === "/federation" ? pathname === "/federation" : pathname.startsWith(item.href);
                  return (
                    <Link
                      key={item.href}
                      href={item.href}
                      onClick={onClose}
                      className={cn(
                        "flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all",
                        isActive
                          ? "bg-cyan-500/10 text-cyan-300"
                          : "text-white/50 hover:bg-white/5"
                      )}
                    >
                      {Icon && <Icon className="size-4 shrink-0" />}
                      <span>{item.label}</span>
                    </Link>
                  );
                })}
              </div>
            </div>
          ))}
        </nav>
      </div>
    </div>
  );
}

/* ── Main Topbar ── */
export function FederationTopbar() {
  const { theme, setTheme } = useTheme();
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <>
      <header className="animate-fed-fade-in flex items-center justify-between gap-4 rounded-2xl border border-white/10 bg-white/[0.03] px-4 py-3 backdrop-blur-xl">
        <div className="flex items-center gap-3">
          <button
            onClick={() => setMobileOpen(true)}
            className="rounded-lg p-2 hover:bg-white/5 lg:hidden"
          >
            <Menu className="size-5 text-white/60" />
          </button>
          <Breadcrumb />
        </div>

        <div className="flex items-center gap-2">
          <button className="flex items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-sm text-white/40 transition-colors hover:border-cyan-400/20">
            <Search className="size-4" />
            <span className="hidden sm:inline">Tìm kiếm...</span>
            <kbd className="ml-2 hidden rounded bg-white/5 px-1.5 py-0.5 text-[0.6rem] font-mono sm:inline">⌘K</kbd>
          </button>

          <NotificationDropdown />

          <button
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            className="rounded-xl border border-white/10 p-2.5 transition-all hover:bg-white/5"
          >
            <Sun className="hidden size-4 text-amber-400 dark:block" />
            <Moon className="block size-4 text-slate-400 dark:hidden" />
          </button>

          <UserMenu />
        </div>
      </header>
      <MobileDrawer open={mobileOpen} onClose={() => setMobileOpen(false)} />
    </>
  );
}
