package trigger

import (
	"context"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
	"wraith/internal/memory"
)

// TriggerType classifies what conditions triggered the whisper
type TriggerType string

const (
	TriggerNone           TriggerType = "none"
	TriggerFailureStreak  TriggerType = "failure_streak"
	TriggerDanger         TriggerType = "danger"
	TriggerLateNight      TriggerType = "late_night"
	TriggerStaleDirectory TriggerType = "stale_directory"
	TriggerGeneralFailure TriggerType = "general_failure"
	TriggerRandom         TriggerType = "random"
)

var dangerousCommands = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(rm\s+-[rf]{1,2}|remove-item\s+-recurse)\b`),
	regexp.MustCompile(`(?i)\b(git\s+push\s+.*--force|git\s+reset\s+--hard)\b`),
	regexp.MustCompile(`(?i)\b(sudo|chmod\s+777|chown\s+-R)\b`),
	regexp.MustCompile(`(?i)\b(drop\s+table|drop\s+database|truncate\s+table)\b`),
	regexp.MustCompile(`(?i)\b(format\s+[a-z]:|diskpart)\b`),
}

// Decision represents whether and why a whisper should be spoken
type Decision struct {
	ShouldSpeak   bool
	Trigger       TriggerType
	MemorySummary string
}

// Engine evaluates deterministic rules to trigger whispers
type Engine struct {
	Store       *memory.Store
	BaseChance  float64
	CooldownSec int
}

// NewEngine creates a new deterministic trigger engine
func NewEngine(store *memory.Store, baseChance float64) *Engine {
	if baseChance <= 0 {
		baseChance = 0.25
	}
	return &Engine{
		Store:       store,
		BaseChance:  baseChance,
		CooldownSec: 10, // Avoid back-to-back chatter
	}
}

// Evaluate determines if the ghost should speak
func (e *Engine) Evaluate(ctx context.Context, cmd string, exitCode int, cwd string, hushed bool) (Decision, error) {
	if hushed {
		return Decision{ShouldSpeak: false, Trigger: TriggerNone}, nil
	}

	now := time.Now()

	// 1. Check cooldown: never speak if we spoke within the cooldown window
	if e.Store != nil {
		lastWhisper, err := e.Store.GetLastWhisperTime(ctx)
		if err == nil && !lastWhisper.IsZero() {
			if time.Since(lastWhisper) < time.Duration(e.CooldownSec)*time.Second {
				return Decision{ShouldSpeak: false, Trigger: TriggerNone}, nil
			}
		}
	}

	// 2. Failure streak: repeated failure of the same command
	if exitCode != 0 && e.Store != nil && cmd != "" {
		streak, err := e.Store.GetCommandFailureStreak(ctx, cmd)
		if err == nil && streak >= 2 {
			return Decision{
				ShouldSpeak:   true,
				Trigger:       TriggerFailureStreak,
				MemorySummary: "The user has failed this exact command multiple times in a row.",
			}, nil
		}
	}

	// 3. Dangerous command detection (rm -rf, git push --force, etc.)
	for _, pattern := range dangerousCommands {
		if pattern.MatchString(cmd) {
			return Decision{
				ShouldSpeak:   true,
				Trigger:       TriggerDanger,
				MemorySummary: "A perilous and destructive command was attempted.",
			}, nil
		}
	}

	// 4. Stale directory check (untouched for 30+ days)
	if cwd != "" {
		if stat, err := os.Stat(cwd); err == nil {
			modTime := stat.ModTime()
			if time.Since(modTime) > 30*24*time.Hour {
				return Decision{
					ShouldSpeak:   true,
					Trigger:       TriggerStaleDirectory,
					MemorySummary: "This directory was slumbering, untouched for over a month.",
				}, nil
			}
		}
	}

	// 5. Late night hours (between midnight and 5:00 AM)
	hour := now.Hour()
	if hour >= 0 && hour < 5 {
		// Raised chance late at night (50%)
		if rand.Float64() < 0.50 {
			summary := "The hour is late, deep past midnight."
			if strings.HasPrefix(strings.TrimSpace(cmd), "git commit") {
				summary = "A git commit made in the dead of night between midnight and dawn."
			}
			return Decision{
				ShouldSpeak:   true,
				Trigger:       TriggerLateNight,
				MemorySummary: summary,
			}, nil
		}
	}

	// 6. Non-zero exit code (general failure)
	if exitCode != 0 {
		// Higher chance for errors (60%)
		if rand.Float64() < 0.60 {
			return Decision{
				ShouldSpeak:   true,
				Trigger:       TriggerGeneralFailure,
				MemorySummary: "The previous command ended in failure.",
			}, nil
		}
	}

	// 7. Base probabilistic chance
	if rand.Float64() < e.BaseChance {
		return Decision{
			ShouldSpeak:   true,
			Trigger:       TriggerRandom,
			MemorySummary: "",
		}, nil
	}

	return Decision{ShouldSpeak: false, Trigger: TriggerNone}, nil
}
