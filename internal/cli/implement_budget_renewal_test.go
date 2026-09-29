package cli

// Suite: Implement Command renewing Run Budget
// Invariant: setup uses the start allowance and post-cycle work uses the last Task settlement allowance.
// Boundary IN: public Implement Command outcome plus integration context propagation.
// Boundary OUT: Task-cycle renewal and watchdog mechanics, covered by internal/daemon/task_budget_renewal_test.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

type cliBudgetClock struct {
	mu  sync.Mutex
	now time.Time
}

func (clock *cliBudgetClock) Now() time.Time {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.now
}

func (clock *cliBudgetClock) Set(now time.Time) {
	clock.mu.Lock()
	clock.now = now
	clock.mu.Unlock()
}

func configureImplementBudget(t *testing.T, repoDir string, maximum time.Duration) {
	t.Helper()
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "budget:\n  max_run_duration: "+maximum.String()+"\n")
	gitImplement(t, repoDir, "add", ".roundfixrc.yml")
	gitImplement(t, repoDir, "commit", "-m", "configure renewing implement budget")
}

func TestImplementSerialRunLongerThanItsBudgetEndsClean(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Settle the first slice"},
		{id: "task_02", title: "Settle the second slice", needs: []string{"task_01"}},
	})
	const maximum = time.Hour
	configureImplementBudget(t, repoDir, maximum)
	startedAt := time.Now()
	clock := &cliBudgetClock{now: startedAt}
	settlements := 0
	runner := &implementFakeRunner{
		gitRoot: repoDir,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
		onTask: func(agent.ExecuteRequest, string) error {
			settlements++
			clock.Set(startedAt.Add(time.Duration(settlements) * 54 * time.Minute))
			return nil
		},
	}
	withImplementCollaborators(t, runner)
	withVersionFreshnessFakeDeps(t, versionFreshnessDependencies{currentVersion: func() string { return "dev" }})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.implementBudgetNow = clock.Now
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("implement exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	run := implementRunFromStore(t, homeDir, implementRunIDFromStderr(t, stderr.String()))
	if run.State != store.StateClean {
		t.Fatalf("Run state = %q, want %q", run.State, store.StateClean)
	}
	if elapsed := clock.Now().Sub(startedAt); elapsed <= maximum {
		t.Fatalf("fixture elapsed = %s, want more than one %s allowance", elapsed, maximum)
	}
}

func TestImplementIntegrationRunsUnderTheRenewedDeadline(t *testing.T) {
	_, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Renew before integration"}})
	const maximum = time.Hour
	configureImplementBudget(t, repoDir, maximum)
	startedAt := time.Now()
	settledAt := startedAt.Add(30 * time.Minute)
	clock := &cliBudgetClock{now: startedAt}
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		onTask: func(agent.ExecuteRequest, string) error {
			clock.Set(settledAt)
			return nil
		},
	}
	withImplementCollaborators(t, runner)
	withVersionFreshnessFakeDeps(t, versionFreshnessDependencies{currentVersion: func() string { return "dev" }})
	var integrationDeadline time.Time
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.implementBudgetNow = clock.Now
		inner := dependencies.integrateRunWorktree
		dependencies.integrateRunWorktree = func(ctx context.Context, ref runworktree.Ref, branch string, message string) (runworktree.IntegrationResult, error) {
			integrationDeadline, _ = ctx.Deadline()
			return inner(ctx, ref, branch, message)
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("implement exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if want := settledAt.Add(maximum); !integrationDeadline.Equal(want) {
		t.Fatalf("integration deadline = %s, want renewed deadline %s", integrationDeadline, want)
	}
}

func TestImplementEndsBudgetExceededWhenTheRenewedDeadlinePasses(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Renew before integration"}})
	const maximum = time.Hour
	configureImplementBudget(t, repoDir, maximum)
	startedAt := time.Now()
	settledAt := startedAt.Add(30 * time.Minute)
	clock := &cliBudgetClock{now: startedAt}
	runner := &implementFakeRunner{
		gitRoot:      repoDir,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		onTask: func(agent.ExecuteRequest, string) error {
			clock.Set(settledAt)
			return nil
		},
	}
	withImplementCollaborators(t, runner)
	withVersionFreshnessFakeDeps(t, versionFreshnessDependencies{currentVersion: func() string { return "dev" }})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.implementBudgetNow = clock.Now
		dependencies.integrateRunWorktree = func(ctx context.Context, _ runworktree.Ref, _ string, _ string) (runworktree.IntegrationResult, error) {
			deadline, ok := ctx.Deadline()
			if !ok {
				return runworktree.IntegrationResult{}, errors.New("integration context has no Run Budget deadline")
			}
			clock.Set(deadline.Add(time.Second))
			return runworktree.IntegrationResult{}, errors.New("integration crossed the renewed deadline")
		}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("implement exit = %d, want %d; stderr=%q stdout=%q", code, exitRunFailed, stderr.String(), stdout.String())
	}
	runID := implementRunIDFromAnyStderrLine(t, stderr.String())
	run := implementRunFromStore(t, homeDir, runID)
	if run.State != store.StateBudgetExceeded {
		t.Fatalf("Run state = %q, want %q", run.State, store.StateBudgetExceeded)
	}
	var reason string
	for _, entry := range runEventsForRun(t, homeDir, runID) {
		if entry.Event.Kind != runevent.KindDaemonOutcome {
			continue
		}
		var outcome runevent.OutcomePayload
		if err := json.Unmarshal(entry.Event.Payload, &outcome); err != nil {
			t.Fatalf("decode outcome payload: %v", err)
		}
		reason = outcome.Reason
	}
	if !strings.HasSuffix(reason, " since Task task_01 settled.") {
		t.Fatalf("BudgetExceeded reason = %q, want renewing Task suffix", reason)
	}
}
