package github

import (
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	githubservice "github.com/NoaLand/the-noaland/the-noaland/internal/service/github"
)

func (m Screen) errorView() tea.View {
	title, detail := "GitHub is unavailable", fmt.Sprintf("Error: %v", m.err)
	switch {
	case errors.Is(m.err, githubservice.ErrCLINotInstalled):
		title = "GitHub CLI is not installed"
		detail = "Install gh from https://cli.github.com\nMake sure gh is on PATH, then restart NoaLand or retry."
	case errors.Is(m.err, githubservice.ErrNotAuthenticated):
		title = "GitHub CLI is not logged in"
		detail = "Run gh auth login in another terminal.\nThen return here and press r to retry."
	}
	heading := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7DCFFF")).Render(title)
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#A9B1D6")).Render("→  Flip Clock    r  retry    q  quit")
	content := lipgloss.JoinVertical(lipgloss.Center, heading, "", detail, "", hint)
	if m.context.Width > 0 && m.context.Height > 0 {
		content = lipgloss.NewStyle().MaxWidth(m.context.Width).MaxHeight(m.context.Height).Render(content)
		content = lipgloss.Place(m.context.Width, m.context.Height, lipgloss.Center, lipgloss.Center, content)
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
