package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL       = "https://phimapi.com"
	SearchAPIURL  = "https://phimapi.com/v1/api/tim-kiem"
	EpisodeAPIURL = "https://phimapi.com/phim"
)

// --- Search Response Structs ---

type SearchResult struct {
	Status string     `json:"status"`
	Data   SearchData `json:"data"`
}

type SearchData struct {
	Items          []AnimeItem `json:"items"`
	CDNImageDomain string      `json:"APP_DOMAIN_CDN_IMAGE"`
}

type AnimeItem struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Year      int    `json:"year"`
	PosterURL string `json:"poster_url"`
}

// --- Episode Response Structs ---

type EpisodeResponse struct {
	Status   bool            `json:"status"`
	Movie    MovieInfo       `json:"movie"`
	Episodes []EpisodeServer `json:"episodes"`
}

type MovieInfo struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	OriginName  string `json:"origin_name"`
	Description string `json:"content"`
	PosterURL   string `json:"poster_url"`
	ThumbURL    string `json:"thumb_url"`
	Year        int    `json:"year"`
	Episode     string `json:"episode_current"`
	Time        string `json:"time"`
	Quality     string `json:"quality"`
	Lang        string `json:"lang"`
}

type EpisodeServer struct {
	ServerName string        `json:"server_name"`
	ServerData []EpisodeData `json:"server_data"`
}

type EpisodeData struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	LinkM3U8  string `json:"link_m3u8"`
	LinkEmbed string `json:"link_embed"`
}

// --- Client ---

type AnimeClient struct {
	HTTPClient *http.Client
}

func NewAnimeClient() *AnimeClient {
	return &AnimeClient{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AnimeClient) doGet(reqURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoNime/1.0.0")
	req.Header.Set("Accept", "application/json")
	return c.HTTPClient.Do(req)
}

// Search tìm kiếm anime theo keyword, trả về danh sách AnimeItem
func (c *AnimeClient) Search(query string) ([]AnimeItem, string, error) {
	encodedQuery := url.QueryEscape(query)
	reqURL := SearchAPIURL + "?keyword=" + encodedQuery + "&limit=20"

	resp, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("lỗi kết nối API: " + resp.Status)
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	if result.Status != "success" {
		return nil, "", errors.New("API trả về lỗi")
	}

	return result.Data.Items, result.Data.CDNImageDomain, nil
}

// GetEpisodes lấy chi tiết phim và danh sách tập theo slug (server đầu tiên có m3u8)
func (c *AnimeClient) GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error) {
	movie, servers, err := c.GetAllServers(slug)
	if err != nil {
		return nil, nil, err
	}
	// Ưu tiên server có link_m3u8 thật (không phải embed)
	for _, server := range servers {
		for _, ep := range server.ServerData {
			if ep.LinkM3U8 != "" {
				return movie, server.ServerData, nil
			}
		}
	}
	// Fallback: lấy server đầu tiên có data
	for _, server := range servers {
		if len(server.ServerData) > 0 {
			return movie, server.ServerData, nil
		}
	}
	return movie, nil, nil
}

// GetAllServers lấy tất cả servers/chất lượng của một phim
func (c *AnimeClient) GetAllServers(slug string) (*MovieInfo, []EpisodeServer, error) {
	encodedSlug := url.PathEscape(slug)
	reqURL := EpisodeAPIURL + "/" + encodedSlug

	resp, err := c.doGet(reqURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("không tìm thấy phim: " + resp.Status)
	}

	var result EpisodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	if !result.Status {
		return nil, nil, errors.New("không tìm thấy thông tin phim")
	}

	return &result.Movie, result.Episodes, nil
}

// --- Trending Structs ---

type TrendingResponse struct {
	Status string       `json:"status"`
	Data   TrendingData `json:"data"`
}

type TrendingData struct {
	Items          []AnimeItem `json:"items"`
	CDNImageDomain string      `json:"APP_DOMAIN_CDN_IMAGE"`
}

// GetTrending lấy danh sách anime mới nhất/xu hướng
func (c *AnimeClient) GetTrending(page int) ([]AnimeItem, string, error) {
	reqURL := fmt.Sprintf("%s/v1/api/danh-sach/phim-bo?page=%d&sort_field=modified.time&limit=20", BaseURL, page)

	resp, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("lỗi kết nối API trending: " + resp.Status)
	}

	var result TrendingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	if result.Status != "success" {
		return nil, "", errors.New("API trending trả về lỗi")
	}

	return result.Data.Items, result.Data.CDNImageDomain, nil
}
