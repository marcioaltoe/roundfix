---
task: task_06
spec: 0168-deliver-one-worktree-per-item
status: completed
type: backend
complexity: medium
---

# Task 06: Provision every item worktree, and let migrated merged items rest

## Overview

Corrective Task from the pre-PR review of 2026-09-25. `worktree.copy` and bootstrap run only in `CreateItem`: a worktree recreated on resume, or one whose owner was killed mid-bootstrap, is used without copied files or bootstrap, so the Run loses files like `.env` and the gate can park the item for no real reason. And every `engine.Run` cleans up `merged` items; a merged item migrated from v18 has no recorded worktree, so cleanup fails and blocks the queue.

## Requirements

1. MUST record, on the item, whether its worktree finished provisioning, and run copy and bootstrap whenever a worktree is recreated or reused without completed provisioning.
2. MUST skip worktree removal for a merged item with no recorded worktree, deleting its local branch only when it is not checked out anywhere.
3. MUST keep a resume of an all-merged queue a no-op.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, a recreated item worktree carries the copied files and the bootstrap output.
- [ ] A worktree whose provisioning was interrupted is completed on reuse.
- [ ] A migrated merged item without a worktree does not block the next item.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRecreatedItemWorktreeIsProvisioned|TestUnfinishedProvisioningIsCompletedOnReuse|TestMigratedMergedItemDoesNotBlockResume)$" ./internal/worktree ./internal/delivery ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRecreatedItemWorktreeIsProvisioned TestUnfinishedProvisioningIsCompletedOnReuse TestMigratedMergedItemDoesNotBlockResume; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implemented durable item-worktree provisioning state as schema version 20.
Each item records `worktree_provisioned`; schema version 19 databases migrate
with the marker false, and ordinary stage updates cannot overwrite the marker
changed by provisioning recovery.

Item creation records the marker only after copy and bootstrap succeed. Resume
clears it before recreating a missing worktree, provisions every recreated or
previously unfinished worktree through `internal/worktree`, and records success
only after both steps finish. Merged items without a recorded worktree use
branch-only cleanup: a branch checked out anywhere is preserved, an unchecked
branch is deleted, and an absent branch is already clean.

Acceptance evidence:

- Recreated worktree carries copied files and bootstrap output:
  `TestRecreatedItemWorktreeIsProvisioned` removed a provisioned real-Git
  worktree, changed both provisioning inputs, resumed it, and observed the
  current copied file, current bootstrap output and a true durable marker.
- Interrupted provisioning is completed on reuse:
  `TestUnfinishedProvisioningIsCompletedOnReuse` failed a real bootstrap after
  it wrote partial output, observed a false marker, reused the same branch and
  worktree, and observed the completed copy, bootstrap output and true marker.
- Migrated merged item does not block the next item:
  `TestMigratedMergedItemDoesNotBlockResume` seeded merged items with branches
  but no worktrees, preserved a checked-out branch, deleted an unchecked
  branch, merged the next queued item and observed no new effects on an
  all-merged resume.

Focused checks:

- `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 ./internal/store ./internal/worktree ./internal/delivery` — passed.
- `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 ./internal/cli -run 'Test(Deliver|EachItem|ItemBranch|EachDelivery|Resume|Park|AMerged|Recreated|Unfinished|Migrated)'` — passed.
- `GOCACHE=/tmp/roundfix-task06-gocache make fmt-check vet skills-sync-check skills-check build` — passed.
- `GOCACHE=/tmp/roundfix-task06-gocache make verify-incremental` — formatting
  and vet passed; the unfiltered Go suite was blocked when an existing test
  attempted unauthorized network access to `api.github.com`. The configured
  `make verify-changed` reached the same environment block. The focused
  offline package checks above passed afterward.

The Task's declared `## Verification` command was not run; Daemon Verification
owns that command and the terminal Task status.
