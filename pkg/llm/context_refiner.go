package llm

import (
	"fmt"
	"time"
)

type ContextRefiner struct {
	perplexity *PerplexityAI
	context    string
}

func NewContextRefiner(p *PerplexityAI, context string) *ContextRefiner {
	return &ContextRefiner{
		perplexity: p,
		context:    context,
	}
}

func (cr *ContextRefiner) StartRefinement() (*RefinementChat, error) {
	messages := []Message{
		{
			Role: "system",
			Content: `You are an advanced context refinement specialist operating within Tomatick. Your role is to analyze context deeply and create actionable blueprints in ONE SHOT.

ABOUT TOMATICK:
• Advanced CLI-based productivity system evolving beyond traditional Pomodoro
• Leverages AI-driven cognitive optimization and pattern recognition
• Core cycle: 40-minute focused sessions, 5-minute breaks, 15-minute breaks after 4 sessions

YOUR ROLE:
1. Analyze user context through multiple lenses WITHOUT asking questions
2. Transform raw input into structured, time-bound blueprints
3. Ensure optimal cognitive performance through strategic planning
4. Create measurable success metrics for each session
5. Facilitate deep work states through clear objective setting

ANALYSIS REQUIREMENTS:
• Perform comprehensive internal analysis using provided framework
• Consider all temporal and cognitive factors
• Evaluate workload and sustainability
• Create detailed, time-aware execution plan
• Ensure clear end goals and success criteria

NO QUESTIONS - just analyze and create the blueprint based on available information.`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf("Here's the initial context to refine:\n\nCurrent Date/Time: %s\nDay of Week: %s\nHour Category: %s\n\n%s",
				time.Now().Format("2006-01-02 15:04 Z07:00"),
				time.Now().Weekday().String(),
				getHourCategory(time.Now().Hour()),
				cr.context),
		},
	}

	return NewRefinementChat(cr.perplexity, messages), nil
}

func getHourCategory(hour int) string {
	switch {
	case hour >= 5 && hour < 8:
		return "Early Morning (Optimal for: personal excellence - meditation, exercise, reading, planning the day ahead)"

	case hour >= 8 && hour < 12:
		return "Peak Morning (Optimal for: complex problem-solving, creative work, critical thinking, important meetings, learning new skills)"

	case hour >= 12 && hour < 14:
		return "Midday (Ideal for: lunch break, light exercise, family meal, quick errands, social connections, mindful rest)"

	case hour >= 14 && hour < 17:
		return "Afternoon (Suitable for: collaborative work, routine tasks, administrative duties, follow-ups, mentoring)"

	case hour >= 17 && hour < 19:
		return "Early Evening (Priority for: family time, children's activities, household management, meal preparation, light chores)"

	case hour >= 19 && hour < 21:
		return "Evening (Focus on: family bonding, personal hobbies, relationship building, gentle exercise, planning next day)"

	case hour >= 21 && hour < 23:
		return "Late Evening (Transition to: relaxation, reflection, light reading, mindfulness, preparing for rest)"

	case hour >= 23 || hour < 5:
		return "Night (Protected time for: sleep, recovery, restoration - avoid scheduling tasks unless absolutely necessary)"

	default:
		// This should never happen. just complying with the switch statement.
		return "Night (Protected time for: sleep, recovery, restoration)"
	}
}
