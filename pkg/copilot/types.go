package copilot

import (
	"context"
	"time"
)

// Mode represents the copilot's operational mode
type Mode uint8

const (
	FollowMode Mode = iota
	LeadMode
)

// SuggestionType represents different types of copilot interactions
type SuggestionType uint8

const (
	TaskSuggestion SuggestionType = iota
	Insight
	Warning
	Clarification
)

// Suggestion represents a single copilot suggestion
type Suggestion struct {
	ID        string         // Unique identifier
	Content   string         // The actual suggestion
	Type      SuggestionType // Type of suggestion
	CreatedAt time.Time
	Priority  uint8 // Suggestion priority
}

// Service defines the main copilot interface
type Service interface {
	// Initialize sets up the copilot service
	Initialize(ctx context.Context) error

	// SetMode changes the copilot's operational mode
	SetMode(mode Mode)

	// AnalyzeContext processes the current context
	AnalyzeContext(ctx context.Context, content string) error

	// GetSuggestions returns current suggestions
	GetSuggestions() []Suggestion

	// ProcessInput handles user input and returns relevant suggestions
	ProcessInput(ctx context.Context, input string) ([]Suggestion, error)

	// Cleanup performs necessary cleanup
	Cleanup() error
}
