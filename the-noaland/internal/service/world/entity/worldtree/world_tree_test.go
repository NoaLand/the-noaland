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

		if a.state != b.state {
			t.Fatalf(
				"step %d: expected same state, got %+v and %+v",
				step,
				a.state,
				b.state,
			)
		}

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

func TestWorldTreesWithDifferentIDsEvolveDifferently(t *testing.T) {
	a := New("world-tree-A", 42)
	b := New("world-tree-B", 42)
	diverged := false

	for range 100 {
		a.Step()
		b.Step()

		if a.state.growth != b.state.growth {
			diverged = true
			break
		}
	}

	if !diverged {
		t.Fatalf(
			"Expected different state, got %+v and %+v",
			a.state.growth,
			b.state.growth,
		)
	}
}
