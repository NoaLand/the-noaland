package avatar

import (
	"image"
	"strings"

	imagerenderer "github.com/NoaLand/the-noaland/the-noaland/internal/ui/internal/renderer/image"
)

// Placeholder reserves terminal cells for an avatar.
func Placeholder(
	width int,
	height int,
) string {
	line := strings.Repeat(" ", width)

	lines := make([]string, height)

	for i := range lines {
		lines[i] = line
	}

	return strings.Join(lines, "\n")
}

// Render produces an avatar image using the shared terminal image renderer.
func Render(img image.Image, cols, rows int) string {
	return imagerenderer.Render(img, cols, rows)
}
