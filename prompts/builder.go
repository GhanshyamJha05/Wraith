package model

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed whisper_v1.txt
var promptTemplateText string

var tmpl = template.Must(template.New("whisper_v1").Parse(promptTemplateText))

// PromptData holds variables passed to the prompt template
type PromptData struct {
	PersonaPrompt   string
	Shell           string
	ExitCode        int
	Hour            int
	Cwd             string
	MemorySummary   string
	RedactedCommand string
	DurationMs      int64
}

// BuildPrompt creates the full prompt string, encapsulating untrusted command text in boundary tags
func BuildPrompt(data PromptData) (string, error) {
	// Defensively escape any accidental closing tag within the command text
	safeCmd := strings.ReplaceAll(data.RedactedCommand, "</user_shell_event>", "[TAG_DEFUSED]")
	data.RedactedCommand = safeCmd

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("rendering prompt template: %w", err)
	}

	return buf.String(), nil
}
