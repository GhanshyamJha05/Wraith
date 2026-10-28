package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// InstallResult holds the outcome of an install action
type InstallResult struct {
	ProfilePath string
	BackupPath  string
	DiffLines   string
	AlreadyDone bool
}

// GetProfilePath returns the standard profile path for the target shell
func GetProfilePath(shell ShellType) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine home directory: %w", err)
	}

	switch shell {
	case ShellPwsh:
		// PowerShell Core (pwsh)
		return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
	case ShellWinPS:
		// Windows PowerShell 5.1
		return filepath.Join(home, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"), nil
	case ShellBash:
		return filepath.Join(home, ".bashrc"), nil
	case ShellZsh:
		return filepath.Join(home, ".zshrc"), nil
	default:
		return "", fmt.Errorf("unknown shell: %s", shell)
	}
}

// InstallHook installs the hook snippet into the specified profile
// If dryRun is true, it does not write changes but returns diff/preview.
func InstallHook(profilePath string, shell ShellType, binaryPath string, port int, budgetMs int, dryRun bool) (*InstallResult, error) {
	hookCode, err := GenerateHookScript(shell, binaryPath, port, budgetMs)
	if err != nil {
		return nil, err
	}

	res := &InstallResult{
		ProfilePath: profilePath,
		DiffLines:   hookCode,
	}

	var existingContent string
	if _, err := os.Stat(profilePath); err == nil {
		bytes, err := os.ReadFile(profilePath)
		if err != nil {
			return nil, fmt.Errorf("reading profile %s: %w", profilePath, err)
		}
		existingContent = string(bytes)
	}

	if strings.Contains(existingContent, MarkerStart) {
		res.AlreadyDone = true
		return res, nil
	}

	if dryRun {
		return res, nil
	}

	// Make parent directories if needed
	if err := os.MkdirAll(filepath.Dir(profilePath), 0755); err != nil {
		return nil, fmt.Errorf("creating directory for profile: %w", err)
	}

	// Create backup if file exists
	if existingContent != "" {
		backupPath := fmt.Sprintf("%s.wraith_bak_%s", profilePath, time.Now().Format("20060102_150405"))
		if err := os.WriteFile(backupPath, []byte(existingContent), 0644); err != nil {
			return nil, fmt.Errorf("creating profile backup: %w", err)
		}
		res.BackupPath = backupPath
	}

	// Append hook
	newContent := existingContent
	if newContent != "" && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += hookCode

	if err := os.WriteFile(profilePath, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("writing updated profile: %w", err)
	}

	return res, nil
}

// ExorciseHook removes the hook from the profile cleanly
func ExorciseHook(profilePath string) (bool, error) {
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return false, nil
	}

	bytes, err := os.ReadFile(profilePath)
	if err != nil {
		return false, fmt.Errorf("reading profile %s: %w", profilePath, err)
	}
	content := string(bytes)

	startIdx := strings.Index(content, MarkerStart)
	endIdx := strings.Index(content, MarkerEnd)

	if startIdx == -1 || endIdx == -1 || endIdx < startIdx {
		return false, nil // Nothing to remove
	}

	// Create backup before removing
	backupPath := fmt.Sprintf("%s.wraith_exorcise_%s", profilePath, time.Now().Format("20060102_150405"))
	_ = os.WriteFile(backupPath, bytes, 0644)

	// Cut out the block including MarkerEnd and newline
	before := content[:startIdx]
	after := content[endIdx+len(MarkerEnd):]
	after = strings.TrimPrefix(after, "\r\n")
	after = strings.TrimPrefix(after, "\n")

	cleaned := strings.TrimRight(before+after, "\r\n ") + "\n"
	if err := os.WriteFile(profilePath, []byte(cleaned), 0644); err != nil {
		return false, fmt.Errorf("writing cleaned profile: %w", err)
	}

	return true, nil
}
