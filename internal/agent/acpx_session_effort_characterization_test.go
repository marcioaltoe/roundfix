package agent

import (
	"context"
	"testing"
)

// Suite: OpenCode Agent Session effort characterization
// Invariant: OpenCode uses the generic effort key, while today's Run still reaches its work prompt without setting it.
// Boundary IN: ACPXRunner Run lifecycle and the acpx command sequence it executes.
// Boundary OUT: profile configuration parsing and future warm-up behavior.

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

func TestACPXSessionEffortCharacterizationTodayRunDoesNotSetReasoningBeforePrompt(t *testing.T) {
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
		RuntimeSpec{ID: "opencode", Protocol: ProtocolACP, Model: model},
		"roundfix-run-opencode-effort-characterization",
	)
	if err != nil {
		t.Fatalf("run OpenCode Agent Session: %v", err)
	}

	invocations := readJSONInvocations(t, harness.invocationsPath)
	sessionIndex := -1
	promptIndex := -1
	for index, invocation := range invocations {
		switch fakeACPXCommandKey(invocation) {
		case "sessions ensure":
			if sessionIndex == -1 {
				sessionIndex = index
			}
		case "prompt":
			if promptIndex == -1 {
				promptIndex = index
			}
		}
	}
	if sessionIndex == -1 || promptIndex == -1 || sessionIndex >= promptIndex {
		t.Fatalf("command sequence must contain session ensure before prompt: %#v", invocations)
	}
	for _, invocation := range invocations[sessionIndex+1 : promptIndex] {
		command := fakeACPXCommandKey(invocation)
		if command == "set effort" || command == "set reasoning_effort" {
			t.Fatalf("reasoning config set appeared between session and prompt: %#v", invocations)
		}
	}
}
