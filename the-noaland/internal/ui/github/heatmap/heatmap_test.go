package heatmap

import (
	"strings"
	"testing"

	githubservice "github.com/NoaLand/the-noaland/the-noaland/internal/service/github"
)

func TestMonthLabelsWithMissingDatesAndClipping(t *testing.T) {
	weeks := []githubservice.ContributionWeek{
		{ContributionDays: []githubservice.ContributionDay{{Date: ""}, {Date: "invalid"}, {Date: "2026-01-01"}}},
		{},
		{ContributionDays: []githubservice.ContributionDay{{Date: "2026-02-01"}}},
	}
	if got := buildMonthLabels(weeks); got != "Jan Fe" {
		t.Fatalf("labels = %q, want %q", got, "Jan Fe")
	}
	if got := buildMonthLabels(nil); got != "" {
		t.Fatalf("empty labels = %q", got)
	}
}

func TestRenderOptionalMonthsAndSparseWeeks(t *testing.T) {
	weeks := []githubservice.ContributionWeek{
		{ContributionDays: []githubservice.ContributionDay{{Date: "2026-01-01", ContributionCount: 3}}},
		{},
	}
	without := Render("Contributions", weeks, false)
	with := Render("Contributions", weeks, true)
	if strings.Count(with, string(byte(10))) != strings.Count(without, string(byte(10)))+1 {
		t.Fatal("month labels must add exactly one row")
	}
	if !strings.Contains(with, "Jan") || strings.Contains(without, "Jan") {
		t.Fatal("month label visibility changed")
	}
	for _, output := range []string{without, with, Render("Empty", nil, false)} {
		if !strings.Contains(output, "Mon") || !strings.Contains(output, "Less") || !strings.Contains(output, "More") {
			t.Fatal("weekday labels or legend missing")
		}
	}
}
