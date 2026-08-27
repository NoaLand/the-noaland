package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type redrawAvatarMsg struct{}

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
	profile Profile
	avatar  image.Image
}

type refreshMsg struct{}

type errMsg struct {
	err error
}

type model struct {
	width  int
	height int

	profile *Profile
	avatar  image.Image

	lastUpdated time.Time

	err error
}

const kittyChunkSize = 4096

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchProfileCmd(),
		refreshTick(),
	)
}

func refreshTick() tea.Cmd {
	return tea.Tick(
		5*time.Minute,
		func(time.Time) tea.Msg {
			return refreshMsg{}
		},
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit

		case "r":
			return m, fetchProfileCmd()
		}

	case refreshMsg:
		return m, tea.Batch(
			fetchProfileCmd(),
			refreshTick(),
		)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.avatar != nil {
			return m, m.renderAvatarCmd()
		}

	case profileLoadedMsg:
		m.profile = &msg.profile
		m.avatar = msg.avatar
		m.lastUpdated = time.Now()

		return m, redrawAvatarLater()

	case redrawAvatarMsg:
		return m, m.renderAvatarCmd()

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.err != nil {
		return tea.NewView(
			fmt.Sprintf(
				"Error: %v\n\nPress q to quit.",
				m.err,
			),
		)
	}

	if m.profile == nil {
		return tea.NewView(
			"Loading GitHub profile...",
		)
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

func (m model) renderAvatarCmd() tea.Cmd {
	if m.avatar == nil ||
		!isWide(m.width, m.height) {
		return nil
	}

	var row, col int

	if isCompactWide(m.width, m.height) {
		row, col = m.compactAvatarPosition()
	} else {
		row, col = m.fullAvatarPosition()
	}

	avatar := renderKittyImage(
		m.avatar,
		avatarWidth(m.width),
		avatarHeight(m.height),
	)

	raw :=
		moveCursor(row, col) +
			avatar +
			"\x1b[H"

	return tea.Raw(raw)
}

func isWide(width, height int) bool {
	return width >= 90 && width > height*2
}

func isVeryWide(width int) bool {
	return width >= 150
}

func isCompactWide(width, height int) bool {
	return width >= 90 && height < 24
}

func (m model) renderWide() string {
	if isCompactWide(m.width, m.height) {
		return m.renderCompactWide()
	}

	return m.renderFullWide()
}

func (m model) renderCompactWide() string {
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

	leftWidth := m.width * 2 / 5
	rightWidth := m.width - leftWidth

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

	avatarCols := avatarWidth(m.width)
	avatarRows := avatarHeight(m.height)

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
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(profile)

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(heatmap)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)
}

func (m model) renderFullWide() string {
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

	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth

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

	avatarCols := avatarWidth(m.width)
	avatarRows := avatarHeight(m.height)

	profileHeader := lipgloss.JoinVertical(
		lipgloss.Center,
		avatarPlaceholder(
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

func (m model) renderRecentHeatmap(
	weekCount int,
) string {
	if len(m.profile.Weeks) == 0 {
		return "No contribution data"
	}

	weeks := m.profile.Weeks

	if len(weeks) > weekCount {
		weeks = weeks[len(weeks)-weekCount:]
	}

	return renderCompactHeatmap(
		fmt.Sprintf(
			"Recent contributions · %d weeks",
			weekCount,
		),
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

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		"",
		"     "+monthLabels,
		matrix,
		"",
		"     "+renderLegend(),
	)
}

func buildMonthLabels(
	weeks []ContributionWeek,
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

func fetchProfileCmd() tea.Cmd {
	return func() tea.Msg {
		profile, err := fetchProfile()
		if err != nil {
			return errMsg{err}
		}

		avatar, err := downloadAvatar(
			profile.AvatarURL,
		)
		if err != nil {
			return errMsg{err}
		}

		return profileLoadedMsg{
			profile: profile,
			avatar:  avatar,
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
	  totalCommitContributions
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
) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf(
				"avatar download returned %s",
				resp.Status,
			)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
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

func avatarPlaceholder(
	width int,
	height int,
) string {
	line := strings.Repeat(" ", width)

	lines := make([]string, height)

	for i := range lines {
		lines[i] = line
	}

	return strings.Join(lines, "\n")
}

func redrawAvatarLater() tea.Cmd {
	return tea.Tick(
		50*time.Millisecond,
		func(time.Time) tea.Msg {
			return redrawAvatarMsg{}
		},
	)
}

func renderKittyImage(
	img image.Image,
	cols int,
	rows int,
) string {
	if img == nil {
		return ""
	}

	var pngBuf bytes.Buffer

	if err := png.Encode(&pngBuf, img); err != nil {
		return ""
	}

	data := base64.StdEncoding.EncodeToString(
		pngBuf.Bytes(),
	)

	var out strings.Builder

	first := true

	for len(data) > 0 {
		n := kittyChunkSize
		if len(data) < n {
			n = len(data)
		}

		chunk := data[:n]
		data = data[n:]

		more := len(data) > 0

		m := 0
		if more {
			m = 1
		}

		if first {
			fmt.Fprintf(
				&out,
				"\x1b_Ga=T,f=100,c=%d,r=%d,q=2,m=%d;%s\x1b\\",
				cols,
				rows,
				m,
				chunk,
			)

			first = false
		} else {
			fmt.Fprintf(
				&out,
				"\x1b_Gm=%d;%s\x1b\\",
				m,
				chunk,
			)
		}
	}

	return out.String()
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

func (m model) compactAvatarPosition() (int, int) {
	leftWidth := m.width * 2 / 5

	avatarCols := avatarWidth(m.width)
	avatarRows := avatarHeight(m.height)

	profileHeight :=
		avatarRows +
			1 + // name
			1 + // meta
			1 + // blank
			3 // stats

	row := (m.height-profileHeight)/2 + 1
	col := (leftWidth-avatarCols)/2 + 1

	return row, col
}

func (m model) fullAvatarPosition() (int, int) {
	leftWidth := m.width / 3

	avatarCols := avatarWidth(m.width)
	avatarRows := avatarHeight(m.height)

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

	row := (m.height-profileHeight)/2 + 1
	col := (leftWidth-avatarCols)/2 + 1

	return row, col
}

