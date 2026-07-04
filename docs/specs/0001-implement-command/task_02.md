---
task: task_02
spec: 0001-implement-command
status: pending
type: data
complexity: medium
---

# Task 02: Migrate the Run Database to work-target locks

## Overview

Implement ADR-0016: schema v4 generalizes the Run Database from PR-only identity to work targets, so a spec Run can hold an Active Run lock keyed by `(git_root, spec_slug)` while review Runs keep their existing behavior. Verifiable on its own through store unit tests including a populated v3 migration fixture.

## Requirements

1. MUST re-key `active_run_locks` to `(target_kind, target_key)`: `("pr", "<head_repository>#<pr_head_branch>")` for review Kinds, `("spec", "<git_root>#<spec_slug>")` for spec Runs.
2. MUST add the Run Kind `implement` and a nullable `spec_slug` column; the PR-shaped columns (`head_repository`, `head_branch`, `pr_number`, `head_sha`) become required-by-Kind — review Kinds keep today's validation, `implement` requires `git_root`, `local_branch`, and `spec_slug` instead.
3. MUST migrate v3 → v4 in place, preserving all existing run rows and rewriting existing lock rows to the new key shape; `PRAGMA user_version` becomes 4.
4. MUST add an `ActiveRunInGitRoot` query returning the Active Run (any Kind) sharing a git root, for the ADR-0012 single-working-tree check.
5. MUST keep the review-path Active Run error text unchanged; for spec targets the error names the work target and the blocking run id so the next useful action is `roundfix stop <run-id>`.
6. MUST leave all existing store behavior and the review-path test suite passing unchanged.

## Subtasks

- [ ] Schema v4 DDL and the v3 → v4 migration
- [ ] Work-target lock key derivation per Run Kind
- [ ] `implement` Kind with by-Kind create-run validation and `spec_slug`
- [ ] `ActiveRunInGitRoot` query
- [ ] Active Run error wording for spec targets

## Acceptance Criteria

- [ ] A migration test opens a populated v3 fixture (runs in several states plus one active lock) and asserts every row survives, the lock is re-keyed, and `user_version` is 4.
- [ ] Creating a second Active Run for the same spec target fails with an error naming the target and the blocking run id; creating one for a different spec slug in the same repository succeeds at the lock level.
- [ ] Review-Kind lock collisions behave exactly as before, including the error text.
- [ ] `ActiveRunInGitRoot` finds Active Runs of any Kind by git root and returns nothing once the Run completes.
- [ ] Create-run validation rejects an `implement` request missing `spec_slug` and rejects review requests missing PR fields, each with a named-field error.

## Verification

- `rtk go test ./internal/store/` — expected: all tests pass, including the migration fixture test.
- `rtk go test ./...` — expected: full suite passes; no review-path test changed.

## References

`_prd.md` → User Story 4; Core Feature 2; Decisions (work targets). `_techspec.md` → Data Models (Run Database schema v4), Build Order 2. ADR-0012, ADR-0016.
