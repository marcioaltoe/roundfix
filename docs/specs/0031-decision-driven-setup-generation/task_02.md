---
task: task_02
spec: 0031-decision-driven-setup-generation
status: completed
type: backend
complexity: high
---

# Task 02: Expose conditional setup preview

## Overview

Deliver the first observable Decision Plan path through the real CLI. A
maintainer can select a profile in an empty repository and receive the profile,
skill setup, unresolved decisions, definite operations, and conditional
operations without allowing audit or blocked apply to write any file.

## Requirements

1. MUST resolve profile entry decisions and activated dependent decisions to a
   deterministic fixed point.
2. MUST expose the selected profile, canonical setup, and active or conditional
   modules in text and JSON results.
3. MUST produce definite and conditional `plannedChanges`, with every
   conditional operation naming its decision ID and triggering value.
4. MUST let `audit --profile` preview a missing-manifest repository instead of
   returning before profile resolution.
5. MUST return the same selection and preview from blocked apply while
   retaining the existing exit-code precedence and JSON envelope.
6. MUST perform no writes when the manifest or any required decision is
   unresolved, invalid, or missing.

## Subtasks

- [x] Implement fixed-point resolution for entry and dependent decisions.
- [x] Build prospective artifact plans for definite and conditional branches.
- [x] Add additive selection and planned-change fields to text and JSON output.
- [x] Route missing-manifest audit through prospective profile resolution.
- [x] Reuse the prospective plan in blocked apply without a second planner.
- [x] Add real-process preview, determinism, exit-code, and no-write tests.

## Acceptance Criteria

- [x] The QA-03 TypeScript/Bun first-run command returns exit code `3`, all
      unresolved entry decisions, and a non-empty managed-change preview.
- [x] A conditional artifact reports `state: conditional` and the exact
      decision/value that would activate it; definite artifacts omit a
      condition.
- [x] Audit and blocked apply report semantically identical selection and
      planned operations for the same profile and answers.
- [x] Repeating preview against unchanged input produces identical semantic
      output.
- [x] File snapshots taken before missing-manifest audit and unanswered apply
      match byte-for-byte afterward.
- [x] Existing text/JSON finding fields and exit codes remain compatible.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_audit.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`
- interface: `docs/specs/0030-context-driven-agent-instructions/qa/evidence/2026-07-16-manual-variations-03/qa-03-first-apply-questions-and-preview.json`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_preview*.py'` — expected: selection, definite/conditional preview, audit/apply agreement, deterministic output, exit codes, and no-write behavior pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — expected: the existing read-only audit contract remains passing.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → Goals 2, 3, 5; Core Features 5–6; Success Criteria QA-03.
- `_techspec.md` → Interfaces: DecisionPlan; API contracts: Read-only preview and Apply; Build Order 2.
- ADR-0047.

## Result

Implemented the CLI-visible conditional setup preview slice:

- Added a shared Decision Plan resolver that starts from profile entry decisions, adds dependent decisions only after activating effects match confirmed answers, and returns deterministic active or conditional module selection.
- Added prospective planned changes for definite artifacts and conditional branches, including `state` and branch conditions in JSON.
- Routed `audit --profile` with a missing manifest through profile resolution instead of returning `manifest.missing`.
- Reused the same Decision Plan and planned-change builder for blocked `apply`, so missing-answer audit and apply previews agree and perform no writes.
- Preserved existing finding objects, JSON envelope fields, and exit-code precedence.

Acceptance evidence:

- QA-03 first-run behavior: `test_typescript_first_apply_returns_entry_decisions_and_preview` passed with exit code `3`, six unresolved entry decisions, and non-empty planned changes.
- Conditional and definite planned changes: the same test verified `guide.autonomous-work` reports `state: conditional` with `{"decisionId": "autonomous.enabled", "equals": true}`, and a definite change omits `condition`.
- Audit/apply agreement and no writes: `test_audit_and_blocked_apply_share_preview_without_writes` passed and compared file snapshots before and after both commands.
- Determinism: `test_preview_output_is_deterministic` passed by comparing repeated JSON payloads.
- Text output: `test_text_preview_names_selection_and_planned_changes` passed for profile/setup/module and planned-change rendering.
- Existing audit compatibility: `test_audit*.py` passed without changing existing finding shapes or exit-code assertions.

Verification:

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_preview*.py'`: passed, 4 tests.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'`: passed, 7 tests.
- `rtk git diff --check`: passed.
- `rtk make verify`: passed; 54 setup-context Python tests, 1272 Go tests, canonical and embedded asset loading, Roundfix skill check, and build passed.

Follow-up:

- Concrete omission/removal of decision-controlled artifacts, template rendering, rendered digests, and semantic audit repair remain in later task slices.
