---
task: task_05
spec: 0012-npm-distribution
status: completed
type: docs
complexity: low
---

# Task 05: Docs (install, release runbook, skill bundle) and skill sync

## Overview

Document the two audiences the distribution serves: users who install via
`npx`/`bunx`/global npm, and maintainers who cut a release by pushing a `v*`
tag. Document the shipped Roundfix skill bundle and the `skills` command surface
(`check`, `install`, `list`), sync the skills for any changed guidance, and
close the SKILL.md-matches-CLI gate.

## Requirements

1. MUST add user install docs covering `npx roundfix`, `bunx roundfix`, and
   global install, noting that behavior and exit codes are identical to the
   built binary.
2. MUST add a maintainer release runbook: the version-agreement rule, required
   tokens/secrets, the tag-push flow, and how assets feed the Upgrade Command.
3. MUST document the shipped Roundfix skill bundle (the 14 owned skills) and the
   `skills` command surface — `check`, `install --target/--dir`, and `list` —
   including that external `skills-lock.json` skills are recommendations, not
   shipped.
4. MUST update `.agents/skills/roundfix/SKILL.md` if the `skills` command
   guidance it carries changed, and regenerate the embedded bundle with `make
   skills-sync`.
5. MUST leave no skill drift and no doc/behavior mismatch.

## Subtasks

- [x] User install docs (npx/bunx/global; exit-code parity)
- [x] Maintainer release runbook (version agreement, secrets, tag flow, upgrade channel)
- [x] Skill-bundle docs: the 14 owned skills and the `skills check/install/list` surface
- [x] SKILL.md/skill manifest update if needed + `make skills-sync`
- [x] Verify no skill drift

## Acceptance Criteria

- [x] Install docs accurately describe the three npm install paths and exit-code parity.
- [x] The release runbook documents version agreement, required secrets, and the tag-triggered flow.
- [x] The docs describe the 14 owned skills and the `skills check/install/list` surface, and that external skills are recommendations only.
- [x] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories, Core Feature 6 (documentation). `_techspec.md` → Build
Order 6, Roundfix skill bundle. ADR-0031. Builds on the skill bundle (task_06).
CLAUDE.md SKILL.md-matches-CLI and skill-ownership gates.
