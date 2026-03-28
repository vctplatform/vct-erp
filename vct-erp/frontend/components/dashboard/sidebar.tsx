import {
  Activity,
  CircleDot,
  Command,
  Fingerprint,
  Search,
  ShieldCheck,
  Sparkles,
  Wifi,
} from "lucide-react";

import { DashboardNavLink } from "@/components/dashboard/dashboard-nav-link";
import { getFinanceNavigation } from "@/lib/finance/navigation";
import { getServerLocale } from "@/lib/i18n/server";

export async function Sidebar() {
  const locale = await getServerLocale();
  const { primary, secondary } = getFinanceNavigation(locale);
  const copy =
    locale === "vi"
      ? {
          module: "Mô-đun tài chính VCT",
          title: "Trung tâm điều hành",
          description:
            "Tầm nhìn cấp ban điều hành xuyên suốt sổ cái, mảng kinh doanh, khóa sổ và kiểm soát.",
          search: "Tìm chứng từ, báo cáo, tín hiệu doanh nghiệp",
          command: "Điều hành",
          assurance: "Bảo đảm",
          readyState: "Cụm điều hành đã sẵn sàng",
          readySub: "Redis, cache và realtime đang ở trạng thái sẵn sàng cho bản alpha.",
          liveSync: "Đồng bộ live",
          auditReady: "Audit ready",
          alphaCluster: "Cụm alpha",
          alphaNote:
            "Sổ cái đã đồng bộ. Cache dashboard đã ấm. Finance hub sẵn sàng cho lưu lượng điều hành.",
          posture: "Tư thế hiện tại",
          postureNote:
            "Khóa sổ tài chính, xem cơ cấu mảng và đẩy báo cáo từ một bề mặt duy nhất.",
          signalPanel: "Tư thế kiểm soát",
          signalBullets: [
            "Realtime feed ưu tiên cho Ban điều hành.",
            "Cache 60 giây để cân bằng tốc độ và ổn định.",
            "Sẵn sàng đẩy báo cáo và kiểm toán từ cùng một bề mặt.",
          ],
        }
      : {
          module: "VCT Finance Module",
          title: "Command Center",
          description:
            "Board-grade visibility across ledger, segments, closing, and control.",
          search: "Search voucher, report, company signal",
          command: "Command",
          assurance: "Assurance",
          readyState: "Executive cluster ready",
          readySub: "Redis, cache, and realtime are primed for the alpha board flow.",
          liveSync: "Live sync",
          auditReady: "Audit ready",
          alphaCluster: "Alpha Cluster",
          alphaNote:
            "Ledger synced. Dashboard cache warm. Finance hub ready for board traffic.",
          posture: "Current posture",
          postureNote:
            "Run finance close, review segment mix, and push reports from one surface.",
          signalPanel: "Control posture",
          signalBullets: [
            "Realtime feed prioritized for the board.",
            "60-second cache to balance speed and stability.",
            "Ready to push reporting and audit workflows from one surface.",
          ],
        };

  return (
    <aside className="sticky top-4 hidden h-[calc(100vh-2rem)] w-[21.5rem] shrink-0 lg:block">
      <div className="relative flex h-full flex-col overflow-hidden rounded-[2rem] border border-white/12 bg-[radial-gradient(circle_at_top_left,rgba(14,165,233,0.18),transparent_30%),radial-gradient(circle_at_bottom_right,rgba(16,185,129,0.12),transparent_28%),linear-gradient(180deg,rgba(14,30,58,0.99),rgba(8,19,39,0.99))] text-white shadow-[0_30px_90px_rgba(6,16,32,0.34)]">
        <div className="absolute -left-14 top-10 h-36 w-36 rounded-full bg-sky-400/15 blur-3xl" />
        <div className="absolute -right-8 bottom-16 h-32 w-32 rounded-full bg-emerald-400/12 blur-3xl" />

        <div className="relative flex h-full flex-col px-5 py-5">
          <div className="space-y-5">
            <div className="rounded-[1.7rem] border border-white/10 bg-white/6 p-4 shadow-[inset_0_1px_0_rgba(255,255,255,0.06)]">
              <div className="flex items-start gap-4">
                <span className="inline-flex h-13 w-13 shrink-0 items-center justify-center rounded-[1.15rem] border border-white/14 bg-[linear-gradient(135deg,rgba(255,255,255,0.18),rgba(255,255,255,0.06))] text-white shadow-[0_16px_34px_rgba(7,16,32,0.28)]">
                  <Command className="size-5" />
                </span>
                <div className="min-w-0">
                  <p className="text-[0.64rem] uppercase tracking-[0.34em] text-white/42">
                    {copy.module}
                  </p>
                  <h2 className="mt-2 text-[1.65rem] font-semibold tracking-tight text-white">
                    {copy.title}
                  </h2>
                  <p className="mt-2 text-sm leading-6 text-white/62">
                    {copy.description}
                  </p>
                </div>
              </div>

              <div className="mt-4 grid grid-cols-2 gap-2.5">
                <div className="rounded-[1.15rem] border border-white/8 bg-white/8 px-3 py-3">
                  <div className="flex items-center gap-2 text-emerald-200">
                    <Wifi className="size-4" />
                    <span className="text-[0.68rem] uppercase tracking-[0.24em]">
                      {copy.liveSync}
                    </span>
                  </div>
                  <p className="mt-2 text-sm font-semibold text-white">
                    WebSocket
                  </p>
                </div>
                <div className="rounded-[1.15rem] border border-white/8 bg-white/8 px-3 py-3">
                  <div className="flex items-center gap-2 text-sky-200">
                    <Fingerprint className="size-4" />
                    <span className="text-[0.68rem] uppercase tracking-[0.24em]">
                      {copy.auditReady}
                    </span>
                  </div>
                  <p className="mt-2 text-sm font-semibold text-white">
                    Traceable
                  </p>
                </div>
              </div>
            </div>

            <div className="rounded-[1.45rem] border border-white/10 bg-white/6 p-3">
              <div className="flex items-center gap-3 rounded-[1.15rem] border border-white/8 bg-[linear-gradient(135deg,rgba(255,255,255,0.08),rgba(255,255,255,0.03))] px-3.5 py-3 text-sm text-white/72">
                <Search className="size-4 text-white/54" />
                <span className="line-clamp-1">{copy.search}</span>
              </div>
              <div className="mt-3 flex items-center gap-2 rounded-[1.05rem] border border-emerald-400/16 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">
                <Sparkles className="size-4" />
                <div>
                  <p className="font-medium">{copy.readyState}</p>
                  <p className="text-xs text-emerald-100/75">{copy.readySub}</p>
                </div>
              </div>
            </div>

            <div className="space-y-4">
              <div className="rounded-[1.55rem] border border-white/10 bg-white/6 p-3.5">
                <p className="px-1.5 text-[0.64rem] uppercase tracking-[0.3em] text-white/42">
                  {copy.command}
                </p>
                <nav className="mt-3 space-y-2">
                  {primary.map((item) => (
                    <DashboardNavLink
                      key={item.href}
                      href={item.href}
                      label={item.label}
                      caption={item.caption}
                      icon={item.icon}
                    />
                  ))}
                </nav>
              </div>

              <div className="rounded-[1.55rem] border border-white/10 bg-white/6 p-3.5">
                <p className="px-1.5 text-[0.64rem] uppercase tracking-[0.3em] text-white/42">
                  {copy.assurance}
                </p>
                <nav className="mt-3 space-y-2">
                  {secondary.map((item) => (
                    <DashboardNavLink
                      key={item.href}
                      href={item.href}
                      label={item.label}
                      caption={item.caption}
                      icon={item.icon}
                    />
                  ))}
                </nav>
              </div>
            </div>
          </div>

          <div className="mt-auto space-y-3 pt-4">
            <div className="rounded-[1.5rem] border border-emerald-400/18 bg-[linear-gradient(180deg,rgba(16,185,129,0.18),rgba(16,185,129,0.08))] p-4 text-sm text-emerald-100 shadow-[0_16px_34px_rgba(5,12,24,0.18)]">
              <div className="flex items-center justify-between gap-3">
                <div className="flex items-center gap-2 text-emerald-200">
                  <Activity className="size-4" />
                  <p className="text-xs uppercase tracking-[0.24em]">
                    {copy.alphaCluster}
                  </p>
                </div>
                <span className="rounded-full border border-emerald-300/18 bg-emerald-200/10 px-2.5 py-1 text-[0.64rem] uppercase tracking-[0.22em] text-emerald-100">
                  Live
                </span>
              </div>
              <p className="mt-3 font-medium leading-6">{copy.alphaNote}</p>
            </div>

            <div className="rounded-[1.45rem] border border-white/10 bg-white/6 p-4 text-sm text-white/72">
              <div className="flex items-center gap-2 text-white/84">
                <ShieldCheck className="size-4" />
                <p className="text-xs uppercase tracking-[0.24em]">
                  {copy.signalPanel}
                </p>
              </div>
              <p className="mt-3 font-medium text-white">{copy.postureNote}</p>
              <div className="mt-4 space-y-2.5">
                {copy.signalBullets.map((item) => (
                  <div key={item} className="flex items-start gap-2">
                    <CircleDot className="mt-1 size-3.5 shrink-0 text-white/38" />
                    <p className="leading-6 text-white/64">{item}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </aside>
  );
}
