---
task: task_05
spec: 0055-owner-identity-without-fork
status: completed
type: docs
complexity: low
---

# Task 05: Document the diagnostics, the marker, and the exit

## Overview

The escape hatch is only usable if its two failure conditions and its one
supervised exit are documented where an operator looks. Update the authorized
roundfix Skill pair and the user guide, then let the sanctioned command own the
derived digest fallout.

## Requirements

1. MUST document both Force Stop refusal conditions in the roundfix Skill pair:
   proven mismatch, and unreadable identity with its next action.
2. MUST document the supervised `--owner-identity-unreadable` flag as an operator
   action of last resort, stating that it never applies to a proven mismatch.
3. MUST document the startup warning and the Run marker for a Run running with
   PID-only reuse protection.
4. MUST state that owner identity is read from the kernel with no subprocess, so
   the proof works on a host that cannot fork.
5. MUST confine protected-tooling edits to exactly the paths the PRD's Tooling
   authority row authorizes, and obtain every derived pin from
   `make baseline-digests`.
6. MUST update any `CONTEXT.md` glossary entry whose shipped behavior these
   Tasks changed, never ahead of the behavior landing.

## Subtasks

- [ ] Update the roundfix Skill pair with the diagnostics, the flag, and the
      marker.
- [ ] Update the user guide's Force Stop section.
- [ ] Run `make skills-sync` and `make baseline-digests`, committing only
      authorized paths.

## Acceptance Criteria

- [ ] The Skill states both refusal conditions and their next actions.
- [ ] The Skill states the flag's single valid precondition.
- [ ] The Skill and guide state that identity capture spawns no subprocess.
- [ ] `roundfix skills check` and `make skills-sync-check` pass.
- [ ] Every changed derived pin came from `make baseline-digests`.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `CONTEXT.md`

## Verification

- `make skills-sync-check` — expected: no drift.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: pass.
- `make baseline-digests` — expected: `ok: true`, only authorized paths changed.
- `make verify` — expected: exit 0.

## References

`_prd.md` → Goal 2 Story 4, User Experience, Tooling authority;
`_techspec.md` → Build Order 5.

## Result

### Implementation

- The canonical Roundfix Skill now distinguishes a proven owner-identity
  mismatch from an unreadable identity, gives the operator's next action for
  each refusal, and documents the single condition that permits the
  `--owner-identity-unreadable` last-resort path. It also states that the flag
  never overrides a proven mismatch.
- The Skill and both user-guide Stop sections now state that owner identity
  comes directly from the kernel without a subprocess, so the proof remains
  available when the host cannot fork. They document the one PID-only startup
  warning and the durable `owner_identity_unproven=true` marker rendered by
  `roundfix runs list`.
- `make skills-sync` copied the canonical Skill to the embedded
  `skills/roundfix/SKILL.md`. The existing `Force Stop` glossary entry in
  `CONTEXT.md` now includes the shipped refusal and supervised-exit behavior.
- No derived baseline pin was edited by hand. The declared
  `make baseline-digests` command remains for Daemon Verification, as required
  by the Daemon-assigned execution contract.

### Focused checks

- Red starting point: targeted searches of the pre-change Skill and user-guide
  Stop sections found only the generic Force Stop failure guidance; they did
  not document `owner_identity_unproven=true`, the supervised flag, or the
  no-subprocess identity proof.
- `rtk make skills-sync` exited `0`.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md
  skills/roundfix/SKILL.md` exited `0`, proving the Skill pair is
  byte-identical after synchronization.
- Focused `rtk grep -n` checks exited `0` and found both refusal classes and
  their actions, the unreadable-only flag precondition, the no-subprocess
  statement, the startup warning, and `owner_identity_unproven=true` across
  the Skill pair and user guides.
- `rtk git diff --check` exited `0`.
- The post-edit changed-path inspection found only the canonical and embedded
  Roundfix Skills, `docs/user-guide/commands.md`, `docs/user-guide/usage.md`,
  `CONTEXT.md`, and this Task file. No derived pin path changed during the
  Agent turn.

### Acceptance evidence

- Refusal conditions and next actions: the Skill directs a proven-mismatch
  operator to investigate PID reuse without signaling that process; an
  unreadable kernel-read diagnostic directs the operator to resolve the host
  resource failure and retry.
- Flag precondition: the Skill and guides permit
  `--owner-identity-unreadable` only after normal Force Stop specifically
  reports an unreadable identity. They state that readable identity or proven
  mismatch exits `2` without signaling and that no configuration, environment
  variable, default, or timeout activates the path.
- No subprocess: the byte-identical Skill pair and both guide pages state that
  identity comes from a direct kernel read and spawns no subprocess, including
  on a host that cannot fork.
- Skill checks: `roundfix skills check` and `make skills-sync-check` were not
  run because both are declared `## Verification` commands owned by the
  Daemon. The focused byte comparison establishes synchronization only, not
  the broader checks.
- Derived pins: none changed by the Agent. The Daemon-owned
  `make baseline-digests` Verification command must generate any deterministic
  fallout and establish this criterion.

The commands under `## Verification` were not run; the Daemon owns those
commands, derived digest generation, and the Task verdict.
