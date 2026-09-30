// Suite: Delivery Queue owner warning persistence.
// Invariant: an owner warning is recorded and logged without blocking the item's transition to running.
// Boundary IN: the real SQLite Delivery Queue store and delivery engine transitions.
// Boundary OUT: item worktrees, revalidation, and the Run executor are represented by existing fakes.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestAdvanceItemRecordsTheOwnerWarningAndKeepsRunning(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const (
		gitRoot  = "/repo-owner-warning"
		slug     = "owner-warning-spec"
		premise  = "premise-changed: internal/premise.go (merge earlier-merge)"
		owner    = "owner-older-than-main: owner build 111111111111 predates starting main 222222222222"
		combined = premise + "; " + owner
	)
	queue, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{slug})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := newFakeDeliveryWorkflow()
	stopAtRun := errors.New("stop after running transition")
	workflow.runErrors[slug] = stopAtRun
	revalidator := &recordingItemRevalidator{results: map[string]Revalidation{
		slug: {
			ChangedPremises: []string{"internal/premise.go"},
			ChangedBy:       []string{"earlier-merge"},
			OwnerWarning:    owner,
		},
	}}
	var log bytes.Buffer
	engine := newRevalidationTestEngine(runStore, workflow, revalidator, &log)
	item := queue.Items[0]

	err = engine.advanceItem(ctx, gitRoot, &item, []string{})
	if !errors.Is(err, stopAtRun) {
		t.Fatalf("advance item error = %v, want Run boundary stop", err)
	}
	if item.Warning != combined || item.Stage != store.DeliveryStageRunning || item.Blocker != "" {
		t.Fatalf("owner-warned item = %+v, want running with %q", item, combined)
	}
	if !slicesContain(workflow.events[slug], "run") {
		t.Fatalf("item events = %v, want Run reached after owner warning", workflow.events[slug])
	}
	wantLog := "roundfix: warning: Delivery Queue item " + slug + ": " + combined + "\n"
	if log.String() != wantLog {
		t.Fatalf("owner warning log = %q, want %q", log.String(), wantLog)
	}
	stored := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if stored.Warning != combined || stored.Stage != store.DeliveryStageRunning {
		t.Fatalf("stored owner-warned item = %+v, want running with warning", stored)
	}
}

func TestRetryKeepsTheRecordedOwnerWarning(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openDeliveryEngineStore(t, ctx)
	const gitRoot = "/repo-owner-warning-retry"
	item := seedNoRunRetryItem(t, ctx, runStore, gitRoot, "owner-warning-retry")
	item.Warning = "owner-older-than-main: owner build old predates starting main base"
	if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
		t.Fatalf("seed owner warning: %v", err)
	}
	recovery := &fakeItemRecovery{
		states:      []ItemState{{UnfinishedTasks: []string{"task_01"}}, {UnfinishedTasks: []string{"task_01"}}},
		carryResult: CarryForwardResult{RunID: "run-retried"},
	}
	revalidator := &recordingItemRevalidator{results: map[string]Revalidation{
		item.SpecSlug: {OwnerWarning: "owner-older-than-main: recomputed warning"},
	}}
	engine := newRetryRevalidationTestEngine(runStore, recovery, revalidator)

	if _, err := engine.Retry(ctx, gitRoot, item.SpecSlug); err != nil {
		t.Fatalf("retry Delivery Queue item: %v", err)
	}
	got := readDeliveryQueue(t, ctx, runStore, gitRoot).Items[0]
	if got.Warning != item.Warning || got.Stage != store.DeliveryStageRunning {
		t.Fatalf("retried item = %+v, want unchanged owner warning and running stage", got)
	}
	if strings.Contains(got.Warning, "recomputed") {
		t.Fatalf("retried owner warning = %q, want recorded warning preserved", got.Warning)
	}
}
