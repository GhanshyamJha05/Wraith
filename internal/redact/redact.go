package redact

import (
	"regexp"
	"strings"
)

// Redactor handles stripping secrets, sensitive flags, and credentials
// Rule 6: Redact before anything is stored or sent to model.
// Skip haunting entirely for commands containing secret-like words.

var (
	// High-risk keywords that warrant completely skipping haunting
	skipKeywords = []string{
		"BEGIN RSA PRIVATE KEY",
		"BEGIN OPENSSH PRIVATE KEY",
		"BEGIN PRIVATE KEY",
		"id_rsa",
		"id_ed25519",
		".pgpass",
		".netrc",
		".npmrc",
	}

	// Secret tokens patterns (sk-..., ghp_..., gho_..., Bearer ..., AWS, etc.)
	tokenPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(bearer\s+)[a-zA-Z0-9_\-\.]{10,}`),
		regexp.MustCompile(`(sk-[a-zA-Z0-9]{20,})`),
		regexp.MustCompile(`(ghp_[a-zA-Z0-9]{20,})`),
		regexp.MustCompile(`(gho_[a-zA-Z0-9]{20,})`),
		regexp.MustCompile(`(glpat-[a-zA-Z0-9_\-]{20,})`),
		regexp.MustCompile(`(xox[baprs]-[a-zA-Z0-9_\-]{10,})`),
		regexp.MustCompile(`(AKIA[0-9A-Z]{16})`),
	}

	// URL with credentials pattern: https://user:pass@host
	urlCredsPattern = regexp.MustCompile(`([a-zA-Z]+://)([^:\s]+):([^@\s]+)@`)

	// Flag patterns: --password foo, -p bar, --token baz, --api-key secret
	flagPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(--(?:password|passwd|token|api[-_]?key|secret|auth[-_]?token)(?:\s+|=))([^\s]+)`),
		regexp.MustCompile(`(-p\s+)([^\s]+)`),
	}

	// Environment variable assignment: API_KEY=secret, TOKEN=xxx, PASSWORD=xxx
	envPatterns = regexp.MustCompile(`(?i)\b([A-Z0-9_]*(?:KEY|TOKEN|SECRET|PASSWORD|PASSWD|AUTH)[A-Z0-9_]*=)([^\s]+)`)
)

// Result represents the outcome of redaction
type Result struct {
	RedactedText string
	HasSecret    bool
	SkipHaunting bool
	Redactions   int
}

// Redact processes command text and strips sensitive content.
// It reports if the command contains secrets that should cause the whisper to be skipped entirely.
func Redact(cmd string) Result {
	if cmd == "" {
		return Result{RedactedText: "", HasSecret: false, SkipHaunting: false}
	}

	// 1. Check for fatal skip keywords
	for _, kw := range skipKeywords {
		if strings.Contains(cmd, kw) {
			return Result{
				RedactedText: "[REDACTED_COMMAND_WITH_HIGH_RISK_SECRET]",
				HasSecret:    true,
				SkipHaunting: true,
				Redactions:   1,
			}
		}
	}

	cleaned := cmd
	redactions := 0
	hasSecret := false

	// 2. Redact URLs with credentials
	if urlCredsPattern.MatchString(cleaned) {
		hasSecret = true
		redactions++
		cleaned = urlCredsPattern.ReplaceAllString(cleaned, `${1}[REDACTED_CREDS]@`)
	}

	// 3. Redact API tokens
	for _, pat := range tokenPatterns {
		if pat.MatchString(cleaned) {
			hasSecret = true
			redactions++
			cleaned = pat.ReplaceAllString(cleaned, "[REDACTED_TOKEN]")
		}
	}

	// 4. Redact command-line flags
	for _, pat := range flagPatterns {
		if pat.MatchString(cleaned) {
			hasSecret = true
			redactions++
			cleaned = pat.ReplaceAllString(cleaned, `${1}[REDACTED_FLAG]`)
		}
	}

	// 5. Redact environment assignments
	if envPatterns.MatchString(cleaned) {
		hasSecret = true
		redactions++
		cleaned = envPatterns.ReplaceAllString(cleaned, `${1}[REDACTED_ENV]`)
	}

	// If sensitive auth flags were found, mark SkipHaunting true for maximum safety
	skip := false
	lower := strings.ToLower(cmd)
	if strings.Contains(lower, "password") || strings.Contains(lower, "api_key") || strings.Contains(lower, "token") {
		if hasSecret {
			skip = true
		}
	}

	return Result{
		RedactedText: cleaned,
		HasSecret:    hasSecret,
		SkipHaunting: skip,
		Redactions:   redactions,
	}
}
