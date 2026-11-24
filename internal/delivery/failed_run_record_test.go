// Suite: Implement executor Run provenance.
// Invariant: a failed or stopped executor leaves its created Run on the item.
// Boundary IN: production Engine and disposable Run Database.
// Boundary OUT: Implement child represented by the package runner fake.
package delivery

import (
	"context"
	"errors"
	"testing"

	"roundfix/internal/store"
)

func TestTheItemRecordsTheRunOfAFailedImplementExecutor(t *testing.T) {
	t.Parallel()
	for _, stopped := range []bool{false, true} {
		name := "failed_executor_parks_with_its_run"
		if stopped {
			name = "stopped_owner_still_records_the_run"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			db := openDeliveryEngineStore(t, t.Context())
			queue, err := db.CreateDeliveryQueue(t.Context(), "/repo", []string{"spec-one"})
			if err != nil {
				t.Fatal(err)
			}
			item := queue.Items[0]
			item.Stage = store.DeliveryStageRunning
			item.RunID = "run-earlier"
			item.Branch = "deliver/spec-one"
			item.Worktree = "/worktrees/spec-one"
			recordDeliveryItemWorkspace(t, t.Context(), db, "/repo", &item)
			if err := db.UpdateDeliveryQueueItem(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			workflow := newFakeDeliveryWorkflow()
			engine := newTestDeliveryEngine(db, workflow, newFakeDeliveryBoundary())
			engine.runner = slotRunnerFunc(func(context.Context, string, string) (RunResult, error) {
				if stopped {
					cancel()
				}
				return RunResult{RunID: "run-failed"}, errors.New("roundfix implement failed with exit code 1")
			})
			_, _ = engine.Run(ctx, "/repo")
			got := readDeliveryQueue(t, t.Context(), db, "/repo").Items[0]
			if got.RunID != "run-failed" {
				t.Fatalf("Run = %q, want run-failed", got.RunID)
			}
			if stopped {
				if got.Stage == store.DeliveryStageParked {
					t.Fatalf("stopped item parked: %+v", got)
				}
			} else if got.Stage != store.DeliveryStageParked || got.Blocker != "delivery-error: run Implement executor: roundfix implement failed with exit code 1" || ClassifyPark(queue, got).Class != ParkClassEnvironment {
				t.Fatalf("failed item = %+v", got)
			}
		})
	}
}
