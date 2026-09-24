// Suite: Delivery Queue persistence.
// Invariant: reopening the real Run Database preserves ordered items and exposes every intent without a receipt.
// Boundary IN: the SQLite-backed store API and schema migration from the preceding version.
// Boundary OUT: delivery stage transitions and external actions, owned by internal/delivery.
package store

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func TestDeliveryQueueRoundTripsItemsAndReceipts(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()

	legacy := openTestStore(t, ctx, homeDir)
	legacyRun, err := legacy.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create existing Run: %v", err)
	}
	closeStore(t, legacy)
	downgradeDeliverySchemaFixture(t, ctx, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	queue, err := runStore.CreateDeliveryQueue(ctx, "/tmp/roundfix", []string{
		"0156-a-delivery-loop-that-outlives-the-session",
		"0157-next-spec",
	})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}

	first := queue.Items[0]
	first.Stage = DeliveryStagePublishing
	first.RunID = legacyRun.ID
	first.CandidateCommits = []string{"candidate-one", "candidate-two"}
	first.PullRequestNumber = "42"
	if err := runStore.UpdateDeliveryQueueItem(ctx, queue.GitRoot, first); err != nil {
		t.Fatalf("update first Delivery Queue item: %v", err)
	}

	second := queue.Items[1]
	second.Stage = DeliveryStageParked
	second.Blocker = "review-stale"
	second.RunID = "run_second"
	second.CandidateCommits = []string{"candidate-three"}
	second.PullRequestNumber = "43"
	second.MergeCommit = "merge-three"
	if err := runStore.UpdateDeliveryQueueItem(ctx, queue.GitRoot, second); err != nil {
		t.Fatalf("update second Delivery Queue item: %v", err)
	}

	pushIntent, err := runStore.RecordDeliveryActionIntent(ctx, queue.GitRoot, first.SpecSlug, DeliveryActionPush)
	if err != nil {
		t.Fatalf("record push intent: %v", err)
	}
	if _, err := runStore.RecordDeliveryActionReceipt(ctx, pushIntent.ID, "candidate-two"); err != nil {
		t.Fatalf("record push receipt: %v", err)
	}
	mergeIntent, err := runStore.RecordDeliveryActionIntent(ctx, queue.GitRoot, second.SpecSlug, DeliveryActionMerge)
	if err != nil {
		t.Fatalf("record merge intent: %v", err)
	}
	closeStore(t, runStore)

	reopened := openTestStore(t, ctx, homeDir)
	defer closeStore(t, reopened)
	persisted, found, err := reopened.DeliveryQueue(ctx, queue.GitRoot)
	if err != nil || !found {
		t.Fatalf("read persisted Delivery Queue: found=%v err=%v", found, err)
	}
	wantItems := []DeliveryQueueItem{first, second}
	if !reflect.DeepEqual(persisted.Items, wantItems) {
		t.Fatalf("persisted Delivery Queue items = %#v, want %#v", persisted.Items, wantItems)
	}
	if _, found, err := reopened.Run(ctx, legacyRun.ID); err != nil || !found {
		t.Fatalf("read existing Run after Delivery Queue migration: found=%v err=%v", found, err)
	}

	receipt, found, err := reopened.DeliveryActionReceipt(ctx, pushIntent.ID)
	if err != nil || !found {
		t.Fatalf("read persisted Delivery Action receipt: found=%v err=%v", found, err)
	}
	if receipt.IntentID != pushIntent.ID || receipt.Result != "candidate-two" {
		t.Fatalf("persisted Delivery Action receipt = %#v", receipt)
	}
	unmatched, err := reopened.UnmatchedDeliveryActionIntents(ctx, queue.GitRoot)
	if err != nil {
		t.Fatalf("list unmatched Delivery Action intents: %v", err)
	}
	if len(unmatched) != 1 || unmatched[0].ID != mergeIntent.ID || unmatched[0].Action != DeliveryActionMerge {
		t.Fatalf("unmatched Delivery Action intents = %#v, want merge intent %#v", unmatched, mergeIntent)
	}
}

func downgradeDeliverySchemaFixture(t *testing.T, ctx context.Context, homeDir string) {
	t.Helper()
	db, err := sql.Open("sqlite", writerDSN(DatabasePath(homeDir)))
	if err != nil {
		t.Fatalf("open Delivery Queue migration fixture: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close Delivery Queue migration fixture: %v", err)
		}
	}()

	statements := []string{
		`DROP TABLE delivery_action_receipts`,
		`DROP TABLE delivery_action_intents`,
		`DROP TABLE delivery_queue_items`,
		`DROP TABLE delivery_queues`,
		`PRAGMA user_version = 13`,
	}
	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("build schema 13 Delivery Queue fixture: %v", err)
		}
	}
}
