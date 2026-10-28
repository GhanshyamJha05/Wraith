package persona

import (
	"testing"
)

func TestEmbeddedPersonasLoaded(t *testing.T) {
	reg, err := NewRegistry("")
	if err != nil {
		t.Fatalf("failed to initialize registry: %v", err)
	}

	expectedPersonas := []string{"victorian_ghost", "sysadmin", "cobol", "marigold"}
	for _, name := range expectedPersonas {
		p, found := reg.Get(name)
		if !found {
			t.Errorf("expected persona %q to be loaded", name)
			continue
		}
		if p.Title == "" || p.SystemPrompt == "" {
			t.Errorf("persona %q has empty title or system prompt", name)
		}
	}
}
