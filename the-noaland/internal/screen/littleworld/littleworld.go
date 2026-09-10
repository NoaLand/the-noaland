package littleworld

import (
	tea "charm.land/bubbletea/v2"
	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

type Screen struct{
	context screen.LayoutContext
	err error
}

func New() *Screen {
	return &Screen{}
}

func (m Screen) Init() tea.Cmd {
	return nil
}

func (m Screen) Update(tea.Msg, screen.LayoutContext) tea.Cmd {
	return nil
}

func (m Screen) view(render func() string) tea.View {
	v := tea.NewView(render())
	v.AltScreen = true

	return v
}
