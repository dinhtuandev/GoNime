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
	"sort"
	"strconv"
	"strings"
	"time"

	"gonime/internal/httpx"
	"gonime/internal/sanitize"
)

// MangaDex scraper — JSON API, không scrape HTML.
//
//	GET /manga?title=...                    → tìm kiếm
//	GET /manga?order[followedCount]=desc    → trending
//	GET /manga/{id}/feed                    → danh sách chapter (vi)
//	GET /at-home/server/{chapterID}         → ảnh chapter
const (
	mangadexAPIBase     = "https://api.mangadex.org"
	mangadexUploadsBase = "https://uploads.mangadex.org"
)

// MangaDex dùng UUID làm slug.
type MangaDexScraper struct {
	Client *http.Client
}

func NewMangaDexScraper() *MangaDexScraper {
	return &MangaDexScraper{
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *MangaDexScraper) Name() string { return "MangaDex" }

type mdRelationship struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		FileName string `json:"fileName"`
		Name     string `json:"name"`
	} `json:"attributes"`
}

type mdManga struct {
	ID         string `json:"id"`
	Attributes struct {
		Title       map[string]string `json:"title"`
		Description map[string]string `json:"description"`
		Status      string            `json:"status"`
		Tags        []struct {
			Attributes struct {
				Name map[string]string `json:"name"`
			} `json:"attributes"`
		} `json:"tags"`
	} `json:"attributes"`
	Relationships []mdRelationship `json:"relationships"`
}

type mdListResponse struct {
	Result string    `json:"result"`
	Data   []mdManga `json:"data"`
	Total  int       `json:"total"`
}

type mdChapter struct {
	ID         string `json:"id"`
	Attributes struct {
		Chapter     string `json:"chapter"`
		Title       string `json:"title"`
		PublishAt   string `json:"publishAt"`
		ExternalURL string `json:"externalUrl"`
	} `json:"attributes"`
}

type mdChapterList struct {
	Result string      `json:"result"`
	Data   []mdChapter `json:"data"`
	Total  int         `json:"total"`
}

func (s *MangaDexScraper) getJSON(ctx context.Context, reqURL string, out interface{}) error {
	resp, err := httpx.Get(ctx, reqURL,
		httpx.WithHeader("Accept", "application/json"),
		httpx.WithHeader("User-Agent", "GoNime/1.0 (https://github.com/dinhtuandev/GoNime)"),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mangadex: HTTP %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// pickText ưu tiên tiếng Việt, rồi tiếng Anh, rồi bất kỳ.
func pickText(m map[string]string) string {
	for _, lang := range []string{"vi", "en"} {
		if v := strings.TrimSpace(m[lang]); v != "" {
			return v
		}
	}
	for _, v := range m {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func (s *MangaDexScraper) toManga(m mdManga) Manga {
	cover := ""
	for _, rel := range m.Relationships {
		if rel.Type == "cover_art" && rel.Attributes.FileName != "" {
			cover = fmt.Sprintf("%s/covers/%s/%s.512.jpg", mangadexUploadsBase, m.ID, rel.Attributes.FileName)
			break
		}
	}
	author := ""
	for _, rel := range m.Relationships {
		if rel.Type == "author" && rel.Attributes.Name != "" {
			author = rel.Attributes.Name
			break
		}
	}
	var genres []string
	for _, tag := range m.Attributes.Tags {
		if name := pickText(tag.Attributes.Name); name != "" {
			genres = append(genres, name)
		}
	}
	return Manga{
		Name:        sanitize.Text(pickText(m.Attributes.Title)),
		Slug:        m.ID,
		Cover:       cover,
		Description: sanitize.Text(pickText(m.Attributes.Description)),
		Author:      sanitize.Text(author),
		Status:      sanitize.Text(m.Attributes.Status),
		Genres:      genres,
	}
}

func (s *MangaDexScraper) list(ctx context.Context, params url.Values) ([]Manga, error) {
	reqURL := mangadexAPIBase + "/manga?" + params.Encode()
	var result mdListResponse
	if err := s.getJSON(ctx, reqURL, &result); err != nil {
		return nil, err
	}
	out := make([]Manga, 0, len(result.Data))
	for _, m := range result.Data {
		out = append(out, s.toManga(m))
	}
	return out, nil
}

func (s *MangaDexScraper) trendingParams(limit int) url.Values {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(limit))
	params.Set("order[followedCount]", "desc")
	params.Add("includes[]", "cover_art")
	params.Add("includes[]", "author")
	params.Add("availableTranslatedLanguage[]", "vi")
	params.Add("contentRating[]", "safe")
	params.Add("contentRating[]", "suggestive")
	params.Add("contentRating[]", "erotica")
	return params
}

func (s *MangaDexScraper) GetTrending() ([]Manga, error) {
	return s.list(context.Background(), s.trendingParams(24))
}

func (s *MangaDexScraper) Search(query string) ([]Manga, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}
	params := s.trendingParams(24)
	params.Del("order[followedCount]")
	params.Set("title", query)
	params.Set("order[relevance]", "desc")
	return s.list(context.Background(), params)
}

func (s *MangaDexScraper) GetDetails(slug string) (*Manga, error) {
	var result struct {
		Data mdManga `json:"data"`
	}
	reqURL := fmt.Sprintf("%s/manga/%s?includes[]=cover_art&includes[]=author", mangadexAPIBase, url.PathEscape(slug))
	if err := s.getJSON(context.Background(), reqURL, &result); err != nil {
		return nil, err
	}
	if result.Data.ID == "" {
		return nil, fmt.Errorf("mangadex: không tìm thấy manga %q", slug)
	}

	manga := s.toManga(result.Data)
	chapters, err := s.fetchChapters(slug)
	if err != nil {
		return nil, err
	}
	manga.Chapters = chapters
	return &manga, nil
}

func (s *MangaDexScraper) fetchChapters(mangaID string) ([]Chapter, error) {
	// availableTranslatedLanguages có thể lỗi thời: thử tiếng Việt, rồi tiếng Anh,
	// cuối cùng lấy mọi ngôn ngữ thay vì trả màn chapter rỗng.
	for _, lang := range []string{"vi", "en", ""} {
		chapters, err := s.fetchChaptersLang(mangaID, lang)
		if err != nil {
			return nil, err
		}
		if len(chapters) > 0 {
			return chapters, nil
		}
	}
	return nil, nil
}

func (s *MangaDexScraper) fetchChaptersLang(mangaID, lang string) ([]Chapter, error) {
	const pageSize = 500
	const maxChapters = 2000

	var chapters []Chapter
	for offset := 0; offset < maxChapters; offset += pageSize {
		params := url.Values{}
		params.Set("limit", strconv.Itoa(pageSize))
		params.Set("offset", strconv.Itoa(offset))
		if lang != "" {
			params.Add("translatedLanguage[]", lang)
		}
		params.Add("order[chapter]", "desc")
		params.Add("includeExternalUrl", "0")

		reqURL := fmt.Sprintf("%s/manga/%s/feed?%s", mangadexAPIBase, url.PathEscape(mangaID), params.Encode())
		var result mdChapterList
		if err := s.getJSON(context.Background(), reqURL, &result); err != nil {
			return nil, err
		}

		for _, ch := range result.Data {
			if ch.Attributes.ExternalURL != "" {
				continue
			}
			label := strings.TrimSpace(ch.Attributes.Chapter)
			if label == "" {
				label = "Oneshot"
			} else if !strings.Contains(strings.ToLower(label), "chương") {
				label = "Chương " + label
			}
			if t := strings.TrimSpace(ch.Attributes.Title); t != "" {
				label += " - " + t
			}
			date := ""
			if len(ch.Attributes.PublishAt) >= 10 {
				date = ch.Attributes.PublishAt[:10]
			}
			chapters = append(chapters, Chapter{
				Number: ch.Attributes.Chapter,
				Title:  label,
				URL:    fmt.Sprintf("https://mangadex.org/chapter/%s", ch.ID),
				Date:   date,
			})
		}

		if offset+len(result.Data) >= result.Total || len(result.Data) == 0 {
			break
		}
	}

	// Feed trả mới nhất trước; đọc truyện thì chapter cũ trước.
	sort.SliceStable(chapters, func(i, j int) bool {
		a, errA := strconv.ParseFloat(chapters[i].Number, 64)
		b, errB := strconv.ParseFloat(chapters[j].Number, 64)
		if errA == nil && errB == nil {
			return a < b
		}
		return chapters[i].Number < chapters[j].Number
	})
	return chapters, nil
}

func (s *MangaDexScraper) GetChapterImages(chapterURL string) ([]string, error) {
	chapterID := chapterIDFromURL(chapterURL)

	var result struct {
		BaseURL string `json:"baseUrl"`
		Chapter struct {
			Hash      string   `json:"hash"`
			Data      []string `json:"data"`
			DataSaver []string `json:"dataSaver"`
		} `json:"chapter"`
	}
	reqURL := fmt.Sprintf("%s/at-home/server/%s", mangadexAPIBase, url.PathEscape(chapterID))
	if err := s.getJSON(context.Background(), reqURL, &result); err != nil {
		return nil, err
	}

	files := result.Chapter.DataSaver
	quality := "data-saver"
	if len(files) == 0 {
		files = result.Chapter.Data
		quality = "data"
	}
	if len(files) == 0 || result.Chapter.Hash == "" {
		return nil, fmt.Errorf("mangadex: chapter %q không có ảnh", chapterID)
	}

	images := make([]string, 0, len(files))
	for _, file := range files {
		images = append(images, fmt.Sprintf("%s/%s/%s/%s", result.BaseURL, quality, result.Chapter.Hash, file))
	}
	return images, nil
}

// chapterIDFromURL lấy UUID chapter từ URL web hoặc chuỗi UUID trần.
func chapterIDFromURL(chapterURL string) string {
	id := strings.TrimSpace(chapterURL)
	if i := strings.IndexAny(id, "?#"); i >= 0 {
		id = id[:i]
	}
	id = strings.TrimSuffix(id, "/")
	if i := strings.LastIndex(id, "/"); i >= 0 {
		id = id[i+1:]
	}
	return id
}

func (s *MangaDexScraper) DownloadImage(ctx context.Context, imageURL string, destDir string) (string, error) {
	p := imageURL
	if i := strings.IndexAny(p, "?#"); i >= 0 {
		p = p[:i]
	}
	destPath := filepath.Join(destDir, filepath.Base(p))

	if _, err := os.Stat(destPath); err == nil {
		return destPath, nil
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "GoNime/1.0 (https://github.com/dinhtuandev/GoNime)")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mangadex image: HTTP %s", resp.Status)
	}
	if resp.ContentLength > maxImageSize {
		return "", fmt.Errorf("image too large: %d bytes (max: %d)", resp.ContentLength, maxImageSize)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, io.LimitReader(resp.Body, maxImageSize)); err != nil {
		os.Remove(destPath)
		return "", err
	}
	return destPath, nil
}
