---
task: task_07
spec: 0259-proven-signals-and-a-shared-build-cache
status: completed
type: test
complexity: medium
---

# Task 07: The Run Budget test no longer races the machine's load

## Overview

Corrective Task 2 of 2. `TestBudgetExceededRunRecordsItsReason` in
`internal/daemon/task_engine_test.go` gives the whole Task cycle a real 500 ms
Run Budget measured from `RunStartedAt`, and requires task_01 to finish its Git
work and commit inside it. Under load the Git work alone takes longer, and the
test fails with
`completed before budget = 0 with 0 commits, want 1 each`. It failed that way
twice on 2026-10-12 in `make verify-changed`: in the task_06 attempt 1 of Run
`run_20261013T004148Z_5319d413bd50ad25`, and in the starting verification of
Run `run_20261013T013741Z_3a6796ce5e63cbc3` (load average 16–18). It passes 20
times out of 20 when run alone. Requirement 6 of task_05, the acceptance of F6,
needs ten consecutive full suites that exit 0, so a test that fails with load
stops that measurement for a reason this Spec does not measure.

## Requirements

1. MUST make `TestBudgetExceededRunRecordsItsReason` independent of how long
   task_01's Git work takes: the Run Budget may expire only after task_01 has
   completed and been committed. Use the engine's existing seams
   (`engine.deps.Now`, the budget renewal after a settled Task, the fake
   runner's per-Task blocking) — for example, a test clock that stays at
   `RunStartedAt` until task_01 settles — and MUST keep every assertion of the
   test: the `context.DeadlineExceeded` error, `store.StateBudgetExceeded`,
   the reason naming the configured maximum and the elapsed time within
   `publicOutcomeReasonMaxRunes`, one completed Task with one commit, task_01
   `completed` and task_02 `in_progress` on disk.
2. MUST replace the test's "Sequential" comment with `t.Parallel()` as its
   first statement, and MUST lower the `internal/daemon` entry of
   `internal/testfixture/parallel_tests_test.go` by one if that record counts
   sequential tests, so the record stays exact.
3. MUST NOT change production code in `internal/daemon` unless a seam is
   missing; any production change MUST keep the Run Budget behavior and every
   other test unchanged.
4. MUST record in the Result one sabotage (for example, letting the test clock
   advance before task_01 settles) and the failure it produced.
5. MUST NOT change any file outside its Context.

## Subtasks

- [ ] Hold the test clock until task_01 settles, then let the budget expire.
- [ ] Make the test parallel and keep the parallel record exact.
- [ ] Prove the test still fails when the budget expires early.

## Acceptance Criteria

- [ ] The test calls `t.Parallel()` first and decides nothing on elapsed wall time before task_01 settles.
- [ ] It passes 50 consecutive times with `-race -parallel 16`.

## Context

- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/testfixture/parallel_tests_test.go`

## Verification

- `grep -A1 '^func TestBudgetExceededRunRecordsItsReason(' internal/daemon/task_engine_test.go | grep -q 't.Parallel()' || { printf 'the budget test does not call t.Parallel first\n' >&2; exit 1; }; out="$(go test -count=50 -race -parallel 16 -v -run '^TestBudgetExceededRunRecordsItsReason$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; test "$(printf '%s\n' "$out" | grep -c -- '--- PASS: TestBudgetExceededRunRecordsItsReason ')" -eq 50 || { printf 'fewer than 50 passes\n' >&2; exit 1; }; go test -count=1 -run '^TestEveryTestIn' ./internal/testfixture || exit 1`

## References

- Runs `run_20261013T004148Z_5319d413bd50ad25` and `run_20261013T013741Z_3a6796ce5e63cbc3`, `make verify-changed` logs
- `_techspec.md` → Testing Approach
- ADR-0259

## Result

The test now calls `t.Parallel()` first. Its injected clock stays at
`RunStartedAt` through task_01's Agent work, Verification, settlement, commit
and budget renewal. The logical start is the test harness deadline (or a day
in the future when the harness supplies no deadline), so the real watchdog
cannot expire the 500 ms budget during that work. The existing fake runner's
`afterTask` hook advances the clock to the budget deadline when task_02's
Agent returns; the engine's existing post-Agent budget checkpoint stops the
Run before task_02 can settle or commit.

Every original assertion remains: `context.DeadlineExceeded`,
`store.StateBudgetExceeded`, configured maximum and elapsed time in the
bounded reason, one completed Task and one commit, task_01 `completed`, and
task_02 `in_progress` on disk. An additional assertion requires
`BudgetRenewedBy` to name task_01. The `internal/daemon` Sequential Test
ceiling falls from 16 to 15. Production code is unchanged.

Focused evidence:

- Acceptance criterion 1: source review confirms first-statement
  `t.Parallel()` and clock advancement only after task_02 runs. A temporary
  Go overlay added `time.Sleep(2 * maximum)` to task_01's `afterTask` hook,
  before settlement and commit. Running
  `go test -overlay=/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task07-slow-n2c0q_vv/overlay.json -count=1 -race -parallel 16 -v -run '^TestBudgetExceededRunRecordsItsReason$' ./internal/daemon`
  exited 0: the test passed in 1.28 s despite the one-second delay exceeding
  the configured 500 ms maximum.
- Acceptance criterion 2, focused evidence only:
  `go test -count=10 -race -parallel 16 -shuffle=on -v -run '^(TestBudgetExceededRunRecordsItsReason|TestTaskCycleDisabledBudgetDerivesNoDeadline|TestTaskCycleFinishesBeforeBudgetDeadline|TestTaskCycleSettlesBudgetOutcomeWithoutError)$' ./internal/daemon`
  exited 0, with ten consecutive passes for each of the four tests and no
  race reports. Shuffle seed: `1791856278834806000`. The declared 50-pass
  Verification remains for the Daemon; this focused run does not claim it.
- Parallel record:
  `go test -count=1 -race -v -run '^TestEveryTestInAParallelTestPackageRunsInParallel$/repository$' ./internal/testfixture`
  exited 0 and reported `internal/daemon: 15 Sequential Tests, ceiling 15`.
- Starting reproduction: an overlay made the original wall-clock test start
  with its budget already expired. A single focused run exited 1 with
  `completed before budget = 0 with 0 commits, want 1 each`.
- Required sabotage: a separate overlay changed the new `afterTask` clock
  trigger from task_02 to task_01, advancing time before task_01's settlement
  and commit. Running
  `go test -overlay=/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task07-sabotage-i6anypko/overlay.json -count=1 -race -parallel 16 -v -run '^TestBudgetExceededRunRecordsItsReason$' ./internal/daemon`
  exited 1 with
  `task_engine_test.go:7778: completed before budget = 0 with 0 commits, want 1 each`.
  Both probe edits exist only in temporary overlays; the repository retains
  the task_02 trigger and has no added delay.
- `git diff --check` exited 0. Scope review found only the two declared test
  files and this Result changed by the Agent. The pre-existing
  `status: in_progress` edit is Daemon-owned and preserved.

Declared Task Verification and repository Verification were not run. Task
status, authored checkboxes, the Task Graph and other Task files are unchanged.
No commit, push or Pull Request was made. No follow-up work was found.

## Carry-forward provenance

- Source Run: `run_20261013T014653Z_1fc0f49baff7f737`
- Source commit: `d8e2b1b58b157c9522e304192b291d6f8b67f649`
