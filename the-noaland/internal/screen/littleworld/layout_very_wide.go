package littleworld

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

func (m Screen) RenderVeryWide(ctx screen.LayoutContext) tea.View {
	m.context = ctx
	return m.view(m.renderVeryWide)
}

func (m Screen) renderVeryWide() string {
	if m.context.Width <= 0 || m.context.Height <= 0 {
		return ""
	}
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#9ECE6A")).Bold(true)
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#565F89"))
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("#A9B1D6"))
	value := lipgloss.NewStyle().Foreground(lipgloss.Color("#C0CAF5"))

	width := min(m.context.Width-4, 176)
	bodyHeight := m.context.Height - 6
	if width < 70 || bodyHeight < 8 {
		return lipgloss.NewStyle().MaxWidth(m.context.Width).MaxHeight(m.context.Height).
			Render(fmt.Sprintf("Little World · Time #%d", m.ticks))
	}

	const sidebarWidth = 26
	annalWidth := width - sidebarWidth - 3
	records := m.world.Annal().Records()
	// Keep the most recent records visible, preserving their chronological order.
	start := max(0, len(records)-(bodyHeight-3))
	count := fmt.Sprintf("%d records", len(records))
	if start > 0 {
		count = fmt.Sprintf("latest %d / %d records", len(records)-start, len(records))
	}

	const timeWidth, entityWidth = 10, 24
	messageWidth := annalWidth - timeWidth - entityWidth
	cell := func(text string, width int) string {
		return lipgloss.NewStyle().Width(width).Render(ansi.Truncate(text, width-1, "…"))
	}
	rows := []string{
		accent.Render("ANNAL") + muted.Render("  ·  "+count),
		label.Render(cell("TIME", timeWidth) + cell("ENTITY", entityWidth) + "EVENT"),
		muted.Render(strings.Repeat("─", annalWidth)),
	}
	for _, record := range records[start:] {
		rows = append(rows,
			muted.Render(cell(fmt.Sprintf("#%d", record.Time), timeWidth))+
				label.Render(cell(string(record.EntityName), entityWidth))+
				value.Render(ansi.Truncate(record.Message, messageWidth, "…")),
		)
	}
	if len(records) == 0 {
		rows = append(rows, muted.Render("Waiting for the first change…"))
	}
	annal := lipgloss.NewStyle().Width(annalWidth).Height(bodyHeight).
		Render(strings.Join(rows, "\n"))
	sidebar := lipgloss.NewStyle().Width(sidebarWidth).Height(bodyHeight).
		Render(strings.Join([]string{
			accent.Render("WORLD"),
			"",
			label.Render("Seed"),
			value.Render(fmt.Sprint(m.world.Seed())),
			"",
			label.Render("Time"),
			value.Render(fmt.Sprintf("#%d", m.ticks)),
			"",
			muted.Render("1 step / second"),
		}, "\n"))
	divider := muted.Render(strings.TrimSuffix(strings.Repeat(" │ \n", bodyHeight), "\n"))
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, divider, annal)
	content := lipgloss.JoinVertical(lipgloss.Left,
		accent.Render("LITTLE WORLD"), "", body, "",
		muted.Render("← / →  switch screen    q  quit"),
	)
	return lipgloss.Place(m.context.Width, m.context.Height, lipgloss.Center, lipgloss.Center, content)
}
