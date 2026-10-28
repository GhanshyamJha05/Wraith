package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"wraith/internal/config"
	"wraith/internal/memory"
	"wraith/internal/theme"
	personas "wraith/personas"
)

// runSeance produces a spooky retrospective from the SQLite log
func runSeance(args []string) {
	cfg := config.DefaultConfig()
	dbPath := filepath.Join(cfg.DataDir, "memory.db")
	store, err := memory.NewStore(dbPath)
	if err != nil {
		fmt.Printf("No spirits could be summoned (could not read %s): %v\n", dbPath, err)
		return
	}
	defer store.Close()

	events, err := store.GetRecentSources(context.Background(), 100)
	if err != nil || len(events) == 0 {
		fmt.Println("🕯️ The seance finds only silence. No past commands have been recorded yet.")
		return
	}

	totalWhispers := len(events)
	modelCount := 0
	grimoireCount := 0
	failures := 0

	for _, e := range events {
		if e.Source == memory.SourceModel {
			modelCount++
		} else if e.Source == memory.SourceGrimoire {
			grimoireCount++
		}
		if e.ExitCode != 0 {
			failures++
		}
	}

	fmt.Println("🕯️ ===== WRAITH SÉANCE: The Spirits Speak =====")
	fmt.Printf("Observed Hauntings:      %d total encounters\n", totalWhispers)
	fmt.Printf("Summoned via Model:      %d whispers\n", modelCount)
	fmt.Printf("Drawn from Grimoire:     %d ancient incantations\n", grimoireCount)
	fmt.Printf("Blunders Witnessed:      %d failed executions\n", failures)
	fmt.Println("-----------------------------------------------")
	fmt.Println("Recent spectral whisper remembered:")
	if len(events) > 0 {
		fmt.Printf("  \"%s\"\n", events[0].WhisperText)
	}
	fmt.Println("🕯️ =============================================")
}

// runDemo replays synthetic shell events demonstrating every trigger type
func runDemo(args []string) {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	debug := fs.Bool("debug", false, "Display whisper provenance debug overlay")
	personaFlag := fs.String("persona", "victorian_ghost", "Persona to demonstrate, or 'all' for montage")
	fx := fs.Bool("fx", false, "Ghost typing animation (capped at 300ms per whisper)")
	plain := fs.Bool("plain", false, "Disable color and ANSI escapes")
	_ = fs.Parse(args)

	cfg := config.DefaultConfig()
	reg, _ := personas.NewRegistry(filepath.Join(cfg.DataDir, "personas"))
	colorLevel := theme.DetectColorLevel(*plain)

	// Demo scenarios
	demoEvents := []struct {
		persona string
		cmd     string
		overlay string
		whisper string
	}{
		{
			persona: "victorian_ghost",
			cmd:     "git status",
			overlay: " (gemma3:4b, 420ms)",
			whisper: "A cold draught passes through the terminal; your files remain unsettled.",
		},
		{
			persona: "sysadmin",
			cmd:     "rm -rf ./node_modules",
			overlay: " (gemma3:4b, 510ms)",
			whisper: "That command just gave your filesystem a viking funeral. Hope backups exist.",
		},
		{
			persona: "cobol",
			cmd:     "go test ./...",
			overlay: " (gemma3:4b, 612ms)",
			whisper: "ABEND S0C4 IN PARAGRAPH-01. RE-READ 80-COLUMN SPECIFICATION BEFORE RETRYING.",
		},
		{
			persona: "marigold",
			cmd:     "git commit -m 'fix: urgent patch at 3am'",
			overlay: " (gemma3:4b, 490ms)",
			whisper: "The altar accepts your effort; breathing softly brings clearer thoughts to the craft.",
		},
		{
			persona: "victorian_ghost",
			cmd:     "docker compose up -d",
			overlay: " (grimoire)",
			whisper: "Cobwebs and dust welcome you back to this long-forsaken chamber.",
		},
	}

	// Filter or montage
	selectedPersona := strings.ToLower(strings.TrimSpace(*personaFlag))

	// Themed Header
	fmt.Println()
	if *plain {
		fmt.Printf("=== 👻 WRAITH SPECTRAL REPLAY (Mode: %s) ===\n", selectedPersona)
	} else {
		fmt.Printf("\033[38;2;179;157;219;1m=== 👻 WRAITH SPECTRAL REPLAY (Mode: %s) ===\033[0m\n", selectedPersona)
	}
	time.Sleep(200 * time.Millisecond)

	for _, d := range demoEvents {
		activePName := d.persona
		if selectedPersona != "all" && selectedPersona != "" {
			activePName = selectedPersona
		}

		p, ok := reg.Get(activePName)
		if !ok {
			p, _ = reg.Get("victorian_ghost")
		}

		// Prompt line
		fmt.Printf("\n> %s\n", d.cmd)
		time.Sleep(150 * time.Millisecond)

		overlay := ""
		if *debug {
			overlay = d.overlay
		}

		// Style whisper with persona skin
		styled := theme.StyleWhisper(d.whisper, p.Skin, colorLevel)

		// Optional Ghost-typing effect (Strictly capped at 300ms max per whisper)
		if *fx && colorLevel != theme.ColorNone {
			printGhostTyping(styled, 300*time.Millisecond)
		} else {
			fmt.Print(styled)
		}

		if overlay != "" {
			fmt.Printf(" %s", overlay)
		}
		fmt.Println()
		time.Sleep(250 * time.Millisecond)
	}

	// Themed Footer
	fmt.Println()
	if *plain {
		fmt.Println("=== 🕯️  Replay complete. The parlor quiets. ===")
	} else {
		fmt.Println("\033[38;2;179;157;219;3m=== 🕯️  Replay complete. The parlor quiets. ===\033[0m")
	}
}

// printGhostTyping outputs text rune-by-rune within a hard maximum budget of 300ms
func printGhostTyping(text string, maxDuration time.Duration) {
	runes := []rune(text)
	if len(runes) == 0 {
		return
	}

	// Calculate delay per rune so total time <= maxDuration
	delay := maxDuration / time.Duration(len(runes))
	if delay > 12*time.Millisecond {
		delay = 12 * time.Millisecond
	}

	inEscape := false
	for _, r := range runes {
		if r == '\033' {
			inEscape = true
		}

		fmt.Printf("%c", r)

		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}

		time.Sleep(delay)
	}
}
