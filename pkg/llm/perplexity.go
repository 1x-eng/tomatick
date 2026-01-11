package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/1x-eng/tomatick/config"
)

type PerplexityAI struct {
	client *http.Client
	config *config.Config
}

type PerplexityRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PerplexityResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewPerplexityAI(cfg *config.Config) *PerplexityAI {
	return &PerplexityAI{
		client: &http.Client{},
		config: cfg,
	}
}

func (p *PerplexityAI) GetResponse(messages []Message) (string, error) {
	url := "https://api.perplexity.ai/chat/completions"

	reqBody := PerplexityRequest{
		Model:    "sonar-reasoning-pro",
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.config.PerplexityAPIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// If not 200, try to get error message
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var perplexityResp PerplexityResponse
	if err := json.Unmarshal(body, &perplexityResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w\nResponse body: %s", err, string(body))
	}

	if len(perplexityResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response: %s", string(body))
	}

	return perplexityResp.Choices[0].Message.Content, nil
}
