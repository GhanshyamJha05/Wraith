package memory

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// WhisperSource indicates where a whisper originated
type WhisperSource string

const (
	SourceModel    WhisperSource = "model"
	SourceGrimoire WhisperSource = "grimoire"
	SourceSilent   WhisperSource = "silent"
)

// ShellEvent represents an observed command execution
type ShellEvent struct {
	ID              int64
	Timestamp       time.Time
	Shell           string
	ExitCode        int
	DurationMs      int64
	Cwd             string
	RedactedCommand string
	WhisperText     string
	Source          WhisperSource
	ModelTag        string
	LatencyMs       int64
	Tokens          int
}

// Store handles persistent memory for Wraith in SQLite
type Store struct {
	db *sql.DB
}

// NewStore initializes the SQLite database
func NewStore(dbPath string) (*Store, error) {
	if dbPath == "" {
		home, _ := os.UserHomeDir()
		dbPath = filepath.Join(home, ".wraith", "memory.db")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("creating db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite: %w", err)
	}

	// Schema setup
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		shell TEXT NOT NULL,
		exit_code INTEGER NOT NULL,
		duration_ms INTEGER NOT NULL,
		cwd TEXT NOT NULL,
		redacted_command TEXT NOT NULL,
		whisper_text TEXT NOT NULL,
		source TEXT NOT NULL,
		model_tag TEXT,
		latency_ms INTEGER,
		tokens INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_events_cwd ON events(cwd);
	`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database handle
func (s *Store) Close() error {
	return s.db.Close()
}

// RecordEvent writes an event and its whisper outcome to SQLite
func (s *Store) RecordEvent(ctx context.Context, e ShellEvent) error {
	query := `
	INSERT INTO events (
		timestamp, shell, exit_code, duration_ms, cwd, redacted_command,
		whisper_text, source, model_tag, latency_ms, tokens
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query,
		e.Timestamp.UTC(), e.Shell, e.ExitCode, e.DurationMs, e.Cwd, e.RedactedCommand,
		e.WhisperText, string(e.Source), e.ModelTag, e.LatencyMs, e.Tokens,
	)
	return err
}

// GetRecentSources returns the sources of the last N whispers
func (s *Store) GetRecentSources(ctx context.Context, limit int) ([]ShellEvent, error) {
	query := `
	SELECT id, timestamp, shell, exit_code, redacted_command, whisper_text, source, model_tag, latency_ms, tokens
	FROM events
	WHERE source != 'silent'
	ORDER BY id DESC
	LIMIT ?
	`
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []ShellEvent
	for rows.Next() {
		var e ShellEvent
		var ts time.Time
		if err := rows.Scan(&e.ID, &ts, &e.Shell, &e.ExitCode, &e.RedactedCommand, &e.WhisperText, &e.Source, &e.ModelTag, &e.LatencyMs, &e.Tokens); err != nil {
			return nil, err
		}
		e.Timestamp = ts
		events = append(events, e)
	}
	return events, nil
}

// GetCommandFailureStreak counts consecutive recent failures of identical/similar command
func (s *Store) GetCommandFailureStreak(ctx context.Context, cmd string) (int, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return 0, nil
	}

	query := `
	SELECT exit_code FROM events
	WHERE redacted_command = ?
	ORDER BY id DESC
	LIMIT 10
	`
	rows, err := s.db.QueryContext(ctx, query, cmd)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	streak := 0
	for rows.Next() {
		var exitCode int
		if err := rows.Scan(&exitCode); err != nil {
			return streak, err
		}
		if exitCode != 0 {
			streak++
		} else {
			break
		}
	}
	return streak, nil
}

// GetLastWhisperTime returns when Wraith last whispered
func (s *Store) GetLastWhisperTime(ctx context.Context) (time.Time, error) {
	query := `SELECT timestamp FROM events WHERE source != 'silent' ORDER BY id DESC LIMIT 1`
	row := s.db.QueryRowContext(ctx, query)
	var ts time.Time
	if err := row.Scan(&ts); err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}
	return ts, nil
}

// Purge removes all records from the database
func (s *Store) Purge(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM events; VACUUM;")
	return err
}
