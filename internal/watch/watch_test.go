package watch

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"roundfix/internal/reviewsource"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

func TestRunReviewEvidenceSharedByPreFetchAndMergeReady(t *testing.T) {
	t.Parallel()
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	evidence := reviewsource.Evidence{
		State:           reviewsource.EvidenceVerified,
		Kind:            reviewsource.EvidenceKindReviewApproval,
		Identity:        "review:9001",
		ExpectedHeadSHA: "abc123",
		ObservedHeadSHA: "abc123",
		Conclusion:      "approved",
		Detail:          "CodeRabbit approved the expected head",
	}
	source := &fakeReviewEvidenceSource{evidence: evidence}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), validRequest(), Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        &fakeSleeper{clock: clock},
		Sink:           sink,
	})
	if err != nil {
		t.Fatalf("watch run: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("outcome = %q, want Clean", result.Outcome)
	}
	if !result.ReviewIssuesKnown {
		t.Fatal("successful zero-issue fetch left Review Issues unknown")
	}
	if result.VerifiedHeadSHA != evidence.ObservedHeadSHA {
		t.Fatalf("verified head = %q, want %q", result.VerifiedHeadSHA, evidence.ObservedHeadSHA)
	}
	if result.Evidence != evidence {
		t.Fatalf("terminal Evidence = %#v, want %#v", result.Evidence, evidence)
	}
	if len(source.requests) != 2 {
		t.Fatalf("evidence calls = %d, want pre-fetch and Merge-Ready calls", len(source.requests))
	}
	if source.requests[0] != source.requests[1] {
		t.Fatalf("watch phases received different requests: %#v", source.requests)
	}

	reviewEvents := sink.eventsOfKind(runevent.KindDaemonReviewStatus)
	if len(reviewEvents) != 3 {
		t.Fatalf("wait phase events = %d, want one event per phase and no unchanged duplicate", len(reviewEvents))
	}
	var payload runevent.ReviewStatusPayload
	if err := json.Unmarshal(reviewEvents[2].Payload, &payload); err != nil {
		t.Fatalf("decode review status event: %v", err)
	}
	if payload.State != string(evidence.State) ||
		payload.Kind != string(evidence.Kind) ||
		payload.Identity != evidence.Identity ||
		payload.ExpectedHeadSHA != evidence.ExpectedHeadSHA ||
		payload.ObservedHeadSHA != evidence.ObservedHeadSHA ||
		payload.Conclusion != evidence.Conclusion ||
		payload.Detail != evidence.Detail {
		t.Fatalf("review status payload = %#v, want evidence %#v", payload, evidence)
	}
}

func TestRunReviewIssuesUnknownWhenStatusDiscoveryFailsBeforeFetch(t *testing.T) {
	t.Parallel()
	sourceErr := errors.New("discover Review Source status")
	fetcher := &fakeFetcher{}

	result, err := Run(context.Background(), validRequest(), Dependencies{
		ReviewEvidence: ReviewEvidenceFunc(func(context.Context, ReviewEvidenceRequest) (reviewsource.Evidence, error) {
			return reviewsource.Evidence{}, sourceErr
		}),
		Fetcher:  fetcher,
		Resolver: &fakeResolver{},
	})

	if !errors.Is(err, sourceErr) {
		t.Fatalf("status discovery error = %v, want %v", err, sourceErr)
	}
	if result.Outcome != store.StateFailed || result.ReviewIssuesKnown {
		t.Fatalf("pre-fetch failure result = %+v", result)
	}
	if result.TerminalReason == "" || result.NextAction == "" {
		t.Fatalf("pre-fetch failure lacks actionable terminal context: %+v", result)
	}
	if fetcher.calls != 0 {
		t.Fatalf("pre-fetch failure called fetch %d time(s)", fetcher.calls)
	}
}

func TestRunReviewIssuesKnownAfterFetchedZero(t *testing.T) {
	t.Parallel()
	req := validRequest()
	evidence := evidenceForHead(reviewsource.EvidenceVerified, reviewsource.EvidenceKindCheckRun, req.HeadSHA)

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: &fakeReviewEvidenceSource{results: []reviewEvidenceResult{
			{evidence: evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
			{evidence: evidence},
		}},
		Fetcher:  &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver: &fakeResolver{},
	})

	if err != nil {
		t.Fatalf("watch fetched zero Review Issues: %v", err)
	}
	if result.Outcome != store.StateClean || !result.ReviewIssuesKnown {
		t.Fatalf("fetched-zero result = %+v", result)
	}
	if result.Evidence != evidence || result.VerifiedHeadSHA != req.HeadSHA {
		t.Fatalf("fetched-zero terminal context = %+v, want Evidence %#v and head %q", result, evidence, req.HeadSHA)
	}
}

func TestRunTransientReviewEvidenceRecoversWithinExistingBounds(t *testing.T) {
	t.Parallel()
	req := validRequest()
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{
			{err: transientReviewError("discover Review Source evidence")},
			{evidence: evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
			{evidence: evidenceForHead(reviewsource.EvidenceVerified, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
		},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if err != nil {
		t.Fatalf("watch transient recovery: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("transient recovery outcome = %q, want Clean", result.Outcome)
	}
	if len(source.requests) != 3 {
		t.Fatalf("Review Source calls = %d, want transient, recovered, and Merge-Ready calls", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval, req.QuietPeriod)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry), "started", "recovered")
}

func TestRunTransientReviewEvidenceRecoversDuringMergeReadyWait(t *testing.T) {
	t.Parallel()
	req := validRequest()
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{
			{evidence: evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
			{err: transientReviewError("confirm Merge-Ready Evidence")},
			{evidence: evidenceForHead(reviewsource.EvidenceVerified, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
		},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if err != nil {
		t.Fatalf("Merge-Ready transient recovery: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("Merge-Ready transient recovery outcome = %q, want Clean", result.Outcome)
	}
	if len(source.requests) != 3 {
		t.Fatalf("Review Source calls = %d, want pre-fetch, transient, and recovered calls", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry), "started", "recovered")
}

func TestRunTransientReviewEvidenceExhaustsReviewTimeout(t *testing.T) {
	t.Parallel()
	req := validRequest()
	req.ReviewTimeout = 2 * time.Second
	req.PollInterval = time.Second
	req.BudgetEnabled = false
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{{err: transientReviewError("discover Review Source evidence")}},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if err != nil {
		t.Fatalf("transient timeout must settle as a bounded outcome: %v", err)
	}
	if result.Outcome != store.StateTimedOut {
		t.Fatalf("transient timeout outcome = %q, want TimedOut", result.Outcome)
	}
	if len(source.requests) != 2 {
		t.Fatalf("Review Source calls = %d, want no call at the timeout boundary", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps, time.Second, time.Second)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry), "started", "exhausted")
}

func TestRunTransientReviewEvidenceExhaustsMergeReadyTimeout(t *testing.T) {
	t.Parallel()
	req := validRequest()
	req.ReviewTimeout = 2 * time.Second
	req.PollInterval = time.Second
	req.BudgetEnabled = false
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{
			{evidence: evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
			{err: transientReviewError("confirm Merge-Ready Evidence")},
		},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if err != nil {
		t.Fatalf("Merge-Ready transient timeout must settle as a bounded outcome: %v", err)
	}
	if result.Outcome != store.StateTimedOut {
		t.Fatalf("Merge-Ready transient timeout outcome = %q, want TimedOut", result.Outcome)
	}
	if len(source.requests) != 3 {
		t.Fatalf("Review Source calls = %d, want one pre-fetch and two retry episode calls", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps, time.Second, time.Second)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry), "started", "exhausted")
}

func TestRunTransientReviewEvidenceAfterMissingCheckExhaustsReviewTimeout(t *testing.T) {
	t.Parallel()
	req := validRequest()
	req.ReviewTimeout = 2 * time.Second
	req.CheckGracePeriod = 10 * time.Second
	req.PollInterval = time.Second
	req.BudgetEnabled = false
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{
			{evidence: evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
			{evidence: evidenceForHead(reviewsource.EvidencePending, reviewsource.EvidenceKindNone, req.HeadSHA)},
			{err: transientReviewError("confirm Merge-Ready Evidence")},
		},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if err != nil {
		t.Fatalf("transient failure after a missing check must settle as a bounded outcome: %v", err)
	}
	if result.Outcome != store.StateTimedOut {
		t.Fatalf("transient failure after missing check outcome = %q, want TimedOut", result.Outcome)
	}
	if len(source.requests) != 3 {
		t.Fatalf("Review Source calls = %d, want pre-fetch, missing, and one transient call", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps, time.Second, time.Second)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry), "started", "exhausted")
}

func TestRunTransientReviewEvidenceExhaustsRunBudgetBeforeTimeout(t *testing.T) {
	t.Parallel()
	req := validRequest()
	req.ReviewTimeout = 5 * time.Second
	req.PollInterval = time.Second
	req.MaxRunDuration = time.Second
	startedAt := time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: startedAt}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{{err: transientReviewError("discover Review Source evidence")}},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if err != nil {
		t.Fatalf("transient Run Budget exhaustion must settle as a bounded outcome: %v", err)
	}
	if result.Outcome != store.StateBudgetExceeded {
		t.Fatalf("transient budget outcome = %q, want BudgetExceeded", result.Outcome)
	}
	if len(source.requests) != 1 {
		t.Fatalf("Review Source calls = %d, want no call at the Run Budget boundary", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps, time.Second)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry), "started", "exhausted")
	reviewEvents := sink.eventsOfKind(runevent.KindDaemonReviewStatus)
	if len(reviewEvents) == 0 {
		t.Fatal("Run Budget retry emitted no wait projection")
	}
	var payload runevent.ReviewStatusPayload
	if err := json.Unmarshal(reviewEvents[0].Payload, &payload); err != nil {
		t.Fatalf("decode Run Budget wait event: %v", err)
	}
	wantDeadline := startedAt.Add(req.MaxRunDuration)
	if !payload.Deadline.Equal(wantDeadline) {
		t.Fatalf("wait deadline = %s, want Run Budget deadline %s", payload.Deadline, wantDeadline)
	}
}

func TestRunPermanentReviewEvidenceFailureDoesNotRetry(t *testing.T) {
	t.Parallel()
	req := validRequest()
	permanent := errors.New("authentication failed: temporary Review Source failure")
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{{err: permanent}},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if !errors.Is(err, permanent) {
		t.Fatalf("permanent failure error = %v, want %v", err, permanent)
	}
	if result.Outcome != store.StateFailed {
		t.Fatalf("permanent failure outcome = %q, want Failed", result.Outcome)
	}
	if len(source.requests) != 1 {
		t.Fatalf("permanent failure calls = %d, want one", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry))
}

func TestRunPermanentMergeReadyEvidenceFailureDoesNotRetry(t *testing.T) {
	t.Parallel()
	req := validRequest()
	permanent := errors.New("validation failed")
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{
			{evidence: evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
			{err: permanent},
		},
	}
	sleeper := &fakeSleeper{clock: clock}
	sink := &recordingEventSink{}

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        sleeper,
		Sink:           sink,
	})
	if !errors.Is(err, permanent) {
		t.Fatalf("permanent Merge-Ready error = %v, want %v", err, permanent)
	}
	if result.Outcome != store.StateFailed {
		t.Fatalf("permanent Merge-Ready outcome = %q, want Failed", result.Outcome)
	}
	if len(source.requests) != 2 {
		t.Fatalf("permanent Merge-Ready calls = %d, want pre-fetch and one confirmation", len(source.requests))
	}
	assertSleeps(t, sleeper.sleeps)
	assertRetryPhases(t, sink.eventsOfKind(runevent.KindDaemonRetry))
}

func TestRunWaitPhaseProjectionDeduplicatesUnchangedEvidence(t *testing.T) {
	t.Parallel()
	req := validRequest()
	req.QuietPeriod = 0
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	pending := evidenceForHead(reviewsource.EvidencePending, reviewsource.EvidenceKindNone, req.HeadSHA)
	reviewed := evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)
	verified := evidenceForHead(reviewsource.EvidenceVerified, reviewsource.EvidenceKindCheckRun, req.HeadSHA)
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{
			{evidence: pending},
			{evidence: pending},
			{evidence: reviewed},
			{evidence: pending},
			{evidence: pending},
			{evidence: verified},
		},
	}
	sink := &recordingEventSink{}
	var progress []WaitProgress

	result, err := Run(context.Background(), req, Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        &fakeSleeper{clock: clock},
		Sink:           sink,
		Progress: func(update WaitProgress) {
			progress = append(progress, update)
		},
	})
	if err != nil {
		t.Fatalf("watch wait projection: %v", err)
	}
	if result.Outcome != store.StateClean {
		t.Fatalf("wait projection outcome = %q, want Clean", result.Outcome)
	}

	events := sink.eventsOfKind(runevent.KindDaemonReviewStatus)
	if len(events) != 6 {
		t.Fatalf("review wait events = %d, want phase entry plus changed Evidence only", len(events))
	}
	if len(progress) != len(events) {
		t.Fatalf("direct progress updates = %d, want one per persisted wait event (%d)", len(progress), len(events))
	}
	wantPhases := []string{
		"WaitingForReview",
		"WaitingForReview",
		"WaitingForReview",
		"WaitingForReviewCheck",
		"WaitingForReviewCheck",
		"WaitingForReviewCheck",
	}
	for index, event := range events {
		var payload runevent.ReviewStatusPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatalf("decode wait event %d: %v", index, err)
		}
		if payload.Phase != wantPhases[index] {
			t.Fatalf("wait event %d phase = %q, want %q", index, payload.Phase, wantPhases[index])
		}
		if payload.ExpectedHeadSHA != req.HeadSHA {
			t.Fatalf("wait event %d expected head = %q, want %q", index, payload.ExpectedHeadSHA, req.HeadSHA)
		}
		if payload.StartedAt.IsZero() || payload.Deadline.IsZero() || !payload.Deadline.After(payload.StartedAt) {
			t.Fatalf("wait event %d has invalid bounds: start=%s deadline=%s", index, payload.StartedAt, payload.Deadline)
		}
		if payload.EvidenceState == "" || payload.EvidenceKind == "" || payload.RetryStatus == "" {
			t.Fatalf("wait event %d missing Evidence or retry projection: %#v", index, payload)
		}
		if progress[index].Phase != payload.Phase ||
			progress[index].ExpectedHeadSHA != payload.ExpectedHeadSHA ||
			string(progress[index].Evidence.State) != payload.EvidenceState ||
			string(progress[index].Evidence.Kind) != payload.EvidenceKind ||
			progress[index].RetryStatus != payload.RetryStatus {
			t.Fatalf("direct progress %d diverged from persisted payload: progress=%#v payload=%#v", index, progress[index], payload)
		}
	}
}

func TestRunReviewSkippedStopsBeforeFetch(t *testing.T) {
	t.Parallel()
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	evidence := reviewsource.Evidence{
		State:           reviewsource.EvidenceSkipped,
		Kind:            reviewsource.EvidenceKindCheckRun,
		Identity:        "check_run:42",
		ExpectedHeadSHA: "abc123",
		ObservedHeadSHA: "abc123",
		Conclusion:      "success",
		Detail:          "CodeRabbit skipped the review",
		Reason:          "Pull request is too large to review",
	}
	source := &fakeReviewEvidenceSource{evidence: evidence}
	fetcher := &fakeFetcher{}
	resolver := &fakeResolver{}

	result, err := Run(context.Background(), validRequest(), Dependencies{
		ReviewEvidence: source,
		Fetcher:        fetcher,
		Resolver:       resolver,
		Clock:          clock,
		Sleeper:        &fakeSleeper{clock: clock},
	})
	if err != nil {
		t.Fatalf("watch Review Skipped: %v", err)
	}
	if result.Outcome != store.StateReviewSkipped || result.Rounds != 0 {
		t.Fatalf("Review Skipped result = %+v", result)
	}
	if result.ReviewIssuesKnown {
		t.Fatal("pre-fetch Review Skipped marked Review Issues known")
	}
	if result.TerminalReason != evidence.Reason {
		t.Fatalf("terminal reason = %q, want %q", result.TerminalReason, evidence.Reason)
	}
	if result.NextAction != "Reduce or split the pull request, then request another Review Source review." {
		t.Fatalf("next action = %q", result.NextAction)
	}
	if result.Evidence != evidence {
		t.Fatalf("terminal Evidence = %#v, want %#v", result.Evidence, evidence)
	}
	if fetcher.calls != 0 || resolver.calls != 0 {
		t.Fatalf("Review Skipped reached later work: fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
	if len(source.requests) != 1 {
		t.Fatalf("Review Skipped Evidence calls = %d, want 1", len(source.requests))
	}
}

func TestRunReviewSkippedDuringMergeReadyPreservesTerminalEvidence(t *testing.T) {
	t.Parallel()
	clock := &fakeClock{now: time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC)}
	calls := 0
	skipped := reviewsource.Evidence{
		State:           reviewsource.EvidenceSkipped,
		Kind:            reviewsource.EvidenceKindCheckRun,
		Identity:        "check_run:42",
		ExpectedHeadSHA: "abc123",
		ObservedHeadSHA: "abc123",
		Conclusion:      "success",
		Reason:          "Review Source size limit was exceeded",
	}
	source := ReviewEvidenceFunc(func(_ context.Context, _ ReviewEvidenceRequest) (reviewsource.Evidence, error) {
		calls++
		if calls == 1 {
			return reviewsource.Evidence{
				State:           reviewsource.EvidenceReviewed,
				Kind:            reviewsource.EvidenceKindCheckRun,
				ExpectedHeadSHA: "abc123",
				ObservedHeadSHA: "abc123",
			}, nil
		}
		return skipped, nil
	})

	result, err := Run(context.Background(), validRequest(), Dependencies{
		ReviewEvidence: source,
		Fetcher:        &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}},
		Resolver:       &fakeResolver{},
		Clock:          clock,
		Sleeper:        &fakeSleeper{clock: clock},
	})
	if err != nil {
		t.Fatalf("watch Review Skipped during Merge-Ready: %v", err)
	}
	if result.Outcome != store.StateReviewSkipped || result.Rounds != 1 || result.Evidence != skipped {
		t.Fatalf("Merge-Ready Review Skipped result = %+v", result)
	}
	if !result.ReviewIssuesKnown {
		t.Fatal("post-fetch Review Skipped lost Review Issue knowledge")
	}
	if calls != 2 {
		t.Fatalf("Review Evidence calls = %d, want pre-fetch and Merge-Ready", calls)
	}
}

func TestRunWaitsFetchesResolvesToClean(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestRunStopRequestDuringStatusWaitStopsAtNextPoll(t *testing.T) {
	t.Parallel()
	req := validRequest()
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	stops := &fakeStopRequestSource{}
	sleeper := &fakeSleeper{
		clock: clock,
		afterSleep: func(time.Duration) {
			stops.requested = true
		},
	}
	status := &fakeStatusSource{
		statuses: []Status{
			{State: StatusPending},
			{State: StatusSettled},
		},
	}
	fetcher := &fakeFetcher{}
	resolver := &fakeResolver{}

	result, err := Run(context.Background(), req, Dependencies{
		StopRequests: stops,
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected Stop Request classification, got result=%#v err=%v", result, err)
	}
	if result.Outcome != store.StateStopped || result.Rounds != 0 {
		t.Fatalf("expected Stopped before a fetched Round, got %#v", result)
	}
	if stops.calls != 3 || status.calls != 1 {
		t.Fatalf("expected Stop Request after one status poll, got stop=%d status=%d", stops.calls, status.calls)
	}
	if fetcher.calls != 0 || resolver.calls != 0 {
		t.Fatalf("Stop Request must block later work, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval)
}

func TestRunStopRequestDuringQuietPeriodStopsBeforeFetch(t *testing.T) {
	t.Parallel()
	req := validRequest()
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	stops := &fakeStopRequestSource{}
	sleeper := &fakeSleeper{
		clock: clock,
		afterSleep: func(duration time.Duration) {
			if duration == req.QuietPeriod {
				stops.requested = true
			}
		},
	}
	status := &fakeStatusSource{
		statuses: []Status{
			{State: StatusPending},
			{State: StatusSettled},
		},
	}
	fetcher := &fakeFetcher{}
	resolver := &fakeResolver{}

	result, err := Run(context.Background(), req, Dependencies{
		StopRequests: stops,
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected Stop Request classification, got result=%#v err=%v", result, err)
	}
	if result.Outcome != store.StateStopped || result.Rounds != 0 {
		t.Fatalf("expected Stopped before fetch, got %#v", result)
	}
	if stops.calls != 6 || status.calls != 2 {
		t.Fatalf("expected Stop Request after quiet-period sleep, got stop=%d status=%d", stops.calls, status.calls)
	}
	if fetcher.calls != 0 || resolver.calls != 0 {
		t.Fatalf("quiet-period Stop Request must block later work, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval, req.QuietPeriod)
}

func TestRunStopRequestDuringTransientRetryStopsBeforeNextCheck(t *testing.T) {
	t.Parallel()
	req := validRequest()
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	stops := &fakeStopRequestSource{}
	sleeper := &fakeSleeper{
		clock: clock,
		afterSleep: func(time.Duration) {
			stops.requested = true
		},
	}
	source := &fakeReviewEvidenceSource{
		results: []reviewEvidenceResult{
			{evidence: evidenceForHead(reviewsource.EvidenceReviewed, reviewsource.EvidenceKindCheckRun, req.HeadSHA)},
			{err: transientReviewError("check Merge-Ready Evidence")},
		},
	}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}}
	resolver := &fakeResolver{}

	result, err := Run(context.Background(), req, Dependencies{
		StopRequests:   stops,
		ReviewEvidence: source,
		Fetcher:        fetcher,
		Resolver:       resolver,
		Clock:          clock,
		Sleeper:        sleeper,
	})

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected Stop Request classification, got result=%#v err=%v", result, err)
	}
	if result.Outcome != store.StateStopped || result.Rounds != 1 {
		t.Fatalf("expected Stopped after one fetched Round, got %#v", result)
	}
	if !stops.observed || len(source.requests) != 2 {
		t.Fatalf("expected Stop Request after one retry sleep, got observed=%v evidence=%d", stops.observed, len(source.requests))
	}
	if fetcher.calls != 1 || resolver.calls != 0 {
		t.Fatalf("retry Stop Request must block later work, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval)
}

func TestRunStopRequestDuringMergeReadyWaitStopsBeforeNextCheck(t *testing.T) {
	t.Parallel()
	req := validRequest()
	clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	stops := &fakeStopRequestSource{}
	sleeper := &fakeSleeper{
		clock: clock,
		afterSleep: func(time.Duration) {
			stops.requested = true
		},
	}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
	fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}}
	check := &fakeCheckSource{states: []HeadCheckState{CheckPending, CheckSuccess}}
	resolver := &fakeResolver{}

	result, err := Run(context.Background(), req, Dependencies{
		StopRequests: stops,
		StatusSource: status,
		Fetcher:      fetcher,
		Resolver:     resolver,
		CheckSource:  check,
		Clock:        clock,
		Sleeper:      sleeper,
	})

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected Stop Request classification, got result=%#v err=%v", result, err)
	}
	if result.Outcome != store.StateStopped || result.Rounds != 1 {
		t.Fatalf("expected Stopped after one fetched Round, got %#v", result)
	}
	if stops.calls != 5 || check.calls != 1 {
		t.Fatalf("expected Stop Request after one Merge-Ready sleep, got stop=%d check=%d", stops.calls, check.calls)
	}
	if fetcher.calls != 1 || resolver.calls != 0 {
		t.Fatalf("Merge-Ready Stop Request must block later work, got fetch=%d resolve=%d", fetcher.calls, resolver.calls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval)
}

func TestRunStopRequestSourceFailureIncludesRunAndOperation(t *testing.T) {
	t.Parallel()
	req := validRequest()
	sourceErr := errors.New("Run Database unavailable")
	stops := &fakeStopRequestSource{err: sourceErr, errAtCall: 1}
	status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}

	result, err := Run(context.Background(), req, Dependencies{
		StopRequests: stops,
		StatusSource: status,
		Fetcher:      &fakeFetcher{},
		Resolver:     &fakeResolver{},
	})

	if !errors.Is(err, sourceErr) {
		t.Fatalf("expected Store observation failure, got result=%#v err=%v", result, err)
	}
	if errors.Is(err, ErrStopRequested) {
		t.Fatalf("Store observation failure must not classify as requested stop: %v", err)
	}
	if result.Outcome != store.StateFailed {
		t.Fatalf("expected Failed observation result, got %#v", result)
	}
	if !strings.Contains(err.Error(), req.RunID) || !strings.Contains(err.Error(), "before Review Source status") {
		t.Fatalf("expected Run and operation context, got %q", err)
	}
	if status.calls != 0 {
		t.Fatalf("observation failure must block Review Source status, got %d calls", status.calls)
	}
}

func TestRunTimesOutAndOffersManualReviewTrigger(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestRunWithoutStopRequestKeepsRunBudgetBehavior(t *testing.T) {
	t.Parallel()
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
		StopRequests: &fakeStopRequestSource{},
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
	t.Parallel()
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

func TestRunConfirmsMergeReadyThroughGraceWindow(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		states           []HeadCheckState
		maxRounds        int
		reviewTimeout    time.Duration
		checkGracePeriod time.Duration
		wantOutcome      string
		wantManualReview bool
		wantCheckCalls   int
		wantSleeps       []time.Duration
	}{
		{
			name:           "success is clean immediately",
			states:         []HeadCheckState{CheckSuccess},
			wantOutcome:    store.StateClean,
			wantCheckCalls: 1,
		},
		{
			name:           "missing check appears late and succeeds within grace window",
			states:         []HeadCheckState{CheckMissing, CheckMissing, CheckSuccess},
			wantOutcome:    store.StateClean,
			wantCheckCalls: 3,
			wantSleeps:     []time.Duration{time.Second, time.Second},
		},
		{
			name:           "missing check appears late and fails within grace window",
			states:         []HeadCheckState{CheckMissing, CheckFailure},
			maxRounds:      1,
			wantOutcome:    store.StateMaxRoundsReached,
			wantCheckCalls: 2,
			wantSleeps:     []time.Duration{time.Second},
		},
		{
			name:             "missing check exhausts grace window as clean unverified",
			states:           []HeadCheckState{CheckMissing},
			reviewTimeout:    time.Second,
			checkGracePeriod: 2 * time.Second,
			wantOutcome:      store.StateCleanUnverified,
			wantCheckCalls:   2,
			wantSleeps:       []time.Duration{time.Second, time.Second},
		},
		{
			name:             "pending check still uses review timeout",
			states:           []HeadCheckState{CheckPending},
			reviewTimeout:    2 * time.Second,
			checkGracePeriod: 5 * time.Second,
			wantOutcome:      store.StateTimedOut,
			wantManualReview: true,
			wantCheckCalls:   2,
			wantSleeps:       []time.Duration{time.Second, time.Second},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validRequest()
			if tt.maxRounds > 0 {
				req.MaxRounds = tt.maxRounds
			}
			if tt.reviewTimeout > 0 {
				req.ReviewTimeout = tt.reviewTimeout
			}
			if tt.checkGracePeriod > 0 {
				req.CheckGracePeriod = tt.checkGracePeriod
			}
			clock := &fakeClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
			sleeper := &fakeSleeper{clock: clock}
			status := &fakeStatusSource{statuses: []Status{{State: StatusSettled}}}
			fetcher := &fakeFetcher{results: []FetchResult{{Round: 1, Issues: 0}}}
			check := &fakeCheckSource{states: tt.states}

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
			if result.Outcome != tt.wantOutcome {
				t.Fatalf("expected %s, got %q", tt.wantOutcome, result.Outcome)
			}
			if result.Rounds != 1 {
				t.Fatalf("expected 1 Round, got %d", result.Rounds)
			}
			if tt.wantManualReview && result.ManualReviewCommand != "@coderabbitai review" {
				t.Fatalf("expected manual review trigger guidance, got %q", result.ManualReviewCommand)
			}
			if check.calls != tt.wantCheckCalls {
				t.Fatalf("expected %d check calls, got %d", tt.wantCheckCalls, check.calls)
			}
			assertSleeps(t, sleeper.sleeps, tt.wantSleeps...)
		})
	}
}

func TestRunKeepsPollingWhenLocalQueueIsEmptyUntilHeadCheckSucceeds(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	if check.calls != 2 {
		t.Fatalf("expected no check call at the timeout boundary, got %d", check.calls)
	}
	assertSleeps(t, sleeper.sleeps, req.PollInterval, req.PollInterval)
}

func TestRunDoesNotConfirmMergeReadinessWithoutUntilClean(t *testing.T) {
	t.Parallel()
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
		RunID:            "run_123",
		PRNumber:         "123",
		HeadSHA:          "abc123",
		UntilClean:       true,
		MaxRounds:        3,
		PollInterval:     time.Second,
		QuietPeriod:      2 * time.Second,
		ReviewTimeout:    5 * time.Second,
		CheckGracePeriod: 3 * time.Second,
		BudgetEnabled:    true,
		MaxRunDuration:   time.Minute,
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
	clock      *fakeClock
	err        error
	afterSleep func(time.Duration)
	sleeps     []time.Duration
}

func (sleeper *fakeSleeper) Sleep(_ context.Context, duration time.Duration) error {
	sleeper.sleeps = append(sleeper.sleeps, duration)
	if sleeper.err != nil {
		return sleeper.err
	}
	if sleeper.clock != nil {
		sleeper.clock.Advance(duration)
	}
	if sleeper.afterSleep != nil {
		sleeper.afterSleep(duration)
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

type fakeReviewEvidenceSource struct {
	evidence reviewsource.Evidence
	err      error
	results  []reviewEvidenceResult
	requests []ReviewEvidenceRequest
}

// Evidence returns queued results in order and repeats the final result for
// every later call, so steady-state phases need only one trailing entry.
func (source *fakeReviewEvidenceSource) Evidence(_ context.Context, req ReviewEvidenceRequest) (reviewsource.Evidence, error) {
	source.requests = append(source.requests, req)
	if len(source.results) > 0 {
		result := source.results[0]
		if len(source.results) > 1 {
			source.results = source.results[1:]
		}
		return result.evidence, result.err
	}
	return source.evidence, source.err
}

type reviewEvidenceResult struct {
	evidence reviewsource.Evidence
	err      error
}

func evidenceForHead(state reviewsource.EvidenceState, kind reviewsource.EvidenceKind, headSHA string) reviewsource.Evidence {
	return reviewsource.Evidence{
		State:           state,
		Kind:            kind,
		ExpectedHeadSHA: headSHA,
		ObservedHeadSHA: headSHA,
	}
}

func transientReviewError(operation string) error {
	return &reviewsource.TransientError{
		Operation: operation,
		Err:       errors.New("temporary network failure"),
	}
}

type recordingEventSink struct {
	events []runevent.RunEvent
}

func (sink *recordingEventSink) Publish(_ context.Context, event runevent.RunEvent) error {
	sink.events = append(sink.events, event)
	return nil
}

func (sink *recordingEventSink) eventsOfKind(kind runevent.Kind) []runevent.RunEvent {
	var events []runevent.RunEvent
	for _, event := range sink.events {
		if event.Kind == kind {
			events = append(events, event)
		}
	}
	return events
}

func assertRetryPhases(t *testing.T, events []runevent.RunEvent, want ...string) {
	t.Helper()
	if len(events) != len(want) {
		t.Fatalf("retry events = %d, want %d phases %v", len(events), len(want), want)
	}
	for index, event := range events {
		var payload runevent.RetryPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatalf("decode retry event %d: %v", index, err)
		}
		if payload.Phase != want[index] {
			t.Fatalf("retry event %d phase = %q, want %q", index, payload.Phase, want[index])
		}
		if payload.Operation == "" || payload.Reason == "" {
			t.Fatalf("retry event %d missing bounded context: %#v", index, payload)
		}
	}
}

type fakeStopRequestSource struct {
	err       error
	errAtCall int
	requested bool
	observed  bool
	calls     int
}

func (source *fakeStopRequestSource) StopRequested(context.Context, string) (bool, error) {
	source.calls++
	if source.err != nil && source.calls == source.errAtCall {
		return false, source.err
	}
	if source.requested {
		source.observed = true
	}
	return source.requested, nil
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
