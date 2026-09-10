package littleworld

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

func (m Screen) RenderTall(ctx screen.LayoutContext) tea.View {
	m.context = ctx
	return m.view(m.renderTall)
}

func (m Screen) renderTall() string {
	return lipgloss.NewStyle().
		Width(m.context.Width).
		Height(m.context.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Render("World Simulator: Tall layout coming soon...")
}
