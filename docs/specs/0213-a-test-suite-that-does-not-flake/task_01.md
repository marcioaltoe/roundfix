---
task: task_01
spec: 0213-a-test-suite-that-does-not-flake
status: pending
type: test
complexity: low
---

# Task 01: The linger and Run Budget tests wait on events, not elapsed time

## Overview

Two `internal/store` and `internal/daemon` tests read a proxy that elapsed time
can move. The linger subtest waits for an empty pending batch, and the linger
flush empties it before its commit lands. The two short Run Budget tests give
the first Task a 200 ms or 250 ms real allowance, which can expire before it
settles. This Task makes each test wait on the event its assertion reads. It
changes only these two test files.

## Requirements

1. MUST make the `linger closes a quiet batch` subtest of
   `TestBatchClosesOnCountLingerAndImmediate` wait, through `testwait.Poll`
   bounded by the test deadline, until `RunEventsAfter` returns the one
   published event. It MUST then assert that `pendingCount()` is zero. It MUST
   remove the 2-second `time.Now` deadline and the `time.Sleep` loop, and keep
   batch size 100 and linger 50 ms, so only the linger can commit the event.
2. MUST make `TestTaskBudgetReasonNamesTheSettlementThatRenewedIt` and
   `TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement` follow
   the TechSpec section "The Run Budget tests":
   - independent `task_01` and `task_02` with `Concurrency: 2`;
   - `engine.deps.Now` set to a `budgetTestClock` that starts at
     `RunStartedAt`, and a one-hour `MaxRunDuration`;
   - a `task_01` that settles only after `task_02` has started stalling, with
     the clock set to the real present minus the maximum just before it
     settles.
3. MUST keep each test's existing assertions: `context.DeadlineExceeded`, the
   reason suffix ` since Task task_01 settled.`, one completed Task, and
   `task_01` completed on disk.
4. MUST NOT configure a Run Budget maximum shorter than one hour anywhere in
   `internal/daemon/task_budget_renewal_test.go`, and MUST leave the other
   tests in that file unchanged.
5. MUST NOT change any production file, add a retry, or widen a deadline.

## Subtasks

- [ ] Make the linger subtest wait for the committed event.
- [ ] Move both short-budget tests onto the fake clock with the ordered renewal.
- [ ] Repeat both files' tests under `-count` and `-cpu 1,4`.

## Acceptance Criteria

- [ ] The linger subtest holds no `time.Sleep`, no 2-second deadline and no
      `pendingCount() == 0` wait, and passes 100 repeated runs at `-cpu 1,4`.
- [ ] Neither Run Budget test configures a maximum below one hour, and both pass
      50 repeated runs at `-cpu 1,4`. Real elapsed time can no longer decide
      their outcome.
- [ ] Every other test in both files passes unedited.

## Context

- interface: `internal/store/journal_batch_test.go`
- interface: `internal/daemon/task_budget_renewal_test.go`
- instruction: `internal/store/journal_writer.go`
- instruction: `internal/daemon/task_engine.go`
- instruction: `internal/testwait/testwait.go`

## Verification

- `! grep -nE 'time\.Sleep|2 \* time\.Second|pendingCount\(\) == 0' internal/store/journal_batch_test.go || exit 1; out="$(go test -count=100 -cpu 1,4 -v -run '^TestBatchClosesOnCountLingerAndImmediate$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestBatchClosesOnCountLingerAndImmediate/linger_closes_a_quiet_batch' || { printf '%s\n' "$out"; exit 1; }` — expected: exit 0 with no failed iteration; before this Task the linger loop is present, so the command fails.
- `! grep -nE '[0-9]+ \* time\.Millisecond' internal/daemon/task_budget_renewal_test.go || exit 1; out="$(go test -count=50 -cpu 1,4 -v -run '^(TestTaskBudgetReasonNamesTheSettlementThatRenewedIt|TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTaskBudgetReasonNamesTheSettlementThatRenewedIt TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0 with no failed iteration; before this Task the 200 ms and 250 ms budgets are present, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Core Features 1-2; Success Metrics 1-2
- [_techspec.md](_techspec.md) — The linger wait; The Run Budget tests; Testing Approach 1; Build Order 1
- [references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md](references/2026-09-30-time-bound-tests-and-a-leaked-detached-child.md)
- ADR-0098; ADR-0158; ADR-0164
