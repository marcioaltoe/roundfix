// Suite: guarded Delivery Retry rollback.
// Invariant: a failed owner start spends no retry and preserves Run provenance.
// Boundary IN: temporary SQLite Run Database.
// Boundary OUT: owner startup and process proof.
package store

import (
	"reflect"
	"testing"
)

func TestRevertDeliveryQueueItemRetryRestoresTheParkedItem(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"reverts_a_retried_item", "refuses_while_an_owner_is_recorded", "refuses_when_the_item_moved", "refuses_when_the_blocker_changed"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := openTestStore(t, ctx, t.TempDir())
			defer closeStore(t, db)
			const root = "/tmp/revert-retry"
			queue, err := db.CreateDeliveryQueue(ctx, root, []string{"0173-retry"})
			if err != nil {
				t.Fatal(err)
			}
			if _, _, _, err := db.RecordDeliveryQueueItemWorktree(ctx, root, queue.Items[0].SpecSlug, "roundfix/deliver-0173-retry", "/tmp/item-worktree"); err != nil {
				t.Fatal(err)
			}
			before := readRetryDeliveryQueueItem(t, ctx, db, root)
			before.Stage, before.Blocker, before.RunID = DeliveryStageParked, "run-unresolved", "run-original"
			before.CandidateCommits = []string{"candidate-original"}
			if err := db.UpdateDeliveryQueueItem(ctx, root, before); err != nil {
				t.Fatal(err)
			}
			retry := before
			retry.Stage = DeliveryStageRunning
			retry.RunID = "new-run"
			retry.CandidateCommits = []string{"new-candidate"}
			if _, _, err := db.RetryDeliveryQueueItem(ctx, root, retry, before.Blocker); err != nil {
				t.Fatal(err)
			}
			if name == "refuses_while_an_owner_is_recorded" {
				if err := db.ClaimDeliveryQueueOwner(ctx, root, 42, "recorded"); err != nil {
					t.Fatal(err)
				}
			}
			if name == "refuses_when_the_item_moved" {
				retry.Stage = DeliveryStageReviewing
				if err := db.UpdateDeliveryQueueItem(ctx, root, retry); err != nil {
					t.Fatal(err)
				}
			}
			if name == "refuses_when_the_blocker_changed" {
				retry.Blocker = "new-blocker"
				if err := db.UpdateDeliveryQueueItem(ctx, root, retry); err != nil {
					t.Fatal(err)
				}
			}
			unchanged, _, err := db.DeliveryQueue(ctx, root)
			if err != nil {
				t.Fatal(err)
			}
			reverted, err := db.RevertDeliveryQueueItemRetry(ctx, root, before, DeliveryStageRunning)
			if err != nil {
				t.Fatal(err)
			}
			after, _, err := db.DeliveryQueue(ctx, root)
			if err != nil {
				t.Fatal(err)
			}
			var links int
			if err := db.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM delivery_queue_runs WHERE git_root = ? AND run_id IN ('run-original', 'new-run')", root).Scan(&links); err != nil {
				t.Fatal(err)
			}
			if links != 2 {
				t.Fatalf("Run links after revert=%d, want both original and retry provenance", links)
			}
			if name != "reverts_a_retried_item" {
				if reverted || !reflect.DeepEqual(after, unchanged) {
					t.Fatalf("guard changed queue: %+v", after)
				}
				return
			}
			want := unchanged.Items[0]
			want.Stage, want.Blocker, want.RetryCount, want.CandidateCommits = before.Stage, before.Blocker, before.RetryCount, before.CandidateCommits
			if !reverted || !reflect.DeepEqual(after.Items[0], want) {
				t.Fatalf("revert=%v queue=%+v want item=%+v", reverted, after, want)
			}
		})
	}
}
