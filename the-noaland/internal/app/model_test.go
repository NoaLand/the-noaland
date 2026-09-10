package app

import (
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

type recordingScreen struct {
	context screen.LayoutContext
	message tea.Msg
	called  string
}

func (s *recordingScreen) Init() tea.Cmd {
	return func() tea.Msg { return "init" }
}

func (s *recordingScreen) Update(msg tea.Msg, ctx screen.LayoutContext) tea.Cmd {
	s.message, s.context = msg, ctx
	return func() tea.Msg { return "update" }
}

func (s *recordingScreen) record(name string, ctx screen.LayoutContext) tea.View {
	s.called, s.context = name, ctx
	return tea.NewView(name)
}

func (s *recordingScreen) RenderTall(ctx screen.LayoutContext) tea.View {
	return s.record("tall", ctx)
}
func (s *recordingScreen) RenderCompactWide(ctx screen.LayoutContext) tea.View {
	return s.record("compact", ctx)
}
func (s *recordingScreen) RenderWide(ctx screen.LayoutContext) tea.View {
	return s.record("wide", ctx)
}
func (s *recordingScreen) RenderVeryWide(ctx screen.LayoutContext) tea.View {
	return s.record("very wide", ctx)
}

func TestModelDispatchesLayoutAfterResize(t *testing.T) {
	recorded := &recordingScreen{}
	m := model{screens: []screen.Screen{recorded}}
	if got := m.Init()().(tea.BatchMsg)[0]().(screenMsg).message; got != "init" {
		t.Fatalf("Init command = %v", got)
	}
	for _, tt := range []struct {
		width, height int
		layout        screen.Layout
		renderer      string
	}{
		{80, 24, screen.LayoutTall, "tall"},
		{150, 23, screen.LayoutCompactWide, "compact"},
		{90, 24, screen.LayoutWide, "wide"},
		{150, 24, screen.LayoutVeryWide, "very wide"},
		{150, 75, screen.LayoutTall, "tall"},
	} {
		msg := tea.WindowSizeMsg{Width: tt.width, Height: tt.height}
		next, cmd := m.Update(msg)
		m = next.(model)
		wantContext := screen.LayoutContext{Width: tt.width, Height: tt.height, Layout: tt.layout}
		if recorded.context != wantContext || recorded.message != msg {
			t.Fatalf("resize forwarded with stale context: %+v", recorded.context)
		}
		if cmd == nil || cmd().(screenMsg).message != "update" {
			t.Fatal("screen command was not forwarded")
		}
		view := m.View()
		if recorded.called != tt.renderer || recorded.context != wantContext {
			t.Fatalf("wrong renderer/context: %s, %+v", recorded.called, recorded.context)
		}
		if !reflect.DeepEqual(view, tea.NewView(tt.renderer)) {
			t.Fatal("screen view was changed by model")
		}
	}
}

func TestModelHandlesQuitAndForwardsOtherMessages(t *testing.T) {
	recorded := &recordingScreen{}
	m := model{screens: []screen.Screen{recorded}}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})
	if cmd == nil {
		t.Fatal("quit command missing")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok || recorded.message != nil {
		t.Fatal("quit should be handled by model")
	}
	msg := tea.KeyPressMsg{Code: 'r'}
	_, cmd = m.Update(msg)
	if recorded.message != msg || cmd == nil || cmd().(screenMsg).message != "update" {
		t.Fatal("refresh key was not forwarded")
	}
}

func TestScreenCyclingAndEmptyList(t *testing.T) {
	first, second := &recordingScreen{}, &recordingScreen{}
	m := model{screens: []screen.Screen{first, second}, width: 150, height: 24}
	for _, tt := range []struct {
		key   rune
		index int
	}{{tea.KeyRight, 1}, {tea.KeyRight, 0}, {tea.KeyLeft, 1}, {tea.KeyLeft, 0}} {
		next, cmd := m.Update(tea.KeyPressMsg{Code: tt.key})
		m = next.(model)
		if m.active != tt.index || cmd == nil {
			t.Fatalf("switch failed: %d", m.active)
		}
		selected := m.screens[m.active].(*recordingScreen)
		if _, ok := selected.message.(screen.ActivatedMsg); !ok {
			t.Fatal("activation not delivered")
		}
		if selected.context.Width != 150 {
			t.Fatal("activation has stale size")
		}
	}
	single := model{screens: []screen.Screen{first}}
	if _, cmd := single.Update(tea.KeyPressMsg{Code: tea.KeyRight}); cmd != nil {
		t.Fatal("single page should not switch")
	}
	empty := model{}
	if empty.Init() != nil {
		t.Fatal("empty init")
	}
	empty.View()
	if _, cmd := empty.Update(tea.KeyPressMsg{Code: tea.KeyLeft}); cmd != nil {
		t.Fatal("empty switch")
	}
}

func TestAsyncResultsStayWithOwningScreen(t *testing.T) {
	first, second := &recordingScreen{}, &recordingScreen{}
	m := model{screens: []screen.Screen{first, second}, active: 1, generation: 2, width: 90, height: 24}
	_, cmd := m.Update(screenMsg{index: 0, generation: 0, message: "loaded"})
	if first.message != "loaded" || second.message != nil || cmd == nil {
		t.Fatal("result routed to wrong page")
	}
	raw := tea.Raw("old image")()
	for _, msg := range []screenMsg{{0, 2, raw}, {1, 1, raw}} {
		if _, cmd := m.Update(msg); cmd != nil {
			t.Fatal("stale/hidden graphics escaped")
		}
	}
	if _, cmd := m.Update(screenMsg{1, 2, raw}); cmd == nil {
		t.Fatal("active graphics suppressed")
	}
}

func TestBatchCommandsKeepOwnership(t *testing.T) {
	m := model{screens: []screen.Screen{&recordingScreen{}, &recordingScreen{}}, active: 1}
	batch := tea.BatchMsg{func() tea.Msg { return "one" }, func() tea.Msg { return "two" }}
	_, cmd := m.Update(screenMsg{0, 0, batch})
	commands := cmd().(tea.BatchMsg)
	for _, command := range commands {
		msg := command().(screenMsg)
		if msg.index != 0 {
			t.Fatal("batch lost originating screen")
		}
	}
}

func TestRegisteredClockAndRoundTripNavigation(t *testing.T) {
	m := New().(model)
	if len(m.screens) != 2 {
		t.Fatal("clock not registered")
	}
	m.width, m.height = 90, 24
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	m = next.(model)
	if m.active != 1 || !m.View().AltScreen {
		t.Fatal("clock not reachable")
	}
	next, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	m = next.(model)
	if m.active != 0 {
		t.Fatal("GitHub not reachable on return")
	}
}

func TestLeavingScreenReceivesDeactivation(t *testing.T) {
	first, second := &recordingScreen{}, &recordingScreen{}
	m := model{screens: []screen.Screen{first, second}}
	m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if _, ok := first.message.(screen.DeactivatedMsg); !ok {
		t.Fatal("background work not notified")
	}
	if _, ok := second.message.(screen.ActivatedMsg); !ok {
		t.Fatal("new screen not activated")
	}
}
