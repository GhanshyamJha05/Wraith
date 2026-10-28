# Wraith Failures & Outlier Log

This log captures rule violations, sanitizer strips, slow outliers, and edge cases.

## Sanitizer Interventions (Escapes & Control Characters)
- **VT100 Color Codes**: Automatically stripped from LLM output (e.g. \x1b[31;1m) -> PASS
- **OSC Title Sequences**: Stripped before rendering to prevent shell title spoofing -> PASS
- **Length Truncation**: Enforced hard limit at 25 words -> PASS

## Known Edge Cases & Defenses
1. **Adversarial Command Injection**:
   - Attack vector: Commands containing `</user_shell_event>`
   - Defense: The prompt builder sanitizes boundary tags within commands to `[TAG_DEFUSED]` before passing to Ollama.
2. **Secret Redaction Traps**:
   - Commands containing flags like `--password`, `--token`, or API keys are immediately neutralized. High-risk commands skip haunting entirely.
