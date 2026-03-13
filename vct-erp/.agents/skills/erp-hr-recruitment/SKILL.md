---
name: erp-hr-recruitment
description: Nhân sự & Tuyển dụng - Báo cáo nhân sự tổng hợp, tuyển dụng, lương & phúc lợi, đào tạo, đánh giá hiệu quả, chấm công, turnover analysis cho VCT Platform ERP.
---

# Nhân sự & Tuyển dụng - Báo cáo Doanh nghiệp

## Tổng quan
Module báo cáo Nhân sự & Tuyển dụng cung cấp toàn bộ hệ thống báo cáo phục vụ quản trị nhân lực.

## Danh mục Báo cáo

### 1. Báo cáo Nhân sự Tổng hợp (Headcount)

#### Biến động Nhân sự
```
Đầu kỳ:                     xxx người
+ Tuyển mới:                 xxx người
+ Chuyển đến (nội bộ):      xxx người
- Nghỉ việc:                xxx người
- Sa thải / Kỷ luật:       xxx người
= Cuối kỳ:                  xxx người
```

#### Cơ cấu Nhân sự
| Tiêu chí | Phân loại | Số lượng | Tỷ lệ % |
|----------|----------|---------|---------|
| Phòng ban | KD, MKT, HR, TC, IT, SX... | xxx | x% |
| Giới tính | Nam / Nữ | xxx | x% |
| Độ tuổi | <25, 25-35, 35-45, >45 | xxx | x% |
| Thâm niên | <1 năm, 1-3, 3-5, >5 năm | xxx | x% |
| Trình độ | ĐH, ThS, TS, TC, CĐ | xxx | x% |
| Loại HĐ | Chính thức, Thử việc, CTV | xxx | x% |

```sql
-- Headcount by department, gender
SELECT d.name AS department, e.gender, e.employment_type,
    COUNT(*) AS headcount,
    AVG(EXTRACT(YEAR FROM AGE(CURRENT_DATE, e.date_of_birth))) AS avg_age,
    AVG(EXTRACT(YEAR FROM AGE(CURRENT_DATE, e.hire_date))) AS avg_tenure_years
FROM employees e
JOIN departments d ON d.id = e.department_id
WHERE e.status = 'active' AND e.org_id = $1
GROUP BY GROUPING SETS ((d.name, e.gender, e.employment_type), (d.name), (e.gender), ())
ORDER BY d.name;
```

### 2. Báo cáo Tuyển dụng (Recruitment)

#### KPIs Tuyển dụng
| KPI | Công thức | Target |
|-----|----------|--------|
| Time-to-Fill | Ngày từ mở yêu cầu → nhận việc | < 30 ngày |
| Time-to-Hire | Ngày từ ứng tuyển → offer | < 15 ngày |
| Cost-per-Hire | Tổng chi TD / Số tuyển | < 5 triệu |
| Quality of Hire | % NV mới đạt KPI sau 6 tháng | > 80% |
| Offer Acceptance Rate | Offers accepted / Offers sent | > 85% |
| Source Effectiveness | Hires per source / Cost per source | Theo kênh |

```sql
-- Recruitment funnel
SELECT jr.title AS position, jr.department,
    COUNT(DISTINCT a.id) AS total_applications,
    COUNT(DISTINCT CASE WHEN a.stage = 'screening_passed' THEN a.id END) AS screened,
    COUNT(DISTINCT CASE WHEN a.stage = 'interview_1' THEN a.id END) AS interview_1,
    COUNT(DISTINCT CASE WHEN a.stage = 'offer_accepted' THEN a.id END) AS accepted,
    COUNT(DISTINCT CASE WHEN a.stage = 'hired' THEN a.id END) AS hired,
    AVG(CASE WHEN a.stage = 'hired' THEN EXTRACT(DAY FROM a.hired_date - jr.created_at) END) AS avg_time_to_fill
FROM job_requisitions jr
LEFT JOIN applications a ON a.requisition_id = jr.id
WHERE jr.org_id = $1 AND jr.created_at BETWEEN $2 AND $3
GROUP BY jr.id, jr.title, jr.department;
```

### 3. Báo cáo Lương & Phúc lợi (Payroll)

#### Bảng Lương Tháng
```
Lương cơ bản (gross)                    xxx
+ Phụ cấp (ăn trưa, đi lại, chức vụ)   xxx
+ Thưởng KPI / OT / Bonus               xxx
= TỔNG THU NHẬP                         xxx
- BHXH (8%) / BHYT (1.5%) / BHTN (1%)  (xxx)
- Thuế TNCN                            (xxx)
- Khấu trừ khác                        (xxx)
= THỰC LĨNH (NET)                       xxx
```

#### Tổng hợp Quỹ lương theo Phòng ban
| Khoản mục | Tháng trước | Tháng này | Lũy kế năm | Ngân sách | % NS |
|-----------|-----------|---------|----------|---------|------|
| Lương cơ bản | xxx | xxx | xxx | xxx | x% |
| Phụ cấp | xxx | xxx | xxx | xxx | x% |
| Thưởng | xxx | xxx | xxx | xxx | x% |
| BHXH/BHYT/BHTN (DN đóng) | xxx | xxx | xxx | xxx | x% |
| **Tổng chi phí nhân sự** | **xxx** | **xxx** | **xxx** | **xxx** | **x%** |

### 4. Báo cáo Đào tạo & Phát triển (L&D)

| Chỉ số | Mục tiêu | Thực tế | % Đạt |
|--------|---------|---------|-------|
| Tổng giờ đào tạo / NV / năm | 40h | xxx | x% |
| Tỷ lệ NV được đào tạo | 100% | x% | x% |
| Chi phí đào tạo / NV | 2 triệu | xxx | x% |
| Tỷ lệ hoàn thành khóa | > 90% | x% | x% |
| ROI đào tạo | > 150% | x% | x% |

### 5. Đánh giá Hiệu quả (Performance Review)

#### Thang đánh giá
| Mức | Điểm | Mô tả | Tỷ lệ % |
|-----|------|-------|---------|
| Outstanding | 5 | Vượt xa kỳ vọng | 5-10% |
| Exceeds | 4 | Vượt kỳ vọng | 20-25% |
| Meets | 3 | Đạt kỳ vọng | 50-55% |
| Below | 2 | Dưới kỳ vọng | 10-15% |
| Unsatisfactory | 1 | Không đạt | 0-5% |

```sql
-- Performance distribution
SELECT d.name AS department, COUNT(*) AS total,
    COUNT(CASE WHEN pr.rating = 5 THEN 1 END) AS outstanding,
    COUNT(CASE WHEN pr.rating = 4 THEN 1 END) AS exceeds,
    COUNT(CASE WHEN pr.rating = 3 THEN 1 END) AS meets,
    COUNT(CASE WHEN pr.rating <= 2 THEN 1 END) AS below,
    ROUND(AVG(pr.rating), 2) AS avg_rating
FROM performance_reviews pr
JOIN employees e ON e.id = pr.employee_id
JOIN departments d ON d.id = e.department_id
WHERE pr.review_period = $1 AND pr.org_id = $2
GROUP BY d.name ORDER BY avg_rating DESC;
```

### 6. Chấm công & Nghỉ phép

#### Tổng hợp Chấm công
| NV | Ngày công | Đi muộn | Về sớm | Nghỉ phép | Nghỉ ốm | OT (giờ) |
|----|----------|---------|--------|----------|---------|---------|
| [Tên] | xx/26 | x | x | x | x | xx |

#### Leave Balance Report
```sql
SELECT e.employee_code, e.full_name, d.name AS department, lt.name AS leave_type,
    lb.annual_entitlement, lb.used, lb.pending_approval,
    (lb.annual_entitlement + lb.carried_forward - lb.used - lb.pending_approval) AS remaining
FROM leave_balances lb
JOIN employees e ON e.id = lb.employee_id
JOIN departments d ON d.id = e.department_id
JOIN leave_types lt ON lt.id = lb.leave_type_id
WHERE lb.year = $1 AND lb.org_id = $2 AND e.status = 'active';
```

### 7. Turnover Analysis

#### KPIs
| Chỉ số | Công thức | Target |
|--------|----------|--------|
| Voluntary Turnover | NV tự nghỉ / Tổng NV BQ | < 10%/năm |
| Involuntary Turnover | NV bị sa thải / Tổng NV BQ | < 3%/năm |
| New Hire Turnover | NV nghỉ < 1 năm / Tổng tuyển mới | < 20% |
| Retention Rate | NV ở lại / Tổng NV BQ | > 85% |

#### Lý do Nghỉ việc
| Lý do | Số lượng | Tỷ lệ % | Hành động |
|-------|---------|---------|----------|
| Lương/thưởng thấp | xx | x% | Review C&B |
| Thiếu cơ hội thăng tiến | xx | x% | Career path |
| Quản lý trực tiếp | xx | x% | Leadership training |
| Work-life balance | xx | x% | Flexible working |

### 8. Bảo hiểm Bắt buộc

| Loại BH | NV đóng | DN đóng | Tổng | Cơ sở tính |
|---------|---------|---------|------|-----------|
| BHXH | 8% | 17.5% | 25.5% | Lương đóng BH |
| BHYT | 1.5% | 3% | 4.5% | Lương đóng BH |
| BHTN | 1% | 1% | 2% | Lương đóng BH |
| KPCĐ | - | 2% | 2% | Quỹ lương |

## Tần suất Báo cáo

| Báo cáo | Tháng | Quý | Năm |
|---------|-------|-----|-----|
| Bảng lương & Chấm công | ✅ | | |
| BHXH/BHYT tổng hợp | ✅ | | |
| Headcount & biến động | ✅ | ✅ | |
| Tuyển dụng pipeline | ✅ | ✅ | |
| Performance review | | ✅ | ✅ |
| Đào tạo & L&D | | ✅ | ✅ |
| Turnover analysis | | ✅ | ✅ |

## Quyền hạn

| Báo cáo | CHRO | TP HR | C&B | Recruitment | Trưởng phòng khác |
|---------|------|-------|-----|-------------|------------------|
| Headcount toàn cty | ✅ | ✅ | Xem | Xem | Phòng mình |
| Lương chi tiết | ✅ | ✅ | ✅ | Không | Không |
| Tuyển dụng | ✅ | ✅ | Không | ✅ | Phòng mình |
| Performance | ✅ | ✅ | Không | Không | Phòng mình |
| Turnover | ✅ | ✅ | Xem | Xem | Phòng mình |
