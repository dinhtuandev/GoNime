package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	minMangaImageSize = 50
	maxMangaImageSize = 200
	maxHistoryItems   = 100
)

type Config struct {
	Player         string `json:"player"`
	DownloadDir    string `json:"download_dir"`
	MangaImageDir  string `json:"manga_image_dir"`
	MangaImageSize int    `json:"manga_image_size"`
	MangaURL       string `json:"manga_url"`
}

func (c *Config) Validate() error {
	if c.Player == "" {
		return errors.New("player không được để trống")
	}
	if c.DownloadDir == "" {
		return errors.New("thư mục tải không được để trống")
	}
	if c.MangaImageDir == "" {
		return errors.New("thư mục manga không được để trống")
	}
	if c.MangaImageSize < minMangaImageSize || c.MangaImageSize > maxMangaImageSize {
		return errors.New("kích thước ảnh manga phải từ 50 đến 200")
	}
	return nil
}

type HistoryItem struct {
	Name      string `json:"name"`
	Episode   string `json:"episode"`
	URL       string `json:"url"`
	ThumbURL  string `json:"thumb_url,omitempty"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

type FavoriteItem struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Type     string `json:"type"`
	ThumbURL string `json:"thumb_url,omitempty"`
}

type Manager struct {
	ConfigDir     string
	ConfigFile    string
	HistoryFile   string
	FavoritesFile string
}

func NewManager() (*Manager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			configDir = ".gonime"
		} else {
			configDir = filepath.Join(homeDir, ".gonime")
		}
	} else {
		configDir = filepath.Join(configDir, "gonime")
	}

	m := &Manager{
		ConfigDir:     configDir,
		ConfigFile:    filepath.Join(configDir, "config.json"),
		HistoryFile:   filepath.Join(configDir, "history.json"),
		FavoritesFile: filepath.Join(configDir, "favorites.json"),
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) LoadConfig() (*Config, error) {
	data, err := os.ReadFile(m.ConfigFile)
	if err != nil {
		cfg := m.getDefaultConfig()
		_ = m.SaveConfig(cfg)
		return cfg, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		cfg := m.getDefaultConfig()
		_ = m.SaveConfig(cfg)
		return cfg, nil
	}

	if cfg.MangaImageSize < minMangaImageSize || cfg.MangaImageSize > maxMangaImageSize {
		cfg.MangaImageSize = 100
	}

	return &cfg, nil
}

func (m *Manager) SaveConfig(cfg *Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := m.ConfigFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, m.ConfigFile)
}

// GetDefaultConfigForManager trả về config mặc định (dùng cho Settings reset)
func GetDefaultConfigForManager(m *Manager) *Config {
	return m.getDefaultConfig()
}

func (m *Manager) getDefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = "."
	}
	downloadDir := filepath.Join(homeDir, "Downloads", "gonime")

	return &Config{
		Player:         "mpv",
		DownloadDir:    downloadDir,
		MangaImageDir:  filepath.Join(downloadDir, "manga"),
		MangaImageSize: 100,
		MangaURL:       "https://truyenvn.shop",
	}
}

func (m *Manager) LoadHistory() ([]HistoryItem, error) {
	data, err := os.ReadFile(m.HistoryFile)
	if err != nil {
		return []HistoryItem{}, nil
	}

	var items []HistoryItem
	if err := json.Unmarshal(data, &items); err != nil {
		return []HistoryItem{}, nil
	}

	return items, nil
}

func (m *Manager) SaveHistory(items []HistoryItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := m.HistoryFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, m.HistoryFile)
}

func (m *Manager) AddToHistory(item HistoryItem) error {
	items, err := m.LoadHistory()
	if err != nil {
		return err
	}

	for i, h := range items {
		if h.Name == item.Name && h.Episode == item.Episode {
			items = append(items[:i], items[i+1:]...)
			break
		}
	}

	items = append([]HistoryItem{item}, items...)
	if len(items) > maxHistoryItems {
		items = items[:maxHistoryItems]
	}

	return m.SaveHistory(items)
}

func (m *Manager) ClearHistory() error {
	return m.SaveHistory([]HistoryItem{})
}

func (m *Manager) LoadFavorites() ([]FavoriteItem, error) {
	data, err := os.ReadFile(m.FavoritesFile)
	if err != nil {
		return []FavoriteItem{}, nil
	}

	var items []FavoriteItem
	if err := json.Unmarshal(data, &items); err != nil {
		return []FavoriteItem{}, nil
	}

	return items, nil
}

func (m *Manager) SaveFavorites(items []FavoriteItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := m.FavoritesFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, m.FavoritesFile)
}

func (m *Manager) AddToFavorites(item FavoriteItem) error {
	items, err := m.LoadFavorites()
	if err != nil {
		return err
	}

	for _, f := range items {
		if f.Slug == item.Slug {
			return nil
		}
	}

	items = append(items, item)
	return m.SaveFavorites(items)
}

func (m *Manager) RemoveFromFavorites(slug string) error {
	items, err := m.LoadFavorites()
	if err != nil {
		return err
	}

	for i, f := range items {
		if f.Slug == slug {
			items = append(items[:i], items[i+1:]...)
			break
		}
	}

	return m.SaveFavorites(items)
}

func (m *Manager) ClearFavorites() error {
	return m.SaveFavorites([]FavoriteItem{})
}
