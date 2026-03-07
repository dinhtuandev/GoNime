package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OPhim API Client - https://ophim.cc
type OPhimClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

type OPhimResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Home    OPhimHome  `json:"home"`
	Movie   OPhimMovie `json:"movie"`
}

type OPhimHome struct {
	SEO      OPhimSEO     `json:"seo"`
	Featured []OPhimMovie `json:"featured"`
	Trending []OPhimMovie `json:"trending"`
	Latest   []OPhimMovie `json:"latest"`
}

type OPhimSEO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type OPhimMovie struct {
	ID             string          `json:"_id"`
	Name           string          `json:"name"`
	Slug           string          `json:"slug"`
	OriginName     string          `json:"origin_name"`
	PosterURL      string          `json:"poster_url"`
	ThumbURL       string          `json:"thumb_url"`
	Year           int             `json:"year"`
	Categories     []OPhimCategory `json:"category"`
	Status         string          `json:"status"`
	Time           string          `json:"time"`
	Description    string          `json:"description"`
	EpisodeCurrent string          `json:"episode_current"`
	EpisodeTotal   string          `json:"episode_total"`
	Quality        string          `json:"quality"`
}

type OPhimCategory struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type OPhimDetailResponse struct {
	Status   int              `json:"status"`
	Message  string           `json:"message"`
	Movie    OPhimMovieDetail `json:"movie"`
	Episodes []OPhimServer    `json:"episodes"`
}

type OPhimMovieDetail struct {
	Name           string          `json:"name"`
	Slug           string          `json:"slug"`
	OriginName     string          `json:"origin_name"`
	PosterURL      string          `json:"poster_url"`
	ThumbURL       string          `json:"thumb_url"`
	Year           int             `json:"year"`
	Categories     []OPhimCategory `json:"category"`
	Actors         []string        `json:"actor"`
	Directors      []string        `json:"director"`
	Duration       string          `json:"duration"`
	Quality        string          `json:"quality"`
	Status         string          `json:"status"`
	Description    string          `json:"description"`
	EpisodeCurrent string          `json:"episode_current"`
}

type OPhimServer struct {
	ServerName string         `json:"server_name"`
	ServerID   int            `json:"server_id"`
	Items      []OPhimEpisode `json:"items"`
}

type OPhimEpisode struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Link string `json:"link"`
}

func NewOPhimClient() *OPhimClient {
	return &OPhimClient{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: "https://ophim.cc",
	}
}

func (c *OPhimClient) Name() string {
	return "OPhim"
}

func (c *OPhimClient) GetTrending() ([]AnimeItem, string, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var result OPhimResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	var animes []AnimeItem
	for _, m := range result.Home.Trending {
		animes = append(animes, AnimeItem{
			Name:       m.Name,
			OriginName: m.OriginName,
			Slug:       m.Slug,
			PosterURL:  m.PosterURL,
			Year:       m.Year,
		})
	}

	return animes, "", nil
}

func (c *OPhimClient) Search(query string) ([]AnimeItem, string, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/")
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var result OPhimResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	query = strings.ToLower(query)
	var animes []AnimeItem

	// Search in all movies from home
	allMovies := append(result.Home.Trending, result.Home.Latest...)
	for _, m := range allMovies {
		if strings.Contains(strings.ToLower(m.Name), query) {
			animes = append(animes, AnimeItem{
				Name:       m.Name,
				OriginName: m.OriginName,
				Slug:       m.Slug,
				PosterURL:  m.PosterURL,
				Year:       m.Year,
			})
		}
		if len(animes) >= 20 {
			break
		}
	}

	return animes, "", nil
}

func (c *OPhimClient) GetDetails(slug string) (*MovieInfo, []EpisodeData, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/phim/%s", c.BaseURL, slug))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var result OPhimDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	detail := &MovieInfo{
		Name:        result.Movie.Name,
		Slug:        result.Movie.Slug,
		OriginName:  result.Movie.OriginName,
		PosterURL:   result.Movie.PosterURL,
		ThumbURL:    result.Movie.ThumbURL,
		Year:        result.Movie.Year,
		Description: result.Movie.Description,
		Episode:     result.Movie.EpisodeCurrent,
		Quality:     result.Movie.Quality,
	}

	var episodes []EpisodeData
	for _, server := range result.Episodes {
		for _, ep := range server.Items {
			episodes = append(episodes, EpisodeData{
				Name:     ep.Name,
				Slug:     ep.Slug,
				LinkM3U8: ep.Link,
			})
		}
	}

	return detail, episodes, nil
}

func (c *OPhimClient) GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error) {
	detail, episodes, err := c.GetDetails(slug)
	if err != nil {
		return nil, nil, err
	}

	movieInfo := &MovieInfo{
		Name:    detail.Name,
		Slug:    detail.Slug,
		Episode: detail.Episode,
	}

	return movieInfo, episodes, nil
}
