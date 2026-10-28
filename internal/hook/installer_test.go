package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallAndExorciseHook(t *testing.T) {
	tempDir := t.TempDir()
	profilePath := filepath.Join(tempDir, "profile.ps1")

	initialContent := "# Existing user profile configuration\n$env:MY_VAR = 'hello'\n"
	err := os.WriteFile(profilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test profile: %v", err)
	}

	// 1. Dry run
	res, err := InstallHook(profilePath, ShellPwsh, "wraith.exe", 23847, 800, true)
	if err != nil {
		t.Fatalf("dry run failed: %v", err)
	}
	if !strings.Contains(res.DiffLines, MarkerStart) {
		t.Errorf("expected diff to contain marker start")
	}

	contentAfterDryRun, _ := os.ReadFile(profilePath)
	if string(contentAfterDryRun) != initialContent {
		t.Errorf("dry run modified file!")
	}

	// 2. Real Install
	res, err = InstallHook(profilePath, ShellPwsh, "wraith.exe", 23847, 800, false)
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if res.BackupPath == "" {
		t.Errorf("expected backup file to be created")
	}
	if _, err := os.Stat(res.BackupPath); os.IsNotExist(err) {
		t.Errorf("backup file does not exist on disk")
	}

	installedBytes, _ := os.ReadFile(profilePath)
	installedContent := string(installedBytes)
	if !strings.Contains(installedContent, MarkerStart) || !strings.Contains(installedContent, MarkerEnd) {
		t.Errorf("installed profile missing marker tags")
	}
	if !strings.Contains(installedContent, "$env:MY_VAR = 'hello'") {
		t.Errorf("original content was wiped or lost")
	}

	// 3. Repeat Install (Idempotency)
	resRepeat, err := InstallHook(profilePath, ShellPwsh, "wraith.exe", 23847, 800, false)
	if err != nil {
		t.Fatalf("repeat install failed: %v", err)
	}
	if !resRepeat.AlreadyDone {
		t.Errorf("expected AlreadyDone to be true on second install")
	}

	// 4. Exorcise (Clean removal)
	removed, err := ExorciseHook(profilePath)
	if err != nil {
		t.Fatalf("exorcise failed: %v", err)
	}
	if !removed {
		t.Errorf("expected removed to be true")
	}

	afterExorciseBytes, _ := os.ReadFile(profilePath)
	afterExorciseContent := strings.TrimSpace(string(afterExorciseBytes))
	expectedContent := strings.TrimSpace(initialContent)

	if afterExorciseContent != expectedContent {
		t.Errorf("profile after exorcism differs.\nGot:\n%s\nExpected:\n%s", afterExorciseContent, expectedContent)
	}
}
