# VCT Platform - Backend CMS

Đây là hệ thống Backend CMS xây dựng bằng **Golang (Fiber)** dành cho VCT Platform.

## 1. Stack Công nghệ
- **Framework**: [Fiber v2](https://gofiber.io/) (Nhanh, nhẹ, routing giống Express.js).
- **Database**: PostgreSQL (qua GORM) & thao tác json tĩnh (git-backed).
- **Docs**: Swagger (mặc định tại `/api/v1/swagger`).

## 2. Cấu trúc Thư mục (Clean Architecture)
```text
cms-backend/
├── cmd/
│   └── api/
│       └── main.go         # Entry point, setup Fiber app & Middleware (CORS)
├── internal/
│   ├── auth/               # Module Xác thực (JWT, RBAC)
│   ├── blog/               # Module Quản lý Blog (CRUD vào file posts.json)
│   ├── media/              # Module Xử lý file upload
│   └── settings/           # Module Cấu hình & Ngôn ngữ (i18n)
├── pkg/                    # Shared libraries (JWT config, Git utilities, Database conn)
├── docker-compose.yaml     # Script chạy Database PostgreSQL qua Docker
├── go.mod                  # Chứa dependencies của dự án Go
└── README.md
```

## 3. Tính năng cốt lõi (Khởi tạo đợt 1)
- **CORS API**: Đã cấu hình Middleware cho phép tên miền `http://localhost:3000` và `https://vct-platform.github.io` gọi API.
- **Git Integration Concept**: Module `blog` được thiết kế để đọc file tĩnh `posts.json` và lưu tự động sử dụng `go-git` trong các bản cập nhật sau.

## 4. Hướng dẫn Khởi chạy
1. **Khởi động Database (Tuỳ chọn nếu dùng GORM):**
   ```bash
   docker-compose up -d
   ```
2. **Chạy Server Backend:**
   ```bash
   go run cmd/api/main.go
   ```
3. Backend sẽ chạy ở cổng `http://localhost:8080`.
