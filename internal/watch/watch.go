package watch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

const (
	StatusPending   = "pending"
	StatusReviewing = "reviewing"
	StatusSettled   = "settled"
)

type HeadCheckState string

const (
	CheckPending HeadCheckState = "pending"
	CheckSuccess HeadCheckState = "success"
	CheckFailure HeadCheckState = "failure"
	CheckMissing HeadCheckState = "missing"
)

type Request struct {
	RunID          string
	PRNumber       string
	HeadSHA        string
	UntilClean     bool
	MaxRounds      int
	PollInterval   time.Duration
	QuietPeriod    time.Duration
	ReviewTimeout  time.Duration
	BudgetEnabled  bool
	MaxRunDuration time.Duration
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

type Result struct {
	Outcome             string
	Rounds              int
	Remaining           int
	ManualReviewCommand string
	CheckMissing        bool
}

type StatusSource interface {
	Status(context.Context, StatusRequest) (Status, error)
}

type StatusFunc func(context.Context, StatusRequest) (Status, error)

func (fn StatusFunc) Status(ctx context.Context, req StatusRequest) (Status, error) {
	return fn(ctx, req)
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
	StatusSource StatusSource
	Fetcher      Fetcher
	Resolver     Resolver
	CheckSource  CheckSource
	Clock        Clock
	Sleeper      Sleeper
	// Sink receives watch-loop Run Events: review status waits, quiet
	// periods, fetch results, and merge-readiness checks. Nil means
	// events are discarded.
	Sink runevent.Sink
}

func Run(ctx context.Context, req Request, deps Dependencies) (Result, error) {
	if err := validateRequest(req, deps); err != nil {
		return Result{}, err
	}
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
	publisher := watchEventPublisher{sink: deps.Sink, runID: req.RunID, clock: clock}

	startedAt := clock.Now()
	currentHeadSHA := req.HeadSHA
	for round := 1; round <= req.MaxRounds; round++ {
		if budgetExceeded(req, startedAt, clock.Now()) {
			return Result{Outcome: store.StateBudgetExceeded, Rounds: round - 1}, nil
		}
		settledWait, err := waitForSettled(ctx, req, currentHeadSHA, deps.StatusSource, clock, sleeper, publisher)
		if err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		status := settledWait.status
		if status.State != StatusSettled {
			return Result{
				Outcome:             store.StateTimedOut,
				Rounds:              round - 1,
				ManualReviewCommand: "@coderabbitai review",
			}, nil
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
		if err := publisher.publish(ctx, runevent.KindDaemonFetch,
			fmt.Sprintf("Fetched Round %03d with %d Review Issue(s).", fetched.Round, fetched.Issues),
			map[string]any{"phase": "completed", "round": fetched.Round, "issues": fetched.Issues},
		); err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		if fetched.Issues == 0 {
			confirm, err := confirmMergeReady(ctx, req, deps.CheckSource, currentHeadSHA, clock, sleeper, publisher)
			if err != nil {
				return Result{Outcome: store.StateFailed, Rounds: round}, err
			}
			if confirm.ready {
				return Result{Outcome: store.StateClean, Rounds: round, CheckMissing: confirm.missing}, nil
			}
			if confirm.timedOut {
				return Result{
					Outcome:             store.StateTimedOut,
					Rounds:              round,
					ManualReviewCommand: "@coderabbitai review",
				}, nil
			}
			if round == req.MaxRounds {
				return Result{Outcome: store.StateMaxRoundsReached, Rounds: round}, nil
			}
			continue
		}

		resolved, err := deps.Resolver.Resolve(ctx)
		if err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round}, err
		}
		if resolved.Outcome != "" {
			return Result{Outcome: resolved.Outcome, Rounds: round, Remaining: resolved.Remaining}, nil
		}
		if resolved.HeadSHA != "" {
			currentHeadSHA = resolved.HeadSHA
		}
		if resolved.Remaining == 0 {
			confirm, err := confirmMergeReady(ctx, req, deps.CheckSource, currentHeadSHA, clock, sleeper, publisher)
			if err != nil {
				return Result{Outcome: store.StateFailed, Rounds: round}, err
			}
			if confirm.ready {
				return Result{Outcome: store.StateClean, Rounds: round, CheckMissing: confirm.missing}, nil
			}
			if confirm.timedOut {
				return Result{
					Outcome:             store.StateTimedOut,
					Rounds:              round,
					ManualReviewCommand: "@coderabbitai review",
				}, nil
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

type settledWaitResult struct {
	status       Status
	statusChecks int
}

func waitForSettled(ctx context.Context, req Request, headSHA string, source StatusSource, clock Clock, sleeper Sleeper, publisher watchEventPublisher) (settledWaitResult, error) {
	startedAt := clock.Now()
	statusChecks := 0
	for {
		if err := ctx.Err(); err != nil {
			return settledWaitResult{}, err
		}
		status, err := source.Status(ctx, StatusRequest{
			PRNumber: req.PRNumber,
			HeadSHA:  headSHA,
		})
		statusChecks++
		if err != nil {
			return settledWaitResult{}, err
		}
		if err := publisher.publish(ctx, runevent.KindDaemonReviewStatus,
			fmt.Sprintf("Review Source status: %s", status.State),
			map[string]any{"state": status.State, "detail": status.Detail},
		); err != nil {
			return settledWaitResult{}, err
		}
		if status.State == StatusSettled {
			return settledWaitResult{status: status, statusChecks: statusChecks}, nil
		}
		if clock.Now().Sub(startedAt) >= req.ReviewTimeout {
			return settledWaitResult{
				status:       Status{State: store.StateTimedOut, Detail: status.Detail},
				statusChecks: statusChecks,
			}, nil
		}
		if err := sleeper.Sleep(ctx, req.PollInterval); err != nil {
			return settledWaitResult{}, err
		}
	}
}

type confirmResult struct {
	ready    bool
	missing  bool
	timedOut bool
}

func confirmMergeReady(ctx context.Context, req Request, source CheckSource, headSHA string, clock Clock, sleeper Sleeper, publisher watchEventPublisher) (confirmResult, error) {
	if !req.UntilClean || source == nil {
		return confirmResult{ready: true}, nil
	}
	startedAt := clock.Now()
	for {
		if err := ctx.Err(); err != nil {
			return confirmResult{}, err
		}
		state, err := source.Check(ctx, headSHA)
		if err == nil {
			if err := publisher.publish(ctx, runevent.KindDaemonReviewStatus,
				fmt.Sprintf("Review Source check: %s", state),
				map[string]any{"state": state, "head_sha": headSHA},
			); err != nil {
				return confirmResult{}, err
			}
			switch state {
			case CheckSuccess:
				return confirmResult{ready: true}, nil
			case CheckMissing:
				return confirmResult{ready: true, missing: true}, nil
			case CheckFailure:
				return confirmResult{}, nil
			case CheckPending:
			default:
				return confirmResult{}, fmt.Errorf("unknown Review Source check state %q", state)
			}
		} else if err := publisher.publish(ctx, runevent.KindDaemonReviewStatus,
			"Review Source check poll failed; retrying.",
			map[string]any{"head_sha": headSHA, "error": err.Error()},
		); err != nil {
			return confirmResult{}, err
		}
		if clock.Now().Sub(startedAt) >= req.ReviewTimeout {
			return confirmResult{timedOut: true}, nil
		}
		if err := sleeper.Sleep(ctx, req.PollInterval); err != nil {
			return confirmResult{}, err
		}
	}
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
	if req.BudgetEnabled && req.MaxRunDuration <= 0 {
		return errors.New("watch max run duration must be greater than 0 when Run Budget is enabled")
	}
	if deps.StatusSource == nil {
		return errors.New("watch requires Review Source status boundary")
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

// watchEventPublisher appends watch-loop Run Events: status waits, quiet
// periods, and fetch results. Publication is part of the Run state
// contract, so a critical sink failure fails the Run.
type watchEventPublisher struct {
	sink  runevent.Sink
	runID string
	clock Clock
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
