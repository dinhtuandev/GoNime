# GoNime CLI

<div align="center">

![GoNime](https://img.shields.io/badge/GoNime-CLI-blue?style=for-the-badge)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey?style=for-the-badge)

**Trình xem Anime và Manga thế hệ mới trên Terminal**

[Features](#-tính-năng) • [Cài đặt](#-cài-đặt) • [Sử dụng](#-sử-dụng) • [Đóng góp](#-đóng-góp)

</div>

---

## 📖 Giới thiệu

GoNime là một ứng dụng CLI native được viết bằng Go, cho phép bạn:
- 🎬 **Xem Anime** - Tìm kiếm và xem anime vietsub trực tiếp trên terminal
- 📚 **Đọc Manga** - Đọc manga từ các nguồn trực tuyến
- ⬇️ **Tải xuống** - Tải anime và manga về máy
- 📜 **Lịch sử** - Theo dõi những gì bạn đã xem
- ⭐ **Yêu thích** - Lưu danh sách anime/manga yêu thích

## ✨ Tính năng

| Tính năng | Mô tả |
|-----------|-------|
| 🔍 Tìm kiếm thông minh | Tìm kiếm anime và manga nhanh chóng |
| 🎬 Phát video chất lượng cao | Hỗ trợ nhiều độ phân giải |
| 📚 Đọc manga | Hiển thị ảnh manga ngay trong terminal |
| 💾 Tải xuống | Tải video và hình ảnh về máy |
| 📜 Lịch sử | Lưu lịch sử xem để tiếp tục sau |
| ⭐ Yêu thích | Đánh dấu anime/manga yêu thích |
| ⚙️ Cấu hình | Tùy chỉnh trình phát, thư mục tải |

## 📋 Yêu cầu hệ thống

- **Hệ điều hành**: Windows 10/11
- **Go**: 1.21 trở lên (chỉ cần khi build từ source)

### Phụ thuộc bên ngoài

| Công cụ | Mục đích | Bắt buộc |
|---------|----------|----------|
| `mpv` hoặc `vlc` | Phát video | ✅ Có |
| `yt-dlp` | Tải video | ❌ Tùy chọn |
| `chafa` | Hiển thị ảnh manga trong terminal | ❌ Tùy chọn |

## 🚀 Cài đặt

### Cách 1: Tải file thực thi (Khuyến nghị)

Tải file `gonime.exe` từ [Releases](https://github.com/dinhtuandev/gonime/releases) và đặt vào thư mục trong PATH.

### Cách 2: Build từ source

```bash
# Clone repository
git clone https://github.com/dinhtuandev/gonime.git
cd gonime

# Cài đặt dependencies
go mod download

# Build
go build -o gonime.exe ./cmd/gonime
```

### Cài đặt phụ thuộc

```powershell
# Cài đặt mpv (khuyến nghị)
winget install mpv

# Hoặc cài đặt VLC
winget install VideoLAN.VLC

# Cài đặt yt-dlp (cho tính năng tải)
winget install yt-dlp

# Cài đặt chafa (cho đọc manga)
# Tải từ: https://github.com/hpjansson/chafa/releases
```

## 📖 Sử dụng

### Khởi động ứng dụng

```bash
gonime
```

Hoặc sử dụng giao diện TUI:

```bash
gonime tui
```

### Các lệnh cơ bản

```bash
# Tìm kiếm anime
gonime anime "Tên anime"

# Tìm kiếm manga
gonime manga "Tên manga"

# Xem phiên bản
gonime version

# Mở cài đặt
gonime settings
```

### Phím tắt trong TUI

| Phím | Chức năng |
|------|-----------|
| `↑/k` | Di chuyển lên |
| `↓/j` | Di chuyển xuống |
| `Enter` | Chọn/Xác nhận |
| `Esc` | Quay lại |
| `q` | Thoát |
| `d` | Xóa (trong lịch sử/yêu thích) |

## 📁 Cấu trúc dự án

```
gonime/
├── cmd/
│   └── gonime/
│       └── main.go          # Điểm khởi đầu
├── internal/
│   ├── api/                 # Client API anime
│   ├── config/              # Quản lý cấu hình
│   ├── player/              # Phát video và tải xuống
│   ├── scraper/             # Scraper manga
│   └── tui/                 # Giao diện TUI
│       ├── main.go          # Model chính
│       ├── anime.go         # Model anime
│       ├── manga.go         # Model manga
│       ├── history.go       # Model lịch sử
│       ├── favorites.go     # Model yêu thích
│       └── settings.go      # Model cài đặt
├── assets/                  # Dữ liệu và tài nguyên
├── docs/                    # Tài liệu
├── go.mod
├── go.sum
└── README.md
```

## ⚙️ Cấu hình

File cấu hình được lưu tại:
- Windows: `%APPDATA%\gonime\config.json`

### Cấu hình mặc định

```json
{
  "player": "mpv",
  "download_dir": "C:\\Users\\YourName\\Downloads\\gonime",
  "manga_image_dir": "C:\\Users\\YourName\\Downloads\\gonime\\manga",
  "manga_image_size": 100
}
```

## 🔧 Phát triển

### Chạy thử nghiệm

```bash
go run ./cmd/gonime
```

### Build với tối ưu hóa

```bash
go build -ldflags="-s -w" -o gonime.exe ./cmd/gonime
```

### Chạy tests

```bash
go test ./...
```

## 📝 Lưu ý

1. **Quảng cáo**: Một số nguồn video có thể chứa quảng cáo
2. **Bản quyền**: Ứng dụng chỉ thu thập dữ liệu từ các nguồn công khai
3. **Hiệu suất**: Cần kết nối internet ổn định để xem/tải nội dung

## 🤝 Đóng góp

Mọi đóng góp đều được chào đón! Vui lòng:
1. Fork repository
2. Tạo branch mới (`git checkout -b feature/amazing-feature`)
3. Commit thay đổi (`git commit -m 'Add amazing feature'`)
4. Push lên branch (`git push origin feature/amazing-feature`)
5. Tạo Pull Request

## 📜 Giấy phép

Dự án được phân phối dưới giấy phép MIT. Xem file `LICENSE` để biết thêm chi tiết.

## 🙏 Cảm ơn

- [bubbletea](https://github.com/charmbracelet/bubbletea) - Framework TUI
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Styling cho terminal
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing

---

<div align="center">

**Made with ❤️ for Anime & Manga lovers**

</div>
