package main

// Layout identifies the presentation selected for the terminal dimensions.
type Layout int

const (
	LayoutTall Layout = iota
	LayoutCompactWide
	LayoutWide
	LayoutVeryWide
)

// resolveLayout keeps terminal breakpoints and their precedence in one place.
func resolveLayout(width, height int) Layout {
	if !(width >= 90 && width > height*2) {
		return LayoutTall
	}
	if height < 24 {
		return LayoutCompactWide
	}
	if width >= 150 {
		return LayoutVeryWide
	}
	return LayoutWide
}
