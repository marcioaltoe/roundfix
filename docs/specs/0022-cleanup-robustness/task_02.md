---
task: task_02
spec: 0022-cleanup-robustness
status: pending
type: docs
complexity: low
---

# Task 02: Docs and skill rules: settle recovery and authoring gates

## Overview

Document the sanctioned recovery for an unsatisfiable task Verification and
record the two field-report authoring rules in the owned skills, closing the
SKILL-matches-CLI gate for the cleanup behavior change.

## Requirements

1. MUST document in the roundfix skill's Settle section: when a task's
   Verification is unsatisfiable (non-hermetic gate), the recovery is to fix
   the task file's `## Verification` and re-run Settle, which re-reads the
   task file; a skip-verification flag is rejected because verification is
   the only gate.
2. MUST add to the write-tasks skill's rules: task Verification commands
   must be hermetic and satisfiable in a fresh worktree, and commit/push
   never appear in requirements or acceptance criteria (the Daemon owns
   them).
3. MUST document the new cleanup behavior (forced removal;
   warn-and-continue after integration with the kept-path warning shape) in
   the README Command Boundaries and the roundfix skill.
4. MUST re-sync the embedded bundle so the drift check passes.

## Subtasks

- [ ] roundfix SKILL.md: Settle recovery + cleanup behavior
- [ ] write-tasks SKILL.md: hermetic-verification and no-commit-criteria
      rules
- [ ] README Command Boundaries cleanup note
- [ ] `make skills-sync`; drift and skills checks pass

## Acceptance Criteria

- [ ] The roundfix skill documents the Verification-edit recovery and the
      kept-worktree warning shape truthfully.
- [ ] The write-tasks skill carries both authoring rules.
- [ ] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Core Feature 3. `_techspec.md` → Build Order 2; Decisions.
CLAUDE.md SKILL.md-matches-CLI HARD RULE.
