package image

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi/sixel"
)

func TestSixelPayloadDimensionsAndPixels(t *testing.T) {
	t.Setenv("NOALAND_IMAGE_PROTOCOL", "sixel")
	t.Setenv("NOALAND_IMAGE_CELL_SIZE", "10x18")
	img := image.NewNRGBA(image.Rect(5, 7, 7, 9))
	for y := 7; y < 9; y++ {
		for x := 5; x < 7; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: 255, A: 255})
		}
	}
	result := Render(img, 14, 7)
	if result.Inline != "" || !strings.HasPrefix(result.Raw, "\x1bP0;1;0q") || !strings.HasSuffix(result.Raw, "\x1b\\") {
		t.Fatal("invalid Sixel framing or fallback selected")
	}
	payload := strings.TrimSuffix(strings.TrimPrefix(result.Raw, "\x1bP0;1;0q"), "\x1b\\")
	decoder := sixel.Decoder{}
	decoded, err := decoder.Decode(strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds().Dx() != 140 || decoded.Bounds().Dy() != 126 {
		t.Fatalf("wrong pixel bounds: %v", decoded.Bounds())
	}
	if c := color.NRGBAModel.Convert(decoded.At(0, 0)).(color.NRGBA); c.R != 255 || c.G != 0 || c.B != 0 {
		t.Fatalf("source pixels changed: %v", c)
	}
}

func TestCellPixels(t *testing.T) {
	for _, tt := range []struct {
		input string
		w, h  int
	}{
		{"", 8, 16}, {"10x20", 10, 20}, {"9X18", 9, 18}, {"0x20", 8, 16}, {"-1x20", 8, 16}, {"1000x2000", 8, 16}, {"bad", 8, 16},
	} {
		w, h := cellPixels(tt.input)
		if w != tt.w || h != tt.h {
			t.Fatalf("%q = %dx%d", tt.input, w, h)
		}
	}
}

func TestSixelOverrideAndMultiplexer(t *testing.T) {
	for _, tt := range []struct {
		env  map[string]string
		want string
	}{
		{map[string]string{"NOALAND_IMAGE_PROTOCOL": "sixel"}, "sixel"},
		{map[string]string{"WT_SESSION": "session", "NOALAND_IMAGE_PROTOCOL": "blocks"}, "blocks"},
		{map[string]string{"WT_SESSION": "session", "TMUX": "session"}, "blocks"},
	} {
		if got := selectProtocol(func(k string) string { return tt.env[k] }); got != tt.want {
			t.Fatalf("got %q", got)
		}
	}
}
