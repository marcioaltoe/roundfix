---
task: task_06
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
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

- [x] Declare the extension decision, module, path, and scaffold template.
- [x] Add the compact root pointer and future-tree reference validation.
- [x] Plan and atomically apply only first creation.
- [x] Exclude repository-owned bytes from the managed inventory.
- [x] Add absent, existing, modified, disabled, and reapply fixtures.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] A true decision and absent extension produce one visible create operation
      and a confirmable plan digest.
- [x] Confirmed apply creates an unmarked extension that is absent from the
      Setup Manifest's managed inventory.
- [x] An existing extension produces no content mutation, even when its bytes
      differ from the original scaffold.
- [x] Audit and reapply preserve existing extension bytes exactly.
- [x] A false decision neither creates nor removes the extension.
- [x] Root-reference validation accepts the planned scaffold in the future
      tree and reports a missing selected extension without rewriting it.
- [x] Canonical and distributed setup skill trees are byte-identical.

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

## Result

Implemented the decision-gated Repository-Owned Extension contract. A resolved
true decision selects the compact root pointer and authorizes one exact,
unmarked scaffold creation through the existing digest-confirmed Change Plan.
The Setup Manifest records that the one-time boundary was crossed without
placing the file in `managedArtifacts` or recording a content digest. Later
audit, reapply, decision changes, and profile transitions treat the extension
as repository-authored bytes and never recreate, rewrite, or remove it.

### Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_repository_extension.py'` — passed, 6 tests.
- `rtk make skills-sync-check` — passed; canonical and distributed skill trees match.
- `rtk make verify` — passed after granting the required Go build-cache access; 1,694 Go tests, both 162-test setup suites, asset loading, skill checks, and the build passed.

### Acceptance evidence

- `test_confirmed_first_creation_is_unmarked_and_outside_managed_inventory`
  proves the single visible create operation, confirmable digest, exact
  preimage/postimage, unmarked scaffold, and managed-inventory exclusion.
- `test_existing_extension_is_never_replaced_or_inventoried` proves differing
  existing bytes produce no extension mutation and never become managed.
- `test_audit_reapply_and_profile_transition_preserve_modified_bytes` proves
  byte-for-byte preservation across audit, reapply, profile transition, and a
  decision change.
- `test_false_decision_neither_creates_nor_removes_extension` proves both
  absent and existing false-decision paths.
- `test_missing_previously_selected_extension_is_reported_without_recreation`
  proves future-tree validation accepts initial creation, then reports a
  later missing selected target without granting recreation authority.
- `test_atomic_failure_rolls_back_the_first_creation` proves the scaffold is
  part of the same rollback-protected atomic write boundary as the manifest
  and managed artifacts.

No follow-up work was discovered within this Task's slice.
