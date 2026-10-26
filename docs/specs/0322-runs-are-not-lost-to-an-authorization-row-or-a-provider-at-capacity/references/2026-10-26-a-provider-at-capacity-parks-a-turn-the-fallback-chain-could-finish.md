---
type: fix
status: promoted
created: 2026-10-26
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
reason: null
---

# A provider at capacity parks a turn the Fallback Chain could finish

## Symptom

`codex-acp` answers a turn with `threadStatus` `systemError`, the text
"Selected model is at capacity. Please try a different model." and
`stopReason` `end_turn`. Since ADR-0274 Roundfix reads that as a Provider
System Error and lets the next Fallback Selection take the Task, but only
before the Task's First Handoff, or in a QA Task while its seeded report
still says `verdict: pending`. Outside that window the recovery is
`exhausted` even though the profile names a next selection: the Task settles
`failed` as runtime infrastructure, the Run ends Unresolved and the Delivery
Queue parks the item `runtime-infrastructure`. The retry is not counted, but
it waits for an operator.

It happened twice in one night on this repository's queue, both in a QA turn
whose agent had already written `verdict: partial` into the report and was
still checking rows:

- Spec 0300, Run `run_20261023T022528Z_d6dffdb8e615fb09`, 2026-10-23 00:04
  BRT (03:04 UTC). The QA turn ran 13 minutes on `codex/gpt-6.1-sol/high`;
  a script in it replaced `verdict: pending` with `verdict: partial`; then
  the capacity reply arrived. Event `daemon.task` at cursor 2496: "qa qa hit a
  provider system error: Selected model is at capacity. Please try a
  different model.; recovery: exhausted". The `qa` profile's fallback
  `claude/opus/high` never ran. Parked 00:09 BRT.
- Spec 0302, Run `run_20261023T070317Z_62af7ae2274ed43c`, 2026-10-23 05:07
  BRT (08:07 UTC). The QA turn had written `verdict: partial` and
  `status: closed` at cursor 4124 and was in its final review when the
  capacity reply arrived at cursor 4222; recovery `exhausted` at cursor 4225.
  Parked 05:11 BRT. In the same Run, Task 01's first turn hit the same reply
  before its First Handoff and fell back to `claude/opus/high` as designed
  (cursor 566).

Earlier, on 2026-10-13, Spec 0265's QA lost Runs 2 to 5 to the same reply
before ADR-0274 existed; those losses are why ADR-0274 exists.

## Where

`internal/daemon/agent_session_owner.go` `recoverProviderSystemError`: the
fallback condition is `(scope task && !handedOff) || (scope qa &&
reportPending())`, where `reportPending` reads the seeded report and is true
only for `verdict: pending`. A verdict the failed turn wrote is read as a
finished gate, although the turn did not finish. A Task's repair turn after
its First Handoff has the same gap. The Delivery Queue's
`runtime-infrastructure` park (`internal/cli/deliver_workflow.go` `runResult`,
phase `provider_system_error`, recovery `exhausted`) is correct for an
exhausted chain and is not the defect.

## Expected

A Provider System Error during any Task or QA turn moves the work to the next
selection of the Fallback Chain, with the existing event naming the cause,
and the item parks `runtime-infrastructure` only when the chain has no next
selection. A QA verdict written by the failed turn is not a conclusion: the
fallback session starts from a report reopened to `verdict: pending`, with
the discarded verdict recorded. Light Tier rules and review Batches keep
their behavior.

## Evidence

Run Event Journal rows of the two Runs above, read on 2026-10-26 with
`sqlite3 -readonly` and `roundfix events <run-id> --filter
agent-selection,outcome,task-status`. Queue interventions of 2026-10-23 00:10
and 05:12 BRT in the operator log.
