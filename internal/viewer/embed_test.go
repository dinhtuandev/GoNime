package viewer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestEmbeddedFrontend ensures the viewer UI ships inside the binary.
func TestEmbeddedFrontend(t *testing.T) {
	for _, name := range []string{"frontend/index.html", "frontend/app.js"} {
		if _, err := fs.Stat(frontendFS, name); err != nil {
			t.Fatalf("embedded %s missing: %v", name, err)
		}
	}
}

// TestViewerServerServesUI actually starts the HTTP server and fetches the
// embedded page and chapter state, catching regressions that a compile-only
// check would miss.
func TestViewerServerServesUI(t *testing.T) {
	w := &WailsViewerReal{}
	w.SetChapters([]ChapterInfo{{Title: "Chương 1", URL: "https://example.com/1"}})

	if err := w.startServer(); err != nil {
		t.Fatalf("startServer: %v", err)
	}

	base := fmt.Sprintf("http://127.0.0.1:%d", w.serverPort)
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(base + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "GoNime Manga Viewer") {
		t.Fatalf("GET / = %d, body does not contain viewer title", resp.StatusCode)
	}

	resp, err = client.Get(base + "/api/images")
	if err != nil {
		t.Fatalf("GET /api/images: %v", err)
	}
	defer resp.Body.Close()

	var state struct {
		Chapters     []ChapterInfo `json:"chapters"`
		ChapterIndex int           `json:"chapterIndex"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		t.Fatalf("decode /api/images: %v", err)
	}
	if len(state.Chapters) != 1 || state.Chapters[0].Title != "Chương 1" {
		t.Fatalf("chapters not exposed: %+v", state.Chapters)
	}
}
