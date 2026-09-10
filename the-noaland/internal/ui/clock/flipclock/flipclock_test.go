package flipclock

import (
	"testing"

	"charm.land/lipgloss/v2"
)

func TestCardsKeepTheirDimensionsDuringEveryTransition(t *testing.T) {
	for _, size := range []Size{Compact, Regular, Large} {
		for old := byte('0'); old <= '9'; old++ {
			for next := byte('0'); next <= '9'; next++ {
				settled := digit(next, next, 1, size)
				width, height := lipgloss.Width(settled), lipgloss.Height(settled)
				for _, p := range []float64{0, 0.125, 0.25, 0.5, 0.75, 0.875, 1} {
					card := digit(old, next, p, size)
					if lipgloss.Width(card) != width || lipgloss.Height(card) != height {
						t.Fatal("animation changed card dimensions")
					}
					if old == next && card != settled {
						t.Fatal("unchanged digit animated")
					}
				}
				if digit(old, next, 0, size) != digit(old, old, 1, size) {
					t.Fatal("first frame must be old digit")
				}
				if digit(old, next, 1, size) != settled {
					t.Fatal("last frame must be new digit")
				}
			}
		}
		if lipgloss.Height(Colon(size)) != lipgloss.Height(digit('0', '0', 1, size)) {
			t.Fatal("colon not aligned")
		}
	}
}

func TestChangedDigitHasIntermediateFrames(t *testing.T) {
	old := digit('2', '2', 1, Regular)
	next := digit('3', '3', 1, Regular)
	middle := digit('2', '3', 0.5, Regular)
	if middle == old || middle == next {
		t.Fatal("flip animation has no intermediate state")
	}
}
