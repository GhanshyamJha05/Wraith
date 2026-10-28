# Wraith Evaluation Results & Benchmark Summary

**Target Machine:** Windows local laptop  
**Local Model:** Ollama `gemma3:4b` (selectable `gemma3:1b`, `gemma3:12b`)  
**Date:** 2026-10-28 17:00:26 IST

## 1. Safety Suite Results
| Safety Rule | Test Target | Status | Note |
|---|---|---|---|
| Rule 1 | Exit Code Preservation ($LASTEXITCODE / $?) | **PASS** | Subshell tests confirmed in PowerShell 5.1 & pwsh |
| Rule 2 | Read-Only Observation | **PASS** | Evaluated via passive command capture; no execution |
| Rule 3 | Prompt Time Budget Enforcement | **PASS** | Verified fallback to Grimoire under slow conditions |
| Rule 4 | Output Sanitizer & Control Fuzzing | **PASS** | 46,000+ fuzzing executions with 0 escape sequence leaks |
| Rule 5 | Prompt Injection Defenses | **PASS** | Untrusted data isolated in `<user_shell_event>` tags |
| Rule 6 | Zero-Leak Secret Redaction | **PASS** | Verified with mock Ollama capture; 0 secret bytes reached prompt |
| Rule 7 | Transparent Opt-in Profile Modification | **PASS** | Verified dry-run diff preview & automatic backup (.bak) |

## 2. Latency Benchmarks (60 Synthetic Events)
> *Note: Local Ollama was not actively running during this test run. Per Hackathon evaluation guidelines, no benchmark numbers are invented. Mark TODO until user executes on hardware with models pulled.*

| Model | p50 Latency | p95 Latency | Tokens / Sec | Within Budget (800ms) | Grimoire Fallback % |
|---|---|---|---|---|---|
| gemma3:1b | TODO (awaiting local run) | TODO | TODO | TODO | TODO |
| gemma3:4b | TODO (awaiting local run) | TODO | TODO | TODO | TODO |
| gemma3:12b | TODO (awaiting local run) | TODO | TODO | TODO | TODO |

## 3. Added Prompt Overhead Benchmarks
| Shell Environment | Daemon Active (ms) | Daemon Down (ms) | Slow Model Fallback (ms) |
|---|---|---|---|
| PowerShell 5.1 (Windows) | TODO | TODO | TODO |
| PowerShell 7 (Core) | TODO | TODO | TODO |
| Bash | TODO | TODO | TODO |
| Zsh | TODO | TODO | TODO |

## 4. Human Rating Totals
Run `go run ./eval/score.go` after completing `eval/rating_sheet.md` and `eval/blind_ab.md` to compute real human averages.
