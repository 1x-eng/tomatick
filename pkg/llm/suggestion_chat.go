package llm

import (
	"fmt"
	"strings"
)

type SuggestionChat struct {
	assistant      *Assistant
	history        []Message
	context        string
	suggestions    []string
	lastAnalysis   string
	acceptedTasks  []string
	completedTasks string
	reflections    string
}

func NewSuggestionChat(assistant *Assistant, initialContext string, suggestions []string, lastAnalysis string, acceptedTasks []string, completedTasks string, reflections string) *SuggestionChat {
	return &SuggestionChat{
		assistant:      assistant,
		history:        make([]Message, 0),
		context:        initialContext,
		suggestions:    suggestions,
		lastAnalysis:   lastAnalysis,
		acceptedTasks:  acceptedTasks,
		completedTasks: completedTasks,
		reflections:    reflections,
	}
}

func (sc *SuggestionChat) Chat(userInput string) (string, error) {
	// Add user message to history
	sc.history = append(sc.history, Message{
		Role:    "user",
		Content: userInput,
	})

	var systemPrompt string
	if len(sc.suggestions) > 0 {
		systemPrompt = `You are an advanced task optimization assistant engaged in a discussion about specific task suggestions. Your core responsibilities:

CONTEXT AWARENESS:
- Maintain strict relevance to the session context and current suggestions
- Detect and flag off-topic or digressing questions
- Guide users back to relevant discussion points

SUGGESTION CLARIFICATION:
- Provide detailed, actionable explanations for suggestions
- Break down complex tasks into clear, achievable steps
- Highlight dependencies and prerequisites
- Explain the reasoning behind each suggestion
- Focus on practical implementation details

TASK COMPLETION STATUS:
- Be aware of the completion status of tasks from the last session
- Use the task completion status to guide your suggestions
- Address user questions and concerns about task completion and hence, the need for suggested tasks
- Maintain focus on performance optimization

USER REFLECTIONS:
- Use user reflections to guide your suggestions
- Explain the reasoning behind suggestions based on user reflections
- Offer concrete examples and evidence
- Address user questions and concerns about the need for suggested tasks
- Maintain focus on performance optimization

RESPONSE GUIDELINES:
1. If question is relevant:
   - Provide clear, structured response
   - Include specific details and clarifications
   - Reference context when applicable
   - Maintain focus on improvement

2. If question seems off-topic:
   - Politely flag the digression
   - Explain why it seems unrelated
   - Offer to hear user's perspective
   - Guide back to relevant discussion

3. For implementation queries:
   - Break down into concrete steps
   - Highlight potential challenges
   - Suggest specific approaches
   - Focus on actionability

Last Session Analysis:
"""
%s
"""

Current Session Context:
"""
%s
"""
`
	} else {
		systemPrompt = `You are an advanced performance analysis assistant engaged in a discussion about the session analysis. Your core responsibilities:

ANALYSIS CLARIFICATION:
- Provide detailed explanations of analysis points
- Explain the reasoning behind observations
- Offer concrete examples and evidence
- Address user questions and concerns
- Maintain focus on performance optimization

TASK CONTEXT:
- Consider the specific tasks that were undertaken
- Analyze completion patterns and challenges
- Reference specific tasks when discussing insights
- Connect analysis points to actual task outcomes

USER REFLECTIONS:
- Incorporate user's original reflections
- Connect analysis insights to user observations
- Address any gaps between user reflections and analysis
- Provide deeper insights into user's observations

RESPONSE GUIDELINES:
1. If question is relevant:
   - Provide clear, structured response
   - Include specific details and clarifications
   - Reference context when applicable
   - Maintain focus on improvement

2. If question seems off-topic:
   - Politely flag the digression
   - Explain why it seems unrelated
   - Offer to hear user's perspective
   - Guide back to relevant discussion

3. For implementation queries:
   - Break down into concrete steps
   - Highlight potential challenges
   - Suggest specific approaches
   - Focus on actionability

Last Session Analysis:
"""
%s
"""

Current Session Context:
"""
%s
"""

Tasks and Their Status:
"""
%s
"""

User's Original Reflections:
"""
%s
"""`
	}

	var args []interface{}
	if len(sc.suggestions) > 0 {
		args = []interface{}{sc.lastAnalysis, sc.context, strings.Join(sc.suggestions, "\n")}
	} else {
		args = []interface{}{sc.lastAnalysis, sc.context, sc.completedTasks, sc.reflections}
	}

	messages := []Message{
		{
			Role:    "system",
			Content: fmt.Sprintf(systemPrompt, args...),
		},
	}

	messages = append(messages, sc.history...)

	response, err := sc.assistant.perplexity.GetResponse(messages)
	if err != nil {
		return "", err
	}

	sc.history = append(sc.history, Message{
		Role:    "assistant",
		Content: response,
	})

	return response, nil
}
