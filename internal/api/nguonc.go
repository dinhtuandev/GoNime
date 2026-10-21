package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// NguonC JSON API — https://phim.nguonc.com/api/
// Docs:
//   GET /api/films/phim-moi-cap-nhat?page=N    → trending
//   GET /api/films/search?keyword=...           → search
//   GET /api/film/{slug}                        → detail + episodes

const nguoncBase = "https://phim.nguonc.com"

type NguonCScraper struct {
	client *http.Client
}

func NewNguonCScraper() *NguonCScraper {
	return &NguonCScraper{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *NguonCScraper) Name() string { return "nguonc" }

func (s *NguonCScraper) doGet(reqURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoNime/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", nguoncBase+"/")
	return s.client.Do(req)
}

// --- Response types ---

// nguoncListResponse — phim.nguonc.com/api/films/*
// API trả về items[] ở top-level, không phải paginate.data
type nguoncListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Items   []struct {
		Name         string `json:"name"`
		Slug         string `json:"slug"`
		OriginalName string `json:"original_name"`
		ThumbURL     string `json:"thumb_url"`
		PosterURL    string `json:"poster_url"`
	} `json:"items"`
}

type nguoncDetailResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Movie   struct {
		Name        string  `json:"name"`
		Slug        string  `json:"slug"`
		OriginName  string  `json:"original_name"`
		PosterURL   string  `json:"poster_url"`
		Description string  `json:"description"`
		Year        flexInt `json:"year"`
		Quality     string  `json:"quality"`
		Total       int     `json:"total_episodes"`
		CurrentEp   string  `json:"current_episode"`
		Episodes    []struct {
			ServerName string `json:"server_name"`
			Items      []struct {
				Name  string `json:"name"`
				Slug  string `json:"slug"`
				Embed string `json:"embed"`
				M3U8  string `json:"m3u8"`
			} `json:"items"`
		} `json:"episodes"`
	} `json:"movie"`
}

// --- Interface implementation ---

func (s *NguonCScraper) GetTrending(page int) ([]AnimeItem, string, error) {
	if page < 1 {
		page = 1
	}
	reqURL := fmt.Sprintf("%s/api/films/phim-moi-cap-nhat?page=%d", nguoncBase, page)
	resp, err := s.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("nguonc trending: HTTP %s", resp.Status)
	}

	var result nguoncListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	var items []AnimeItem
	for _, m := range result.Items {
		items = append(items, AnimeItem{
			Name:       m.Name,
			OriginName: m.OriginalName,
			Slug:       m.Slug,
			PosterURL:  m.PosterURL,
		})
	}
	return items, "", nil
}

func (s *NguonCScraper) Search(query string) ([]AnimeItem, string, error) {
	encoded := url.QueryEscape(query)
	reqURL := fmt.Sprintf("%s/api/films/search?keyword=%s", nguoncBase, encoded)
	resp, err := s.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("nguonc search: HTTP %s", resp.Status)
	}

	var result nguoncListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	var items []AnimeItem
	for _, m := range result.Items {
		items = append(items, AnimeItem{
			Name:       m.Name,
			OriginName: m.OriginalName,
			Slug:       m.Slug,
			PosterURL:  m.PosterURL,
		})
	}
	return items, "", nil
}

func (s *NguonCScraper) GetDetails(slug string) (*MovieInfo, []EpisodeData, error) {
	reqURL := fmt.Sprintf("%s/api/film/%s", nguoncBase, slug)
	resp, err := s.doGet(reqURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("nguonc detail: HTTP %s", resp.Status)
	}

	var result nguoncDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	m := result.Movie
	year := int(m.Year)

	info := &MovieInfo{
		Name:        m.Name,
		Slug:        m.Slug,
		OriginName:  m.OriginName,
		PosterURL:   m.PosterURL,
		Description: m.Description,
		Year:        year,
		Quality:     m.Quality,
		Episode:     m.CurrentEp,
	}

	// Lấy server đầu tiên có dữ liệu, ưu tiên m3u8
	var episodes []EpisodeData
	for _, server := range m.Episodes {
		for _, ep := range server.Items {
			link := ep.M3U8
			if link == "" {
				link = ep.Embed
			}
			episodes = append(episodes, EpisodeData{
				Name:      ep.Name,
				Slug:      ep.Slug,
				LinkM3U8:  link,
				LinkEmbed: ep.Embed,
			})
		}
		if len(episodes) > 0 {
			break // Lấy server đầu tiên là đủ
		}
	}

	return info, episodes, nil
}

func (s *NguonCScraper) GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error) {
	return s.GetDetails(slug)
}

// flexInt chấp nhận year dạng số, chuỗi số, rỗng hoặc null từ API NguonC.
type flexInt int

func (f *flexInt) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		*f = 0
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		*f = 0
		return nil
	}
	*f = flexInt(v)
	return nil
}
