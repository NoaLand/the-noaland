package app

import "github.com/NoaLand/the-noaland/the-noaland/internal/screen"

// resolveLayout keeps terminal breakpoints and their precedence in one place.
func resolveLayout(width, height int) screen.Layout {
	if !(width >= 90 && width > height*2) {
		return screen.LayoutTall
	}
	if height < 24 {
		return screen.LayoutCompactWide
	}
	if width >= 150 {
		return screen.LayoutVeryWide
	}
	return screen.LayoutWide
}
