---
task: task_01
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
type: backend
complexity: high
---

# Task 01: Establish upgrade compatibility contracts

## Overview

Create the validated asset boundary needed by every upgrade-compatibility
slice without changing the behavior of currently bundled profiles. This
prefactoring task makes the new declarative shapes independently testable
before live modules and planners adopt them.

## Requirements

1. MUST add the versioned clause, baseline transition, formatter,
   Repository-Owned Extension, delegation-alias, and stable dispatch-trigger
   contracts defined by the TechSpec.
2. MUST normalize valid contracts into immutable catalog values with
   deterministic ordering.
3. MUST reject missing, duplicate, unknown, unsafe, or structurally invalid
   values with stable diagnostics before target-repository inspection.
4. MUST retain deterministic loading for the currently supported asset schema
   versions while later Tasks migrate bundled assets.
5. MUST keep catalog loading local, read-only, network-free, and free of
   profile-specific imperative migration branches.
6. MUST keep canonical and distributed setup skill trees synchronized.

## Subtasks

- [ ] Add typed catalog values for every new declarative contract.
- [ ] Add schema-version routing that preserves supported existing assets.
- [ ] Add stable validation diagnostics for each new invariant.
- [ ] Create a valid isolated fixture spanning all new contract shapes.
- [ ] Add one-field mutation cases for invalid and unsafe shapes.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] A valid fixture loads the same normalized clause, transition, formatter,
      extension, delegation, and dispatch values on repeated reads.
- [ ] Duplicate clause or trigger IDs and invalid enforcement or disposition
      values fail with stable diagnostics.
- [ ] Unsafe extension paths, malformed formatter metadata, incomplete
      transitions, and invalid delegation aliases fail before repository
      inspection.
- [ ] Existing canonical profiles and setup snapshots still load without
      adopting the new runtime behavior prematurely.
- [ ] Contract loading performs no writes, process execution, or network
      access.
- [ ] Canonical and distributed setup skill trees are byte-identical.

## Context

- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`
- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/assets/contract-v1.json`
- interface: `.agents/skills/setup-context-driven/tests/test_assets.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_upgrade_contracts.py'` — expected: valid upgrade contracts load deterministically and every targeted mutation fails with its asserted diagnostic.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Core Features 1, 6–10, 12; Non-Goals / Out of Scope.
- `_techspec.md` → Implementation Design: Interfaces and Data Models; Testing
  Approach; Build Order 1.
- ADR-0046 → declarative setup ownership and stable managed identities.
- ADR-0047 → declarative Decision Plan inputs.
