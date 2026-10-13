---
task: task_01
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: completed
type: test
complexity: high
---

# Task 01: The detach and other CLI tests bound by real time only what they assert

## Overview

Two detach tests failed Spec 0259's QA gate because a 500 ms liveness budget
they do not assert expired while a loaded machine started their child
(`_prd.md` → the origin). This Task fixes them and the other `internal/cli`
tests of the inventory, rows 1 to 14 of `_techspec.md` → Current behavior, so
each fires the phase it asserts whatever the load. It is verifiable on its own:
the structural checks show no unasserted short budget remains in these tests,
and every changed test passes under the race detector with its assertions.

## Requirements

1. MUST answer Spec 0259's QA failure of 2026-10-12 (`qa-report-2026-10-12-03.md`,
   Run `run_20261013T014653Z_1fc0f49baff7f737`), in which
   `TestRunDetachedCommandTimesOutWaitingForRunCreation` and
   `TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline`
   failed with `Detached Run child produced no liveness signal within 500ms`.
2. MUST give every phase budget a detach test does not assert the value of
   `testwait.Bound(t)`, and give
   `TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline`
   the production budgets `detachHandshakeTimeout` and
   `detachRunCreationTimeout`, as `_techspec.md` → API Contract 1 states. The
   asserted 50 ms budgets of rows 2 and 3 MUST stay, and every assertion of
   rows 1 to 8 (exit status, stdout, stderr, each diagnostic substring) MUST
   stay byte-identical.
3. MUST change the `sleep-before-liveness` and `slow-run-creation` child modes
   to block through `waitForDetachTestSentinel` on a sentinel path the test
   never writes, and make the `run-created` mode wait for the duration in
   `ROUNDFIX_DETACH_TEST_RUN_CREATION_DELAY` (750 ms when unset) after its
   liveness byte, with `runDetachParentForTest` setting it to the liveness
   budget, as API Contract 2 states. A blocking child MUST still end when its
   owner is gone (ADR-0213).
4. MUST NOT change, in `internal/cli/detach_test.go`, the import block, the
   constant block other than adding the new variable name, the survivor struct,
   `TestDetachedChildIsTerminatedAtTeardown`, the owner-environment line of
   `runDetachParentForTest`, the survivor identity read after
   `runDetachedCommand`, the `ignore-sentinel` mode, the owner check inside
   `waitForDetachTestSentinel`, or `terminate`: Spec 0259 rewrites those lines
   (`_techspec.md` → 0259 overlap).
5. MUST fix rows 9 to 14 as API Contract 3 states: attach is cancelled once its
   output holds `Following live events.`; the "inside budget" Run Budget is one
   hour; the Force Stop tests wait for the owner and the monitor through
   `testwait`; the forge probe test keeps `context.DeadlineExceeded` and the
   250 ms budget with the child's 30 s sleep as its elapsed ceiling and cancels
   its second probe after the child's PID file appears; the dispose barriers
   are bounded by the test deadline the parent passes to its children. In
   `internal/cli/implement_test.go` it MUST NOT change the caller kill of
   `TestRunImplementDetachSurvivesCallerProcessGroupKill`, and in
   `internal/cli/orphan_unix_test.go` it MUST NOT change
   `TestCLIForceStopOwnerProcessHelper`.
6. MUST keep every changed test calling `t.Parallel()` first, add no Sequential
   Test, send no signal and write no executable (ADR-0259, ADR-0266, ADR-0125).
7. MUST NOT change any production file.

## Subtasks

- [ ] Give the unasserted detach budgets the test deadline or the production constants.
- [ ] Make the racing detach children block, and delay Run creation past the liveness deadline.
- [ ] Cancel attach on its follow-loop line and give the inside-budget Run an hour.
- [ ] Wait on the owner, the monitor, the probe child and the barrier through events and the test deadline.
- [ ] Run each changed test once under the race detector.

## Acceptance Criteria

- [ ] No detach test passes a budget it does not assert below the test deadline, except row 1's production budgets.
- [ ] Rows 1 to 8 keep every assertion and pass under `-race`.
- [ ] Rows 9 to 14 pass under `-race` with no 2 s, 3 s, 5 s or 10 s wait and no 10 s Run Budget left on work they do not assert.
- [ ] No production file and none of the 0259-owned regions changed.

## Context

- instruction: `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `internal/cli/detach.go`
- instruction: `internal/cli/attach.go`
- instruction: `internal/testwait/testwait.go`
- interface: `internal/cli/detach_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/orphan_unix_test.go`
- interface: `internal/cli/readiness_forge_test.go`
- interface: `internal/cli/review_dispose_lock_test.go`

## Verification

- `grep -qF 'ROUNDFIX_DETACH_TEST_RUN_CREATION_DELAY' internal/cli/detach_test.go || { printf 'the run-created child has no delay variable\n' >&2; exit 1; }; grep -qF 'detachHandshakeTimeout, detachRunCreationTimeout' internal/cli/detach_test.go || { printf 'row 1 does not use the production budgets\n' >&2; exit 1; }; ! grep -qF '500*time.Millisecond' internal/cli/detach_test.go || { printf 'detach_test.go still injects a 500 ms budget\n' >&2; exit 1; }; ! grep -qF ', time.Second, time.Second)' internal/cli/detach_test.go || { printf 'detach_test.go still injects 1 s budgets\n' >&2; exit 1; }; ! grep -Eq '(liveness|runCreation): +time[.]Second,' internal/cli/detach_test.go || { printf 'the review refusal test still injects 1 s budgets\n' >&2; exit 1; }; out="$(go test -count=1 -race -v -run '^(TestRunDetachedCommand.*|TestRunDetachedReviewRefusalLeavesArtifactDirectoryAbsent)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunDetachedCommandTimesOutWaitingForLiveness TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline TestRunDetachedCommandTimesOutWaitingForRunCreation TestRunDetachedCommandReportsInvalidLivenessHandshake TestRunDetachedCommandReportsMalformedRunCreationHandshake TestRunDetachedCommandReportsChildExitBeforeHandshakeWithOutput TestRunDetachedReviewRefusalLeavesArtifactDirectoryAbsent TestRunDetachedCommandReportsSilentChildExitBeforeHandshake; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the delay variable is absent and the 500 ms and 1 s budgets are present, so the command fails.
- `grep -qF 'Following live events.' internal/cli/implement_test.go || { printf 'attach is not cancelled on its follow-loop line\n' >&2; exit 1; }; ! grep -qF 'max_run_duration: 10s' internal/cli/implement_test.go || { printf 'the inside-budget Run still has a 10 s budget\n' >&2; exit 1; }; test "$(grep -cF 'time.After(2 * time.Second)' internal/cli/orphan_unix_test.go)" = 1 || { printf 'a Force Stop test still waits 2 s for the owner\n' >&2; exit 1; }; ! grep -qF 'context.WithTimeout(t.Context(), 5*time.Second)' internal/cli/orphan_unix_test.go || { printf 'the Force Stop monitor still has a 5 s budget\n' >&2; exit 1; }; ! grep -qF 'time.Since(started) > 3*time.Second' internal/cli/readiness_forge_test.go || { printf 'the forge probe test still has a 3 s ceiling\n' >&2; exit 1; }; ! grep -qF 'time.Now().Add(10 * time.Second)' internal/cli/review_dispose_lock_test.go || { printf 'the dispose barrier still has a 10 s deadline\n' >&2; exit 1; }; out="$(go test -count=1 -race -v -run '^(TestRunImplementDetachSurvivesCallerProcessGroupKill|TestImplementRunBudgetBoundsSetupAndIntegration|TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion|TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner|TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken|TestConcurrentDisposeOfOneFindingAppendsOnce|TestConcurrentDisposeOfTwoFindingsBothSucceed)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunImplementDetachSurvivesCallerProcessGroupKill TestImplementRunBudgetBoundsSetupAndIntegration TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken TestConcurrentDisposeOfOneFindingAppendsOnce TestConcurrentDisposeOfTwoFindingsBothSucceed; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the follow-loop line is absent from the test and the short waits are present, so the command fails.

## References

- `_prd.md` → Goals; Core Feature 1; Core Feature 2; Success Metric 1; Non-Goals / Out of Scope
- `_techspec.md` → Current behavior; 0259 overlap; API Contract 1; API Contract 2; API Contract 3; Invariants 2-4; Build Order 1
- ADR-0271; ADR-0213; ADR-0259; ADR-0266; ADR-0125


## Result

Implemented the task_01 test slice for Daemon Verification. The pre-existing
`status: in_progress` remains unchanged; Subtasks, Acceptance Criteria and
Verification remain authored as received.

### Implementation

- Detach phases that are not asserted use `testwait.Bound(t)` through
  `detachTestBound`, defined in `implement_test.go` to preserve the protected
  detach import block. The monitor-report test uses the production handshake
  and Run-creation constants. Both asserted 50 ms limits remain.
- The two timeout child modes wait on an unwritten sentinel through the existing
  owner-watching wait. Run creation reads
  `ROUNDFIX_DETACH_TEST_RUN_CREATION_DELAY`, defaults to 750 ms, and receives the
  liveness budget from the parent, after writing its liveness byte.
- Attach cancels synchronously when its output contains `Following live events.`;
  the inside-budget Run uses one hour. Force Stop observes the terminal Run and
  owner exit through `testwait`, and joins its cancelled monitor at teardown.
- The forge probe keeps its asserted 250 ms deadline and
  `context.DeadlineExceeded`, with the child's 30 s sleep as its elapsed ceiling.
  Its second probe waits for a nonempty PID file, cancels, asserts
  `context.Canceled` and the child's death, and joins the probe goroutine.
- Dispose parents wait through `testwait.Poll`; children receive the parent's
  test deadline minus the `testwait` margin as an RFC3339Nano environment value.
  All fifteen selected top-level tests call `t.Parallel()` first.

### Focused checks

The following focused implementation check ran once, independently of the
Daemon's authored Verification:

```sh
rtk proxy go test -race -count=1 -parallel 16 -v -run '^(TestRunDetachedCommandTimesOutWaitingForLiveness|TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline|TestRunDetachedCommandTimesOutWaitingForRunCreation|TestRunDetachedCommandReportsInvalidLivenessHandshake|TestRunDetachedCommandReportsMalformedRunCreationHandshake|TestRunDetachedCommandReportsChildExitBeforeHandshakeWithOutput|TestRunDetachedReviewRefusalLeavesArtifactDirectoryAbsent|TestRunDetachedCommandReportsSilentChildExitBeforeHandshake|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestImplementRunBudgetBoundsSetupAndIntegration|TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion|TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner|TestForgeProbesAreBoundedAndNeverPromptOrPrintAToken|TestConcurrentDisposeOfOneFindingAppendsOnce|TestConcurrentDisposeOfTwoFindingsBothSucceed)$' ./internal/cli
```

Exit 0; all fifteen selected tests and their subtests reported PASS, with no
race diagnostic or Helper Residue. Package duration: 12.217 s. The monitor-report
test took 10.08 s; the inside-budget subtest took 0.40 s.

A read-only `rtk proxy python3 -` audit compared the working tree to `HEAD` and
exited 0. It found the pre-change 500 ms/1 s detach budgets, finite 2 s child
sleeps and short CLI waits in `HEAD`, then checked the resulting contracts and
protected regions byte-for-byte.

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Unasserted detach budgets use the test deadline, except production budgets in row 1 | Audit confirmed the deadline helper delegates to `testwait.Bound(t)`, production constants in row 1, removal of the old injected budgets and retention of both 50 ms limits. |
| Rows 1–8 preserve every assertion and pass under the race detector | Audit compared all eight test bodies against `HEAD`, permitting only the budget substitutions; every other byte matched. All eight selected tests reported PASS in the focused race check. |
| Rows 9–14 pass without the unasserted short waits or 10 s Run Budget | All seven selected tests reported PASS. Audit confirmed event-driven attach/probe cancellation, deadline-bound Force Stop and barriers, and the one-hour Run Budget. The existing 2 s owner-helper cleanup remains outside the changed waits. |
| No production file or protected 0259 region changes | Audit confirmed exactly five test files plus this Task changed; the detach imports, existing constants, survivor struct, teardown test, owner environment line, survivor identity read, ignore-sentinel mode, sentinel owner check and terminate method matched `HEAD`. The caller kill and Force Stop helper also matched. No signal, executable or Sequential Test was added. |

The sandbox could not initialize because its writable worktree root contains a
symlink. Approved escalated execution was used for the scoped reads, edits and
focused check. No declared Verification, repository gate, commit, push or pull
request was run. No follow-up implementation was added.

## Carry-forward provenance

- Source Run: `run_20261013T045342Z_a78df1f6c39a67e2`
- Source commit: `40c3edddce2ab676d9a3cbf7f4f1432b82e8b74b`
