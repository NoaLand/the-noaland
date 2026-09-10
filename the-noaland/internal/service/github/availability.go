package github

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var (
	ErrCLINotInstalled  = errors.New("GitHub CLI is not installed or not on PATH")
	ErrNotAuthenticated = errors.New("GitHub CLI has no authentication credentials")
)

func githubHost() string {
	if host := os.Getenv("GH_HOST"); host != "" {
		return host
	}
	return "github.com"
}

// Check verifies installation and locally available credentials.
// API requests still validate token access and network connectivity.
func Check() error {
	_, err := checkCLI(exec.LookPath, runCredentialCheck, githubHost())
	return err
}

func runCredentialCheck(path string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// auth token prints a secret. Run discards stdout/stderr; never capture or log it.
	err := exec.CommandContext(ctx, path, args...).Run()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

func checkCLI(lookup func(string) (string, error), run func(string, ...string) error, host string) (string, error) {
	path, err := lookup("gh")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", ErrCLINotInstalled
		}
		return "", fmt.Errorf("locate GitHub CLI: %w", err)
	}
	if err := run(path, "auth", "token", "--hostname", host); err != nil {
		var exit interface{ ExitCode() int }
		if errors.As(err, &exit) && (exit.ExitCode() == 1 || exit.ExitCode() == 4) {
			return "", ErrNotAuthenticated
		}
		return "", fmt.Errorf("check GitHub CLI authentication: %w", err)
	}
	return path, nil
}
