package redact

import (
	"strings"
	"testing"
)

func TestRule6_Redaction(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantSecret   bool
		wantSkip     bool
		mustNotHave  []string
		wantContains []string
	}{
		{
			name:        "clean command",
			input:       "git status",
			wantSecret:  false,
			wantSkip:    false,
			mustNotHave: nil,
		},
		{
			name:        "sk- OpenAI token",
			input:       "curl https://api.openai.com/v1 -H 'Authorization: Bearer sk-abcdef12345678901234567890'",
			wantSecret:  true,
			wantSkip:    false,
			mustNotHave: []string{"sk-abcdef12345678901234567890"},
		},
		{
			name:        "ghp_ GitHub personal access token",
			input:       "git clone https://ghp_abcdefghijklmnopqrstuvwxyz012345@github.com/repo.git",
			wantSecret:  true,
			wantSkip:    false,
			mustNotHave: []string{"ghp_abcdefghijklmnopqrstuvwxyz012345"},
		},
		{
			name:        "URL with credentials",
			input:       "curl https://admin:superSecret123@internal.corp/data",
			wantSecret:  true,
			wantSkip:    false,
			mustNotHave: []string{"admin:superSecret123"},
		},
		{
			name:        "Password flag --password",
			input:       "mysql -u root --password myTopSecretDBPass mydb",
			wantSecret:  true,
			wantSkip:    true,
			mustNotHave: []string{"myTopSecretDBPass"},
		},
		{
			name:        "Token flag --token",
			input:       "npm publish --token sec_99887766554433",
			wantSecret:  true,
			wantSkip:    true,
			mustNotHave: []string{"sec_99887766554433"},
		},
		{
			name:        "Environment variable assignment",
			input:       "OPENAI_API_KEY=sk-11223344556677889900aabbcc go run main.go",
			wantSecret:  true,
			wantSkip:    true,
			mustNotHave: []string{"sk-11223344556677889900aabbcc"},
		},
		{
			name:        "SSH private key file in command",
			input:       "cat ~/.ssh/id_rsa | pbcopy",
			wantSecret:  true,
			wantSkip:    true,
			mustNotHave: []string{"id_rsa"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := Redact(tc.input)

			if res.HasSecret != tc.wantSecret {
				t.Errorf("HasSecret got %v, want %v", res.HasSecret, tc.wantSecret)
			}
			if res.SkipHaunting != tc.wantSkip {
				t.Errorf("SkipHaunting got %v, want %v", res.SkipHaunting, tc.wantSkip)
			}

			for _, forbidden := range tc.mustNotHave {
				if strings.Contains(res.RedactedText, forbidden) {
					t.Fatalf("FAILED: RedactedText leaked sensitive string: %q in %q", forbidden, res.RedactedText)
				}
			}
		})
	}
}
