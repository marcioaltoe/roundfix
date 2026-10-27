//go:build unix

// Suite: Implement Command Unix FIFOs and process groups.
// Invariant: Verification Capacity, cancellation and detached ownership preserve their public outcomes.
// Boundary IN: temporary Spec repositories, FIFO releases and owned child process groups.
// Boundary OUT: Windows process control and real runtime providers.
package cli

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func newImplementNamedPipe(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatalf("create named pipe %s: %v", path, err)
	}
	return path
}

func releaseImplementNamedPipe(t *testing.T, path string) {
	t.Helper()
	pipe, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open named pipe %s: %v", path, err)
	}
	if _, err := io.WriteString(pipe, "release\n"); err != nil {
		_ = pipe.Close()
		t.Fatalf("release named pipe %s: %v", path, err)
	}
	if err := pipe.Close(); err != nil {
		t.Fatalf("close named pipe %s: %v", path, err)
	}
}

func blockingVerificationCommand(taskID string, logPath string, startedPath string, releasePath string) string {
	return fmt.Sprintf(
		"%s || exit 1; printf 'start %s\\n' >> %s; printf 'started\\n' > %s; IFS= read -r _ < %s; printf 'end %s\\n' >> %s",
		implementFixtureVerificationCommand(taskID),
		taskID,
		shellQuoteImplement(logPath),
		shellQuoteImplement(startedPath),
		shellQuoteImplement(releasePath),
		taskID,
		shellQuoteImplement(logPath),
	)
}

func TestRunImplementDetachSurvivesCallerProcessGroupKill(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	// Use a passing repository check in this temporary repository.
	writeUserConfig(t, homeDir, "defaults:\n  verification: true\n")
	ownerPID := 0
	ownerIdentity := ""
	reapOwner := func() {
		if ownerPID > 0 {
			killDetachFixtureGroup(t, ownerPID, ownerIdentity)
			waitForDetachFixtureExit(t, ownerPID, ownerIdentity)
			ownerPID = 0
		}
	}
	t.Cleanup(reapOwner)
	dir := t.TempDir()
	t.Cleanup(reapOwner)
	promptStarted := filepath.Join(dir, "prompt-started")
	releasePrompt := filepath.Join(dir, "release")
	t.Cleanup(func() {
		_ = os.WriteFile(releasePrompt, []byte("release\n"), 0o644)
	})
	fakeACPX := fakeACPXCommand(t)
	t.Cleanup(reapOwner)
	cmd := exec.Command(os.Args[0], "implement", "--spec", implementTestSlug, "--detach")
	cmd.Dir = repoDir
	cmd.Env = cliHelperEnv(t, fakeACPX, map[string]string{
		"ROUNDFIX_FAKE_ACPX_PROMPT_STARTED": promptStarted,
		"ROUNDFIX_FAKE_ACPX_RELEASE":        releasePrompt,
		cliTestHoldAfterRunEnv:              "1",
	})
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open caller stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start detach caller: %v", err)
	}

	callerIdentity, err := store.OwnerProcessIdentity(t.Context(), cmd.Process.Pid)
	if err != nil {
		t.Fatalf("record detach caller identity: %v", err)
	}
	callerReaped := false
	t.Cleanup(func() {
		if !callerReaped {
			killDetachFixtureGroup(t, cmd.Process.Pid, callerIdentity)
			_ = waitProcessForTest(t, cmd)
		}
	})
	firstLine := readLineWithTimeout(t, bufio.NewReader(stdoutPipe))
	runID, ok := strings.CutPrefix(strings.TrimSpace(firstLine), "Run ID: ")
	if !ok || strings.TrimSpace(runID) == "" {
		t.Fatalf("expected first detach line with Run id, got %q stderr=%q", firstLine, stderr.String())
	}
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open detached Run reader: %v", err)
	}
	ownedRun, found, runErr := reader.Run(context.Background(), runID)
	closeErr := reader.Close()
	if runErr != nil || closeErr != nil || !found || ownedRun.OwnerPID == nil || *ownedRun.OwnerPID <= 0 {
		t.Fatalf("read detached Run owner: run=%+v found=%v error=%v", ownedRun, found, errors.Join(runErr, closeErr))
	}
	ownerPID = *ownedRun.OwnerPID
	ownerIdentity, err = store.OwnerProcessIdentity(t.Context(), ownerPID)
	if err != nil {
		t.Fatalf("record detached Run owner %d start identity: %v", ownerPID, err)
	}
	if _, err := store.SignalProvenProcessGroup(t.Context(), cmd.Process.Pid, callerIdentity, syscall.SIGKILL); err != nil {
		t.Fatalf("kill caller process group: %v", err)
	}
	_ = waitProcessForTest(t, cmd)
	callerReaped = true
	waitForFile(t, promptStarted, nil)
	if record := os.Getenv(detachTestOwnerRecordEnv); record != "" {
		mustWrite(t, record, strconv.Itoa(ownerPID))
		waitForDetachDeathTestKill(t, record)
	}

	attachCtx, cancel := context.WithCancel(t.Context())
	defer cancel()
	attachStdout := attachFollowCancelWriter{cancel: cancel}
	var attachStderr bytes.Buffer
	attachCode := runCLIContext(t, attachCtx, []string{"attach", runID}, &attachStdout, &attachStderr)
	if attachCode != exitOK {
		t.Fatalf("expected attach to detach cleanly from active Run, got %d stderr=%q stdout=%q", attachCode, attachStderr.String(), attachStdout.String())
	}
	if !strings.Contains(attachStdout.String(), "Detached; Run "+runID+" keeps going.") {
		t.Fatalf("expected attach to prove active Run is attachable, got stdout=%q stderr=%q", attachStdout.String(), attachStderr.String())
	}

	mustWrite(t, releasePrompt, "release\n")
	run := waitForRunState(t, homeDir, runID, store.StateClean, nil)
	if run.State != store.StateClean {
		t.Fatalf("expected detached child to reach Clean after caller kill, got %s", run.State)
	}
	waitForCleanOutcomeEvent(t, homeDir, runID, nil)
}

// Invariant: Task Capacity 2 overlaps Agent turns while Verification Capacity
// 1 serializes real shell gates, and Agent-authored terminal status never
// bypasses Daemon settlement.
// Owning layer: public Implement Command integration.
// Existing canonical suite: TestRunImplementExecutesSpecEndToEnd.
func TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow(t *testing.T) {
	t.Parallel()
	verificationDir := t.TempDir()
	logPath := filepath.Join(verificationDir, "verification.log")
	startedPaths := map[string]string{
		"task_01": filepath.Join(verificationDir, "task_01.started"),
		"task_02": filepath.Join(verificationDir, "task_02.started"),
	}
	releasePaths := map[string]string{
		"task_01": newImplementNamedPipe(t, "task_01.release"),
		"task_02": newImplementNamedPipe(t, "task_02.release"),
	}
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Build the first capacity slice",
			verification: []string{blockingVerificationCommand("task_01", logPath, startedPaths["task_01"], releasePaths["task_01"])},
		},
		{
			id:           "task_02",
			title:        "Build the second capacity slice",
			verification: []string{blockingVerificationCommand("task_02", logPath, startedPaths["task_02"], releasePaths["task_02"])},
		},
	})
	// Disable settlement verification to preserve the call or commit sequence this test measures.
	writeUserConfig(t, homeDir, "verification:\n  repository_at_settlement: false\n")
	configureImplementCapacities(t, repoDir, 2, 1)
	overlap := newImplementAgentOverlapProbe()
	t.Cleanup(overlap.releaseAgents)
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusFailed,
		},
		onTask: overlap.onTask,
	}
	withAgentRunner(t, runner)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	resultCh := runImplementCommandAsync(t, ctx, "implement", "--spec", implementTestSlug, "--no-input")
	assertImplementTaskSet(t, waitImplementAgentStarts(t, overlap, 2, resultCh), "task_01", "task_02")
	if got := overlap.maxObservedActive(); got != 2 {
		t.Fatalf("expected two simultaneous Agent turns, got max active %d", got)
	}
	overlap.releaseAgents()

	runID := onlyImplementRunID(t, homeDir)
	events := waitForImplementJournal(t, homeDir, runID, resultCh, func(events []store.JournalEvent) bool {
		waiting := 0
		started := 0
		for _, event := range implementVerificationEvidenceFromEvents(t, events) {
			if event.attempt != 1 || event.retry != 0 {
				continue
			}
			switch event.phase {
			case string(runevent.VerificationPhaseWaiting):
				waiting++
			case string(runevent.VerificationPhaseStarted):
				started++
			}
		}
		return waiting == 2 && started == 1
	})
	firstEvidence := implementVerificationEvidenceFromEvents(t, events)
	firstTask := ""
	for _, event := range firstEvidence {
		if event.phase == string(runevent.VerificationPhaseStarted) {
			firstTask = event.task
			break
		}
	}
	if firstTask == "" {
		t.Fatal("expected one started Verification Task")
	}
	secondTask := "task_01"
	if firstTask == secondTask {
		secondTask = "task_02"
	}
	waitForFile(t, startedPaths[firstTask], resultCh)
	if _, err := os.Stat(startedPaths[secondTask]); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected %s to remain queued before capacity release, stat error %v", secondTask, err)
	}
	releaseImplementNamedPipe(t, releasePaths[firstTask])

	events = waitForImplementJournal(t, homeDir, runID, resultCh, func(events []store.JournalEvent) bool {
		started := 0
		for _, event := range implementVerificationEvidenceFromEvents(t, events) {
			if event.phase == string(runevent.VerificationPhaseStarted) {
				started++
			}
		}
		return started == 2
	})
	waitForFile(t, startedPaths[secondTask], resultCh)
	releaseImplementNamedPipe(t, releasePaths[secondTask])

	result := waitImplementCommandResult(t, resultCh)
	if result.code != exitOK {
		t.Fatalf("expected Clean exit %d, got %d stderr=%q stdout=%q", exitOK, result.code, result.stderr, result.stdout)
	}
	wantStdout := "task_01 completed — Build the first capacity slice\n" +
		"task_02 completed — Build the second capacity slice\n" +
		"Clean: all 2 Task(s) completed.\n" +
		"Tokens: none reported by 2 prompt(s); cost not reported\n"
	if result.stdout != wantStdout {
		t.Fatalf("normal Implement stdout changed\nwant: %q\n got: %q", wantStdout, result.stdout)
	}
	if !strings.Contains(result.stderr, "Implement Run: "+runID) || !strings.Contains(result.stderr, "reached Clean") {
		t.Fatalf("expected Run identity and terminal progress on stderr, got %q", result.stderr)
	}
	for _, progress := range []string{"Waiting for Verification", "retry", "diagnostics"} {
		if strings.Contains(result.stdout, progress) {
			t.Fatalf("expected %q progress outside stdout, got %q", progress, result.stdout)
		}
	}

	finalEvents := runEventsForRun(t, homeDir, runID)
	verification := implementVerificationEvidenceFromEvents(t, finalEvents)
	active := 0
	maxActive := 0
	waitingCursor := map[string]int64{}
	startedCursor := map[string]int64{}
	for _, event := range verification {
		switch event.phase {
		case string(runevent.VerificationPhaseWaiting):
			waitingCursor[event.task] = event.cursor
			if event.mode != "shared" {
				t.Fatalf("expected shared waiting mode for %s, got %+v", event.task, event)
			}
		case string(runevent.VerificationPhaseStarted):
			startedCursor[event.task] = event.cursor
			active++
			if active > maxActive {
				maxActive = active
			}
		case string(runevent.VerificationPhaseCommandPassed), string(runevent.VerificationPhaseFailed):
			active--
		}
	}
	if active != 0 || maxActive != 1 {
		t.Fatalf("expected real gates to drain with max active 1, got active=%d max=%d events=%+v", active, maxActive, verification)
	}
	for _, taskID := range []string{"task_01", "task_02"} {
		if waitingCursor[taskID] == 0 || startedCursor[taskID] == 0 || waitingCursor[taskID] >= startedCursor[taskID] {
			t.Fatalf("expected waiting before started for %s, waiting=%d started=%d", taskID, waitingCursor[taskID], startedCursor[taskID])
		}
		content := mustRead(t, implementTaskPath(repoDir, taskID))
		if !strings.Contains(content, "status: completed") {
			t.Fatalf("expected only Daemon settlement to leave %s completed, got:\n%s", taskID, content)
		}
	}
	if got := mustRead(t, logPath); got != "start "+firstTask+"\nend "+firstTask+"\nstart "+secondTask+"\nend "+secondTask+"\n" {
		t.Fatalf("expected serialized real-shell lifecycle, got %q", got)
	}
	if runner.calls != 2 {
		t.Fatalf("expected one Agent turn per Task with no repair, got %d", runner.calls)
	}
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Run Clean, got %s", run.State)
	}
	assertRunWorktreeRemoved(t, run.WorkDir)
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}

// Invariant: cancellation while a Task waits for Verification Capacity starts
// no child for that Task, returns every worker, and preserves resumable
// in_progress Task Worktrees without a terminal settlement.
// Owning layer: public Implement Command integration.
// Existing canonical suite: TestRunImplementStopRequestEndsStoppedWithInterruptMapping.
func TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks(t *testing.T) {
	t.Parallel()
	verificationDir := t.TempDir()
	logPath := filepath.Join(verificationDir, "verification.log")
	startedPaths := map[string]string{
		"task_01": filepath.Join(verificationDir, "task_01.started"),
		"task_02": filepath.Join(verificationDir, "task_02.started"),
	}
	releasePaths := map[string]string{
		"task_01": newImplementNamedPipe(t, "task_01.release"),
		"task_02": newImplementNamedPipe(t, "task_02.release"),
	}
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{
			id:           "task_01",
			title:        "Queue the first cancellation gate",
			verification: []string{blockingVerificationCommand("task_01", logPath, startedPaths["task_01"], releasePaths["task_01"])},
		},
		{
			id:           "task_02",
			title:        "Queue the second cancellation gate",
			verification: []string{blockingVerificationCommand("task_02", logPath, startedPaths["task_02"], releasePaths["task_02"])},
		},
	})
	configureImplementCapacities(t, repoDir, 2, 1)
	overlap := newImplementAgentOverlapProbe()
	t.Cleanup(overlap.releaseAgents)
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusFailed,
		},
		onTask: overlap.onTask,
	}
	withAgentRunner(t, runner)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	resultCh := runImplementCommandAsync(t, ctx, "implement", "--spec", implementTestSlug, "--no-input")
	assertImplementTaskSet(t, waitImplementAgentStarts(t, overlap, 2, resultCh), "task_01", "task_02")
	overlap.releaseAgents()
	runID := onlyImplementRunID(t, homeDir)
	events := waitForImplementJournal(t, homeDir, runID, resultCh, func(events []store.JournalEvent) bool {
		waiting := 0
		started := 0
		for _, event := range implementVerificationEvidenceFromEvents(t, events) {
			switch event.phase {
			case string(runevent.VerificationPhaseWaiting):
				waiting++
			case string(runevent.VerificationPhaseStarted):
				started++
			}
		}
		return waiting == 2 && started == 1
	})
	activeTask := ""
	for _, event := range implementVerificationEvidenceFromEvents(t, events) {
		if event.phase == string(runevent.VerificationPhaseStarted) {
			activeTask = event.task
			break
		}
	}
	if activeTask == "" {
		t.Fatal("expected one active Verification before cancellation")
	}
	queuedTask := "task_01"
	if activeTask == queuedTask {
		queuedTask = "task_02"
	}
	waitForFile(t, startedPaths[activeTask], resultCh)
	if _, err := os.Stat(startedPaths[queuedTask]); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected queued Task %s to start no child process, stat error %v", queuedTask, err)
	}

	cancel()
	result := waitImplementCommandResult(t, resultCh)
	if result.code != exitOK {
		t.Fatalf("expected Stopped command exit %d before interrupt mapping, got %d stderr=%q stdout=%q", exitOK, result.code, result.stderr, result.stdout)
	}
	wantStdout := "task_01 pending — Queue the first cancellation gate\n" +
		"task_02 pending — Queue the second cancellation gate\n" +
		"Stopped: 0 completed, 0 failed, 0 skipped, 2 pending.\n" +
		"Tokens: none reported by 2 prompt(s); cost not reported\n"
	if result.stdout != wantStdout {
		t.Fatalf("queued cancellation stdout changed\nwant: %q\n got: %q", wantStdout, result.stdout)
	}
	if !strings.Contains(result.stderr, "reached Stopped") || strings.Contains(result.stdout, "Waiting for Verification") {
		t.Fatalf("expected stop progress on stderr and stable stdout, stderr=%q stdout=%q", result.stderr, result.stdout)
	}
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateStopped {
		t.Fatalf("expected Run Stopped, got %s", run.State)
	}
	for _, taskID := range []string{activeTask, queuedTask} {
		taskRef := taskRefForImplementRun(t, run, repoDir, taskID)
		assertRunWorktreeExists(t, taskRef.Path)
		content := mustRead(t, filepath.Join(taskRef.Path, "docs", "specs", implementTestSlug, taskID+".md"))
		if !strings.Contains(content, "status: in_progress") {
			t.Fatalf("expected resumable in_progress status for %s, got:\n%s", taskID, content)
		}
	}
	if _, err := os.Stat(startedPaths[queuedTask]); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("queued Task %s started a child after cancellation, stat error %v", queuedTask, err)
	}
	finalEvents := runEventsForRun(t, homeDir, runID)
	for _, entry := range finalEvents {
		if entry.Event.Kind != runevent.KindDaemonTask || entry.Event.ReviewIssue != queuedTask {
			continue
		}
		if strings.Contains(string(entry.Event.Payload), `"phase":"settled"`) {
			t.Fatalf("expected no terminal settlement for queued Task %s, got %+v", queuedTask, entry.Event)
		}
	}
	if got := mustRead(t, logPath); got != "start "+activeTask+"\n" {
		t.Fatalf("expected only the active child to start before cancellation, got %q", got)
	}
	assertRunWorktreeExists(t, run.WorkDir)
	assertNoActiveRunInGitRoot(t, homeDir, repoDir)
}
