// Package render cung cấp khả năng render ảnh ra terminal bằng chafa.
package render

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// maxCachedImageSize giới hạn dung lượng ảnh tải về cache.
const maxCachedImageSize = 50 * 1024 * 1024

// ChafaRenderer download ảnh và render ra terminal bằng chafa
type ChafaRenderer struct {
	chafaPath string // đường dẫn tới chafa.exe
	cacheDir  string // thư mục cache ảnh
	Width     int    // chiều rộng terminal (chars)
	Height    int    // chiều cao terminal (chars)
	Mode      string // auto | native | chafa
}

// NewChafaRenderer tạo renderer, tự tìm chafa trong PATH và common paths
func NewChafaRenderer(cacheDir string, width, height int) (*ChafaRenderer, error) {
	return NewImageRenderer(cacheDir, width, height, "chafa")
}

// NewImageRenderer tạo renderer với image mode: auto/native/halfblock/chafa.
// Chỉ mode chafa bắt buộc phải có binary chafa.
func NewImageRenderer(cacheDir string, width, height int, mode string) (*ChafaRenderer, error) {
	mode = NormalizeMode(mode)
	path, findErr := findChafa()
	if findErr != nil && mode == "chafa" {
		return nil, findErr
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("không tạo được cache dir: %v", err)
	}
	return &ChafaRenderer{
		chafaPath: path,
		cacheDir:  cacheDir,
		Width:     width,
		Height:    height,
		Mode:      mode,
	}, nil
}

func NormalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "auto", "native", "chafa", "halfblock":
		return strings.ToLower(strings.TrimSpace(mode))
	case "truecolor", "blocks":
		return "halfblock"
	default:
		return "auto"
	}
}

func SupportsNativeImage() bool {
	term := strings.ToLower(os.Getenv("TERM"))
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	if os.Getenv("KITTY_WINDOW_ID") != "" || strings.Contains(term, "kitty") {
		return true
	}
	if strings.Contains(termProgram, "iterm") {
		return true
	}
	if strings.Contains(term, "sixel") {
		return true
	}
	return false
}

func DetectNativeFormat() string {
	term := strings.ToLower(os.Getenv("TERM"))
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	if os.Getenv("KITTY_WINDOW_ID") != "" || strings.Contains(term, "kitty") {
		return "kitty"
	}
	if strings.Contains(termProgram, "iterm") {
		return "iterm"
	}
	if strings.Contains(term, "sixel") {
		return "sixels"
	}
	return ""
}

// ResolveEffectiveMode chọn mode thực tế sau khi xét môi trường terminal.
func ResolveEffectiveMode(mode string) string {
	mode = NormalizeMode(mode)
	switch mode {
	case "native":
		if SupportsNativeImage() {
			return "native"
		}
		return "halfblock"
	case "chafa":
		return "chafa"
	case "halfblock":
		return "halfblock"
	default: // auto
		if SupportsNativeImage() {
			return "native"
		}
		return "halfblock"
	}
}

// IsAvailable kiểm tra chafa có cài không (không trả về renderer)
func IsAvailable() bool {
	_, err := findChafa()
	return err == nil
}

// findChafa tìm chafa trong PATH và common install paths trên Windows
func findChafa() (string, error) {
	// Thử trong PATH trước
	if path, err := exec.LookPath("chafa"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("chafa.exe"); err == nil {
		return path, nil
	}

	// Tìm trong vị trí WinGet cài
	home, _ := os.UserHomeDir()
	commonPaths := []string{
		filepath.Join(home, "AppData", "Local", "Microsoft", "WinGet", "Packages"),
		`C:\Program Files\chafa`,
		`C:\Program Files (x86)\chafa`,
		filepath.Join(home, "scoop", "shims"),
	}

	for _, base := range commonPaths {
		// Đệ quy tìm chafa.exe trong thư mục
		var found string
		_ = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil || found != "" {
				return nil
			}
			if !info.IsDir() && strings.EqualFold(info.Name(), "chafa.exe") {
				found = path
			}
			return nil
		})
		if found != "" {
			return found, nil
		}
	}

	return "", fmt.Errorf("không tìm thấy chafa. Cài bằng: winget install hpjansson.Chafa")
}

// RenderURL download URL ảnh, cache lại, render bằng chafa, trả về chuỗi ANSI
func (r *ChafaRenderer) RenderURL(ctx context.Context, imageURL string) (string, error) {
	// B1: Download (hoặc lấy cache)
	localPath, err := r.downloadWithCache(ctx, imageURL)
	if err != nil {
		return "", fmt.Errorf("không tải được ảnh: %v", err)
	}

	// B2: Render bằng chafa
	return r.RenderFile(localPath)
}

// RenderFile render file ảnh local bằng chafa
func (r *ChafaRenderer) RenderFile(path string) (string, error) {
	w := r.Width
	h := r.Height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}

	size := fmt.Sprintf("%dx%d", w, h)
	mode := ResolveEffectiveMode(r.Mode)

	if mode == "native" {
		if out, err := r.renderNative(path, size); err == nil {
			return out, nil
		}
		mode = "halfblock"
	}

	if mode == "halfblock" {
		if out, err := RenderHalfBlock(path, w, h); err == nil {
			return out, nil
		}
	}

	return r.renderSymbols(path, size)
}

func (r *ChafaRenderer) renderNative(path, size string) (string, error) {
	format := DetectNativeFormat()
	if format == "" {
		return "", fmt.Errorf("native format not supported")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, r.chafaPath,
		"--size", size,
		"--format", format,
		"--animate", "off",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (r *ChafaRenderer) renderSymbols(path, size string) (string, error) {
	if r.chafaPath == "" {
		return "", fmt.Errorf("chafa không khả dụng. Cài bằng: winget install hpjansson.Chafa")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Dùng --symbols block để tương thích tốt với Windows Terminal
	cmd := exec.CommandContext(ctx, r.chafaPath,
		"--size", size,
		"--colors", "256",
		"--symbols", "block+border+space",
		"--animate", "off",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		// Thử lại với mode đơn giản hơn
		cmd2 := exec.CommandContext(context.Background(), r.chafaPath,
			"--size", size,
			"--symbols", "ascii",
			path,
		)
		out2, err2 := cmd2.Output()
		if err2 != nil {
			return "", fmt.Errorf("chafa render lỗi: %v", err)
		}
		return string(out2), nil
	}
	return string(out), nil
}

// RenderThumbnail render a thumbnail version of an image URL using chafa.
// It downloads the image (with cache) if needed, then renders with smaller dims.
func (r *ChafaRenderer) RenderThumbnail(ctx context.Context, imageURL string, width int, height int) (string, error) {
	// Ensure image is cached locally
	localPath, err := r.downloadWithCache(ctx, imageURL)
	if err != nil {
		return "", err
	}

	// Temporarily override size for thumbnail rendering
	oldW, oldH := r.Width, r.Height
	r.Width, r.Height = width, height
	defer func() { r.Width, r.Height = oldW, oldH }()

	return r.RenderFile(localPath)
}

// downloadWithCache tải ảnh về cache dir, trả về đường dẫn local
func (r *ChafaRenderer) downloadWithCache(ctx context.Context, imageURL string) (string, error) {
	// Tạo tên file từ URL
	hash := sanitizeFilename(imageURL)
	// Giữ extension
	ext := filepath.Ext(imageURL)
	if ext == "" || len(ext) > 5 {
		ext = ".jpg"
	}
	// Lấy tối đa 40 ký tự cuối URL làm tên
	cachePath := filepath.Join(r.cacheDir, hash+ext)

	// Nếu đã cache → dùng luôn
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	// Download
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if resp.ContentLength > maxCachedImageSize {
		return "", fmt.Errorf("ảnh quá lớn: %d bytes (tối đa: %d)", resp.ContentLength, maxCachedImageSize)
	}

	f, err := os.Create(cachePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, io.LimitReader(resp.Body, maxCachedImageSize)); err != nil {
		os.Remove(cachePath)
		return "", err
	}

	return cachePath, nil
}

// sanitizeFilename tạo tên file an toàn từ URL
func sanitizeFilename(url string) string {
	// Lấy 60 ký tự cuối, thay ký tự đặc biệt
	if len(url) > 60 {
		url = url[len(url)-60:]
	}
	replacer := strings.NewReplacer(
		"/", "_",
		":", "_",
		"?", "_",
		"&", "_",
		"=", "_",
		".", "_",
		" ", "_",
	)
	return replacer.Replace(url)
}
