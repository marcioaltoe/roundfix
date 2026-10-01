// Suite: prerequisite delivery transitions.
// Boundary IN: engine and real SQLite queue. Boundary OUT: Git and delivery actions.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"roundfix/internal/store"
)

type prerequisiteReaderFunc func(context.Context, string, string) ([]string, error)

func (f prerequisiteReaderFunc) UnmetPrerequisites(ctx context.Context, root, slug string) ([]string, error) {
	return f(ctx, root, slug)
}

func prerequisiteEngine(t *testing.T, slugs ...string) (*Engine, *fakeDeliveryWorkflow) {
	t.Helper()
	db := openDeliveryEngineStore(t, t.Context())
	if _, err := db.CreateDeliveryQueue(t.Context(), "/repo", slugs); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	return newTestDeliveryEngine(db, workflow, newFakeDeliveryBoundary()), workflow
}

func TestAnItemWaitsForAPrerequisiteAheadInTheQueue(t *testing.T) {
	engine, workflow := prerequisiteEngine(t, "dependent", "prerequisite")
	engine.prerequisites = prerequisiteReaderFunc(func(_ context.Context, _, slug string) ([]string, error) {
		if slug == "dependent" {
			return []string{"prerequisite"}, nil
		}
		return nil, nil
	})
	var log bytes.Buffer
	engine.log = &log
	if _, err := engine.Run(t.Context(), "/repo"); err != nil {
		t.Fatal(err)
	}
	queue := readDeliveryQueue(t, t.Context(), engine.store, "/repo")
	if queue.Items[0].Stage != store.DeliveryStageQueued || queue.Items[0].Worktree != "" || len(workflow.events["dependent"]) != 0 {
		t.Fatalf("dependent = %+v actions=%v", queue.Items[0], workflow.events)
	}
	if queue.Items[1].Stage != store.DeliveryStageMerged || !strings.Contains(log.String(), "prerequisite wait") {
		t.Fatalf("queue=%+v log=%s", queue, &log)
	}
}

func TestAnItemParksWhenItsPrerequisiteIsParked(t *testing.T) {
	for _, prerequisiteStage := range []store.DeliveryStage{store.DeliveryStageParked, store.DeliveryStageMerged, "absent"} {
		t.Run(string(prerequisiteStage), func(t *testing.T) {
			engine, workflow := prerequisiteEngine(t, "dependent", "prerequisite")
			queue := readDeliveryQueue(t, t.Context(), engine.store, "/repo")
			item := queue.Items[1]
			item.Stage = store.DeliveryStageParked
			item.Blocker = BlockerRunUnresolved
			if prerequisiteStage == store.DeliveryStageMerged {
				item.Stage = prerequisiteStage
				item.Blocker = ""
			}
			if err := engine.store.UpdateDeliveryQueueItem(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			required := "prerequisite"
			if prerequisiteStage == "absent" {
				required = "outside"
			}
			engine.prerequisites = prerequisiteReaderFunc(func(_ context.Context, _, slug string) ([]string, error) {
				if slug == "dependent" {
					return []string{required}, nil
				}
				return nil, nil
			})
			if _, err := engine.Run(t.Context(), "/repo"); err != nil {
				t.Fatal(err)
			}
			item = readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]
			if item.Stage != store.DeliveryStageParked || item.Blocker != BlockerPrerequisiteUnmerged+": "+required || item.Worktree != "" || len(workflow.events["dependent"]) != 0 {
				t.Fatalf("item=%+v actions=%v", item, workflow.events)
			}
		})
	}
}

func TestTheOwnerReleasesADependencyParkWhenThePrerequisiteMerges(t *testing.T) {
	for _, releaseAtStart := range []bool{false, true} {
		t.Run(map[bool]string{true: "start of pass", false: "after merge"}[releaseAtStart], func(t *testing.T) {
			engine, _ := prerequisiteEngine(t, "dependent", "prerequisite")
			queue := readDeliveryQueue(t, t.Context(), engine.store, "/repo")
			item := queue.Items[0]
			item.Stage = store.DeliveryStageParked
			item.Blocker = BlockerPrerequisiteUnmerged + ": prerequisite"
			// Retry counts belong to the retry transition, not ordinary updates.
			for retry := 0; retry < 3; retry++ {
				if err := engine.store.UpdateDeliveryQueueItem(t.Context(), "/repo", item); err != nil {
					t.Fatal(err)
				}
				queued := item
				queued.Stage = store.DeliveryStageQueued
				if _, _, err := engine.store.RetryDeliveryQueueItem(t.Context(), "/repo", queued, item.Blocker); err != nil {
					t.Fatal(err)
				}
			}
			if err := engine.store.UpdateDeliveryQueueItem(t.Context(), "/repo", item); err != nil {
				t.Fatal(err)
			}
			engine.prerequisites = prerequisiteReaderFunc(func(ctx context.Context, _, slug string) ([]string, error) {
				if slug != "dependent" || releaseAtStart {
					return nil, nil
				}
				queue, _, err := engine.store.DeliveryQueue(ctx, "/repo")
				if err != nil {
					return nil, err
				}
				if queue.Items[1].Stage == store.DeliveryStageMerged {
					return nil, nil
				}
				return []string{"prerequisite"}, nil
			})
			var log bytes.Buffer
			engine.log = &log
			if _, err := engine.Run(t.Context(), "/repo"); err != nil {
				t.Fatal(err)
			}
			item = readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]
			want := store.DeliveryStageQueued
			if releaseAtStart {
				want = store.DeliveryStageMerged
			}
			if item.Stage != want || item.Blocker != "" || item.RetryCount != 3 || !strings.Contains(log.String(), "prerequisite release") {
				t.Fatalf("item=%+v log=%s", item, &log)
			}
			if !releaseAtStart {
				if _, err := engine.Run(t.Context(), "/repo"); err != nil {
					t.Fatal(err)
				}
				if readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0].Stage != store.DeliveryStageMerged {
					t.Fatal("released item did not start next pass")
				}
			}
		})
	}
}

func TestRetryReturnsADependencyParkToTheQueue(t *testing.T) {
	engine, workflow := prerequisiteEngine(t, "dependent")
	queue := readDeliveryQueue(t, t.Context(), engine.store, "/repo")
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = BlockerPrerequisiteUnmerged + ": prerequisite"
	if err := engine.store.UpdateDeliveryQueueItem(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	// Nil workspace/recovery/revalidator proves none is needed for this retry.
	engine.workspace = nil
	engine.recovery = nil
	engine.revalidator = nil
	result, err := engine.Retry(t.Context(), "/repo", "dependent")
	if err != nil {
		t.Fatal(err)
	}
	item = readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]
	if result.Stage != store.DeliveryStageQueued || item.Stage != store.DeliveryStageQueued || item.Blocker != "" || len(workflow.events) != 0 {
		t.Fatalf("result=%+v item=%+v", result, item)
	}
}

func TestAnItemWithoutPrerequisitesStartsAsBefore(t *testing.T) {
	for _, readerEnabled := range []bool{false, true} {
		engine, workflow := prerequisiteEngine(t, "plain")
		if readerEnabled {
			engine.prerequisites = prerequisiteReaderFunc(func(context.Context, string, string) ([]string, error) { return nil, nil })
		}
		if _, err := engine.Run(t.Context(), "/repo"); err != nil {
			t.Fatal(err)
		}
		item := readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]
		if item.Stage != store.DeliveryStageMerged || workflow.events["plain"][0] != "create-branch" {
			t.Fatalf("item=%+v actions=%v", item, workflow.events)
		}
	}
}

func TestPrerequisiteReadErrorDoesNotStartAnItem(t *testing.T) {
	engine, workflow := prerequisiteEngine(t, "dependent")
	engine.prerequisites = prerequisiteReaderFunc(func(context.Context, string, string) ([]string, error) { return nil, errors.New("fetch failed") })
	if _, err := engine.Run(t.Context(), "/repo"); err != nil {
		t.Fatal(err)
	}
	item := readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]
	if !strings.HasPrefix(item.Blocker, BlockerDeliveryError) || item.Worktree != "" || len(workflow.events) != 0 {
		t.Fatalf("item=%+v", item)
	}
}

func TestPrerequisiteParkGuidesTheOperator(t *testing.T) {
	item := store.DeliveryQueueItem{SpecSlug: "dependent", Blocker: BlockerPrerequisiteUnmerged + ": prerequisite"}
	for _, parked := range []bool{false, true} {
		queue := store.DeliveryQueue{}
		want := "deliver or merge prerequisite, then run roundfix deliver retry dependent"
		if parked {
			queue.Items = []store.DeliveryQueueItem{{SpecSlug: "prerequisite", Stage: store.DeliveryStageParked}}
			want = "run roundfix deliver retry prerequisite; dependent returns to the queue once prerequisite is merged"
		}
		got := ClassifyPark(queue, item)
		if got.Class != ParkClassDependency || got.Next != want {
			t.Fatalf("classification=%+v", got)
		}
	}
}

func TestDependencyParkRemainsUntilEveryPrerequisiteIsMet(t *testing.T) {
	engine, _ := prerequisiteEngine(t, "dependent")
	item := readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = BlockerPrerequisiteUnmerged + ": first, second"
	if err := engine.store.UpdateDeliveryQueueItem(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	unmet := []string{"first", "second"}
	engine.prerequisites = prerequisiteReaderFunc(func(context.Context, string, string) ([]string, error) { return unmet, nil })
	for _, remaining := range [][]string{{"first", "second"}, {"second"}} {
		unmet = remaining
		if _, err := engine.Run(t.Context(), "/repo"); err != nil {
			t.Fatal(err)
		}
		got := readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]
		if got.Stage != store.DeliveryStageParked || got.Blocker != item.Blocker {
			t.Fatalf("premature release with %v unmet: %+v", unmet, got)
		}
	}
	engine.prerequisites = nil
	if _, err := engine.Run(t.Context(), "/repo"); err != nil {
		t.Fatal(err)
	}
	if got := readDeliveryQueue(t, t.Context(), engine.store, "/repo").Items[0]; got.Stage != store.DeliveryStageParked {
		t.Fatalf("nil reader released %+v", got)
	}
}
