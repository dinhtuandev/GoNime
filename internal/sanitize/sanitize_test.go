package sanitize

import "testing"

func TestTextStripsEscapes(t *testing.T) {
	tests := map[string]string{
		"plain":                     "plain",
		"\x1b[31mred\x1b[0m":        "red",
		"a\x1b]8;;http://x\x07b":    "ab",
		"tab\tkept\nline":           "tab\tkept\nline",
		"bell\x07gone":              "bellgone",
		"del\x7fchar":               "delchar",
		"c1\u009bcontrol":           "c1control",
		"osc\x1b]0;title\x1b\\tail": "osctail",
	}

	for in, want := range tests {
		if got := Text(in); got != want {
			t.Fatalf("Text(%q) = %q, want %q", in, got, want)
		}
	}
}
