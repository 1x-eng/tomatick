package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type ProgressModel struct {
	progress    progress.Model
	total       time.Duration
	elapsed     time.Duration
	done        bool
	theme       *Theme
	description string
}

func NewProgressModel(duration time.Duration, description string, theme *Theme) ProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return ProgressModel{
		progress:    p,
		total:       duration,
		theme:       theme,
		description: description,
	}
}

func (m ProgressModel) Init() tea.Cmd {
	return tick()
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tickMsg:
		if m.done {
			return m, tea.Quit
		}

		m.elapsed += time.Second
		if m.elapsed >= m.total {
			m.done = true
			return m, tea.Quit
		}
		return m, tick()
	}
	return m, nil
}

func (m ProgressModel) View() string {
	if m.done {
		return m.theme.Styles.SuccessText.Render("Done!")
	}

	remainingTime := m.total - m.elapsed
	progress := float64(m.elapsed) / float64(m.total)

	str := strings.Builder{}
	str.WriteString(m.theme.Styles.InfoText.Render(m.description))
	str.WriteString("\n")
	str.WriteString(m.progress.ViewAs(progress))
	str.WriteString("\n")
	str.WriteString(m.theme.Styles.Timer.Render(
		fmt.Sprintf("Remaining: %02d:%02d",
			int(remainingTime.Minutes()),
			int(remainingTime.Seconds())%60,
		),
	))

	return m.theme.Styles.Subtitle.Render(str.String())
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
