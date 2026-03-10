package tui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gonime/internal/player"
)

// playWithOptions chạy player với phụ đề, tốc độ và HTTP headers nếu có hỗ trợ.
func playWithOptions(p *player.Player, url, title, subPath string, speed float64, referer string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
	defer cancel()

	baseName := strings.TrimSuffix(filepath.Base(p.Path), ".exe")
	baseName = strings.ToLower(baseName)

	var args []string

	switch baseName {
	case "mpv":
		args = append(args,
			url,
			"--force-window=yes",
			"--title="+title,
			"--no-terminal",
			"--geometry=640x480",
		)
		// Gửi Referer + User-Agent cho CDN chấp nhận stream
		if referer != "" {
			args = append(args,
				"--http-header-fields=Referer: "+referer+",User-Agent: GoNime/1.0",
			)
		}
		if speed > 0 && speed != 1.0 {
			args = append(args, fmt.Sprintf("--speed=%g", speed))
		}
		if subPath != "" {
			args = append(args, "--sub-file="+subPath)
		}
	case "vlc":
		args = append(args,
			url,
			"--meta-title="+title,
			"--no-video-title-show",
			"--play-and-exit",
		)
		// Gửi Referer + User-Agent cho CDN chấp nhận stream
		if referer != "" {
			args = append(args, "--http-referrer="+referer)
			args = append(args, "--http-user-agent=GoNime/1.0")
		}
		if speed > 0 && speed != 1.0 {
			// VLC dùng rate khoảng 0.25–4.0
			args = append(args, fmt.Sprintf("--rate=%g", speed))
		}
		if subPath != "" {
			args = append(args, "--sub-file="+subPath)
		}
	default:
		args = append(args, url)
	}

	cmd := exec.CommandContext(ctx, p.Path, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	// Detach the player process so it runs independently of the TUI's terminal.
	// On Windows this is critical: without DETACHED_PROCESS, mpv installed via
	// Windows Store / winget may fail to open its GUI window on first launch.
	setDetachedProcess(cmd)
	return cmd.Start()
}

// findLocalSubtitle tìm phụ đề .srt/.ass trong subtitlesDir dựa trên slug + tên tập.
// Quy ước tên: {slug}-{epName}.srt/.ass hoặc chỉ slug.srt.
func findLocalSubtitle(subtitlesDir, slug, epName string) string {
	if subtitlesDir == "" || slug == "" {
		return ""
	}

	_ = os.MkdirAll(subtitlesDir, 0o755)

	safeSlug := sanitizeSlug(slug)
	safeEp := sanitizeSlug(epName)

	candidates := []string{
		filepath.Join(subtitlesDir, safeSlug+"-"+safeEp+".srt"),
		filepath.Join(subtitlesDir, safeSlug+"-"+safeEp+".ass"),
		filepath.Join(subtitlesDir, safeSlug+".srt"),
		filepath.Join(subtitlesDir, safeSlug+".ass"),
	}

	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

func sanitizeSlug(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	replacer := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"?", "_",
		"&", "_",
	)
	return replacer.Replace(s)
}
