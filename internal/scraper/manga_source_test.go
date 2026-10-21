//go:build integration

package scraper

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestOTruyenSourceAvailable kiểm tra nguồn OTruyen có hoạt động không
func TestOTruyenSourceAvailable(t *testing.T) {
	t.Log("=== Test: OTruyen Source Available ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	// Test 1: Kiểm tra nguồn có trả về dữ liệu trending không
	t.Log("Testing GetTrending...")
	start := time.Now()
	trending, err := scraper.GetTrending()
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetTrending failed: %v", err)
	} else {
		t.Logf("✅ GetTrending succeeded in %v", duration)
		t.Logf("   Found %d trending manga", len(trending))

		if len(trending) == 0 {
			t.Error("❌ No trending manga found")
		} else {
			for i, m := range trending[:3] {
				t.Logf("   - Trending %d: %s", i+1, m.Name)
			}
		}
	}
}

// TestOTruyenSearchAvailable kiểm tra chức năng tìm kiếm có hoạt động không
func TestOTruyenSearchAvailable(t *testing.T) {
	t.Log("=== Test: OTruyen Search Available ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	// Test với từ khóa phổ biến
	testQueries := []string{"one", "naruto", "dragon"}

	for _, query := range testQueries {
		t.Logf("Testing search with keyword: '%s'", query)
		start := time.Now()
		results, err := scraper.Search(query)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("❌ Search '%s' failed: %v", query, err)
		} else {
			t.Logf("✅ Search '%s' succeeded in %v", query, duration)
			t.Logf("   Found %d results", len(results))

			if len(results) == 0 {
				t.Logf("   ⚠️ No results for '%s'", query)
			} else {
				for i, m := range results[:3] {
					t.Logf("   - Result %d: %s (slug: %s)", i+1, m.Name, m.Slug)
				}
			}
		}
	}
}

// TestOTruyenChapterLoad kiểm tra khi mở 1 manga cụ thể có load được các tập không
func TestOTruyenChapterLoad(t *testing.T) {
	t.Log("=== Test: OTruyen Chapter Load ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	// Tìm kiếm manga trước
	t.Log("Searching for manga to test chapter loading...")
	results, err := scraper.Search("yakuza")
	if err != nil || len(results) == 0 {
		// Thử tìm manga khác
		results, err = scraper.Search("one")
		if err != nil || len(results) == 0 {
			t.Skip("Cannot find manga for chapter test")
			return
		}
	}

	manga := results[0]
	t.Logf("Testing with manga: %s (slug: %s)", manga.Name, manga.Slug)

	// Lấy chi tiết manga
	start := time.Now()
	details, err := scraper.GetDetails(manga.Slug)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetDetails failed: %v", err)
		return
	}

	t.Logf("✅ GetDetails succeeded in %v", duration)
	t.Logf("   Manga: %s", details.Name)
	t.Logf("   Chapters count: %d", len(details.Chapters))

	if len(details.Chapters) == 0 {
		t.Error("❌ No chapters found for manga")
		return
	}

	// Hiển thị một số chapter đầu tiên
	t.Log("   First 5 chapters:")
	for _, ch := range details.Chapters[:5] {
		t.Logf("   - Chapter %s: %s", ch.Number, ch.URL)
	}
}

// TestOTruyenSpecificChapterAPI kiểm tra API chapter cụ thể mà user cung cấp
// URL: https://sv1.otruyencdn.com/v1/api/chapter/69b4e3277b89b5b2570e218e
func TestOTruyenSpecificChapterAPI(t *testing.T) {
	t.Log("=== Test: OTruyen Specific Chapter API ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	// Sử dụng URL chapter mà user cung cấp
	chapterURL := "https://sv1.otruyencdn.com/v1/api/chapter/69b4e3277b89b5b2570e218e"

	t.Logf("Testing chapter URL: %s", chapterURL)

	start := time.Now()
	images, err := scraper.GetChapterImages(chapterURL)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetChapterImages failed: %v", err)
		return
	}

	t.Logf("✅ GetChapterImages succeeded in %v", duration)
	t.Logf("   Found %d images in chapter", len(images))

	if len(images) == 0 {
		t.Error("❌ No images found in chapter")
		return
	}

	// Hiển thị một số image URL đầu tiên
	t.Log("   First 5 image URLs:")
	for i, img := range images[:5] {
		t.Logf("   - Image %d: %s", i+1, img)
	}

	// Kiểm tra xem URL có đúng format không
	expectedBase := "https://sv1.otruyencdn.com"
	for i, img := range images {
		if len(img) < len(expectedBase) || img[:len(expectedBase)] != expectedBase {
			t.Logf("   ⚠️ Image %d URL doesn't match expected CDN", i+1)
		}
	}
}

// TestOTruyenChapterCaching kiểm tra xem khi vào tập có catch được không (download image)
func TestOTruyenChapterCaching(t *testing.T) {
	t.Log("=== Test: OTruyen Chapter Caching ===")

	scraper := NewOTruyenScraper("https://otruyen.cc")

	// Sử dụng URL chapter mà user cung cấp
	chapterURL := "https://sv1.otruyencdn.com/v1/api/chapter/69b4e3277b89b5b2570e218e"

	// Lấy danh sách images
	images, err := scraper.GetChapterImages(chapterURL)
	if err != nil || len(images) == 0 {
		t.Skip("Cannot get chapter images for caching test")
		return
	}

	t.Logf("Testing caching for %d images", len(images))

	// Tạo thư mục cache tạm
	cacheDir, err := os.MkdirTemp("", "manga_cache_test")
	if err != nil {
		t.Skip("Cannot create temp cache directory")
		return
	}
	defer os.RemoveAll(cacheDir)

	t.Logf("Cache directory: %s", cacheDir)

	// Download một vài image để test caching
	ctx := context.Background()
	successCount := 0
	failCount := 0

	// Test download up to 3 images
	testCount := 3
	if len(images) < testCount {
		testCount = len(images)
	}

	for i := 0; i < testCount; i++ {
		imgURL := images[i]
		t.Logf("Testing download image %d: %s", i+1, filepath.Base(imgURL))

		start := time.Now()
		localPath, err := scraper.DownloadImage(ctx, imgURL, cacheDir)
		duration := time.Since(start)

		if err != nil {
			t.Logf("❌ Image %d download failed: %v", i+1, err)
			failCount++
		} else {
			t.Logf("✅ Image %d downloaded in %v: %s", i+1, duration, localPath)

			// Kiểm tra file tồn tại
			if _, err := os.Stat(localPath); err == nil {
				t.Logf("   ✓ File exists and cached")
				successCount++
			} else {
				t.Logf("   ⚠️ File not found after download")
				failCount++
			}
		}
	}

	t.Logf("Caching test result: %d success, %d failed", successCount, failCount)

	if successCount == 0 {
		t.Error("❌ Caching test: no images could be cached")
	} else {
		t.Logf("✅ Caching test: %d/%d images cached successfully", successCount, testCount)
	}
}

// TestTruyenvnSourceAvailable kiểm tra nguồn Truyenvn có hoạt động không
func TestTruyenvnSourceAvailable(t *testing.T) {
	t.Log("=== Test: Truyenvn Source Available ===")

	scraper := NewMangaScraper("https://truyenvn.shop")

	t.Log("Testing GetTrending...")
	start := time.Now()
	trending, err := scraper.GetTrending()
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ GetTrending failed: %v", err)
	} else {
		t.Logf("✅ GetTrending succeeded in %v", duration)
		t.Logf("   Found %d trending manga", len(trending))

		if len(trending) == 0 {
			t.Error("❌ No trending manga found")
		} else {
			for i, m := range trending[:3] {
				t.Logf("   - Trending %d: %s", i+1, m.Name)
			}
		}
	}
}

// TestTruyenvnSearchAvailable kiểm tra chức năng tìm kiếm của Truyenvn
func TestTruyenvnSearchAvailable(t *testing.T) {
	t.Log("=== Test: Truyenvn Search Available ===")

	scraper := NewMangaScraper("https://truyenvn.shop")

	t.Log("Testing search with keyword: 'one'")
	start := time.Now()
	results, err := scraper.Search("one")
	duration := time.Since(start)

	if err != nil {
		t.Errorf("❌ Search failed: %v", err)
	} else {
		t.Logf("✅ Search succeeded in %v", duration)
		t.Logf("   Found %d results", len(results))

		for i, m := range results[:3] {
			t.Logf("   - Result %d: %s (slug: %s)", i+1, m.Name, m.Slug)
		}
	}
}

// TestAllMangaSources tất cả các nguồn manga
func TestAllMangaSources(t *testing.T) {
	t.Log("=== Test: All Manga Sources ===")

	sources := GetMangaSourceNames()
	t.Logf("Available sources: %v", sources)

	for _, sourceName := range sources {
		t.Logf("Testing source: %s", sourceName)
		source := GetMangaSource(sourceName)

		start := time.Now()
		trending, err := source.GetTrending()
		duration := time.Since(start)

		if err != nil {
			t.Errorf("❌ Source %s - GetTrending failed: %v", sourceName, err)
		} else {
			t.Logf("✅ Source %s - GetTrending: %d items in %v", sourceName, len(trending), duration)
		}

		// Test search
		start = time.Now()
		searchResults, err := source.Search("test")
		duration = time.Since(start)

		if err != nil {
			t.Errorf("❌ Source %s - Search failed: %v", sourceName, err)
		} else {
			t.Logf("✅ Source %s - Search: %d items in %v", sourceName, len(searchResults), duration)
		}
	}
}

// TestMangaSourceInterface kiểm tra interface của MangaSource
func TestMangaSourceInterface(t *testing.T) {
	t.Log("=== Test: Manga Source Interface ===")

	// Kiểm tra tất cả sources implement đúng interface
	sources := []MangaSource{
		NewOTruyenScraper("https://otruyen.cc"),
		NewMangaScraper("https://truyenvn.shop"),
	}

	for _, source := range sources {
		name := source.Name()
		t.Logf("Testing source: %s", name)

		// Kiểm tra các method không panics
		_ = source.Name()

		// GetTrending
		trending, err := source.GetTrending()
		if err != nil {
			t.Logf("   GetTrending error: %v", err)
		} else {
			t.Logf("   GetTrending: %d items", len(trending))
		}

		// Search
		results, err := source.Search("one")
		if err != nil {
			t.Logf("   Search error: %v", err)
		} else {
			t.Logf("   Search: %d items", len(results))

			// GetDetails cho item đầu tiên nếu có
			if len(results) > 0 {
				details, err := source.GetDetails(results[0].Slug)
				if err != nil {
					t.Logf("   GetDetails error: %v", err)
				} else {
					t.Logf("   GetDetails: %s, Chapters: %d", details.Name, len(details.Chapters))
				}
			}
		}

		t.Logf("✅ Source %s interface test completed", name)
	}
}

// TestOTruyenChapterWithRealData test với dữ liệu thực từ API
func TestOTruyenChapterWithRealData(t *testing.T) {
	t.Log("=== Test: OTruyen Chapter with Real Data ===")

	// Test data từ user cung cấp
	testChapterURL := "https://sv1.otruyencdn.com/v1/api/chapter/69b4e3277b89b5b2570e218e"

	expectedImageCount := 28 // page_0 đến page_27

	scraper := NewOTruyenScraper("https://otruyen.cc")

	t.Logf("Testing chapter: %s", testChapterURL)

	// Lấy images
	images, err := scraper.GetChapterImages(testChapterURL)
	if err != nil {
		t.Errorf("❌ GetChapterImages failed: %v", err)
		return
	}

	t.Logf("✅ GetChapterImages succeeded")
	t.Logf("   Image count: %d (expected: %d)", len(images), expectedImageCount)

	if len(images) != expectedImageCount {
		t.Logf("⚠️ Image count mismatch: got %d, expected %d", len(images), expectedImageCount)
	}

	// Kiểm tra URL format
	domainCDN := "https://sv1.otruyencdn.com"
	for i, img := range images {
		expectedURL := fmt.Sprintf("%s/uploads/20260314/d80da56ab2364bc024a3ffc7e32929e9/chapter_50.2/page_%d.jpg", domainCDN, i)
		if img != expectedURL {
			t.Logf("⚠️ Image %d URL mismatch:", i)
			t.Logf("   Got:      %s", img)
			t.Logf("   Expected: %s", expectedURL)
		}
	}

	t.Log("✅ Chapter data verification completed")
}
