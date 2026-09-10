package world

import "testing"

func TestWorldEvolves(t *testing.T) {
	world := New()
	seed := NewTree("seed")

	world.Add(seed)

	world.Step()
	world.Step()
	world.Step()

	if seed.state != 3 {
		t.Fatalf("Expected seed state to be 3 after 3 steps, got %d", seed.state)
	}
}
