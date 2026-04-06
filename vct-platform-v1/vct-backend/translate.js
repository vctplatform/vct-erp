const fs = require('fs');

const replacements = {
  "Giai Vo Co Truyen Toan Quoc 2026": "Giải Võ Cổ Truyền Toàn Quốc 2026",
  "Nha thi dau Quy Nhon": "Nhà thi đấu Quy Nhơn",
  "Doan Binh Dinh": "Đoàn Bình Định",
  "Binh Dinh": "Bình Định",
  "Nguyen Van Trong": "Nguyễn Văn Trọng",
  "Doan Ha Noi": "Đoàn Hà Nội",
  "Ha Noi": "Hà Nội",
  "Tran Tuan Anh": "Trần Tuấn Anh",
  "Pham Hoang Nam": "Phạm Hoàng Nam",
  "Le Thu Huong": "Lê Thu Hương",
  "Thang diem": "Thắng điểm",
  "San 1": "Sân 1",
  "Nha thi dau A": "Nhà thi đấu A",
  "San 2": "Sân 2",
  "Dang Quoc Minh": "Đặng Quốc Minh",
  "Vo Hai Yen": "Võ Hải Yến",
  "De nghi xem lai diem ky thuat": "Đề nghị xem lại điểm kỹ thuật",
  "Ngoc Tran Quyen": "Ngọc Trản Quyền",
  "Thanh nien": "Thanh niên",
  "Ban To Chuc": "Ban Tổ Chức",
  "Truong Doan Demo": "Trưởng Đoàn Demo",
  "Lien Doan VCT": "Liên Đoàn VCT",
  "Noi dung quyen nu": "Nội dung quyền nữ",
  "Noi dung doi khang nam": "Nội dung đối kháng nam",
  "Van dong vien doi khang": "Vận động viên đối kháng",
  "Van dong vien quyen": "Vận động viên quyền",
  "Dang ky hop le": "Đăng ký hợp lệ",
  "Lich thi dau sang ngay khai mac": "Lịch thi đấu sáng ngày khai mạc",
  "Can bo y te": "Cán bộ y tế",
  "Khoi tao du lieu demo": "Khởi tạo dữ liệu demo",
  "Du lieu demo da duoc nap thanh cong.": "Dữ liệu demo đã được nạp thành công."
};

const files = [
  'db/seeds/dev/0001_seed_entity_records.sql',
  'db/seeds/dev/0002_seed_relational_core.sql'
];

files.forEach(file => {
  if (fs.existsSync(file)) {
    let content = fs.readFileSync(file, 'utf8');
    for (const [key, value] of Object.entries(replacements)) {
      content = content.replace(new RegExp(key, 'g'), value);
    }
    fs.writeFileSync(file, content);
    console.log(`Updated ${file}`);
  } else {
    console.log(`File not found: ${file}`);
  }
});
