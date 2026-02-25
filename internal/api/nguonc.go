package api

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NguonCScraper struct {
	BaseURL string
	Client  *http.Client
}

func NewNguonCScraper() *NguonCScraper {
	return &NguonCScraper{
		BaseURL: "https://phim.nguonc.com",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *NguonCScraper) Name() string {
	return "NguonC"
}

func (s *NguonCScraper) doRequest(reqURL string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Referer", s.BaseURL+"/")

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

func (s *NguonCScraper) GetTrending() ([]AnimeItem, string, error) {
	doc, err := s.doRequest(s.BaseURL + "/")
	if err != nil {
		return nil, "", err
	}

	var items []AnimeItem
	doc.Find(".film-list .film-item").Each(func(i int, sel *goquery.Selection) {
		link := sel.Find("a").First()
		href, _ := link.Attr("href")
		slug := extractSlug(href)

		name := strings.TrimSpace(link.Text())
		poster, _ := sel.Find("img").Attr("data-src")
		if poster == "" {
			poster, _ = sel.Find("img").Attr("src")
		}

		yearStr := strings.TrimSpace(sel.Find(".meta .year").Text())
		year := 0
		fmt.Sscanf(yearStr, "%d", &year)

		if name != "" && slug != "" {
			items = append(items, AnimeItem{
				Name:      name,
				Slug:      slug,
				PosterURL: poster,
				Year:      year,
			})
		}
	})

	if len(items) == 0 {
		doc.Find(".tab-content .film-item").Each(func(i int, sel *goquery.Selection) {
			link := sel.Find("a").First()
			href, _ := link.Attr("href")
			slug := extractSlug(href)

			name := strings.TrimSpace(link.Text())
			poster, _ := sel.Find("img").Attr("data-src")
			if poster == "" {
				poster, _ = sel.Find("img").Attr("src")
			}

			yearStr := strings.TrimSpace(sel.Find(".year").Text())
			year := 0
			fmt.Sscanf(yearStr, "%d", &year)

			if name != "" && slug != "" {
				items = append(items, AnimeItem{
					Name:      name,
					Slug:      slug,
					PosterURL: poster,
					Year:      year,
				})
			}
		})
	}

	return items, "", nil
}

func (s *NguonCScraper) Search(query string) ([]AnimeItem, string, error) {
	encodedQuery := url.QueryEscape(query)
	doc, err := s.doRequest(s.BaseURL + "/?search=" + encodedQuery)
	if err != nil {
		return nil, "", err
	}

	var items []AnimeItem
	queryLower := strings.ToLower(query)

	doc.Find(".film-list .film-item, .result-search .film-item").Each(func(i int, sel *goquery.Selection) {
		link := sel.Find("a").First()
		href, _ := link.Attr("href")
		slug := extractSlug(href)

		name := strings.TrimSpace(link.Text())
		if name == "" {
			name = strings.TrimSpace(sel.Find(".film-name").Text())
		}

		poster, _ := sel.Find("img").Attr("data-src")
		if poster == "" {
			poster, _ = sel.Find("img").Attr("src")
		}

		yearStr := strings.TrimSpace(sel.Find(".year").Text())
		year := 0
		fmt.Sscanf(yearStr, "%d", &year)

		nameLower := strings.ToLower(name)
		if (strings.Contains(nameLower, queryLower) || slug != "") && name != "" {
			items = append(items, AnimeItem{
				Name:      name,
				Slug:      slug,
				PosterURL: poster,
				Year:      year,
			})
		}

		if len(items) >= 20 {
			return
		}
	})

	return items, "", nil
}

func (s *NguonCScraper) GetDetails(slug string) (*MovieInfo, []EpisodeData, error) {
	doc, err := s.doRequest(s.BaseURL + "/" + slug)
	if err != nil {
		return nil, nil, err
	}

	name := strings.TrimSpace(doc.Find(".heading-name").Text())
	if name == "" {
		name = strings.TrimSpace(doc.Find("h1").Text())
	}

	poster, _ := doc.Find(".poster img").Attr("src")
	if poster == "" {
		poster, _ = doc.Find(".film-poster img").Attr("src")
	}

	thumb, _ := doc.Find(".film-poster img").Attr("data-src")
	if thumb == "" {
		thumb = poster
	}

	yearStr := strings.TrimSpace(doc.Find(".meta .year").Text())
	year := 0
	fmt.Sscanf(yearStr, "%d", &year)

	description := strings.TrimSpace(doc.Find(".content .desc").Text())
	if description == "" {
		description = strings.TrimSpace(doc.Find(".film-content").Text())
	}

	episodeStr := strings.TrimSpace(doc.Find(".meta .episode").Text())

	quality := strings.TrimSpace(doc.Find(".meta .quality").Text())

	info := &MovieInfo{
		Name:        name,
		Slug:        slug,
		PosterURL:   poster,
		ThumbURL:    thumb,
		Year:        year,
		Description: description,
		Episode:     episodeStr,
		Quality:     quality,
	}

	var episodes []EpisodeData
	doc.Find(".episode-list .episode-item, .episodes .episode a").Each(func(i int, sel *goquery.Selection) {
		epLink, _ := sel.Attr("href")
		epName := strings.TrimSpace(sel.Text())

		if epLink != "" && epName != "" {
			episodes = append(episodes, EpisodeData{
				Name:     epName,
				Slug:     extractSlug(epLink),
				LinkM3U8: epLink,
			})
		}
	})

	if len(episodes) == 0 {
		doc.Find(".server-list .server-item a").Each(func(i int, sel *goquery.Selection) {
			epLink, _ := sel.Attr("href")
			epName := strings.TrimSpace(sel.Text())

			if epLink != "" {
				episodes = append(episodes, EpisodeData{
					Name:     epName,
					Slug:     extractSlug(epLink),
					LinkM3U8: epLink,
				})
			}
		})
	}

	return info, episodes, nil
}

func (s *NguonCScraper) GetEpisodes(slug string) (*MovieInfo, []EpisodeData, error) {
	return s.GetDetails(slug)
}

func extractSlug(href string) string {
	if href == "" {
		return ""
	}
	parts := strings.Split(strings.TrimSuffix(href, "/"), "/")
	return parts[len(parts)-1]
}

var domainRegex = regexp.MustCompile(`https?://[^/]+`)

func extractDomain(href string) string {
	matches := domainRegex.FindStringSubmatch(href)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}
