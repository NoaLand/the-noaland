package app

import (
	"testing"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

func TestResolveLayout(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		want          screen.Layout
	}{
		{"initial size", 0, 0, screen.LayoutTall},
		{"below minimum width", 89, 23, screen.LayoutTall},
		{"minimum compact width", 90, 23, screen.LayoutCompactWide},
		{"compact height boundary", 90, 24, screen.LayoutWide},
		{"above aspect ratio", 90, 44, screen.LayoutWide},
		{"equal aspect ratio", 90, 45, screen.LayoutTall},
		{"below aspect ratio", 90, 46, screen.LayoutTall},
		{"below very wide width", 149, 24, screen.LayoutWide},
		{"very wide width boundary", 150, 24, screen.LayoutVeryWide},
		{"compact takes precedence", 150, 23, screen.LayoutCompactWide},
		{"very wide above aspect ratio", 150, 74, screen.LayoutVeryWide},
		{"tall takes precedence", 150, 75, screen.LayoutTall},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveLayout(tt.width, tt.height); got != tt.want {
				t.Errorf("resolveLayout(%d, %d) = %v, want %v", tt.width, tt.height, got, tt.want)
			}
		})
	}
}
