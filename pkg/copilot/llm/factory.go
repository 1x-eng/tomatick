package llm

import (
	"fmt"

	"github.com/1x-eng/tomatick/config"
	"github.com/1x-eng/tomatick/pkg/copilot/llm/providers"
)

func NewProvider(config config.LLMConfig) (Provider, error) {
	switch config.Provider {
	case "openai":
		return providers.NewOpenAIProvider(), nil
	case "anthropic":
		return nil, fmt.Errorf("anthropic provider not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider: %s", config.Provider)
	}
}
