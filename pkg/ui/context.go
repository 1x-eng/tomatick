package ui

import (
	"fmt"
	"strings"
)

type ContextPresenter struct {
	theme *Theme
}

func NewContextPresenter(theme *Theme) *ContextPresenter {
	return &ContextPresenter{theme: theme}
}

func (cp *ContextPresenter) PresentContextMenu() string {
	var sb strings.Builder

	// Header with decorative border
	border := strings.Repeat("═", 60)
	sb.WriteString(cp.theme.Styles.Title.Render(border))
	sb.WriteString(cp.theme.Styles.Title.Render("📋 Context Management"))
	sb.WriteString(cp.theme.Styles.SystemInstruction.Render(`
	───────────────────────────────────────────────
				Why Context Matters
	───────────────────────────────────────────────
	
	Your session context helps the AI:
	• Understand your current focus areas
	• Provide more relevant task suggestions
	• Track progress across sessions
	• Prevent context switching
	• Optimize for sustainable progress
	
	Your copilot uses this information to:
	• Calibrate task difficulty
	• Manage cognitive load
	• Maintain strategic momentum
	• Prevent burnout
	───────────────────────────────────────────────
	`))
	sb.WriteString(cp.theme.Styles.Title.Render(border))

	sb.WriteString(cp.theme.Styles.Subtitle.Render("\nOptions:"))

	options := []struct {
		emoji string
		title string
	}{
		{cp.theme.Emoji.Context, "Load Existing Context"},
		{cp.theme.Emoji.Brain, "Create New Context"},
	}

	for _, opt := range options {
		sb.WriteString(fmt.Sprintf("\n%s %s\n",
			opt.emoji,
			cp.theme.Styles.TaskNumber.Render(opt.title)))
	}

	return sb.String()
}

func (cp *ContextPresenter) PresentContextList(contexts []string) string {
	var sb strings.Builder

	sb.WriteString("\n\n📚 Available Contexts")
	sb.WriteString("\n" + strings.Repeat("─", 40))

	if len(contexts) == 0 {
		sb.WriteString("\n" + cp.theme.Styles.InfoText.Render("No saved contexts found"))
		return sb.String()
	}

	for _, ctx := range contexts {
		sb.WriteString(fmt.Sprintf("\n%s %s",
			cp.theme.Emoji.Bullet,
			cp.theme.Styles.TaskItem.Render(strings.TrimSpace(ctx))))
	}

	sb.WriteString("\n" + strings.Repeat("─", 40))
	sb.WriteString("\n\n")

	return sb.String()
}

func (cp *ContextPresenter) PresentContextInput() string {
	var sb strings.Builder

	// Header
	sb.WriteString(cp.theme.Styles.Title.Render("\n✨ New Context Creation"))
	sb.WriteString(cp.theme.Styles.Subtitle.Render("\n" + strings.Repeat("─", 40)))

	// Guidelines
	guidelines := []string{
		"What are you working on?",
		"What are your main objectives?",
		"Any specific challenges to address?",
		"Expected outcomes or deliverables?",
	}

	for _, guide := range guidelines {
		sb.WriteString(fmt.Sprintf("\n%s %s",
			cp.theme.Emoji.Bullet,
			cp.theme.Styles.InfoText.Render(guide)))
	}

	sb.WriteString("\n\n" + cp.theme.Styles.SystemInstruction.Render(
		"Type your context below (type 'done' when finished)"))
	sb.WriteString("\n> ")

	return sb.String()
}

func (cp *ContextPresenter) PresentDeltaContextInput() string {
	var sb strings.Builder

	// Header
	sb.WriteString(cp.theme.Styles.Title.Render("\n✨ Additional Context"))
	sb.WriteString(cp.theme.Styles.Subtitle.Render("\n" + strings.Repeat("─", 40)))

	// Guidelines
	guidelines := []string{
		"What additional context would you like to add?",
		"Any new objectives or focus areas?",
		"Specific details for this session?",
		"Recent developments or changes?",
	}

	for _, guide := range guidelines {
		sb.WriteString(fmt.Sprintf("\n%s %s",
			cp.theme.Emoji.Bullet,
			cp.theme.Styles.InfoText.Render(guide)))
	}

	sb.WriteString("\n\n" + cp.theme.Styles.SystemInstruction.Render(
		"Type your additional context below (type 'done' when finished)"))
	sb.WriteString("\n> ")

	return sb.String()
}

func (cp *ContextPresenter) PresentRefinementOption() string {
	var sb strings.Builder

	// Header
	sb.WriteString(cp.theme.Styles.Title.Render("\n🧠 Context Refinement"))
	sb.WriteString(cp.theme.Styles.Subtitle.Render("\n" + strings.Repeat("─", 40)))

	// Description
	description := `
Your copilot can help fine-tune or refine your context by:
• Analyzing scope and objectives
• Identifying missing details
• Clarifying ambiguities
• Ensuring completeness
• Structuring information effectively`

	sb.WriteString(cp.theme.Styles.InfoText.Render(description))
	sb.WriteString("\n\n")
	return sb.String()
}

func (cp *ContextPresenter) PresentRefinementStart() string {
	var sb strings.Builder

	// Header with emoji and title
	sb.WriteString(cp.theme.Styles.ChatHeader.Render("\n🤖 Context Refinement Session"))
	sb.WriteString("\n" + cp.theme.Styles.Subtitle.Render(strings.Repeat("─", 40)))

	// Instructions in a styled box
	instructions := `
Your AI copilot will guide you through refining your context:

📝 Answer each question thoroughly
⌨️  Type your response naturally
✅ Type 'done' to finish your answer
❌ Type 'exit' to keep original context

The copilot will analyze your responses and help improve the context.`

	// Create a styled box for instructions
	instructionBox := cp.theme.Styles.InfoBox.Render(instructions)
	sb.WriteString("\n" + instructionBox + "\n")

	return sb.String()
}

func (cp *ContextPresenter) PresentRefinedContext(refinedContext string) string {
	var sb strings.Builder

	sb.WriteString(cp.theme.Styles.Title.Render("\n✨ Refined Context ✨"))
	sb.WriteString("\n" + cp.theme.Styles.ChatBorder.Render(strings.Repeat("─", 40)))

	cleanContext := strings.TrimSpace(refinedContext)
	if strings.HasPrefix(strings.ToLower(cleanContext), "context refinement complete") {
		if idx := strings.Index(cleanContext, "•"); idx != -1 {
			cleanContext = cleanContext[idx:]
		}
	}

	sb.WriteString("\n\n" + cp.theme.Styles.RefinedContext.Render(cleanContext))

	reviewPrompt := `
Please review the refined context above:
• Is it accurate and complete?
• Does it capture all key points?
• Are there any missing details?`

	sb.WriteString("\n" + cp.theme.Styles.SystemInstruction.Render(reviewPrompt))
	return sb.String()
}

func (cp *ContextPresenter) GetTheme() *Theme {
	return cp.theme
}
