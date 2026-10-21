//go:build integration

package api

import (
	"testing"
	"time"
)

func TestAnimeTrending(t *testing.T) {
	t.Log("=== Testing Anime Trending ===")

	client := NewAnimeClient()

	t.Log("Testing GetTrending...")
	start := time.Now()
	animes, cdnURL, err := client.GetTrending(1)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("GetTrending failed: %v", err)
	} else {
		t.Logf("GetTrending succeeded in %v", duration)
		t.Logf("Found %d trending anime", len(animes))
		t.Logf("CDN URL: %s", cdnURL)

		if len(animes) == 0 {
			t.Log("Warning: No trending anime found")
		}

		for i, anime := range animes {
			if i >= 3 {
				break
			}
			t.Logf("Trending %d: %s (slug: %s)", i+1, anime.Name, anime.Slug)
			if anime.Name == "" || anime.Slug == "" {
				t.Errorf("Trending anime %d has empty fields: %+v", i+1, anime)
			}
		}
	}
}

func TestAnimeSearch(t *testing.T) {
	t.Log("=== Testing Anime Search ===")

	client := NewAnimeClient()

	t.Log("Testing Search with 'naruto'...")
	start := time.Now()
	animes, cdnURL, err := client.Search("naruto")
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Search failed: %v", err)
	} else {
		t.Logf("Search succeeded in %v", duration)
		t.Logf("Found %d anime results", len(animes))
		t.Logf("CDN URL: %s", cdnURL)

		for i, anime := range animes {
			if i >= 3 {
				break
			}
			t.Logf("Search %d: %s (slug: %s)", i+1, anime.Name, anime.Slug)
			if anime.Name == "" || anime.Slug == "" {
				t.Errorf("Search result %d has empty fields: %+v", i+1, anime)
			}
		}
	}
}

func TestAnimeEpisodes(t *testing.T) {
	t.Log("=== Testing Anime Episodes ===")

	client := NewAnimeClient()

	// Get some anime first
	animes, _, err := client.Search("naruto")
	if err != nil || len(animes) == 0 {
		t.Logf("Cannot get anime for episodes test: %v", err)
		return
	}

	t.Log("Testing GetEpisodes for first anime...")
	start := time.Now()
	movieInfo, episodes, err := client.GetEpisodes(animes[0].Slug)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("GetEpisodes failed: %v", err)
	} else {
		t.Logf("GetEpisodes succeeded in %v", duration)
		t.Logf("Anime: %s", movieInfo.Name)
		t.Logf("Episodes: %d", len(episodes))

		if movieInfo.Name == "" {
			t.Errorf("MovieInfo has empty name: %+v", movieInfo)
		}

		if len(episodes) > 0 {
			firstEp := episodes[0]
			t.Logf("First episode: %s (slug: %s)", firstEp.Name, firstEp.Slug)
		}
	}
}
