package world

import (
	"testing"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

func TestWorldEvolves(t *testing.T) {
	world := New()
	worldTree := worldtree.New("world-tree")

	world.Add(worldTree)

	world.Step()
	world.Step()

	if worldTree.Appearance() != worldtree.Seed {
		t.Fatalf(
			"expected world tree appearance to be %s, got %s",
			worldtree.Seed,
			worldTree.Appearance(),
		)
	}

	world.Step()
	if worldTree.Appearance() != worldtree.Sprout {
		t.Fatalf(
			"expected world tree appearance to be %s, got %s",
			worldtree.Sprout,
			worldTree.Appearance(),
		)
	}
}
