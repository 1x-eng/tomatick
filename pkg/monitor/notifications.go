package monitor

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/1x-eng/tomatick/pkg/llm"
	"github.com/charmbracelet/lipgloss"
)

// NotificationType represents different types of break violation notifications
type NotificationType int

const (
	GentleReminder NotificationType = iota
	ConcernedNotice
	SupportiveIntervention
)

// NotificationStyle holds the styling for break notifications
var NotificationStyle = lipgloss.NewStyle().
	Width(50).
	Align(lipgloss.Left)

// displayNotification shows a clean, inline notification
func displayNotification(message string) {
	border := strings.Repeat("═", 50)

	fmt.Printf("\n%s\n", border)
	fmt.Printf("🔔  Break Time\n")
	fmt.Printf("%s\n\n", strings.Repeat("─", 50))

	// Split and print message with proper padding
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			fmt.Println()
		} else {
			fmt.Printf("  %s\n", strings.TrimSpace(line))
		}
	}

	fmt.Printf("%s\n\n", border)
}

// BreakViolation represents a collection of activity events during a break
type BreakViolation struct {
	Events    []ActivityEvent
	StartTime time.Time
	EndTime   time.Time
}

// NotificationManager handles generating appropriate notifications for break violations
type NotificationManager struct {
	llmClient           *llm.PerplexityAI
	userName            string
	breakViolationCount int
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(llmClient *llm.PerplexityAI, userName string) *NotificationManager {
	return &NotificationManager{
		llmClient:           llmClient,
		userName:            userName,
		breakViolationCount: 0,
	}
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

// createViolationContext creates a context string for the LLM based on the violation
func createViolationContext(violation BreakViolation) string {
	var details []string
	apps := make(map[string]bool)
	hasKeyboardActivity := false
	activityDuration := time.Since(violation.StartTime)

	for _, event := range violation.Events {
		if event.Type == AppFocusChange {
			appName := strings.TrimPrefix(event.Details, "Work-related app detected: ")
			apps[appName] = true
		}
		if event.Type == KeyboardActivity {
			hasKeyboardActivity = true
		}
	}

	// Add work apps used during break
	if len(apps) > 0 {
		appList := make([]string, 0, len(apps))
		for app := range apps {
			appList = append(appList, app)
		}
		details = append(details, fmt.Sprintf("Work apps used: %s", strings.Join(appList, ", ")))
	}

	// Add keyboard/mouse activity if detected
	if hasKeyboardActivity {
		details = append(details, "Active keyboard/mouse usage detected")
	}

	// Add activity duration
	if len(violation.Events) > 0 {
		details = append(details, fmt.Sprintf("Activity duration: %.1f minutes", activityDuration.Minutes()))
	}

	return strings.Join(details, "\n")
}

// GenerateNotification creates a contextual, supportive notification based on break violations
func (nm *NotificationManager) GenerateNotification(violation BreakViolation) (string, error) {
	nm.breakViolationCount++
	context := createViolationContext(violation)

	messages := []llm.Message{
		{
			Role: "system",
			Content: fmt.Sprintf(`You are a supportive productivity assistant that helps %s maintain healthy work-life balance.
Your role is to provide gentle, understanding reminders about the importance of taking proper breaks.
Always maintain a supportive, non-judgmental tone. Focus on well-being and long-term productivity.
Keep responses concise (max 2 sentences) and encouraging.
When addressing the user, use their name: %s`, nm.userName, nm.userName),
		},
		{
			Role: "user",
			Content: fmt.Sprintf(`Generate a gentle reminder for %s who has been working during their break:

Context:
%s

Requirements:
- Address them by name (%s)
- Be supportive and understanding
- Avoid any negative or judgmental language
- Keep it brief (max 2 sentences)
- Focus on well-being and long-term productivity
- Suggest a simple action they can take right now`, nm.userName, context, nm.userName),
		},
	}

	response, err := nm.llmClient.GetResponse(messages)
	if err != nil {
		return getDefaultNotification(violation), nil
	}

	cleaned := cleanResponse(response)
	message := fmt.Sprintf("%s\n\nViolation count: %d", cleaned, nm.breakViolationCount)

	displayNotification(message)
	return "", nil
}

// getDefaultNotification returns a data-driven default notification if LLM fails
func getDefaultNotification(violation BreakViolation) string {
	context := createViolationContext(violation)
	return fmt.Sprintf("%s\n\nConsistent breaks are essential for sustained productivity.", context)
}
