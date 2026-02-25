package scraper

import (
	"context"
)

type MangaSource interface {
	Name() string
	GetTrending() ([]Manga, error)
	Search(query string) ([]Manga, error)
	GetDetails(slug string) (*Manga, error)
	GetChapterImages(chapterURL string) ([]string, error)
	DownloadImage(ctx context.Context, imageURL string, destDir string) (string, error)
}

var AllMangaSources = map[string]MangaSource{
	"truyenvn": NewMangaScraper("https://truyenvn.shop"),
	"truyenqq": NewTruyenQQScraper(),
}

func GetMangaSourceNames() []string {
	return []string{"truyenvn", "truyenqq"}
}

func GetMangaSource(name string) MangaSource {
	if s, ok := AllMangaSources[name]; ok {
		return s
	}
	return AllMangaSources["truyenvn"]
}
