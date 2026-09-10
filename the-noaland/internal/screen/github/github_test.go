package github

import (
	"fmt"
	"image"
	"strings"
	"testing"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	githubservice "github.com/NoaLand/the-noaland/the-noaland/internal/service/github"
)

func TestGitHubScreenRetainsLoadedState(t *testing.T) {
	instance := New()
	ctx := screen.LayoutContext{Width: 90, Height: 24, Layout: screen.LayoutWide}
	if instance.RenderWide(ctx).AltScreen {
		t.Fatal("loading view should not enter alternate screen")
	}
	cmd := instance.Update(profileLoadedMsg{profile: githubservice.Profile{Login: "noaland"}}, ctx)
	if cmd == nil || instance.profile == nil || instance.profile.Login != "noaland" || instance.lastUpdated.IsZero() {
		t.Fatal("profile update did not persist screen state or schedule redraw")
	}
	if !instance.RenderWide(ctx).AltScreen {
		t.Fatal("loaded view should enter alternate screen")
	}
}

func TestSetupMessagesAndRetryRecovery(t *testing.T) {
	ctx := screen.LayoutContext{Width: 90, Height: 24, Layout: screen.LayoutWide}
	for _, tt := range []struct {
		err     error
		message string
	}{
		{githubservice.ErrCLINotInstalled, "https://cli.github.com"},
		{fmt.Errorf("wrapped: %w", githubservice.ErrNotAuthenticated), "gh auth login"},
	} {
		instance := New()
		instance.Update(errMsg{tt.err}, ctx)
		for _, view := range []string{instance.RenderWide(ctx).Content, instance.RenderCompactWide(ctx).Content, instance.RenderVeryWide(ctx).Content, instance.RenderTall(ctx).Content} {
			for _, text := range []string{tt.message, "Flip Clock", "retry", "→"} {
				if !strings.Contains(view, text) {
					t.Fatalf("missing setup hint %q", text)
				}
			}
		}
		instance.avatar = image.NewRGBA(image.Rect(0, 0, 2, 2))
		if instance.renderAvatarCmd() != nil {
			t.Fatal("old avatar must not cover setup message")
		}
		instance.Update(profileLoadedMsg{profile: githubservice.Profile{Login: "ready"}}, ctx)
		if instance.err != nil || strings.Contains(instance.RenderWide(ctx).Content, tt.message) {
			t.Fatal("successful retry did not clear error")
		}
	}
}
