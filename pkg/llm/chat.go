package llm

import (
	"fmt"
	"strings"
)

type SuggestionChat struct {
	assistant    *Assistant
	history      []Message
	context      string
	suggestions  []string
	lastAnalysis string
}

func NewSuggestionChat(assistant *Assistant, initialContext string, suggestions []string, lastAnalysis string) *SuggestionChat {
	return &SuggestionChat{
		assistant:    assistant,
		history:      make([]Message, 0),
		context:      initialContext,
		suggestions:  suggestions,
		lastAnalysis: lastAnalysis,
	}
}

func (sc *SuggestionChat) Chat(userInput string) (string, error) {
	// Add user message to history
	sc.history = append(sc.history, Message{
		Role:    "user",
		Content: userInput,
	})

	systemPrompt := `You are an advanced task optimization assistant engaged in a discussion about specific task suggestions. Your core responsibilities:

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

RESPONSE GUIDELINES:
1. If question is relevant:
   - Provide clear, structured response
   - Include specific steps or clarifications
   - Reference context when applicable
   - Maintain focus on task completion

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

Current Suggestions Under Discussion:
"""
%s
"""`

	messages := []Message{
		{
			Role:    "system",
			Content: fmt.Sprintf(systemPrompt, sc.lastAnalysis, sc.context, strings.Join(sc.suggestions, "\n")),
		},
	}

	// Add chat history
	messages = append(messages, sc.history...)

	// Get response from llm
	response, err := sc.assistant.perplexity.GetResponse(messages)
	if err != nil {
		return "", err
	}

	// Add assistant's response to history
	sc.history = append(sc.history, Message{
		Role:    "assistant",
		Content: response,
	})

	return response, nil
}
