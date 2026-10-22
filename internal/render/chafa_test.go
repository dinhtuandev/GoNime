package render

import "testing"

func TestNormalizeMode(t *testing.T) {
	tests := map[string]string{
		"auto":      "auto",
		"native":    "native",
		"chafa":     "chafa",
		"halfblock": "halfblock",
		"truecolor": "halfblock",
		"AUTO":      "auto",
		"foo":       "auto",
		"":          "auto",
	}

	for in, want := range tests {
		if got := NormalizeMode(in); got != want {
			t.Fatalf("NormalizeMode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestResolveEffectiveMode(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("TERM_PROGRAM", "")
	if got := ResolveEffectiveMode("native"); got != "halfblock" {
		t.Fatalf("expected native to fallback to halfblock when unsupported, got %q", got)
	}
	if got := ResolveEffectiveMode("auto"); got != "halfblock" {
		t.Fatalf("expected auto -> halfblock without native support, got %q", got)
	}

	t.Setenv("KITTY_WINDOW_ID", "1")
	if got := ResolveEffectiveMode("auto"); got != "native" {
		t.Fatalf("expected auto -> native when kitty detected, got %q", got)
	}

	if got := ResolveEffectiveMode("chafa"); got != "chafa" {
		t.Fatalf("expected chafa mode unchanged, got %q", got)
	}
	if got := ResolveEffectiveMode("halfblock"); got != "halfblock" {
		t.Fatalf("expected halfblock mode unchanged, got %q", got)
	}
}

func TestDetectNativeFormat(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("TERM_PROGRAM", "")
	if got := DetectNativeFormat(); got != "" {
		t.Fatalf("expected no native format, got %q", got)
	}

	t.Setenv("KITTY_WINDOW_ID", "1")
	if got := DetectNativeFormat(); got != "kitty" {
		t.Fatalf("expected kitty format, got %q", got)
	}
}
