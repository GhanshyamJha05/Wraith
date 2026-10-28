package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultPort        = 23847 // "WRAITH" in vanity port range
	DefaultBudgetMs    = 800
	DefaultOllamaURL   = "http://127.0.0.1:11434"
	DefaultModel       = "gemma3:4b"
	DefaultPersona     = "victorian_ghost"
	DefaultTriggerRate = 0.25
)

// Config represents runtime configuration for Wraith
type Config struct {
	Port         int     `json:"port"`
	BudgetMs     int     `json:"budget_ms"`
	OllamaURL    string  `json:"ollama_url"`
	Model        string  `json:"model"`
	Persona      string  `json:"persona"`
	TriggerRate  float64 `json:"trigger_rate"`
	Hushed       bool    `json:"hushed"`
	DataDir      string  `json:"data_dir"`
	DebugOverlay bool    `json:"debug_overlay"`
}

// GetDataDir returns the standard data directory for Wraith
func GetDataDir() string {
	if custom := os.Getenv("WRAITH_DATA_DIR"); custom != "" {
		return custom
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".wraith"
	}
	return filepath.Join(home, ".wraith")
}

// DefaultConfig returns safe defaults
func DefaultConfig() *Config {
	debug := os.Getenv("WRAITH_DEBUG") == "1"
	return &Config{
		Port:         DefaultPort,
		BudgetMs:     DefaultBudgetMs,
		OllamaURL:    DefaultOllamaURL,
		Model:        DefaultModel,
		Persona:      DefaultPersona,
		TriggerRate:  DefaultTriggerRate,
		Hushed:       false,
		DataDir:      GetDataDir(),
		DebugOverlay: debug,
	}
}
