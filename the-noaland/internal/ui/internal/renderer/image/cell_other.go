//go:build !windows

package image

// Kitty uses cell dimensions directly. Sixel retains its documented estimate
// when pixel metrics are unavailable; do not read stdin alongside Bubble Tea.
func systemCellPixels() (int, int) { return 0, 0 }
