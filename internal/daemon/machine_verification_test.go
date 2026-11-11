package daemon

// Suite: Machine Verification Capacity at the Daemon boundary.
// Boundary: real Run gate and Task entry points, fake machine grants and verifiers.
// No child Verification command or Agent is started.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
	"roundfix/internal/testwait"
)

type machineGrantRequest struct {
	holder store.MachineVerificationHolder
	grant  chan error
}

type recordingMachineGate struct {
	mu                       sync.Mutex
	requests                 chan machineGrantRequest
	run                      verificationGate
	wait                     store.MachineVerificationWait
	active, grants, releases int
	violations               []string
}

func newRecordingMachineGate(run verificationGate) *recordingMachineGate {
	return &recordingMachineGate{run: run, requests: make(chan machineGrantRequest, 8)}
}

func (gate *recordingMachineGate) assertRunHeld() {
	if release, available := gate.run.TryAcquire(verificationExclusive); available {
		release()
		gate.mu.Lock()
		gate.violations = append(gate.violations, "machine acquired or released without the Run unit")
		gate.mu.Unlock()
	}
}

func (gate *recordingMachineGate) Acquire(ctx context.Context, holder store.MachineVerificationHolder, onWait func(store.MachineVerificationWait)) (func(), error) {
	gate.assertRunHeld()
	// Repeated snapshots must produce only one waiting event and one pause.
	onWait(gate.wait)
	onWait(gate.wait)
	req := machineGrantRequest{holder: holder, grant: make(chan error, 1)}
	select {
	case gate.requests <- req:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	select {
	case err := <-req.grant:
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	gate.mu.Lock()
	gate.active++
	gate.grants++
	gate.mu.Unlock()
	return func() {
		gate.assertRunHeld()
		gate.mu.Lock()
		gate.active--
		gate.releases++
		gate.mu.Unlock()
	}, nil
}

func (gate *recordingMachineGate) check(t *testing.T, grants int) {
	t.Helper()
	gate.mu.Lock()
	defer gate.mu.Unlock()
	if len(gate.violations) != 0 || gate.active != 0 || gate.grants != grants || gate.releases != grants {
		t.Fatalf("machine active=%d grants=%d releases=%d violations=%v; want %d grants and releases", gate.active, gate.grants, gate.releases, gate.violations, grants)
	}
}

type machineCheckingVerifier struct {
	mu        sync.Mutex
	gate      *recordingMachineGate
	requests  []VerifyRequest
	temporary bool
	hold      <-chan struct{}
	started   chan VerifyRequest
}

func (v *machineCheckingVerifier) Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error) {
	v.mu.Lock()
	v.requests = append(v.requests, req)
	number := len(v.requests)
	v.mu.Unlock()
	v.gate.mu.Lock()
	held := v.gate.active
	v.gate.mu.Unlock()
	if held != 1 {
		return VerifyResult{}, fmt.Errorf("Verification started with %d machine units", held)
	}
	if v.started != nil {
		select {
		case v.started <- req:
		case <-ctx.Done():
			return VerifyResult{}, ctx.Err()
		}
	}
	if v.hold != nil {
		select {
		case <-v.hold:
		case <-ctx.Done():
			return VerifyResult{}, ctx.Err()
		}
	}
	result := VerifyResult{OutputPath: req.OutputPath}
	if v.temporary && number == 1 {
		return result, &TemporaryVerificationFailureError{CommandFailure: &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("exit status 75")}}
	}
	return result, nil
}

func (v *machineCheckingVerifier) count() int {
	v.mu.Lock()
	defer v.mu.Unlock()
	return len(v.requests)
}

func machineFixture(t *testing.T, seeds []taskSpecSeed) (*taskCycleFixture, TaskPlan, *Engine, *recordingMachineGate, *machineCheckingVerifier) {
	t.Helper()
	fixture := newTaskCycleFixture(t, seeds)
	plan := fixture.plan()
	plan.verificationGate = newVerificationGate(1)
	machine := newRecordingMachineGate(plan.verificationGate)
	verifier := &machineCheckingVerifier{gate: machine}
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MachineVerification = machine
	return fixture, plan, engine, machine, verifier
}

func TestEveryDaemonVerificationStepTakesTheMachineCapacityAfterTheRunsOwn(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name  string
		steps []string
	}{
		{"pre-work", []string{"Task task_01 pre-work probe"}},
		{"precondition", []string{"Task task_01 repository precondition"}},
		{"settlement", []string{"Task task_01 Verification attempt 1"}},
		{"retry", []string{"Task task_01 Verification attempt 1", "Task task_01 Verification attempt 1 retry 1"}},
		{"integrated", []string{"Integrated Repository Verification"}},
		{"integrated-retry", []string{"Integrated Repository Verification", "Integrated Repository Verification"}},
		{"qa", []string{"QA Task task_01 repository Verification"}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			_, plan, engine, machine, verifier := machineFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"repository check"}}})
			// The retry owns the entire Run capacity but only one machine unit.
			plan.VerificationConcurrency = 2
			plan.verificationGate = newVerificationGate(2)
			machine.run = plan.verificationGate
			plan.RepositoryVerification = "repository check"
			verifier.temporary = len(row.steps) == 2
			task := plan.Tasks[0]
			if row.name == "settlement" || row.name == "retry" {
				task.Verification = []string{"declared check", "repository check"}
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan error, 1)
			go func() {
				var err error
				switch row.name {
				case "pre-work":
					_, err = engine.verifyTaskPreWork(ctx, plan, task, 1)
				case "precondition":
					_, _, err = engine.verifyRepositoryPrecondition(ctx, plan, task, 1)
				case "integrated", "integrated-retry":
					var failure string
					failure, err = engine.verifyIntegratedRepository(ctx, plan, 1)
					if err == nil && failure != "" {
						err = errors.New(failure)
					}
				case "qa":
					_, _, err = engine.runQARepositoryVerification(ctx, plan, task, 1, &speccheck.MechanicalResult{})
				default:
					retry := false
					var outcome verificationAttemptOutcome
					outcome, err = engine.verifyTask(ctx, plan, task, 1, 1, &retry)
					if err == nil && outcome.Failure != "" {
						err = errors.New(outcome.Failure)
					}
				}
				done <- err
			}()
			for i, step := range row.steps {
				request := testwait.Until(t, step, machine.requests, done)
				want := store.MachineVerificationHolder{Kind: store.MachineVerificationHolderRun, RunID: plan.RunID, SpecSlug: plan.Spec.Slug, Step: step}
				if request.holder != want {
					t.Fatalf("holder=%+v, want %+v", request.holder, want)
				}
				if got := verifier.count(); got != i {
					t.Fatalf("commands before grant=%d, want %d", got, i)
				}
				if i == 1 {
					if release, available := plan.verificationGate.TryAcquire(verificationShared); available {
						release()
						t.Fatal("exclusive retry left a Run unit free")
					}
				}
				request.grant <- nil
			}
			if err := testwait.Until(t, "Verification to return", done, (<-chan struct{})(nil)); err != nil {
				t.Fatal(err)
			}
			machine.check(t, len(row.steps))
			release, available := plan.verificationGate.TryAcquire(verificationExclusive)
			if !available {
				t.Fatal("Run capacity leaked")
			}
			release()
			select {
			case extra := <-machine.requests:
				t.Fatalf("extra machine request: %+v", extra.holder)
			default:
			}
			if row.name == "settlement" || row.name == "retry" {
				if got := verifier.requests[len(verifier.requests)-1].Command; got != "repository check" {
					t.Fatalf("last settlement command=%q", got)
				}
			}
		})
	}
	t.Run("two-ready-Tasks", testMachineCapacityTakesRunFirstForTwoReadyTasks)
}

type observedRunVerificationGate struct {
	verificationGate
	queued chan struct{}
}

func (gate *observedRunVerificationGate) Acquire(ctx context.Context, mode verificationMode) (func(), error) {
	gate.queued <- struct{}{}
	return gate.verificationGate.Acquire(ctx, mode)
}

func testMachineCapacityTakesRunFirstForTwoReadyTasks(t *testing.T) {
	t.Helper()
	_, plan, engine, machine, verifier := machineFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"first"}}, {id: "task_02", verification: []string{"second"}},
	})
	runGate := &observedRunVerificationGate{verificationGate: plan.verificationGate, queued: make(chan struct{}, 1)}
	plan.verificationGate = runGate
	machine.run = runGate
	commandDone := make(chan struct{})
	verifier.hold = commandDone
	verifier.started = make(chan VerifyRequest, 2)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 2)
	runTask := func(task spec.Task, ordinal int) {
		retry := false
		_, err := engine.verifyTask(ctx, plan, task, ordinal, 1, &retry)
		done <- err
	}
	go runTask(plan.Tasks[0], 1)
	first := testwait.Until(t, "first machine request", machine.requests, done)
	first.grant <- nil
	testwait.Until(t, "first command to start", verifier.started, done)
	go runTask(plan.Tasks[1], 2)
	testwait.Until(t, "second Task to queue for its Run", runGate.queued, done)
	select {
	case request := <-machine.requests:
		t.Fatalf("second machine request while first held: %+v", request.holder)
	default:
	}
	close(commandDone)
	second := testwait.Until(t, "second machine request after first release", machine.requests, (<-chan struct{})(nil))
	second.grant <- nil
	for range 2 {
		if err := testwait.Until(t, "Task Verification", done, (<-chan struct{})(nil)); err != nil {
			t.Fatal(err)
		}
	}
	machine.check(t, 2)
}

func TestAStopRequestWhileWaitingForTheMachineCapacityStartsNoCommand(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"stop", "cancel", "error", "cancel-at-grant", "stop-at-grant"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			fixture, plan, engine, machine, verifier := machineFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"check"}}})
			clock := &budgetTestClock{now: time.Now()}
			engine.deps.Now = clock.Now
			plan.runBudget = newTaskCycleBudget(clock.Now(), time.Minute, true)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			done := make(chan error, 1)
			go func() { _, err := engine.verifyTaskPreWork(ctx, plan, plan.Tasks[0], 1); done <- err }()
			request := testwait.Until(t, "machine wait", machine.requests, done)
			clock.Set(clock.Now().Add(2 * time.Minute))
			scripted := errors.New("machine store unavailable")
			switch mode {
			case "stop", "stop-at-grant":
				if err := fixture.store.RequestStop(t.Context(), plan.RunID); err != nil {
					t.Fatal(err)
				}
				if mode == "stop-at-grant" {
					request.grant <- nil
				}
			case "error":
				request.grant <- scripted
			case "cancel":
				cancel()
			case "cancel-at-grant":
				cancel()
				request.grant <- nil
			}
			err := testwait.Until(t, "wait to end", done, (<-chan struct{})(nil))
			want := context.Canceled
			if strings.HasPrefix(mode, "stop") {
				want = ErrStopRequested
			}
			if mode == "error" {
				want = scripted
			}
			if !errors.Is(err, want) {
				t.Fatalf("error=%v, want %v", err, want)
			}
			if mode == "error" && !strings.Contains(err.Error(), fmt.Sprintf("acquire Machine Verification Capacity for run %q Task task_01 attempt 1 retry 0", plan.RunID)) {
				t.Fatalf("unwrapped error=%v", err)
			}
			if verifier.count() != 0 {
				t.Fatal("command started while waiting")
			}
			release, available := plan.verificationGate.TryAcquire(verificationExclusive)
			if !available {
				t.Fatal("Run unit leaked")
			}
			release()
			deadline, _, _, _ := plan.runBudget.snapshot()
			if !deadline.Equal(clock.Now().Add(time.Minute)) {
				t.Fatalf("budget not resumed on %s: %s", mode, deadline)
			}
			if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != "pending" {
				t.Fatalf("status=%s, want pending", got)
			}
		})
	}
}

func TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity(t *testing.T) {
	t.Parallel()
	started := time.Now()
	budget := newTaskCycleBudget(started, time.Minute, true)
	resume := budget.pause(started.Add(30 * time.Second))
	if budget.expiredAt(started.Add(2 * time.Minute)) {
		t.Fatal("paused budget expired")
	}
	resume(started.Add(150 * time.Second))
	if budget.expiredAt(started.Add(179 * time.Second)) {
		t.Fatal("budget expired before three minutes")
	}
	if !budget.expiredAt(started.Add(3 * time.Minute)) {
		t.Fatal("budget did not expire at three minutes")
	}

	t.Run("overlap-renew-and-disabled", func(t *testing.T) {
		budget := newTaskCycleBudget(started, time.Minute, true)
		first := budget.pause(started.Add(30 * time.Second))
		second := budget.pause(started.Add(40 * time.Second))
		first(started.Add(50 * time.Second))
		first(started.Add(60 * time.Second)) // duplicate resume cannot end the second wait
		if !budget.renew(started.Add(2*time.Minute), "task_02") {
			t.Fatal("settlement during pause did not renew")
		}
		if budget.expiredAt(started.Add(5 * time.Minute)) {
			t.Fatal("renewed pause expired")
		}
		second(started.Add(5 * time.Minute))
		deadline, _, task, _ := budget.snapshot()
		if !deadline.Equal(started.Add(6*time.Minute)) || task != "task_02" {
			t.Fatalf("renewed deadline=%s task=%s", deadline, task)
		}
		disabled := newTaskCycleBudget(started, time.Minute, false)
		disabled.pause(started)(started.Add(time.Hour))
		deadline, _, _, _ = disabled.snapshot()
		if !deadline.IsZero() {
			t.Fatal("disabled budget gained a deadline")
		}
	})
	t.Run("watch-timer-fires-during-pause", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithCancelCause(t.Context())
			defer cancel(nil)
			budget := newTaskCycleBudget(time.Now(), time.Minute, true)
			timer := make(chan time.Time, 1)
			done := make(chan struct{})
			go func() {
				defer close(done)
				budget.watch(ctx, cancel, func(time.Duration) (<-chan time.Time, func() bool) {
					return timer, func() bool { return true }
				})
			}()
			synctest.Wait()
			resume := budget.pause(time.Now())
			timer <- time.Now().Add(time.Minute)
			synctest.Wait()
			if ctx.Err() != nil {
				t.Fatalf("watch cancelled paused budget: %v", ctx.Err())
			}
			resume(time.Now().Add(2 * time.Minute))
			synctest.Wait()
			timer <- time.Now().Add(3 * time.Minute)
			synctest.Wait()
			if !errors.Is(context.Cause(ctx), context.DeadlineExceeded) {
				t.Fatalf("resumed watch cause=%v", context.Cause(ctx))
			}
			<-done
		})
	})
}

func TestTheMachineCapacityWaitIsJournaledAndPrinted(t *testing.T) {
	t.Parallel()
	fixture, plan, engine, machine, _ := machineFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"check"}}})
	started := time.Now()
	clock := &budgetTestClock{now: started}
	engine.deps.Now = clock.Now
	plan.runBudget = newTaskCycleBudget(started, time.Minute, true)
	machine.wait = store.MachineVerificationWait{
		Capacity: 2,
		Holding: []store.MachineVerificationEntry{
			{Holder: store.MachineVerificationHolder{Kind: store.MachineVerificationHolderRun, RunID: "run-other", Repository: "/other", SpecSlug: "other-spec", Step: "other Verification"}},
			{Holder: store.MachineVerificationHolder{Kind: store.MachineVerificationHolderDelivery, Repository: "/delivery", SpecSlug: "delivery-spec", Step: "repository gate"}},
		},
		Ahead: []store.MachineVerificationEntry{{Ticket: 4}},
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := engine.verifyTaskPreWork(ctx, plan, plan.Tasks[0], 1); done <- err }()
	request := testwait.Until(t, "waiting event", machine.requests, done)
	if plan.runBudget.expiredAt(started.Add(2 * time.Minute)) {
		t.Fatal("machine wait did not pause the budget")
	}
	clock.Set(started.Add(12 * time.Second))
	request.grant <- nil
	if err := testwait.Until(t, "grant and event", done, (<-chan struct{})(nil)); err != nil {
		t.Fatal(err)
	}
	holders := []string{store.DescribeMachineVerificationHolder(machine.wait.Holding[0].Holder) + ", other Verification", store.DescribeMachineVerificationHolder(machine.wait.Holding[1].Holder) + ", repository gate"}
	step := "Task task_01 pre-work probe"
	waiting := step + " waits for the Machine Verification Capacity: 2 of 2 held by " + strings.Join(holders, "; ") + "."
	acquired := step + " took the Machine Verification Capacity after waiting 12s."
	events := eventsOfKind(fixture.sink, runevent.KindDaemonStatus)
	if len(events) != 2 || events[0].Summary != waiting || events[1].Summary != acquired {
		t.Fatalf("events=%+v", events)
	}
	for i, want := range []map[string]any{
		{"machine_verification": "waiting", "step": step, "capacity": float64(2), "held": float64(2), "waiting_ahead": float64(1), "holders": []any{holders[0], holders[1]}},
		{"machine_verification": "acquired", "step": step, "waited_seconds": float64(12)},
	} {
		var got map[string]any
		if err := json.Unmarshal(events[i].Payload, &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("payload=%v, want %v", got, want)
		}
	}
	if fixture.progress.String() != waiting+"\n" {
		t.Fatalf("progress=%q", fixture.progress.String())
	}
	machine.check(t, 1)
}

// Ensure the optional nil dependency preserves immediate acquisition without
// requiring a Run store, event sink or clock.
func TestNilMachineVerificationPreservesTheImmediateRunGrant(t *testing.T) {
	t.Parallel()
	engine := &Engine{deps: Dependencies{Progress: io.Discard}}
	plan := TaskPlan{verificationGate: newVerificationGate(1)}
	release, err := engine.acquireVerificationCapacity(t.Context(), plan, spec.Task{}, verificationAttemptRequest{Mode: verificationShared})
	if err != nil {
		t.Fatal(err)
	}
	release()
}
