---
task: task_02
spec: 0171-a-deterministic-suite-under-load
status: pending
type: backend
complexity: high
---

# Task 02: Implement tests wait on events

## Overview

Spec 0163's QA gate failed twice on Implement tests that pass in isolation:
`TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks` waited
90 s for a second Agent start without looking at the command it had started, and
`TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork` configures
a one-second bootstrap timeout it does not examine. This Task moves
`internal/cli/implement_test.go` onto `internal/testwait`.

## Requirements

1. MUST remove `implementWaitBudget`, `detachStartupBudget` and every literal
   90-second wait from `internal/cli/implement_test.go`, and route every wait —
   `waitImplementCommandResult`, `waitImplementAgentStarts`,
   `waitForImplementJournal`, `waitForFile`, `waitForFileContains`,
   `waitForRunState`, `waitForCleanOutcomeEvent`, `readLineWithTimeout` and
   `waitProcessForTest` — through `testwait.Until` or `testwait.Poll`. No
   `time.After` or `time.NewTimer` may remain in the file.
2. MUST make every wait that runs while `runImplementCommandAsync` holds a
   command in flight watch that command's result channel, so a command that
   ends first fails the wait at once with its exit code, stdout and stderr.
   This covers at least the waits in
   `TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks`.
3. MUST replace `bootstrap_timeout: 1s` in
   `TestRunImplementBootstrapFailureEndsFailedBeforeAgentWork`,
   `TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification` and
   `TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork` with a
   timeout derived from `testwait.Bound(t)`, keeping each test's assertions on
   exit code, stderr and Run state unchanged.
4. MUST leave `attachDetachBudget` and the attach follow loop it bounds
   unchanged; Spec 0124 owns them.
5. MUST keep every top-level test function name recorded in
   `docs/references/coverage-record.json`, and MUST NOT edit that file or
   `internal/cli/cli_test.go`, both governed.
6. MUST change no production file: the behaviour under test is unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] `implement_test.go` holds no `time.After`, `time.NewTimer`,
      `implementWaitBudget` or `bootstrap_timeout: 1s`, and it uses `testwait`.
- [ ] The queued-cancellation test and the three bootstrap tests pass every
      iteration of `-count=5 -cpu 1,4`.

## Context

- interface: `internal/cli/implement_test.go`

## Verification

- `out="$(go test -count=5 -cpu 1,4 -v -run "^(TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks|TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork|TestRunImplementBootstrapFailureEndsFailedBeforeAgentWork|TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork TestRunImplementBootstrapFailureEndsFailedBeforeAgentWork TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done; grep -q "testwait\." internal/cli/implement_test.go && ! grep -q "time.After(" internal/cli/implement_test.go && ! grep -q "time.NewTimer(" internal/cli/implement_test.go && ! grep -q "implementWaitBudget" internal/cli/implement_test.go && ! grep -q "bootstrap_timeout: 1s" internal/cli/implement_test.go && ! grep -q "command, \"1s\")" internal/cli/implement_test.go` — expected: exit 0; before this Task `implement_test.go` does not use `testwait` and still holds `implementWaitBudget`, so the command fails.

## References

- [_prd.md](_prd.md) — Core Feature 2; Success Metrics 1-2
- [_techspec.md](_techspec.md) — The Implement tests
