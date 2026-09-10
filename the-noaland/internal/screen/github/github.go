package github

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"net/http"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
	githubservice "github.com/NoaLand/the-noaland/the-noaland/internal/service/github"
	"github.com/NoaLand/the-noaland/the-noaland/internal/ui/heatmap"
)

type redrawAvatarMsg struct{}

type refreshMsg struct{}

type errMsg struct {
	err error
}

// Screen displays a GitHub profile and contribution activity.
type Screen struct {
	context screen.LayoutContext

	profile *githubservice.Profile
	avatar  image.Image

	lastUpdated time.Time

	err error
}

// New creates a GitHub screen with no profile loaded.
func New() *Screen {
	return &Screen{}
}

func (m Screen) Init() tea.Cmd {
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

func (m *Screen) Update(msg tea.Msg, ctx screen.LayoutContext) tea.Cmd {
	m.context = ctx
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "r":
			return fetchProfileCmd()
		}

	case refreshMsg:
		return tea.Batch(
			fetchProfileCmd(),
			refreshTick(),
		)

	case tea.WindowSizeMsg:
		if m.avatar != nil {
			return m.renderAvatarCmd()
		}

	case profileLoadedMsg:
		m.profile = &msg.profile
		m.avatar = msg.avatar
		m.lastUpdated = time.Now()

		return redrawAvatarLater()

	case redrawAvatarMsg:
		return m.renderAvatarCmd()

	case errMsg:
		m.err = msg.err
	}

	return nil
}

func (m Screen) view(render func() string) tea.View {
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

	v := tea.NewView(render())
	v.AltScreen = true

	return v
}

func (m Screen) todayContributions() int {
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

func (m Screen) thisWeekContributions() int {
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

func (m Screen) renderRecentHeatmap(
	weekCount int,
) string {
	if len(m.profile.Weeks) == 0 {
		return "No contribution data"
	}

	weeks := m.profile.Weeks

	if len(weeks) > weekCount {
		weeks = weeks[len(weeks)-weekCount:]
	}

	return heatmap.Render(
		fmt.Sprintf(
			"Recent contributions · %d weeks",
			weekCount,
		),
		weeks,
		false,
	)
}

func (m Screen) renderYearHeatmap() string {
	if len(m.profile.Weeks) == 0 {
		return "No contribution data"
	}

	return heatmap.Render(
		"Contribution activity",
		m.profile.Weeks,
		true,
	)
}

func fetchProfileCmd() tea.Cmd {
	return func() tea.Msg {
		profile, err := githubservice.FetchProfile()
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

func redrawAvatarLater() tea.Cmd {
	return tea.Tick(
		50*time.Millisecond,
		func(time.Time) tea.Msg {
			return redrawAvatarMsg{}
		},
	)
}
