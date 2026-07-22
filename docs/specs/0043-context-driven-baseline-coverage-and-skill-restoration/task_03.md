---
task: task_03
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: pending
type: backend
complexity: high
---

# Task 03: Reject unresolved Decision Plan references

## Overview

Make audit prove references against the exact artifacts selected by the
Decision Plan rather than incidental files or Markdown syntax. The slice
blocks managed pointers to excluded targets and repository-owned pointers that
are missing, unsafe, or outside the repository.

## Requirements

1. MUST bind every declared reference token to either one managed target or one
   repository-owned relative path.
2. MUST reject a definite setup-owned source whose target is absent from the
   exact definite artifact set, even when a stale target file exists on disk.
3. MUST validate repository-owned targets inside the repository without
   treating them as setup-owned or generating their content.
4. MUST validate the future tree represented by the Decision Plan before apply
   can create a broken pointer.
5. MUST retain relative Markdown-link scanning as defense in depth while using
   typed declarations as the primary path-pointer contract.
6. MUST cover every finite decision branch that includes or excludes a
   referenced artifact across all bundled profiles.

## Subtasks

- [ ] Resolve declared reference tokens into selected target paths.
- [ ] Validate managed targets against definite Decision Plan artifacts.
- [ ] Validate repository-owned paths against repository boundaries and
      existence.
- [ ] Integrate future-tree reference findings into audit and apply planning.
- [ ] Add profile and decision-transition regressions for missing, excluded,
      stale, and external targets.

## Acceptance Criteria

- [ ] Audit blocks a managed reference when its target is excluded or absent
      from the selected Decision Plan.
- [ ] A stale on-disk guide cannot satisfy a managed reference omitted from the
      selected artifact set.
- [ ] The single-context monorepo case audits clean because its referenced
      monorepo guide is selected.
- [ ] A frontend profile with no repository-owned `DESIGN.md` reports one
      blocking path-specific finding and does not create that file.
- [ ] Absolute, escaping, and repository-external declared paths fail without
      traversal.
- [ ] Every artifact-changing boolean or enum branch either resolves all
      references or produces the expected stable blocking finding.
- [ ] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`
- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_plan_contracts.py`
- interface: `.agents/skills/setup-context-driven/tests/test_audit.py`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_decision_plan_contracts.py`
  — expected: every selected reference resolves across all finite artifact
  branches and excluded targets fail deterministically.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_audit.py`
  — expected: managed and repository-owned broken pointers block audit without
  writes or path traversal.
- `rtk make verify` — expected: the full repository gate passes with exact
  Decision Plan reference validation.

## References

- `_prd.md` → Goal 3; User Story 3; Core Features 4, 9; Success Metrics.
- `_techspec.md` → Data Models: typed references; API Contracts: audit;
  Testing Approach; Build Order 3.
- ADR-0046 → managed identity and repository-owned boundary.
- ADR-0047 → exact selected artifact composition.
