package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/logrusorgru/aurora"
)

type Theme struct {
	Aurora aurora.Aurora
	Styles ThemeStyles
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

func NewTheme() *Theme {
	return &Theme{
		Aurora: aurora.NewAurora(true),
		Styles: ThemeStyles{
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F92672")).
				MarginBottom(1),

			Subtitle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A6E22E")).
				MarginBottom(1),

			TaskItem: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#66D9EF")).
				PaddingLeft(2),

			TaskNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FD971F")).
				Bold(true),

			TaskPrompt: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E6DB74")).
				PaddingLeft(2),

			InfoText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#75715E")),

			ErrorText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F92672")).
				Bold(true),

			SuccessText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A6E22E")).
				Bold(true),

			Timer: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FD971F")).
				Bold(true).
				Padding(0, 1),

			Progress: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#AE81FF")).
				Bold(true),

			Spinner: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#66D9EF")).
				Bold(true),

			SystemInstruction: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E69F66")).
				Bold(true),

			SystemMessage: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#75715E")).
				Bold(true),

			AIMessage: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#007ACC")).
				Bold(true),
		},
	}
}
