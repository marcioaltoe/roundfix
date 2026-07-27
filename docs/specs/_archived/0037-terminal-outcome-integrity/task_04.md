---
task: task_04
spec: 0037-terminal-outcome-integrity
status: completed
type: backend
complexity: medium
---

# Task 04: Interrupt every Review Source wait on Stop Request

## Overview

Make the Run Database Stop Request visible at every Review Source wait
boundary. A graceful stop interrupts status, retry, quiet-period, and
Merge-Ready waits by the next polling boundary and starts no later Review
Source or repository mutation.

## Requirements

1. MUST expose Stop Request observation through a narrow context-aware source.
2. MUST supply the Store-backed source to every operational watch Run.
3. MUST check the source before each Review Source status access.
4. MUST check again after every interruptible wait, including quiet, retry, and
   Merge-Ready sleeps.
5. MUST return the existing stop classification when a request is observed.
6. MUST not perform another fetch, check, artifact write, commit, push, or
   Review Source mutation after observation.
7. MUST propagate Store observation failures with Run and operation context.

## Subtasks

- [x] Add the Stop Request source contract to watch dependencies.
- [x] Wire the Store-backed source into operational watch Runs.
- [x] Cover every status and sleep boundary.
- [x] Preserve the existing stopped classification and CLI exit behavior.
- [x] Add fake-clock tests for each waiting phase.
- [x] Prove no downstream operation starts after observation.

## Acceptance Criteria

- [x] A request during status wait reaches Stopped by the next poll.
- [x] Requests during quiet period, transient retry sleep, and Merge-Ready wait
      each interrupt their respective boundary.
- [x] No later Review Source call or repository mutation occurs.
- [x] Operational runs always provide the Store source; only isolated tests may
      omit it.
- [x] Store read failure is distinguishable from a requested stop.
- [x] Existing Run Budget and timeout behavior remains unchanged without a stop.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/store/store.go`

## Verification

- `rtk go test ./internal/watch -run 'TestRun.*StopRequest' -count=1`
  — expected: every wait boundary observes Stop Request and starts no later
  operation.
- `rtk go test ./internal/cli -run 'TestRunWatch.*StopRequest' -count=1`
  — expected: Store-backed stop observation reaches the existing Stopped CLI
  contract.
- `rtk go test -race ./internal/watch ./internal/cli -run 'Test.*StopRequest' -count=1`
  — expected: cancellation and polling remain race-free.

## References

- `_prd.md` → Goal 3; User Story 2; Core Feature 4; User Experience; Success
  Metrics.
- `_techspec.md` → Interfaces: StopRequestSource; API Contracts: Graceful Stop
  Requests; Build Order 4.
- `../../adr/0022-stop-requests-travel-through-the-run-database.md` → durable
  Stop Request transport.

## Result

Implemented a context-aware Stop Request source in the watch dependencies and
wired the operational Watch Run to its Run Database Store. The watch loop now
observes Stop Requests before and after Review Source status and Merge-Ready
accesses and after status-poll, quiet-period, transient-retry, and Merge-Ready
sleeps. Observation returns the stopped classification before any later fetch,
check, resolve, artifact, commit, push, or Review Source mutation can start.
Store observation failures retain their cause and add the Run ID plus boundary
operation.

Acceptance evidence:

- Status wait: `TestRunStopRequestDuringStatusWaitStopsAtNextPoll` observed the
  request after one fake-clock poll, returned Stopped with zero fetched Rounds,
  and made no second status, fetch, or resolve call.
- Quiet, retry, and Merge-Ready waits:
  `TestRunStopRequestDuringQuietPeriodStopsBeforeFetch`,
  `TestRunStopRequestDuringTransientRetryStopsBeforeNextCheck`, and
  `TestRunStopRequestDuringMergeReadyWaitStopsBeforeNextCheck` each observed
  the request immediately after the relevant fake-clock sleep and proved no
  later operation began.
- Operational source and CLI contract:
  `TestRunWatchStopRequestBeforeAgentMarksStopped` recorded a real Run Database
  Stop Request without canceling the context, reached Stopped with exit code
  zero, and performed no fetch or Agent work. The sole operational
  `watch.Run` call supplies `StopRequests: runStore`; isolated package tests
  may leave the source nil.
- Failure distinction:
  `TestRunStopRequestSourceFailureIncludesRunAndOperation` preserved the Store
  error chain, reported Failed rather than Stopped, named `run_123` and the
  status operation, and made no Review Source call.
- Unchanged safeguards:
  `TestRunWithoutStopRequestKeepsRunBudgetBehavior` retained BudgetExceeded
  behavior with a non-requesting source, and the complete watch and CLI package
  suites passed.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test ./internal/watch -run 'TestRun.*StopRequest' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test ./internal/cli -run 'TestRunWatch.*StopRequest' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -race ./internal/watch ./internal/cli -run 'Test.*StopRequest' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test ./internal/watch ./internal/cli -count=1`
  — passed with process-control permission. The first sandboxed attempt reached
  the unrelated detached-process test and failed with `operation not
  permitted`; the identical permitted rerun passed both packages.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

Follow-ups: none.
