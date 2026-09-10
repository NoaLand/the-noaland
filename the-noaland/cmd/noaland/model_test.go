package main

import (
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"
)

type recordingScreen struct {
	context LayoutContext
	message tea.Msg
	called  string
}

func (s *recordingScreen) Init() tea.Cmd {
	return func() tea.Msg { return "init" }
}

func (s *recordingScreen) Update(msg tea.Msg, ctx LayoutContext) tea.Cmd {
	s.message, s.context = msg, ctx
	return func() tea.Msg { return "update" }
}

func (s *recordingScreen) record(name string, ctx LayoutContext) tea.View {
	s.called, s.context = name, ctx
	return tea.NewView(name)
}

func (s *recordingScreen) RenderTall(ctx LayoutContext) tea.View {
	return s.record("tall", ctx)
}
func (s *recordingScreen) RenderCompactWide(ctx LayoutContext) tea.View {
	return s.record("compact", ctx)
}
func (s *recordingScreen) RenderWide(ctx LayoutContext) tea.View {
	return s.record("wide", ctx)
}
func (s *recordingScreen) RenderVeryWide(ctx LayoutContext) tea.View {
	return s.record("very wide", ctx)
}

func TestModelDispatchesLayoutAfterResize(t *testing.T) {
	screen := &recordingScreen{}
	m := model{screen: screen}
	if got := m.Init()(); got != "init" {
		t.Fatalf("Init command = %v", got)
	}
	for _, tt := range []struct {
		width, height int
		layout        Layout
		renderer      string
	}{
		{80, 24, LayoutTall, "tall"},
		{150, 23, LayoutCompactWide, "compact"},
		{90, 24, LayoutWide, "wide"},
		{150, 24, LayoutVeryWide, "very wide"},
		{150, 75, LayoutTall, "tall"},
	} {
		msg := tea.WindowSizeMsg{Width: tt.width, Height: tt.height}
		next, cmd := m.Update(msg)
		m = next.(model)
		wantContext := LayoutContext{Width: tt.width, Height: tt.height, Layout: tt.layout}
		if screen.context != wantContext || screen.message != msg {
			t.Fatalf("resize forwarded with stale context: %+v", screen.context)
		}
		if cmd == nil || cmd() != "update" {
			t.Fatal("screen command was not forwarded")
		}
		view := m.View()
		if screen.called != tt.renderer || screen.context != wantContext {
			t.Fatalf("wrong renderer/context: %s, %+v", screen.called, screen.context)
		}
		if !reflect.DeepEqual(view, tea.NewView(tt.renderer)) {
			t.Fatal("screen view was changed by model")
		}
	}
}

func TestModelHandlesQuitAndForwardsOtherMessages(t *testing.T) {
	screen := &recordingScreen{}
	m := model{screen: screen}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'q'})
	if cmd == nil {
		t.Fatal("quit command missing")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok || screen.message != nil {
		t.Fatal("quit should be handled by model")
	}
	msg := tea.KeyPressMsg{Code: 'r'}
	_, cmd = m.Update(msg)
	if screen.message != msg || cmd == nil || cmd() != "update" {
		t.Fatal("refresh key was not forwarded")
	}
}

func TestGitHubScreenRetainsLoadedState(t *testing.T) {
	screen := &githubScreen{}
	ctx := LayoutContext{Width: 90, Height: 24, Layout: LayoutWide}
	if screen.RenderWide(ctx).AltScreen {
		t.Fatal("loading view should not enter alternate screen")
	}
	cmd := screen.Update(profileLoadedMsg{profile: Profile{Login: "noaland"}}, ctx)
	if cmd == nil || screen.profile == nil || screen.profile.Login != "noaland" || screen.lastUpdated.IsZero() {
		t.Fatal("profile update did not persist screen state or schedule redraw")
	}
	if !screen.RenderWide(ctx).AltScreen {
		t.Fatal("loaded view should enter alternate screen")
	}
}
