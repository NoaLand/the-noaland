package heatmap

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"

	githubservice "github.com/NoaLand/the-noaland/the-noaland/internal/service/github"
)

// Render draws GitHub contribution weeks, optionally including month labels.
func Render(
	title string,
	weeks []githubservice.ContributionWeek,
	showMonths bool,
) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5"))

	weekLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89"))

	rows := make([]string, 7)

	for _, week := range weeks {
		for dayIndex := range 7 {
			if dayIndex >= len(week.ContributionDays) {
				rows[dayIndex] += "  "
				continue
			}

			day := week.ContributionDays[dayIndex]

			if day.Date == "" {
				rows[dayIndex] += "  "
				continue
			}

			rows[dayIndex] += contributionCell(
				day.ContributionCount,
			)
		}
	}

	weekdayLabels := []string{
		"   ",
		"Mon",
		"   ",
		"Wed",
		"   ",
		"Fri",
		"   ",
	}

	var heatmapRows []string

	for i, row := range rows {
		heatmapRows = append(
			heatmapRows,
			fmt.Sprintf(
				"%s  %s",
				weekLabelStyle.Render(
					weekdayLabels[i],
				),
				row,
			),
		)
	}

	matrix := lipgloss.JoinVertical(
		lipgloss.Left,
		heatmapRows...,
	)

	parts := []string{titleStyle.Render(title), ""}
	if showMonths {
		parts = append(parts, "     "+buildMonthLabels(weeks))
	}
	parts = append(parts, matrix, "", "     "+renderLegend())
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func buildMonthLabels(
	weeks []githubservice.ContributionWeek,
) string {
	if len(weeks) == 0 {
		return ""
	}

	type monthMarker struct {
		weekIndex int
		label     string
	}

	var markers []monthMarker

	lastMonth := time.Month(0)

	for weekIndex, week := range weeks {
		for _, day := range week.ContributionDays {
			if day.Date == "" {
				continue
			}

			t, err := time.Parse(
				"2006-01-02",
				day.Date,
			)
			if err != nil {
				continue
			}

			if t.Month() != lastMonth {
				markers = append(
					markers,
					monthMarker{
						weekIndex: weekIndex,
						label:     t.Format("Jan"),
					},
				)

				lastMonth = t.Month()
			}

			break
		}
	}

	totalWidth := len(weeks) * 2

	line := make([]rune, totalWidth)

	for i := range line {
		line[i] = ' '
	}

	for _, marker := range markers {
		pos := marker.weekIndex * 2

		for i, r := range marker.label {
			if pos+i >= len(line) {
				break
			}

			line[pos+i] = r
		}
	}

	return string(line)
}

func renderLegend() string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89"))

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		labelStyle.Render("Less "),
		contributionLegendCell("#292E42"),
		contributionLegendCell("#3B4261"),
		contributionLegendCell("#73DACA"),
		contributionLegendCell("#41A6B5"),
		contributionLegendCell("#7DCFFF"),
		labelStyle.Render(" More"),
	)
}

func contributionLegendCell(
	color string,
) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render("■ ")
}

func contributionCell(count int) string {
	style := lipgloss.NewStyle()

	switch {
	case count == 0:
		style = style.Foreground(
			lipgloss.Color("#292E42"),
		)

	case count <= 2:
		style = style.Foreground(
			lipgloss.Color("#3B4261"),
		)

	case count <= 5:
		style = style.Foreground(
			lipgloss.Color("#73DACA"),
		)

	case count <= 10:
		style = style.Foreground(
			lipgloss.Color("#41A6B5"),
		)

	default:
		style = style.Foreground(
			lipgloss.Color("#7DCFFF"),
		)
	}

	return style.Render("■ ")
}
