import { DataTable } from "@/components/finance/data-table";
import { InsightCard } from "@/components/finance/insight-card";
import { ModuleHeader } from "@/components/finance/module-header";
import { SectionPanel } from "@/components/finance/section-panel";
import { StatusPill } from "@/components/finance/status-pill";
import { getFinanceDashboardSnapshot } from "@/lib/api/finance";
import { buildControlSignals, statusTone } from "@/lib/finance/insights";
import {
  formatFinanceDate,
  translateFinanceDataMode,
  translateFinanceRefreshMode,
  translateFinanceStatus,
} from "@/lib/i18n/finance";
import { getServerLocale } from "@/lib/i18n/server";

export default async function ControlRoomPage() {
  const locale = await getServerLocale();
  const snapshot = await getFinanceDashboardSnapshot(locale);
  const signals = buildControlSignals(snapshot, locale);

  const copy =
    locale === "vi"
      ? {
          kicker: "Bề mặt kiểm soát",
          title: "Phòng điều khiển",
          description:
            "Nơi lãnh đạo tài chính kiểm tra liệu mô-đun có đủ đáng tin để trình Ban điều hành hay chưa: chế độ runtime, cache, audit và giao thức phản ứng.",
          auditSurface: "Bề mặt kiểm toán",
          recentEvents: "Sự kiện kiểm soát gần đây",
          recentEventsDescription:
            "Các bản ghi đại diện theo phong cách audit cho thấy UI có chỗ cho quản trị, không chỉ cho biểu đồ.",
          timestamp: "Thời điểm",
          actor: "Tác nhân",
          action: "Hành động",
          result: "Kết quả",
          incidentProtocol: "Giao thức sự cố",
          ifSomethingDrifts: "Nếu có gì đó trôi",
          ifSomethingDriftsDescription:
            "Mẫu phản ứng mà mô-đun mong đợi khi độ tin cậy tài chính giảm xuống.",
          notes: [
            "Đóng băng export trình Ban điều hành nếu dashboard rơi khỏi live một cách bất ngờ.",
            "Chạy lại đối soát ngân quỹ trước khi ghi đè thủ công bất kỳ KPI nào hiển thị cho Ban điều hành.",
            "Ghi log người thao tác, lý do và số trước/sau mỗi khi tài chính vá một chỉ số đang hiển thị.",
          ],
        }
      : {
          kicker: "Control surface",
          title: "Control Room",
          description:
            "Where finance leadership checks whether the module is trustworthy enough to present to the board: runtime mode, cache, audit, and response protocol.",
          auditSurface: "Audit surface",
          recentEvents: "Recent control events",
          recentEventsDescription:
            "Representative audit-style entries that show the UI has a place for governance, not just charts.",
          timestamp: "Timestamp",
          actor: "Actor",
          action: "Action",
          result: "Result",
          incidentProtocol: "Incident protocol",
          ifSomethingDrifts: "If something drifts",
          ifSomethingDriftsDescription:
            "The response pattern the module expects when finance confidence drops.",
          notes: [
            "Freeze board exports if dashboard mode falls back from live unexpectedly.",
            "Re-run treasury reconciliation before manually overriding any board-visible KPI.",
            "Log the operator, the reason, and the before/after number whenever finance patches a displayed metric.",
          ],
        };

  const timestamp = formatFinanceDate(snapshot.generated_at, locale, {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  const auditRows = [
    {
      timestamp,
      actor: "ceo-001",
      action:
        locale === "vi"
          ? "Đã xem dashboard điều hành"
          : "Viewed command center dashboard",
      result: (
        <StatusPill tone="emerald">
          {translateFinanceStatus("logged", locale)}
        </StatusPill>
      ),
    },
    {
      timestamp,
      actor: "command-center",
      action:
        locale === "vi"
          ? "Đã làm mới proxy contract dashboard"
          : "Refreshed dashboard proxy contract",
      result: (
        <StatusPill tone="navy">
          {translateFinanceStatus("recorded", locale)}
        </StatusPill>
      ),
    },
    {
      timestamp,
      actor: "finance-hub",
      action:
        locale === "vi"
          ? "Kênh realtime đã sẵn sàng cho websocket broadcast"
          : "Realtime channel armed for websocket broadcast",
      result: (
        <StatusPill tone="emerald">
          {translateFinanceStatus("ready", locale)}
        </StatusPill>
      ),
    },
  ];

  return (
    <section className="space-y-6">
      <ModuleHeader
        kicker={copy.kicker}
        title={copy.title}
        description={copy.description}
        mode={{
          label:
            translateFinanceDataMode(snapshot.data_mode, locale) ??
            snapshot.data_mode,
          tone: statusTone(snapshot.data_mode),
        }}
        actions={
          <>
            <StatusPill tone="navy">{snapshot.company_code}</StatusPill>
            <StatusPill
              tone={
                snapshot.recommended_refresh === "websocket"
                  ? "emerald"
                  : "amber"
              }
            >
              {translateFinanceRefreshMode(snapshot.recommended_refresh, locale)}
            </StatusPill>
          </>
        }
      />

      <div className="grid gap-4 md:grid-cols-3">
        {signals.map((signal) => (
          <InsightCard
            key={signal.label}
            label={signal.label}
            value={signal.value}
            caption={signal.caption}
            tone={signal.tone}
          />
        ))}
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <SectionPanel
          kicker={copy.auditSurface}
          title={copy.recentEvents}
          description={copy.recentEventsDescription}
        >
          <DataTable
            columns={[
              { key: "timestamp", label: copy.timestamp },
              { key: "actor", label: copy.actor },
              { key: "action", label: copy.action },
              { key: "result", label: copy.result, align: "right" },
            ]}
            rows={auditRows}
          />
        </SectionPanel>

        <SectionPanel
          kicker={copy.incidentProtocol}
          title={copy.ifSomethingDrifts}
          description={copy.ifSomethingDriftsDescription}
        >
          <div className="space-y-3 text-sm leading-6 text-[var(--color-ink-soft)]">
            {copy.notes.map((item) => (
              <div
                key={item}
                className="rounded-[1.2rem] border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4"
              >
                {item}
              </div>
            ))}
          </div>
        </SectionPanel>
      </div>
    </section>
  );
}
