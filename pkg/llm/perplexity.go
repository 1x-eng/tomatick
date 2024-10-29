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
		Model:    "llama-3.1-sonar-huge-128k-online",
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.config.PerplexityAPIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var perplexityResp PerplexityResponse
	err = json.Unmarshal(body, &perplexityResp)
	if err != nil {
		return "", err
	}

	if len(perplexityResp.Choices) > 0 {
		return perplexityResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from Perplexity AI")
}
