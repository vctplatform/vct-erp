export function formatCurrency(value: number) {
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency: "VND",
    maximumFractionDigits: 0,
  }).format(value);
}

export function formatCompactCurrency(value: number) {
  const absolute = Math.abs(value);

  if (absolute >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)} ty VND`;
  }

  if (absolute >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)} tr VND`;
  }

  return formatCurrency(value);
}

export function formatPercent(value: number) {
  return `${value >= 0 ? "+" : ""}${value.toFixed(1)}%`;
}

export function formatRunway(value: number) {
  return `${value.toFixed(1)} thang`;
}
