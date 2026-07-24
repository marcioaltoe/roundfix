---
task: task_01
spec: 0047-context-driven-guidance-composition
status: completed
type: backend
complexity: high
---

# Task 01: Render the Instruction Hierarchy

## Overview

Make the active Baseline catalog express one deterministic Instruction
Hierarchy and one semantic owner for each managed concern. Generated root
instructions become a compact precedence map that no narrower guide can
weaken.

## Requirements

1. MUST declare the confirmed hierarchy order across universal, context,
   Spec, autonomous, stack, surface, and optional knowledge guidance.
2. MUST derive semantic-owner destinations only from artifacts active in the
   resolved Baseline Profile.
3. MUST render each active guide pointer exactly once and omit inactive or
   empty destinations.
4. MUST reject duplicate ownership, invalid precedence, and a narrower clause
   that weakens universal or confirmed project policy.
5. MUST keep the public Plan and Result schemas unchanged.

## Subtasks

- [x] Add hierarchy and semantic-owner catalog declarations.
- [x] Validate active owner uniqueness and dependency order.
- [x] Render the compact root precedence map.
- [x] Add catalog mutation and rendering tests.

## Acceptance Criteria

- [x] Every active managed concern resolves to exactly one semantic owner.
- [x] Generated root pointers follow the confirmed hierarchy order.
- [x] Inactive modules create neither pointers nor semantic destinations.
- [x] Duplicate or weakening declarations fail catalog validation.
- [x] Plan and Result JSON fixtures retain their existing schemas.

## Context

- instruction: `docs/adr/0074-repository-rules-use-hybrid-semantic-ownership.md`
- interface: `internal/baseline/catalog.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/assets/modules/core.json`

## Verification

- `rtk go test -count=1 ./internal/baseline -run 'TestInstructionHierarchy|TestSemanticOwnerRegistry|TestCatalog'` — expected: ordering, ownership, mutation, and schema-preservation cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1 and 5; User Stories 1, 3, and 4; Core Features 1–4, 14–15.
- `_techspec.md` → System Architecture; Implementation Design: Data Models; Build Order 1.
- ADR-0074 → hybrid semantic ownership.

## Result

The Baseline catalog now declares the complete Instruction Hierarchy and one
semantic owner for every managed guide. Plan assembly orders active root blocks
by that hierarchy, renders the precedence contract in the universal root
block, filters semantic destinations to the resolved active modules and
artifacts, and keeps shared guide paths to one root pointer.

Verification:

- `rtk go test -count=1 ./internal/baseline -run 'TestInstructionHierarchy|TestSemanticOwnerRegistry|TestCatalog'`
  passed with 25 tests.
- `rtk make verify` passed: 2,115 repository tests, 4 setup-skill contract
  tests, the Roundfix skill check, and the final build completed.

Acceptance evidence:

- `TestSemanticOwnerRegistry` proves every active guide has one owner with its
  managed ID, module, path, title, and classifications.
- `TestInstructionHierarchyRendersActivePointersOnce` proves root pointers
  follow universal, context, Spec, autonomous, stack, surface, and optional
  knowledge precedence and each active path appears once.
- The same rendering test and registry test prove inactive external-triage,
  Secondbrain, and repository-specific artifacts create no pointer or semantic
  destination.
- `TestCatalogMutation` rejects invalid hierarchy order, reversed dependency
  order, duplicate classifications, and an explicit narrower-clause weakening
  of universal policy.
- `TestInstructionHierarchyPreservesPlanAndResultSchemas` proves the public
  Plan and Result JSON field sets remain unchanged; the regenerated catalog
  identity fixture records only the intentional embedded catalog changes.
