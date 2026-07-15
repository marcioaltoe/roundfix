//go:build unix

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestRunImplementReclaimsDeadOwnerActiveRun(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	withImplementCollaborators(t, runner)
	pid := reapedCLIProcessPID(t)
	blocking := seedImplementActiveRun(t, homeDir, repoDir, "0002-other-spec", pid)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected implement to proceed after reclaim, got %d stderr=%q", code, stderr.String())
	}
	assertOrphanWarning(t, stderr.String(), blocking.ID, pid)
	if !strings.Contains(stdout.String(), "Clean: all 1 Task(s) completed.") {
		t.Fatalf("expected implement success, got %q", stdout.String())
	}
	assertReclaimedRunInCLI(t, homeDir, blocking.ID, pid)
}

func TestRunSettleReclaimsDeadOwnerActiveRun(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:           "task_01",
		status:       string(spec.StatusFailed),
		verification: []string{"test -f done.txt"},
	}})
	mustWrite(t, filepath.Join(repoDir, "done.txt"), "settled work\n")
	pid := reapedCLIProcessPID(t)
	blocking := seedReviewActiveRun(t, homeDir, repoDir, store.KindResolve, pid)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected settle to proceed after reclaim, got %d stderr=%q", code, stderr.String())
	}
	assertOrphanWarning(t, stderr.String(), blocking.ID, pid)
	if !strings.Contains(stdout.String(), "settled task_01 completed") {
		t.Fatalf("expected settle success, got %q", stdout.String())
	}
	assertReclaimedRunInCLI(t, homeDir, blocking.ID, pid)
}

func TestReviewFetchReclaimsDeadOwnerActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	pid := reapedCLIProcessPID(t)
	blocking := seedReviewActiveRun(t, homeDir, repoDir, store.KindWatch, pid)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected fetch to proceed after reclaim, got %d stderr=%q", code, stderr.String())
	}
	assertOrphanWarning(t, stderr.String(), blocking.ID, pid)
	if !strings.Contains(stdout.String(), "Fetch complete") {
		t.Fatalf("expected fetch success, got %q", stdout.String())
	}
	assertReclaimedRunInCLI(t, homeDir, blocking.ID, pid)
}

// A terminal Run with a dead owner must never be "reclaimed" by stop: the
// existing terminal-Run error stands, and no state or lock changes.
func TestStopTerminalRunWithDeadOwnerKeepsTerminalError(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	pid := reapedCLIProcessPID(t)
	seeded := seedReviewActiveRun(t, homeDir, repoDir, store.KindResolve, pid)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if _, err := runStore.CompleteRun(context.Background(), seeded.ID, store.StateFailed); err != nil {
		t.Fatalf("complete seeded Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", seeded.ID}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected terminal-Run stop refusal exit 2, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if strings.Contains(stderr.String(), "reclaimed orphaned Active Run") {
		t.Fatalf("expected no reclamation of a terminal Run, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "cannot record Stop Request for terminal Run") {
		t.Fatalf("expected the terminal-Run stop error, got %q", stderr.String())
	}
}

func seedImplementActiveRun(t *testing.T, homeDir string, repoDir string, specSlug string, pid int) store.Run {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close Run Database: %v", err)
		}
	}()
	run, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/orphan-owner",
		SpecSlug:    specSlug,
		OwnerPID:    pid,
	})
	if err != nil {
		t.Fatalf("seed implement active Run: %v", err)
	}
	return run
}

func seedReviewActiveRun(t *testing.T, homeDir string, repoDir string, kind string, pid int) store.Run {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close Run Database: %v", err)
		}
	}()
	run, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           kind,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
		OwnerPID:       pid,
	})
	if err != nil {
		t.Fatalf("seed review active Run: %v", err)
	}
	return run
}

func reapedCLIProcessPID(t *testing.T) int {
	t.Helper()
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start owner process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for owner process: %v", err)
	}
	if store.ProcessAlive(pid) {
		t.Skipf("owner pid %d is still alive or was reused", pid)
	}
	return pid
}

func assertOrphanWarning(t *testing.T, stderr string, runID string, pid int) {
	t.Helper()
	for _, want := range []string{
		"reclaimed orphaned Active Run " + runID,
		fmt.Sprintf("owner process %d not running; lock reclaimed", pid),
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("expected orphan warning to contain %q, got %q", want, stderr)
		}
	}
}

func assertReclaimedRunInCLI(t *testing.T, homeDir string, runID string, pid int) {
	t.Helper()
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("close Run Database reader: %v", err)
		}
	}()
	run, found, err := reader.Run(context.Background(), runID)
	if err != nil || !found {
		t.Fatalf("read reclaimed Run: found=%v err=%v", found, err)
	}
	if run.State != store.StateFailed || run.CompletedAt == nil {
		t.Fatalf("expected Failed reclaimed Run, got %#v", run)
	}
	events, err := reader.RunEventsAfter(context.Background(), runID, 0, 10)
	if err != nil {
		t.Fatalf("read Run Events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one reclamation event, got %d", len(events))
	}
	var payload struct {
		State    string `json:"state"`
		Reason   string `json:"reason"`
		OwnerPID int    `json:"owner_pid"`
	}
	if err := json.Unmarshal(events[0].Event.Payload, &payload); err != nil {
		t.Fatalf("decode reclamation event: %v", err)
	}
	wantReason := fmt.Sprintf("owner process %d not running; lock reclaimed", pid)
	if payload.State != store.StateFailed || payload.Reason != wantReason || payload.OwnerPID != pid {
		t.Fatalf("unexpected reclamation payload: %#v", payload)
	}
}
