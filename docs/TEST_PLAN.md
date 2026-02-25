# GoNime — Kế hoạch Kiểm thử (Test Plan)

Tài liệu này hướng dẫn cách kiểm thử toàn bộ chức năng của dự án GoNime, bao gồm kiểm thử tự động (Automated Tests) và kiểm thử thủ công (Manual Testing).

## 1. Kiểm thử tự động (Unit Tests)

Chạy các lệnh sau trong terminal để kiểm tra logic xử lý dữ liệu:

### Kiểm tra Cấu hình (Config)

Kiểm tra xem ứng dụng có lưu/load thiết lập, lịch sử và yêu thích đúng không.

```powershell
go test ./internal/config/... -v
```

### Kiểm tra Trình phát (Player)

Kiểm tra logic tìm kiếm đường dẫn mpv/vlc trên Windows.

```powershell
go test ./internal/player/... -v
```

### Kiểm tra API & Scraper

Kiểm tra kết nối và khả năng bóc tách dữ liệu từ web.

```powershell
go test ./internal/api/... -v
go test ./internal/scraper/... -v
```

---

## 2. Kiểm thử thủ công (Manual TUI Testing)

Vì đây là ứng dụng giao diện terminal (TUI), việc kiểm thử bằng tay là quan trọng nhất để đảm bảo trải nghiệm người dùng.

### Case 1: Khởi động & Trending

- [x] Chạy `.\gonime.exe`.
- [] Chọn **Xem Anime** → Phải hiển thị danh sách "Xu hướng" (Trending).
- [x] Nhấn `Esc` quay lại menu chính.
- [ ] Chọn **Đọc Manga** → Phải hiển thị danh sách truyện mới nhất từ web.

### Case 2: Tìm kiếm & Giao diện (UI)

- [x] Tại màn hình Anime/Manga, nhấn phím `s`.
- [x] Kiểm tra **Khung nhập liệu**: Viền phải khít, label "Nhập tên..." nằm trên khung, có con trỏ `▌`.
- [x] Nhập "Naruto" và nhấn `Enter`.
- [x] Kết quả tìm kiếm phải hiển thị danh sách phim/truyện đúng từ khóa.
- [x] Nhấn `Ctrl+U` để xóa nhanh dòng tìm kiếm.

### Case 3: Phát Video (Anime)

- [x] Chọn một bộ phim → Chọn một tập phim bất kỳ.
- [x] Nhấn `Enter` để phát.
- [x] **Kết quả**:
  - [x] Nếu máy có mpv/vlc: Trình phát phải tự bật lên.
  - [x] Nếu máy không có: Trình duyệt web mặc định phải mở Link M3U8/Embed.
- [x] Thử phím `n` (Next) và `p` (Previous) để chuyển tập.
- [x] Thử phím `s` để đổi Server (nếu phim có nhiều server).

### Case 4: Đọc Truyện (Manga)

- [x] Chọn một bộ truyện → Chọn một Chapter.
- [x] Nhấn `Enter`.
- [ ] **Kết quả**:
  - [ ] Nếu có `chafa`: Ảnh phải hiển thị dưới dạng ASCII/Sixel trong terminal. (có lẽ không đúng ý tưởng, thay vì vậy chỉ cần load ảnh lên terminal)
  - [x] Nhấn `o` (Open): Ảnh phải được tải về và mở bằng trình xem ảnh của Windows.
  - [x] Nhấn `w` (Web): Mở link chapter đó trên trình duyệt.

### Case 5: Cài đặt (Settings)

- [ ] Vào mục **Cài đặt**.
- [ ] Kiểm tra phần **🔍 KIỂM TRA HỆ THỐNG**:
  - Phải hiện dấu `✓` cho các công cụ đã cài (VLC, mpv, chafa).
- [ ] Thử đổi **URL nguồn Manga**: Thay bằng một URL khác và lưu lại.
- [ ] Thử nhấn `r` để Reset toàn bộ cài đặt về mặc định.

### Case 6: Lịch sử & Yêu thích

- [ ] Xem một phim hoặc đọc một truyện bất kỳ.
- [ ] Quay lại menu chính → Chọn **Lịch sử**: Phải thấy mục đã xem ở đầu danh sách.
- [ ] Tại màn hình danh sách tập/chapter, nhấn `f` để yêu thích.
- [ ] Quay lại menu chính → Chọn **Yêu thích**: Phải thấy mục đó đã xuất hiện.

---

## 3. Lệnh Build nhanh để Test

Nếu bạn vừa sửa code nhanh, hãy dùng lệnh này để build lại và chạy ngay:

```powershell
go build -o gonime.exe ./cmd/gonime/ ; .\gonime.exe
```
