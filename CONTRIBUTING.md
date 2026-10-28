# Contributing to Wraith

Thank you for your interest in contributing to Wraith!

---

## Code of Conduct & Ground Rules

1. **Non-Negotiable Safety Rules:** Every pull request must preserve the 7 Non-Negotiable Safety Rules detailed in the [README](README.md) and [docs/THREAT_MODEL.md](docs/THREAT_MODEL.md).
2. **Read-Only Non-Interference:** Wraith must never alter, delay, intercept, or execute commands on the user's behalf.
3. **Budget First:** Shell prompt budgets are sacred. No feature may block interactive prompt rendering beyond the configured budget (default 800 ms).
4. **Offline First:** Wraith never transmits telemetry, code, or prompts across the network. All inference occurs locally via Ollama.

---

## Development Setup

### Prerequisites
- Go 1.22+ (tested with Go 1.24)
- (Optional) [Ollama](https://ollama.com) running locally with `gemma3:4b` (`ollama serve`)

### Building and Testing
```bash
# Run tests
go test -v ./...

# Run fuzz testing against the output sanitizer
go test -fuzz=FuzzSanitizer -fuzztime=10s ./internal/sanitize

# Build standalone binary
go build -o bin/wraith ./cmd/wraith
```

---

## Submitting Pull Requests

1. Fork the repository and create a feature branch (`git checkout -b feat/my-improvement`).
2. Ensure `go vet ./...` passes without warnings.
3. Add unit tests for any new trigger, parser, or security boundary.
4. Open a pull request with a clear description of the motivation and verified shells/platforms.

