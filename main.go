package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
	Color             string `json:"color"`
}

type ContributionWeek struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

type Profile struct {
	Login              string
	Name               string
	AvatarURL          string
	TotalContributions int
	Weeks              []ContributionWeek
}

type githubResponse struct {
	Data struct {
		Viewer struct {
			Login                   string `json:"login"`
			Name                    string `json:"name"`
			AvatarURL               string `json:"avatarUrl"`
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int                `json:"totalContributions"`
					Weeks              []ContributionWeek `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"viewer"`
	} `json:"data"`
}

type profileLoadedMsg struct {
	profile    Profile
	avatarPath string
}

type avatarRenderedMsg struct {
	content string
}

type errMsg struct {
	err error
}

type model struct {
	width  int
	height int

	profile    *Profile
	avatarPath string
	avatar     string

	err error
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return fetchProfileCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.avatarPath != "" {
			return m, renderAvatarCmd(
				m.avatarPath,
				avatarWidth(m.width),
				avatarHeight(m.height),
			)
		}

	case profileLoadedMsg:
		m.profile = &msg.profile
		m.avatarPath = msg.avatarPath

		return m, renderAvatarCmd(
			m.avatarPath,
			avatarWidth(m.width),
			avatarHeight(m.height),
		)

	case avatarRenderedMsg:
		m.avatar = msg.content

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.err != nil {
		return tea.NewView(
			fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err),
		)
	}

	if m.profile == nil {
		return tea.NewView("Loading GitHub profile...")
	}

	var content string

	if isWide(m.width, m.height) {
		content = m.renderWide()
	} else {
		content = m.renderTall()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func isWide(width, height int) bool {
	return width >= 90 && width > height*2
}

func isVeryWide(width int) bool {
	return width >= 150
}

func (m model) renderWide() string {
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

	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth

	today := m.todayContributions()
	thisWeek := m.thisWeekContributions()
	thisYear := m.profile.TotalContributions

	stats := lipgloss.JoinVertical(
		lipgloss.Left,
		statLine("Today", today, statLabelStyle, statValueStyle),
		statLine("This week", thisWeek, statLabelStyle, statValueStyle),
		statLine("This year", thisYear, statLabelStyle, statValueStyle),
	)

	profileHeader := lipgloss.JoinVertical(
		lipgloss.Center,
		m.avatar,
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
	)

	var right string

	if isVeryWide(m.width) {
		right = m.renderYearHeatmap()
	} else {
		right = m.renderRecentHeatmap(12)
	}

	leftPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(left)

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(right)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)
}

func (m model) renderTall() string {
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render("Tall layout coming soon...")
}

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

func (m model) todayContributions() int {
	today := time.Now().Format("2006-01-02")

	for _, week := range m.profile.Weeks {
		for _, day := range week.ContributionDays {
			if day.Date == today {
				return day.ContributionCount
			}
		}
	}

	return 0
}

func (m model) thisWeekContributions() int {
	if len(m.profile.Weeks) == 0 {
		return 0
	}

	total := 0
	currentWeek := m.profile.Weeks[len(m.profile.Weeks)-1]

	for _, day := range currentWeek.ContributionDays {
		total += day.ContributionCount
	}

	return total
}

func (m model) renderRecentHeatmap(weekCount int) string {
	if len(m.profile.Weeks) == 0 {
		return "No contribution data"
	}

	weeks := m.profile.Weeks

	if len(weeks) > weekCount {
		weeks = weeks[len(weeks)-weekCount:]
	}

	return renderCompactHeatmap(
		fmt.Sprintf("Recent contributions · %d weeks", weekCount),
		weeks,
	)
}

func renderCompactHeatmap(
	title string,
	weeks []ContributionWeek,
) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5"))

	weekLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89"))

	rows := make([]string, 7)

	for _, week := range weeks {
		for dayIndex := 0; dayIndex < 7; dayIndex++ {
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
				weekLabelStyle.Render(weekdayLabels[i]),
				row,
			),
		)
	}

	matrix := lipgloss.JoinVertical(
		lipgloss.Left,
		heatmapRows...,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		"",
		matrix,
		"",
		"     "+renderLegend(),
	)
}

func (m model) renderYearHeatmap() string {
	if len(m.profile.Weeks) == 0 {
		return "No contribution data"
	}

	return renderYearHeatmap(
		"Contribution activity",
		m.profile.Weeks,
	)
}

func renderYearHeatmap(
	title string,
	weeks []ContributionWeek,
) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5"))

	weekLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89"))

	monthLabels := buildMonthLabels(weeks)

	rows := make([]string, 7)

	for weekIndex, week := range weeks {
		for dayIndex := 0; dayIndex < 7; dayIndex++ {
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

		_ = weekIndex
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
				weekLabelStyle.Render(weekdayLabels[i]),
				row,
			),
		)
	}

	matrix := lipgloss.JoinVertical(
		lipgloss.Left,
		heatmapRows...,
	)

	legend := renderLegend()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		"",
		"     "+monthLabels,
		matrix,
		"",
		"     "+legend,
	)
}

func buildMonthLabels(weeks []ContributionWeek) string {
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

			t, err := time.Parse("2006-01-02", day.Date)
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

func contributionLegendCell(color string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render("■ ")
}

func renderHeatmap(
	title string,
	weeks []ContributionWeek,
) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#C0CAF5"))

	rows := make([]string, 7)

	for _, week := range weeks {
		for dayIndex := 0; dayIndex < 7; dayIndex++ {
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

	matrix := lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		"",
		matrix,
	)
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

func fetchProfileCmd() tea.Cmd {
	return func() tea.Msg {
		profile, err := fetchProfile()
		if err != nil {
			return errMsg{err}
		}

		path, err := downloadAvatar(profile.AvatarURL)
		if err != nil {
			return errMsg{err}
		}

		return profileLoadedMsg{
			profile:    profile,
			avatarPath: path,
		}
	}
}

func renderAvatarCmd(
	path string,
	width,
	height int,
) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(
			"chafa",
			"--format=symbols",
			"--size",
			fmt.Sprintf("%dx%d", width, height),
			path,
		)

		output, err := cmd.Output()
		if err != nil {
			return errMsg{
				err: fmt.Errorf(
					"render avatar: %w",
					err,
				),
			}
		}

		return avatarRenderedMsg{
			content: string(output),
		}
	}
}

func fetchProfile() (Profile, error) {
	query := `
query {
  viewer {
    login
    name
    avatarUrl
    contributionsCollection {
      contributionCalendar {
        totalContributions
        weeks {
          contributionDays {
            date
            contributionCount
            color
          }
        }
      }
    }
  }
}`

	cmd := exec.Command(
		"gh",
		"api",
		"graphql",
		"-f",
		"query="+query,
	)

	output, err := cmd.Output()
	if err != nil {
		return Profile{},
			fmt.Errorf(
				"gh api graphql: %w",
				err,
			)
	}

	var response githubResponse

	if err := json.Unmarshal(
		output,
		&response,
	); err != nil {
		return Profile{},
			fmt.Errorf(
				"decode response: %w",
				err,
			)
	}

	viewer := response.Data.Viewer
	calendar :=
		viewer.
			ContributionsCollection.
			ContributionCalendar

	return Profile{
		Login:              viewer.Login,
		Name:               viewer.Name,
		AvatarURL:          viewer.AvatarURL,
		TotalContributions: calendar.TotalContributions,
		Weeks:              calendar.Weeks,
	}, nil
}

func downloadAvatar(
	url string,
) (string, error) {
	cacheDir := filepath.Join(
		os.TempDir(),
		"dev-dashboard",
	)

	if err := os.MkdirAll(
		cacheDir,
		0o755,
	); err != nil {
		return "", err
	}

	path := filepath.Join(
		cacheDir,
		"avatar.png",
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "",
			fmt.Errorf(
				"avatar download returned %s",
				resp.Status,
			)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(
		file,
		resp.Body,
	); err != nil {
		return "", err
	}

	return path, nil
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

func main() {
	p := tea.NewProgram(
		initialModel(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"dev-dashboard: %v\n",
			err,
		)

		os.Exit(1)
	}
}
