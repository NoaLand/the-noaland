package littleworld

import (
	"math/rand/v2"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/service/world"
)

type tickMsg struct{}

// Screen owns one world shared by every layout. Rendering never advances it.
type Screen struct {
	context screen.LayoutContext
	world   *world.World
	started bool
}

func New() *Screen {
	w := world.New(rand.Uint64())
	w.Generate()
	return &Screen{world: w}
}

func (m *Screen) Init() tea.Cmd {
	if m.started {
		return nil
	}
	m.started = true
	return nextTick()
}

func (m *Screen) Update(msg tea.Msg, ctx screen.LayoutContext) tea.Cmd {
	m.context = ctx
	if _, ok := msg.(tickMsg); ok && m.started {
		m.world.Step()
		return nextTick()
	}
	return nil
}

// The simulation continues in the background, independent of the active layout.
func nextTick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg { return tickMsg{} })
}

func (m Screen) view(render func() string) tea.View {
	v := tea.NewView(render())
	v.AltScreen = true
	return v
}
