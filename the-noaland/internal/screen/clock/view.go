package clock

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/ui/clock/flipclock"
)

func (s *Screen) pair(start int, size flipclock.Size) string {
	old, current := s.previous.Format("150405"), s.current.Format("150405")
	return flipclock.Pair(old[start:start+2], current[start:start+2], s.progress, size)
}

func (s *Screen) horizontal(size flipclock.Size) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, s.pair(0, size), flipclock.Colon(size),
		s.pair(2, size), flipclock.Colon(size), s.pair(4, size))
}

func (s *Screen) view(ctx screen.LayoutContext, body string) tea.View {
	if ctx.Width <= 0 || ctx.Height <= 0 {
		v := tea.NewView("")
		v.AltScreen = true
		return v
	}
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#7DCFFF")).Bold(true).Render("FLIP CLOCK")
	date := lipgloss.NewStyle().Foreground(lipgloss.Color("#A9B1D6")).Render(s.current.Format("Monday · 02 January 2006"))
	zone := lipgloss.NewStyle().Foreground(lipgloss.Color("#565F89")).Render(s.current.Format("MST · UTC-07:00"))
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#565F89")).Render("← / →  switch screen    q  quit")
	content := lipgloss.JoinVertical(lipgloss.Center, title, "", body, "", date, zone, "", hint)
	// Keep tiny windows usable without clipping a digit card in half.
	if lipgloss.Width(content) > ctx.Width || lipgloss.Height(content) > ctx.Height {
		content = s.current.Format("15:04:05")
	}
	width, height := max(0, ctx.Width), max(0, ctx.Height)
	content = lipgloss.NewStyle().MaxWidth(width).MaxHeight(height).Render(content)
	v := tea.NewView(lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content))
	v.AltScreen = true
	return v
}
