---
task: task_04
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
type: backend
complexity: high
---

# Task 04: Retry transient Review Source failures and project waits

## Overview

Retry only positively typed transient Review Source failures within existing
timeout and Run Budget bounds while projecting every wait phase and retry
episode. Fake-clock coverage proves bounded recovery, exhaustion, and graceful
Stop Request interruption without real sleeps.

## Requirements

1. MUST retry only typed transient errors and never infer retryability from
   output text.
2. MUST reuse the configured poll interval, Review Source timeout, and Run
   Budget without a new retry setting.
3. MUST publish one retry-start event and one recovery or exhaustion event per
   episode.
4. MUST expose phase, expected head, start, deadline, Evidence, and retry status
   for both Review Source wait phases.
5. MUST publish progress only on phase, Evidence, or retry changes.
6. MUST let the Store-backed Stop Request interrupt retry sleep and win over
   another attempt.
7. MUST keep authentication, validation, and permanent failures terminal.

## Subtasks

- [x] Add bounded transient retry episodes to watch.
- [x] Reuse existing time and budget boundaries.
- [x] Add wait-phase and deadline projection.
- [x] Deduplicate unchanged progress and episode events.
- [x] Integrate graceful Stop Request interruption.
- [x] Add recovery, exhaustion, permanent, and cancellation matrices.

## Acceptance Criteria

- [x] One transient failure retries and can recover without ending the Run.
- [x] Exhaustion occurs at the existing timeout or Run Budget boundary.
- [x] Each episode contains at most one start and one terminal episode event.
- [x] Permanent failure performs zero retry sleeps.
- [x] Both wait phases expose the documented head, deadline, Evidence, and retry
      fields.
- [x] Unchanged polling appends no duplicate progress event.
- [x] Stop Request during retry sleep starts no later Review Source call.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/watch/watch.go`
- interface: `internal/watch/watch_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/event_test.go`

## Verification

- `rtk go test ./internal/watch -run 'TestRun.*(Transient|Retry|WaitPhase|StopRequest)' -count=1`
  — expected: fake-clock recovery, exhaustion, projection, deduplication, and
  stop cases pass.
- `rtk go test ./internal/cli ./internal/runevent -run 'Test.*(ReviewRetry|WaitingForReview|ReviewStatusEvent)' -count=1`
  — expected: CLI and stream expose bounded wait evidence.
- `rtk go test -race ./internal/watch ./internal/cli -run 'Test.*(Transient|Retry|StopRequest)' -count=1`
  — expected: retry and cancellation have no race or leaked waiter.

## References

- `_prd.md` → Goals 2–3; User Stories 2–3; Core Features 5–6; User Experience;
  Success Metrics.
- `_techspec.md` → API Contracts: Transient retry and Wait projection; Testing
  Approach; Build Order 4.
- `../0037-terminal-outcome-integrity/_techspec.md` → stop-aware Review Source
  waits.

## Result

Implemented typed transient retry episodes in both Review Source wait phases.
The watch loop now clamps each poll to the earlier Review Source or Run Budget
deadline, rejects permanent errors without sleeping, publishes one start and
one recovered or exhausted event per episode, and checks the Store-backed Stop
Request before any later Review Source call.

`daemon.review_status` now projects the wait phase, expected head, start,
deadline, Evidence, and retry status. The same deduplicated projection drives
CLI progress, so unchanged polls append neither another Run Event nor another
progress line.

Verification:

- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/watch -run 'TestRun.*(Transient|Retry|WaitPhase|StopRequest)' -count=1`
  — passed: 15 tests.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/cli ./internal/runevent -run 'Test.*(ReviewRetry|WaitingForReview|ReviewStatusEvent)' -count=1`
  — passed: 3 tests in 2 packages.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test -race ./internal/watch ./internal/cli -run 'Test.*(Transient|Retry|StopRequest)' -count=1`
  — passed: 22 tests in 2 packages.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/watch ./internal/runevent -count=1`
  — passed: 70 tests in 2 packages.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/cli -run 'Test(RunWatchPrintsDeterministicStdoutReport|RunWatchTimeoutOffersManualReviewWithoutFetching|RunWatchMissingHeadCheckEndsCleanUnverified|RunWatchNoAgentConsoleSuppressesAgentDisplayOnly|WatchRunJournalsOrderedLoopNarrative|AttachRendersWatchDaemonEventsInTimeline)$' -count=1`
  — passed: 10 tests.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

Acceptance evidence:

- Transient recovery: pre-fetch and Merge-Ready recovery tests both reach
  Clean after one typed transient failure.
- Existing bounds: timeout tests stop before a call at the deadline, and the
  Run Budget test projects and exhausts at the earlier budget deadline.
- Episode cardinality: repeated transient failures produce only `started` plus
  one `recovered` or `exhausted` terminal event.
- Permanent failures: authentication-like text containing
  `temporary Review Source failure` and validation failures remain terminal
  and record zero retry sleeps, proving retryability is type-based.
- Wait projection and deduplication: both `WaitingForReview` and
  `WaitingForReviewCheck` expose all documented fields; unchanged Evidence
  produces no duplicate persisted or direct progress update.
- Stop Request priority: a Stop Request observed after retry sleep returns
  Stopped with exactly one failed Merge-Ready access and no later Review Source
  call.

Follow-up evidence outside this Task's slice:

- A diagnostic full `internal/cli` run reached 790 passing tests but retained
  two unrelated failures: a legacy migration assertion expects schema version
  10 while the current Store schema is 11, and the Unix owner-process test
  could not read the spawned process start time. Task-related CLI regression
  tests pass in the focused commands above.
