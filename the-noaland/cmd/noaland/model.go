package main

import tea "charm.land/bubbletea/v2"

type model struct {
	width, height int
	screen        Screen
}

func initialModel() model {
	return model{screen: &githubScreen{}}
}

func (m model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m model) layoutContext() LayoutContext {
	return LayoutContext{
		Width:  m.width,
		Height: m.height,
		Layout: resolveLayout(m.width, m.height),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	return m, m.screen.Update(msg, m.layoutContext())
}

func (m model) View() tea.View {
	ctx := m.layoutContext()
	switch ctx.Layout {
	case LayoutCompactWide:
		return m.screen.RenderCompactWide(ctx)
	case LayoutWide:
		return m.screen.RenderWide(ctx)
	case LayoutVeryWide:
		return m.screen.RenderVeryWide(ctx)
	default:
		return m.screen.RenderTall(ctx)
	}
}
