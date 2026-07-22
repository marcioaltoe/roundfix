---
task: task_05
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
type: backend
complexity: high
---

# Task 05: Enforce the Upgrade Retention Contract

## Overview

Add the runtime boundary that prevents an upgrade from silently removing or
weakening previously managed clauses. Resolved preview and apply must present
the same ordered normative accounting and bind it to the existing Change Plan
authorization.

## Requirements

1. MUST resolve a source baseline from an additive Setup Manifest identity or
   an exact declared legacy fingerprint, never from prose inference.
2. MUST account for every prior mandatory clause as retained, moved, replaced,
   or rejected with the shape and reason required by the transition contract.
3. MUST require every accepted setup-owned target to remain reachable in the
   selected future artifact graph with equivalent enforcement strength.
4. MUST expose ordered retention accounting in text and JSON while preserving
   existing command names, arguments, output discipline, and exit codes.
5. MUST include the normative accounting in `planDigest` so changed mappings
   invalidate prior confirmation.
6. MUST block unknown, incomplete, or weakened transitions before any write
   and add the current baseline identity only after a confirmed successful
   apply.
7. MUST keep canonical and distributed setup skill trees synchronized.

## Subtasks

- [ ] Resolve current and declared legacy baseline identities.
- [ ] Evaluate transition mappings against the future Decision Plan graph.
- [ ] Add retention accounting to Change Plan text and JSON output.
- [ ] Bind ordered accounting into the canonical plan digest.
- [ ] Add the sanitized pre-0.9 upgrade and legacy-manifest fixtures.
- [ ] Prove blocking transitions leave repository bytes unchanged.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] One unaccounted prior mandatory clause blocks preview and apply with
      exit `1`, names the clause, and performs no writes.
- [ ] An unknown legacy fingerprint blocks instead of selecting the nearest
      known baseline.
- [ ] A complete transition lists retained, moved, replaced, and rejected
      entries with reasons in identical text and JSON order.
- [ ] A target with weaker enforcement or no selected carrier blocks the
      transition.
- [ ] Changing only retention accounting changes `planDigest`, and the old
      confirmation is rejected as stale.
- [ ] Confirmed apply records the current baseline identity additively and a
      second resolved preview has no file changes.
- [ ] Existing exits `0`, `1`, `2`, and `3` retain their documented meanings.
- [ ] Canonical and distributed setup skill trees are byte-identical.

## Context

- instruction: `docs/adr/0058-baseline-upgrades-fail-closed-on-unaccounted-rule-removal.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_manifest_migration.py`
- interface: `.agents/skills/setup-context-driven/tests/test_preview.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_upgrade_retention.py'` — expected: complete transitions are authorizable and every unknown, missing, weakened, or stale case blocks without writes.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_manifest_migration.py` — expected: only declared legacy fingerprints receive the additive baseline identity.
- `rtk make verify` — expected: the full repository gate passes with synchronized skill trees.

## References

- `_prd.md` → Goals 1–2; User Story 1; User Story 2; Core Features 1–3 and 11;
  User Experience; Success Metrics.
- `_techspec.md` → Interfaces; Data Models; API Contracts; Testing Approach;
  Build Order 3.
- ADR-0058 → fail-closed clause accounting and Change Plan authorization.
