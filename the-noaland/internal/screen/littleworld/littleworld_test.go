package littleworld

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

func TestLayoutsShareWorldWithoutAdvancingIt(t *testing.T) {
	w := world.New(42)
	tree := worldtree.New("world-tree", "Silver-Oak", w.Seed())
	w.Add(tree)
	s := &Screen{world: w}
	ctx := screen.LayoutContext{Width: 160, Height: 30}
	if s.Init() == nil || s.Init() != nil {
		t.Fatal("initialization must start exactly one tick loop")
	}
	for range 3 {
		if s.Update(tickMsg{}, ctx) == nil {
			t.Fatal("tick must schedule the next step")
		}
	}
	before := tree.Express().Appearance
	if before != worldtree.Sprout {
		t.Fatal("ticks did not advance the world")
	}
	for _, render := range []func(screen.LayoutContext) tea.View{
		s.RenderTall, s.RenderCompactWide, s.RenderWide, s.RenderVeryWide,
	} {
		render(ctx)
		if s.world != w || s.ticks != 3 || tree.Express().Appearance != before {
			t.Fatal("rendering changed the shared world")
		}
	}
	for _, msg := range []tea.Msg{
		screen.DeactivatedMsg{}, screen.ActivatedMsg{}, tea.WindowSizeMsg{Width: 150, Height: 24},
	} {
		if s.Update(msg, ctx) != nil || s.world != w || s.ticks != 3 {
			t.Fatal("page or layout changes restarted the world or timer")
		}
	}
}

func TestWorldContinuesWhileHidden(t *testing.T) {
	s := New()
	ctx := screen.LayoutContext{}
	if s.Update(tickMsg{}, ctx) != nil || s.ticks != 0 {
		t.Fatal("an uninitialized screen must not advance")
	}
	s.Init()
	s.Update(screen.DeactivatedMsg{}, ctx)
	if s.Update(tickMsg{}, ctx) == nil || s.ticks != 1 {
		t.Fatal("hidden world must keep advancing")
	}
}
