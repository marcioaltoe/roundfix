package watch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"roundfix/internal/store"
)

const (
	StatusPending   = "pending"
	StatusReviewing = "reviewing"
	StatusSettled   = "settled"
)

type Request struct {
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
}

type Result struct {
	Outcome             string
	Rounds              int
	Remaining           int
	ManualReviewCommand string
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
	Clock        Clock
	Sleeper      Sleeper
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

	startedAt := clock.Now()
	for round := 1; round <= req.MaxRounds; round++ {
		if budgetExceeded(req, startedAt, clock.Now()) {
			return Result{Outcome: store.StateBudgetExceeded, Rounds: round - 1}, nil
		}
		status, err := waitForSettled(ctx, req, deps.StatusSource, clock, sleeper)
		if err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		if status.State != StatusSettled {
			return Result{
				Outcome:             store.StateTimedOut,
				Rounds:              round - 1,
				ManualReviewCommand: "@coderabbitai review",
			}, nil
		}
		if err := sleeper.Sleep(ctx, req.QuietPeriod); err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		if budgetExceeded(req, startedAt, clock.Now()) {
			return Result{Outcome: store.StateBudgetExceeded, Rounds: round - 1}, nil
		}

		fetched, err := deps.Fetcher.Fetch(ctx, round)
		if err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round - 1}, err
		}
		if fetched.Issues == 0 {
			return Result{Outcome: store.StateClean, Rounds: round}, nil
		}

		resolved, err := deps.Resolver.Resolve(ctx)
		if err != nil {
			return Result{Outcome: store.StateFailed, Rounds: round}, err
		}
		if resolved.Remaining == 0 {
			return Result{Outcome: store.StateClean, Rounds: round}, nil
		}
		if !resolved.Progress {
			return Result{Outcome: store.StateFailed, Rounds: round, Remaining: resolved.Remaining}, fmt.Errorf("watch made no progress with %d Unresolved Review Issue(s) remaining", resolved.Remaining)
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

func waitForSettled(ctx context.Context, req Request, source StatusSource, clock Clock, sleeper Sleeper) (Status, error) {
	startedAt := clock.Now()
	for {
		if err := ctx.Err(); err != nil {
			return Status{}, err
		}
		status, err := source.Status(ctx, StatusRequest{
			PRNumber: req.PRNumber,
			HeadSHA:  req.HeadSHA,
		})
		if err != nil {
			return Status{}, err
		}
		if status.State == StatusSettled {
			return status, nil
		}
		if clock.Now().Sub(startedAt) >= req.ReviewTimeout {
			return Status{State: store.StateTimedOut, Detail: status.Detail}, nil
		}
		if err := sleeper.Sleep(ctx, req.PollInterval); err != nil {
			return Status{}, err
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
