package api

type MovieSource interface {
	Name() string
	GetTrending() ([]AnimeItem, string, error)
	Search(query string) ([]AnimeItem, string, error)
	GetDetails(slug string) (*MovieInfo, []EpisodeData, error)
	GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error)
}

type MovieSourceWithPage interface {
	Name() string
	GetTrending(page int) ([]AnimeItem, string, error)
	Search(query string) ([]AnimeItem, string, error)
	GetDetails(slug string) (*MovieInfo, []EpisodeData, error)
	GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error)
}

var AllMovieSources = map[string]MovieSource{
	"ophim":  NewOPhimClient(),
	"nguonc": NewNguonCScraper(),
}

var AllMovieSourcesWithPage = map[string]MovieSourceWithPage{
	"phimapi": NewAnimeClient(),
}

func GetMovieSourceNames() []string {
	return []string{"phimapi", "ophim", "nguonc"}
}

func GetMovieSource(name string) MovieSource {
	if s, ok := AllMovieSources[name]; ok {
		return s
	}
	return AllMovieSources["ophim"]
}

func GetMovieSourceWithPage(name string) MovieSourceWithPage {
	if s, ok := AllMovieSourcesWithPage[name]; ok {
		return s
	}
	return AllMovieSourcesWithPage["phimapi"]
}
