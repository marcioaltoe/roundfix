---
task: task_05
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
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

- [x] Resolve current and declared legacy baseline identities.
- [x] Evaluate transition mappings against the future Decision Plan graph.
- [x] Add retention accounting to Change Plan text and JSON output.
- [x] Bind ordered accounting into the canonical plan digest.
- [x] Add the sanitized pre-0.9 upgrade and legacy-manifest fixtures.
- [x] Prove blocking transitions leave repository bytes unchanged.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] One unaccounted prior mandatory clause blocks preview and apply with
      exit `1`, names the clause, and performs no writes.
- [x] An unknown legacy fingerprint blocks instead of selecting the nearest
      known baseline.
- [x] A complete transition lists retained, moved, replaced, and rejected
      entries with reasons in identical text and JSON order.
- [x] A target with weaker enforcement or no selected carrier blocks the
      transition.
- [x] Changing only retention accounting changes `planDigest`, and the old
      confirmation is rejected as stale.
- [x] Confirmed apply records the current baseline identity additively and a
      second resolved preview has no file changes.
- [x] Existing exits `0`, `1`, `2`, and `3` retain their documented meanings.
- [x] Canonical and distributed setup skill trees are byte-identical.

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

## Result

Implemented exact source-baseline resolution from `generator.baseline` or a
declared canonical fingerprint of legacy managed-artifact identities. The
runtime now evaluates every ordered transition mapping against the definite
Decision Plan graph, blocks missing, unreachable, or weakened targets before
mutation, renders the same `retentionAccounting` in text and JSON, and binds
that accounting into `planDigest`. Confirmed apply records the current
baseline without replacing other generator metadata.

Added a sanitized pre-0.9 corpus and legacy Setup Manifest fixture. Focused
tests cover all four dispositions, unknown and incomplete transitions,
enforcement and carrier failures, stale confirmation, additive migration,
idempotent re-preview, stable exits, and byte preservation on every blocked
path. Regenerated the distributed setup skill from the canonical tree.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_upgrade_retention.py'`: passed, 7 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_manifest_migration.py`: passed, 4 tests.
- `rtk make verify`: passed after rerunning the unchanged command with access
  to the existing Go build cache; 1,694 Go tests passed, both setup skill
  trees passed 156 tests, assets and skill synchronization passed, and the
  CLI build completed. The first managed-sandbox attempt could not read the
  external Go build cache and made no repository changes.

Acceptance evidence:

- `test_unaccounted_clause_blocks_preview_and_apply_without_writes` names the
  omitted clause and proves exit `1` plus an unchanged repository snapshot.
- `test_unknown_legacy_fingerprint_blocks_preview_and_apply_without_writes`
  and the manifest-migration negative case prove exact-match, fail-closed
  legacy identity handling without adding a baseline.
- `test_complete_transition_has_identical_ordered_text_and_json_accounting`
  proves stable clause order, reasons, and retained, moved, replaced, and
  rejected dispositions; the confirmed-apply case proves preview/apply parity.
- `test_weaker_target_and_missing_selected_carrier_block` proves equivalent
  enforcement and selected-carrier reachability are mandatory.
- `test_retention_only_change_invalidates_old_confirmation` holds the catalog
  and file delta constant, changes only accounting, observes a new digest,
  and proves the prior confirmation is stale.
- `test_confirmed_apply_records_current_baseline_and_repreview_is_empty`
  proves additive baseline metadata, successful apply, and an empty second
  Change Plan.
- `test_exit_code_meanings_remain_stable` plus the successful apply cover
  exits `0`, `1`, `2`, and `3` with their existing meanings.
- `rtk make verify` ran both mirrored suites and the skill-sync check, proving
  the canonical and distributed trees are byte-identical.

Follow-ups: none.
