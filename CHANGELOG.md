# Changelog

All notable changes to Wraith will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.2.0] - 2026-10-28

### Added
- **Visual Theme & Skin Engine:** User-editable skins in persona YAMLs supporting TrueColor (`#RRGGBB`), 256-color, 16-color ANSI fallbacks, and plain ASCII mode.
- **Strict ANSI Isolation Rule:** Test suite verifying model output cannot inject terminal escape sequences or control codes (`theme_test.go`).
- **Windows Terminal Color Scheme:** Added `wraith theme` command and `docs/windows-terminal-scheme.json`.
- **System Doctor Command (`wraith doctor`):** Diagnostic verification for Ollama reachability, model pull state, shell hook installation, daemon status, data directory permissions, and exit-code self-test.
- **OS User Data Directory Integration:** Standardized storage using `%LOCALAPPDATA%\wraith` on Windows, `$XDG_DATA_HOME/wraith` on Linux, and `~/Library/Application Support/wraith` on macOS.
- **The Plain Helper Persona:** Non-spooky, objective error explainer strictly avoiding destructive command suggestions or false action claims.
- **Ghost Ceremony:** ASCII ghost banner on `wraith haunt`, quiet farewell on `wraith exorcise`, and themed status header on `wraith status`.
- **Demo Enhancements:** Multi-persona montage with `--persona all`, plain fallback with `--plain`, and `--fx` ghost-typing animation capped at 300ms.
- **GitHub Actions Workflows:** Multi-platform CI (Windows, Linux, macOS) running `go vet`, unit tests, and fuzzing; cross-compiled release pipeline.
- **Comprehensive Documentation:** Added `docs/ARCHITECTURE.md`, `docs/THREAT_MODEL.md`, `SECURITY.md`, and `CONTRIBUTING.md`.

---

## [0.1.0] - 2026-10-27

### Added
- Initial release of Wraith haunted terminal daemon and prompt hook.
- Deterministic trigger engine supporting failure streaks, dangerous commands (`rm -rf`), late-night coding, and stale directories.
- Zero-leak secret redaction engine for API keys (`sk-*`, `ghp_*`), credentials, and auth flags.
- Static fallback Grimoire ensuring zero prompt lag if Ollama times out.
- Four starter ghost personas: Victorian Ghost, Sysadmin, COBOL, and Marigold.
- Full opt-in transparency with diff preview on `wraith haunt`.

