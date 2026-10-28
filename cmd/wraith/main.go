package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"wraith/internal/config"
	"wraith/internal/daemon"
	"wraith/internal/hook"
	"wraith/internal/memory"
	"wraith/internal/theme"
	personas "wraith/personas"
)

const version = "0.2.0-visual"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	subcommand := os.Args[1]
	subArgs := os.Args[2:]

	switch subcommand {
	case "daemon":
		runDaemon(subArgs)
	case "hook":
		runHook(subArgs)
	case "status":
		runStatus(subArgs)
	case "haunt":
		runHaunt(subArgs)
	case "exorcise":
		runExorcise(subArgs)
	case "purge":
		runPurge(subArgs)
	case "hush":
		runHush(subArgs)
	case "wake":
		runWake(subArgs)
	case "persona":
		runPersona(subArgs)
	case "demo":
		runDemo(subArgs)
	case "seance":
		runSeance(subArgs)
	case "doctor":
		runDoctor(subArgs)
	case "theme":
		runTheme(subArgs)
	case "version":
		fmt.Printf("wraith v%s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n\n", subcommand)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(theme.GhostBanner(false))
	fmt.Println(`Usage:
  wraith <command> [options]

Commands:
  haunt       Install shell hook into profile (transparent opt-in)
  exorcise    Remove shell hook from profile cleanly
  status      Show daemon state, Ollama reachability, and recent whispers
  hush        Quiet mode: silence all whispers
  wake        Re-enable whispers
  persona     View or switch the active ghost persona
  daemon      Run the local background whisper daemon (127.0.0.1)
  hook        Internal hook dispatcher called by shell prompt
  demo        Replay a scripted demonstration of ghost triggers
  doctor      Run system environment diagnostics and exit-code self-test
  theme       Print Windows Terminal color scheme JSON and instructions
  seance      Generate a spooky summary of your past shell sessions
  purge       Delete all stored telemetry and memory database
  version     Display Wraith version info`)
}

func runDaemon(args []string) {
	fs := flag.NewFlagSet("daemon", flag.ExitOnError)
	port := fs.Int("port", config.DefaultPort, "Port to listen on (127.0.0.1)")
	model := fs.String("model", config.DefaultModel, "Ollama model name (gemma3:4b, gemma3:1b, gemma3:12b)")
	persona := fs.String("persona", config.DefaultPersona, "Initial ghost persona")
	budget := fs.Int("budget-ms", config.DefaultBudgetMs, "Hard prompt response budget in milliseconds")
	_ = fs.Parse(args)

	cfg := config.DefaultConfig()
	cfg.Port = *port
	cfg.Model = *model
	cfg.Persona = *persona
	cfg.BudgetMs = *budget

	srv, err := daemon.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[wraith daemon] Failed to initialize: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[wraith daemon] Ghost listening on 127.0.0.1:%d (model: %s, budget: %dms)\n", cfg.Port, cfg.Model, cfg.BudgetMs)
	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "[wraith daemon] Exited: %v\n", err)
		os.Exit(1)
	}
}

func runHook(args []string) {
	if len(args) == 0 || args[0] != "ingest" {
		return
	}

	fs := flag.NewFlagSet("ingest", flag.ContinueOnError)
	shell := fs.String("shell", "pwsh", "Shell type")
	lastExit := fs.Int("last-exit", 0, "Last command exit code")
	cmd := fs.String("cmd", "", "Executed command")
	port := fs.Int("port", config.DefaultPort, "Daemon port")
	budget := fs.Int("budget-ms", config.DefaultBudgetMs, "Max wait budget")
	plain := fs.Bool("plain", false, "Disable color and ANSI escapes")
	_ = fs.Parse(args[1:])

	daemonURL := fmt.Sprintf("http://127.0.0.1:%d/ingest", *port)
	cwd, _ := os.Getwd()

	reqPayload := daemon.IngestRequest{
		Shell:      *shell,
		ExitCode:   *lastExit,
		Command:    *cmd,
		Cwd:        cwd,
		DurationMs: 0,
		BudgetMs:   *budget,
	}

	payloadBytes, _ := json.Marshal(reqPayload)

	client := &http.Client{
		Timeout: time.Duration(*budget+100) * time.Millisecond,
	}

	resp, err := client.Post(daemonURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		spawnDaemonInBackground()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var ingestResp daemon.IngestResponse
		if err := json.NewDecoder(resp.Body).Decode(&ingestResp); err == nil {
			if ingestResp.Whisper != "" {
				// Load persona skin from registry for styling
				reg, _ := personas.NewRegistry("")
				var skin theme.Skin
				if p, ok := reg.Get(ingestResp.Persona); ok {
					skin = p.Skin
				} else {
					skin = theme.Skin{
						Marker:      "👻",
						MarkerPlain: "[*]",
						Color:       theme.ColorDef{TrueColorHex: "#B39DDB", Color256: 183, Color16: 95},
						Style:       theme.TextStyle{Italic: true},
					}
				}

				level := theme.DetectColorLevel(*plain)
				styled := theme.StyleWhisper(ingestResp.Whisper, skin, level)
				overlay := ingestResp.Overlay
				fmt.Printf("%s%s\n", styled, overlay)
			}
		}
	}
}

func spawnDaemonInBackground() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	cmd := exec.Command(exe, "daemon")
	cmd.SysProcAttr = getSysProcAttr()
	_ = cmd.Start()
}

func runStatus(args []string) {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	plain := fs.Bool("plain", false, "Plain ASCII output without color")
	_ = fs.Parse(args)

	cfg := config.DefaultConfig()
	url := fmt.Sprintf("http://127.0.0.1:%d/status", cfg.Port)

	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(theme.StatusHeader(*plain, cfg.Persona, true))
		fmt.Printf("Daemon: NOT RUNNING (offline at 127.0.0.1:%d)\n", cfg.Port)
		fmt.Printf("Run 'wraith daemon' or execute any shell command with hook installed to wake it.\n")
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("Error reading status: %v\n", err)
		return
	}

	hushed := false
	if h, ok := data["hushed"].(bool); ok {
		hushed = h
	}
	activePersona := fmt.Sprintf("%v", data["persona"])

	fmt.Println(theme.StatusHeader(*plain, activePersona, hushed))
	fmt.Printf("  Daemon:          %v\n", data["daemon"])
	fmt.Printf("  Ollama URL:      %v\n", data["ollama_url"])
	fmt.Printf("  Model:           %v (Reachable: %v)\n", data["model"], data["model_reachable"])
	fmt.Printf("  Active Persona:  %v\n", data["persona"])
	fmt.Printf("  Hushed:          %v\n", data["hushed"])

	if recent, ok := data["recent_whispers"].([]interface{}); ok && len(recent) > 0 {
		fmt.Printf("\nRecent Whisper Log (Last %d):\n", len(recent))
		for _, item := range recent {
			if m, ok := item.(map[string]interface{}); ok {
				fmt.Printf("  [%s] (%s) %s\n", m["source"], m["model_tag"], m["whisper_text"])
			}
		}
	} else {
		fmt.Printf("\nNo whispers logged yet.\n")
	}
}

func runHaunt(args []string) {
	fs := flag.NewFlagSet("haunt", flag.ContinueOnError)
	plain := fs.Bool("plain", false, "Disable color in preview")
	_ = fs.Parse(args)

	shellType := detectShell()
	profilePath, err := hook.GetProfilePath(shellType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting profile: %v\n", err)
		os.Exit(1)
	}

	binPath, _ := os.Executable()
	cfg := config.DefaultConfig()

	// 1. Dry run preview (Rule 7: full transparency)
	res, err := hook.InstallHook(profilePath, shellType, binPath, cfg.Port, cfg.BudgetMs, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating hook: %v\n", err)
		os.Exit(1)
	}

	if res.AlreadyDone {
		fmt.Printf("✨ Wraith is already haunting your profile at:\n  %s\n", profilePath)
		return
	}

	// ASCII ghost banner and themed prompt ceremony
	fmt.Print(theme.GhostBanner(*plain))
	fmt.Println("=================================================================")
	fmt.Println("👻 WRAITH PROFILE MODIFICATION PREVIEW (Opt-In Transparency)")
	fmt.Println("=================================================================")
	fmt.Printf("Target Profile: %s\n\n", profilePath)
	fmt.Println("The following lines will be appended to your profile:")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Print(res.DiffLines)
	fmt.Println("-----------------------------------------------------------------")
	fmt.Print("🕯️  Awaken the specter within this shell parlor? [y/N]: ")

	var answer string
	fmt.Scanln(&answer)
	answer = strings.ToLower(strings.TrimSpace(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("Aborted. The ghosts remain in slumber. No modifications were made.")
		return
	}

	// 2. Perform real install
	res, err = hook.InstallHook(profilePath, shellType, binPath, cfg.Port, cfg.BudgetMs, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Installation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Successfully haunted!\n")
	if res.BackupPath != "" {
		fmt.Printf("   A backup of your original profile was saved at:\n   %s\n", res.BackupPath)
	}
	fmt.Printf("   Restart your shell or run: . \"%s\"\n", profilePath)
}

func runExorcise(args []string) {
	fs := flag.NewFlagSet("exorcise", flag.ContinueOnError)
	plain := fs.Bool("plain", false, "Disable colored farewell ceremony")
	_ = fs.Parse(args)

	shellType := detectShell()
	profilePath, err := hook.GetProfilePath(shellType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding profile: %v\n", err)
		os.Exit(1)
	}

	removed, err := hook.ExorciseHook(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exorcism failed: %v\n", err)
		os.Exit(1)
	}

	if removed {
		fmt.Print(theme.FarewellCeremony(*plain, profilePath))
	} else {
		fmt.Printf("No Wraith hook was found in: %s\n", profilePath)
	}
}

func runPurge(args []string) {
	cfg := config.DefaultConfig()
	dbPath := filepath.Join(cfg.DataDir, "memory.db")

	store, err := memory.NewStore(dbPath)
	if err == nil {
		_ = store.Purge(context.Background())
		_ = store.Close()
	}
	_ = os.Remove(dbPath)

	fmt.Println("🧹 All memory, history logs, and sqlite records have been purged.")
}

func runHush(args []string) {
	cfg := config.DefaultConfig()
	url := fmt.Sprintf("http://127.0.0.1:%d/hush", cfg.Port)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Quiet mode enabled (daemon was offline).")
		return
	}
	resp.Body.Close()
	fmt.Println("🤫 The ghosts have been hushed. Wraith will remain silent.")
}

func runWake(args []string) {
	cfg := config.DefaultConfig()
	url := fmt.Sprintf("http://127.0.0.1:%d/wake", cfg.Port)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Whispers re-enabled.")
		return
	}
	resp.Body.Close()
	fmt.Println("👁️ The ghosts are awake. Whispers resumed.")
}

func runPersona(args []string) {
	cfg := config.DefaultConfig()
	reg, _ := personas.NewRegistry(filepath.Join(cfg.DataDir, "personas"))

	if len(args) == 0 {
		fmt.Println("Available Ghost Personas:")
		for _, p := range reg.List() {
			fmt.Printf("  • %-16s %s\n    %s\n", p.Name, p.Title, p.Description)
		}
		fmt.Printf("\nSwitch persona using: wraith persona <name>\n")
		return
	}

	target := args[0]
	if _, ok := reg.Get(target); !ok {
		fmt.Fprintf(os.Stderr, "Unknown persona: %q. Run 'wraith persona' to see available.\n", target)
		os.Exit(1)
	}

	url := fmt.Sprintf("http://127.0.0.1:%d/persona", cfg.Port)
	payload, _ := json.Marshal(map[string]string{"persona": target})
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("Switched default persona to %q (daemon offline).\n", target)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("🎭 Active ghost persona changed to: %q\n", target)
}

func runTheme(args []string) {
	fmt.Println("🎨 WRAITH SPECTRAL THEME (Windows Terminal)")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println(theme.ThemeInstructions())
	fmt.Println("JSON Snippet to copy:")
	fmt.Println(theme.WindowsTerminalSchemeJSON)
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("File also available at: docs/windows-terminal-scheme.json")
}

func detectShell() hook.ShellType {
	if runtime.GOOS == "windows" {
		if os.Getenv("PSModulePath") != "" {
			if strings.Contains(strings.ToLower(os.Getenv("POWERSHELL_DISTRIBUTION_CHANNEL")), "core") {
				return hook.ShellPwsh
			}
			return hook.ShellWinPS
		}
		return hook.ShellWinPS
	}
	shellEnv := os.Getenv("SHELL")
	if strings.Contains(shellEnv, "zsh") {
		return hook.ShellZsh
	}
	return hook.ShellBash
}
