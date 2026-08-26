package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Profile struct {
	Login     string
	Name      string
	AvatarURL string
}

type githubResponse struct {
	Data struct {
		Viewer struct {
			Login     string `json:"login"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatarUrl"`
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

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7DCFFF"))

	handleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A9B1D6"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#565F89"))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.avatar,
		"",
		titleStyle.Render(m.profile.Name),
		handleStyle.Render("@"+m.profile.Login),
		"",
		hintStyle.Render("q to quit"),
	)

	container := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)

	v := tea.NewView(container)
	v.AltScreen = true

	return v
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

func renderAvatarCmd(path string, width, height int) tea.Cmd {
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
				fmt.Errorf("render avatar: %w", err),
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
		return Profile{}, fmt.Errorf("gh api graphql: %w", err)
	}

	var response githubResponse

	if err := json.Unmarshal(output, &response); err != nil {
		return Profile{}, fmt.Errorf("decode response: %w", err)
	}

	return Profile{
		Login:     response.Data.Viewer.Login,
		Name:      response.Data.Viewer.Name,
		AvatarURL: response.Data.Viewer.AvatarURL,
	}, nil
}

func downloadAvatar(url string) (string, error) {
	cacheDir := filepath.Join(os.TempDir(), "dev-dashboard")

	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}

	path := filepath.Join(cacheDir, "avatar.png")

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"avatar download returned %s",
			resp.Status,
		)
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}

	return path, nil
}

func avatarWidth(width int) int {
	switch {
	case width >= 100:
		return 28
	case width >= 70:
		return 22
	default:
		return 16
	}
}

func avatarHeight(height int) int {
	switch {
	case height >= 40:
		return 14
	case height >= 28:
		return 10
	default:
		return 8
	}
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "dev-dashboard: %v\n", err)
		os.Exit(1)
	}
}
