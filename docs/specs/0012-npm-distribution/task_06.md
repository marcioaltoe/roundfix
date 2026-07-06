---
task: task_06
spec: 0012-npm-distribution
status: pending
type: backend
complexity: high
---

# Task 06: Roundfix skill bundle — embed, sync, check, install, list

## Overview

Widen the embedded skill set from the single operational `roundfix` skill to the
whole Roundfix-owned bundle — `roundfix` plus the 13 authorial workflow skills —
so the distributed binary carries the spec-driven method it drives. `make
skills-sync`/`skills-sync-check`, `skills check`, and `skills install` operate
over the full set, and a new `skills list` distinguishes the owned bundle from
the recommended external skills.

## Requirements

1. MUST define an owned-skills manifest of exactly the 14 Roundfix-owned skills
   (`roundfix` plus `write-idea`, `write-prd`, `write-techspec`, `write-tasks`,
   `setup-workflow`, `implement-task`, `implement-spec`, `brainstorming`,
   `council`, `business-analyst`, `archive-spec`, `qa-gate`, `evidence-gate`),
   and embed all of them from a synced bundle directory isolated from the
   package's Go files.
2. MUST update `make skills-sync` to sync every owned skill from the canonical
   `.agents/skills/<name>/` and `skills-sync-check` to fail `make verify` on any
   drift across the whole bundle.
3. MUST apply a per-skill `Check` policy: `roundfix` keeps its strict
   contract-wording, Roundfix-branding, and openai-manifest checks; the
   authorial skills get structural validation only (SKILL.md present, frontmatter
   parses, `name` matches the directory, no banned "reference project"
   branding), with no Roundfix-branding and no version-tracks-tag requirement.
4. MUST have `skills install` write all owned skills to the chosen target with
   the existing `--target`/`--dir` semantics unchanged, and MUST add a read-only
   `skills list` that prints the owned bundle and, separately, the recommended
   external skills (names derived from `skills-lock.json` at sync time) with a
   note that external skills install through the user's own skills tooling.
5. MUST NOT embed, vendor, or modify any externally-managed (`skills-lock.json`)
   skill — external skills are listed by name only.

## Subtasks

- [ ] Owned-skills manifest + widened embed from a synced bundle directory
- [ ] `make skills-sync`/`skills-sync-check` over the whole bundle
- [ ] Per-skill `Check` policy (strict roundfix, structural authorial)
- [ ] `skills install` over the full set (semantics unchanged)
- [ ] `skills list` (owned vs recommended external) + recommended-external manifest
- [ ] Tests: install writes 14 skills; check passes for the set and fails on a broken authorial skill; list separates owned from external

## Acceptance Criteria

- [ ] `roundfix skills install --target <t> --dir <tmp>` writes all 14 owned skills to the target.
- [ ] `roundfix skills check` passes for the whole owned set and fails when an authorial skill is missing its SKILL.md or carries banned reference branding.
- [ ] `roundfix skills list` prints the 14 owned skills and the recommended external skills separately, noting external ones install via the user's tooling.
- [ ] `make skills-sync` then `skills-sync-check` reports no drift across the bundle; editing one authorial skill without re-sync fails `make verify`.
- [ ] No `skills-lock.json` skill is embedded or modified.

## Verification

- `rtk go test ./skills/` — expected: the widened skills package tests pass.
- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: passes for the full bundle, no drift.
- `rtk go run ./cmd/roundfix skills list` — expected: owned bundle and recommended external skills printed separately.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → User Stories 7-8; Core Feature 6; Decisions. `_techspec.md` →
Roundfix skill bundle, Build Order 4. CLAUDE.md skill-ownership HARD RULE and
SKILL.md-matches-CLI gate. CONTEXT.md → Roundfix Skill.
