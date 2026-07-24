---
task: task_01
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
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

- [x] Add typed catalog values for every new declarative contract.
- [x] Add schema-version routing that preserves supported existing assets.
- [x] Add stable validation diagnostics for each new invariant.
- [x] Create a valid isolated fixture spanning all new contract shapes.
- [x] Add one-field mutation cases for invalid and unsafe shapes.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] A valid fixture loads the same normalized clause, transition, formatter,
      extension, delegation, and dispatch values on repeated reads.
- [x] Duplicate clause or trigger IDs and invalid enforcement or disposition
      values fail with stable diagnostics.
- [x] Unsafe extension paths, malformed formatter metadata, incomplete
      transitions, and invalid delegation aliases fail before repository
      inspection.
- [x] Existing canonical profiles and setup snapshots still load without
      adopting the new runtime behavior prematurely.
- [x] Contract loading performs no writes, process execution, or network
      access.
- [x] Canonical and distributed setup skill trees are byte-identical.

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

## Result

Implemented a version-routed, read-only asset boundary for module/profile v3,
coverage v2, and upgrade-transition v1 documents. The loader now exposes
frozen clause, transition, formatter, Repository-Owned Extension, delegation,
and stable dispatch-trigger values in deterministic order while canonical v1
and v2 profiles retain their existing behavior. An isolated upgrade-contract
fixture and one-field mutation cases exercise the new shapes without migrating
bundled profiles.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_upgrade_contracts.py'`: passed, 6 tests.
- `rtk make skills-sync-check`: passed; canonical and distributed skill trees match.
- `rtk make setup-context-check`: passed, 133 tests in each skill tree plus direct catalog loading.
- `rtk make verify`: passed after allowing access to the existing Go build cache; 1,694 Go tests, 133 tests in each setup skill tree, skill validation, and build passed. The first sandboxed attempt could not read the external Go build cache and produced no test assertion failure.

Acceptance evidence:

- Repeated fixture loads compare equal; clauses, transition mappings, formatter paths, extension paths, delegation aliases, and dispatch triggers are sorted tuples held by frozen dataclasses.
- Mutation tests assert stable diagnostics for duplicate clause/trigger IDs, invalid enforcement/disposition enums, unsafe extension paths, malformed formatter versions, incomplete mappings, invalid aliases, missing clause fields, unknown targets, and malformed trigger collections.
- Canonical and existing v2 fixture loads return empty upgrade-transition, extension, and formatter catalogs, proving the new contracts do not activate runtime behavior prematurely.
- Side-effect guards fail on any `Path.write_text`, `subprocess.run`, or `urllib.request.urlopen` call during catalog loading; fixture bytes remain unchanged.
- `make skills-sync` regenerated the distributed setup skill, and both the focused byte comparison and `make skills-sync-check` prove parity.

Follow-ups: none for this Task slice. Planner, renderer, manifest migration,
delegation scanning, and formatter execution remain assigned to later Tasks.
