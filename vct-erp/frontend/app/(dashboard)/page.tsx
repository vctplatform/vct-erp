import { DashboardClient } from "@/components/dashboard/dashboard-client";
import { DataTable } from "@/components/finance/data-table";
import { InsightCard } from "@/components/finance/insight-card";
import { ModuleHeader } from "@/components/finance/module-header";
import { SectionPanel } from "@/components/finance/section-panel";
import { StatusPill } from "@/components/finance/status-pill";
import { getFinanceDashboardSnapshot } from "@/lib/api/finance";
import {
  averageSeriesValue,
  buildSegmentShares,
  latestLabel,
  latestRunwayEnding,
  latestSeriesValue,
  requireDashboardCard,
  statusTone,
  strongestSegment,
  totalRevenue,
} from "@/lib/finance/insights";
import { formatCompactCurrency, formatPercent } from "@/lib/formatters";
import {
  translateFinanceDataMode,
  translateFinanceSegment,
  translateFinanceStatus,
} from "@/lib/i18n/finance";
import { getServerLocale } from "@/lib/i18n/server";

export default async function DashboardPage() {
  const locale = await getServerLocale();
  const snapshot = await getFinanceDashboardSnapshot(locale);
  const cashCard = requireDashboardCard(snapshot, "cash_assets");
  const revenueCard = requireDashboardCard(snapshot, "quarter_net_revenue");
  const runwayCard = requireDashboardCard(snapshot, "runway_index");
  const segmentShares = buildSegmentShares(snapshot);
  const leadSegment = strongestSegment(snapshot);
  const latestProfit = latestSeriesValue(snapshot, "profit");
  const averageExpense = averageSeriesValue(snapshot, "expense");
  const futureCash = latestRunwayEnding(snapshot);
  const totalBookedRevenue = totalRevenue(snapshot);

  const copy =
    locale === "vi"
      ? {
          kicker: "Trí tuệ tài chính Tập đoàn VCT",
          title: "Trung tâm điều hành",
          description:
            "Bề mặt điều hành cho mô-đun tài chính: theo dõi tư thế tiền mặt trực tiếp, động lượng từng mảng, độ sẵn sàng của sổ cái và những gì cần can thiệp trước khi Ban điều hành hỏi.",
          runway: "runway",
          operatingRhythm: "Nhịp vận hành",
          boardRead: "Điều Ban điều hành nên đọc trước",
          boardReadDescription:
            "Một lớp vận hành ngắn gọn giúp dịch tín hiệu live của dashboard thành ưu tiên tài chính trong cửa sổ khóa sổ hiện tại.",
          bookedRevenue: "Doanh thu đã ghi nhận",
          topContributorPrefix: "Mảng đóng góp lớn nhất",
          cashPosture: "Tư thế tiền mặt",
          cashPostureCaption:
            "Số dư ngân quỹ đang sẵn dùng trên tài khoản tiền gửi và tiền mặt.",
          latestProfitPulse: "Xung lợi nhuận mới nhất",
          latestProfitCaption:
            "Tín hiệu lợi nhuận tháng của chu kỳ hiện tại sau hấp thụ chi phí.",
          forwardEndingCash: "Tiền cuối kỳ dự phóng",
          forwardEndingCashCaption:
            "Số dư tiền mặt dự kiến ở cuối chân trời dự báo 6 tháng hiện tại.",
          boardAgenda: "Chương trình Ban điều hành",
          tonightMoves: "Các động tác tài chính tối nay",
          tonightMovesDescription:
            "Những hành động ưu tiên để chuyển tín hiệu live thành quyết định vận hành.",
          activePlays: "4 ưu tiên đang chạy",
          focus: "Trọng tâm",
          owner: "Phụ trách",
          window: "Cửa sổ",
          posture: "Tư thế",
          segmentPulse: "Nhịp mảng kinh doanh",
          concentration: "Mức tập trung danh mục",
          concentrationDescription:
            "Đọc nhanh xem mảng nào đang gánh quý này và quản trị nên tập trung vào đâu.",
          segmentShare: "Tỷ trọng mảng",
          segmentShareSuffix: "trong doanh thu danh mục hiện tại.",
          costDiscipline: "Kỷ luật chi phí",
          expenseFrame: "Khung chi phí",
          expenseFrameDescription:
            "Mức chi bình quân lấy từ đồ thị dòng tiền live, hữu ích khi bàn về chất lượng burn thay vì chỉ nhìn runway.",
          averageMonthlyExpense: "Chi phí bình quân tháng",
          averageMonthlyExpenseCaption:
            "Trung bình 6 tháng gần nhất từ biểu đồ dashboard.",
          quarterRevenue: "Doanh thu quý",
          quarterRevenueCaption:
            "Doanh thu thuần đã được ghi nhận trong quý hiện tại.",
          runwayPosture: "Tư thế runway",
          runwayPostureCaption:
            "Góc nhìn cấp điều hành cho các thảo luận về khẩu vị rủi ro và nhịp tăng trưởng.",
          noData: "chưa có",
        }
      : {
          kicker: "VCT Group Financial Intelligence",
          title: "Command Center",
          description:
            "Executive surface for the finance module: watch live cash posture, segment momentum, ledger readiness, and what needs intervention before the board asks.",
          runway: "runway",
          operatingRhythm: "Operating rhythm",
          boardRead: "What the board should read first",
          boardReadDescription:
            "A compact operating layer that translates the live dashboard into finance priorities for the current close window.",
          bookedRevenue: "Booked revenue",
          topContributorPrefix: "Top contributor",
          cashPosture: "Cash posture",
          cashPostureCaption:
            "Treasury balance available across bank and cash accounts.",
          latestProfitPulse: "Latest profit pulse",
          latestProfitCaption:
            "Monthly profit signal for the current cycle after expense absorption.",
          forwardEndingCash: "Forward ending cash",
          forwardEndingCashCaption:
            "Projected ending balance at the end of the current 6-month forecast horizon.",
          boardAgenda: "Board agenda",
          tonightMoves: "Tonight's finance moves",
          tonightMovesDescription:
            "Priority actions that translate live signals into operational decisions.",
          activePlays: "4 active plays",
          focus: "Focus",
          owner: "Owner",
          window: "Window",
          posture: "Posture",
          segmentPulse: "Segment pulse",
          concentration: "Portfolio concentration",
          concentrationDescription:
            "A quick read on which business lines are carrying the quarter and where management attention should concentrate.",
          segmentShare: "Segment share",
          segmentShareSuffix: "of current portfolio revenue.",
          costDiscipline: "Cost discipline",
          expenseFrame: "Expense frame",
          expenseFrameDescription:
            "Average expense load from the live cashflow chart, useful for talking about burn quality rather than only raw runway.",
          averageMonthlyExpense: "Average monthly expense",
          averageMonthlyExpenseCaption:
            "Trailing 6-month average from the dashboard chart.",
          quarterRevenue: "Quarter revenue",
          quarterRevenueCaption:
            "Net revenue already recognized in the current quarter.",
          runwayPosture: "Runway posture",
          runwayPostureCaption:
            "A high-level view for board discussions around risk appetite and growth pacing.",
          noData: "n/a",
        };

  const agendaSource =
    locale === "vi"
      ? [
          ["Khóa doanh thu chưa thực hiện", "Bàn doanh thu", latestLabel(snapshot), "ready"],
          ["Đối soát tài khoản ngân quỹ", "Kiểm soát sổ cái", "1121 + 1111", "watch"],
          ["Rà soát giảm trừ bán lẻ", "Tài chính bán lẻ", "Biến động 5211", "queue"],
          ["Làm mới dự báo cho Ban điều hành", "Văn phòng CFO", "Runway 6 tháng", "live"],
        ]
      : [
          ["Close deferred revenue", "Revenue desk", latestLabel(snapshot), "ready"],
          ["Reconcile treasury accounts", "GL control", "1121 + 1111", "watch"],
          ["Review retail deductions", "Retail finance", "5211 variance", "queue"],
          ["Refresh board forecast", "CFO office", "6-month runway", "live"],
        ];

  const agendaRows = agendaSource.map(([focus, owner, window, posture]) => ({
    focus,
    owner,
    window,
    posture: (
      <StatusPill tone={statusTone(posture)}>
        {translateFinanceStatus(posture, locale)}
      </StatusPill>
    ),
  }));

  const leadSegmentLabel =
    translateFinanceSegment(leadSegment?.label, locale) ?? copy.noData;

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
            <StatusPill tone={statusTone(runwayCard.status)}>
              {copy.runway}{" "}
              {translateFinanceStatus(runwayCard.status, locale) ?? copy.noData}
            </StatusPill>
            <StatusPill tone="navy">{snapshot.company_code}</StatusPill>
          </>
        }
      />

      <DashboardClient initialData={snapshot} />

      <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <SectionPanel
          kicker={copy.operatingRhythm}
          title={copy.boardRead}
          description={copy.boardReadDescription}
        >
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            <InsightCard
              label={copy.bookedRevenue}
              value={formatCompactCurrency(totalBookedRevenue, locale)}
              caption={
                locale === "vi"
                  ? `${copy.topContributorPrefix} ${leadSegmentLabel} với ${formatPercent(leadSegment?.share ?? 0)} cơ cấu doanh thu.`
                  : `${copy.topContributorPrefix} ${leadSegmentLabel} with ${formatPercent(leadSegment?.share ?? 0)} of revenue mix.`
              }
              tone="emerald"
              trend={{
                direction: revenueCard.trend.direction,
                label: formatPercent(revenueCard.trend.percentage),
              }}
            />
            <InsightCard
              label={copy.cashPosture}
              value={formatCompactCurrency(cashCard.value, locale)}
              caption={copy.cashPostureCaption}
              tone="navy"
              trend={{
                direction: cashCard.trend.direction,
                label: formatPercent(cashCard.trend.percentage),
              }}
            />
            <InsightCard
              label={copy.latestProfitPulse}
              value={formatCompactCurrency(latestProfit, locale)}
              caption={
                locale === "vi"
                  ? `${copy.latestProfitCaption} Kỳ ${latestLabel(snapshot)}.`
                  : `${copy.latestProfitCaption} For ${latestLabel(snapshot)}.`
              }
              tone={latestProfit >= 0 ? "emerald" : "rose"}
            />
            <InsightCard
              label={copy.forwardEndingCash}
              value={formatCompactCurrency(futureCash, locale)}
              caption={copy.forwardEndingCashCaption}
              tone="amber"
            />
          </div>
        </SectionPanel>

        <SectionPanel
          kicker={copy.boardAgenda}
          title={copy.tonightMoves}
          description={copy.tonightMovesDescription}
          aside={<StatusPill tone="amber">{copy.activePlays}</StatusPill>}
        >
          <DataTable
            compact
            columns={[
              { key: "focus", label: copy.focus },
              { key: "owner", label: copy.owner },
              { key: "window", label: copy.window },
              { key: "posture", label: copy.posture, align: "right" },
            ]}
            rows={agendaRows}
          />
        </SectionPanel>
      </div>

      <SectionPanel
        kicker={copy.segmentPulse}
        title={copy.concentration}
        description={copy.concentrationDescription}
      >
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {segmentShares.map((segment) => (
            <div
              key={segment.label}
              className="rounded-[1.4rem] border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4"
            >
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-sm font-medium text-[var(--color-ink)]">
                    {translateFinanceSegment(segment.label, locale)}
                  </p>
                  <p className="mt-1 text-xs uppercase tracking-[0.22em] text-[var(--color-ink-soft)]">
                    {copy.segmentShare}
                  </p>
                </div>
                <span
                  className="size-3 rounded-full"
                  style={{ backgroundColor: segment.color }}
                />
              </div>
              <p className="mt-4 text-2xl font-semibold tracking-tight text-[var(--color-ink)]">
                {formatCompactCurrency(segment.value, locale)}
              </p>
              <p className="mt-2 text-sm text-[var(--color-ink-soft)]">
                {formatPercent(segment.share)} {copy.segmentShareSuffix}
              </p>
            </div>
          ))}
        </div>
      </SectionPanel>

      <SectionPanel
        kicker={copy.costDiscipline}
        title={copy.expenseFrame}
        description={copy.expenseFrameDescription}
        aside={<StatusPill tone="navy">{latestLabel(snapshot)}</StatusPill>}
      >
        <div className="grid gap-4 md:grid-cols-3">
          <InsightCard
            label={copy.averageMonthlyExpense}
            value={formatCompactCurrency(averageExpense, locale)}
            caption={copy.averageMonthlyExpenseCaption}
            tone="amber"
          />
          <InsightCard
            label={copy.quarterRevenue}
            value={formatCompactCurrency(revenueCard.value, locale)}
            caption={copy.quarterRevenueCaption}
            tone="emerald"
          />
          <InsightCard
            label={copy.runwayPosture}
            value={
              locale === "vi"
                ? `${runwayCard.value.toFixed(1)} tháng`
                : `${runwayCard.value.toFixed(1)} months`
            }
            caption={copy.runwayPostureCaption}
            tone={statusTone(runwayCard.status)}
          />
        </div>
      </SectionPanel>
    </section>
  );
}
