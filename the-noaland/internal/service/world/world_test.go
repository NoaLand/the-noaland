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

	if appearance := worldTree.Express().Appearance; appearance != worldtree.Seed {
		t.Fatalf(
			"expected world tree appearance to be %s, got %s",
			worldtree.Seed,
			appearance,
		)
	}

	world.Step()
	if appearance := worldTree.Express().Appearance; appearance != worldtree.Sprout {
		t.Fatalf(
			"expected world tree appearance to be %s, got %s",
			worldtree.Sprout,
			appearance,
		)
	}
}
