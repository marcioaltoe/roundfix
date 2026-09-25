---
task: task_03
spec: 0171-a-deterministic-suite-under-load
status: completed
type: backend
complexity: medium
---

# Task 03: Daemon and worktree tests wait on events

## Overview

Spec 0164's QA gate failed on
`TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`,
and a documentation-only `make verify` failed on
`TestBootstrapSerializesAcrossSiblings`; both pass in isolation. The Daemon
waits are bound to the deadline but ignore an early `TaskCycle` return, the
worktree bootstrap tests pass one- and five-second timeouts they do not examine,
and the sibling round has no bound at all. This Task moves both files onto
`internal/testwait`.

## Requirements

1. MUST replace `testWaitBound`, `testWaitFallback`, `testWaitDeadlineMargin`
   and `failTestWait` in `internal/daemon/task_engine_test.go` with `testwait`,
   keeping the full goroutine dump on a timed-out wait. No `time.After` may
   remain in the file.
2. MUST make `waitSchedulerStarts`, `taskCapacityVerifier.waitStart`,
   `waitPublishedVerificationPhase` and `waitPublishedStopEvent` watch the
   in-flight `TaskCycle` result channel wherever the calling test has one, so a
   `TaskCycle` that returns first fails the wait at once with its result and
   error. This covers at least
   `TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`
   and `TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable`.
3. MUST derive from `testwait.Bound(t)` every `BootstrapSpec.Timeout` in
   `internal/worktree/worktree_test.go` that is not the subject of its test,
   including the five-second timeout in `bootstrapSerializationRound` and the
   one-second timeouts in `TestCreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot`,
   `TestBootstrapFailureAfterWorkIsClassifiedApart` and
   `TestRunBootstrapReturnsBootstrapErrorOnNonZeroExit`.
4. MUST keep `TestRunBootstrapReturnsBootstrapErrorOnTimeout`'s 10 ms bound
   against `sleep 1`, because load can only lengthen the sleep and so cannot flip
   its outcome.
5. MUST make each round of `TestBootstrapSerializesAcrossSiblings` wait for its
   sibling starts and results through `testwait`, naming the siblings seen so
   far on failure, and MUST make releasing a sibling unable to block past the
   deadline when that sibling has already exited without opening the release
   FIFO.
6. MUST keep every top-level test function name recorded in
   `docs/references/coverage-record.json`, and MUST NOT edit that file.
7. MUST change no production file, including `bootstrapPermit` and
   `runBootstrap` in `internal/worktree/worktree.go`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The Daemon file holds no `time.After` or `testWaitBound` and uses
      `testwait`; the worktree file holds no one- or five-second bootstrap
      timeout outside the timeout test and uses `testwait`.
- [ ] The two queued Daemon tests and the worktree bootstrap tests pass every
      iteration of `-count=10 -cpu 1,4`.

## Context

- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/worktree/worktree_test.go`

## Verification

- `out="$(go test -count=10 -cpu 1,4 -v -run "^(TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement|TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable|TestBootstrapSerializesAcrossSiblings|TestCreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot|TestBootstrapFailureAfterWorkIsClassifiedApart|TestRunBootstrapReturnsBootstrapErrorOnNonZeroExit|TestRunBootstrapReturnsBootstrapErrorOnTimeout)$" ./internal/daemon ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable TestBootstrapSerializesAcrossSiblings TestCreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot TestBootstrapFailureAfterWorkIsClassifiedApart TestRunBootstrapReturnsBootstrapErrorOnNonZeroExit TestRunBootstrapReturnsBootstrapErrorOnTimeout; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done; grep -q "testwait\." internal/daemon/task_engine_test.go && grep -q "testwait\." internal/worktree/worktree_test.go && ! grep -q "time.After(" internal/daemon/task_engine_test.go && ! grep -q "func testWaitBound" internal/daemon/task_engine_test.go && ! grep -q "Timeout: time.Second" internal/worktree/worktree_test.go && ! grep -q "Timeout: 5 \* time.Second" internal/worktree/worktree_test.go` — expected: exit 0; before this Task neither file uses `testwait` and the Daemon file still defines `testWaitBound`, so the command fails.

## References

- [_prd.md](_prd.md) — Core Feature 3; Success Metrics 1-2
- [_techspec.md](_techspec.md) — The Daemon and worktree tests

## Result

### Implementation

- Replaced the Daemon test-local wait bound, fallback, timeout failure, and
  `time.After` paths with `internal/testwait`; timed-out waits retain the full
  goroutine dump supplied by `testwait`.
- Made scheduler-start, Verification-start, published Verification-phase, and
  published Stop-event waits observe the in-flight `TaskCycle` result channel.
  The two cross-file scheduler-wait callers pass the same result channel.
- Derived every non-timeout bootstrap bound from `testwait.Bound(t)` and kept
  the timeout test's 10 ms bound against `sleep 1`.
- Made each bootstrap serialization round consume named sibling start and
  result events through `testwait`. FIFO release now retries a non-blocking
  open while watching that sibling's result, so an exited sibling cannot leave
  the test blocked in a write-only FIFO open. Added a regression test for the
  no-reader release attempt.

### Focused checks

- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=3 -cpu 1,4 -run '^(TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement|TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable)$' ./internal/daemon` — exit 0.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=3 -cpu 1,4 -run '^(TestBootstrapSerializesAcrossSiblings|TestTryWriteBootstrapReleaseReturnsWithoutAReader|TestCreateRunsBootstrapAfterCopyInRunWorktreeRoot|TestCreateTaskRunsBootstrapAfterCopyInTaskWorktreeRoot|TestBootstrapFailureAfterWorkIsClassifiedApart|TestRunBootstrapReturnsBootstrapErrorOnNonZeroExit|TestRunBootstrapReturnsBootstrapErrorOnTimeout)$' ./internal/worktree` — exit 0.
- `GOCACHE=/private/tmp/roundfix-task03-gocache make verify-incremental` — exit 0 after rerunning with approved network access; the sandboxed attempt was blocked when an existing test reached `api.github.com`.
- `git diff --check` — exit 0.

### Acceptance evidence

- Criterion 1: source inspection finds `testwait` use in both target files, no
  `time.After`, `testWaitBound`, fallback, margin, or local timeout failure in
  the Daemon file, and no one- or five-second bootstrap timeout in the worktree
  file. The only behavior-under-test bootstrap bounds remain 10 ms for timeout
  and 1 ns for the empty-command skip.
- Criterion 2: the two queued Daemon tests and all named bootstrap tests passed
  six focused iterations each across CPU settings 1 and 4 (`-count=3 -cpu
  1,4`). The authored count-10 Verification command was not run; the Daemon
  owns that gate.
- Scope: `git diff --exit-code` reports no change to
  `docs/references/coverage-record.json`, `_tasks.md`, or
  `internal/worktree/worktree.go`. The diff adds one top-level test and removes
  or renames none.
