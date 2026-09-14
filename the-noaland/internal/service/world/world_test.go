package world

import (
	"testing"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

func TestWorldEvolves(t *testing.T) {
	world := New(42)
	worldTree := worldtree.New("world-tree", world.Seed())

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

	for _, record := range world.historiographer.annal.records {
		t.Logf(
			"Time #%d - %s: %s\n",
			record.time,
			record.entityId,
			record.message,
		)
	}
}

func TestWorldDeterministic(t *testing.T) {
	worldA := New(42)
	worldTreeA := worldtree.New("world-tree", worldA.Seed())
	worldA.Add(worldTreeA)

	worldB := New(42)
	worldTreeB := worldtree.New("world-tree", worldB.Seed())
	worldB.Add(worldTreeB)

	for range 100 {
		worldA.Step()
		worldB.Step()
	}

	if len(worldA.historiographer.annal.records) != len(worldB.historiographer.annal.records) {
		t.Fatalf(
			"expected same number of records, got %d and %d",
			len(worldA.historiographer.annal.records),
			len(worldB.historiographer.annal.records),
		)
	}

	for i := range worldA.historiographer.annal.records {
		a := worldA.historiographer.annal.records[i]
		b := worldB.historiographer.annal.records[i]

		if a != b {
			t.Fatalf(
				"record %d differs: got %+v and %+v",
				i,
				a,
				b,
			)
		}
	}
}
