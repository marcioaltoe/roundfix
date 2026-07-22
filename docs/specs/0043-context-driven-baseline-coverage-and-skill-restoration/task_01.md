---
task: task_01
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: pending
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

- [ ] Add typed catalog values for rules, coverage, dispatch, references, and
      external source integrity.
- [ ] Add stable validation diagnostics for each new invariant.
- [ ] Add valid isolated fixtures for every versioned contract shape.
- [ ] Add mutation cases that remove or corrupt one required field at a time.
- [ ] Preserve deterministic loading of the canonical and embedded catalogs.

## Acceptance Criteria

- [ ] A valid fixture loads normalized rule, dispatch, reference, and external
      source contracts deterministically.
- [ ] Removing one required rule carrier or dispatch mapping makes asset
      validation fail with a stable diagnostic.
- [ ] An unknown managed reference, absolute repository path, mutable external
      ref, unsafe source path, or malformed digest fails before repository
      inspection.
- [ ] Duplicate rule, coverage, dispatch, and reference identifiers are
      rejected rather than resolved by input order.
- [ ] Current canonical and embedded setup catalogs still load successfully.
- [ ] Contract validation performs no filesystem writes, command execution, or
      network access.
- [ ] Canonical and embedded setup skill trees are synchronized after the
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
