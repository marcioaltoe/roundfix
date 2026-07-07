---
task: task_02
spec: 0022-cleanup-robustness
status: completed
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

- [x] roundfix SKILL.md: Settle recovery + cleanup behavior
- [x] write-tasks SKILL.md: hermetic-verification and no-commit-criteria
      rules
- [x] README Command Boundaries cleanup note
- [x] `make skills-sync`; drift and skills checks pass

## Acceptance Criteria

- [x] The roundfix skill documents the Verification-edit recovery and the
      kept-worktree warning shape truthfully.
- [x] The write-tasks skill carries both authoring rules.
- [x] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Core Feature 3. `_techspec.md` → Build Order 2; Decisions.
CLAUDE.md SKILL.md-matches-CLI HARD RULE.

## Result

Updated the canonical Roundfix skill, write-tasks skill, and README command
boundary docs, then regenerated the embedded `skills/` bundle from
`.agents/skills/`.

Evidence:

- `.agents/skills/roundfix/SKILL.md` now documents the Clean cleanup behavior:
  `git worktree remove --force`, unchanged Clean outcome/stdout/exit code
  after a post-integration cleanup failure, the warning shape
  `roundfix: Run Worktree cleanup failed; kept <path>: <reason>`, and one
  Daemon Run Event. Its Settle section now says to fix a failed task file's
  `## Verification` and re-run Settle when Verification is unsatisfiable, and
  explicitly rejects a skip-verification flag because Verification is the only
  gate.
- `.agents/skills/write-tasks/SKILL.md` now requires hermetic, satisfiable
  task Verification commands for fresh worktrees and forbids commit, push, PR,
  or branch-publishing criteria in task Requirements, Subtasks, Acceptance
  Criteria, or Verification commands.
- `README.md` Command Boundaries now documents forced Roundfix-owned worktree
  cleanup and the post-integration kept-path warning shape.
- `rtk make skills-sync`: passed and regenerated `skills/roundfix/` and
  `skills/write-tasks/`.
- `rtk make skills-sync-check`: passed with no drift output.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check`: passed for
  roundfix and the bundled authorial workflow skills.
- `rtk make verify`: passed. It ran `rtk go test ./...` with 887 tests in
  19 packages, `rtk go run -buildvcs=false ./cmd/roundfix skills check`, and
  `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`; the full
  target exited 0.
