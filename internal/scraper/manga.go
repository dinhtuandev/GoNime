package scraper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	maxImageSize = 50 * 1024 * 1024 // 50MB max image size
)

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

type MangaScraper struct {
	BaseURL string
	Client  *http.Client
}

func NewMangaScraper(baseURL string) *MangaScraper {
	return &MangaScraper{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *MangaScraper) doRequest(reqURL string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", s.BaseURL)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch: %s", resp.Status)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *MangaScraper) Search(query string) ([]Manga, error) {
	encodedQuery := url.QueryEscape(query)
	reqURL := s.BaseURL + "/?s=" + encodedQuery

	doc, err := s.doRequest(reqURL)
	if err != nil {
		return nil, err
	}

	var mangas []Manga
	doc.Find(".c-tabs-item__content .post-title a").Each(func(i int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		title := strings.TrimSpace(sel.Text())
		slug := extractSlug(link)

		if slug != "" && title != "" {
			mangas = append(mangas, Manga{
				Name: title,
				Slug: slug,
			})
		}
	})

	return mangas, nil
}

func (s *MangaScraper) GetDetails(slug string) (*Manga, error) {
	encodedSlug := url.PathEscape(slug)
	reqURL := s.BaseURL + "/truyen-tranh/" + encodedSlug + "/"

	doc, err := s.doRequest(reqURL)
	if err != nil {
		return nil, err
	}

	manga := &Manga{
		Name: strings.TrimSpace(doc.Find(".post-title").Text()),
		Slug: slug,
	}

	doc.Find(".summary__content").Each(func(i int, sel *goquery.Selection) {
		manga.Description = strings.TrimSpace(sel.Text())
	})

	doc.Find(".summary_image img").Each(func(i int, sel *goquery.Selection) {
		if cover, exists := sel.Attr("data-src"); exists && cover != "" {
			manga.Cover = cover
		} else if cover, exists := sel.Attr("src"); exists {
			manga.Cover = cover
		}
	})

	doc.Find(".author-content a").Each(func(i int, sel *goquery.Selection) {
		manga.Author = strings.TrimSpace(sel.Text())
	})

	doc.Find(".post-status .summary-content").Each(func(i int, sel *goquery.Selection) {
		if i == 0 {
			manga.Status = strings.TrimSpace(sel.Text())
		}
	})

	doc.Find(".genres-content a").Each(func(i int, sel *goquery.Selection) {
		manga.Genres = append(manga.Genres, strings.TrimSpace(sel.Text()))
	})

	doc.Find(".version-chap .chapter-item a, .chapter-list .chapter-item a, .wp-manga-chapter a").Each(func(i int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		title := strings.TrimSpace(sel.Text())
		number := extractChapterNumber(title)

		if link != "" {
			manga.Chapters = append(manga.Chapters, Chapter{
				Number: number,
				Title:  title,
				URL:    link,
			})
		}
	})

	return manga, nil
}

func (s *MangaScraper) GetChapterImages(chapterURL string) ([]string, error) {
	doc, err := s.doRequest(chapterURL)
	if err != nil {
		return nil, err
	}

	var images []string
	seen := make(map[string]bool)

	doc.Find(".reading-content img, .page-break img, .blocks-gallery-item img").Each(func(i int, sel *goquery.Selection) {
		var src string
		if dataSrc, exists := sel.Attr("data-src"); exists && dataSrc != "" {
			src = dataSrc
		} else if dataSrc, exists := sel.Attr("data-original"); exists && dataSrc != "" {
			src = dataSrc
		} else if srcAttr, exists := sel.Attr("src"); exists {
			src = srcAttr
		}

		src = strings.TrimSpace(src)
		if src != "" && strings.HasPrefix(src, "http") && !seen[src] {
			seen[src] = true
			images = append(images, src)
		}
	})

	return images, nil
}

func (s *MangaScraper) DownloadImage(ctx context.Context, imageURL string, destDir string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", s.BaseURL)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: %s", resp.Status)
	}

	if resp.ContentLength > maxImageSize {
		return "", fmt.Errorf("image too large: %d bytes (max: %d)", resp.ContentLength, maxImageSize)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	filename := filepath.Base(imageURL)
	if ext := filepath.Ext(filename); ext == "" {
		filename += ".jpg"
	}
	filename = sanitizeFilename(filename)

	destPath := filepath.Join(destDir, filename)

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	limitedReader := io.LimitReader(resp.Body, maxImageSize)
	_, err = io.Copy(out, limitedReader)
	return destPath, err
}

func (s *MangaScraper) GetNewManga() ([]Manga, error) {
	reqURL := s.BaseURL + "/truyen-tranh/"

	doc, err := s.doRequest(reqURL)
	if err != nil {
		return nil, err
	}

	var mangas []Manga
	seen := make(map[string]bool)

	doc.Find(".c-tabs-item__content .post-title a, .page-item-detail .post-title a").Each(func(i int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		title := strings.TrimSpace(sel.Text())
		slug := extractSlug(link)

		if slug != "" && title != "" && !seen[slug] {
			seen[slug] = true
			mangas = append(mangas, Manga{
				Name: title,
				Slug: slug,
			})
		}
	})

	return mangas, nil
}

func extractSlug(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		parts := strings.Split(strings.TrimSuffix(rawURL, "/"), "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] != "" {
				return parts[i]
			}
		}
		return ""
	}

	parts := strings.Split(strings.TrimSuffix(parsedURL.Path, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return ""
}

func extractChapterNumber(title string) string {
	re := regexp.MustCompile(`(?i)chapter\s*(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(title)
	if len(matches) > 1 {
		return matches[1]
	}
	return title
}

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	return re.ReplaceAllString(name, "_")
}
