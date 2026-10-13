---
task: task_03
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: pending
type: test
complexity: medium
---

# Task 03: The daemon, worktree, store and spec tests stop racing their work

## Overview

Five tests outside `internal/cli` and `internal/agent` bound Git, a shell or an
SQLite operation by a real budget they do not assert, or race an asserted
budget against setup work that can finish late (`_techspec.md` → Current
behavior, rows 22 to 26). The Run Budget test cannot reach the watchdog's timer
today, so the Daemon engine gains a timer dependency that defaults to the real
timer. It is verifiable on its own: the default is proved, the short budgets
are gone, and each changed test passes under the race detector.

## Requirements

1. MUST answer Spec 0259's QA failures of 2026-10-12 by fixing rows 22 to 26 of
   `_techspec.md` → Current behavior, the members of the class those failures
   exposed outside `internal/cli` and `internal/agent`.
2. MUST add `BudgetTimer` to the Daemon engine's `Deps` and make the Run Budget
   watchdog start its timer through it when it is set and through
   `time.NewTimer` otherwise, as `_techspec.md` → Interfaces and API Contract 5
   state. The duration passed MUST stay `time.Until(deadline)`; no production
   value changes.
3. MUST add `TestTaskBudgetWatchdogUsesTheRealTimerByDefault`, which runs a
   Task Cycle without `BudgetTimer` and with a deadline already past when
   task_01 starts, and proves the watchdog still cancels the running Task with
   the Run Budget deadline.
4. MUST make `TestTaskCycleEndsRunAtBudgetDeadline` expire the budget only after
   `budgetDeadlineRunner` reports that task_01 started, through a fake
   `BudgetTimer`, and keep its four assertions (the Run Budget error, the
   unprovable termination in its second case, only task_01 started with one
   cancellation, task_02 pending). It MUST NOT change
   `TestBudgetExceededRunRecordsItsReason`, which Spec 0259 rewrites, or the
   lines between the two tests.
5. MUST make `TestSourceRevisionRecoveryCancelsTheFetchProcess` cancel the fetch
   after the fetch fixture reports its start, asserting the cancellation it
   caused, the marker fields and the dead fetch process; make
   `TestRunBootstrapReturnsBootstrapErrorOnTimeout` accept a timeout that ends
   before the shell writes its marker while keeping every `BootstrapError`
   assertion; and bound the retention scan of
   `TestRetentionScanOutsideWriteTransaction` and the prune of
   `TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents` by the test's
   context, as API Contract 6 states. Each changed assertion is a declared
   change: the source revision test asserts the cancellation it caused instead
   of a 2 s deadline.
6. MUST keep every changed test parallel, add no Sequential Test, send no
   signal and write no executable (ADR-0259, ADR-0266, ADR-0125).

## Subtasks

- [ ] Add the engine's budget timer dependency and prove its default.
- [ ] Expire the Run Budget test's budget after task_01 starts.
- [ ] Cancel the fetch after it starts and accept an early bootstrap timeout.
- [ ] Bound the two store operations by the test's context.
- [ ] Run each changed test once under the race detector.

## Acceptance Criteria

- [ ] An engine without `BudgetTimer` still ends a Run at a real deadline, proved by a test.
- [ ] The Run Budget deadline test no longer depends on how long Git setup takes.
- [ ] No row 22 to 26 test bounds work it does not assert by 2 s, 1 s, 500 ms or 250 ms.
- [ ] Each changed test passes under `-race`.

## Context

- instruction: `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- instruction: `internal/daemon/task_budget_renewal_test.go`
- instruction: `internal/testwait/testwait.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/worktree/worktree_test.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/store/prune_reclaimable_test.go`
- interface: `internal/spec/source_revision_test.go`

## Verification

- `grep -q 'BudgetTimer' internal/daemon/engine.go || { printf 'the engine has no budget timer dependency\n' >&2; exit 1; }; ! grep -qF 'plan.MaxRunDuration = 3 * time.Second' internal/daemon/task_engine_test.go || { printf 'the Run Budget deadline test still races a 3 s budget\n' >&2; exit 1; }; out="$(go test -count=1 -race -v -run '^(TestTaskBudgetWatchdogUsesTheRealTimerByDefault|TestTaskCycleEndsRunAtBudgetDeadline)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTaskBudgetWatchdogUsesTheRealTimerByDefault TestTaskCycleEndsRunAtBudgetDeadline; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go build -buildvcs=false ./internal/daemon || exit 1` — expected: exit 0. Before this Task the dependency and the default test are absent and the 3 s budget is present, so the command fails.
- `! grep -qF 'context.WithTimeout(t.Context(), 2*time.Second)' internal/spec/source_revision_test.go || { printf 'the fetch cancellation test still has a 2 s deadline\n' >&2; exit 1; }; ! grep -qF 'context.WithTimeout(ctx, 500*time.Millisecond)' internal/store/journal_test.go || { printf 'the retention scan test still has a 500 ms budget\n' >&2; exit 1; }; ! grep -qF 'context.WithTimeout(ctx, 250*time.Millisecond)' internal/store/prune_reclaimable_test.go || { printf 'the prune test still has a 250 ms budget\n' >&2; exit 1; }; out="$(go test -count=1 -race -v -run '^TestSourceRevisionRecoveryCancelsTheFetchProcess$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestSourceRevisionRecoveryCancelsTheFetchProcess ' || { printf 'the fetch cancellation test did not pass\n' >&2; exit 1; }; out="$(go test -count=1 -race -v -run '^(TestRetentionScanOutsideWriteTransaction|TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRetentionScanOutsideWriteTransaction TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -race -v -run '^TestRunBootstrapReturnsBootstrapErrorOnTimeout$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestRunBootstrapReturnsBootstrapErrorOnTimeout ' || { printf 'the bootstrap timeout test did not pass\n' >&2; exit 1; }` — expected: exit 0. Before this Task the 2 s, 500 ms and 250 ms budgets are present, so the command fails.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 1; Success Metric 4
- `_techspec.md` → Current behavior; 0259 overlap; Interfaces; API Contract 5; API Contract 6; Invariant 1; Testing Approach; Build Order 3
- ADR-0271; ADR-0259; ADR-0266; ADR-0125
