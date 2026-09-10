# The NoaLand

[![NoaLand CI](https://github.com/NoaLand/the-noaland/actions/workflows/noaland-ci.yml/badge.svg)](https://github.com/NoaLand/the-noaland/actions/workflows/noaland-ci.yml)

[![NoaLand Release](https://github.com/NoaLand/the-noaland/actions/workflows/noaland-release.yml/badge.svg)](https://github.com/NoaLand/the-noaland/actions/workflows/noaland-release.yml)

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**The NoaLand** is a personal terminal space built with Go and Bubble Tea.

It started as a lightweight developer dashboard, but is gradually evolving into something broader: a collection of interactive terminal screens for information, tools, ambient visuals, and small experiments.

The goal is to keep it lightweight, responsive, and enjoyable to leave running as part of a terminal workspace.

---

## Screens

### GitHub

The first screen in The NoaLand.

It displays personal GitHub activity inside the terminal, including:

- GitHub profile information
- Daily, weekly, and yearly contribution statistics
- Contribution heatmap
- Responsive layouts for different terminal sizes
- GitHub avatar rendering on supported terminals
- Periodic refresh of GitHub activity

GitHub data is queried through the GitHub CLI and GraphQL API.

---

### Flip Clock

A terminal split-flap clock built as the second screen in The NoaLand.

The clock renders fixed-size digit cards and animates transitions between digits using Bubble Tea's update loop.

Example frame:

```text
╭───────╮ ╭───────╮     ╭───────╮ ╭───────╮
│  ███  │ │   █   │     │ █████ │ │  ███  │
│ █   █ │ │  ██   │     │     █ │ │ █   █ │
│ █   █ │ │   █   │  •  │   ██  │ │    ██ │
├───────┤ ├───────┤     ├───────┤ ├───────┤
│ █   █ │ │   █   │  •  │  █    │ │     █ │
│ █   █ │ │   █   │     │ █     │ │ █   █ │
│  ███  │ │  ███  │     │ █████ │ │  ███  │
╰───────╯ ╰───────╯     ╰───────╯ ╰───────╯
```

The rendering component itself is timer-independent: the screen owns animation timing and supplies the previous digit, current digit, and animation progress.

---

## Architecture

The NoaLand is organized around three main concepts.

### Application

The top-level application acts as the coordinator.

It is responsible for things such as:

- Registering screens
- Selecting the active screen
- Handling global keyboard input
- Switching between screens
- Future automatic screen rotation
- Tracking terminal dimensions
- Selecting the active layout

The application should not need to know how an individual screen implements its behavior.

### Screens

Each screen owns its own state, behavior, and rendering logic.

Examples include:

```text
GitHub
Flip Clock
Mini World
Weather
News
Agents
Music
...
```

Screens implement a common layout contract defined by The NoaLand.

Layouts are intentionally treated as a small, fixed set of application-level contracts. If a new layout is introduced, every screen is expected to implement it.

This keeps layout support explicit and lets the compiler expose missing implementations.

### UI Components

Screens can be composed from smaller reusable UI components.

Current and future examples include:

```text
GitHub contribution heatmap
GitHub avatar
Flip-clock digit cards
Panels
Indicators
ASCII visualizations
```

Components do not own application-level navigation or screen lifecycle.

They are kept close to the feature that uses them until a genuine shared abstraction emerges.

---

## Responsive Layouts

The NoaLand adapts to the dimensions of the current terminal.

Layout selection is controlled by the application rather than individual screens.

Conceptually:

```text
Terminal Size
     │
     ▼
Layout Resolver
     │
     ▼
Active Layout
     │
     ▼
Active Screen
     │
     ▼
RenderWide / RenderCompactWide / RenderTall / ...
```

This keeps layout behavior consistent across every screen.

---

## Requirements

### Go

The project is written in Go.

Check the required Go version in:

```text
go.mod
```

### GitHub CLI

The GitHub screen currently uses the GitHub CLI:

```bash
gh
```

Authenticate before running:

```bash
gh auth login
```

The screen queries GitHub through:

```bash
gh api graphql
```

### Terminal

Most of The NoaLand is rendered with standard terminal text, Unicode, ANSI styling, Bubble Tea, and Lip Gloss.

Some richer rendering features may depend on additional terminal capabilities.

The GitHub avatar renderer currently uses the Kitty Graphics Protocol where supported.

Terminal capabilities may differ between Ghostty, Windows Terminal, Kitty, WezTerm, and other terminal emulators.

---

## Build

Clone the repository:

```bash
git clone https://github.com/NoaLand/the-noaland.git
cd the-noaland
```

Build everything:

```bash
go build ./...
```

Run tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Run the application:

```bash
go run ./the-noaland/cmd/noaland
```

---

## CI

The project uses GitHub Actions for continuous integration.

The CI pipeline verifies:

```text
gofmt
go vet
go test
go build
```

on changes to the main development branch and pull requests.

---

## Releases

Release builds are produced through GitHub Actions when a version tag is pushed.

For example:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds distributable binaries for supported platforms and publishes them to the corresponding GitHub Release.

---

## v0.1.0

`v0.1.0` is the first milestone after the project evolved from `dev-dashboard` into **The NoaLand**.

It introduces the first two screens:

```text
GitHub
Flip Clock
```

and establishes the direction of the project as a multi-screen terminal environment rather than a single-purpose developer dashboard.

---

## Roadmap

Some experiments currently planned for future screens include:

```text
Mini World Simulator
Daily News
ASCII Weather
AI Agents Monitor
Lyrics with Music
Retro Holographic Display
```

The list is intentionally open-ended.

The NoaLand is also a place to experiment with terminal UI architecture, animation, rendering, interaction, and developer tooling.

---

## Tech Stack

Core dependencies:

```text
Go
Bubble Tea
Lip Gloss
GitHub CLI
GitHub GraphQL API
```

Terminal-specific rendering may additionally make use of protocols such as the Kitty Graphics Protocol.

---

## License

Licensed under the Apache License 2.0.

See: [LICENSE](./LICENSE) for details.

---

*Copyright © 2026 NoaLand*
