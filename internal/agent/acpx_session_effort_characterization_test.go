package agent

import (
	"context"
	"testing"
)

// Suite: OpenCode Agent Session effort characterization
// Invariant: OpenCode uses the generic effort key and applies a deferred effort after setup but before work.
// Boundary IN: ACPXRunner Run lifecycle and the acpx command sequence it executes.
// Boundary OUT: profile configuration parsing and adapter queue-owner lifetime.

func TestACPXSessionEffortCharacterizationDeclaredBreakReasoningKeyMapsOpenCode(t *testing.T) {
	t.Parallel()

	key, err := acpxReasoningEffortConfigKey(RuntimeSpec{ID: "opencode"})
	if err != nil {
		t.Fatalf("reasoning key for OpenCode: %v", err)
	}
	if key != acpxGenericReasoningEffortKey {
		t.Fatalf("reasoning key = %q, want %q", key, acpxGenericReasoningEffortKey)
	}
}

func TestACPXSessionEffortCharacterizationDeclaredBreakWarmSessionOrdersEffortBeforeWork(t *testing.T) {
	t.Parallel()

	harness := newFakeACPXHarness(t)
	const model = "openrouter/x-ai/grok-4.5"
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": sessionCapabilitySnapshotFixture(
			t,
			model,
			[]string{model},
			"effort",
			"low",
			[]string{"low", "medium", "high"},
		),
	}))

	_, err := harness.run(
		context.Background(),
		RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: model, ReasoningEffort: "high"},
		"roundfix-run-opencode-effort-characterization",
	)
	if err != nil {
		t.Fatalf("run OpenCode Agent Session: %v", err)
	}

	invocations := readJSONInvocations(t, harness.invocationsPath)
	commands := make([]string, 0, 4)
	for _, invocation := range invocations {
		command := fakeACPXCommandKey(invocation)
		switch command {
		case "sessions ensure", "prompt", "set effort":
			commands = append(commands, command)
		}
	}
	wantCommands := []string{"sessions ensure", "prompt", "set effort", "prompt"}
	if len(commands) != len(wantCommands) {
		t.Fatalf("command sequence = %#v, want %#v; invocations: %#v", commands, wantCommands, invocations)
	}
	for index := range wantCommands {
		if commands[index] != wantCommands[index] {
			t.Fatalf("command sequence = %#v, want %#v; invocations: %#v", commands, wantCommands, invocations)
		}
	}
	prompts := readJSONStrings(t, harness.promptsPath)
	if len(prompts) != 2 || prompts[0] != acpxDeferredEffortWarmupPrompt || prompts[1] != "prompt" {
		t.Fatalf("prompt bytes = %#v, want setup %q then work prompt", prompts, acpxDeferredEffortWarmupPrompt)
	}
}
