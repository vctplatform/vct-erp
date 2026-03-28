import { DataTable } from "@/components/finance/data-table";
import { InsightCard } from "@/components/finance/insight-card";
import { ModuleHeader } from "@/components/finance/module-header";
import { SectionPanel } from "@/components/finance/section-panel";
import { StatusPill } from "@/components/finance/status-pill";
import { getFinanceDashboardSnapshot } from "@/lib/api/finance";
import {
  buildSegmentShares,
  latestSeriesValue,
  statusTone,
  strongestSegment,
  weakestSegment,
} from "@/lib/finance/insights";
import { formatCompactCurrency, formatPercent } from "@/lib/formatters";
import {
  translateFinanceDataMode,
  translateFinanceSegment,
  translateFinanceStatus,
} from "@/lib/i18n/finance";
import { getServerLocale } from "@/lib/i18n/server";

export default async function SegmentsPage() {
  const locale = await getServerLocale();
  const snapshot = await getFinanceDashboardSnapshot(locale);
  const segments = buildSegmentShares(snapshot);
  const lead = strongestSegment(snapshot);
  const tail = weakestSegment(snapshot);
  const latestProfit = latestSeriesValue(snapshot, "profit");

  const segmentNotes: Record<string, string> =
    locale === "vi"
      ? {
          SaaS: "Động cơ thuê bao với doanh thu chưa thực hiện chuyển dần thành biên lợi nhuận dự báo được.",
          Dojo: "Dòng tiền hội viên đi kèm kỷ luật công nợ và độ nhạy cao với khóa sổ hàng tháng.",
          Retail: "Doanh thu tốc độ cao nhưng có giảm trừ và bút toán đảo POS cần theo dõi sát.",
          Rental: "Dòng thu gắn với tiền cọc, nơi khấu trừ hư hại tạo ra tín hiệu thu nhập khác.",
        }
      : {
          SaaS: "Subscription engine with deferred revenue turning into predictable margin.",
          Dojo: "Membership cashflow with receivable discipline and monthly close sensitivity.",
          Retail: "High velocity revenue with deductions and POS reversals needing close review.",
          Rental: "Deposit-driven income stream where damage offsets create other income signals.",
        };

  const segmentLevers: Record<string, string> =
    locale === "vi"
      ? {
          SaaS: "Đẩy nhanh gia hạn và giữ nhịp ghi nhận doanh thu.",
          Dojo: "Giảm tuổi nợ học phí và ổn định tỷ lệ sử dụng lớp.",
          Retail: "Thu hẹp giảm trừ và siết cửa sổ phê duyệt trả hàng.",
          Rental: "Nâng kỷ luật giải phóng cọc và tốc độ thu hồi hư hại.",
        }
      : {
          SaaS: "Accelerate renewals and protect recognition cadence.",
          Dojo: "Reduce aging receivables and stabilize class utilization.",
          Retail: "Trim deductions and tighten return authorization windows.",
          Rental: "Improve deposit release discipline and damage recovery speed.",
        };

  const copy =
    locale === "vi"
      ? {
          kicker: "Cơ cấu kinh doanh",
          title: "Studio mảng kinh doanh",
          description:
            "Góc nhìn danh mục cho ban quản trị: mảng nào đang gánh quý này, biên lợi nhuận tập trung ở đâu và mỗi mảng cần đòn bẩy gì tiếp theo.",
          profitPulse: "xung lợi nhuận",
          portfolioShare: "Tỷ trọng danh mục",
          playbook: "Sổ tay hành động",
          marginPlays: "Đòn bẩy biên lợi nhuận theo mảng",
          marginPlaysDescription:
            "Danh sách vận hành sắc hơn một biểu đồ tròn: nơi tài chính nên đẩy, giữ hay bảo vệ.",
          segment: "Mảng",
          revenue: "Doanh thu",
          share: "Tỷ trọng",
          nextLever: "Đòn bẩy tiếp theo",
          signal: "Tín hiệu",
          portfolioPosture: "Tư thế danh mục",
          managementKnow: "Điều ban quản trị cần biết",
          managementKnowDescription:
            "Một đoạn tóm tắt ngắn cho bộ hồ sơ trình ban điều hành.",
          leadContributor: "Mảng dẫn đầu",
          fragileSegment: "Mảng mong manh nhất",
          fragileCaption:
            "Đóng góp nhỏ nhất hiện tại nhưng vẫn quan trọng để đa dạng tín hiệu.",
          portfolioNote: "Ghi chú danh mục",
          portfolioNoteCaption:
            "Cơ cấu mảng hiện đã được nối vào cùng contract dashboard mà Ban điều hành đang xem.",
          noData: "chưa có",
        }
      : {
          kicker: "Business mix",
          title: "Segment Studio",
          description:
            "A portfolio view for management: which line is carrying the quarter, where margin is concentrated, and what lever each segment needs next.",
          profitPulse: "profit pulse",
          portfolioShare: "Portfolio share",
          playbook: "Playbook",
          marginPlays: "Margin plays by segment",
          marginPlaysDescription:
            "A sharper operational list than a pie chart: where finance should press, hold, or protect.",
          segment: "Segment",
          revenue: "Revenue",
          share: "Share",
          nextLever: "Next lever",
          signal: "Signal",
          portfolioPosture: "Portfolio posture",
          managementKnow: "What management should know",
          managementKnowDescription: "A concise narrative for the board packet.",
          leadContributor: "Lead contributor",
          fragileSegment: "Most fragile segment",
          fragileCaption:
            "Smallest contribution today, but still important for signal diversity.",
          portfolioNote: "Portfolio note",
          portfolioNoteCaption:
            "Segment mix is now wired into the same dashboard contract the board sees.",
          noData: "n/a",
        };

  const rows = segments.map((segment) => ({
    segment: (
      <div>
        <p className="font-medium text-[var(--color-ink)]">
          {translateFinanceSegment(segment.label, locale)}
        </p>
        <p className="mt-1 text-xs text-[var(--color-ink-soft)]">
          {segmentNotes[segment.label] ?? copy.portfolioShare}
        </p>
      </div>
    ),
    revenue: formatCompactCurrency(segment.value, locale),
    share: formatPercent(segment.share),
    play: segmentLevers[segment.label] ?? copy.nextLever,
    signal: (
      <StatusPill tone={segment === lead ? "emerald" : "navy"}>
        {translateFinanceStatus(
          segment === lead ? "lead" : segment === tail ? "watch" : "stable",
          locale,
        )}
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
            translateFinanceDataMode(snapshot.data_mode, locale) ??
            snapshot.data_mode,
          tone: statusTone(snapshot.data_mode),
        }}
        actions={
          <>
            <StatusPill tone="emerald">
              {translateFinanceSegment(lead?.label, locale) ?? copy.noData}{" "}
              {translateFinanceStatus("lead", locale)}
            </StatusPill>
            <StatusPill tone={latestProfit >= 0 ? "emerald" : "rose"}>
              {copy.profitPulse}{" "}
              {translateFinanceStatus(latestProfit >= 0 ? "up" : "down", locale)}
            </StatusPill>
          </>
        }
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {segments.map((segment) => (
          <div
            key={segment.label}
            className="overflow-hidden rounded-[1.55rem] border border-[var(--color-border)] bg-[var(--color-panel)] shadow-[0_14px_48px_rgba(13,26,44,0.08)]"
          >
            <div
              className="h-2"
              style={{
                background: `linear-gradient(90deg, ${segment.color}, rgba(255,255,255,0.18))`,
              }}
            />
            <div className="p-5">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.24em] text-[var(--color-ink-soft)]">
                    {translateFinanceSegment(segment.label, locale)}
                  </p>
                  <p className="mt-2 text-3xl font-semibold tracking-tight text-[var(--color-ink)]">
                    {formatCompactCurrency(segment.value, locale)}
                  </p>
                </div>
                <StatusPill
                  tone={
                    segment === lead
                      ? "emerald"
                      : segment === tail
                        ? "amber"
                        : "navy"
                  }
                >
                  {translateFinanceStatus(
                    segment === lead ? "lead" : segment === tail ? "watch" : "steady",
                    locale,
                  )}
                </StatusPill>
              </div>
              <p className="mt-4 text-sm leading-6 text-[var(--color-ink-soft)]">
                {segmentNotes[segment.label]}
              </p>
              <div className="mt-5 flex items-center justify-between text-sm">
                <span className="text-[var(--color-ink-soft)]">
                  {copy.portfolioShare}
                </span>
                <span className="font-semibold text-[var(--color-ink)]">
                  {formatPercent(segment.share)}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
        <SectionPanel
          kicker={copy.playbook}
          title={copy.marginPlays}
          description={copy.marginPlaysDescription}
        >
          <DataTable
            columns={[
              { key: "segment", label: copy.segment },
              { key: "revenue", label: copy.revenue, align: "right" },
              { key: "share", label: copy.share, align: "right" },
              { key: "play", label: copy.nextLever },
              { key: "signal", label: copy.signal, align: "right" },
            ]}
            rows={rows}
          />
        </SectionPanel>

        <SectionPanel
          kicker={copy.portfolioPosture}
          title={copy.managementKnow}
          description={copy.managementKnowDescription}
        >
          <div className="grid gap-4">
            <InsightCard
              label={copy.leadContributor}
              value={translateFinanceSegment(lead?.label, locale) ?? copy.noData}
              caption={
                locale === "vi"
                  ? `Hiện đang đóng góp ${formatPercent(lead?.share ?? 0)} doanh thu đã ghi nhận.`
                  : `Currently driving ${formatPercent(lead?.share ?? 0)} of booked revenue.`
              }
              tone="emerald"
            />
            <InsightCard
              label={copy.fragileSegment}
              value={translateFinanceSegment(tail?.label, locale) ?? copy.noData}
              caption={copy.fragileCaption}
              tone="amber"
            />
            <InsightCard
              label={copy.portfolioNote}
              value={snapshot.company_code}
              caption={copy.portfolioNoteCaption}
              tone="navy"
            />
          </div>
        </SectionPanel>
      </div>
    </section>
  );
}
