# 🎌 GoNime

**Xem Anime, Phim & Đọc Manga ngay trên Terminal — Viết bằng Go**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue?style=flat-square)]()

```
    ██████╗   ██████╗  ███╗   ██╗ ██╗ ███╗   ███╗ ███████╗
   ██╔════╝  ██╔═══██╗ ████╗  ██║ ██║ ████╗ ████║ ██╔════╝
 ██║  ███╗ ██║   ██║ ██╔██╗ ██║ ██║ ██╔████╔██║ █████╗
 ██║   ██║ ██║   ██║ ██║╚██╗██║ ██║ ██║╚██╔╝██║ ██╔══╝
   ╚██████╔╝ ╚██████╔╝ ██║ ╚████║ ██║ ██║ ╚═╝ ██║ ███████╗
    ╚═════╝   ╚═════╝  ╚═╝  ╚═══╝ ╚═╝ ╚═╝     ╚═╝ ╚══════╝
```

---

## ✨ Tính năng chính

| **Tính năng** | **Chi tiết** |
| ------------- | ------------ |
| 🎬 **Xem Phim/Anime** | Tìm kiếm & phát từ 3 nguồn: **PhimAPI**, **OPhim**, **NguonC** (có hỗ trợ TMDB ID từ PhimAPI). |
| 📚 **Đọc Manga** | Tìm kiếm & đọc từ **TruyenVN** và **OTruyen** (HTML + API CDN chapter). |
| 🖼 **Thumbnail trong TUI** | Dùng **chafa** để hiển thị poster/thumbnail cho Episodes/Chapters, History, Favorites. |
| 📖 **Manga Viewer riêng** | Mở manga trong cửa sổ riêng (Wails + frontend web). |
| 🔖 **Chapter / Episode Navigation** | Chuyển chapter/tập bằng phím tắt (n/p, g/G, PgUp/Dn, …). |
| 📜 **Lịch sử** | Tự động lưu lịch sử xem anime/manga. |
| ⭐ **Yêu thích** | Thêm/xóa phim & manga yêu thích, có poster preview. |
| ⬇ **Download episode** | Tải tập hiện tại bằng **yt-dlp** (phím `d` trong màn Episodes). |
| ⚙️ **Cài đặt trong TUI** | Chỉnh sửa player, thư mục download/cache trực tiếp từ Settings. |
| 🎮 **Giao diện Modern Dark** | BubbleTea + Lipgloss, hỗ trợ fuzzy search cho Anime/Manga. |

---

## 🚀 Cài đặt

### Yêu cầu

| **Công cụ** | **Dùng cho** | **Cài đặt gợi ý (Windows)** |
| ----------- | ------------ | ---------------------------- |
| **Go 1.24+** | Build & chạy GoNime | `winget install GoLang.Go` |
| **mpv** hoặc **VLC** | Phát video (bắt buộc) | `winget install shinchiro.mpv` hoặc `winget install VideoLAN.VLC` |
| **chafa** | Hiển thị thumbnail trong terminal (khuyến nghị) | `winget install hpjansson.Chafa` |
| **yt-dlp** | Tải episode về máy (tùy chọn, khuyến nghị) | `winget install yt-dlp` |

> **Khuyến nghị:** Dùng Windows Terminal hoặc terminal hỗ trợ Unicode + ANSI để có trải nghiệm tốt nhất.

### Build từ source

```bash
# Clone repo
git clone https://github.com/dinhtuandev/GoNime.git
cd GoNime

# Cài dependencies (nếu cần)
go mod tidy

# Build (Windows)
go build -ldflags="-s -w" -o gonime.exe ./cmd/gonime

# Chạy
./gonime.exe
```

---

## 🎮 Sử dụng nhanh

### Menu chính

```
╔══════════════════════════════════════════════════════╗
║  🎌 GoNime                                           ║
║                                                      ║
║  ▶ 🎬 Xem Phim/Anime                                 ║
║    📚 Đọc Manga                                      ║
║    📜 Lịch sử                                        ║
║    ⭐ Yêu thích                                       ║
║    ⚙️ Cài đặt                                        ║
║    🚪 Thoát                                           ║
║                                                      ║
║  ↑/↓: Điều hướng  Enter: Chọn  q: Thoát             ║
╚══════════════════════════════════════════════════════╝
```

### CLI Commands

```bash
# Xem version
gonime version

# Mở TUI xem phim/anime (vào thẳng màn Anime)
gonime anime

# Mở TUI đọc manga (vào thẳng màn Manga)
gonime manga

# Xem lịch sử
gonime history

# Xem / chỉnh cài đặt
gonime settings
```

---

## ⌨️ Phím tắt chi tiết

### Xem Phim/Anime (Anime TUI)

**Trending & Search**

| **Phím** | **Tác dụng** |
| -------- | ------------ |
| `m` | Đổi nguồn phim (`phimapi` → `ophim` → `nguonc`). |
| `s` | Chuyển sang màn hình tìm kiếm. |
| Gõ ký tự | Nhập tên phim/anime trong ô search. |
| `Enter` | Tìm kiếm (khi đang nhập) / Chọn phim / Vào màn tập. |
| `↑/↓` hoặc `j/k` | Di chuyển trong danh sách. |
| `Esc` | Quay lại menu chính. |

**Episodes**

| **Phím** | **Tác dụng** |
| -------- | ------------ |
| `Enter` | Phát tập đang chọn. |
| `n` | Phát tập tiếp theo. |
| `p` | Phát tập trước đó. |
| `s` | Đổi server (nếu có nhiều server). |
| `g` / `G` | Nhảy về tập đầu/cuối. |
| `PgUp` / `PgDn` | Cuộn danh sách ±10 tập. |
| `f` | Thêm phim vào danh sách yêu thích. |
| `d` | **Tải tập hiện tại bằng `yt-dlp` về `download_dir`.** |
| `Esc` | Quay lại màn Trending. |

> **Downloader:** Phím `d` yêu cầu `yt-dlp` có trong PATH và `download_dir` đã cấu hình hợp lệ.

---

### Đọc Manga (Manga TUI)

**Trending & Search**

| **Phím** | **Tác dụng** |
| -------- | ------------ |
| `m` | Đổi nguồn manga (`otruyen` ↔ `truyenvn`). |
| `s` | Chuyển sang màn hình tìm kiếm manga. |
| Gõ ký tự | Nhập tên manga. |
| `Enter` | Tìm kiếm / Chọn manga / Vào màn Chapters. |
| `↑/↓` hoặc `j/k` | Di chuyển trong danh sách. |
| `Esc` | Quay lại menu chính. |

**Chapters**

| **Phím** | **Tác dụng** |
| -------- | ------------ |
| `Enter` | Đọc chapter đang chọn. |
| `n` / `p` | Chapter tiếp theo / trước. |
| `PgUp` / `PgDn` | Cuộn danh sách ±10 chapter. |
| `g` / `G` | Nhảy về chapter đầu/cuối. |
| `f` | Thêm manga vào yêu thích. |
| `Esc` | Quay lại Trending. |

**Reading (đọc từng trang trong terminal)**  

| **Phím** | **Tác dụng** |
| -------- | ------------ |
| `←/→` hoặc `↑/↓` | Trang trước / sau. |
| `Enter` | Trang sau. |
| `n` / `p` | Chapter tiếp theo / trước. |
| `r` | Render lại ảnh hiện tại. |
| `o` | Mở bằng viewer ngoài (Wails viewer nếu có). |
| `w` | Mở chapter trong trình duyệt hệ thống. |
| `f` | Thêm vào yêu thích. |
| `Esc` | Quay lại danh sách chapters. |

---

## 🎨 Giao diện (Modern Dark)

Ví dụ màn Anime Trending:

```
╔═══════════════════════════════════════════════════════════════╗
║  🎬 ANIME / MOVIE                  [phimapi]   [🔥 Xu hướng]  ║
╠═══════════════════════════════════════════════════════════════╣
║  ▶ Vạt Áo Đồng Phục                            (2026)        ║
║    Siêu Nhân Điện Quang Geed                   (2017)        ║
║    Trót Yêu Quý Cô Lạnh Lùng                   (2026)        ║
║    ...                                                         ║
║                                                               ║
║  m: Đổi nguồn │ s: Tìm kiếm │ Enter: Xem │ ↑↓/jk: Di chuyển  ║
╚═══════════════════════════════════════════════════════════════╝
```

Màn Episodes, Chapters, History, Favorites đều có **poster panel bên phải** khi có thumbnail (từ PhimAPI/OPhim/NguonC/OTruyen/TruyenVN).

---

## ⚙️ Cấu hình

File cấu hình:

- **Windows:** `%APPDATA%\gonime\config.json`  
- **Linux/macOS:** `~/.config/gonime/config.json`

Ví dụ:

```json
{
  "player": "mpv",
  "download_dir": "C:/Users/you/Downloads/gonime",
  "manga_image_dir": "C:/Users/you/Downloads/gonime/manga",
  "manga_image_size": 100,
  "manga_url": "https://truyenvn.shop"
}
```

- **player**: tên binary player (`mpv`, `vlc`, ...).
- **download_dir**: thư mục để yt-dlp tải video.
- **manga_image_dir**: thư mục cache ảnh manga (cho chafa/viewer).
- **manga_image_size**: kích thước ảnh manga (50–200).
- **manga_url**: URL mặc định cho scraper TruyenVN (vẫn dùng nội bộ, nhưng nguồn chính đang là `otruyen`).

---

## 🔗 Nguồn dữ liệu

### Phim/Anime

| **Nguồn** | **Loại** | **Mô tả** |
| --------- | -------- | --------- |
| [`phimapi.com`](https://phimapi.com) | Anime/Phim | API chính; dùng `/v1/api/tim-kiem`, `/phim/{slug}`, `/danh-sach/phim-moi-cap-nhat`, `/tmdb/{type}/{id}`. |
| [`ophim.cc`](https://ophim.cc) | Phim | API dự phòng qua `OPhimClient`. |
| [`phim.nguonc.com`](https://phim.nguonc.com) | Phim | Scraper HTML qua `NguonCScraper`. |

### Manga

| **Nguồn** | **Loại** | **Mô tả** |
| --------- | -------- | --------- |
| [`truyenvn.shop`](https://truyenvn.shop) | Manga | Scraper HTML (Madara theme) — `MangaScraper`. |
| [`otruyen.cc`](https://otruyen.cc) | Manga | Scraper HTML + CDN chapter API (`sv*.otruyencdn.com/v1/api/chapter/...`) — `OTruyenScraper`. |

---

## 🏗 Cấu trúc dự án (rút gọn)

```text
gonime/
├── cmd/
│   └── gonime/
│       └── main.go              # Entry point, CLI commands (anime/manga/history/settings)
├── internal/
│   ├── api/
│   │   ├── client.go            # PhimAPI client (search, episodes, TMDB, trending)
│   │   ├── ophim.go             # OPhim client
│   │   ├── nguonc.go            # NguonC scraper
│   │   └── movie.go             # Movie source interfaces + registry
│   ├── config/
│   │   └── config.go            # Config, history, favorites (+ defaults)
│   ├── player/
│   │   └── player.go            # mpv/vlc launcher, yt-dlp Downloader wrapper
│   ├── render/
│   │   └── chafa.go             # Chafa-based image renderer & cache
│   ├── scraper/
│   │   ├── manga.go             # TruyenVN scraper
│   │   ├── otruyen.go           # OTruyen scraper (HTML + chapter CDN)
│   │   └── sources.go           # MangaSource interface + registry
│   ├── viewer/
│   │   ├── viewer.go            # Viewer abstraction
│   │   ├── wails_viewer.go      # Wails-based viewer implementation
│   │   └── frontend/            # Viewer frontend (HTML/JS/CSS)
│   └── tui/
│       ├── main.go              # Main menu TUI
│       ├── anime.go             # Anime/Movie TUI (trending/search/episodes)
│       ├── manga.go             # Manga TUI (trending/search/chapters/reading)
│       ├── history.go           # History TUI (poster + metadata)
│       ├── favorites.go         # Favorites TUI (poster + metadata)
│       ├── scraper_bridge.go    # Bridge cho manga favorites → scraper
│       └── settings.go          # Settings TUI
├── docs/                        # Kế hoạch test, pipeline, ideas
├── go.mod
├── go.sum
└── README.md
```

---

## 🗺 Roadmap

- [x] TUI chính với BubbleTea.
- [x] Tìm kiếm & phát anime (PhimAPI).
- [x] Multiple movie sources (PhimAPI, OPhim, NguonC).
- [x] Đọc manga với chafa rendering.
- [x] Multiple manga sources (TruyenVN, OTruyen).
- [x] Wails full manga viewer.
- [x] Chapter navigation trong viewer + TUI.
- [x] Thumbnail rendering trong CLI (episodes, chapters, history, favorites).
- [x] Lịch sử xem/đọc.
- [x] Danh sách yêu thích.
- [x] Cài đặt có thể chỉnh sửa trong TUI.
- [x] CLI subcommands cơ bản.
- [x] **Tải episode về máy bằng yt-dlp (phím `d`).**
- [ ] Subtitles support.
- [ ] Prefetch & cache ảnh manga nâng cao.
- [ ] Config profile (đa người dùng).

---

## 🤝 Đóng góp

PR và Issues luôn được chào đón!

```bash
git clone https://github.com/dinhtuandev/GoNime.git
cd GoNime
go mod tidy
go run ./cmd/gonime
```

Khi gửi PR, vui lòng:

- Giữ code style tương tự project (Go fmt, idiomatic Go).
- Chạy `go test ./...` trước khi push.
- Tránh commit file nhị phân lớn nếu không thực sự cần thiết.

---

## 📄 License

MIT License — xem [LICENSE](LICENSE)

---

<div align="center">

Made with ❤️ and Go

</div>
