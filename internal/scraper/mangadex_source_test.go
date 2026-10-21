//go:build integration

package scraper

import (
	"context"
	"os"
	"testing"
)

// TestMangaDexSourceAvailable kiểm tra API MangaDex còn sống và trả dữ liệu.
func TestMangaDexSourceAvailable(t *testing.T) {
	s := NewMangaDexScraper()

	trending, err := s.GetTrending()
	if err != nil {
		t.Fatalf("GetTrending: %v", err)
	}
	if len(trending) == 0 {
		t.Fatal("GetTrending trả danh sách rỗng")
	}
	t.Logf("trending: %d truyện, đầu tiên: %s (%s)", len(trending), trending[0].Name, trending[0].Slug)

	manga, err := s.GetDetails(trending[0].Slug)
	if err != nil {
		t.Fatalf("GetDetails: %v", err)
	}
	if len(manga.Chapters) == 0 {
		t.Fatal("GetDetails không có chapter tiếng Việt")
	}
	t.Logf("chapters: %d, đầu tiên: %s", len(manga.Chapters), manga.Chapters[0].Title)

	images, err := s.GetChapterImages(manga.Chapters[0].URL)
	if err != nil {
		t.Fatalf("GetChapterImages: %v", err)
	}
	if len(images) == 0 {
		t.Fatal("GetChapterImages trả rỗng")
	}
	t.Logf("images: %d, đầu tiên: %s", len(images), images[0])

	dir := t.TempDir()
	path, err := s.DownloadImage(context.Background(), images[0], dir)
	if err != nil {
		t.Fatalf("DownloadImage: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		t.Fatalf("ảnh tải về rỗng: %v", err)
	}
	t.Logf("downloaded: %s (%d bytes)", path, info.Size())

}
