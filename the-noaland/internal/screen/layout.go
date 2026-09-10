package screen

// Layout identifies the presentation selected for the terminal dimensions.
type Layout int

const (
	LayoutTall Layout = iota
	LayoutCompactWide
	LayoutWide
	LayoutVeryWide
)
