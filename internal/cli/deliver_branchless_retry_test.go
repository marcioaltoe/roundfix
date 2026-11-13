// Boundary: production delivery engine/workflow, SQLite and local Git worktrees.
// Origin is a disposable bare repository; no network or real user configuration.
package cli

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"no workspace", "worktree without branch"} {
		t.Run(name, func(t *testing.T) {
			home, repo := newItemBranchStartRepository(t)
			workflow := newItemRecoveryWorkflowForRepository(t, home, repo)
			workflow.loaded.Config.Worktree.Location = filepath.Join(filepath.Dir(home), "worktrees")
			ctx := t.Context()
			queue, err := workflow.store.CreateDeliveryQueue(ctx, repo, []string{implementTestSlug})
			if err != nil {
				t.Fatal(err)
			}
			oldBranch := deliveryBranchPrefix + implementTestSlug + "-1111111111111111"
			oldWorktree := filepath.Join(filepath.Dir(home), "unfinished-worktree")
			if name == "worktree without branch" {
				if _, _, _, err := workflow.store.RecordDeliveryQueueItemWorktree(ctx, repo, implementTestSlug, oldBranch, oldWorktree); err != nil {
					t.Fatal(err)
				}
				if err := workflow.store.SetDeliveryQueueItemWorktreeProvisioned(ctx, repo, implementTestSlug, true); err != nil {
					t.Fatal(err)
				}
			}
			item := queue.Items[0]
			item.Stage = store.DeliveryStageParked
			item.Blocker = `delivery-error: read prerequisites: refresh prerequisite default branch "main": git fetch origin main: Operation timed out`
			// Like a failed first start, park the in-memory item with no branch
			// after the store may have recorded the workspace pair.
			if err := workflow.store.UpdateDeliveryQueueItem(ctx, repo, item); err != nil {
				t.Fatal(err)
			}
			before, found, err := workflow.store.DeliveryQueue(ctx, repo)
			if err != nil || !found {
				t.Fatalf("read parked queue: found=%v err=%v", found, err)
			}
			parked := before.Items[0]
			if parked.Stage != store.DeliveryStageParked || parked.Branch != "" || parked.RetryCount != 0 || (name == "worktree without branch" && (parked.Worktree != oldWorktree || !parked.WorktreeProvisioned)) {
				t.Fatalf("parked item = %+v", parked)
			}
			engine := delivery.NewEngine(workflow.store, commandDeliveryEngineDependencies(workflow.store, workflow.loaded, io.Discard))
			result, retryErr := engine.Retry(ctx, repo, implementTestSlug)
			after, found, err := workflow.store.DeliveryQueue(ctx, repo)
			if err != nil || !found {
				t.Fatalf("read retried queue: found=%v err=%v", found, err)
			}
			stored := after.Items[0]
			if retryErr != nil {
				t.Fatalf("production retry: %v; stored stage=%s blocker=%q retry_count=%d branch=%q worktree=%q provisioned=%v", retryErr, stored.Stage, stored.Blocker, stored.RetryCount, stored.Branch, stored.Worktree, stored.WorktreeProvisioned)
			}
			if result.Stage != store.DeliveryStageQueued || result.Blocker != item.Blocker || stored.Stage != store.DeliveryStageQueued || stored.Blocker != "" || stored.RetryCount != 1 || stored.Branch != "" || stored.Worktree != "" || stored.WorktreeProvisioned {
				t.Fatalf("retry result=%+v stored=%+v, want queued, one retry and no workspace", result, stored)
			}
			branch, path, err := workflow.CreateItemBranch(ctx, repo, implementTestSlug)
			if err != nil {
				t.Fatalf("create item branch after retry: %v", err)
			}
			if !strings.HasPrefix(branch, deliveryBranchPrefix+implementTestSlug+"-") || branch == oldBranch || path == "" || path == oldWorktree {
				t.Fatalf("new branch=%q worktree=%q", branch, path)
			}
			if head, want := strings.TrimSpace(gittest.Run(t, path, "rev-parse", "HEAD")), strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "origin/main")); head != want {
				t.Fatalf("new worktree head=%q, want origin/main %q", head, want)
			}
		})
	}
}
