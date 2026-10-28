package theme

import (
	"regexp"
	"strings"
	"testing"
	"wraith/internal/sanitize"
)

func TestHardRule_StyledOutputContainsOnlyWraithsEscapes(t *testing.T) {
	// Adversarial model output attempting ANSI escapes, OSC commands, bold, and cursor jumps
	maliciousModelOutputs := []string{
		"\x1b[31;1mHacked RED\x1b[0m",
		"\x1b]0;Evil Terminal Title\x07echo exploit",
		"\x1b[2J\x1b[HClear screen attack",
		"Normal text with \x1b[38;2;255;0;0mhidden injection",
		"Simple benign whisper from model",
	}

	testSkin := Skin{
		Marker:      "👻",
		MarkerPlain: "[*]",
		LinePrefix:  "[ghost] ",
		Color: ColorDef{
			TrueColorHex: "#B39DDB",
			Color256:     183,
			Color16:      95,
		},
		Style: TextStyle{
			Italic: true,
		},
		QuoteStyle: true,
	}

	for _, raw := range maliciousModelOutputs {
		// 1. Model output MUST be sanitized FIRST
		sanitized := sanitize.Sanitize(raw)

		// Verify model output alone has zero escapes
		if strings.Contains(sanitized, "\x1b") {
			t.Fatalf("Sanitization failed to remove escape code from %q: got %q", raw, sanitized)
		}

		// 2. Wraith applies theme styling AFTER sanitization
		styled := StyleWhisper(sanitized, testSkin, ColorTrueColor)

		// 3. Extract all ANSI sequences present in styled output
		re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
		allMatches := re.FindAllString(styled, -1)

		// Verify that ANY escape sequence present matches exactly Wraith's own approved tokens:
		// Expected: \x1b[38;2;179;157;219;3m (Truecolor Lavender + Italic) and \x1b[0m (Reset)
		for _, m := range allMatches {
			if m != "\x1b[38;2;179;157;219;3m" && m != "\x1b[0m" {
				t.Errorf("Unexpected escape code found in styled output: %q (full output: %q)", m, styled)
			}
		}

		// Verify no OSC escapes or non-SGR sequences exist
		if strings.Contains(styled, "\x1b]") || strings.Contains(styled, "\x07") {
			t.Errorf("OSC or BEL escape found in styled output: %q", styled)
		}
	}
}

func TestNoColorAndPlainFlag(t *testing.T) {
	skin := Skin{
		Marker:      "👻",
		MarkerPlain: "[ghost]",
		LinePrefix:  "[root] ",
		Color: ColorDef{
			Color16: 92,
		},
	}

	// Plain color mode
	plain := StyleWhisper("Something failed", skin, ColorNone)
	if strings.Contains(plain, "\x1b") {
		t.Errorf("Plain output must not contain ANSI escapes, got: %q", plain)
	}
	if !strings.HasPrefix(plain, "[ghost] [root] Something failed") {
		t.Errorf("Expected ASCII fallback marker and prefix, got: %q", plain)
	}
}

