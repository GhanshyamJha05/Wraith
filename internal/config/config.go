package config

import (
	"os"
	"path/filepath"
	"runtime"
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

// GetDataDir returns the standard OS user data directory for Wraith
// Windows: %LOCALAPPDATA%\wraith
// Linux:   $XDG_DATA_HOME/wraith or ~/.local/share/wraith
// macOS:   ~/Library/Application Support/wraith
func GetDataDir() string {
	if custom := os.Getenv("WRAITH_DATA_DIR"); custom != "" {
		return custom
	}

	switch runtime.GOOS {
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "wraith")
		}
	case "darwin":
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "Library", "Application Support", "wraith")
		}
	default: // Linux / Unix / BSD
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, "wraith")
		}
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ".local", "share", "wraith")
		}
	}

	// Fallback to user home
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".wraith")
	}
	return ".wraith"
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
