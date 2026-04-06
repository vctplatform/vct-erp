# VCT Platform - Design System & UI/UX Philosophy (V1.0)
*Tài liệu chuẩn hóa Ngôn ngữ Thiết kế lõi dành cho các dự án: vct-website, vct-erp, vct-platform, vct-business.*

---

## 1. Triết Lý Thiết Kế (Design Philosophy)
VCT Platform định vị là một "Hệ sinh thái Quản trị Võ thuật Công nghệ cao" (Fintech & Martial Arts). Do đó, ngôn ngữ thiết kế phải cộng hưởng được 2 yếu tố: **Sự Uy Nghiêm (Heritage)** của Võ thuật và **Công Nghệ Đột Phá (Cyber/Fintech)** của Tương lai.

*   **Tên Ngôn ngữ Thiết kế:** `Neo-Glassmorphism Cyber-Teal`
*   **Từ khóa (Keywords):** Quyền lực, Vĩ mô, Sắc nét, Xuyên thấu, Neon Cyan, High-End Enterprise.

---

## 2. Hệ Màu Tone Lõi (Color Palette)
Xóa bỏ sự phân mảnh màu sắc rực rỡ (Cam/Đỏ/Xanh Lá). Toàn bộ hệ sinh thái phải chia sẻ 1 nhịp đập màu sắc duy nhất để tạo tính đồng bộ thương hiệu (Brand Consistency).

*   **Background (Nền):** **Pitch Dark / Deep Slate** `#020617` (Nền đen sâu thẳm, không dùng xám sáng).
*   **Primary Accent (Màu Nhấn Chủ Đạo):** **Cyber Cyan / Electric Teal** (`#06B6D4` hoặc `#0D9488`). Dùng cho Menu, Button chính, Icon Highlight, Outline.
*   **Text Primary (Chữ Chính):** **Chrome White** `#F8FAFC` (Trắng sáng, độ tương phản cao trên nền đen).
*   **Text Secondary (Chữ Phụ/Mô tả):** **Delicate Gray** `#94A3B8` (Xám tinh tế, text không bị chói mắt khi đọc bảng biểu ERP dài).
*   **Card Background (Nền Khối):** `rgba(255, 255, 255, 0.03)` kết hợp hiệu ứng kính (Glassmorphism).

---

## 3. Typography (Nghệ Thuật Chữ)
Sử dụng quy tắc "Dual-Font" (2 Phông chữ) để tách biệt Marketing và Data ERP.
*   **Tiêu đề Lớn (Marketing / Landing Page):** Serif (Có chân) -> Gợi ý: `Playfair Display`, `Lora` hoặc Serif cao cấp. Mang lại sự quyền lực, uy nghiêm cho các Liên đoàn / Chủ tịch.
*   **Text Bảng Biểu, Dashboard (ERP / Mobile App):** Sans-Serif (Không chân) -> Gợi ý: `Inter`, `Roboto`, `SF Pro`. Định dạng Clean, chữ số thẳng hàng, dễ đọc dữ liệu tài chính dày đặc.

---

## 4. UI Components & Hiệu ứng (Effects)
Tuyệt đối tuân thủ quy tắc `Neo-Glassmorphism` (Kính mờ Tương lai).

*   **Borders (Viền):** Luôn dùng viền tàng hình mỏng `border: 1px solid rgba(255, 255, 255, 0.05)`. Khi Hover, đổi sang viền Cyan `border-color: rgba(6, 182, 212, 0.3)`.
*   **Glow & Box-Shadow (Phát Sáng):** Thay vì dùng đổ bóng đen (Drop shadow), sử dụng phát sáng Neon Cyan: `box-shadow: 0 0 25px rgba(6, 182, 212, 0.2)`. Tránh làm gắt, chỉ tỏa sáng mờ ảo (Soft Glow).
*   **Blur (Làm Mờ Kính):** Lớp phủ Modal, Navbar, Sidebar ERP phải có `backdrop-filter: blur(20px)` để tạo chiều sâu giao diện.

---

## 5. Quy Chuẩn Hình Ảnh & Graphic (Assets / AI Gen)
Hình ảnh minh họa, Background trên ERP hoặc Website phải tuân thủ chuẩn mức "Hologram/Cyberpunk".

*   **Tuyệt đối KHÔNG:** Dùng ảnh người thật chụp sáng trưng ghép vào nền đen. Không lạm dụng ảnh vector Flat Design (Kiểu hoạt hình). Không sử dụng CSS `filter: hue-rotate(...)` để bẻ màu ảnh gốc, dẫn đến rác pixel.
*   **Quy ước Gen AI/Graphic:** Khi yêu cầu AI (Midjourney/Gemini) hoặc Designer vẽ hình, phải kèm theo Keywords gốc: `NEON CYAN and ELECTRIC BLUE, floating, holographic, glassmorphism UI, pitch black background, hyper-realistic 3D`. Hình ảnh tự mang ánh sáng Xanh Cyan sẽ hòa quyện hoàn hảo vào hệ thống gốc.

---

## 6. Cách Áp Dụng Cho Team ERP & Mobile App (vct-erp / vct-business)
*Khi Javis (Team DB/ERP) và Jack (Team Business Logic) xây dựng Frontend:*
1. Thiết lập Dark Mode mặc định, không làm Light Mode.
2. Nút "Lưu", "Thanh Toán", "Đăng Ký" trên phần mềm ERP -> Màu `#06B6D4` (Dùng Tailwind: `bg-cyan-500`).
3. Card thống kê Doanh thu (Dashboard Võ Đường) -> Áp dụng Glassmorphism Card (Tailwind: `bg-slate-900/40 border border-slate-800 backdrop-blur-md`). Trang tính bảng Table dùng Text `text-slate-200`.
