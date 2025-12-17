package webhook

import (
	"time"
)

// EventType represents the type of event occurring in the system
type EventType string

const (
	EventWorkStart    EventType = "work_start"
	EventWorkComplete EventType = "work_complete"
	EventBreakStart   EventType = "break_start"
	EventBreakEnd     EventType = "break_end"

	// AI & Intelligence Events
	EventContextRefined EventType = "context_refined"
	EventAISuggestions  EventType = "ai_suggestions"
	EventAIAnalysis     EventType = "ai_analysis"
	EventAIChatExchange EventType = "ai_chat_exchange"

	// Lifecycle Events
	EventSessionSummary EventType = "session_summary"
)

// EventPayload represents the standard structure sent to webhooks
type EventPayload struct {
	Type      EventType         `json:"type"`
	Timestamp time.Time         `json:"timestamp"`
	Data      map[string]string `json:"data,omitempty"`
}

// Dispatcher defines the interface for sending events
type Dispatcher interface {
	Dispatch(eventType EventType, data map[string]string)
	Wait()
}
