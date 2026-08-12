---
task: task_03
spec: 0016-worktree-bootstrap
status: completed
type: docs
complexity: low
---

# Task 03: Docs and skill sync (bootstrap + env-file recipe)

## Overview

Document `worktree.bootstrap` and tie it together with the existing
`worktree.copy` into a coherent recipe for stateful monorepos, and sync the
Roundfix Skill for the new config surface. Runs after the behavior tasks so the
documented contract matches shipped behavior and the SKILL.md-matches-CLI gate
closes.

## Requirements

1. MUST document `worktree.bootstrap` and `worktree.bootstrap_timeout` (purpose,
   when it runs, failure behavior) and present the combined recipe with
   `worktree.copy` (env files) and `worktree.concurrency: 1` for shared-database
   monorepos, including the gitignore-safety note for copied files.
2. MUST update `.agents/skills/roundfix/SKILL.md` (and manifest) for the new
   config surface and regenerate the embedded copy with `make skills-sync`.
3. MUST note that Roundfix runs the bootstrap command but does not own database
   provisioning or dependency strategy (that lives in the command).
4. MUST leave no skill drift and no doc/behavior mismatch; CONTEXT.md already
   carries the Worktree Bootstrap term.

## Subtasks

- [x] Docs for `worktree.bootstrap`/`bootstrap_timeout` + the monorepo recipe
- [x] Document `worktree.copy` for env files with the gitignore-safety note
- [x] Update SKILL.md/manifest + `make skills-sync`
- [x] Verify no skill drift

## Acceptance Criteria

- [x] Docs describe `worktree.bootstrap`, its failure behavior, and the combined `copy`/`bootstrap`/`concurrency: 1` recipe for stateful monorepos.
- [x] SKILL.md matches the shipped CLI/config surface; the embedded copy is regenerated.
- [x] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories, Core Feature 5 (documentation). `_techspec.md` → Env-file
recipe, Build Order 3. ADR-0034. CONTEXT.md → Worktree Bootstrap. CLAUDE.md
SKILL.md-matches-CLI gate.

## Result

- README documents `worktree.bootstrap`, `worktree.bootstrap_timeout`, when the
  command runs, start/non-zero/timeout failure behavior, stderr/journal output,
  and the shared-database monorepo recipe with `worktree.copy` and
  `worktree.concurrency: 1`.
- README and the Roundfix skill document that copied env files must already be
  gitignored, and that Roundfix only invokes and times the bootstrap command;
  dependency installation, database provisioning, migrations, seeding, and cache
  strategy live in the configured command.
- `.agents/skills/roundfix/SKILL.md` and
  `.agents/skills/roundfix/agents/openai.yaml` were updated for the new config
  surface. `rtk make skills-sync` regenerated `skills/roundfix/SKILL.md` and
  `skills/roundfix/agents/openai.yaml`.
- Verification passed:
  `rtk go run ./cmd/roundfix skills check` reported
  `Roundfix skill check passed: roundfix`; `rtk make skills-sync-check` reported
  no drift; `rtk make verify` passed tests, skills check, and build.
