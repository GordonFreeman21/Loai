package backend

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OllamaClient struct {
	BaseURL string
}

type ModelListResponse struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name string `json:"name"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponseChunk struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

func NewOllamaClient(baseURL string) *OllamaClient {
	return &OllamaClient{BaseURL: baseURL}
}

func (c *OllamaClient) GetModels() ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/tags", c.BaseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ModelListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range result.Models {
		models = append(models, m.Name)
	}
	return models, nil
}

func (c *OllamaClient) StreamChat(ctx context.Context, model string, messages []Message) error {
	reqBody, err := json.Marshal(ChatRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/chat", c.BaseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama error: %s", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk ChatResponseChunk
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}

		if chunk.Message.Content != "" {
			runtime.EventsEmit(ctx, "ollama_chunk", chunk.Message.Content)
		}

		if chunk.Done {
			runtime.EventsEmit(ctx, "ollama_done", nil)
			break
		}
	}

	return scanner.Err()
}

func (c *OllamaClient) PullModel(ctx context.Context, model string) error {
	reqBody, _ := json.Marshal(map[string]string{"name": model})
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/pull", c.BaseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		runtime.EventsEmit(ctx, "pull_progress", scanner.Text())
	}
	runtime.EventsEmit(ctx, "pull_done", nil)
	return nil
}
