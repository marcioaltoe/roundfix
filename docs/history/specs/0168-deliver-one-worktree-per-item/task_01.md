---
task: task_01
spec: 0168-deliver-one-worktree-per-item
status: completed
type: backend
complexity: low
---

# Task 01: Record each item's worktree

## Overview

A Delivery Queue item records its branch and a starting branch, but no worktree. Resuming or inspecting an item that runs in its own worktree needs the worktree path recorded before the worktree is created, in the same write as the branch.

## Requirements

1. MUST add `Worktree` to the Delivery Queue item, persisted in a `worktree` column added by a migration from the current schema version.
2. MUST open a database at the previous schema version with every existing item readable and its worktree empty.
3. MUST record the item branch and its worktree path in one write, and never overwrite the first recorded pair with a later proposal.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The branch and worktree are recorded together and read back.
- [ ] A second proposal for the same item keeps the first pair.
- [ ] A database at the previous schema version migrates with an empty worktree.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/delivery.go`
- interface: `internal/store/store.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDeliveryItemRecordsItsWorktreeWithItsBranch|TestOpenMigratesDeliveryQueueAddingItemWorktree)$" ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliveryItemRecordsItsWorktreeWithItsBranch TestOpenMigratesDeliveryQueueAddingItemWorktree; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The item record

## Result

Implemented the Delivery Queue item worktree record as schema version 19. New
databases create the `worktree` column, schema version 18 databases add it with
an empty default, and queue reads expose it as `DeliveryQueueItem.Worktree`.
`RecordDeliveryQueueItemWorktree` writes the proposed branch and worktree in
one conditional update and returns the first durable pair on retries.

Acceptance evidence:

- Branch and worktree recorded together and read back:
  `TestDeliveryItemRecordsItsWorktreeWithItsBranch` passed against the real
  SQLite store.
- Second proposal keeps the first pair: the same test proposed a different
  branch and worktree, then observed the original pair in both the return value
  and the persisted queue item.
- Previous schema migrates with an empty worktree:
  `TestOpenMigratesDeliveryQueueAddingItemWorktree` opened a schema version 18
  fixture, preserved its existing item, observed an empty worktree, and matched
  the fresh schema.

Focused checks:

- `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 -run 'TestDeliveryItemRecordsItsWorktreeWithItsBranch|TestOpenMigratesDeliveryQueueAddingItemWorktree' ./internal/store` — passed.
- `GOCACHE=/tmp/roundfix-task01-gocache go test -count=1 ./internal/store` — passed.
- `GOCACHE=/tmp/roundfix-task01-gocache make verify-incremental` — passed after
  rerunning outside the network sandbox; the first attempt was blocked from
  reaching `api.github.com`.

The Task's declared `## Verification` command was not run; Daemon Verification
owns that command and the terminal Task status.
