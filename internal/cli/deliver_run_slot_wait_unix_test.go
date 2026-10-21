//go:build unix

// Boundary: real owner-process liveness and read-only Run Database queries.
package cli

import (
	"os"
	"os/exec"
	"reflect"
	"sort"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

func TestRunSlotsExcludeAHolderWhoseOwnerIsGone(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	child := exec.CommandContext(t.Context(), "/bin/sh", "-c", "exit 0")
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if child.ProcessState == nil {
			_ = child.Process.Kill()
			_ = child.Wait()
		}
	})
	if err := child.Wait(); err != nil {
		t.Fatal(err)
	}
	deadPID := child.Process.Pid
	if store.ProcessAlive(deadPID) {
		t.Fatalf("reaped child PID %d still alive", deadPID)
	}
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	var want []delivery.RunSlotHolder
	for _, seed := range []struct {
		slug string
		pid  int
	}{{"live", os.Getpid()}, {"gone", deadPID}, {"unrecorded", 0}} {
		run, err := db.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: repo, LocalBranch: "feat/" + seed.slug, SpecSlug: seed.slug, OwnerPID: seed.pid})
		if err != nil {
			t.Fatal(err)
		}
		if seed.slug != "gone" {
			want = append(want, delivery.RunSlotHolder{RunID: run.ID, Repository: run.GitRoot, SpecSlug: seed.slug})
		}
	}
	sort.Slice(want, func(i, j int) bool { return want[i].RunID < want[j].RunID })
	workflow := &commandDeliveryWorkflow{store: db}
	workflow.loaded.Config.Runs.MaxActive = 2
	slots, err := workflow.RunSlots(t.Context())
	if err != nil || slots.MaxActive != 2 || !reflect.DeepEqual(slots.Holders, want) {
		t.Fatalf("slots=%+v error=%v want=%+v", slots, err, want)
	}
	runs, err := db.ActiveImplementRuns(t.Context())
	if err != nil || len(runs) != 3 {
		t.Fatalf("slot read reclaimed Runs: %+v %v", runs, err)
	}
}
