package viewer

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
)

//go:embed frontend
var frontendFS embed.FS

// OpenPayload defines payload format for viewer Open API
type OpenPayload struct {
	Images []string `json:"images"`
	Start  int      `json:"start"`
}

// ViewerClient defines the abstraction for opening a set of images in a viewer.
type ViewerClient interface {
	Open(images []string, start int) error
}

// NoopViewer is a placeholder implementation used during initial wiring.
type NoopViewer struct{}

func (n *NoopViewer) Open(images []string, start int) error {
	// No-op: viewer not implemented yet
	return nil
}

func NewViewer() ViewerClient {
	return &NoopViewer{}
}

// Http-based viewer skeleton for Phase 2.2
type HttpViewer struct {
	started bool
	port    int
}

func (h *HttpViewer) Open(images []string, start int) error {
	// Start a tiny HTTP server on a random port if not started
	if !h.started {
		// Simple local HTTP server to accept /open
		go func() {
			mux := http.NewServeMux()
			mux.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				var payload OpenPayload
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				log.Printf("[viewer-http-skeleton] Received Open: images=%d start=%d", len(payload.Images), payload.Start)
				w.WriteHeader(http.StatusOK)
			})
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				log.Println("viewer http listen error:", err)
				return
			}
			h.port = ln.Addr().(*net.TCPAddr).Port
			log.Println("[viewer-http-skeleton] listening on port", h.port)
			http.Serve(ln, mux)
		}()
		h.started = true
	}
	payload := OpenPayload{Images: images, Start: start}
	data, _ := json.Marshal(payload)
	url := fmt.Sprintf("http://127.0.0.1:%d/open", h.port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("viewer http open failed: status %d", resp.StatusCode)
	}
	log.Printf("[viewer-http-skeleton] Open request delivered to viewer on port %d", h.port)
	return nil
}

func NewHttpViewer() ViewerClient {
	return &HttpViewer{}
}

// Wails skeleton viewer for Phase 2.2 preparation
type WailsViewer struct {
	started bool
	images  []string
	start   int
}

func (w *WailsViewer) Open(images []string, start int) error {
	w.started = true
	w.images = images
	w.start = start
	// Lightweight log to indicate viewer wiring is active
	fmt.Printf("[viewer-skeleton] Open called with %d images, start=%d\n", len(images), start)
	return nil
}

func NewWailsViewer() ViewerClient {
	return &WailsViewer{}
}
