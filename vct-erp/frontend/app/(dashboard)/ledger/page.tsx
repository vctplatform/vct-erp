import { DataTable } from "@/components/finance/data-table";
import { InsightCard } from "@/components/finance/insight-card";
import { ModuleHeader } from "@/components/finance/module-header";
import { SectionPanel } from "@/components/finance/section-panel";
import { StatusPill } from "@/components/finance/status-pill";
import { getFinanceDashboardSnapshot } from "@/lib/api/finance";
import {
  buildLedgerLanes,
  latestLabel,
  requireDashboardCard,
  statusTone,
} from "@/lib/finance/insights";
import { formatCompactCurrency, formatPercent } from "@/lib/formatters";
import {
  translateFinanceRefreshMode,
  translateFinanceSegment,
  translateFinanceStatus,
} from "@/lib/i18n/finance";
import { getServerLocale } from "@/lib/i18n/server";

export default async function LedgerPage() {
  const locale = await getServerLocale();
  const snapshot = await getFinanceDashboardSnapshot(locale);
  const cashCard = requireDashboardCard(snapshot, "cash_assets");
  const revenueCard = requireDashboardCard(snapshot, "quarter_net_revenue");
  const runwayCard = requireDashboardCard(snapshot, "runway_index");
  const lanes = buildLedgerLanes(snapshot, locale);

  const copy =
    locale === "vi"
      ? {
          kicker: "Kiểm soát sổ cái lõi",
          title: "Vận hành sổ cái",
          description:
            "Bề mặt làm việc cho đội tài chính: quan sát các làn ghi sổ, rà soát chất lượng chứng từ và nhìn nơi khóa sổ có thể trôi.",
          journalLive: "nhật ký trực tiếp",
          runway: "runway",
          treasuryPosted: "Ngân quỹ đã ghi sổ",
          treasuryPostedCaption:
            "Số dư tiền và ngân hàng hiện đang phản ánh trong chế độ xem sổ cái live.",
          quarterVouchers: "Chứng từ quý",
          quarterVouchersCaption:
            "Chỉ số đếm kiểm soát mô phỏng cho các làn ghi sổ đang hoạt động trên bề mặt này.",
          revenueWatch: "Doanh thu cần theo dõi",
          revenueWatchCaption:
            "Doanh thu đã ghi nhận hiện đang định hình câu chuyện khóa sổ.",
          partitionWatch: "Theo dõi partition",
          partitionWatchCaption:
            "Lượt seed hiện tại đang rơi vào journal_items partition đang hoạt động.",
          postingLanes: "Làn ghi sổ",
          cashEnters: "Tiền đi vào sổ cái như thế nào",
          cashEntersDescription:
            "Mỗi làn vận hành ánh xạ một dòng tiền kinh doanh sang mẫu chứng từ mà tài chính cần theo dõi.",
          clean: "sạch",
          lane: "luồng",
          reviewQueue: "Hàng đợi rà soát",
          voucherPackets: "Gói chứng từ cần kiểm tra",
          voucherPacketsDescription:
            "Một phiên bản an toàn cho Ban điều hành của màn hình rà soát nhật ký mà đội tài chính dùng trước khi khóa sổ.",
          voucher: "Chứng từ",
          laneColumn: "Làn",
          amount: "Giá trị",
          controlRule: "Quy tắc kiểm soát",
          confidence: "Độ tin cậy",
          closeNotes: "Ghi chú khóa sổ",
          matters: "Những điều vẫn quan trọng",
          mattersDescription:
            "Những lời nhắc ngắn giữ chất lượng sổ cái không suy giảm khi tốc độ tăng lên.",
          notes: [
            "Bảo vệ idempotency trên các route capture trước khi retry lại bất kỳ giao dịch nào từ retail hoặc võ đường.",
            "Giữ hiển thị semantics reversal trong suốt quá trình khóa sổ để số tiền báo cáo không trôi khỏi góc nhìn kiểm toán.",
            "Dùng trang đối soát trước khi xuất các báo cáo hướng thuế từ bộ báo cáo.",
          ],
        }
      : {
          kicker: "Core ledger control",
          title: "Ledger Ops",
          description:
            "A working surface for finance operations: watch posting lanes, review voucher quality, and see where the close can still drift.",
          journalLive: "journal live",
          runway: "runway",
          treasuryPosted: "Treasury posted",
          treasuryPostedCaption:
            "Cash and bank balances currently represented in the live ledger view.",
          quarterVouchers: "Quarter vouchers",
          quarterVouchersCaption:
            "A synthetic control count for active posting lanes in this UI surface.",
          revenueWatch: "Revenue under watch",
          revenueWatchCaption:
            "Recognized revenue currently shaping the close narrative.",
          partitionWatch: "Partition watch",
          partitionWatchCaption:
            "Current seed run is landing into the active journal_items partition.",
          postingLanes: "Posting lanes",
          cashEnters: "How cash enters the ledger",
          cashEntersDescription:
            "Each operating lane maps a business flow to the voucher pattern finance needs to monitor.",
          clean: "clean",
          lane: "lane",
          reviewQueue: "Review queue",
          voucherPackets: "Voucher packets to inspect",
          voucherPacketsDescription:
            "A board-safe representation of the journal review surface the finance team would use before close.",
          voucher: "Voucher",
          laneColumn: "Lane",
          amount: "Amount",
          controlRule: "Control rule",
          confidence: "Confidence",
          closeNotes: "Close notes",
          matters: "What still matters",
          mattersDescription:
            "Short reminders that keep ledger quality from degrading when speed goes up.",
          notes: [
            "Protect idempotency on capture routes before retrying anything from retail or dojo.",
            "Keep reversal semantics visible during close so reported cash does not drift from the audit view.",
            "Use the reconciliation page before exporting tax-facing reports from the reporting suite.",
          ],
        };

  const reviewRows = lanes.map((lane, index) => ({
    voucher: `VC-${latestLabel(snapshot).replace("-", "")}-${String(index + 1).padStart(3, "0")}`,
    lane: translateFinanceSegment(lane.segment, locale),
    amount: formatCompactCurrency(lane.amount * 0.11, locale),
    rule: lane.checkpoint,
    confidence: (
      <StatusPill tone={lane.confidence >= 98 ? "emerald" : "amber"}>
        {lane.confidence}% {locale === "vi" ? "tự động" : "auto"}
      </StatusPill>
    ),
  }));

  return (
    <section className="space-y-6">
      <ModuleHeader
        kicker={copy.kicker}
        title={copy.title}
        description={copy.description}
        mode={{
          label:
            translateFinanceRefreshMode(snapshot.recommended_refresh, locale) ??
            snapshot.recommended_refresh,
          tone: statusTone(snapshot.data_mode),
        }}
        actions={
          <>
            <StatusPill tone="emerald">{copy.journalLive}</StatusPill>
            <StatusPill tone={statusTone(runwayCard.status)}>
              {copy.runway} {translateFinanceStatus(runwayCard.status, locale)}
            </StatusPill>
          </>
        }
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <InsightCard
          label={copy.treasuryPosted}
          value={formatCompactCurrency(cashCard.value, locale)}
          caption={copy.treasuryPostedCaption}
          tone="navy"
          trend={{
            direction: cashCard.trend.direction,
            label: formatPercent(cashCard.trend.percentage),
          }}
        />
        <InsightCard
          label={copy.quarterVouchers}
          value={`${lanes.length * 48}`}
          caption={copy.quarterVouchersCaption}
          tone="emerald"
        />
        <InsightCard
          label={copy.revenueWatch}
          value={formatCompactCurrency(revenueCard.value, locale)}
          caption={copy.revenueWatchCaption}
          tone="amber"
        />
        <InsightCard
          label={copy.partitionWatch}
          value="2026-Q1"
          caption={copy.partitionWatchCaption}
          tone="navy"
        />
      </div>

      <SectionPanel
        kicker={copy.postingLanes}
        title={copy.cashEnters}
        description={copy.cashEntersDescription}
      >
        <div className="grid gap-4 md:grid-cols-2">
          {lanes.map((lane) => (
            <div
              key={lane.segment}
              className="rounded-[1.4rem] border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4"
            >
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-xs uppercase tracking-[0.24em] text-[var(--color-ink-soft)]">
                    {translateFinanceSegment(lane.segment, locale)}
                  </p>
                  <p className="mt-2 text-xl font-semibold text-[var(--color-ink)]">
                    {lane.voucherType} {copy.lane}
                  </p>
                </div>
                <StatusPill tone={lane.confidence >= 98 ? "emerald" : "amber"}>
                  {lane.confidence}% {copy.clean}
                </StatusPill>
              </div>
              <p className="mt-4 text-2xl font-semibold tracking-tight text-[var(--color-ink)]">
                {formatCompactCurrency(lane.amount, locale)}
              </p>
              <p className="mt-3 text-sm leading-6 text-[var(--color-ink-soft)]">
                {lane.checkpoint}
              </p>
            </div>
          ))}
        </div>
      </SectionPanel>

      <div className="grid gap-6 xl:grid-cols-[1.3fr_0.7fr]">
        <SectionPanel
          kicker={copy.reviewQueue}
          title={copy.voucherPackets}
          description={copy.voucherPacketsDescription}
        >
          <DataTable
            columns={[
              { key: "voucher", label: copy.voucher },
              { key: "lane", label: copy.laneColumn },
              { key: "amount", label: copy.amount, align: "right" },
              { key: "rule", label: copy.controlRule },
              { key: "confidence", label: copy.confidence, align: "right" },
            ]}
            rows={reviewRows}
          />
        </SectionPanel>

        <SectionPanel
          kicker={copy.closeNotes}
          title={copy.matters}
          description={copy.mattersDescription}
        >
          <div className="space-y-4 text-sm leading-6 text-[var(--color-ink-soft)]">
            {copy.notes.map((note) => (
              <div
                key={note}
                className="rounded-[1.2rem] border border-[var(--color-border)] bg-[var(--color-canvas-soft)] p-4"
              >
                {note}
              </div>
            ))}
          </div>
        </SectionPanel>
      </div>
    </section>
  );
}
