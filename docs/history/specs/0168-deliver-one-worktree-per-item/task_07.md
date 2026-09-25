---
task: task_07
spec: 0168-deliver-one-worktree-per-item
status: completed
type: backend
complexity: low
---

# Task 07: A half-removed item worktree never blocks the queue

## Overview

Corrective Task from the second pre-PR review of 2026-09-25. When `git worktree remove --force` fails part-way (for example on a read-only directory a bootstrap wrote inside the item worktree), the registration is gone but the directory remains; every later cleanup refuses with "unregistered path still exists", and `engine.Run` stops before the next queued item.

## Requirements

1. MUST remove, in the no-registration case, a leftover item path whose `.git` file points at this repository's `worktrees/` admin directory, making it writable first.
2. MUST report a failed cleanup of a merged item as a warning on that item and continue with the next queued item, instead of stopping the queue.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, a half-removed item worktree containing a read-only directory is cleaned up on the next attempt.
- [ ] A merged item whose cleanup fails does not stop the next queued item.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/delivery/engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCleanupRecoversAHalfRemovedItemWorktree|TestAFailedMergedCleanupDoesNotStopTheQueue)$" ./internal/worktree ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCleanupRecoversAHalfRemovedItemWorktree TestAFailedMergedCleanupDoesNotStopTheQueue; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implemented guarded recovery for a half-removed item worktree. When Git no
longer registers the recorded path, cleanup now requires a real directory with
a regular `.git` pointer whose target is directly under this repository's Git
`worktrees/` administration directory. It then adds owner permissions to the
remaining tree, removes it and deletes the local item branch. An unregistered
path without that repository marker remains untouched.

Merged-item cleanup failures now persist `warning: cleanup failed: ...` on the
merged item and let the Delivery Queue advance. A later successful cleanup
clears that warning without replaying the merge.

Acceptance evidence:

- Half-removed real-Git worktree recovery:
  `TestCleanupRecoversAHalfRemovedItemWorktree` created an item worktree with a
  read-only directory, observed `git worktree remove --force` drop the
  registration while leaving the path and `.git` pointer, then observed
  `CleanupItem` remove the path and branch. The negative companion
  `TestCleanupPreservesAnUnregisteredPathWithoutThisRepositorysMarker`
  confirmed cleanup refuses and preserves a path whose marker points outside
  this repository's Git administration directory.
- Failed merged cleanup does not stop the queue:
  `TestAFailedMergedCleanupDoesNotStopTheQueue` injected cleanup failure for
  the first of two items, observed its persisted cleanup warning, and observed
  the second item merge and clean up. The existing
  `TestDeliveryEngineRetriesMergedItemCleanupWithoutReplayingMerge` now also
  proves a later retry clears the warning without a second merge effect.

Focused checks:

- `rtk env GOCACHE=/tmp/roundfix-task-07-gocache go test -count=1 -run '^TestCleanupRecoversAHalfRemovedItemWorktree$' ./internal/worktree` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-07-gocache go test -count=1 -run '^TestAFailedMergedCleanupDoesNotStopTheQueue$' ./internal/delivery` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-07-gocache go test -count=1 -run '^TestDeliveryEngineRetriesMergedItemCleanupWithoutReplayingMerge$' ./internal/delivery` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-07-gocache go test -count=1 -run '^TestCleanupPreservesAnUnregisteredPathWithoutThisRepositorysMarker$' ./internal/worktree` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-07-gocache go test -count=1 ./internal/worktree ./internal/delivery` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task-07-gocache make verify-incremental` — the sandboxed attempt was blocked by the existing suite's access to `api.github.com`; the permitted rerun passed vet, the complete Go suite, skill checks and the build.

The Task's declared `## Verification` command was not run; Daemon Verification
owns that command and the terminal Task status.
