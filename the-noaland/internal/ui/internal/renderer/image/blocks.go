package image

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// renderBlocks samples two pixels per cell and resets colors at each row.
func renderBlocks(img image.Image, cols, rows int) string {
	bounds := img.Bounds()
	pixel := func(x, y int) color.NRGBA {
		c := color.NRGBAModel.Convert(img.At(bounds.Min.X+x*bounds.Dx()/cols, bounds.Min.Y+y*bounds.Dy()/(rows*2))).(color.NRGBA)
		// Composite transparent pixels against the dashboard's dark background.
		blend := func(v, bg uint8) uint8 { return uint8((uint32(v)*uint32(c.A) + uint32(bg)*(255-uint32(c.A))) / 255) }
		return color.NRGBA{R: blend(c.R, 26), G: blend(c.G, 27), B: blend(c.B, 38), A: 255}
	}
	var out strings.Builder
	for y := 0; y < rows; y++ {
		if y > 0 {
			out.WriteByte('\n')
		}
		for x := 0; x < cols; x++ {
			top, bottom := pixel(x, y*2), pixel(x, y*2+1)
			fmt.Fprintf(&out, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", top.R, top.G, top.B, bottom.R, bottom.G, bottom.B)
		}
		out.WriteString("\x1b[0m")
	}
	return out.String()
}
