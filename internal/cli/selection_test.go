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
	if claude.Model != "opus" || claude.ReasoningEffort != "" {
		t.Fatalf("expected Claude opus with model-managed reasoning, got %#v", claude)
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

func TestResolveSelectionRejectsExplicitEmptyModel(t *testing.T) {
	defaults := roundconfig.RuntimeDefaults{
		Model:           "configured-model",
		ReasoningEffort: "configured-reasoning",
	}

	_, err := ResolveSelection("codex", defaults, InvocationSelection{ModelSet: true})
	if err == nil {
		t.Fatal("expected explicit empty model to fail")
	}
	if !strings.Contains(err.Error(), "model") {
		t.Fatalf("expected error containing model, got %q", err.Error())
	}
}

func TestResolveSelectionAllowsEmptyReasoningEffort(t *testing.T) {
	defaults := roundconfig.RuntimeDefaults{Model: "configured-model"}

	got, err := ResolveSelection("codex", defaults, InvocationSelection{})
	if err != nil {
		t.Fatalf("expected configured empty reasoning effort to pass, got %v", err)
	}
	if got.Model != "configured-model" || got.ReasoningEffort != "" {
		t.Fatalf("expected configured model with empty reasoning effort, got %#v", got)
	}

	got, err = ResolveSelection("codex", roundconfig.RuntimeDefaults{
		Model:           "configured-model",
		ReasoningEffort: "configured-reasoning",
	}, InvocationSelection{ReasoningEffortSet: true})
	if err != nil {
		t.Fatalf("expected explicit empty reasoning effort to pass, got %v", err)
	}
	if got.Model != "configured-model" || got.ReasoningEffort != "" {
		t.Fatalf("expected explicit empty reasoning effort to override config, got %#v", got)
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
