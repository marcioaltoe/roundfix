---
task: task_02
spec: 0008-worktree-isolation
status: completed
type: data
complexity: medium
---

# Task 02: Record the execution workspace and its provisioning

## Overview

The persistence and configuration ground for isolation: the Run row records
where the Run executes (schema v6 `work_dir`), config gains the
`worktree.copy` provisioning list, and the Artifact Directory's builtin
default leaves the repository tree. Verifiable through store migration and
config tests.

## Requirements

1. MUST migrate the Run Database to schema v6: `runs` gains nullable
   `work_dir` (empty for legacy rows), populated v5 fixture survives
   intact, `user_version` becomes 6, fresh databases create v6 directly;
   `CreateRunRequest` carries the value and reads expose it.
2. MUST add `worktree.copy` to config: a list of repository-relative paths
   (default empty) validated as non-absolute and clean (no `..`); generated
   config output documents it.
3. MUST move the Artifact Directory *builtin* default out of the repo tree
   to Roundfix Home (`artifacts/<repo-id>`); explicitly configured values —
   including repo-relative ones — resolve exactly as today, and the
   generated config comment explains the default.
4. MUST leave every existing store behavior, lock semantics, and
   review-path test unchanged.

## Subtasks

- [x] Schema v6 migration with populated v5 fixture test
- [x] work_dir through CreateRun and reads
- [x] worktree.copy config with validation
- [x] Artifact Directory builtin default relocation

## Acceptance Criteria

- [x] Migration test proves row/lock survival and version 6; legacy rows
      read back with empty work_dir.
- [x] Config tests cover the copy list (valid, absolute rejected, dot-dot
      rejected) and the new builtin artifact default, with explicit-value
      resolution byte-identical to today.
- [x] Full suite passes with zero review-path assertion changes beyond the
      builtin-default cases deliberately updated.

## Verification

- `rtk go test ./internal/store/ ./internal/config/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 7; Core Feature 6; Decisions. `_techspec.md` → Data
Models, Build Order 2. ADR-0023.

## Result

- Schema v6 is implemented in `internal/store`: fresh databases set
  `user_version = 6`, migrations from v3/v4/v5 flow to v6, and `runs.work_dir`
  is nullable. `TestOpenMigratesV5RunDatabasePreservingRunsLocksAndAddingWorkDir`
  proves populated v5 rows and locks survive while legacy `Run.WorkDir` reads
  back empty; v3/v4 migration tests now also assert version 6 and empty legacy
  WorkDir. `TestCreateRunPersistsWorkDir` proves `CreateRunRequest.WorkDir`
  is returned by create/read/active lookup paths.
- Config now supports `worktree.copy`; tests cover project override with valid
  repo-relative entries plus absolute and `..` rejection. Generated config
  includes `worktree.copy` and documents that empty `defaults.artifact_dir`
  uses Roundfix Home `artifacts/<repo-id>`.
- Artifact Directory builtin default now resolves to Roundfix Home
  `.roundfix/artifacts/<repo-id>`. Config tests prove the new builtin default
  and keep explicit relative, home-expanded, and absolute values resolving as
  before. CLI review-path tests were updated only where they deliberately
  asserted the old builtin default artifact location.
- Verification passed:
  `rtk go test ./internal/store/ ./internal/config/` → 55 passed in 2 packages;
  `rtk go test ./...` → 670 passed in 17 packages;
  `rtk make verify` → full gate passed (`go test ./...`, `roundfix skills check`,
  and build).
