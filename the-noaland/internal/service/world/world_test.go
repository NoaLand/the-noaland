package world

import (
	"testing"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

func TestWorldEvolves(t *testing.T) {
	world := New(42)
	worldTree := worldtree.New("world-tree", "Silver-Oak", world.Seed())

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

	for range 100 {
		world.Step()
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
	worldA.Generate()

	worldB := New(42)
	worldB.Generate()

	if len(worldA.entities) != len(worldB.entities) {
		t.Fatalf(
			"expected same number of entities when generates the world with the same seed, but got %d and %d",
			len(worldA.entities),
			len(worldB.entities),
		)
	}

	for i := range worldA.entities {
		a := worldA.entities[i]
		b := worldB.entities[i]

		if a.ID() != b.ID() || a.Name() != b.Name() {
			t.Fatalf(
				"expected same entity when generates the world with the same seed, got %s and %s",
				a.Name(),
				b.Name(),
			)
		}
	}

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

	t.Logf("World A\n")
	for _, record := range worldA.historiographer.annal.records {
		t.Logf(
			"Time #%d - %s: %s\n",
			record.time,
			record.entityId,
			record.message,
		)
	}
	t.Logf("\n\n")
	t.Logf("World B\n")
	for _, record := range worldB.historiographer.annal.records {
		t.Logf(
			"Time #%d - %s: %s\n",
			record.time,
			record.entityId,
			record.message,
		)
	}
}
