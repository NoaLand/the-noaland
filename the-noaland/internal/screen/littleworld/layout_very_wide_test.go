package littleworld

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world/entity/worldtree"
)

func TestVeryWideShowsAnnalWithoutChangingWorld(t *testing.T) {
	w := world.New(42)
	w.Add(worldtree.New("world-tree", "Silver-Oak", w.Seed()))
	s := &Screen{world: w}
	s.Init()
	ctx := screen.LayoutContext{Width: 150, Height: 24}
	if content := s.RenderVeryWide(ctx).Content; !strings.Contains(content, "Waiting for the first change") {
		t.Fatal("empty annal must explain that it is waiting for changes")
	}
	for range 100 {
		s.Update(tickMsg{}, ctx)
	}
	records := w.Annal().Records()
	view := s.RenderVeryWide(ctx)
	content := ansi.Strip(view.Content)
	lastPosition := -1
	for _, record := range records {
		position := strings.Index(content, record.Message)
		if position <= lastPosition || !strings.Contains(content, fmt.Sprintf("#%d", record.Time)) {
			t.Fatalf("missing or out-of-order record: %+v", record)
		}
		lastPosition = position
	}
	if !strings.Contains(content, "#100") || !strings.Contains(content, "42") {
		t.Fatal("world time and seed must be visible")
	}
	if !view.AltScreen || !slices.Equal(records, w.Annal().Records()) || s.world.Time() != 100 {
		t.Fatal("rendering changed the world or disabled the alternate screen")
	}
}

func TestVeryWideKeepsLatestRecordsWithinViewport(t *testing.T) {
	w := world.New(42)
	for i := range 30 {
		name := fmt.Sprintf("Tree-%02d", i)
		w.Add(worldtree.New(entity.EntityID(name), entity.EntityName(name), w.Seed()))
	}
	// Exercise cell-width truncation for wide Unicode names as well.
	w.Add(worldtree.New("long-name", entity.EntityName(strings.Repeat("世界树", 30)), w.Seed()))
	for range 100 {
		w.Step()
	}
	s := &Screen{world: w}
	records := w.Annal().Records()
	for _, size := range [][2]int{{150, 24}, {224, 57}, {300, 80}, {1, 1}, {0, 0}} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			ctx := screen.LayoutContext{Width: size[0], Height: size[1]}
			view := s.RenderVeryWide(ctx)
			if lipgloss.Width(view.Content) > ctx.Width || (ctx.Height > 0 && lipgloss.Height(view.Content) > ctx.Height) {
				t.Fatalf("view exceeds %dx%d: got %dx%d", ctx.Width, ctx.Height,
					lipgloss.Width(view.Content), lipgloss.Height(view.Content))
			}
			if ctx.Width == 0 && view.Content != "" {
				t.Fatal("uninitialized viewport must be empty")
			}
			if ctx.Width >= 150 {
				content := ansi.Strip(view.Content)
				latest := records[len(records)-1]
				namePrefix := string([]rune(latest.EntityName)[:7])
				found := false
				for _, line := range strings.Split(content, "\n") {
					if strings.Contains(line, fmt.Sprintf("#%d", latest.Time)) && strings.Contains(line, namePrefix) {
						found = true
					}
				}
				if !strings.Contains(content, "latest") || !found {
					t.Fatal("overflow must retain the latest record and indicate hidden history")
				}
				if strings.Contains(content, records[0].Message) {
					t.Fatal("old records must give way to the latest records")
				}
			}
		})
	}
}
