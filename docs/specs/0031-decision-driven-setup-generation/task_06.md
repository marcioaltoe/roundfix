---
task: task_06
spec: 0031-decision-driven-setup-generation
status: completed
type: backend
complexity: high
---

# Task 06: Migrate existing manifests and synchronize the workflow

## Overview

Complete the existing-repository journey and ship the coherent portable skill.
A repository configured by spec 0030 will reuse its compatible answers, replace
only obsolete setup-owned inventory and blocks, reach a clean audit, and rerun
idempotently while the skill asks only newly activated dependent questions.

## Requirements

1. MUST keep Setup Manifest schema version 1 and preserve every compatible
   stored `{value, confirmedAt}` decision without reconfirmation.
2. MUST migrate old mixed managed IDs, template versions, rendered digests, and
   module inventory to the concrete Decision Plan.
3. MUST remove an obsolete block or guide only when the manifest and ownership
   markers prove setup ownership, preserving all repository-authored bytes.
4. MUST stop with a blocking finding instead of guessing when legacy ownership
   or marker state is ambiguous.
5. MUST update the skill workflow to present selection and definite or
   conditional preview, then ask one unresolved decision at a time, including
   dependent questions introduced by enabled capabilities.
6. MUST synchronize the canonical and embedded skill copies and retain
   portable execution without a canonical setups checkout.
7. MUST cover all three profiles, all nine decision effects, TypeScript/Bun
   priority variations, required/extra skills, and repeat apply through real
   temporary-repository subprocesses.

## Subtasks

- [x] Implement schema-v1 inventory and managed-artifact migration from spec
      0030 output.
- [x] Preserve compatible answers and route only missing dependent decisions.
- [x] Add ambiguous-ownership failure and owner-byte preservation coverage.
- [x] Update the setup skill's preview, question, confirmation, and apply flow.
- [x] Add cross-profile and every-decision macro regression cases.
- [x] Regenerate the embedded skill and keep snapshot checks portable.
- [x] Ensure the repository gate includes the focused migration coverage.

## Acceptance Criteria

- [x] A spec 0030 manifest with all compatible answers migrates without a
      decision finding, preserves confirmation dates, and reaches a clean audit.
- [x] Enabling a capability whose dependent value is absent requests only that
      missing decision and persists it for the next run.
- [x] Legacy mixed blocks are removed or replaced only when ownership is proven;
      ambiguous markers block migration without partial writes.
- [x] Repository-authored bytes before, between, and after migrated managed
      blocks remain byte-for-byte unchanged.
- [x] TypeScript/Bun, Go CLI/TUI, and Rust CLI repositories each apply, audit,
      and reapply successfully with representative decision combinations.
- [x] Missing required skills remain blocking, extra skills remain opt-in
      information, and no removal command or file deletion is introduced.
- [x] Canonical and embedded skill copies match, the normal workflow needs no
      developer-specific checkout, and the full repository gate passes.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/setup-context-driven/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`
- interface: `.agents/skills/setup-context-driven/tests/test_workflow.py`
- interface: `docs/specs/0030-context-driven-agent-instructions/qa/manual_variation_qa.py`
- interface: `Makefile`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_manifest_migration*.py'` — expected: compatible-answer reuse, inventory migration, ambiguous ownership, no-partial-write behavior, owner-byte preservation, and idempotency pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_macro_profiles*.py'` — expected: all profiles and representative decision combinations apply, audit, and reapply successfully with correct skill semantics.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_workflow*.py'` — expected: conditional preview, one-question routing, confirmation, and portable orchestration pass.
- `rtk make skills-sync-check` — expected: canonical and embedded workflow skill copies match.
- `rtk make verify` — expected: formatting, Go tests, all setup-context-driven tests, asset validation, skill checks, sync check, and build pass.

## References

- `_prd.md` → Goals 4–5; Core Feature 7; Success Criteria for spec 0030 migration, idempotency, and owner-byte preservation; Non-goals.
- `_techspec.md` → Data models: Setup Manifest compatibility; Integration points; Testing approach; Build Order 5–6; Risks and considerations.
- ADR-0046 and ADR-0047.

## Result

- Implemented schema-version 1 migration behavior that preserves every
  compatible catalog decision already stored in the manifest, including
  inactive dependent answers and their `confirmedAt` dates.
- Added stale managed-artifact ownership validation so obsolete inventory is
  removed only when the target file still contains the matching setup marker;
  ambiguous legacy ownership blocks before writes.
- Adjusted guide adoption so split setup-owned guides can be appended to files
  that already contain setup-managed blocks without adopting repository-authored
  bytes.
- Updated the setup workflow skill to report active or conditional modules,
  conditional planned-change triggers, and one unresolved decision at a time,
  including dependent questions introduced by enabled capabilities.
- Synchronized the canonical `.agents/skills/setup-context-driven` skill and
  embedded `skills/setup-context-driven` copy.

Acceptance evidence:

- Compatible spec 0030 migration, preserved confirmation dates, obsolete
  inventory removal, clean audit, and idempotent reapply are covered by
  `test_spec0030_manifest_migrates_answers_inventory_and_owned_blocks`.
- Missing dependent decision routing and persistence are covered by
  `test_enabled_capability_routes_only_missing_dependent_decision` and
  `test_enabled_autonomous_work_routes_one_missing_dependent_question`.
- Ambiguous legacy ownership and no-partial-write behavior are covered by
  `test_ambiguous_legacy_ownership_blocks_without_partial_writes`.
- Cross-profile apply, audit, reapply, representative decision combinations,
  required-skill blocking, and extra-skill opt-in reporting are covered by
  `test_supported_profiles_cover_representative_decision_combinations` and
  `test_required_skill_failure_and_extra_reporting_keep_exit_semantics`.
- Canonical/embedded sync and portable execution are covered by
  `rtk make skills-sync-check` and the bundled asset load inside
  `rtk make verify`.

Verification evidence:

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_manifest_migration*.py'` — passed, 3 tests.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_macro_profiles*.py'` — passed, 4 tests.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_workflow*.py'` — passed, 5 tests.
- `rtk make skills-sync-check` — passed.
- `rtk make verify` — passed: setup-context-driven Python suite, `rtk go test ./...`, bundled asset validation for both skill copies, `roundfix skills check`, and Go build.
- `rtk git diff --check` — passed.
