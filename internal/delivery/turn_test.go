// Suite: Delivery Queue turns.
// Invariant: each unfinished item has an honest turn, and recording never decides delivery.
// Boundary IN: TurnLines and Engine.Run over a real temporary queue store.
// Boundary OUT: synchronized workflow fakes and an in-memory TurnRecorder.
package delivery

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/store"
	"roundfix/internal/testwait"
)

func TestTurnLinesNameEachItemsTurn(t *testing.T) {
	t.Parallel()
	stages := []store.DeliveryStage{"queued", "running", "gating", "queued", "reviewing", "reviewing", "reviewing", "queued", "archiving", "parked"}
	states := []TurnState{TurnWait, TurnAdvancing, TurnWait, TurnWait, TurnReviewing, TurnReviewSlotWait, TurnReviewSlotWait, TurnPrerequisiteWait, TurnReviewing}
	var queue store.DeliveryQueue
	turn := Turn{ReviewCapacity: 2}
	for index, stage := range stages {
		slug := string(rune('a' + index))
		queue.Items = append(queue.Items, store.DeliveryQueueItem{SpecSlug: slug, Stage: stage, Position: index})
		if index < len(states) {
			turn.Items = append(turn.Items, ItemTurn{SpecSlug: slug, Stage: stage, State: states[index]})
		}
	}
	turn.Items[7].Prerequisites = []string{"0090-x", "0091-y"}
	turn.Items[8].Stage = store.DeliveryStageReviewing
	want := []string{
		"Turn: a waiting for the owner's turn (3 of 3 waiting)",
		"Turn: b advancing at running",
		"Turn: c waiting for the owner's turn (1 of 3 waiting)",
		"Turn: d waiting for the owner's turn (2 of 3 waiting)",
		"Turn: e reviewing; 1 of 2 review slot(s) held",
		"Turn: f waiting for a review slot (1 of 2 waiting); 1 of 2 review slot(s) held",
		"Turn: g waiting for a review slot (2 of 2 waiting); 1 of 2 review slot(s) held",
		"Turn: h waiting for prerequisites: 0090-x, 0091-y",
		"Turn: i waiting for the owner to read it",
	}
	if got := TurnLines(queue, turn); !reflect.DeepEqual(got, want) {
		t.Fatalf("lines = %q, want %q", got, want)
	}
	t.Run("nothing advancing", func(t *testing.T) {
		turn.Items = slices.Delete(slices.Clone(turn.Items), 1, 2)
		want := slices.Clone(want)
		want[0] = "Turn: a waiting for the owner's turn (1 of 3 waiting)"
		want[1] = "Turn: b waiting for the owner to read it"
		want[2] = "Turn: c waiting for the owner's turn (2 of 3 waiting)"
		want[3] = "Turn: d waiting for the owner's turn (3 of 3 waiting)"
		if got := TurnLines(queue, turn); !reflect.DeepEqual(got, want) {
			t.Fatalf("lines = %q, want %q", got, want)
		}
	})
	t.Run("advancing stage moves", func(t *testing.T) {
		turn := Turn{Items: []ItemTurn{{SpecSlug: "b", Stage: store.DeliveryStageQueued, State: TurnAdvancing}}}
		if got := TurnLines(queue, turn)[1]; got != "Turn: b advancing at running" {
			t.Fatal(got)
		}
	})
}

func TestTurnLinesNameTheChecksAnItemWaitsFor(t *testing.T) {
	t.Parallel()
	queue := store.DeliveryQueue{Items: []store.DeliveryQueueItem{
		{SpecSlug: "0149-example", Stage: store.DeliveryStageChecking, Position: 0},
		{SpecSlug: "0150-other", Stage: store.DeliveryStageQueued, Position: 1},
	}}
	turn := Turn{Items: []ItemTurn{
		{SpecSlug: "0149-example", Stage: store.DeliveryStageChecking, State: TurnAdvancing, Checks: []string{
			"CI | Backend / Lint, Type Check, Test & Build",
			".github/workflows/ci-frontends.yml (no run yet)",
		}},
		{SpecSlug: "0150-other", Stage: store.DeliveryStageQueued, State: TurnWait},
	}}
	want := []string{
		"Turn: 0149-example advancing at checking; waiting for 2 check(s): CI | Backend / Lint, Type Check, Test & Build; .github/workflows/ci-frontends.yml (no run yet)",
		"Turn: 0150-other waiting for the owner's turn (1 of 1 waiting)",
	}
	if got := TurnLines(queue, turn); !slices.Equal(got, want) {
		t.Fatalf("lines = %q, want %q", got, want)
	}
	t.Run("no checks", func(t *testing.T) {
		turn.Items[0].Checks = nil
		want[0] = "Turn: 0149-example advancing at checking"
		if got := TurnLines(queue, turn); !slices.Equal(got, want) {
			t.Fatalf("lines = %q, want %q", got, want)
		}
	})
}

func TestTheOwnerRecordsTheChecksAnItemWaitsFor(t *testing.T) {
	t.Parallel()
	engine, runStore, item, _, boundary := newCheckWaitFixture(t, time.Minute)
	var mu sync.Mutex
	var turns []Turn
	engine.turnRecorder = turnRecorderFunc(func(_ context.Context, turn Turn) error {
		mu.Lock()
		defer mu.Unlock()
		turns = append(turns, turn)
		return nil
	})
	polls := 0
	boundary.currentHeadChecks = func(context.Context, string) (CheckReport, error) {
		polls++
		if polls <= 2 {
			return recoveryReport("pending"), nil
		}
		return recoveryReport("pass"), nil
	}
	engine.advancingTurn(t.Context(), "/repo", item.SpecSlug)
	if err := engine.checkCandidate(t.Context(), "/repo", item); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	found := false
	for _, turn := range turns {
		for _, entry := range turn.Items {
			if entry.SpecSlug == item.SpecSlug && entry.State == TurnAdvancing && slices.Equal(entry.Checks, []string{"CI / verify"}) {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("no advancing turn with pending checks: %+v", turns)
	}
	last := turns[len(turns)-1]
	if len(last.Items) != 1 || last.Items[0].SpecSlug != item.SpecSlug || len(last.Items[0].Checks) != 0 {
		t.Fatalf("last turn still has checks or lost the item: %+v", last)
	}
	stored := readDeliveryQueue(t, t.Context(), runStore, "/repo").Items[0]
	if polls != 3 || stored.Stage != store.DeliveryStageMerging || last.Items[0].Stage != store.DeliveryStageMerging {
		t.Fatalf("polls=%d item=%+v last=%+v", polls, stored, last)
	}
}

func TestTheOwnerRecordsEachItemsTurn(t *testing.T) {
	t.Parallel()
	fixture := newReviewSlotFixture(t, "implementing", "reviewed", "waiting", "later")
	fixture.seed(t, 0, store.DeliveryStageRunning)
	fixture.seed(t, 1, store.DeliveryStageReviewing)
	fixture.seed(t, 2, store.DeliveryStageReviewing)
	recorded := make(chan struct{})
	recorder := &keptTurns{matches: func(turn Turn) bool {
		return turn.ReviewCapacity == 1 && turnHas(turn, "implementing", TurnAdvancing) && turnHas(turn, "reviewed", TurnReviewing) && turnHas(turn, "waiting", TurnReviewSlotWait) && turnHas(turn, "later", TurnWait)
	}, matched: recorded}
	fixture.engine = newTurnTestEngine(fixture, recorder)
	// Run synchronously: the RunSpec fake invokes testwait on the test goroutine.
	fixture.workflow.run = func(ctx context.Context, slug string) error {
		if slug == "implementing" {
			testwait.Until(t, "recorded concurrent turns", recorded, ctx.Done())
		}
		return ctx.Err()
	}
	fixture.workflow.review = func(ctx context.Context, slug string) (ReviewOutcome, error) {
		if slug == "reviewed" {
			select {
			case <-recorded:
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
		return ReviewOutcomeReviewed, nil
	}
	if _, err := fixture.engine.Run(t.Context(), fixture.root); err != nil {
		t.Fatal(err)
	}
	fixture.assertMerged(t)
	recorder.mu.Lock()
	if len(recorder.turns) == 0 {
		t.Error("no turns recorded")
	}
	last := recorder.turns[len(recorder.turns)-1]
	recorder.mu.Unlock()
	if len(last.Items) != 0 {
		t.Errorf("last turn = %+v, want no unfinished items", last)
	}

	t.Run("prerequisite wait", func(t *testing.T) {
		t.Parallel()
		fixture := newReviewSlotFixture(t, "blocker", "dependent")
		fixture.seed(t, 0, store.DeliveryStageReviewing)
		recorded := make(chan struct{})
		recorder := &keptTurns{matched: recorded, matches: func(turn Turn) bool {
			for _, item := range turn.Items {
				if item.SpecSlug == "dependent" && item.State == TurnPrerequisiteWait && slices.Equal(item.Prerequisites, []string{"blocker"}) {
					return true
				}
			}
			return false
		}}
		fixture.engine = newTurnTestEngine(fixture, recorder)
		fixture.engine.prerequisites = prerequisiteReaderFunc(func(ctx context.Context, _, slug string) ([]string, error) {
			if slug == "dependent" {
				queue, _, err := fixture.store.DeliveryQueue(ctx, fixture.root)
				if err != nil {
					return nil, err
				}
				if queue.Items[0].Stage != store.DeliveryStageMerged {
					return []string{"blocker"}, nil
				}
			}
			return nil, nil
		})
		release := make(chan struct{})
		fixture.workflow.review = func(ctx context.Context, slug string) (ReviewOutcome, error) {
			if slug == "blocker" {
				select {
				case <-release:
				case <-ctx.Done():
					return "", ctx.Err()
				}
			}
			return ReviewOutcomeReviewed, nil
		}
		done := fixture.start(t)
		testwait.Until(t, "dependent prerequisite turn before review release", recorded, done)
		close(release)
		fixture.finish(t, done)
		// Existing prerequisite deferral deliberately leaves this item for the next Run.
		fixture.finish(t, fixture.start(t))
		fixture.assertMerged(t)
	})
	t.Run("recorder failure", func(t *testing.T) {
		t.Parallel()
		fixture := newReviewSlotFixture(t, "first", "second")
		fixture.engine = newTurnTestEngine(fixture, turnRecorderFunc(func(context.Context, Turn) error { return errors.New("disk full") }))
		fixture.finish(t, fixture.start(t))
		fixture.assertMerged(t)
		if !strings.Contains(fixture.log.String(), "roundfix: warning: Delivery Queue turn record: disk full\n") {
			t.Fatal(fixture.log.String())
		}
	})
}

type turnRecorderFunc func(context.Context, Turn) error

func (record turnRecorderFunc) RecordTurn(ctx context.Context, turn Turn) error {
	return record(ctx, turn)
}

type keptTurns struct {
	mu      sync.Mutex
	turns   []Turn
	matches func(Turn) bool
	matched chan struct{}
	once    sync.Once
}

func (recorder *keptTurns) RecordTurn(_ context.Context, turn Turn) error {
	recorder.mu.Lock()
	defer recorder.mu.Unlock()
	recorder.turns = append(recorder.turns, turn)
	if recorder.matches(turn) {
		recorder.once.Do(func() { close(recorder.matched) })
	}
	return nil
}
func turnHas(turn Turn, slug string, state TurnState) bool {
	for _, item := range turn.Items {
		if item.SpecSlug == slug {
			return item.State == state
		}
	}
	return false
}

func newTurnTestEngine(fixture *reviewSlotFixture, recorder TurnRecorder) *Engine {
	wf := fixture.workflow
	return NewEngine(fixture.store, EngineDependencies{
		Workspace: wf, Runner: wf, Reviewer: wf, Archiver: wf, Gate: wf,
		Authorizer: wf, Publication: wf, PullRequests: fixture.boundary,
		Revalidator: cleanTestRevalidator{}, Log: fixture.log,
		ReviewCapacity: 1, TurnPollInterval: 10 * time.Millisecond, TurnRecorder: recorder,
	})
}
