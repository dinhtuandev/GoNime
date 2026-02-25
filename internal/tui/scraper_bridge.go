package tui

import "gonime/internal/scraper"

// scraper_iface là interface tối thiểu để FavoritesModel gọi scraper
// mà không bị circular dependency
type scraper_iface interface {
	GetDetails(slug string) (*scraper.Manga, error)
}

type mangaScraperAdapter struct {
	s scraper.MangaSource
}

func (a *mangaScraperAdapter) GetDetails(slug string) (*scraper.Manga, error) {
	return a.s.GetDetails(slug)
}

func newScraper() scraper_iface {
	return &mangaScraperAdapter{
		s: scraper.GetMangaSource("truyenvn"),
	}
}
