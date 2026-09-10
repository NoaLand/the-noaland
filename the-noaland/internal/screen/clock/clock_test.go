package clock

import (
	"math"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

func TestMidnightAnimationAndCatchUp(t *testing.T) {
	now := time.Date(2026, 12, 31, 23, 59, 59, 800000000, time.UTC)
	s := New()
	s.now = func() time.Time { return now }
	ctx := screen.LayoutContext{}
	if s.Init() != nil {
		t.Fatal("hidden clock must not start timers")
	}
	if s.Update(screen.ActivatedMsg{}, ctx) == nil {
		t.Fatal("activation must start clock")
	}
	generation := s.generation
	now = now.Add(200 * time.Millisecond)
	if cmd := s.Update(tickMsg{generation}, ctx); cmd == nil {
		t.Fatal("tick must schedule animation and next second")
	}
	if s.previous.Format("150405") != "235959" || s.current.Format("150405") != "000000" || s.current.Year() != 2027 || s.progress != 0 {
		t.Fatal("midnight transition lost old/new time")
	}
	now = now.Add(animationDuration / 2)
	if s.Update(frameMsg{generation}, ctx) == nil || math.Abs(s.progress-0.5) > 0.001 {
		t.Fatal("animation midpoint incorrect")
	}
	now = now.Add(time.Second)
	if s.Update(frameMsg{generation}, ctx) != nil || s.progress != 1 {
		t.Fatal("late frame should finish animation")
	}
	// A delayed second tick must read wall time rather than increment a counter.
	now = now.Add(10 * time.Second)
	s.Update(tickMsg{generation}, ctx)
	if !s.current.Equal(now) {
		t.Fatal("clock drifted after delayed tick")
	}
}

func TestDeactivationStopsTimersAndRejectsOldMessages(t *testing.T) {
	now := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	s := New()
	s.now = func() time.Time { return now }
	ctx := screen.LayoutContext{}
	s.Update(screen.ActivatedMsg{}, ctx)
	old := s.generation
	s.Update(screen.DeactivatedMsg{}, ctx)
	if s.active || s.Update(tickMsg{old}, ctx) != nil || s.Update(frameMsg{old}, ctx) != nil {
		t.Fatal("hidden clock still animates")
	}
	now = now.Add(time.Hour)
	s.Update(screen.ActivatedMsg{}, ctx)
	if !s.current.Equal(now) || s.progress != 1 {
		t.Fatal("reactivation did not resync")
	}
	if s.Update(tickMsg{old}, ctx) != nil {
		t.Fatal("old timer started duplicate loop")
	}
}

func TestClockLayoutsFitViewport(t *testing.T) {
	s := New()
	s.current = time.Date(2026, 9, 10, 12, 34, 56, 0, time.UTC)
	s.previous = s.current.Add(-time.Second)
	for _, tt := range []struct {
		name          string
		width, height int
		render        func(screen.LayoutContext) tea.View
	}{
		{"wide", 90, 24, s.RenderWide},
		{"compact", 90, 20, s.RenderCompactWide},
		{"very wide", 150, 24, s.RenderVeryWide},
		{"tall", 40, 35, s.RenderTall},
		{"short tall", 80, 24, s.RenderTall},
		{"tiny", 5, 2, s.RenderTall},
		{"initial", 0, 0, s.RenderTall},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx := screen.LayoutContext{Width: tt.width, Height: tt.height}
			for _, progress := range []float64{0, 0.25, 0.5, 0.75, 1} {
				s.progress = progress
				v := tt.render(ctx)
				if !v.AltScreen {
					t.Fatal("clock must use alternate screen")
				}
				if tt.width == 0 {
					if v.Content != "" {
						t.Fatal("initial viewport not empty")
					}
					continue
				}
				if lipgloss.Width(v.Content) > tt.width || lipgloss.Height(v.Content) > tt.height {
					t.Fatal("clock exceeds viewport")
				}
				if tt.width >= 40 && !strings.Contains(ansi.Strip(v.Content), "FLIP CLOCK") {
					t.Fatal("normal layout fell back to plain text")
				}
			}
		})
	}
}
