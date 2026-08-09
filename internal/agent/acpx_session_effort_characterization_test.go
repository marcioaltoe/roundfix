package agent

import (
	"context"
	"errors"
	"testing"
)

// Suite: OpenCode Agent Session effort characterization
// Invariant: today's OpenCode Run reaches its work prompt without setting a reasoning option.
// Boundary IN: ACPXRunner Run lifecycle and the acpx command sequence it executes.
// Boundary OUT: profile configuration parsing and future warm-up behavior.

func TestACPXSessionEffortCharacterizationTodayReasoningKeyRefusesOpenCode(t *testing.T) {
	t.Parallel()

	_, err := acpxReasoningEffortConfigKey(RuntimeSpec{ID: "opencode"})
	var managed *ModelManagedReasoningError
	if !errors.As(err, &managed) {
		t.Fatalf("error = %T %v, want *ModelManagedReasoningError", err, err)
	}
	if managed.Runtime != "opencode" {
		t.Fatalf("managed runtime = %q, want opencode", managed.Runtime)
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
