package daemon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"wraith/internal/config"
	"wraith/internal/memory"
)

// TestRule3_BudgetFallback verifies that when the LLM takes longer than the budget
// (or is down), the daemon falls back to Grimoire and tracks source as "grimoire"
// within the time budget.
func TestRule3_BudgetFallback(t *testing.T) {
	tempDir := t.TempDir()

	// Mock slow LLM that sleeps 500ms
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Too slow whisper"}`))
	}))
	defer slowServer.Close()

	cfg := &config.Config{
		Port:         23999,
		BudgetMs:     150, // 150ms hard budget
		OllamaURL:    slowServer.URL,
		Model:        "gemma3:4b",
		Persona:      "victorian_ghost",
		TriggerRate:  1.0,
		DataDir:      tempDir,
		DebugOverlay: true,
	}

	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	defer server.Shutdown(context.Background())

	reqBody := `{"shell":"pwsh","exit_code":1,"command":"go test","cwd":"/test","duration_ms":100,"budget_ms":150}`
	httpReq := httptest.NewRequest("POST", "/ingest", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	start := time.Now()
	server.handleIngest(rr, httpReq)
	elapsed := time.Since(start)

	// Ingestion must complete close to the budget (under 350ms including overhead)
	if elapsed > 400*time.Millisecond {
		t.Fatalf("Rule 3 violated: prompt hung longer than budget: %v", elapsed)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"source":"grimoire"`) {
		t.Fatalf("Expected fallback to source 'grimoire' due to budget overrun, got: %s", body)
	}

	// Verify overlay includes (grimoire)
	if !strings.Contains(body, "(grimoire)") {
		t.Fatalf("Expected debug overlay to indicate (grimoire), got: %s", body)
	}

	// Verify SQLite has recorded the source as grimoire
	recent, err := server.store.GetRecentSources(context.Background(), 1)
	if err != nil || len(recent) == 0 {
		t.Fatalf("Expected recent events in store: %v", err)
	}
	if recent[0].Source != memory.SourceGrimoire {
		t.Fatalf("Store recorded source as %q, want %q", recent[0].Source, memory.SourceGrimoire)
	}
}
