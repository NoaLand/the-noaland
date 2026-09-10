package clock

import (
	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/ui/clock/flipclock"
)

func (s *Screen) RenderCompactWide(ctx screen.LayoutContext) tea.View {
	return s.view(ctx, s.horizontal(flipclock.Compact))
}
