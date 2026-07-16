---
task: task_01
spec: 0031-decision-driven-setup-generation
status: pending
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

- [ ] Define the decision-effect and profile-entry-decision asset shapes.
- [ ] Add immutable models for validated conditions, effects, and bindings.
- [ ] Validate target ownership, value types, tokens, duplicates, and cycles.
- [ ] Declare effects for all nine existing decision IDs without adding new
      user-facing decisions.
- [ ] Add mutation-based contract tests and deterministic ordering tests.
- [ ] Keep the canonical and existing portable assets loadable together.

## Acceptance Criteria

- [ ] Loading the canonical catalog yields one validated effect contract for
      every existing decision and the expected entry decisions for every
      profile.
- [ ] Unknown effect targets, incompatible conditions, duplicate bindings,
      undeclared tokens, and dependency cycles each produce a stable diagnostic
      before repository inspection.
- [ ] Loading the same assets twice produces the same ordered profiles,
      modules, decisions, and effects.
- [ ] Existing setup snapshot validation and module-skill reference validation
      remain unchanged and passing.
- [ ] Tests exercise the contracts through public catalog loading rather than
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
