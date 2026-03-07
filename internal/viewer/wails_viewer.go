package viewer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// WailsViewerReal is a real Wails-based viewer integration (Phase 3)
// It serves a web-based viewer frontend and handles IPC for manga viewing
type WailsViewerReal struct {
	started             bool
	images              []string
	start               int
	serverPort          int
	server              *http.Server
	chapters            []ChapterInfo
	currentChapterIndex int
	scraper             interface {
		GetChapterImages(url string) ([]string, error)
	}
}

// ChapterInfo represents a manga chapter
type ChapterInfo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Open starts the viewer with the given images and opens it in the default browser
func (w *WailsViewerReal) Open(images []string, start int) error {
	w.images = images
	w.start = start

	openedNow := false
	if !w.started {
		if err := w.startServer(); err != nil {
			return fmt.Errorf("failed to start viewer server: %v", err)
		}
		w.started = true
		openedNow = true
	}

	// Update the viewer with new images
	w.updateViewer()

	// Open browser only on first run. The frontend polls /api/images for updates.
	if openedNow {
		url := fmt.Sprintf("http://localhost:%d", w.serverPort)
		log.Printf("[viewer-wails] Opening viewer at %s", url)
		if err := w.openBrowser(url); err != nil {
			log.Printf("[viewer-wails] Warning: could not open browser: %v", err)
		}
	}

	return nil
}

// startServer starts the HTTP server to serve the viewer frontend
func (w *WailsViewerReal) startServer() error {
	// Find available port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to find available port: %v", err)
	}
	w.serverPort = ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	// Get the directory where the executable is located
	frontendDir := w.getFrontendDir()

	mux := http.NewServeMux()

	// Serve static files from frontend directory
	fs := http.FileServer(http.Dir(frontendDir))
	mux.Handle("/", fs)

	// API endpoint to get current images
	mux.HandleFunc("/api/images", w.handleImages)

	// API endpoint to update viewer
	mux.HandleFunc("/api/update", w.handleUpdate)

	// API endpoint to change chapter
	mux.HandleFunc("/api/chapter", w.handleChapterChange)

	w.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", w.serverPort),
		Handler: mux,
	}

	go func() {
		log.Printf("[viewer-wails] Server starting on port %d", w.serverPort)
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[viewer-wails] Server error: %v", err)
		}
	}()

	return nil
}

// getFrontendDir returns the directory containing the frontend files
func (w *WailsViewerReal) getFrontendDir() string {
	// Try multiple locations for the frontend directory
	possibleDirs := []string{
		// Current directory (when running from project root)
		"internal/viewer/frontend",
		// Relative to executable (production)
		filepath.Join(filepath.Dir(os.Args[0]), "internal/viewer/frontend"),
		// Absolute path from source
		"D:\\WorkingSpace\\gonime\\internal\\viewer\\frontend",
	}

	for _, dir := range possibleDirs {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}

	// Fallback to current working directory
	return "internal/viewer/frontend"
}

// handleImages returns the current images as JSON
func (w *WailsViewerReal) handleImages(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")

	chapterTitle := ""
	if w.currentChapterIndex >= 0 && w.currentChapterIndex < len(w.chapters) {
		chapterTitle = w.chapters[w.currentChapterIndex].Title
	}

	json.NewEncoder(wr).Encode(map[string]interface{}{
		"images":       w.images,
		"start":        w.start,
		"chapters":     w.chapters,
		"chapterIndex": w.currentChapterIndex,
		"chapterTitle": chapterTitle,
	})
}

// handleChapterChange handles chapter navigation requests
func (w *WailsViewerReal) handleChapterChange(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		wr.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Index int `json:"index"`
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wr).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Validate chapter index
	if data.Index < 0 || data.Index >= len(w.chapters) {
		wr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wr).Encode(map[string]string{"error": "invalid chapter index"})
		return
	}

	// Get images for the new chapter
	w.currentChapterIndex = data.Index
	chapterURL := w.chapters[data.Index].URL

	var images []string
	var err error

	// Try to get images from scraper if available
	if w.scraper != nil {
		images, err = w.scraper.GetChapterImages(chapterURL)
		if err != nil {
			log.Printf("[viewer-wails] Failed to get chapter images: %v", err)
		}
	}

	w.images = images
	w.start = 0

	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]interface{}{
		"images":       w.images,
		"start":        w.start,
		"chapterTitle": w.chapters[data.Index].Title,
	})
}

// handleUpdate handles POST requests to update viewer data
func (w *WailsViewerReal) handleUpdate(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		wr.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Images []string `json:"images"`
		Start  int      `json:"start"`
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		wr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wr).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.images = data.Images
	w.start = data.Start

	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]string{"status": "ok"})
}

// updateViewer sends updated data to the browser if open
func (w *WailsViewerReal) updateViewer() {
	if w.serverPort == 0 {
		return
	}

	updateURL := fmt.Sprintf("http://127.0.0.1:%d/api/update", w.serverPort)
	payload, _ := json.Marshal(map[string]interface{}{
		"images": w.images,
		"start":  w.start,
	})

	_, _ = http.Post(updateURL, "application/json", bytes.NewBuffer(payload))
}

// openBrowser opens the URL in the default browser
func (w *WailsViewerReal) openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// WailsViewerReal creates a new Wails viewer instance
func NewWailsViewerReal() ViewerClient {
	return &WailsViewerReal{}
}

// SetChapters sets the chapters for navigation
func (w *WailsViewerReal) SetChapters(chapters []ChapterInfo) {
	w.chapters = chapters
	w.currentChapterIndex = 0
}

// SetScraper sets the scraper for fetching chapter images
func (w *WailsViewerReal) SetScraper(scraper interface {
	GetChapterImages(url string) ([]string, error)
}) {
	w.scraper = scraper
}
