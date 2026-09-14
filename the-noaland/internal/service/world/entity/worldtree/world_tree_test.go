package worldtree

import (
	"testing"
)

func TestWorldTreeDeterministic(t *testing.T) {
	a := New("world-tree", 42)
	b := New("world-tree", 42)

	for step := range 100 {
		a.Step()
		b.Step()

		aExpr := a.Express()
		bExpr := b.Express()

		if aExpr.Appearance != bExpr.Appearance {
			t.Fatalf(
				"step %d: expected same appearance, got %s and %s",
				step,
				aExpr.Appearance,
				bExpr.Appearance,
			)
		}
	}
}
