package api

import (
	"encoding/json"
	"testing"
)

// TestFlexIntYear covers the NguonC responses that broke GetDetails:
// the API sends year as a JSON number, the struct used to expect a string.
func TestFlexIntYear(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		want    int
	}{
		{"number", `{"movie":{"year":2024}}`, 2024},
		{"string", `{"movie":{"year":"2021"}}`, 2021},
		{"missing", `{"movie":{}}`, 0},
		{"invalid", `{"movie":{"year":"unknown"}}`, 0},
		{"null", `{"movie":{"year":null}}`, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var result nguoncDetailResponse
			if err := json.Unmarshal([]byte(tc.payload), &result); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got := int(result.Movie.Year); got != tc.want {
				t.Fatalf("year = %d, want %d", got, tc.want)
			}
		})
	}
}
