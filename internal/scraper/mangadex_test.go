package scraper

import "testing"

func TestChapterIDFromURL(t *testing.T) {
	const id = "0f9e3a3e-1b2c-4d5e-8f90-123456789abc"
	cases := []struct {
		in   string
		want string
	}{
		{"https://mangadex.org/chapter/" + id, id},
		{"https://mangadex.org/chapter/" + id + "/", id},
		{"https://mangadex.org/chapter/" + id + "?foo=1", id},
		{"  " + id + "  ", id},
		{id, id},
		{"", ""},
	}
	for _, tc := range cases {
		if got := chapterIDFromURL(tc.in); got != tc.want {
			t.Fatalf("chapterIDFromURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
