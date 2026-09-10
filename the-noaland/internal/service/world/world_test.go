package world

import "testing"

func TestWorldEvolves(t *testing.T) {
	world := New()
	world_tree := NewTree("world-tree")

	world.Add(world_tree)

	world.Step()
	world.Step()
	world.Step()

	if world_tree.state != 3 {
		t.Fatalf("Expected seed state to be 3 after 3 steps, got %d", world_tree.state)
	}
}
