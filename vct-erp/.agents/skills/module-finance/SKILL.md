---
name: module-finance
description: >-
  Mega-Skill for Finance & Accounting ERP module. Vietnamese Accounting Standards (VAS),
  double-entry bookkeeping, chart of accounts, journal entries, financial reporting,
  bank reconciliation, tax compliance, and SaaS revenue recognition.
metadata:
  author: VCT Platform
  version: "1.0.0"
  type: "Mega-Skill"
  locale: vi-VN
---

# MODULE-FINANCE — MEGA-SKILL

> Domain expertise cho module Tài chính & Kế toán. VAS-compliant, double-entry, enterprise-grade.

---

## 🔹 NĂNG LỰC: KẾ TOÁN TRƯỞNG (Chief Accountant)

> *"Mỗi đồng phải có nguồn gốc. Mỗi bút toán phải cân bằng."*

### Vietnamese Accounting Standards (VAS TT200)

#### Hệ thống Tài khoản (Chart of Accounts)
```
Loại 1: TÀI SẢN (asset, debit side)
├── 111 - Tiền mặt
│   └── 1111 - Tiền mặt Việt Nam
├── 112 - Tiền gửi ngân hàng
│   └── 1121 - TGNH Việt Nam đồng
├── 131 - Phải thu khách hàng
├── 133 - Thuế GTGT được khấu trừ
│   └── 1331 - Thuế GTGT hàng hóa dịch vụ
├── 141 - Tạm ứng
└── 242 - Chi phí trả trước

Loại 3: NỢ PHẢI TRẢ (liability, credit side)
├── 331 - Phải trả người bán
├── 333 - Thuế phải nộp Nhà nước
│   ├── 33311 - Thuế GTGT đầu ra
│   └── 3334 - Thuế TNDN
├── 334 - Phải trả người lao động
└── 338 - Phải trả, phải nộp khác
    ├── 3387 - Doanh thu chưa thực hiện (deferred revenue)
    └── 3388 - Phải trả, phải nộp khác (deposit holding)

Loại 4: VỐN CHỦ SỞ HỮU (equity, credit side)
├── 411 - Vốn đầu tư của chủ sở hữu
├── 421 - Lợi nhuận sau thuế chưa phân phối
└── 911 - Xác định kết quả kinh doanh

Loại 5: DOANH THU (revenue, credit side)
├── 511 - Doanh thu bán hàng và cung cấp dịch vụ
│   ├── 5111 - Doanh thu bán hàng hóa
│   └── 5113 - Doanh thu cung cấp dịch vụ (SaaS + Võ đường)
├── 515 - Doanh thu hoạt động tài chính
├── 521 - Các khoản giảm trừ doanh thu
│   └── 5211 - Chiết khấu thương mại
└── 711 - Thu nhập khác

Loại 6-8: CHI PHÍ (expense, debit side)
├── 632 - Giá vốn hàng bán
├── 635 - Chi phí tài chính
├── 641 - Chi phí bán hàng
├── 642 - Chi phí quản lý doanh nghiệp
│   ├── 6421 - Chi phí nhân viên quản lý
│   └── 6422 - Chi phí vật liệu quản lý
├── 811 - Chi phí khác
└── 821 - Chi phí thuế TNDN
    └── 8211 - CP thuế TNDN hiện hành
```

### Nguyên tắc Bút toán Kép (Double-Entry)
```
RULE #1: Mọi giao dịch → Nợ = Có (Debit = Credit)
RULE #2: Tài sản tăng → Nợ (Debit)
         Tài sản giảm → Có (Credit)
RULE #3: Nợ phải trả tăng → Có (Credit)
         Nợ phải trả giảm → Nợ (Debit)
RULE #4: Doanh thu → Có (Credit)
         Chi phí → Nợ (Debit)
```

### Bút toán Mẫu (VCT Business Scenarios)

#### 1. Thu tiền SaaS subscription (Khách trả trước 12 tháng)
```
Nợ 1121 (TGNH VND)           12.000.000
    Có 3387 (Doanh thu chưa TH)          12.000.000
→ Hàng tháng ghi nhận doanh thu:
Nợ 3387 (Doanh thu chưa TH)   1.000.000
    Có 5113 (DT cung cấp DV)              1.000.000
```

#### 2. Võ đường đóng học phí
```
Nợ 1111 (Tiền mặt)            500.000
    Có 5113 (DT cung cấp DV)              500.000
→ Nếu có thuế GTGT:
Nợ 1111 (Tiền mặt)            550.000
    Có 5113 (DT cung cấp DV)              500.000
    Có 33311 (Thuế GTGT đầu ra)            50.000
```

#### 3. Chi lương nhân viên
```
Nợ 6421 (CP nhân viên QL)     15.000.000
    Có 334 (Phải trả NLĐ)                15.000.000
→ Khi chuyển khoản:
Nợ 334 (Phải trả NLĐ)         15.000.000
    Có 1121 (TGNH VND)                    15.000.000
```

#### 4. Nhận tiền đặt cọc thuê thiết bị
```
Nợ 1121 (TGNH VND)            2.000.000
    Có 3388 (Deposit holding)              2.000.000
→ Trả lại đặt cọc:
Nợ 3388 (Deposit holding)      2.000.000
    Có 1121 (TGNH VND)                     2.000.000
```

---

## 🔹 NĂNG LỰC: FINANCIAL REPORTING

### Báo cáo Tài chính (VAS Format)

#### B01-DN: Bảng Cân đối Kế toán (Balance Sheet)
```
TÀI SẢN = NỢ PHẢI TRẢ + VỐN CHỦ SỞ HỮU

I. Tài sản ngắn hạn
   1. Tiền (111 + 112)
   2. Phải thu (131)
   3. Tạm ứng (141)
   4. Chi phí trả trước (242)

II. Tài sản dài hạn
   (Tương lai: TSCĐ, BĐSĐT...)

III. Nợ phải trả
   1. Phải trả người bán (331)
   2. Thuế phải nộp (333)
   3. Phải trả NLĐ (334)
   4. Doanh thu chưa thực hiện (3387)

IV. Vốn chủ sở hữu
   1. Vốn đầu tư (411)
   2. LNST chưa phân phối (421)
```

#### B02-DN: Báo cáo Kết quả Kinh doanh (P&L)
```
1. Doanh thu bán hàng và CCDV (511)
2. Các khoản giảm trừ DT (521)
3. Doanh thu thuần = (1) - (2)
4. Giá vốn hàng bán (632)
5. Lợi nhuận gộp = (3) - (4)
6. Doanh thu hoạt động TC (515)
7. Chi phí tài chính (635)
8. Chi phí bán hàng (641)
9. Chi phí QLDN (642)
10. LNTT từ HĐKD = (5) + (6) - (7) - (8) - (9)
11. Thu nhập khác (711)
12. Chi phí khác (811)
13. LNTT = (10) + (11) - (12)
14. Chi phí thuế TNDN (821)
15. LNST = (13) - (14)
```

#### Sổ Cái (General Journal)
```sql
-- SQL pattern for general journal
SELECT
  je.posting_date,
  je.entry_no,
  je.description,
  a.code AS account_code,
  a.name AS account_name,
  CASE WHEN ji.side = 'debit' THEN ji.amount END AS debit,
  CASE WHEN ji.side = 'credit' THEN ji.amount END AS credit
FROM journal_entries je
JOIN journal_items ji ON ji.journal_entry_id = je.id
JOIN accounts a ON a.id = ji.account_id
WHERE je.status = 'posted'
ORDER BY je.posting_date, je.entry_no, ji.line_no
```

---

## 🔹 NĂNG LỰC: BANK RECONCILIATION

### Quy trình Đối chiếu Ngân hàng
```
1. Import bank statement (CSV/API)
   → bank_statement_lines table

2. Auto-match rules:
   ├── Exact amount + reference → matched
   ├── Amount match + date ±3 days → suggested
   └── No match → manual review

3. Status flow: open → matched | manual | ignored

4. Reconciliation report:
   ├── Book balance (from account_balances)
   ├── Bank balance (from bank_statement_lines)
   ├── Outstanding deposits (book not bank)
   ├── Outstanding payments (bank not book)
   └── Adjusted balance (should match)
```

---

## 🔹 NĂNG LỰC: SAAS REVENUE RECOGNITION

### Deferred Revenue Flow
```
1. Contract created → saas_contracts table
2. Initial journal: Nợ Cash / Có Deferred Revenue
3. Monthly schedule → saas_revenue_schedules table
4. Monthly recognition: Nợ Deferred / Có Service Revenue
5. Status: scheduled → recognized

Automation:
├── Cron job checks due schedules
├── Creates journal entries automatically
├── Marks schedule as recognized
└── Publishes event via outbox
```

---

## Development Checklist (Finance Module)

- [ ] Double-entry balance enforced (DB constraint)
- [ ] Voucher numbering VAS-compliant
- [ ] Chart of accounts follows TT200
- [ ] Financial reports match VAS templates
- [ ] Bank reconciliation auto-matching
- [ ] SaaS revenue recognition automated
- [ ] Void/reverse instead of delete
- [ ] Audit trail for all operations
- [ ] Multi-currency support (future)
- [ ] Tax calculation accurate

## Trigger Patterns

- "tài chính", "kế toán", "bút toán", "journal entry"
- "doanh thu", "chi phí", "báo cáo tài chính"
- "đối chiếu ngân hàng", "reconciliation"
- "thuế", "VAT", "GTGT", "TNDN"
- "chart of accounts", "hệ thống tài khoản"
- "VAS", "TT200", "chuẩn mực kế toán"
