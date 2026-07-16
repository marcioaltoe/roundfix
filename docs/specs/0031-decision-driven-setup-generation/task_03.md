---
task: task_03
spec: 0031-decision-driven-setup-generation
status: completed
type: backend
complexity: medium
---

# Task 03: Make Spec scaffolding and external triage decision-driven

## Overview

Apply declarative boolean effects to the current mixed context-workflow
guidance. Fresh repositories will include local Spec or external-triage rules
only when the corresponding durable decision enables them, while retaining
independent domain and docs-layout guidance.

## Requirements

1. MUST split mixed root and guide ownership so local Spec, domain, docs-layout,
   and external-triage guidance can be selected independently.
2. MUST activate local Spec routing, issue tracking, root rules, and
   Spec-specific docs-layout text only when `spec.scaffold=true`.
3. MUST omit every local Spec rule and guide controlled by
   `spec.scaffold` when its value is false.
4. MUST activate an English external-triage pointer and supporting guide only
   when `triage.external=true`, and omit both when false.
5. MUST derive fresh-repository audit and apply expectations from the shared
   Decision Plan rather than profile-specific branches.
6. MUST preserve compact root pointers, repository-authored bytes, and the
   existing skill setup selection policy.

## Subtasks

- [x] Separate mixed context-workflow artifacts into independently owned units.
- [x] Declare local Spec module/artifact effects for both boolean values.
- [x] Add the conditional external-triage module, pointer, and guide.
- [x] Make docs-layout content consistent with the selected Spec decision.
- [x] Add fresh-repository apply/audit tests for all four boolean combinations.
- [x] Prove false values omit controlled content without changing owner bytes.

## Acceptance Criteria

- [x] `spec.scaffold=false` produces no local Spec root rule,
      Spec-routing guide, issue-tracker guide, or Spec-only docs-layout text.
- [x] `spec.scaffold=true` produces the complete local Spec guidance and a clean
      audit.
- [x] `triage.external=false` produces no external-triage managed artifact;
      `true` produces one compact pointer and its English supporting guide.
- [x] Domain and general docs-layout guidance remain present for every tested
      Spec/triage combination.
- [x] Applying either false value does not remove repository-authored text or
      alter the selected canonical setup snapshot.
- [x] Repeating apply for each combination produces no second file change.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/assets/modules/context-workflow.json`
- interface: `.agents/skills/setup-context-driven/assets/templates/root/context-workflow.md`
- interface: `.agents/skills/setup-context-driven/assets/templates/guides/docs-layout.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_spec_triage_decisions*.py'` — expected: true/false Spec and triage combinations, independent base guidance, owner-byte preservation, clean audit, and idempotency pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_apply*.py'` — expected: existing safe apply, adoption, atomicity, and managed removal behavior remain passing.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → Goals 1, 3, 5; Core Features 1–2, 6; Success Criteria QA-07.
- `_techspec.md` → System architecture: mixed workflow ownership; Data models: Decision effects and Managed artifacts; Build Order 3.
- ADR-0047.

## Result

- Split base `context-workflow` ownership from conditional `spec-workflow` and `external-triage` modules in both canonical and embedded skill assets.
- Added Spec-specific root/docs-layout artifacts and external-triage root/guide artifacts, with boolean effects driven by `spec.scaffold` and `triage.external`.
- Changed concrete audit/apply artifact expectations to use the shared Decision Plan, including inactive artifact ownership for setup-owned removals.
- Added `test_spec_triage_decisions.py` real-process coverage for all four Spec/triage combinations, clean audit, independent domain/docs-layout guidance, owner-byte preservation, and idempotency.
- Evidence: `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_spec_triage_decisions*.py'` passed with 2 tests.
- Evidence: `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_apply*.py'` passed with 8 tests.
- Evidence: `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_assets*.py'`, `test_audit*.py`, `test_preview*.py`, and `test_decision_plan*.py` all passed.
- Evidence: `rtk git diff --check` passed.
- Evidence: `rtk make verify` passed, including 56 Python tests, `rtk go test ./...`, asset loading for canonical and embedded skills, `roundfix skills check`, and Go build.
