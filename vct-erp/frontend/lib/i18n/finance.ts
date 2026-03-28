import type { FinanceDashboardSnapshot } from "@/lib/contracts/finance";
import type { AppLocale } from "@/lib/i18n/shared";

type LocalizedText = Record<AppLocale, string>;

const financeLocaleMap: Record<AppLocale, string> = {
  vi: "vi-VN",
  en: "en-US",
};

const cardTitles: Record<string, LocalizedText> = {
  cash_assets: {
    vi: "Tổng tài sản hiện có",
    en: "Current Cash Assets",
  },
  quarter_net_revenue: {
    vi: "Doanh thu thuần quý",
    en: "Quarter Net Revenue",
  },
  runway_index: {
    vi: "Chỉ số runway",
    en: "Runway Index",
  },
};

const cardDescriptions: Record<string, LocalizedText> = {
  cash_assets: {
    vi: "Tổng hợp số dư tiền mặt và tiền gửi đang sẵn sàng cho điều hành.",
    en: "Combined cash and bank balances currently available for decisions.",
  },
  quarter_net_revenue: {
    vi: "Doanh thu đã được ghi nhận trong quý hiện tại sau các khoản giảm trừ.",
    en: "Revenue already recognized in the current quarter after deductions.",
  },
  runway_index: {
    vi: "Số tháng vận hành ước tính theo mức đốt tiền hiện tại.",
    en: "Estimated operating months remaining under the current burn rate.",
  },
};

const fallbackDescriptions: Record<AppLocale, string> = {
  vi: "Chờ kết nối đến backend tài chính",
  en: "Waiting for the finance backend connection",
};

const statusLabels: Record<string, LocalizedText> = {
  healthy: { vi: "ổn định", en: "healthy" },
  warning: { vi: "cảnh báo", en: "warning" },
  critical: { vi: "nguy cơ", en: "critical" },
  ready: { vi: "sẵn sàng", en: "ready" },
  live: { vi: "trực tiếp", en: "live" },
  watch: { vi: "theo dõi", en: "watch" },
  risk: { vi: "rủi ro", en: "risk" },
  queue: { vi: "hàng đợi", en: "queue" },
  balanced: { vi: "cân", en: "balanced" },
  scheduled: { vi: "đã lịch", en: "scheduled" },
  "on track": { vi: "đúng tiến độ", en: "on track" },
  review: { vi: "rà soát", en: "review" },
  recorded: { vi: "đã ghi", en: "recorded" },
  logged: { vi: "đã log", en: "logged" },
  stable: { vi: "ổn định", en: "stable" },
  steady: { vi: "ổn định", en: "steady" },
  lead: { vi: "dẫn đầu", en: "lead" },
  today: { vi: "hôm nay", en: "today" },
  up: { vi: "tăng", en: "up" },
  down: { vi: "giảm", en: "down" },
};

const periodLabels: Record<string, LocalizedText> = {
  "vs thang truoc": {
    vi: "so với tháng trước",
    en: "vs previous month",
  },
  "vs quy truoc": {
    vi: "so với quý trước",
    en: "vs previous quarter",
  },
  "vs previous month": {
    vi: "so với tháng trước",
    en: "vs previous month",
  },
  "vs previous quarter": {
    vi: "so với quý trước",
    en: "vs previous quarter",
  },
};

const dataModeLabels: Record<string, LocalizedText> = {
  live: { vi: "trực tiếp", en: "live" },
  mock: { vi: "mô phỏng", en: "mock" },
  fallback: { vi: "dự phòng", en: "fallback" },
};

const refreshModeLabels: Record<string, LocalizedText> = {
  websocket: { vi: "websocket", en: "websocket" },
  polling: { vi: "thăm dò định kỳ", en: "polling" },
};

const segmentLabels: Record<string, LocalizedText> = {
  SaaS: { vi: "SaaS", en: "SaaS" },
  Dojo: { vi: "Võ đường", en: "Dojo" },
  Retail: { vi: "Bán lẻ", en: "Retail" },
  Rental: { vi: "Cho thuê", en: "Rental" },
};

const seriesLabels: Record<string, LocalizedText> = {
  revenue: { vi: "Doanh thu", en: "Revenue" },
  expense: { vi: "Chi phí", en: "Expense" },
  profit: { vi: "Lợi nhuận", en: "Profit" },
};

function localizeValue(
  value: string | undefined,
  locale: AppLocale,
  dictionary: Record<string, LocalizedText>,
) {
  if (!value) {
    return value;
  }

  const key = value.trim();
  const localized = dictionary[key] ?? dictionary[key.toLowerCase()];
  return localized ? localized[locale] : value;
}

export function getFinanceLocaleCode(locale: AppLocale) {
  return financeLocaleMap[locale];
}

export function formatFinanceDate(
  value: Date | string,
  locale: AppLocale,
  options?: Intl.DateTimeFormatOptions,
) {
  return new Intl.DateTimeFormat(
    getFinanceLocaleCode(locale),
    options,
  ).format(new Date(value));
}

export function translateFinanceStatus(
  status: string | undefined,
  locale: AppLocale,
) {
  return localizeValue(status, locale, statusLabels);
}

export function translateFinanceTrendPeriod(
  period: string | undefined,
  locale: AppLocale,
) {
  return localizeValue(period, locale, periodLabels);
}

export function translateFinanceDataMode(
  mode: string | undefined,
  locale: AppLocale,
) {
  return localizeValue(mode, locale, dataModeLabels);
}

export function translateFinanceRefreshMode(
  mode: string | undefined,
  locale: AppLocale,
) {
  return localizeValue(mode, locale, refreshModeLabels);
}

export function translateFinanceSegment(
  segment: string | undefined,
  locale: AppLocale,
) {
  return localizeValue(segment, locale, segmentLabels);
}

export function translateFinanceSeriesLabel(
  key: string,
  label: string,
  locale: AppLocale,
) {
  return localizeValue(key, locale, seriesLabels) ?? label;
}

export function localizeFinanceDashboardSnapshot(
  snapshot: FinanceDashboardSnapshot,
  locale: AppLocale,
) {
  return {
    ...snapshot,
    cards: snapshot.cards.map((card) => ({
      ...card,
      title: cardTitles[card.key]?.[locale] ?? card.title,
      description:
        cardDescriptions[card.key]?.[locale] ??
        (card.description === "Cho ket noi den backend tai chinh"
          ? fallbackDescriptions[locale]
          : card.description),
    })),
  };
}
