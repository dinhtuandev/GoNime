package i18n

import "gonime/internal/config"

var vi = map[string]string{
	// ── Anime screens ─────────────────────────────────────────────
	"anime_title":        "ANIME / PHIM",
	"anime_tab_trending": "Xu hướng",
	"anime_tab_search":   "Tìm kiếm",

	"anime_loading_trending": "Đang tải xu hướng...",
	"anime_loading_search":   "Đang tìm kiếm...",
	"anime_loading_eps":      "Đang tải danh sách tập...",
	"anime_loading_player":   "Đang mở trình phát...",

	"anime_empty_trending":        "Không có xu hướng.",
	"anime_empty_trending_hint_1": "Thử đổi nguồn bằng phím",
	"anime_empty_trending_hint_2": "để xem nguồn khác.",
	"anime_empty_eps":             "Không có tập nào.",

	"anime_search_prompt":   "Nhập tên anime:",
	"anime_search_hint_1":   "Gõ tên anime rồi nhấn",
	"anime_search_hint_2":   "để tìm kiếm (≥2 ký tự).",
	"anime_search_results":  "kết quả",
	"anime_search_fuzzy":    "khớp",
	"anime_search_no_match": "Không tìm thấy kết quả cho",

	"anime_ep_list_title": "DANH SÁCH TẬP",
	"anime_ep_progress":   "Tiến độ:",
	"anime_ep_total":      "Tổng",
	"anime_ep_unit":       "tập",
	"anime_ep_watching":   "Đang xem",
	"anime_ep_prefix":     "Tập",

	"anime_msg_last_ep":    "Đây là tập cuối cùng.",
	"anime_msg_first_ep":   "Đây là tập đầu tiên.",
	"anime_msg_one_server": "Chỉ có 1 server.",

	"anime_playing_ep": "Đang phát:",
	"anime_dl_sent":    "Đang tải:",

	// ── Manga screens ─────────────────────────────────────────────
	"manga_title":        "MANGA",
	"manga_tab_trending": "Xu hướng",
	"manga_tab_search":   "Tìm kiếm",

	"manga_loading_trending": "Đang tải xu hướng...",
	"manga_loading_search":   "Đang tìm kiếm...",
	"manga_loading_chaps":    "Đang tải danh sách chương...",
	"manga_loading_pages":    "Đang tải trang...",
	"manga_rendering":        "Đang render ảnh...",

	"manga_empty_trending":        "Không có xu hướng.",
	"manga_empty_trending_hint_1": "Thử đổi nguồn bằng phím",
	"manga_empty_trending_hint_2": "để xem nguồn khác.",
	"manga_empty_chaps":           "Không có chương nào.",
	"manga_empty_pages":           "Không tải được ảnh.",

	"manga_search_prompt":   "Nhập tên manga:",
	"manga_search_hint_1":   "Gõ tên manga rồi nhấn",
	"manga_search_hint_2":   "để tìm kiếm (≥2 ký tự).",
	"manga_search_results":  "kết quả",
	"manga_search_fuzzy":    "khớp",
	"manga_search_no_match": "Không tìm thấy kết quả cho",

	"manga_chap_list_title": "DANH SÁCH CHƯƠNG",
	"manga_chap_total":      "Tổng số chương",
	"manga_tag_reading":     "đang đọc",

	"manga_reading_title": "ĐANG ĐỌC",
	"manga_page_prefix":   "Trang",

	"manga_msg_last_chap":    "Đây là chương mới nhất.",
	"manga_msg_first_chap":   "Đây là chương cuối trong danh sách.",
	"manga_msg_browser_open": "Đã mở trình duyệt.",

	// ── Common status messages ────────────────────────────────────
	"status_fav_added":       "Đã thêm vào yêu thích",
	"status_source_switched": "Đã đổi nguồn",

	// ── Hint shortcuts ────────────────────────────────────────────
	"hint_switch_source": "🔀 Đổi nguồn",
	"hint_search":        "Tìm kiếm",
	"hint_start_end":     "⏮ g/G: Đầu/Cuối",

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

	"tab_anime": "🎬 Anime / Phim",
	"tab_manga": "📖 Manga",

	"hint_esc_back":       "↩ Esc: Quay lại",
	"hint_esc_menu":       "Esc: Về Menu",
	"hint_nav":            "🧭 ↑↓/jk: Di chuyển",
	"hint_enter_open":     "📂 Enter: Mở",
	"hint_enter_play":     "▶ Enter: Phát",
	"hint_enter_read":     "📖 Enter: Đọc",
	"hint_enter_select":   "🔍 Enter: Tìm/Chọn",
	"hint_enter_select_2": "✅ Enter: Chọn",
	"hint_enter_edit":     "✏ Enter: Sửa",
	"hint_enter_reopen":   "🔄 Enter: Mở lại",
	"hint_source":         "m: Đổi nguồn",
	"hint_clear":          "🧹 Ctrl+U: Xóa",
	"hint_next_prev":      "🔁 n/p: Tiếp/Trước",
	"hint_pgup_dn":        "📄 PgUp/Dn: ±10",
	"hint_first_last":     "⏮ g/G: Đầu/Cuối",
	"hint_fav":            "f: ⭐",
	"hint_server":         "🖥 s: Server",
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

	"msg_added_fav":           "Đã thêm vào yêu thích",
	"msg_removed_fav":         "Đã xóa khỏi yêu thích",
	"msg_del_item":            "Đã xóa khỏi danh sách",
	"msg_del_all":             "Đã xóa toàn bộ danh sách",
	"err_no_play_url":         "Không có link phát",
	"err_fav_manga_not_found": "Không mở được manga này từ dữ liệu đã lưu",

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
	"label_anime":      "🎬 Anime",
	"label_manga":      "📚 Manga",

	"set_player":                  "🎬 Trình phát video",
	"set_down_dir":                "📁 Thư mục tải xuống",
	"set_manga_dir":               "🖼  Thư mục ảnh manga",
	"set_manga_size":              "📐 Kích thước ảnh manga (%)",
	"set_manga_src":               "🌐 URL nguồn Manga (sẽ được update trong tương lai)",
	"set_image_mode":              "🖼 Chế độ hiển thị ảnh (auto/native/halfblock/chafa)",
	"set_lang":                    "🈯 Ngôn ngữ giao diện (vi/en)",
	"set_sys_check":               "KIỂM TRA HỆ THỐNG",
	"set_config":                  "CẤU HÌNH",
	"sys_ready":                   "sẵn sàng",
	"sys_not_found":               "không tìm thấy",
	"set_saved":                   "✅ Đã lưu: %s",
	"set_reset_ok":                "✅ Đặt lại mặc định thành công",
	"set_need_chafa":              "(cần để xem manga)",
	"set_need_ytdlp":              "(cần cho tải tập bằng phím d)",
	"set_need_native_terminal":    "(cần terminal hỗ trợ kitty/sixel/iTerm)",
	"set_warn_native_unsupported": "Image mode=native nhưng terminal hiện tại không hỗ trợ; app sẽ fallback",
	"set_err_load":                "Không thể tải cấu hình.",
	"set_err_save":                "Không lưu được cấu hình: %v",
	"set_err_cfg_unavailable":     "Cấu hình không khả dụng",
	"set_err_player_empty":        "Tên trình phát không được để trống",
	"set_err_download_dir_empty":  "Thư mục tải không được để trống",
	"set_err_manga_dir_empty":     "Thư mục manga không được để trống",
	"set_err_size_number":         "Vui lòng nhập số từ 50 đến 200",
	"set_err_size_range":          "Kích thước phải từ 50 đến 200",
	"set_err_url_empty":           "URL không được để trống",
	"set_err_image_mode":          "Image mode hỗ trợ 'auto', 'native', 'halfblock' hoặc 'chafa'",
	"set_err_lang":                "Ngôn ngữ chỉ hỗ trợ 'vi' hoặc 'en'",

	"hint_nav_page":   "📄 ←→/↑↓: Trang",
	"hint_wheel_page": "🖱 Wheel: Lên/Xuống trang",
	"hint_next_page":  "Enter: Trang sau",
	"hint_np_chapter": "n/p: Chapter",
	"hint_viewer":     "o: Viewer",
	"hint_browser":    "w: Trình duyệt",
	"hint_retry":      "Thử lại",
}

var en = map[string]string{
	// ── Anime screens ─────────────────────────────────────────────
	"anime_title":        "ANIME / MOVIES",
	"anime_tab_trending": "Trending",
	"anime_tab_search":   "Search",

	"anime_loading_trending": "Loading trending...",
	"anime_loading_search":   "Searching...",
	"anime_loading_eps":      "Loading episodes...",
	"anime_loading_player":   "Starting player...",

	"anime_empty_trending":        "No trending items.",
	"anime_empty_trending_hint_1": "Try switching source with",
	"anime_empty_trending_hint_2": "to see another source.",
	"anime_empty_eps":             "No episodes found.",

	"anime_search_prompt":   "Enter anime name:",
	"anime_search_hint_1":   "Type an anime name then press",
	"anime_search_hint_2":   "to search (≥2 chars).",
	"anime_search_results":  "results",
	"anime_search_fuzzy":    "matches",
	"anime_search_no_match": "No results for",

	"anime_ep_list_title": "EPISODES",
	"anime_ep_progress":   "Progress:",
	"anime_ep_total":      "Total",
	"anime_ep_unit":       "eps",
	"anime_ep_watching":   "Watching",
	"anime_ep_prefix":     "Ep",

	"anime_msg_last_ep":    "This is the last episode.",
	"anime_msg_first_ep":   "This is the first episode.",
	"anime_msg_one_server": "Only 1 server available.",

	"anime_playing_ep": "Playing:",
	"anime_dl_sent":    "Downloading:",

	// ── Manga screens ─────────────────────────────────────────────
	"manga_title":        "MANGA",
	"manga_tab_trending": "Trending",
	"manga_tab_search":   "Search",

	"manga_loading_trending": "Loading trending...",
	"manga_loading_search":   "Searching...",
	"manga_loading_chaps":    "Loading chapters...",
	"manga_loading_pages":    "Loading pages...",
	"manga_rendering":        "Rendering image...",

	"manga_empty_trending":        "No trending manga.",
	"manga_empty_trending_hint_1": "Try switching source with",
	"manga_empty_trending_hint_2": "to see another source.",
	"manga_empty_chaps":           "No chapters found.",
	"manga_empty_pages":           "Could not load images.",

	"manga_search_prompt":   "Enter manga name:",
	"manga_search_hint_1":   "Type a manga name then press",
	"manga_search_hint_2":   "to search (≥2 chars).",
	"manga_search_results":  "results",
	"manga_search_fuzzy":    "matches",
	"manga_search_no_match": "No results for",

	"manga_chap_list_title": "CHAPTERS",
	"manga_chap_total":      "Total chapters",
	"manga_tag_reading":     "reading",

	"manga_reading_title": "READING",
	"manga_page_prefix":   "Page",

	"manga_msg_last_chap":    "This is the latest chapter.",
	"manga_msg_first_chap":   "This is the last chapter in list.",
	"manga_msg_browser_open": "Opened in browser.",

	// ── Common status messages ────────────────────────────────────
	"status_fav_added":       "Added to favorites",
	"status_source_switched": "Source switched",

	// ── Hint shortcuts ────────────────────────────────────────────
	"hint_switch_source": "🔀 Switch source",
	"hint_search":        "Search",
	"hint_start_end":     "⏮ g/G: Top/Bottom",

	"app.subtitle":          "🎌  Anime & Manga CLI Player  ·  English Interface  🎌",
	"menu.watch_anime":      "Watch Anime",
	"menu.watch_anime_desc": "Browse & watch anime/movies",
	"menu.read_manga":       "Read Manga",
	"menu.read_manga_desc":  "Read manga from multiple sources",

	"menu.history":        "History",
	"menu.history_desc":   "Recently watched/Read items",
	"menu.favorites":      "Favorites",
	"menu.favorites_desc": "Your favorite anime/manga list",
	"menu.settings":       "Settings",
	"menu.settings_desc":  "Configure GoNime",
	"menu.exit":           "Exit",
	"menu.exit_desc":      "Quit application",

	"title_trending":  "TRENDING",
	"title_search":    "SEARCH",
	"title_episodes":  "EPISODES",
	"title_chapters":  "CHAPTERS",
	"title_history":   "HISTORY",
	"title_favorites": "FAVORITES",
	"title_settings":  "SETTINGS",

	"tab_anime": "🎬 Anime / Movies",
	"tab_manga": "📖 Manga",

	"hint_esc_back":       "↩ Esc: Back",
	"hint_esc_menu":       "Esc: Menu",
	"hint_nav":            "🧭 ↑↓/jk: Navigate",
	"hint_enter_open":     "📂 Enter: Open",
	"hint_enter_play":     "▶ Enter: Play",
	"hint_enter_read":     "📖 Enter: Read",
	"hint_enter_select":   "🔍 Enter: Search/Select",
	"hint_enter_select_2": "✅ Enter: Select",
	"hint_enter_edit":     "✏ Enter: Edit",
	"hint_enter_reopen":   "🔄 Enter: Reopen",
	"hint_source":         "m: Source",
	"hint_clear":          "🧹 Ctrl+U: Clear",
	"hint_next_prev":      "🔁 n/p: Next/Prev",
	"hint_pgup_dn":        "📄 PgUp/Dn: ±10",
	"hint_first_last":     "⏮ g/G: Top/Bottom",
	"hint_fav":            "f: ⭐",
	"hint_server":         "🖥 s: Server",
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

	"msg_added_fav":           "Added to favorites",
	"msg_removed_fav":         "Removed from favorites",
	"msg_del_item":            "Removed from list",
	"msg_del_all":             "Cleared entire list",
	"err_no_play_url":         "No playable URL",
	"err_fav_manga_not_found": "Could not open this saved manga entry",

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
	"label_anime":      "🎬 Anime",
	"label_manga":      "📚 Manga",

	"set_player":                  "🎬 Video Player",
	"set_down_dir":                "📁 Download Directory",
	"set_manga_dir":               "🖼  Manga Image Directory",
	"set_manga_size":              "📐 Manga Image Size (%)",
	"set_manga_src":               "🌐 Manga Source URL (future updates)",
	"set_image_mode":              "🖼 Image Display Mode (auto/native/halfblock/chafa)",
	"set_lang":                    "🈯 Interface Language (vi/en)",
	"set_sys_check":               "SYSTEM CHECK",
	"set_config":                  "CONFIGURATION",
	"sys_ready":                   "ready",
	"sys_not_found":               "not found",
	"set_saved":                   "✅ Saved: %s",
	"set_reset_ok":                "✅ Reset to defaults",
	"set_need_chafa":              "(required for manga preview)",
	"set_need_ytdlp":              "(required for episode download with key d)",
	"set_need_native_terminal":    "(requires terminal with kitty/sixel/iTerm support)",
	"set_warn_native_unsupported": "Image mode=native is not supported by current terminal; app will fallback",
	"set_err_load":                "Could not load configuration.",
	"set_err_save":                "Could not save configuration: %v",
	"set_err_cfg_unavailable":     "Configuration is unavailable",
	"set_err_player_empty":        "Player name cannot be empty",
	"set_err_download_dir_empty":  "Download directory cannot be empty",
	"set_err_manga_dir_empty":     "Manga directory cannot be empty",
	"set_err_size_number":         "Please enter a number from 50 to 200",
	"set_err_size_range":          "Size must be between 50 and 200",
	"set_err_url_empty":           "URL cannot be empty",
	"set_err_image_mode":          "Image mode supports 'auto', 'native', 'halfblock' or 'chafa'",
	"set_err_lang":                "Language only supports 'vi' or 'en'",

	"hint_nav_page":   "📄 ←→/↑↓: Page",
	"hint_wheel_page": "🖱 Wheel: Prev/Next page",
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
