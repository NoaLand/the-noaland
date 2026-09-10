package main

import "testing"

func TestResolveLayout(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		want          Layout
	}{
		{"initial size", 0, 0, LayoutTall},
		{"below minimum width", 89, 23, LayoutTall},
		{"minimum compact width", 90, 23, LayoutCompactWide},
		{"compact height boundary", 90, 24, LayoutWide},
		{"above aspect ratio", 90, 44, LayoutWide},
		{"equal aspect ratio", 90, 45, LayoutTall},
		{"below aspect ratio", 90, 46, LayoutTall},
		{"below very wide width", 149, 24, LayoutWide},
		{"very wide width boundary", 150, 24, LayoutVeryWide},
		{"compact takes precedence", 150, 23, LayoutCompactWide},
		{"very wide above aspect ratio", 150, 74, LayoutVeryWide},
		{"tall takes precedence", 150, 75, LayoutTall},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveLayout(tt.width, tt.height); got != tt.want {
				t.Errorf("resolveLayout(%d, %d) = %v, want %v", tt.width, tt.height, got, tt.want)
			}
		})
	}
}
