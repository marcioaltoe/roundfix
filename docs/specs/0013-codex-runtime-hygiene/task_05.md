---
task: task_05
spec: 0013-codex-runtime-hygiene
status: pending
type: docs
complexity: low
---

# Task 05: Docs and skill sync

## Overview

Document the Doctor Command and the codex-hygiene behavior, and sync the Roundfix
Skill for the new command surface. Runs after the behavior tasks so the
documented contract matches shipped behavior and the SKILL.md-matches-CLI gate
closes.

## Requirements

1. MUST document `roundfix doctor` — what it checks, its diagnosis-only nature,
   the codex-hygiene check, and the macOS-only quarantine/notarization behavior
   with the curl-reinstall remediation.
2. MUST note the verified-clean codex spawn behavior where codex/acpx setup is
   documented.
3. MUST update `.agents/skills/roundfix/SKILL.md` (and manifest) for the new
   `doctor` command surface and regenerate the embedded copy with
   `make skills-sync`.
4. MUST leave no skill drift and no doc/behavior mismatch; CONTEXT.md already
   carries the Doctor Command term.

## Subtasks

- [ ] Doctor Command docs incl. codex hygiene and platform behavior
- [ ] Note verified-clean codex spawn where relevant
- [ ] Update SKILL.md/manifest for `doctor` + `make skills-sync`
- [ ] Verify no skill drift

## Acceptance Criteria

- [ ] Docs accurately describe `roundfix doctor`, its checks, and the macOS-only codex hygiene remediation.
- [ ] SKILL.md matches the shipped CLI surface; the embedded copy is regenerated.
- [ ] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories (documentation). `_techspec.md` → Build Order 5. ADR-0032.
CONTEXT.md → Doctor Command. CLAUDE.md SKILL.md-matches-CLI gate.
