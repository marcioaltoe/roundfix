---
task: task_01
spec: 0171-a-deterministic-suite-under-load
status: completed
type: backend
complexity: medium
---

# Task 01: Deadline-bound waits that watch the work

## Overview

The Implement, Daemon and worktree tests each carry their own wait helpers, and
most of them wait a fixed wall-clock budget on one channel. A loaded machine
exceeds the budget, and a command that already ended is reported as slow work
long after its real error was written. This Task adds the one helper package the
other Tasks move onto.

## Requirements

1. MUST create `internal/testwait/testwait.go`, an ordinary package with no
   production caller, exporting `Margin`, `Bound`, `Until` and `Poll` as the
   TechSpec section "The wait helper" describes.
2. MUST make `Bound(t)` return the time until `t.Deadline()` minus `Margin`, read
   through an interface assertion on `Deadline() (time.Time, bool)`, and return
   the ten-minute fallback only when the test has no deadline. It MUST never
   return a smaller fixed duration.
3. MUST make `Until` return the first value from `ready`; fail the test at once
   when `ended` delivers first, printing that value; and fail at `Bound(t)` with
   `what` and every goroutine's stack. When `ready` and `ended` are both ready,
   the `ready` value MUST win. A nil `ended` watches nothing.
4. MUST make `Poll` re-evaluate its condition at a short interval until it
   holds; fail at once when `ended` delivers first, printing that value, after
   evaluating the condition one last time so work that ends right after its
   last effect does not fail the wait; and at `Bound(t)` fail with `what` and
   the condition's last observation.
5. MUST document that a failing wait calls `t.Fatalf`, so it runs on the test
   goroutine, and that it cancels nothing: the watched work is cancelled by the
   test's own context and cleanups.
6. MUST spawn no process from the package or its tests, so the suite-guard
   installation contract of ADR-0126 is not engaged.
7. MUST prove each behavior with its own named test in
   `internal/testwait/testwait_test.go`, using a `testing.TB` wrapper that
   overrides `Deadline` and records `Fatalf` instead of failing the real test.
   Each negative case is a separate test, and each test asserts the value it
   pins — the returned duration, the returned value, the recorded message and
   the elapsed time — not only that a call returned.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With a deadline ten minutes away, `Bound` returns ten minutes minus `Margin`
      within one second; with no deadline it returns ten minutes.
- [ ] `Until` returns a ready value, prefers it over a simultaneous end, fails
      within one second naming the ended work's value, and with a deadline
      `Margin` plus 300 ms away fails no earlier than 300 ms, naming `what`.
- [ ] `Poll` returns when its condition holds, fails within one second when the
      work ends, and fails at the deadline with the last observation.

## Context

- creates: `internal/testwait/testwait.go`
- creates: `internal/testwait/testwait_test.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/testfixture/fixture.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestBoundFollowsTheTestDeadline|TestBoundFallsBackWithoutATestDeadline|TestUntilReturnsTheReadyValue|TestUntilPrefersTheReadyValueOverAnEndedWork|TestUntilFailsAtOnceWhenTheWorkEnds|TestUntilFailsAtTheDeadlineAndNotBefore|TestPollReturnsWhenTheConditionHolds|TestPollFailsAtOnceWhenTheWorkEnds|TestPollFailsAtTheDeadlineWithTheLastObservation)$" ./internal/testwait 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestBoundFollowsTheTestDeadline TestBoundFallsBackWithoutATestDeadline TestUntilReturnsTheReadyValue TestUntilPrefersTheReadyValueOverAnEndedWork TestUntilFailsAtOnceWhenTheWorkEnds TestUntilFailsAtTheDeadlineAndNotBefore TestPollReturnsWhenTheConditionHolds TestPollFailsAtOnceWhenTheWorkEnds TestPollFailsAtTheDeadlineWithTheLastObservation; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the package `internal/testwait` does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Core Feature 1; Success Metric 3
- [_techspec.md](_techspec.md) — The wait helper; API Contracts 1-2

## Result

Implemented `internal/testwait` as a test-only package. `Bound` reads an
optional `Deadline() (time.Time, bool)` capability and otherwise uses only the
ten-minute fallback. `Until` gives a ready value priority over ended work and
reports ended work or a deadline stack dump through `t.Fatalf`. `Poll` checks
its condition immediately and at a 10 ms interval, performs a final check when
watched work ends, and reports the last observation at the deadline. Package
documentation records the test-goroutine, caller-owned cancellation and
cleanup contract. Neither the package nor its tests imports or invokes a
process API, so ADR-0126's suite-guard installation contract is not engaged.

Focused evidence by acceptance criterion:

- Deadline bound: `TestBoundFollowsTheTestDeadline` asserts the returned bound
  is ten minutes minus `Margin` within one second;
  `TestBoundFallsBackWithoutATestDeadline` asserts the exact ten-minute
  fallback. Both passed in `rtk env
  GOCACHE=/private/tmp/roundfix-task-01-gocache go test -count=1 -v -run
  '^TestBound' ./internal/testwait` (exit 0).
- Channel wait: `TestUntilReturnsTheReadyValue` pins the first queued value;
  `TestUntilPrefersTheReadyValueOverAnEndedWork` pins ready-over-ended
  priority; `TestUntilFailsAtOnceWhenTheWorkEnds` checks the recorded ended
  value and an elapsed time below one second; and
  `TestUntilFailsAtTheDeadlineAndNotBefore` checks at least 300 ms elapsed,
  the named wait and another goroutine's stack. All passed in `rtk env
  GOCACHE=/private/tmp/roundfix-task-01-gocache go test -count=1 -v -run
  '^TestUntil' ./internal/testwait` (exit 0; the deadline case took 0.30 s).
- Polling wait: `TestPollReturnsWhenTheConditionHolds` pins repeated condition
  evaluation; `TestPollFailsAtOnceWhenTheWorkEnds` pins the final false check,
  ended value and elapsed time below one second;
  `TestPollReturnsWhenEndedWorkPublishedItsLastEffect` pins the final true
  check; and `TestPollFailsAtTheDeadlineWithTheLastObservation` pins at least
  300 ms elapsed, the named wait and the final observation. All passed in `rtk
  env GOCACHE=/private/tmp/roundfix-task-01-gocache go test -count=1 -v -run
  '^TestPoll' ./internal/testwait` (exit 0; the deadline case took 0.30 s).

Additional focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-task-01-gocache go test
  ./internal/testwait` — exit 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task-01-gocache go vet
  ./internal/testwait` — exit 0.
- `rtk env GOCACHE=/private/tmp/roundfix-task-01-gocache go test -race
  ./internal/testwait` — exit 0.

The Daemon-owned command in `## Verification` was not run.
