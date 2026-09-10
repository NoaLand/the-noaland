package clock

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/ui/clock/flipclock"
)

func (s *Screen) RenderTall(ctx screen.LayoutContext) tea.View {
	if ctx.Height >= 29 {
		return s.view(ctx, lipgloss.JoinVertical(lipgloss.Center, s.pair(0, flipclock.Compact),
			s.pair(2, flipclock.Compact), s.pair(4, flipclock.Compact)))
	}
	body := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Top, s.pair(0, flipclock.Compact), flipclock.Colon(flipclock.Compact), s.pair(2, flipclock.Compact)),
		s.current.Format("05")+" seconds")
	return s.view(ctx, body)
}
