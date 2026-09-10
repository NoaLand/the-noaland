package image

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi/sixel"
)

// Sixel is pixel-sized. Until cell-size queries are integrated with Bubble Tea,
// use a documented estimate with an explicit override for font/DPI differences.
func cellPixels(value string) (int, int) {
	parts := strings.Split(strings.ToLower(value), "x")
	if len(parts) == 2 {
		w, e1 := strconv.Atoi(parts[0])
		h, e2 := strconv.Atoi(parts[1])
		if e1 == nil && e2 == nil && w > 0 && h > 0 && w <= 128 && h <= 256 {
			return w, h
		}
	}
	return 8, 16
}

func renderSixel(img image.Image, cols, rows int) string {
	cw, ch := cellPixels(os.Getenv("NOALAND_IMAGE_CELL_SIZE"))
	// Avoid oversized allocations for malformed dimensions.
	if cols > 4096/cw || rows > 4096/ch {
		return ""
	}
	width, height := cols*cw, rows*ch
	source := img.Bounds()
	scaled := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := color.NRGBAModel.Convert(img.At(source.Min.X+x*source.Dx()/width, source.Min.Y+y*source.Dy()/height)).(color.NRGBA)
			// Flatten alpha so redraws replace old pixels, including transparent areas.
			blend := func(v, bg uint8) uint8 { return uint8((uint32(v)*uint32(c.A) + uint32(bg)*(255-uint32(c.A))) / 255) }
			scaled.SetNRGBA(x, y, color.NRGBA{R: blend(c.R, 26), G: blend(c.G, 27), B: blend(c.B, 38), A: 255})
		}
	}
	var out bytes.Buffer
	out.WriteString("\x1bP0;1;0q")
	encoder := sixel.Encoder{}
	if err := encoder.Encode(&out, scaled); err != nil {
		return ""
	}
	out.WriteString("\x1b\\")
	return out.String()
}
