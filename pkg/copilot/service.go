package copilot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/1x-eng/tomatick/config"
	"github.com/1x-eng/tomatick/pkg/copilot/llm"
	"github.com/1x-eng/tomatick/pkg/copilot/prompts"
	"github.com/google/uuid"
)

type CopilotService struct {
	config        *config.CopilotConfig
	provider      llm.Provider
	mode          Mode
	suggestions   []Suggestion
	promptBuilder *prompts.PromptBuilder
	mu            sync.RWMutex
}

func NewCopilotService(cfg *config.CopilotConfig) (Service, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	mode := FollowMode
	if cfg.DefaultMode == "lead" {
		mode = LeadMode
	}

	provider, err := llm.NewProvider(cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize LLM provider: %w", err)
	}

	return &CopilotService{
		config:        cfg,
		provider:      provider,
		mode:          mode,
		promptBuilder: prompts.NewPromptBuilder(mode),
		suggestions:   make([]Suggestion, 0),
	}, nil
}

func (cs *CopilotService) Initialize(ctx context.Context) error {
	return cs.provider.Initialize(llm.ProviderConfig{
		Name:       cs.config.LLM.Provider,
		Model:      cs.config.LLM.Model,
		APIKey:     cs.config.LLM.APIKey,
		Parameters: cs.config.LLM.Parameters,
	})
}

func (cs *CopilotService) SetMode(mode Mode) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.mode = mode
	cs.promptBuilder = prompts.NewPromptBuilder(mode)
}

func (cs *CopilotService) AnalyzeContext(ctx context.Context, content string) error {
	prompt := cs.promptBuilder.BuildContextAnalysisPrompt(content, nil)

	response, err := cs.provider.Generate(ctx, prompt, llm.GenerationParams{
		Temperature:   0.7,
		MaxTokens:     150,
		TopP:          1.0,
		StopSequences: []string{"###"},
	})

	if err != nil {
		return fmt.Errorf("context analysis failed: %w", err)
	}

	suggestion := Suggestion{
		ID:        uuid.New().String(),
		Content:   response,
		Type:      Insight,
		CreatedAt: time.Now(),
		Priority:  1,
	}

	cs.mu.Lock()
	cs.suggestions = append(cs.suggestions, suggestion)
	cs.mu.Unlock()

	return nil
}

func (cs *CopilotService) GetSuggestions() []Suggestion {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.suggestions
}

func (cs *CopilotService) ProcessInput(ctx context.Context, input string) ([]Suggestion, error) {
	_ = cs.promptBuilder.GetSystemPrompt()

	response, err := cs.provider.Generate(ctx, input, llm.GenerationParams{
		Temperature:   0.7,
		MaxTokens:     150,
		TopP:          1.0,
		StopSequences: []string{"###"},
	})

	if err != nil {
		return nil, fmt.Errorf("input processing failed: %w", err)
	}

	suggestion := Suggestion{
		ID:        uuid.New().String(),
		Content:   response,
		Type:      TaskSuggestion,
		CreatedAt: time.Now(),
		Priority:  1,
	}

	cs.mu.Lock()
	cs.suggestions = append(cs.suggestions, suggestion)
	cs.mu.Unlock()

	return []Suggestion{suggestion}, nil
}

func (cs *CopilotService) Cleanup() error {
	if cs.provider != nil {
		return cs.provider.Cleanup()
	}
	return nil
}
