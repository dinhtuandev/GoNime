<div align="center">

# 🎌 GoNime

**Xem Anime & Đọc Manga ngay trên Terminal — Viết bằng Go**

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

_Lấy cảm hứng từ [anisub](https://github.com/niyakipham/anisub) — nâng cấp lên Go native với TUI đẹp hơn, nhanh hơn, đa nền tảng_

</div>

---

## ✨ Tính năng

| Tính năng              | Chi tiết                                                                            |
| ---------------------- | ----------------------------------------------------------------------------------- |
| 📺 **Xem Anime**       | Tìm kiếm & phát qua `phimapi.com`, stream M3U8 trực tiếp bằng mpv/VLC               |
| 📚 **Đọc Manga**       | Tìm kiếm & đọc từ `truyenvn.shop`, render ảnh trực tiếp trong terminal bằng `chafa` |
| 🖼 **Chafa Rendering** | Hiển thị ảnh manga dưới dạng ASCII/Unicode art ngay trong terminal                  |
| 📜 **Lịch sử**         | Tự động lưu lịch sử xem/đọc, phát lại bằng một phím                                 |
| ⭐ **Yêu thích**       | Thêm/xóa anime & manga yêu thích, mở nhanh                                          |
| ⚙️ **Cài đặt**         | Chỉnh sửa player, thư mục tải, kích thước ảnh trực tiếp trong TUI                   |
| 🎮 **TUI đẹp**         | Giao diện terminal đầy đủ bằng BubbleTea + Lipgloss                                 |
| 🖥️ **CLI commands**    | `gonime anime`, `gonime manga`, `gonime history list/play/clear`                    |

---

## 🖼 Giao diện

```
╔══════════════════════════════════════════════════════╗
║  🎌 GoNime                                           ║
║                                                      ║
║  ▶ 📺 Xem Anime                                      ║
║    📚 Đọc Manga                                      ║
║    📜 Lịch sử                                        ║
║    ⭐ Yêu thích                                      ║
║    ⚙️  Cài đặt                                       ║
║    🚪 Thoát                                          ║
║                                                      ║
║  ↑/↓: Điều hướng  Enter: Chọn  q: Thoát             ║
╚══════════════════════════════════════════════════════╝
```

---

## 🚀 Cài đặt

### Yêu cầu

| Công cụ              | Dùng cho                                        | Cài đặt                                                 |
| -------------------- | ----------------------------------------------- | ------------------------------------------------------- |
| **mpv** hoặc **VLC** | Phát anime (bắt buộc)                           | `winget install mpv` hoặc `winget install VideoLAN.VLC` |
| **chafa**            | Hiển thị ảnh manga trong terminal (khuyến nghị) | `winget install hpjansson.Chafa`                        |
| **yt-dlp**           | Tải video (tùy chọn)                            | `winget install yt-dlp`                                 |

> **Windows Terminal** hoặc terminal hỗ trợ Unicode/ANSI để có trải nghiệm tốt nhất.

---

### Option 1: Tải binary sẵn (Windows)

1. Vào [Releases](https://github.com/dinhtuandev/GoNime/releases)
2. Tải `gonime_windows_amd64.exe`
3. Đặt vào thư mục trong `PATH` (ví dụ `C:\Tools\`)
4. Chạy `gonime`

### Option 2: Build từ source

```bash
# Clone repo
git clone https://github.com/dinhtuandev/GoNime.git
cd GoNime

# Build
go build -ldflags="-s -w" -o gonime.exe ./cmd/gonime

# Chạy
.\gonime.exe
```

---

## 🎮 Sử dụng

### Mở TUI chính

```bash
gonime
```

### CLI Commands

```bash
# Xem version
gonime version

# Tìm anime (mở TUI với search sẵn)
gonime anime "one piece"

# Xem lịch sử
gonime history list

# Phát lại mục gần nhất trong lịch sử
gonime history play

# Phát lại mục thứ 3 trong lịch sử
gonime history play 3

# Xóa toàn bộ lịch sử
gonime history clear

# Xem/chỉnh cài đặt
gonime settings
```

---

## ⌨️ Phím tắt

### Menu chính

| Phím                     | Tác dụng   |
| ------------------------ | ---------- |
| `↑` / `↓` hoặc `k` / `j` | Điều hướng |
| `Enter`                  | Chọn       |
| `q` hoặc `Ctrl+C`        | Thoát      |

### Xem Anime

| Phím     | Tác dụng                         |
| -------- | -------------------------------- |
| Gõ ký tự | Nhập tên anime                   |
| `Enter`  | Tìm kiếm / Chọn anime / Phát tập |
| `n`      | Tập tiếp theo                    |
| `p`      | Tập trước                        |
| `f`      | Thêm vào yêu thích               |
| `Esc`    | Quay lại                         |

### Đọc Manga

| Phím                  | Tác dụng                         |
| --------------------- | -------------------------------- |
| `↑` / `↓` / `←` / `→` | Chuyển trang                     |
| `Enter`               | Trang tiếp theo                  |
| `n`                   | Chapter tiếp theo                |
| `p`                   | Chapter trước                    |
| `r`                   | Reload/re-render ảnh             |
| `o`                   | Mở ảnh bằng Windows Photo Viewer |
| `w`                   | Mở chapter trong trình duyệt     |
| `f`                   | Thêm vào yêu thích               |
| `Esc`                 | Quay lại                         |

### Lịch sử & Yêu thích

| Phím    | Tác dụng          |
| ------- | ----------------- |
| `Enter` | Mở lại / phát lại |
| `d`     | Xóa mục đang chọn |
| `c`     | Xóa tất cả        |
| `Esc`   | Quay lại          |

### Cài đặt

| Phím    | Tác dụng                 |
| ------- | ------------------------ |
| `Enter` | Bắt đầu chỉnh sửa trường |
| `r`     | Đặt lại về mặc định      |
| `Esc`   | Hủy / Quay lại           |

---

## ⚙️ Cấu hình

File cấu hình được lưu tại:

- **Windows**: `%APPDATA%\gonime\config.json`
- **Linux/macOS**: `~/.config/gonime/config.json`

```json
{
  "player": "mpv",
  "download_dir": "C:\\Users\\<user>\\Downloads\\gonime",
  "manga_image_dir": "C:\\Users\\<user>\\Downloads\\gonime\\manga",
  "manga_image_size": 100
}
```

| Trường             | Mô tả                               | Mặc định                   |
| ------------------ | ----------------------------------- | -------------------------- |
| `player`           | Trình phát video (`mpv` hoặc `vlc`) | `mpv`                      |
| `download_dir`     | Thư mục tải video                   | `~/Downloads/gonime`       |
| `manga_image_dir`  | Thư mục cache ảnh manga             | `~/Downloads/gonime/manga` |
| `manga_image_size` | Kích thước ảnh chafa (%)            | `100`                      |

---

## 🏗 Cấu trúc dự án

```
gonime/
├── cmd/
│   └── gonime/
│       └── main.go          # Entry point, CLI commands
├── internal/
│   ├── api/
│   │   └── client.go        # Anime API (phimapi.com)
│   ├── config/
│   │   └── config.go        # Config, history, favorites
│   ├── player/
│   │   └── player.go        # mpv/vlc launcher, yt-dlp downloader
│   ├── render/
│   │   └── chafa.go         # Chafa image renderer
│   ├── scraper/
│   │   └── manga.go         # Manga scraper (truyenvn.shop)
│   └── tui/
│       ├── main.go          # Main menu TUI
│       ├── anime.go         # Anime search & playback TUI
│       ├── manga.go         # Manga reader TUI (với chafa)
│       ├── history.go       # History TUI
│       ├── favorites.go     # Favorites TUI
│       ├── settings.go      # Settings editor TUI
│       └── scraper_bridge.go # Interface bridge
├── anisub/
│   └── anisub.sh            # Shell script gốc (tham khảo)
├── go.mod
├── go.sum
└── README.md
```

---

## 🔧 Nguồn dữ liệu

| Nguồn                                  | Dùng cho            | API/URL                               |
| -------------------------------------- | ------------------- | ------------------------------------- |
| [phimapi.com](https://phimapi.com)     | Anime stream (M3U8) | `https://phimapi.com/v1/api/tim-kiem` |
| [truyenvn.shop](https://truyenvn.shop) | Manga images        | Web scraping với goquery              |

---

## 📦 Dependencies

| Package                                                               | Dùng cho      |
| --------------------------------------------------------------------- | ------------- |
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | TUI framework |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)   | TUI styling   |
| [spf13/cobra](https://github.com/spf13/cobra)                         | CLI framework |
| [PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)         | HTML scraping |

**External tools:**

- [mpv](https://mpv.io) hoặc [VLC](https://www.videolan.org) — phát video
- [chafa](https://hpjansson.org/chafa/) — render ảnh manga trong terminal
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) — tải video (tùy chọn)

---

## 🗺 Roadmap

- [x] TUI chính với BubbleTea
- [x] Tìm kiếm & phát anime qua phimapi.com
- [x] Đọc manga với chafa rendering
- [x] Lịch sử xem/đọc
- [x] Danh sách yêu thích
- [x] Cài đặt có thể chỉnh sửa
- [x] CLI subcommands
- [ ] Tải episode về máy (yt-dlp integration)
- [ ] Subtitles support
- [ ] Prefetch ảnh manga (tải trước trang tiếp theo)
- [ ] Multiple anime API sources
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

Made with ❤️ and Go | Lấy cảm hứng từ **anisub**

</div>
