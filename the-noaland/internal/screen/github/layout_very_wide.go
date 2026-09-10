package github

import (
	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

func (m Screen) RenderVeryWide(ctx screen.LayoutContext) tea.View {
	m.context = ctx
	return m.view(func() string { return m.renderFullWide(m.renderYearHeatmap()) })
}
