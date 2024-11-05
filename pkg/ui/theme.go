package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/logrusorgru/aurora"
)

type Theme struct {
	Aurora aurora.Aurora
	Styles ThemeStyles
	Emoji  ThemeEmoji
}

type ThemeStyles struct {
	Title             lipgloss.Style
	Subtitle          lipgloss.Style
	TaskItem          lipgloss.Style
	TaskNumber        lipgloss.Style
	TaskPrompt        lipgloss.Style
	InfoText          lipgloss.Style
	ErrorText         lipgloss.Style
	SuccessText       lipgloss.Style
	Timer             lipgloss.Style
	SystemInstruction lipgloss.Style
	SystemMessage     lipgloss.Style
	Progress          lipgloss.Style
	Spinner           lipgloss.Style
	AIMessage         lipgloss.Style
	Break             lipgloss.Style
	ChatHeader        lipgloss.Style
	ChatBorder        lipgloss.Style
	UserMessage       lipgloss.Style
	ChatPrompt        lipgloss.Style
	ChatSession       lipgloss.Style
	ChatDivider       lipgloss.Style
}

type ThemeEmoji struct {
	TaskComplete   string
	TaskInProgress string
	TaskPending    string
	Reflection     string
	Timer          string
	Break          string
	Analysis       string
	Warning        string
	Success        string
	Error          string
	Suggestion     string
	Help           string
	Stats          string
	Context        string
	Sound          string
	Brain          string
	Bullet         string
	Section        string
	ChatStart      string
	UserInput      string
	AIResponse     string
	ChatEnd        string
	ChatDivider    string
}

func NewTheme() *Theme {
	return &Theme{
		Aurora: aurora.NewAurora(true),
		Styles: ThemeStyles{
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#E0B0FF")).
				MarginBottom(1).
				Padding(1, 0),

			Subtitle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9DC8C8")).
				MarginBottom(1),

			TaskItem: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D4E2D4")).
				PaddingLeft(2),

			TaskNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#DEBACE")).
				Bold(true),

			InfoText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#B4C8EA")),

			ErrorText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFB4B4")).
				Bold(true),

			SuccessText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#BADEB3")).
				Bold(true),

			Timer: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E6CCB2")).
				Bold(true).
				Padding(0, 1),

			SystemInstruction: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C7B7A3")).
				Italic(true),

			SystemMessage: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#B5B9FF")).
				Italic(true),

			Progress: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D7C0AE")).
				Bold(true),

			Spinner: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#C4B5FD")).
				Bold(true),

			AIMessage: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#B8E7E1")).
				Italic(true),

			Break: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#86EFAC")).
				Bold(true),

			ChatHeader: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9D8CFF")).
				Bold(true).
				Padding(1, 0),

			ChatBorder: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9D8CFF")).
				Bold(true).
				Padding(1, 0),

			UserMessage: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9D8CFF")).
				Bold(true).
				Padding(1, 0),

			ChatPrompt: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9D8CFF")).
				Bold(true).
				Padding(1, 0),

			ChatSession: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9D8CFF")).
				Bold(true).
				Padding(1, 0),

			ChatDivider: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9D8CFF")).
				Bold(true).
				Padding(1, 0),
		},
		Emoji: ThemeEmoji{
			TaskComplete:   "✅",
			TaskInProgress: "⏳",
			TaskPending:    "📝",
			Reflection:     "💭",
			Timer:          "⏰",
			Break:          "🌿",
			Analysis:       "🔍",
			Warning:        "⚠️",
			Success:        "✨",
			Error:          "❌",
			Suggestion:     "💡",
			Help:           "ℹ️",
			Stats:          "📊",
			Context:        "🎯",
			Sound:          "🔔",
			Brain:          "🧠",
			Bullet:         "•",
			Section:        "📋",
			ChatStart:      "👋",
			UserInput:      "👤",
			AIResponse:     "🤖",
			ChatEnd:        "👋",
			ChatDivider:    "🔀",
		},
	}
}
