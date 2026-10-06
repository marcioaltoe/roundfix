---
task: task_03
spec: 0240-a-lost-rollout-is-infrastructure
status: completed
type: backend
complexity: high
---

# Task 03: The Daemon recovers a Lost Rollout and the queue parks an unrecovered one without counting its retry

## Overview

Today a Lost Rollout settles the Task `failed`, or halts a QA gate before its
report. The Delivery Queue then parks `run-unresolved`, and each attempt
spends a retry. This Task makes the Agent Session owner recover the loss in a
new Agent Session without repair. The Fallback Chain takes the Task before its
First Handoff or before a pending QA report. Each recovery is recorded in the
Run and, for a QA fallback, in the seeded report. A third loss settles the
Task as runtime infrastructure, and the queue parks that Run as
`runtime-infrastructure` with an uncounted retry. It answers the Backlog Entry
"A lost Codex rollout fails the Task and spends a retry" of 2026-10-06. It is
verifiable with the Daemon's fake runner and the delivery fakes.

## Requirements

1. MUST handle a Lost Rollout in `agentSessionOwner.Run` per Invariants 5 to 8,
   using `agent.DescribeLostRollout` from task_02. The owner gains the fields
   of `_techspec.md` → Interfaces. `handedOff` is set when a prompt returns
   without error. A recovery ends the lost session and never resumes its id.
   Then it either activates the next Fallback Selection (`fallback`) or
   re-activates the same selection under `<name>-rollout-NN`
   (`new_session`), and it sends the prompt of Invariant 8. It never starts a
   Verification Feedback repair for the lost turn.
2. MUST publish, for each loss, the `daemon.task` Run Event and the stderr line
   of Invariant 9. A `fallback` recovery also publishes the existing
   `agent_selection_fallback` notification, and `selectionReasonCode` returns
   `rollout_lost` for a Lost Rollout. The selection attempt is persisted as
   failed, as `fallbackAfterSelectionFailure` does.
3. MUST stop after `maxLostRolloutRecoveries` recoveries per owner, per
   Invariant 10. The third loss publishes `recovery: exhausted` and returns an
   error whose text starts `runtime infrastructure: lost rollout`, which
   `agentFailureReason` keeps as the Task's failure reason. In the QA step it
   settles the QA Task `failed` with that reason instead of halting the cycle,
   so the Run ends Unresolved.
4. MUST give the QA owner a `reportPending` predicate in `task_engine.go`. It
   reports true while the seeded report at the QA step's report path reads
   `verdict: pending`. Before a QA fallback session starts, the Daemon appends
   the section of Invariant 11 to that report.
5. MUST make the delivery workflow set `RunResult.RuntimeInfrastructure` for an
   Unresolved Run whose journal holds an exhausted `rollout_lost` event, read
   through `Store.RunEventsOfKinds`, per Invariant 12. `runCandidate` parks
   `runtime-infrastructure: <scope_id> lost its rollout at <step>` before the
   `qa-environment-partial` and `run-unresolved` decisions.
6. MUST classify `runtime-infrastructure` in `ClassifyPark` as `environment`
   with the next action of Invariant 13. A Delivery Retry from that blocker
   skips the retry-limit refusal in `internal/delivery/engine.go` and in
   `RetryDeliveryQueueItem`, leaves `retry_count` unchanged, and re-enters as a
   `run-unresolved` retry does.
7. MUST add the tests named in Verification:
   - `internal/daemon/lost_rollout_recovery_test.go`, with the package's fake
     runner returning an error that `DescribeLostRollout` places:
     - a loss on the first turn with a fallback configured runs the next prompt
       on fallback 1 under a `-rollout-01` session name, with no Verification
       Feedback prompt, and the Task completes;
     - a loss on the Verification Feedback turn after a handoff continues on
       the same selection with the Task prompt, the notice and the repair
       prompt;
     - a loss on the first turn with no fallback continues on the same
       selection;
     - a QA loss while the report is pending runs on fallback 1 and leaves the
       section with `retry_spent: false` in the report;
     - three losses settle the Task `failed` with the
       `runtime infrastructure: lost rollout` reason and publish `exhausted`;
     - every case asserts its `rollout_lost` events.
   - `internal/delivery/runtime_infrastructure_test.go`: the park from an
     executor result, its Park Class and next action, and a retry at the
     queue's retry limit that succeeds without raising `retry_count`.
   - `internal/cli/deliver_runtime_infrastructure_test.go`: `runResult` reads
     an exhausted event from a temporary store's journal and ignores a
     recovered one.
8. MUST NOT change the runner, the review command, Runs without Agent
   Selection Profiles, selection failures that are not Lost Rollouts, the
   light-tier escalation, or any existing test's expectations.

## Subtasks

- [ ] Recognize and recover a Lost Rollout in the owner, with its prompt and
      session name.
- [ ] Publish the Run Event, the notice and the fallback notification.
- [ ] Bound the recoveries and settle the third as runtime infrastructure.
- [ ] Add the QA report predicate and section.
- [ ] Park `runtime-infrastructure` and retry it without counting.
- [ ] Add the Daemon, delivery and CLI tests.

## Acceptance Criteria

- [ ] A Lost Rollout before the First Handoff runs on the next Fallback
      Selection; after it, on the same selection; neither starts a repair.
- [ ] A QA fallback is recorded in the seeded report.
- [ ] A third loss parks the queue item as `runtime-infrastructure`, and its
      retry leaves the retry count unchanged.
- [ ] Every existing Daemon, delivery and CLI test passes unchanged.

## Context

- creates: `internal/daemon/lost_rollout_recovery_test.go`
- creates: `internal/delivery/runtime_infrastructure_test.go`
- creates: `internal/cli/deliver_runtime_infrastructure_test.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/store/delivery.go`
- interface: `internal/cli/deliver_workflow.go`
- instruction: `internal/agent/protocol_failure.go`
- instruction: `internal/daemon/task_engine_test.go`
- instruction: `internal/delivery/retry_test.go`
- instruction: `internal/delivery/park_class_test.go`
- instruction: `docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md`
- instruction: `docs/adr/0114-opening-an-agent-session-is-not-agent-work.md`
- instruction: `docs/adr/0057-daemon-exclusively-owns-implement-task-status.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLostRolloutBeforeFirstHandoffActivatesTheFallback|TestLostRolloutAfterFirstHandoffResumesInANewSession|TestLostRolloutWithoutFallbackResumesTheSameSelection|TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt|TestThirdLostRolloutSettlesRuntimeInfrastructure)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLostRolloutBeforeFirstHandoffActivatesTheFallback TestLostRolloutAfterFirstHandoffResumesInANewSession TestLostRolloutWithoutFallbackResumesTheSameSelection TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt TestThirdLostRolloutSettlesRuntimeInfrastructure; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestRuntimeInfrastructureParksTheItem|TestRuntimeInfrastructureRetryIsNotCounted|TestClassifyParkRuntimeInfrastructure|TestDeliveryRunResultReadsAnExhaustedLostRollout)$' ./internal/delivery ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRuntimeInfrastructureParksTheItem TestRuntimeInfrastructureRetryIsNotCounted TestClassifyParkRuntimeInfrastructure TestDeliveryRunResultReadsAnExhaustedLostRollout; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four tests do not exist, so the command fails.

## References

- `_prd.md` → Core Features 2-4; User Stories 1-3; Success Metric 2; Success Metric 3; Success Metric 4; Goals 2-4
- `_techspec.md` → Interfaces; Invariants 5 to 13; API Contract 2; API Contract 3; API Contract 4; Testing Approach 2; Testing Approach 3; Build Order 3
- ADR-0245
- ADR-0114
- ADR-0057


## Result

Implemented this Task's recovery, QA and Delivery Queue slice. Status remains
Daemon-owned; authored Verification has not been run in this turn.

- First Handoff boundary: the session owner detects `DescribeLostRollout`
  before ordinary selection failure handling, ends the lost session, persists
  the failed selection attempt and recovers at most twice. Before handoff it
  takes the next fallback; after handoff or without a fallback it retains the
  selection. Replacement names use `-rollout-NN`, and prompts preserve the
  initial Task prompt, recovery notice and any lost Verification Feedback
  prompt. The first-turn fake emits Agent output before the loss, proving that
  work-started does not prevent this recovery. The new Daemon tests assert
  selection persistence, distinct session names, recovery events and Task
  settlement without adding a repair for a lost turn.
- QA fallback: the QA owner reads the seeded report's pending verdict. Before
  fallback it appends the runtime fallback section without replacing the
  frontmatter or rows. The QA test inspects the pending report before the
  fallback runs and the final report for `retry_spent: false`. An additional
  QA exhaustion test proves the third loss settles failed with the runtime
  infrastructure reason instead of halting the cycle.
- Exhaustion and retry: the third loss publishes `rollout_lost` with
  `recovery: exhausted`. Delivery reads this event through
  `Store.RunEventsOfKinds`, parks `runtime-infrastructure` ahead of QA partial
  and generic unresolved decisions, and classifies it as `environment` with
  the specified next action. The real-store retry test reaches the retry
  limit with an ordinary retry, then proves the infrastructure retry carries
  forward and re-enters Running without increasing `retry_count`. CLI journal
  tests distinguish exhausted events from both recovered outcomes.
- Existing expectations: no existing test files or expectations changed.
  The affected Daemon, Delivery, CLI and store suites have passing results
  from the checks below. Runner, review execution, legacy Runs, ordinary
  selection failure and light-tier escalation code remain outside this diff.

Focused implementation checks:

- `GOCACHE=/private/tmp/roundfix-task03-cache rtk proxy go test ./internal/daemon ./internal/delivery ./internal/cli -run 'LostRollout|RuntimeInfrastructure' -count=1`
  — passed.
- `GOCACHE=/private/tmp/roundfix-task03-cache rtk proxy go test ./internal/daemon -run 'LostRollout' -count=1`
  — passed after adding persistence and QA exhaustion assertions.
- `GOCACHE=/private/tmp/roundfix-task03-cache rtk proxy go test ./internal/daemon ./internal/delivery ./internal/cli ./internal/store -count=1`
  — Delivery and store passed. Daemon and CLI repository guards detected a
  concurrent edit to the new Daemon test; CLI also reported process-table
  permission failures in two existing force-stop integration tests.
- `GOCACHE=/private/tmp/roundfix-task03-cache rtk proxy go test ./internal/daemon -count=1`
  — passed with stable files after strengthening the first-turn output case
  (46.568s).
- `GOCACHE=/private/tmp/roundfix-task03-cache rtk proxy go test ./internal/cli -count=1`
  — passed with stable files and elevated process-table access (184.010s).
- `rtk proxy git -c core.fsmonitor=false diff --check` — passed.

No commits, pushes or Pull Requests were created. The Task Graph and other
Task files were not edited. The incoming `status: in_progress` is preserved.
