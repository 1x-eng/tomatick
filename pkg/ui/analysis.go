package ui

import (
	"fmt"
	"regexp"
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
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "##") {
			if currentSection != "" {
				sections[currentSection] = currentContent
			}
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "##"))
			currentContent = []string{}
		} else if line != "" {
			// Normalize bullet points and numbering
			line = ap.normalizeListItem(line)
			currentContent = append(currentContent, line)
		}
	}
	if currentSection != "" {
		sections[currentSection] = currentContent
	}
	return sections
}

func (ap *AnalysisPresenter) normalizeListItem(line string) string {
	// Remove any existing bullets or numbers
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "•")
	line = strings.TrimPrefix(line, "-")
	line = strings.TrimPrefix(line, "*")

	// Remove numbered prefixes (e.g., "1.", "2.")
	if matched, _ := regexp.MatchString(`^\d+\.`, line); matched {
		line = regexp.MustCompile(`^\d+\.`).ReplaceAllString(line, "")
	}

	return strings.TrimSpace(line)
}

func (ap *AnalysisPresenter) formatSections(sections map[string][]string) string {
	var sb strings.Builder

	// Initial title with single newline
	sb.WriteString(ap.theme.Styles.Subtitle.Render("🤖 Your copilot's analysis"))

	for section, content := range sections {
		// Add section header with consistent spacing
		sb.WriteString(fmt.Sprintf("\n\n%s %s",
			ap.theme.Emoji.Section,
			ap.theme.Styles.TaskNumber.Render(section)))

		// Process content items
		for _, line := range content {
			line = strings.TrimSpace(line)
			line = strings.ReplaceAll(line, "**", "")

			if line != "" {
				sb.WriteString(fmt.Sprintf("\n%s %s",
					ap.theme.Emoji.Bullet,
					ap.theme.Styles.InfoText.Render(line)))
			}
		}
	}

	return sb.String()
}
