package watch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"roundfix/internal/preflight"
	"roundfix/internal/reviewsource"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

const (
	StatusPending   = "pending"
	StatusReviewing = "reviewing"
	StatusSettled   = "settled"
)

const (
	WaitPhaseReview      = "WaitingForReview"
	WaitPhaseReviewCheck = "WaitingForReviewCheck"

	RetryStatusNone      = "none"
	RetryStatusRetrying  = "retrying"
	RetryStatusRecovered = "recovered"
	RetryStatusExhausted = "exhausted"
)

var ErrStopRequested = errors.New("stop requested")

const reviewSkippedNextAction = "Reduce or split the pull request, then request another Review Source review."

type HeadCheckState string

const (
	CheckPending HeadCheckState = "pending"
	CheckSuccess HeadCheckState = "success"
	CheckFailure HeadCheckState = "failure"
	CheckMissing HeadCheckState = "missing"
)

type Request struct {
	RunID            string
	PRNumber         string
	BaseRepository   string
	HeadSHA          string
	RequestReview    bool
	ReviewCommand    string
	UntilClean       bool
	MaxRounds        int
	PollInterval     time.Duration
	QuietPeriod      time.Duration
	ReviewTimeout    time.Duration
	CheckGracePeriod time.Duration
	BudgetEnabled    bool
	MaxRunDuration   time.Duration
}

type StatusRequest struct {
	PRNumber string
	HeadSHA  string
}

type Status struct {
	State  string
	Detail string
}

type FetchResult struct {
	Round  int
	Issues int
}

type ResolveResult struct {
	Remaining int
	Progress  bool
	HeadSHA   string
	Outcome   string
}

type ArtifactPublishRequest struct {
	Round          int
	ParentHeadSHA  string
	ParentEvidence reviewsource.Evidence
}

type ArtifactPublication struct {
	Commit   reviewsource.ArtifactCommit
	Evidence reviewsource.Evidence
}

type ArtifactPublisher interface {
	PublishArtifacts(context.Context, ArtifactPublishRequest) (ArtifactPublication, error)
}

type ArtifactPublishFunc func(context.Context, ArtifactPublishRequest) (ArtifactPublication, error)

func (fn ArtifactPublishFunc) PublishArtifacts(ctx context.Context, req ArtifactPublishRequest) (ArtifactPublication, error) {
	return fn(ctx, req)
}

type Result struct {
	Outcome                    string
	Rounds                     int
	Remaining                  int
	ManualReviewCommand        string
	ReviewIssuesKnown          bool
	TerminalReason             string
	TerminalReasonIsDiagnostic bool
	NextAction                 string
	Evidence                   reviewsource.Evidence
	VerifiedHeadSHA            string
}

type StatusSource interface {
	Status(context.Context, StatusRequest) (Status, error)
}

type StopRequestSource interface {
	StopRequested(context.Context, string) (bool, error)
}

type StatusFunc func(context.Context, StatusRequest) (Status, error)

func (fn StatusFunc) Status(ctx context.Context, req StatusRequest) (Status, error) {
	return fn(ctx, req)
}

type ReviewEvidenceRequest struct {
	PRNumber        string
	ExpectedHeadSHA string
}

type ReviewEvidenceSource interface {
	Evidence(context.Context, ReviewEvidenceRequest) (reviewsource.Evidence, error)
}

type ReviewEvidenceFunc func(context.Context, ReviewEvidenceRequest) (reviewsource.Evidence, error)

func (fn ReviewEvidenceFunc) Evidence(ctx context.Context, req ReviewEvidenceRequest) (reviewsource.Evidence, error) {
	return fn(ctx, req)
}

// WaitProgress projects one changed Review Source wait observation to
// interactive and non-interactive consumers.
type WaitProgress struct {
	Phase           string
	ExpectedHeadSHA string
	StartedAt       time.Time
	Deadline        time.Time
	Evidence        reviewsource.Evidence
	RetryStatus     string
}

type Fetcher interface {
	Fetch(context.Context, int) (FetchResult, error)
}

type FetchFunc func(context.Context, int) (FetchResult, error)

func (fn FetchFunc) Fetch(ctx context.Context, round int) (FetchResult, error) {
	return fn(ctx, round)
}

type Resolver interface {
	Resolve(context.Context) (ResolveResult, error)
}

type ResolveFunc func(context.Context) (ResolveResult, error)

func (fn ResolveFunc) Resolve(ctx context.Context) (ResolveResult, error) {
	return fn(ctx)
}

// Checkout is the branch and revision observed in the user checkout at one
// write boundary.
type Checkout struct {
	Branch   string
	Revision string
}

// CheckoutReader re-reads the user checkout without moving it.
type CheckoutReader interface {
	ReadCheckout(context.Context, string) (Checkout, error)
}

type CheckoutReadFunc func(context.Context, string) (Checkout, error)

func (fn CheckoutReadFunc) ReadCheckout(ctx context.Context, workDir string) (Checkout, error) {
	return fn(ctx, workDir)
}

// GitCheckoutReader reads the current branch and revision through the same
// Git inspection path used by Preflight.
type GitCheckoutReader struct {
	Runner preflight.GitRunner
}

func (reader GitCheckoutReader) ReadCheckout(ctx context.Context, workDir string) (Checkout, error) {
	state, err := preflight.InspectGit(ctx, workDir, reader.Runner)
	if err != nil {
		return Checkout{}, fmt.Errorf("read checkout target: %w", err)
	}
	return Checkout{Branch: state.Branch, Revision: state.HEAD}, nil
}

// CheckoutMovedError reports a user checkout that no longer matches the
// target recorded on the Run at Preflight.
type CheckoutMovedError struct {
	ExpectedBranch   string
	ExpectedRevision string
	FoundBranch      string
	FoundRevision    string
}

func (moved CheckoutMovedError) Error() string {
	return fmt.Sprintf(
		"checkout moved: expected PR Head Branch %q at revision %q, found branch %q at revision %q",
		moved.ExpectedBranch,
		displayRevision(moved.ExpectedRevision),
		moved.FoundBranch,
		displayRevision(moved.FoundRevision),
	)
}

func displayRevision(revision string) string {
	if revision = strings.TrimSpace(revision); revision != "" {
		return revision
	}
	return "<unknown>"
}

// WriteBoundaryGuard stops a Run before a commit or push when its checkout no
// longer matches the recorded target.
type WriteBoundaryGuard interface {
	Check(context.Context) error
}

type WriteBoundaryGuardFunc func(context.Context) error

func (fn WriteBoundaryGuardFunc) Check(ctx context.Context) error {
	return fn(ctx)
}

type TargetGuard struct {
	WorkDir   string
	Target    Checkout
	Checkouts CheckoutReader
	mu        sync.Mutex
}

func (guard *TargetGuard) Check(ctx context.Context) error {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	return guard.checkLocked(ctx)
}

func (guard *TargetGuard) checkLocked(ctx context.Context) error {
	reader := guard.Checkouts
	if reader == nil {
		reader = GitCheckoutReader{}
	}
	found, err := reader.ReadCheckout(ctx, guard.WorkDir)
	if err != nil {
		return err
	}
	expectedBranch := strings.TrimSpace(guard.Target.Branch)
	expectedRevision := strings.TrimSpace(guard.Target.Revision)
	foundBranch := strings.TrimSpace(found.Branch)
	foundRevision := strings.TrimSpace(found.Revision)
	if foundBranch == expectedBranch && foundRevision == expectedRevision {
		return nil
	}
	return CheckoutMovedError{
		ExpectedBranch:   expectedBranch,
		ExpectedRevision: expectedRevision,
		FoundBranch:      foundBranch,
		FoundRevision:    foundRevision,
	}
}

// RecordWriteBoundary advances the expected revision after a successful
// Roundfix-owned commit while retaining the recorded PR Head Branch.
func (guard *TargetGuard) RecordWriteBoundary(ctx context.Context) error {
	guard.mu.Lock()
	defer guard.mu.Unlock()
	reader := guard.Checkouts
	if reader == nil {
		reader = GitCheckoutReader{}
	}
	found, err := reader.ReadCheckout(ctx, guard.WorkDir)
	if err != nil {
		return err
	}
	expectedBranch := strings.TrimSpace(guard.Target.Branch)
	foundBranch := strings.TrimSpace(found.Branch)
	if foundBranch != expectedBranch {
		return CheckoutMovedError{
			ExpectedBranch:   expectedBranch,
			ExpectedRevision: guard.Target.Revision,
			FoundBranch:      foundBranch,
			FoundRevision:    found.Revision,
		}
	}
	guard.Target.Revision = strings.TrimSpace(found.Revision)
	return nil
}

type CheckSource interface {
	Check(context.Context, string) (HeadCheckState, error)
}

type CheckFunc func(context.Context, string) (HeadCheckState, error)

func (fn CheckFunc) Check(ctx context.Context, headSHA string) (HeadCheckState, error) {
	return fn(ctx, headSHA)
}

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

type Sleeper interface {
	Sleep(context.Context, time.Duration) error
}

type realSleeper struct{}

func (realSleeper) Sleep(ctx context.Context, duration time.Duration) error {
	if duration <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

type Dependencies struct {
	StopRequests    StopRequestSource
	ReviewEvidence  ReviewEvidenceSource
	ReviewRequester reviewsource.ReviewRequester
	Artifacts       ArtifactPublisher
	StatusSource    StatusSource
	Fetcher         Fetcher
	Resolver        Resolver
	WriteGuard      WriteBoundaryGuard
	CheckSource     CheckSource
	Clock           Clock
	Sleeper         Sleeper
	// Sink receives watch-loop Run Events: review status waits, quiet
	// periods, fetch results, and merge-readiness checks. Nil means
	// events are discarded.
	Sink runevent.Sink
	// Progress receives the same deduplicated Review Source wait projection
	// that is persisted as daemon.review_status. Nil disables direct progress.
	Progress func(WaitProgress)
}

func Run(ctx context.Context, req Request, deps Dependencies) (result Result, resultErr error) {
	if err := validateRequest(req, deps); err != nil {
		return Result{}, err
	}
	reviewIssuesKnown := false
	terminalEvidence := reviewsource.Evidence{}
	verifiedHeadSHA := ""
	defer func() {
		finalizeResult(&result, resultErr, reviewIssuesKnown, terminalEvidence, verifiedHeadSHA)
	}()
	clock := deps.Clock
	if clock == nil {
		clock = realClock{}
	}
	sleeper := deps.Sleeper
	if sleeper == nil {
		sleeper = realSleeper{}
	}
	if deps.Sink == nil {
		deps.Sink = runevent.Discard
	}
	publisher := &watchEventPublisher{
		sink:     deps.Sink,
		runID:    req.RunID,
		clock:    clock,
		progress: deps.Progress,
	}

	startedAt := clock.Now()
	runDeadline := time.Time{}
	if req.BudgetEnabled {
		runDeadline = startedAt.Add(req.MaxRunDuration)
	}
	currentHeadSHA := req.HeadSHA
	for round := 1; round <= req.MaxRounds; round++ {
		if budgetExceeded(req, startedAt, clock.Now()) {
			return Result{Outcome: store.StateBudgetExceeded, Rounds: round - 1}, nil
		}
		settledWait, err := waitForSettled(ctx, req, currentHeadSHA, runDeadline, deps.StopRequests, deps.ReviewEvidence, deps.StatusSource, clock, sleeper, publisher)
		if err != nil {
			return resultForError(round-1, err), err
		}
		terminalEvidence = settledWait.evidence
		if settledWait.budgetExceeded {
			return Result{Outcome: store.StateBudgetExceeded, Rounds: round - 1}, nil
		}
		status := settledWait.status
		if status.State != StatusSettled {
			return resultForTimedOut(round-1, settledWait.evidence), nil
		}
		if settledWait.evidence.State == reviewsource.EvidenceSkipped {
			return resultForReviewSkipped(round-1, settledWait.evidence), nil
		}
		settledBeforeRun := round == 1 && settledWait.statusChecks == 1
		if !settledBeforeRun {
			if req.QuietPeriod > 0 {
				if err := publisher.publish(ctx, runevent.KindDaemonQuietPeriod,
					fmt.Sprintf("Quiet period: waiting %s before fetching Round %03d.", req.QuietPeriod, round),
					map[string]any{"seconds": req.QuietPeriod.Seconds(), "round": round},
				); err != nil {
					return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
				}
			}
			if err := sleeper.Sleep(ctx, req.QuietPeriod); err != nil {
				return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
			}
			if err := observeStopRequest(ctx, deps.StopRequests, req.RunID, "after quiet-period wait"); err != nil {
				return resultForError(round-1, err), err
			}
		}
		if budgetExceeded(req, startedAt, clock.Now()) {
			return Result{Outcome: store.StateBudgetExceeded, Rounds: round - 1}, nil
		}

		if err := publisher.publish(ctx, runevent.KindDaemonFetch,
			fmt.Sprintf("Fetching Round %03d from the Review Source.", round),
			map[string]any{"phase": "started", "round": round},
		); err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		fetched, err := deps.Fetcher.Fetch(ctx, round)
		if err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		reviewIssuesKnown = true
		if err := publisher.publish(ctx, runevent.KindDaemonFetch,
			fmt.Sprintf("Fetched Round %03d with %d Review Issue(s).", fetched.Round, fetched.Issues),
			map[string]any{"phase": "completed", "round": fetched.Round, "issues": fetched.Issues},
		); err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		if fetched.Issues == 0 {
			confirm, confirmedHeadSHA, err := confirmMergeReadyWithArtifacts(ctx, req, fetched.Round, runDeadline, deps, currentHeadSHA, clock, sleeper, publisher)
			if err != nil {
				return resultForError(round, err), err
			}
			if confirm.evidence.State != "" {
				terminalEvidence = confirm.evidence
			}
			if confirm.budgetExceeded {
				return Result{Outcome: store.StateBudgetExceeded, Rounds: round}, nil
			}
			if confirm.skipped {
				return resultForReviewSkipped(round, confirm.evidence), nil
			}
			if confirm.ready {
				verifiedHeadSHA = confirmedHeadSHA
				return Result{Outcome: store.StateClean, Rounds: round}, nil
			}
			if confirm.unverified {
				return Result{Outcome: store.StateCleanUnverified, Rounds: round}, nil
			}
			if confirm.timedOut {
				return resultForTimedOut(round, confirm.evidence), nil
			}
			if round == req.MaxRounds {
				return Result{Outcome: store.StateMaxRoundsReached, Rounds: round}, nil
			}
			continue
		}

		resolved, err := deps.Resolver.Resolve(ctx)
		if err != nil {
			return resultForError(round, err), err
		}
		if resolved.Outcome != "" {
			return Result{Outcome: resolved.Outcome, Rounds: round, Remaining: resolved.Remaining}, nil
		}
		if err := requestReviewForResolvedHead(ctx, req, deps.ReviewRequester, currentHeadSHA, resolved.HeadSHA); err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round, Remaining: resolved.Remaining}, err
		}
		if resolved.HeadSHA != "" {
			currentHeadSHA = resolved.HeadSHA
		}
		if resolved.Remaining == 0 {
			confirm, confirmedHeadSHA, err := confirmMergeReadyWithArtifacts(ctx, req, fetched.Round, runDeadline, deps, currentHeadSHA, clock, sleeper, publisher)
			if err != nil {
				return resultForError(round, err), err
			}
			if confirm.evidence.State != "" {
				terminalEvidence = confirm.evidence
			}
			if confirm.budgetExceeded {
				return Result{Outcome: store.StateBudgetExceeded, Rounds: round}, nil
			}
			if confirm.skipped {
				return resultForReviewSkipped(round, confirm.evidence), nil
			}
			if confirm.ready {
				verifiedHeadSHA = confirmedHeadSHA
				return Result{Outcome: store.StateClean, Rounds: round}, nil
			}
			if confirm.unverified {
				return Result{Outcome: store.StateCleanUnverified, Rounds: round}, nil
			}
			if confirm.timedOut {
				return resultForTimedOut(round, confirm.evidence), nil
			}
			if round == req.MaxRounds {
				return Result{Outcome: store.StateMaxRoundsReached, Rounds: round}, nil
			}
			continue
		}
		if !resolved.Progress {
			// A Round that settles nothing will not improve by repeating:
			// end the Run as Unresolved instead of burning more Rounds.
			return Result{Outcome: store.StateUnresolved, Rounds: round, Remaining: resolved.Remaining}, nil
		}
		if !req.UntilClean {
			return Result{Outcome: store.StateMaxRoundsReached, Rounds: round, Remaining: resolved.Remaining}, nil
		}
		if round == req.MaxRounds {
			return Result{Outcome: store.StateMaxRoundsReached, Rounds: round, Remaining: resolved.Remaining}, nil
		}
	}
	return Result{Outcome: store.StateMaxRoundsReached, Rounds: req.MaxRounds}, nil
}

func requestReviewForResolvedHead(ctx context.Context, req Request, requester reviewsource.ReviewRequester, currentHeadSHA string, resolvedHeadSHA string) error {
	resolvedHeadSHA = strings.TrimSpace(resolvedHeadSHA)
	if !req.RequestReview || requester == nil || resolvedHeadSHA == "" || resolvedHeadSHA == strings.TrimSpace(currentHeadSHA) {
		return nil
	}
	if _, err := requester.RequestReview(ctx, reviewsource.ReviewRequest{
		RunID:          req.RunID,
		BaseRepository: req.BaseRepository,
		PRNumber:       req.PRNumber,
		HeadSHA:        resolvedHeadSHA,
		Command:        req.ReviewCommand,
	}); err != nil {
		return fmt.Errorf("request Review Source review for pushed head: %w", err)
	}
	return nil
}

type settledWaitResult struct {
	status         Status
	statusChecks   int
	evidence       reviewsource.Evidence
	budgetExceeded bool
}

func waitForSettled(ctx context.Context, req Request, headSHA string, runDeadline time.Time, stops StopRequestSource, evidenceSource ReviewEvidenceSource, statusSource StatusSource, clock Clock, sleeper Sleeper, publisher *watchEventPublisher) (settledWaitResult, error) {
	startedAt := clock.Now()
	deadline := boundedDeadline(startedAt.Add(req.ReviewTimeout), runDeadline)
	statusChecks := 0
	retry := retryEpisode{}
	progress := publisher.beginWait(WaitPhaseReview, headSHA, startedAt, deadline)
	if err := publisher.publishWait(ctx, progress); err != nil {
		return settledWaitResult{}, err
	}
	for {
		if err := ctx.Err(); err != nil {
			return settledWaitResult{}, err
		}
		if err := observeStopRequest(ctx, stops, req.RunID, "before Review Source status"); err != nil {
			return settledWaitResult{}, err
		}
		if !clock.Now().Before(deadline) {
			if err := retry.exhaust(ctx, publisher, &progress); err != nil {
				return settledWaitResult{}, err
			}
			return settledWaitResult{
				status:         Status{State: store.StateTimedOut, Detail: progress.Evidence.Detail},
				statusChecks:   statusChecks,
				evidence:       progress.Evidence,
				budgetExceeded: deadlineUsesRunBudget(deadline, runDeadline),
			}, nil
		}
		var status Status
		var evidence reviewsource.Evidence
		var err error
		if evidenceSource != nil {
			evidence, err = evidenceSource.Evidence(ctx, ReviewEvidenceRequest{
				PRNumber:        req.PRNumber,
				ExpectedHeadSHA: headSHA,
			})
			if err == nil {
				status = statusFromEvidence(evidence)
			}
		} else {
			status, err = statusSource.Status(ctx, StatusRequest{
				PRNumber: req.PRNumber,
				HeadSHA:  headSHA,
			})
			if err == nil {
				evidence = evidenceFromStatus(status, headSHA)
			}
		}
		statusChecks++
		if err != nil {
			if stopErr := observeStopRequest(ctx, stops, req.RunID, "after failed Review Source status access"); stopErr != nil {
				return settledWaitResult{}, stopErr
			}
			if ctxErr := ctx.Err(); ctxErr != nil {
				return settledWaitResult{}, ctxErr
			}
			if !reviewsource.IsTransient(err) {
				return settledWaitResult{}, err
			}
			if startErr := retry.start(ctx, err, publisher, &progress); startErr != nil {
				return settledWaitResult{}, startErr
			}
		} else {
			if err := observeStopRequest(ctx, stops, req.RunID, "after Review Source status access"); err != nil {
				return settledWaitResult{}, err
			}
			if recoverErr := retry.recover(ctx, publisher, &progress); recoverErr != nil {
				return settledWaitResult{}, recoverErr
			}
			progress.Evidence = normalizedEvidence(evidence, headSHA)
			if err := publisher.publishWait(ctx, progress); err != nil {
				return settledWaitResult{}, err
			}
			if status.State == StatusSettled {
				return settledWaitResult{status: status, statusChecks: statusChecks, evidence: evidence}, nil
			}
		}
		if !clock.Now().Before(deadline) {
			if err := retry.exhaust(ctx, publisher, &progress); err != nil {
				return settledWaitResult{}, err
			}
			return settledWaitResult{
				status:         Status{State: store.StateTimedOut, Detail: progress.Evidence.Detail},
				statusChecks:   statusChecks,
				evidence:       progress.Evidence,
				budgetExceeded: deadlineUsesRunBudget(deadline, runDeadline),
			}, nil
		}
		if err := sleeper.Sleep(ctx, nextPollDelay(clock.Now(), deadline, req.PollInterval)); err != nil {
			return settledWaitResult{}, err
		}
		if err := observeStopRequest(ctx, stops, req.RunID, "after Review Source status wait"); err != nil {
			return settledWaitResult{}, err
		}
	}
}

func reviewSkippedReason(evidence reviewsource.Evidence) string {
	if reason := strings.TrimSpace(evidence.Reason); reason != "" {
		return reviewsource.BoundEvidenceDetail(reason)
	}
	if detail := strings.TrimSpace(evidence.Detail); detail != "" {
		return reviewsource.BoundEvidenceDetail(detail)
	}
	return "Review Source explicitly skipped the review."
}

func resultForReviewSkipped(rounds int, evidence reviewsource.Evidence) Result {
	return Result{
		Outcome:        store.StateReviewSkipped,
		Rounds:         rounds,
		TerminalReason: reviewSkippedReason(evidence),
		NextAction:     reviewSkippedNextAction,
		Evidence:       evidence,
	}
}

func resultForTimedOut(rounds int, evidence reviewsource.Evidence) Result {
	terminalReason := unrecognisedEvidenceReason(evidence)
	return Result{
		Outcome:                    store.StateTimedOut,
		Rounds:                     rounds,
		ManualReviewCommand:        "@coderabbitai review",
		TerminalReason:             terminalReason,
		TerminalReasonIsDiagnostic: terminalReason != "",
		Evidence:                   evidence,
	}
}

func unrecognisedEvidenceReason(evidence reviewsource.Evidence) string {
	if evidence.State != reviewsource.EvidencePending || evidence.Kind == reviewsource.EvidenceKindNone {
		return ""
	}
	detail := strings.TrimSpace(evidence.Detail)
	if detail == "" {
		return "Review Source signal was not recognised."
	}
	return reviewsource.BoundEvidenceDetail("Review Source signal was not recognised: " + detail)
}

type confirmResult struct {
	ready          bool
	unverified     bool
	timedOut       bool
	budgetExceeded bool
	skipped        bool
	evidence       reviewsource.Evidence
}

func confirmMergeReadyWithArtifacts(
	ctx context.Context,
	req Request,
	round int,
	runDeadline time.Time,
	deps Dependencies,
	headSHA string,
	clock Clock,
	sleeper Sleeper,
	publisher *watchEventPublisher,
) (confirmResult, string, error) {
	artifactPublished := false
	for {
		confirm, err := confirmMergeReady(
			ctx,
			req,
			runDeadline,
			deps.StopRequests,
			deps.ReviewEvidence,
			deps.CheckSource,
			headSHA,
			clock,
			sleeper,
			publisher,
		)
		if err != nil || !confirm.ready {
			return confirm, "", err
		}
		verifiedHeadSHA := verifiedHeadFromEvidence(confirm.evidence, headSHA)
		if artifactPublished ||
			deps.Artifacts == nil ||
			!acceptedVerifiedEvidenceForHead(confirm.evidence, headSHA) {
			return confirm, verifiedHeadSHA, nil
		}
		if deps.WriteGuard != nil {
			if err := deps.WriteGuard.Check(ctx); err != nil {
				return confirmResult{}, "", fmt.Errorf("guard review artifact commit: %w", err)
			}
		}

		publication, err := deps.Artifacts.PublishArtifacts(ctx, ArtifactPublishRequest{
			Round:          round,
			ParentHeadSHA:  headSHA,
			ParentEvidence: confirm.evidence,
		})
		if err != nil {
			return confirmResult{}, "", err
		}
		if publication.Commit.CommitSHA == "" {
			return confirm, verifiedHeadSHA, nil
		}
		artifactPublished = true
		if acceptedArtifactEvidence(publication) {
			if err := publisher.publishAcceptedEvidence(ctx, publication.Evidence); err != nil {
				return confirmResult{}, "", err
			}
			return confirmResult{ready: true, evidence: publication.Evidence}, verifiedHeadSHA, nil
		}

		// A created commit that cannot inherit falls through to the ordinary
		// current-head Evidence poll. The publisher is not called again.
		headSHA = publication.Commit.CommitSHA
	}
}

func acceptedVerifiedEvidenceForHead(evidence reviewsource.Evidence, headSHA string) bool {
	return evidence.State == reviewsource.EvidenceVerified &&
		evidence.ExpectedHeadSHA == headSHA &&
		evidence.ObservedHeadSHA == headSHA
}

func acceptedArtifactEvidence(publication ArtifactPublication) bool {
	commit := publication.Commit
	evidence := publication.Evidence
	return commit.CommitSHA != "" &&
		commit.ParentSHA != "" &&
		evidence.State == reviewsource.EvidenceVerified &&
		evidence.Kind == reviewsource.EvidenceKindArtifactOnlyDescendant &&
		evidence.ExpectedHeadSHA == commit.CommitSHA &&
		evidence.ObservedHeadSHA == commit.CommitSHA &&
		evidence.ParentHeadSHA == commit.ParentSHA
}

func confirmMergeReady(ctx context.Context, req Request, runDeadline time.Time, stops StopRequestSource, evidenceSource ReviewEvidenceSource, checkSource CheckSource, headSHA string, clock Clock, sleeper Sleeper, publisher *watchEventPublisher) (confirmResult, error) {
	if !req.UntilClean || evidenceSource == nil && checkSource == nil {
		return confirmResult{ready: true}, nil
	}
	startedAt := clock.Now()
	deadline := boundedDeadline(startedAt.Add(req.ReviewTimeout), runDeadline)
	retry := retryEpisode{}
	progress := publisher.beginWait(WaitPhaseReviewCheck, headSHA, startedAt, deadline)
	if err := publisher.publishWait(ctx, progress); err != nil {
		return confirmResult{}, err
	}
	for {
		if err := ctx.Err(); err != nil {
			return confirmResult{}, err
		}
		if err := observeStopRequest(ctx, stops, req.RunID, "before Review Source Merge-Ready check"); err != nil {
			return confirmResult{}, err
		}
		if !clock.Now().Before(deadline) {
			wasRetrying := retry.active
			if err := retry.exhaust(ctx, publisher, &progress); err != nil {
				return confirmResult{}, err
			}
			if deadlineUsesRunBudget(deadline, runDeadline) {
				return confirmResult{budgetExceeded: true, evidence: progress.Evidence}, nil
			}
			if !wasRetrying &&
				progress.Evidence.State == reviewsource.EvidencePending &&
				progress.Evidence.Kind == reviewsource.EvidenceKindNone {
				return confirmResult{unverified: true, evidence: progress.Evidence}, nil
			}
			return confirmResult{timedOut: true, evidence: progress.Evidence}, nil
		}
		missingCheck := false
		state := CheckMissing
		var evidence reviewsource.Evidence
		var err error
		if evidenceSource != nil {
			evidence, err = evidenceSource.Evidence(ctx, ReviewEvidenceRequest{
				PRNumber:        req.PRNumber,
				ExpectedHeadSHA: headSHA,
			})
			if err == nil {
				state = headCheckFromEvidence(evidence)
			}
		} else {
			state, err = checkSource.Check(ctx, headSHA)
			if err == nil {
				evidence = evidenceFromHeadCheck(state, headSHA)
			}
		}
		if err != nil {
			if stopErr := observeStopRequest(ctx, stops, req.RunID, "after failed Review Source Merge-Ready check"); stopErr != nil {
				return confirmResult{}, stopErr
			}
			if ctxErr := ctx.Err(); ctxErr != nil {
				return confirmResult{}, ctxErr
			}
			if !reviewsource.IsTransient(err) {
				return confirmResult{}, err
			}
			deadline = boundedDeadline(startedAt.Add(req.ReviewTimeout), runDeadline)
			progress.Deadline = deadline
			if startErr := retry.start(ctx, err, publisher, &progress); startErr != nil {
				return confirmResult{}, startErr
			}
		} else {
			if err := observeStopRequest(ctx, stops, req.RunID, "after Review Source Merge-Ready check"); err != nil {
				return confirmResult{}, err
			}
			if recoverErr := retry.recover(ctx, publisher, &progress); recoverErr != nil {
				return confirmResult{}, recoverErr
			}
			evidence = normalizedEvidence(evidence, headSHA)
			if evidence.State == reviewsource.EvidencePending && evidence.Kind == reviewsource.EvidenceKindNone {
				deadline = boundedDeadline(startedAt.Add(req.CheckGracePeriod), runDeadline)
			} else {
				deadline = boundedDeadline(startedAt.Add(req.ReviewTimeout), runDeadline)
			}
			progress.Deadline = deadline
			progress.Evidence = evidence
			if err := publisher.publishWait(ctx, progress); err != nil {
				return confirmResult{}, err
			}
			if evidence.State == reviewsource.EvidenceSkipped {
				return confirmResult{skipped: true, evidence: evidence}, nil
			}
			switch state {
			case CheckSuccess:
				return confirmResult{ready: true, evidence: evidence}, nil
			case CheckMissing:
				missingCheck = true
			case CheckFailure:
				return confirmResult{}, nil
			case CheckPending:
			default:
				return confirmResult{}, fmt.Errorf("unknown Review Source check state %q", state)
			}
		}
		if !clock.Now().Before(deadline) {
			wasRetrying := retry.active
			if err := retry.exhaust(ctx, publisher, &progress); err != nil {
				return confirmResult{}, err
			}
			if deadlineUsesRunBudget(deadline, runDeadline) {
				return confirmResult{budgetExceeded: true, evidence: progress.Evidence}, nil
			}
			if !wasRetrying && missingCheck {
				return confirmResult{unverified: true, evidence: progress.Evidence}, nil
			}
			return confirmResult{timedOut: true, evidence: progress.Evidence}, nil
		}
		if err := sleeper.Sleep(ctx, nextPollDelay(clock.Now(), deadline, req.PollInterval)); err != nil {
			return confirmResult{}, err
		}
		operation := "after Merge-Ready wait"
		if retry.active {
			operation = "after transient Review Source retry wait"
		}
		if err := observeStopRequest(ctx, stops, req.RunID, operation); err != nil {
			return confirmResult{}, err
		}
	}
}

func observeStopRequest(ctx context.Context, source StopRequestSource, runID string, operation string) error {
	if source == nil {
		return nil
	}
	requested, err := source.StopRequested(ctx, runID)
	if err != nil {
		return fmt.Errorf("observe Stop Request for Run %q %s: %w", runID, operation, err)
	}
	if requested {
		return ErrStopRequested
	}
	return nil
}

func resultForError(rounds int, err error) Result {
	outcome := store.StateFailed
	var checkoutMoved CheckoutMovedError
	switch {
	case errors.Is(err, ErrStopRequested):
		outcome = store.StateStopped
	case errors.As(err, &checkoutMoved):
		outcome = store.StateCheckoutMoved
	}
	return Result{Outcome: outcome, Rounds: rounds}
}

func finalizeResult(result *Result, resultErr error, reviewIssuesKnown bool, evidence reviewsource.Evidence, verifiedHeadSHA string) {
	if result == nil {
		return
	}
	result.ReviewIssuesKnown = reviewIssuesKnown
	if result.Evidence.State == "" {
		result.Evidence = evidence
	}
	if result.VerifiedHeadSHA == "" {
		result.VerifiedHeadSHA = verifiedHeadSHA
	}
	if result.Outcome == "" || result.Outcome == store.StateClean {
		return
	}
	if result.TerminalReason == "" && resultErr != nil {
		result.TerminalReason = boundedTerminalText(resultErr.Error())
	}
	defaultReason, defaultAction := terminalDefaults(result.Outcome)
	if result.TerminalReason == "" {
		result.TerminalReason = defaultReason
	}
	if result.NextAction == "" {
		result.NextAction = defaultAction
	}
	result.TerminalReason = boundedTerminalText(result.TerminalReason)
	result.NextAction = boundedTerminalText(result.NextAction)
}

func terminalDefaults(outcome string) (string, string) {
	switch outcome {
	case store.StateCleanUnverified:
		return "Merge-Ready was not confirmed for the fetched head.", "Confirm the pull request's Review Source Evidence before merging."
	case store.StateReviewSkipped:
		return "The Review Source explicitly skipped the review.", reviewSkippedNextAction
	case store.StateMaxRoundsReached:
		return "The configured maximum number of Rounds was reached.", "Review the remaining Review Issues before deciding whether to start another Run."
	case store.StateBudgetExceeded:
		return "The Run Budget was exhausted.", "Inspect the Run Event Stream before starting another Run with an appropriate budget."
	case store.StateTimedOut:
		return "Review Source Evidence did not arrive before the timeout.", "Request another Review Source review, then start another watch Run."
	case store.StateUnresolved:
		return "The last Round settled no Review Issues.", "Review the remaining Review Issues and address them before starting another Run."
	case store.StateStopped:
		return "A Stop Request ended the Run.", "Inspect the preserved work before starting another Run."
	case store.StateCheckoutMoved:
		return "The checkout no longer matches the target recorded at Preflight.", "Restore the PR Head Branch checkout, then start another Run."
	case store.StateFailed:
		return "The watch Run failed before it could complete.", "Inspect the diagnostics, correct the failure, and start another Run."
	default:
		return "The Run reached a non-Clean outcome.", "Inspect the Run Event Stream before deciding the next recovery step."
	}
}

func boundedTerminalText(text string) string {
	return reviewsource.BoundEvidenceDetail(strings.Join(strings.Fields(text), " "))
}

func verifiedHeadFromEvidence(evidence reviewsource.Evidence, fallback string) string {
	if evidence.State != reviewsource.EvidenceVerified {
		return ""
	}
	if evidence.ObservedHeadSHA != "" {
		return evidence.ObservedHeadSHA
	}
	if evidence.ExpectedHeadSHA != "" {
		return evidence.ExpectedHeadSHA
	}
	return fallback
}

func validateRequest(req Request, deps Dependencies) error {
	if req.PRNumber == "" {
		return errors.New("watch requires Open Pull Request number")
	}
	if req.HeadSHA == "" {
		return errors.New("watch requires HEAD SHA")
	}
	if req.MaxRounds < 1 {
		return errors.New("watch max rounds must be greater than 0")
	}
	if req.PollInterval <= 0 {
		return errors.New("watch poll interval must be greater than 0")
	}
	if req.QuietPeriod < 0 {
		return errors.New("watch quiet period must not be negative")
	}
	if req.ReviewTimeout <= 0 {
		return errors.New("watch review timeout must be greater than 0")
	}
	if req.UntilClean && (deps.ReviewEvidence != nil || deps.CheckSource != nil) && req.CheckGracePeriod <= 0 {
		return errors.New("watch check grace period must be greater than 0")
	}
	if req.BudgetEnabled && req.MaxRunDuration <= 0 {
		return errors.New("watch max run duration must be greater than 0 when Run Budget is enabled")
	}
	if deps.ReviewEvidence == nil && deps.StatusSource == nil {
		return errors.New("watch requires Review Source evidence boundary")
	}
	if deps.Fetcher == nil {
		return errors.New("watch requires fetch boundary")
	}
	if deps.Resolver == nil {
		return errors.New("watch requires resolve boundary")
	}
	return nil
}

func budgetExceeded(req Request, startedAt time.Time, now time.Time) bool {
	return req.BudgetEnabled && !now.Before(startedAt.Add(req.MaxRunDuration))
}

func boundedDeadline(phaseDeadline time.Time, runDeadline time.Time) time.Time {
	if !runDeadline.IsZero() && runDeadline.Before(phaseDeadline) {
		return runDeadline
	}
	return phaseDeadline
}

func deadlineUsesRunBudget(deadline time.Time, runDeadline time.Time) bool {
	return !runDeadline.IsZero() && deadline.Equal(runDeadline)
}

func nextPollDelay(now time.Time, deadline time.Time, pollInterval time.Duration) time.Duration {
	remaining := deadline.Sub(now)
	if remaining < pollInterval {
		return remaining
	}
	return pollInterval
}

type retryEpisode struct {
	active    bool
	operation string
	reason    string
}

func (episode *retryEpisode) start(ctx context.Context, err error, publisher *watchEventPublisher, progress *WaitProgress) error {
	if episode.active {
		return nil
	}
	var transient *reviewsource.TransientError
	if !errors.As(err, &transient) {
		return fmt.Errorf("start Review Source retry episode for permanent failure: %w", err)
	}
	episode.active = true
	episode.operation = reviewsource.BoundEvidenceDetail(strings.TrimSpace(transient.Operation))
	if episode.operation == "" {
		episode.operation = "access Review Source evidence"
	}
	episode.reason = reviewsource.BoundEvidenceDetail(err.Error())
	if err := publisher.publishRetry(ctx, "started", episode.operation, episode.reason); err != nil {
		return err
	}
	progress.RetryStatus = RetryStatusRetrying
	return publisher.publishWait(ctx, *progress)
}

func (episode *retryEpisode) recover(ctx context.Context, publisher *watchEventPublisher, progress *WaitProgress) error {
	if !episode.active {
		return nil
	}
	if err := publisher.publishRetry(ctx, "recovered", episode.operation, episode.reason); err != nil {
		return err
	}
	episode.active = false
	progress.RetryStatus = RetryStatusRecovered
	return nil
}

func (episode *retryEpisode) exhaust(ctx context.Context, publisher *watchEventPublisher, progress *WaitProgress) error {
	if !episode.active {
		return nil
	}
	if err := publisher.publishRetry(ctx, "exhausted", episode.operation, episode.reason); err != nil {
		return err
	}
	episode.active = false
	progress.RetryStatus = RetryStatusExhausted
	return publisher.publishWait(ctx, *progress)
}

// watchEventPublisher appends watch-loop Run Events: status waits, quiet
// periods, and fetch results. Publication is part of the Run state
// contract, so a critical sink failure fails the Run.
type watchEventPublisher struct {
	sink           runevent.Sink
	runID          string
	clock          Clock
	progress       func(WaitProgress)
	latestEvidence reviewsource.Evidence
	lastWait       *WaitProgress
}

func (publisher watchEventPublisher) publish(ctx context.Context, kind runevent.Kind, summary string, payload any) error {
	if publisher.runID == "" {
		return nil
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode watch event payload: %w", err)
	}
	if err := publisher.sink.Publish(ctx, runevent.RunEvent{
		RunID:   publisher.runID,
		Source:  runevent.SourceDaemon,
		Kind:    kind,
		Summary: runevent.BoundSummary(summary),
		Time:    publisher.clock.Now(),
		Payload: raw,
	}); err != nil {
		return fmt.Errorf("publish watch event %s: %w", kind, err)
	}
	return nil
}

func (publisher *watchEventPublisher) beginWait(phase string, headSHA string, startedAt time.Time, deadline time.Time) WaitProgress {
	evidence := publisher.latestEvidence
	if evidence.ExpectedHeadSHA != "" && evidence.ExpectedHeadSHA != headSHA {
		evidence = reviewsource.Evidence{}
	}
	evidence = normalizedEvidence(evidence, headSHA)
	return WaitProgress{
		Phase:           phase,
		ExpectedHeadSHA: headSHA,
		StartedAt:       startedAt,
		Deadline:        deadline,
		Evidence:        evidence,
		RetryStatus:     RetryStatusNone,
	}
}

func (publisher *watchEventPublisher) publishWait(ctx context.Context, progress WaitProgress) error {
	progress.Evidence = normalizedEvidence(progress.Evidence, progress.ExpectedHeadSHA)
	if publisher.lastWait != nil && *publisher.lastWait == progress {
		return nil
	}
	payload := runevent.ReviewStatusPayload{
		Phase:           progress.Phase,
		StartedAt:       progress.StartedAt,
		Deadline:        progress.Deadline,
		EvidenceState:   string(progress.Evidence.State),
		EvidenceKind:    string(progress.Evidence.Kind),
		RetryStatus:     progress.RetryStatus,
		State:           string(progress.Evidence.State),
		Kind:            string(progress.Evidence.Kind),
		Identity:        progress.Evidence.Identity,
		ExpectedHeadSHA: progress.ExpectedHeadSHA,
		ObservedHeadSHA: progress.Evidence.ObservedHeadSHA,
		ParentHeadSHA:   progress.Evidence.ParentHeadSHA,
		Conclusion:      progress.Evidence.Conclusion,
		Detail:          progress.Evidence.Detail,
		Reason:          progress.Evidence.Reason,
	}
	if err := publisher.publish(
		ctx,
		runevent.KindDaemonReviewStatus,
		fmt.Sprintf("Review Source Evidence: %s (%s); phase %s; retry %s.",
			progress.Evidence.State,
			progress.Evidence.Kind,
			progress.Phase,
			progress.RetryStatus,
		),
		payload,
	); err != nil {
		return err
	}
	publisher.latestEvidence = progress.Evidence
	current := progress
	publisher.lastWait = &current
	if publisher.progress != nil {
		publisher.progress(progress)
	}
	return nil
}

func (publisher *watchEventPublisher) publishAcceptedEvidence(ctx context.Context, evidence reviewsource.Evidence) error {
	payload := runevent.ReviewStatusPayload{
		EvidenceState:   string(evidence.State),
		EvidenceKind:    string(evidence.Kind),
		State:           string(evidence.State),
		Kind:            string(evidence.Kind),
		Identity:        evidence.Identity,
		ExpectedHeadSHA: evidence.ExpectedHeadSHA,
		ObservedHeadSHA: evidence.ObservedHeadSHA,
		ParentHeadSHA:   evidence.ParentHeadSHA,
		Conclusion:      evidence.Conclusion,
		Detail:          evidence.Detail,
		Reason:          evidence.Reason,
	}
	if err := publisher.publish(
		ctx,
		runevent.KindDaemonReviewStatus,
		fmt.Sprintf("Review Source Evidence: %s (%s).", evidence.State, evidence.Kind),
		payload,
	); err != nil {
		return err
	}
	publisher.latestEvidence = evidence
	return nil
}

func (publisher *watchEventPublisher) publishRetry(ctx context.Context, phase string, operation string, reason string) error {
	return publisher.publish(
		ctx,
		runevent.KindDaemonRetry,
		fmt.Sprintf("Review Source retry %s: %s.", phase, operation),
		runevent.RetryPayload{
			Phase:     phase,
			Operation: operation,
			Reason:    reason,
		},
	)
}

func normalizedEvidence(evidence reviewsource.Evidence, headSHA string) reviewsource.Evidence {
	if evidence.State == "" {
		evidence.State = reviewsource.EvidencePending
	}
	if evidence.Kind == "" {
		evidence.Kind = reviewsource.EvidenceKindNone
	}
	if evidence.ExpectedHeadSHA == "" {
		evidence.ExpectedHeadSHA = headSHA
	}
	return evidence
}

func evidenceFromStatus(status Status, headSHA string) reviewsource.Evidence {
	state := reviewsource.EvidencePending
	switch status.State {
	case StatusReviewing:
		state = reviewsource.EvidenceReviewing
	case StatusSettled:
		state = reviewsource.EvidenceReviewed
	}
	return reviewsource.Evidence{
		State:           state,
		Kind:            reviewsource.EvidenceKindNone,
		ExpectedHeadSHA: headSHA,
		Detail:          status.Detail,
	}
}

func evidenceFromHeadCheck(state HeadCheckState, headSHA string) reviewsource.Evidence {
	evidenceState := reviewsource.EvidencePending
	switch state {
	case CheckPending:
		evidenceState = reviewsource.EvidenceReviewing
	case CheckSuccess:
		evidenceState = reviewsource.EvidenceVerified
	case CheckFailure:
		evidenceState = reviewsource.EvidenceFailed
	}
	return reviewsource.Evidence{
		State:           evidenceState,
		Kind:            reviewsource.EvidenceKindNone,
		ExpectedHeadSHA: headSHA,
	}
}

func statusFromEvidence(evidence reviewsource.Evidence) Status {
	state := StatusPending
	switch evidence.State {
	case reviewsource.EvidenceReviewing:
		state = StatusReviewing
	case reviewsource.EvidenceReviewed, reviewsource.EvidenceVerified, reviewsource.EvidenceSkipped, reviewsource.EvidenceFailed:
		state = StatusSettled
	}
	return Status{State: state, Detail: evidence.Detail}
}

func headCheckFromEvidence(evidence reviewsource.Evidence) HeadCheckState {
	switch evidence.State {
	case reviewsource.EvidenceVerified:
		return CheckSuccess
	case reviewsource.EvidencePending:
		if evidence.Kind == reviewsource.EvidenceKindNone {
			return CheckMissing
		}
		return CheckPending
	case reviewsource.EvidenceReviewing:
		return CheckPending
	default:
		return CheckFailure
	}
}
