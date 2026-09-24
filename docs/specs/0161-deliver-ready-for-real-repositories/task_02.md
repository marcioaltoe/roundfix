---
task: task_02
spec: 0161-deliver-ready-for-real-repositories
status: completed
type: backend
complexity: medium
---

# Task 02: Untracked, per-delivery item branches

## Overview

The item branch is created tracking `origin/<default>`, so with `implement.auto_push: true` a Clean Run pushes to the default branch; its name is fixed per slug and created with `git switch -c`, so re-delivery and resume after a crash park with "branch already exists".

## Requirements

1. MUST create the item branch with no upstream.
2. MUST give the branch a per-delivery suffix and record it on the item before any Git command creates it.
3. MUST reuse the recorded branch on resume when it already exists.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git, the item branch has no upstream.
- [ ] Two deliveries of one slug get different branches.
- [ ] A resume after a crash between creation and the stage write reuses the branch.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/store/delivery.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestItemBranchHasNoUpstream|TestEachDeliveryGetsItsOwnBranch|TestResumeReusesTheRecordedItemBranch)$" ./internal/cli ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestItemBranchHasNoUpstream TestEachDeliveryGetsItsOwnBranch TestResumeReusesTheRecordedItemBranch; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Branches

## Result

Implemented durable per-delivery item branches without changing the Task status
or running the Daemon-owned Verification command:

- Branch creation records `roundfix/deliver-<slug>-<random suffix>` on the
  Delivery Queue item before inspecting or changing the Git checkout.
- The first creation refreshes the default branch and uses
  `git switch --no-track -c`, so the local item branch has no upstream.
- A retry keeps the first recorded branch. If that branch already exists
  locally, the workflow switches to it without fetching or creating another
  branch.

Focused-check evidence:

- Before the implementation,
  `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test ./internal/cli -run '^TestItemBranchHasNoUpstream$' -count=1`
  failed because the item branch upstream was `origin/main`.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/cli -run '^(TestDeliveryWorkflowCreatesAnItemBranchFromTheRefreshedDefault|TestItemBranchHasNoUpstream|TestEachDeliveryGetsItsOwnBranch|TestResumeReusesTheRecordedItemBranch)$'`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/store -run '^TestRecordDeliveryQueueItemBranchKeepsTheFirstBranch$'`
  — passed.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 ./internal/cli ./internal/store`
  — the sandboxed run reached two unrelated force-stop integration tests and
  was blocked from reading the process table; the rerun with process-table
  permission passed both affected package suites.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go vet ./internal/cli ./internal/store`
  — passed.
- `rtk git diff --check` — passed.

Acceptance evidence:

1. `TestItemBranchHasNoUpstream` creates the branch with real Git and observes
   an empty `%(upstream:short)` value.
2. `TestEachDeliveryGetsItsOwnBranch` replaces a terminal queue with a second
   delivery of the same slug and observes two different suffixed branch names.
3. `TestResumeReusesTheRecordedItemBranch` leaves the item at `queued` after
   branch creation, makes its remote unavailable, and observes that retry uses
   the recorded local branch. `TestRecordDeliveryQueueItemBranchKeepsTheFirstBranch`
   also proves a later proposal cannot overwrite the durable branch value.
