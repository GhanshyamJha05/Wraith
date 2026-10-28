package main

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
	"wraith/internal/config"
	"wraith/internal/memory"
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
	debug := false
	for _, a := range args {
		if a == "--debug" {
			debug = true
		}
	}

	demoEvents := []struct {
		desc    string
		cmd     string
		exit    int
		trigger string
		overlay string
	}{
		{
			desc:    "Routine successful command",
			cmd:     "git status",
			exit:    0,
			trigger: "random",
			overlay: " (gemma3:4b, 420ms)",
		},
		{
			desc:    "Dangerous recursive delete",
			cmd:     "rm -rf ./node_modules",
			exit:    0,
			trigger: "danger",
			overlay: " (gemma3:4b, 510ms)",
		},
		{
			desc:    "Repeated failed test run",
			cmd:     "go test ./...",
			exit:    1,
			trigger: "failure_streak",
			overlay: " (gemma3:4b, 612ms)",
		},
		{
			desc:    "Late night git commit after midnight",
			cmd:     "git commit -m 'fix: urgent patch at 3am'",
			exit:    0,
			trigger: "late_night",
			overlay: " (gemma3:4b, 490ms)",
		},
		{
			desc:    "Model budget timeout fallback",
			cmd:     "docker compose up -d",
			exit:    0,
			trigger: "grimoire_fallback",
			overlay: " (grimoire)",
		},
	}

	whispersSample := []string{
		"A cold draught passes through the terminal; your files remain unsettled.",
		"Such reckless excision would terrify even the graveyard sexton.",
		"Four times you have summoned this error; it still refuses to leave.",
		"Midnight has passed, and your eyes carry the dim pallor of tombs.",
		"Cobwebs and dust welcome you back to this long-forsaken chamber.",
	}

	fmt.Println("🎬 [wraith demo] Simulating terminal session haunting...")
	time.Sleep(400 * time.Millisecond)

	for i, d := range demoEvents {
		fmt.Printf("\n> %s\n", d.cmd)
		time.Sleep(300 * time.Millisecond)
		tag := ""
		if debug {
			tag = d.overlay
		}
		fmt.Printf("\033[90m👻 %s%s\033[0m\n", whispersSample[i%len(whispersSample)], tag)
		time.Sleep(400 * time.Millisecond)
	}

	fmt.Println("\n✨ [wraith demo] Replay complete.")
}
