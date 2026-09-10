package world

import "testing"

func TestWorldEvolves(t *testing.T) {
	world := New()
	worldTree := NewWorldTree("world-tree")

	world.Add(worldTree)

	world.Step()
	world.Step()

	if worldTree.Appearance() != WorldTreeSeed {
		t.Fatalf(
			"expected world tree state to be %s, got %s",
			WorldTreeSeed,
			worldTree.Appearance(),
		)
	}

	world.Step()
	if worldTree.Appearance() != WorldTreeSprout {
		t.Fatalf(
			"expected world tree state to be %s, got %s",
			WorldTreeSprout,
			worldTree.Appearance(),
		)
	}
}
