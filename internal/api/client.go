package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type Anime struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	Episodes    []Episode `json:"episodes"`
}

type Episode struct {
	Number string `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

type AnimeClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewAnimeClient(baseURL string) *AnimeClient {
	return &AnimeClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AnimeClient) Search(query string) ([]Anime, error) {
	encodedQuery := url.QueryEscape(query)
	reqURL := c.BaseURL + "/search?q=" + encodedQuery

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoNime/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to search anime: " + resp.Status)
	}

	var result []Anime
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *AnimeClient) GetDetails(slug string) (*Anime, error) {
	encodedSlug := url.PathEscape(slug)
	reqURL := c.BaseURL + "/anime/" + encodedSlug

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoNime/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get anime details: " + resp.Status)
	}

	var anime Anime
	if err := json.NewDecoder(resp.Body).Decode(&anime); err != nil {
		return nil, err
	}

	return &anime, nil
}
