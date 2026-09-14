package world_test

import (
	"testing"

	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

func TestAnnalRecordsAreChronologicalSnapshots(t *testing.T) {
	w := world.New(42)
	w.Add(worldtree.New("world-tree", "Silver-Oak", w.Seed()))
	if len(w.Annal().Records()) != 0 {
		t.Fatal("a new tree must not produce a change record")
	}
	for range 3 {
		w.Step()
	}
	snapshot := w.Annal().Records()
	want := world.Record{
		Time: 3, EntityID: "world-tree", EntityName: "Silver-Oak",
		Message: "Silver-Oak changed from Seed to Sprout.",
	}
	if len(snapshot) != 1 || snapshot[0] != want {
		t.Fatalf("unexpected public record: %+v", snapshot)
	}
	snapshot[0].Message = "changed by the caller"
	if got := w.Annal().Records()[0]; got != want {
		t.Fatalf("the caller changed the annal: %+v", got)
	}
	for range 100 {
		w.Step()
	}
	if len(snapshot) != 1 {
		t.Fatal("later steps changed an existing snapshot")
	}
	records := w.Annal().Records()
	if len(records) != 3 {
		t.Fatalf("expected three growth records, got %d", len(records))
	}
	for i := 1; i < len(records); i++ {
		if records[i].Time <= records[i-1].Time {
			t.Fatal("records are not chronological")
		}
	}
}
