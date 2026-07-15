//go:build unix

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"testing"
)

func TestReclaimOrphanedRunCompletesFailedReleasesLockAndJournalsReason(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	pid := reapedOwnerPID(t)
	req := sampleCreateRunRequest()
	req.OwnerPID = pid
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create active Run: %v", err)
	}
	reason := fmt.Sprintf("owner process %d not running; lock reclaimed", pid)

	if err := runStore.ReclaimOrphanedRun(ctx, run, reason); err != nil {
		t.Fatalf("reclaim orphaned Run: %v", err)
	}
	reclaimed, found, err := runStore.Run(ctx, run.ID)
	if err != nil || !found {
		t.Fatalf("read reclaimed Run: found=%v err=%v", found, err)
	}
	if reclaimed.State != StateFailed || reclaimed.CompletedAt == nil {
		t.Fatalf("expected Failed reclaimed Run with completion timestamp, got %#v", reclaimed)
	}
	if _, found, err := runStore.ActiveRun(ctx, req.HeadRepository, req.HeadBranch); err != nil || found {
		t.Fatalf("expected reclaimed Run lock released, found=%v err=%v", found, err)
	}
	assertReclamationEvent(t, ctx, runStore, run.ID, reason, pid)

	if err := runStore.ReclaimOrphanedRun(ctx, run, reason); err != nil {
		t.Fatalf("reclaim orphaned Run twice: %v", err)
	}
	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events after second reclaim: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected second reclaim to journal no extra events, got %d", len(events))
	}
}

func TestReclaimOrphanedRunConcurrentAttemptsJournalOnce(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	pid := reapedOwnerPID(t)
	req := sampleCreateRunRequest()
	req.OwnerPID = pid
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create active Run: %v", err)
	}
	reason := fmt.Sprintf("owner process %d not running; lock reclaimed", pid)

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- runStore.ReclaimOrphanedRun(ctx, run, reason)
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent reclaim: %v", err)
		}
	}

	reclaimed, found, err := runStore.Run(ctx, run.ID)
	if err != nil || !found {
		t.Fatalf("read reclaimed Run: found=%v err=%v", found, err)
	}
	if reclaimed.State != StateFailed {
		t.Fatalf("expected one terminal Failed Run, got %#v", reclaimed)
	}
	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected exactly one reclamation event, got %d", len(events))
	}
}

func reapedOwnerPID(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start owner process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for owner process: %v", err)
	}
	if ProcessAlive(pid) {
		t.Skipf("owner pid %d is still alive or was reused", pid)
	}
	return pid
}

func assertReclamationEvent(t *testing.T, ctx context.Context, runStore *Store, runID string, reason string, pid int) {
	t.Helper()
	events, err := runStore.RunEventsAfter(ctx, runID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one reclamation event, got %d", len(events))
	}
	event := events[0].Event
	if event.Kind != "daemon.outcome" || event.Summary == "" {
		t.Fatalf("expected daemon outcome reclamation event, got %#v", event)
	}
	var payload struct {
		State    string `json:"state"`
		Reason   string `json:"reason"`
		OwnerPID int    `json:"owner_pid"`
	}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("decode reclamation payload: %v", err)
	}
	if payload.State != StateFailed || payload.Reason != reason || payload.OwnerPID != pid {
		t.Fatalf("unexpected reclamation payload: %#v", payload)
	}
}
