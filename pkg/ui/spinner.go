package ui

import "github.com/charmbracelet/lipgloss"

type Spinner struct {
	frames  []string
	current int
	style   lipgloss.Style
}

func NewSpinner(style lipgloss.Style) *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		style:  style,
	}
}

func (s *Spinner) Next() string {
	frame := s.style.Render(s.frames[s.current])
	s.current = (s.current + 1) % len(s.frames)
	return frame
}
