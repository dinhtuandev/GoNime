package player

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	maxDownloadSize = 500 * 1024 * 1024 // 500MB max download
	downloadTimeout = 30 * time.Minute
)

type Player struct {
	Name string
	Path string
}

type Downloader struct {
	YTDLPPath string
}

// NewPlayer tạo player với detection nâng cao
func NewPlayer(name string) (*Player, error) {
	// Thử tìm trước bằng exec.LookPath
	path, err := exec.LookPath(name)
	if err == nil {
		return &Player{Name: name, Path: path}, nil
	}

	// Thêm detection cho Windows Store/MSIX packages
	if runtime.GOOS == "windows" {
		// Thử tìm trong WindowsApps (MSIX packages)
		windowsApps := filepath.Join(os.Getenv("ProgramFiles"), "WindowsApps")
		if checkFileExists(windowsApps) {
			// Tìm kiếm các ứng dụng mpv và vlc trong WindowsApps
			searchPaths := []string{
				filepath.Join(windowsApps, "mpv*.exe"),
				filepath.Join(windowsApps, "vlc*.exe"),
			}
			for _, pattern := range searchPaths {
				matches, _ := filepath.Glob(pattern)
				if len(matches) > 0 {
					return &Player{Name: name, Path: matches[0]}, nil
				}
			}
		}
	}

	// Thử tìm trong các đường dẫn phổ biến trên Windows
	searchPaths := getCommonPlayerPaths(name)
	for _, p := range searchPaths {
		if checkFileExists(p) {
			return &Player{Name: name, Path: p}, nil
		}
	}

	// Thêm fallback cho Windows Store apps
	if runtime.GOOS == "windows" {
		// Thử tìm trong Microsoft Store apps
		storePaths := []string{
			`C:\Program Files\WindowsApps\Microsoft.Mpv_1.0.0.0_x64__8wekyb3d8bbwe\mpv.exe`,
			`C:\Program Files\WindowsApps\VideoLAN.VLC_3.0.16.0_x64__ql6h4bf0d3p0t\vlc.exe`,
		}
		for _, p := range storePaths {
			if checkFileExists(p) {
				return &Player{Name: name, Path: p}, nil
			}
		}
	}

	return nil, fmt.Errorf("không tìm thấy trình phát %s: %v", name, err)
}

// getCommonPlayerPaths trả về danh sách đường dẫn phổ biến cho player trên Windows
func getCommonPlayerPaths(name string) []string {
	localAppData := os.Getenv("LOCALAPPDATA")
	userProfile := os.Getenv("USERPROFILE")

	paths := []string{
		// Standard Program Files
		`C:\Program Files\mpv\mpv.exe`,
		`C:\Program Files (x86)\mpv\mpv.exe`,
		`C:\Program Files\VideoLAN\VLC\vlc.exe`,
		`C:\Program Files (x86)\VideoLAN\VLC\vlc.exe`,
	}

	if localAppData != "" {
		paths = append(paths,
			// WinGet symlinks — đây là nơi winget tạo symlink
			filepath.Join(localAppData, "Microsoft", "WinGet", "Links", "mpv.exe"),
			filepath.Join(localAppData, "Microsoft", "WinGet", "Links", "vlc.exe"),
			// WindowsApps aliases — MSIX/Store apps
			filepath.Join(localAppData, "Microsoft", "WindowsApps", "mpv.exe"),
			filepath.Join(localAppData, "Microsoft", "WindowsApps", "vlc.exe"),
			// WinGet Packages (portable) — shinchiro.mpv
			filepath.Join(localAppData, "Microsoft", "WinGet", "Packages", "shinchiro.mpv_Microsoft.Winget.Source_8wekyb3d8bbwe", "mpv.exe"),
		)
	}

	if userProfile != "" {
		paths = append(paths,
			// Scoop
			filepath.Join(userProfile, "scoop", "shims", "mpv.exe"),
			filepath.Join(userProfile, "scoop", "shims", "vlc.exe"),
			// Chocolatey
			filepath.Join(userProfile, "AppData", "Local", "Programs", "mpv", "mpv.exe"),
		)
	}

	// Filter theo tên player nếu cần
	if name == "mpv" || name == "mpv.exe" {
		var filtered []string
		for _, p := range paths {
			if strings.Contains(strings.ToLower(p), "mpv") {
				filtered = append(filtered, p)
			}
		}
		return filtered
	}
	if name == "vlc" || name == "vlc.exe" {
		var filtered []string
		for _, p := range paths {
			if strings.Contains(strings.ToLower(p), "vlc") {
				filtered = append(filtered, p)
			}
		}
		return filtered
	}

	return paths
}

// FindAnyPlayer tìm bất kỳ trình phát nào có sẵn
func FindAnyPlayer() (*Player, error) {
	names := []string{"mpv", "vlc", "mpv.exe", "vlc.exe", "ffplay"}
	for _, name := range names {
		if p, err := NewPlayer(name); err == nil {
			return p, nil
		}
	}
	return nil, fmt.Errorf("không tìm thấy trình phát video nào (mpv, vlc).\n\n💡 Cài đặt:\n  winget install shinchiro.mpv\n  winget install VideoLAN.VLC\n  scoop install mpv\n  choco install mpv\n\nĐể xem anime, bạn cần cài đặt một trong các trình phát này.")
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (p *Player) Play(url string, title string) error {
	ctx := context.Background()
	return p.PlayWithContext(ctx, url, title)
}

// PlayWithContext cải thiện flags cho Windows
func (p *Player) PlayWithContext(ctx context.Context, url string, title string) error {
	var cmd *exec.Cmd

	baseName := strings.TrimSuffix(filepath.Base(p.Path), ".exe")
	baseName = strings.ToLower(baseName)

	switch baseName {
	case "mpv":
		// Cải thiện flags cho Windows
		cmd = exec.CommandContext(ctx, p.Path,
			url,
			"--force-window=yes",
			"--title="+title,
			"--no-terminal",
			"--geometry=640x480",
		)
	case "vlc":
		// Cải thiện flags cho Windows
		cmd = exec.CommandContext(ctx, p.Path,
			url,
			"--meta-title="+title,
			"--qt-start-minimized",
			"--qt-notification=0",
		)
	default:
		cmd = exec.CommandContext(ctx, p.Path, url)
	}

	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

// PlayWithBrowser mở URL bằng trình duyệt mặc định (fallback cuối)
func PlayWithBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", "", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

func (p *Player) IsRunning() bool {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+p.Name+".exe")
		output, _ := cmd.Output()
		return strings.Contains(string(output), p.Name)
	case "darwin", "linux":
		cmd := exec.Command("pgrep", "-x", p.Name)
		_, err := cmd.Output()
		return err == nil
	default:
		return false
	}
}

func (p *Player) Kill() error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("taskkill", "/IM", p.Name+".exe", "/F")
		return cmd.Run()
	case "darwin", "linux":
		cmd := exec.Command("pkill", "-x", p.Name)
		return cmd.Run()
	default:
		return fmt.Errorf("không hỗ trợ hệ điều hành này")
	}
}

func NewDownloader() (*Downloader, error) {
	path, err := exec.LookPath("yt-dlp")
	if err != nil {
		path, err = exec.LookPath("yt-dlp.exe")
		if err != nil {
			return nil, fmt.Errorf("không tìm thấy yt-dlp: %v", err)
		}
	}
	return &Downloader{YTDLPPath: path}, nil
}

func (d *Downloader) Download(url string, outputPath string, progressCallback func(progress string)) error {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	outputDir := filepath.Dir(outputPath)
	outputTemplate := filepath.Join(outputDir, "%(title)s.%(ext)s")

	cmd := exec.CommandContext(ctx, d.YTDLPPath,
		url,
		"-o", outputTemplate,
		"--no-playlist",
		"--newline",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if progressCallback != nil {
			progressCallback(string(buf[:n]))
		}
	}

	return cmd.Wait()
}

func (d *Downloader) DownloadWithFormat(url string, outputPath string, format string) error {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, d.YTDLPPath,
		url,
		"-o", outputPath,
		"-f", format,
		"--no-playlist",
	)
	return cmd.Run()
}

func DownloadFile(url string, outputPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "GoNime/1.0")

	client := &http.Client{
		Timeout: downloadTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lỗi tải file: %s", resp.Status)
	}

	if resp.ContentLength > maxDownloadSize {
		return fmt.Errorf("file quá lớn: %d bytes (tối đa: %d)", resp.ContentLength, maxDownloadSize)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	limitedReader := io.LimitReader(resp.Body, maxDownloadSize)
	_, err = io.Copy(out, limitedReader)
	return err
}

func GetDefaultPlayer() string {
	players := []string{"mpv", "vlc", "ffplay"}
	for _, p := range players {
		if _, err := exec.LookPath(p); err == nil {
			return p
		}
	}
	return ""
}
