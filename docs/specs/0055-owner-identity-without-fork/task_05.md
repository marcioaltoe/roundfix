---
task: task_05
spec: 0055-owner-identity-without-fork
status: pending
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
- interface: `docs/user-guide/`
- interface: `CONTEXT.md`

## Verification

- `make skills-sync-check` — expected: no drift.
- `roundfix skills check` — expected: pass.
- `make baseline-digests` — expected: `ok: true`, only authorized paths changed.
- `make verify` — expected: exit 0.

## References

`_prd.md` → Goal 2 Story 4, User Experience, Tooling authority;
`_techspec.md` → Build Order 5.
