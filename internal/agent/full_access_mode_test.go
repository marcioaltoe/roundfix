package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestChooseFullAccessModeTakesTheFirstAdvertisedCandidate(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name                   string
		candidates, advertised []string
		want                   string
		ok                     bool
	}{
		{"candidate order", []string{"agent-full-access", "full-access"}, []string{"full-access", "agent-full-access"}, "agent-full-access", true},
		{"legacy candidate", []string{"agent-full-access", "full-access"}, []string{"full-access"}, "full-access", true},
		{"trim spaces", []string{" agent-full-access "}, []string{" agent-full-access "}, "agent-full-access", true},
		{"exact comparison", []string{"full-access"}, []string{"FULL-ACCESS", "agent-full-access"}, "", false},
		{"no values", []string{"full-access"}, nil, "", false},
		{"no candidates", nil, []string{"full-access"}, "", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := chooseFullAccessMode(tt.candidates, tt.advertised)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("choice = %q, %t; want %q, %t", got, ok, tt.want, tt.ok)
			}
		})
	}
	harness := newFakeACPXHarness(t)
	runtime := fullAccessRuntimeForTest(t, "codex")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": fullAccessSnapshotForTest(t, runtime, []string{"full-access", "agent-full-access"}, false),
	}))
	if _, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot}); err != nil {
		t.Fatal(err)
	}
	assertFullAccessCalls(t, harness, []string{"set-mode agent-full-access"})

}

func TestProofSelectsTheAdvertisedCodexFullAccessMode(t *testing.T) {
	t.Parallel()
	harness := newFakeACPXHarness(t)
	runtime := fullAccessRuntimeForTest(t, "codex")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": fullAccessSnapshotForTest(t, runtime, []string{"read-only", "workspace-write", "agent", "agent-full-access"}, false),
	}))
	proof, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})
	if err != nil {
		t.Fatal(err)
	}
	if proof.EffectiveAccessPolicy != AccessPolicyFullAccess {
		t.Fatalf("effective policy = %q", proof.EffectiveAccessPolicy)
	}
	assertFullAccessCalls(t, harness, []string{"set-mode agent-full-access"})
	assertLastInvocationClosesDisposable(t, harness)
}

func TestRunAppliesTheSandboxPresetWhenTheSessionAdvertisesIt(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"full-access", "agent-full-access"} {
		t.Run(mode, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			runtime := fullAccessRuntimeForTest(t, "codex")
			harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
				"sessions show": fullAccessSnapshotForTest(t, runtime, []string{"read-only", mode}, true),
			}))
			if _, err := harness.run(context.Background(), runtime, "full-access-run"); err != nil {
				t.Fatal(err)
			}
			assertFullAccessCalls(t, harness, []string{"set-mode " + mode, "set sandbox_mode value=danger-full-access"})
		})
	}
}

func TestRunSelectsBypassPermissionsForClaude(t *testing.T) {
	t.Parallel()
	harness := newFakeACPXHarness(t)
	runtime := fullAccessRuntimeForTest(t, "claude")
	harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
		"sessions show": fullAccessSnapshotForTest(t, runtime, []string{"default", "bypassPermissions"}, false),
	}))
	if _, err := harness.run(context.Background(), runtime, "claude-full-access-run"); err != nil {
		t.Fatal(err)
	}
	assertFullAccessCalls(t, harness, []string{"set-mode bypassPermissions"})
}

func TestFullAccessWithoutAnAdvertisedModeIsDegraded(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name  string
		modes []string
		list  string
	}{
		{"unmatched", []string{"read-only", "workspace-write", "agent"}, "read-only, workspace-write, agent"},
		{"no mode option", nil, "none"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			runtime := fullAccessRuntimeForTest(t, "codex")
			harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
				"sessions show": fullAccessSnapshotForTest(t, runtime, tt.modes, true),
			}))
			proof, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})
			if err != nil {
				t.Fatal(err)
			}
			want := "full-access (degraded: no advertised mode grants full access; advertised modes: " + tt.list + ")"
			if string(proof.EffectiveAccessPolicy) != want {
				t.Fatalf("policy = %q, want %q", proof.EffectiveAccessPolicy, want)
			}
			assertFullAccessCalls(t, harness, nil)
			assertLastInvocationClosesDisposable(t, harness)
			var warnings []string
			harness.runner.warnf = func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }
			sink := newCaptureSink("")
			if _, err := harness.runWithSink(context.Background(), runtime, "degraded-run", sink); err != nil {
				t.Fatal(err)
			}
			if !sink.HasStatus(AccessModeNotAdvertisedStatus) {
				t.Fatalf("missing degraded status: %+v", sink.Events())
			}
			warning := "full-access mode not advertised for Agent Session degraded-run; advertised modes: " + tt.list
			if !containsWarning(warnings, warning) {
				t.Fatalf("warnings = %#v, want %q", warnings, warning)
			}
			assertFullAccessCalls(t, harness, nil)
			if !containsCommandKey(readJSONInvocations(t, harness.invocationsPath), "prompt") {
				t.Fatal("Run did not reach prompt")
			}
		})
	}
}

func fullAccessRuntimeForTest(t *testing.T, id string) RuntimeSpec {
	t.Helper()
	runtime, err := RuntimeFor(RuntimeOptions{Agent: id, EnableFullAccess: true})
	if err != nil {
		t.Fatal(err)
	}
	return selectedRuntime(runtime)
}

func fullAccessSnapshotForTest(t *testing.T, runtime RuntimeSpec, modes []string, sandbox bool) string {
	t.Helper()
	reasoningID := "reasoning_effort"
	if runtime.ID == "claude" {
		reasoningID = "effort"
	}
	var snapshot struct {
		Schema string `json:"schema"`
		ACPX   struct {
			Model   string           `json:"current_model_id"`
			Options []map[string]any `json:"config_options"`
		} `json:"acpx"`
	}
	payload := sessionCapabilitySnapshotFixture(t, runtime.Model, []string{runtime.Model}, reasoningID, runtime.ReasoningEffort, []string{runtime.ReasoningEffort})
	if err := json.Unmarshal([]byte(payload), &snapshot); err != nil {
		t.Fatal(err)
	}
	options := snapshot.ACPX.Options[:0]
	for _, option := range snapshot.ACPX.Options {
		if option["id"] != "mode" && option["id"] != "sandbox_mode" {
			options = append(options, option)
		}
	}
	if len(modes) > 0 {
		values := make([]map[string]string, 0, len(modes))
		for _, mode := range modes {
			values = append(values, map[string]string{"value": mode})
		}
		options = append(options, map[string]any{"id": "mode", "category": "mode", "type": "select", "currentValue": modes[0], "options": values})
	}
	if sandbox {
		options = append(options, map[string]any{"id": "sandbox_mode", "type": "select", "currentValue": "workspace-write", "options": []map[string]string{{"value": "workspace-write"}, {"value": "danger-full-access"}}})
	}
	snapshot.ACPX.Options = options
	return mustJSONForTest(t, snapshot)
}

func assertFullAccessCalls(t *testing.T, harness *fakeACPXHarness, want []string) {
	t.Helper()
	var got []string
	shows := make(map[string]int)
	for _, call := range readJSONInvocations(t, harness.invocationsPath) {
		switch fakeACPXCommandKey(call) {
		case "set-mode":
			for i, arg := range call {
				if arg == "set-mode" {
					got = append(got, "set-mode "+call[i+1])
					break
				}
			}
		case "set sandbox_mode":
			for i, arg := range call {
				if arg == "set" {
					got = append(got, "set sandbox_mode value="+call[i+2])
					break
				}
			}
		case "sessions show":
			shows[call[len(call)-1]]++
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("access calls = %#v, want %#v", got, want)
	}
	for session, count := range shows {
		if count > 2 {
			t.Fatalf("unexpected extra sessions show calls for %s: %d", session, count)
		}
	}
}
