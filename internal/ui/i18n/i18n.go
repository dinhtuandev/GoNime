package i18n

import "gonime/internal/config"

// Keyed translations for vi/en. Chỉ dùng cho text lớn, không phải mọi string.
var vi = map[string]string{
	"app.subtitle":        "🎌  Anime & Manga CLI Player  ·  Vietnamese Subtitles  🎌",
	"menu.watch_anime":    "Xem Anime",
	"menu.watch_anime_desc": "Khám phá & xem anime vietsub",
	"menu.read_manga":     "Đọc Manga",
	"menu.read_manga_desc": "Đọc manga từ nhiều nguồn",
	"menu.history":        "Lịch sử",
	"menu.history_desc":   "Xem lại các anime/manga đã xem",
	"menu.favorites":      "Yêu thích",
	"menu.favorites_desc": "Danh sách anime/manga yêu thích",
	"menu.settings":       "Cài đặt",
	"menu.settings_desc":  "Cấu hình ứng dụng",
	"menu.exit":           "Thoát",
	"menu.exit_desc":      "Thoát khỏi ứng dụng",
}

var en = map[string]string{
	"app.subtitle":        "🎌  Anime & Manga CLI Player  ·  English Interface  🎌",
	"menu.watch_anime":    "Watch Anime",
	"menu.watch_anime_desc": "Browse & watch anime/movies",
	"menu.read_manga":     "Read Manga",
	"menu.read_manga_desc": "Read manga from multiple sources",
	"menu.history":        "History",
	"menu.history_desc":   "Recently watched/Read items",
	"menu.favorites":      "Favorites",
	"menu.favorites_desc": "Your favorite anime/manga list",
	"menu.settings":       "Settings",
	"menu.settings_desc":  "Configure GoNime",
	"menu.exit":           "Exit",
	"menu.exit_desc":      "Quit application",
}

// T returns translation for key based on config language. Fallback: vi -> en -> key.
func T(cfg *config.Config, key string) string {
	lang := "vi"
	if cfg != nil && cfg.Language != "" {
		lang = cfg.Language
	}
	var dict map[string]string
	switch lang {
	case "en":
		dict = en
	default:
		dict = vi
	}
	if v, ok := dict[key]; ok {
		return v
	}
	// Fallback cross-dict
	if v, ok := vi[key]; ok {
		return v
	}
	if v, ok := en[key]; ok {
		return v
	}
	return key
}

