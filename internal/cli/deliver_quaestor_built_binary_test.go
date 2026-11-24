// Suite: measured delivery recovery through the built executable.
// Boundary: temporary HOME, repository and Run Database; no Agent or network.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func TestTheBuiltBinaryReportsTheQueueQuaestorMeasured(t *testing.T) {
	t.Parallel()
	binary := buildRoundfixBinaryForMacro(t)
	t.Run("status_of_the_measured_queue", func(t *testing.T) {
		t.Parallel()
		home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
		db, err := store.Open(t.Context(), home)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = db.Close() })
		const slug = "0008-livros-e-apuracao-do-icms"
		queue, err := db.CreateDeliveryQueue(t.Context(), repo, []string{slug})
		if err != nil {
			t.Fatal(err)
		}
		item := queue.Items[0]
		item.Stage, item.Branch, item.Worktree, item.RetryCount = store.DeliveryStageRunning, "roundfix/deliver-"+slug, "<item worktree>", 1
		if _, _, _, err := db.RecordDeliveryQueueItemWorktree(t.Context(), repo, slug, item.Branch, item.Worktree); err != nil {
			t.Fatal(err)
		}
		run, err := db.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, LocalBranch: item.Branch, SpecSlug: slug})
		if err != nil {
			t.Fatal(err)
		}
		// Quaestor measured five completed Tasks, the QA gate failed, and a Failed Run.
		for i := 1; i <= 6; i++ {
			status := "completed"
			if i == 6 {
				status = "failed"
			}
			payload := fmt.Sprintf(`{"task":"task_%02d","phase":"settled","status":%q}`, i, status)
			if _, err := db.AppendRunEvent(t.Context(), runevent.RunEvent{RunID: run.ID, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonTask, Payload: []byte(payload)}); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := db.CompleteRun(t.Context(), run.ID, store.StateFailed); err != nil {
			t.Fatal(err)
		}
		item.RunID = run.ID
		item.Stage, item.Blocker = store.DeliveryStageParked, "delivery-error: run Implement executor: roundfix implement failed with exit code 1"
		if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
			t.Fatal(err)
		}
		blocker := item.Blocker
		item.Stage, item.Blocker = store.DeliveryStageRunning, ""
		if _, _, err := db.RetryDeliveryQueueItem(t.Context(), repo, item, blocker); err != nil {
			t.Fatal(err)
		}
		// The nine prompts total 16,587,262 tokens; no cost was reported.
		for i := 0; i < 9; i++ {
			tokens := int64(1843029)
			if i == 8 {
				tokens++
			}
			if err := db.AppendTokenUsage(t.Context(), store.TokenUsageRecord{RunID: run.ID, ScopeKind: "task", ScopeID: fmt.Sprintf("task_%02d", i+1), Session: fmt.Sprintf("session-%d", i), Basis: "turn", TotalTokens: &tokens}); err != nil {
				t.Fatal(err)
			}
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
		out, diag, code := runRoundfixBinaryMacro(t, binary, repo, home, t.TempDir(), "", nil, "deliver", "status")
		want := slug + "\trunning\t-\t<item worktree>\nOwner: Delivery Queue has no recorded owner and 1 item(s) can still advance; run roundfix deliver resume\nRun: " + slug + " no Active Run\nLimits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none\nUsage: 16587262 tokens from 9 of 9 prompt(s) across 1 Run(s); cost not reported\n"
		if code != 0 || out != want || diag != "" {
			t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, out, diag, want)
		}
	})
	t.Run("retry_whose_owner_cannot_start", func(t *testing.T) {
		t.Parallel()
		home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
		setImplementFixtureAuthorizationOperations(t, repo, allDeliveryOperations...)
		db, err := store.Open(t.Context(), home)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = db.Close() })
		queue, err := db.CreateDeliveryQueue(t.Context(), repo, []string{implementTestSlug})
		if err != nil {
			t.Fatal(err)
		}
		item := queue.Items[0]
		const blocker = "delivery-error: run Implement executor: roundfix implement failed with exit code 1"
		item.Stage, item.Blocker = store.DeliveryStageParked, blocker
		if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
			t.Fatal(err)
		}
		for range 2 {
			retried := item
			retried.Stage, retried.Blocker = store.DeliveryStageQueued, ""
			if _, _, err := db.RetryDeliveryQueueItem(t.Context(), repo, retried, blocker); err != nil {
				t.Fatal(err)
			}
			if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
				t.Fatal(err)
			}
		}
		item.RetryCount = 2
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
		artifacts := filepath.Join(home, ".roundfix", "artifacts")
		if err := os.WriteFile(artifacts, []byte("not a directory"), 0600); err != nil {
			t.Fatal(err)
		}
		artifactDir, err := roundconfig.ResolveArtifactDirectory("", repo, home)
		if err != nil {
			t.Fatal(err)
		}
		out, diag, code := runRoundfixBinaryMacro(t, binary, repo, home, t.TempDir(), "", nil, "deliver", "retry", implementTestSlug)
		want := fmt.Sprintf("Preflight failed\n\nReason:\n  stat Artifact Directory %q: stat %s: not a directory\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix deliver --help' for usage.\n", artifactDir, artifactDir) + "Retry failed\n\nReason:\n  the Queue Owner did not start\n\nItem:\n  stage: parked; blocker: " + blocker + "\n\nNext:\n  roundfix deliver retry " + implementTestSlug + "\n"
		if code != exitPreflight || out != "" || diag != want {
			t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, out, diag, want)
		}
		after := openDeliveryQueueForCLI(t, home, repo)
		got := after.Items[0]
		if after.OwnerPID != 0 || got.Stage != store.DeliveryStageParked || got.Blocker != blocker || got.RetryCount != item.RetryCount {
			t.Fatalf("retry changed parked state: %+v", after)
		}
	})
}
