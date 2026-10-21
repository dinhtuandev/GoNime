//go:build integration

package api

import (
	"testing"
	"time"
)

func logItems(t *testing.T, prefix string, items []AnimeItem) {
	t.Helper()
	if len(items) == 0 {
		return
	}
	n := len(items)
	if n > 3 {
		n = 3
	}
	for i := 0; i < n; i++ {
		t.Logf("   - %s %d: %s (slug: %s)", prefix, i+1, items[i].Name, items[i].Slug)
	}
}

// TestVSMOVSourceAvailable kiểm tra nguồn VSMOV có hoạt động không
func TestVSMOVSourceAvailable(t *testing.T) {
	client := NewVSMOVClient()

	start := time.Now()
	animes, cdnURL, err := client.GetTrending(1)
	if err != nil {
		t.Fatalf("GetTrending failed: %v", err)
	}
	t.Logf("GetTrending: %d items in %v, cdn=%q", len(animes), time.Since(start), cdnURL)
	if len(animes) == 0 {
		t.Fatal("no trending anime found")
	}
	logItems(t, "Trending", animes)
}

// TestVSMOVSearchAvailable kiểm tra tìm kiếm của VSMOV
func TestVSMOVSearchAvailable(t *testing.T) {
	client := NewVSMOVClient()

	for _, query := range []string{"avatar", "one piece", "conan"} {
		start := time.Now()
		results, _, err := client.Search(query)
		if err != nil {
			t.Errorf("Search %q failed: %v", query, err)
			continue
		}
		t.Logf("Search %q: %d results in %v", query, len(results), time.Since(start))
		logItems(t, "Result", results)
	}
}

// TestVSMOVEpisodes kiểm tra chi tiết + danh sách tập của VSMOV
func TestVSMOVEpisodes(t *testing.T) {
	client := NewVSMOVClient()

	results, _, err := client.Search("avatar")
	if err != nil || len(results) == 0 {
		results, _, err = client.GetTrending(1)
		if err != nil || len(results) == 0 {
			t.Skipf("không tìm được phim để test: %v", err)
		}
	}

	anime := results[0]
	movie, episodes, err := client.GetEpisodes(anime.Slug)
	if err != nil {
		t.Fatalf("GetEpisodes(%s) failed: %v", anime.Slug, err)
	}
	t.Logf("anime=%s movie=%s year=%d episodes=%d current=%q",
		anime.Name, movie.Name, movie.Year, len(episodes), movie.Episode)
	if len(episodes) == 0 {
		t.Fatal("no episodes returned")
	}

	n := len(episodes)
	if n > 5 {
		n = 5
	}
	for _, ep := range episodes[:n] {
		t.Logf("   - Episode %s: M3U8=%t Embed=%t", ep.Name, ep.LinkM3U8 != "", ep.LinkEmbed != "")
	}
}

// TestNguonCSourceAvailable kiểm tra nguồn NguonC có hoạt động không
func TestNguonCSourceAvailable(t *testing.T) {
	client := NewNguonCScraper()

	start := time.Now()
	animes, _, err := client.GetTrending(1)
	if err != nil {
		t.Fatalf("GetTrending failed: %v", err)
	}
	t.Logf("GetTrending: %d items in %v", len(animes), time.Since(start))
	if len(animes) == 0 {
		t.Fatal("no trending anime found")
	}
	logItems(t, "Trending", animes)
}

// TestNguonCSearchAvailable kiểm tra chức năng tìm kiếm của NguonC
func TestNguonCSearchAvailable(t *testing.T) {
	client := NewNguonCScraper()

	for _, query := range []string{"naruto", "one piece", "conan"} {
		start := time.Now()
		results, _, err := client.Search(query)
		if err != nil {
			t.Errorf("Search %q failed: %v", query, err)
			continue
		}
		t.Logf("Search %q: %d results in %v", query, len(results), time.Since(start))
		logItems(t, "Result", results)
	}
}

// TestNguonCEpisodes kiểm tra khi mở 1 anime cụ thể có load được các tập không
func TestNguonCEpisodes(t *testing.T) {
	client := NewNguonCScraper()

	results, _, err := client.Search("naruto")
	if err != nil || len(results) == 0 {
		t.Skipf("cannot find anime for episode test: %v", err)
	}

	movie, episodes, err := client.GetEpisodes(results[0].Slug)
	if err != nil {
		t.Fatalf("GetEpisodes failed: %v", err)
	}
	t.Logf("anime=%s year=%d episodes=%d current=%q",
		movie.Name, movie.Year, len(episodes), movie.Episode)
}

// TestAllAnimeSources kiểm tra mọi nguồn đang đăng ký
func TestAllAnimeSources(t *testing.T) {
	for _, name := range GetMovieSourceNames() {
		source := GetMovieSource(name)

		start := time.Now()
		trending, _, err := source.GetTrending(1)
		if err != nil {
			t.Errorf("%s GetTrending failed: %v", name, err)
		} else {
			t.Logf("%s GetTrending: %d items in %v", name, len(trending), time.Since(start))
		}

		start = time.Now()
		results, _, err := source.Search("one piece")
		if err != nil {
			t.Errorf("%s Search failed: %v", name, err)
		} else {
			t.Logf("%s Search: %d items in %v", name, len(results), time.Since(start))
		}
	}
}
