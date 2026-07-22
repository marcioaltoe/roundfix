---
spec: 0037-terminal-outcome-integrity
status: active
created: 2026-07-17
surfaces: [backend, cli, data, docs]
---

# Terminal outcome integrity

A force-stopped Run can currently be completed again by its still-running owner, a Stop Request can remain unnoticed throughout a Review Source wait, and cleanup can target an Agent Session that never reached the active lifecycle. The resulting state is unsafe for users and Supervisors: the Run Database can contradict the Stop Command, the released lock can coexist with live work, and secondary cleanup noise can obscure the primary failure. Prior dogfood evidence was absorbed into this Spec and remains in Git history; the still-open behavior is reproduced by the [Vortex detached-watch finding](../../findings/2026-07-16-vortex-pr87-detached-watch-notification.md).

## Goals

- A terminal Run completion cannot be replaced by a competing terminal outcome.
- Force Stop reports completion and releases the Active Run lock only after the owning process has exited.
- A Stop Request interrupts Review Source waits without waiting for the full review or check deadline.
- Agent Session cleanup acts only on durably registered active sessions and never hides the primary failure.

## User Stories

1. As a user force-stopping a runaway Run, I want the command to prove the owner exited before reporting Stopped, so that no process continues mutating state or files behind a released lock.
2. As a user stopping a watch Run that is waiting on the Review Source, I want the Run to observe the Stop Request by its next poll, so that it does not remain Active until a long grace period expires.
3. As a Supervisor following a terminal Run, I want one stable outcome and one matching outcome event, so that later owner activity cannot contradict the result I already received.
4. As a user diagnosing a failure before Agent work began, I want the primary failure reported first without a nonexistent-session warning, so that recovery starts from the real cause.
5. As a developer integrating an Integration Pending Run, I want the existing recovery flow preserved through an explicit guarded transition, so that general terminal immutability does not strand completed work.

## Core Features

1. Run completion is compare-and-set. A completion may move a non-terminal Run to one terminal outcome; repeating the same outcome is idempotent, while attempting a different terminal outcome is rejected without changing the Run, its completion timestamp, its Active Run lock, or its terminal outcome event.
2. The only permitted transition from one terminal outcome to another is the explicit reconciliation of Integration Pending to Clean after integration evidence has been established. The transition records the prior outcome and the reconciliation evidence.
3. Force Stop first cancels registered Agent Sessions, then terminates the owning process and waits for positive exit confirmation within a bounded interval. If exit cannot be proven, the command fails actionably, leaves the Run Active, and keeps its lock.
4. Review Source status, retry, quiet-period, and merge-readiness waits check the Stop Request before each poll and after each sleep. A waiting Run reaches Stopped no later than its next configured poll boundary.
5. The durable Agent Selection lifecycle is the Agent Session registry. Cleanup targets only scopes whose latest lifecycle state is active; closing a session records the closed state, and an already-absent registered session is idempotent.
6. Primary Run failures print and journal before secondary cleanup warnings. Cleanup failures remain visible but never replace the primary reason.
7. The terminal completion winner alone publishes the terminal outcome and notification. A losing owner observes the stored terminal outcome and exits without another completion event or notification.

## User Experience

- Successful Force Stop reports Stopped only after owner exit is confirmed and the Active Run lock is released.
- Failed Force Stop names the Run, owner process, failed termination step, and the exact command to inspect or retry. It does not claim the Run stopped.
- Graceful stop during a Review Source wait produces the normal Stopped report at the next polling boundary.
- A Run that fails before Agent work contains no session-close warning. If cleanup of a registered session also fails, that warning follows the primary failure and is labeled secondary.

## Non-Goals / Out of Scope

- Pause, resume, or checkpoint recovery for Active Runs.
- Replacing the Run Database as the Stop Request control channel.
- Force-stopping arbitrary processes that are not the recorded Run owner.
- Changing the rule that an in-flight Work Item settles before a graceful Stop Request completes.
- Defining terminal Run Worktree classification or cleanup; spec 0038 owns that behavior.
- Changing Review Source evidence, retry, or notification content; spec 0039 owns those contracts.

## Success Metrics

- A deterministic owner-versus-force-stop race produces zero terminal outcome overwrites and exactly one terminal outcome event.
- A failed owner-termination proof leaves the Run Active and its Active Run lock present in every test case.
- A Stop Request during every Review Source wait phase reaches Stopped by the next configured poll boundary.
- A Run that starts no Agent work produces zero Agent Session close attempts and zero absent-session warnings.
- Integration Pending recovery remains supported only through the guarded reconciliation transition.

## Decisions

- Force Stop fails closed when owner exit cannot be proven; it never releases the lock on a warning-only basis.
- Agent Selection lifecycle records are the Agent Session registry; no second session table is introduced.
- General terminal completion is immutable, with one explicit Integration Pending reconciliation transition. See [ADR-0052](../../adr/0052-run-completion-is-compare-and-set.md).

## Open Questions

None.
