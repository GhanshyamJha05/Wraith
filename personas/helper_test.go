package persona

import (
	"strings"
	"testing"
	"wraith/internal/grimoire"
)

func TestHelperPersonaProperties(t *testing.T) {
	r, err := NewRegistry("")
	if err != nil {
		t.Fatalf("Failed to initialize registry: %v", err)
	}

	helper, ok := r.Get("helper")
	if !ok {
		t.Fatalf("Expected 'helper' persona to be loaded in registry")
	}

	if helper.Title != "The Plain Helper" {
		t.Errorf("Unexpected title: %q", helper.Title)
	}

	if helper.Skin.Marker != "ℹ️" {
		t.Errorf("Expected marker ℹ️, got %q", helper.Skin.Marker)
	}

	if helper.Skin.MarkerPlain != "[info]" {
		t.Errorf("Expected plain marker [info], got %q", helper.Skin.MarkerPlain)
	}

	// Verify SystemPrompt enforces critical safety constraints
	prompt := strings.ToLower(helper.SystemPrompt)
	if !strings.Contains(prompt, "never suggest destructive commands") {
		t.Errorf("Helper system prompt must strictly prohibit destructive suggestions")
	}
	if !strings.Contains(prompt, "never claim that you executed") {
		t.Errorf("Helper system prompt must strictly forbid claiming actions")
	}
}

func TestHelperGrimoireSafety(t *testing.T) {
	triggers := []string{"failure_streak", "danger", "late_night", "stale_directory", "general_failure"}
	destructivePhrases := []string{"rm -rf", "format", "drop table", "force", "--hard", "delete everything"}
	actionClaims := []string{"i ran", "i executed", "i fixed", "i deleted", "i tested"}

	for _, tr := range triggers {
		for i := 0; i < 20; i++ {
			whisper := strings.ToLower(grimoire.GetWhisper("helper", tr))
			for _, d := range destructivePhrases {
				if strings.Contains(whisper, d) {
					t.Errorf("Helper grimoire contained destructive suggestion %q in %q", d, whisper)
				}
			}
			for _, a := range actionClaims {
				if strings.Contains(whisper, a) {
					t.Errorf("Helper grimoire made false action claim %q in %q", a, whisper)
				}
			}
		}
	}
}

