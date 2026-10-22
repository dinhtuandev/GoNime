package render

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

// RenderHalfBlock renders an image file as truecolor half-blocks.
func RenderHalfBlock(path string, width, height int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}
	return HalfBlockString(img, width, height), nil
}

// HalfBlockString converts an image to ANSI truecolor using "▀" cells:
// foreground paints the top pixel, background the bottom pixel.
func HalfBlockString(img image.Image, width, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	bounds := img.Bounds()
	scaleX := float64(bounds.Dx()) / float64(width)
	scaleY := float64(bounds.Dy()) / float64(height*2)

	var b strings.Builder
	b.Grow(width * height * 24)

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			x0 := bounds.Min.X + int(float64(col)*scaleX)
			x1 := bounds.Min.X + int(float64(col+1)*scaleX)
			if x1 <= x0 {
				x1 = x0 + 1
			}

			topY0 := bounds.Min.Y + int(float64(row*2)*scaleY)
			topY1 := bounds.Min.Y + int(float64(row*2+1)*scaleY)
			if topY1 <= topY0 {
				topY1 = topY0 + 1
			}
			bottomY0 := topY1
			bottomY1 := bounds.Min.Y + int(float64(row*2+2)*scaleY)
			if bottomY1 <= bottomY0 {
				bottomY1 = bottomY0 + 1
			}

			tr, tg, tb := averageColor(img, x0, topY0, x1, topY1)
			br, bg, bb := averageColor(img, x0, bottomY0, x1, bottomY1)

			fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀",
				tr, tg, tb, br, bg, bb)
		}
		b.WriteString("\x1b[0m\n")
	}

	return b.String()
}

// averageColor returns the mean RGB of a pixel window.
func averageColor(img image.Image, x0, y0, x1, y1 int) (uint8, uint8, uint8) {
	var r, g, b, n uint64
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			cr, cg, cb, _ := img.At(x, y).RGBA()
			r += uint64(cr >> 8)
			g += uint64(cg >> 8)
			b += uint64(cb >> 8)
			n++
		}
	}
	if n == 0 {
		return 0, 0, 0
	}
	return uint8(r / n), uint8(g / n), uint8(b / n)
}
