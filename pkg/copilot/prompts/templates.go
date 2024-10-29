package prompts

import (
	"fmt"
	"strings"

	"github.com/1x-eng/tomatick/pkg/copilot"
)

const (
	systemPromptFollow = `You are an AI assistant helping with task management in a Pomodoro-style work session.
Your role is to:
1. Help break down tasks into manageable chunks
2. Suggest improvements to task descriptions
3. Identify potential blockers
4. Provide gentle reminders about best practices
Be concise and practical in your suggestions.`

	systemPromptLead = `You are an AI assistant actively guiding a Pomodoro-style work session.
Your role is to:
1. Proactively suggest task organization
2. Recommend task prioritization
3. Guide time management
4. Provide strategic work suggestions
Be direct but supportive in your guidance.`
)

type PromptBuilder struct {
	mode copilot.Mode
}

func NewPromptBuilder(mode copilot.Mode) *PromptBuilder {
	return &PromptBuilder{mode: mode}
}

func (pb *PromptBuilder) GetSystemPrompt() string {
	if pb.mode == copilot.LeadMode {
		return systemPromptLead
	}
	return systemPromptFollow
}

func (pb *PromptBuilder) BuildContextAnalysisPrompt(context string, tasks []string) string {
	var b strings.Builder

	b.WriteString("Current work context:\n")
	b.WriteString(context)
	b.WriteString("\n\nCurrent tasks:\n")
	for i, task := range tasks {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, task))
	}

	if pb.mode == copilot.LeadMode {
		b.WriteString("\nProvide strategic suggestions for task organization and execution.")
	} else {
		b.WriteString("\nProvide helpful observations and suggestions for task improvement.")
	}

	return b.String()
}
