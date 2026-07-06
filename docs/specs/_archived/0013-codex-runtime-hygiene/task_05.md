---
task: task_05
spec: 0013-codex-runtime-hygiene
status: completed
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

- [x] Doctor Command docs incl. codex hygiene and platform behavior
- [x] Note verified-clean codex spawn where relevant
- [x] Update SKILL.md/manifest for `doctor` + `make skills-sync`
- [x] Verify no skill drift

## Acceptance Criteria

- [x] Docs accurately describe `roundfix doctor`, its checks, and the macOS-only codex hygiene remediation.
- [x] SKILL.md matches the shipped CLI surface; the embedded copy is regenerated.
- [x] `roundfix skills check` passes and `skills-sync-check` reports no drift.

## Verification

- `rtk make skills-sync` then `rtk go run ./cmd/roundfix skills check` — expected: passes, no drift.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all pass.

## References

`_prd.md` → all stories (documentation). `_techspec.md` → Build Order 5. ADR-0032.
CONTEXT.md → Doctor Command. CLAUDE.md SKILL.md-matches-CLI gate.

## Result

- Documented `roundfix doctor` in `README.md` as a diagnosis-only command that
  checks Node.js, pinned acpx, the configured Agent probe, and codex runtime
  hygiene; the docs now describe the macOS quarantine/Gatekeeper behavior, the
  curl reinstall remediation into `~/.local/bin`, and the non-Darwin skipped
  codex check.
- Added the verified-clean codex spawn note where acpx/codex setup is
  documented: macOS codex-acp launches resolve `CODEX_PATH` first, then
  `PATH`, pass a verified-clean codex through `CODEX_PATH`, and surface the
  hygiene risk when no clean codex is available.
- Updated `.agents/skills/roundfix/SKILL.md` plus
  `.agents/skills/roundfix/agents/openai.yaml` for the Doctor Command surface
  and ran `rtk make skills-sync` to regenerate `skills/roundfix/`.
- Verification passed:
  - `rtk make skills-sync` exited 0 and regenerated `skills/roundfix/`.
  - `rtk go run ./cmd/roundfix skills check` exited 0 with
    `Roundfix skill check passed: roundfix`.
  - `rtk make skills-sync-check` exited 0 with no diff output.
  - `rtk make verify` exited 0; it ran `rtk go test ./...` with 785 passing
    tests in 18 packages, `roundfix skills check`, and the Roundfix build.
