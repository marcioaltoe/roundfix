// Suite: Queue Owner review slots.
// Invariant: bounded reviews run beside implementation, preserve queued start order and stop with their Run.
// Boundary IN: Engine.Run and the real SQLite Delivery Queue store, with provisioned temporary item worktrees.
// Boundary OUT: workspace, implementation, review and publication use synchronized fakes and event signals.
package delivery

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/store"
	"roundfix/internal/testwait"
)

func TestAReviewHandedToTheQueueStartsWhileAnotherItemImplements(t *testing.T) {
	t.Parallel()
	fixture := newReviewSlotFixture(t, "implementing", "retried")
	fixture.seed(t, 0, store.DeliveryStageRunning)
	fixture.seed(t, 1, store.DeliveryStageReviewing)
	if got := fixture.engine.turnPollInterval; got != 30*time.Second {
		t.Fatalf("default poll interval = %s, want 30s", got)
	}
	fixture.engine = fixture.newEngine(2, 10*time.Millisecond)
	retriedStarted := make(chan struct{})
	implementingStarted := make(chan struct{})
	release := make(chan struct{})
	fixture.workflow.run = func(ctx context.Context, slug string) error {
		if slug == "implementing" {
			select {
			case <-retriedStarted:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	}
	fixture.workflow.review = func(ctx context.Context, slug string) (ReviewOutcome, error) {
		if slug == "implementing" {
			close(implementingStarted)
			return ReviewOutcomeReviewed, nil
		}
		close(retriedStarted)
		select {
		case <-release:
			return ReviewOutcomeReviewed, nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	// The next pass archives implementing only after the current pass has
	// visited retried while its review still holds the worktree.
	fixture.workflow.archive = func(slug string) {
		if slug == "implementing" {
			close(release)
		}
	}
	done := fixture.start(t)
	testwait.Until(t, "retried review while implementation blocks", retriedStarted, done)
	testwait.Until(t, "implementing review", implementingStarted, done)
	fixture.finish(t, done)
	fixture.assertMerged(t)
	fixture.workflow.mu.Lock()
	defer fixture.workflow.mu.Unlock()
	if fixture.workflow.used["retried"] != 0 {
		t.Errorf("UseItemBranch for retried = %d, want 0", fixture.workflow.used["retried"])
	}
	for _, slug := range fixture.slugs {
		if fixture.workflow.reviews[slug] != 1 {
			t.Errorf("Review(%s) count = %d, want 1", slug, fixture.workflow.reviews[slug])
		}
	}
}

func TestReviewCapacityBoundsConcurrentReviews(t *testing.T) {
	t.Parallel()
	runReviewCapacityCase(t, 2)
	t.Run("capacity below one", func(t *testing.T) {
		t.Parallel()
		runReviewCapacityCase(t, 0)
	})
}

func runReviewCapacityCase(t *testing.T, capacity int) {
	t.Helper()
	fixture := newReviewSlotFixture(t, "first", "second", "third")
	for index := range fixture.slugs {
		fixture.seed(t, index, store.DeliveryStageReviewing)
	}
	fixture.engine = fixture.newEngine(capacity, 0)
	started := make(chan string, 3)
	release := make(chan struct{})
	fixture.workflow.review = func(ctx context.Context, slug string) (ReviewOutcome, error) {
		started <- slug
		select {
		case <-release:
			return ReviewOutcomeReviewed, nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	done := fixture.start(t)
	limit := max(1, capacity)
	for range limit {
		testwait.Until(t, "initial reviews", started, done)
	}
	// Reaching the capacity log proves the pass admitted the waiting item
	// before the first reviews are released.
	testwait.Poll(t, "third review waiting for capacity", done, func() (bool, string) {
		log := fixture.log.String()
		return strings.Contains(log, "item third: waiting for a review slot"), log
	})
	close(release)
	fixture.finish(t, done)
	fixture.assertMerged(t)
	fixture.workflow.mu.Lock()
	defer fixture.workflow.mu.Unlock()
	if fixture.workflow.largest != limit {
		t.Errorf("largest concurrent reviews = %d, want %d", fixture.workflow.largest, limit)
	}
	events := fixture.workflow.events
	thirdStart := eventIndex(events, "third:review")
	firstEnd := eventIndex(events, "first:review-end")
	secondEnd := eventIndex(events, "second:review-end")
	if thirdStart < 0 || (thirdStart < firstEnd && thirdStart < secondEnd) {
		t.Errorf("third review started before capacity freed: %v", events)
	}
	want := fmt.Sprintf("roundfix: Delivery Queue item third: waiting for a review slot; %d of %d review slot(s) held\n", limit, limit)
	if count := strings.Count(fixture.log.String(), want); count != 1 {
		t.Errorf("wait line count = %d, want 1; log=%s", count, fixture.log.String())
	}
}

func TestAQueuedItemWaitsForAnEarlierItemInReview(t *testing.T) {
	t.Parallel()
	fixture := newReviewSlotFixture(t, "first", "second")
	fixture.engine.reviewCapacity = 1
	started := make(chan struct{})
	release := make(chan struct{})
	fixture.workflow.review = func(ctx context.Context, slug string) (ReviewOutcome, error) {
		if slug == "first" {
			close(started)
			select {
			case <-release:
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
		return ReviewOutcomeReviewed, nil
	}
	secondChecked := make(chan struct{}, 1)
	fixture.engine.prerequisites = prerequisiteReaderFunc(func(_ context.Context, _, slug string) ([]string, error) {
		if slug == "second" {
			select {
			case secondChecked <- struct{}{}:
			default:
			}
		}
		return nil, nil
	})
	done := fixture.start(t)
	testwait.Until(t, "first review", started, done)
	testwait.Until(t, "second queued item visited", secondChecked, done)
	fixture.workflow.mu.Lock()
	created := fixture.workflow.created["second"]
	fixture.workflow.mu.Unlock()
	if created != 0 {
		t.Errorf("second branch created during first review")
	}
	close(release)
	fixture.finish(t, done)
	fixture.assertMerged(t)
	var want []string
	for _, slug := range fixture.slugs {
		for _, event := range []string{"create", "run", "review", "review-end", "archive", "gate", "publication", "push", "pull-request", "checks", "merge"} {
			want = append(want, slug+":"+event)
		}
	}
	if got := fixture.workflow.eventSnapshot(); !reflect.DeepEqual(got, want) {
		t.Errorf("event order = %v, want inline capacity-one order %v", got, want)
	}
	t.Run("an earlier blocker ends during the pass", testEarlierReviewEndsDuringPass)
}

func testEarlierReviewEndsDuringPass(t *testing.T) {
	t.Helper()
	fixture := newReviewSlotFixture(t, "first", "second", "third")
	signal := make(chan struct{})
	fixture.workflow.review = func(ctx context.Context, slug string) (ReviewOutcome, error) {
		if slug == "first" {
			select {
			case <-signal:
				return ReviewOutcomeFindings, nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
		return ReviewOutcomeReviewed, nil
	}
	var once sync.Once
	// Run synchronously so Poll, like every testwait helper, stays on the
	// test goroutine even though the engine calls this fake reader.
	fixture.engine.prerequisites = prerequisiteReaderFunc(func(ctx context.Context, _, slug string) ([]string, error) {
		if slug == "third" {
			once.Do(func() {
				close(signal)
				testwait.Poll(t, "first review parks during this pass", (<-chan error)(nil), func() (bool, string) {
					queue, _, err := fixture.store.DeliveryQueue(ctx, fixture.root)
					if err != nil {
						return false, err.Error()
					}
					return queue.Items[0].Stage == store.DeliveryStageParked, string(queue.Items[0].Stage)
				})
			})
		}
		return nil, nil
	})
	if _, err := fixture.engine.Run(t.Context(), fixture.root); err != nil {
		t.Fatal(err)
	}
	events := fixture.workflow.eventSnapshot()
	if second, third := eventIndex(events, "second:create"), eventIndex(events, "third:create"); second < 0 || third < 0 || second > third {
		t.Errorf("queued start order after earlier park = %v", events)
	}
}

func TestAFailedReviewParksOnlyItsItem(t *testing.T) {
	t.Parallel()
	fixture := newReviewSlotFixture(t, "broken", "next")
	fixture.seed(t, 0, store.DeliveryStageReviewing)
	fixture.workflow.review = func(_ context.Context, slug string) (ReviewOutcome, error) {
		if slug == "broken" {
			return "", errors.New("review exploded")
		}
		return ReviewOutcomeReviewed, nil
	}
	fixture.finish(t, fixture.start(t))
	queue := readDeliveryQueue(t, t.Context(), fixture.store, fixture.root)
	if item := queue.Items[0]; item.Stage != store.DeliveryStageParked || item.Blocker != "delivery-error: run pre-PR review: review exploded" {
		t.Errorf("broken review = %+v", item)
	}
	if queue.Items[1].Stage != store.DeliveryStageMerged {
		t.Errorf("next item = %+v, want merged", queue.Items[1])
	}
}

func TestStoppingTheOwnerLeavesAnItemInReviewAtReviewing(t *testing.T) {
	t.Parallel()
	fixture := newReviewSlotFixture(t, "stopping")
	fixture.seed(t, 0, store.DeliveryStageReviewing)
	started := make(chan struct{})
	returned := make(chan struct{})
	fixture.workflow.review = func(ctx context.Context, _ string) (ReviewOutcome, error) {
		close(started)
		<-ctx.Done()
		defer close(returned)
		return "", ctx.Err()
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := fixture.startContext(t, ctx)
	testwait.Until(t, "review starts", started, done)
	cancel()
	err := testwait.Until(t, "Run stops after its review", done, (<-chan error)(nil))
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Run error = %v, want context.Canceled", err)
	}
	select {
	case <-returned:
	default:
		t.Error("Run returned before its review")
	}
	item := readDeliveryQueue(t, t.Context(), fixture.store, fixture.root).Items[0]
	if item.Stage != store.DeliveryStageReviewing || item.Blocker != "" {
		t.Errorf("cancelled review = %+v, want reviewing without blocker", item)
	}
}

type reviewSlotFixture struct {
	store    *store.Store
	engine   *Engine
	workflow *slotWorkflow
	boundary *slotBoundary
	root     string
	slugs    []string
	log      *slotLog
}

func newReviewSlotFixture(t *testing.T, slugs ...string) *reviewSlotFixture {
	t.Helper()
	root := t.TempDir()
	fixture := &reviewSlotFixture{store: openDeliveryEngineStore(t, t.Context()), root: root, slugs: slugs, log: &slotLog{}}
	fixture.workflow = &slotWorkflow{root: root, created: make(map[string]int), used: make(map[string]int), reviews: make(map[string]int)}
	fixture.boundary = &slotBoundary{fake: newFakeDeliveryBoundary()}
	fixture.boundary.fake.recordEvent = fixture.workflow.record
	if _, err := fixture.store.CreateDeliveryQueue(t.Context(), root, slugs); err != nil {
		t.Fatal(err)
	}
	for _, slug := range slugs {
		if err := os.Mkdir(filepath.Join(root, slug), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	fixture.engine = fixture.newEngine(1, 0)
	return fixture
}

func (fixture *reviewSlotFixture) newEngine(capacity int, pollInterval time.Duration) *Engine {
	wf := fixture.workflow
	return NewEngine(fixture.store, EngineDependencies{Workspace: wf, Runner: wf, Reviewer: wf, Archiver: wf, Gate: wf, Authorizer: wf, Publication: wf, PullRequests: fixture.boundary, Revalidator: cleanTestRevalidator{}, ReviewCapacity: capacity, TurnPollInterval: pollInterval, Log: fixture.log})
}

func (fixture *reviewSlotFixture) seed(t *testing.T, index int, stage store.DeliveryStage) {
	t.Helper()
	item := readDeliveryQueue(t, t.Context(), fixture.store, fixture.root).Items[index]
	item.Stage, item.Branch = stage, "roundfix/deliver-"+item.SpecSlug
	item.Worktree, item.WorktreeProvisioned = filepath.Join(fixture.root, item.SpecSlug), true
	recordDeliveryItemWorkspace(t, t.Context(), fixture.store, fixture.root, &item)
	if err := fixture.store.SetDeliveryQueueItemWorktreeProvisioned(t.Context(), fixture.root, item.SpecSlug, true); err != nil {
		t.Fatal(err)
	}
	item.CandidateCommits = []string{"reviewed-" + item.SpecSlug}
	if err := fixture.store.UpdateDeliveryQueueItem(t.Context(), fixture.root, item); err != nil {
		t.Fatal(err)
	}
}

func (fixture *reviewSlotFixture) start(t *testing.T) <-chan error {
	t.Helper()
	return fixture.startContext(t, t.Context())
}

func (fixture *reviewSlotFixture) startContext(t *testing.T, parent context.Context) <-chan error {
	t.Helper()
	ctx, cancel := context.WithCancel(parent)
	done := make(chan error, 1)
	ended := make(chan struct{})
	go func() {
		defer close(ended)
		_, err := fixture.engine.Run(ctx, fixture.root)
		done <- err
	}()
	t.Cleanup(func() {
		cancel()
		testwait.Until(t, "review fixture stops", ended, (<-chan error)(nil))
	})
	return done
}

func (fixture *reviewSlotFixture) finish(t *testing.T, done <-chan error) {
	t.Helper()
	if err := testwait.Until(t, "Delivery Queue Run", done, (<-chan error)(nil)); err != nil {
		t.Fatal(err)
	}
}

func (fixture *reviewSlotFixture) assertMerged(t *testing.T) {
	t.Helper()
	for _, item := range readDeliveryQueue(t, t.Context(), fixture.store, fixture.root).Items {
		if item.Stage != store.DeliveryStageMerged {
			t.Errorf("item %s = %+v, want merged", item.SpecSlug, item)
		}
	}
}

func eventIndex(events []string, event string) int {
	for index, value := range events {
		if value == event {
			return index
		}
	}
	return -1
}

type slotWorkflow struct {
	mu                     sync.Mutex
	root                   string
	created, used, reviews map[string]int
	events                 []string
	active, largest        int
	run                    func(context.Context, string) error
	review                 func(context.Context, string) (ReviewOutcome, error)
	archive                func(string)
}

func (wf *slotWorkflow) record(slug, event string) {
	wf.mu.Lock()
	defer wf.mu.Unlock()
	wf.events = append(wf.events, slug+":"+event)
}

func (wf *slotWorkflow) eventSnapshot() []string {
	wf.mu.Lock()
	defer wf.mu.Unlock()
	return append([]string(nil), wf.events...)
}

func (wf *slotWorkflow) CreateItemBranch(_ context.Context, _, slug string) (string, string, error) {
	wf.mu.Lock()
	defer wf.mu.Unlock()
	wf.created[slug]++
	wf.events = append(wf.events, slug+":create")
	return "roundfix/deliver-" + slug, filepath.Join(wf.root, slug), nil
}
func (wf *slotWorkflow) UseItemBranch(_ context.Context, _, slug, _, worktree string, _ bool) (string, error) {
	wf.mu.Lock()
	defer wf.mu.Unlock()
	wf.used[slug]++
	return worktree, nil
}
func (wf *slotWorkflow) ReleaseMergedRuns(context.Context, string, store.DeliveryQueueItem) error {
	return nil
}
func (wf *slotWorkflow) RemoveItemBranch(context.Context, string, string, string) error { return nil }
func (wf *slotWorkflow) RunSpec(ctx context.Context, _, slug string) (RunResult, error) {
	wf.record(slug, "run")
	if wf.run != nil {
		if err := wf.run(ctx, slug); err != nil {
			return RunResult{}, err
		}
	}
	return RunResult{Outcome: RunOutcomeClean, CandidateCommits: []string{"reviewed-" + slug}}, nil
}
func (wf *slotWorkflow) ReviewPolicy(context.Context, string, string) (ReviewPolicy, error) {
	return ReviewPolicyEnabled, nil
}
func (wf *slotWorkflow) RecordReviewOmission(context.Context, string, string, string) error {
	return nil
}
func (wf *slotWorkflow) Review(ctx context.Context, _, slug, head string) (ReviewResult, error) {
	wf.mu.Lock()
	wf.reviews[slug]++
	wf.active++
	wf.largest = max(wf.largest, wf.active)
	wf.events = append(wf.events, slug+":review")
	wf.mu.Unlock()
	outcome := ReviewOutcomeReviewed
	var err error
	if wf.review != nil {
		outcome, err = wf.review(ctx, slug)
	}
	wf.mu.Lock()
	wf.active--
	wf.events = append(wf.events, slug+":review-end")
	wf.mu.Unlock()
	return ReviewResult{Outcome: outcome, Head: head}, err
}
func (wf *slotWorkflow) Archive(_ context.Context, _, slug, head string) (ArchiveResult, error) {
	wf.record(slug, "archive")
	if wf.archive != nil {
		wf.archive(slug)
	}
	return ArchiveResult{Parent: head, Head: "archived-" + slug, ExactSpecMove: true}, nil
}
func (wf *slotWorkflow) Gate(_ context.Context, _, slug, _ string) (GateResult, error) {
	wf.record(slug, "gate")
	return GateResult{Passed: true}, nil
}
func (wf *slotWorkflow) Authorization(context.Context, string, string) (Authorization, error) {
	return Authorization{Operations: []string{"push", "pull_request", "merge"}}, nil
}
func (wf *slotWorkflow) Publication(_ context.Context, _, slug, branch string) (Publication, error) {
	wf.record(slug, "publication")
	return Publication{Remote: "origin", HeadBranch: branch, BaseBranch: "main", Title: slug}, nil
}

type slotLog struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (log *slotLog) Write(p []byte) (int, error) {
	log.mu.Lock()
	defer log.mu.Unlock()
	return log.buffer.Write(p)
}
func (log *slotLog) String() string { log.mu.Lock(); defer log.mu.Unlock(); return log.buffer.String() }

type slotBoundary struct {
	mu   sync.Mutex
	fake *fakeDeliveryBoundary
}

func (b *slotBoundary) WithWorkDir(dir string) PullRequestBoundary {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fake.WithWorkDir(dir)
	return b
}
func (b *slotBoundary) RemoteBranchHead(ctx context.Context, remote, branch string) (RemoteHead, bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.fake.RemoteBranchHead(ctx, remote, branch)
}
func (b *slotBoundary) PushBranch(ctx context.Context, remote, branch, head string) (RemoteHead, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.fake.PushBranch(ctx, remote, branch, head)
}
func (b *slotBoundary) FindOrCreatePullRequest(ctx context.Context, req PullRequestRequest) (PullRequestResult, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.fake.FindOrCreatePullRequest(ctx, req)
}
func (b *slotBoundary) CurrentHeadChecks(ctx context.Context, number string) (CheckReport, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.fake.CurrentHeadChecks(ctx, number)
}
func (b *slotBoundary) MergePullRequest(ctx context.Context, number, head string) (MergeResult, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.fake.MergePullRequest(ctx, number, head)
}
