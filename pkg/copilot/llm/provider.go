package llm

import (
	"context"
)

type Provider interface {
	Initialize(config ProviderConfig) error
	Generate(ctx context.Context, prompt string, params GenerationParams) (string, error)
	Cleanup() error
}

type ProviderConfig struct {
	Name       string
	Model      string
	APIKey     string
	Parameters map[string]interface{}
}

type GenerationParams struct {
	Temperature   float64
	MaxTokens     int
	TopP          float64
	StopSequences []string
}
