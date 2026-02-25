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

func (s *MangaScraper) Name() string {
	return "TruyenVN"
}

func (s *MangaScraper) doRequest(reqURL string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", s.BaseURL+"/")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lỗi HTTP %d khi truy cập %s", resp.StatusCode, reqURL)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

// Search tìm manga theo query — thử nhiều selector để tương thích với thay đổi HTML
func (s *MangaScraper) Search(query string) ([]Manga, error) {
	encodedQuery := url.QueryEscape(query)
	reqURL := s.BaseURL + "/?s=" + encodedQuery

	doc, err := s.doRequest(reqURL)
	if err != nil {
		return nil, err
	}

	var mangas []Manga
	seen := make(map[string]bool)

	// Danh sách selector ưu tiên — thử từ trên xuống
	selectorChains := []string{
		// Madara theme cũ
		".c-tabs-item__content .post-title a",
		// Madara theme mới — search results
		".page-item-detail .post-title a",
		".row.c-tabs-item__content .post-title a",
		// Madara generic
		".post-title h3 a",
		".post-title a",
		// Fallback rộng: mọi link dẫn đến /truyen-tranh/
		"a[href*='/truyen-tranh/']",
	}

	for _, sel := range selectorChains {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			// Lấy title: ưu tiên attr title, rồi text nội dung
			title := strings.TrimSpace(s.AttrOr("title", ""))
			if title == "" {
				title = strings.TrimSpace(s.Text())
			}
			slug := extractSlug(link)

			// Bỏ qua link trống, title trống, hoặc slug chung chung
			if slug == "" || title == "" || seen[slug] {
				return
			}
			// Bỏ qua nếu slug quá ngắn (thường là thẻ page/category)
			if len(slug) < 3 {
				return
			}
			seen[slug] = true
			mangas = append(mangas, Manga{
				Name: title,
				Slug: slug,
			})
		})

		// Nếu đã tìm thấy kết quả, dừng
		if len(mangas) > 0 {
			break
		}
	}

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

	doc.Find(".summary__content, .description-summary .summary__content").Each(func(i int, sel *goquery.Selection) {
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

	// Thử nhiều selector cho chapter list
	chapterSelectors := []string{
		".version-chap .chapter-item a",
		".chapter-list .chapter-item a",
		".wp-manga-chapter a",
		"li.wp-manga-chapter a",
	}
	for _, sel := range chapterSelectors {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			title := strings.TrimSpace(s.Text())
			number := extractChapterNumber(title)

			if link != "" {
				manga.Chapters = append(manga.Chapters, Chapter{
					Number: number,
					Title:  title,
					URL:    link,
				})
			}
		})
		if len(manga.Chapters) > 0 {
			break
		}
	}

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

// GetTrending lấy danh sách manga nổi bật từ trang chủ
func (s *MangaScraper) GetTrending() ([]Manga, error) {
	doc, err := s.doRequest(s.BaseURL + "/")
	if err != nil {
		return nil, err
	}

	var mangas []Manga
	seen := make(map[string]bool)

	// Selector chuẩn của Madara theme cho trang chủ/listing
	// Updated selectors for truyenvn.shop structure
	selectors := []string{
		".page-item-detail .post-title a",
		".post-title a",
		".slider__thumb_item a",
		".c-tabs-item__content .post-title a",
		".manga-title-badges a",
		"h3 a",
		".manga-item .manga-title a",
		".popular-manga-list .manga-item a",
	}

	for _, sel := range selectors {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			if link == "" || !strings.Contains(link, "/truyen-tranh/") {
				return
			}

			// Lấy title từ nhiều nguồn
			var title string
			if t, ok := s.Attr("title"); ok && t != "" {
				title = t
			} else if t, ok := s.Attr("alt"); ok && t != "" {
				title = t
			} else {
				title = strings.TrimSpace(s.Text())
			}

			slug := extractSlug(link)
			if slug == "" || title == "" || seen[slug] || len(slug) < 3 {
				return
			}

			seen[slug] = true
			mangas = append(mangas, Manga{
				Name: title,
				Slug: slug,
			})
		})
		if len(mangas) >= 20 {
			break
		}
	}

	// If still no results, try to parse from the main content structure we saw in the HTML
	if len(mangas) < 5 {
		doc.Find(".page-item-detail, .manga-item, .post-item").Each(func(i int, sel *goquery.Selection) {
			link := sel.Find("a").First()
			href, _ := link.Attr("href")
			if href == "" || !strings.Contains(href, "/truyen-tranh/") {
				return
			}

			title := strings.TrimSpace(link.AttrOr("title", ""))
			if title == "" {
				title = strings.TrimSpace(link.Text())
			}

			slug := extractSlug(href)
			if slug != "" && title != "" && !seen[slug] && len(slug) > 3 {
				seen[slug] = true
				mangas = append(mangas, Manga{
					Name: title,
					Slug: slug,
				})
			}
		})
	}

	// Limit results to 20
	if len(mangas) > 20 {
		mangas = mangas[:20]
	}

	return mangas, nil
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
