package github

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/ui/avatar"
)

func (m Screen) RenderWide(ctx screen.LayoutContext) tea.View {
	m.context = ctx
	return m.view(func() string { return m.renderFullWide(m.renderRecentHeatmap(12)) })
}

func (m Screen) renderFullWide(right string) string {
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7DCFFF"))

	handleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9B1D6"))

	statLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89"))

	statValueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5"))

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#292E42"))

	updatedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B4261"))

	leftWidth := m.context.Width / 3
	rightWidth := m.context.Width - leftWidth

	today := m.todayContributions()
	thisWeek := m.thisWeekContributions()
	thisYear := m.profile.TotalContributions

	updated := updatedStyle.Render(
		"Updated " + m.lastUpdated.Format("15:04"),
	)

	stats := lipgloss.JoinVertical(
		lipgloss.Left,
		statLine(
			"Today",
			today,
			statLabelStyle,
			statValueStyle,
		),
		statLine(
			"This week",
			thisWeek,
			statLabelStyle,
			statValueStyle,
		),
		statLine(
			"This year",
			thisYear,
			statLabelStyle,
			statValueStyle,
		),
	)

	avatarCols := avatarWidth(m.context.Width)
	avatarRows := avatarHeight(m.context.Height)

	profileHeader := lipgloss.JoinVertical(
		lipgloss.Center,
		avatar.Placeholder(
			avatarCols,
			avatarRows,
		),
		"",
		nameStyle.Render(m.profile.Name),
		handleStyle.Render("@"+m.profile.Login),
	)

	left := lipgloss.JoinVertical(
		lipgloss.Center,
		profileHeader,
		"",
		separatorStyle.Render("──────────────"),
		"",
		stats,
		"",
		updated,
	)

	leftPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.context.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(left)

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.context.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(right)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)
}

func (m Screen) fullAvatarPosition() (int, int) {
	leftWidth := m.context.Width / 3

	avatarCols := avatarWidth(m.context.Width)
	avatarRows := avatarHeight(m.context.Height)

	profileHeight :=
		avatarRows +
			1 + // blank after avatar
			1 + // name
			1 + // handle
			1 + // blank
			1 + // separator
			1 + // blank
			3 + // stats
			1 + // blank
			1 // updated

	row := (m.context.Height-profileHeight)/2 + 1
	col := (leftWidth-avatarCols)/2 + 1

	return row, col
}
