package scraper

import (
	"context"
	"encoding/json"
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
	otruyenDefaultBaseURL = "https://otruyen.cc"
)

// OTruyenScraper implements MangaSource bằng cách đọc HTML từ otruyen.cc,
// tương tự như MangaScraper cho truyenvn.
type OTruyenScraper struct {
	BaseURL string
	Client  *http.Client
}

func NewOTruyenScraper(baseURL string) *OTruyenScraper {
	if baseURL == "" {
		baseURL = otruyenDefaultBaseURL
	}
	return &OTruyenScraper{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *OTruyenScraper) Name() string {
	return "OTruyen"
}

func (s *OTruyenScraper) doRequest(reqURL string) (*goquery.Document, error) {
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

// Search: OTruyen không có trang search HTML riêng, nên ta dựa
// trên endpoint API /tim-kiem nhưng chỉ để lấy danh sách slug,
// sau đó sử dụng luồng HTML cho details/chapters.
func (s *OTruyenScraper) Search(query string) ([]Manga, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, fmt.Errorf("từ khóa tìm kiếm không được để trống")
	}

	apiURL := "https://otruyenapi.com/v1/api/tim-kiem"
	params := url.Values{}
	params.Set("keyword", q)

	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) GoNime/1.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối tới API OTruyen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API OTruyen trả về mã lỗi %d", resp.StatusCode)
	}

	var result otruyenListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("không parse được JSON từ API OTruyen: %w", err)
	}

	var mangas []Manga
	for _, item := range result.Data.Items {
		m := Manga{
			Name: strings.TrimSpace(item.Name),
			Slug: strings.TrimSpace(item.Slug),
		}
		mangas = append(mangas, m)
	}
	return mangas, nil
}

func (s *OTruyenScraper) GetTrending() ([]Manga, error) {
	apiURL := "https://otruyenapi.com/v1/api/home"

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) GoNime/1.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối tới API OTruyen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API OTruyen trả về mã lỗi %d", resp.StatusCode)
	}

	var result otruyenListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("không parse được JSON từ API OTruyen: %w", err)
	}

	var mangas []Manga
	for _, item := range result.Data.Items {
		m := Manga{
			Name: strings.TrimSpace(item.Name),
			Slug: strings.TrimSpace(item.Slug),
		}
		if item.ThumbURL != "" && result.Data.AppDomainCdnImage != "" {
			m.Cover = fmt.Sprintf("%s/uploads/comics/%s", result.Data.AppDomainCdnImage, item.ThumbURL)
		}
		mangas = append(mangas, m)
		if len(mangas) >= 20 {
			break
		}
	}
	return mangas, nil
}

func (s *OTruyenScraper) GetDetails(slug string) (*Manga, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("slug không hợp lệ")
	}

	encodedSlug := url.PathEscape(slug)
	reqURL := s.BaseURL + "/truyen-tranh/" + encodedSlug

	doc, err := s.doRequest(reqURL)
	if err != nil {
		return nil, err
	}

	manga := &Manga{
		Name: strings.TrimSpace(doc.Find("h1").First().Text()),
		Slug: slug,
	}

	// Description: đoạn văn ngay dưới tiêu đề.
	doc.Find("p").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		if strings.Contains(strings.ToLower(text), "truyện tranh") && manga.Description == "" {
			manga.Description = text
		}
	})

	// Trạng thái, tác giả, thể loại: dựa vào các đoạn có prefix tương ứng.
	doc.Find("p").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		lower := strings.ToLower(text)
		switch {
		case strings.HasPrefix(lower, "trạng thái:"):
			manga.Status = strings.TrimSpace(strings.TrimPrefix(text, "Trạng thái:"))
		case strings.HasPrefix(lower, "tác giả:"):
			manga.Author = strings.TrimSpace(strings.TrimPrefix(text, "Tác giả:"))
		case strings.HasPrefix(lower, "thể loại:"):
			parts := strings.Split(strings.TrimSpace(strings.TrimPrefix(text, "Thể loại:")), ",")
			for _, p := range parts {
				if g := strings.TrimSpace(p); g != "" {
					manga.Genres = append(manga.Genres, g)
				}
			}
		}
	})

	// Cover: ưu tiên img có src/thumb_url trong nội dung.
	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		if manga.Cover != "" {
			return
		}
		if src, ok := sel.Attr("src"); ok && strings.Contains(src, "otruyencdn.com") {
			manga.Cover = strings.TrimSpace(src)
		}
	})

	// Chapters: trang OTruyen hiển thị dạng text:
	// "12|https://sv1.otruyencdn.com/v1/api/chapter/...| ...".
	// HTML không có thẻ <a>, nên dùng regex trên toàn bộ source.
	if len(manga.Chapters) == 0 {
		resp, err := s.Client.Get(reqURL)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				if body, err := io.ReadAll(resp.Body); err == nil {
					text := string(body)
					re := regexp.MustCompile(`(\d+)\|(https://[^|]+/v1/api/chapter/[^|]+)\|`)
					matches := re.FindAllStringSubmatch(text, -1)
					for _, msub := range matches {
						if len(msub) < 3 {
							continue
						}
						number := strings.TrimSpace(msub[1])
						href := strings.TrimSpace(msub[2])
						if href == "" {
							continue
						}
						title := "Chapter " + number
						manga.Chapters = append(manga.Chapters, Chapter{
							Number: number,
							Title:  title,
							URL:    href,
						})
					}
				}
			}
		}
	}

	return manga, nil
}

func (s *OTruyenScraper) GetChapterImages(chapterURL string) ([]string, error) {
	chapterURL = strings.TrimSpace(chapterURL)
	if chapterURL == "" {
		return nil, fmt.Errorf("URL chapter không hợp lệ")
	}

	var resp otruyenChapterResponse
	if err := s.getJSONAbsolute(chapterURL, &resp); err != nil {
		return nil, err
	}

	base := strings.TrimRight(resp.Data.DomainCDN, "/")
	path := strings.Trim(resp.Data.Item.ChapterPath, "/")

	var images []string
	for _, img := range resp.Data.Item.ChapterImage {
		file := strings.TrimLeft(img.ImageFile, "/")
		if file == "" {
			continue
		}
		u := fmt.Sprintf("%s/%s/%s", base, path, file)
		images = append(images, u)
	}

	return images, nil
}

func (s *OTruyenScraper) DownloadImage(ctx context.Context, imageURL string, destDir string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", err
	}
	// Use a real browser User-Agent — CDN blocks non-browser agents
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://otruyen.cc/")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: %s", resp.Status)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
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

	_, err = out.ReadFrom(resp.Body)
	return destPath, err
}

// --- JSON helper chỉ dùng cho chapter API CDN ---

func (s *OTruyenScraper) getJSONAbsolute(fullURL string, v interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) GoNime/1.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("lỗi kết nối tới CDN OTruyen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CDN OTruyen trả về mã lỗi %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("không parse được JSON từ CDN OTruyen: %w", err)
	}
	return nil
}

// --- DTOs tối thiểu để parse search + chapter ---

type otruyenListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AppDomainCdnImage string `json:"APP_DOMAIN_CDN_IMAGE"`
		Items             []struct {
			Name     string `json:"name"`
			Slug     string `json:"slug"`
			ThumbURL string `json:"thumb_url"`
		} `json:"items"`
	} `json:"data"`
}

type otruyenChapterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		DomainCDN string `json:"domain_cdn"`
		Item      struct {
			ChapterPath  string `json:"chapter_path"`
			ChapterImage []struct {
				ImagePage int    `json:"image_page"`
				ImageFile string `json:"image_file"`
			} `json:"chapter_image"`
		} `json:"item"`
	} `json:"data"`
}
