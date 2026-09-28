// Suite: Delivery Queue retry command and owner lifecycle.
// Invariant: one explicit retry reaches the proven recorded owner or a replacement, and an owner releases only an idle queue.
// Boundary IN: public CLI dispatch and the real SQLite Delivery Queue store.
// Boundary OUT: Delivery Engine execution, owner process control, and detached owner launch.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

func TestDeliverRetryStartsAnOwnerWhenNoneIsRunning(t *testing.T) {
	homeDir, repoDir := newParkedDeliveryQueueForRetry(t, 0, "")
	engine := &retryCommandDeliveryEngine{
		retryResult: delivery.RetryResult{
			SpecSlug: implementTestSlug,
			Blocker:  "run-unresolved",
			Stage:    store.DeliveryStageRunning,
			CarriedFrom: delivery.CarryForwardResult{
				RunID:   "run_123",
				Carried: []string{"task_02", "task_01"},
			},
		},
	}
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(runStore *store.Store, _ roundconfig.Loaded) deliveryEngine {
			if runStore == nil {
				t.Fatal("retry received a nil Run Database")
			}
			return engine
		}
		dependencies.startDeliveryOwner = func(_ context.Context, _ roundconfig.Loaded, _ commandEnvironment, stdout, _ io.Writer) int {
			started++
			fmt.Fprintln(stdout, "Delivery Owner: test-owner")
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "retry", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver retry exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if started != 1 {
		t.Fatalf("Delivery Queue owner starts = %d, want 1", started)
	}
	if engine.retryCalls != 1 || engine.gitRoot != repoDir || engine.specSlug != implementTestSlug {
		t.Fatalf("retry call = count %d root %q slug %q", engine.retryCalls, engine.gitRoot, engine.specSlug)
	}
	want := "Carried forward from Run run_123: task_02, task_01\n" +
		"Retried " + implementTestSlug + ": run-unresolved -> running\n" +
		"Delivery Owner: test-owner\n"
	if stdout.String() != want {
		t.Fatalf("deliver retry stdout = %q, want %q", stdout.String(), want)
	}
	if _, err := os.Stat(store.DatabasePath(homeDir)); err != nil {
		t.Fatalf("real Run Database missing after retry: %v", err)
	}
}

func TestDeliverRetryHandsTheItemToALiveOwner(t *testing.T) {
	ownerPID := os.Getpid()
	const ownerIdentity = "recorded-owner"
	_, _ = newParkedDeliveryQueueForRetry(t, ownerPID, ownerIdentity)
	engine := &retryCommandDeliveryEngine{
		retryResult: delivery.RetryResult{
			SpecSlug:      implementTestSlug,
			Blocker:       "review-stale",
			Stage:         store.DeliveryStageReviewing,
			OwnerPID:      ownerPID,
			OwnerIdentity: ownerIdentity,
		},
	}
	processes := &fakeDeliveryOwnerProcesses{}
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine { return engine }
		dependencies.ownerProcesses = processes
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "retry", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("deliver retry exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if started != 0 {
		t.Fatalf("live owner started %d replacement owners", started)
	}
	if processes.provedPID != ownerPID || processes.identity != ownerIdentity {
		t.Fatalf("owner proof = %+v", processes)
	}
	want := fmt.Sprintf(
		"Retried %s: review-stale -> reviewing\nHanded %s to Delivery Queue owner PID %d.\n",
		implementTestSlug,
		implementTestSlug,
		ownerPID,
	)
	if stdout.String() != want {
		t.Fatalf("deliver retry stdout = %q, want %q", stdout.String(), want)
	}
}

func TestDeliverRetryReclaimsAStaleOwnerAndStartsOne(t *testing.T) {
	ownerPID := os.Getpid()
	const ownerIdentity = "stale-owner"
	homeDir, repoDir := newParkedDeliveryQueueForRetry(t, ownerPID, ownerIdentity)
	engine := &retryCommandDeliveryEngine{
		retryResult: delivery.RetryResult{
			SpecSlug:      implementTestSlug,
			Blocker:       "review-stale",
			Stage:         store.DeliveryStageReviewing,
			OwnerPID:      ownerPID,
			OwnerIdentity: ownerIdentity,
		},
	}
	processes := &fakeDeliveryOwnerProcesses{proveErr: fmt.Errorf("identity mismatch: %w", store.ErrOwnerProcessIdentityUnproven)}
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine { return engine }
		dependencies.ownerProcesses = processes
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "retry", implementTestSlug}, &stdout, &stderr)

	if code != exitOK || started != 1 {
		t.Fatalf("deliver retry exit=%d starts=%d stderr=%q stdout=%q", code, started, stderr.String(), stdout.String())
	}
	wantNotice := fmt.Sprintf("roundfix: Delivery Queue owner PID %d has a different process identity; reclaimed its owner record.\n", ownerPID)
	if stderr.String() != wantNotice {
		t.Fatalf("stale owner notice = %q, want %q", stderr.String(), wantNotice)
	}
	queue := openDeliveryQueueForCLI(t, homeDir, repoDir)
	if queue.OwnerPID != 0 || queue.OwnerIdentity != "" {
		t.Fatalf("stale owner remains recorded: %+v", queue)
	}
}

func TestDeliverRetryRequiresOneSlug(t *testing.T) {
	assertDeliverRetryArgumentsRefused(t, []string{"deliver", "retry"}, "missing required Spec slug")
}

func TestDeliverRetryRefusesAnEmptySlug(t *testing.T) {
	assertDeliverRetryArgumentsRefused(t, []string{"deliver", "retry", "  "}, "Spec slug cannot be empty")
}

func TestDeliverRetryRefusesAnExtraArgument(t *testing.T) {
	assertDeliverRetryArgumentsRefused(t, []string{"deliver", "retry", implementTestSlug, "extra"}, "unexpected argument")
}

func TestDeliverRetryRefusesUnknownFlag(t *testing.T) {
	assertDeliverRetryArgumentsRefused(t, []string{"deliver", "retry", "--unknown", implementTestSlug}, "flag provided but not defined")
}

func TestDeliverRetryRefusalStartsNoOwner(t *testing.T) {
	_, _ = newParkedDeliveryQueueForRetry(t, 0, "")
	engine := &retryCommandDeliveryEngine{retryErr: errors.New("item is not parked")}
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine { return engine }
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "retry", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight || started != 0 || stdout.Len() != 0 {
		t.Fatalf("refused retry exit=%d starts=%d stdout=%q stderr=%q", code, started, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "Preflight failed") || !strings.Contains(stderr.String(), "item is not parked") {
		t.Fatalf("refused retry diagnostic = %q", stderr.String())
	}
}

// Sequential: detached-child setup owns a raw process file descriptor.
func TestDeliveryOwnerRunsAgainForAnItemRetriedDuringItsPass(t *testing.T) {
	homeDir, repoDir := newParkedDeliveryQueueForRetry(t, 0, "")
	var runStore *store.Store
	engine := &retryCommandDeliveryEngine{}
	engine.run = func(ctx context.Context, gitRoot string) (delivery.EngineResult, error) {
		engine.runCalls++
		queue, found, err := runStore.DeliveryQueue(ctx, gitRoot)
		if err != nil || !found {
			return delivery.EngineResult{}, fmt.Errorf("read queue: found=%t: %w", found, err)
		}
		item := queue.Items[0]
		if engine.runCalls == 1 {
			item.Stage = store.DeliveryStageRunning
			item.Blocker = ""
		} else {
			item.Stage = store.DeliveryStageMerged
		}
		if err := runStore.UpdateDeliveryQueueItem(ctx, gitRoot, item); err != nil {
			return delivery.EngineResult{}, err
		}
		return delivery.EngineResult{}, nil
	}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(opened *store.Store, _ roundconfig.Loaded) deliveryEngine {
			runStore = opened
			return engine
		}
	})
	configureDeliveryOwnerChild(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "resume"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("delivery owner exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if engine.runCalls != 2 {
		t.Fatalf("Delivery Engine passes = %d, want 2", engine.runCalls)
	}
	queue := openDeliveryQueueForCLI(t, homeDir, repoDir)
	if queue.OwnerPID != 0 || queue.Items[0].Stage != store.DeliveryStageMerged {
		t.Fatalf("queue after second pass = %+v", queue)
	}
}

// Sequential: detached-child setup owns a raw process file descriptor.
func TestDeliveryOwnerReleasesAnIdleQueueAfterOnePass(t *testing.T) {
	homeDir, repoDir := newParkedDeliveryQueueForRetry(t, 0, "")
	engine := &retryCommandDeliveryEngine{}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine { return engine }
	})
	configureDeliveryOwnerChild(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "resume"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("delivery owner exit=%d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if engine.runCalls != 1 {
		t.Fatalf("Delivery Engine passes = %d, want 1", engine.runCalls)
	}
	queue := openDeliveryQueueForCLI(t, homeDir, repoDir)
	if queue.OwnerPID != 0 || queue.OwnerIdentity != "" {
		t.Fatalf("idle owner remains recorded: %+v", queue)
	}
}

// Sequential: detached-child setup owns a raw process file descriptor.
func TestDeliveryOwnerFailsWhenIdleReleaseFails(t *testing.T) {
	_, repoDir := newParkedDeliveryQueueForRetry(t, 0, "")
	var runStore *store.Store
	engine := &retryCommandDeliveryEngine{}
	engine.run = func(ctx context.Context, gitRoot string) (delivery.EngineResult, error) {
		engine.runCalls++
		queue, found, err := runStore.DeliveryQueue(ctx, gitRoot)
		if err != nil || !found {
			return delivery.EngineResult{}, fmt.Errorf("read queue: found=%t: %w", found, err)
		}
		if _, err := runStore.ReleaseDeliveryQueueOwner(ctx, gitRoot, queue.OwnerPID, queue.OwnerIdentity); err != nil {
			return delivery.EngineResult{}, err
		}
		if err := runStore.ClaimDeliveryQueueOwner(ctx, gitRoot, 4242, "replacement-owner"); err != nil {
			return delivery.EngineResult{}, err
		}
		return delivery.EngineResult{}, nil
	}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(opened *store.Store, _ roundconfig.Loaded) deliveryEngine {
			runStore = opened
			return engine
		}
	})
	configureDeliveryOwnerChild(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "resume"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("delivery owner release failure exit=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "release idle Delivery Queue owner") || !strings.Contains(stderr.String(), "recorded owner is PID 4242") {
		t.Fatalf("delivery owner release diagnostic = %q", stderr.String())
	}
	queue := openDeliveryQueueForCLI(t, commandEnvironmentForTest(t).homeDir, repoDir)
	if queue.OwnerPID != 4242 {
		t.Fatalf("replacement owner PID = %d, want 4242", queue.OwnerPID)
	}
}

// Sequential: detached-child setup owns a raw process file descriptor.
func TestDeliveryOwnerReleasesItsClaimAfterEngineFailure(t *testing.T) {
	homeDir, repoDir := newParkedDeliveryQueueForRetry(t, 0, "")
	engine := &retryCommandDeliveryEngine{
		run: func(context.Context, string) (delivery.EngineResult, error) {
			return delivery.EngineResult{}, errors.New("engine pass failed")
		},
	}
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine { return engine }
	})
	configureDeliveryOwnerChild(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"deliver", "resume"}, &stdout, &stderr)

	if code != exitPreflight || !strings.Contains(stderr.String(), "engine pass failed") {
		t.Fatalf("delivery owner engine failure exit=%d stderr=%q", code, stderr.String())
	}
	queue := openDeliveryQueueForCLI(t, homeDir, repoDir)
	if queue.OwnerPID != 0 || queue.OwnerIdentity != "" {
		t.Fatalf("failed owner claim remains recorded: %+v", queue)
	}
}

func TestTopLevelUsageNamesDeliverRetry(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--help"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("top-level help exit=%d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "roundfix deliver <start|status|resume|retry|stop> [<slug> ...]") {
		t.Fatalf("top-level help does not name deliver retry:\n%s", stdout.String())
	}
}

func assertDeliverRetryArgumentsRefused(t *testing.T, args []string, diagnostic string) {
	t.Helper()
	_, _ = newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	engineBuilt := 0
	started := 0
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine {
			engineBuilt++
			return &retryCommandDeliveryEngine{}
		}
		dependencies.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
			started++
			return exitOK
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, args, &stdout, &stderr)

	if code != exitPreflight || stdout.Len() != 0 {
		t.Fatalf("invalid retry args %q exit=%d stdout=%q stderr=%q", args, code, stdout.String(), stderr.String())
	}
	if engineBuilt != 0 || started != 0 {
		t.Fatalf("invalid retry args built engine %d times and started owner %d times", engineBuilt, started)
	}
	if !strings.Contains(stderr.String(), diagnostic) {
		t.Fatalf("invalid retry args diagnostic = %q, want %q", stderr.String(), diagnostic)
	}
}

func newParkedDeliveryQueueForRetry(t *testing.T, ownerPID int, ownerIdentity string) (string, string) {
	t.Helper()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repoDir, []string{implementTestSlug})
	if err != nil {
		t.Fatalf("create Delivery Queue: %v", err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageParked
	item.Blocker = "run-unresolved"
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repoDir, item); err != nil {
		t.Fatalf("park Delivery Queue item: %v", err)
	}
	if ownerPID > 0 {
		if err := runStore.ClaimDeliveryQueueOwner(t.Context(), repoDir, ownerPID, ownerIdentity); err != nil {
			t.Fatalf("claim Delivery Queue owner: %v", err)
		}
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database: %v", err)
	}
	return homeDir, repoDir
}

func configureDeliveryOwnerChild(t *testing.T) {
	t.Helper()
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("create detached child pipe: %v", err)
	}
	t.Cleanup(func() {
		_ = readPipe.Close()
		_ = writePipe.Close()
	})
	tempPath := filepath.Join(t.TempDir(), "console.tmp")
	if err := os.WriteFile(tempPath, nil, 0o600); err != nil {
		t.Fatalf("create detached console temp file: %v", err)
	}
	setDetachChildEnvironmentForTest(t, strconv.Itoa(int(writePipe.Fd())), tempPath)
}

type retryCommandDeliveryEngine struct {
	retryResult delivery.RetryResult
	retryErr    error
	retryCalls  int
	runCalls    int
	gitRoot     string
	specSlug    string
	run         func(context.Context, string) (delivery.EngineResult, error)
}

func (engine *retryCommandDeliveryEngine) Run(ctx context.Context, gitRoot string) (delivery.EngineResult, error) {
	if engine.run != nil {
		return engine.run(ctx, gitRoot)
	}
	engine.runCalls++
	return delivery.EngineResult{}, nil
}

func (engine *retryCommandDeliveryEngine) Retry(_ context.Context, gitRoot, specSlug string) (delivery.RetryResult, error) {
	engine.retryCalls++
	engine.gitRoot = gitRoot
	engine.specSlug = specSlug
	return engine.retryResult, engine.retryErr
}
