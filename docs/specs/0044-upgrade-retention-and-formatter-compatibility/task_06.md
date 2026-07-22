---
task: task_06
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
type: backend
complexity: medium
---

# Task 06: Create the Repository-Owned Extension

## Overview

Give project-specific hard rules an explicit repository-owned home without
turning that content into setup-managed policy. The existing Decision Plan and
Change Plan must authorize the one-time scaffold while every later run
preserves repository-authored bytes.

## Requirements

1. MUST add one boolean entry decision that activates the declared extension
   contract and compact managed root pointer.
2. MUST plan creation of the safe unmarked extension path only when the
   decision is true and the file is absent.
3. MUST include the initial scaffold's exact preimage and postimage in the
   existing Change Plan authorization and atomic write boundary.
4. MUST exclude the extension from `managedArtifacts` and every managed digest
   or content audit after creation.
5. MUST preserve an existing extension byte-for-byte across audit, apply,
   profile transition, and decision changes.
6. MUST keep the typed root reference valid without granting setup authority
   to recreate or rewrite existing content.
7. MUST keep canonical and distributed setup skill trees synchronized.

## Subtasks

- [ ] Declare the extension decision, module, path, and scaffold template.
- [ ] Add the compact root pointer and future-tree reference validation.
- [ ] Plan and atomically apply only first creation.
- [ ] Exclude repository-owned bytes from the managed inventory.
- [ ] Add absent, existing, modified, disabled, and reapply fixtures.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] A true decision and absent extension produce one visible create operation
      and a confirmable plan digest.
- [ ] Confirmed apply creates an unmarked extension that is absent from the
      Setup Manifest's managed inventory.
- [ ] An existing extension produces no content mutation, even when its bytes
      differ from the original scaffold.
- [ ] Audit and reapply preserve existing extension bytes exactly.
- [ ] A false decision neither creates nor removes the extension.
- [ ] Root-reference validation accepts the planned scaffold in the future
      tree and reports a missing selected extension without rewriting it.
- [ ] Canonical and distributed setup skill trees are byte-identical.

## Context

- interface: `.agents/skills/setup-context-driven/assets/decisions.json`
- interface: `.agents/skills/setup-context-driven/assets/templates/index.json`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_plan_contracts.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_repository_extension.py'` — expected: only confirmed first creation mutates the unmarked extension and every later flow preserves its bytes.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 5; User Story 6; Core Feature 9; User Experience; Non-Goals
  / Out of Scope.
- `_techspec.md` → Data Models; API Contracts; Risks & Considerations; Build
  Order 4.
- ADR-0046 → preservation outside setup ownership markers.
