---
task: task_02
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: completed
type: test
complexity: high
---

# Task 02: The agent tests stop racing the real cancel grace and cleanup timeout

## Overview

The `internal/agent` tests re-execute the test binary as a fake ACPX, and two
real budgets bound those children: the 2 s stop grace the fake harness sets and
the production 5 s timeout around the disposable-session close
(`_techspec.md` → Current behavior, rows 15 to 21). Neither is what the tests
assert. This Task gives the runner its cleanup timeout as a dependency with the
production default, gives the tests the test deadline, and fixes the three
other agent rows. It is verifiable on its own: the cleanup sites read the
runner's timeout, the default is proved, and every changed test passes under
the race detector.

## Requirements

1. MUST answer Spec 0259's QA failures of 2026-10-12 by fixing rows 15 to 21 of
   `_techspec.md` → Current behavior, the agent members of the class those
   failures exposed.
2. MUST add the unexported `cleanupTimeout` field and the
   `disposableCleanupTimeout` method to `ACPXRunner` as `_techspec.md` →
   Interfaces and API Contract 4 state, and make the cleanup contexts in
   `fallback.go`, `selection_assignment.go` and `sealed.go` use the method. A
   runner without the field MUST keep the 5 s `acpxPreflightCleanupTimeout`;
   no production value changes.
3. MUST add `TestDisposableCleanupTimeoutDefaultsToTheProductionValue`, which
   proves that a runner whose `cleanupTimeout` is zero or negative returns
   `acpxPreflightCleanupTimeout`, and that a positive value is returned as set.
4. MUST set `cleanupTimeout` to `testwait.Bound(t)` in every test that builds an
   `ACPXRunner` whose cleanup spawns a child, and set the fake ACPX harness stop
   grace to `testwait.Bound(t)`; the fake cancellation clock still drives the
   grace timer the cancel tests assert, and their assertions stay.
5. MUST give `TestRealACPXCommandOverrideRoundTripCancelCrashResume`'s stop
   grace and its `ps` context the test deadline, start
   `TestSpawnedFixtureDeathFailsFast`'s measured span at the fixture's first
   output instead of before its startup while keeping its 9 s asserted bound,
   and make `TestConsoleDisplaySinkSerializesConcurrentPublish` wait for its
   first write through `testwait.Until`, as API Contract 6 states.
6. MUST keep every changed test parallel, add no Sequential Test, send no
   signal and write no executable (ADR-0259, ADR-0266, ADR-0125).

## Subtasks

- [ ] Add the runner's cleanup timeout dependency and use it at the three cleanup sites.
- [ ] Prove its default.
- [ ] Give the tests' runners and the fake harness the test deadline.
- [ ] Fix the integration, fixture-death and console-sink rows.
- [ ] Run each changed test once under the race detector.

## Acceptance Criteria

- [ ] A runner without `cleanupTimeout` still uses 5 s, proved by a test.
- [ ] No agent test bounds a fake ACPX child by a 2 s or 5 s real timeout it does not assert.
- [ ] The cancel, cleanup, sealed, fixture-death and console-sink tests pass under `-race`.

## Context

- instruction: `docs/adr/0271-a-test-bounds-by-real-time-only-what-it-asserts.md`
- instruction: `docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`
- instruction: `internal/agent/agent.go`
- instruction: `internal/testwait/testwait.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/fallback.go`
- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/sealed.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/sealed_test.go`
- interface: `internal/agent/node_options_test.go`
- interface: `internal/agent/selection_assignment_test.go`
- interface: `internal/agent/selection_catalogue_characterization_test.go`
- interface: `internal/agent/codex_spawn_test.go`
- interface: `internal/agent/acpx_integration_test.go`
- interface: `internal/agent/agent_test.go`

## Verification

- `grep -qF 'func (runner *ACPXRunner) disposableCleanupTimeout() time.Duration' internal/agent/acpx_runner.go || { printf 'the runner has no cleanup timeout method\n' >&2; exit 1; }; ! grep -q 'acpxPreflightCleanupTimeout' internal/agent/fallback.go internal/agent/selection_assignment.go internal/agent/sealed.go || { printf 'a cleanup site still reads the constant directly\n' >&2; exit 1; }; ! grep -qF 'harness.stopGrace = 2 * time.Second' internal/agent/acpx_runner_test.go || { printf 'the fake harness still has a 2 s stop grace\n' >&2; exit 1; }; ! grep -qF 'StopGrace: 2 * time.Second' internal/agent/acpx_integration_test.go || { printf 'the real ACPX test still has a 2 s stop grace\n' >&2; exit 1; }; ! grep -qF 'time.After(time.Second)' internal/agent/agent_test.go || { printf 'the console sink test still waits 1 s\n' >&2; exit 1; }; out="$(go test -count=1 -race -v -run '^(TestDisposableCleanupTimeoutDefaultsToTheProductionValue|TestACPXRunCancelsPromptCooperatively|TestACPXRunClosesSessionAfterCancelGracePeriod|TestACPXRunCancellationCommandFailuresWarnAndContinue|TestProveExactSelectionCancelCleanup|TestProveExactSelectionTimeoutCleanup|TestACPXProbeCancellationStillClosesDisposableSession|TestACPXProbeFallbackCancellationStillClosesCandidateSession|TestSealedACPXPromptCancellationCancelsAndClosesSession|TestSpawnedFixtureDeathFailsFast|TestConsoleDisplaySinkSerializesConcurrentPublish)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDisposableCleanupTimeoutDefaultsToTheProductionValue TestACPXRunCancelsPromptCooperatively TestACPXRunClosesSessionAfterCancelGracePeriod TestACPXRunCancellationCommandFailuresWarnAndContinue TestProveExactSelectionCancelCleanup TestProveExactSelectionTimeoutCleanup TestACPXProbeCancellationStillClosesDisposableSession TestACPXProbeFallbackCancellationStillClosesCandidateSession TestSealedACPXPromptCancellationCancelsAndClosesSession TestSpawnedFixtureDeathFailsFast TestConsoleDisplaySinkSerializesConcurrentPublish; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go build -buildvcs=false ./internal/agent || exit 1` — expected: exit 0. Before this Task the method and the default test are absent and the short budgets are present, so the command fails.

## References

- `_prd.md` → Goals; Core Feature 3; Success Metric 1; Success Metric 4
- `_techspec.md` → Current behavior; Interfaces; API Contract 4; API Contract 6; Invariant 1; Testing Approach; Build Order 2
- ADR-0271; ADR-0089; ADR-0259; ADR-0266; ADR-0125


## Result

Implemented the Task 02 slice for TechSpec inventory rows 15–21. The Daemon
retains Task status and declared Verification ownership.

### Implementation

- `ACPXRunner.cleanupTimeout` supplies the three disposable-session cleanup
  contexts through `disposableCleanupTimeout()`. Zero and negative values use
  the existing `acpxPreflightCleanupTimeout`; the production constant remains
  5 s. A positive value is returned unchanged.
- Child-backed test runners receive `testwait.Bound(t)`, including the shared
  fake ACPX harness and its fresh-runner replacement. The shared harness and
  blocking harness use the test deadline as stop grace; the unused 20 ms
  harness fallback is removed. Existing fake-clock cancellation assertions,
  close ordering, warning assertions and cleanup assertions remain intact.
- The opt-in real ACPX request helper receives the test so its stop grace uses
  `testwait.Bound(t)`. Its process-table inspection also uses that bound.
- The immediate-exit fixture emits its startup output before exiting with
  status 23. The nested waiter reads that first output before starting the
  measured span; the existing 9 s death-detection assertion remains. The test
  now calls `t.Parallel()` first.
- The console sink test waits for its first write with `testwait.Until`, also
  observing an early publish failure. Its serialized-output assertion remains.
  No Sequential Test, signal operation or written executable was added.

### Focused checks

- Initial regression check:
  `go test -race -count=1 -timeout=2m -run '^TestDisposableCleanupTimeoutDefaultsToTheProductionValue$' ./internal/agent`
  exited 1 before the runner implementation: unknown `cleanupTimeout` field
  and absent `disposableCleanupTimeout` method.
- After implementation:
  `go test -count=1 -race -shuffle=on -timeout=10m -v ./internal/agent`
  exited 0; package time 83.791 s; shuffle seed `1791867458742324000`.
  Raw output: `/tmp/roundfix-task02-agent-race.log`.
- `git diff --check` exited 0.
- The Task's declared Verification, `make verify`, and `make verify-changed`
  were not run: the Daemon owns those checks.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| A runner without the dependency keeps the production 5 s cleanup timeout | `TestDisposableCleanupTimeoutDefaultsToTheProductionValue` passed under `-race`, including parallel zero, negative and positive cases. The default cases assert the unchanged production constant. All three cleanup sites call the new method. |
| Fake ACPX children have no unasserted real 2 s or 5 s stop/cleanup bound | The shared harness, child-backed runner constructors and fresh-runner replacement receive `testwait.Bound(t)`. The cancellation tests retain their fake-clock assertions. The package race run exercises the selection, probe, fallback and sealed tests sharing those fixtures. |
| Cancel, cleanup, sealed, fixture-death and console-sink tests pass under `-race` | The package run passed `TestACPXRunCancelsPromptCooperatively`, `TestACPXRunClosesSessionAfterCancelGracePeriod`, `TestACPXRunCancellationCommandFailuresWarnAndContinue`, `TestProveExactSelectionCancelCleanup`, `TestProveExactSelectionTimeoutCleanup`, `TestACPXProbeCancellationStillClosesDisposableSession`, `TestACPXProbeFallbackCancellationStillClosesCandidateSession`, `TestSealedACPXPromptCancellationCancelsAndClosesSession`, `TestSpawnedFixtureDeathFailsFast`, and `TestConsoleDisplaySinkSerializesConcurrentPublish`. |

### Limits and follow-up

`TestRealACPXCommandOverrideRoundTripCancelCrashResume` compiled in the race run
but kept its existing opt-in skip because `ROUNDFIX_REAL_ACPX` was not enabled;
no live ACPX round trip is claimed. The initial shell sandbox rejected the
worktree's symlinked writable root; focused work and checks used its approved
resolved filesystem path. The assigned task file's pre-existing
`status: in_progress` delta is preserved. No other Task file or Task Graph was
edited, and no commit, push or pull request was created. No out-of-slice work
was identified.

## Carry-forward provenance

- Source Run: `run_20261013T045342Z_a78df1f6c39a67e2`
- Source commit: `c3857f7877ce607f863c4a62656601d772a9d9fe`
