package trigger

import (
	"context"
	"path/filepath"
	"testing"
	"time"
	"wraith/internal/memory"
)

func TestTriggerRules(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_trigger.db")
	store, err := memory.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}
	defer store.Close()

	engine := NewEngine(store, 0.0) // 0 base chance to test pure deterministic rules

	// 1. Dangerous command test
	dangerCmds := []string{
		"rm -rf /tmp/data",
		"git push origin main --force",
		"DROP TABLE users;",
		"chmod 777 script.sh",
	}

	for _, cmd := range dangerCmds {
		dec, err := engine.Evaluate(context.Background(), cmd, 0, "", false)
		if err != nil {
			t.Fatalf("Evaluate failed: %v", err)
		}
		if !dec.ShouldSpeak || dec.Trigger != TriggerDanger {
			t.Errorf("Expected danger trigger for %q, got: %+v", cmd, dec)
		}
	}

	// 2. Failure streak test
	failCmd := "go build ./..."
	// Insert two past failures into store
	_ = store.RecordEvent(context.Background(), memory.ShellEvent{
		Timestamp:       time.Now().Add(-2 * time.Minute),
		Shell:           "pwsh",
		ExitCode:        1,
		RedactedCommand: failCmd,
		Source:          memory.SourceModel,
	})
	_ = store.RecordEvent(context.Background(), memory.ShellEvent{
		Timestamp:       time.Now().Add(-1 * time.Minute),
		Shell:           "pwsh",
		ExitCode:        1,
		RedactedCommand: failCmd,
		Source:          memory.SourceModel,
	})

	dec, err := engine.Evaluate(context.Background(), failCmd, 1, "", false)
	if err != nil {
		t.Fatalf("Evaluate streak failed: %v", err)
	}
	if !dec.ShouldSpeak || dec.Trigger != TriggerFailureStreak {
		t.Errorf("Expected failure streak trigger, got: %+v", dec)
	}

	// 3. Quiet mode (hushed)
	decHushed, _ := engine.Evaluate(context.Background(), "rm -rf /", 1, "", true)
	if decHushed.ShouldSpeak {
		t.Errorf("Expected quiet mode to suppress speech")
	}
}
