# Kế hoạch triển khai ứng dụng GoNime CLI bằng Go

Dự án này sẽ viết lại `anisub` thành một ứng dụng CLI native cho Windows bằng ngôn ngữ Go. Mục tiêu là tạo ra một tệp `.exe` duy nhất, hiệu năng cao và chạy mượt mà trên Command Prompt hoặc PowerShell.

## Giai đoạn 1: Thiết lập dự án và Cấu trúc

1.  **Khởi tạo dự án Go:**
    *   [ ] Tạo một thư mục dự án mới.
    *   [ ] Chạy `go mod init gonime` để khởi tạo Go module.
    *   [ ] Cài đặt các thư viện cần thiết:
        *   `github.com/charmbracelet/bubbletea` (cho giao diện TUI)
        *   `github.com/charmbracelet/lipgloss` (để tạo kiểu cho giao diện)
        *   `github.com/spf13/cobra` (để quản lý lệnh và cờ)
        *   `github.com/PuerkitoBio/goquery` (để scrape dữ liệu web)

2.  **Tổ chức cấu trúc thư mục:**
    *   [ ] `cmd/gonime/main.go`: Điểm khởi đầu của ứng dụng.
    *   [ ] `internal/api/`: Chứa code để tương tác với các API bên ngoài (phimapi.com).
    *   [ ] `internal/scraper/`: Chứa code để scrape dữ liệu từ các trang web manga.
    *   [ ] `internal/tui/`: Chứa các "model" của Bubble Tea, định nghĩa các thành phần giao diện (menu, list, spinner...).
    *   [ ] `internal/config/`: Quản lý cấu hình của người dùng (lưu/tải file config).

## Giai đoạn 2: Tương tác với API và Scraper

1.  **Client cho API Anime:**
    *   [ ] Trong `internal/api/`, tạo các hàm để tìm kiếm và lấy chi tiết anime từ `phimapi.com`.
    *   [ ] Định nghĩa các `struct` trong Go để ánh xạ dữ liệu JSON trả về, giúp code an toàn và dễ quản lý hơn.

2.  **Scraper cho Manga:**
    *   [ ] Trong `internal/scraper/`, sử dụng `net/http` để tải trang HTML và `goquery` để phân tích, trích xuất dữ liệu (danh sách truyện, chapter, link ảnh).
    *   [ ] Viết các hàm để lấy truyện mới, tìm kiếm truyện và lấy danh sách ảnh của một chapter.

## Giai đoạn 3: Xây dựng Giao diện Người dùng (TUI)

1.  **Tạo Model chính (Main Model):**
    *   [ ] Trong `internal/tui/`, tạo một model Bubble Tea chính để hiển thị menu ban đầu: `Xem Anime`, `Đọc Manga`, `Cài đặt`, `Thoát`.
    *   [ ] Xử lý logic điều hướng khi người dùng chọn một mục.

2.  **Model danh sách (List Model):**
    *   [ ] Tạo một model có thể tái sử dụng để hiển thị danh sách (kết quả tìm kiếm, danh sách tập, danh sách chapter).
    *   [ ] Thêm các tính năng như phân trang, loading spinner khi đang tải dữ liệu.

3.  **Model đọc truyện (Reader Model):**
    *   [ ] Đây là phần phức tạp nhất. Model này sẽ hiển thị từng ảnh của chapter manga.
    *   [ ] Tích hợp việc gọi `chafa.exe` (một công cụ bên ngoài) để render ảnh ngay trong terminal. Người dùng sẽ cần cài đặt `chafa` và thêm nó vào PATH.
    *   [ ] Xử lý việc tải ảnh về một thư mục tạm trước khi hiển thị.

4.  **Model cài đặt (Settings Model):**
    *   [ ] Tạo giao diện cho phép người dùng chỉnh sửa các cấu hình, ví dụ như đường dẫn tới trình phát video.

## Giai đoạn 4: Triển khai các tính năng Cốt lõi

1.  **Phát video:**
    *   [ ] Sử dụng package `os/exec` của Go để gọi một trình phát video bên ngoài (ví dụ `mpv.exe`, `vlc.exe`).
    *   [ ] Truyền URL của video làm tham số cho trình phát.

2.  **Tải xuống:**
    *   [ ] Tương tự, sử dụng `os/exec` để thực thi `yt-dlp.exe` với các tham số cần thiết.
    *   [ ] (Nâng cao) Hiển thị tiến trình tải xuống bằng cách đọc output từ `yt-dlp`.

3.  **Lịch sử và Yêu thích:**
    *   [ ] Sử dụng package `encoding/json` để lưu và đọc dữ liệu từ một file `history.json` và `favorites.json`.
    *   [ ] Lưu file này trong thư mục cấu hình của người dùng (có thể lấy đường dẫn bằng `os.UserConfigDir()`).

## Giai đoạn 5: Hoàn thiện và Phân phối

1.  **Biên dịch:**
    *   [ ] Sử dụng `go build -o gonime.exe ./cmd/gonime` để biên dịch toàn bộ dự án thành một file `gonime.exe`.
    *   [ ] (Nâng cao) Sử dụng `ldflags` để loại bỏ thông tin gỡ lỗi, giúp giảm kích thước file thực thi.

2.  **Đóng gói và Hướng dẫn:**
    *   [ ] Viết một tệp `README.md` rõ ràng.
    *   [ ] Hướng dẫn người dùng cách cài đặt và sử dụng ứng dụng.
    *   **Quan trọng:** Hướng dẫn người dùng cài đặt các phụ thuộc bên ngoài và đảm bảo chúng có trong `PATH` của hệ thống:
        *   Một trình phát video (khuyến nghị `mpv`).
        *   `yt-dlp.exe` (cho tính năng tải video).
        *   `chafa.exe` (cho tính năng đọc manga).
