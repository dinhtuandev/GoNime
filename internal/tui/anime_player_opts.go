package tui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gonime/internal/player"
)

// playWithOptions chạy player với phụ đề, tốc độ và HTTP headers nếu có hỗ trợ.
func playWithOptions(p *player.Player, url, title, subPath string, speed float64, referer string) error {

	baseName := strings.TrimSuffix(filepath.Base(p.Path), ".exe")
	baseName = strings.ToLower(baseName)

	var args []string

	switch baseName {
	case "mpv":
		args = append(args,
			url,
			"--force-window=yes",
			"--title="+title,
			"--really-quiet",
			"--geometry=640x480",
		)
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
		args = append(args, url)
		// VLC metadata
		args = append(args, "--meta-title="+title)
		args = append(args, "--no-video-title-show")
		args = append(args, "--play-and-exit")
		// Network caching để tránh buffering
		args = append(args, "--network-caching=5000")
		// HTTP headers - VLC không hỗ trợ tốt cho m3u8, thử nhiều cách
		if referer != "" {
			args = append(args, "--http-referrer="+referer)
		}
		// User-Agent cần thiết cho một số CDN
		args = append(args, "--http-user-agent=GoNime/1.0")
		// Thêm options để VLC xử lý tốt hơn streaming
		args = append(args, "--file-caching=3000")
		args = append(args, "--live-caching=3000")
		// Disable deinterlace cho m3u8 streams
		args = append(args, "--deinterlace=auto")
		if speed > 0 && speed != 1.0 {
			args = append(args, fmt.Sprintf("--rate=%g", speed))
		}
		if subPath != "" {
			args = append(args, "--sub-file="+subPath)
		}
	default:
		args = append(args, url)
	}

	log.Printf("[player] Starting %s with args: %v", baseName, args)

	cmd := exec.Command(p.Path, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	setDetachedProcess(cmd)

	err := cmd.Start()
	if err != nil {
		log.Printf("[player] Failed to start %s: %v", baseName, err)
		return fmt.Errorf("không thể khởi động %s: %v", baseName, err)
	}

	log.Printf("[player] Successfully started %s", baseName)
	return nil
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
