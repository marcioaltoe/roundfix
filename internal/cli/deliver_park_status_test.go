// Suite: public delivery park status.
// Boundary IN: CLI status and the real queue store.
// Boundary OUT: no owner or GitHub action is started by status.
package cli

import (
	"bytes"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

func TestDeliverStatusPrintsAParkLinePerParkedItem(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repoDir, []string{"first", "active", "last"})
	if err != nil {
		t.Fatal(err)
	}
	for i := range queue.Items {
		item := queue.Items[i]
		if item.SpecSlug == "active" {
			continue
		}
		item.Stage = store.DeliveryStageParked
		item.Blocker = delivery.BlockerReviewStale
		if item.SpecSlug == "first" {
			item.Blocker = delivery.BlockerChecksTimeout
			item.Warning = "premise-changed: internal/example.go"
		}
		if err := runStore.UpdateDeliveryQueueItem(t.Context(), repoDir, item); err != nil {
			t.Fatal(err)
		}
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "first\tparked\tchecks-timeout\t-\nactive\tqueued\t-\t-\nlast\tparked\treview-stale\t-\n" +
		"Warning: first premise-changed: internal/example.go\n" +
		"Park: first environment: resolve the blocker, then run roundfix deliver retry first\n" +
		"Park: last review: resolve the blocker, then run roundfix deliver retry last\n" +
		"Limits: deadline none, retries per item none, concurrency 1, tokens none\n" +
		"Usage: no Runs recorded\n" +
		"Pending question: first parked checks-timeout\nAnswer: resolve the blocker, then run roundfix deliver retry first\nWaiting behind it: 1 parked item(s)\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q want=%q", code, stderr.String(), stdout.String(), want)
	}
}

func TestDeliverStatusReproducesSurfaceTranscriptOne(t *testing.T) {
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repo, []string{"0204-first", "0205-second"})
	if err != nil {
		t.Fatal(err)
	}
	const path = "internal/baseline/assets/profiles/standard-typescript-monorepo.json"
	for index, item := range queue.Items {
		item.Stage = store.DeliveryStageParked
		if index == 0 {
			item.Blocker = "pull-request-conflict: " + path
			item.Branch = "roundfix/deliver-first"
			item.Worktree = "/worktrees/0204-first"
			if _, _, _, err := runStore.RecordDeliveryQueueItemWorktree(t.Context(), repo, item.SpecSlug, item.Branch, item.Worktree); err != nil {
				t.Fatal(err)
			}
		} else {
			item.Blocker = "prerequisite-unmerged: 0204-first"
		}
		if err := runStore.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
			t.Fatal(err)
		}
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	next := "merge the default branch into the item branch in /worktrees/0204-first, resolve " + path + ", commit, then run roundfix deliver retry 0204-first"
	want := "0204-first\tparked\tpull-request-conflict: " + path + "\t/worktrees/0204-first\n" +
		"0205-second\tparked\tprerequisite-unmerged: 0204-first\t-\n" +
		"Park: 0204-first conflict: " + next + "\n" +
		"Park: 0205-second dependency: run roundfix deliver retry 0204-first; 0205-second returns to the queue once 0204-first is merged\n" +
		"Limits: deadline none, retries per item none, concurrency 1, tokens none\n" +
		"Usage: no Runs recorded\n" +
		"Pending question: 0204-first parked pull-request-conflict: " + path + "\n" +
		"Answer: " + next + "\nWaiting behind it: 1 parked item(s)\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q want=%q", code, stderr.String(), stdout.String(), want)
	}
}
