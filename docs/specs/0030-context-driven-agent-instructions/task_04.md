---
task: task_04
spec: 0030-context-driven-agent-instructions
status: completed
type: backend
complexity: medium
---

# Task 04: Validate canonical skill setups and installed skills

## Overview

Connect selected instruction profiles to portable canonical skill requirements and the repository's installed skills. The slice blocks missing capabilities and references outside the selected setup while offering a strictly informational, opt-in cleanup report that never removes anything.

## Requirements

1. MUST validate that every skill referenced by a selected module exists in the selected profile's canonical setup snapshot.
2. MUST discover installed skills from `.agents/skills/*/SKILL.md` and enrich their origin from `skills-lock.json` when available.
3. MUST report missing required skills and generated rules that reference skills outside the selected setup as blocking findings.
4. MUST expose locked skills outside the selected setup only when extra-skill reporting is requested and MUST keep them informational.
5. MUST exempt manifest-declared repository-owned skills and classify untracked skill directories separately rather than treating them as removable.
6. MUST NOT delete skills, modify `skills-lock.json`, or emit an executable removal action.
7. MUST implement `sync-setups` against an explicitly supplied canonical setups directory, including check-only behavior, normalized paths, deterministic digests, and atomic snapshot updates.
8. MUST allow normal audit and apply operations to run without access to the canonical setups checkout or network.

## Subtasks

- [x] Implement installed-skill and lockfile discovery.
- [x] Implement required, extra locked, local, and untracked classification.
- [x] Add module-reference-to-setup validation and stable skill finding codes.
- [x] Implement canonical setup parsing, normalization, digesting, synchronization, and check-only behavior.
- [x] Add tests for missing skills, extras, local exemptions, malformed lockfiles, and setup drift.
- [x] Verify all three initial setup snapshots round-trip deterministically.

## Acceptance Criteria

- [x] Missing required skills and referenced skills outside the setup block compliance with precise next actions.
- [x] Extra installed skills appear only when requested, remain informational, and never change the exit code from `0` by themselves.
- [x] Repository-owned and untracked skills are not mislabeled as automatic removal candidates.
- [x] No script path removes a skill or edits `skills-lock.json`.
- [x] Snapshot synchronization from unchanged canonical setup files produces no diff.
- [x] Snapshot check reports drift when canonical normalized content differs and succeeds when it matches.
- [x] Auditing a repository remains fully portable when no canonical setups directory is available.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `skills-lock.json`

## Verification

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_skills*.py'` — expected: required, extra, local, untracked, and lockfile classifications pass.
- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_sync_setups*.py'` — expected: setup normalization, drift detection, atomic update, and idempotency pass.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → User Stories 5–6 and 8; Core Features 9–11 and 14; Non-Goals.
- `_techspec.md` → Data Models: Setup snapshots and installed-skill classification; Integration Points; Build Order 4.

## Result

Implemented skill capability validation and canonical setup synchronization in `context_setup.py`, then mirrored the setup-context-driven skill bundle. Audit now discovers installed `.agents/skills/*/SKILL.md` directories, enriches locked origins from `skills-lock.json`, blocks missing required setup skills, reports module skill references outside the selected snapshot, and exposes extra locked or untracked skills only through `--show-extra-skills` as informational findings. `sync-setups` now reads an explicit canonical setups directory, normalizes setup paths, recomputes deterministic digests, supports `--check`, and writes snapshot updates through temporary sibling files and atomic replacement.

Evidence by acceptance criterion:

- Missing required skills and outside-setup module references block compliance: `test_missing_required_skill_blocks_compliance` and `test_module_skill_reference_outside_setup_is_blocking`.
- Extra installed skills are opt-in informational findings and preserve exit code `0`: `test_extra_locked_skills_are_informational_only_when_requested`.
- Repository-owned local skills and untracked skills are classified separately, not as automatic removal candidates: `test_local_and_untracked_skills_are_not_removal_candidates`.
- Audit does not remove skills or edit `skills-lock.json`, and no finding action emits an executable removal command: `test_extra_locked_skills_are_informational_only_when_requested`, `test_local_and_untracked_skills_are_not_removal_candidates`, and `test_malformed_lockfile_is_invalid_input_without_writes`.
- Unchanged canonical setup snapshots produce no diff: `test_check_succeeds_when_canonical_snapshots_match`.
- Snapshot check reports drift and matching snapshots pass: `test_check_reports_drift_when_canonical_content_differs` and `test_check_succeeds_when_canonical_snapshots_match`.
- Normal audit remains portable without a canonical setups directory: `test_audit_is_portable_without_canonical_setups_directory`.

Verification:

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_skills*.py'` — passed, 6 tests.
- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_sync_setups*.py'` — passed, 4 tests.
- `rtk git diff --check` — passed.
- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — passed, 6 tests.
- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_apply*.py'` — passed, 8 tests.
- `rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py --help` — passed and documents implemented `sync-setups`.
- `rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py sync-setups --help` — passed and documents `--source-dir`, `--check`, and output format.
- `rtk make verify` — passed; Go tests passed in 19 packages, Roundfix skill check passed, and the binary built.
