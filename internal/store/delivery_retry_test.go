// Suite: Delivery Queue retry persistence.
// Invariant: retry and idle-owner release are guarded, atomic Run Database transitions.
// Boundary IN: the real SQLite-backed store.
// Boundary OUT: the delivery engine and owner process lifecycle.
package store

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestRetryDeliveryQueueItemReentersAParkedItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/retry-parked-item"
	item := seedRetryDeliveryQueueItem(t, ctx, runStore, gitRoot, DeliveryStageParked, "run-unresolved")
	if err := runStore.ClaimDeliveryQueueOwner(ctx, gitRoot, 4242, "owner-identity"); err != nil {
		t.Fatalf("claim Delivery Queue owner: %v", err)
	}
	item.Stage = DeliveryStageRunning
	item.RunID = "run-retried"
	item.CandidateCommits = []string{"candidate-one", "candidate-two"}

	ownerPID, ownerIdentity, err := runStore.RetryDeliveryQueueItem(ctx, gitRoot, item, "run-unresolved")
	if err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}
	if ownerPID != 4242 || ownerIdentity != "owner-identity" {
		t.Fatalf("retry owner = pid:%d identity:%q, want pid:4242 identity:%q", ownerPID, ownerIdentity, "owner-identity")
	}
	got := readRetryDeliveryQueueItem(t, ctx, runStore, gitRoot)
	if got.Stage != DeliveryStageRunning || got.Blocker != "" || got.RunID != "run-retried" || !reflect.DeepEqual(got.CandidateCommits, item.CandidateCommits) {
		t.Fatalf("retried item = %+v, want running with retry evidence", got)
	}
}

func TestRetryDeliveryQueueItemRefusesAnItemThatIsNotParked(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/retry-non-parked-item"
	before := seedRetryDeliveryQueueItem(t, ctx, runStore, gitRoot, DeliveryStageReviewing, "")
	requested := before
	requested.Stage = DeliveryStageRunning
	requested.RunID = "replacement-run"

	_, _, err := runStore.RetryDeliveryQueueItem(ctx, gitRoot, requested, "run-unresolved")
	if err == nil || !strings.Contains(err.Error(), `stage "reviewing"`) || !strings.Contains(err.Error(), `blocker ""`) {
		t.Fatalf("retry non-parked item error = %v, want stored stage and blocker", err)
	}
	assertRetryDeliveryQueueItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestRetryDeliveryQueueItemRefusesAChangedBlocker(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/retry-changed-blocker"
	before := seedRetryDeliveryQueueItem(t, ctx, runStore, gitRoot, DeliveryStageParked, "review-stale")
	requested := before
	requested.Stage = DeliveryStageReviewing

	_, _, err := runStore.RetryDeliveryQueueItem(ctx, gitRoot, requested, "run-unresolved")
	if err == nil || !strings.Contains(err.Error(), `stage "parked"`) || !strings.Contains(err.Error(), `blocker "review-stale"`) {
		t.Fatalf("retry changed blocker error = %v, want stored stage and blocker", err)
	}
	assertRetryDeliveryQueueItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestRetryDeliveryQueueItemRefusesATerminalStage(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/retry-terminal-stage"
	before := seedRetryDeliveryQueueItem(t, ctx, runStore, gitRoot, DeliveryStageParked, "gate-failed")
	requested := before
	requested.Stage = DeliveryStageMerged

	_, _, err := runStore.RetryDeliveryQueueItem(ctx, gitRoot, requested, "gate-failed")
	if err == nil || !strings.Contains(err.Error(), `target stage "merged"`) {
		t.Fatalf("retry terminal stage error = %v, want target-stage refusal", err)
	}
	assertRetryDeliveryQueueItemUnchanged(t, ctx, runStore, gitRoot, before)
}

func TestReleaseIdleDeliveryQueueOwnerReleasesAnIdleQueue(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/release-idle-delivery-owner"
	seedRetryDeliveryQueueItem(t, ctx, runStore, gitRoot, DeliveryStageParked, "run-unresolved")
	if err := runStore.ClaimDeliveryQueueOwner(ctx, gitRoot, 4242, "owner-identity"); err != nil {
		t.Fatalf("claim Delivery Queue owner: %v", err)
	}

	released, err := runStore.ReleaseIdleDeliveryQueueOwner(ctx, gitRoot, 4242, "owner-identity")
	if err != nil || !released {
		t.Fatalf("release idle Delivery Queue owner: released=%v err=%v", released, err)
	}
	queue, found, err := runStore.DeliveryQueue(ctx, gitRoot)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue: found=%v err=%v", found, err)
	}
	if queue.OwnerPID != 0 || queue.OwnerIdentity != "" {
		t.Fatalf("released owner = pid:%d identity:%q", queue.OwnerPID, queue.OwnerIdentity)
	}
}

func TestReleaseIdleDeliveryQueueOwnerKeepsAQueueWithAnAdvanceableItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/keep-active-delivery-owner"
	seedRetryDeliveryQueueItem(t, ctx, runStore, gitRoot, DeliveryStageRunning, "")
	if err := runStore.ClaimDeliveryQueueOwner(ctx, gitRoot, 4242, "owner-identity"); err != nil {
		t.Fatalf("claim Delivery Queue owner: %v", err)
	}

	released, err := runStore.ReleaseIdleDeliveryQueueOwner(ctx, gitRoot, 4242, "owner-identity")
	if err != nil || released {
		t.Fatalf("release active Delivery Queue owner: released=%v err=%v", released, err)
	}
	queue, _, _ := runStore.DeliveryQueue(ctx, gitRoot)
	if queue.OwnerPID != 4242 || queue.OwnerIdentity != "owner-identity" {
		t.Fatalf("retained owner = pid:%d identity:%q", queue.OwnerPID, queue.OwnerIdentity)
	}
}

func TestReleaseIdleDeliveryQueueOwnerRefusesAnotherOwner(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/refuse-other-delivery-owner"
	seedRetryDeliveryQueueItem(t, ctx, runStore, gitRoot, DeliveryStageParked, "gate-failed")
	if err := runStore.ClaimDeliveryQueueOwner(ctx, gitRoot, 4242, "owner-identity"); err != nil {
		t.Fatalf("claim Delivery Queue owner: %v", err)
	}

	released, err := runStore.ReleaseIdleDeliveryQueueOwner(ctx, gitRoot, 4343, "other-identity")
	if err == nil || released || !strings.Contains(err.Error(), "PID 4242") || !strings.Contains(err.Error(), `"owner-identity"`) {
		t.Fatalf("release another owner: released=%v err=%v, want recorded owner", released, err)
	}
	queue, _, _ := runStore.DeliveryQueue(ctx, gitRoot)
	if queue.OwnerPID != 4242 || queue.OwnerIdentity != "owner-identity" {
		t.Fatalf("owner after refusal = pid:%d identity:%q", queue.OwnerPID, queue.OwnerIdentity)
	}
}

func seedRetryDeliveryQueueItem(
	t *testing.T,
	ctx context.Context,
	runStore *Store,
	gitRoot string,
	stage DeliveryStage,
	blocker string,
) DeliveryQueueItem {
	t.Helper()
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"0173-retry"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = stage
	item.Blocker = blocker
	item.Branch = "roundfix/deliver-0173-retry"
	item.RunID = "run-original"
	item.CandidateCommits = []string{"candidate-original"}
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed Delivery Queue item: %v", err)
	}
	return item
}

func readRetryDeliveryQueueItem(t *testing.T, ctx context.Context, runStore *Store, gitRoot string) DeliveryQueueItem {
	t.Helper()
	queue, found, err := runStore.DeliveryQueue(ctx, gitRoot)
	if err != nil || !found || len(queue.Items) != 1 {
		t.Fatalf("read Delivery Queue item: found=%v items=%d err=%v", found, len(queue.Items), err)
	}
	return queue.Items[0]
}

func assertRetryDeliveryQueueItemUnchanged(
	t *testing.T,
	ctx context.Context,
	runStore *Store,
	gitRoot string,
	want DeliveryQueueItem,
) {
	t.Helper()
	if got := readRetryDeliveryQueueItem(t, ctx, runStore, gitRoot); !reflect.DeepEqual(got, want) {
		t.Fatalf("Delivery Queue item changed after refusal:\n got: %#v\nwant: %#v", got, want)
	}
}
