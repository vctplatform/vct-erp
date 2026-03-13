---
name: erp-procurement-inventory
description: Mua hàng & Kho vận - Đơn đặt hàng, quản lý nhà cung cấp, tồn kho ABC, xuất nhập tồn, logistics & vận chuyển cho VCT Platform ERP.
---

# Mua hàng & Kho vận - Báo cáo Doanh nghiệp

## Tổng quan
Báo cáo Mua hàng & Kho vận bao phủ toàn bộ quy trình procurement-to-pay, quản lý tồn kho, đánh giá nhà cung cấp và logistics.

## Danh mục Báo cáo

### 1. Báo cáo Mua hàng (Procurement)

#### KPIs Mua hàng
| KPI | Công thức | Target |
|-----|----------|--------|
| Cost Savings | (Giá thị trường - Giá mua) / Giá TT | > 5% |
| PO Cycle Time | Ngày từ PR → PO phê duyệt | < 3 ngày |
| Supplier On-time Delivery | Giao đúng hạn / Tổng đơn | > 95% |
| PO Accuracy | PO đúng / Tổng PO | > 98% |
| Maverick Spend | Chi ngoài HĐ / Tổng chi mua | < 5% |

```sql
-- Procurement spend analysis by category & supplier
SELECT s.name AS supplier, c.name AS category,
    COUNT(po.id) AS total_pos,
    SUM(po.total_amount) AS total_spend,
    AVG(po.total_amount) AS avg_po_value,
    ROUND(AVG(EXTRACT(DAY FROM po.delivery_date - po.order_date)), 1) AS avg_lead_time,
    COUNT(CASE WHEN po.delivered_on_time THEN 1 END)::numeric / 
        NULLIF(COUNT(po.id), 0) * 100 AS on_time_pct
FROM purchase_orders po
JOIN suppliers s ON s.id = po.supplier_id
JOIN categories c ON c.id = po.category_id
WHERE po.org_id = $1 AND po.order_date BETWEEN $2 AND $3
GROUP BY s.name, c.name ORDER BY total_spend DESC;
```

### 2. Báo cáo Tồn kho (Inventory)

#### Phân tích ABC
| Nhóm | Tiêu chí | % SKU | % Giá trị | Chiến lược |
|------|---------|-------|----------|-----------|
| A | Top 80% giá trị | ~20% | ~80% | Kiểm soát chặt |
| B | Tiếp theo 15% | ~30% | ~15% | Kiểm soát vừa |
| C | Còn lại 5% | ~50% | ~5% | Kiểm soát lỏng |

```sql
-- ABC analysis
WITH item_values AS (
    SELECT i.id, i.sku, i.name, i.quantity * i.unit_cost AS total_value,
        SUM(i.quantity * i.unit_cost) OVER () AS grand_total
    FROM inventory_items i WHERE i.org_id = $1 AND i.quantity > 0
),
ranked AS (
    SELECT *, SUM(total_value) OVER (ORDER BY total_value DESC) / grand_total * 100 AS cumulative_pct
    FROM item_values
)
SELECT sku, name, total_value,
    CASE WHEN cumulative_pct <= 80 THEN 'A'
         WHEN cumulative_pct <= 95 THEN 'B' ELSE 'C' END AS abc_class
FROM ranked ORDER BY total_value DESC;
```

#### Báo cáo Xuất Nhập Tồn
| Mã HH | Tên | ĐVT | Tồn đầu kỳ | Nhập trong kỳ | Xuất trong kỳ | Tồn cuối kỳ | Giá trị tồn |
|-------|-----|-----|-----------|-------------|-------------|-----------|-----------|
| HH001 | [Tên] | Cái | xxx | xxx | xxx | xxx | xxx VNĐ |

#### KPIs Kho
| KPI | Công thức | Target |
|-----|----------|--------|
| Inventory Turnover | GVHB / Tồn kho BQ | > 6 lần/năm |
| Days of Inventory | 365 / Inventory Turnover | < 60 ngày |
| Stock Accuracy | Kiểm kê đúng / Tổng SKU | > 99% |
| Fill Rate | Đáp ứng đúng / Yêu cầu xuất | > 98% |
| Dead Stock % | Hàng không luân chuyển > 6 tháng | < 5% |
| Shrinkage Rate | Hao hụt / Tổng tồn | < 1% |

### 3. Đánh giá Nhà cung cấp

| Tiêu chí | Trọng số | Điểm (1-5) |
|----------|---------|-----------|
| Chất lượng sản phẩm | 30% | xxx |
| Giao hàng đúng hạn | 25% | xxx |
| Giá cạnh tranh | 20% | xxx |
| Dịch vụ hậu mãi | 15% | xxx |
| Năng lực tài chính | 10% | xxx |
| **Tổng điểm** | **100%** | **x.xx** |

### 4. Báo cáo Logistics

| Chỉ số | Tháng trước | Tháng này | Chênh lệch |
|--------|-----------|---------|-----------|
| Tổng đơn vận chuyển | xxx | xxx | ±x% |
| Chi phí vận chuyển | xxx VNĐ | xxx VNĐ | ±x% |
| Tỷ lệ giao đúng hạn | x% | x% | ±x% |
| Tỷ lệ đơn hỏng/trễ | x% | x% | ±x% |
| Lead time trung bình | x ngày | x ngày | ±x |

## Tần suất

| Báo cáo | Ngày | Tuần | Tháng | Quý | Năm |
|---------|-----|------|-------|-----|-----|
| Xuất nhập tồn | ✅ | | | | |
| Cảnh báo tồn kho thấp | ✅ | | | | |
| PO tracking | | ✅ | | | |
| Spend analysis | | | ✅ | | |
| ABC analysis | | | | ✅ | |
| Supplier scorecard | | | | ✅ | ✅ |
| Kiểm kê | | | | | ✅ |
