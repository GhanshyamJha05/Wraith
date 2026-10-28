package theme

import (
	"fmt"
	"os"
	"strings"
)

// ColorLevel represents terminal color capabilities
type ColorLevel int

const (
	ColorNone ColorLevel = iota
	Color16
	Color256
	ColorTrueColor
)

// TextStyle configures text rendition
type TextStyle struct {
	Italic    bool
	Bold      bool
	Dim       bool
	Uppercase bool
}

// ColorDef holds fallbacks across color support depths
type ColorDef struct {
	TrueColorHex string // e.g. "#B39DDB"
	Color256     int    // e.g. 183
	Color16      int    // ANSI 30-37 or 90-97 (e.g. 95 for bright magenta, 92 bright green)
}

// Skin defines visual skinning for a persona
type Skin struct {
	Marker      string    `json:"marker" yaml:"marker"`             // Unicode or width-safe glyph
	MarkerPlain string    `json:"marker_plain" yaml:"marker_plain"` // ASCII fallback marker
	LinePrefix  string    `json:"line_prefix" yaml:"line_prefix"`   // e.g. "[root@beyond] "
	Color       ColorDef  `json:"color" yaml:"color"`
	Style       TextStyle `json:"style" yaml:"style"`
	QuoteStyle  bool      `json:"quote_style" yaml:"quote_style"`   // wrap whisper in matching quotes like “...”
}

// DetectColorLevel checks terminal environment and NO_COLOR standard
func DetectColorLevel(plainFlag bool) ColorLevel {
	if plainFlag {
		return ColorNone
	}

	// 1. NO_COLOR standard: if set and not empty, turn off all colors
	if os.Getenv("NO_COLOR") != "" {
		return ColorNone
	}

	// 2. TERM=dumb check
	term := strings.ToLower(os.Getenv("TERM"))
	if term == "dumb" {
		return ColorNone
	}

	// 3. Windows Terminal, ConEmu, WezTerm, Alacritty, Ghostty or COLORTERM check
	colorTerm := strings.ToLower(os.Getenv("COLORTERM"))
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return ColorTrueColor
	}

	if os.Getenv("WT_SESSION") != "" { // Windows Terminal
		return ColorTrueColor
	}

	// 4. Check 256-color support
	if strings.Contains(term, "256color") || strings.Contains(term, "xterm") {
		return Color256
	}

	// 5. Windows standard console support or default basic terminal
	return Color16
}

// StyleWhisper applies Wraith's theme and ANSI escapes to an already sanitized whisper string.
// Model output must NEVER supply control or color codes; styling is applied strictly here.
func StyleWhisper(whisper string, skin Skin, level ColorLevel) string {
	if whisper == "" {
		return ""
	}

	text := whisper
	if skin.Style.Uppercase {
		text = strings.ToUpper(text)
	}
	if skin.QuoteStyle {
		text = fmt.Sprintf("“%s”", text)
	}

	// Plain mode or no color
	if level == ColorNone {
		marker := skin.MarkerPlain
		if marker == "" {
			marker = skin.Marker
		}
		if skin.LinePrefix != "" {
			return fmt.Sprintf("%s %s%s", marker, skin.LinePrefix, text)
		}
		return fmt.Sprintf("%s %s", marker, text)
	}

	// Build ANSI sequence
	var openCodes []string

	// Color sequence
	switch level {
	case ColorTrueColor:
		if skin.Color.TrueColorHex != "" {
			r, g, b, err := hexToRGB(skin.Color.TrueColorHex)
			if err == nil {
				openCodes = append(openCodes, fmt.Sprintf("38;2;%d;%d;%d", r, g, b))
			} else if skin.Color.Color256 > 0 {
				openCodes = append(openCodes, fmt.Sprintf("38;5;%d", skin.Color.Color256))
			} else if skin.Color.Color16 > 0 {
				openCodes = append(openCodes, fmt.Sprintf("%d", skin.Color.Color16))
			}
		} else if skin.Color.Color256 > 0 {
			openCodes = append(openCodes, fmt.Sprintf("38;5;%d", skin.Color.Color256))
		} else if skin.Color.Color16 > 0 {
			openCodes = append(openCodes, fmt.Sprintf("%d", skin.Color.Color16))
		}
	case Color256:
		if skin.Color.Color256 > 0 {
			openCodes = append(openCodes, fmt.Sprintf("38;5;%d", skin.Color.Color256))
		} else if skin.Color.Color16 > 0 {
			openCodes = append(openCodes, fmt.Sprintf("%d", skin.Color.Color16))
		}
	case Color16:
		if skin.Color.Color16 > 0 {
			openCodes = append(openCodes, fmt.Sprintf("%d", skin.Color.Color16))
		}
	}

	// Text styles
	if skin.Style.Bold {
		openCodes = append(openCodes, "1")
	}
	if skin.Style.Dim {
		openCodes = append(openCodes, "2")
	}
	if skin.Style.Italic {
		openCodes = append(openCodes, "3")
	}

	prefix := ""
	if skin.LinePrefix != "" {
		prefix = skin.LinePrefix
	}

	marker := skin.Marker
	if marker == "" {
		marker = "👻"
	}

	if len(openCodes) == 0 {
		return fmt.Sprintf("%s %s%s", marker, prefix, text)
	}

	ansiStart := fmt.Sprintf("\033[%sm", strings.Join(openCodes, ";"))
	ansiEnd := "\033[0m"

	return fmt.Sprintf("%s %s%s%s%s", marker, ansiStart, prefix, text, ansiEnd)
}

func hexToRGB(hex string) (int, int, int, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex length")
	}
	var r, g, b int
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return 0, 0, 0, err
	}
	return r, g, b, nil
}
