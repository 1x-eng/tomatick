package llm

import (
	"fmt"
	"strings"
)

type Assistant struct {
	perplexity *PerplexityAI
	context    string
}

func NewAssistant(p *PerplexityAI, context string) *Assistant {
	return &Assistant{
		perplexity: p,
		context:    context,
	}
}

func (a *Assistant) GetTaskSuggestions(currentTasks []string) ([]string, error) {
	tasksStr := strings.Join(currentTasks, "\n")
	prompt := fmt.Sprintf(`Given the following context and current tasks, suggest 3 additional focused tasks that align with the context and current work:

Context:
%s

Current Tasks:
%s

Provide only the tasks, one per line, without numbers or bullets.`, a.context, tasksStr)

	messages := []Message{
		{Role: "system", Content: "You are a focused task management assistant. Keep suggestions concise and actionable."},
		{Role: "user", Content: prompt},
	}

	response, err := a.perplexity.GetResponse(messages)
	if err != nil {
		return nil, err
	}

	suggestions := strings.Split(strings.TrimSpace(response), "\n")
	return suggestions, nil
}

func (a *Assistant) AnalyzeProgress(completedTasks []string, reflections string) (string, error) {
	prompt := fmt.Sprintf(`Given the following context, completed tasks, and reflections, provide a brief analysis of progress and suggestions for the next cycle:

Context:
%s

Completed Tasks:
%s

Reflections:
%s

Provide a concise analysis focusing on alignment with context and actionable next steps.`, a.context, strings.Join(completedTasks, "\n"), reflections)

	messages := []Message{
		{Role: "system", Content: "You are a productivity analysis assistant. Keep feedback constructive and forward-looking."},
		{Role: "user", Content: prompt},
	}

	return a.perplexity.GetResponse(messages)
}
