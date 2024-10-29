package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	frames       []string
	current      int
	style        lipgloss.Style
	messages     []SpinnerMessage
	messageIndex int
}

type SpinnerMessage struct {
	Text  string
	Emoji string
}

func NewSpinner(style lipgloss.Style) *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		style:  style,
		messages: []SpinnerMessage{
			{"Analyzing patterns", "🧠"},
			{"Processing insights", "✨"},
			{"Optimizing flow", "🌊"},
			{"Calibrating focus", "🎯"},
			{"Enhancing clarity", "💫"},
			{"Synthesizing data", "📊"},
			{"Refining suggestions", "💡"},
			{"Mapping connections", "🔄"},
			{"Elevating performance", "📈"},
			{"Harmonizing workflow", "🎵"},
		},
	}
}

func (s *Spinner) Next() string {
	frame := s.style.Render(s.frames[s.current])
	message := s.messages[s.messageIndex]

	// Update indices
	s.current = (s.current + 1) % len(s.frames)
	if s.current == 0 {
		s.messageIndex = (s.messageIndex + 1) % len(s.messages)
	}

	return fmt.Sprintf("%s %s %s",
		frame,
		message.Emoji,
		s.style.Render(message.Text))
}
