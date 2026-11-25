// Boundary: engine stages and SQLite persistence, with fake workflow and forge.
package delivery

import (
	"context"
	"path/filepath"
	"roundfix/internal/store"
	"testing"
)

func TestAPassedGateLeavesNoLaneWarningAndNoPullRequestParagraph(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := openDeliveryEngineStore(t, ctx)
	root := filepath.Join(t.TempDir(), "repo")
	if _, err := db.CreateDeliveryQueue(ctx, root, []string{"spec-one"}); err != nil {
		t.Fatal(err)
	}
	workflow := newFakeDeliveryWorkflow()
	boundary := newFakeDeliveryBoundary()
	publications := 0
	boundary.findPullRequest = func(req PullRequestRequest) PullRequestResult {
		publications++
		item := readDeliveryQueue(t, ctx, db, root).Items[0]
		if item.Stage != store.DeliveryStagePublishing || item.Warning != "" || req.Body != "deliver spec-one" {
			t.Fatalf("item=%+v body=%q", item, req.Body)
		}
		pr := PullRequest{Number: "1", State: "OPEN", HeadBranch: req.HeadBranch, HeadSHA: boundary.remoteHeads[req.HeadBranch]}
		boundary.pullRequests[req.HeadBranch] = pr
		return PullRequestResult{PullRequest: pr, Created: true}
	}
	if _, err := newTestDeliveryEngine(db, workflow, boundary).Run(ctx, root); err != nil {
		t.Fatal(err)
	}
	if publications != 1 {
		t.Fatalf("publications=%d", publications)
	}
}

func TestAStoredLaneBlockerKeepsItsParkClass(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, blocker, next string
		class               ParkClass
	}{
		{"container_blocker", "linux-lane-container: Docker unavailable", "run roundfix deliver retry spec-one; the Linux Lane was retired, so the gate no longer runs it", ParkClassEnvironment},
		{"change_failure_blocker", "gate-failed: Linux Lane change failure: go test exited 1", "", ParkClassFinding},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			item := store.DeliveryQueueItem{SpecSlug: "spec-one", Blocker: row.blocker}
			got := ClassifyPark(store.DeliveryQueue{}, item)
			want := row.next
			if want == "" {
				item.Blocker = BlockerGateFailed
				want = ClassifyPark(store.DeliveryQueue{}, item).Next
			}
			if got.Class != row.class || got.Next != want {
				t.Fatalf("park=%+v", got)
			}
		})
	}
	t.Run("skip_segment_survives_the_low_disk_rewrite", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()
		db := openDeliveryEngineStore(t, ctx)
		root := filepath.Join(t.TempDir(), "repo")
		q, err := db.CreateDeliveryQueue(ctx, root, []string{"spec-one"})
		if err != nil {
			t.Fatal(err)
		}
		item := q.Items[0]
		item.Warning = "premise-changed: file; low-disk: old; linux-lane-skipped: image absent; detail"
		if err := db.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
			t.Fatal(err)
		}
		engine := NewEngine(db, EngineDependencies{FreeSpace: freeSpaceCheckerFunc(func(_ context.Context) []string { return []string{"new location"} })})
		if err := engine.checkFreeSpace(ctx, root, &item, "repository gate"); err != nil {
			t.Fatal(err)
		}
		if got := readDeliveryQueue(t, ctx, db, root).Items[0].Warning; got != "premise-changed: file; linux-lane-skipped: image absent; detail; low-disk: new location" {
			t.Fatalf("warning=%q", got)
		}
	})
}
