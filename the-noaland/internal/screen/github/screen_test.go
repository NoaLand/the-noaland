package github

import (
	"testing"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

func TestGitHubScreenRetainsLoadedState(t *testing.T) {
	instance := New()
	ctx := screen.LayoutContext{Width: 90, Height: 24, Layout: screen.LayoutWide}
	if instance.RenderWide(ctx).AltScreen {
		t.Fatal("loading view should not enter alternate screen")
	}
	cmd := instance.Update(profileLoadedMsg{profile: Profile{Login: "noaland"}}, ctx)
	if cmd == nil || instance.profile == nil || instance.profile.Login != "noaland" || instance.lastUpdated.IsZero() {
		t.Fatal("profile update did not persist screen state or schedule redraw")
	}
	if !instance.RenderWide(ctx).AltScreen {
		t.Fatal("loaded view should enter alternate screen")
	}
}
