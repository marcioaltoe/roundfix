package config

import (
	"strings"
	"testing"
)

func TestVerificationFormatDefaultsEmpty(t *testing.T) {
	for _, tc := range []struct {
		name, content string
	}{
		{"unset", ""},
		{"empty", "verification:\n  format: \"\"\n"},
		{"whitespace", "verification:\n  format: \" \\t \\n\"\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := ResolveConfigProposal(nil, []byte(tc.content))
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Verification.Format != "" {
				t.Fatalf("format = %q, want empty", cfg.Verification.Format)
			}
		})
	}
}

func TestVerificationFormatProjectReplacesUser(t *testing.T) {
	user := []byte("verification:\n  format: '  user-format --write  '\n")
	for _, tc := range []struct {
		name, project, want string
	}{
		{"user retained", "", "user-format --write"},
		{"project replaces", "verification:\n  format: '  project-format --write  '\n", "project-format --write"},
		{"empty replaces", "verification:\n  format: ''\n", ""},
		{"whitespace replaces", "verification:\n  format: '   '\n", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := ResolveConfigProposal(user, []byte(tc.project))
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Verification.Format != tc.want {
				t.Fatalf("format = %q, want %q", cfg.Verification.Format, tc.want)
			}
		})
	}
}

func TestVerificationFormatRejectsANonString(t *testing.T) {
	for _, tc := range []struct {
		name, value string
	}{
		{"list", "[gofmt]"},
		{"mapping", "{command: gofmt}"},
		{"integer", "1"},
		{"boolean", "true"},
		{"null", "null"},
	} {
		for _, source := range []string{"user", "project"} {
			t.Run(source+"/"+tc.name, func(t *testing.T) {
				content := []byte("verification:\n  format: " + tc.value + "\n")
				var user, project []byte
				if source == "user" {
					user = content
				} else {
					project = content
				}
				_, err := ResolveConfigProposal(user, project)
				if err == nil || !strings.Contains(err.Error(), "verification.format must be a string") {
					t.Fatalf("invalid format error = %v", err)
				}
			})
		}
	}
}

func TestVerificationFormatRendersInTheDefaultConfig(t *testing.T) {
	content := DefaultConfigYAML()
	const entry = "  # Formats the files a QA step commits under the Spec's qa/ directory; empty runs nothing.\n  format: \"\"\n"
	start := strings.Index(content, "\nverification:\n")
	if start < 0 {
		t.Fatal("default config lacks verification mapping")
	}
	section := strings.SplitN(content[start+len("\nverification:\n"):], "\n\n", 2)[0]
	if !strings.Contains(section+"\n", entry) {
		t.Fatal("default verification mapping lacks empty format command and comment")
	}
	cfg, err := ResolveConfigProposal(nil, []byte(content))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Verification.Format != "" {
		t.Fatalf("rendered format = %q, want empty", cfg.Verification.Format)
	}
}
