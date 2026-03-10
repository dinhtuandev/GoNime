package api

// MovieSource là interface duy nhất cho mọi nguồn phim.
// GetTrending nhận page (1-based). Nguồn nào không hỗ trợ page thì bỏ qua.
type MovieSource interface {
	Name() string
	GetTrending(page int) ([]AnimeItem, string, error)
	Search(query string) ([]AnimeItem, string, error)
	GetDetails(slug string) (*MovieInfo, []EpisodeData, error)
	GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error)
}

var allMovieSources = map[string]MovieSource{
	"phimapi": NewAnimeClient(),
	"ophim":   NewOPhimClient(),
	"nguonc":  NewNguonCScraper(),
}

func GetMovieSourceNames() []string {
	return []string{"phimapi", "ophim", "nguonc"}
}

func GetMovieSource(name string) MovieSource {
	if s, ok := allMovieSources[name]; ok {
		return s
	}
	return allMovieSources["phimapi"]
}

// GetRefererForSource trả về Referer phù hợp cho CDN của từng nguồn phim.
// Player cần gửi header này để CDN cho phép stream m3u8.
func GetRefererForSource(name string) string {
	switch name {
	case "nguonc":
		return "https://phim.nguonc.com/"
	case "ophim":
		return "https://ophim1.com/"
	case "phimapi":
		return "https://phimapi.com/"
	default:
		return ""
	}
}

// ── backward compat aliases (used only inside tui package) ─────────────────
// Đã xóa MovieSourceWithPage — dùng MovieSource trực tiếp.
