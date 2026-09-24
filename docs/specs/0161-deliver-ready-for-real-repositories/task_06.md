---
task: task_06
spec: 0161-deliver-ready-for-real-repositories
status: completed
type: backend
complexity: medium
---

# Task 06: Park only what the item touched

## Overview

Corrective Task from the pre-PR review of 2026-09-24. The clean-checkout check reads `git status --porcelain`, which honours `status.showUntrackedFiles=no`, so park's `git clean -fd` can delete an untracked user file; `startingBranch` lives on the shared workflow and leaks from one item to the next, so a park after a dirty-checkout refusal resets exactly the changes that refusal protected; and after a crash the starting branch is taken from the item branch.

## Requirements

1. MUST detect untracked files regardless of `status.showUntrackedFiles`, and on park remove only untracked paths that are new since the item started.
2. MUST let park discard and switch only when the current item passed the clean-checkout check and entered its branch; otherwise park without touching the checkout.
3. MUST record the starting branch on the item in the same write that records the item branch, and restore that recorded branch on park, including after a crash.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With real Git and `status.showUntrackedFiles=no`, a pre-existing untracked file survives a park.
- [ ] A park after a dirty-checkout refusal leaves the checkout untouched.
- [ ] After a crash, park returns the checkout to the recorded starting branch.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/store/delivery.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestParkKeepsUntrackedFilesTheItemDidNotCreate|TestParkNeverTouchesACheckoutTheItemRefused|TestParkRestoresTheRecordedStartingBranchAfterACrash)$" ./internal/cli ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestParkKeepsUntrackedFilesTheItemDidNotCreate TestParkNeverTouchesACheckoutTheItemRefused TestParkRestoresTheRecordedStartingBranchAfterACrash; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Park and resume

## Result

Implemented item-scoped park ownership. Git inspection now forces complete
untracked-file enumeration, the Delivery Queue records the item branch and its
starting branch atomically, and park mutates the checkout only when the current
branch belongs to a nonterminal item with a recorded starting branch. An owned
park resets tracked changes, cleans only the explicitly observed untracked
paths with literal pathspecs, and restores the persisted starting branch.

Acceptance evidence:

- `TestParkKeepsUntrackedFilesTheItemDidNotCreate` passed against real Git with
  `status.showUntrackedFiles=no`; branch creation refused the hidden untracked
  file and park preserved its contents.
- `TestParkNeverTouchesACheckoutTheItemRefused` passed against real Git; after a
  prior item populated workflow history, the next item's dirty-checkout refusal
  preserved branch, HEAD, porcelain status, and tracked file contents.
- `TestParkRestoresTheRecordedStartingBranchAfterACrash` passed after closing
  and reopening the real Run Database; park removed the item-created untracked
  file and returned from the item branch to the persisted `main` branch.

Focused checks:

- `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 -run 'TestPark(KeepsUntrackedFilesTheItemDidNotCreate|NeverTouchesACheckoutTheItemRefused|RestoresTheRecordedStartingBranchAfterACrash)' ./internal/cli` — passed.
- `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 -run '^(TestRecordDeliveryQueueItemBranchKeepsTheFirstBranch|TestOpenMigratesV14DeliveryQueueAddingOwnerAndItemBranch)$' ./internal/store` — passed.
- `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 -run '^TestInspectGit' ./internal/preflight` — passed.
- `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 ./internal/cli` — passed with host process-table permission. The sandboxed run reached two unrelated force-stop integration tests and was blocked by `operation not permitted`; rerunning with the required permission exited 0.
- `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 ./internal/store` and `GOCACHE=/tmp/roundfix-task06-gocache go test -count=1 ./internal/preflight` — passed.
- `git diff --check` — passed.

The Daemon-owned `## Verification` command was not run.
