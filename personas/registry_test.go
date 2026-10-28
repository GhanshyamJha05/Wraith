package persona

import (
	"testing"
)

func TestEmbeddedPersonasLoaded(t *testing.T) {
	r, err := NewRegistry("")
	if err != nil {
		t.Fatalf("Failed to initialize registry: %v", err)
	}

	expected := []string{"victorian_ghost", "sysadmin", "cobol", "marigold"}
	for _, name := range expected {
		p, ok := r.Get(name)
		if !ok {
			t.Errorf("Expected persona %q to be loaded", name)
			continue
		}
		if p.SystemPrompt == "" {
			t.Errorf("Persona %q has empty system prompt", name)
		}
		if p.Skin.Marker == "" {
			t.Errorf("Persona %q missing skin marker", name)
		}
	}

	// Verify Victorian skin
	vic, _ := r.Get("victorian_ghost")
	if vic.Skin.Color.TrueColorHex != "#B39DDB" || !vic.Skin.Style.Italic || !vic.Skin.QuoteStyle {
		t.Errorf("Victorian skin properties mismatch: %+v", vic.Skin)
	}

	// Verify Sysadmin skin
	sys, _ := r.Get("sysadmin")
	if sys.Skin.LinePrefix != "[root@beyond] " || sys.Skin.Color.TrueColorHex != "#00FF66" {
		t.Errorf("Sysadmin skin properties mismatch: %+v", sys.Skin)
	}

	// Verify COBOL skin
	cob, _ := r.Get("cobol")
	if cob.Skin.Color.TrueColorHex != "#FFB000" || !cob.Skin.Style.Uppercase {
		t.Errorf("COBOL skin properties mismatch: %+v", cob.Skin)
	}

	// Verify Marigold skin
	mari, _ := r.Get("marigold")
	if mari.Skin.Color.TrueColorHex != "#FFA726" || mari.Skin.Marker != "🌼" {
		t.Errorf("Marigold skin properties mismatch: %+v", mari.Skin)
	}
}
