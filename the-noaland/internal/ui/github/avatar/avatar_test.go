package avatar

import (
	"image"
	"strings"
	"testing"
)

func TestPlaceholderAndMissingImage(t *testing.T) {
	if got := Placeholder(3, 2); got != "   \n   " {
		t.Fatalf("placeholder = %q", got)
	}
	if got := Render(nil, 14, 7); got != "" {
		t.Fatalf("missing image output = %q", got)
	}
}

func TestViewUsesFallbackWithoutRawOutput(t *testing.T) {
	t.Setenv("NOALAND_IMAGE_PROTOCOL", "blocks")
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	if got := View(img, 3, 2); strings.Count(got, "▀") != 6 {
		t.Fatal("inline avatar missing")
	}
	if Render(img, 3, 2) != "" {
		t.Fatal("fallback must not trigger a raw redraw")
	}
	if View(nil, 3, 2) != Placeholder(3, 2) {
		t.Fatal("loading placeholder changed")
	}
	t.Setenv("NOALAND_IMAGE_PROTOCOL", "kitty")
	if View(img, 3, 2) != Placeholder(3, 2) || Render(img, 3, 2) == "" {
		t.Fatal("Kitty placement contract changed")
	}
}
