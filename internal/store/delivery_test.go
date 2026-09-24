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
	if err := runStore.ClaimDeliveryQueueOwner(ctx, queue.GitRoot, 4242, "owner-identity"); err != nil {
		t.Fatalf("claim Delivery Queue owner: %v", err)
	}
	if err := runStore.ClaimDeliveryQueueOwner(ctx, queue.GitRoot, 4343, "different-owner"); err == nil {
		t.Fatal("different Delivery Queue owner claim succeeded")
	}

	first := queue.Items[0]
	first.Stage = DeliveryStagePublishing
	first.Branch = "roundfix/deliver-0156"
	first.RunID = legacyRun.ID
	first.CandidateCommits = []string{"candidate-one", "candidate-two"}
	first.PullRequestNumber = "42"
	if err := runStore.UpdateDeliveryQueueItem(ctx, queue.GitRoot, first); err != nil {
		t.Fatalf("update first Delivery Queue item: %v", err)
	}

	second := queue.Items[1]
	second.Stage = DeliveryStageParked
	second.Branch = "roundfix/deliver-0157"
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
	if persisted.OwnerPID != 4242 || persisted.OwnerIdentity != "owner-identity" {
		t.Fatalf("persisted Delivery Queue owner = pid:%d identity:%q", persisted.OwnerPID, persisted.OwnerIdentity)
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
	released, err := reopened.ReleaseDeliveryQueueOwner(ctx, queue.GitRoot, 4242, "owner-identity")
	if err != nil || !released {
		t.Fatalf("release Delivery Queue owner: released=%v err=%v", released, err)
	}
	withoutOwner, found, err := reopened.DeliveryQueue(ctx, queue.GitRoot)
	if err != nil || !found {
		t.Fatalf("read released Delivery Queue: found=%v err=%v", found, err)
	}
	if withoutOwner.OwnerPID != 0 || withoutOwner.OwnerIdentity != "" {
		t.Fatalf("released Delivery Queue owner = pid:%d identity:%q", withoutOwner.OwnerPID, withoutOwner.OwnerIdentity)
	}
}

func TestCreateDeliveryQueueReplacesOnlyATerminalUnownedQueue(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/replace-delivery"
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"unfinished"})
	if err != nil {
		t.Fatalf("create unfinished Delivery Queue: %v", err)
	}
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"replacement"}); err == nil {
		t.Fatal("unfinished Delivery Queue was replaced")
	}
	item := queue.Items[0]
	item.Stage = DeliveryStageMerged
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("finish Delivery Queue: %v", err)
	}
	if err := runStore.ClaimDeliveryQueueOwner(ctx, gitRoot, 4242, "live-owner"); err != nil {
		t.Fatalf("claim terminal Delivery Queue: %v", err)
	}
	if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"replacement"}); err == nil {
		t.Fatal("owned terminal Delivery Queue was replaced")
	}
	if released, err := runStore.ReleaseDeliveryQueueOwner(ctx, gitRoot, 4242, "live-owner"); err != nil || !released {
		t.Fatalf("release terminal Delivery Queue: released=%v err=%v", released, err)
	}
	replaced, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"replacement"})
	if err != nil {
		t.Fatalf("replace terminal unowned Delivery Queue: %v", err)
	}
	if len(replaced.Items) != 1 || replaced.Items[0].SpecSlug != "replacement" {
		t.Fatalf("replacement Delivery Queue = %+v", replaced)
	}
}

func TestOpenMigratesV14DeliveryQueueAddingOwnerAndItemBranch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	freshHomeDir := t.TempDir()
	fresh := openTestStore(t, ctx, freshHomeDir)
	freshSchema := readRunDatabaseSchema(t, fresh)
	closeStore(t, fresh)

	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	queue, err := runStore.CreateDeliveryQueue(ctx, "/tmp/existing-delivery", []string{"existing-spec"})
	if err != nil {
		t.Fatalf("create existing Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = DeliveryStageParked
	item.Blocker = "review-stale"
	if err := runStore.UpdateDeliveryQueueItem(ctx, queue.GitRoot, item); err != nil {
		t.Fatalf("update existing Delivery Queue item: %v", err)
	}
	closeStore(t, runStore)
	downgradeDeliveryOwnerSchemaFixture(t, ctx, homeDir)

	reopened := openTestStore(t, ctx, homeDir)
	defer closeStore(t, reopened)
	persisted, found, err := reopened.DeliveryQueue(ctx, queue.GitRoot)
	if err != nil || !found {
		t.Fatalf("read migrated Delivery Queue: found=%v err=%v", found, err)
	}
	if len(persisted.Items) != 1 || persisted.Items[0].Stage != DeliveryStageParked || persisted.Items[0].Blocker != "review-stale" {
		t.Fatalf("migrated Delivery Queue = %+v", persisted)
	}
	if persisted.Items[0].Branch != "" {
		t.Fatalf("migrated Delivery Queue item branch = %q, want empty", persisted.Items[0].Branch)
	}
	if persisted.OwnerPID != 0 || persisted.OwnerIdentity != "" {
		t.Fatalf("migrated Delivery Queue owner = pid:%d identity:%q", persisted.OwnerPID, persisted.OwnerIdentity)
	}
	if migratedSchema := readRunDatabaseSchema(t, reopened); migratedSchema != freshSchema {
		t.Fatalf("migrated schema differs from fresh schema:\n--- migrated ---\n%s\n--- fresh ---\n%s", migratedSchema, freshSchema)
	}
}

func readRunDatabaseSchema(t *testing.T, store *Store) string {
	t.Helper()
	rows, err := store.db.Query(`
SELECT type, name, tbl_name, sql
FROM sqlite_schema
WHERE sql IS NOT NULL
ORDER BY type, name`)
	if err != nil {
		t.Fatalf("read Run Database schema: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Fatalf("close Run Database schema rows: %v", err)
		}
	}()

	schema := ""
	for rows.Next() {
		var objectType string
		var name string
		var tableName string
		var sqlText string
		if err := rows.Scan(&objectType, &name, &tableName, &sqlText); err != nil {
			t.Fatalf("scan Run Database schema: %v", err)
		}
		schema += objectType + "\x00" + name + "\x00" + tableName + "\x00" + sqlText + "\n"
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate Run Database schema: %v", err)
	}
	return schema
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

func downgradeDeliveryOwnerSchemaFixture(t *testing.T, ctx context.Context, homeDir string) {
	t.Helper()
	db, err := sql.Open("sqlite", writerDSN(DatabasePath(homeDir)))
	if err != nil {
		t.Fatalf("open Delivery Queue owner migration fixture: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close Delivery Queue owner migration fixture: %v", err)
		}
	}()

	for _, statement := range []string{
		`ALTER TABLE delivery_queues DROP COLUMN owner_identity`,
		`ALTER TABLE delivery_queues DROP COLUMN owner_pid`,
		`PRAGMA user_version = 14`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("build schema 14 Delivery Queue owner fixture: %v", err)
		}
	}
}
