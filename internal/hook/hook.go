package hook

import (
	_ "embed"
	"fmt"
	"strings"
)

// Embedded hook templates
//
//go:embed templates/wraith_hook.ps1
var pwshHookTemplate string

//go:embed templates/wraith_hook.bash
var bashHookTemplate string

//go:embed templates/wraith_hook.zsh
var zshHookTemplate string

const (
	MarkerStart = "# >>> WRAITH HOOK START >>>"
	MarkerEnd   = "# <<< WRAITH HOOK END <<<"
)

// ShellType identifies supported shells
type ShellType string

const (
	ShellPwsh ShellType = "pwsh"
	ShellWinPS ShellType = "powershell"
	ShellBash ShellType = "bash"
	ShellZsh  ShellType = "zsh"
)

// GenerateHookScript returns the shell-specific hook with markers
func GenerateHookScript(shell ShellType, binaryPath string, port int, budgetMs int) (string, error) {
	if binaryPath == "" {
		binaryPath = "wraith"
	}

	var rawTemplate string
	switch shell {
	case ShellPwsh, ShellWinPS:
		rawTemplate = pwshHookTemplate
	case ShellBash:
		rawTemplate = bashHookTemplate
	case ShellZsh:
		rawTemplate = zshHookTemplate
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}

	// Normalize binary path for PowerShell/bash
	escapedBin := strings.ReplaceAll(binaryPath, `\`, `/`)
	script := strings.ReplaceAll(rawTemplate, "__WRAITH_BIN__", escapedBin)
	script = strings.ReplaceAll(script, "__WRAITH_PORT__", fmt.Sprintf("%d", port))
	script = strings.ReplaceAll(script, "__WRAITH_BUDGET_MS__", fmt.Sprintf("%d", budgetMs))

	wrapped := fmt.Sprintf("%s\n%s\n%s\n", MarkerStart, strings.TrimSpace(script), MarkerEnd)
	return wrapped, nil
}
