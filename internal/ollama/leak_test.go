package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wraith/internal/redact"
	"wraith/prompts"
)

// TestRule6_ZeroLeakPromptCapture verifies that when sensitive inputs are redacted,
// NO sensitive values reach the prompt sent to the LLM backend.
func TestRule6_ZeroLeakPromptCapture(t *testing.T) {
	capturedPrompts := make([]string, 0)

	// Mock Ollama server recording all incoming prompt bodies
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req GenerateRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		capturedPrompts = append(capturedPrompts, req.Prompt)

		resp := GenerateResponse{
			Model:    "gemma3:4b",
			Response: "The shadows notice your presence.",
			Done:     true,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	sensitiveSecret := "ghp_SuperSecretGitHubToken998877665544"
	sensitivePass := "TopSecretDatabasePass!123"

	rawCommand := "git push --force https://user:" + sensitiveSecret + "@github.com/repo.git --password " + sensitivePass

	// Redaction pipeline
	redacted := redact.Redact(rawCommand)

	promptData := model.PromptData{
		PersonaPrompt:   "You are a ghost.",
		Shell:           "pwsh",
		ExitCode:        0,
		Hour:            2,
		Cwd:             "C:/work",
		MemorySummary:   "",
		RedactedCommand: redacted.RedactedText,
		DurationMs:      50,
	}

	prompt, err := model.BuildPrompt(promptData)
	if err != nil {
		t.Fatalf("BuildPrompt failed: %v", err)
	}

	client := NewClient(mockServer.URL)
	_, err = client.Generate(context.Background(), "gemma3:4b", prompt, 42)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(capturedPrompts) == 0 {
		t.Fatal("Expected at least one prompt captured by mock Ollama server")
	}

	sentPrompt := capturedPrompts[0]

	// Assert zero secret leakage
	if strings.Contains(sentPrompt, sensitiveSecret) {
		t.Fatalf("LEAK DETECTED: prompt sent to Ollama contained GitHub secret token!")
	}
	if strings.Contains(sentPrompt, sensitivePass) {
		t.Fatalf("LEAK DETECTED: prompt sent to Ollama contained database password!")
	}
}
