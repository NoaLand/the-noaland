package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type Response struct {
	Data struct {
		Viewer struct {
			Login     string `json:"login"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"viewer"`
	} `json:"data"`
}

func main() {
	profile, err := fetchProfile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch profile: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", profile.Data.Viewer.Name)
	fmt.Printf("@%s\n\n", profile.Data.Viewer.Login)

	avatarPath := "/tmp/dev-dashboard-avatar.png"

	if err := download(profile.Data.Viewer.AvatarURL, avatarPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to download avatar: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("chafa", "--size", "20x10", avatarPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to render avatar: %v\n", err)
		os.Exit(1)
	}
}

func fetchProfile() (*Response, error) {
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
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func download(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
