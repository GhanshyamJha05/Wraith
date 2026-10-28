package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client interacts with the local Ollama instance
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// GenerateRequest represents Ollama /api/generate request
type GenerateRequest struct {
	Model     string                 `json:"model"`
	Prompt    string                 `json:"prompt"`
	Stream    bool                   `json:"stream"`
	KeepAlive string                 `json:"keep_alive"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

// GenerateResponse represents Ollama /api/generate response
type GenerateResponse struct {
	Model              string `json:"model"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

// NewClient returns a new Ollama client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:11434"
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Ping checks if Ollama daemon is reachable
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama unreachable at %s: %w\nRun: 'ollama serve' to start it", c.BaseURL, err)
	}
	defer resp.Body.Close()
	return nil
}

// Generate sends prompt to Ollama with keep_alive and token limits
func (c *Client) Generate(ctx context.Context, model string, prompt string, seed int) (*GenerateResponse, error) {
	reqBody := GenerateRequest{
		Model:     model,
		Prompt:    prompt,
		Stream:    false,
		KeepAlive: "10m", // Keep model warm in memory
		Options: map[string]interface{}{
			"num_predict": 60,   // Max tokens
			"temperature": 0.65, // Low-ish for consistency
		},
	}
	if seed != 0 {
		reqBody.Options["seed"] = seed
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/generate", c.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama generate failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("model %q not found in Ollama (HTTP 404). Run: 'ollama pull %s'", model, model)
		}
		return nil, fmt.Errorf("ollama returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return nil, fmt.Errorf("decoding Ollama response: %w", err)
	}

	return &genResp, nil
}
