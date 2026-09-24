// Suite: Delivery Queue command family.
// Invariant: only validated Spec slugs are persisted, and the persisted queue remains observable and controllable across owner processes.
// Boundary IN: public CLI dispatch, the real SQLite Delivery Queue store, and Spec loading.
// Boundary OUT: detached process launch and process signals, injected as operating-system boundaries.
package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/store"
)

func TestDeliverCommandRecordsAQueueAndReportsIt(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(
			context.Context,
			roundconfig.Loaded,
			commandEnvironment,
			io.Writer,
			io.Writer,
		) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver start exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if started != 1 {
		t.Fatalf("detached owner starts = %d, want 1", started)
	}
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	queue, found, err := runStore.DeliveryQueue(context.Background(), repoDir)
	if err != nil || !found || len(queue.Items) != 1 {
		t.Fatalf("recorded Delivery Queue: found=%v items=%d err=%v", found, len(queue.Items), err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = "review-stale"
	if err := runStore.UpdateDeliveryQueueItem(context.Background(), repoDir, item); err != nil {
		t.Fatalf("park Delivery Queue item: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	stdout.Reset()
	stderr.Reset()

	code = runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver status exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if got, want := stdout.String(), implementTestSlug+"\tparked\treview-stale\n"; got != want {
		t.Fatalf("deliver status = %q, want %q", got, want)
	}
}

func TestDeliverCommandRejectsUnknownSlugBeforeRecordingQueue(t *testing.T) {
	t.Parallel()
	homeDir, _ := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	started := false
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started = true
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "start", "9999-unknown-spec"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unknown slug exit = %d, want %d stderr=%q", code, exitPreflight, stderr.String())
	}
	if started {
		t.Fatal("unknown slug started the Delivery Queue owner")
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); !os.IsNotExist(err) {
		t.Fatalf("unknown slug created Run Database: err=%v", err)
	}
	if !strings.Contains(stderr.String(), "9999-unknown-spec") {
		t.Fatalf("unknown slug diagnostic = %q", stderr.String())
	}
}

func TestDeliverCommandStopsAndResumesPersistedQueue(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	if _, err := runStore.CreateDeliveryQueue(ctx, repoDir, []string{implementTestSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	if err := runStore.ClaimDeliveryQueueOwner(ctx, repoDir, 4242, "owner-identity"); err != nil {
		t.Fatalf("claim Delivery Queue owner: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	processes := &fakeDeliveryOwnerProcesses{}
	resumed := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.ownerProcesses = processes
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			resumed++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "stop"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver stop exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if processes.provedPID != 4242 || processes.terminatedPID != 4242 || processes.identity != "owner-identity" {
		t.Fatalf("owner process calls = %+v", processes)
	}
	stdout.Reset()
	stderr.Reset()

	code = runCLI(t, []string{"deliver", "resume"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver resume exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if resumed != 1 {
		t.Fatalf("resume owner starts = %d, want 1", resumed)
	}
	persisted := openDeliveryQueueForCLI(t, homeDir, repoDir)
	if persisted.OwnerPID != 0 || persisted.OwnerIdentity != "" || len(persisted.Items) != 1 {
		t.Fatalf("persisted queue after stop/resume = %+v", persisted)
	}
}

func TestATerminalQueueIsReplacedByANewStart(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	queue, err := runStore.CreateDeliveryQueue(ctx, repoDir, []string{"old-spec"})
	if err != nil {
		t.Fatalf("create old Delivery Queue: %v", err)
	}
	old := queue.Items[0]
	old.Stage = store.DeliveryStageParked
	old.Blocker = "old blocker"
	if err := runStore.UpdateDeliveryQueueItem(ctx, repoDir, old); err != nil {
		t.Fatalf("park old Delivery Queue item: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("replacement deliver start exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if started != 1 {
		t.Fatalf("detached owner starts = %d, want 1", started)
	}
	replaced := openDeliveryQueueForCLI(t, homeDir, repoDir)
	if len(replaced.Items) != 1 || replaced.Items[0].SpecSlug != implementTestSlug || replaced.Items[0].Stage != store.DeliveryStageQueued {
		t.Fatalf("replacement Delivery Queue = %+v", replaced)
	}
}

func TestDeliveryWorkflowCreatesAnItemBranchFromTheRefreshedDefault(t *testing.T) {
	t.Parallel()
	origin := t.TempDir()
	gittest.InitRepo(t, origin, "--initial-branch=main")
	gittest.PersistIdentity(t, origin)
	seedPath := filepath.Join(origin, "seed.txt")
	if err := os.WriteFile(seedPath, []byte("seed\n"), 0o644); err != nil {
		t.Fatalf("write origin seed: %v", err)
	}
	gittest.Run(t, origin, "add", "seed.txt")
	gittest.Run(t, origin, "commit", "-m", "chore: seed")

	checkout := filepath.Join(t.TempDir(), "checkout")
	gittest.Run(t, "", "clone", origin, checkout)
	gittest.Harden(t, checkout)
	if err := os.WriteFile(seedPath, []byte("refreshed\n"), 0o644); err != nil {
		t.Fatalf("refresh origin seed: %v", err)
	}
	gittest.Run(t, origin, "add", "seed.txt")
	gittest.Run(t, origin, "commit", "-m", "fix: refresh default")
	wantHead := strings.TrimSpace(gittest.Run(t, origin, "rev-parse", "main"))
	workflow := &commandDeliveryWorkflow{
		loaded: roundconfig.Loaded{GitRoot: checkout},
		git:    preflight.ExecGitRunner{},
	}

	branch, err := workflow.CreateItemBranch(t.Context(), checkout, "0156-delivery")

	if err != nil {
		t.Fatalf("create item branch: %v", err)
	}
	if branch != "roundfix/deliver-0156-delivery" {
		t.Fatalf("item branch = %q", branch)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "branch", "--show-current")); got != branch {
		t.Fatalf("current branch = %q, want %q", got, branch)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("item branch head = %q, want refreshed default %q", got, wantHead)
	}
}

func TestDeliverCommandRefusesUnknownFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args []string
	}{
		{name: "start", args: []string{"deliver", "start", implementTestSlug, "--unknown"}},
		{name: "status", args: []string{"deliver", "status", "--unknown"}},
		{name: "resume", args: []string{"deliver", "resume", "--unknown"}},
		{name: "stop", args: []string{"deliver", "stop", "--unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLI(t, tt.args, &stdout, &stderr)

			if code != exitPreflight || stdout.Len() != 0 {
				t.Fatalf("unknown flag exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stderr.String(), "flag provided but not defined") {
				t.Fatalf("unknown flag diagnostic = %q", stderr.String())
			}
		})
	}
}

type fakeDeliveryOwnerProcesses struct {
	provedPID     int
	terminatedPID int
	identity      string
}

func (fake *fakeDeliveryOwnerProcesses) ProveOwner(_ context.Context, pid int, identity string) error {
	fake.provedPID = pid
	fake.identity = identity
	return nil
}

func (fake *fakeDeliveryOwnerProcesses) TerminateTreeAndWait(_ context.Context, pid int, identity string) ([]store.TerminationOutcome, error) {
	fake.terminatedPID = pid
	fake.identity = identity
	return []store.TerminationOutcome{{PID: pid, Proven: true}}, nil
}

func openDeliveryQueueForCLI(t *testing.T, homeDir, gitRoot string) store.DeliveryQueue {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close Run Database: %v", err)
		}
	}()
	queue, found, err := runStore.DeliveryQueue(context.Background(), gitRoot)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue: found=%v err=%v", found, err)
	}
	return queue
}
