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
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	TaskItem    lipgloss.Style
	TaskNumber  lipgloss.Style
	TaskPrompt  lipgloss.Style
	InfoText    lipgloss.Style
	ErrorText   lipgloss.Style
	SuccessText lipgloss.Style
	Timer       lipgloss.Style
	Border      lipgloss.Style
	Progress    lipgloss.Style
}

func NewTheme() *Theme {
	return &Theme{
		Aurora: aurora.NewAurora(true),
		Styles: ThemeStyles{
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF6B6B")).
				MarginBottom(1),

			Subtitle: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				MarginBottom(1),

			TaskItem: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#95E1D3")).
				PaddingLeft(2),

			TaskNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFE66D")).
				Bold(true),

			TaskPrompt: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A8E6CF")).
				PaddingLeft(2),

			InfoText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#89ABE3")),

			ErrorText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true),

			SuccessText: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#A8E6CF")).
				Bold(true),

			Timer: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFE66D")).
				Bold(true).
				Padding(0, 1),

			Border: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#4ECDC4")).
				Padding(1),

			Progress: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true),
		},
	}
}
