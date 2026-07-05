---
task: task_02
spec: 0008-worktree-isolation
status: pending
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

- [ ] Schema v6 migration with populated v5 fixture test
- [ ] work_dir through CreateRun and reads
- [ ] worktree.copy config with validation
- [ ] Artifact Directory builtin default relocation

## Acceptance Criteria

- [ ] Migration test proves row/lock survival and version 6; legacy rows
      read back with empty work_dir.
- [ ] Config tests cover the copy list (valid, absolute rejected, dot-dot
      rejected) and the new builtin artifact default, with explicit-value
      resolution byte-identical to today.
- [ ] Full suite passes with zero review-path assertion changes beyond the
      builtin-default cases deliberately updated.

## Verification

- `rtk go test ./internal/store/ ./internal/config/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 7; Core Feature 6; Decisions. `_techspec.md` → Data
Models, Build Order 2. ADR-0023.
