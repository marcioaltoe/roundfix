---
task: task_08
spec: 0168-deliver-one-worktree-per-item
status: completed
type: backend
complexity: low
---

# Task 08: Remove only a truly half-removed item worktree

## Overview

Third pre-PR review of 2026-09-25, data-loss finding: `CleanupItem` only asks whether the item branch is checked out anywhere, so an item worktree that is still registered but on another branch, or detached, takes the "unregistered" route and `removeUnregisteredItemWorktree` deletes it with `os.RemoveAll`, uncommitted files included, leaving a dangling registration. Every `deliver resume` retries cleanup of merged items, so this would erase a user's work in that worktree silently.

## Requirements

1. MUST remove a leftover item path only when no registered worktree matches it by canonical path and the admin directory named in its `.git` file no longer exists.
2. MUST refuse removal, keeping the path and reporting the merged item's cleanup as a warning, when the worktree is still registered, on any branch or detached.
3. MUST only remove a path equal to the one `ItemRefFor` derives from the configured location for that item.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, a registered item worktree switched to another branch keeps its uncommitted file through cleanup.
- [ ] A detached item worktree is kept.
- [ ] A truly half-removed item worktree is still recovered.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCleanupKeepsARegisteredItemWorktreeOnAnotherBranch|TestCleanupKeepsADetachedItemWorktree|TestCleanupRecoversAHalfRemovedItemWorktree)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCleanupKeepsARegisteredItemWorktreeOnAnotherBranch TestCleanupKeepsADetachedItemWorktree TestCleanupRecoversAHalfRemovedItemWorktree; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first two cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implemented path-first item worktree cleanup. `CleanupItem` now resolves every
registered worktree by canonical path before selecting a cleanup route. A path
registered on another branch or with detached HEAD returns an error without
removing the path or deleting the recorded item branch; the Delivery Queue's
existing merged-cleanup handling records that error as a warning and continues.

Raw recovery now accepts only an `ItemRef` produced by `ItemRefFor` for the
configured `worktree.location`, with a recorded path exactly equal to the
derived item path. The remaining `.git` file must point directly under this
repository's `worktrees/` administration directory, and that named admin
directory must already be absent before Roundfix makes the leftover tree
writable and removes it.

Acceptance evidence:

- Registered worktree on another branch:
  `TestCleanupKeepsARegisteredItemWorktreeOnAnotherBranch` failed before the
  production change because cleanup returned nil and removed the path. It now
  observes a branch-mismatch refusal and confirms the registered path and its
  uncommitted file retain their original content.
- Detached worktree:
  `TestCleanupKeepsADetachedItemWorktree` failed before the production change
  because cleanup returned nil and removed the path. It now observes the
  detached-registration refusal and confirms the registered path and its
  uncommitted file remain intact.
- Truly half-removed worktree:
  `TestCleanupRecoversAHalfRemovedItemWorktree` now proves that real Git has
  removed both the registration and the admin directory named by the leftover
  `.git` file before `CleanupItem` removes the path and local item branch.

Negative companions:

- `TestCleanupPreservesAnUnregisteredPathWhileItsAdminDirectoryExists` keeps a
  path whose `.git` file names a live administration directory.
- `TestCleanupPreservesAPathNotDerivedForTheItem` keeps a path that differs
  from the configured `ItemRefFor` derivation.
- `TestCleanupPreservesAnUnregisteredPathWithoutThisRepositorysMarker` keeps a
  path whose `.git` file points outside this repository's administration tree.

Focused checks:

- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache go test -count=1 -run '^TestCleanupKeepsARegisteredItemWorktreeOnAnotherBranch$' ./internal/worktree` — passed after the fix.
- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache go test -count=1 -run '^TestCleanupKeepsADetachedItemWorktree$' ./internal/worktree` — passed after the fix.
- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache go test -count=1 -run '^TestCleanupPreservesAPathNotDerivedForTheItem$' ./internal/worktree` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache go test -count=1 -run '^TestCleanupRecoversAHalfRemovedItemWorktree$' ./internal/worktree` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache go test -count=1 -run '^TestAMergedItemLeavesNoWorktreeOrBranch$' ./internal/cli` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache go test -count=1 -run '^TestAFailedMergedCleanupDoesNotStopTheQueue$' ./internal/delivery` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache go test -count=1 ./internal/worktree ./internal/cli` — passed after the sandbox-blocked `api.github.com` dependency was rerun with network permission.
- `rtk env GOCACHE=/tmp/roundfix-task-08-gocache make verify-incremental` — the sandboxed attempt stalled in the existing baseline integration suite; the permitted rerun passed vet, the complete Go suite, skill checks and the build.

The Task's declared `## Verification` command was not run; Daemon Verification
owns that command and the terminal Task status.
