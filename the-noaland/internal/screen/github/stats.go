package github

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

func statLine(
	label string,
	value int,
	labelStyle lipgloss.Style,
	valueStyle lipgloss.Style,
) string {
	labelText := labelStyle.Render(
		fmt.Sprintf("%-10s", label),
	)

	valueText := valueStyle.Render(
		fmt.Sprintf("%5d", value),
	)

	return labelText + valueText
}
