---
task: task_02
spec: 0171-a-deterministic-suite-under-load
status: completed
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

## Result

Implemented the Implement-test wait migration without changing production
behavior. Every named helper now delegates to `testwait.Until` or
`testwait.Poll`; fixed wait budgets, timer selects and literal 90-second waits
were removed. The two tests that run `runImplementCommandAsync` pass its result
channel through every intervening Agent-start, journal and file wait. An ended
command is formatted with its exit code, stdout and stderr for `testwait`'s
immediate failure diagnostic. The three bootstrap tests derive their
non-subject product timeout from `testwait.Bound(t)`. `attachDetachBudget` and
the attach follow loop remain unchanged, and no production or governed file
was edited.

Focused evidence by acceptance criterion:

- Wait migration and source constraints: `rtk rg -n
  'implementWaitBudget|detachStartupBudget|time\.After\(|time\.NewTimer\(|bootstrap_timeout:
  1s|90[[:space:]]*\*[[:space:]]*time\.Second|command, "1s"\)'
  internal/cli/implement_test.go` found no matches (exit 1), while `rtk rg -n
  'testwait\.(Until|Poll|Bound)|attachDetachBudget'
  internal/cli/implement_test.go` found all nine migrated waits, the three
  deadline-derived bootstrap timeouts, and the unchanged attach budget (exit
  0). `rtk git diff --check` also exited 0.
- Queued cancellation and bootstrap behavior: `rtk env
  GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -ldflags '-X
  roundfix/internal/app.Version=dev' -run
  '^(TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks|TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork|TestRunImplementBootstrapFailureEndsFailedBeforeAgentWork|TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification)$'
  ./internal/cli` passed (exit 0, 1.756 s). A focused stress variant with the
  same selection and link flag at `-count=2 -cpu 1,4` also passed (exit 0,
  2.847 s). The development-version link flag prevents the live release lookup
  assigned to Task 04 and changes no repository file.

Additional focused checks:

- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -ldflags '-X
  roundfix/internal/app.Version=dev' -run
  '^(TestRunImplementDetachPrintsReportAndCompletesRun|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow)$'
  ./internal/cli` passed (exit 0, 3.177 s), exercising the other converted
  line, process, file, Run-state, outcome and journal wait paths.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache make verify-incremental`
  passed `go vet ./...` and began the repository-wide test phase, then the
  sandbox blocked the pre-existing live `api.github.com` release lookup. That
  lookup is explicitly owned by Task 04, so this Task did not alter or bypass
  it.

The Daemon-owned command in `## Verification`, including its `-count=5`
stress, was not run.
