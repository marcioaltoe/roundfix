---
task: task_04
spec: 0014-run-store-retention
status: pending
type: docs
complexity: low
---

# Task 04: Docs and skill sync

## Overview

Document Run store retention and the GC Command, and sync the Roundfix Skill for
the new command surface. Runs after the behavior tasks so the documented
contract matches shipped behavior and the SKILL.md-matches-CLI gate closes.

## Requirements

1. MUST document the `store.journal_retention` config key (default, `0`
   semantics), the `roundfix gc` command with `--dry-run`, and that retention
   never touches Active Runs, `runs` rows, or locks.
2. MUST note the self-healing preflight-sweep prune and that review artifacts
   under the spec tree are out of retention scope.
3. MUST update `.agents/skills/roundfix/SKILL.md` (and manifest) for the new
   `gc` command surface and regenerate the embedded copy with `make
   skills-sync`.
4. MUST leave no skill drift and no doc/behavior mismatch; CONTEXT.md already
   carries GC Command and Journal Retention.

## Subtasks

- [ ] Docs for `store.journal_retention` + `roundfix gc [--dry-run]` + safety scope
- [ ] Note self-healing sweep prune and out-of-scope review artifacts
- [ ] Update SKILL.md/manifest for `gc` + `make skills-sync`
- [ ] Verify no skill drift

## Acceptance Criteria

- [ ] Docs accurately describe retention, `roundfix gc`/`--dry-run`, and the never-prune-Active-Runs guarantee.
- [ ] SKILL.md matches the shipped CLI surface; the embedded copy is regenerated.
- [ ] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories (documentation). `_techspec.md` → Build Order 4. ADR-0033.
CONTEXT.md → GC Command, Journal Retention. CLAUDE.md SKILL.md-matches-CLI gate.
