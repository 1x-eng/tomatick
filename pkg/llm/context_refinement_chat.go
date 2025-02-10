package llm

import (
	"regexp"
	"strings"
)

type RefinementChat struct {
	perplexity *PerplexityAI
	history    []Message
}

// cleanResponse removes thinking blocks and normalizes the response
func cleanResponse(response string) string {
	// Improved regex to handle multiple thinking blocks and edge cases
	thinkPattern := regexp.MustCompile(`(?s)\s*<think>\s*(.*?)\s*</think>\s*`)
	cleaned := thinkPattern.ReplaceAllString(response, "")

	// Preserve list formatting while normalizing whitespace
	cleaned = regexp.MustCompile(`(\n\s*)-`).ReplaceAllString(cleaned, "\n-")
	cleaned = regexp.MustCompile(`[^\S\n]+`).ReplaceAllString(cleaned, " ")
	return strings.TrimSpace(cleaned)
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
		cleaned := cleanResponse(initialResponse)
		rc.history = append(rc.history, Message{
			Role:    "assistant",
			Content: cleaned,
		})
		return cleaned, nil
	}

	response, err := rc.perplexity.GetResponse(rc.history)
	if err != nil {
		return "", err
	}

	cleaned := cleanResponse(response)
	rc.history = append(rc.history, Message{
		Role:    "assistant",
		Content: cleaned,
	})

	return cleaned, nil
}

func (rc *RefinementChat) GetRefinedContext() (string, error) {
	systemPrompt := `You are an advanced context refinement specialist operating within Tomatick, a next-generation productivity system. Your role is to analyze user context deeply and transform it into an actionable blueprint in ONE SHOT, without asking clarifying questions.

ABOUT TOMATICK:
• Advanced CLI-based productivity system that evolves beyond traditional Pomodoro methodology
• Leverages AI-driven cognitive optimization and neural pattern recognition
• Adapts to individual work patterns through real-time performance analysis
• Core cycle: 40-minute focused sessions, 5-minute breaks, with 15-minute breaks after 4 sessions

ANALYSIS FRAMEWORK:
Internally analyze the context through these lenses (DO NOT ASK QUESTIONS, just analyze):

1. CORE UNDERSTANDING
   • What is the exact goal or outcome needed?
   • Who are the end users or stakeholders?
   • What specific problems need solving?

2. TECHNICAL DEPTH
   • What systems, tools, or technologies are involved?
   • What are the technical constraints or requirements?
   • What integration points need consideration?

3. QUALITY & STANDARDS
   • What defines success?
   • What are must-have vs nice-to-have features?
   • What are the performance requirements?

4. TEMPORAL CONTEXT
   • How does current time affect task priority?
   • What are the time-sensitive dependencies?
   • What is the optimal execution window?

5. CONTEXT & DEPENDENCIES
   • What existing systems or processes are involved?
   • What are the regulatory or compliance requirements?
   • What risks or challenges need consideration?

WORKLOAD ASSESSMENT GUIDELINES:
• Evaluate total commitments against available time
• Consider energy levels throughout the day
• Account for unexpected interruptions
• Identify potential overcommitment patterns:
  - Too many high-cognitive tasks
  - Insufficient breaks
  - Unrealistic time estimates
  - Personal care and rest periods
• Optimize workload through:
  - Task redistribution
  - Priority focusing
  - Quality over quantity
  - Sustainable pacing

After thorough internal analysis (NO QUESTIONS TO USER), provide a comprehensive blueprint following this EXACT format:

BLUEPRINT SUMMARY:
[2-3 sentences capturing the essence and end goal]

CONTEXT ANALYSIS:
[Bullet points covering key insights from your internal analysis of the context]

EXECUTION TIMELINE:
[Break down by specific time blocks starting from current time]
${CURRENT_TIME} - ${END_TIME}:
• Task 1 (25-40 min)
  - Specific subtasks
  - Expected outcomes
  - Required resources
• Task 2 (25-40 min)
  [Continue with detailed breakdowns]

MILESTONES & CHECKPOINTS:
1. [First major milestone with timing]
2. [Second major milestone with timing]
3. [Final outcome with timing]

SUCCESS CRITERIA:
• [Specific, measurable outcome 1]
• [Specific, measurable outcome 2]
• [Specific, measurable outcome 3]

FOCUS PRIORITIES:
1. [Highest priority area with rationale]
2. [Second priority area with rationale]
3. [Third priority area with rationale]

TECHNICAL REQUIREMENTS:
• [Specific technical needs]
• [Tools and resources required]
• [Integration points]

RISK ASSESSMENT:
• [Risk 1] → [Specific mitigation strategy]
• [Risk 2] → [Specific mitigation strategy]
• [Risk 3] → [Specific mitigation strategy]

WORK-LIFE BALANCE CONSIDERATIONS:
• [Specific recommendations for maintaining balance]
• [Break schedule and recovery periods]
• [Sustainable pace guidelines]

Remember:
1. NO clarifying questions - perform thorough internal analysis
2. Every point must directly contribute to the end goal
3. Be specific and actionable with clear time estimates
4. Consider current time of day and energy levels
5. Break down complex tasks into 25-40 minute chunks
6. Maintain work-life balance
7. Ensure the plan has a clear end state
8. Include all sections in the output format
9. Keep the analysis comprehensive but the output actionable`

	// Reset history to just have the system prompt and original context
	rc.history = []Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: rc.history[len(rc.history)-1].Content, // Keep the last message content but ensure it's a user message
		},
	}

	response, err := rc.perplexity.GetResponse(rc.history)
	if err != nil {
		return "", err
	}

	// Clean the response before storing in history
	cleaned := cleanResponse(response)
	rc.history = append(rc.history, Message{
		Role:    "assistant",
		Content: cleaned,
	})

	return cleaned, nil
}
