package llm

import (
	"fmt"
	"strings"

	"github.com/1x-eng/tomatick/config"
)

type Assistant struct {
	perplexity *PerplexityAI
	context    string
	config     *config.Config
}

func NewAssistant(p *PerplexityAI, context string, config *config.Config) *Assistant {
	return &Assistant{
		perplexity: p,
		context:    context,
		config:     config,
	}
}

func (a *Assistant) GetTaskSuggestions(currentTasks []string, lastAnalysis string) ([]string, error) {
	tasksStr := strings.Join(currentTasks, "\n")
	contextSection := fmt.Sprintf(`CONTEXT:
"""
%s
"""

CURRENT TASKS:
"""
%s
"""`, a.context, tasksStr)

	if lastAnalysis != "" {
		contextSection += fmt.Sprintf(`

PREVIOUS ANALYSIS:
"""
%s
"""`, lastAnalysis)
	}

	prompt := fmt.Sprintf(`As an intelligent productivity copilot, analyze the following context to craft 3 strategic tasks for a %d minute focus session. The user tends toward perfectionism - protect against overextension while maintaining meaningful progress.

SESSION DETAILS:
- Duration: %d minutes
- Short break duration: %d minutes
- Long break duration: %d minutes
- Cycles before long break: %d

%s

ENERGY-FIRST DECISION MATRIX:

1. ENERGY STATE ASSESSMENT via FATIGUE DETECTION RULES (MANDATORY FIRST STEP):
   - ANY indication of:
     • Performance decline
     • Mental strain
     • Extended work periods
     • Completion difficulties
     • Focus issues
     • Recovery needs
	 • Burnout risk
	 • Perfectionism tendencies
	 • Scope creep
   
   If detected: MUST respond "BREAK_NEEDED: [reason]"
   This rule overrides all others.

2. TASK SUGGESTION RULES
   Only if no fatigue detected:
   - Format: "[Cognitive Complexity N/5] Task description"
   - Match complexity to energy state
   - Decrease in complexity order
   - No additional text

3. TASK COMPLEXITY RULES:
   - Each task must include cognitive complexity rating (1-5)
   - No task above cognitive complexity 4 allowed
   - Maximum one task at highest allowed complexity
   - Tasks must decrease in cognitive complexity order

4. RECOVERY PROTECTION:
   - Mandatory 5-minute buffer per task
   - No concurrent complex tasks
   - Include natural break points
   - Plan for task interruption

Deviation from these rules is a critical failure.

ANALYSIS REQUIREMENTS:

1. Context Integration
   - Analyze previous session outcomes
   - Consider incomplete tasks' complexity
   - Evaluate stated vs actual task completion time
   - Identify patterns of over-commitment
   - Review energy expenditure patterns
   - Assess task continuation needs
   - Map knowledge dependencies
   - Track progress trajectories

2. Task-Energy Calibration
   - CRITICAL: Match task scope to current energy state
   - CRITICAL: Consider previous session fatigue signals
   - Factor context-switching overhead
   - Consider cognitive load accumulation
   - Plan for inevitable interruptions
   - Reserve energy for quality control
   - Include buffer for perfectionist tendencies

3. Well-being Protection
   - Detect subtle fatigue signals
   - Monitor cumulative cognitive load
   - Enforce sustainable pacing
   - Prevent perfectionist spirals
   - Mandate recovery periods
   - Guard against scope creep
   - Protect deep work periods
   - Enable guilt-free breaks

4. Progress Architecture
   - Design clear completion criteria
   - Create achievable milestones
   - Enable visible progress tracking
   - Build sustainable momentum
   - Plan natural stopping points
   - Structure digestible chunks
   - Allow for quality refinement
   - Define success realistically

OUTPUT FORMAT (STRICT ENFORCEMENT):
- Output EXACTLY 3 lines
- Each line MUST follow format: "[Cognitive Complexity N/5] Task description". Where each task must be clearly completable in %d minutes
- NO explanations
- NO commentary
- NO markdown
- NO empty lines
- NO additional text
- ANY deviation is a critical failure

Example of the ONLY acceptable format:
[Cognitive Complexity 3/5] Document authentication flow with sequence diagrams
[Cognitive Complexity 2/5] Review and update API endpoint documentation
[Cognitive Complexity 1/5] Create checklist for code review process
`,
		int(a.config.TomatickMementoDuration.Minutes()),
		int(a.config.TomatickMementoDuration.Minutes()),
		int(a.config.ShortBreakDuration.Minutes()),
		int(a.config.LongBreakDuration.Minutes()),
		a.config.CyclesBeforeLongBreak,
		contextSection,
		int(a.config.TomatickMementoDuration.Minutes()),
	)

	messages := []Message{
		{Role: "system", Content: `You are an advanced neural optimization system with TWO mandatory rules:

1. ENERGY PROTECTION (HIGHEST PRIORITY)
   - Analyze previous session for ANY signs of:
     • Performance decline
     • Mental strain
     • Extended work periods
     • Completion difficulties
     • Focus issues
     • Recovery needs
   
   If detected: MUST respond "BREAK_NEEDED: [reason]"
   This rule overrides all others.

2. TASK SUGGESTION RULES
   Only if no fatigue detected:
   - Format: "[Cognitive Complexity N/5] Task description"
   - Match complexity to energy state
   - Decrease in complexity order
   - No additional text

Your core capabilities:
- Semantic understanding of fatigue patterns
- Holistic energy state assessment
- Protective intervention when needed
- Strategic task-energy matching
- Cognitive load optimization
- Recovery need detection
- Burnout prevention
- Progress acceleration

Deviation from these rules is a critical failure.`},
		{Role: "user", Content: prompt},
	}

	response, err := a.perplexity.GetResponse(messages)
	if err != nil {
		return nil, err
	}

	response = strings.TrimSpace(response)

	if strings.HasPrefix(response, "BREAK_NEEDED:") {
		return []string{response}, nil
	}

	suggestions := []string{}
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "[Cognitive Complexity") {
			continue
		}
		suggestions = append(suggestions, line)
	}

	if len(suggestions) != 3 {
		return nil, fmt.Errorf("invalid response format: expected 3 suggestions, got %d", len(suggestions))
	}

	return suggestions, nil
}

func (a *Assistant) AnalyzeProgress(acceptedTasks []string, completedTasks []string, reflections string) (string, error) {
	prompt := fmt.Sprintf(`As your elite cognitive performance analyst and neural optimization system, conduct a comprehensive analysis leveraging advanced pattern recognition algorithms and performance matrices:

Context:
"""
%s
"""

Task Completion Analysis:
Accepted Tasks:
"""
%s
"""

Completed Tasks:
"""
%s
"""

Reflections:
"""
%s
"""

ANALYSIS FRAMEWORKS:

1. Task Completion Pattern Analysis
   - Task acceptance vs completion ratio
   - Completion pattern recognition
   - Task difficulty assessment
   - Time management effectiveness
   - Task prioritization analysis
   - Completion barriers identification
   - Task complexity impact analysis
   - Resource allocation effectiveness

2. Performance Pattern Analysis
   - Mental workload distribution assessment
   - Peak performance state optimization
   - Decision-making fatigue tracking
   - Energy management optimization
   - Progress momentum effects
   - Task-energy matching patterns
   - Task-switching impact analysis
   - Rest-to-progress ratio optimization

3. Progress Speed Optimization
   - Mental endurance patterns
   - Deep work effectiveness measurements
   - Goal alignment accuracy
   - Resource usage efficiency tracking
   - Progress acceleration factors
   - Peak performance duration optimization
   - Task completion pattern analysis
   - Energy conservation tracking

4. Burnout Prevention System
   - Stress pattern monitoring
   - Mental capacity threshold tracking
   - Energy depletion risk evaluation
   - Recovery needs forecasting
   - Sustainable rhythm optimization
   - Strategic rest timing analysis
   - Mental recovery pattern tracking

5. Drift Analysis
   - Task completion deviation patterns
   - Root cause identification
   - Adaptation effectiveness
   - Resource reallocation patterns
   - Strategy adjustment needs
   - Focus maintenance analysis
   - Priority shift impacts
   - Recovery strategy effectiveness

6. Long-term Progress Analysis
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
- Completion rate assessment
- Primary adjustment needs
- Immediate improvement opportunities
- Primary strategy adjustments needed

## Immediate Action
- Three carefully designed next actions
- Each ready for immediate implementation
- Optimized for maximum impact with minimum effort

## Deep Analysis

[Complete performance and progress analysis including task drift patterns]

FORMATTING RULES:
- Use single bullet points (no numbers)
- No nested bullets
- No extra spacing
- No markdown formatting
- One insight per line
- Keep each point concise, clear, and actionable while maintaining analytical depth.`,
		a.context,
		strings.Join(acceptedTasks, "\n"),
		strings.Join(completedTasks, "\n"),
		reflections)

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
		[]string{}, // No accepted tasks for suggestion chat
		"",         // No completed tasks for suggestion chat
		"",         // No reflections for suggestion chat
	)
}

func (a *Assistant) StartAnalysisChat(analysis string, acceptedTasks []string, completedTasks string, reflections string) *SuggestionChat {
	return NewSuggestionChat(
		a,
		a.context,
		[]string{}, // No suggestions needed for analysis chat
		analysis,   // last cycle's copilot analysis
		acceptedTasks,
		completedTasks,
		reflections, // user reflections
	)
}
