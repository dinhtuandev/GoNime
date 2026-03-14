package api

import (
	"testing"
	"time"
)

// TestOPhimSourceAvailable kiểm tra nguồn OPhim có hoạt động không
func TestOPhimSourceAvailable(t *testing.T) {
	t.Log("=== Test: OPhim Source Available ===")

	client := NewOPhimClient()

	t.Log("Testing GetTrending...")
	start := time.Now()
	animes, cdnURL, err := client.GetTrending(1)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetTrending failed: %v", err)
	} else {
		t.Logf("✅ GetTrending succeeded in %v", duration)
		t.Logf("   Found %d trending anime", len(animes))
		t.Logf("   CDN URL: %s", cdnURL)

		if len(animes) == 0 {
			t.Error("❌ No trending anime found")
		} else {
			for i, anime := range animes[:3] {
				t.Logf("   - Trending %d: %s (slug: %s)", i+1, anime.Name, anime.Slug)
			}
		}
	}
}

// TestOPhimSearchAvailable kiểm tra chức năng tìm kiếm của OPhim
func TestOPhimSearchAvailable(t *testing.T) {
	t.Log("=== Test: OPhim Search Available ===")

	client := NewOPhimClient()

	testQueries := []string{"naruto", "one piece", "dragon"}

	for _, query := range testQueries {
		t.Logf("Testing search with keyword: '%s'", query)
		start := time.Now()
		results, cdnURL, err := client.Search(query)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("❌ Search '%s' failed: %v", query, err)
		} else {
			t.Logf("✅ Search '%s' succeeded in %v", query, duration)
			t.Logf("   Found %d results", len(results))
			t.Logf("   CDN URL: %s", cdnURL)

			for i, anime := range results[:3] {
				t.Logf("   - Result %d: %s (slug: %s)", i+1, anime.Name, anime.Slug)
			}
		}
	}
}

// TestOPhimEpisodes kiểm tra khi mở 1 anime cụ thể có load được các tập không
func TestOPhimEpisodes(t *testing.T) {
	t.Log("=== Test: OPhim Episodes Load ===")

	client := NewOPhimClient()

	// Tìm kiếm anime trước
	t.Log("Searching for anime to test episode loading...")
	results, _, err := client.Search("naruto")
	if err != nil || len(results) == 0 {
		// Thử tìm anime khác
		results, _, err = client.Search("one piece")
		if err != nil || len(results) == 0 {
			t.Skip("Cannot find anime for episode test")
			return
		}
	}

	anime := results[0]
	t.Logf("Testing with anime: %s (slug: %s)", anime.Name, anime.Slug)

	// Lấy chi tiết anime
	start := time.Now()
	movieInfo, episodes, err := client.GetEpisodes(anime.Slug)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetEpisodes failed: %v", err)
		return
	}

	t.Logf("✅ GetEpisodes succeeded in %v", duration)
	t.Logf("   Anime: %s", movieInfo.Name)
	t.Logf("   Episodes count: %d", len(episodes))
	t.Logf("   Episode current: %s", movieInfo.Episode)

	if len(episodes) > 0 {
		// Hiển thị một số tập đầu tiên
		t.Log("   First 5 episodes:")
		for _, ep := range episodes[:5] {
			t.Logf("   - Episode %s: M3U8=%t, Embed=%t", ep.Name, ep.LinkM3U8 != "", ep.LinkEmbed != "")
		}
	}
}

// TestNguonCSourceAvailable kiểm tra nguồn NguonC có hoạt động không
func TestNguonCSourceAvailable(t *testing.T) {
	t.Log("=== Test: NguonC Source Available ===")

	client := NewNguonCScraper()

	t.Log("Testing GetTrending...")
	start := time.Now()
	animes, cdnURL, err := client.GetTrending(1)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetTrending failed: %v", err)
	} else {
		t.Logf("✅ GetTrending succeeded in %v", duration)
		t.Logf("   Found %d trending anime", len(animes))
		t.Logf("   CDN URL: %s", cdnURL)

		if len(animes) == 0 {
			t.Error("❌ No trending anime found")
		} else {
			for i, anime := range animes[:3] {
				t.Logf("   - Trending %d: %s (slug: %s)", i+1, anime.Name, anime.Slug)
			}
		}
	}
}

// TestNguonCSearchAvailable kiểm tra chức năng tìm kiếm của NguonC
func TestNguonCSearchAvailable(t *testing.T) {
	t.Log("=== Test: NguonC Search Available ===")

	client := NewNguonCScraper()

	testQueries := []string{"naruto", "one piece", "conan"}

	for _, query := range testQueries {
		t.Logf("Testing search with keyword: '%s'", query)
		start := time.Now()
		results, cdnURL, err := client.Search(query)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("❌ Search '%s' failed: %v", query, err)
		} else {
			t.Logf("✅ Search '%s' succeeded in %v", query, duration)
			t.Logf("   Found %d results", len(results))
			t.Logf("   CDN URL: %s", cdnURL)

			for i, anime := range results[:3] {
				t.Logf("   - Result %d: %s (slug: %s)", i+1, anime.Name, anime.Slug)
			}
		}
	}
}

// TestNguonCEpisodes kiểm tra khi mở 1 anime cụ thể có load được các tập không
func TestNguonCEpisodes(t *testing.T) {
	t.Log("=== Test: NguonC Episodes Load ===")

	client := NewNguonCScraper()

	// Tìm kiếm anime trước
	t.Log("Searching for anime to test episode loading...")
	results, _, err := client.Search("naruto")
	if err != nil || len(results) == 0 {
		t.Skip("Cannot find anime for episode test")
		return
	}

	anime := results[0]
	t.Logf("Testing with anime: %s (slug: %s)", anime.Name, anime.Slug)

	// Lấy chi tiết anime
	start := time.Now()
	movieInfo, episodes, err := client.GetEpisodes(anime.Slug)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetEpisodes failed: %v", err)
		return
	}

	t.Logf("✅ GetEpisodes succeeded in %v", duration)
	t.Logf("   Anime: %s", movieInfo.Name)
	t.Logf("   Episodes count: %d", len(episodes))
	t.Logf("   Episode current: %s", movieInfo.Episode)

	if len(episodes) > 0 {
		// Hiển thị một số tập đầu tiên
		t.Log("   First 5 episodes:")
		for _, ep := range episodes[:5] {
			t.Logf("   - Episode %s: M3U8=%t, Embed=%t", ep.Name, ep.LinkM3U8 != "", ep.LinkEmbed != "")
		}
	}
}

// TestAllAnimeSources kiểm tra tất cả các nguồn anime
func TestAllAnimeSources(t *testing.T) {
	t.Log("=== Test: All Anime Sources ===")

	// Test từng nguồn anime
	sources := []struct {
		name   string
		client interface{}
	}{
		{"PhimAPI", NewAnimeClient()},
		{"OPhim", NewOPhimClient()},
		{"NguonC", NewNguonCScraper()},
	}

	for _, source := range sources {
		name := source.name
		t.Logf("Testing source: %s", name)

		var getTrending func(int) ([]AnimeItem, string, error)
		var search func(string) ([]AnimeItem, string, error)

		switch c := source.client.(type) {
		case *AnimeClient:
			getTrending = c.GetTrending
			search = c.Search
		case *OPhimClient:
			getTrending = c.GetTrending
			search = c.Search
		case *NguonCScraper:
			getTrending = c.GetTrending
			search = c.Search
		}

		// Test GetTrending
		start := time.Now()
		trending, _, err := getTrending(1)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("❌ Source %s - GetTrending failed: %v", name, err)
		} else {
			t.Logf("✅ Source %s - GetTrending: %d items in %v", name, len(trending), duration)
		}

		// Test Search
		start = time.Now()
		searchResults, _, err := search("test")
		duration = time.Since(start)

		if err != nil {
			t.Errorf("❌ Source %s - Search failed: %v", name, err)
		} else {
			t.Logf("✅ Source %s - Search: %d items in %v", name, len(searchResults), duration)
		}
	}
}
