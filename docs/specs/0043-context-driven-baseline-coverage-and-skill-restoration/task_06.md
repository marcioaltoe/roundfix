---
task: task_06
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: completed
type: backend
complexity: high
---

# Task 06: Restore external Repository Skill Set members

## Overview

Add the explicit restoration surface for missing or drifted external skills in
the selected Repository Skill Set. The command acquires only the declared
immutable source, previews every directory and lock mutation, requires the
same plan digest for apply, and rolls back the repository on any failed proof.

## Requirements

1. MUST add the non-interactive `restore-skills` command contract, structured
   schema, categorized exits, repeatable skill filter, offline source option,
   and plan confirmation defined by the TechSpec.
2. MUST acquire supported GitHub sources by exact commit with argv-only Git,
   prompting disabled, commit identity verified, and one acquisition per
   unique source/ref pair.
3. MUST verify each acquired source subtree against snapshot authority before
   planning any repository write and MUST never fall back to a branch or
   default revision.
4. MUST preview every created, refreshed, and removed skill path plus the exact
   targeted lock edit under one digest-bound Change Plan.
5. MUST atomically swap staged skill directories and the lock file, preserve
   unrelated skills and lock entries, reject stale confirmation, verify final
   tree digests, and roll back every target on failure.
6. MUST keep persisted lock provenance portable and use the isolated lock
   adapter only on already verified temporary bytes; writes remain disabled if
   its compatibility fixture disagrees with Spec 0036 Task 01.
7. MUST bound acquired file count and bytes, reject traversal, links, devices,
   and unsupported providers, and never execute downloaded skill content.

## Subtasks

- [x] Add restore selection, JSON/text output, exit, and help contracts.
- [x] Add exact-commit Git and offline-object-store acquisition adapters.
- [x] Build complete directory and lock operations under the shared plan
      digest.
- [x] Add isolated lock compatibility normalization with portable provenance.
- [x] Add staged directory/lock swap, postwrite proof, and rollback.
- [x] Add security limits and unsafe-tree rejection before mutation.
- [x] Add disposable-Git integration flows and injected-boundary failure cases.

## Acceptance Criteria

- [x] Preview for a drifted nested skill names every file that will be created,
      refreshed, or removed and the exact lock entry that will change.
- [x] Confirmation restores bytes from the declared full commit, produces the
      expected complete-tree digest, and writes no absolute or machine-local
      lock value.
- [x] Multiple selected skills from one source/ref use one acquisition while
      retaining separate path operations and digest proofs.
- [x] Missing confirmation or stale state performs no repository mutation and
      returns the current structured plan; malformed input exits `2`.
- [x] Missing Git, unreachable commit, commit mismatch, unsupported provider,
      digest mismatch, unsafe tree, size breach, or incompatible lock adapter
      exits non-zero before mutation with a specific next action.
- [x] An injected directory-swap, lock-write, or postwrite-verification failure
      restores every targeted directory and original lock byte.
- [x] Unrelated skills and lock entries remain byte-identical after success.
- [x] A second restoration against the resulting repository has an empty plan,
      final audit exits `0`, and no confirmation is required.
- [x] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/specs/0036-doctor-skill-readiness/task_01.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `skills-lock.json`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_restore_skills.py`
  — expected: preview, exact acquisition, confirmation, portable lock, security,
  rollback, final audit, and idempotency flows pass against disposable local
  Git sources.
- `rtk make verify` — expected: the full repository gate passes with the
  canonical and embedded restoration surface synchronized.

## References

- `_prd.md` → User Story 5; Core Features 8–9; User Experience; Non-Goals;
  Success Metrics.
- `_techspec.md` → Interfaces: SkillSource and SkillLockAdapter; API Contracts:
  restore-skills; Integration Points; Risks & Considerations; Build Order 6.
- Spec 0036 Task 01 → external lock compatibility fixture and Doctor ownership
  boundary.
- `docs/agents/skill-governance.md` → prohibition on authorial changes to
  upstream-managed skill content.

## Result

Implemented the non-interactive `restore-skills` surface with deterministic
text and `setup-context-driven/restore-v1` JSON output, categorized exits,
repeatable skill selection, offline Git object-store acquisition, and exact
plan-digest confirmation. Acquisition fetches only the declared full commit,
groups skills by source/ref, reads regular Git blobs without executing source
content, enforces file/byte limits, and verifies each complete-tree digest
before repository planning.

Confirmed plans stage sibling skill directories and lock bytes, atomically
swap only selected targets, verify final tree and lock bytes, and restore every
original target on directory, lock, or postwrite proof failure. The isolated
lock adapter checks the Spec 0036 compatibility fixture and persists only the
portable GitHub repository, full ref, source-relative `SKILL.md` path, and
compatibility hash. Unrelated skill trees and lock entries remain unchanged.

Acceptance evidence:

- `test_preview_names_nested_file_and_exact_lock_operations` proves nested
  create, refresh, remove, and exact lock-entry preview operations.
- `test_confirmation_restores_exact_commit_with_portable_lock_and_is_idempotent`
  proves exact committed bytes, portable lock provenance, unrelated-state
  preservation, an empty second plan, and a final clean audit.
- `test_multiple_skills_from_one_ref_report_one_acquisition` proves one grouped
  acquisition with separate skill changes and digest proofs.
- `test_stale_and_malformed_confirmation_never_mutate` and
  `test_invalid_skill_filter_and_malformed_lock_exit_two_without_mutation`
  prove confirmation and usage exits without repository writes.
- `test_source_and_security_failures_happen_before_mutation`,
  `test_missing_git_and_unreachable_commit_are_actionable_and_read_only`,
  `test_commit_identity_mismatch_is_rejected`, and
  `test_incompatible_lock_adapter_and_unsupported_provider_block_before_write`
  prove the acquisition, provider, digest, unsafe-tree, size, Git, and adapter
  failure contracts with specific next actions.
- `test_injected_swap_lock_and_postwrite_failures_roll_back_all_targets` proves
  byte-exact rollback across all three injected mutation boundaries.

Verification:

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_restore_skills.py`
  — passed, 12 tests.
- `rtk python3 -B skills/setup-context-driven/tests/test_restore_skills.py` —
  passed, 12 tests against the embedded tree.
- `rtk make skills-sync-check` — passed after `rtk make skills-sync` refreshed
  the embedded repo-owned skill.
- `rtk make verify` — passed with 1,687 Go tests, 122 setup-context Python tests,
  both asset catalogs, the Roundfix skill check, and the binary build. The
  first sandboxed attempt could not access the standard Go build cache; the
  unchanged command passed when rerun with that filesystem access.

Follow-ups: none. Doctor behavior remains outside this Task and owned by Spec
0036 Task 01.
