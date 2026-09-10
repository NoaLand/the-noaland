package world

import "testing"

func TestWorldEvolves(t *testing.T) {
	world := New()
	worldTree := NewWorldTree("world-tree")

	world.Add(worldTree)

	world.Step()
	world.Step()
	world.Step()

	if worldTree.state != 3 {
		t.Fatalf("Expected seed state to be 3 after 3 steps, got %d", worldTree.state)
	}
}
