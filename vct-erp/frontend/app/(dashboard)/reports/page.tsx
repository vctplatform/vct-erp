import { DataTable } from "@/components/finance/data-table";
import { InsightCard } from "@/components/finance/insight-card";
import { ModuleHeader } from "@/components/finance/module-header";
import { SectionPanel } from "@/components/finance/section-panel";
import { StatusPill } from "@/components/finance/status-pill";
import { getFinanceDashboardSnapshot } from "@/lib/api/finance";
import {
  latestSeriesValue,
  requireDashboardCard,
  statusTone,
  totalRevenue,
} from "@/lib/finance/insights";
import { formatCompactCurrency } from "@/lib/formatters";
import {
  translateFinanceDataMode,
  translateFinanceRefreshMode,
  translateFinanceStatus,
} from "@/lib/i18n/finance";
import { getServerLocale } from "@/lib/i18n/server";

export default async function ReportsPage() {
  const locale = await getServerLocale();
  const snapshot = await getFinanceDashboardSnapshot(locale);
  const cashCard = requireDashboardCard(snapshot, "cash_assets");
  const runwayCard = requireDashboardCard(snapshot, "runway_index");
  const revenueTotal = totalRevenue(snapshot);
  const latestProfit = latestSeriesValue(snapshot, "profit");

  const copy =
    locale === "vi"
      ? {
          kicker: "Bộ báo cáo",
          title: "Báo cáo & Thuế",
          description:
            "Lớp xuất bản của mô-đun tài chính: biến sự thật live của sổ cái thành bộ quản trị, báo cáo pháp định và một kỳ khóa sổ bớt áp lực.",
          balanced: "bảng cân đối sạch",
          runway: "runway",
          reportedCash: "Tiền đã báo cáo",
          reportedCashCaption:
            "Nền ngân quỹ hiện đã sẵn sàng cho các lớp báo cáo.",
          bookedRevenue: "Doanh thu đã ghi nhận",
          bookedRevenueCaption:
            "Dấu chân doanh thu liên mảng sẵn cho người dùng báo cáo.",
          latestProfit: "Lợi nhuận mới nhất",
          latestProfitCaption:
            "Dòng lợi nhuận gần nhất sẵn cho lớp báo cáo diễn giải.",
          runwayNarrative: "Câu chuyện runway",
          runwayNarrativeCaption:
            "Một tín hiệu chỉ dành cho quản trị, bổ trợ nhưng không thay thế báo cáo pháp định.",
          publishStack: "Ngăn xếp xuất bản",
          reportDeck: "Bộ báo cáo",
          reportDeckDescription:
            "Mọi người dùng báo cáo nên biết artifact nào đang live, đã soát hay còn chờ ký duyệt tài chính.",
          report: "Báo cáo",
          cadence: "Nhịp",
          owner: "Phụ trách",
          status: "Trạng thái",
          signal: "Tín hiệu",
          closingCalendar: "Lịch khóa sổ",
          whatNext: "Việc phải diễn ra tiếp theo",
          whatNextDescription:
            "Một lịch thân thiện với Ban điều hành để quyết định tài chính hạ cánh trước hạn chót.",
          milestone: "Mốc",
          due: "Hạn",
          posture: "Tư thế",
          noData: "chưa có",
        }
      : {
          kicker: "Reporting suite",
          title: "Reports & Tax",
          description:
            "The publishing layer of the finance module: turn live ledger truth into management packs, statutory statements, and a calmer close.",
          balanced: "trial balance clean",
          runway: "runway",
          reportedCash: "Reported cash",
          reportedCashCaption:
            "Current treasury base already ready for reporting surfaces.",
          bookedRevenue: "Booked revenue",
          bookedRevenueCaption:
            "Cross-segment revenue footprint available to report consumers.",
          latestProfit: "Latest profit",
          latestProfitCaption:
            "Most recent profit line available to narrative reporting.",
          runwayNarrative: "Runway narrative",
          runwayNarrativeCaption:
            "A management-only signal that complements but does not replace statutory reporting.",
          publishStack: "Publish stack",
          reportDeck: "Report deck",
          reportDeckDescription:
            "Every reporting consumer should know which artifact is live, reviewed, or still waiting on finance sign-off.",
          report: "Report",
          cadence: "Cadence",
          owner: "Owner",
          status: "Status",
          signal: "Signal",
          closingCalendar: "Closing calendar",
          whatNext: "What has to happen next",
          whatNextDescription:
            "A board-friendly calendar so finance decisions land before the deadline, not after.",
          milestone: "Milestone",
          due: "Due",
          posture: "Posture",
          noData: "n/a",
        };

  const reportRows = [
    {
      report:
        locale === "vi" ? "Bảng cân đối số phát sinh" : "Trial Balance",
      cadence: locale === "vi" ? "Khóa sổ hằng ngày" : "Daily close",
      owner: locale === "vi" ? "Bàn sổ cái" : "GL desk",
      status: (
        <StatusPill tone="emerald">
          {translateFinanceStatus("balanced", locale)}
        </StatusPill>
      ),
      signal: formatCompactCurrency(cashCard.value, locale),
    },
    {
      report: "P&L B02-DN",
      cadence: locale === "vi" ? "Hàng tháng" : "Monthly",
      owner: locale === "vi" ? "Văn phòng CFO" : "CFO office",
      status: (
        <StatusPill tone="emerald">
          {translateFinanceStatus("ready", locale)}
        </StatusPill>
      ),
      signal: formatCompactCurrency(latestProfit, locale),
    },
    {
      report: locale === "vi" ? "Sổ nhật ký chung" : "General Journal",
      cadence: locale === "vi" ? "Liên tục" : "Continuous",
      owner: locale === "vi" ? "Vận hành tài chính" : "Finance ops",
      status: (
        <StatusPill tone="navy">
          {translateFinanceStatus("live", locale)}
        </StatusPill>
      ),
      signal: (
        translateFinanceRefreshMode(snapshot.recommended_refresh, locale) ??
        snapshot.recommended_refresh
      ).toUpperCase(),
    },
    {
      report: locale === "vi" ? "Hồ sơ thuế" : "Tax Pack",
      cadence: locale === "vi" ? "Cuối tháng" : "Month-end",
      owner: locale === "vi" ? "Tuân thủ" : "Compliance",
      status: (
        <StatusPill tone="amber">
          {translateFinanceStatus("review", locale)}
        </StatusPill>
      ),
      signal: "VAT + TNDN",
    },
  ];

  const calendarRows = [
    {
      milestone:
        locale === "vi"
          ? "Chốt ảnh chụp cho Ban điều hành"
          : "Board snapshot freeze",
      due: "2026-03-25 18:00",
      owner: locale === "vi" ? "PMO tài chính" : "Finance PMO",
      posture: (
        <StatusPill tone="navy">
          {translateFinanceStatus("scheduled", locale)}
        </StatusPill>
      ),
    },
    {
      milestone:
        locale === "vi"
          ? "Ký duyệt bảng cân đối số phát sinh"
          : "Trial balance sign-off",
      due: "2026-03-27 11:00",
      owner: locale === "vi" ? "Kế toán trưởng" : "Chief accountant",
      posture: (
        <StatusPill tone="emerald">
          {translateFinanceStatus("on track", locale)}
        </StatusPill>
      ),
    },
    {
      milestone:
        locale === "vi" ? "Rà soát bộ hồ sơ thuế" : "Tax file package review",
      due: "2026-03-29 16:30",
      owner: locale === "vi" ? "Tuân thủ" : "Compliance",
      posture: (
        <StatusPill tone="amber">
          {translateFinanceStatus("watch", locale)}
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
            <StatusPill tone="emerald">{copy.balanced}</StatusPill>
            <StatusPill tone={statusTone(runwayCard.status)}>
              {copy.runway}{" "}
              {translateFinanceStatus(runwayCard.status, locale) ?? copy.noData}
            </StatusPill>
          </>
        }
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <InsightCard
          label={copy.reportedCash}
          value={formatCompactCurrency(cashCard.value, locale)}
          caption={copy.reportedCashCaption}
          tone="navy"
        />
        <InsightCard
          label={copy.bookedRevenue}
          value={formatCompactCurrency(revenueTotal, locale)}
          caption={copy.bookedRevenueCaption}
          tone="emerald"
        />
        <InsightCard
          label={copy.latestProfit}
          value={formatCompactCurrency(latestProfit, locale)}
          caption={copy.latestProfitCaption}
          tone={latestProfit >= 0 ? "emerald" : "rose"}
        />
        <InsightCard
          label={copy.runwayNarrative}
          value={
            locale === "vi"
              ? `${runwayCard.value.toFixed(1)} tháng`
              : `${runwayCard.value.toFixed(1)} months`
          }
          caption={copy.runwayNarrativeCaption}
          tone={statusTone(runwayCard.status)}
        />
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
        <SectionPanel
          kicker={copy.publishStack}
          title={copy.reportDeck}
          description={copy.reportDeckDescription}
        >
          <DataTable
            columns={[
              { key: "report", label: copy.report },
              { key: "cadence", label: copy.cadence },
              { key: "owner", label: copy.owner },
              { key: "status", label: copy.status, align: "right" },
              { key: "signal", label: copy.signal, align: "right" },
            ]}
            rows={reportRows}
          />
        </SectionPanel>

        <SectionPanel
          kicker={copy.closingCalendar}
          title={copy.whatNext}
          description={copy.whatNextDescription}
        >
          <DataTable
            compact
            columns={[
              { key: "milestone", label: copy.milestone },
              { key: "due", label: copy.due },
              { key: "owner", label: copy.owner },
              { key: "posture", label: copy.posture, align: "right" },
            ]}
            rows={calendarRows}
          />
        </SectionPanel>
      </div>
    </section>
  );
}
