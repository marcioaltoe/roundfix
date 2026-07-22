---
task: task_04
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: completed
type: backend
complexity: high
---

# Task 04: Authorize exact setup Change Plans

## Overview

Replace the split preview/apply projections with one concrete Change Plan that
owns both machine-readable operations and executable bytes. Maintainers can
authorize that exact plan by digest; stale state, ambiguous removal authority,
or an observed delta mismatch produces no unapproved mutation.

## Requirements

1. MUST derive definite public `plannedChanges` and apply mutations from the
   same concrete Change Plan after decisions resolve.
2. MUST report every create, refresh, remove, rename, and reference edit with
   path, managed identity, state, applicable condition, reason, and exact path
   preimage/postimage digests.
3. MUST calculate a deterministic plan digest over the selection, decision
   values, catalog, ordered operations, and path digests while excluding
   volatile timestamps.
4. MUST let read-only audit accept repeated decision inputs and expose an
   authorizable digest only for a fully resolved concrete plan.
5. MUST require `--confirm-plan` for a non-empty apply, preserve categorized
   exits `0/1/2/3`, and keep structured results on stdout and diagnostics on
   stderr.
6. MUST preserve every unmarked file outside proven manifest ownership unless
   an explicit adoption or removal decision appears in the presented plan.
7. MUST stage writes, verify the observed affected-path delta against the
   authorized plan, roll back on failure, and leave a second apply empty.

## Subtasks

- [x] Consolidate public operations and executable file mutations under one
      Change Plan.
- [x] Add exact operation metadata and canonical plan digest calculation.
- [x] Add audit decision inputs and digest-bound apply confirmation.
- [x] Add explicit preserve/remove decisions for ambiguous old inventory.
- [x] Add rename and reference-edit planning without double-counting path
      deltas.
- [x] Add postwrite delta proof and rollback across all affected paths.
- [x] Reproduce the observed omitted-removal and false-positive-removal cases.

## Acceptance Criteria

- [x] The resolved public plan's unique path/digest triples equal the complete
      observed before/after tree delta after apply.
- [x] The Go-to-Rust transition previews every actual managed removal before
      confirmation.
- [x] An unmarked conditional guide outside prior manifest ownership is
      preserved and never appears as a removal candidate.
- [x] Missing confirmation and stale confirmation exit `3` with the current
      plan and no writes; malformed confirmation exits `2`.
- [x] A matching digest applies only the listed mutations and returns
      structured success without mixing diagnostics into stdout.
- [x] Missing or invalid old ownership markers require an explicit
      `preserve|remove` answer before any affected path can change.
- [x] A forced postwrite mismatch or I/O failure restores all original bytes.
- [x] Repeating confirmed apply against the resulting repository reports an
      empty plan and exits `0` without another confirmation.
- [x] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`
- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_preview.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_preview.py`
  — expected: resolved previews expose exact operations and digests; unresolved
  decisions remain conditional and deterministic.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_apply.py`
  — expected: confirmation, stale-state rejection, removal authority, exact
  tree-delta parity, rollback, and idempotency cases pass.
- `rtk make verify` — expected: the full repository gate passes with the
  preview-first apply contract and stable CLI behavior.

## References

- `_prd.md` → Goal 4; User Story 4; Core Features 5–6, 9; User Experience;
  Success Metrics.
- `_techspec.md` → Interfaces: FileMutation and ChangePlan; Data Models:
  PlannedChange and removal authority; API Contracts: audit and apply; Build
  Order 4.
- ADR-0046 → explicit ownership and confirmation boundaries.
- ADR-0047 → one shared Decision Plan across preview, audit, and apply.

## Result

Implemented one concrete, deterministic Change Plan as the authority for both
`plannedChanges` and filesystem writes. Every operation now carries its reason
and exact path preimage/postimage digests; resolved plans carry a canonical
digest over the selected profile, decision values, catalog, ordered operations,
and path digests. Audit accepts repeated decisions, while non-empty apply
requires the matching `--confirm-plan` digest and rejects missing, malformed,
or stale confirmation without writes.

Ambiguous old inventory now introduces an explicit
`removal.<managed-id>=preserve|remove` decision. Proven marker ownership still
permits automatic managed removal, while unmarked paths outside the old Setup
Manifest remain outside the plan. Apply stages exact bytes, verifies every
affected path against the authorized postimage, and restores all original
bytes after an I/O failure or postwrite mismatch. Rename and managed-reference
edits are represented as logical operations on the same path mutations.

Acceptance evidence:

- Exact delta parity, matching confirmation, stdout/stderr separation, and an
  empty second apply pass in
  `test_confirmed_plan_matches_complete_tree_delta_and_second_apply_is_empty`.
- The Go-to-Rust managed-removal path and preservation of unrelated content
  pass in `test_obsolete_managed_artifacts_are_removed_without_unowned_deletions`.
- The false-positive conditional removal case passes in
  `test_unmarked_conditional_guide_outside_manifest_is_never_removed`.
- Missing, stale, and malformed confirmation behavior passes in
  `test_non_empty_apply_requires_matching_plan_digest_without_writes`.
- Explicit preserve/remove authority passes in
  `test_ambiguous_old_inventory_requires_explicit_preserve_or_remove`.
- I/O rollback and forced postwrite-mismatch rollback pass in
  `test_failure_before_commit_preserves_target_files` and
  `test_postwrite_delta_mismatch_restores_every_original_byte`.
- Rename, reference-edit metadata, and complete delta parity pass in
  `test_catalog_path_move_plans_rename_and_reference_edit_without_delta_gaps`.
- Repeated audit decisions produce the same authorizable digest in
  `test_audit_repeated_decisions_exposes_authorizable_concrete_digest`.
- `make skills-sync` refreshed the embedded setup skill from the canonical
  tree; the repository verification gate confirmed the trees match.

Verification:

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_preview.py`
  passed: 5 tests.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_apply.py`
  passed: 17 tests.
- `rtk make verify` passed: 1,687 Go tests across 20 packages, 106 setup
  tests, both asset catalogs, the Roundfix skill check, and the Go build.

Follow-ups: none.

Verification feedback repair:

- Attempts 1 and 2 passed the Go and Python suites but failed the skill distribution
  equality check because generated Python bytecode existed only under the
  canonical setup skill tree.
- The preview and apply subprocess helpers now propagate
  `PYTHONDONTWRITEBYTECODE=1` to the CLI processes they own. After removing the
  generated cache and synchronizing the embedded tree, each focused test leaves
  the trees byte-identical and `rtk make skills-sync-check` exits `0`.
