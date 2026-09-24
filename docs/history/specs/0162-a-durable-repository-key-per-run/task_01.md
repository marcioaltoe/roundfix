---
task: task_01
spec: 0162-a-durable-repository-key-per-run
status: completed
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

## Result

Implemented schema version 17 with a non-null `runs.repository_root` key. New
Runs record `config.RepositoryRoot(git_root)`, and migration from every
supported earlier version now adds the column, backfills rows whose checkout
still resolves, and preserves an empty key when resolution fails. Fresh and
migrated databases use the same schema representation.

`ListRuns` now matches the current repository key first. It consults the
existing `git_root` alias set only for rows whose recorded key is empty, so a
legacy Run stays reachable without letting an alias override a recorded
repository identity.

Focused evidence:

- `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test ./internal/store` —
  exited 0 after the final Go test edit. The package migration suite also
  exercised the supported historical schema paths.
- `TestCreateRunRecordsTheRepositoryKey` creates a Run from a real linked
  worktree and observes the main checkout path in both the returned Run and
  `runs.repository_root`.
- `TestMigrationBackfillsResolvableRepositoryKeys` migrates a schema-v16
  database, observes the main checkout key for a resolvable linked-worktree
  row, an empty key for malformed Git metadata, and byte-equivalent fresh and
  migrated schema declarations.
- `TestARemovedWorktreeRunStaysListed` creates and removes a real Git worktree,
  then observes its Run from the main checkout. The legacy-listing coverage
  also proves aliases apply only to empty-key rows.
- `rtk make verify-incremental GOCACHE=/tmp/roundfix-task01-gocache` — exited
  0 with host process-table access. The sandboxed attempt reached all packages
  but exited 2 because two existing force-stop integration tests could not
  enumerate their owned process trees; `internal/store` passed in that attempt.

The Daemon-owned `## Verification` command was not run.
