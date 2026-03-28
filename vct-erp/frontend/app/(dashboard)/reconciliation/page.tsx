import { DataTable } from "@/components/finance/data-table";
import { InsightCard } from "@/components/finance/insight-card";
import { ModuleHeader } from "@/components/finance/module-header";
import { SectionPanel } from "@/components/finance/section-panel";
import { StatusPill } from "@/components/finance/status-pill";
import { getFinanceDashboardSnapshot } from "@/lib/api/finance";
import { requireDashboardCard, statusTone } from "@/lib/finance/insights";
import { formatCompactCurrency } from "@/lib/formatters";
import {
  translateFinanceDataMode,
  translateFinanceStatus,
} from "@/lib/i18n/finance";
import { getServerLocale } from "@/lib/i18n/server";

export default async function ReconciliationPage() {
  const locale = await getServerLocale();
  const snapshot = await getFinanceDashboardSnapshot(locale);
  const cashCard = requireDashboardCard(snapshot, "cash_assets");
  const runwayCard = requireDashboardCard(snapshot, "runway_index");

  const bankAmount = cashCard.value * 0.65;
  const cashAmount = cashCard.value - bankAmount;
  const matchedAmount = cashCard.value * 0.968;
  const openAmount = cashCard.value - matchedAmount;

  const copy =
    locale === "vi"
      ? {
          kicker: "Khóa sổ ngân quỹ",
          title: "Đối soát",
          description:
            "Bề mặt UI cho khóa sổ tài chính: ghép dòng tiền, cô lập ngoại lệ và ngăn báo cáo điều hành chạy nhanh hơn sự thật ngân hàng.",
          matched: "đã khớp 96.8%",
          runway: "runway",
          matchedAmount: "Giá trị đã khớp",
          matchedAmountCaption:
            "Phần ngân quỹ đã được đối soát với góc nhìn khóa sổ của sổ cái.",
          openExceptions: "Ngoại lệ còn mở",
          openExceptionsCaption: "Phần còn lại vẫn cần điều tra thủ công.",
          bankLedger: "Sổ ngân hàng",
          bankLedgerCaption:
            "Chế độ xem mô phỏng của tài khoản 1121 trong tư thế ngân quỹ hiện tại.",
          cashDesk: "Quỹ tiền mặt",
          cashDeskCaption:
            "Chế độ xem mô phỏng của tài khoản 1111 vẫn cần kỷ luật chốt quỹ.",
          exceptionQueue: "Hàng đợi ngoại lệ",
          distortClose: "Các mục vẫn có thể làm lệch khóa sổ",
          distortCloseDescription:
            "Không phải mọi sai lệch đều nguy hiểm, nhưng mọi sai lệch đều cần lý do trước khi tài chính ký duyệt.",
          reference: "Tham chiếu",
          source: "Nguồn",
          amount: "Giá trị",
          action: "Hành động",
          priority: "Ưu tiên",
          closeChecklist: "Checklist khóa sổ",
          signoff: "Trình tự ký duyệt",
          signoffDescription: "Runbook tối thiểu cho kiểm soát viên ngân quỹ.",
          checklist: [
            "Xác nhận mọi dòng sao kê ngân hàng đều có source reference.",
            "Rà soát phiếu quỹ thủ công trước khi khóa ngày.",
            "Xử lý mọi trường hợp giải phóng tiền cọc thuê mà chưa map hư hại.",
            "Làm mới cache dashboard sau khi các điều chỉnh ngân quỹ đã hạ cánh.",
          ],
          exceptionRows: [
            [
              "BANK-1121-00048",
              "Internet banking",
              openAmount * 0.38,
              "Kiểm tra source_ref còn thiếu",
              "today",
              "amber",
            ],
            [
              "CASH-1111-00017",
              "Counter cash",
              openAmount * 0.22,
              "Xác thực phiếu chốt quỹ",
              "queue",
              "navy",
            ],
            [
              "BANK-1121-00063",
              "Giải phóng cọc thuê",
              openAmount * 0.4,
              "Map tất toán cọc sang bút toán 3388",
              "risk",
              "rose",
            ],
          ] as const,
        }
      : {
          kicker: "Treasury close",
          title: "Reconciliation",
          description:
            "A UI surface for finance close: match cash movement, isolate exceptions, and prevent board reporting from outrunning the bank truth.",
          matched: "96.8% matched",
          runway: "runway",
          matchedAmount: "Matched amount",
          matchedAmountCaption:
            "Treasury already reconciled against the ledger close view.",
          openExceptions: "Open exceptions",
          openExceptionsCaption:
            "Residual amount still requiring manual investigation.",
          bankLedger: "Bank ledger",
          bankLedgerCaption:
            "Modeled view of account 1121 within the current treasury posture.",
          cashDesk: "Cash desk",
          cashDeskCaption:
            "Modeled view of account 1111 that still needs cashier close discipline.",
          exceptionQueue: "Exception queue",
          distortClose: "Items that can still distort close",
          distortCloseDescription:
            "Not every mismatch is dangerous, but every mismatch needs a reason before finance signs off.",
          reference: "Reference",
          source: "Source",
          amount: "Amount",
          action: "Action",
          priority: "Priority",
          closeChecklist: "Close checklist",
          signoff: "Sign-off sequence",
          signoffDescription: "A minimal runbook for the treasury controller.",
          checklist: [
            "Confirm all bank statement lines have source references.",
            "Review manual cash vouchers before locking the day.",
            "Resolve any rental deposit release without damage mapping.",
            "Re-run dashboard cache after treasury corrections land.",
          ],
          exceptionRows: [
            [
              "BANK-1121-00048",
              "Internet banking",
              openAmount * 0.38,
              "Check missing journal source_ref",
              "today",
              "amber",
            ],
            [
              "CASH-1111-00017",
              "Counter cash",
              openAmount * 0.22,
              "Validate cashier close slip",
              "queue",
              "navy",
            ],
            [
              "BANK-1121-00063",
              "Rental deposit release",
              openAmount * 0.4,
              "Map deposit settlement to 3388 release",
              "risk",
              "rose",
            ],
          ] as const,
        };

  const exceptionRows = copy.exceptionRows.map(
    ([reference, source, amount, action, priority, tone]) => ({
      reference,
      source,
      amount: formatCompactCurrency(amount, locale),
      action,
      priority: (
        <StatusPill tone={tone}>
          {translateFinanceStatus(priority, locale)}
        </StatusPill>
      ),
    }),
  );

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
            <StatusPill tone="emerald">{copy.matched}</StatusPill>
            <StatusPill tone={statusTone(runwayCard.status)}>
              {copy.runway} {translateFinanceStatus(runwayCard.status, locale)}
            </StatusPill>
          </>
        }
      />

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <InsightCard
          label={copy.matchedAmount}
          value={formatCompactCurrency(matchedAmount, locale)}
          caption={copy.matchedAmountCaption}
          tone="emerald"
        />
        <InsightCard
          label={copy.openExceptions}
          value={formatCompactCurrency(openAmount, locale)}
          caption={copy.openExceptionsCaption}
          tone="amber"
        />
        <InsightCard
          label={copy.bankLedger}
          value={formatCompactCurrency(bankAmount, locale)}
          caption={copy.bankLedgerCaption}
          tone="navy"
        />
        <InsightCard
          label={copy.cashDesk}
          value={formatCompactCurrency(cashAmount, locale)}
          caption={copy.cashDeskCaption}
          tone="navy"
        />
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <SectionPanel
          kicker={copy.exceptionQueue}
          title={copy.distortClose}
          description={copy.distortCloseDescription}
        >
          <DataTable
            columns={[
              { key: "reference", label: copy.reference },
              { key: "source", label: copy.source },
              { key: "amount", label: copy.amount, align: "right" },
              { key: "action", label: copy.action },
              { key: "priority", label: copy.priority, align: "right" },
            ]}
            rows={exceptionRows}
          />
        </SectionPanel>

        <SectionPanel
          kicker={copy.closeChecklist}
          title={copy.signoff}
          description={copy.signoffDescription}
        >
          <div className="space-y-3">
            {copy.checklist.map((item) => (
              <div
                key={item}
                className="rounded-[1.2rem] border border-[var(--color-border)] bg-[var(--color-canvas-soft)] px-4 py-3 text-sm leading-6 text-[var(--color-ink-soft)]"
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
