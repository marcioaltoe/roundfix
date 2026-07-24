---
task: task_01
spec: 0047-context-driven-guidance-composition
status: pending
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

- [ ] Add hierarchy and semantic-owner catalog declarations.
- [ ] Validate active owner uniqueness and dependency order.
- [ ] Render the compact root precedence map.
- [ ] Add catalog mutation and rendering tests.

## Acceptance Criteria

- [ ] Every active managed concern resolves to exactly one semantic owner.
- [ ] Generated root pointers follow the confirmed hierarchy order.
- [ ] Inactive modules create neither pointers nor semantic destinations.
- [ ] Duplicate or weakening declarations fail catalog validation.
- [ ] Plan and Result JSON fixtures retain their existing schemas.

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
