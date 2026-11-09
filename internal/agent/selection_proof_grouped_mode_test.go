package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestProofOverTheRecordedCodexACPPayloads(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		version string
		policy  AccessPolicy
	}{
		{"2.6.0/full-access", "2.6.0", AccessPolicyFullAccess},
		{"2.6.0/runtime-default", "2.6.0", AccessPolicyRuntimeDefault},
		{"2.5.0/full-access", "2.5.0", AccessPolicyFullAccess},
		{"2.5.0/runtime-default", "2.5.0", AccessPolicyRuntimeDefault},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fixture, err := os.ReadFile(filepath.Join("testdata", "codex-acp-"+tt.version+"-config-options.json"))
			if err != nil {
				t.Fatalf("read recorded config options: %v", err)
			}
			snapshot := `{"schema":"acpx.session.v1","acpx":{"current_model_id":"gpt-6.1-sol","config_options":` + string(fixture) + `}}`
			harness := newFakeACPXHarness(t)
			harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
				"sessions show": snapshot,
			}))

			var runtime RuntimeSpec
			if tt.policy == AccessPolicyFullAccess {
				runtime = fullAccessRuntimeForTest(t, "codex")
			} else {
				runtime = selectedRuntime(RuntimeSpec{ID: "codex", Protocol: ProtocolACP})
			}
			runtime.Model = "gpt-6.1-sol"
			runtime.ReasoningEffort = "high"

			proof, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{
				Runtime: runtime,
				WorkDir: harness.gitRoot,
			})
			if err != nil {
				t.Fatalf("prove exact selection: %v", err)
			}
			if proof.Status != SelectionProofStatusProven || proof.Model != runtime.Model || proof.ReasoningEffort != runtime.ReasoningEffort {
				t.Errorf("unexpected proof: %#v", proof)
			}
			if proof.Assignment.Model != runtime.Model || proof.Assignment.ReasoningEffort != runtime.ReasoningEffort || proof.Assignment.Encoding != SelectionEncodingIndependent {
				t.Errorf("unexpected assignment: %#v", proof.Assignment)
			}
			if proof.EffectiveAccessPolicy != tt.policy {
				t.Errorf("effective policy = %q, want %q", proof.EffectiveAccessPolicy, tt.policy)
			}
			assertLastInvocationClosesDisposable(t, harness)
			if tt.policy == AccessPolicyFullAccess {
				assertFullAccessCalls(t, harness, []string{"set-mode agent-full-access"})
			} else {
				assertFullAccessCalls(t, harness, nil)
			}
		})
	}
}
