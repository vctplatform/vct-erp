---
name: erp-production-quality
description: Sản xuất & Chất lượng - Kế hoạch sản xuất, năng suất OEE, kiểm soát chất lượng, bảo trì thiết bị, chi phí sản xuất, BOM cho VCT Platform ERP.
---

# Sản xuất & Chất lượng - Báo cáo Doanh nghiệp

## Tổng quan
Báo cáo Sản xuất & Chất lượng bao phủ toàn bộ quy trình từ lập kế hoạch sản xuất, theo dõi năng suất, kiểm soát chất lượng đến bảo trì thiết bị và phân tích chi phí.

## Danh mục Báo cáo

### 1. Kế hoạch Sản xuất vs Thực tế

| Sản phẩm | KH sản lượng | TT sản lượng | Đạt % | KH thời gian | TT thời gian | Đúng hạn? |
|----------|-------------|-------------|-------|-------------|-------------|----------|
| SP-001 | xxx | xxx | x% | dd/mm | dd/mm | ✅/❌ |

```sql
-- Production plan vs actual
SELECT p.product_code, p.product_name,
    pp.planned_quantity, pp.actual_quantity,
    ROUND(pp.actual_quantity::numeric / NULLIF(pp.planned_quantity, 0) * 100, 1) AS achievement_pct,
    pp.planned_end_date, pp.actual_end_date,
    CASE WHEN pp.actual_end_date <= pp.planned_end_date THEN true ELSE false END AS on_time,
    pp.planned_cost, pp.actual_cost,
    pp.actual_cost - pp.planned_cost AS cost_variance
FROM production_plans pp
JOIN products p ON p.id = pp.product_id
WHERE pp.org_id = $1 AND pp.period = $2
ORDER BY achievement_pct ASC;
```

### 2. Hiệu suất Thiết bị Tổng thể (OEE)

#### Công thức OEE
```
OEE = Availability × Performance × Quality

Availability = Thời gian chạy thực tế / Thời gian kế hoạch
Performance  = Sản lượng thực tế / Sản lượng lý thuyết
Quality      = Sản phẩm đạt / Tổng sản lượng
```

| Thiết bị | Availability | Performance | Quality | OEE | Benchmark |
|----------|-------------|-------------|---------|-----|-----------|
| Máy 01 | x% | x% | x% | x% | ≥ 85% |
| Máy 02 | x% | x% | x% | x% | ≥ 85% |

#### Phân tích Downtime
| Lý do dừng máy | Tổng thời gian | Tần suất | TB mỗi lần | % Tổng DT |
|----------------|---------------|---------|-----------|----------|
| Hỏng máy | xxx phút | xx lần | xx phút | x% |
| Chuyển đổi sản phẩm | xxx phút | xx lần | xx phút | x% |
| Thiếu nguyên liệu | xxx phút | xx lần | xx phút | x% |
| Bảo trì định kỳ | xxx phút | xx lần | xx phút | x% |

### 3. Kiểm soát Chất lượng (QC Report)

#### KPIs Chất lượng
| KPI | Công thức | Target |
|-----|----------|--------|
| First Pass Yield | SP đạt lần đầu / Tổng SP | > 98% |
| Defect Rate | SP lỗi / Tổng SP | < 2% |
| Customer Complaint Rate | Khiếu nại / Đơn giao | < 0.5% |
| Return Rate | SP trả lại / SP giao | < 1% |
| Cost of Quality | Chi phí CL / Doanh thu | < 3% |
| COPQ (Cost of Poor Quality) | Phế phẩm + Sửa chữa + Trả hàng | < 1% |

#### Pareto Analysis - Lỗi phổ biến
| Loại lỗi | Số lượng | Tỷ lệ % | Lũy kế % | Hành động |
|----------|---------|---------|---------|----------|
| Kích thước sai | xxx | x% | x% | SPC control |
| Trầy xước | xxx | x% | x% | Process review |
| Thiếu chi tiết | xxx | x% | x% | Checklist |
| Sai màu | xxx | x% | x% | Color match |

### 4. Bảo trì Thiết bị (Maintenance)

| KPI | Công thức | Target |
|-----|----------|--------|
| MTBF (Mean Time Between Failures) | Tổng giờ chạy / Số lần hỏng | > 500 giờ |
| MTTR (Mean Time To Repair) | Tổng giờ sửa / Số lần sửa | < 4 giờ |
| Planned Maintenance % | Bảo trì KH / Tổng bảo trì | > 80% |
| Maintenance Cost Ratio | Chi phí BT / Giá trị TB | < 5%/năm |

#### Lịch Bảo trì Định kỳ
| Thiết bị | Loại BT | Tần suất | Lần cuối | Lần tiếp | Trạng thái |
|----------|---------|---------|---------|---------|-----------|
| Máy 01 | Thay dầu | 3 tháng | dd/mm | dd/mm | Đúng hạn |
| Máy 02 | Kiểm tra | 1 tháng | dd/mm | dd/mm | ⚠️ Quá hạn |

### 5. Chi phí Sản xuất (Production Cost)

#### BOM Cost Analysis
| Thành phần | ĐVT | SL/SP | Đơn giá | Thành tiền | % Tổng |
|-----------|-----|------|---------|-----------|--------|
| Nguyên liệu chính | kg | x | xxx | xxx | x% |
| Nguyên liệu phụ | kg | x | xxx | xxx | x% |
| Nhân công trực tiếp | giờ | x | xxx | xxx | x% |
| Chi phí SX chung | | | | xxx | x% |
| **Tổng giá thành** | | | | **xxx** | **100%** |

#### So sánh Giá thành Định mức vs Thực tế
| Khoản mục | Định mức | Thực tế | Chênh lệch | Nguyên nhân |
|-----------|---------|---------|-----------|-----------|
| NVL | xxx | xxx | ±xxx | Giá NVL tăng |
| Nhân công | xxx | xxx | ±xxx | OT tăng |
| SXC | xxx | xxx | ±xxx | Điện tăng |

### 6. Báo cáo Năng suất Lao động

| Ca/Line | Sản lượng | Nhân công | Năng suất/người | Target | Đạt? |
|---------|----------|----------|----------------|--------|-----|
| Ca sáng | xxx | xx | xxx SP/người | xxx | ✅/❌ |
| Ca chiều | xxx | xx | xxx SP/người | xxx | ✅/❌ |

## Tần suất

| Báo cáo | Ca | Ngày | Tuần | Tháng | Quý |
|---------|---|------|------|-------|-----|
| Sản lượng | ✅ | | | | |
| Chất lượng (defect) | | ✅ | | | |
| OEE | | ✅ | ✅ | | |
| KH vs Thực tế | | | ✅ | | |
| Chi phí SX | | | | ✅ | |
| Bảo trì | | | | ✅ | |
| Pareto lỗi | | | | | ✅ |
