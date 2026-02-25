# 🎌 GoNime

**Xem Anime, Phim & Đọc Manga ngay trên Terminal — Viết bằng Go**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue?style=flat-square)]()
[![Build](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square)]()

```
    ██████╗   ██████╗  ███╗   ██╗ ██╗ ███╗   ███╗ ███████╗
   ██╔════╝  ██╔═══██╗ ████╗  ██║ ██║ ████╗ ████║ ██╔════╝
 ██║  ███╗ ██║   ██║ ██╔██╗ ██║ ██║ ██╔████╔██║ █████╗
 ██║   ██║ ██║   ██║ ██║╚██╗██║ ██║ ██║╚██╔╝██║ ██╔══╝
   ╚██████╔╝ ╚██████╔╝ ██║ ╚████║ ██║ ██║ ╚═╝ ██║ ███████╗
    ╚═════╝   ╚═════╝  ╚═╝  ╚═══╝ ╚═╝ ╚═╝     ╚═╝ ╚══════╝
```

---

## ✨ Tính năng

| Tính năng | Chi tiết |
| ---------- | --------- |
| 🎬 **Xem Phim/Anime** | Tìm kiếm & phát từ 3 nguồn: PhimAPI, OPhim, NguonC |
| 📚 **Đọc Manga** | Tìm kiếm & đọc từ 2 nguồn: TruyenVN, TruyenQQ |
| 🖼 **Thumbnail Rendering** | Hiển thị thumbnail trong CLI bằng chafa |
| 📖 **Full Manga Viewer** | Mở manga trong cửa sổ riêng với Wails |
| 🔖 **Chapter Navigation** | Chuyển chapter dễ dàng trong viewer |
| 📜 **Lịch sử** | Tự động lưu lịch sử xem/đọc |
| ⭐ **Yêu thích** | Thêm/xóa phim & manga yêu thích |
| ⚙️ **Cài đặt** | Chỉnh sửa player, thư mục cache trong TUI |
| 🎮 **TUI Hiện đại** | Giao diện Modern Dark với BubbleTea + Lipgloss |
| 💻 **CLI Commands** | `gonime anime`, `gonime manga`, `gonime history` |

---

## 🚀 Cài đặt

### Yêu cầu

| Công cụ | Dùng cho | Cài đặt |
| ------- | -------- | ------- |
| **mpv** hoặc **VLC** | Phát video (bắt buộc) | `winget install mpv` hoặc `winget install VideoLAN.VLC` |
| **chafa** | Hiển thị thumbnail (khuyến nghị) | `winget install hpjansson.Chafa` |
| **yt-dlp** | Tải video (tùy chọn) | `winget install yt-dlp` |

> **Windows Terminal** hoặc terminal hỗ trợ Unicode/ANSI để có trải nghiệm tốt nhất.

### Build từ source

```bash
# Clone repo
git clone https://github.com/dinhtuandev/GoNime.git
cd GoNime

# Build
go build -ldflags="-s -w" -o gonime.exe ./cmd/gonime

# Chạy
./gonime.exe
```

---

## 🎮 Sử dụng

### Menu chính

```
╔══════════════════════════════════════════════════════╗
║  🎌 GoNime                                        ║
║                                                 ║
║  ▶ 🎬 Xem Phim/Anime                            ║
║    📚 Đọc Manga                                  ║
║    📜 Lịch sử                                    ║
║    ⭐ Yêu thích                                   ║
║    ⚙️ Cài đặt                                     ║
║    🚪 Thoát                                      ║
║                                                 ║
║  ↑/↓: Điều hướng  Enter: Chọn  q: Thoát       ║
╚══════════════════════════════════════════════════════╝
```

### CLI Commands

```bash
# Xem version
gonime version

# Xem phim/anime (mở TUI)
gonime anime "one piece"

# Đọc manga (mở TUI)
gonime manga "naruto"

# Xem lịch sử
gonime history list

# Phát lại mục gần nhất
gonime history play

# Xóa toàn bộ lịch sử
gonime history clear

# Xem/chỉnh cài đặt
gonime settings
```

---

## ⌨️ Phím tắt

### Xem Phim/Anime

| Phím | Tác dụng |
| ---- | -------- |
| `m` | Đổi nguồn (PhimAPI → OPhim → NguonC) |
| `s` | Chuyển sang màn hình tìm kiếm |
| Gõ ký tự | Nhập tên phim |
| `Enter` | Tìm kiếm / Chọn phim / Phát tập |
| `n` | Tập tiếp theo |
| `p` | Tập trước |
| `f` | Thêm vào yêu thích |
| `Esc` | Quay lại |

### Đọc Manga

| Phím | Tác dụng |
| ---- | -------- |
| `m` | Đổi nguồn (TruyenVN → TruyenQQ) |
| `s` | Chuyển sang màn hình tìm kiếm |
| Gõ ký tự | Nhập tên manga |
| `Enter` | Tìm kiếm / Chọn manga / Đọc chapter |
| `n` | Chapter tiếp theo |
| `p` | Chapter trước |
| `v` | Mở trong Wails viewer (nếu có) |
| `r` | Reload/re-render ảnh |
| `f` | Thêm vào yêu thích |
| `Esc` | Quay lại |

---

## 🎨 Giao diện

### Modern Dark Theme

```
╔═══════════════════════════════════════════════════════════════╗
║  🎬 ANIME/MOVIE              [phimapi]        [🔥 Xu hướng] ║
╠═══════════════════════════════════════════════════════════════╣
║  ▶ Naruto Shippuden                      (2007)             ║
║    One Piece                                            (1999)║
║    Dragon Ball Super                                     (2015)║
║    Attack on Titan                                       (2013)║
║    My Hero Academia                                      (2016)║
║                                                           ║
║  m: Đổi nguồn │ s: Tìm kiếm │ Enter: Xem │ ↑↓: Di chuyển ║
╚═══════════════════════════════════════════════════════════════╝
```

---

## ⚙️ Cấu hình

File cấu hình: `%APPDATA%\gonime\config.json` (Windows) hoặc `~/.config/gonime/config.json` (Linux/macOS)

```json
{
  "player": "mpv",
  "download_dir": "~/Downloads/gonime",
  "manga_image_dir": "~/Downloads/gonime/manga",
  "manga_image_size": 100
}
```

---

## 🔗 Nguồn dữ liệu

### Phim/Anime

| Nguồn | Loại | API/URL |
| ------ | ---- | -------- |
| [phimapi.com](https://phimapi.com) | Anime/Phim | API chính |
| [ophim.cc](https://ophim.cc) | Phim | API dự phòng |
| [phim.nguonc.com](https://phim.nguonc.com) | Phim | Scraper |

### Manga

| Nguồn | Loại | URL |
| ------ | ---- | --- |
| [truyenvn.shop](https://truyenvn.shop) | Manga | Web scraping |
| [truyenqqno.com](https://truyenqqno.com) | Manga | Web scraping |

---

## 🏗 Cấu trúc dự án

```
gonime/
├── cmd/
│   └── gonime/
│       └── main.go              # Entry point, CLI commands
├── internal/
│   ├── api/
│   │   ├── client.go            # PhimAPI client
│   │   ├── ophim.go            # OPhim client
│   │   ├── nguonc.go           # NguonC scraper
│   │   └── movie.go            # Movie source interface
│   ├── config/
│   │   └── config.go           # Config, history, favorites
│   ├── player/
│   │   └── player.go           # mpv/vlc launcher
│   ├── render/
│   │   └── chafa.go            # Chafa image renderer
│   ├── scraper/
│   │   ├── manga.go            # TruyenVN scraper
│   │   ├── truyenqq.go         # TruyenQQ scraper
│   │   └── sources.go          # Manga source interface
│   ├── viewer/
│   │   ├── viewer.go           # Viewer abstraction
│   │   ├── wails_viewer.go     # Wails viewer
│   │   └── frontend/           # Wails frontend files
│   └── tui/
│       ├── main.go             # Main menu TUI
│       ├── anime.go            # Anime/Movie TUI
│       ├── manga.go            # Manga reader TUI
│       ├── history.go          # History TUI
│       ├── favorites.go         # Favorites TUI
│       └── settings.go         # Settings TUI
├── go.mod
├── go.sum
└── README.md
```

---

## 🗺 Roadmap

- [x] TUI chính với BubbleTea
- [x] Tìm kiếm & phát anime (PhimAPI)
- [x] Đọc manga với chafa rendering
- [x] Multiple movie sources (PhimAPI, OPhim, NguonC)
- [x] Multiple manga sources (TruyenVN, TruyenQQ)
- [x] Wails full manga viewer
- [x] Chapter navigation trong viewer
- [x] Thumbnail rendering trong CLI
- [x] Lịch sử xem/đọc
- [x] Danh sách yêu thích
- [x] Cài đặt có thể chỉnh sửa
- [x] CLI subcommands
- [ ] Tải episode về máy (yt-dlp)
- [ ] Subtitles support
- [ ] Prefetch ảnh manga
- [ ] Config profile (nhiều người dùng)

---

## 🤝 Đóng góp

PR và Issues luôn được chào đón!

```bash
git clone https://github.com/dinhtuandev/GoNime.git
cd GoNime
go mod tidy
go run ./cmd/gonime
```

---

## 📄 License

MIT License — xem [LICENSE](LICENSE)

---

<div align="center">

Made with ❤️ and Go

</div>
