package render

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestHalfBlockStringDimensions(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 30), G: uint8(y * 30), B: 120, A: 255})
		}
	}

	out := HalfBlockString(img, 4, 2)
	if got := strings.Count(out, "▀"); got != 8 {
		t.Fatalf("cells = %d, want 8", got)
	}
	if got := strings.Count(out, "\x1b[0m\n"); got != 2 {
		t.Fatalf("rows = %d, want 2", got)
	}
	if !strings.Contains(out, "\x1b[38;2;") || !strings.Contains(out, "\x1b[48;2;") {
		t.Fatal("missing truecolor foreground/background sequences")
	}
}

func TestHalfBlockStringClampsTinySize(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	out := HalfBlockString(img, 0, -1)
	if strings.Count(out, "▀") != 1 {
		t.Fatalf("expected a single cell, got %q", out)
	}
}
