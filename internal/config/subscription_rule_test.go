package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestCheckSubscriptionRule(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, runtime, model, refusal string }{
		{"openai", "opencode", "openrouter/openai/gpt-6.1-sol", "open"},
		{"anthropic", "opencode", "openrouter/anthropic/claude-opus-5.5", "open"},
		{"custom", "opencode-custom", "openrouter/openai/gpt-6.1-sol", "open"},
		{"trimmed and uppercase model", " opencode-custom ", " OPENROUTER/ANTHROPIC/claude-opus-5.5 ", "open"},
		{"alias", "opencode", "openrouter/~openai/gpt-6.1-sol", "open"},
		{"author suffix", "opencode", "openrouter/anthropic:extended/claude", "open"},
		{"alias and suffix", "opencode", "openrouter/~OPENAI:extended/gpt", "open"},
		{"auto", "opencode", "openrouter/openrouter/auto", "open"},
		{"typesafe router", "opencode", "openrouter/typesafe/jev-router", "open"},
		{"preset", "opencode", "openrouter/@preset/x", "open"},
		{"preset alias suffix", "opencode", "openrouter/~@saved:free/x", "open"},
		{"retired", "opencode", "roundfix-openrouter/typesafe/jev-router", "retired"},
		{"retired any model", "opencode-custom", "ROUNDFIX-OPENROUTER/deepseek/other", "retired"},
		{"retired provider only", "opencode", "roundfix-openrouter", "retired"},
		{"deepseek", "opencode", "openrouter/deepseek/deepseek-v4-pro", ""},
		{"x-ai", "opencode", "openrouter/x-ai/grok-4.5", ""},
		{"unknown router author", "opencode", "openrouter/third-party/auto", ""},
		{"different author prefix", "opencode", "openrouter/openai-other/x", ""},
		{"openai provider", "opencode", "openai/gpt-6.1-sol", ""},
		{"anthropic provider", "opencode", "anthropic/claude-opus-5.5", ""},
		{"opencode provider", "opencode", "opencode/model", ""},
		{"opencode-go", "opencode", "opencode-go/model", ""},
		{"no author", "opencode", "openrouter", ""},
		{"empty model", "opencode", "", ""},
		{"runtime case retained", "OpenCode", "openrouter/openai/gpt", ""},
	}
	for _, runtime := range []string{"codex", "claude", "cursor", "codex-custom", "claude-custom", "cursor-custom", "", "other"} {
		for _, model := range []string{"openrouter/openai/gpt", "roundfix-openrouter/typesafe/jev-router"} {
			tests = append(tests, struct{ name, runtime, model, refusal string }{runtime + "/" + model, runtime, model, ""})
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckSubscriptionRule("profiles.backend.preferred.model", tt.runtime, tt.model)
			if tt.refusal == "" {
				if err != nil {
					t.Fatalf("allowed selection refused: %v", err)
				}
				return
			}
			wording := "can reach OpenAI or Anthropic models through OpenRouter"
			if tt.refusal == "retired" {
				wording = "names the retired Jev Router"
			}
			want := fmt.Sprintf("profiles.backend.preferred.model %q %s; %s", strings.TrimSpace(tt.model), wording, SubscriptionRule)
			if err == nil || err.Error() != want {
				t.Fatalf("error = %v, want %q", err, want)
			}
		})
	}
}
