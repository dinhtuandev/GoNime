package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OPhim sử dụng cùng API structure với PhimAPI (đã xác nhận từ ophim1.com).
// Domain: https://ophim1.com/v1/api

const ophimBaseURL = "https://ophim1.com/v1/api"
const ophimCDN = "https://img.ophim.live/uploads/movies/"

type OPhimClient struct {
	HTTPClient *http.Client
}

func NewOPhimClient() *OPhimClient {
	return &OPhimClient{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *OPhimClient) Name() string { return "ophim" }

func (c *OPhimClient) doGet(reqURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoNime/1.0")
	req.Header.Set("Accept", "application/json")
	return c.HTTPClient.Do(req)
}

// ophimListResponse là response chung cho list/trending/search
type ophimListResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items []struct {
			Name       string `json:"name"`
			Slug       string `json:"slug"`
			OriginName string `json:"origin_name"`
			ThumbURL   string `json:"thumb_url"`
			PosterURL  string `json:"poster_url"`
			Year       int    `json:"year"`
		} `json:"items"`
		AppCDN string `json:"APP_DOMAIN_CDN_IMAGE"`
	} `json:"data"`
}

func (c *OPhimClient) GetTrending(page int) ([]AnimeItem, string, error) {
	if page < 1 {
		page = 1
	}
	reqURL := fmt.Sprintf("%s/danh-sach/phim-moi-cap-nhat?page=%d", ophimBaseURL, page)
	resp, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("ophim trending: HTTP " + resp.Status)
	}

	var result ophimListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	cdn := result.Data.AppCDN
	if cdn == "" {
		cdn = ophimCDN
	}
	var items []AnimeItem
	for _, m := range result.Data.Items {
		items = append(items, AnimeItem{
			Name:       m.Name,
			OriginName: m.OriginName,
			Slug:       m.Slug,
			PosterURL:  cdn + m.PosterURL,
			Year:       m.Year,
		})
	}
	return items, cdn, nil
}

func (c *OPhimClient) Search(query string) ([]AnimeItem, string, error) {
	encoded := url.QueryEscape(query)
	reqURL := fmt.Sprintf("%s/tim-kiem?keyword=%s&limit=40", ophimBaseURL, encoded)
	resp, err := c.doGet(reqURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("ophim search: HTTP " + resp.Status)
	}

	var result ophimListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	cdn := result.Data.AppCDN
	if cdn == "" {
		cdn = ophimCDN
	}
	var items []AnimeItem
	for _, m := range result.Data.Items {
		items = append(items, AnimeItem{
			Name:       m.Name,
			OriginName: m.OriginName,
			Slug:       m.Slug,
			PosterURL:  cdn + m.PosterURL,
			Year:       m.Year,
		})
	}
	return items, cdn, nil
}

// ophimDetailResponse là response chi tiết 1 phim
type ophimDetailResponse struct {
	Status string `json:"status"`
	Movie  struct {
		Name           string `json:"name"`
		Slug           string `json:"slug"`
		OriginName     string `json:"origin_name"`
		PosterURL      string `json:"poster_url"`
		ThumbURL       string `json:"thumb_url"`
		Year           int    `json:"year"`
		Description    string `json:"content"`
		EpisodeCurrent string `json:"episode_current"`
		Quality        string `json:"quality"`
		Lang           string `json:"lang"`
	} `json:"movie"`
	Episodes []struct {
		ServerName string `json:"server_name"`
		ServerData []struct {
			Name      string `json:"name"`
			Slug      string `json:"slug"`
			LinkM3U8  string `json:"link_m3u8"`
			LinkEmbed string `json:"link_embed"`
		} `json:"server_data"`
	} `json:"episodes"`
}

func (c *OPhimClient) GetDetails(slug string) (*MovieInfo, []EpisodeData, error) {
	reqURL := fmt.Sprintf("%s/phim/%s", ophimBaseURL, slug)
	resp, err := c.doGet(reqURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("ophim detail: HTTP " + resp.Status)
	}

	var result ophimDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	info := &MovieInfo{
		Name:        result.Movie.Name,
		Slug:        result.Movie.Slug,
		OriginName:  result.Movie.OriginName,
		PosterURL:   ophimCDN + result.Movie.PosterURL,
		ThumbURL:    ophimCDN + result.Movie.ThumbURL,
		Year:        result.Movie.Year,
		Description: result.Movie.Description,
		Episode:     result.Movie.EpisodeCurrent,
		Quality:     result.Movie.Quality,
		Lang:        result.Movie.Lang,
	}

	// Ưu tiên server có link_m3u8 thật
	var episodes []EpisodeData
	for _, server := range result.Episodes {
		for _, ep := range server.ServerData {
			episodes = append(episodes, EpisodeData{
				Name:      ep.Name,
				Slug:      ep.Slug,
				LinkM3U8:  ep.LinkM3U8,
				LinkEmbed: ep.LinkEmbed,
			})
		}
		// Lấy server đầu tiên có dữ liệu là đủ
		if len(episodes) > 0 {
			break
		}
	}

	// Nếu không có server nào có m3u8, fallback lấy tất cả (có thể toàn embed)
	if len(episodes) == 0 {
		for _, server := range result.Episodes {
			for _, ep := range server.ServerData {
				link := ep.LinkEmbed
				if ep.LinkM3U8 != "" {
					link = ep.LinkM3U8
				}
				episodes = append(episodes, EpisodeData{
					Name:      ep.Name,
					Slug:      ep.Slug,
					LinkM3U8:  link,
					LinkEmbed: ep.LinkEmbed,
				})
			}
		}
	}

	return info, episodes, nil
}

func (c *OPhimClient) GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error) {
	return c.GetDetails(slug)
}

// GetTrending_old kept for reference — removed, use GetTrending(page int).
// Search attempts to find the best matching episodes by preferring m3u8 servers.
func pickBestLink(m3u8, embed string) string {
	if m3u8 != "" && strings.HasPrefix(m3u8, "http") {
		return m3u8
	}
	return embed
}
