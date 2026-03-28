import type { AppLocale } from "@/lib/i18n/shared";

export type FinanceNavIcon =
  | "command"
  | "segments"
  | "ledger"
  | "reconciliation"
  | "reports"
  | "control";

export type FinanceNavItem = {
  href: string;
  label: string;
  caption: string;
  icon: FinanceNavIcon;
};

type FinanceNavItemCopy = {
  href: string;
  icon: FinanceNavIcon;
  label: Record<AppLocale, string>;
  caption: Record<AppLocale, string>;
};

const financePrimaryNavCopy: FinanceNavItemCopy[] = [
  {
    href: "/",
    icon: "command",
    label: {
      vi: "Trung tâm điều hành",
      en: "Command Center",
    },
    caption: {
      vi: "Tín hiệu điều hành, runway và trạng thái tiền mặt theo thời gian thực",
      en: "Board signals, runway, and realtime cash posture",
    },
  },
  {
    href: "/segments",
    icon: "segments",
    label: {
      vi: "Mảng kinh doanh",
      en: "Segments",
    },
    caption: {
      vi: "Hiệu suất của SaaS, võ đường, bán lẻ và cho thuê",
      en: "SaaS, dojo, retail, and rental performance mix",
    },
  },
  {
    href: "/ledger",
    icon: "ledger",
    label: {
      vi: "Vận hành sổ cái",
      en: "Ledger Ops",
    },
    caption: {
      vi: "Luồng chứng từ, rà soát bút toán và chất lượng ghi sổ",
      en: "Voucher lanes, journal review, and posting quality",
    },
  },
  {
    href: "/reconciliation",
    icon: "reconciliation",
    label: {
      vi: "Đối soát",
      en: "Reconciliation",
    },
    caption: {
      vi: "Đối chiếu ngân hàng, ngoại lệ chưa xử lý và checklist khóa sổ",
      en: "Bank matching, unresolved exceptions, and close checklists",
    },
  },
];

const financeSecondaryNavCopy: FinanceNavItemCopy[] = [
  {
    href: "/reports",
    icon: "reports",
    label: {
      vi: "Báo cáo",
      en: "Reports",
    },
    caption: {
      vi: "Bảng cân đối số phát sinh, P&L, hồ sơ thuế và lịch khóa sổ",
      en: "Trial balance, P&L, tax pack, and closing calendar",
    },
  },
  {
    href: "/control-room",
    icon: "control",
    label: {
      vi: "Phòng điều khiển",
      en: "Control Room",
    },
    caption: {
      vi: "Lan can chính sách, truy cập, cache và tín hiệu kiểm toán",
      en: "Policy guardrails, access, cache, and audit signals",
    },
  },
];

function localizeNav(
  items: FinanceNavItemCopy[],
  locale: AppLocale,
): FinanceNavItem[] {
  return items.map((item) => ({
    href: item.href,
    icon: item.icon,
    label: item.label[locale],
    caption: item.caption[locale],
  }));
}

export function getFinanceNavigation(locale: AppLocale) {
  const primary = localizeNav(financePrimaryNavCopy, locale);
  const secondary = localizeNav(financeSecondaryNavCopy, locale);

  return {
    primary,
    secondary,
    all: [...primary, ...secondary],
  };
}
