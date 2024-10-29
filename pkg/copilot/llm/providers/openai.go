package providers

import (
	"context"
	"fmt"

	"github.com/1x-eng/tomatick/pkg/copilot/llm"
	"github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	config llm.ProviderConfig
}

func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{}
}

func (p *OpenAIProvider) Initialize(config llm.ProviderConfig) error {
	if config.APIKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}
	p.client = openai.NewClient(config.APIKey)
	p.config = config
	return nil
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt string, params llm.GenerationParams) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: p.config.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: params.Temperature,
		MaxTokens:   params.MaxTokens,
		TopP:        params.TopP,
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func (p *OpenAIProvider) Cleanup() error {
	return nil
}
