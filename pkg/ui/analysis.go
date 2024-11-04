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

		// Match both ## and **Section** headers
		if strings.HasPrefix(line, "##") || (strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**")) {
			if currentSection != "" {
				sections[currentSection] = currentContent
			}
			// Clean up section header
			currentSection = strings.TrimSpace(strings.Trim(strings.TrimPrefix(strings.TrimPrefix(line, "##"), "**"), "**"))
			currentContent = []string{}
		} else if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•") {
			// Normalize bullet points and add to current content
			line = ap.normalizeListItem(line)
			if line != "" {
				currentContent = append(currentContent, line)
			}
		} else if line != "" && currentSection != "" {
			// Add non-empty lines that aren't headers or bullets
			currentContent = append(currentContent, line)
		}
	}

	// Don't forget the last section
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

	// Initial title with double newline and border
	sb.WriteString("\n" + ap.theme.Styles.Subtitle.Render("🤖 Your copilot's analysis"))
	sb.WriteString("\n\n" + ap.theme.Styles.Subtitle.Render(strings.Repeat("─", 50)) + "\n")

	// If sections is empty, add a message
	if len(sections) == 0 {
		sb.WriteString("\n" + ap.theme.Styles.InfoText.Render("Analysis processing..."))
		return sb.String()
	}

	for section, content := range sections {
		// Add section header with consistent spacing
		sb.WriteString(fmt.Sprintf("\n%s %s\n",
			ap.theme.Emoji.Section,
			ap.theme.Styles.TaskNumber.Render(section)))

		// Process content items
		for _, line := range content {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Remove any markdown formatting
			line = strings.ReplaceAll(line, "**", "")
			line = strings.ReplaceAll(line, "*", "")
			line = strings.ReplaceAll(line, "[", "")
			line = strings.ReplaceAll(line, "]", "")

			sb.WriteString(fmt.Sprintf("%s %s\n",
				ap.theme.Emoji.Bullet,
				ap.theme.Styles.InfoText.Render(line)))
		}
	}

	// Add bottom border
	sb.WriteString("\n" + ap.theme.Styles.Subtitle.Render(strings.Repeat("─", 50)))

	return sb.String()
}
