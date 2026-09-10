package github

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

func (m Screen) RenderCompactWide(ctx screen.LayoutContext) tea.View {
	m.context = ctx
	return m.view(m.renderCompactWide)
}

func (m Screen) renderCompactWide() string {
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

	updatedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3B4261"))

	updated := updatedStyle.Render(
		"Updated " + m.lastUpdated.Format("15:04"),
	)

	leftWidth := m.context.Width * 2 / 5
	rightWidth := m.context.Width - leftWidth

	today := m.todayContributions()
	thisWeek := m.thisWeekContributions()
	thisYear := m.profile.TotalContributions

	stats := lipgloss.JoinVertical(
		lipgloss.Left,
		statLine(
			"Today",
			today,
			statLabelStyle,
			statValueStyle,
		),
		statLine(
			"Week",
			thisWeek,
			statLabelStyle,
			statValueStyle,
		),
		statLine(
			"Year",
			thisYear,
			statLabelStyle,
			statValueStyle,
		),
	)

	meta := lipgloss.JoinHorizontal(
		lipgloss.Center,
		handleStyle.Render("@"+m.profile.Login),
		updatedStyle.Render(" · "),
		updated,
	)

	avatarCols := avatarWidth(m.context.Width)
	avatarRows := avatarHeight(m.context.Height)

	profile := lipgloss.JoinVertical(
		lipgloss.Center,
		avatarPlaceholder(
			avatarCols,
			avatarRows,
		),
		nameStyle.Render(m.profile.Name),
		meta,
		"",
		stats,
	)

	heatmap := m.renderRecentHeatmap(12)

	leftPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.context.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(profile)

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.context.Height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(heatmap)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)
}

func (m Screen) compactAvatarPosition() (int, int) {
	leftWidth := m.context.Width * 2 / 5

	avatarCols := avatarWidth(m.context.Width)
	avatarRows := avatarHeight(m.context.Height)

	profileHeight :=
		avatarRows +
			1 + // name
			1 + // meta
			1 + // blank
			3 // stats

	row := (m.context.Height-profileHeight)/2 + 1
	col := (leftWidth-avatarCols)/2 + 1

	return row, col
}
