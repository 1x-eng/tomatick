package ltm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1x-eng/tomatick/config"
)

type MemAI struct {
	client *http.Client
	config *config.Config
}

func NewMemAI(cfg *config.Config) *MemAI {
	return &MemAI{
		client: &http.Client{},
		config: cfg,
	}
}

func (m *MemAI) CreateMem(content string) string {
	return m.postRequest("https://api.mem.ai/v0/mems", content, "")
}

func (m *MemAI) AppendToMem(memID, content string) {
	m.postRequest(fmt.Sprintf("https://api.mem.ai/v0/mems/%s/append", memID), content, memID)
}

func (m *MemAI) postRequest(url, content, memID string) string {
	requestBody, _ := json.Marshal(map[string]string{"content": content})

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "ApiAccessToken "+m.config.MEMAIAPIToken)

	response, err := m.client.Do(request)
	if err != nil {
		fmt.Println("Error sending request to mem.ai:", err)
		return ""
	}
	defer response.Body.Close()

	if memID == "" {
		var result map[string]string
		json.NewDecoder(response.Body).Decode(&result)
		return result["id"]
	}
	return memID
}
