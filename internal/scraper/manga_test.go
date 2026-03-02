package scraper

import (
	"testing"
	"time"
)

func TestMangaTrending(t *testing.T) {
	t.Log("=== Testing Manga Trending ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	t.Log("Testing GetTrending...")
	trendingStart := time.Now()
	trendingList, err := scraper.GetTrending()
	trendingDuration := time.Since(trendingStart)

	if err != nil {
		t.Errorf("GetTrending failed: %v", err)
	} else {
		t.Logf("GetTrending succeeded in %v", trendingDuration)
		t.Logf("Found %d trending manga", len(trendingList))

		if len(trendingList) == 0 {
			t.Log("Warning: No trending manga found - may be normal if site structure changed")
		}

		for i, manga := range trendingList {
			if i >= 3 {
				break
			}
			t.Logf("Trending %d: %s (slug: %s)", i+1, manga.Name, manga.Slug)
			if manga.Name == "" || manga.Slug == "" {
				t.Errorf("Trending manga %d has empty fields: %+v", i+1, manga)
			}
		}
	}
}

func TestMangaSearch(t *testing.T) {
	t.Log("=== Testing Manga Search ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	t.Log("Testing search with keyword 'one'...")
	searchStart := time.Now()
	searchResults, err := scraper.Search("one")
	searchDuration := time.Since(searchStart)

	if err != nil {
		t.Errorf("Search failed: %v", err)
	} else {
		t.Logf("Search succeeded in %v", searchDuration)
		t.Logf("Found %d search results", len(searchResults))

		for i, manga := range searchResults {
			if i >= 3 {
				break
			}
			t.Logf("Search %d: %s (slug: %s)", i+1, manga.Name, manga.Slug)
			if manga.Name == "" || manga.Slug == "" {
				t.Errorf("Search result %d has empty fields: %+v", i+1, manga)
			}
		}
	}
}

func TestMangaDetails(t *testing.T) {
	t.Log("=== Testing Manga Details ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	mangas, err := scraper.Search("one")
	if err != nil {
		t.Logf("Cannot get manga for details test: %v", err)
		return
	}

	if len(mangas) == 0 {
		t.Log("No manga available for details test")
		return
	}

	t.Log("Testing GetDetails for first manga...")
	detailsStart := time.Now()
	details, err := scraper.GetDetails(mangas[0].Slug)
	detailsDuration := time.Since(detailsStart)

	if err != nil {
		t.Errorf("GetDetails failed: %v", err)
	} else {
		t.Logf("GetDetails succeeded in %v", detailsDuration)
		t.Logf("Manga: %s", details.Name)
		t.Logf("Chapters: %d", len(details.Chapters))

		if details.Name == "" {
			t.Errorf("Manga details has empty name: %+v", details)
		}

		if len(details.Chapters) > 0 {
			firstChapter := details.Chapters[0]
			t.Logf("First chapter: %s (URL: %s)", firstChapter.Title, firstChapter.URL)
		}
	}
}

func TestChapterImages(t *testing.T) {
	t.Log("=== Testing Chapter Images ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	mangas, err := scraper.Search("one")
	if err != nil || len(mangas) == 0 {
		t.Logf("Cannot get manga for chapter test: %v", err)
		return
	}

	details, err := scraper.GetDetails(mangas[0].Slug)
	if err != nil || len(details.Chapters) == 0 {
		t.Logf("Cannot get chapters for image test: %v", err)
		return
	}

	t.Log("Testing GetChapterImages for first chapter...")
	imagesStart := time.Now()
	images, err := scraper.GetChapterImages(details.Chapters[0].URL)
	imagesDuration := time.Since(imagesStart)

	if err != nil {
		t.Errorf("GetChapterImages failed: %v", err)
	} else {
		t.Logf("GetChapterImages succeeded in %v", imagesDuration)
		t.Logf("Found %d images in chapter", len(images))

		for i, imgURL := range images {
			if i >= 3 {
				break
			}
			t.Logf("Image %d: %s", i+1, imgURL)
		}
	}
}
