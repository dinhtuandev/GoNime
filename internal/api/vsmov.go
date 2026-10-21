package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// VSMOV JSON API — https://vsmov.com/api
//
// Endpoints:
//
//	GET /danh-sach/phim-moi-cap-nhat?page=N  → trending
//	GET /tim-kiem?keyword=...&limit=N        → search
//	GET /phim/{slug}                         → detail + episodes
//
// Danh sách trả `{status, items[], pagination{currentPage,totalPages}}`.
// Chi tiết trả `{movie{...}, episodes[]}` giống phimapi.
const vsmovBaseURL = "https://vsmov.com/api"

type VSMOVClient struct {
	HTTPClient *http.Client
}

func NewVSMOVClient() *VSMOVClient {
	return &VSMOVClient{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *VSMOVClient) Name() string { return "vsmov" }

// jsonText chấp nhận string, null, object hoặc array (VSMOV trả poster_url
// dạng {} khi phim không có poster).
type jsonText string

func (t *jsonText) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		*t = jsonText(s)
		return nil
	}
	*t = ""
	return nil
}

func (t jsonText) String() string { return string(t) }

type vsmovListItem struct {
	Name           string   `json:"name"`
	OriginName     string   `json:"origin_name"`
	Slug           string   `json:"slug"`
	PosterURL      jsonText `json:"poster_url"`
	ThumbURL       jsonText `json:"thumb_url"`
	Year           int      `json:"year"`
	EpisodeCurrent jsonText `json:"episode_current"`
	Quality        string   `json:"quality"`
	Lang           string   `json:"lang"`
	Content        jsonText `json:"content"`
}

type vsmovListResponse struct {
	Status bool            `json:"status"`
	Items  []vsmovListItem `json:"items"`
	Paging *struct {
		CurrentPage int `json:"currentPage"`
		TotalPages  int `json:"totalPages"`
	} `json:"pagination"`
}

type vsmovDetailResponse struct {
	Status bool          `json:"status"`
	Movie  vsmovListItem `json:"movie"`
	// Episodes dùng chung shape với phimapi: server_name + server_data.
	Episodes []EpisodeServer `json:"episodes"`
}

func (c *VSMOVClient) getJSON(reqURL string, out interface{}) error {
	resp, err := c.HTTPClient.Get(reqURL) //nolint:noctx
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vsmov: HTTP %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *VSMOVClient) GetTrending(page int) ([]AnimeItem, string, error) {
	if page < 1 {
		page = 1
	}
	reqURL := fmt.Sprintf("%s/danh-sach/phim-moi-cap-nhat?page=%d", vsmovBaseURL, page)

	var result vsmovListResponse
	if err := c.getJSON(reqURL, &result); err != nil {
		return nil, "", err
	}
	return vsmovItems(result.Items), "", nil
}

func (c *VSMOVClient) Search(query string) ([]AnimeItem, string, error) {
	reqURL := fmt.Sprintf("%s/tim-kiem?keyword=%s&limit=40", vsmovBaseURL, url.QueryEscape(query))

	var result vsmovListResponse
	if err := c.getJSON(reqURL, &result); err != nil {
		return nil, "", err
	}
	return vsmovItems(result.Items), "", nil
}

func vsmovItems(items []vsmovListItem) []AnimeItem {
	out := make([]AnimeItem, 0, len(items))
	for _, m := range items {
		poster := m.PosterURL.String()
		if poster == "" {
			poster = m.ThumbURL.String()
		}
		out = append(out, AnimeItem{
			Name:       m.Name,
			OriginName: m.OriginName,
			Slug:       m.Slug,
			Year:       m.Year,
			PosterURL:  poster,
		})
	}
	return out
}

func (c *VSMOVClient) GetDetails(slug string) (*MovieInfo, []EpisodeData, error) {
	reqURL := fmt.Sprintf("%s/phim/%s", vsmovBaseURL, url.PathEscape(slug))

	var result vsmovDetailResponse
	if err := c.getJSON(reqURL, &result); err != nil {
		return nil, nil, err
	}
	if result.Movie.Slug == "" {
		return nil, nil, fmt.Errorf("vsmov: không tìm thấy phim %q", slug)
	}

	m := result.Movie
	poster := m.PosterURL.String()
	thumb := m.ThumbURL.String()

	movie := &MovieInfo{
		Name:        m.Name,
		Slug:        m.Slug,
		OriginName:  m.OriginName,
		PosterURL:   poster,
		ThumbURL:    thumb,
		Year:        m.Year,
		Episode:     m.EpisodeCurrent.String(),
		Quality:     m.Quality,
		Lang:        m.Lang,
		Description: m.Content.String(),
	}

	var episodes []EpisodeData
	for _, server := range result.Episodes {
		for _, ep := range server.ServerData {
			link := ep.LinkM3U8
			if link == "" {
				link = ep.LinkEmbed
			}
			episodes = append(episodes, EpisodeData{
				Name:      ep.Name,
				Slug:      ep.Slug,
				LinkM3U8:  ep.LinkM3U8,
				LinkEmbed: link,
			})
		}
	}
	return movie, episodes, nil
}

func (c *VSMOVClient) GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error) {
	return c.GetDetails(slug)
}
