---
spec: 0131-a-failed-gate-accepts-its-repair
status: archived
created: 2026-09-10
surfaces: [backend]
archived: "2026-09-10"
source_slug: 0131-a-failed-gate-accepts-its-repair
---


# A failed gate accepts its repair

When a Spec's terminal QA gate fails, the documented recovery is to author
corrective Tasks and let the gate run again. Today that recovery cannot be
executed: adding a Task under a gate that has already settled makes the Task
Graph refuse to load, so the Supervisor can neither check the Spec nor dispatch
the corrective work. The gate that found the defect is the reason the defect
cannot be fixed.

Measured on 2026-09-10 while closing Spec 0119: its QA gate settled `failed`
with four findings, and adding two corrective Tasks as gate dependencies made
`roundfix spec check` and `roundfix implement` fail with a stale-gate load
error naming those Tasks.

## Project Constraints

- Identifier strategy: applicable — preserve existing Spec slugs, Task IDs and error identities; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0091 keeps the authored QA gate terminal, and this Spec preserves that; ADR-0057 keeps the Daemon the exclusive writer of Task status, so recovery must not ask a human or an Agent to edit a status; ADR-0117 places a check at the stage that can establish it, and gate staleness is established only where a passing verdict exists; ADR-0096 requires mechanical facts before the QA Agent turn and is preserved unchanged, because this Spec alters when a graph loads and not which facts the gate proves. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The change is ordinary source in `internal/spec` with its test. Source: `docs/agents/agent-instructions.md`.

## Goals

- A Spec whose QA gate failed accepts corrective Tasks and can be checked and
  dispatched without a human editing any Task status.
- A gate that passed still cannot stand above work it never saw.

## User Stories

1. As the Supervisor, I want to author corrective Tasks after a failed QA gate
   and dispatch them, so that a gate finding leads to a repair instead of a
   deadlock.
2. As a reviewer, I want a passing gate to remain invalid above incomplete
   dependencies, so that a stale pass can never certify unseen work.

## Core Features

1. Gate staleness invalidates only a gate whose recorded verdict is a pass. A
   gate recorded as failed accepts new incomplete dependencies, because it
   carries no passing verdict to protect.
2. The Task Graph loads in that state, so checking, dispatch and settlement all
   proceed and the corrective Tasks run before the gate runs again.
3. A passing gate above an incomplete dependency keeps refusing with the same
   error and the same identity it reports today.
4. No Task status is changed by this behavior; the Daemon remains the only
   writer of status, and recovery needs no manual edit.

## User Experience

The Supervisor adds corrective Tasks to a Spec whose gate failed, runs the
checker, sees a clean load, and dispatches. Nothing instructs anyone to edit a
status field. A stale passing gate still refuses and still names the
dependencies that invalidated it.

## Non-Goals / Out of Scope

- A general gate-recovery contract, carry-forward of prior QA evidence, or
  re-validation policy. Spec 0129 owns those.
- Changing gate terminality, Task status ownership, QA verdict semantics, or
  the archive contract.
- Relaxing the requirement that the gate depends on every non-QA leaf.

## Success Metrics

- A characterization records today's refusal for both a completed and a failed
  gate above an incomplete dependency; after the change the completed case
  still refuses and the failed case loads.
- Adding a corrective Task to a Spec whose gate failed lets the checker load
  and report normally.

## Decisions

- The staleness rule protects a passing verdict. A failed gate has no verdict
  to protect, so extending the rule to it bought no safety and cost the only
  documented recovery path.
- Scope is deliberately one condition and its characterization. The broader
  recovery contract stays with Spec 0129 so this repair can land immediately.

### Declared intentional breaks

1. A Task Graph whose QA gate is recorded failed now loads with incomplete
   dependencies, where it previously refused with a stale-gate error.

Everything else keeps behaving as it does: a completed gate above an incomplete
dependency still refuses with the same error type and message, the gate still
must cover every non-QA leaf, and Task status ownership is unchanged.

### Regression locks

- Today's behavior for both gate verdicts is captured as characterization before
  the change, so the completed case cannot regress silently.
- The new acceptance is fail-closed on evidence: it applies only to a gate whose
  recorded status is exactly failed, never to an absent, malformed or unknown
  status.

## Open Questions

None.
