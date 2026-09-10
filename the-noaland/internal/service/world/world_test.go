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
		t.Fatalf("Expected world tree state to be mature after 3 steps, got %s", worldTree.State())
	}
}
