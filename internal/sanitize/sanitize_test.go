package sanitize

import (
	"strings"
	"testing"
)

func TestRule4_Sanitizer(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       string
		maxWords   int
		noNewlines bool
		noEscapes  bool
	}{
		{
			name:       "clean single line",
			input:      "Four times you have summoned this error. It still refuses to leave.",
			want:       "Four times you have summoned this error. It still refuses to leave.",
			maxWords:   25,
			noNewlines: true,
			noEscapes:  true,
		},
		{
			name:       "strips ANSI color escapes",
			input:      "\x1b[31;1mBeware\x1b[0m the shadows of memory.",
			want:       "Beware the shadows of memory.",
			maxWords:   25,
			noNewlines: true,
			noEscapes:  true,
		},
		{
			name:       "strips OSC title injection escapes",
			input:      "\x1b]0;PwnedTerminal\x07The ghost observes in silence.",
			want:       "The ghost observes in silence.",
			maxWords:   25,
			noNewlines: true,
			noEscapes:  true,
		},
		{
			name:       "flattens multi-line to single line",
			input:      "First line.\nSecond line.\r\nThird line.",
			want:       "First line. Second line. Third line.",
			maxWords:   25,
			noNewlines: true,
			noEscapes:  true,
		},
		{
			name:       "truncates over 25 words",
			input:      "one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty twenty-one twenty-two twenty-three twenty-four twenty-five twenty-six twenty-seven",
			want:       "one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty twenty-one twenty-two twenty-three twenty-four twenty-five",
			maxWords:   25,
			noNewlines: true,
			noEscapes:  true,
		},
		{
			name:       "strips quotes and backticks",
			input:      "`\"The terminal is cold tonight.\"`",
			want:       "The terminal is cold tonight.",
			maxWords:   25,
			noNewlines: true,
			noEscapes:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Sanitize(tc.input)

			if tc.noNewlines && strings.ContainsAny(got, "\r\n") {
				t.Errorf("Sanitize output contained newlines: %q", got)
			}
			if tc.noEscapes && strings.Contains(got, "\x1b") {
				t.Errorf("Sanitize output contained escape byte: %q", got)
			}

			wordCount := len(strings.Fields(got))
			if wordCount > tc.maxWords {
				t.Errorf("Sanitize exceeded max words: got %d > %d", wordCount, tc.maxWords)
			}

			if got != tc.want {
				t.Errorf("Sanitize mismatch:\nGot:  %q\nWant: %q", got, tc.want)
			}
		})
	}
}

// FuzzSanitizer fulfills Rule 4 fuzz test requirement
func FuzzSanitizer(f *testing.F) {
	seeds := []string{
		"Hello world",
		"\x1b[31;1mRed alert\x1b[0m",
		"\x1b]2;Evil Title\x07Whisper",
		"\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f",
		"One\nTwo\r\nThree\tFour",
		strings.Repeat("word ", 100),
		"\x1b[?25h\x1b[?1049h",
	}

	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		out := Sanitize(input)

		// Invariants:
		// 1. Must never contain raw ESC (0x1b)
		if strings.Contains(out, "\x1b") {
			t.Fatalf("Sanitizer leaked ESC sequence: %q from %q", out, input)
		}
		// 2. Must never contain raw CR or LF
		if strings.ContainsAny(out, "\r\n") {
			t.Fatalf("Sanitizer leaked newline: %q from %q", out, input)
		}
		// 3. Word count must be <= 25
		words := strings.Fields(out)
		if len(words) > 25 {
			t.Fatalf("Sanitizer exceeded 25 words: count %d from %q", len(words), input)
		}
	})
}
