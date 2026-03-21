# Command Center JSON Contract

This document defines the alpha dashboard payload for VCT Command Center.

## Mock Endpoint

`GET /api/v1/finance/dashboard/mock`

Headers:

- `X-App-Role: cfo|ceo|system_admin`

Returns:

```json
{
  "company_code": "VCT_GROUP",
  "generated_at": "2026-03-22T08:30:00Z",
  "data_mode": "mock",
  "recommended_refresh": "websocket",
  "cards": [
    {
      "key": "cash_assets",
      "title": "Tong tai san hien co",
      "value": 8420000000,
      "formatted_value": "8.42 ty VND",
      "unit": "VND",
      "trend": {
        "direction": "up",
        "percentage": 6.8,
        "delta": 536000000,
        "period": "vs thang truoc"
      },
      "chart_data": [
        { "label": "2026-03-16", "value": 7860000000 }
      ]
    }
  ],
  "revenue_mix": [
    { "label": "SaaS", "value": 1320000000, "color": "#0F766E" }
  ],
  "cashflow_chart": {
    "granularity": "month",
    "x_axis": ["2025-10", "2025-11", "2025-12", "2026-01", "2026-02", "2026-03"],
    "series": [
      {
        "key": "revenue",
        "label": "Revenue",
        "color": "#0F766E",
        "values": [1480000000, 1560000000, 1620000000, 1710000000, 1790000000, 1880000000]
      }
    ]
  },
  "runway_projection": [
    {
      "label": "2026-04",
      "opening_cash": 8420000000,
      "contracted_inflow": 640000000,
      "projected_burn": 950000000,
      "projected_ending": 8110000000
    }
  ]
}
```

## KPI Card Contract

Each dashboard card uses:

```json
{
  "key": "cash_assets",
  "title": "Tong tai san hien co",
  "value": 8420000000,
  "formatted_value": "8.42 ty VND",
  "unit": "VND",
  "status": "healthy",
  "description": "Tong 1111 va 1121 tai thoi diem hien tai",
  "trend": {
    "direction": "up",
    "percentage": 6.8,
    "delta": 536000000,
    "period": "vs thang truoc"
  },
  "chart_data": [
    { "label": "2026-03-16", "value": 7860000000 }
  ]
}
```

## Pie Chart Contract

```json
{ "label": "SaaS", "value": 1320000000, "color": "#0F766E" }
```

## Multi-Line Chart Contract

```json
{
  "granularity": "month",
  "x_axis": ["2025-10", "2025-11", "2025-12", "2026-01", "2026-02", "2026-03"],
  "series": [
    {
      "key": "revenue",
      "label": "Revenue",
      "color": "#0F766E",
      "values": [1480000000, 1560000000, 1620000000, 1710000000, 1790000000, 1880000000]
    },
    {
      "key": "expense",
      "label": "Expense",
      "color": "#D97706",
      "values": [890000000, 910000000, 930000000, 955000000, 980000000, 1010000000]
    },
    {
      "key": "profit",
      "label": "Profit",
      "color": "#2563EB",
      "values": [590000000, 650000000, 690000000, 755000000, 810000000, 870000000]
    }
  ]
}
```

## Planned Real Endpoints

- `GET /api/v1/finance/dashboard`
- `GET /api/v1/finance/summary`
- `GET /api/v1/finance/cash-runway`
- `GET /api/v1/finance/segments`

These remain the real-data endpoints. During alpha UI work, frontend can bind to `/api/v1/finance/dashboard/mock` first, then swap to `/api/v1/finance/dashboard` or the widget endpoints without changing the chart schema.

Realtime integration notes live in `finance-realtime-dashboard.md`.
