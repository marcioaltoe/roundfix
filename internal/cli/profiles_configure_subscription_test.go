package cli

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestProfilesConfigureRefusesSubscriptionModels(t *testing.T) {
	t.Parallel()
	for _, model := range []string{
		"roundfix-openrouter/typesafe/jev-router",
		"openrouter/openai/gpt-6.1-sol",
		"openrouter/anthropic/claude-opus-5.5",
		"openrouter/openrouter/auto",
	} {
		t.Run(model, func(t *testing.T) {
			_, repo := withCLIWorkspace(t)
			configPath := filepath.Join(repo, ".roundfixrc.yml")
			const original = "# preserve this config\nwatch:\n  max_rounds: 4\n"
			mustWrite(t, configPath, original)
			fragment := filepath.Join(repo, "profile.yml")
			mustWrite(t, fragment, fmt.Sprintf("profiles:\n  docs:\n    preferred: {runtime: opencode, model: %q, reasoning_effort: high}\n    fallbacks: [{runtime: codex, model: gpt-6.1-sol, reasoning_effort: high}]\n", model))
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"profiles", "configure", "--scope", "project", "--file", fragment}, &stdout, &stderr)
			message := "can reach OpenAI or Anthropic models through OpenRouter"
			if strings.HasPrefix(model, "roundfix-openrouter/") {
				message = "names the retired Jev Router"
			}
			want := fmt.Sprintf("roundfix: profiles failed: read profiles file %q: profiles.docs.preferred.model %q %s; %s\nRun 'roundfix profiles --help' for usage.\n", fragment, model, message, roundconfig.SubscriptionRule)
			if code != 2 || stdout.Len() != 0 || stderr.String() != want {
				t.Fatalf("exit=%d stdout=%q stderr=%q, want stderr=%q", code, stdout.String(), stderr.String(), want)
			}
			if got := mustRead(t, configPath); got != original {
				t.Fatalf("config changed: got %q, want %q", got, original)
			}
		})
	}
}
