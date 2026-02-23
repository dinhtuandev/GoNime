package player

import (
	"runtime"
	"strings"
	"testing"
)

func TestPlayerDetection(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Chỉ test trên Windows")
	}

	// Test tìm kiếm path mặc định (không phụ thuộc vào việc máy đã cài hay chưa, chỉ test logic lọc tên)
	paths := getCommonPlayerPaths("mpv")
	for _, p := range paths {
		if !strings.Contains(strings.ToLower(p), "mpv") {
			t.Errorf("Path tìm cho mpv lại chứa tên khác: %s", p)
		}
	}

	pathsVLC := getCommonPlayerPaths("vlc")
	for _, p := range pathsVLC {
		if !strings.Contains(strings.ToLower(p), "vlc") {
			t.Errorf("Path tìm cho vlc lại chứa tên khác: %s", p)
		}
	}
}

func TestBrowserFallback(t *testing.T) {
	// Chỉ kiểm tra hàm không bị lỗi khi gọi (không thực sự mở browser trong unit test)
	// Hàm này dùng Start() nên trả về nil nếu process bắt chạy (không đợi kết thúc)
	url := "https://google.com"
	err := PlayWithBrowser(url)
	if err != nil {
		t.Errorf("PlayWithBrowser lỗi: %v", err)
	}
}
