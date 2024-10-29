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
}

func NewTheme() *Theme {
	return &Theme{
		Aurora: aurora.NewAurora(true),
		Styles: ThemeStyles{
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF79C6")).
				MarginBottom(1).
				Padding(1, 0),

			Subtitle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#50FA7B")).
				MarginBottom(1),

			TaskItem: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8BE9FD")).
				PaddingLeft(2),

			TaskNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFB86C")).
				Bold(true),

			InfoText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6272A4")),

			ErrorText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF5555")).
				Bold(true),

			SuccessText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#50FA7B")).
				Bold(true),

			Timer: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#BD93F9")).
				Bold(true).
				Padding(0, 1),

			Progress: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF79C6")).
				Bold(true),

			Spinner: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8BE9FD")).
				Bold(true),
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
		},
	}
}
