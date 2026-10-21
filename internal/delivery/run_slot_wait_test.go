// Suite: bounded Run Slot Waits.
// Invariant: only ceiling refusals wait; polling never grants capacity or extends a deadline.
// Boundary IN: Delivery Engine, real Run Database, fake clock and event-driven sleepers.
// Boundary OUT: Implement and slot observations are injected; race admission uses store.CreateRun.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/store"
	"roundfix/internal/testwait"
)

type slotRunnerFunc func(context.Context, string, string) (RunResult, error)

func (f slotRunnerFunc) RunSpec(ctx context.Context, root, slug string) (RunResult, error) {
	return f(ctx, root, slug)
}

type slotReaderFunc func(context.Context) (RunSlots, error)

func (f slotReaderFunc) RunSlots(ctx context.Context) (RunSlots, error) { return f(ctx) }

type slotSleepFunc func(context.Context, time.Duration) error

func (f slotSleepFunc) Sleep(ctx context.Context, d time.Duration) error { return f(ctx, d) }

func heldRunSlot() RunSlots {
	return RunSlots{MaxActive: 1, Holders: []RunSlotHolder{{RunID: "run_holder", Repository: "/holder", SpecSlug: "holder"}}}
}

func newSlotWaitEngine(t *testing.T, db *store.Store, root string, deadline time.Time) (*Engine, *fakeDeliveryClock) {
	t.Helper()
	if _, err := db.CreateDeliveryQueueWithLimits(t.Context(), root, []string{"waiting"}, store.DeliveryQueueLimits{Deadline: deadline}); err != nil {
		t.Fatal(err)
	}
	clock := &fakeDeliveryClock{now: time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)}
	workflow := newFakeDeliveryWorkflow()
	engine := NewEngine(db, EngineDependencies{
		Workspace: workflow, Reviewer: workflow, Archiver: workflow, Gate: workflow,
		Authorizer: workflow, Publication: workflow, PullRequests: newFakeDeliveryBoundary(),
		Revalidator: cleanTestRevalidator{}, Clock: clock, Sleeper: clock,
		RunSlotWait: time.Minute,
		RunSlots:    slotReaderFunc(func(context.Context) (RunSlots, error) { return heldRunSlot(), nil }),
		Runner: slotRunnerFunc(func(context.Context, string, string) (RunResult, error) {
			return RunResult{}, &RunSlotRefusal{Err: errors.New("ceiling refusal\nReason: busy"), Slots: heldRunSlot()}
		}),
	})
	return engine, clock
}

func TestARunCeilingRefusalWaitsAndStartsWhenASlotFrees(t *testing.T) {
	t.Parallel()
	db := openDeliveryEngineStore(t, t.Context())
	engine, clock := newSlotWaitEngine(t, db, "/queue", time.Time{})
	var log bytes.Buffer
	engine.log = &log
	calls, polls := 0, 0
	engine.runner = slotRunnerFunc(func(context.Context, string, string) (RunResult, error) {
		calls++
		if calls == 1 {
			return RunResult{}, fmt.Errorf("preflight: %w", &RunSlotRefusal{Err: errors.New("ceiling"), Slots: heldRunSlot()})
		}
		got := readDeliveryQueue(t, t.Context(), db, "/queue").Items[0]
		if got.Stage != store.DeliveryStageRunning || got.Blocker != "" {
			t.Fatalf("before retry = %+v", got)
		}
		return RunResult{RunID: "new-run", Outcome: RunOutcomeUnresolved}, nil
	})
	engine.sleeper = slotSleepFunc(func(ctx context.Context, d time.Duration) error {
		got := readDeliveryQueue(t, ctx, db, "/queue").Items[0]
		slots := heldRunSlot()
		if polls > 0 {
			slots.Holders[0].RunID = "replacement"
		}
		want := "run-slot-wait: until 2026-10-21T00:01:00Z; 1 Run(s) hold 1 slot(s): " + slots.Holders[0].RunID + " in /holder for holder"
		if got.Stage != store.DeliveryStageRunning || got.Blocker != want {
			t.Fatalf("waiting = %+v, want %q", got, want)
		}
		return clock.Sleep(ctx, d)
	})
	engine.runSlotPollInterval = 10 * time.Second
	engine.runSlots = slotReaderFunc(func(context.Context) (RunSlots, error) {
		polls++
		if polls == 1 {
			slots := heldRunSlot()
			slots.Holders[0].RunID = "replacement"
			return slots, nil
		}
		return RunSlots{MaxActive: 1}, nil
	})
	if _, err := engine.Run(t.Context(), "/queue"); err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, t.Context(), db, "/queue").Items[0]
	if calls != 2 || polls != 2 || got.RunID != "new-run" || got.Blocker != BlockerRunUnresolved {
		t.Fatalf("calls=%d polls=%d item=%+v", calls, polls, got)
	}
	wantLog := "roundfix: Delivery Queue item waiting: waiting for a Run slot until 2026-10-21T00:01:00Z; 1 Run(s) hold 1 slot(s)\nroundfix: Delivery Queue item waiting: a Run slot is free; starting the Implement Run\n"
	if log.String() != wantLog {
		t.Fatalf("log = %q, want %q", log.String(), wantLog)
	}
	// Formatting also handles an unlimited ceiling and a deadline with a non-UTC zone.
	if !(RunSlots{}).Free() || (heldRunSlot()).Free() {
		t.Fatal("slot availability disagrees with ceiling")
	}
	if got := RunSlotWaitText(clock.Now().In(time.FixedZone("BRT", -3*60*60)), RunSlots{}); got != "run-slot-wait: until 2026-10-21T00:00:20Z; 0 Run(s) hold 0 slot(s)" {
		t.Fatalf("empty wait text = %q", got)
	}
}

func TestARunSlotWaitParksAsTodayAtItsDeadline(t *testing.T) {
	t.Parallel()
	for _, lostRace := range []bool{false, true} {
		t.Run(fmt.Sprint(lostRace), func(t *testing.T) {
			db := openDeliveryEngineStore(t, t.Context())
			noWait, _ := newSlotWaitEngine(t, db, "/no-wait", time.Time{})
			noWait.runSlotWait = 0
			if lostRace {
				noWait.runner = slotRunnerFunc(func(context.Context, string, string) (RunResult, error) {
					return RunResult{}, &RunSlotRefusal{Err: errors.New("latest ceiling refusal"), Slots: heldRunSlot()}
				})
			}
			if _, err := noWait.Run(t.Context(), "/no-wait"); err != nil {
				t.Fatal(err)
			}
			want := readDeliveryQueue(t, t.Context(), db, "/no-wait").Items[0].Blocker
			engine, clock := newSlotWaitEngine(t, db, "/wait", time.Time{})
			start := clock.Now()
			calls := 0
			base := engine.runner
			engine.runner = slotRunnerFunc(func(ctx context.Context, root, slug string) (RunResult, error) {
				calls++
				if lostRace && calls > 1 {
					return noWait.runner.RunSpec(ctx, root, slug)
				}
				return base.RunSpec(ctx, root, slug)
			})
			if lostRace {
				engine.runSlots = slotReaderFunc(func(context.Context) (RunSlots, error) { return RunSlots{MaxActive: 1}, nil })
			}
			engine.runSlotWait = 45 * time.Second // final sleep must be shortened to 15s.
			if _, err := engine.Run(t.Context(), "/wait"); err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueue(t, t.Context(), db, "/wait").Items[0]
			if got.Stage != store.DeliveryStageParked || got.Blocker != want {
				t.Fatalf("item=%+v, no-wait reason=%q", got, want)
			}
			if clock.Now().Sub(start) != 45*time.Second || clock.sleeps != 2 {
				t.Fatalf("clock=%v sleeps=%d", clock.Now(), clock.sleeps)
			}
			wantCalls := 1
			if lostRace {
				wantCalls = 2
			}
			if calls != wantCalls {
				t.Fatalf("RunSpec calls=%d, want %d", calls, wantCalls)
			}
		})
	}
}

func TestTheQueueDeadlineEndsARunSlotWait(t *testing.T) {
	t.Parallel()
	db := openDeliveryEngineStore(t, t.Context())
	deadline := time.Date(2026, 10, 21, 0, 0, 15, 0, time.UTC)
	engine, clock := newSlotWaitEngine(t, db, "/queue", deadline)
	if _, err := engine.Run(t.Context(), "/queue"); err != nil {
		t.Fatal(err)
	}
	got := readDeliveryQueue(t, t.Context(), db, "/queue").Items[0]
	if got.Stage != store.DeliveryStageParked || got.Blocker != BlockerQueueDeadline || !clock.Now().Equal(deadline) {
		t.Fatalf("item=%+v now=%v", got, clock.Now())
	}
}

func TestAZeroRunSlotWaitParksAtOnce(t *testing.T) {
	t.Parallel()
	for _, unsetReader := range []bool{false, true} {
		t.Run(fmt.Sprint(unsetReader), func(t *testing.T) {
			db := openDeliveryEngineStore(t, t.Context())
			engine, clock := newSlotWaitEngine(t, db, "/queue", time.Time{})
			if unsetReader {
				engine.runSlots = nil
			} else {
				engine.runSlotWait = 0
			}
			if _, err := engine.Run(t.Context(), "/queue"); err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueue(t, t.Context(), db, "/queue").Items[0]
			if clock.sleeps != 0 || got.Blocker != "delivery-error: run Implement executor: ceiling refusal\nReason: busy" {
				t.Fatalf("item=%+v sleeps=%d", got, clock.sleeps)
			}
		})
	}
}

func TestOtherRunSpecFailuresDoNotWait(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                   string
		afterWait, readerError bool
	}{{name: "other refusal"}, {name: "failure after free", afterWait: true}, {name: "reader failure", afterWait: true, readerError: true}} {
		t.Run(tc.name, func(t *testing.T) {
			db := openDeliveryEngineStore(t, t.Context())
			engine, clock := newSlotWaitEngine(t, db, "/queue", time.Time{})
			base := engine.runner
			calls := 0
			engine.runner = slotRunnerFunc(func(ctx context.Context, root, slug string) (RunResult, error) {
				calls++
				if tc.afterWait && calls == 1 {
					return base.RunSpec(ctx, root, slug)
				}
				return RunResult{}, errors.New("another failure")
			})
			engine.runSlots = slotReaderFunc(func(context.Context) (RunSlots, error) {
				if tc.readerError {
					return RunSlots{}, errors.New("unreadable")
				}
				return RunSlots{MaxActive: 1}, nil
			})
			if _, err := engine.Run(t.Context(), "/queue"); err != nil {
				t.Fatal(err)
			}
			want := "delivery-error: run Implement executor: another failure"
			if tc.readerError {
				want = "delivery-error: wait for a Run slot: unreadable"
			}
			got := readDeliveryQueue(t, t.Context(), db, "/queue").Items[0]
			if got.Blocker != want {
				t.Fatalf("blocker=%q, want %q", got.Blocker, want)
			}
			if !tc.afterWait && clock.sleeps != 0 {
				t.Fatalf("slept %d times", clock.sleeps)
			}
		})
	}
}

func TestAStoppedOwnerClearsTheRunSlotWait(t *testing.T) {
	t.Parallel()
	for _, cancelAt := range []string{"sleep", "slot reader"} {
		t.Run(cancelAt, func(t *testing.T) {
			db := openDeliveryEngineStore(t, t.Context())
			engine, clock := newSlotWaitEngine(t, db, "/queue", time.Time{})
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			stop := func() error {
				got := readDeliveryQueue(t, t.Context(), db, "/queue").Items[0]
				if !strings.HasPrefix(got.Blocker, WaitRunSlot+":") {
					t.Fatalf("no wait recorded: %+v", got)
				}
				cancel()
				return ctx.Err()
			}
			engine.sleeper = slotSleepFunc(func(ctx context.Context, d time.Duration) error {
				if cancelAt == "sleep" {
					return stop()
				}
				return clock.Sleep(ctx, d)
			})
			engine.runSlots = slotReaderFunc(func(context.Context) (RunSlots, error) { return RunSlots{}, stop() })
			if _, err := engine.Run(ctx, "/queue"); !errors.Is(err, context.Canceled) {
				t.Fatalf("error=%v", err)
			}
			got := readDeliveryQueue(t, t.Context(), db, "/queue").Items[0]
			if got.Stage != store.DeliveryStageRunning || got.Blocker != "" {
				t.Fatalf("cancelled item=%+v", got)
			}
		})
	}
}

func TestTwoQueuesWaitingForOneSlotStartOneRun(t *testing.T) {
	t.Parallel()
	db := openDeliveryEngineStore(t, t.Context())
	holder, err := db.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: "/holder", SpecSlug: "holder", LocalBranch: "feat/holder", HeadSHA: "holder", Agent: "codex", MaxActiveImplementRuns: 1})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	var wg sync.WaitGroup
	t.Cleanup(func() { cancel(); wg.Wait() })
	type pollEvent struct {
		root    string
		attempt int
	}
	waiting := make(chan pollEvent, 8)
	created := make(chan string, 2)
	activeCounts := make(chan int, 8)
	ended := make(chan error, 2)
	release := make(chan struct{})
	// Both readers take a free observation before either runner can retry.
	observed := make(chan struct{}, 2)
	admit := make(chan struct{})
	reader := slotReaderFunc(func(ctx context.Context) (RunSlots, error) {
		runs, err := db.ActiveImplementRuns(ctx)
		if err != nil {
			return RunSlots{}, err
		}
		slots := RunSlots{MaxActive: 1}
		for _, run := range runs {
			slots.Holders = append(slots.Holders, RunSlotHolder{RunID: run.ID, Repository: run.GitRoot, SpecSlug: run.SpecSlug})
		}
		if slots.Free() {
			observed <- struct{}{}
			select {
			case <-admit:
			case <-ctx.Done():
				return RunSlots{}, ctx.Err()
			}
		}
		return slots, nil
	})
	for _, root := range []string{"/queue-a", "/queue-b"} {
		engine, _ := newSlotWaitEngine(t, db, root, time.Time{})
		engine.runSlots = reader
		sleeps := 0
		engine.sleeper = slotSleepFunc(func(ctx context.Context, _ time.Duration) error {
			sleeps++
			waiting <- pollEvent{root: root, attempt: sleeps}
			if sleeps == 1 {
				select {
				case <-release:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			<-ctx.Done()
			return ctx.Err()
		})
		engine.runner = slotRunnerFunc(func(ctx context.Context, work, slug string) (RunResult, error) {
			run, err := db.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: root, SpecSlug: slug, LocalBranch: "feat/waiting", HeadSHA: "head", Agent: "codex", MaxActiveImplementRuns: 1})
			active, readErr := db.ActiveImplementRuns(ctx)
			if readErr != nil {
				return RunResult{}, readErr
			}
			activeCounts <- len(active)
			if err != nil {
				var ceiling store.ActiveImplementRunCeilingError
				if errors.As(err, &ceiling) {
					slots := RunSlots{MaxActive: ceiling.MaxActive}
					for _, held := range ceiling.Holders {
						slots.Holders = append(slots.Holders, RunSlotHolder{RunID: held.ID, Repository: held.GitRoot, SpecSlug: held.SpecSlug})
					}
					return RunResult{}, &RunSlotRefusal{Err: err, Slots: slots}
				}
				return RunResult{}, err
			}
			created <- run.ID
			return RunResult{RunID: run.ID, Outcome: RunOutcomeUnresolved}, nil
		})
		wg.Go(func() { _, err := engine.Run(ctx, root); ended <- err })
	}
	for range 2 {
		event := testwait.Until(t, "initial slot wait", waiting, ended)
		if event.attempt != 1 {
			t.Fatalf("event=%+v", event)
		}
	}
	if _, err := db.CompleteRun(t.Context(), holder.ID, store.StateClean); err != nil {
		t.Fatal(err)
	}
	close(release)
	for range 2 {
		testwait.Until(t, "free slot observations", observed, ended)
	}
	close(admit)
	runID := testwait.Until(t, "one admitted Run", created, (<-chan error)(nil))
	loser := testwait.Until(t, "losing queue waits again", waiting, (<-chan error)(nil))
	if loser.attempt != 2 {
		t.Fatalf("loser event=%+v", loser)
	}
	got := readDeliveryQueue(t, t.Context(), db, loser.root).Items[0]
	if got.Stage != store.DeliveryStageRunning || !strings.HasPrefix(got.Blocker, WaitRunSlot+":") || !strings.Contains(got.Blocker, runID) {
		t.Fatalf("losing item=%+v", got)
	}
	runs, err := db.ActiveImplementRuns(t.Context())
	if err != nil || len(runs) != 1 || runs[0].ID != runID {
		t.Fatalf("active Runs=%+v err=%v", runs, err)
	}
	cancel()
	wg.Wait()
	// The two initial refusals and simultaneous free observations forced both
	// queues through CreateRun; only one created a Run, and it stays active.
	if len(created) != 0 {
		t.Fatalf("extra Runs created: %d", len(created))
	}
	if len(activeCounts) != 4 {
		t.Fatalf("admission attempts = %d, want two refusals and two retries", len(activeCounts))
	}
	for range 4 {
		if count := <-activeCounts; count > 1 {
			t.Fatalf("Active Implement Runs exceeded ceiling: %d", count)
		}
	}
}
