package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetDataDirStandardOS(t *testing.T) {
	dir := GetDataDir()
	if dir == "" {
		t.Fatalf("GetDataDir() returned empty string")
	}

	switch runtime.GOOS {
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			expected := filepath.Join(local, "wraith")
			if dir != expected {
				t.Errorf("Expected %q, got %q", expected, dir)
			}
		}
	case "darwin":
		if home, err := os.UserHomeDir(); err == nil {
			expected := filepath.Join(home, "Library", "Application Support", "wraith")
			if dir != expected {
				t.Errorf("Expected %q, got %q", expected, dir)
			}
		}
	}
}

func TestGetDataDirOverride(t *testing.T) {
	custom := filepath.Join(os.TempDir(), "wraith-custom-test")
	t.Setenv("WRAITH_DATA_DIR", custom)
	if got := GetDataDir(); got != custom {
		t.Errorf("Expected override %q, got %q", custom, got)
	}
}

