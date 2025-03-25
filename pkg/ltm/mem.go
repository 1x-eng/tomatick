package ltm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MemAI struct {
	client *http.Client
	config interface {
		GetMemAIToken() string
	}
}

type MemResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// NewMemAI creates a new MemAI client
func NewMemAI(cfg interface{ GetMemAIToken() string }) *MemAI {
	return &MemAI{
		client: &http.Client{},
		config: cfg,
	}
}

// CreateMem implements LongTermMemory interface
func (m *MemAI) CreateMem(content string) (string, error) {
	return m.postRequest("https://api.mem.ai/v0/mems", content, "")
}

// AppendToMem implements LongTermMemory interface
func (m *MemAI) AppendToMem(memID, content string) (string, error) {
	return m.postRequest(fmt.Sprintf("https://api.mem.ai/v0/mems/%s/append", memID), content, memID)
}

func (m *MemAI) postRequest(url, content, memID string) (string, error) {
	var memResponse MemResponse

	requestBody, _ := json.Marshal(map[string]string{"content": content})

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "ApiAccessToken "+m.config.GetMemAIToken())

	response, err := m.client.Do(request)
	if err != nil {
		fmt.Println("Error sending request to mem.ai:", err)
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&memResponse); err != nil {
		return "", err
	}

	if memID == "" {
		return memResponse.ID, nil
	}
	return memID, nil
}
