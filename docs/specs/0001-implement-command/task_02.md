---
task: task_02
spec: 0001-implement-command
status: completed
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

- [x] Schema v4 DDL and the v3 → v4 migration
- [x] Work-target lock key derivation per Run Kind
- [x] `implement` Kind with by-Kind create-run validation and `spec_slug`
- [x] `ActiveRunInGitRoot` query
- [x] Active Run error wording for spec targets

## Acceptance Criteria

- [x] A migration test opens a populated v3 fixture (runs in several states plus one active lock) and asserts every row survives, the lock is re-keyed, and `user_version` is 4.
- [x] Creating a second Active Run for the same spec target fails with an error naming the target and the blocking run id; creating one for a different spec slug in the same repository succeeds at the lock level.
- [x] Review-Kind lock collisions behave exactly as before, including the error text.
- [x] `ActiveRunInGitRoot` finds Active Runs of any Kind by git root and returns nothing once the Run completes.
- [x] Create-run validation rejects an `implement` request missing `spec_slug` and rejects review requests missing PR fields, each with a named-field error.

## Verification

- `rtk go test ./internal/store/` — expected: all tests pass, including the migration fixture test.
- `rtk go test ./...` — expected: full suite passes; no review-path test changed.

## References

`_prd.md` → User Story 4; Core Feature 2; Decisions (work targets). `_techspec.md` → Data Models (Run Database schema v4), Build Order 2. ADR-0012, ADR-0016.

## Result

### What changed

- The Run Database is now schema v4. `active_run_locks` is keyed by
  `(target_kind, target_key)`: review Kinds lock `("pr",
  "<head_repository>#<head_branch>")`, the new `implement` Kind locks
  `("spec", "<git_root>#<spec_slug>")`. Opening a v3 database migrates it in
  place (adds `spec_slug`, rebuilds the lock table re-keying every row as a
  `pr` target, sets `user_version = 4`); a fresh database creates v4
  directly. An unsupported schema version now fails `Open` loudly.
- `CreateRun` accepts Kind `implement` and validates required fields by
  Kind: review Kinds keep the exact v3 required set; `implement` requires
  Git root, local branch, and Spec slug, and does not require the PR-shaped
  fields. `Run` and `CreateRunRequest` gained `SpecSlug`.
- `ActiveRunError` keeps the review-path text byte-identical; for spec
  targets it names the repository, Spec slug, and blocking run id, and
  points at `roundfix stop <run-id>`.
- New query `Store.ActiveRunInGitRoot(ctx, gitRoot)` returns the Active Run
  of any Kind sharing a Git root (`(Run, bool, error)` like `ActiveRun`),
  backing the ADR-0012 single-working-tree Preflight Validation.
- `spec_slug` and the relaxed PR columns use the store's existing
  empty-string sentinel (`NOT NULL DEFAULT ''`) instead of SQL NULL,
  matching `base_repository`/`completed_at`; optionality is enforced by
  Kind in `CreateRun`. Migrated v3 databases keep plain `NOT NULL` on the
  PR columns, which the store's explicit `''` bindings always satisfy.

### Verification

- `rtk go test ./internal/store/` — pass (35 tests, including the new
  migration fixture, spec-lock, error-text, git-root, and validation tests).
- `rtk go test ./...` — pass (303 tests, 16 packages); no review-path test
  changed.
- `make verify` — pass (fmt-check, full test suite, `roundfix skills
  check`, build).
- `rtk go test -race ./internal/store/` — pass (extra insurance; no new
  goroutines).

### Evidence per acceptance criterion

1. Migration fixture: `TestOpenMigratesV3RunDatabasePreservingRunsAndRekeyingLocks`
   builds a populated v3 database via raw SQL (Active/Clean/Fetched runs +
   one active lock), opens it, and asserts all 3 rows survive field-by-field,
   the single lock row becomes `("pr", "owner/project#feature/review")`,
   `user_version` is 4, and an `implement` Run can be created on the
   migrated database.
2. Spec-target collisions: `TestCreateRunRejectsSecondActiveRunForSameSpecTarget`
   asserts the duplicate fails with the full error text naming the
   repository, Spec slug, and blocking run id, while a different Spec slug
   in the same repository and the same slug in a different repository both
   pass the lock. `TestCompletedImplementRunReleasesSpecTargetLock` covers
   release on completion.
3. Review-path text: `TestReviewKindActiveRunErrorTextUnchanged` asserts the
   exact v3 error string; the pre-existing duplicate-rejection test runs
   unchanged.
4. Git-root query: `TestActiveRunInGitRootFindsActiveRunsOfAnyKind` finds a
   review Run, returns nothing after completion, and finds an `implement`
   Run in the same root; unrelated roots return nothing.
5. By-Kind validation: `TestCreateRunValidatesRequiredFieldsByKind` covers
   `implement` missing Spec slug / Git root / local branch and review
   missing pull request / Head Repository / HEAD, each with the exact
   named-field message, plus unknown and empty Kind.

### Follow-ups for downstream tasks

- task_05/task_06: an `implement` `CreateRunRequest` needs `Kind:
  store.KindImplement`, `GitRoot`, `LocalBranch`, `SpecSlug`; the PR-shaped
  fields stay empty. `ArtifactDir` is NOT required by the store for
  `implement` Runs (the task named only the three fields) — decide there
  whether the daemon must always pass it and tighten validation if so.
- ADR-0012's single-working-tree check is NOT enforced by `CreateRun` (per
  ADR-0016 it is a Preflight Validation, not a lock): callers must invoke
  `ActiveRunInGitRoot` during preflight before creating an `implement` Run.
- The spec-target error text is now public contract: `Active Run already
  exists for repository "<git_root>" and Spec "<slug>"; existing
  run_id=<id> state=<state>; stop it with: roundfix stop <id>`.
