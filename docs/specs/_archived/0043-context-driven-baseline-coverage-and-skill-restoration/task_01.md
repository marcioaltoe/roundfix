---
task: task_01
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: completed
type: backend
complexity: high
---

# Task 01: Enforce versioned setup asset contracts

## Overview

Create the validated contract boundary for semantic rules, skill dispatch,
typed references, immutable skill sources, and complete-tree digests. This
prefactoring slice keeps the current setup workflow operational while making
each new declarative shape independently testable before bundled assets adopt
it.

## Requirements

1. MUST extend the versioned asset loader with the rule, coverage, dispatch,
   reference, source, and digest contracts from the TechSpec.
2. MUST reject missing, duplicate, unknown, unsafe, mutable, or structurally
   incompatible values with stable diagnostics.
3. MUST prove that required rules belong to selected modules and are reachable
   through declared supporting-guide rule carriers.
4. MUST prove exact equality between each module's required skills and its
   dispatch declarations.
5. MUST validate setup-owned reference targets by managed identity and
   repository-owned targets by safe repository-relative path without
   inspecting a target repository.
6. MUST keep asset loading deterministic, local, read-only, and free of
   profile-specific imperative branches.

## Subtasks

- [x] Add typed catalog values for rules, coverage, dispatch, references, and
      external source integrity.
- [x] Add stable validation diagnostics for each new invariant.
- [x] Add valid isolated fixtures for every versioned contract shape.
- [x] Add mutation cases that remove or corrupt one required field at a time.
- [x] Preserve deterministic loading of the canonical and embedded catalogs.

## Acceptance Criteria

- [x] A valid fixture loads normalized rule, dispatch, reference, and external
      source contracts deterministically.
- [x] Removing one required rule carrier or dispatch mapping makes asset
      validation fail with a stable diagnostic.
- [x] An unknown managed reference, absolute repository path, mutable external
      ref, unsafe source path, or malformed digest fails before repository
      inspection.
- [x] Duplicate rule, coverage, dispatch, and reference identifiers are
      rejected rather than resolved by input order.
- [x] Current canonical and embedded setup catalogs still load successfully.
- [x] Contract validation performs no filesystem writes, command execution, or
      network access.
- [x] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`
- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/tests/test_assets.py`
- interface: `.agents/skills/setup-context-driven/assets/contract-v1.json`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_assets.py`
  — expected: valid contracts load and every targeted malformed contract fails
  with the asserted diagnostic.
- `rtk make verify` — expected: the full repository gate passes with canonical
  and embedded asset catalogs valid.

## References

- `_prd.md` → Goals 1–2; Core Features 1–3, 7; Success Metrics.
- `_techspec.md` → Implementation Design: Interfaces and Data Models; Testing
  Approach; Build Order 1.
- ADR-0046 → declarative ownership and stable managed identities.
- ADR-0047 → declarative Decision Plan effects and shared planning inputs.

## Result

Implemented a dual-version asset boundary that keeps the bundled v1 catalogs
operational while strictly loading isolated v2 coverage, rule, skill-dispatch,
typed-reference, immutable-source, and complete-tree-digest declarations into
frozen normalized values. Validation now rejects missing or duplicate contract
values, unreachable required rules, dispatch inequality, unresolved managed
identities, unsafe repository/source paths, mutable Git revisions, and malformed
digests with stable diagnostics before any target-repository inspection.

Verification:

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_assets.py`:
  passed 11 tests, including the valid v2 fixture and all targeted mutations.
- `rtk make setup-context-check`: passed 84 setup tests and loaded both bundled
  asset catalogs.
- `rtk make verify`: passed after granting access to Go's external build cache;
  1,687 Go tests, 84 setup tests, asset validation, skill synchronization,
  shipped-skill checks, and the CLI build completed successfully.
- `rtk git diff --check`: passed.
- `rtk diff -r .agents/skills/setup-context-driven skills/setup-context-driven`:
  produced no differences after `rtk make skills-sync`.

Acceptance evidence:

- `test_versioned_contract_fixture_loads_normalized_values_deterministically`
  loads the same frozen coverage, rule, dispatch, reference, and external-source
  values on repeated reads.
- `test_versioned_contract_mutations_fail_with_stable_diagnostics` proves rule
  ownership and carrier reachability, the required `artifact.rules` binding,
  exact dispatch equality, managed-reference identity, safe repository paths,
  immutable source refs, safe source paths, and lowercase SHA-256 tree digests.
- `test_versioned_contract_duplicate_identifiers_are_rejected` covers rule,
  coverage, dispatch, and reference identifiers independently.
- `test_canonical_and_embedded_catalogs_load_successfully` preserves current v1
  compatibility, and the full setup check exercises the same load path.
- `test_contract_loading_has_no_write_command_or_network_side_effects` blocks
  writes, subprocess execution, and network calls while loading the valid v2
  fixture, then proves the fixture bytes are unchanged.
- The canonical-to-embedded recursive diff is empty after synchronization.

Follow-up: later Tasks can migrate bundled profile, module, template, and setup
content onto the validated v2 boundary; this slice intentionally does not adopt
or render that content.
