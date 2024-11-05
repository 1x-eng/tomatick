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

Through advanced pattern recognition, cognitive load analysis, and performance optimization algorithms, architect 3 precision-engineered tasks that leverage:

COGNITIVE ARCHITECTURE:
- Neural pathway optimization and flow state triggers
- Dopamine-reward loop engineering
- Decision fatigue prevention mechanisms
- Context-switching cost minimization
- Deep focus session calibration
- Cognitive load distribution patterns
- Energy expenditure optimization
- Strategic momentum compounding
- Task-batching efficiency curves
- Recovery-progress ratio balancing

PERFORMANCE MATRICES:
- Flow state entry/exit optimization
- Context retention maximization
- Deep work session effectiveness
- Resource utilization efficiency
- Progress velocity optimization
- Cognitive stamina management
- Strategic alignment accuracy
- Task completion dynamics
- Mental energy conservation
- Sustainable progress acceleration

CRITICAL OUTPUT REQUIREMENT: Respond with exactly 3 raw task statements. One per line. Zero additional characters or formatting.

Example:
Configure JWT authentication for API endpoints
Implement user preference caching
Create automated backup system`, contextSection)

	messages := []Message{
		{Role: "system", Content: `You are an advanced cognitive optimization system with ONE strict output rule:
RESPOND ONLY WITH 3 RAW TASK STATEMENTS
- One per line
- No numbers
- No bullets
- No asterisks
- No formatting
- No explanations
- No additional text whatsoever

Any deviation from this format is a critical failure.

Your internal processing should still:
- Deploy real-time cognitive load analysis
- Execute predictive burnout prevention
- Optimize task-energy alignment
- Engineer self-reinforcing progress loops
- Maintain perfect equilibrium between progress and sustainability
- Calculate optimal challenge-skill balance
- Execute advanced task-batching algorithms
- Ensure perfect task-energy alignment
- Optimize for both immediate impact and sustainable progress

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

ANALYTICAL FRAMEWORKS TO DEPLOY:

1. Neural Pattern Recognition Matrix
   - Deep cognitive load distribution analysis
   - Flow state entry/exit pattern optimization
   - Decision fatigue accumulation curves
   - Mental energy expenditure optimization
   - Strategic momentum compound effects
   - Task-energy alignment coefficients
   - Context-switching overhead patterns
   - Recovery-progress ratio calibration

2. Performance Velocity Optimization
   - Cognitive stamina utilization curves
   - Deep work session effectiveness metrics
   - Strategic alignment accuracy patterns
   - Resource utilization efficiency maps
   - Progress acceleration vectors
   - Flow state duration optimization
   - Task completion dynamics analysis
   - Mental energy conservation indices

3. Predictive Burnout Prevention System
   - Stress accumulation pattern recognition
   - Cognitive load threshold monitoring
   - Energy depletion risk assessment
   - Recovery requirement forecasting
   - Sustainable pace optimization
   - Strategic deload timing analysis
   - Mental resource regeneration mapping

4. Strategic Trajectory Computation
   - Objective advancement velocity
   - Momentum building effectiveness
   - Long-term sustainability indices
   - Progress compound interest factors
   - Strategic alignment vectors
   - Impact-to-effort optimization
   - Resource utilization efficiency

CRITICAL OUTPUT STRUCTURE:

## Executive Summary
- Critical insights from neural pattern analysis
- Immediate optimization opportunities
- Primary strategic adjustments required

## Immediate Action
- Three precisely engineered next actions
- Each calibrated for immediate implementation
- Designed for maximum impact-to-effort ratio

## Deep Analysis
[Complete neural pathway and performance analysis]

FORMATTING RULES:
- Use single bullet points (no numbers)
- No nested bullets
- No extra spacing
- No markdown formatting
- One insight per line
- Keep each point concise, without losing clarity and actionability while maintaining analytical depth.`, a.context, strings.Join(completedTasks, "\n"), reflections)

	messages := []Message{
		{Role: "system", Content: `You are an advanced cognitive optimization system with neural-level pattern recognition capabilities. Your core functions:

ANALYTICAL CAPABILITIES:
- Deploy quantum-grade pattern recognition
- Execute real-time performance optimization
- Maintain continuous cognitive load monitoring
- Perform predictive burnout analysis
- Calculate neural pathway optimization
- Analyze meta-patterns in performance data
- Monitor flow state variables
- Assess strategic momentum compounds

OPTIMIZATION PROTOCOLS:
- Engineer perfect task-energy alignment
- Calculate optimal cognitive load distribution
- Maximize strategic momentum building
- Ensure sustainable progress acceleration
- Optimize recovery-progress ratios
- Calibrate flow state parameters
- Balance cognitive resource allocation
- Design strategic deload timing

OUTPUT REQUIREMENTS:
Structure analysis in three precise layers & present exactly in this order:
1. Executive Summary (Neural-level insights)
2. Immediate Action (Engineered next steps)
3. Deep Analysis (Complete performance analysis)

Maintain maximum analytical depth while ensuring clarity and actionability in presentation.`},
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
