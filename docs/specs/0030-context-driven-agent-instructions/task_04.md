---
task: task_04
spec: 0030-context-driven-agent-instructions
status: pending
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

- [ ] Implement installed-skill and lockfile discovery.
- [ ] Implement required, extra locked, local, and untracked classification.
- [ ] Add module-reference-to-setup validation and stable skill finding codes.
- [ ] Implement canonical setup parsing, normalization, digesting, synchronization, and check-only behavior.
- [ ] Add tests for missing skills, extras, local exemptions, malformed lockfiles, and setup drift.
- [ ] Verify all three initial setup snapshots round-trip deterministically.

## Acceptance Criteria

- [ ] Missing required skills and referenced skills outside the setup block compliance with precise next actions.
- [ ] Extra installed skills appear only when requested, remain informational, and never change the exit code from `0` by themselves.
- [ ] Repository-owned and untracked skills are not mislabeled as automatic removal candidates.
- [ ] No script path removes a skill or edits `skills-lock.json`.
- [ ] Snapshot synchronization from unchanged canonical setup files produces no diff.
- [ ] Snapshot check reports drift when canonical normalized content differs and succeeds when it matches.
- [ ] Auditing a repository remains fully portable when no canonical setups directory is available.

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
