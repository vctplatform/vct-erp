---
name: erp-finance-accounting
description: Tài chính & Kế toán - Báo cáo P/L, Bảng cân đối kế toán, Lưu chuyển tiền tệ, Thuế VAT/TNDN/TNCN, Công nợ, Tỷ số tài chính, Ngân sách vs Thực tế, Sổ cái & Nhật ký cho VCT Platform ERP.
---

# Tài chính & Kế toán - Báo cáo Doanh nghiệp

## Tổng quan
Module báo cáo Tài chính & Kế toán cung cấp toàn bộ hệ thống báo cáo theo chuẩn VAS (Vietnamese Accounting Standards) và IFRS, phục vụ quản trị nội bộ, nộp thuế, và ra quyết định tài chính.

## Danh mục Báo cáo Tài chính

### 1. Báo cáo Kết quả Hoạt động Kinh doanh (P/L - Income Statement)

#### Cấu trúc báo cáo theo VAS
```
┌─────────────────────────────────────────────────────────────┐
│   BÁO CÁO KẾT QUẢ HOẠT ĐỘNG KINH DOANH                   │
│   Kỳ báo cáo: [Tháng/Quý/Năm]                             │
├─────────────────────────────────────────────────────────────┤
│ 1. Doanh thu bán hàng và cung cấp dịch vụ        xxx       │
│ 2. Các khoản giảm trừ doanh thu                  (xxx)     │
│ 3. Doanh thu thuần (= 1 - 2)                     xxx       │
│ 4. Giá vốn hàng bán                              (xxx)     │
│ 5. Lợi nhuận gộp (= 3 - 4)                       xxx       │
│ 6. Doanh thu hoạt động tài chính                  xxx       │
│ 7. Chi phí tài chính                             (xxx)     │
│    - Trong đó: Chi phí lãi vay                   (xxx)     │
│ 8. Chi phí bán hàng                              (xxx)     │
│ 9. Chi phí quản lý doanh nghiệp                  (xxx)     │
│10. Lợi nhuận thuần từ HĐKD (= 5+6-7-8-9)         xxx       │
│11. Thu nhập khác                                   xxx       │
│12. Chi phí khác                                  (xxx)     │
│13. Lợi nhuận khác (= 11 - 12)                     xxx       │
│14. Tổng lợi nhuận trước thuế (= 10 + 13)          xxx       │
│15. Chi phí thuế TNDN hiện hành                   (xxx)     │
│16. Chi phí thuế TNDN hoãn lại                    (xxx)     │
│17. Lợi nhuận sau thuế TNDN (= 14 - 15 - 16)      xxx       │
└─────────────────────────────────────────────────────────────┘
```

#### SQL Template - P/L Report
```sql
-- Báo cáo P/L theo kỳ
WITH revenue AS (
    SELECT
        COALESCE(SUM(CASE WHEN account_code LIKE '511%' THEN amount END), 0) AS gross_revenue,
        COALESCE(SUM(CASE WHEN account_code LIKE '521%' THEN amount END), 0) AS revenue_deductions,
        COALESCE(SUM(CASE WHEN account_code LIKE '515%' THEN amount END), 0) AS financial_revenue,
        COALESCE(SUM(CASE WHEN account_code LIKE '711%' THEN amount END), 0) AS other_income
    FROM journal_entries je
    JOIN accounts a ON a.id = je.account_id
    WHERE je.posting_date BETWEEN $1 AND $2
    AND je.status = 'posted'
    AND je.org_id = $3
),
expenses AS (
    SELECT
        COALESCE(SUM(CASE WHEN account_code LIKE '632%' THEN amount END), 0) AS cogs,
        COALESCE(SUM(CASE WHEN account_code LIKE '635%' THEN amount END), 0) AS financial_expense,
        COALESCE(SUM(CASE WHEN account_code LIKE '641%' THEN amount END), 0) AS selling_expense,
        COALESCE(SUM(CASE WHEN account_code LIKE '642%' THEN amount END), 0) AS admin_expense,
        COALESCE(SUM(CASE WHEN account_code LIKE '811%' THEN amount END), 0) AS other_expense,
        COALESCE(SUM(CASE WHEN account_code LIKE '821%' THEN amount END), 0) AS cit_expense
    FROM journal_entries je
    JOIN accounts a ON a.id = je.account_id
    WHERE je.posting_date BETWEEN $1 AND $2
    AND je.status = 'posted'
    AND je.org_id = $3
)
SELECT
    r.gross_revenue,
    r.revenue_deductions,
    r.gross_revenue - r.revenue_deductions AS net_revenue,
    e.cogs,
    (r.gross_revenue - r.revenue_deductions - e.cogs) AS gross_profit,
    r.financial_revenue,
    e.financial_expense,
    e.selling_expense,
    e.admin_expense,
    (r.gross_revenue - r.revenue_deductions - e.cogs + r.financial_revenue 
     - e.financial_expense - e.selling_expense - e.admin_expense) AS operating_profit,
    r.other_income,
    e.other_expense,
    (r.gross_revenue - r.revenue_deductions - e.cogs + r.financial_revenue 
     - e.financial_expense - e.selling_expense - e.admin_expense 
     + r.other_income - e.other_expense) AS profit_before_tax,
    e.cit_expense,
    (r.gross_revenue - r.revenue_deductions - e.cogs + r.financial_revenue 
     - e.financial_expense - e.selling_expense - e.admin_expense 
     + r.other_income - e.other_expense - e.cit_expense) AS net_profit
FROM revenue r, expenses e;
```

### 2. Bảng Cân đối Kế toán (Balance Sheet)

#### Cấu trúc theo VAS
```
┌─────────────────────────────────────────────────────────────┐
│   BẢNG CÂN ĐỐI KẾ TOÁN                                    │
│   Tại ngày: [dd/mm/yyyy]                                    │
├─────────────────────────────────────────────────────────────┤
│ A. TÀI SẢN                                                  │
│ I. Tài sản ngắn hạn                                         │
│    1. Tiền và tương đương tiền              (111-113)        │
│    2. Đầu tư tài chính ngắn hạn            (121)            │
│    3. Phải thu ngắn hạn                     (131-139)        │
│    4. Hàng tồn kho                          (151-157)        │
│    5. Tài sản ngắn hạn khác                (141,242)        │
│ II. Tài sản dài hạn                                         │
│    1. Phải thu dài hạn                      (131*)           │
│    2. Tài sản cố định                       (211-214)        │
│    3. Bất động sản đầu tư                   (217)            │
│    4. Tài sản dở dang dài hạn              (241)            │
│    5. Đầu tư tài chính dài hạn             (221-228)        │
│    6. Tài sản dài hạn khác                 (242*)           │
│                                                              │
│ B. NGUỒN VỐN                                                │
│ I. Nợ phải trả                                               │
│    1. Nợ ngắn hạn                                            │
│       - Phải trả người bán                  (331)            │
│       - Thuế và các khoản phải nộp NN       (333)            │
│       - Phải trả người lao động             (334)            │
│       - Vay và nợ thuê tài chính ngắn hạn  (341)            │
│       - Dự phòng phải trả ngắn hạn         (352)            │
│    2. Nợ dài hạn                                             │
│       - Vay và nợ thuê tài chính dài hạn   (341*)           │
│       - Quỹ dự phòng trợ cấp              (351)            │
│ II. Vốn chủ sở hữu                                          │
│    1. Vốn đầu tư của CSH                    (411)            │
│    2. Thặng dư vốn cổ phần                  (412)            │
│    3. Lợi nhuận sau thuế chưa phân phối     (421)            │
│    4. Quỹ đầu tư phát triển                 (414)            │
└─────────────────────────────────────────────────────────────┘
```

#### SQL Template - Balance Sheet
```sql
-- Bảng cân đối kế toán tại thời điểm
SELECT
    a.account_code,
    a.account_name,
    a.account_type,
    CASE 
        WHEN a.account_code LIKE '1%' THEN 'ASSET'
        WHEN a.account_code LIKE '2%' THEN 'ASSET'
        WHEN a.account_code LIKE '3%' THEN 'LIABILITY'
        WHEN a.account_code LIKE '4%' THEN 'EQUITY'
    END AS bs_category,
    CASE
        WHEN a.account_code LIKE '1%' THEN 'SHORT_TERM_ASSET'
        WHEN a.account_code LIKE '2%' THEN 'LONG_TERM_ASSET'
        WHEN a.account_code BETWEEN '311' AND '339' THEN 'SHORT_TERM_LIABILITY'
        WHEN a.account_code BETWEEN '341' AND '359' THEN 'LONG_TERM_LIABILITY'
        WHEN a.account_code LIKE '4%' THEN 'EQUITY'
    END AS bs_sub_category,
    SUM(je.debit_amount - je.credit_amount) AS balance
FROM accounts a
LEFT JOIN journal_entries je ON je.account_id = a.id
WHERE je.posting_date <= $1
AND je.status = 'posted'
AND je.org_id = $2
GROUP BY a.account_code, a.account_name, a.account_type
HAVING SUM(je.debit_amount - je.credit_amount) != 0
ORDER BY a.account_code;
```

### 3. Báo cáo Lưu chuyển Tiền tệ (Cash Flow Statement)

#### Phương pháp Trực tiếp
```
┌─────────────────────────────────────────────────────────────┐
│ I. LƯU CHUYỂN TIỀN TỪ HOẠT ĐỘNG KINH DOANH                │
│    1. Tiền thu từ bán hàng, CCDV                   xxx      │
│    2. Tiền chi trả cho NCC hàng hóa, DV           (xxx)    │
│    3. Tiền chi trả cho NLĐ                        (xxx)    │
│    4. Tiền lãi vay đã trả                         (xxx)    │
│    5. Thuế TNDN đã nộp                            (xxx)    │
│    6. Tiền thu/chi khác từ HĐKD                    xxx      │
│    → Lưu chuyển tiền thuần từ HĐKD                 xxx      │
│                                                              │
│ II. LƯU CHUYỂN TIỀN TỪ HOẠT ĐỘNG ĐẦU TƯ                   │
│    1. Tiền chi mua sắm TSCĐ, TSDT                (xxx)    │
│    2. Tiền thu từ thanh lý TSCĐ                    xxx      │
│    3. Tiền chi cho vay, mua CCTC                  (xxx)    │
│    4. Tiền thu hồi cho vay, bán CCTC               xxx      │
│    5. Cổ tức, lợi nhuận được chia                   xxx      │
│    → Lưu chuyển tiền thuần từ HĐĐT                 xxx      │
│                                                              │
│ III. LƯU CHUYỂN TIỀN TỪ HOẠT ĐỘNG TÀI CHÍNH               │
│    1. Tiền thu từ phát hành CP, nhận vốn góp        xxx      │
│    2. Tiền chi trả vốn góp cho CSH                (xxx)    │
│    3. Tiền vay nhận được                            xxx      │
│    4. Tiền trả nợ gốc vay                        (xxx)    │
│    5. Cổ tức, lợi nhuận đã trả cho CSH           (xxx)    │
│    → Lưu chuyển tiền thuần từ HĐTC                 xxx      │
│                                                              │
│ IV. TĂNG/GIẢM TIỀN THUẦN TRONG KỲ (I+II+III)      xxx      │
│     Tiền đầu kỳ                                     xxx      │
│     Tiền cuối kỳ                                     xxx      │
└─────────────────────────────────────────────────────────────┘
```

### 4. Sổ Cái & Sổ Chi tiết

#### Sổ Cái Tổng hợp (General Ledger)
```sql
-- Sổ cái tài khoản theo kỳ
SELECT
    je.posting_date,
    je.entry_number,
    je.description,
    je.debit_amount,
    je.credit_amount,
    SUM(je.debit_amount - je.credit_amount) OVER (
        ORDER BY je.posting_date, je.entry_number
    ) AS running_balance
FROM journal_entries je
WHERE je.account_id = $1
AND je.posting_date BETWEEN $2 AND $3
AND je.status = 'posted'
ORDER BY je.posting_date, je.entry_number;
```

#### Sổ Nhật ký Chung
| Ngày | Số CT | Diễn giải | TK Nợ | TK Có | Số tiền |
|------|-------|-----------|--------|-------|---------|
| dd/mm | JE001 | Mua hàng | 156 | 331 | xxx |
| dd/mm | JE002 | Thu tiền | 111 | 131 | xxx |

### 5. Báo cáo Thuế

#### Tờ khai Thuế GTGT (VAT) - Mẫu 01/GTGT
| Chỉ tiêu | Nội dung | Giá trị hàng hóa | Thuế GTGT |
|-----------|----------|-------------------|-----------|
| [21] | Hàng hóa, DV bán ra chịu thuế suất 0% | xxx | 0 |
| [22] | Hàng hóa, DV bán ra chịu thuế suất 5% | xxx | xxx |
| [23] | Hàng hóa, DV bán ra chịu thuế suất 8% | xxx | xxx |
| [24] | Hàng hóa, DV bán ra chịu thuế suất 10% | xxx | xxx |
| [25] | **Tổng thuế GTGT đầu ra** | | **xxx** |
| [26] | Thuế GTGT đầu vào được khấu trừ | | xxx |
| [40] | **Thuế GTGT phải nộp (= 25 - 26)** | | **xxx** |

#### Thuế TNDN - Tạm tính quý
```sql
-- Tạm tính thuế TNDN theo quý
SELECT
    EXTRACT(QUARTER FROM je.posting_date) AS quarter,
    EXTRACT(YEAR FROM je.posting_date) AS year,
    SUM(CASE WHEN a.account_type = 'REVENUE' THEN je.credit_amount - je.debit_amount ELSE 0 END) AS total_revenue,
    SUM(CASE WHEN a.account_type = 'EXPENSE' THEN je.debit_amount - je.credit_amount ELSE 0 END) AS total_expense,
    SUM(CASE WHEN a.account_type = 'REVENUE' THEN je.credit_amount - je.debit_amount ELSE 0 END) -
    SUM(CASE WHEN a.account_type = 'EXPENSE' THEN je.debit_amount - je.credit_amount ELSE 0 END) AS taxable_income,
    (SUM(CASE WHEN a.account_type = 'REVENUE' THEN je.credit_amount - je.debit_amount ELSE 0 END) -
     SUM(CASE WHEN a.account_type = 'EXPENSE' THEN je.debit_amount - je.credit_amount ELSE 0 END)) * 0.20 AS cit_amount
FROM journal_entries je
JOIN accounts a ON a.id = je.account_id
WHERE je.posting_date BETWEEN $1 AND $2
AND je.status = 'posted'
AND je.org_id = $3
GROUP BY EXTRACT(QUARTER FROM je.posting_date), EXTRACT(YEAR FROM je.posting_date);
```

### 6. Báo cáo Công nợ

#### Bảng Tuổi nợ Phải thu (Aging Report)
```sql
-- Aging report - Phải thu khách hàng
SELECT
    c.name AS customer_name,
    c.tax_code,
    SUM(CASE WHEN age_days BETWEEN 0 AND 30 THEN balance ELSE 0 END) AS "0-30 ngày",
    SUM(CASE WHEN age_days BETWEEN 31 AND 60 THEN balance ELSE 0 END) AS "31-60 ngày",
    SUM(CASE WHEN age_days BETWEEN 61 AND 90 THEN balance ELSE 0 END) AS "61-90 ngày",
    SUM(CASE WHEN age_days BETWEEN 91 AND 180 THEN balance ELSE 0 END) AS "91-180 ngày",
    SUM(CASE WHEN age_days > 180 THEN balance ELSE 0 END) AS ">180 ngày",
    SUM(balance) AS total_outstanding
FROM (
    SELECT
        inv.customer_id,
        inv.invoice_number,
        inv.total_amount - COALESCE(SUM(pay.amount), 0) AS balance,
        CURRENT_DATE - inv.due_date AS age_days
    FROM invoices inv
    LEFT JOIN payments pay ON pay.invoice_id = inv.id
    WHERE inv.type = 'receivable'
    AND inv.status != 'paid'
    AND inv.org_id = $1
    GROUP BY inv.id
    HAVING inv.total_amount - COALESCE(SUM(pay.amount), 0) > 0
) aged
JOIN customers c ON c.id = aged.customer_id
GROUP BY c.name, c.tax_code
ORDER BY total_outstanding DESC;
```

#### Đối chiếu Công nợ
| Nội dung | Số liệu DN | Số liệu Đối tác | Chênh lệch |
|----------|-----------|-----------------|------------|
| Dư đầu kỳ | xxx | xxx | xxx |
| Phát sinh Nợ | xxx | xxx | xxx |
| Phát sinh Có | xxx | xxx | xxx |
| Dư cuối kỳ | xxx | xxx | xxx |

### 7. Phân tích Tỷ số Tài chính

#### Nhóm tỷ số Thanh khoản
| Tỷ số | Công thức | Ý nghĩa | Benchmark |
|-------|----------|---------|-----------|
| Current Ratio | TSNH / Nợ NH | Khả năng thanh toán NH | > 1.5 |
| Quick Ratio | (TSNH - HTK) / Nợ NH | Thanh toán nhanh | > 1.0 |
| Cash Ratio | Tiền / Nợ NH | Thanh toán tức thời | > 0.5 |

#### Nhóm tỷ số Hiệu quả
| Tỷ số | Công thức | Ý nghĩa | Benchmark |
|-------|----------|---------|-----------|
| ROA | Lợi nhuận ròng / Tổng TS | Hiệu quả sử dụng TS | > 5% |
| ROE | Lợi nhuận ròng / VCSH | Sinh lời vốn CSH | > 15% |
| ROS | Lợi nhuận ròng / Doanh thu thuần | Biên lợi nhuận ròng | > 10% |
| Gross Margin | Lợi nhuận gộp / Doanh thu thuần | Biên lợi nhuận gộp | > 30% |
| EBITDA Margin | EBITDA / Doanh thu thuần | HC trước lãi vay, thuế, KH | > 20% |

#### Nhóm tỷ số Đòn bẩy
| Tỷ số | Công thức | Ý nghĩa | Benchmark |
|-------|----------|---------|-----------|
| D/E Ratio | Tổng Nợ / VCSH | Cấu trúc vốn | < 2.0 |
| D/A Ratio | Tổng Nợ / Tổng TS | Tỷ lệ nợ | < 60% |
| ICR | EBIT / Chi phí lãi vay | Khả năng trả lãi | > 3.0 |

#### Nhóm tỷ số Hoạt động
| Tỷ số | Công thức | Ý nghĩa | Benchmark |
|-------|----------|---------|-----------|
| Inventory Turnover | GVHB / HTK bình quân | Vòng quay HTK | > 6 |
| DSO | Phải thu BQ × 365 / DT thuần | Kỳ thu tiền BQ | < 45 ngày |
| DPO | Phải trả BQ × 365 / GVHB | Kỳ thanh toán BQ | 30-60 ngày |
| Asset Turnover | DT thuần / Tổng TS BQ | Hiệu suất TS | > 1.0 |

```sql
-- Tính tỷ số tài chính tự động
WITH bs AS (
    SELECT
        SUM(CASE WHEN account_code LIKE '1%' THEN balance END) AS current_assets,
        SUM(CASE WHEN account_code LIKE '15%' THEN balance END) AS inventory,
        SUM(CASE WHEN account_code IN ('111','112','113') THEN balance END) AS cash,
        SUM(CASE WHEN account_code LIKE '1%' OR account_code LIKE '2%' THEN balance END) AS total_assets,
        SUM(CASE WHEN account_code LIKE '3%' AND account_code < '340' THEN balance END) AS current_liabilities,
        SUM(CASE WHEN account_code LIKE '3%' THEN balance END) AS total_liabilities,
        SUM(CASE WHEN account_code LIKE '4%' THEN balance END) AS equity
    FROM account_balances
    WHERE period_end = $1 AND org_id = $2
),
pl AS (
    SELECT net_revenue, gross_profit, operating_profit, net_profit, financial_expense
    FROM pl_report WHERE period = $3 AND org_id = $2
)
SELECT
    bs.current_assets / NULLIF(bs.current_liabilities, 0) AS current_ratio,
    (bs.current_assets - bs.inventory) / NULLIF(bs.current_liabilities, 0) AS quick_ratio,
    pl.net_profit / NULLIF(bs.total_assets, 0) AS roa,
    pl.net_profit / NULLIF(bs.equity, 0) AS roe,
    pl.net_profit / NULLIF(pl.net_revenue, 0) AS ros,
    pl.gross_profit / NULLIF(pl.net_revenue, 0) AS gross_margin,
    bs.total_liabilities / NULLIF(bs.equity, 0) AS de_ratio
FROM bs, pl;
```

### 8. Ngân sách vs Thực tế (Budget vs Actual)

#### Báo cáo Variance Analysis
| Khoản mục | Ngân sách | Thực tế | Chênh lệch | % Chênh lệch | Đánh giá |
|-----------|----------|---------|-----------|-------------|---------|
| Doanh thu | xxx | xxx | xxx | x% | ▲/▼ |
| COGS | xxx | xxx | xxx | x% | ▲/▼ |
| Chi phí bán hàng | xxx | xxx | xxx | x% | ▲/▼ |
| Chi phí QLDN | xxx | xxx | xxx | x% | ▲/▼ |
| Lợi nhuận | xxx | xxx | xxx | x% | ▲/▼ |

```sql
-- Budget vs Actual comparison
SELECT
    b.category,
    b.account_code,
    a.account_name,
    b.budget_amount,
    COALESCE(act.actual_amount, 0) AS actual_amount,
    COALESCE(act.actual_amount, 0) - b.budget_amount AS variance,
    ROUND(((COALESCE(act.actual_amount, 0) - b.budget_amount) / NULLIF(b.budget_amount, 0)) * 100, 2) AS variance_pct,
    CASE 
        WHEN b.category = 'REVENUE' AND act.actual_amount >= b.budget_amount THEN 'FAVORABLE'
        WHEN b.category = 'EXPENSE' AND act.actual_amount <= b.budget_amount THEN 'FAVORABLE'
        ELSE 'UNFAVORABLE'
    END AS evaluation
FROM budgets b
JOIN accounts a ON a.account_code = b.account_code
LEFT JOIN (
    SELECT account_code, SUM(amount) AS actual_amount
    FROM journal_entries
    WHERE posting_date BETWEEN $1 AND $2
    AND status = 'posted' AND org_id = $3
    GROUP BY account_code
) act ON act.account_code = b.account_code
WHERE b.period = $4 AND b.org_id = $3
ORDER BY b.category, b.account_code;
```

### 9. Báo cáo Tài sản Cố định & Khấu hao

| STT | Tên TSCĐ | Mã TS | Ngày mua | Nguyên giá | KH lũy kế | GTCL | PP khấu hao | Thời gian KH |
|-----|---------|-------|---------|-----------|----------|------|------------|-------------|
| 1 | [Tên] | TS001 | dd/mm/yy | xxx | xxx | xxx | Đường thẳng | 5 năm |

### 10. Phân tích Dòng tiền (Cash Flow Analysis)

#### Dự báo Dòng tiền Tuần/Tháng
```sql
-- Cash flow forecast
SELECT
    week_start,
    opening_balance,
    SUM(expected_inflow) AS inflow,
    SUM(expected_outflow) AS outflow,
    opening_balance + SUM(expected_inflow) - SUM(expected_outflow) AS closing_balance
FROM cash_flow_forecast
WHERE org_id = $1
AND week_start BETWEEN $2 AND $3
GROUP BY week_start, opening_balance
ORDER BY week_start;
```

## Tần suất Báo cáo

| Báo cáo | Ngày | Tuần | Tháng | Quý | Năm |
|---------|-----|------|-------|-----|-----|
| Sổ quỹ tiền mặt | ✅ | | | | |
| Bảng kê ngân hàng | ✅ | | | | |
| Dòng tiền dự báo | | ✅ | | | |
| P/L tóm tắt | | | ✅ | | |
| Công nợ phải thu/trả | | | ✅ | | |
| Báo cáo VAT | | | ✅ | ✅ | |
| P/L chi tiết | | | | ✅ | ✅ |
| Bảng CĐKT | | | | ✅ | ✅ |
| Lưu chuyển tiền tệ | | | | ✅ | ✅ |
| Tỷ số tài chính | | | | ✅ | ✅ |
| Budget vs Actual | | | ✅ | ✅ | ✅ |
| Thuế TNDN | | | | ✅ | ✅ |
| Báo cáo tài chính kiểm toán | | | | | ✅ |

## Quyền hạn Xem Báo cáo

| Báo cáo | CEO | CFO | Kế toán trưởng | Kế toán viên |
|---------|-----|-----|---------------|-------------|
| P/L | ✅ | ✅ | ✅ | Phần hành |
| Bảng CĐKT | ✅ | ✅ | ✅ | Không |
| Cash Flow | ✅ | ✅ | ✅ | Không |
| Công nợ | ✅ | ✅ | ✅ | Phần hành |
| Tỷ số TC | ✅ | ✅ | ✅ | Không |
| Sổ cái/chi tiết | ✅ | ✅ | ✅ | Phần hành |
| Thuế | ✅ | ✅ | ✅ | Phần hành |
| Budget vs Actual | ✅ | ✅ | ✅ | Không |
