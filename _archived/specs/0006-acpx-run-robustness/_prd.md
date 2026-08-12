---
spec: 0006-acpx-run-robustness
status: active
created: 2026-07-05
surfaces: [cli, infra]
---

# ACPX Run Robustness

Two dogfood Runs in one day lost a finished Task to the same transport
failure: the Agent completed and verified its work, acpx's 10 MiB message
buffer blew on a large adapter message at end of turn, the process exited 1,
and Roundfix — classifying Batches by exit code alone — settled the Task
failed and ended the Run Unresolved with correct work stranded uncommitted.
Recovery meant hand-playing the Daemon: re-verify, rewrite the status, craft
the commit. This Spec makes the agent layer classify honestly (the parsed
result outranks the exit code), gives failed-but-done Tasks a one-command
recovery, and applies whatever mitigation exists for the buffer limit itself.

## Goals

- A Batch whose prompt result already arrived never fails on transport
  teardown noise; the Daemon's verbatim verification stays the only gate.
  See ADR-0020.
- A failed Task whose work is preserved in the working tree is recoverable
  with one command — verification first, settlement and commit only on pass.
- The acpx message-buffer limit is mitigated where acpx allows it, and the
  upstream report is ready where it does not.
- Every anomaly stays loud: teardown noise is journaled with the tool's
  stderr tail, never silently swallowed.

## User Stories

1. As a developer whose runtime transport dies after the Agent finished, I
   want the Batch to proceed to verification instead of failing, so that
   finished work is not discarded over teardown noise.
2. As a developer with a failed-but-done Task preserved in the tree, I want
   one command that re-runs its Verification and settles it on pass, so that
   recovery does not require hand-playing the Daemon.
3. As a developer hitting the acpx message-buffer limit, I want Roundfix to
   apply any available mitigation and hand me upstream-ready evidence, so
   that the limit stops killing Runs.
4. As a developer auditing a Run after teardown noise, I want the anomaly in
   the Run Event Journal with the stderr tail, so that nothing about the
   transport failure is invisible.

## Core Features

1. **Result-over-exit classification.** When the stream delivered a parsed
   prompt result before a nonzero acpx exit, the Batch proceeds — the exit is
   journaled as an anomaly (existing daemon event kinds, stderr tail
   included) and verification decides the outcome as always. Without a
   parsed result, behavior is unchanged. See ADR-0020.
2. **Settle Command.** `roundfix settle` targets one failed Task of a Spec in
   the current repository: Preflight Validation (Spec and Task valid, task
   status `failed`, no Active Run for the target or working tree), then the
   Task's Verification commands verbatim; on pass it settles `completed` and
   creates the standard Task commit; on failure it changes nothing and exits
   1. The commit stages the working tree's current changes plus the task
   file — the command's documented contract is that the developer reviews
   the tree before settling (the multi-writer stance applies).
3. **Buffer mitigation.** Apply acpx's message-buffer configuration when one
   exists (config or invocation surface, verified against the pinned
   version); otherwise document the limit in the shipped docs and produce
   the upstream issue evidence. Either way the findings log records the
   outcome.
4. **No silent anomalies.** The teardown-noise path and the settle path both
   report deterministically: journal events for the Run case, stdout report
   lines for the command case.

## User Experience

One new command (`roundfix settle --spec <slug> --task <task_id>`) with the
house exit codes (0 settled, 1 verification failed, 2 Preflight Validation)
and a deterministic stdout report: one line per Verification command with its
outcome, then one settle line. Runs behave identically except that transport
teardown noise no longer fails finished Batches.

## Non-Goals / Out of Scope

- Retry budgets or Agent escalation (work-plan item 7).
- Any change to verification semantics, commit contracts, or ADR-0014
  ownership — settlement still happens only after passing verification.
- Automatic settle: recovery stays an explicit developer command.
- Fixing acpx upstream ourselves; we mitigate, evidence, and report.

## Success Metrics

- A simulated post-result nonzero exit (fake acpx rig) proceeds to
  verification and commits on pass, with the anomaly journaled.
- A real failed-but-done Task settles with one command producing a commit
  byte-shaped like the Daemon's (message, trailers, staged content).
- The buffer investigation ends with either a working mitigation proven
  against acpx 0.12.0 or a filed-ready upstream report recorded in the
  findings log.
- The full existing suite passes unchanged apart from deliberate
  classification-test updates.

## Decisions

- The parsed prompt result outranks the exit code; verification remains the
  gate. See ADR-0020.
- Settle stages the current working tree plus the task file — reviewing the
  tree first is the developer's documented responsibility (consistent with
  the multi-writer worktree stance).
- Settle is a local recovery command like stop: it creates no Run and writes
  no journal events.
- The settle target must be `failed` — settling `pending` or `in_progress`
  Tasks is refused (they belong to a Run or a re-run).

## Open Questions

None.
