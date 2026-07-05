package watch

import (
	"context"
	"errors"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestRunWaitsFetchesResolvesToClean(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	status := &fakeStatusSource{
		statuses: []Status{
			{State: StatusPending, Detail: "CodeRabbit is queued"},
			{State: StatusSettled, Detail: "CodeRabbit is settled"},
		},
	}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 2}}}
	resolver := &fakeResolver{results: []ResolveResult{{Remaining: 0, Progress: true}}}

	result, err := Run(context.Background(), validRequest(), Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("expected Clean, got %q", result.Outcome)
	}
	if status.calls != 2 {
		t.Fatalf("expected polling until settled, got %d calls", status.calls)
	}
	if !sleeper.saw(validRequest().PollInterval) {
		t.Fatalf("expected poll interval sleep, got %#v", sleeper.sleeps)
	}
	if !sleeper.saw(validRequest().QuietPeriod) {
		t.Fatalf("expected quiet period sleep, got %#v", sleeper.sleeps)
	}
	if fetcher.calls != 1 {
		t.Fatalf("expected one fetch, got %d", fetcher.calls)
	}
	if resolver.calls != 1 {
		t.Fatalf("expected one resolve, got %d", resolver.calls)
	}
}

func TestRunSkipsQuietPeriodWhenReviewAlreadySettledAtStart(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
	fetchCalls := 0

	result, err := Run(context.Background(), validRequest(), Dependencies{
		StatusSource: status,
		Fetcher: FetchFunc(func(_ context.Context, round int) (FetchResult, error) {
			fetchCalls++
			assertSleeps(t, sleeper.sleeps)
			return FetchResult{Round: round, Issues: 0}, nil
		}),
		Resolver: &fakeResolver{},
		Clock:    clock,
		Sleeper:  sleeper,
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("expected Clean, got %q", result.Outcome)
	}
	if status.calls != 1 {
		t.Fatalf("expected one status check, got %d", status.calls)
	}
	if fetchCalls != 1 {
		t.Fatalf("expected one fetch, got %d", fetchCalls)
	}
	assertSleeps(t, sleeper.sleeps)
}

func TestRunSleepsBetweenStatusChecksAndKeepsQuietPeriodWhenReviewSettlesDuringRun(t *testing.T) {
	req := validRequest()
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	statusCalls := 0

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: StatusFunc(func(_ context.Context, _ StatusRequest) (Status, error) {
			statusCalls++
			switch statusCalls {
			case 1:
				assertSleeps(t, sleeper.sleeps)
				return Status{State: StatusPending}, nil
			case 2:
				assertSleeps(t, sleeper.sleeps, req.PollInterval)
				return Status{State: StatusPending}, nil
			case 3:
				assertSleeps(t, sleeper.sleeps, req.PollInterval, req.PollInterval)
				return Status{State: StatusSettled}, nil
			default:
				t.Fatalf("unexpected status call %d", statusCalls)
				return Status{}, nil
			}
		}),
		Fetcher:  &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver: &fakeResolver{},
		Clock:    clock,
		Sleeper:  sleeper,
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("expected Clean, got %q", result.Outcome)
	}
	if statusCalls != 3 {
		t.Fatalf("expected three status checks, got %d", statusCalls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval, req.PollInterval, req.QuietPeriod)
}

func TestRunTimesOutAndOffersManualReviewTrigger(t *testing.T) {
	req := validRequest()
	req.ReviewTimeout = 2 * time.Second
	req.PollInterval = time.Second
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	status := &fakeStatusSource{
		statuses: []Status{
			{State: StatusPending},
			{State: StatusPending},
			{State: StatusPending},
		},
	}
	fetcher := &fakeFetcher{}
	resolver := &fakeResolver{}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if err != nil {
		t.Fatalf("watch timeout should be terminal result, got %v", err)
	}
	if result.Outcome != store.StateTimedOut {
		t.Fatalf("expected TimedOut, got %q", result.Outcome)
	}
	if result.ManualReviewCommand != "@coderabbitai review" {
		t.Fatalf("expected manual review trigger guidance, got %q", result.ManualReviewCommand)
	}
	if fetcher.calls != 0 || resolver.calls != 0 {
		t.Fatalf("timeout must not fetch or resolve, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
}

func TestRunStopsAtMaxRoundsReached(t *testing.T) {
	req := validRequest()
	req.MaxRounds = 2
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	status := &fakeStatusSource{
		statuses: []Status{{State: StatusSettled}, {State: StatusSettled}},
	}
	fetcher := &fakeFetcher{
		results: []FetchResult{
			{Round: 1, Issues: 1},
			{Round: 2, Issues: 1},
		},
	}
	resolver := &fakeResolver{
		results: []ResolveResult{
			{Remaining: 1, Progress: true},
			{Remaining: 2, Progress: true},
		},
	}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if err != nil {
		t.Fatalf("watch max rounds should be terminal result, got %v", err)
	}
	if result.Outcome != store.StateMaxRoundsReached {
		t.Fatalf("expected MaxRoundsReached, got %q", result.Outcome)
	}
	if result.Remaining != 2 {
		t.Fatalf("expected remaining issues to be reported, got %d", result.Remaining)
	}
	if result.Rounds != 2 {
		t.Fatalf("expected 2 rounds, got %d", result.Rounds)
	}
}

func TestRunStopsWhenBudgetExceeded(t *testing.T) {
	req := validRequest()
	req.BudgetEnabled = true
	req.MaxRunDuration = 2 * time.Second
	req.PollInterval = time.Second
	req.QuietPeriod = 2 * time.Second
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	status := &fakeStatusSource{
		statuses: []Status{
			{State: StatusPending},
			{State: StatusSettled},
		},
	}
	fetcher := &fakeFetcher{}
	resolver := &fakeResolver{}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if err != nil {
		t.Fatalf("watch budget should be terminal result, got %v", err)
	}
	if result.Outcome != store.StateBudgetExceeded {
		t.Fatalf("expected BudgetExceeded, got %q", result.Outcome)
	}
	if fetcher.calls != 0 || resolver.calls != 0 {
		t.Fatalf("budget exceeded after quiet period must not fetch or resolve, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
}

func TestRunReturnsUnresolvedWhenResolveMakesNoProgress(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 1}}}
	resolver := &fakeResolver{results: []ResolveResult{{Remaining: 1, Progress: false}}}

	result, err := Run(context.Background(), validRequest(), Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		Clock:        clock,
		Sleeper:      &fakeSleeper{clock: clock},
	})

	if err != nil {
		t.Fatalf("a no-progress Round is a terminal outcome, not an error, got %v", err)
	}
	if result.Outcome != store.StateUnresolved {
		t.Fatalf("expected Unresolved, got %q", result.Outcome)
	}
	if result.Remaining != 1 {
		t.Fatalf("expected remaining Unresolved Review Issue count, got %d", result.Remaining)
	}
	if fetcher.calls != 1 || resolver.calls != 1 {
		t.Fatalf("expected exactly one Round before stopping, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
}

func TestRunConfirmsSuccessAndMissingHeadCheckAsClean(t *testing.T) {
	tests := []struct {
		name        string
		checkState  HeadCheckState
		wantMissing bool
	}{
		{name: "success", checkState: CheckSuccess},
		{name: "missing", checkState: CheckMissing, wantMissing: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validRequest()
			clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
			sleeper := &fakeSleeper{clock: clock}
			status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
			fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 1}}}
			resolver := &fakeResolver{results: []ResolveResult{{Remaining: 0, Progress: true}}}
			check := &fakeCheckSource{states: []HeadCheckState{tt.checkState}}

			result, err := Run(context.Background(), req, Dependencies{
				StatusSource: status,
				Fetcher:      fetcher,
				Resolver:     resolver,
				CheckSource:  check,
				Clock:        clock,
				Sleeper:      sleeper,
			})

			if err != nil {
				t.Fatalf("watch run: %v", err)
			}
			if result.Outcome != store.StateClean {
				t.Fatalf("expected Clean, got %q", result.Outcome)
			}
			if result.Rounds != 1 {
				t.Fatalf("expected 1 Round, got %d", result.Rounds)
			}
			if result.CheckMissing != tt.wantMissing {
				t.Fatalf("expected CheckMissing %t, got %t", tt.wantMissing, result.CheckMissing)
			}
			if check.calls != 1 {
				t.Fatalf("expected one check call, got %d", check.calls)
			}
			assertSleeps(t, sleeper.sleeps)
		})
	}
}

func TestRunKeepsPollingWhenLocalQueueIsEmptyUntilHeadCheckSucceeds(t *testing.T) {
	req := validRequest()
	req.ReviewTimeout = 5 * time.Second
	req.PollInterval = time.Second
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}}
	check := &fakeCheckSource{states: []HeadCheckState{CheckPending, CheckPending, CheckSuccess}}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     &fakeResolver{},
		CheckSource:  check,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("expected Clean after check success, got %q", result.Outcome)
	}
	if fetcher.calls != 1 {
		t.Fatalf("expected one fetch before confirmation, got %d", fetcher.calls)
	}
	if check.calls != 3 {
		t.Fatalf("expected polling until the third check, got %d calls", check.calls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval, req.PollInterval)
}

func TestRunReentersFetchWhenHeadCheckFails(t *testing.T) {
	req := validRequest()
	req.MaxRounds = 2
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	statusHeadSHAs := []string{}
	check := &fakeCheckSource{states: []HeadCheckState{CheckFailure, CheckSuccess}}
	fetcher := &fakeFetcher{
		results: []FetchResult{
			{Round: 1, Issues: 1},
			{Round: 2, Issues: 1},
		},
	}
	resolver := &fakeResolver{
		results: []ResolveResult{
			{Remaining: 0, Progress: true, HeadSHA: "def456"},
			{Remaining: 0, Progress: true, HeadSHA: "fedcba"},
		},
	}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: StatusFunc(func(_ context.Context, req StatusRequest) (Status, error) {
			statusHeadSHAs = append(statusHeadSHAs, req.HeadSHA)
			return Status{State: StatusSettled}, nil
		}),
		Fetcher:     fetcher,
		Resolver:    resolver,
		CheckSource: check,
		Clock:       clock,
		Sleeper:     sleeper,
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("expected Clean after second Round succeeds, got %q", result.Outcome)
	}
	if result.Rounds != 2 {
		t.Fatalf("expected 2 Rounds, got %d", result.Rounds)
	}
	if fetcher.calls != 2 || resolver.calls != 2 {
		t.Fatalf("expected failure check to re-enter fetch and resolve, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
	if check.calls != 2 {
		t.Fatalf("expected two check calls, got %d", check.calls)
	}
	wantStatusHeads := []string{"abc123", "def456"}
	if len(statusHeadSHAs) != len(wantStatusHeads) {
		t.Fatalf("expected status head SHAs %#v, got %#v", wantStatusHeads, statusHeadSHAs)
	}
	for i := range wantStatusHeads {
		if statusHeadSHAs[i] != wantStatusHeads[i] {
			t.Fatalf("expected status head SHAs %#v, got %#v", wantStatusHeads, statusHeadSHAs)
		}
	}
	wantCheckHeads := []string{"def456", "fedcba"}
	if len(check.headSHAs) != len(wantCheckHeads) {
		t.Fatalf("expected check head SHAs %#v, got %#v", wantCheckHeads, check.headSHAs)
	}
	for i := range wantCheckHeads {
		if check.headSHAs[i] != wantCheckHeads[i] {
			t.Fatalf("expected check head SHAs %#v, got %#v", wantCheckHeads, check.headSHAs)
		}
	}
}

func TestRunReturnsMaxRoundsReachedWhenHeadCheckFailsOnFinalRound(t *testing.T) {
	req := validRequest()
	req.MaxRounds = 1
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 1}}}
	resolver := &fakeResolver{results: []ResolveResult{{Remaining: 0, Progress: true}}}
	check := &fakeCheckSource{states: []HeadCheckState{CheckFailure}}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		CheckSource:  check,
		Clock:        clock,
		Sleeper:      &fakeSleeper{clock: clock},
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateMaxRoundsReached {
		t.Fatalf("expected MaxRoundsReached, got %q", result.Outcome)
	}
	if result.Rounds != 1 {
		t.Fatalf("expected 1 Round, got %d", result.Rounds)
	}
	if fetcher.calls != 1 || resolver.calls != 1 || check.calls != 1 {
		t.Fatalf("expected one Round and one check, got fetch=%d resolve=%d check=%d", fetcher.calls, resolver.calls, check.calls)
	}
}

func TestRunTimesOutWhileHeadCheckStaysPending(t *testing.T) {
	req := validRequest()
	req.ReviewTimeout = 2 * time.Second
	req.PollInterval = time.Second
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	sleeper := &fakeSleeper{clock: clock}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}}
	check := &fakeCheckSource{states: []HeadCheckState{CheckPending, CheckPending, CheckPending}}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     &fakeResolver{},
		CheckSource:  check,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateTimedOut {
		t.Fatalf("expected TimedOut, got %q", result.Outcome)
	}
	if result.Rounds != 1 {
		t.Fatalf("expected timeout after 1 fetched Round, got %d", result.Rounds)
	}
	if result.ManualReviewCommand != "@coderabbitai review" {
		t.Fatalf("expected manual review trigger guidance, got %q", result.ManualReviewCommand)
	}
	if check.calls != 3 {
		t.Fatalf("expected three check calls before timeout, got %d", check.calls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval, req.PollInterval)
}

func TestRunDoesNotConfirmMergeReadinessWithoutUntilClean(t *testing.T) {
	req := validRequest()
	req.UntilClean = false
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 1}}}
	resolver := &fakeResolver{results: []ResolveResult{{Remaining: 0, Progress: true}}}

	result, err := Run(context.Background(), req, Dependencies{
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		CheckSource: CheckFunc(func(context.Context, string) (HeadCheckState, error) {
			t.Fatal("non-until-clean watch must not poll the merge-readiness check")
			return CheckFailure, nil
		}),
		Clock:   clock,
		Sleeper: &fakeSleeper{clock: clock},
	})

	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("expected legacy Clean, got %q", result.Outcome)
	}
}

func validRequest() Request {
	return Request{
		PRNumber:       "123",
		HeadSHA:        "abc123",
		UntilClean:     true,
		MaxRounds:      3,
		PollInterval:   time.Second,
		QuietPeriod:    2 * time.Second,
		ReviewTimeout:  5 * time.Second,
		BudgetEnabled:  true,
		MaxRunDuration: time.Minute,
	}
}

type fakeClock struct {
	now time.Time
}

func (clock *fakeClock) Now() time.Time {
	return clock.now
}

func (clock *fakeClock) Advance(duration time.Duration) {
	clock.now = clock.now.Add(duration)
}

type fakeSleeper struct {
	clock  *fakeClock
	err    error
	sleeps []time.Duration
}

func (sleeper *fakeSleeper) Sleep(_ context.Context, duration time.Duration) error {
	sleeper.sleeps = append(sleeper.sleeps, duration)
	if sleeper.err != nil {
		return sleeper.err
	}
	if sleeper.clock != nil {
		sleeper.clock.Advance(duration)
	}
	return nil
}

func (sleeper *fakeSleeper) saw(duration time.Duration) bool {
	for _, sleep := range sleeper.sleeps {
		if sleep == duration {
			return true
		}
	}
	return false
}

func assertSleeps(t *testing.T, got []time.Duration, want ...time.Duration) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected sleeps %#v, got %#v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected sleeps %#v, got %#v", want, got)
		}
	}
}

type fakeStatusSource struct {
	err      error
	calls    int
	statuses []Status
}

func (source *fakeStatusSource) Status(context.Context, StatusRequest) (Status, error) {
	source.calls++
	if source.err != nil {
		return Status{}, source.err
	}
	if len(source.statuses) == 0 {
		return Status{State: StatusSettled}, nil
	}
	status := source.statuses[0]
	if len(source.statuses) > 1 {
		source.statuses = source.statuses[1:]
	}
	return status, nil
}

type fakeCheckSource struct {
	err      error
	calls    int
	states   []HeadCheckState
	headSHAs []string
}

func (source *fakeCheckSource) Check(_ context.Context, headSHA string) (HeadCheckState, error) {
	source.calls++
	source.headSHAs = append(source.headSHAs, headSHA)
	if source.err != nil {
		return "", source.err
	}
	if len(source.states) == 0 {
		return CheckSuccess, nil
	}
	state := source.states[0]
	if len(source.states) > 1 {
		source.states = source.states[1:]
	}
	return state, nil
}

type fakeFetcher struct {
	err     error
	calls   int
	results []FetchResult
}

func (fetcher *fakeFetcher) Fetch(context.Context, int) (FetchResult, error) {
	fetcher.calls++
	if fetcher.err != nil {
		return FetchResult{}, fetcher.err
	}
	if len(fetcher.results) == 0 {
		return FetchResult{}, nil
	}
	result := fetcher.results[0]
	if len(fetcher.results) > 1 {
		fetcher.results = fetcher.results[1:]
	}
	return result, nil
}

type fakeResolver struct {
	err     error
	calls   int
	results []ResolveResult
}

func (resolver *fakeResolver) Resolve(context.Context) (ResolveResult, error) {
	resolver.calls++
	if resolver.err != nil {
		return ResolveResult{}, resolver.err
	}
	if len(resolver.results) == 0 {
		return ResolveResult{}, errors.New("missing fake resolve result")
	}
	result := resolver.results[0]
	if len(resolver.results) > 1 {
		resolver.results = resolver.results[1:]
	}
	return result, nil
}
