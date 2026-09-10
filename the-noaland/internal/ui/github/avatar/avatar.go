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
	return imagerenderer.Render(img, cols, rows).Raw
}

// View renders an inline avatar or reserves cells for out-of-band graphics.
func View(img image.Image, cols, rows int) string {
	result := imagerenderer.Render(img, cols, rows)
	if result.Inline != "" {
		return result.Inline
	}
	return Placeholder(cols, rows)
}
