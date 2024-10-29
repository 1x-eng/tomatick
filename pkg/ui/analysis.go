package ui

import (
	"fmt"
	"strings"
)

type AnalysisPresenter struct {
	theme *Theme
}

func NewAnalysisPresenter(theme *Theme) *AnalysisPresenter {
	return &AnalysisPresenter{theme: theme}
}

func (ap *AnalysisPresenter) Present(analysis string) string {
	sections := ap.parseAnalysis(analysis)
	return ap.formatSections(sections)
}

func (ap *AnalysisPresenter) parseAnalysis(analysis string) map[string][]string {
	sections := make(map[string][]string)
	var currentSection string
	var currentContent []string

	lines := strings.Split(analysis, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			if currentSection != "" {
				sections[currentSection] = currentContent
			}
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "##"))
			currentContent = []string{}
		} else if strings.TrimSpace(line) != "" {
			currentContent = append(currentContent, strings.TrimSpace(line))
		}
	}
	if currentSection != "" {
		sections[currentSection] = currentContent
	}
	return sections
}

func (ap *AnalysisPresenter) formatSections(sections map[string][]string) string {
	var sb strings.Builder

	sb.WriteString(ap.theme.Styles.Subtitle.Render("\n🤖 Your copilot's analysis\n"))

	for section, content := range sections {
		sb.WriteString(ap.theme.Styles.TaskNumber.Render(
			fmt.Sprintf("%s %s\n",
				ap.theme.Emoji.Section,
				(section))))

		for _, line := range content {
			line = strings.TrimSpace(line)
			line = strings.ReplaceAll(line, "**", "")
			line = strings.TrimLeft(line, "#-")

			if line != "" {
				if strings.HasPrefix(line, "•") {
					sb.WriteString(line + "\n")
				} else {
					sb.WriteString(fmt.Sprintf("%s %s\n",
						ap.theme.Emoji.Bullet,
						line))
				}
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
