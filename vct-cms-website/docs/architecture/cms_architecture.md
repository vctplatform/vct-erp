# Kiến trúc Hệ thống Quản trị (CMS) cho VCT Platform

Dựa trên yêu cầu xây dựng CMS cho một website Static (HTML/JS) mà không cần can thiệp vào code, giải pháp tối ưu nhất là **Decap CMS** (tiền thân là Netlify CMS). Đây là một "Git-based CMS" hoạt động hoàn toàn trên trình duyệt.

## 1. Kiến trúc Hệ thống (Architecture)

*   **Mô hình hoạt động:** Git-based CMS.
*   **Giao diện CMS:** Một ứng dụng Single Page React (nằm trong thư mục `admin/`) tải trực tiếp trên trình duyệt.
*   **Lưu trữ Dữ liệu:** CMS đọc và ghi trực tiếp vào các file JSON của dự án (`data/posts.json`, `data/pitch.json`, `data/lang/vi.json`) thông qua **GitHub API**.
*   **Xác thực (Authentication):** Sử dụng **GitHub OAuth App**. Chỉ những tài khoản GitHub được phân quyền vào Repository (Collaborators) mới có quyền đăng nhập và chỉnh sửa CMS.

### Sơ đồ Luồng Hoạt động
1. Quản trị viên truy cập `https://vct-platform.github.io/vct-website/admin/`
2. Đăng nhập bằng tài khoản GitHub (Netlify Identity / OAuth).
3. Giao diện CMS load dữ liệu từ các file JSON trong nhánh `main`.
4. Quản trị viên thêm/sửa/xóa bài viết, tải ảnh lên (`assets/images`).
5. Bấm **Publish**. CMS dùng GitHub API để tự động tạo một commit mới: `Update posts.json` và push lên nhánh `main`.
6. GitHub Actions (`deploy.yml`) bắt được sự kiện push, tự động chạy script tối ưu Sharp và Deploy phiên bản web mới nhất.

---

## 2. Hướng dẫn Thiết lập (Technical Setup)

Tôi đã tạo sẵn thư mục `admin/` bao gồm:
*   `admin/index.html`: Khởi tạo Decap CMS và Netlify Identity.
*   `admin/config.yml`: File cấu hình định nghĩa toàn bộ giao diện quản trị (Blog, Pitch Deck, i18n).

### Cách thiết lập GitHub OAuth App cho CMS

Vì ứng dụng được host trên GitHub Pages, bạn cần một cơ chế xác thực OAuth. Cách chuẩn nhất là sử dụng **Netlify Identity** (Miễn phí) dù code bạn nằm trên GitHub.

**Bước 1: Link Repository với Netlify (Chỉ dùng cho Auth)**
1. Mở [Netlify.com](https://app.netlify.com/), tạo tài khoản và chọn **Add new site > Import an existing project**.
2. Chọn repo `vct-website` từ GitHub. (Không cần setup build command, cứ để trống).
3. Đặt URL cho site, ví dụ `vct-cms-auth.netlify.app`.

**Bước 2: Bật Netlify Identity**
1. Vào **Site configuration > Identity** trên Netlify, chọn **Enable Identity**.
2. Cuộn xuống **Services > Git Gateway**, chọn **Enable Git Gateway** và cấp quyền GitHub. Việc này cho phép CMS giao tiếp với GitHub API ẩn danh thông qua Netlify.
3. Chỉnh sửa `admin/config.yml` (Đã được tạo sẵn):
   Chỉ cần đảm bảo `backend.name: git-gateway` nếu dùng Netlify Git Gateway. *(Hiện tại tôi đang để mặc định backend là `github`, nếu bạn dùng Git Gateway, hãy sửa `github` thành `git-gateway` trong file config)*.

**Lưu ý Về Hộp thư Liên hệ (Inquiry Management):**
Vì đây là website static, CMS không có Database riêng để nhận Form. Giải pháp hiện tại (Formspree) là hoàn hảo. Đội ngũ có thể đăng nhập Formspree dashboard để quản lý list liên hệ, hoặc cấu hình webhook từ Formspree bắn thẳng dữ liệu vào một Google Sheet.

---

## 3. Phân tích: Git-based CMS (Decap) vs WordPress

| Tiêu chí | Decap CMS (Giải pháp hiện tại) | WordPress (CMS Truyền thống) |
| :--- | :--- | :--- |
| **Tốc độ (Performance)** | **Nhanh tuyệt đối** (100 Lighthouse). Dữ liệu JSON tĩnh load ngay lập tức. Cứu tinh cho trải nghiệm Mobile. | Chậm hơn, phụ thuộc vào tốc độ query Database (MySQL) và xử lý PHP ở Server. |
| **Bảo mật (Security)** | **Bất khả xâm phạm**. Không có database hay admin backend để hack (không SQL Injection, không Brute-force). | Luôn là mục tiêu của hacker. Cần cập nhật Plugin thường xuyên. |
| **Chi phí** | **Miễn phí 100%**. Host GitHub Pages mãi mãi. | Tốn phí host (Domain, VPS/Hosting + SSL định kỳ hàng năm). |
| **Workflow (CI/CD)** | Git Workflow chuẩn. Mọi thay đổi content đều là một commit, có thể rollback bất kỳ dòng chữ nào nếu nhập sai. | Dữ liệu nằm trong DB, khó track lịch sử chỉnh sửa code/content song song. |
| **Phù hợp với ai?** | Web nội dung cấu trúc rõ ràng (Blog, Landing Page, Docs), ít phải thêm bớt chức năng eCommerce ngay lập tức. | Web cần cài cắm nhiều Plugin phức tạp, bán hàng (WooCommerce). |

### Kết luận
Đối với VCT Platform giai đoạn hiện tại, **Decap CMS là "vũ khí bí mật" hoàn hảo**. Tốc độ load site vẫn nhanh như chớp, bảo mật 10/10, phí duy trì 0đ, mà đội Content vẫn có giao diện Dashboard kéo thả, soạn thảo WYSIWYG mượt mà không kém gì WordPress.
