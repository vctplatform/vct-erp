import { BarChart3, BellRing, Building2, ChartPie, Search } from "lucide-react";
import Link from "next/link";

const navItems = [
  { label: "Tong quan", href: "/", icon: BarChart3 },
  { label: "Doanh thu", href: "/", icon: ChartPie },
  { label: "Cong ty", href: "/", icon: Building2 },
  { label: "Canh bao", href: "/", icon: BellRing },
];

export function Sidebar() {
  return (
    <details
      open
      className="group/sidebar hidden w-[18rem] shrink-0 rounded-[1.75rem] border border-[var(--color-border)] bg-[linear-gradient(180deg,rgba(16,36,70,0.96),rgba(10,24,50,0.96))] text-white shadow-[0_24px_80px_rgba(10,24,50,0.32)] lg:block"
    >
      <summary className="flex cursor-pointer list-none items-center justify-between border-b border-white/10 px-5 py-5">
        <div>
          <p className="text-[0.65rem] uppercase tracking-[0.32em] text-white/45">
            VCT
          </p>
          <h2 className="mt-2 text-xl font-semibold tracking-tight">
            Command Center
          </h2>
        </div>
        <span className="rounded-full border border-white/15 bg-white/10 px-3 py-1 text-xs text-white/70 group-open/sidebar:hidden">
          Open
        </span>
      </summary>

      <div className="flex h-[calc(100vh-3rem)] flex-col justify-between p-4 group-open/sidebar:p-5">
        <div className="space-y-6">
          <div className="rounded-[1.35rem] border border-white/10 bg-white/6 p-3">
            <div className="flex items-center gap-3 rounded-2xl bg-white/8 px-3 py-3 text-sm text-white/70">
              <Search className="size-4" />
              <span>Scan reports, vouchers, risk</span>
            </div>
          </div>

          <nav className="space-y-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              return (
                <Link
                  key={item.label}
                  href={item.href}
                  className="flex items-center gap-3 rounded-2xl px-3 py-3 text-sm text-white/72 transition hover:bg-white/8 hover:text-white"
                >
                  <Icon className="size-4" />
                  <span>{item.label}</span>
                </Link>
              );
            })}
          </nav>
        </div>

        <div className="rounded-[1.4rem] border border-emerald-400/20 bg-emerald-500/12 p-4 text-sm text-emerald-100">
          <p className="text-xs uppercase tracking-[0.24em] text-emerald-200/70">
            System Health
          </p>
          <p className="mt-2 font-medium">
            Ledger synced. Realtime dashboard listening for new cash movement.
          </p>
        </div>
      </div>
    </details>
  );
}
