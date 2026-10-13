---
task: task_03
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
status: pending
type: test
complexity: medium
---

# Task 03: Test helpers that run a child Roundfix binary name its Signal Exit

## Overview

Seven test call sites reduce a child Roundfix binary's end through
`exitCodeFromWait`, which reports a child that a signal ended as exit `1`, and
two helpers return `ExitError.ExitCode()`, which reports it as `-1`; neither
names the signal (`_techspec.md` → Current behavior). This Task adds
`testfixture.ChildExitCode`, moves every one of them onto it, and proves with
a fixture that kills itself that `runRoundfixBinaryMacro` now names the
signal. It is verifiable on its own through the sabotage test and a search
that finds no remaining call.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [a test child is killed by a signal from an unknown sender](references/2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md):
   a test that runs a child Roundfix binary reports a Signal Exit by name.
2. MUST add `internal/testfixture/child_exit.go` with `ChildExitCode` as
   `_techspec.md` → Interfaces and Invariant 2 state, built on
   `childexit.From` and `Exit.String`.
3. MUST move onto `ChildExitCode`, following Invariant 3, every site
   `_techspec.md` → Current behavior lists: `runCLIHelper` and
   `runRoundfixBinaryMacro` in `implement_test.go`, the `run` closure of
   `reopen_test.go`, the call in `reopen_late_dependency_test.go`,
   `waitForReviewDisposeProcesses` in `review_dispose_lock_test.go`, the call
   in `version_freshness_isolation_test.go`, `assertScriptFixture` in
   `script_fixture_test.go`, `runBaselineReleaseCLI` in
   `baseline_release_gate_test.go` and `baselineCitationGoRun` in
   `baseline_update_citations_test.go`. A helper that returns standard error
   appends the note as its last line; a site that only compares a code puts
   the note in its failure message. No assertion's expected value changes.
4. MUST leave no `_test.go` file in the repository calling
   `exitCodeFromWait`; the production function and its `detach.go` callers
   stay unchanged.
5. MUST add `internal/cli/child_signal_report_unix_test.go`, with a darwin or
   linux build constraint, holding
   `TestRoundfixBinaryHelperNamesTheSignalThatKilledItsChild`: it compiles,
   with `testfixture.FixtureBinary`, a program that sends itself `SIGKILL`
   (its own PID only, ADR-0266), runs it through `runRoundfixBinaryMacro` in
   a temporary home and directory, and asserts code `-1` and the last line of
   the returned standard error
   `child was killed by signal SIGKILL (9), no core dump`. It calls
   `t.Parallel()` first (ADR-0259) and waits on no clock (ADR-0271).
6. MUST change `internal/cli/baseline_release_gate_test.go` only inside
   `runBaselineReleaseCLI`, the Governed Path `_authorization.md` bounds.

## Subtasks

- [ ] Add `ChildExitCode`.
- [ ] Move the nine sites onto it.
- [ ] Add the sabotage test through `runRoundfixBinaryMacro`.

## Acceptance Criteria

- [ ] The sabotage test passes and names `SIGKILL (9)` and the core state.
- [ ] A repository search finds no `exitCodeFromWait` call in a test file.
- [ ] Every swept test file still compiles and its tests keep their expected values.

## Consumers

- `testfixture.ChildExitCode` and the swept helpers' returned standard error — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- creates: `internal/testfixture/child_exit.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/reopen_test.go`
- interface: `internal/cli/reopen_late_dependency_test.go`
- interface: `internal/cli/review_dispose_lock_test.go`
- interface: `internal/cli/version_freshness_isolation_test.go`
- interface: `internal/cli/script_fixture_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- interface: `internal/cli/baseline_update_citations_test.go`
- creates: `internal/cli/child_signal_report_unix_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestRoundfixBinaryHelperNamesTheSignalThatKilledItsChild$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestRoundfixBinaryHelperNamesTheSignalThatKilledItsChild ' || { printf 'missing pass: TestRoundfixBinaryHelperNamesTheSignalThatKilledItsChild\n' >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, so no PASS line is printed.
- `matches="$(grep -rln 'exitCodeFromWait(' --include='*_test.go' internal cmd skills)"; test -z "$matches" || { printf 'test files still call exitCodeFromWait:\n%s\n' "$matches" >&2; exit 1; }; for f in internal/cli/implement_test.go internal/cli/reopen_test.go internal/cli/reopen_late_dependency_test.go internal/cli/review_dispose_lock_test.go internal/cli/version_freshness_isolation_test.go internal/cli/script_fixture_test.go internal/cli/baseline_release_gate_test.go internal/cli/baseline_update_citations_test.go; do grep -q 'testfixture.ChildExitCode(' "$f" || { printf '%s does not use testfixture.ChildExitCode\n' "$f" >&2; exit 1; }; done; go vet ./internal/testfixture && go test -count=1 -run '^$' ./internal/cli` — expected: exit 0; before this Task seven calls remain in six test files.

## References

- [_prd.md](_prd.md) — Goal 2; User Story 1; Core Feature 3; Success Metric 1; Success Metric 2
- [_techspec.md](_techspec.md) — API Contract 3; Interfaces; Invariant 2; Invariant 3; Testing Approach; Build Order 3
- [2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md](references/2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md)
- ADR-0266; ADR-0259; ADR-0271
