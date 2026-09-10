package image

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestProtocolSelection(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"ghostty", map[string]string{"TERM_PROGRAM": "ghostty"}, "kitty"},
		{"ghostty zellij", map[string]string{"TERM": "xterm-ghostty", "ZELLIJ": "0"}, "kitty"},
		{"kitty", map[string]string{"TERM": "xterm-kitty"}, "kitty"},
		{"windows terminal", map[string]string{"WT_SESSION": "session"}, "sixel"},
		{"vscode", map[string]string{"TERM_PROGRAM": "vscode"}, "blocks"},
		{"unknown", nil, "blocks"},
		{"tmux", map[string]string{"TERM_PROGRAM": "ghostty", "TMUX": "session"}, "blocks"},
		{"override", map[string]string{"NOALAND_IMAGE_PROTOCOL": "kitty", "TMUX": "session"}, "kitty"},
		{"force blocks", map[string]string{"NOALAND_IMAGE_PROTOCOL": "blocks", "TERM_PROGRAM": "ghostty"}, "blocks"},
		{"disabled", map[string]string{"NOALAND_IMAGE_PROTOCOL": "none"}, "none"},
		{"invalid override", map[string]string{"NOALAND_IMAGE_PROTOCOL": "typo"}, "blocks"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectProtocolForOS(func(key string) string { return tt.env[key] }, "linux"); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBlocksUseInlineCellsAndRespectBounds(t *testing.T) {
	t.Setenv("NOALAND_IMAGE_PROTOCOL", "blocks")
	img := image.NewNRGBA(image.Rect(5, 7, 7, 11))
	for y := 7; y < 11; y++ {
		for x := 5; x < 7; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: 255, A: 255})
		}
	}
	result := Render(img, 3, 2)
	if result.Raw != "" {
		t.Fatal("fallback must not send graphics outside normal UI")
	}
	lines := strings.Split(result.Inline, "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d rows", len(lines))
	}
	for _, line := range lines {
		if strings.Count(line, "▀") != 3 || !strings.HasSuffix(line, "\x1b[0m") {
			t.Fatalf("incorrect row width/reset: %q", line)
		}
		if !strings.Contains(line, "\x1b[38;2;255;0;0m") {
			t.Fatal("source bounds or color sampling changed")
		}
	}
}

func TestDisabledInvalidAndKittyResults(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	t.Setenv("NOALAND_IMAGE_PROTOCOL", "none")
	if Render(img, 2, 2) != (Result{}) {
		t.Fatal("disabled renderer produced output")
	}
	t.Setenv("NOALAND_IMAGE_PROTOCOL", "kitty")
	result := Render(img, 2, 2)
	if result.Inline != "" || result.Raw != renderKitty(img, 2, 2) {
		t.Fatal("Kitty output changed")
	}
	for _, size := range [][2]int{{0, 2}, {2, 0}, {-1, 2}} {
		if Render(img, size[0], size[1]) != (Result{}) {
			t.Fatal("invalid dimensions produced output")
		}
	}
	if Render(nil, 2, 2) != (Result{}) || Render(image.NewRGBA(image.Rectangle{}), 2, 2) != (Result{}) {
		t.Fatal("missing image produced output")
	}
}

func TestWindowsWithoutTerminalIdentity(t *testing.T) {
	for _, tt := range []struct {
		name, goos string
		env        map[string]string
		want       string
	}{
		{"reported PowerShell environment", "windows", nil, "sixel"},
		{"explicit auto", "windows", map[string]string{"NOALAND_IMAGE_PROTOCOL": "auto"}, "sixel"},
		{"manual blocks", "windows", map[string]string{"NOALAND_IMAGE_PROTOCOL": "blocks"}, "blocks"},
		{"manual kitty", "windows", map[string]string{"NOALAND_IMAGE_PROTOCOL": "kitty"}, "kitty"},
		{"disabled", "windows", map[string]string{"NOALAND_IMAGE_PROTOCOL": "none"}, "none"},
		{"vscode", "windows", map[string]string{"TERM_PROGRAM": "vscode"}, "blocks"},
		{"dumb terminal", "windows", map[string]string{"TERM": "dumb"}, "blocks"},
		{"tmux", "windows", map[string]string{"TMUX": "session"}, "blocks"},
		{"zellij", "windows", map[string]string{"ZELLIJ": "session"}, "blocks"},
		{"remote", "windows", map[string]string{"SSH_CONNECTION": "connection"}, "blocks"},
		{"macOS unknown", "darwin", nil, "blocks"},
		{"Linux unknown", "linux", nil, "blocks"},
		{"Ghostty", "darwin", map[string]string{"TERM_PROGRAM": "ghostty"}, "kitty"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := selectProtocolForOS(func(k string) string { return tt.env[k] }, tt.goos)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
