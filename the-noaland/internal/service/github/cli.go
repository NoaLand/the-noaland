package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

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

// FetchProfile loads the authenticated user's profile using the GitHub CLI.
func FetchProfile() (Profile, error) {
	host := githubHost()
	path, err := checkCLI(exec.LookPath, runCredentialCheck, host)
	if err != nil {
		return Profile{}, err
	}
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
		path,
		"api",
		"graphql",
		"--hostname", host,
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
