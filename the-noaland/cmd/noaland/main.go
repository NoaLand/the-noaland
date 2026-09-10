package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/app"
)

func main() {
	p := tea.NewProgram(
		app.New(),
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
