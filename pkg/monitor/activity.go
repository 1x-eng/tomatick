package monitor

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/1x-eng/tomatick/config"
	"github.com/1x-eng/tomatick/pkg/llm"
)

// ActivityEvent represents a user activity event
type ActivityEvent struct {
	Timestamp time.Time
	Type      ActivityType
	Details   string
}

// ActivityType represents different types of user activity
type ActivityType int

const (
	KeyboardActivity ActivityType = iota
	MouseActivity
	AppFocusChange
)

// ActivityMonitor handles monitoring user activity during breaks
type ActivityMonitor struct {
	mu             sync.RWMutex
	isBreak        bool
	lastActivity   time.Time
	violations     []ActivityEvent
	config         *config.Config
	llmClient      *llm.PerplexityAI
	stopChan       chan struct{}
	violationsChan chan ActivityEvent
	userName       string
}

// NewActivityMonitor creates a new activity monitor
func NewActivityMonitor(cfg *config.Config, llmClient *llm.PerplexityAI) *ActivityMonitor {
	return &ActivityMonitor{
		config:         cfg,
		llmClient:      llmClient,
		stopChan:       make(chan struct{}),
		violationsChan: make(chan ActivityEvent, 100),
		lastActivity:   time.Now(),
		userName:       cfg.UserName,
	}
}

// StartBreak signals the start of a break period
func (am *ActivityMonitor) StartBreak() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.isBreak = true
	am.violations = nil // Reset violations for new break
	go am.monitorActivity()
}

// EndBreak signals the end of a break period and returns any violations
func (am *ActivityMonitor) EndBreak() []ActivityEvent {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.isBreak = false
	am.stopChan <- struct{}{} // Signal monitoring goroutine to stop
	return am.violations
}

// IsOnBreak returns whether the user is currently on break
func (am *ActivityMonitor) IsOnBreak() bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.isBreak
}

// monitorActivity starts monitoring user activity
func (am *ActivityMonitor) monitorActivity() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-am.stopChan:
			return
		case event := <-am.violationsChan:
			am.mu.Lock()
			am.violations = append(am.violations, event)
			am.mu.Unlock()
		case <-ticker.C:
			if !am.IsOnBreak() {
				return
			}
			// Check for activity
			if activity := am.checkActivity(); activity != nil {
				am.violationsChan <- *activity
			}
		}
	}
}

// checkActivity checks for current user activity
func (am *ActivityMonitor) checkActivity() *ActivityEvent {
	// Get the frontmost app
	appName, err := getForegroundApp()
	if err != nil {
		log.Printf("Error getting foreground app: %v", err)
		return nil
	}

	// First check if there's continuous activity
	hasActivity := hasRecentActivity()
	isWorkRelated := isWorkApp(appName)

	// Only report violations if:
	// 1. A work-related app is active AND there's continuous activity, or
	// 2. There's continuous activity in any app that exceeds our threshold
	if isWorkRelated && hasActivity {
		return &ActivityEvent{
			Timestamp: time.Now(),
			Type:      AppFocusChange,
			Details:   fmt.Sprintf("Active work in %s detected during break", appName),
		}
	} else if hasActivity {
		return &ActivityEvent{
			Timestamp: time.Now(),
			Type:      KeyboardActivity,
			Details:   fmt.Sprintf("Sustained activity in %s during break", appName),
		}
	}

	return nil
}

// GetViolationSummary returns a summary of break violations
func (am *ActivityMonitor) GetViolationSummary() string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if len(am.violations) == 0 {
		return "No break violations detected."
	}

	summary := fmt.Sprintf("Detected %d break violations:\n", len(am.violations))
	for _, v := range am.violations {
		summary += fmt.Sprintf("- %s: %s\n", v.Timestamp.Format("15:04:05"), v.Details)
	}
	return summary
}
