package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

// RenderNativeImage encodes an image file with the terminal's native graphics
// protocol so the terminal draws real pixels instead of text cells.
func RenderNativeImage(path string, cols, rows int) (string, error) {
	switch DetectNativeFormat() {
	case "kitty":
		return renderKitty(path, cols, rows)
	case "iterm":
		return renderITerm2(path, cols, rows)
	}
	return "", fmt.Errorf("no native image protocol detected")
}

// renderKitty sends the image as a PNG through the Kitty graphics protocol,
// sized to cols x rows terminal cells.
func renderKitty(path string, cols, rows int) (string, error) {
	img, err := decodeImage(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, downscale(img, cols*8, rows*16)); err != nil {
		return "", err
	}

	data := base64.StdEncoding.EncodeToString(buf.Bytes())
	var b strings.Builder
	const chunk = 4096
	for first := true; data != ""; first = false {
		n := chunk
		if n > len(data) {
			n = len(data)
		}
		part, rest := data[:n], data[n:]
		more := 1
		if rest == "" {
			more = 0
		}
		if first {
			fmt.Fprintf(&b, "\x1b_Ga=T,f=100,c=%d,r=%d,m=%d;%s\x1b\\", cols, rows, more, part)
		} else {
			fmt.Fprintf(&b, "\x1b_Gm=%d;%s\x1b\\", more, part)
		}
		data = rest
	}
	return b.String(), nil
}

// renderITerm2 sends the original file bytes through the iTerm2 inline image
// protocol, which iTerm2, WezTerm and Ghostty understand.
func renderITerm2(path string, cols, rows int) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\x1b]1337;File=inline=1;width=%d;height=%d;preserveAspectRatio=1:%s\x07",
		cols, rows, base64.StdEncoding.EncodeToString(raw)), nil
}

func decodeImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

// downscale resizes with box averaging, which keeps small thumbnails stable.
func downscale(img image.Image, width, height int) *image.NRGBA {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	src := img.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		y0 := src.Min.Y + y*src.Dy()/height
		y1 := src.Min.Y + (y+1)*src.Dy()/height
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for x := 0; x < width; x++ {
			x0 := src.Min.X + x*src.Dx()/width
			x1 := src.Min.X + (x+1)*src.Dx()/width
			if x1 <= x0 {
				x1 = x0 + 1
			}

			var r, g, b, n uint32
			for yy := y0; yy < y1; yy++ {
				for xx := x0; xx < x1; xx++ {
					cr, cg, cb, _ := img.At(xx, yy).RGBA()
					r += cr >> 8
					g += cg >> 8
					b += cb >> 8
					n++
				}
			}
			dst.SetNRGBA(x, y, color.NRGBA{R: uint8(r / n), G: uint8(g / n), B: uint8(b / n), A: 255})
		}
	}
	return dst
}
