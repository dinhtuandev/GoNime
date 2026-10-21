package api

import (
	"encoding/json"
	"testing"
)

// TestVSMOVJSONText covers poster_url values seen live: string, null, {}.
func TestVSMOVJSONText(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`"https://img.example/a.jpg"`, "https://img.example/a.jpg"},
		{`null`, ""},
		{`{}`, ""},
		{`[]`, ""},
	}
	for _, tc := range cases {
		var v jsonText
		if err := json.Unmarshal([]byte(tc.in), &v); err != nil {
			t.Fatalf("unmarshal %s: %v", tc.in, err)
		}
		if got := v.String(); got != tc.want {
			t.Fatalf("jsonText(%s) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestVSMOVDetailDecode(t *testing.T) {
	const payload = `{
      "status": true,
      "movie": {"name": "A", "slug": "a", "year": 2024, "poster_url": {}, "content": "desc"},
      "episodes": [{"server_name": "Vietsub", "server_data": [
        {"name": "1", "slug": "1", "link_m3u8": "", "link_embed": "https://v1.streamvsmov.com/video/x"}
      ]}]
    }`

	var result vsmovDetailResponse
	if err := json.Unmarshal([]byte(payload), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.Movie.Name != "A" || result.Movie.Year != 2024 {
		t.Fatalf("movie = %+v", result.Movie)
	}
	if result.Movie.PosterURL.String() != "" {
		t.Fatalf("poster should be empty, got %q", result.Movie.PosterURL)
	}
	if result.Movie.Content.String() != "desc" {
		t.Fatalf("content = %q", result.Movie.Content)
	}
	if len(result.Episodes) != 1 || len(result.Episodes[0].ServerData) != 1 {
		t.Fatalf("episodes = %+v", result.Episodes)
	}
	if got := result.Episodes[0].ServerData[0].LinkEmbed; got != "https://v1.streamvsmov.com/video/x" {
		t.Fatalf("link_embed = %q", got)
	}
}
