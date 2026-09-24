---
task: task_01
spec: 0162-a-durable-repository-key-per-run
status: pending
type: backend
complexity: medium
---

# Task 01: Record the repository key

## Overview

Spec 0157 finds a Run's repository at query time through `.git/worktrees/*/gitdir`, which `git worktree remove` deletes, so a linked-worktree Run drops out of the repository's listings once its worktree is gone.

## Requirements

1. MUST add `runs.repository_root` in a migration after the current schema version, recording `config.RepositoryRoot(gitRoot)` for every new Run.
2. MUST backfill existing rows whose `git_root` still resolves, and leave the others empty.
3. MUST keep fresh and migrated schemas equivalent, and migrate from every earlier version.
4. MUST list a Run by its recorded key, falling back to the existing alias match only for rows without one.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A new Run records its key.
- [ ] Migration backfills a resolvable row and leaves an unresolvable one empty.
- [ ] With real Git, a Run from a linked worktree is still listed from the main checkout after `git worktree remove`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/config/config.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCreateRunRecordsTheRepositoryKey|TestMigrationBackfillsResolvableRepositoryKeys|TestARemovedWorktreeRunStaysListed)$" ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCreateRunRecordsTheRepositoryKey TestMigrationBackfillsResolvableRepositoryKeys TestARemovedWorktreeRunStaysListed; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The key
