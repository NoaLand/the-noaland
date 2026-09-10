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
	m := model{screen: recorded}
	if got := m.Init()(); got != "init" {
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
		if cmd == nil || cmd() != "update" {
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
	m := model{screen: recorded}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})
	if cmd == nil {
		t.Fatal("quit command missing")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok || recorded.message != nil {
		t.Fatal("quit should be handled by model")
	}
	msg := tea.KeyPressMsg{Code: 'r'}
	_, cmd = m.Update(msg)
	if recorded.message != msg || cmd == nil || cmd() != "update" {
		t.Fatal("refresh key was not forwarded")
	}
}
