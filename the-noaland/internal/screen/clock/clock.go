package clock

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/NoaLand/the-noaland/the-noaland/internal/screen"
)

const animationDuration = 360 * time.Millisecond
const frameInterval = 45 * time.Millisecond

type tickMsg struct{ generation uint64 }
type frameMsg struct{ generation uint64 }

// Screen owns the clock and animation state. UI components only render frames.
type Screen struct {
	now               func() time.Time
	current, previous time.Time
	animationStart    time.Time
	progress          float64
	active            bool
	generation        uint64
}

func New() *Screen {
	now := time.Now()
	return &Screen{now: time.Now, current: now, previous: now, progress: 1}
}

// Timers start on activation so background screens do no animation work.
func (s *Screen) Init() tea.Cmd { return nil }

func (s *Screen) Update(msg tea.Msg, _ screen.LayoutContext) tea.Cmd {
	switch msg := msg.(type) {
	case screen.ActivatedMsg:
		s.generation++
		s.active = true
		s.current = s.now()
		s.previous = s.current
		s.progress = 1
		return s.nextSecond()
	case screen.DeactivatedMsg:
		s.generation++
		s.active = false
		s.progress = 1
		return nil
	case tickMsg:
		if !s.active || msg.generation != s.generation {
			return nil
		}
		now := s.now()
		s.previous, s.current = s.current, now
		s.progress = 1
		if s.previous.Format("150405") != now.Format("150405") {
			s.animationStart = now
			s.progress = 0
			return tea.Batch(s.nextSecond(), s.nextFrame())
		}
		return s.nextSecond()
	case frameMsg:
		if !s.active || msg.generation != s.generation || s.progress >= 1 {
			return nil
		}
		s.progress = max(0, min(1, float64(s.now().Sub(s.animationStart))/float64(animationDuration)))
		if s.progress < 1 {
			return s.nextFrame()
		}
	}
	return nil
}

func (s *Screen) nextSecond() tea.Cmd {
	now := s.now()
	delay := now.Truncate(time.Second).Add(time.Second).Sub(now)
	generation := s.generation
	return tea.Tick(delay, func(time.Time) tea.Msg { return tickMsg{generation} })
}

func (s *Screen) nextFrame() tea.Cmd {
	generation := s.generation
	return tea.Tick(frameInterval, func(time.Time) tea.Msg { return frameMsg{generation} })
}
