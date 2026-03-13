---
name: erp-admin-legal
description: Hành chính & Pháp chế - Quản lý tài sản cố định, hợp đồng, giấy phép, tuân thủ pháp luật, chi phí văn phòng, bảo hiểm cho VCT Platform ERP.
---

# Hành chính & Pháp chế - Báo cáo Doanh nghiệp

## Tổng quan
Báo cáo Hành chính & Pháp chế bao phủ quản lý tài sản, hợp đồng, tuân thủ pháp luật, chi phí văn phòng và quản trị rủi ro pháp lý.

## Danh mục Báo cáo

### 1. Quản lý Tài sản Cố định (Fixed Assets)

#### Bảng Tài sản & Khấu hao
| Mã TS | Tên TSCĐ | Ngày mua | Nguyên giá | KH lũy kế | GTCL | PP khấu hao | Bộ phận |
|-------|---------|---------|-----------|----------|------|------------|--------|
| TS001 | [Tên] | dd/mm/yy | xxx | xxx | xxx | Đường thẳng | [BP] |

```sql
-- Fixed asset depreciation schedule
SELECT fa.asset_code, fa.name, fa.acquisition_date, fa.original_cost,
    fa.useful_life_months, fa.depreciation_method,
    SUM(fd.amount) AS accumulated_depreciation,
    fa.original_cost - SUM(fd.amount) AS net_book_value,
    d.name AS department, fa.status
FROM fixed_assets fa
LEFT JOIN fixed_asset_depreciation fd ON fd.asset_id = fa.id
JOIN departments d ON d.id = fa.department_id
WHERE fa.org_id = $1
GROUP BY fa.id, d.name ORDER BY fa.asset_code;
```

#### KPIs Tài sản
| KPI | Công thức | Ý nghĩa |
|-----|----------|---------|
| Tỷ lệ hao mòn | KH lũy kế / Nguyên giá | Mức độ sử dụng |
| Asset Utilization | Doanh thu / Tổng TSCĐ | Hiệu suất TS |
| Maintenance Cost Ratio | Chi phí bảo trì / Nguyên giá | < 5%/năm |

### 2. Quản lý Hợp đồng (Contract Management)

#### Dashboard Hợp đồng
| Trạng thái | Số lượng | Giá trị |
|-----------|---------|---------|
| Đang hiệu lực | xxx | xxx VNĐ |
| Sắp hết hạn (30 ngày) | xxx | xxx VNĐ |
| Đã hết hạn | xxx | xxx VNĐ |
| Đang đàm phán | xxx | xxx VNĐ |
| Tranh chấp | xxx | xxx VNĐ |

```sql
-- Contract expiration alerts
SELECT c.contract_number, c.title, c.counterparty, c.contract_type,
    c.start_date, c.end_date, c.total_value,
    c.end_date - CURRENT_DATE AS days_remaining,
    CASE WHEN c.end_date < CURRENT_DATE THEN 'EXPIRED'
         WHEN c.end_date - CURRENT_DATE <= 30 THEN 'EXPIRING_SOON'
         WHEN c.end_date - CURRENT_DATE <= 90 THEN 'REVIEW_NEEDED'
         ELSE 'ACTIVE' END AS alert_status
FROM contracts c
WHERE c.org_id = $1 AND c.status != 'terminated'
ORDER BY c.end_date ASC;
```

### 3. Tuân thủ Pháp luật (Legal Compliance)

#### Checklist Giấy phép & Chứng nhận
| Loại | Số GP | Ngày cấp | Ngày hết hạn | Cơ quan cấp | Trạng thái |
|------|------|---------|-------------|-----------|-----------|
| Giấy ĐKKD | xxx | dd/mm/yy | Không hạn | Sở KHĐT | ✅ |
| Giấy phép con | xxx | dd/mm/yy | dd/mm/yy | [CQ] | ⚠️ Sắp HH |
| ISO 9001 | xxx | dd/mm/yy | dd/mm/yy | [TC] | ✅ |
| PCCC | xxx | dd/mm/yy | dd/mm/yy | CA PCCC | ✅ |

#### Theo dõi Vụ việc Pháp lý
| Vụ việc | Loại | Bên liên quan | Giá trị tranh chấp | Tiến độ | Rủi ro |
|---------|------|-------------|-------------------|---------|--------|
| [Mô tả] | Lao động | [Tên] | xxx VNĐ | Đang xử lý | Cao |

### 4. Chi phí Văn phòng & Hành chính

| Khoản mục | Tháng | Lũy kế | Ngân sách | % NS | So cùng kỳ |
|-----------|-------|--------|---------|------|-----------|
| Thuê văn phòng | xxx | xxx | xxx | x% | ±x% |
| Điện / Nước | xxx | xxx | xxx | x% | ±x% |
| Viễn thông / Internet | xxx | xxx | xxx | x% | ±x% |
| Văn phòng phẩm | xxx | xxx | xxx | x% | ±x% |
| Bảo trì / Sửa chữa | xxx | xxx | xxx | x% | ±x% |
| Bảo vệ / Vệ sinh | xxx | xxx | xxx | x% | ±x% |
| Đi lại / Công tác | xxx | xxx | xxx | x% | ±x% |
| Tiếp khách | xxx | xxx | xxx | x% | ±x% |
| **Tổng** | **xxx** | **xxx** | **xxx** | **x%** | **±x%** |

### 5. Quản lý Xe cộ & Phương tiện

| Xe | Biển số | Km đầu kỳ | Km cuối kỳ | Nhiên liệu | Bảo trì | Phí khác | Tổng CP |
|----|---------|----------|----------|-----------|---------|---------|---------|
| [Loại] | xxx | xxx | xxx | xxx | xxx | xxx | xxx |

### 6. Báo cáo Bảo hiểm Tài sản

| Tài sản bảo hiểm | Hãng BH | Giá trị BH | Phí BH/năm | Hết hạn | Claims |
|------------------|---------|-----------|----------|--------|--------|
| Văn phòng | [Hãng] | xxx VNĐ | xxx VNĐ | dd/mm | 0 |
| Xe cộ | [Hãng] | xxx VNĐ | xxx VNĐ | dd/mm | 1 |
| Hàng hóa | [Hãng] | xxx VNĐ | xxx VNĐ | dd/mm | 0 |

## Tần suất

| Báo cáo | Tháng | Quý | Năm |
|---------|-------|-----|-----|
| Chi phí hành chính | ✅ | | |
| Hợp đồng sắp hết hạn | ✅ | | |
| Khấu hao TSCĐ | ✅ | | |
| Compliance checklist | | ✅ | |
| Kiểm kê tài sản | | | ✅ |
| Tổng hợp pháp lý | | ✅ | ✅ |
