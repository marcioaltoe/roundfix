//go:build unix

package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	pid, ownerWait := startCLIForceStopOwnerProcess(t)
	ownerIdentity, err := store.OwnerProcessIdentity(context.Background(), pid)
	if err != nil {
		t.Fatalf("read genuine owner process identity: %v", err)
	}
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:          store.KindImplement,
		GitRoot:       repoDir,
		LocalBranch:   "ma/force-stop-owner",
		HeadSHA:       "abc123",
		SpecSlug:      "0001-widget-flow",
		Agent:         "codex",
		OwnerPID:      pid,
		OwnerIdentity: ownerIdentity,
	})
	if err != nil {
		t.Fatalf("create active Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	type observation struct {
		ownerAlive bool
		err        error
	}
	observed := make(chan observation, 1)
	monitorCtx, cancelMonitor := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancelMonitor()
	go func() {
		reader, err := store.OpenReader(monitorCtx, homeDir)
		if err != nil {
			observed <- observation{err: err}
			return
		}
		defer func() {
			_ = reader.Close()
		}()
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		for {
			current, found, err := reader.Run(monitorCtx, active.ID)
			if err != nil {
				observed <- observation{err: err}
				return
			}
			if !found {
				observed <- observation{err: fmt.Errorf("Run %s disappeared", active.ID)}
				return
			}
			if current.State == store.StateStopped {
				observed <- observation{ownerAlive: store.ProcessAlive(pid)}
				return
			}
			select {
			case <-monitorCtx.Done():
				observed <- observation{err: monitorCtx.Err()}
				return
			case <-ticker.C:
			}
		}
	}()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("force stop exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	terminal := <-observed
	if terminal.err != nil {
		t.Fatalf("observe terminal Run: %v", terminal.err)
	}
	if terminal.ownerAlive {
		t.Fatalf("Run %s reached Stopped while owner process %d was alive", active.ID, pid)
	}
	select {
	case err := <-ownerWait:
		if err != nil {
			t.Fatalf("owner process exit: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("owner process %d did not exit", pid)
	}
	if store.ProcessAlive(pid) {
		t.Fatalf("owner process %d remained alive after force stop", pid)
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

// TestRunForceStopOwnerPIDReuseFailsClosed exercises the real identity
// comparison: the Run records the PID of a live scratch process together
// with the identity token of a different (exited) owner, so Force Stop must
// refuse before sending any signal, exactly as it must for a reused PID.
func TestRunForceStopOwnerPIDReuseFailsClosed(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	pid, _ := startCLIForceStopOwnerProcess(t)
	request := store.CreateRunRequest{
		Kind:          store.KindImplement,
		GitRoot:       repoDir,
		LocalBranch:   "ma/force-stop-reused-pid",
		HeadSHA:       "abc123",
		SpecSlug:      "0001-widget-flow",
		Agent:         "codex",
		OwnerPID:      pid,
		OwnerIdentity: "identity-token-of-exited-owner-process",
	}
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	active, err := runStore.CreateRun(context.Background(), request)
	if err != nil {
		t.Fatalf("create active Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("force stop exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("failed force stop printed success output %q", stdout.String())
	}
	for _, want := range []string{
		active.ID,
		strconv.Itoa(pid),
		"prove owner process identity",
		"remains Active",
		"Active Run lock retained",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("force stop diagnostic missing %q: %q", want, stderr.String())
		}
	}
	if !store.ProcessAlive(pid) {
		t.Fatalf("refusal must not signal the live process holding reused PID %d", pid)
	}
	assertRunState(t, homeDir, active.ID, store.StateActive)
	runStore, err = store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store after reused PID refusal: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store after reused PID refusal: %v", err)
		}
	}()
	if _, err := runStore.CreateRun(context.Background(), request); err == nil {
		t.Fatal("reused PID refusal released the Active Run lock")
	}
}

// TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner covers Run
// rows created before owner identity recording existed: an absent stored
// token keeps the legacy PID-only proof instead of bricking the manual
// escape hatch.
func TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	pid, ownerWait := startCLIForceStopOwnerProcess(t)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/force-stop-legacy-owner",
		HeadSHA:     "abc123",
		SpecSlug:    "0001-widget-flow",
		Agent:       "codex",
		OwnerPID:    pid,
	})
	if err != nil {
		t.Fatalf("create active Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("legacy force stop exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	select {
	case err := <-ownerWait:
		if err != nil {
			t.Fatalf("owner process exit: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("owner process %d did not exit", pid)
	}
	if store.ProcessAlive(pid) {
		t.Fatalf("owner process %d remained alive after legacy force stop", pid)
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

func TestCLIForceStopOwnerProcessHelper(t *testing.T) {
	if os.Getenv("ROUNDFIX_CLI_FORCE_STOP_OWNER_HELPER") == "" {
		return
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	defer signal.Stop(signals)
	fmt.Fprintln(os.Stdout, "ready")
	<-signals
}

func startCLIForceStopOwnerProcess(t *testing.T) (int, <-chan error) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestCLIForceStopOwnerProcessHelper$")
	cmd.Env = append(os.Environ(), "ROUNDFIX_CLI_FORCE_STOP_OWNER_HELPER=1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open owner process stdout: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start owner process: %v", err)
	}
	wait := make(chan error, 1)
	go func() {
		wait <- cmd.Wait()
		close(wait)
	}()
	t.Cleanup(func() {
		if store.ProcessAlive(cmd.Process.Pid) {
			_ = cmd.Process.Kill()
		}
		select {
		case <-wait:
		case <-time.After(2 * time.Second):
		}
	})
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatalf("owner process did not become ready: %v", scanner.Err())
	}
	if scanner.Text() != "ready" {
		t.Fatalf("owner process readiness = %q, want ready", scanner.Text())
	}
	return cmd.Process.Pid, wait
}

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

	code := RunContext(context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

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

func TestReviewFetchBlocksOlderLiveRunAfterReclaimingNewerOrphan(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	live := seedReviewActiveRun(t, homeDir, repoDir, store.KindWatch, os.Getpid())
	pid := reapedCLIProcessPID(t)
	orphan := seedReviewActiveRunSkippingLock(t, homeDir, repoDir, store.KindResolve, pid)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected fetch to block on older live Run after reclaim, got %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	assertOrphanWarning(t, stderr.String(), orphan.ID, pid)
	if !strings.Contains(stderr.String(), "Active Run "+live.ID) {
		t.Fatalf("expected Branch Integrity Preflight to block on live Run %s, got %q", live.ID, stderr.String())
	}
	assertReclaimedRunInCLI(t, homeDir, orphan.ID, pid)
	assertRunState(t, homeDir, live.ID, store.StateActive)
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

func seedReviewActiveRunSkippingLock(t *testing.T, homeDir string, repoDir string, kind string, pid int) store.Run {
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
	run, err := runStore.CreateRunSkippingActiveLock(context.Background(), store.CreateRunRequest{
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
		t.Fatalf("seed bypassed review active Run: %v", err)
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
