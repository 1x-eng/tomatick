package llm

import (
	"fmt"
)

type RefinementChat struct {
	perplexity *PerplexityAI
	history    []Message
}

func NewRefinementChat(p *PerplexityAI, initialMessages []Message) *RefinementChat {
	return &RefinementChat{
		perplexity: p,
		history:    initialMessages,
	}
}

func (rc *RefinementChat) Chat(userInput string) (string, error) {
	if userInput != "" {
		rc.history = append(rc.history, Message{
			Role:    "user",
			Content: userInput,
		})
	}

	if len(rc.history) <= 2 {
		initialResponse, err := rc.perplexity.GetResponse(rc.history)
		if err != nil {
			return "", err
		}
		rc.history = append(rc.history, Message{
			Role:    "assistant",
			Content: initialResponse,
		})
		return initialResponse, nil
	}

	response, err := rc.perplexity.GetResponse(rc.history)
	if err != nil {
		return "", err
	}

	rc.history = append(rc.history, Message{
		Role:    "assistant",
		Content: response,
	})

	return response, nil
}

func (rc *RefinementChat) GetRefinedContext() (string, error) {
	systemPrompt := `You are an expert at creating focused, goal-oriented session blueprints for Tomatick, a next-generation productivity system. Your task is to take user context and transform it into a clear, actionable plan with specific objectives and time-bound outcomes.

CORE RESPONSIBILITIES:
1. Analyze and restructure user context into clear, actionable plans
2. Ensure work-life balance by setting realistic goals and timeframes
3. Create time-bound roadmaps with specific milestones
4. Identify potential blockers and provide mitigation strategies
5. Help users maintain sustainable work patterns

CONTEXT ANALYSIS GUIDELINES:
• Consider work-life balance implications
• Account for personal wellbeing and family time
• Set realistic expectations and timeframes
• Identify signs of overwork or unsustainable patterns
• Suggest breaks and boundaries where needed

OUTPUT REQUIREMENTS:
Your response must follow this EXACT format:

CORE OBJECTIVE:
[One clear, measurable end goal that this session aims to achieve]

CONTEXT ESSENCE:
[3-5 bullet points distilling the key information and priorities]

TIME-BOUND ROADMAP:
[Break down by specific time blocks, e.g.:
2:00 PM - 2:40 PM:
• Task 1 (25 min)
• Task 2 (15 min)
...]

SUCCESS CRITERIA:
[2-3 specific, measurable outcomes that define success]

FOCUS AREAS:
[Key areas requiring attention, in priority order]

WORK-LIFE BALANCE CONSIDERATIONS:
[Specific recommendations for maintaining balance]

POTENTIAL BLOCKERS:
[Identify challenges and mitigation strategies]

Remember:
1. Every point must directly contribute to the core objective
2. Be specific and actionable
3. Include clear time estimates
4. Focus on measurable outcomes
5. Keep it concise but comprehensive
6. Emphasize sustainable work patterns
7. Consider work-life balance in all recommendations`

	// Reset history to just have the system prompt and original context
	rc.history = []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		rc.history[len(rc.history)-1], // Keep the last message which contains the context
	}

	response, err := rc.perplexity.GetResponse(rc.history)
	if err != nil {
		return "", err
	}

	// Update history with the response
	rc.history = append(rc.history, Message{
		Role:    "assistant",
		Content: response,
	})

	return response, nil
}

func (rc *RefinementChat) RequestContextModification(originalContext, userFeedback string) (string, error) {
	// Create modification request
	modificationRequest := fmt.Sprintf(`Previous blueprint:
%s

User feedback for modifications:
%s

Please provide an updated blueprint that incorporates this feedback. Maintain the same structured format with all sections (CORE OBJECTIVE, CONTEXT ESSENCE, etc.) while addressing the feedback.`, originalContext, userFeedback)

	// Set up a fresh conversation with proper role alternation
	rc.history = []Message{
		{
			Role:    "system",
			Content: rc.history[0].Content, // Keep the system prompt
		},
		{
			Role:    "user",
			Content: modificationRequest,
		},
	}

	// Get the response
	response, err := rc.perplexity.GetResponse(rc.history)
	if err != nil {
		return "", err
	}

	// Add the assistant's response to history
	rc.history = append(rc.history, Message{
		Role:    "assistant",
		Content: response,
	})

	return response, nil
}
