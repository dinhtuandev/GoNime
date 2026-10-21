package viewer

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"
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

	// mu guards all mutable viewer state; HTTP handlers run on their own
	// goroutines while the TUI mutates the same fields.
	mu sync.Mutex
}

// ChapterInfo represents a manga chapter
type ChapterInfo struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Open starts the viewer with the given images and opens it in the default browser
func (w *WailsViewerReal) Open(images []string, start int) error {
	w.mu.Lock()
	w.images = images
	w.start = start
	started := w.started
	w.mu.Unlock()

	if !started {
		if err := w.startServer(); err != nil {
			return fmt.Errorf("failed to start viewer server: %v", err)
		}
		w.mu.Lock()
		w.started = true
		port := w.serverPort
		w.mu.Unlock()
		url := fmt.Sprintf("http://127.0.0.1:%d", port)
		log.Printf("[viewer-wails] Opening viewer at %s", url)
		if err := w.openBrowser(url); err != nil {
			log.Printf("[viewer-wails] Warning: could not open browser: %v", err)
		}
	}

	return nil
}

// startServer starts the HTTP server to serve the viewer frontend
func (w *WailsViewerReal) startServer() error {
	// Bind synchronously so the port is live before Open returns.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to find available port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port

	sub, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		ln.Close()
		return fmt.Errorf("failed to load embedded frontend: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/api/images", w.handleImages)
	mux.HandleFunc("/api/update", w.handleUpdate)
	mux.HandleFunc("/api/chapter", w.handleChapterChange)

	w.mu.Lock()
	w.serverPort = port
	srv := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	w.server = srv
	w.mu.Unlock()

	go func() {
		log.Printf("[viewer-wails] Server starting on port %d", port)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("[viewer-wails] Server error: %v", err)
		}
	}()

	return nil
}

// handleImages returns the current images as JSON
func (w *WailsViewerReal) handleImages(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-Type", "application/json")

	w.mu.Lock()
	images := w.images
	start := w.start
	chapters := w.chapters
	chapterIndex := w.currentChapterIndex
	w.mu.Unlock()

	chapterTitle := ""
	if chapterIndex >= 0 && chapterIndex < len(chapters) {
		chapterTitle = chapters[chapterIndex].Title
	}

	json.NewEncoder(wr).Encode(map[string]interface{}{
		"images":       images,
		"start":        start,
		"chapters":     chapters,
		"chapterIndex": chapterIndex,
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

	w.mu.Lock()
	if data.Index < 0 || data.Index >= len(w.chapters) {
		w.mu.Unlock()
		wr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wr).Encode(map[string]string{"error": "invalid chapter index"})
		return
	}
	title := w.chapters[data.Index].Title
	chapterURL := w.chapters[data.Index].URL
	w.currentChapterIndex = data.Index
	scraper := w.scraper
	w.mu.Unlock()

	var images []string
	if scraper != nil {
		var err error
		images, err = scraper.GetChapterImages(chapterURL)
		if err != nil {
			log.Printf("[viewer-wails] Failed to get chapter images: %v", err)
		}
	}
	if images == nil {
		images = []string{}
	}

	w.mu.Lock()
	w.images = images
	w.start = 0
	w.mu.Unlock()

	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]interface{}{
		"images":       images,
		"start":        0,
		"chapterTitle": title,
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

	w.mu.Lock()
	w.images = data.Images
	w.start = data.Start
	w.mu.Unlock()

	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]string{"status": "ok"})
}

// openBrowser opens the URL in the default browser
func (w *WailsViewerReal) openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
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
	w.mu.Lock()
	w.chapters = chapters
	w.currentChapterIndex = 0
	w.mu.Unlock()
}

// SetScraper sets the scraper for fetching chapter images
func (w *WailsViewerReal) SetScraper(scraper interface {
	GetChapterImages(url string) ([]string, error)
}) {
	w.mu.Lock()
	w.scraper = scraper
	w.mu.Unlock()
}
