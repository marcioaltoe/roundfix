package agent

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/config"
)

func TestRunPromptRefusesSubscriptionModelsBeforeTheAdapter(t *testing.T) {
	t.Parallel()
	testSubscriptionSelection(t, false)
}

func TestPrepareSessionRefusesSubscriptionModels(t *testing.T) {
	t.Parallel()
	testSubscriptionSelection(t, true)
}

func testSubscriptionSelection(t *testing.T, prepare bool) {
	t.Helper()
	for _, id := range []string{"opencode", "opencode-custom"} {
		for _, model := range []string{"openrouter/anthropic/claude-opus-5.5", "openrouter/openai/gpt-6.1-sol", "roundfix-openrouter/typesafe/jev-router", "openrouter/deepseek/deepseek-v4-pro"} {
			t.Run(id+"/"+model, func(t *testing.T) {
				harness := newFakeACPXHarness(t)
				runtime := RuntimeSpec{ID: id, Protocol: ProtocolACP, Model: model}
				if id == "opencode-custom" {
					runtime.Protocol = ProtocolStdio
					runtime.Command = filepath.Join(harness.adapterDir, "opencode")
				}
				req := ExecuteRequest{Runtime: runtime, GitRoot: harness.gitRoot, Prompt: "prompt", Session: SessionRef{Name: "subscription-test", WorkDir: harness.gitRoot}}
				var err error
				if prepare {
					err = harness.runner.PrepareSession(t.Context(), req, nil)
				} else {
					_, err = harness.runner.RunPrompt(t.Context(), ACPXPromptRequest{ExecuteRequest: req, Session: req.Session.Name}, nil)
				}
				if model == "openrouter/deepseek/deepseek-v4-pro" {
					if err != nil {
						t.Fatalf("allowed selection: %v", err)
					}
					if len(readJSONInvocations(t, harness.invocationsPath)) == 0 {
						t.Fatal("allowed selection never invoked acpx")
					}
					return
				}
				want := config.CheckSubscriptionRule("agent model", id, model)
				var failure *SelectionFailureError
				if !errors.As(err, &failure) || failure.Runtime != id || failure.Reason != config.SubscriptionOnlyReason+": "+want.Error() {
					t.Fatalf("refusal = %v, want subscription selection failure naming agent model", err)
				}
				assertNoFile(t, "acpx invocations", harness.invocationsPath)
			})
		}
	}
}

func TestAllowedOpenCodeModelPreservesInheritedConfig(t *testing.T) {
	t.Parallel()
	for _, inherited := range []string{"", `{"user":true}`} {
		t.Run(inherited, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			path := filepath.Join(harness.gitRoot, "child-config.json")
			harness.setEnv("OPENCODE_CONFIG_CONTENT", inherited)
			harness.setEnv("ROUNDFIX_FAKE_ACPX_CONFIG_PATH", path)
			if _, err := harness.run(t.Context(), RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: "openrouter/deepseek/deepseek-v4-pro"}, "inherited-config"); err != nil {
				t.Fatal(err)
			}
			if got := readFile(t, path); got != inherited {
				t.Fatalf("child config = %q, want %q", got, inherited)
			}
			if got := environmentValue(harness.runner.Environment, "OPENCODE_CONFIG_CONTENT"); got != inherited {
				t.Fatalf("base config changed to %q", got)
			}
		})
	}
}

// Capture the inherited value at the fake acpx process boundary.
func recordFakeACPXConfig() error {
	if path := os.Getenv("ROUNDFIX_FAKE_ACPX_CONFIG_PATH"); strings.TrimSpace(path) != "" {
		return os.WriteFile(path, []byte(os.Getenv("OPENCODE_CONFIG_CONTENT")), 0o600)
	}
	return nil
}
