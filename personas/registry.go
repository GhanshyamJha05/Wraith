package persona

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Embedded default personas so Wraith runs self-contained
//
//go:embed *.yaml
var embeddedPersonas embed.FS

// Persona holds metadata and prompting instructions for a ghost identity
type Persona struct {
	Name         string
	Title        string
	Description  string
	SystemPrompt string
}

// Registry manages loaded personas
type Registry struct {
	personas map[string]Persona
}

// NewRegistry loads personas from embedded assets and optional user directory
func NewRegistry(userDir string) (*Registry, error) {
	r := &Registry{
		personas: make(map[string]Persona),
	}

	// 1. Load embedded defaults
	entries, err := embeddedPersonas.ReadDir(".")
	if err == nil {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".yaml") {
				data, err := embeddedPersonas.ReadFile(e.Name())
				if err == nil {
					p, err := parseSimpleYAML(string(data))
					if err == nil && p.Name != "" {
						r.personas[p.Name] = p
					}
				}
			}
		}
	}

	// 2. Load custom user personas from disk if available
	if userDir != "" {
		files, err := os.ReadDir(userDir)
		if err == nil {
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml") {
					path := filepath.Join(userDir, f.Name())
					bytes, err := os.ReadFile(path)
					if err == nil {
						p, err := parseSimpleYAML(string(bytes))
						if err == nil && p.Name != "" {
							r.personas[p.Name] = p
						}
					}
				}
			}
		}
	}

	return r, nil
}

// Get returns the persona by name or fallback
func (r *Registry) Get(name string) (Persona, bool) {
	p, ok := r.personas[name]
	return p, ok
}

// List returns all loaded persona names
func (r *Registry) List() []Persona {
	list := make([]Persona, 0, len(r.personas))
	for _, p := range r.personas {
		list = append(list, p)
	}
	return list
}

// parseSimpleYAML parses basic key-value and multiline strings without third-party YAML deps
func parseSimpleYAML(content string) (Persona, error) {
	var p Persona
	scanner := bufio.NewScanner(strings.NewReader(content))
	var multilineKey string
	var multilineVal strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if multilineKey != "" {
			if strings.HasPrefix(line, "  ") || strings.TrimSpace(line) == "" {
				multilineVal.WriteString(strings.TrimPrefix(line, "  ") + "\n")
				continue
			} else {
				if multilineKey == "system_prompt" {
					p.SystemPrompt = strings.TrimSpace(multilineVal.String())
				}
				multilineKey = ""
				multilineVal.Reset()
			}
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if val == "|" {
			multilineKey = key
			continue
		}

		switch key {
		case "name":
			p.Name = val
		case "title":
			p.Title = val
		case "description":
			p.Description = val
		case "system_prompt":
			p.SystemPrompt = val
		}
	}

	if multilineKey == "system_prompt" {
		p.SystemPrompt = strings.TrimSpace(multilineVal.String())
	}

	if p.Name == "" {
		return p, fmt.Errorf("missing 'name' field in persona")
	}
	return p, nil
}
