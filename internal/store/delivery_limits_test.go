// Suite: Delivery Queue limits and warnings.
// Invariant: queue limits and item retry metadata survive process boundaries and guard retries atomically.
// Boundary IN: the real SQLite-backed Delivery Queue store and its schema migration.
// Boundary OUT: Delivery Engine behavior, owned by internal/delivery/limits_test.go.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDeliveryQueueRecordsItsLimits(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	deadline := time.Date(2026, time.September, 29, 3, 4, 5, 987654321, time.FixedZone("west", -3*60*60))

	runStore := openTestStore(t, ctx, homeDir)
	queue, err := runStore.CreateDeliveryQueueWithLimits(ctx, "/tmp/limited-delivery", []string{"limited-spec"}, DeliveryQueueLimits{
		Deadline:   deadline,
		MaxRetries: 3,
	})
	if err != nil {
		t.Fatalf("create limited Delivery Queue: %v", err)
	}
	closeStore(t, runStore)

	reopened := openTestStore(t, ctx, homeDir)
	defer closeStore(t, reopened)
	persisted, found, err := reopened.DeliveryQueue(ctx, queue.GitRoot)
	if err != nil || !found {
		t.Fatalf("read limited Delivery Queue after reopening: found=%v err=%v", found, err)
	}
	wantDeadline := time.Unix(deadline.Unix(), 0).UTC()
	if !persisted.Limits.Deadline.Equal(wantDeadline) || persisted.Limits.Deadline.Location() != time.UTC {
		t.Fatalf("persisted deadline = %s in %s, want %s in UTC", persisted.Limits.Deadline, persisted.Limits.Deadline.Location(), wantDeadline)
	}
	if persisted.Limits.MaxRetries != 3 {
		t.Fatalf("persisted max retries = %d, want 3", persisted.Limits.MaxRetries)
	}
}

func TestDeliveryQueueWithoutLimitsRecordsNone(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()

	runStore := openTestStore(t, ctx, homeDir)
	queue, err := runStore.CreateDeliveryQueue(ctx, "/tmp/unlimited-delivery", []string{"unlimited-spec"})
	if err != nil {
		t.Fatalf("create Delivery Queue without limits: %v", err)
	}
	closeStore(t, runStore)

	reopened := openTestStore(t, ctx, homeDir)
	defer closeStore(t, reopened)
	persisted, found, err := reopened.DeliveryQueue(ctx, queue.GitRoot)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue without limits after reopening: found=%v err=%v", found, err)
	}
	if !persisted.Limits.Deadline.IsZero() || persisted.Limits.MaxRetries != 0 {
		t.Fatalf("persisted limits = %+v, want zero limits", persisted.Limits)
	}
}

func TestDeliveryQueueRefusesANegativeRetryLimit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/negative-retry-limit"

	_, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{"negative-limit"}, DeliveryQueueLimits{MaxRetries: -1})
	if err == nil || !strings.Contains(err.Error(), "Max retries") {
		t.Fatalf("negative retry limit error = %v, want named refusal", err)
	}
	if _, found, readErr := runStore.DeliveryQueue(ctx, gitRoot); readErr != nil || found {
		t.Fatalf("queue after negative retry limit: found=%v err=%v", found, readErr)
	}
}

func TestRetryIncrementsTheItemRetryCount(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/retry-count"
	queue, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{"retry-count"}, DeliveryQueueLimits{MaxRetries: 2})
	if err != nil {
		t.Fatalf("create limited Delivery Queue: %v", err)
	}
	parked := queue.Items[0]
	parked.Stage = DeliveryStageParked
	parked.Blocker = "review-findings"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, parked); err != nil {
		t.Fatalf("park Delivery Queue item: %v", err)
	}
	target := parked
	target.Stage = DeliveryStageRunning
	if _, _, err := runStore.RetryDeliveryQueueItem(ctx, gitRoot, target, parked.Blocker); err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}

	got := readDeliveryQueueItem(t, ctx, runStore, gitRoot)
	if got.RetryCount != 1 || got.Stage != DeliveryStageRunning || got.Blocker != "" {
		t.Fatalf("retried item = %+v, want running with retry count 1", got)
	}
}

func TestRetryAtTheLimitIsRefusedAndLeavesTheItemUnchanged(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	const gitRoot = "/tmp/retry-limit"
	queue, err := runStore.CreateDeliveryQueueWithLimits(ctx, gitRoot, []string{"retry-limit"}, DeliveryQueueLimits{MaxRetries: 1})
	if err != nil {
		t.Fatalf("create limited Delivery Queue: %v", err)
	}
	parked := queue.Items[0]
	parked.Stage = DeliveryStageParked
	parked.Blocker = "first-blocker"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, parked); err != nil {
		t.Fatalf("park Delivery Queue item for first retry: %v", err)
	}
	target := parked
	target.Stage = DeliveryStageRunning
	if _, _, err := runStore.RetryDeliveryQueueItem(ctx, gitRoot, target, parked.Blocker); err != nil {
		t.Fatalf("use permitted retry: %v", err)
	}

	parked = readDeliveryQueueItem(t, ctx, runStore, gitRoot)
	parked.Stage = DeliveryStageParked
	parked.Blocker = "second-blocker"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, parked); err != nil {
		t.Fatalf("park Delivery Queue item at its retry limit: %v", err)
	}
	before := readDeliveryQueueItem(t, ctx, runStore, gitRoot)
	target = before
	target.Stage = DeliveryStageReviewing

	_, _, err = runStore.RetryDeliveryQueueItem(ctx, gitRoot, target, before.Blocker)
	if !errors.Is(err, ErrDeliveryRetryLimit) {
		t.Fatalf("retry at limit error = %v, want ErrDeliveryRetryLimit", err)
	}
	if after := readDeliveryQueueItem(t, ctx, runStore, gitRoot); !reflect.DeepEqual(after, before) {
		t.Fatalf("item changed after retry-limit refusal:\n got: %#v\nwant: %#v", after, before)
	}
}

func TestThePreviousSchemaGainsTheLimitColumns(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	queue, err := runStore.CreateDeliveryQueue(ctx, "/tmp/previous-limit-schema", []string{"existing-spec"})
	if err != nil {
		t.Fatalf("create existing Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = DeliveryStageParked
	item.Blocker = "review-findings"
	if err := runStore.UpdateDeliveryQueueItem(ctx, queue.GitRoot, item); err != nil {
		t.Fatalf("update existing Delivery Queue item: %v", err)
	}
	closeStore(t, runStore)
	downgradeDeliveryLimitsSchemaFixture(t, ctx, homeDir)

	reopened := openTestStore(t, ctx, homeDir)
	defer closeStore(t, reopened)
	version, err := reopened.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migrated schema version: %v", err)
	}
	if version != schemaVersion {
		t.Fatalf("migrated schema version = %d, want %d", version, schemaVersion)
	}
	persisted, found, err := reopened.DeliveryQueue(ctx, queue.GitRoot)
	if err != nil || !found {
		t.Fatalf("read migrated Delivery Queue: found=%v err=%v", found, err)
	}
	if len(persisted.Items) != 1 || persisted.Items[0].Stage != DeliveryStageParked || persisted.Items[0].Blocker != item.Blocker {
		t.Fatalf("migrated Delivery Queue lost its row: %+v", persisted)
	}
	if !persisted.Limits.Deadline.IsZero() || persisted.Limits.MaxRetries != 0 || persisted.Items[0].RetryCount != 0 || persisted.Items[0].Warning != "" {
		t.Fatalf("migrated defaults = limits:%+v item:%+v", persisted.Limits, persisted.Items[0])
	}
}

func TestDeliveryQueueItemWarningRoundTrips(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	const gitRoot = "/tmp/item-warning"
	runStore := openTestStore(t, ctx, homeDir)
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"warned-spec", "quiet-spec"})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	warned := queue.Items[0]
	warned.Stage = DeliveryStageParked
	warned.Blocker = "review-findings"
	warned.Warning = "premise-changed: internal/delivery/engine.go (merge abc123)"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, warned); err != nil {
		t.Fatalf("record Delivery Queue item warning: %v", err)
	}
	closeStore(t, runStore)

	reopened := openTestStore(t, ctx, homeDir)
	persisted, found, err := reopened.DeliveryQueue(ctx, gitRoot)
	if err != nil || !found {
		t.Fatalf("read warned Delivery Queue after reopening: found=%v err=%v", found, err)
	}
	if persisted.Items[0].Warning != warned.Warning || persisted.Items[1].Warning != "" {
		t.Fatalf("persisted warnings = %q and %q", persisted.Items[0].Warning, persisted.Items[1].Warning)
	}
	target := persisted.Items[0]
	target.Stage = DeliveryStageRunning
	if _, _, err := reopened.RetryDeliveryQueueItem(ctx, gitRoot, target, warned.Blocker); err != nil {
		t.Fatalf("retry warned Delivery Queue item: %v", err)
	}
	closeStore(t, reopened)

	reopenedAgain := openTestStore(t, ctx, homeDir)
	defer closeStore(t, reopenedAgain)
	afterRetry, found, err := reopenedAgain.DeliveryQueue(ctx, gitRoot)
	if err != nil || !found {
		t.Fatalf("read warned Delivery Queue after retry: found=%v err=%v", found, err)
	}
	if afterRetry.Items[0].Warning != warned.Warning || afterRetry.Items[1].Warning != "" {
		t.Fatalf("warnings after retry = %q and %q", afterRetry.Items[0].Warning, afterRetry.Items[1].Warning)
	}
}

func readDeliveryQueueItem(t *testing.T, ctx context.Context, runStore *Store, gitRoot string) DeliveryQueueItem {
	t.Helper()
	queue, found, err := runStore.DeliveryQueue(ctx, gitRoot)
	if err != nil || !found || len(queue.Items) != 1 {
		t.Fatalf("read one-item Delivery Queue: found=%v items=%d err=%v", found, len(queue.Items), err)
	}
	return queue.Items[0]
}

func downgradeDeliveryLimitsSchemaFixture(t *testing.T, ctx context.Context, homeDir string) {
	t.Helper()
	db, err := sql.Open("sqlite", writerDSN(DatabasePath(homeDir)))
	if err != nil {
		t.Fatalf("open Delivery Queue limits migration fixture: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close Delivery Queue limits migration fixture: %v", err)
		}
	}()

	statements := []string{
		`ALTER TABLE delivery_queue_items DROP COLUMN warning`,
		`ALTER TABLE delivery_queue_items DROP COLUMN retry_count`,
		`ALTER TABLE delivery_queues DROP COLUMN max_retries`,
		`ALTER TABLE delivery_queues DROP COLUMN deadline_unix`,
		fmt.Sprintf(`PRAGMA user_version = %d`, schemaVersion-1),
	}
	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			t.Fatalf("build previous Delivery Queue limits schema fixture with %q: %v", statement, err)
		}
	}
}
