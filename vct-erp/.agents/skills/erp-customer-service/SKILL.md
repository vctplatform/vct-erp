---
name: erp-customer-service
description: CSKH & Hậu mãi - Ticket support, SLA compliance, NPS/CSAT/CES, bảo hành & trả hàng, churn analysis, knowledge base cho VCT Platform ERP.
---

# CSKH & Hậu mãi - Báo cáo Doanh nghiệp

## Tổng quan
Báo cáo CSKH & Hậu mãi bao phủ toàn bộ dịch vụ sau bán hàng, từ ticket support, SLA, khảo sát hài lòng đến bảo hành và phân tích churn.

## Danh mục Báo cáo

### 1. Báo cáo Ticket / Support

#### Dashboard Support
| Chỉ số | Hôm nay | Tuần này | Tháng này |
|--------|---------|---------|---------|
| Ticket mới | xxx | xxx | xxx |
| Ticket đang xử lý | xxx | xxx | xxx |
| Ticket đã giải quyết | xxx | xxx | xxx |
| Ticket quá hạn SLA | xxx | xxx | xxx |
| Avg Response Time | x phút | x phút | x phút |
| Avg Resolution Time | x giờ | x giờ | x giờ |

```sql
-- Ticket analytics
SELECT
    DATE_TRUNC('week', t.created_at) AS week,
    t.priority, t.category,
    COUNT(*) AS total_tickets,
    COUNT(CASE WHEN t.status = 'resolved' THEN 1 END) AS resolved,
    COUNT(CASE WHEN t.sla_breached THEN 1 END) AS sla_breached,
    ROUND(AVG(EXTRACT(EPOCH FROM t.first_response_at - t.created_at) / 60), 1) AS avg_response_min,
    ROUND(AVG(EXTRACT(EPOCH FROM t.resolved_at - t.created_at) / 3600), 1) AS avg_resolution_hours,
    ROUND(COUNT(CASE WHEN t.status = 'resolved' AND NOT t.sla_breached THEN 1 END)::numeric /
          NULLIF(COUNT(*), 0) * 100, 1) AS sla_compliance_pct
FROM support_tickets t
WHERE t.org_id = $1 AND t.created_at BETWEEN $2 AND $3
GROUP BY DATE_TRUNC('week', t.created_at), t.priority, t.category
ORDER BY week DESC, t.priority;
```

#### SLA Definitions
| Ưu tiên | First Response | Resolution | Escalation |
|---------|---------------|-----------|-----------|
| Critical (P1) | 15 phút | 4 giờ | VP ngay |
| High (P2) | 1 giờ | 8 giờ | Manager sau 4h |
| Medium (P3) | 4 giờ | 24 giờ | Manager sau 12h |
| Low (P4) | 8 giờ | 72 giờ | Supervisor sau 48h |

#### KPIs Support
| KPI | Công thức | Target |
|-----|----------|--------|
| First Response Time | Avg thời gian phản hồi đầu | < 1 giờ |
| Resolution Time | Avg thời gian giải quyết | < 24 giờ |
| First Contact Resolution | Giải quyết lần đầu / Tổng ticket | > 70% |
| SLA Compliance | Đúng SLA / Tổng ticket | > 95% |
| Customer Effort Score | Survey score | < 3/7 |
| Reopened Rate | Ticket mở lại / Tổng resolved | < 5% |
| Backlog | Ticket chưa giải quyết | Giảm WoW |
| Agent Utilization | Ticket / Agent / ngày | 15-25 |

### 2. Khảo sát Hài lòng Khách hàng

#### NPS (Net Promoter Score)
```
NPS = % Promoters (9-10) - % Detractors (0-6)
Target: > 50
```

| Nhóm | Điểm | Tỷ lệ | Hành động |
|------|------|-------|----------|
| Promoters | 9-10 | x% | Referral program |
| Passives | 7-8 | x% | Engage & upsell |
| Detractors | 0-6 | x% | Recovery program |

#### CSAT (Customer Satisfaction Score)
| Kênh | Responses | Avg Score | CSAT % (4-5) |
|------|----------|----------|-------------|
| Phone | xxx | x.x/5 | x% |
| Email | xxx | x.x/5 | x% |
| Chat | xxx | x.x/5 | x% |
| In-person | xxx | x.x/5 | x% |

#### CES (Customer Effort Score)
| Touchpoint | Score (1-7) | Target |
|-----------|-----------|--------|
| Đặt hàng | x.x | < 3.0 |
| Thanh toán | x.x | < 3.0 |
| Hỗ trợ kỹ thuật | x.x | < 3.0 |
| Trả hàng / Đổi | x.x | < 3.0 |

### 3. Bảo hành & Trả hàng

#### Dashboard Bảo hành
| Chỉ số | Tháng trước | Tháng này | Trend |
|--------|-----------|---------|-------|
| Claims mới | xxx | xxx | ±x% |
| Claims đang xử lý | xxx | xxx | ±x% |
| Claims hoàn tất | xxx | xxx | ±x% |
| Avg processing time | x ngày | x ngày | ±x |
| Warranty cost | xxx VNĐ | xxx VNĐ | ±x% |

#### KPIs Bảo hành
| KPI | Công thức | Target |
|-----|----------|--------|
| Warranty Claim Rate | Claims / SP bán | < 2% |
| Return Rate | SP trả / SP bán | < 1% |
| Avg Warranty Cost | Tổng chi BH / Claims | Giảm QoQ |
| Processing Time | Avg ngày xử lý claim | < 5 ngày |
| Repeat Claim Rate | Claims lặp / Tổng claims | < 10% |

### 4. Phân tích Churn (Customer Retention)

```sql
-- Churn analysis
WITH customer_activity AS (
    SELECT c.id, c.name, c.segment,
        MAX(o.order_date) AS last_order_date,
        COUNT(o.id) AS total_orders,
        SUM(o.total_amount) AS lifetime_value,
        CURRENT_DATE - MAX(o.order_date) AS days_since_last_order
    FROM customers c
    LEFT JOIN orders o ON o.customer_id = c.id AND o.status = 'completed'
    WHERE c.org_id = $1
    GROUP BY c.id
)
SELECT segment,
    COUNT(*) AS total_customers,
    COUNT(CASE WHEN days_since_last_order <= 90 THEN 1 END) AS active,
    COUNT(CASE WHEN days_since_last_order BETWEEN 91 AND 180 THEN 1 END) AS at_risk,
    COUNT(CASE WHEN days_since_last_order > 180 THEN 1 END) AS churned,
    ROUND(COUNT(CASE WHEN days_since_last_order > 180 THEN 1 END)::numeric / 
          NULLIF(COUNT(*), 0) * 100, 1) AS churn_rate
FROM customer_activity GROUP BY segment;
```

#### Retention Metrics
| Metric | Công thức | Target |
|--------|----------|--------|
| Customer Retention Rate | KH giữ lại / KH đầu kỳ | > 85% |
| Revenue Retention Rate | MRR giữ lại / MRR đầu kỳ | > 90% |
| Churn Rate | KH mất / KH đầu kỳ | < 5%/quý |
| Win-back Rate | KH quay lại / KH churned | > 15% |
| CLV (Customer Lifetime Value) | ARPU × Avg lifespan | Tăng YoY |

### 5. Knowledge Base Analytics

| Chỉ số | Giá trị | Trend |
|--------|---------|-------|
| Tổng bài viết | xxx | +x/tháng |
| Lượt xem/tháng | xxx | ±x% |
| Search success rate | x% | ±x% |
| Helpful rate (thumbs up) | x% | ±x% |
| Ticket deflection rate | x% | ±x% |

## Tần suất

| Báo cáo | Ngày | Tuần | Tháng | Quý |
|---------|-----|------|-------|-----|
| Ticket dashboard | ✅ | | | |
| SLA compliance | | ✅ | | |
| NPS/CSAT | | | ✅ | ✅ |
| Bảo hành | | | ✅ | |
| Churn analysis | | | | ✅ |
| Knowledge base | | | ✅ | |

## Quyền hạn

| Báo cáo | VP CSKH | Manager | Team Lead | Agent |
|---------|---------|---------|----------|-------|
| Toàn bộ BC | ✅ | ✅ | Xem | Không |
| NPS/CSAT | ✅ | ✅ | ✅ | Không |
| Ticket analytics | ✅ | ✅ | Team | Cá nhân |
| Churn | ✅ | ✅ | Không | Không |
