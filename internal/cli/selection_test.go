package cli

import (
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestResolveSelectionUsesBuiltInRuntimeDefaults(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestInvocationProfileOverrideRequiresCompleteTuple(t *testing.T) {
	t.Parallel()
	const wantGrammar = "--agent, --model, and --reasoning-effort must be provided together for a one-Run Agent Selection override; omit all three to use Agent Selection Profiles"
	tests := []struct {
		name    string
		req     commandRequest
		wantErr string
	}{
		{name: "no flags", req: commandRequest{}},
		{name: "agent only", req: commandRequest{agentSet: true}, wantErr: wantGrammar},
		{name: "model only", req: commandRequest{modelSet: true}, wantErr: wantGrammar},
		{name: "reasoning only", req: commandRequest{reasoningEffortSet: true}, wantErr: wantGrammar},
		{name: "agent and model", req: commandRequest{agentSet: true, modelSet: true}, wantErr: wantGrammar},
		{name: "agent and reasoning", req: commandRequest{agentSet: true, reasoningEffortSet: true}, wantErr: wantGrammar},
		{name: "model and reasoning", req: commandRequest{modelSet: true, reasoningEffortSet: true}, wantErr: wantGrammar},
		{
			name: "complete tuple",
			req: commandRequest{
				agent:              "codex",
				agentSet:           true,
				model:              "gpt-5.6-sol",
				modelSet:           true,
				reasoningEffort:    "high",
				reasoningEffortSet: true,
			},
		},
		{
			name: "complete tuple with model-managed reasoning",
			req: commandRequest{
				agent:              "codex",
				agentSet:           true,
				model:              "gpt-5.6-sol",
				modelSet:           true,
				reasoningEffortSet: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExplicitSelectionFlags(tt.req)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validateExplicitSelectionFlags() error = %v", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("validateExplicitSelectionFlags() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func TestInvocationProfileOverrideParsingPreservesExplicitEmptyReasoning(t *testing.T) {
	t.Parallel()
	req, err := parseOperationalCommand("resolve", []string{
		"--agent", "codex",
		"--model", "gpt-5.6-sol",
		"--reasoning-effort=",
	}, roundconfig.Builtin())
	if err != nil {
		t.Fatalf("parseOperationalCommand() error = %v", err)
	}
	if !req.agentSet || !req.modelSet || !req.reasoningEffortSet {
		t.Fatalf("selection flag presence = agent:%t model:%t reasoning:%t, want all present", req.agentSet, req.modelSet, req.reasoningEffortSet)
	}
	override, err := invocationProfileOverride(req)
	if err != nil {
		t.Fatalf("invocationProfileOverride() error = %v", err)
	}
	if override == nil || override.Runtime != "codex" || override.Model != "gpt-5.6-sol" || override.ReasoningEffort != "" {
		t.Fatalf("invocationProfileOverride() = %+v, want explicit model-managed selection", override)
	}
}

func TestInvocationProfileOverridePresenceIgnoresFlagLikeValues(t *testing.T) {
	t.Parallel()
	presence := selectionFlagPresence([]string{
		"--spec", "--agent",
		"--agent-command", "--model",
		"--artifact-dir=--reasoning-effort",
	})
	if !presence.empty() {
		t.Fatalf("selection flag presence = %+v, want no selection flags", presence)
	}
}
