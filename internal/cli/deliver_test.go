// Suite: Delivery Queue command family.
// Invariant: only validated Spec slugs are persisted, and the persisted queue remains observable and controllable across owner processes.
// Boundary IN: public CLI dispatch, the real SQLite Delivery Queue store, and Spec loading.
// Boundary OUT: detached process launch and process signals, injected as operating-system boundaries.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestDeliverStatusPrintsTheItemWorktree(t *testing.T) {
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
	stdout.Reset()
	stderr.Reset()
	code = runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver status without worktree exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if got, want := stdout.String(), implementTestSlug+"\tqueued\t-\t-\n"; got != want {
		t.Fatalf("deliver status without worktree = %q, want %q", got, want)
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
	if _, _, err := runStore.RecordDeliveryQueueItemWorktree(
		context.Background(),
		repoDir,
		item.SpecSlug,
		"roundfix/deliver-status",
		"/worktrees/delivery-item",
	); err != nil {
		t.Fatalf("record Delivery Queue item worktree: %v", err)
	}
	item.Stage = store.DeliveryStageParked
	item.Blocker = "review-stale"
	item.Branch = "roundfix/deliver-status"
	item.Worktree = "/worktrees/delivery-item"
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
	if got, want := stdout.String(), implementTestSlug+"\tparked\treview-stale\t/worktrees/delivery-item\n"; got != want {
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

func TestResumeReleasesAStaleOwner(t *testing.T) {
	tests := []struct {
		name           string
		proofErr       error
		wantExit       int
		wantResumed    int
		wantOwnerPID   int
		wantDiagnostic string
	}{
		{
			name:           "reused PID",
			proofErr:       fmt.Errorf("identity mismatch: %w", store.ErrOwnerProcessIdentityUnproven),
			wantExit:       exitOK,
			wantResumed:    1,
			wantDiagnostic: "different process identity",
		},
		{
			name:           "unreadable identity",
			proofErr:       fmt.Errorf("identity read failed: %w", store.ErrOwnerIdentityUnreadable),
			wantExit:       exitPreflight,
			wantOwnerPID:   os.Getpid(),
			wantDiagnostic: "identity read failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			ctx := context.Background()
			runStore, err := store.Open(ctx, homeDir)
			if err != nil {
				t.Fatalf("open Run Database: %v", err)
			}
			if _, err := runStore.CreateDeliveryQueue(ctx, repoDir, []string{implementTestSlug}); err != nil {
				t.Fatalf("create Delivery Queue: %v", err)
			}
			if err := runStore.ClaimDeliveryQueueOwner(ctx, repoDir, os.Getpid(), "recorded-owner-identity"); err != nil {
				t.Fatalf("claim stale Delivery Queue owner: %v", err)
			}
			if err := runStore.Close(); err != nil {
				t.Fatalf("close Run Database: %v", err)
			}

			processes := &fakeDeliveryOwnerProcesses{proveErr: tt.proofErr}
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

			code := runCLI(t, []string{"deliver", "resume"}, &stdout, &stderr)

			if code != tt.wantExit || resumed != tt.wantResumed {
				t.Fatalf("deliver resume exit=%d resumed=%d stderr=%q, want exit=%d resumed=%d", code, resumed, stderr.String(), tt.wantExit, tt.wantResumed)
			}
			if processes.provedPID != os.Getpid() || processes.identity != "recorded-owner-identity" {
				t.Fatalf("owner identity proof = %+v", processes)
			}
			if !strings.Contains(stderr.String(), tt.wantDiagnostic) {
				t.Fatalf("deliver resume diagnostic = %q, want %q", stderr.String(), tt.wantDiagnostic)
			}
			persisted := openDeliveryQueueForCLI(t, homeDir, repoDir)
			if persisted.OwnerPID != tt.wantOwnerPID {
				t.Fatalf("persisted owner PID = %d, want %d", persisted.OwnerPID, tt.wantOwnerPID)
			}
		})
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

func TestEachItemRunsInItsOwnWorktree(t *testing.T) {
	t.Parallel()
	origin, checkout := newDeliveryBranchRepository(t)
	seedPath := filepath.Join(origin, "seed.txt")
	if err := os.WriteFile(seedPath, []byte("refreshed\n"), 0o644); err != nil {
		t.Fatalf("refresh origin seed: %v", err)
	}
	gittest.Run(t, origin, "add", "seed.txt")
	gittest.Run(t, origin, "commit", "-m", "fix: refresh default")
	wantHead := strings.TrimSpace(gittest.Run(t, origin, "rev-parse", "main"))
	const specSlug = "0168-item-worktree"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)
	copyPath := filepath.Join(checkout, "delivery.env")
	mustWrite(t, copyPath, "copied from user checkout\n")
	workflow.loaded.Config.Worktree.Copy = []string{"delivery.env"}
	workflow.loaded.Config.Worktree.Bootstrap = "printf 'bootstrapped\\n' > delivery-bootstrap.txt"
	workflow.loaded.Config.Worktree.BootstrapTimeout = time.Minute

	branch, itemWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)

	if err != nil {
		t.Fatalf("create item worktree: %v", err)
	}
	if !strings.HasPrefix(branch, "roundfix/deliver-"+specSlug+"-") {
		t.Fatalf("item branch = %q, want per-delivery suffix", branch)
	}
	if !strings.HasPrefix(itemWorktree, workflow.loaded.Config.Worktree.Location+string(filepath.Separator)) {
		t.Fatalf("item worktree = %q, want under %q", itemWorktree, workflow.loaded.Config.Worktree.Location)
	}
	if got := strings.TrimSpace(gittest.Run(t, itemWorktree, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("item branch head = %q, want refreshed default %q", got, wantHead)
	}
	upstream := strings.TrimSpace(gittest.Run(
		t,
		itemWorktree,
		"for-each-ref",
		"--format=%(upstream:short)",
		"refs/heads/"+branch,
	))
	if upstream != "" {
		t.Fatalf("item branch upstream = %q, want none", upstream)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "branch", "--show-current")); got != "main" {
		t.Fatalf("user checkout branch = %q, want unchanged main", got)
	}
	if got := mustRead(t, filepath.Join(itemWorktree, "delivery.env")); got != "copied from user checkout\n" {
		t.Fatalf("provisioned copy = %q", got)
	}
	if got := mustRead(t, filepath.Join(itemWorktree, "delivery-bootstrap.txt")); got != "bootstrapped\n" {
		t.Fatalf("bootstrap output = %q", got)
	}
	queue, found, err := workflow.store.DeliveryQueue(t.Context(), checkout)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue: found=%v err=%v", found, err)
	}
	if queue.Items[0].Branch != branch || queue.Items[0].Worktree != itemWorktree {
		t.Fatalf("recorded item workspace = branch %q worktree %q, want %q and %q", queue.Items[0].Branch, queue.Items[0].Worktree, branch, itemWorktree)
	}
}

func TestItemBranchHasNoUpstream(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-no-upstream"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)

	branch, itemWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("create item worktree: %v", err)
	}
	upstream := strings.TrimSpace(gittest.Run(
		t,
		itemWorktree,
		"for-each-ref",
		"--format=%(upstream:short)",
		"refs/heads/"+branch,
	))
	if upstream != "" {
		t.Fatalf("item branch upstream = %q, want none", upstream)
	}
}

func TestEachDeliveryGetsItsOwnBranch(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-repeat"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)

	firstBranch, firstWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("create first item worktree: %v", err)
	}
	queue, found, err := workflow.store.DeliveryQueue(t.Context(), checkout)
	if err != nil || !found {
		t.Fatalf("read first Delivery Queue: found=%v err=%v", found, err)
	}
	firstItem := queue.Items[0]
	firstItem.Stage = store.DeliveryStageParked
	if err := workflow.store.UpdateDeliveryQueueItem(t.Context(), checkout, firstItem); err != nil {
		t.Fatalf("park first delivery: %v", err)
	}
	if _, err := workflow.store.CreateDeliveryQueue(t.Context(), checkout, []string{specSlug}); err != nil {
		t.Fatalf("create second Delivery Queue: %v", err)
	}

	secondBranch, secondWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("create second item worktree: %v", err)
	}
	if firstBranch == secondBranch {
		t.Fatalf("delivery branches = %q and %q, want different branches", firstBranch, secondBranch)
	}
	if firstWorktree == secondWorktree {
		t.Fatalf("delivery worktrees = %q and %q, want different worktrees", firstWorktree, secondWorktree)
	}
}

func TestResumeReusesTheRecordedItemBranch(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-resume-branch"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)

	firstBranch, firstWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("create item worktree before crash: %v", err)
	}
	queue, found, err := workflow.store.DeliveryQueue(t.Context(), checkout)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue after worktree creation: found=%v err=%v", found, err)
	}
	if got := queue.Items[0].Branch; got != firstBranch {
		t.Fatalf("recorded item branch = %q, want %q", got, firstBranch)
	}
	if queue.Items[0].Stage != store.DeliveryStageQueued {
		t.Fatalf("item stage after simulated crash = %q, want queued", queue.Items[0].Stage)
	}
	gittest.Run(t, checkout, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing-origin"))

	resumedBranch, resumedWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("resume item worktree: %v", err)
	}
	if resumedBranch != firstBranch {
		t.Fatalf("resumed item branch = %q, want recorded branch %q", resumedBranch, firstBranch)
	}
	if resumedWorktree != firstWorktree {
		t.Fatalf("resumed item worktree = %q, want recorded worktree %q", resumedWorktree, firstWorktree)
	}
}

func TestDeliverNeverTouchesTheUserCheckout(t *testing.T) {
	origin, checkout := newDeliveryBranchRepository(t)
	gittest.Run(t, checkout, "switch", "-c", "user-work")
	mustWrite(t, filepath.Join(checkout, ".gitignore"), "ignored.txt\n")
	gittest.Run(t, checkout, "add", ".gitignore")
	gittest.Run(t, checkout, "commit", "-m", "chore: ignore user file")
	mustWrite(t, filepath.Join(checkout, "seed.txt"), "dirty user change\n")
	mustWrite(t, filepath.Join(checkout, "ignored.txt"), "ignored user content\n")
	nestedRepository := filepath.Join(checkout, "nested-repository")
	gittest.InitRepo(t, nestedRepository, "--initial-branch=main")
	mustWrite(t, filepath.Join(nestedRepository, "nested.txt"), "nested repository content\n")

	beforeBranch := strings.TrimSpace(gittest.Run(t, checkout, "branch", "--show-current"))
	beforeHead := strings.TrimSpace(gittest.Run(t, checkout, "rev-parse", "HEAD"))
	beforeStatus := gittest.Run(t, checkout, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignored=matching")
	beforeSeed := mustRead(t, filepath.Join(checkout, "seed.txt"))
	beforeIgnored := mustRead(t, filepath.Join(checkout, "ignored.txt"))
	beforeNested := mustRead(t, filepath.Join(nestedRepository, "nested.txt"))

	ctx := t.Context()
	runStore, err := store.Open(ctx, t.TempDir())
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close Run Database: %v", err)
		}
	})
	const parkedSlug = "0168-parked"
	const mergedSlug = "0168-merged"
	if _, err := runStore.CreateDeliveryQueue(ctx, checkout, []string{parkedSlug, mergedSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := &commandDeliveryWorkflow{
		store: runStore,
		loaded: roundconfig.Loaded{
			GitRoot: checkout,
			Config: roundconfig.Config{
				Worktree: roundconfig.Worktree{Location: filepath.Join(t.TempDir(), "worktrees")},
			},
		},
		git: preflight.ExecGitRunner{},
	}
	reviewedHead := strings.TrimSpace(gittest.Run(t, origin, "rev-parse", "main"))
	flow := &parkTestDeliveryFlow{parkSlug: parkedSlug, reviewedHead: reviewedHead}
	engine := delivery.NewEngine(runStore, delivery.EngineDependencies{
		Workspace:    workflow,
		Runner:       flow,
		Reviewer:     flow,
		Archiver:     flow,
		Gate:         flow,
		Authorizer:   flow,
		Publication:  flow,
		PullRequests: flow,
	})

	if _, err := engine.Run(ctx, checkout); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue, found, err := runStore.DeliveryQueue(ctx, checkout)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue: found=%v err=%v", found, err)
	}
	if queue.Items[0].Stage != store.DeliveryStageParked || queue.Items[1].Stage != store.DeliveryStageMerged {
		t.Fatalf("Delivery Queue stages = %q and %q, want parked and merged", queue.Items[0].Stage, queue.Items[1].Stage)
	}
	if queue.Items[0].Worktree == "" || queue.Items[1].Worktree == "" || queue.Items[0].Worktree == queue.Items[1].Worktree {
		t.Fatalf("item worktrees = %q and %q, want distinct recorded paths", queue.Items[0].Worktree, queue.Items[1].Worktree)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "branch", "--show-current")); got != beforeBranch {
		t.Fatalf("user checkout branch = %q, want unchanged %q", got, beforeBranch)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "rev-parse", "HEAD")); got != beforeHead {
		t.Fatalf("user checkout HEAD = %q, want unchanged %q", got, beforeHead)
	}
	if got := gittest.Run(t, checkout, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignored=matching"); got != beforeStatus {
		t.Fatalf("user checkout status = %q, want unchanged %q", got, beforeStatus)
	}
	if got := mustRead(t, filepath.Join(checkout, "seed.txt")); got != beforeSeed {
		t.Fatalf("tracked user file = %q, want unchanged %q", got, beforeSeed)
	}
	if got := mustRead(t, filepath.Join(checkout, "ignored.txt")); got != beforeIgnored {
		t.Fatalf("ignored user file = %q, want unchanged %q", got, beforeIgnored)
	}
	if got := mustRead(t, filepath.Join(nestedRepository, "nested.txt")); got != beforeNested {
		t.Fatalf("nested repository file = %q, want unchanged %q", got, beforeNested)
	}
}

func TestParkLeavesTheItemWorktreeInPlace(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-park-worktree"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)
	flow := &parkTestDeliveryFlow{parkSlug: specSlug}
	engine := delivery.NewEngine(workflow.store, delivery.EngineDependencies{
		Workspace:    workflow,
		Runner:       flow,
		Reviewer:     flow,
		Archiver:     flow,
		Gate:         flow,
		Authorizer:   flow,
		Publication:  flow,
		PullRequests: flow,
	})

	if _, err := engine.Run(t.Context(), checkout); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	queue, found, err := workflow.store.DeliveryQueue(t.Context(), checkout)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue: found=%v err=%v", found, err)
	}
	item := queue.Items[0]
	if item.Stage != store.DeliveryStageParked || item.Worktree == "" {
		t.Fatalf("parked item = %+v, want parked with recorded worktree", item)
	}
	info, err := os.Stat(item.Worktree)
	if err != nil || !info.IsDir() {
		t.Fatalf("parked item worktree %q: info=%v err=%v", item.Worktree, info, err)
	}
	registered := gittest.Run(t, checkout, "worktree", "list", "--porcelain")
	if !strings.Contains(registered, "branch refs/heads/"+item.Branch+"\n") {
		t.Fatalf("registered worktrees = %q, want parked branch %q", registered, item.Branch)
	}
}

func TestResumeUsesTheRecordedItemWorktree(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-resume-existing-worktree"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)
	branch, itemWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("create item worktree: %v", err)
	}
	seedDeliveryItemStage(t, workflow.store, checkout, store.DeliveryStageRunning)
	flow := &parkTestDeliveryFlow{parkSlug: specSlug}
	engine := newDeliveryLifecycleTestEngine(workflow.store, workflow, flow)

	if _, err := engine.Run(t.Context(), checkout); err != nil {
		t.Fatalf("resume Delivery Engine: %v", err)
	}

	if flow.runCalls != 1 || flow.runWorkDir != itemWorktree {
		t.Fatalf("resumed Implement calls = %d at %q, want one at %q", flow.runCalls, flow.runWorkDir, itemWorktree)
	}
	if got := strings.TrimSpace(gittest.Run(t, itemWorktree, "branch", "--show-current")); got != branch {
		t.Fatalf("recorded worktree branch = %q, want %q", got, branch)
	}
}

func TestResumeRecreatesAMissingItemWorktree(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-resume-missing-worktree"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)
	branch, itemWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("create item worktree: %v", err)
	}
	seedDeliveryItemStage(t, workflow.store, checkout, store.DeliveryStageRunning)
	gittest.Run(t, checkout, "worktree", "remove", itemWorktree)
	flow := &parkTestDeliveryFlow{parkSlug: specSlug}
	engine := newDeliveryLifecycleTestEngine(workflow.store, workflow, flow)

	if _, err := engine.Run(t.Context(), checkout); err != nil {
		t.Fatalf("resume Delivery Engine: %v", err)
	}

	if flow.runCalls != 1 || flow.runWorkDir != itemWorktree {
		t.Fatalf("resumed Implement calls = %d at %q, want one at recreated worktree %q", flow.runCalls, flow.runWorkDir, itemWorktree)
	}
	if got := strings.TrimSpace(gittest.Run(t, itemWorktree, "branch", "--show-current")); got != branch {
		t.Fatalf("recreated worktree branch = %q, want %q", got, branch)
	}
}

func TestResumeParksWhenTheItemBranchIsGone(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-resume-missing-branch"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)
	branch, itemWorktree, err := workflow.CreateItemBranch(t.Context(), checkout, specSlug)
	if err != nil {
		t.Fatalf("create item worktree: %v", err)
	}
	seedDeliveryItemStage(t, workflow.store, checkout, store.DeliveryStageRunning)
	gittest.Run(t, checkout, "worktree", "remove", itemWorktree)
	gittest.Run(t, checkout, "branch", "-D", branch)
	flow := &parkTestDeliveryFlow{parkSlug: specSlug}
	engine := newDeliveryLifecycleTestEngine(workflow.store, workflow, flow)

	if _, err := engine.Run(t.Context(), checkout); err != nil {
		t.Fatalf("resume Delivery Engine: %v", err)
	}
	item := readDeliveryItemForCLI(t, workflow.store, checkout)
	if item.Stage != store.DeliveryStageParked || item.Blocker != delivery.BlockerItemWorktreeMissing {
		t.Fatalf("resumed item = %+v, want parked with %q", item, delivery.BlockerItemWorktreeMissing)
	}
	if flow.runCalls != 0 {
		t.Fatalf("Implement calls after missing branch = %d, want none", flow.runCalls)
	}
	if _, err := engine.Run(t.Context(), checkout); err != nil {
		t.Fatalf("resume parked Delivery Engine: %v", err)
	}
	if flow.runCalls != 0 {
		t.Fatalf("Implement calls after second resume = %d, want none", flow.runCalls)
	}
}

func TestAMergedItemLeavesNoWorktreeOrBranch(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	const specSlug = "0168-merged-cleanup"
	workflow := newDeliveryBranchWorkflow(t, checkout, specSlug)
	mustWrite(t, filepath.Join(checkout, "delivery.env"), "copied into item worktree\n")
	workflow.loaded.Config.Worktree.Copy = []string{"delivery.env"}
	workflow.loaded.Config.Worktree.Bootstrap = "printf 'bootstrap residue\\n' > delivery-bootstrap.txt"
	workflow.loaded.Config.Worktree.BootstrapTimeout = time.Minute
	flow := &parkTestDeliveryFlow{}
	engine := newDeliveryLifecycleTestEngine(workflow.store, workflow, flow)

	if _, err := engine.Run(t.Context(), checkout); err != nil {
		t.Fatalf("run Delivery Engine: %v", err)
	}

	item := readDeliveryItemForCLI(t, workflow.store, checkout)
	if item.Stage != store.DeliveryStageMerged {
		t.Fatalf("merged item stage = %q, want merged", item.Stage)
	}
	if _, err := os.Stat(item.Worktree); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("merged item worktree %q still exists: %v", item.Worktree, err)
	}
	exists, err := localItemBranchExists(t.Context(), preflight.ExecGitRunner{}, checkout, item.Branch)
	if err != nil {
		t.Fatalf("inspect merged item branch: %v", err)
	}
	if exists {
		t.Fatalf("merged item branch %q still exists", item.Branch)
	}
}

func seedDeliveryItemStage(t *testing.T, runStore *store.Store, gitRoot string, stage store.DeliveryStage) {
	t.Helper()
	item := readDeliveryItemForCLI(t, runStore, gitRoot)
	item.Stage = stage
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), gitRoot, item); err != nil {
		t.Fatalf("seed Delivery Queue item stage %q: %v", stage, err)
	}
}

func readDeliveryItemForCLI(t *testing.T, runStore *store.Store, gitRoot string) store.DeliveryQueueItem {
	t.Helper()
	queue, found, err := runStore.DeliveryQueue(t.Context(), gitRoot)
	if err != nil || !found || len(queue.Items) != 1 {
		t.Fatalf("read Delivery Queue item: found=%v items=%d err=%v", found, len(queue.Items), err)
	}
	return queue.Items[0]
}

func newDeliveryLifecycleTestEngine(
	runStore *store.Store,
	workflow *commandDeliveryWorkflow,
	flow *parkTestDeliveryFlow,
) *delivery.Engine {
	return delivery.NewEngine(runStore, delivery.EngineDependencies{
		Workspace:    workflow,
		Runner:       flow,
		Reviewer:     flow,
		Archiver:     flow,
		Gate:         flow,
		Authorizer:   flow,
		Publication:  flow,
		PullRequests: flow,
	})
}

func TestResumeAcceptsARealArchiveCommit(t *testing.T) {
	repository, reviewedHead, archiveHead := commitRealArchive(t, nil)

	item := resumeArchivedDelivery(t, repository, reviewedHead)

	if item.Blocker == delivery.BlockerReviewStale {
		t.Fatalf("resumed archive item = %+v, want real archive commit accepted", item)
	}
	if !reflect.DeepEqual(item.CandidateCommits, []string{reviewedHead, archiveHead}) {
		t.Fatalf("candidate commits = %q, want reviewed and archive heads", item.CandidateCommits)
	}
}

func TestResumeRefusesAnArchiveCommitWithExtraChanges(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, string)
	}{
		{
			name: "unrelated path",
			mutate: func(t *testing.T, repository string) {
				t.Helper()
				mustWrite(t, filepath.Join(repository, "unrelated.txt"), "unrelated\n")
			},
		},
		{
			name: "changed PRD body",
			mutate: func(t *testing.T, repository string) {
				t.Helper()
				prdPath := archiveTestRepositoryPath(repository, spec.ArchiveKindSpec, implementTestSlug, "_prd.md")
				mustWrite(t, prdPath, mustRead(t, prdPath)+"\nChanged after archive.\n")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository, reviewedHead, _ := commitRealArchive(t, tt.mutate)

			item := resumeArchivedDelivery(t, repository, reviewedHead)

			if item.Stage != store.DeliveryStageParked || item.Blocker != delivery.BlockerReviewStale {
				t.Fatalf("resumed archive item = %+v, want parked as review-stale", item)
			}
			if !reflect.DeepEqual(item.CandidateCommits, []string{reviewedHead}) {
				t.Fatalf("candidate commits = %q, want only reviewed head", item.CandidateCommits)
			}
		})
	}
}

func commitRealArchive(t *testing.T, mutate func(*testing.T, string)) (string, string, string) {
	t.Helper()
	_, repository := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
	})
	const unprovenAction = "a maintainer publishes the tagged release"
	appendArchiveUnreachableDeclarations(t, repository, []string{unprovenAction})
	writeArchiveQAReport(t, repository, spec.VerdictPartial,
		"rows_blocked_declared: 1",
		"rows_blocked_finding: 0",
		"rows_blocked_environment: 0",
	)
	gittest.Run(t, repository, "add", "docs/specs")
	gittest.Run(t, repository, "commit", "-m", "docs: record passing QA")
	reviewedHead := strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))

	archiveResult, err := spec.Archive(spec.ArchiveRequest{
		SpecsRoot:   filepath.Join(repository, "docs", "specs"),
		BuiltInRoot: true,
		Slug:        implementTestSlug,
		ArchivedAt:  time.Date(2026, time.September, 24, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("archive Spec through real command boundary: %v", err)
	}
	if got := readArchivedUnproven(t, filepath.Join(archiveResult.ArchivedDir, "_prd.md")); !reflect.DeepEqual(got, []string{unprovenAction}) {
		t.Fatalf("real archive unproven actions = %q, want %q", got, []string{unprovenAction})
	}
	if mutate != nil {
		mutate(t, repository)
	}
	gittest.Run(t, repository, "add", "-A")
	gittest.Run(t, repository, "commit", "-m", "docs: archive "+implementTestSlug)
	archiveHead := strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
	return repository, reviewedHead, archiveHead
}

func resumeArchivedDelivery(t *testing.T, repository string, reviewedHead string) store.DeliveryQueueItem {
	t.Helper()
	runStore, err := store.Open(t.Context(), t.TempDir())
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close Run Database: %v", err)
		}
	})
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repository, []string{implementTestSlug})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageArchiving
	item.Branch = strings.TrimSpace(gittest.Run(t, repository, "branch", "--show-current"))
	item.Worktree = repository
	item.CandidateCommits = []string{reviewedHead}
	if _, _, err := runStore.RecordDeliveryQueueItemWorktree(
		t.Context(),
		repository,
		item.SpecSlug,
		item.Branch,
		item.Worktree,
	); err != nil {
		t.Fatalf("record archiving item worktree: %v", err)
	}
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repository, item); err != nil {
		t.Fatalf("seed archiving item: %v", err)
	}
	workflow := &commandDeliveryWorkflow{
		store: runStore,
		loaded: roundconfig.Loaded{
			GitRoot: repository,
			Config:  roundconfig.Config{Specs: roundconfig.Specs{Root: "docs/specs"}},
		},
		git: preflight.ExecGitRunner{},
	}
	flow := &parkTestDeliveryFlow{}
	engine := delivery.NewEngine(runStore, delivery.EngineDependencies{
		Workspace:    flow,
		Runner:       flow,
		Reviewer:     flow,
		Archiver:     workflow,
		Gate:         flow,
		Authorizer:   flow,
		Publication:  flow,
		PullRequests: flow,
	})
	if _, err := engine.Run(t.Context(), repository); err != nil {
		t.Fatalf("resume Delivery Engine after archive commit: %v", err)
	}
	resumed, found, err := runStore.DeliveryQueue(t.Context(), repository)
	if err != nil || !found {
		t.Fatalf("read resumed Delivery Queue: found=%v err=%v", found, err)
	}
	return resumed.Items[0]
}

func newDeliveryBranchRepository(t *testing.T) (string, string) {
	t.Helper()
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
	return origin, checkout
}

func newDeliveryBranchWorkflow(t *testing.T, checkout string, specSlug string) *commandDeliveryWorkflow {
	t.Helper()
	runStore, err := store.Open(t.Context(), t.TempDir())
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close Run Database: %v", err)
		}
	})
	if _, err := runStore.CreateDeliveryQueue(t.Context(), checkout, []string{specSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	return &commandDeliveryWorkflow{
		store: runStore,
		loaded: roundconfig.Loaded{
			GitRoot: checkout,
			Config: roundconfig.Config{
				Worktree: roundconfig.Worktree{Location: filepath.Join(t.TempDir(), "worktrees")},
			},
		},
		git: preflight.ExecGitRunner{},
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

type parkTestDeliveryFlow struct {
	parkSlug     string
	reviewedHead string
	workDir      string
	runWorkDir   string
	runCalls     int
	remoteHeads  map[string]string
	pullRequests map[string]delivery.PullRequest
}

func (flow *parkTestDeliveryFlow) CreateItemBranch(context.Context, string, string) (string, string, error) {
	return "", "", errors.New("unexpected item worktree creation")
}

func (flow *parkTestDeliveryFlow) UseItemBranch(_ context.Context, _ string, _ string, itemWorktree string) (string, error) {
	return itemWorktree, nil
}

func (flow *parkTestDeliveryFlow) RemoveItemBranch(context.Context, string, string, string) error {
	return nil
}

func (flow *parkTestDeliveryFlow) RunSpec(_ context.Context, workDir string, slug string) (delivery.RunResult, error) {
	flow.runWorkDir = workDir
	flow.runCalls++
	if slug == flow.parkSlug {
		return delivery.RunResult{Outcome: delivery.RunOutcomeUnresolved}, nil
	}
	head := flow.reviewedHead
	if head == "" {
		head = "reviewed-" + slug
	}
	return delivery.RunResult{
		RunID:            "run-" + slug,
		Outcome:          delivery.RunOutcomeClean,
		CandidateCommits: []string{head},
	}, nil
}

func (flow *parkTestDeliveryFlow) ReviewPolicy(context.Context, string, string) (delivery.ReviewPolicy, error) {
	return delivery.ReviewPolicyEnabled, nil
}

func (flow *parkTestDeliveryFlow) Review(_ context.Context, _, _, head string) (delivery.ReviewResult, error) {
	return delivery.ReviewResult{Outcome: delivery.ReviewOutcomeReviewed, Head: head}, nil
}

func (flow *parkTestDeliveryFlow) RecordReviewOmission(context.Context, string, string, string) error {
	return errors.New("unexpected review omission")
}

func (flow *parkTestDeliveryFlow) Archive(_ context.Context, _, slug string, reviewedHead string) (delivery.ArchiveResult, error) {
	return delivery.ArchiveResult{Parent: reviewedHead, Head: "archived-" + slug, ExactSpecMove: true}, nil
}

func (flow *parkTestDeliveryFlow) Gate(context.Context, string, string, string) (delivery.GateResult, error) {
	return delivery.GateResult{Passed: true}, nil
}

func (flow *parkTestDeliveryFlow) Authorization(context.Context, string, string) (delivery.Authorization, error) {
	return delivery.Authorization{Operations: []string{"push", "pull_request", "merge"}}, nil
}

func (flow *parkTestDeliveryFlow) Publication(_ context.Context, _ string, slug, branch string) (delivery.Publication, error) {
	return delivery.Publication{Remote: "origin", HeadBranch: branch, BaseBranch: "main", Title: slug}, nil
}

func (flow *parkTestDeliveryFlow) WithWorkDir(workDir string) delivery.PullRequestBoundary {
	flow.workDir = workDir
	if flow.remoteHeads == nil {
		flow.remoteHeads = map[string]string{}
	}
	if flow.pullRequests == nil {
		flow.pullRequests = map[string]delivery.PullRequest{}
	}
	return flow
}

func (flow *parkTestDeliveryFlow) RemoteBranchHead(_ context.Context, remote, branch string) (delivery.RemoteHead, bool, error) {
	head, found := flow.remoteHeads[branch]
	return delivery.RemoteHead{Remote: remote, Branch: branch, SHA: head}, found, nil
}

func (flow *parkTestDeliveryFlow) PushBranch(_ context.Context, remote, branch, head string) (delivery.RemoteHead, error) {
	flow.remoteHeads[branch] = head
	return delivery.RemoteHead{Remote: remote, Branch: branch, SHA: head}, nil
}

func (flow *parkTestDeliveryFlow) FindOrCreatePullRequest(_ context.Context, req delivery.PullRequestRequest) (delivery.PullRequestResult, error) {
	pullRequest := delivery.PullRequest{
		Number:     fmt.Sprintf("%d", len(flow.pullRequests)+1),
		State:      "OPEN",
		HeadBranch: req.HeadBranch,
		HeadSHA:    flow.remoteHeads[req.HeadBranch],
	}
	flow.pullRequests[pullRequest.Number] = pullRequest
	return delivery.PullRequestResult{PullRequest: pullRequest, Created: true}, nil
}

func (flow *parkTestDeliveryFlow) CurrentHeadChecks(_ context.Context, number string) (delivery.CheckReport, error) {
	pullRequest, found := flow.pullRequests[number]
	if !found {
		return delivery.CheckReport{}, errors.New("pull request not found")
	}
	return delivery.CheckReport{
		HeadSHA: pullRequest.HeadSHA,
		Checks:  []delivery.PullRequestCheck{{Name: "test", Bucket: "pass"}},
	}, nil
}

func (flow *parkTestDeliveryFlow) MergePullRequest(_ context.Context, number, expectedHead string) (delivery.MergeResult, error) {
	pullRequest, found := flow.pullRequests[number]
	if !found {
		return delivery.MergeResult{}, errors.New("pull request not found")
	}
	pullRequest.State = "MERGED"
	pullRequest.HeadSHA = expectedHead
	pullRequest.MergeCommit = "merge-" + number
	flow.pullRequests[number] = pullRequest
	return delivery.MergeResult{PullRequest: pullRequest}, nil
}

type fakeDeliveryOwnerProcesses struct {
	provedPID     int
	terminatedPID int
	identity      string
	proveErr      error
}

func (fake *fakeDeliveryOwnerProcesses) ProveOwner(_ context.Context, pid int, identity string) error {
	fake.provedPID = pid
	fake.identity = identity
	return fake.proveErr
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
