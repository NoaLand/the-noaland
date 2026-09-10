package app

import (
	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	githubscreen "github.com/NoaLand/the-noaland/the-noaland/internal/screen/github"
)

var _ screen.Screen = (*githubscreen.Screen)(nil)

type model struct {
	width, height int
	screen        screen.Screen
}

// New creates the application with its initial screen.
func New() tea.Model {
	return model{screen: githubscreen.New()}
}

func (m model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m model) layoutContext() screen.LayoutContext {
	return screen.LayoutContext{
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
	case screen.LayoutCompactWide:
		return m.screen.RenderCompactWide(ctx)
	case screen.LayoutWide:
		return m.screen.RenderWide(ctx)
	case screen.LayoutVeryWide:
		return m.screen.RenderVeryWide(ctx)
	default:
		return m.screen.RenderTall(ctx)
	}
}
