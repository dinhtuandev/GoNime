# GoNime CLI - Pipeline Description

Mục tiêu
- Xây dựng một CLI native Go cho Windows (điểm phát hành .exe duy nhất) cho việc xem anime/phim, đọc manga, tải xuống và quản lý history/yêu thích.
- Đảm bảo hiệu năng cao, dễ bảo trì và có thể mở rộng với TUI (Bubble Tea) và scraper/API modular.

## 1. Kiến trúc hệ thống (High-level)
- Client API: giao tiếp với phimapi.com (HTTP, JSON) để tìm kiếm và lấy chi tiết anime.
- Scraper Manga: thu thập danh sách truyện, danh sách chapter và ảnh từ truyenvn.shop bằng net/http và goquery.
- TUI (Bubble Tea): giao diện người dùng dựa trên terminal với các mô hình riêng cho từng phần (Main menu, Anime, Manga, History, Favorites, Settings).
- Player & Downloader: phát video (mpv/vlc) và tải video bằng yt-dlp; hỗ trợ tải ảnh manga và quản lý tiến trình tải.
- Cấu hình: lưu trữ cấu hình người dùng và lịch sử ở thư mục cấu hình người dùng.
- Build & phân phối: đóng gói một executable duy nhất và đi kèm documentation.

## 2. Dữ liệu và mô hình (Data Models)
- API Anime
  - Anime, Episode: meta, danh sách ep với URL video/stream.
- Manga
  - Manga, Chapter: thông tin cơ bản, danh sách chapters, ảnh chapter.
- History
  - HistoryItem: name, episode, url, timestamp
- Favorites
  - FavoriteItem: name, slug
- Config
  - Config: player, download_dir, manga_image_dir, manga_image_size

## 3. Pipeline triển khai (Phases)

Phase 1 — Thiết lập và scaffolding
- Khởi tạo module Go và thiết kế cấu trúc thư mục:
  - cmd/gonime: entrypoint cho executable
  - internal/api/: client API cho phimapi.com
  - internal/scraper/: scraper cho manga
  - internal/tui/: model và view cho Bubble Tea
  - internal/config/: quản lý config, history, favorites
- Cài đặt dependencies: bubbletea, lipgloss, cobra, goquery
- Tạo skeleton cho các file chính: main.go, các model TUI và API client skeleton

Phase 2 — Implement API client & scraper
- API client (internal/api/client.go):
  - Sử dụng timeout HTTP (ví dụ 30s)
  - Encode query và slug an toàn
  - Định nghĩa structs cho JSON mapping
- Manga scraper (internal/scraper/manga.go):
  - Tìm kiếm manga, lấy chi tiết, danh sách chapters, danh sách hình ảnh
  - Giới hạn kích thước ảnh và xử lý network errors

Phase 3 — Xây dựng TUI (Bubble Tea)
- Mẫu MainModel (internal/tui/main.go): menu chính: Xem Anime, Đọc Manga, Lịch sử, Yêu thích, Cài đặt, Thoát
- AnimeModel (internal/tui/anime.go): tìm kiếm, danh sách anime, chọn anime, xem danh sách tập
- MangaModel (internal/tui/manga.go): tìm kiếm manga, chọn manga, danh sách chapter, xem chapter
- HistoryModel, FavoritesModel, SettingsModel: quản lý các danh sách và cấu hình
- Model Reader (image viewer): hiển thị ảnh chapter với con trỏ người dùng và hỗ trợ lệnh đọc

Phase 4 — Tính năng lõi
- Phát video: dùng os/exec để gọi mpv/vlc với URL video
- Tải xuống: dùng yt-dlp để tải video (hoặc fallback bằng ffmpeg nếu vừa đủ)
- Lịch sử và Yêu thích: lưu ở JSON/Config; cập nhật khi người dùng xem/tải/ký gửi

Phase 5 — Build, phân phối và runtime
- Compile single exe: go build -ldflags "-s -w" -o gonime.exe ./cmd/gonime
- Đóng gói dữ liệu mẫu (assets/aniw_export_2026-01-14.csv) và README đầy đủ
- Thiết lập CI/CD cho build multi-OS (Windows/Linux/macOS) và test

## 4. Quy trình CI/CD (điểm mốc tham khảo)
- GitHub Actions (hoặc CI tương tự):
  - Checkout -> Cache Go modules
  - go vet ./...
  - go test ./...
  - go build -o gonime.exe ./cmd/gonime (Windows target)
  - Build cho Linux/macOS nếu cần (optional)
- Kiểm tra lẫn nhau: căn cứ vào lỗi compile, lint (nếu có), và unit tests
- Bản release: đóng gói binary và tài liệu kèm; upload lên release page

## 5. Quy trình làm việc (Workflow)
- Task phân chia theo module: api, scraper, tui, config, player
- Pull requests kiểm tra chất lượng: compile, test, và chạy thử basic
- Lưu trữ dữ liệu nhạy cảm: không commit secrets; config và history được quản lý tại user config dir

## 6. Bảo mật và ổn định
- Timeout và cân bằng retry cho requests
- Validate input từ người dùng (query, slug, URLs)
- Giới hạn kích thước download/images để tránh DoS và OOM
- Atomic writes cho config/history/favorites
- Cross-platform support: Windows/macOS/Linux; xử lý điều hành tương ứng (ví dụ cho player)

## 7. Kết quả đầu ra và deliverables
- gonime.exe (single binary cho Windows)
- README.md, tài liệu vận hành, và ví dụ usage
- assets dữ liệu mẫu (csv)
- Các file mẫu cho config/history/favorites

## 8. Rủi ro và mitigations
- Rủi ro: footgun khi scraping và phụ thuộc nguồn; Mitigation: error handling, timeouts, and fallbacks
- Rủi ro: dependency drift; Mitigation: lock versions, go mod tidy, và vendor nếu needed
- Rủi ro: user data colliding; Mitigation: atomic writes và validation

## 9. Tiêu chuẩn chất lượng (Quality gates)
- Build success on CI
- 0 lỗi compile; go vet pass
- Basic tests (where applicable)
- Documentation up-to-date

## 10. Quickstart

- Chuẩn bị môi trường:
  - Go 1.21+ hoặc tương thích
  - Trình phát video (mpv/vlc)
  - yt-dlp (tải) và chafa (preview thumbnail) nếu có

- Build từ source:
  go mod download
  go build -ldflags="-s -w" -o gonime.exe ./cmd/gonime

- Chạy:
  gonime.exe

## 11. Liên hệ và đóng góp
- Góp ý, báo lỗi và đóng góp qua PR.
- Tài liệu này có thể được mở rộng theo nhu cầu thêm tính năng.

---

End of plan.
