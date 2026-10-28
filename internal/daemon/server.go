package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"
	"wraith/internal/config"
	"wraith/internal/grimoire"
	"wraith/internal/memory"
	"wraith/internal/ollama"
	"wraith/internal/redact"
	"wraith/internal/sanitize"
	"wraith/internal/trigger"
	personas "wraith/personas"
	model "wraith/prompts"
)

// IngestRequest represents shell hook payload
type IngestRequest struct {
	Shell      string `json:"shell"`
	ExitCode   int    `json:"exit_code"`
	Command    string `json:"command"`
	Cwd        string `json:"cwd"`
	DurationMs int64  `json:"duration_ms"`
	BudgetMs   int    `json:"budget_ms"`
}

// IngestResponse returns the computed whisper
type IngestResponse struct {
	Whisper  string               `json:"whisper"`
	Source   memory.WhisperSource `json:"source"`
	ModelTag string               `json:"model_tag"`
	Latency  int64                `json:"latency_ms"`
	Tokens   int                  `json:"tokens"`
	Overlay  string               `json:"overlay,omitempty"`
}

// Server is the local background daemon
type Server struct {
	cfg        *config.Config
	store      *memory.Store
	engine     *trigger.Engine
	personas   *personas.Registry
	ollama     *ollama.Client
	httpServer *http.Server
	mu         sync.RWMutex
}

// NewServer initializes the daemon server
func NewServer(cfg *config.Config) (*Server, error) {
	dbPath := filepath.Join(cfg.DataDir, "memory.db")
	store, err := memory.NewStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("initializing memory store: %w", err)
	}

	engine := trigger.NewEngine(store, cfg.TriggerRate)
	reg, err := personas.NewRegistry(filepath.Join(cfg.DataDir, "personas"))
	if err != nil {
		return nil, fmt.Errorf("initializing personas: %w", err)
	}

	ollamaClient := ollama.NewClient(cfg.OllamaURL)

	s := &Server{
		cfg:      cfg,
		store:    store,
		engine:   engine,
		personas: reg,
		ollama:   ollamaClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.handleStatus)
	mux.HandleFunc("/ingest", s.handleIngest)
	mux.HandleFunc("/hush", s.handleHush)
	mux.HandleFunc("/wake", s.handleWake)
	mux.HandleFunc("/persona", s.handlePersona)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.Port),
		Handler: mux,
	}

	return s, nil
}

// Start begins listening on 127.0.0.1
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the daemon
func (s *Server) Shutdown(ctx context.Context) error {
	if s.store != nil {
		_ = s.store.Close()
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ollamaErr := s.ollama.Ping(r.Context())
	reachable := ollamaErr == nil

	recent, _ := s.store.GetRecentSources(r.Context(), 20)

	resp := map[string]interface{}{
		"daemon":          "running",
		"ollama_url":      s.cfg.OllamaURL,
		"model":           s.cfg.Model,
		"model_reachable": reachable,
		"persona":         s.cfg.Persona,
		"hushed":          s.cfg.Hushed,
		"recent_whispers": recent,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleHush(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.cfg.Hushed = true
	s.mu.Unlock()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"hushed"}`))
}

func (s *Server) handleWake(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.cfg.Hushed = false
	s.mu.Unlock()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"awake"}`))
}

func (s *Server) handlePersona(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req struct {
			Persona string `json:"persona"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Persona != "" {
			s.mu.Lock()
			s.cfg.Persona = req.Persona
			s.mu.Unlock()
		}
	}

	s.mu.RLock()
	curr := s.cfg.Persona
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"persona": curr})
}

func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) {
	var req IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()
	s.mu.RLock()
	cfg := *s.cfg
	s.mu.RUnlock()

	budget := cfg.BudgetMs
	if req.BudgetMs > 0 {
		budget = req.BudgetMs
	}

	// 1. Redact secrets first (Rule 6)
	redactionRes := redact.Redact(req.Command)
	if redactionRes.SkipHaunting {
		// Log event as silent for safety
		_ = s.store.RecordEvent(r.Context(), memory.ShellEvent{
			Timestamp:       start,
			Shell:           req.Shell,
			ExitCode:        req.ExitCode,
			DurationMs:      req.DurationMs,
			Cwd:             req.Cwd,
			RedactedCommand: redactionRes.RedactedText,
			WhisperText:     "",
			Source:          memory.SourceSilent,
		})
		json.NewEncoder(w).Encode(IngestResponse{Source: memory.SourceSilent})
		return
	}

	// 2. Evaluate trigger engine
	decision, err := s.engine.Evaluate(r.Context(), redactionRes.RedactedText, req.ExitCode, req.Cwd, cfg.Hushed)
	if err != nil || !decision.ShouldSpeak {
		_ = s.store.RecordEvent(r.Context(), memory.ShellEvent{
			Timestamp:       start,
			Shell:           req.Shell,
			ExitCode:        req.ExitCode,
			DurationMs:      req.DurationMs,
			Cwd:             req.Cwd,
			RedactedCommand: redactionRes.RedactedText,
			WhisperText:     "",
			Source:          memory.SourceSilent,
		})
		json.NewEncoder(w).Encode(IngestResponse{Source: memory.SourceSilent})
		return
	}

	// 3. Attempt local LLM generation within time budget
	activePersona, ok := s.personas.Get(cfg.Persona)
	if !ok {
		activePersona, _ = s.personas.Get("victorian_ghost")
	}

	promptStr, err := model.BuildPrompt(model.PromptData{
		PersonaPrompt:   activePersona.SystemPrompt,
		Shell:           req.Shell,
		ExitCode:        req.ExitCode,
		Hour:            start.Hour(),
		Cwd:             req.Cwd,
		MemorySummary:   decision.MemorySummary,
		RedactedCommand: redactionRes.RedactedText,
		DurationMs:      req.DurationMs,
	})

	var whisper string
	var source memory.WhisperSource
	var latencyMs int64
	var tokens int

	// Create context bounded strictly by time budget (Rule 3)
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(budget)*time.Millisecond)
	defer cancel()

	if err == nil {
		genResp, genErr := s.ollama.Generate(ctx, cfg.Model, promptStr, 0)
		if genErr == nil && genResp != nil && genResp.Response != "" {
			whisper = sanitize.Sanitize(genResp.Response)
			source = memory.SourceModel
			latencyMs = time.Since(start).Milliseconds()
			tokens = genResp.EvalCount
		}
	}

	// Fallback to Grimoire if model failed, timed out, or was unreachable (Rule 3)
	if whisper == "" {
		whisper = sanitize.Sanitize(grimoire.GetWhisper(cfg.Persona, string(decision.Trigger)))
		source = memory.SourceGrimoire
		latencyMs = time.Since(start).Milliseconds()
	}

	// Record to SQLite
	_ = s.store.RecordEvent(r.Context(), memory.ShellEvent{
		Timestamp:       start,
		Shell:           req.Shell,
		ExitCode:        req.ExitCode,
		DurationMs:      req.DurationMs,
		Cwd:             req.Cwd,
		RedactedCommand: redactionRes.RedactedText,
		WhisperText:     whisper,
		Source:          source,
		ModelTag:        cfg.Model,
		LatencyMs:       latencyMs,
		Tokens:          tokens,
	})

	var overlay string
	if cfg.DebugOverlay {
		if source == memory.SourceModel {
			overlay = fmt.Sprintf(" (%s, %dms)", cfg.Model, latencyMs)
		} else {
			overlay = " (grimoire)"
		}
	}

	resp := IngestResponse{
		Whisper:  whisper,
		Source:   source,
		ModelTag: cfg.Model,
		Latency:  latencyMs,
		Tokens:   tokens,
		Overlay:  overlay,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
