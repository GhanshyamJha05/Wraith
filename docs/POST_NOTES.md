# Hackathon Post Notes: Wraith ("Ghost in the Machine")

## 1. What We Built
Wraith is a small, shippable terminal companion in Go that turns everyday terminal errors and late-night coding into atmospheric, dryly humorous moments.

Key architectural highlights:
- **Zero Cloud Footprint:** Runs locally via Ollama (`gemma3:4b` default).
- **Single Cross-Compiled Go Binary:** Standard library focused, with only `modernc.org/sqlite` (pure Go, zero CGO) for storage.
- **7 Non-Negotiable Safety Rules:** Exit code preservation, read-only observation, 800ms prompt budget guarantee, ANSI escape stripping fuzzing, prompt injection boundary containment, and zero-leak secret redaction.

## 2. Evaluation & Benchmark Observations
All synthetic benchmark datasets and human rating sheets are stored under `eval/`:
- `eval/RESULTS.md`: Safety test suite pass/fail results.
- `eval/rating_sheet.md`: 30 pre-generated event evaluations with blank human scoring cells.
- `eval/blind_ab.md`: Blind comparison between model generation and static grimoire whispers.
- `eval/failures.md`: Documented edge cases, prompt injection defense mechanisms, and sanitizer rules.

*Note: As per project guidelines, all benchmark latency figures on unrun models remain marked `TODO` until executed on physical host hardware with active model downloads.*

## 3. How to Run Live Benchmarks
To execute live benchmarks against a running Ollama daemon:
```powershell
ollama serve
ollama pull gemma3:4b
.\tasks.ps1 eval
```
To calculate human rating averages after filling the rating sheets:
```powershell
go run .\eval\score.go
```
