// Suite: Delivery Queue merge persistence.
// Invariant: recording a merge requires the observed park and preserves retry and workspace metadata.
// Boundary IN: the real SQLite store in a temporary home.
// Boundary OUT: the merge observer and owner lifecycle.
package store

import (
	"reflect"
	"strings"
	"testing"
)

func TestRecordDeliveryQueueItemMergedRequiresTheObservedPark(t *testing.T) {
	t.Parallel()
	for _, test := range []string{"changed blocker", "changed stage", "changed slug", "missing position"} {
		t.Run(test, func(t *testing.T) {
			ctx := t.Context()
			writer := openTestStore(t, ctx, t.TempDir())
			defer closeStore(t, writer)
			root := "/repo-record-merge-refusal"
			before := seedRetryDeliveryQueueItem(t, ctx, writer, root, DeliveryStageParked, "checks-failed")
			requested := before
			requested.Stage, requested.MergeCommit = DeliveryStageMerged, "merge-commit"
			requested.CandidateCommits = []string{"candidate-original", "merged-head"}
			blocker := before.Blocker
			switch test {
			case "changed blocker":
				blocker = "review-stale"
			case "changed stage":
				before.Stage = DeliveryStageReviewing
				if err := writer.UpdateDeliveryQueueItem(ctx, root, before); err != nil {
					t.Fatal(err)
				}
			case "changed slug":
				requested.SpecSlug = "another-spec"
			case "missing position":
				requested.Position++
			}
			if _, _, err := writer.RecordDeliveryQueueItemMerged(ctx, root, requested, blocker); err == nil {
				t.Fatal("accepted a changed park")
			}
			assertRetryDeliveryQueueItemUnchanged(t, ctx, writer, root, before)
		})
	}
}

func TestRecordDeliveryQueueItemMergedPreservesStoredMetadata(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	writer := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, writer)
	root := "/repo-record-merge"
	before := seedRetryDeliveryQueueItem(t, ctx, writer, root, DeliveryStageParked, "checks-failed")
	before.Warning, before.PullRequestNumber = "preserved warning", "404"
	if err := writer.UpdateDeliveryQueueItem(ctx, root, before); err != nil {
		t.Fatal(err)
	}
	if err := writer.ClaimDeliveryQueueOwner(ctx, root, 4242, "owner-identity"); err != nil {
		t.Fatal(err)
	}
	requested := before
	requested.Stage, requested.MergeCommit = DeliveryStageMerged, "merge-commit"
	requested.CandidateCommits = []string{"candidate-original", "merged-head"}
	// Stale caller metadata must never overwrite persisted metadata.
	requested.RunID, requested.Branch, requested.Worktree, requested.Warning, requested.PullRequestNumber = "stale-run", "stale-branch", "stale-worktree", "stale-warning", "999"
	pid, identity, err := writer.RecordDeliveryQueueItemMerged(ctx, root, requested, before.Blocker)
	if err != nil {
		t.Fatal(err)
	}
	want := before
	want.Stage, want.Blocker, want.MergeCommit, want.CandidateCommits = DeliveryStageMerged, "", requested.MergeCommit, requested.CandidateCommits
	if got := readRetryDeliveryQueueItem(t, ctx, writer, root); !reflect.DeepEqual(got, want) {
		t.Fatalf("merged item = %#v, want %#v", got, want)
	}
	if pid != 4242 || identity != "owner-identity" {
		t.Fatalf("owner = %d %q", pid, identity)
	}
	if _, _, err := writer.RecordDeliveryQueueItemMerged(ctx, root, requested, before.Blocker); err == nil || !strings.Contains(err.Error(), `stage "merged"`) {
		t.Fatalf("second record = %v", err)
	}
	assertRetryDeliveryQueueItemUnchanged(t, ctx, writer, root, want)
}
