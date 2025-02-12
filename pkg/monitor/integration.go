package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/1x-eng/tomatick/config"
	"github.com/1x-eng/tomatick/pkg/llm"
)

// TomatickMonitor provides a high-level interface for monitoring Tomatick breaks
type TomatickMonitor struct {
	activityMonitor  *ActivityMonitor
	notificationMgr  *NotificationManager
	config           *config.Config
	lastNotification time.Time
	notifyThreshold  time.Duration
}

// NewTomatickMonitor creates a new TomatickMonitor instance
func NewTomatickMonitor(cfg *config.Config, llmClient *llm.PerplexityAI) (*TomatickMonitor, error) {
	if err := InitializeMonitoring(cfg); err != nil {
		return nil, fmt.Errorf("failed to initialize monitoring: %w", err)
	}

	return &TomatickMonitor{
		activityMonitor: NewActivityMonitor(cfg, llmClient),
		notificationMgr: NewNotificationManager(llmClient, cfg.UserName),
		config:          cfg,
		notifyThreshold: 60 * time.Second, // Increase threshold to 60 seconds between notifications
	}, nil
}

// OnBreakStart should be called when a Tomatick break starts
func (tm *TomatickMonitor) OnBreakStart() {
	tm.activityMonitor.StartBreak()
	log.Println("Break monitoring started")
}

// OnBreakEnd should be called when a Tomatick break ends
func (tm *TomatickMonitor) OnBreakEnd() BreakSummary {
	violations := tm.activityMonitor.EndBreak()

	if len(violations) == 0 {
		return BreakSummary{
			HasViolations: false,
			Message:       "Great job taking a proper break!",
		}
	}

	// Return summary immediately without additional notification
	return BreakSummary{
		HasViolations: true,
		ViolationDetails: fmt.Sprintf("Break ended with %d violations:\n%s",
			len(violations),
			tm.activityMonitor.GetViolationSummary()),
	}
}

// CheckBreakViolations checks for break violations and returns a notification if needed
func (tm *TomatickMonitor) CheckBreakViolations() *string {
	if !tm.activityMonitor.IsOnBreak() {
		return nil
	}

	// Check if enough time has passed since last notification
	if time.Since(tm.lastNotification) < tm.notifyThreshold {
		return nil
	}

	violations := tm.activityMonitor.violations
	if len(violations) == 0 {
		return nil
	}

	// Only notify if there are new violations since last check
	lastViolationTime := violations[len(violations)-1].Timestamp
	if !lastViolationTime.After(tm.lastNotification) {
		return nil
	}

	violation := BreakViolation{
		Events:    violations,
		StartTime: violations[0].Timestamp,
		EndTime:   violations[len(violations)-1].Timestamp,
	}

	notification, err := tm.notificationMgr.GenerateNotification(violation)
	if err != nil {
		return nil
	}

	tm.lastNotification = time.Now()
	return &notification
}

// BreakSummary contains information about a completed break
type BreakSummary struct {
	HasViolations    bool
	Message          string
	ViolationDetails string
}
