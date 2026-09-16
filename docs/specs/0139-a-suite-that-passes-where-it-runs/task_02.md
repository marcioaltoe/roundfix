---
task: task_02
spec: 0139-a-suite-that-passes-where-it-runs
status: pending
type: test
complexity: medium
---

# Task 02: Bind Task-cycle waits to the test deadline

## Overview

The Task-cycle test helpers wait a fixed two seconds for a scheduler start or a
TaskCycle result. Under load, the Daemon returns the correct result later than
that: fixture copies and git processes take real time. So the tests fail
although the engine is right.

This slice replaces the guessed bound with one derived from the running test's
own deadline. A correct result is no longer failed by load, and a genuine hang
still fails, with the test binary's timeout report.

## Requirements

1. MUST add one helper, `testWaitBound(t testing.TB) time.Duration`, which
   returns:
   - the time remaining before `t.Deadline()`, minus a small margin, when the
     test has a deadline;
   - a generous fixed fallback when it does not.
2. MUST replace every `time.After(2 * time.Second)` wait in the Daemon test file
   with a wait bounded by `testWaitBound`. That is every waiting helper, including
   `waitSchedulerStarts`, `waitIntegratedTask`, `waitTaskCycleResult`,
   `waitVerificationStart`, `waitPublishedStopEvent` and
   `waitPublishedVerificationPhase`.
3. MUST add `TestWaitBoundFollowsTheTestDeadline`, proving the bound tracks the
   deadline and that the no-deadline fallback applies.
4. MUST keep every assertion that nothing extra started as an immediate check,
   not a wait.
5. MUST make a wait that reaches its bound fail the test with a message naming
   what it waited for, followed by every goroutine's stack (for example through
   `runtime.Stack` with all goroutines), so a genuine hang stays diagnosable
   although the helper fails before the test binary's own timeout.
6. MUST keep every existing Task-cycle assertion passing, and MUST NOT raise a
   constant timeout, add a retry or skip, or remove `t.Parallel`.
7. MUST NOT change production code.

## Subtasks

- [ ] Add `testWaitBound` and its test.
- [ ] Replace the eight fixed two-second waits.
- [ ] Confirm the Task-cycle tests pass under concurrent load.

## Acceptance Criteria

- [ ] No `time.After(2 * time.Second)` wait remains in the Daemon test file;
      today there are eight.
- [ ] `TestWaitBoundFollowsTheTestDeadline` passes, and a wait that reaches its
      bound reports every goroutine's stack.
- [ ] `go test -count=3 -parallel 16 ./internal/daemon -run 'TestTaskCycle'`
      passes while `go test ./internal/baseline` runs concurrently.

## Context

- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/task_context_test.go`

## Verification

- `test -f internal/daemon/task_engine_test.go && ! grep -q 'time.After(2 \* time.Second)' internal/daemon/task_engine_test.go` — no fixed two-second wait remains; this fails today.
- `grep -q 'func TestWaitBoundFollowsTheTestDeadline' internal/daemon/task_engine_test.go && grep -q 'runtime.Stack(' internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon -run '^TestWaitBoundFollowsTheTestDeadline$'` — the bound follows the test deadline and an expired wait reports goroutine stacks; this fails today.
- `grep -q 'func testWaitBound(' internal/daemon/task_engine_test.go || exit 1; go test -count=1 ./internal/baseline >/dev/null 2>&1 & load=$!; go test -count=3 -parallel 16 ./internal/daemon -run 'TestTaskCycle'; status=$?; wait "$load"; exit "$status"` — the Task-cycle tests pass under concurrent load. The baseline run only generates load, so its exit status is deliberately not the gate.

## References

- `_prd.md` → Goals 1 and 3; User Story 3; Core Feature 2; Regression locks.
- `_techspec.md` → Implementation Design: Task-cycle waits; Testing Approach 2;
  Build Order 2.
