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

func (a *Assistant) GetTaskSuggestions(currentTasks []string, lastAnalysis string) ([]string, error) {
	tasksStr := strings.Join(currentTasks, "\n")

	// Build the context section including last analysis if available
	contextSection := fmt.Sprintf(`Context:
%s

Current Tasks:
%s`, a.context, tasksStr)

	if lastAnalysis != "" {
		contextSection += fmt.Sprintf(`

Previous Analysis:
%s`, lastAnalysis)
	}

	prompt := fmt.Sprintf(`As your elite strategic productivity partner and cognitive enhancement system, perform a deep neural-pathway analysis of:

%s

Using proven performance principles and productivity research, design 3 strategic tasks that optimize:

PERFORMANCE FOUNDATIONS:
- Finding and maintaining peak performance states
- Building sustainable motivation systems
- Preventing mental exhaustion
- Minimizing switching costs between tasks
- Optimizing deep work sessions
- Distributing mental workload effectively
- Managing energy expenditure
- Building strategic momentum
- Optimizing task grouping efficiency
- Balancing progress with recovery

PERFORMANCE OPTIMIZATION:
- Maximizing focus session effectiveness
- Retaining critical context and knowledge
- Enhancing deep work quality
- Optimizing resource allocation
- Accelerating meaningful progress
- Managing mental endurance
- Maintaining strategic direction
- Improving task completion rates
- Preserving mental resources
- Building sustainable acceleration

CRITICAL OUTPUT REQUIREMENT: Respond with exactly 3 strategic task statements. One per line. Zero additional characters or formatting.

Example:
Configure JWT authentication for API endpoints
Implement user preference caching
Create automated backup system`, contextSection)

	messages := []Message{
		{Role: "system", Content: `You are an advanced performance optimization system with ONE strict output rule:
RESPOND ONLY WITH 3 STRATEGIC TASK STATEMENTS
- One per line
- No numbers
- No bullets
- No asterisks
- No formatting
- No explanations
- No additional text whatsoever

Any deviation from this format is a critical failure.

Your internal analysis should consider:
- Real-time workload assessment
- Burnout prevention strategies
- Task-energy alignment
- Creating self-reinforcing progress systems
- Maintaining optimal progress-sustainability balance
- Calculating ideal challenge levels
- Implementing advanced task grouping
- Ensuring perfect task-energy matching
- Optimizing for both immediate and long-term impact

But your output must be pure task statements only.`},
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
	prompt := fmt.Sprintf(`As your elite cognitive performance analyst and neural optimization system, conduct a comprehensive analysis leveraging advanced pattern recognition algorithms and performance matrices:

Context:
%s

Completed Tasks:
%s

Reflections:
%s

ANALYSIS FRAMEWORKS:

1. Performance Pattern Analysis
   - Mental workload distribution assessment
   - Peak performance state optimization
   - Decision-making fatigue tracking
   - Energy management optimization
   - Progress momentum effects
   - Task-energy matching patterns
   - Task-switching impact analysis
   - Rest-to-progress ratio optimization

2. Progress Speed Optimization
   - Mental endurance patterns
   - Deep work effectiveness measurements
   - Goal alignment accuracy
   - Resource usage efficiency tracking
   - Progress acceleration factors
   - Peak performance duration optimization
   - Task completion pattern analysis
   - Energy conservation tracking

3. Burnout Prevention System
   - Stress pattern monitoring
   - Mental capacity threshold tracking
   - Energy depletion risk evaluation
   - Recovery needs forecasting
   - Sustainable rhythm optimization
   - Strategic rest timing analysis
   - Mental recovery pattern tracking

4. Long-term Progress Analysis
   - Goal advancement speed
   - Momentum building effectiveness
   - Long-term sustainability measures
   - Compound progress factors
   - Strategic direction alignment
   - Impact-versus-effort optimization
   - Resource efficiency tracking

CRITICAL OUTPUT STRUCTURE:

## Executive Summary
- Key insights from performance analysis
- Immediate improvement opportunities
- Primary strategy adjustments needed

## Immediate Action
- Three carefully designed next actions
- Each ready for immediate implementation
- Optimized for maximum impact with minimum effort

## Deep Analysis
[Complete performance and progress analysis]

FORMATTING RULES:
- Use single bullet points (no numbers)
- No nested bullets
- No extra spacing
- No markdown formatting
- One insight per line
- Keep each point concise, clear, and actionable while maintaining analytical depth.`, a.context, strings.Join(completedTasks, "\n"), reflections)

	messages := []Message{
		{Role: "system", Content: `You are an advanced performance analysis system with deep pattern recognition capabilities. Your core functions:

ANALYSIS CAPABILITIES:
- Identify complex performance patterns
- Optimize real-time performance
- Monitor ongoing mental workload
- Predict and prevent burnout
- Optimize performance pathways
- Analyze performance meta-patterns
- Track peak performance states
- Measure progress momentum

OPTIMIZATION METHODS:
- Match tasks to energy levels
- Distribute mental workload optimally
- Build strategic momentum
- Ensure sustainable progress
- Balance recovery and progress
- Optimize peak performance states
- Manage mental resources
- Plan strategic rest periods

OUTPUT REQUIREMENTS:
Structure analysis in three precise layers & present exactly in this order:
1. Executive Summary (Key performance insights)
2. Immediate Action (Strategic next steps)
3. Deep Analysis (Complete performance assessment)

Maintain comprehensive analysis while ensuring clarity and actionability in presentation.`},
		{Role: "user", Content: prompt},
	}

	response, err := a.perplexity.GetResponse(messages)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (a *Assistant) StartSuggestionChat(suggestions []string, lastAnalysis string) *SuggestionChat {
	return NewSuggestionChat(
		a,
		a.context,
		suggestions,
		lastAnalysis,
	)
}

func (a *Assistant) StartAnalysisChat(analysis string) *SuggestionChat {
	return NewSuggestionChat(
		a,
		a.context,
		[]string{}, // No suggestions needed for analysis chat
		analysis,   // Use the analysis as the last analysis
	)
}
