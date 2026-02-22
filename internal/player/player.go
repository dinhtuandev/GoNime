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

func NewPlayer(name string) (*Player, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy trình phát %s: %v", name, err)
	}
	return &Player{Name: name, Path: path}, nil
}

func (p *Player) Play(url string, title string) error {
	ctx := context.Background()
	return p.PlayWithContext(ctx, url, title)
}

func (p *Player) PlayWithContext(ctx context.Context, url string, title string) error {
	var cmd *exec.Cmd

	switch p.Name {
	case "mpv":
		cmd = exec.CommandContext(ctx, p.Path,
			url,
			"--no-terminal",
			"--force-window",
			"--title="+title,
		)
	case "vlc":
		cmd = exec.CommandContext(ctx, p.Path,
			url,
			"--meta-title="+title,
		)
	default:
		cmd = exec.CommandContext(ctx, p.Path, url)
	}

	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Start()
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
