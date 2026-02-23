# GoNime — Kế Hoạch Sửa Lỗi & Tính Năng Mới

## Mô tả vấn đề

Ba lỗi được báo cáo sau khi các phiên fix trước, cùng với một tính năng mới được đề xuất.

---

## Lỗi 1 — Thanh tìm kiếm bị lệch khung

**Nguyên nhân root cause:**  
`lipgloss.Border()` tự thêm padding + newline nội bộ không kiểm soát được. Khi kết hợp với label text bên ngoài, vị trí không thể căn chỉnh pixel-perfect qua chuỗi string thông thường.

**Fix:**  
Bỏ `inputStyle.Border()` — thay bằng tự vẽ box dùng ký tự Unicode box-drawing:

```
  Nhập tên anime:
  ╭──────────────────────────────────────────────────────╮
  │ bleach▌                                              │
  ╰──────────────────────────────────────────────────────╯
```

Dùng hàm `drawInputBox(label, value, width)` dùng chung cho cả anime.go và manga.go.

### Các file thay đổi

#### [MODIFY] [anime.go](file:///d:/WorkingSpace/gonime/internal/tui/anime.go)

- Xóa `inputBoxStyle` dùng lipgloss border
- Thêm hàm `drawInputBox()` render thủ công bằng `strings.Builder`
- Áp dụng cho `mangaScreenSearch` View

#### [MODIFY] [manga.go](file:///d:/WorkingSpace/gonime/internal/tui/manga.go)

- Tương tự anime.go

---

## Lỗi 2 — MPV không phát video

**Nguyên nhân root cause:**  
`winget install mpv` cài bản **Microsoft Store MSIX** (9P3JFR0CLLL6 bởi SilverEzhik), **không phải** CLI mpv. MSIX package không thêm `mpv.exe` vào PATH và nằm trong `C:\Users\...\AppData\Local\Packages\...` — một vị trí không được code hiện tại kiểm tra.

Thêm vào đó, `--no-terminal` flag khiến mpv chạy nhưng hoàn toàn im lặng nếu không có window GUI (MSIX mpv có thể không hỗ trợ flag này).

**Fix:**

1. **Cài đúng mpv CLI**: Hướng dẫn user chạy `winget install shinchiro.mpv` (portable CLI version)
2. **Mở rộng `FindAnyPlayer()`** để tìm thêm:
   - Scan `%LOCALAPPDATA%\Packages\` cho MSIX mpv
   - Scan `%LOCALAPPDATA%\Microsoft\WinGet\Links\` (WinGet symlinks)
   - Path đặc biệt: `%LOCALAPPDATA%\Microsoft\WindowsApps\mpv.exe`
3. **Bỏ `--no-terminal` flag** với mpv — thay bằng `--force-window=yes` để đảm bảo open GUI
4. **Fallback cuối**: nếu không tìm được player nào → dùng `cmd /c start "" <url>` (mở bằng default browser/app)
5. **Hiển thị thông báo lỗi rõ** trong TUI với hướng dẫn cài đặt

### Các file thay đổi

#### [MODIFY] [player.go](file:///d:/WorkingSpace/gonime/internal/player/player.go)

- Thêm WinGet symlinks path vào `FindAnyPlayer()`
- Thêm `%LOCALAPPDATA%\Microsoft\WindowsApps\` vào scan list
- Sửa mpv args: bỏ `--no-terminal`, thêm `--force-window=yes`
- Thêm `PlayURL()` fallback dùng Windows `start` command

---

## Lỗi 3 — Manga Search không hoạt động

**Nguyên nhân root cause:**  
CSS selector trong scraper đã lỗi thời. `truyenvn.shop` đã thay đổi HTML structure:

|                | Selector cũ                           | Hiện tại                           |
| -------------- | ------------------------------------- | ---------------------------------- |
| Search results | `.c-tabs-item__content .post-title a` | **Không tồn tại**                  |
| Manga links    | Dạng trên                             | `a[href*="/truyen-tranh/"][title]` |

Website hiện dùng class `manga-slider`, `widget-manga-search` — cấu trúc Madara theme mới.

**Fix:**

1. Cập nhật selector trong `Search()`:
   ```go
   // Thay vì: .c-tabs-item__content .post-title a
   // Dùng: div.page-listing-item .post-title a
   //    hoặc: div.c-image-inner + .post-title a
   //    hoặc: a[href*="/truyen-tranh/"] với title attr
   ```
2. Thêm **fallback selector chain**: thử nhiều selector theo thứ tự
3. Thêm **debug logging**: ghi ra số kết quả tìm được, báo lỗi cụ thể trong TUI

### Các file thay đổi

#### [MODIFY] [manga.go (scraper)](file:///d:/WorkingSpace/gonime/internal/scraper/manga.go)

- Cập nhật `Search()` với selector mới
- Thêm fallback selector chain
- Thêm thông báo lỗi mạng và "không tìm thấy kết quả" rõ ràng

---

## Tính năng: Trang Trending khi vào Anime/Manga

**Mô tả:**  
Thay vì hiện ô search trống, khi mở Anime/Manga hiển thị ngay danh sách **trending/nổi bật**. Phím `s` để chuyển sang chế độ tìm kiếm.

**Design:**

```
📺 ANIME                                        [🔥 Xu hướng]

  ▶ Naruto Shippuden
    One Piece
    Bleach: TYBW
    Attack on Titan Final Season
    ...

  s: Tìm kiếm  │  Enter: Xem  │  ↑↓: Điều hướng  │  Esc: Về menu
```

**Khi nhấn `s`:**

```
📺 ANIME                                        [🔍 Tìm kiếm]

  Nhập tên anime:
  ╭────────────────────────────────────────────────────╮
  │ ▌                                                  │
  ╰────────────────────────────────────────────────────╯

  Esc: Về xu hướng  │  Enter: Tìm kiếm  │  Ctrl+U: Xóa
```

**Nguồn data:**

- Anime: `https://phimapi.com/v1/api/danh-sach/phim-bo?page=1&sort_field=modified.time`
- Manga: Parse trang chủ `truyenvn.shop` — section manga nổi bật

### Các file thay đổi

#### [MODIFY] [anime.go](file:///d:/WorkingSpace/gonime/internal/tui/anime.go)

- Thêm `animeScreenTrending` screen state và `trendingMode bool`
- Thêm `fetchTrending()` command gọi API trending
- Phím `s` toggle sang search mode
- `Init()` trigger fetch trending ngay khi mở

#### [MODIFY] [manga.go](file:///d:/WorkingSpace/gonime/internal/tui/manga.go)

- Tương tự — fetch trending manga từ trang chủ truyenvn.shop
- Phím `s` toggle sang search mode

#### [MODIFY] [api/client.go](file:///d:/WorkingSpace/gonime/internal/api/client.go)

- Thêm `GetTrending(page int)` gọi danh sách phim bộ mới nhất

---

## Kế hoạch Kiểm tra (Verification)

### Build check (tự động)

```powershell
cd d:\WorkingSpace\gonime
go build ./... 2>&1; echo "BUILD: $LASTEXITCODE"
go vet ./...  2>&1; echo "VET:   $LASTEXITCODE"
```

### Test thủ công — Lỗi 1 (Input box)

1. Chạy `.\gonime.exe`
2. Chọn **Xem Anime** → kiểm tra label và border thẳng hàng
3. Chọn **Đọc Manga** → kiểm tra tương tự
4. Gõ một số chữ kể cả `n`, `p`, `f` → phải hiển thị đúng trong box

### Test thủ công — Lỗi 2 (MPV)

1. Chạy `.\gonime.exe anime`
2. Search "one piece" → chọn một bộ → chọn tập 1 → Enter
3. Kiểm tra mpv/VLC có mở không
4. Nếu không có player: phải hiện thông báo hướng dẫn cài đặt rõ ràng

### Test thủ công — Lỗi 3 (Manga search)

1. Chạy `.\gonime.exe manga`
2. Search "naruto" hoặc "one piece" → Enter
3. Phải hiển thị danh sách kết quả (không còn bị stuck ở hint)

### Test thủ công — Trending feature

1. Chạy `.\gonime.exe`
2. Vào **Xem Anime** → phải thấy danh sách trending ngay lập tức
3. Nhấn `s` → chuyển sang ô search
4. Vào **Đọc Manga** → tương tự

### Push GitHub

```powershell
git add -A
git commit -m "fix: input layout, mpv detection, manga scraper, trending feature"
git push origin main
```
