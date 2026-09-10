package world

import "testing"

func TestWorldEvolves(t *testing.T) {
	world := New()
	worldTree := NewWorldTree("world-tree")

	world.Add(worldTree)

	world.Step()
	world.Step()
	world.Step()

	if worldTree.State() != WorldTreeMature {
		t.Fatalf(
			"expected world tree state to be %s, got %s",
			WorldTreeMature,
			worldTree.State(),
		)
	}
}
