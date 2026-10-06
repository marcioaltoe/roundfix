// Suite: Delivery Run results read exhausted Lost Rollouts from the journal.
package cli

import (
	"encoding/json"
	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
	"testing"
	"time"
)

func TestDeliveryRunResultReadsAnExhaustedLostRollout(t *testing.T) {
	for _, recovery := range []string{"fallback", "new_session", "exhausted"} {
		t.Run(recovery, func(t *testing.T) {
			ctx := t.Context()
			s := openDeliverReleaseStore(t, ctx, t.TempDir())
			root := t.TempDir()
			run, err := s.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: root, LocalBranch: "feat/runtime", HeadSHA: "fixture-head", SpecSlug: "runtime-spec", Agent: "codex"})
			if err != nil {
				t.Fatal(err)
			}
			payload, err := json.Marshal(map[string]any{"phase": "rollout_lost", "scope_id": "qa", "step": "session/resume", "recovery": recovery, "retry_spent": false})
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
			want := ""
			if recovery == "exhausted" {
				want = "qa lost its rollout at session/resume"
			}
			if got.RuntimeInfrastructure != want {
				t.Fatalf("result=%+v want detail=%q", got, want)
			}
		})
	}
}
