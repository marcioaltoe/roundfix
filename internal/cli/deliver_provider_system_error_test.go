// Suite: Delivery Run results classify exhausted Provider System Errors.
// Boundary IN: the Run Event journal and delivery workflow result reader.
package cli

import (
	"encoding/json"
	"testing"
	"time"

	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func TestDeliveryRunResultReadsAnExhaustedProviderSystemError(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		phase    string
		recovery string
		want     string
	}{
		{"provider exhausted", "provider_system_error", "exhausted", "qa hit a provider system error: model at capacity"},
		{"provider fallback", "provider_system_error", "fallback", ""},
		{"lost rollout exhausted", "rollout_lost", "exhausted", "qa lost its rollout at session/resume"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			s := openDeliverReleaseStore(t, ctx, t.TempDir())
			root := t.TempDir()
			run, err := s.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: root, LocalBranch: "feat/runtime", HeadSHA: "fixture-head", SpecSlug: "runtime-spec", Agent: "codex"})
			if err != nil {
				t.Fatal(err)
			}
			payload, err := json.Marshal(map[string]any{"phase": tc.phase, "scope_id": "qa", "step": "session/resume", "detail": "model at capacity", "recovery": tc.recovery, "retry_spent": false})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := s.AppendRunEvent(ctx, runevent.RunEvent{RunID: run.ID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, Time: time.Now(), Payload: payload}); err != nil {
				t.Fatal(err)
			}
			workflow := &commandDeliveryWorkflow{store: s, git: preflight.ExecGitRunner{}}
			run.State = store.StateUnresolved
			got, err := workflow.runResult(ctx, root, "runtime-spec", roundfixCommandResult{exitCode: exitRunFailed}, "", nil, &run)
			if err != nil {
				t.Fatal(err)
			}
			if got.RuntimeInfrastructure != tc.want {
				t.Fatalf("result=%+v want detail=%q", got, tc.want)
			}
		})
	}
}
