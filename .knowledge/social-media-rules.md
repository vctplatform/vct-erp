# Social Media Publishing Rules (Luật Đăng bài Bắt buộc)

Khi Chairman (Người dùng) yêu cầu "đăng bài lên mạng xã hội", các Agent Marketing BẮT BUỘC phải thực hiện đăng bài đồng loạt lên đủ 3 nền tảng hệ sinh thái của VCT Platform:
1. **Facebook Fanpage**
2. **Instagram**
3. **Threads**

## Nguyên tắc viết bản thảo (Drafting Rules):
- **Facebook & Instagram:** Sử dụng Bản đầy đủ (Long-form), luôn đính kèm Hình ảnh.
- **Quy tắc Thiết kế Ảnh (Dành cho Design/Media Agent): BẤT CỨ VĂN BẢN (TEXT) NÀO XUẤT HIỆN TRÊN HÌNH ẢNH BAO GỒM GIAO DIỆN APP ĐỀU PHẢI LÀ TIẾNG VIỆT (ví dụ: "LỊCH THI ĐẤU", "VÕ TRƯỜNG", "ĐĂNG KÝ"). Tuyệt đối cấm dùng tiếng Anh, ký tự vô nghĩa (gibberish) hoặc UI text mặc định. Hình ảnh UI phải hiển thị chuẩn Dark-mode theo ngôn ngữ tiếng Việt.
- **Threads (Cảnh báo Chí mạng):** BẮT BUỘC phải viết bản tóm tắt (Short-form). **Giới hạn dung lượng chữ tối đa là 500 ký tự**. Nếu văn bản dài hơn 500 ký tự, Meta Graph API sẽ ném lỗi 500 (`Param text must be at most 500 chars`) và task sẽ thất bại. Khuyến cáo độ dài lý tưởng cho Threads là ~350 ký tự.

## Format bắt buộc:
1. Phải có Call-to-Action (Trải nghiệm/Sử dụng/Đăng nhập).
2. Tên miền truy cập phần mềm duy nhất: `https://app.vctplatform.vn`
3. Hashtag tiêu chuẩn: `#VCTPlatform #VoCoTruyen #VoThuatVietNam #DigitalTransformation`

## Lệnh tự động hóa hệ thống (Dành cho Dev/Agent):
Sử dụng script `publish_post.js` kết hợp file `config.json` để bắn lệnh API. Tuy nhiên, vì Threads yêu cầu nội dung ngắn dưới 500 chữ và FB/IG yêu cầu nội dung dài, các Agent nên tách làm 2 file config riêng biệt:
- File Long-form: Cấu hình `platforms: ["facebook", "instagram"]`
- File Short-form: Cấu hình `platforms: ["threads"]`
