---
task: task_03
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
status: completed
type: backend
complexity: high
---

# Task 03: The Daemon recovers a Provider System Error through the Fallback Chain and the queue parks an unrecovered one without counting its retry

## Overview

After task_02 a Provider System Error reaches the Agent Session owner as a
failure, but the owner offers the Fallback Chain only to a Selection Failure
before Agent work, so a mid-turn capacity failure would still end the Task as
a plain failure and the Delivery Queue would count its retry. This Task
recovers it with the eligibility the owner already applies to a Lost Rollout,
records it, settles an unrecovered one as runtime infrastructure, and parks
the queue item as `runtime-infrastructure`. It answers the Backlog Entry
[a model-at-capacity reply is accepted as a finished turn](references/2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md)
of 2026-10-13. It is verifiable on its own with the Daemon's scripted runner
and the delivery workflow's event reader.

## Requirements

1. MUST check `agent.DescribeProviderSystemError` in `agentSessionOwner.Run`
   after `DescribeLostRollout` and before the Selection Failure branch, and
   recover per `_techspec.md` → Invariants 7 and 8: the next Fallback
   Selection when the scope is `task` before its First Handoff or `qa` with
   the seeded report pending, `exhausted` otherwise, with no same-selection
   retry and no use of the Lost Rollout budget. A `review` scope returns the
   error unchanged.
2. MUST record each recovery per Invariant 9: the closed session, the attempt
   persisted `failed` with reason code `provider_system_error` (through
   `selectionReasonCode`), the `daemon.task` event of API Contract 2, the
   progress line, and the fallback prompt with `providerSystemErrorNotice`.
3. MUST append the QA report lines of Invariant 10 before a QA fallback session
   starts, leaving `verdict: pending`.
4. MUST return `providerSystemErrorExhaustedError` for an exhausted recovery,
   name it in `agentFailureReason`, and make the QA step settle the QA Task
   `failed` with its text and return a pending verdict without halting, as it
   does for an exhausted Lost Rollout (Invariant 11).
5. MUST make `runResult` in `internal/cli/deliver_workflow.go` read an
   exhausted `provider_system_error` event into `RuntimeInfrastructure` with
   the text of API Contract 3, beside the `rollout_lost` reader, whose text
   stays as it is.
6. MUST add `internal/daemon/provider_system_error_recovery_test.go` with the
   daemon tests named in Verification, built like
   `lost_rollout_recovery_test.go` (a scripted runner whose turn publishes an
   Agent message and returns a Provider System Error built with
   `ProtocolFailure.ProviderSystemError`):
   - `TestProviderSystemErrorBeforeFirstHandoffActivatesTheFallback`: a first
     turn that fails activates the fallback selection, and the Task completes
     on it with one `fallback` event.
   - `TestProviderSystemErrorAfterFirstHandoffSettlesRuntimeInfrastructure`: a
     Verification Feedback turn that fails settles the Task `failed` with the
     exhausted text and opens no new session.
   - `TestProviderSystemErrorWithoutFallbackSettlesRuntimeInfrastructure`: a
     profile with no fallback settles `failed` with one `exhausted` event and
     no second turn on the same selection.
   - `TestQAProviderSystemErrorBeforeTheReportFallsBackAndRecordsIt`: the QA
     fallback runs the gate and the report carries the three lines of
     Invariant 10.
   - `TestQAProviderSystemErrorExhaustionSettlesInsteadOfHalting`: with no
     fallback, the QA Task settles `failed` and the cycle returns no error.
   Each calls `t.Parallel()` first (ADR-0259) and waits on no elapsed time
   (ADR-0271).
7. MUST add `internal/cli/deliver_provider_system_error_test.go` with
   `TestDeliveryRunResultReadsAnExhaustedProviderSystemError`: an exhausted
   event yields the API Contract 3 text, a `fallback` event yields none, and
   an exhausted `rollout_lost` event still yields its own text. It calls
   `t.Parallel()` first.
8. MUST keep every existing test in `internal/daemon/lost_rollout_recovery_test.go`
   and `internal/cli/deliver_runtime_infrastructure_test.go` passing
   unchanged, and MUST NOT change the delivery engine, the park classes, the
   retry counting or any Lost Rollout behavior.

## Subtasks

- [ ] Recover a Provider System Error in the owner with the Lost Rollout
      eligibility and no same-selection retry.
- [ ] Record the attempt, the Run Event, the progress line and the QA report
      lines.
- [ ] Settle an exhausted recovery in Tasks and QA.
- [ ] Read the exhausted event in the delivery workflow.
- [ ] Add the daemon and delivery workflow tests.

## Acceptance Criteria

- [ ] A Provider System Error before the First Handoff, or in QA before its
      report, runs the next Fallback Selection and is recorded once.
- [ ] Any other Provider System Error settles `failed` as runtime
      infrastructure without a same-selection retry.
- [ ] The Delivery Queue parks an unrecovered one as
      `runtime-infrastructure` with its detail.
- [ ] Every Lost Rollout test passes unchanged.

## Context

- creates: `internal/daemon/provider_system_error_recovery_test.go`
- creates: `internal/cli/deliver_provider_system_error_test.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/deliver_workflow.go`
- instruction: `internal/daemon/lost_rollout_recovery_test.go`
- instruction: `internal/cli/deliver_runtime_infrastructure_test.go`
- instruction: `internal/agent/lost_rollout.go`
- instruction: `.agents/skills/roundfix/references/deliver.md`
- instruction: `docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md`
- instruction: `docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md`
- instruction: `docs/adr/0114-opening-an-agent-session-is-not-agent-work.md`
- instruction: `docs/adr/0057-daemon-exclusively-owns-implement-task-status.md`

## Consumers

- The `daemon.task` `provider_system_error` Run Event read by the delivery workflow — consumers: `./internal/daemon`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -v -run '^(TestProviderSystemErrorBeforeFirstHandoffActivatesTheFallback|TestProviderSystemErrorAfterFirstHandoffSettlesRuntimeInfrastructure|TestProviderSystemErrorWithoutFallbackSettlesRuntimeInfrastructure|TestQAProviderSystemErrorBeforeTheReportFallsBackAndRecordsIt|TestQAProviderSystemErrorExhaustionSettlesInsteadOfHalting|TestLostRolloutBeforeFirstHandoffActivatesTheFallback|TestLostRolloutWithoutFallbackResumesTheSameSelection|TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt|TestQALostRolloutExhaustionSettlesInsteadOfHalting)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestProviderSystemErrorBeforeFirstHandoffActivatesTheFallback TestProviderSystemErrorAfterFirstHandoffSettlesRuntimeInfrastructure TestProviderSystemErrorWithoutFallbackSettlesRuntimeInfrastructure TestQAProviderSystemErrorBeforeTheReportFallsBackAndRecordsIt TestQAProviderSystemErrorExhaustionSettlesInsteadOfHalting TestLostRolloutBeforeFirstHandoffActivatesTheFallback TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five new tests do not exist, so their PASS lines are missing and the command fails.
- `out="$(go test -count=1 -v -run '^(TestDeliveryRunResultReadsAnExhaustedProviderSystemError|TestDeliveryRunResultReadsAnExhaustedLostRollout)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryRunResultReadsAnExhaustedProviderSystemError TestDeliveryRunResultReadsAnExhaustedLostRollout; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the new test does not exist, so its PASS line is missing.

## References

- [_prd.md](_prd.md) — Goals 2-3; User Stories 1-2; Core Feature 2; Success Metric 2; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; Invariants 7-11; API Contract 2; API Contract 3; Testing Approach 3; Testing Approach 4; Build Order 3
- [2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md](references/2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md)
- ADR-0274; ADR-0245; ADR-0114; ADR-0050; ADR-0057; ADR-0259; ADR-0271

## Result

Implemented the assigned recovery and delivery-reader slice. The Agent Session
owner checks a Provider System Error after a Lost Rollout and before Selection
Failure handling. Eligible Tasks and pending QA reports move to the next
Fallback Selection; other non-review scopes exhaust immediately. Review scopes
return the original error. Recovery closes the failed session, persists its
failed attempt with `provider_system_error`, records the event and progress
line, and supplies the original Task prompt plus the recovery notice. QA
fallback records its three report lines before the next session opens and
leaves the seeded verdict pending.

Exhaustion carries the specified runtime infrastructure reason through Task
settlement. QA exhaustion settles its Task failed and returns a pending verdict
without halting the cycle. The delivery workflow reads exhausted provider
failure events into `RuntimeInfrastructure` with the specified scope and detail.
The existing delivery engine then uses its runtime-infrastructure park and
uncounted retry path; its engine, park classes and retry accounting were not
changed.

### Focused checks

- Before implementation:
  `rtk proxy go test -count=1 -run 'ProviderSystemError' ./internal/daemon ./internal/cli`
  exited 1. All five new daemon scenarios reproduced missing recovery or
  settlement, and the new delivery-reader scenario returned no provider detail.
- After implementation, the same focused command exited 0 for both packages.
- After adding eligibility and independent-budget coverage:
  `rtk proxy go test -count=1 -run 'ProviderSystemError|LostRollout' ./internal/daemon ./internal/cli`
  exited 0 for both packages. This selects the new tests, existing provider
  review classification tests, and every Lost Rollout test in these packages.
- `rtk proxy git diff --check` exited 0 after this Result was appended.
- `rtk proxy git diff --exit-code -- internal/daemon/lost_rollout_recovery_test.go internal/cli/deliver_runtime_infrastructure_test.go internal/delivery`
  showed no changes to those regression files or delivery engine paths.

### Acceptance evidence

1. **Fallback before First Handoff or while QA is pending, recorded once:**
   `TestProviderSystemErrorBeforeFirstHandoffActivatesTheFallback` proves the
   fallback selection, session closure, failed attempt reason code, single
   provider event and progress line, notification, and recovery prompt.
   `TestQAProviderSystemErrorBeforeTheReportFallsBackAndRecordsIt` proves that
   the fallback runs the gate, the pending seed already carries each required
   report line once, and the final report retains those lines.
2. **Exhaustion without retrying the selection:**
   `TestProviderSystemErrorAfterFirstHandoffSettlesRuntimeInfrastructure` proves
   a failed Verification Feedback turn stays on its original session, settles
   with the exact runtime infrastructure text, and opens no fallback.
   `TestProviderSystemErrorWithoutFallbackSettlesRuntimeInfrastructure` proves
   one turn and one exhausted event.
   `TestQAProviderSystemErrorExhaustionSettlesInsteadOfHalting` proves failed QA
   settlement, a preserved pending report, and a pending cycle verdict without
   an error. `TestProviderSystemErrorRecoveryEligibility` also covers the
   no-output Selection Failure shape, both review error shapes, and refusal of
   QA fallback when the report-pending predicate is false.
   `TestProviderSystemErrorDoesNotUseLostRolloutBudget` proves provider
   recoveries leave Lost Rollout recovery available, and a spent Lost Rollout
   budget does not prevent an eligible provider fallback.
3. **Delivery parks runtime infrastructure with the provider detail:**
   `TestDeliveryRunResultReadsAnExhaustedProviderSystemError` proves the exact
   API Contract 3 text for exhaustion, no infrastructure classification for a
   recovered fallback, and unchanged Lost Rollout text. Source inspection of
   `internal/delivery/engine.go` confirms that a non-empty
   `RuntimeInfrastructure` selects the existing `runtime-infrastructure` park
   and bypasses counted retry enforcement. `internal/store/delivery.go` sets
   the retry increment to zero for that blocker. No queue-engine execution is
   claimed by this focused check.
4. **Lost Rollout regression behavior preserved:** the final focused command
   passes all Lost Rollout tests selected in `internal/daemon` and
   `internal/cli`, including every test in both unchanged regression files.
   The Lost Rollout recovery implementation remains unchanged.

Task status remains Daemon-owned. The declared Verification commands and both
repository Verification commands were not run. No commit, push or pull request
was made. No follow-up outside this Task's slice was implemented.
