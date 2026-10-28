package hook

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

// TestRule1_ExitCodePreservation_PowerShell executes an actual PowerShell child process
// running the hook script after a failing exit code (e.g. exit 42) and asserts
// that $LASTEXITCODE is preserved as 42.
func TestRule1_ExitCodePreservation_PowerShell(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("PowerShell test only running on Windows")
	}

	shells := []string{"powershell.exe", "pwsh.exe"}

	for _, shellExe := range shells {
		shellPath, err := exec.LookPath(shellExe)
		if err != nil {
			t.Logf("Shell %s not found in PATH, skipping", shellExe)
			continue
		}

		t.Run(shellExe, func(t *testing.T) {
			// Generate the hook script
			hookSnippet, err := GenerateHookScript(ShellPwsh, "wraith.exe", 23847, 800)
			if err != nil {
				t.Fatalf("failed to generate hook: %v", err)
			}

			// Script to run:
			// 1. Define hook function
			// 2. Run a command that exits with 42
			// 3. Run prompt function (which calls Invoke-WraithHook)
			// 4. Output the value of $LASTEXITCODE
			psScript := fmt.Sprintf(`
$ErrorActionPreference = 'Continue'
%s

# Simulate previous command failing with exit 42
powershell -NoProfile -Command "exit 42"
$savedBefore = $LASTEXITCODE

# Run the prompt hook
Invoke-WraithHook

# Verify that $LASTEXITCODE is still 42
if ($LASTEXITCODE -eq 42) {
    Write-Output "EXIT_CODE_PRESERVED:42"
} else {
    Write-Output "EXIT_CODE_CORRUPTED:$LASTEXITCODE"
}
`, hookSnippet)

			cmd := exec.Command(shellPath, "-NoProfile", "-NonInteractive", "-Command", psScript)
			out, err := cmd.CombinedOutput()
			outStr := string(out)

			if !strings.Contains(outStr, "EXIT_CODE_PRESERVED:42") {
				t.Fatalf("Rule 1 violated in %s! Output:\n%s (err: %v)", shellExe, outStr, err)
			}
		})
	}
}

// TestRule1_ExitCodePreservation_Bash verifies exit code preservation for bash
func TestRule1_ExitCodePreservation_Bash(t *testing.T) {
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH, skipping bash hook test")
	}

	// Verify bash is functional on Windows (e.g. WSL registration errors)
	checkCmd := exec.Command(bashPath, "-c", "echo ok")
	if out, err := checkCmd.CombinedOutput(); err != nil || !strings.Contains(string(out), "ok") {
		t.Skipf("bash is not functional on this environment (%v, out: %s), skipping", err, string(out))
	}

	hookSnippet, err := GenerateHookScript(ShellBash, "wraith", 23847, 800)
	if err != nil {
		t.Fatalf("failed to generate bash hook: %v", err)
	}

	bashScript := fmt.Sprintf(`
%s

(exit 42)
_wraith_precmd
if [ $? -eq 42 ]; then
    echo "EXIT_CODE_PRESERVED:42"
else
    echo "EXIT_CODE_CORRUPTED:$?"
fi
`, hookSnippet)

	cmd := exec.Command(bashPath, "-c", bashScript)
	out, err := cmd.CombinedOutput()
	outStr := string(out)

	if !strings.Contains(outStr, "EXIT_CODE_PRESERVED:42") {
		t.Fatalf("Rule 1 violated in bash! Output:\n%s (err: %v)", outStr, err)
	}
}
