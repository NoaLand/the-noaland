package github

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	"github.com/NoaLand/the-noaland/the-noaland/internal/ui/github/avatar"
)

func (m Screen) renderAvatarCmd() tea.Cmd {
	layout := m.context.Layout
	if m.avatar == nil ||
		layout == screen.LayoutTall {
		return nil
	}

	var row, col int

	if layout == screen.LayoutCompactWide {
		row, col = m.compactAvatarPosition()
	} else {
		row, col = m.fullAvatarPosition()
	}

	avatar := avatar.Render(
		m.avatar,
		avatarWidth(m.context.Width),
		avatarHeight(m.context.Height),
	)

	if avatar == "" {
		return nil
	}

	raw :=
		moveCursor(row, col) +
			avatar +
			"\x1b[H"

	return tea.Raw(raw)
}

func avatarWidth(width int) int {
	if width >= 150 {
		return 18
	}

	return 14
}

func avatarHeight(height int) int {
	if height >= 35 {
		return 9
	}

	return 7
}

func moveCursor(
	row int,
	col int,
) string {
	return fmt.Sprintf(
		"\x1b[%d;%dH",
		row,
		col,
	)
}
