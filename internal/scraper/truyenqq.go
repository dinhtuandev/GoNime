package scraper

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TruyenQQScraper struct {
	BaseURL string
	Client  *http.Client
}

func NewTruyenQQScraper() *TruyenQQScraper {
	return &TruyenQQScraper{
		BaseURL: "https://truyenqqno.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *TruyenQQScraper) Name() string {
	return "TruyenQQ"
}

func (s *TruyenQQScraper) doRequest(reqURL string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", s.BaseURL+"/")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *TruyenQQScraper) GetTrending() ([]Manga, error) {
	doc, err := s.doRequest(s.BaseURL + "/")
	if err != nil {
		return nil, err
	}

	var mangas []Manga
	seen := make(map[string]bool)

	doc.Find(".item, .story-item, .book-item").Each(func(i int, sel *goquery.Selection) {
		link := sel.Find("a").First()
		href, _ := link.Attr("href")
		slug := extractSlugFromURL(href)

		title := strings.TrimSpace(link.AttrOr("title", ""))
		if title == "" {
			title = strings.TrimSpace(link.Text())
		}
		if title == "" {
			title = strings.TrimSpace(sel.Find(".story-title, .title, .book-title").Text())
		}

		if title != "" && slug != "" && !seen[slug] {
			seen[slug] = true
			mangas = append(mangas, Manga{
				Name: title,
				Slug: slug,
			})
		}
	})

	if len(mangas) < 5 {
		doc.Find("a[href*='/truyen/']").Each(func(i int, sel *goquery.Selection) {
			href, _ := sel.Attr("href")
			slug := extractSlugFromURL(href)

			title := strings.TrimSpace(sel.AttrOr("title", ""))
			if title == "" {
				title = strings.TrimSpace(sel.Text())
			}

			if title != "" && slug != "" && !seen[slug] && len(slug) > 3 {
				seen[slug] = true
				mangas = append(mangas, Manga{
					Name: title,
					Slug: slug,
				})
			}
		})
	}

	if len(mangas) > 20 {
		mangas = mangas[:20]
	}

	return mangas, nil
}

func (s *TruyenQQScraper) Search(query string) ([]Manga, error) {
	doc, err := s.doRequest(s.BaseURL + "/")
	if err != nil {
		return nil, err
	}

	var mangas []Manga
	seen := make(map[string]bool)
	queryLower := strings.ToLower(query)

	doc.Find(".book_avatar, .item, .story-item, .book-item").Each(func(i int, sel *goquery.Selection) {
		link := sel.Find("a").First()
		href, _ := link.Attr("href")
		slug := extractSlugFromURL(href)

		title := strings.TrimSpace(link.AttrOr("title", ""))
		if title == "" {
			title = strings.TrimSpace(link.Text())
		}
		if title == "" {
			title = strings.TrimSpace(sel.Find(".book_name, .story-title, .title, .book-title").Text())
		}

		titleLower := strings.ToLower(title)
		if strings.Contains(titleLower, queryLower) && title != "" && slug != "" && !seen[slug] {
			seen[slug] = true
			mangas = append(mangas, Manga{
				Name: title,
				Slug: slug,
			})
		}

		if len(mangas) >= 20 {
			return
		}
	})

	if len(mangas) < 5 {
		doc.Find("a[href*='/truyen-tranh/']").Each(func(i int, sel *goquery.Selection) {
			href, _ := sel.Attr("href")
			slug := extractSlugFromURL(href)

			title := strings.TrimSpace(sel.AttrOr("title", ""))
			if title == "" {
				title = strings.TrimSpace(sel.Text())
			}

			titleLower := strings.ToLower(title)
			if strings.Contains(titleLower, queryLower) && title != "" && slug != "" && !seen[slug] && len(slug) > 3 {
				seen[slug] = true
				mangas = append(mangas, Manga{
					Name: title,
					Slug: slug,
				})
			}

			if len(mangas) >= 20 {
				return
			}
		})
	}

	if len(mangas) > 20 {
		mangas = mangas[:20]
	}

	return mangas, nil
}

func (s *TruyenQQScraper) GetDetails(slug string) (*Manga, error) {
	doc, err := s.doRequest(s.BaseURL + "/truyen-tranh/" + slug)
	if err != nil {
		return nil, err
	}

	manga := &Manga{
		Slug: slug,
	}

	manga.Name = strings.TrimSpace(doc.Find("h1[itemprop='name'], .book_other h1").First().Text())
	if manga.Name == "" {
		manga.Name = strings.TrimSpace(doc.Find("title").Text())
	}

	manga.Description = strings.TrimSpace(doc.Find(".story-detail-info, .detail-content, .description, .story-description, .content, .info").Text())

	doc.Find(".book_avatar img, .poster, .cover, .thumbnail img").Each(func(i int, sel *goquery.Selection) {
		if cover, exists := sel.Attr("src"); exists && cover != "" {
			manga.Cover = cover
		} else if cover, exists := sel.Attr("data-src"); exists && cover != "" {
			manga.Cover = cover
		}
	})

	doc.Find(".author, .artist, .translator").Each(func(i int, sel *goquery.Selection) {
		if manga.Author == "" {
			manga.Author = strings.TrimSpace(sel.Text())
		}
	})

	doc.Find(".status, .trang-thai").Each(func(i int, sel *goquery.Selection) {
		if manga.Status == "" {
			manga.Status = strings.TrimSpace(sel.Text())
		}
	})

	doc.Find(".list01 a, .category, .genre, .tag-list a").Each(func(i int, sel *goquery.Selection) {
		genre := strings.TrimSpace(sel.Text())
		if genre != "" {
			manga.Genres = append(manga.Genres, genre)
		}
	})

	doc.Find(".works-chapter-item .name-chap a, .chapter-list a, .episode-list a, .chapter-item a").Each(func(i int, sel *goquery.Selection) {
		link, _ := sel.Attr("href")
		title := strings.TrimSpace(sel.Text())

		if link != "" && title != "" {
			manga.Chapters = append(manga.Chapters, Chapter{
				Number: extractChapterNumber(title),
				Title:  title,
				URL:    link,
			})
		}
	})

	return manga, nil
}

func (s *TruyenQQScraper) GetChapterImages(chapterURL string) ([]string, error) {
	doc, err := s.doRequest(chapterURL)
	if err != nil {
		return nil, err
	}

	var images []string
	seen := make(map[string]bool)

	doc.Find(".chapter_content img, .chapter_content_div img").Each(func(i int, sel *goquery.Selection) {
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

	if len(images) == 0 {
		doc.Find("img").Each(func(i int, sel *goquery.Selection) {
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
	}

	return images, nil
}

func extractSlugFromURL(href string) string {
	if href == "" {
		return ""
	}
	parts := strings.Split(strings.TrimSuffix(href, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && parts[i] != "truyen" && parts[i] != "truyen-tranh" && parts[i] != "truyen-tranh" && !strings.HasPrefix(parts[i], "chap-") && !strings.Contains(parts[i], ".html") {
			return parts[i]
		}
	}
	return ""
}

const truyenQQMaxImageSize = 50 * 1024 * 1024 // 50MB

func (s *TruyenQQScraper) DownloadImage(ctx context.Context, imageURL string, destDir string) (string, error) {
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

	if resp.ContentLength > truyenQQMaxImageSize {
		return "", fmt.Errorf("image too large: %d bytes (max: %d)", resp.ContentLength, truyenQQMaxImageSize)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	filename := filepath.Base(imageURL)
	if ext := filepath.Ext(filename); ext == "" {
		filename += ".jpg"
	}

	destPath := filepath.Join(destDir, filename)
	outFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = outFile.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return destPath, nil
}
