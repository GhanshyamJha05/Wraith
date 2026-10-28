```
          .---.
         /     \
        | () () |     W R A I T H
         \  ^  /      The Haunted Terminal
          |||||       (runs 100% offline via Ollama)
          |||||
```

# 👻 WRAITH: Offline AI Terminal Companion

[![CI](https://github.com/GhanshyamJha05/Wraith/actions/workflows/ci.yml/badge.svg)](https://github.com/GhanshyamJha05/Wraith/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-purple.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/GhanshyamJha05/Wraith)](https://goreportcard.com/report/github.com/GhanshyamJha05/Wraith)

> **A single-binary Go CLI that observes your shell prompt and whispers witty, contextual, offline AI commentary after errors, late-night coding sessions, and dangerous operations.**

> *"Four times you have summoned this error. It still refuses to leave."*

Wraith runs **100% locally** via [Ollama](https://ollama.com). Every whisper is private and zero cents per token. Zero command strings or file paths ever leave your machine.

---

## 💻 Terminal Demonstration

```
=== 👻 WRAITH SPECTRAL REPLAY (Mode: all) ===

> git status
🕯️ “A cold draught passes through the terminal; your files remain unsettled.”

> rm -rf ./node_modules
👻 [root@beyond] That command just gave your filesystem a viking funeral. Hope backups exist.

> go test ./...
📟 ABEND S0C4 IN PARAGRAPH-01. RE-READ 80-COLUMN SPECIFICATION BEFORE RETRYING.

> git commit -m 'fix: urgent patch at 3am'
🌼 The altar accepts your effort; breathing softly brings clearer thoughts to the craft.

> docker compose up -d
🕯️ “Cobwebs and dust welcome you back to this long-forsaken chamber.”

=== 🕯️  Replay complete. The parlor quiets. ===
```

Replay the interactive showcase yourself:
```bash
wraith demo --persona all --fx
```

---

## 🔍 Verified Platforms & Shells

To ensure absolute safety and prevent shell corruption, Wraith explicitly tracks shell compatibility:

| Platform / Shell | Shell Version | Hook Compatibility | Exit-Code Test | Status |
|---|:---:|:---:|:---:|:---:|
| **Windows 11 (x64)** | Windows PowerShell 5.1 | Supported (`Microsoft.PowerShell_profile.ps1`) | **PASSED** (`exit 42` preserved) | :white_check_mark: **Verified** |
| **Windows 11 (x64)** | PowerShell Core 7.x (`pwsh`) | Supported (`Microsoft.PowerShell_profile.ps1`) | **PASSED** (`exit 42` preserved) | :white_check_mark: **Verified** |
| **Ubuntu / Linux (x64)** | Bash 5.x | Supported (`.bashrc` precmd) | Tested in CI | :white_check_mark: **Verified** |
| **macOS (arm64 / x64)** | Zsh 5.x | Supported (`.zshrc` precmd) | Tested in CI | :white_check_mark: **Verified** |
| **Fish Shell** | Any | Untested | Untested | ⚠️ *Untested* |
| **Windows CMD (`cmd.exe`)** | Any | Unsupported (No prompt hook mechanism) | N/A | ❌ *Unsupported* |

---

## 🎭 Persona Skin Gallery

Each persona in `personas/` is fully customized with user-editable YAML skins. Output automatically degrades across TrueColor (24-bit), 256-color, 16-color ANSI, or plain ASCII fallback if `NO_COLOR` or `--plain` is active.

| Persona | Marker | Prefix & Styling | Color Palette | Example Whisper |
|---|:---:|---|---|---|
| **The Victorian Ghost** (`victorian_ghost`) | `🕯️` (`[*]`) | *Italic*, wrapped in matching quotation marks | **Lavender** (`#B39DDB`, 256: `183`, 16: `95`) | `🕯️ “Four times you have summoned this error; it still refuses to leave.”` |
| **Grumpy Dead Sysadmin** (`sysadmin`) | `👻` (`[sys]`) | `[root@beyond]` prefix | **Phosphor Green** (`#00FF66`, 256: `48`, 16: `92`) | `👻 [root@beyond] That command just gave your filesystem a viking funeral. Hope backups exist.` |
| **Ghost of COBOL** (`cobol`) | `📟` (`[COBOL]`) | **Bold ALL CAPS**, mainframe phrasing | **Amber** (`#FFB000`, 256: `214`, 16: `93`) | `📟 ABEND S0C4 IN PARAGRAPH-01. RE-READ 80-COLUMN SPECIFICATION BEFORE RETRYING.` |
| **Marigold** (`marigold`) | `🌼` (`[*]`) | *Italic*, tender & comforting phrasing | **Cempasúchil Orange** (`#FFA726`, 256: `215`, 16: `93`) | `🌼 Each stumble is only a petal laid upon the path of learning.` |
| **The Plain Helper** (`helper`) | `ℹ️` (`[info]`) | Neutral, non-spooky, plain language | **Cyan Slate** (`#9CCFD8`, 256: `116`, 16: `96`) | `ℹ️ The command exited with a non-zero status code, indicating an execution error.` |

Switch persona on the fly:
```bash
wraith persona sysadmin
# or switch to non-spooky practical helper:
wraith persona helper
```

---

## 🎨 Windows Terminal Color Scheme

Wraith ships with a custom **"Wraith Spectral"** palette for Windows Terminal:
```bash
wraith theme
```
JSON definition also available in [docs/windows-terminal-scheme.json](docs/windows-terminal-scheme.json).

---

## 📦 Installation & Quickstart

### 1. Download Pre-Compiled Binary
Download the pre-compiled binary for your architecture from [GitHub Releases](https://github.com/GhanshyamJha05/Wraith/releases):
- `wraith-windows-amd64.exe` (or `wraith-windows-arm64.exe`)
- `wraith-linux-amd64` (or `wraith-linux-arm64`)
- `wraith-darwin-amd64` (or `wraith-darwin-arm64`)

Or build from source:
```bash
git clone https://github.com/GhanshyamJha05/Wraith.git
cd Wraith
go build -o bin/wraith ./cmd/wraith
```

### 2. Verify Your Environment (`wraith doctor`)
Wraith includes a built-in doctor command to check all dependencies, data directories, and subshell exit codes:
```bash
wraith doctor
```

Sample output:
```
🩺 WRAITH SYSTEM DIAGNOSTIC (wraith doctor)
=================================================================
[1/6] User Data Location:
      Path: C:\Users\...\AppData\Local\wraith
      ✅ Accessible and writable

[2/6] Ollama Service:
      Target URL: http://127.0.0.1:11434
      ✅ Service is online and responding

[3/6] Model Availability (gemma3:4b):
      ✅ Model "gemma3:4b" is pulled and ready in memory

[4/6] Shell Profile Hook:
      ✅ Hook is installed in C:\Users\...\Microsoft.PowerShell_profile.ps1

[5/6] Background Daemon Status (127.0.0.1:23847):
      ✅ Daemon is actively listening and ready

[6/6] Exit Code Preservation Self-Test (Rule 1):
      ✅ Subshell preserved exit code 42 exactly
=================================================================
✨ Doctor inspection complete: Core subsystem is healthy!
```

### 3. Install the Ghost Hook (`wraith haunt`)
Wraith is completely transparent. It displays an ASCII ghost banner, previews the exact profile changes, and requests explicit confirmation:
```powershell
wraith haunt
```

### 4. Remove / Exorcise (`wraith exorcise`)
To remove Wraith cleanly with a quiet farewell:
```powershell
wraith exorcise
```

---

## 🛡️ The 7 Non-Negotiable Safety Rules

1. **Exit Status Preservation:** The prompt hook captures `$LASTEXITCODE`, `$?`, or `$status` and guarantees exact restoration. A failed command (`exit 42`) remains exit code 42. Verified by automated tests.
2. **Read-Only Observation:** Wraith strictly observes command strings and exit codes. It never executes, intercepts, alters, or reruns commands.
3. **Prompt Budget Guarantee:** Hard deadline (default 800 ms, configurable). If Ollama times out or is slow, Wraith either falls back to an instant pre-written Grimoire whisper or stays silent.
4. **Output Sanitization & Fuzz Testing:** All LLM outputs are stripped of ANSI color escapes, OSC title injections, and control bytes. Strictly enforces at most 1 line and 25 words. Fuzz-tested with continuous randomized VT100 exploit vectors. **Hard styling boundary:** ANSI colors and skins are applied strictly by Wraith AFTER model output is sanitized.
5. **Prompt Injection Boundary:** Command text is treated as untrusted data and quarantined inside `<user_shell_event>` boundary tags with defusal rules against break-outs.
6. **Zero-Leak Secret Redaction:** API keys (`sk-*`, `ghp_*`), passwords, auth flags (`--password`, `-p`, `--token`), and sensitive environment variables are sanitized before storage or LLM inference. High-risk commands skip haunting entirely.
7. **Transparent Opt-In:** Full diff preview and automated `.bak` backups before profile edits.

For full architectural rationale and formal threat analysis, see:
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- [docs/THREAT_MODEL.md](docs/THREAT_MODEL.md)

---

## 📂 Data Storage Locations

Wraith respects standard OS directory specifications and never stores database files in your working project folders:
- **Windows:** `%LOCALAPPDATA%\wraith\memory.db`
- **Linux:** `$XDG_DATA_HOME/wraith/memory.db` (or `~/.local/share/wraith/memory.db`)
- **macOS:** `~/Library/Application Support/wraith/memory.db`

Override with `WRAITH_DATA_DIR=/custom/path` if desired.

---

## 📜 License

Distributed under the [MIT License](LICENSE).
