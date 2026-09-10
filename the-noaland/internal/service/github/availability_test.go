package github

import (
	"context"
	"errors"
	"os/exec"
	"reflect"
	"testing"
)

type commandExit int

func (e commandExit) Error() string { return "command failed" }
func (e commandExit) ExitCode() int { return int(e) }

func TestCheckCLI(t *testing.T) {
	for _, tt := range []struct {
		name                    string
		lookupErr, runErr, want error
	}{
		{"missing", exec.ErrNotFound, nil, ErrCLINotInstalled},
		{"not logged in", nil, commandExit(1), ErrNotAuthenticated},
		{"authentication required", nil, commandExit(4), ErrNotAuthenticated},
		{"authenticated", nil, nil, nil},
		{"lookup failure", errors.New("denied"), nil, nil},
		{"timeout", nil, context.DeadlineExceeded, context.DeadlineExceeded},
		{"cancelled", nil, commandExit(2), nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			path, err := checkCLI(func(name string) (string, error) {
				if name != "gh" {
					t.Fatal("wrong executable")
				}
				return "/tools/gh", tt.lookupErr
			}, func(path string, args ...string) error {
				called = true
				if path != "/tools/gh" || !reflect.DeepEqual(args, []string{"auth", "token", "--hostname", "github.example"}) {
					t.Fatal("wrong credential check")
				}
				return tt.runErr
			}, "github.example")
			if tt.lookupErr != nil && called {
				t.Fatal("ran missing executable")
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Fatalf("got %v, want %v", err, tt.want)
			}
			if tt.name == "authenticated" && (err != nil || path != "/tools/gh") {
				t.Fatal("authenticated CLI rejected")
			}
			if tt.name == "lookup failure" || tt.name == "cancelled" {
				if err == nil || errors.Is(err, ErrNotAuthenticated) || errors.Is(err, ErrCLINotInstalled) {
					t.Fatal("non-authentication failure misclassified")
				}
			}
		})
	}
}

func TestGitHubHost(t *testing.T) {
	t.Setenv("GH_HOST", "")
	if githubHost() != "github.com" {
		t.Fatal("wrong default host")
	}
	t.Setenv("GH_HOST", "github.example")
	if githubHost() != "github.example" {
		t.Fatal("host override ignored")
	}
}
