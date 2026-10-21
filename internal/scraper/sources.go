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

const maxImageSize = 50 * 1024 * 1024 // 50MB max image size

type Manga struct {
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Cover       string    `json:"cover"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Status      string    `json:"status"`
	Genres      []string  `json:"genres"`
	Chapters    []Chapter `json:"chapters"`
}

type Chapter struct {
	Number string `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Date   string `json:"date"`
}

var AllMangaSources = map[string]MangaSource{
	"mangadex": NewMangaDexScraper(),
}

func GetMangaSourceNames() []string {
	return []string{"mangadex"}
}

func GetMangaSource(name string) MangaSource {
	if s, ok := AllMangaSources[name]; ok {
		return s
	}
	return AllMangaSources["mangadex"]
}
