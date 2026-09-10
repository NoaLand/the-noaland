// Package flipclock renders fixed-size split-flap digit cards without timers.
package flipclock

import (
	"math"
	"strings"

	"charm.land/lipgloss/v2"
)

type Size int

const (
	Compact Size = iota
	Regular
	Large
)

var digits = [10][6]string{
	{" ### ", "#   #", "#   #", "#   #", "#   #", " ### "},
	{"  #  ", " ##  ", "  #  ", "  #  ", "  #  ", " ### "},
	{" ### ", "#   #", "    #", " ### ", "#    ", "#####"},
	{"#### ", "    #", " ### ", "    #", "#   #", " ### "},
	{"#  # ", "#  # ", "#####", "   # ", "   # ", "   # "},
	{"#####", "#    ", "#### ", "    #", "#   #", " ### "},
	{" ### ", "#    ", "#### ", "#   #", "#   #", " ### "},
	{"#####", "    #", "   # ", "  #  ", " #   ", " #   "},
	{" ### ", "#   #", " ### ", "#   #", "#   #", " ### "},
	{" ### ", "#   #", "#   #", " ####", "    #", " ### "},
}
var smallDigits = [10][4]string{
	{"###", "# #", "# #", "###"},
	{" # ", "## ", " # ", "###"},
	{"## ", "  #", " # ", "###"},
	{"## ", " ##", "  #", "## "},
	{"# #", "###", "  #", "  #"},
	{"###", "## ", "  #", "## "},
	{"#  ", "###", "# #", "###"},
	{"###", "  #", " # ", " # "},
	{"###", "###", "# #", "###"},
	{"###", "# #", "###", "  #"},
}

// Pair renders two digit cards. Progress is 0..1; unchanged digits stay still.
func Pair(previous, current string, progress float64, size Size) string {
	if len(previous) != 2 || len(current) != 2 {
		return ""
	}
	return lipgloss.JoinHorizontal(lipgloss.Top,
		digit(previous[0], current[0], progress, size), " ",
		digit(previous[1], current[1], progress, size))
}

// Colon has the same height as a card so horizontal joins remain stable.
func Colon(size Size) string {
	height := 9
	if size == Compact {
		height = 7
	}
	if size == Large {
		height = 15
	}
	lines := make([]string, height)
	for i := range lines {
		lines[i] = "   "
	}
	lines[height/3] = " • "
	lines[height*2/3] = " • "
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#7DCFFF")).Render(strings.Join(lines, "\n"))
}

func glyph(value byte, size Size) []string {
	if value < '0' || value > '9' {
		value = '0'
	}
	var lines []string
	if size == Compact {
		lines = append(lines, smallDigits[value-'0'][:]...)
	} else {
		lines = append(lines, digits[value-'0'][:]...)
	}
	if size != Large {
		return lines
	}
	expanded := make([]string, 0, len(lines)*2)
	for _, line := range lines {
		var out strings.Builder
		for _, r := range line {
			out.WriteRune(r)
			out.WriteRune(r)
		}
		expanded = append(expanded, out.String(), out.String())
	}
	return expanded
}

func digit(previous, current byte, progress float64, size Size) string {
	old, next := glyph(previous, size), glyph(current, size)
	if math.IsNaN(progress) {
		progress = 1
	}
	progress = max(0, min(1, progress))
	rows := append([]string(nil), next...)
	half := len(rows) / 2
	if previous != current && progress < 1 {
		if progress < 0.5 {
			copy(rows[half:], old[half:])
			visible := int(math.Ceil(float64(half) * (1 - progress*2)))
			for i := half - visible; i < half; i++ {
				source := (i - (half - visible)) * half / max(visible, 1)
				rows[i] = old[source]
			}
		} else {
			copy(rows[half:], old[half:])
			visible := int(math.Ceil(float64(half) * (progress*2 - 1)))
			for i := 0; i < visible; i++ {
				rows[half+i] = next[half+i*half/max(visible, 1)]
			}
		}
	}
	width := len(rows[0]) + 2
	border := lipgloss.NewStyle().Foreground(lipgloss.Color("#414868"))
	face := lipgloss.NewStyle().Foreground(lipgloss.Color("#C0CAF5")).Background(lipgloss.Color("#24283B"))
	top := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E7FF")).Background(lipgloss.Color("#292E42"))
	out := []string{border.Render("╭" + strings.Repeat("─", width) + "╮")}
	for i, row := range rows {
		if i == half {
			out = append(out, border.Render("├"+strings.Repeat("─", width)+"┤"))
		}
		row = strings.ReplaceAll(row, "#", "█")
		style := face
		if i < half {
			style = top
		}
		out = append(out, border.Render("│")+style.Render(" "+row+" ")+border.Render("│"))
	}
	out = append(out, border.Render("╰"+strings.Repeat("─", width)+"╯"))
	return strings.Join(out, "\n")
}
