package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadSave(t *testing.T) {
	// Tạo thư mục tạm để test
	tempDir, err := os.MkdirTemp("", "gonime-test-*")
	if err != nil {
		t.Fatalf("Không tạo được temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	m := &Manager{
		ConfigDir:     tempDir,
		ConfigFile:    filepath.Join(tempDir, "config.json"),
		HistoryFile:   filepath.Join(tempDir, "history.json"),
		FavoritesFile: filepath.Join(tempDir, "favorites.json"),
	}

	// 1. Kiểm tra load mặc định khi chưa có file
	cfg, err := m.LoadConfig()
	if err != nil {
		t.Errorf("Lỗi load config mặc định: %v", err)
	}
	if cfg.MangaURL == "" {
		t.Error("MangaURL mặc định không được để trống")
	}

	// 2. Chỉnh sửa và lưu
	cfg.Player = "custom-player"
	cfg.MangaURL = "https://test-manga.com"
	err = m.SaveConfig(cfg)
	if err != nil {
		t.Errorf("Lỗi save config: %v", err)
	}

	// 3. Load lại và kiểm tra
	cfg2, err := m.LoadConfig()
	if err != nil {
		t.Errorf("Lỗi load lại config: %v", err)
	}
	if cfg2.Player != "custom-player" {
		t.Errorf("Player không khớp: mong đợi custom-player, nhận được %s", cfg2.Player)
	}
	if cfg2.MangaURL != "https://test-manga.com" {
		t.Errorf("MangaURL không khớp: %s", cfg2.MangaURL)
	}
}

func TestHistory(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "gonime-hist-*")
	defer os.RemoveAll(tempDir)

	m := &Manager{
		ConfigDir:   tempDir,
		HistoryFile: filepath.Join(tempDir, "history.json"),
	}

	item := HistoryItem{
		Name:    "Test Anime",
		Episode: "1",
		Type:    "anime",
	}

	err := m.AddToHistory(item)
	if err != nil {
		t.Errorf("Lỗi add history: %v", err)
	}

	items, err := m.LoadHistory()
	if err != nil || len(items) == 0 {
		t.Fatal("Không load được history")
	}

	if items[0].Name != "Test Anime" {
		t.Errorf("Tên history không khớp: %s", items[0].Name)
	}
}

func TestNormalizeImageMode(t *testing.T) {
	tests := map[string]string{
		"auto":   "auto",
		"native": "native",
		"chafa":  "chafa",
		"foo":    "auto",
		"":       "auto",
	}
	for in, want := range tests {
		if got := NormalizeImageMode(in); got != want {
			t.Fatalf("NormalizeImageMode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestConfigValidateSetsImageModeDefault(t *testing.T) {
	cfg := &Config{
		Player:         "mpv",
		DownloadDir:    "C:/tmp",
		MangaImageDir:  "C:/tmp/manga",
		MangaImageSize: 100,
		MangaURL:       "https://example.com",
		ImageMode:      "invalid",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
	if cfg.ImageMode != "auto" {
		t.Fatalf("expected image mode normalized to auto, got %q", cfg.ImageMode)
	}
}
