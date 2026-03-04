package i18n

import "gonime/internal/config"

var vi = map[string]string{
	"app.subtitle":          "🎌  Anime & Manga CLI Player  ·  Vietnamese Subtitles  🎌",
	"menu.watch_anime":      "Xem Anime",
	"menu.watch_anime_desc": "Khám phá & xem anime vietsub",
	"menu.read_manga":       "Đọc Manga",
	"menu.read_manga_desc":  "Đọc manga từ nhiều nguồn",
	"menu.history":          "Lịch sử xem",
	"menu.history_desc":     "Xem lại các anime/manga đã xem",
	"menu.favorites":        "Yêu thích",
	"menu.favorites_desc":   "Danh sách anime/manga yêu thích",
	"menu.settings":         "Cài đặt",
	"menu.settings_desc":    "Cấu hình ứng dụng",
	"menu.exit":             "Thoát",
	"menu.exit_desc":        "Thoát khỏi ứng dụng",

	"title_trending":  "XU HƯỚNG",
	"title_search":    "TÌM KIẾM",
	"title_episodes":  "DANH SÁCH TẬP",
	"title_chapters":  "DANH SÁCH CHƯƠNG",
	"title_history":   "LỊCH SỬ XEM",
	"title_favorites": "YÊU THÍCH",
	"title_settings":  "CÀI ĐẶT",

	"tab_anime": "Anime / Phim",
	"tab_manga": "Manga",

	"hint_esc_back":       "Esc: Quay lại",
	"hint_esc_menu":       "Esc: Về Menu",
	"hint_nav":            "↑↓/jk: Di chuyển",
	"hint_enter_open":     "Enter: Mở",
	"hint_enter_play":     "Enter: Phát",
	"hint_enter_read":     "Enter: Đọc",
	"hint_enter_select":   "Enter: Tìm/Chọn",
	"hint_enter_select_2": "Enter: Chọn",
	"hint_enter_edit":     "Enter: Sửa",
	"hint_enter_reopen":   "Enter: Mở lại",
	"hint_source":         "m: Đổi nguồn",
	"hint_clear":          "Ctrl+U: Xóa",
	"hint_next_prev":      "n/p: Tiếp/Trước",
	"hint_pgup_dn":        "PgUp/Dn: ±10",
	"hint_first_last":     "g/G: Đầu/Cuối",
	"hint_fav":            "f: ⭐",
	"hint_server":         "s: Server",
	"hint_switch_tab":     "1/2: Anime/Manga",
	"hint_del_item":       "d: Xóa mục",
	"hint_del_all":        "c: Xóa tất cả",
	"hint_default":        "r: Mặc định",
	"hint_quit":           "q: Thoát",

	"loading_trending": "Đang tải xu hướng...",
	"loading_search":   "Đang tìm kiếm...",
	"loading_episodes": "Đang tải danh sách tập...",
	"loading_player":   "Đang khởi động trình phát...",
	"loading_chapters": "Đang tải danh sách chapter...",
	"loading_chapter":  "Đang tải chapter...",
	"loading_image":    "Đang render ảnh...",
	"loading_generic":  "Đang tải...",

	"search_anime_prompt": "Nhập tên anime:",
	"search_manga_prompt": "Nhập tên manga:",
	"search_anime_hint":   "Gõ tên anime rồi nhấn [Enter] để tìm (≥2 ký tự)",
	"search_manga_hint":   "Gõ tên manga rồi nhấn [Enter] để tìm (≥2 ký tự)",
	"search_results":      "%d kết quả",
	"search_matches":      "~%d khớp",

	"err_prefix": "Lỗi: %v",

	"msg_added_fav":   "Đã thêm vào yêu thích",
	"msg_removed_fav": "Đã xóa khỏi yêu thích",
	"msg_del_item":    "Đã xóa khỏi danh sách",
	"msg_del_all":     "Đã xóa toàn bộ danh sách",

	"empty_fav_anime":  "Chưa có anime/phim yêu thích nào.",
	"empty_fav_manga":  "Chưa có manga yêu thích nào.",
	"empty_fav_hint_a": "Nhấn [f] khi xem anime/phim để thêm vào đây.",
	"empty_fav_hint_m": "Nhấn [f] khi đọc manga để thêm vào đây.",
	"empty_hist_anime": "Chưa xem anime/phim nào.",
	"empty_hist_manga": "Chưa đọc manga nào.",
	"empty_hist_hint":  "Khám phá rồi quay lại đây nhé!",

	"count_fav_anime":  "%d anime/phim yêu thích",
	"count_fav_manga":  "%d manga yêu thích",
	"count_hist_anime": "%d mục anime/phim",
	"count_hist_manga": "%d mục manga",
	"count_item_hist":  "목 %d mục trong lịch sử",

	"set_player":     "🎬 Trình phát video",
	"set_down_dir":   "📁 Thư mục tải xuống",
	"set_manga_dir":  "🖼  Thư mục ảnh manga",
	"set_manga_size": "📐 Kích thước ảnh manga (%)",
	"set_manga_src":  "🌐 URL nguồn Manga (sẽ được update trong tương lai)",
	"set_lang":       "🈯 Ngôn ngữ giao diện (vi/en)",
	"set_sys_check":  "KIỂM TRA HỆ THỐNG",
	"set_config":     "CẤU HÌNH",
	"sys_ready":      "sẵn sàng",
	"sys_not_found":  "không tìm thấy",

	"hint_nav_page":   "←→/↑↓: Trang",
	"hint_next_page":  "Enter: Trang sau",
	"hint_np_chapter": "n/p: Chapter",
	"hint_viewer":     "o: Viewer",
	"hint_browser":    "w: Trình duyệt",
	"hint_retry":      "Thử lại",
}

var en = map[string]string{
	"app.subtitle":          "🎌  Anime & Manga CLI Player  ·  English Interface  🎌",
	"menu.watch_anime":      "Watch Anime",
	"menu.watch_anime_desc": "Browse & watch anime/movies",
	"menu.read_manga":       "Read Manga",
	"menu.read_manga_desc":  "Read manga from multiple sources",
	"menu.history":          "History",
	"menu.history_desc":     "Recently watched/Read items",
	"menu.favorites":        "Favorites",
	"menu.favorites_desc":   "Your favorite anime/manga list",
	"menu.settings":         "Settings",
	"menu.settings_desc":    "Configure GoNime",
	"menu.exit":             "Exit",
	"menu.exit_desc":        "Quit application",

	"title_trending":  "TRENDING",
	"title_search":    "SEARCH",
	"title_episodes":  "EPISODES",
	"title_chapters":  "CHAPTERS",
	"title_history":   "HISTORY",
	"title_favorites": "FAVORITES",
	"title_settings":  "SETTINGS",

	"tab_anime": "Anime / Movies",
	"tab_manga": "Manga",

	"hint_esc_back":       "Esc: Back",
	"hint_esc_menu":       "Esc: Menu",
	"hint_nav":            "↑↓/jk: Navigate",
	"hint_enter_open":     "Enter: Open",
	"hint_enter_play":     "Enter: Play",
	"hint_enter_read":     "Enter: Read",
	"hint_enter_select":   "Enter: Search/Select",
	"hint_enter_select_2": "Enter: Select",
	"hint_enter_edit":     "Enter: Edit",
	"hint_enter_reopen":   "Enter: Reopen",
	"hint_source":         "m: Source",
	"hint_clear":          "Ctrl+U: Clear",
	"hint_next_prev":      "n/p: Next/Prev",
	"hint_pgup_dn":        "PgUp/Dn: ±10",
	"hint_first_last":     "g/G: Top/Bottom",
	"hint_fav":            "f: ⭐",
	"hint_server":         "s: Server",
	"hint_switch_tab":     "1/2: Anime/Manga",
	"hint_del_item":       "d: Delete",
	"hint_del_all":        "c: Clear All",
	"hint_default":        "r: Default",
	"hint_quit":           "q: Quit",

	"loading_trending": "Loading trending...",
	"loading_search":   "Searching...",
	"loading_episodes": "Loading episodes...",
	"loading_player":   "Starting player...",
	"loading_chapters": "Loading chapters...",
	"loading_chapter":  "Loading chapter...",
	"loading_image":    "Rendering image...",
	"loading_generic":  "Loading...",

	"search_anime_prompt": "Enter anime name:",
	"search_manga_prompt": "Enter manga name:",
	"search_anime_hint":   "Type anime name and press [Enter] to search (≥2 chars)",
	"search_manga_hint":   "Type manga name and press [Enter] to search (≥2 chars)",
	"search_results":      "%d results",
	"search_matches":      "~%d matches",

	"err_prefix": "Error: %v",

	"msg_added_fav":   "Added to favorites",
	"msg_removed_fav": "Removed from favorites",
	"msg_del_item":    "Removed from list",
	"msg_del_all":     "Cleared entire list",

	"empty_fav_anime":  "No favorite anime/movies yet.",
	"empty_fav_manga":  "No favorite manga yet.",
	"empty_fav_hint_a": "Press [f] while watching anime to add.",
	"empty_fav_hint_m": "Press [f] while reading manga to add.",
	"empty_hist_anime": "No watched anime/movies.",
	"empty_hist_manga": "No read manga.",
	"empty_hist_hint":  "Explore then check back here!",

	"count_fav_anime":  "%d favorite anime/movies",
	"count_fav_manga":  "%d favorite manga",
	"count_hist_anime": "%d anime/movie items",
	"count_hist_manga": "%d manga items",
	"count_item_hist":  "%d items in history",

	"set_player":     "🎬 Video Player",
	"set_down_dir":   "📁 Download Directory",
	"set_manga_dir":  "🖼  Manga Image Directory",
	"set_manga_size": "📐 Manga Image Size (%)",
	"set_manga_src":  "🌐 Manga Source URL (future updates)",
	"set_lang":       "🈯 Interface Language (vi/en)",
	"set_sys_check":  "SYSTEM CHECK",
	"set_config":     "CONFIGURATION",
	"sys_ready":      "ready",
	"sys_not_found":  "not found",

	"hint_nav_page":   "←→/↑↓: Page",
	"hint_next_page":  "Enter: Next Page",
	"hint_np_chapter": "n/p: Chapter",
	"hint_viewer":     "o: Viewer",
	"hint_browser":    "w: Browser",
	"hint_retry":      "Retry",
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
