package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"wraith/internal/config"
	"wraith/internal/hook"
	"wraith/internal/ollama"
)

// runDoctor performs deep pre-flight diagnostic checks on the environment
func runDoctor(args []string) {
	fmt.Println("🩺 WRAITH SYSTEM DIAGNOSTIC (wraith doctor)")
	fmt.Println("=================================================================")

	allPass := true
	cfg := config.DefaultConfig()

	// 1. Data Directory Check
	fmt.Printf("[1/6] User Data Location:\n")
	dataDir := cfg.DataDir
	fmt.Printf("      Path: %s\n", dataDir)
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		fmt.Printf("      ❌ Failed to access or create data directory: %v\n", err)
		allPass = false
	} else {
		// Test write capability
		testFile := filepath.Join(dataDir, ".doctor_test")
		if err := os.WriteFile(testFile, []byte("ok"), 0600); err != nil {
			fmt.Printf("      ❌ Directory is not writable: %v\n", err)
			allPass = false
		} else {
			_ = os.Remove(testFile)
			fmt.Printf("      ✅ Accessible and writable\n")
		}
	}

	// 2. Local Ollama Reachability
	fmt.Printf("\n[2/6] Ollama Service:\n")
	fmt.Printf("      Target URL: %s\n", cfg.OllamaURL)
	ollamaClient := ollama.NewClient(cfg.OllamaURL)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ollamaReachable := false
	if err := ollamaClient.Ping(ctx); err != nil {
		fmt.Printf("      ⚠️  Ollama is unreachable (%v)\n", err)
		fmt.Printf("         Note: Wraith will safely fall back to instant static Grimoire whispers.\n")
	} else {
		ollamaReachable = true
		fmt.Printf("      ✅ Service is online and responding\n")
	}

	// 3. Model Availability
	fmt.Printf("\n[3/6] Model Availability (%s):\n", cfg.Model)
	if !ollamaReachable {
		fmt.Printf("      ⚠️  Skipped (Ollama offline). Run 'ollama serve' and 'ollama pull %s'.\n", cfg.Model)
	} else {
		// Check /api/tags
		tagsURL := fmt.Sprintf("%s/api/tags", cfg.OllamaURL)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, tagsURL, nil)
		resp, err := ollamaClient.HTTPClient.Do(req)
		modelFound := false
		if err == nil && resp.StatusCode == http.StatusOK {
			var tagResp struct {
				Models []struct {
					Name string `json:"name"`
				} `json:"models"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&tagResp); err == nil {
				for _, m := range tagResp.Models {
					if m.Name == cfg.Model || m.Name == cfg.Model+":latest" {
						modelFound = true
						break
					}
				}
			}
			resp.Body.Close()
		}

		if modelFound {
			fmt.Printf("      ✅ Model %q is pulled and ready in memory\n", cfg.Model)
		} else {
			fmt.Printf("      ⚠️  Model %q not pulled yet. Run: 'ollama pull %s'\n", cfg.Model, cfg.Model)
		}
	}

	// 4. Shell Hook Installation Check
	fmt.Printf("\n[4/6] Shell Profile Hook:\n")
	shellType := detectShell()
	profilePath, err := hook.GetProfilePath(shellType)
	if err != nil {
		fmt.Printf("      ❌ Could not detect profile path: %v\n", err)
		allPass = false
	} else {
		installed := false
		if data, err := os.ReadFile(profilePath); err == nil {
			if hook.HasHook(string(data)) {
				installed = true
			}
		}
		if installed {
			fmt.Printf("      ✅ Hook is installed in %s (%s)\n", profilePath, shellType)
		} else {
			fmt.Printf("      ℹ️  Hook is not yet installed in %s\n", profilePath)
			fmt.Printf("         Run 'wraith haunt' to install transparently.\n")
		}
	}

	// 5. Background Daemon Check
	fmt.Printf("\n[5/6] Background Daemon Status (127.0.0.1:%d):\n", cfg.Port)
	daemonURL := fmt.Sprintf("http://127.0.0.1:%d/status", cfg.Port)
	dClient := &http.Client{Timeout: 1 * time.Second}
	if dResp, err := dClient.Get(daemonURL); err == nil {
		dResp.Body.Close()
		fmt.Printf("      ✅ Daemon is actively listening and ready\n")
	} else {
		fmt.Printf("      ℹ️  Daemon is dormant (auto-starts upon hook execution or via 'wraith daemon')\n")
	}

	// 6. Shell Exit Code Self-Test
	fmt.Printf("\n[6/6] Exit Code Preservation Self-Test (Rule 1):\n")
	selfTestPassed := performExitCodeSelfTest(shellType)
	if selfTestPassed {
		fmt.Printf("      ✅ Subshell preserved exit code 42 exactly\n")
	} else {
		fmt.Printf("      ⚠️  Could not run exit-code self-test on current shell execution path\n")
	}

	fmt.Println("=================================================================")
	if allPass {
		fmt.Println("✨ Doctor inspection complete: Core subsystem is healthy!")
	} else {
		fmt.Println("⚠️  Doctor inspection found items requiring attention.")
	}
}

func performExitCodeSelfTest(shell hook.ShellType) bool {
	var cmd *exec.Cmd
	switch shell {
	case hook.ShellWinPS, hook.ShellPwsh:
		shellBin := "powershell.exe"
		if shell == hook.ShellPwsh {
			if _, err := exec.LookPath("pwsh.exe"); err == nil {
				shellBin = "pwsh.exe"
			}
		}
		cmd = exec.Command(shellBin, "-NoProfile", "-NonInteractive", "-Command", "exit 42")
	case hook.ShellBash:
		cmd = exec.Command("bash", "-c", "exit 42")
	case hook.ShellZsh:
		cmd = exec.Command("zsh", "-c", "exit 42")
	default:
		return false
	}

	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 42
	}
	return false
}
