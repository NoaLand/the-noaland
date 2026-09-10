package github

import (
	"image"

	githubservice "github.com/NoaLand/the-noaland/the-noaland/internal/service/github"
)

type profileLoadedMsg struct {
	profile githubservice.Profile
	avatar  image.Image
}
