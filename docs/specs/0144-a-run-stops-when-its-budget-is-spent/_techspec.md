---
spec: 0144-a-run-stops-when-its-budget-is-spent
prd: _prd.md
created: 2026-09-18
---

# A Run stops when its budget is spent — Technical Spec

## Executive Summary

The watch loop already implements this rule: derive a deadline from the Run's
start, compare against it, settle `BudgetExceeded`. This design applies the same
rule at the Implement path, ends the Run through the cancellation the Stop
Command already uses, and widens Task Carry-Forward by one accepted outcome.

The trade-off this design accepts is that the bound is wall-clock and blind: it
cannot tell a stalled Agent Session from slow but healthy work, so it ends both.
That is what makes it useful unattended, and why the Run stays recoverable
instead of being discarded.

## Project Constraints

- Identifier strategy: applicable — the Run outcome vocabulary, Run identifiers and Run Branch names keep their spelling; this Spec produces an existing outcome in a new place and coins none. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — bounding a Run reads local configuration and cancels local Agent Sessions and processes; no credential store or network transport is involved. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — Run settlement, Task status and cancellation are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0057 applies: the Daemon owns Run and Task status, so the bound is the Daemon's to apply and record.
  ADR-0127 applies: process residue is a readiness observation, so a bounded Run proves the processes it cancelled rather than assuming them gone.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the watch loop that already enforces this budget and the Run Database vocabulary that already names the outcome.
  ADR-0158 applies: a budget-expired Run settles `BudgetExceeded` and stays recoverable through carry-forward.
- Tooling authority: applicable — the Run outcome and the carry-forward contract are public CLI behavior, and the repository's hard rule ships the Roundfix skill update with such a change. Express maintainer authorization: granted 2026-09-18, recorded in [_authorization.md](_authorization.md); bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync` regenerates the mirror, and `make baseline-digests` rewrites any derived pin the approved edit moves. The Daemon, the CLI and the Run Database are ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Run deadline | `internal/daemon` run lifecycle | Derive the deadline from the Run's start when the budget is enabled, and end the Run when it passes. |
| Cancellation | `internal/daemon` stop path | Cancel owned Agent Sessions and processes, proving termination. |
| Settlement | `internal/daemon` outcome recording | Settle `BudgetExceeded` with a reason naming the maximum and the elapsed time. |
| Carry-forward acceptance | `internal/worktree` reconciliation | Accept a `BudgetExceeded` Run beside `Stopped` and `Unresolved`. |
| Shipped skill | `.agents/skills/roundfix/SKILL.md` and its mirror | State the bound, the outcome and the widened acceptance. |

No new package, state or configuration key is proposed. The outcome, the
configuration and the cancellation path all exist.

## Implementation Design

### Deriving and honouring the deadline

When the budget is enabled and the maximum is positive, the Run's deadline is
its start plus the maximum, computed once. The lifecycle already checks for a
Stop Request between Tasks and inside the Task cycle; the deadline is checked at
the same points, so the bound lands where the Daemon can act rather than in the
middle of an operation it cannot interrupt safely.

A deadline that has passed ends the Run: no further Task starts, the current
Agent Session and owned processes are cancelled, and the Run settles.

### Cancelling what the Daemon owns

Cancellation reuses the Stop Command's path, which cancels registered Agent
Sessions, terminates the recorded owning process and proves absence before
reporting success. An unprovable termination is reported with its reason, as it
is today, and does not become a silent success.

### Settling the outcome

The Run settles `BudgetExceeded`. Its reason names the configured maximum and
the elapsed time, within the existing public reason length. Tasks that completed
keep their status; a Task interrupted mid-flight settles the way an interrupted
Task settles today. The Run Worktree and Run Branch are preserved, as they are
for every non-integrated outcome.

### Accepting the outcome for carry-forward

Carry-forward's accepted-outcome set gains `BudgetExceeded`. Every other rule it
applies is unchanged: the proof requirements, the refusal of a set whose members
cannot all be proved, and the absence of a force bypass.

### Interfaces

No exported signature changes. The lifecycle gains a deadline it already knows
how to compare, and the accepted-outcome set gains one member.

### Data Models

No schema, stored record or event payload shape changes. `BudgetExceeded` is an
existing Run Database state.

## API Contracts

1. An Implement Run with the budget enabled ends at its configured maximum and
   settles `BudgetExceeded`, with a reason naming the maximum and the elapsed
   time.
2. `reconcile --carry-forward` accepts a `BudgetExceeded` Run beside `Stopped`
   and `Unresolved`, and keeps refusing every other terminal outcome.
3. With the budget disabled, every Run outcome, report and exit code keeps its
   current behavior, including the Run Window report at start.

## Coverage Map

- Goal 1 → Run deadline.
- Goal 2 → Settlement.
- Goal 3 → Carry-forward acceptance.
- Goal 4 → Run deadline (disabled path).
- User Story 1 → Run deadline, Cancellation.
- User Story 2 → Settlement; API Contract 1.
- User Story 3 → Carry-forward acceptance; API Contract 2.
- User Story 4 → Run deadline (disabled path); API Contract 3.
- Core Feature 1 → Run deadline.
- Core Feature 2 → Cancellation.
- Core Feature 3 → Settlement.
- Core Feature 4 → Carry-forward acceptance; Shipped skill.
- Core Feature 5 → Run deadline (disabled path).
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 2.
- API Contracts 1-3 → Run deadline, Carry-forward acceptance, Settlement.

## Integration Points

- **Stop Command.** The bound reuses its cancellation and its proof of
  termination; its own outcome and owner-identity proof are untouched.
- **Run Window.** Unchanged. It still refuses a start outside the window and
  still reports that a Run may run past it.
- **Watch loop.** Unchanged. It keeps its own deadline and its own budget
  handling for Rounds.
- **Spec 0127.** Durable orchestration, queue preparation, delivery order and
  progress records stay there.

## Testing Approach

1. **Bound and settlement, at the task-cycle seam.** A fixture Run whose
   configured maximum expires during its work ends, cancels its Agent Session,
   and settles `BudgetExceeded` with the reason's two facts. The case fails on
   the tree as it stands today, where no deadline exists.
2. **Carry-forward acceptance, at the reconciliation seam.** A `BudgetExceeded`
   Run with proved Tasks is accepted; one whose Tasks cannot be proved is
   refused; an outcome outside the set is still refused.
3. **Disabled budget, at the same seams.** No deadline is derived and the Run
   settles exactly as before.
4. **Outside evidence.** The watch loop's deadline derivation and budget check,
   and the Run Database's outcome vocabulary, are read as they stand and shown to
   be the rule this Spec applies rather than a new one. Where either cannot be
   read, the row records that and does not block.
5. **Repository gate.** The terminal QA Task records the Daemon's `make verify`
   result as a fact.

## Build Order

1. Derive the deadline and end the Run through the existing cancellation, with
   task-cycle tests (depends on: none).
2. Settle `BudgetExceeded` with its reason, with tests (depends on: 1).
3. Accept the outcome for carry-forward, with reconciliation tests (depends on:
   2).
4. State the bound, the outcome and the widened acceptance in the Roundfix skill
   and its mirror (depends on: 1, 2, 3).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A blind bound.** Wall-clock cannot tell a stall from slow work, so a healthy
  long Run with a short maximum ends too. The budget is opt-in and its value is
  the maintainer's; the Run stays recoverable.
- **Cancellation that cannot be proved.** The existing path already reports an
  unprovable termination with its reason instead of claiming success, and this
  design inherits that rather than weakening it.
- **A widened carry-forward.** One more accepted outcome is one more way to hand
  work back; the proof requirements that decide whether a Task may be handed back
  are untouched.

## Decisions

- **Check the deadline where the Stop Request is already checked.** Those are the
  points where the Daemon can end a Run without corrupting an operation.
- **Reuse the cancellation path.** A second way to cancel would be a second way
  to be wrong about what stopped.
- **One more accepted outcome, no force bypass.** Carry-forward's proof rules are
  what make it safe, and they do not change.

## Vocabulary Contract

No token is coined. `BudgetExceeded`, Run Worktree, Run Branch and Task
Carry-Forward are existing terms used with their existing meanings.
