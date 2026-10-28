package sanitize

import (
	"regexp"
	"strings"
	"unicode"
)

// Rule 4: Sanitize all model output before printing:
// strip ANSI/OSC escape sequences and control characters,
// enforce max 25 words and one line.

var (
	// Matches ANSI escape sequences (CSI sequences like \x1b[...m, OSC sequences like \x1b]...BEL/\x1b\\)
	ansiPattern = regexp.MustCompile(`\x1b(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~]|\].*?(?:\x07|\x1b\\))`)
	
	// OSC pattern specifically: ESC ] ... (BEL or ESC \)
	oscPattern = regexp.MustCompile(`\x1b\][^\x07\x1b]*(\x07|\x1b\\)`)
)

const MaxWords = 25

// Sanitize cleans model output string according to Rule 4
func Sanitize(input string) string {
	if input == "" {
		return ""
	}

	// 1. Strip OSC sequences first
	cleaned := oscPattern.ReplaceAllString(input, "")

	// 2. Strip standard ANSI / CSI escape sequences
	cleaned = ansiPattern.ReplaceAllString(cleaned, "")

	// 3. Strip surrounding quotes, markdown bold/italic tags, backticks
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Trim(cleaned, `"'` + "`")

	// 4. Flatten newlines and control characters to spaces
	var sb strings.Builder
	for _, r := range cleaned {
		if r == '\r' || r == '\n' || r == '\t' {
			sb.WriteRune(' ')
		} else if unicode.IsControl(r) {
			// drop control characters (BEL, BS, etc.)
			continue
		} else {
			sb.WriteRune(r)
		}
	}

	// 5. Collapse multiple whitespaces and enforce 1 single line
	words := strings.Fields(sb.String())
	if len(words) == 0 {
		return ""
	}

	// 6. Enforce max 25 words limit
	if len(words) > MaxWords {
		words = words[:MaxWords]
	}

	res := strings.Join(words, " ")
	return res
}
