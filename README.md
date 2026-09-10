# dev-dashboard

A lightweight terminal dashboard for my development environment.

The initial version focuses on displaying GitHub profile and contribution activity inside a TUI, designed to stay open as part of a terminal workspace.

## Goals

- Show GitHub profile information
- Visualize contribution activity
- Adapt to different terminal sizes
- Work well inside Ghostty and Zellij
- Stay lightweight and easy to extend

## Planned

- GitHub avatar
- Daily / weekly / yearly contribution stats
- Contribution heatmap
- Responsive compact and expanded layouts
- Additional developer status panels in the future

## Tech

- Go
- Bubble Tea
- Lip Gloss
- GitHub CLI / GraphQL API

## Flip Clock

Run the dashboard with:

    go run ./the-noaland/cmd/noaland

Use left/right arrow keys to cycle between GitHub and Flip Clock. Press q, Esc,
or Ctrl+C to exit. GitHub remains the initial screen.

Flip Clock shows 24-hour local time (hours, minutes, seconds), date, and timezone.
Changed digits animate for 360 ms using split-flap text cards. The clock requires
no image protocol or external service. Small windows fall back to plain time.

Wide, compact-wide, very-wide, and tall layouts have separate render methods.
The flip-card UI component accepts digits, animation progress, and size; it does
not know the window layout or own timers. Clock timers run only while its screen
is active. Returning to the clock synchronizes immediately to current local time.

Tests cover digit geometry, midnight rollover, delayed ticks, stale timer
messages after switching, all four layouts, and navigation between both screens.
Animation smoothness and native image cleanup on page switches still require
visual verification in the target terminal.