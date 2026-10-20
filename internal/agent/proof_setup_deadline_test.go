// Suite: Exact selection proof setup deadline
// Invariant: only the setup context uses the requested deadline, with the built-in fallback.
// Boundary IN: ProveExactSelection and Probe with a compiled fake acpx.
// Boundary OUT: codex hygiene inspection records the setup context's remaining budget.
package agent

import (
	"context"
	"testing"
	"time"
)

func TestExactSelectionProofSetupUsesTheRequestedDeadline(t *testing.T) {
	t.Parallel()
	for _, entry := range []string{"ProveExactSelection", "Probe"} {
		for _, timeout := range []time.Duration{2 * time.Hour, 0, -time.Second} {
			t.Run(entry+"/"+timeout.String(), func(t *testing.T) {
				t.Parallel()
				harness := newFakeACPXHarness(t)
				harness.setEnv(fakeACPXStdout, MinimumACPXVersion+"\n")
				harness.setEnv(fakeACPXStdoutCall, mustJSONForTest(t, map[string]string{
					"sessions show": sessionCapabilitySnapshotFixture(t, "gpt-5.5", []string{"gpt-5.5"}, "reasoning_effort", "medium", []string{"medium", "xhigh"}),
				}))
				harness.setEnv(codexPathEnv, "/configured/clean/codex")
				probe := &remainingDeadlineCodexProbe{}
				harness.runner.codexSpawn = codexSpawnDependencies{
					goos:       "darwin",
					lookPath:   fakeCodexLookPath("/path/clean/codex", nil),
					quarantine: probe,
					acceptance: probe,
				}
				request := ProbeRequest{
					Runtime:      RuntimeSpec{ID: "codex", Protocol: ProtocolACP, Model: "gpt-5.5", ReasoningEffort: "xhigh"},
					WorkDir:      harness.gitRoot,
					SetupTimeout: timeout,
				}
				var err error
				if entry == "Probe" {
					err = harness.runner.Probe(context.Background(), request)
				} else {
					_, err = harness.runner.ProveExactSelection(context.Background(), request)
				}
				if err != nil {
					t.Fatal(err)
				}
				if len(probe.remaining) == 0 {
					t.Fatal("no setup contexts recorded")
				}
				for index, remaining := range probe.remaining {
					if timeout > 0 {
						if remaining <= time.Hour || remaining > timeout {
							t.Fatalf("context %d remaining = %s, want more than 1h and at most %s", index, remaining, timeout)
						}
					} else if remaining <= 0 || remaining > acpxPreflightSetupTimeout {
						t.Fatalf("context %d remaining = %s, want positive and at most %s", index, remaining, acpxPreflightSetupTimeout)
					}
				}
			})
		}
	}
}

type remainingDeadlineCodexProbe struct{ remaining []time.Duration }

func (probe *remainingDeadlineCodexProbe) record(ctx context.Context) {
	deadline, ok := ctx.Deadline()
	if !ok {
		probe.remaining = append(probe.remaining, 0)
		return
	}
	probe.remaining = append(probe.remaining, time.Until(deadline))
}

func (probe *remainingDeadlineCodexProbe) Quarantined(ctx context.Context, _ string) (bool, error) {
	probe.record(ctx)
	return false, nil
}

func (probe *remainingDeadlineCodexProbe) Accepted(ctx context.Context, _ string) (bool, error) {
	probe.record(ctx)
	return true, nil
}
