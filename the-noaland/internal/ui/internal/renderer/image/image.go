package image

import (
	"image"
	"os"
	"strings"
)

// Result separates inline terminal cells from out-of-band image protocol output.
type Result struct {
	Inline string
	Raw    string
}

// Render uses Kitty in known compatible environments and inline cells elsewhere.
// NOALAND_IMAGE_PROTOCOL can explicitly select auto, kitty, blocks, or none.
func Render(img image.Image, cols, rows int) Result {
	if img == nil || cols <= 0 || rows <= 0 || img.Bounds().Empty() {
		return Result{}
	}
	switch selectProtocol(os.Getenv) {
	case "kitty":
		if raw := renderKitty(img, cols, rows); raw != "" {
			return Result{Raw: raw}
		}
	case "none":
		return Result{}
	}
	return Result{Inline: renderBlocks(img, cols, rows)}
}

// Environment detection is a heuristic, not a terminal capability probe.
// Multiplexers can change capabilities; explicit overrides take priority.
func selectProtocol(getenv func(string) string) string {
	switch strings.ToLower(strings.TrimSpace(getenv("NOALAND_IMAGE_PROTOCOL"))) {
	case "kitty":
		return "kitty"
	case "blocks":
		return "blocks"
	case "none":
		return "none"
	}
	if getenv("TMUX") != "" {
		return "blocks"
	}
	program := strings.ToLower(getenv("TERM_PROGRAM"))
	term := strings.ToLower(getenv("TERM"))
	if program == "ghostty" || program == "kitty" || term == "xterm-kitty" || term == "xterm-ghostty" {
		return "kitty"
	}
	return "blocks"
}
