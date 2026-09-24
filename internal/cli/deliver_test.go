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

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
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

func TestDeliveryWorkflowCreatesAnItemBranchFromTheRefreshedDefault(t *testing.T) {
	t.Parallel()
	origin, checkout := newDeliveryBranchRepository(t)
	seedPath := filepath.Join(origin, "seed.txt")
	if err := os.WriteFile(seedPath, []byte("refreshed\n"), 0o644); err != nil {
		t.Fatalf("refresh origin seed: %v", err)
	}
	gittest.Run(t, origin, "add", "seed.txt")
	gittest.Run(t, origin, "commit", "-m", "fix: refresh default")
	wantHead := strings.TrimSpace(gittest.Run(t, origin, "rev-parse", "main"))
	workflow := newDeliveryBranchWorkflow(t, checkout, "0156-delivery")

	branch, err := workflow.CreateItemBranch(t.Context(), checkout, "0156-delivery")

	if err != nil {
		t.Fatalf("create item branch: %v", err)
	}
	if !strings.HasPrefix(branch, "roundfix/deliver-0156-delivery-") {
		t.Fatalf("item branch = %q, want per-delivery suffix", branch)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "branch", "--show-current")); got != branch {
		t.Fatalf("current branch = %q, want %q", got, branch)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "rev-parse", "HEAD")); got != wantHead {
		t.Fatalf("item branch head = %q, want refreshed default %q", got, wantHead)
	}
}

func TestItemBranchHasNoUpstream(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	workflow := newDeliveryBranchWorkflow(t, checkout, "0161-untracked")

	branch, err := workflow.CreateItemBranch(t.Context(), checkout, "0161-untracked")

	if err != nil {
		t.Fatalf("create item branch: %v", err)
	}
	upstream := strings.TrimSpace(gittest.Run(
		t,
		checkout,
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
	workflow := newDeliveryBranchWorkflow(t, checkout, "0161-repeat")

	firstBranch, err := workflow.CreateItemBranch(t.Context(), checkout, "0161-repeat")
	if err != nil {
		t.Fatalf("create first item branch: %v", err)
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
	gittest.Run(t, checkout, "switch", "main")
	if _, err := workflow.store.CreateDeliveryQueue(t.Context(), checkout, []string{"0161-repeat"}); err != nil {
		t.Fatalf("create second Delivery Queue: %v", err)
	}

	secondBranch, err := workflow.CreateItemBranch(t.Context(), checkout, "0161-repeat")

	if err != nil {
		t.Fatalf("create second item branch: %v", err)
	}
	if firstBranch == secondBranch {
		t.Fatalf("delivery branches = %q and %q, want different branches", firstBranch, secondBranch)
	}
}

func TestResumeReusesTheRecordedItemBranch(t *testing.T) {
	t.Parallel()
	_, checkout := newDeliveryBranchRepository(t)
	workflow := newDeliveryBranchWorkflow(t, checkout, "0161-resume")

	firstBranch, err := workflow.CreateItemBranch(t.Context(), checkout, "0161-resume")
	if err != nil {
		t.Fatalf("create item branch before crash: %v", err)
	}
	queue, found, err := workflow.store.DeliveryQueue(t.Context(), checkout)
	if err != nil || !found {
		t.Fatalf("read Delivery Queue after branch creation: found=%v err=%v", found, err)
	}
	if got := queue.Items[0].Branch; got != firstBranch {
		t.Fatalf("recorded item branch = %q, want %q", got, firstBranch)
	}
	if queue.Items[0].Stage != store.DeliveryStageQueued {
		t.Fatalf("item stage after simulated crash = %q, want queued", queue.Items[0].Stage)
	}
	gittest.Run(t, checkout, "switch", "main")
	gittest.Run(t, checkout, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing-origin"))

	resumedBranch, err := workflow.CreateItemBranch(t.Context(), checkout, "0161-resume")

	if err != nil {
		t.Fatalf("resume item branch: %v", err)
	}
	if resumedBranch != firstBranch {
		t.Fatalf("resumed item branch = %q, want recorded branch %q", resumedBranch, firstBranch)
	}
	if got := strings.TrimSpace(gittest.Run(t, checkout, "branch", "--show-current")); got != firstBranch {
		t.Fatalf("current branch after resume = %q, want %q", got, firstBranch)
	}
}

func TestAParkLeavesACleanCheckout(t *testing.T) {
	_, checkout := newDeliveryBranchRepository(t)
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
	const parkedSlug = "0161-non-exact-archive"
	const nextSlug = "0161-next-item"
	if _, err := runStore.CreateDeliveryQueue(ctx, checkout, []string{parkedSlug, nextSlug}); err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	workflow := &commandDeliveryWorkflow{
		store:  runStore,
		loaded: roundconfig.Loaded{GitRoot: checkout},
		git:    preflight.ExecGitRunner{},
	}
	reviewedHead := strings.TrimSpace(gittest.Run(t, checkout, "rev-parse", "HEAD"))
	flow := &parkTestDeliveryFlow{checkout: checkout, reviewedHead: reviewedHead}
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
	if got, want := flow.runs, []string{parkedSlug, nextSlug}; !reflect.DeepEqual(got, want) {
		t.Fatalf("started items = %v, want %v", got, want)
	}
	if queue.Items[0].Stage != store.DeliveryStageParked || queue.Items[1].Branch == "" {
		t.Fatalf("Delivery Queue after non-exact archive = %+v", queue.Items)
	}
	state, err := preflight.InspectGit(ctx, checkout, preflight.ExecGitRunner{})
	if err != nil {
		t.Fatalf("inspect checkout after parks: %v", err)
	}
	if state.Branch != "main" || len(state.Dirty) != 0 {
		t.Fatalf("checkout after parks = branch %q dirty=%v, want clean main", state.Branch, state.Dirty)
	}
}

func TestResumeAcceptsTheArchiveCommit(t *testing.T) {
	tests := []struct {
		name      string
		extraPath bool
		wantExact bool
	}{
		{name: "exact Spec move", wantExact: true},
		{name: "move with unrelated change", extraPath: true, wantExact: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := t.TempDir()
			gittest.InitRepo(t, repository, "--initial-branch=main")
			gittest.PersistIdentity(t, repository)
			const slug = "0161-archive-resume"
			source := filepath.Join(repository, "docs", "specs", slug)
			if err := os.MkdirAll(source, 0o755); err != nil {
				t.Fatalf("create active Spec: %v", err)
			}
			if err := os.MkdirAll(filepath.Join(repository, "docs", "specs", "still-active"), 0o755); err != nil {
				t.Fatalf("create remaining Spec: %v", err)
			}
			if err := os.WriteFile(filepath.Join(source, "_prd.md"), []byte("# Archived Spec\n"), 0o644); err != nil {
				t.Fatalf("write active Spec: %v", err)
			}
			if err := os.WriteFile(filepath.Join(repository, "docs", "specs", "still-active", "_prd.md"), []byte("# Active Spec\n"), 0o644); err != nil {
				t.Fatalf("write remaining Spec: %v", err)
			}
			gittest.Run(t, repository, "add", "docs/specs")
			gittest.Run(t, repository, "commit", "-m", "docs: add Specs")
			gittest.Run(t, repository, "switch", "-c", "roundfix/deliver-"+slug)
			reviewedHead := strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
			destination := filepath.Join(repository, "docs", "history", "specs", slug)
			if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
				t.Fatalf("create archive root: %v", err)
			}
			if err := os.Rename(source, destination); err != nil {
				t.Fatalf("move Spec to archive root: %v", err)
			}
			if tt.extraPath {
				if err := os.WriteFile(filepath.Join(repository, "unrelated.txt"), []byte("unrelated\n"), 0o644); err != nil {
					t.Fatalf("write unrelated change: %v", err)
				}
			}
			gittest.Run(t, repository, "add", "-A")
			gittest.Run(t, repository, "commit", "-m", "docs: archive "+slug)
			archiveHead := strings.TrimSpace(gittest.Run(t, repository, "rev-parse", "HEAD"))
			workflow := &commandDeliveryWorkflow{
				loaded: roundconfig.Loaded{
					GitRoot: repository,
					Config:  roundconfig.Config{Specs: roundconfig.Specs{Root: "docs/specs"}},
				},
				git: preflight.ExecGitRunner{},
			}

			result, err := workflow.Archive(t.Context(), repository, slug, reviewedHead)

			if err != nil {
				t.Fatalf("resume archive reconciliation: %v", err)
			}
			if result.ExactSpecMove != tt.wantExact {
				t.Fatalf("archive result = %+v, want exact=%v", result, tt.wantExact)
			}
			if tt.wantExact && (result.Parent != reviewedHead || result.Head != archiveHead) {
				t.Fatalf("archive result = %+v, want parent %q head %q", result, reviewedHead, archiveHead)
			}
			if !tt.wantExact {
				return
			}

			runStore, err := store.Open(t.Context(), t.TempDir())
			if err != nil {
				t.Fatalf("open Run Database: %v", err)
			}
			t.Cleanup(func() {
				if err := runStore.Close(); err != nil {
					t.Errorf("close Run Database: %v", err)
				}
			})
			queue, err := runStore.CreateDeliveryQueue(t.Context(), repository, []string{slug})
			if err != nil {
				t.Fatalf("create Delivery Queue: %v", err)
			}
			item := queue.Items[0]
			item.Stage = store.DeliveryStageArchiving
			item.Branch = "roundfix/deliver-" + slug
			item.CandidateCommits = []string{reviewedHead}
			if err := runStore.UpdateDeliveryQueueItem(t.Context(), repository, item); err != nil {
				t.Fatalf("seed archiving item: %v", err)
			}
			workflow.store = runStore
			flow := &parkTestDeliveryFlow{}
			engine := delivery.NewEngine(runStore, delivery.EngineDependencies{
				Workspace:    workflow,
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
			got := resumed.Items[0]
			if got.Blocker == delivery.BlockerReviewStale || !reflect.DeepEqual(got.CandidateCommits, []string{reviewedHead, archiveHead}) {
				t.Fatalf("resumed archive item = %+v, want archive head accepted past archiving", got)
			}
		})
	}
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
		store:  runStore,
		loaded: roundconfig.Loaded{GitRoot: checkout},
		git:    preflight.ExecGitRunner{},
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
	checkout     string
	reviewedHead string
	runs         []string
}

func (flow *parkTestDeliveryFlow) RunSpec(_ context.Context, _ string, slug string) (delivery.RunResult, error) {
	flow.runs = append(flow.runs, slug)
	if len(flow.runs) > 1 {
		return delivery.RunResult{Outcome: delivery.RunOutcomeUnresolved}, nil
	}
	return delivery.RunResult{
		RunID:            "run-" + slug,
		Outcome:          delivery.RunOutcomeClean,
		CandidateCommits: []string{flow.reviewedHead},
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

func (flow *parkTestDeliveryFlow) Archive(_ context.Context, _, _ string, reviewedHead string) (delivery.ArchiveResult, error) {
	if err := os.WriteFile(filepath.Join(flow.checkout, "seed.txt"), []byte("tracked archive change\n"), 0o644); err != nil {
		return delivery.ArchiveResult{}, err
	}
	if err := os.WriteFile(filepath.Join(flow.checkout, "archive-fragment.txt"), []byte("non-exact archive\n"), 0o644); err != nil {
		return delivery.ArchiveResult{}, err
	}
	return delivery.ArchiveResult{Parent: reviewedHead, ExactSpecMove: false}, nil
}

func (flow *parkTestDeliveryFlow) Gate(context.Context, string, string, string) (delivery.GateResult, error) {
	return delivery.GateResult{}, errors.New("unexpected repository gate")
}

func (flow *parkTestDeliveryFlow) Authorization(context.Context, string, string) (delivery.Authorization, error) {
	return delivery.Authorization{}, errors.New("unexpected authorization read")
}

func (flow *parkTestDeliveryFlow) Publication(context.Context, string, string, string) (delivery.Publication, error) {
	return delivery.Publication{}, errors.New("unexpected publication plan")
}

func (flow *parkTestDeliveryFlow) RemoteBranchHead(context.Context, string, string) (delivery.RemoteHead, bool, error) {
	return delivery.RemoteHead{}, false, errors.New("unexpected remote head read")
}

func (flow *parkTestDeliveryFlow) PushBranch(context.Context, string, string, string) (delivery.RemoteHead, error) {
	return delivery.RemoteHead{}, errors.New("unexpected push")
}

func (flow *parkTestDeliveryFlow) FindOrCreatePullRequest(context.Context, delivery.PullRequestRequest) (delivery.PullRequestResult, error) {
	return delivery.PullRequestResult{}, errors.New("unexpected pull request creation")
}

func (flow *parkTestDeliveryFlow) CurrentHeadChecks(context.Context, string) (delivery.CheckReport, error) {
	return delivery.CheckReport{}, errors.New("unexpected check read")
}

func (flow *parkTestDeliveryFlow) MergePullRequest(context.Context, string, string) (delivery.MergeResult, error) {
	return delivery.MergeResult{}, errors.New("unexpected merge")
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
