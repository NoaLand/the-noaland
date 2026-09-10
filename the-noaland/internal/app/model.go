package app

import (
	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	clockscreen "github.com/NoaLand/the-noaland/the-noaland/internal/screen/clock"
	githubscreen "github.com/NoaLand/the-noaland/the-noaland/internal/screen/github"
	littleworldscreen "github.com/NoaLand/the-noaland/the-noaland/internal/screen/littleworld"
)

var _ screen.Screen = (*githubscreen.Screen)(nil)
var _ screen.Screen = (*clockscreen.Screen)(nil)
var _ screen.Screen = (*littleworldscreen.Screen)(nil)

type model struct {
	width, height int
	screens       []screen.Screen
	active        int
	generation    uint64
}

// screenMsg routes asynchronous results back to their owner. The generation
// prevents delayed raw graphics from a previous visit drawing on the new page.
type screenMsg struct {
	index      int
	generation uint64
	message    tea.Msg
}

func screenCommand(index int, generation uint64, cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	return func() tea.Msg { return screenMsg{index, generation, cmd()} }
}

// New creates the application with its registered screens.
func New() tea.Model {
	return model{screens: []screen.Screen{
		githubscreen.New(),
		clockscreen.New(),
		littleworldscreen.New(),
	}}
}

func (m model) Init() tea.Cmd {
	commands := make([]tea.Cmd, 0, len(m.screens))
	for i, s := range m.screens {
		commands = append(commands, screenCommand(i, m.generation, s.Init()))
	}
	if len(m.screens) > 0 {
		commands = append(commands, m.updateScreen(m.active, screen.ActivatedMsg{}))
	}
	return tea.Batch(commands...)
}

func (m model) layoutContext() screen.LayoutContext {
	return screen.LayoutContext{Width: m.width, Height: m.height, Layout: resolveLayout(m.width, m.height)}
}

func (m model) updateScreen(index int, msg tea.Msg) tea.Cmd {
	return screenCommand(index, m.generation, m.screens[index].Update(msg, m.layoutContext()))
}

func (m model) changeScreen(delta int) (tea.Model, tea.Cmd) {
	if len(m.screens) < 2 {
		return m, nil
	}
	deactivate := m.updateScreen(m.active, screen.DeactivatedMsg{})
	m.active = (m.active + delta + len(m.screens)) % len(m.screens)
	m.generation++
	return m, tea.Batch(tea.ClearScreen, deactivate, m.updateScreen(m.active, screen.ActivatedMsg{}))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "left":
			return m.changeScreen(-1)
		case "right":
			return m.changeScreen(1)
		}
	case screenMsg:
		if msg.index < 0 || msg.index >= len(m.screens) {
			return m, nil
		}
		switch result := msg.message.(type) {
		case tea.BatchMsg:
			commands := make([]tea.Cmd, 0, len(result))
			for _, cmd := range result {
				commands = append(commands, screenCommand(msg.index, msg.generation, cmd))
			}
			return m, tea.Batch(commands...)
		case tea.RawMsg:
			if msg.index != m.active || msg.generation != m.generation {
				return m, nil
			}
			return m, func() tea.Msg { return result }
		case nil:
			return m, nil
		default:
			return m, m.updateScreen(msg.index, result)
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}
	if len(m.screens) == 0 {
		return m, nil
	}
	return m, m.updateScreen(m.active, msg)
}

func (m model) View() tea.View {
	if len(m.screens) == 0 {
		return tea.NewView("")
	}
	active := m.screens[m.active]
	ctx := m.layoutContext()
	switch ctx.Layout {
	case screen.LayoutCompactWide:
		return active.RenderCompactWide(ctx)
	case screen.LayoutWide:
		return active.RenderWide(ctx)
	case screen.LayoutVeryWide:
		return active.RenderVeryWide(ctx)
	default:
		return active.RenderTall(ctx)
	}
}
