# UI Refresh + P0/P1 Tracking

Mục tiêu: sửa lỗi logic phím (P0), thêm viewport/chuột/help (P1), làm mới giao diện, và cho **mọi thứ nhìn thấy đều bấm được** — click vào hàng list, tab, hay gợi ý phím ở footer đều chạy đúng hành động như phím tương ứng.

Trạng thái: `[ ]` chưa làm · `[~]` đang làm · `[x]` xong · `[!]` bỏ qua có lý do

## P0 — Logic phím và dữ liệu

- [x] Manga giữ màn quay lại (`returnScreen`/`returnIndex`) khi thoát chapters, không mất kết quả tìm kiếm
- [x] `pgdown`/`pgup` clamp đúng khi danh sách rỗng (anime + manga)
- [x] Settings lưu xong cập nhật `parent.cfg` + `parent.menuItems`, menu đổi ngôn ngữ ngay
- [x] Xoá nhánh reload config chết trong `MainModel.Update`
- [x] `q` hoạt động như `esc` ở mọi màn danh sách; chỉ thoát app ở menu chính
- [x] Giảm debounce play 2s → 300ms; bỏ chặn enter khi mở chapter manga
- [x] Confirm `y/n` trước khi xoá 1 mục (`d`) và xoá sạch (`c`) ở history/favorites
- [x] History entry manga mở bằng image viewer, không đưa URL chapter vào mpv/VLC
- [x] Manga reader thêm `n`/`p` chuyển chapter (đúng chiều: `n` tới, `p` lùi)

## P1 — Viewport, chuột, help

- [x] `windowedRange` cho manga trending/chapters và anime episodes; chỉ báo `start-end/total`
- [x] Mouse wheel cuộn được ở mọi danh sách (menu, anime, manga, history, favorites, settings)
- [x] Click hàng list = chọn hàng đó; click hàng đang chọn = enter
- [x] Click tab (history/favorites) để đổi tab
- [x] Click menu chính để mở mục
- [x] Click nửa trái/phải màn đọc manga = trang trước/sau
- [x] Footer hint bấm được: mỗi hint mang `trigger` key, click phát lại đúng `tea.KeyMsg`
- [x] Overlay help `?`, đóng bằng `esc`/`?`
- [x] Footer tự sinh từ danh sách hint theo màn, không viết tay chuỗi

## UI refresh

- [x] Header thống nhất: icon + tên module + badge (nguồn, trang, module cha) qua `drawHeader`
- [x] List dạng card, hàng chọn tô nền hết chiều rộng ở anime/manga (`listRowWidth`)
- [x] Footer dạng chip phím, divider full width, có phân trang bên phải
- [x] Empty state / loading / error dùng chung style (`drawLoadingLine`, `drawErrorLine`, `theme.Muted`)
- [x] Responsive: ẩn poster khi hẹp (<80 cột), list full width khi rộng
- [x] Settings hiển thị value + mô tả, không chỉ tên field

## Kiến trúc chuột

- `uiState` (pointer, chia sẻ qua `MainModel.ui`): lưu dòng bắt đầu list, hitbox hint/tab.
- Mỗi model ghi hitbox trong `View` (con trỏ dùng chung nên không mất khi receiver là value).
- `Update` nhận `tea.MouseMsg`, tra registry, phát `tea.KeyMsg` tương ứng hoặc đổi index trực tiếp.
- Không thêm dependency mới.

## Verify

- [x] `go build ./...`, `go vet ./...`, `gofmt -l .`
- [x] `go test ./...` pass
- [x] Test mới: `ui_mouse_test.go` (hint hitbox, row/grid mapping, clamp), `manga_test.go` (n/p, esc origin, pgdown rỗng), `view_smoke_test.go` (render 40/60/80/120 cột, help, menu click)
- [!] Smoke `go run ./cmd/gonime` không chạy headless: cần terminal thật. Thay bằng test render + click; kiểm tra tay khi có TTY.

## Ghi chú

- Terminal chỉ gửi mouse event khi app bật `WithMouseCellMotion` (đã bật).
- Click-to-key phụ thuộc terminal report toạ độ; không hỗ trợ terminal tắt mouse.
- Giữ nguyên palette neon hiện tại, chỉ hệ thống lại hierarchy và spacing.
