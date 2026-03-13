---
name: erp-planning-strategy
description: Kế hoạch & Chiến lược - Lập kế hoạch kinh doanh, quản lý OKR/KPI, phân bổ ngân sách, quản lý dự án nội bộ, đánh giá hiệu quả chiến lược cho VCT Platform ERP.
---

# Kế hoạch & Chiến lược - ERP Module

## Tổng quan
Phòng Kế hoạch & Chiến lược chịu trách nhiệm xây dựng kế hoạch kinh doanh, theo dõi mục tiêu chiến lược, phân bổ nguồn lực, và đánh giá hiệu quả hoạt động toàn doanh nghiệp.

## Vai trò & Nhân sự

| Vai trò | Trách nhiệm |
|---------|-------------|
| Giám đốc Chiến lược | Định hướng chiến lược, phối hợp BGĐ |
| Trưởng phòng Kế hoạch | Lập KH kinh doanh, phân bổ ngân sách |
| Chuyên viên Kế hoạch | Tổng hợp kế hoạch phòng ban, theo dõi tiến độ |
| Chuyên viên Phân tích | Phân tích dữ liệu, báo cáo hiệu quả |
| Chuyên viên Dự án | Quản lý dự án nội bộ, PMO |

## Quy trình nghiệp vụ

### 1. Lập Kế hoạch Kinh doanh Hàng năm

```mermaid
flowchart TD
    A[Đánh giá kết quả năm trước] --> B[Phân tích thị trường & xu hướng]
    B --> C[Xác định mục tiêu chiến lược]
    C --> D[Phòng ban đề xuất KH]
    D --> E[Tổng hợp & Cân đối]
    E --> F[Phân bổ ngân sách]
    F --> G[Trình BGĐ/HĐQT]
    G --> H{Phê duyệt?}
    H -->|Duyệt| I[Ban hành KH chính thức]
    H -->|Chỉnh sửa| D
    I --> J[Triển khai theo quý/tháng]
    J --> K[Theo dõi & Điều chỉnh]
```

#### Nội dung Kế hoạch Kinh doanh
```markdown
# Kế hoạch Kinh doanh [Năm]

## 1. Tầm nhìn & Sứ mệnh
## 2. Phân tích Môi trường
  - SWOT Analysis
  - PESTEL Analysis
  - Porter's Five Forces
## 3. Mục tiêu Chiến lược
  - Định lượng (doanh thu, lợi nhuận, thị phần)
  - Định tính (thương hiệu, năng lực, văn hóa)
## 4. Kế hoạch Hành động theo Phòng ban
## 5. Phân bổ Ngân sách
## 6. Timeline & Milestones
## 7. KPIs & Cơ chế Đánh giá
## 8. Phương án Dự phòng
```

### 2. Quản lý OKR/KPI theo Phòng ban

```mermaid
flowchart TD
    A[CEO đặt OKR công ty] --> B[Phòng ban cascade OKR]
    B --> C[Cá nhân đăng ký OKR]
    C --> D[Review & Phê duyệt]
    D --> E[Thực hiện & Cập nhật tiến độ]
    E --> F[Check-in hàng tuần]
    F --> G[Review giữa kỳ]
    G --> H{Cần điều chỉnh?}
    H -->|Có| I[Điều chỉnh OKR]
    H -->|Không| J[Tiếp tục thực hiện]
    I --> E
    J --> K[Đánh giá cuối kỳ]
    K --> L[Scoring & Retrospective]
```

#### KPIs theo Phòng ban
| Phòng ban | KPI chính | Đơn vị | Tần suất đo |
|-----------|----------|--------|------------|
| Sales | Doanh thu, số HĐ mới, conversion rate | VNĐ, %, số | Tuần |
| Marketing | Leads, CAC, brand awareness | Số, VNĐ, % | Tháng |
| Finance | Cash flow, ROI, chi phí/DT | VNĐ, % | Tháng |
| HR | Turnover rate, time-to-hire, eNPS | %, ngày, điểm | Quý |
| Accounting | Accuracy rate, closing time | %, ngày | Tháng |
| Admin | SLA compliance, cost savings | %, VNĐ | Quý |
| Procurement | Cost savings, supplier score | %, điểm | Quý |

### 3. Phân bổ Ngân sách theo Dự án/Phòng ban

```mermaid
flowchart LR
    A[Tổng ngân sách] --> B[Chi phí cố định - 40%]
    A --> C[Chi phí biến đổi - 35%]
    A --> D[Đầu tư phát triển - 15%]
    A --> E[Dự phòng - 10%]
    
    B --> B1[Lương & phúc lợi]
    B --> B2[Thuê VP & hạ tầng]
    B --> B3[Khấu hao TSCĐ]
    
    C --> C1[Marketing & bán hàng]
    C --> C2[Nguyên vật liệu]
    C --> C3[Vận chuyển & logistics]
    
    D --> D1[R&D / Công nghệ]
    D --> D2[Đào tạo nhân sự]
    D --> D3[Mở rộng thị trường]
```

#### Quy trình Phê duyệt Ngân sách
| Bước | Người thực hiện | Thời gian | Output |
|------|----------------|----------|--------|
| 1. Đề xuất | Trưởng phòng | T10 | Bản đề xuất NS |
| 2. Tổng hợp | Phòng KH | T11 đầu | Bản tổng hợp |
| 3. Cân đối | CFO + Phòng KH | T11 giữa | Bản cân đối |
| 4. Phê duyệt | CEO/HĐQT | T12 đầu | Nghị quyết NS |
| 5. Ban hành | Phòng KH | T12 giữa | KH ngân sách chính thức |

### 4. Quản lý Dự án Nội bộ (PMO)

```mermaid
flowchart TD
    A[Đề xuất dự án] --> B[Đánh giá tính khả thi]
    B --> C{Khả thi?}
    C -->|Có| D[Lập kế hoạch dự án]
    C -->|Không| E[Từ chối & Ghi nhận]
    D --> F[Phê duyệt dự án]
    F --> G[Kick-off]
    G --> H[Thực hiện theo phase]
    H --> I[Báo cáo tiến độ tuần]
    I --> J{Đúng tiến độ?}
    J -->|Có| K[Tiếp tục]
    J -->|Không| L[Escalation & Điều chỉnh]
    K --> M[Nghiệm thu]
    L --> H
    M --> N[Đóng dự án & Lessons Learned]
```

#### Template Gantt Chart
| Task | Người PTV | W1 | W2 | W3 | W4 | W5 | W6 | W7 | W8 |
|------|----------|----|----|----|----|----|----|----|----|
| Khởi động | PM | ██ | | | | | | | |
| Phân tích | BA | | ██ | ██ | | | | | |
| Thiết kế | Architect | | | ██ | ██ | | | | |
| Phát triển | Dev | | | | ██ | ██ | ██ | | |
| Kiểm thử | QA | | | | | | ██ | ██ | |
| Triển khai | DevOps | | | | | | | | ██ |

### 5. Đánh giá Hiệu quả Chiến lược

#### Framework đánh giá
| Tiêu chí | Trọng số | Thang điểm | Mô tả |
|----------|---------|-----------|--------|
| Đạt mục tiêu DT | 30% | 1-5 | So với KH ban đầu |
| Hiệu quả chi phí | 20% | 1-5 | Chi phí/doanh thu |
| Tăng trưởng KH | 15% | 1-5 | Số KH mới, giữ chân |
| Phát triển nhân sự | 15% | 1-5 | Năng lực đội ngũ |
| Đổi mới sáng tạo | 10% | 1-5 | Sáng kiến, cải tiến |
| Tuân thủ pháp luật | 10% | 1-5 | Vi phạm, rủi ro |

### 6. Quy trình Phê duyệt Đề xuất

```mermaid
flowchart LR
    A[Nhân viên đề xuất] --> B[Trưởng nhóm review]
    B --> C[Trưởng phòng phê duyệt]
    C --> D{Giá trị > ngưỡng?}
    D -->|Nhỏ| E[Thực hiện]
    D -->|TB| F[Ban GĐ phê duyệt]
    D -->|Lớn| G[HĐQT phê duyệt]
    F --> H[Thực hiện]
    G --> H
    H --> I[Báo cáo kết quả]
```

## Báo cáo Định kỳ

| Báo cáo | Tần suất | Người nhận | Nội dung |
|---------|---------|-----------|---------|
| Flash Report | Hàng ngày | CEO, COO | Doanh thu, vấn đề nổi bật |
| Weekly Report | Tuần | BGĐ | Tiến độ KPI, vấn đề cần giải quyết |
| Monthly Report | Tháng | BGĐ, HĐQT | Kết quả kinh doanh, phân tích |
| Quarterly Review | Quý | HĐQT | OKR progress, điều chỉnh chiến lược |
| Annual Report | Năm | HĐQT, cổ đông | Tổng kết toàn diện |

## Quyền hạn trong ERP

| Chức năng | GĐ Chiến lược | TP Kế hoạch | CV Kế hoạch | CV Phân tích |
|-----------|--------------|-------------|-------------|-------------|
| Xem KPIs toàn công ty | ✅ | ✅ | Phòng ban phụ trách | Phòng ban phụ trách |
| Tạo/sửa OKR | Công ty | Phòng ban | Đề xuất | Không |
| Phê duyệt ngân sách | Đề xuất BGĐ | Tổng hợp | Nhập liệu | Không |
| Quản lý dự án | Portfolio | Chương trình | Dự án đơn lẻ | Không |
| Báo cáo | Tất cả | Tất cả | Phòng ban | Phân tích dữ liệu |
