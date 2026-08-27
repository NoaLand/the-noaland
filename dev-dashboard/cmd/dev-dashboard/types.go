package main

import "image"

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

type profileLoadedMsg struct {
	profile Profile
	avatar  image.Image
}
