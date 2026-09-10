package screen

import tea "charm.land/bubbletea/v2"

// LayoutContext contains dimensions and the layout selected by the application.
type LayoutContext struct {
	Width, Height int
	Layout        Layout
}

// Screen requires an explicit renderer for every supported layout.
// Update mutates screen state and returns any asynchronous work.
type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg, LayoutContext) tea.Cmd

	RenderTall(LayoutContext) tea.View
	RenderCompactWide(LayoutContext) tea.View
	RenderWide(LayoutContext) tea.View
	RenderVeryWide(LayoutContext) tea.View
}

// ActivatedMsg asks a screen to redraw out-of-band content after page switching.
type ActivatedMsg struct{}
