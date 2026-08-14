---
task: task_01
spec: 0031-decision-driven-setup-generation
status: completed
type: backend
complexity: high
---

# Task 01: Introduce validated Decision Plan contracts

## Overview

Establish the catalog contracts that make decision effects data rather than
special-case Python branches. This prefactoring slice is complete when the
canonical asset loader accepts the nine declared effects and rejects an invalid
effect graph before any repository state is inspected.

## Requirements

1. MUST extend the portable asset contract with profile entry decisions,
   typed decision conditions, module activation, dependent decisions, artifact
   inclusion or exclusion, template selection, and render bindings.
2. MUST represent validated conditions and effects as immutable internal data
   consumed by later Decision Plan resolution.
3. MUST validate every module, artifact, template, binding, and dependent
   decision target against the catalog that owns it.
4. MUST reject type-incompatible conditions, duplicate bindings, undeclared
   template tokens, and decision dependency cycles with stable diagnostics.
5. MUST keep profile and module order deterministic and preserve the existing
   setup snapshot contract.
6. MUST use only the Python standard library and perform no target-repository
   reads or writes during asset validation.

## Subtasks

- [x] Define the decision-effect and profile-entry-decision asset shapes.
- [x] Add immutable models for validated conditions, effects, and bindings.
- [x] Validate target ownership, value types, tokens, duplicates, and cycles.
- [x] Declare effects for all nine existing decision IDs without adding new
      user-facing decisions.
- [x] Add mutation-based contract tests and deterministic ordering tests.
- [x] Keep the canonical and existing portable assets loadable together.

## Acceptance Criteria

- [x] Loading the canonical catalog yields one validated effect contract for
      every existing decision and the expected entry decisions for every
      profile.
- [x] Unknown effect targets, incompatible conditions, duplicate bindings,
      undeclared tokens, and dependency cycles each produce a stable diagnostic
      before repository inspection.
- [x] Loading the same assets twice produces the same ordered profiles,
      modules, decisions, and effects.
- [x] Existing setup snapshot validation and module-skill reference validation
      remain unchanged and passing.
- [x] Tests exercise the contracts through public catalog loading rather than
      test-only production branches.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/assets/contract-v1.json`
- interface: `.agents/skills/setup-context-driven/assets/decisions.json`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/tests/test_assets.py`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_decision_plan*.py'` — expected: decision-effect models, target validation, cycle rejection, and deterministic ordering pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_assets*.py'` — expected: the complete portable asset contract and existing setup validation pass.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → Goals 1, 3, 5; Core Features 1, 7; Non-goals.
- `_techspec.md` → System architecture; Interfaces; Data models: Decision effects; Build Order 1.
- ADR-0047.

## Result

Implemented the asset-level Decision Plan contract slice:

- Added immutable validated models for decision conditions, effects, template selections, and render bindings.
- Extended profiles with ordered entry decisions and decisions with declarative effects for the nine existing decision IDs.
- Validated effect target ownership, condition value types, module and artifact targets, dependent decisions, template selections, render tokens, duplicate bindings, and dependency cycles during public catalog loading.
- Kept the canonical `.agents/skills/setup-context-driven` assets and embedded `skills/setup-context-driven` bundle loadable with the same contract.
- Added mutation-based public loader tests and deterministic ordering coverage in `test_decision_plan_contracts.py`.

Acceptance evidence:

- Canonical catalog effects and profile entry decisions: `test_canonical_catalog_declares_entry_decisions_and_effects` passed.
- Stable invalid-contract diagnostics: `test_invalid_effect_contracts_fail_with_stable_diagnostics` passed for unknown module, artifact, template, and dependent-decision targets; incompatible conditions; duplicate bindings; undeclared template tokens; and dependency cycles.
- Deterministic load order: `test_loading_same_assets_twice_is_deterministic` passed.
- Existing setup snapshot and module-skill validation: `test_assets*.py` passed, and `rtk make verify` loaded both canonical and embedded portable assets.
- Public loader coverage: all new tests call `load_asset_catalog(...)`; no test-only production branch was added.

Verification:

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_decision_plan*.py'`: passed, 5 tests.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_assets*.py'`: passed, 6 tests.
- `rtk git diff --check`: passed.
- `rtk make verify`: passed after rerun with filesystem approval for the Go build cache; 50 Python tests, 1272 Go tests, setup-context asset loading, Roundfix skill check, and build all passed.

Follow-up:

- Decision Plan resolution, generated-artifact splitting, safe string rendering, and audit/apply consumption remain in later task slices.
