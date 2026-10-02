---
task: task_01
spec: 0213-a-test-suite-that-does-not-flake
status: completed
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

## Result

Implemented this Task's test slice for Daemon Verification. Task status and
the declared Verification commands remain Daemon-owned.

- The linger subtest now uses deadline-bounded `testwait.Poll` to read the
  committed event through `RunEventsAfter` (via `countPublishableEvents`),
  then asserts an empty pending batch. Batch size 100 and linger 50 ms remain;
  no explicit flush, sleep, or independent 2-second deadline was added.
- Both budget tests now schedule independent Tasks at concurrency 2, using
  the existing fake Task Worktree manager for the parallel scheduler.
  `budgetTestClock` starts at `RunStartedAt`, and each maximum is one hour.
  The shared runner waits for the stalled Task's start before setting the
  clock to the real present minus that maximum and returning the first Task's
  settlement. The production watchdog still performs cancellation.
- Existing DeadlineExceeded, renewing-Task reason, completed-count, and
  on-disk completion assertions remain. Every other test function in both
  files is unchanged; no production file or deadline configuration changed.

Focused checks, using `GOCACHE=/private/tmp/roundfix-0213-task01-go-cache`:

- `rtk proxy go test -count=10 -cpu 1,4 -run '^(TestBatch.*|TestCommitJournalBatchClassifiesAmbiguousCommit)$' ./internal/store`
  — exit 0, 3.817 s; exercises every test in the store file, including the
  linger subtest, 10 times at each CPU setting.
- `rtk proxy go test -count=10 -cpu 1,4 -run '^(TestTaskBudget.*|TestTaskCycleReportsTheRenewedBudgetDeadline|TestTaskCycleReportsNoBudgetDeadlineWhenTheBudgetIsDisabled)$' ./internal/daemon`
  — exit 0, 27.692 s; exercises all nine tests in the budget file, 10 times
  at each CPU setting.
- Source inspection against `HEAD` confirmed every other test function is
  unchanged, the store file contains none of the removed unsafe wait patterns,
  and the budget file contains no numeric millisecond budget.
- `rtk proxy go test -race -count=1 -cpu 1,4 -run '^(TestBatchClosesOnCountLingerAndImmediate|TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement|TestTaskBudgetReasonNamesTheSettlementThatRenewedIt)$' ./internal/store ./internal/daemon`
  — exit 0; store 9.058 s, daemon 20.314 s, no race diagnostics.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.
- `rtk proxy make verify-incremental` — exit 0 with host process access;
  formatting, vet, repository tests, skill checks, and build passed. The
  worktree stayed unchanged throughout this successful run.

The initial checks could not access the shared Go build cache under the
sandbox; the task-scoped cache resolved that limitation. The first concurrent
fixture attempt ended before either runner started because it retained the
serial fixture's default Git Task Worktree manager. Using the existing
parallel fixture resolved that setup error. The earlier incremental check was
stopped because its test binary predated this correction.
The next incremental run's CLI and daemon tests reported PASS, but their
repository guards rejected the concurrent addition of this Result section.
The successful final run above corrected that check ordering by holding the
worktree unchanged and recording its outcome only after the process exited.

Acceptance evidence: the structural requirements for the linger and budget
criteria are present, and both files' edited and unedited tests have focused
repetition evidence. The declared 100-run linger and 50-run budget gates are
reserved for the Daemon; those exact acceptance counts are not claimed here.

## Carry-forward provenance

- Source Run: `run_20261002T064202Z_a22e1dd4644c4dc6`
- Source commit: `572600e5355e0fe97b9a744a2f960ba967cd3a60`
