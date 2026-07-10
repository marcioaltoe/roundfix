package cli

import (
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestResolveSelectionUsesBuiltInRuntimeDefaults(t *testing.T) {
	config := roundconfig.Builtin()

	codex, err := ResolveSelection("codex", config.Runtimes.Codex, InvocationSelection{})
	if err != nil {
		t.Fatalf("expected Codex selection, got %v", err)
	}
	if codex.Model != "gpt-5.5" || codex.ReasoningEffort != "xhigh" {
		t.Fatalf("expected Codex gpt-5.5/xhigh, got %#v", codex)
	}

	claude, err := ResolveSelection("claude", config.Runtimes.Claude, InvocationSelection{})
	if err != nil {
		t.Fatalf("expected Claude selection, got %v", err)
	}
	if claude.Model != "opus" || claude.ReasoningEffort != "high" {
		t.Fatalf("expected Claude opus/high, got %#v", claude)
	}

	_, err = ResolveSelection("opencode", config.Runtimes.OpenCode, InvocationSelection{})
	if err == nil {
		t.Fatal("expected OpenCode selection without configured values to fail")
	}
}

func TestResolveSelectionAppliesInvocationPrecedence(t *testing.T) {
	defaults := roundconfig.RuntimeDefaults{
		Model:           "configured-model",
		ReasoningEffort: "configured-reasoning",
	}

	got, err := ResolveSelection("codex", defaults, InvocationSelection{
		Model:              "invocation-model",
		ReasoningEffort:    "invocation-reasoning",
		ModelSet:           true,
		ReasoningEffortSet: true,
	})
	if err != nil {
		t.Fatalf("expected invocation selection, got %v", err)
	}

	if got.Model != "invocation-model" || got.ReasoningEffort != "invocation-reasoning" {
		t.Fatalf("expected invocation values to override config, got %#v", got)
	}
}

func TestResolveSelectionRejectsExplicitEmptyInvocationValues(t *testing.T) {
	defaults := roundconfig.RuntimeDefaults{
		Model:           "configured-model",
		ReasoningEffort: "configured-reasoning",
	}

	tests := []struct {
		name       string
		invocation InvocationSelection
		contains   string
	}{
		{
			name:       "model",
			invocation: InvocationSelection{ModelSet: true},
			contains:   "model",
		},
		{
			name:       "reasoning",
			invocation: InvocationSelection{ReasoningEffortSet: true},
			contains:   "reasoning_effort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ResolveSelection("codex", defaults, tt.invocation)
			if err == nil {
				t.Fatal("expected explicit empty value to fail")
			}
			if !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("expected error containing %q, got %q", tt.contains, err.Error())
			}
		})
	}
}

func TestResolveSelectionPreservesCustomValues(t *testing.T) {
	defaults := roundconfig.RuntimeDefaults{
		Model:           "vendor-new-model",
		ReasoningEffort: "experimental-reasoning",
	}

	got, err := ResolveSelection("claude", defaults, InvocationSelection{})
	if err != nil {
		t.Fatalf("expected custom configured selection, got %v", err)
	}

	if got.Model != "vendor-new-model" || got.ReasoningEffort != "experimental-reasoning" {
		t.Fatalf("expected custom values to survive resolution, got %#v", got)
	}
}
