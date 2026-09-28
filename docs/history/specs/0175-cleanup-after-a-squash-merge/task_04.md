---
task: task_04
spec: 0175-cleanup-after-a-squash-merge
status: completed
type: backend
complexity: medium
---

# Task 04: Release a merged Spec's Runs after a delivery merge

## Overview

For a merged item, `Engine.Run` in `internal/delivery/engine.go` calls
`RemoveItemBranch` and nothing else. The Spec's Run Worktrees and Run Branches
stay until a later, unrelated `implement` preflight or `stop` happens to prune
the few it can prove. Squash-merged Runs are never among them. The maintainer
asked for every finished squash merge to be followed by a cleanup of completed
Runs and branches. This Task makes the delivery owner release every terminal
Run of the merged Spec that the merged-head proof of task_01 covers, before it
removes the item worktree and item branch. The merge record comes from the
Delivery Queue item the engine itself persisted after `mergeCandidate`
verified the merged head. The release removes Git surfaces in the user's
repository, and `roundfix deliver status` shows the maintainer what stayed.

## Requirements

1. MUST add `ReleaseMergedRuns(ctx context.Context, gitRoot string, item
   store.DeliveryQueueItem) error` to `delivery.ItemWorkspace`. `Engine.Run`
   MUST call it for every merged item before `RemoveItemBranch`, and MUST still
   call `RemoveItemBranch` when the release failed. The two errors join into one
   blocker with the existing `warning: cleanup failed` prefix, and a later clean
   pass clears it as today.
2. MUST add `releaseMergedSpecRuns(ctx, runStore *store.Store, repository
   string, merged runworktree.MergedHead) error` in a new
   `internal/cli/deliver_release.go`, doing exactly what the TechSpec's
   "Automatic release after a delivery merge" says:
   - It lists the repository's Runs and keeps implement Runs of the merged
     Spec slug.
   - It leaves Active Runs untouched.
   - It inspects each terminal Run with `GitRoot` set to the repository through
     `InspectTerminalRunMerged`, and releases `safe` and `superseded` results
     through `ApplyTerminalRun`.
   - It returns one error that names each kept Run and its reason. A Run of
     another Spec MUST never be inspected or changed.
3. MUST implement `commandDeliveryWorkflow.ReleaseMergedRuns` in
   `internal/cli/deliver_workflow.go`. It builds the item's `MergedHead` from
   its Spec slug, item branch, last candidate commit, merge commit and pull
   request number, and refuses without mutation an item that lacks a merge
   commit or a candidate head.
4. MUST implement the new method on the fakes in
   `internal/delivery/engine_test.go` and `internal/cli/deliver_test.go`
   without changing what their existing tests assert.
5. MUST state in `docs/user-guide/commands.md` (deliver) and in the Delivery
   queue section of `.agents/skills/roundfix/SKILL.md` that after an item
   merges, Roundfix releases every Run of that Spec proven at the merged head
   before removing the item worktree and branch, and that a Run it cannot prove
   stays and is named in the item's cleanup warning (the phrase `cleanup
   warning`). Then regenerate `skills/roundfix/SKILL.md` with `make
   skills-sync`.
6. MUST put new tests in `internal/delivery/merged_release_test.go` and
   `internal/cli/deliver_merged_release_test.go`, and MUST NOT remove or
   rename an existing test.

## Subtasks

- [ ] Implement the engine call and the command workflow release.
- [ ] Update the deliver guide and the Roundfix skill.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] The engine releases a merged item's Runs before it removes the item branch, and does not release anything for an unmerged item.
- [ ] A kept Run becomes the item's cleanup warning while the item branch is still removed, and a clean retry clears the warning.
- [ ] The workflow releases a contained Run and the Spec 0172 shape of the merged Spec, keeps an unrepresented Run, and never touches an Active Run or another Spec's Run.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/engine_test.go`
- creates: `internal/delivery/merged_release_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_release.go`
- interface: `internal/cli/deliver_test.go`
- creates: `internal/cli/deliver_merged_release_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEngineReleasesMergedRunsBeforeRemovingTheItemBranch|TestEngineRecordsACleanupWarningWhenAMergedRunIsKept|TestEngineClearsTheCleanupWarningAfterACleanRelease|TestEngineReleasesNothingForAnUnmergedItem|TestDeliveryEngineRetriesMergedItemCleanupWithoutReplayingMerge|TestAFailedMergedCleanupDoesNotStopTheQueue|TestDeliverReleaseRemovesEveryProvenRunOfTheMergedSpec|TestDeliverReleaseKeepsAnUnrepresentedRun|TestDeliverReleaseLeavesAnActiveRunAlone|TestDeliverReleaseNeverTouchesAnotherSpecsRun|TestDeliverReleaseRefusesAnItemWithoutAMergeCommit|TestAMergedItemLeavesNoWorktreeOrBranch|TestMigratedMergedItemDoesNotBlockResume)$" ./internal/delivery ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEngineReleasesMergedRunsBeforeRemovingTheItemBranch TestEngineRecordsACleanupWarningWhenAMergedRunIsKept TestEngineClearsTheCleanupWarningAfterACleanRelease TestEngineReleasesNothingForAnUnmergedItem TestDeliveryEngineRetriesMergedItemCleanupWithoutReplayingMerge TestAFailedMergedCleanupDoesNotStopTheQueue TestDeliverReleaseRemovesEveryProvenRunOfTheMergedSpec TestDeliverReleaseKeepsAnUnrepresentedRun TestDeliverReleaseLeavesAnActiveRunAlone TestDeliverReleaseNeverTouchesAnotherSpecsRun TestDeliverReleaseRefusesAnItemWithoutAMergeCommit TestAMergedItemLeavesNoWorktreeOrBranch TestMigratedMergedItemDoesNotBlockResume; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "cleanup warning" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "cleanup warning" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and neither guide says `cleanup warning`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Automatic release after a delivery merge
- `_prd.md` → Core Feature 2; Success Metric 3
- `_techspec.md` → API Contract 3

## Result

Implemented automatic merged-Spec Run release in the delivery owner. Every
merged item now attempts Run release before item-worktree and item-branch
cleanup, attempts both operations even when release fails, joins both failures
under the existing `warning: cleanup failed` blocker, and clears that warning
after a clean retry. Unmerged items do not enter this cleanup path.

The command workflow now validates the persisted merge commit and last
candidate head before reading Runs. It builds the merged-head record from the
Delivery Queue item, selects only Implement Runs of that Spec in the current
repository, preserves and names Active or unproven Runs, and applies the
existing revalidated cleanup only to terminal `safe` and `superseded` results.
Runs of another Spec are excluded before inspection.

Updated the Deliver Command guide and canonical Roundfix skill to describe
the merged-head release and cleanup warning, then regenerated the shipped
skill with `make skills-sync`.

Acceptance evidence:

- `TestEngineReleasesMergedRunsBeforeRemovingTheItemBranch` observes release
  before item-branch removal, while
  `TestEngineReleasesNothingForAnUnmergedItem` observes no cleanup calls for a
  parked item.
- `TestEngineRecordsACleanupWarningWhenAMergedRunIsKept` observes a kept Run
  and its reason in the blocker while item cleanup is still attempted;
  `TestEngineClearsTheCleanupWarningAfterACleanRelease` observes a clean retry
  clear it. `TestEngineJoinsRunReleaseAndItemBranchCleanupFailures` separately
  observes both failures in one warning.
- `TestDeliverReleaseRemovesEveryProvenRunOfTheMergedSpec` uses real Git
  worktrees and the Spec 0172-shaped fixture to release both a contained Run
  and a superseded Task/QA Run.
- `TestDeliverReleaseKeepsAnUnrepresentedRun`,
  `TestDeliverReleaseLeavesAnActiveRunAlone`, and
  `TestDeliverReleaseNeverTouchesAnotherSpecsRun` separately preserve an
  unrepresented Run, an Active Run, and another Spec's Run.
- `TestDeliverReleaseRefusesAnItemWithoutAMergeCommit` and
  `TestDeliverReleaseRefusesAnItemWithoutACandidateHead` observe validation
  before either Git surface is mutated.

Focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 ./internal/delivery` — passed after the last code edit.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 -run '^Test(Deliver|Delivery)' ./internal/cli` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1 ./internal/cli` — passed outside the sandbox in 93.921s. The sandboxed attempt reached two unrelated force-stop integration tests and was denied process-table access; no delivery test failed.
- `rtk make skills-sync` — passed; `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` confirmed byte-identical skill trees.
- `rtk git diff --check` — passed.

The first focused Go invocations could not write the default macOS Go build
cache; rerunning them with the task-local `GOCACHE` above passed. The Task's
declared `## Verification` command was not run; Daemon Verification owns that
command and the terminal Task status.

## Carry-forward provenance

- Source Run: `run_20260928T174529Z_c0cad0da5734a85f`
- Source commit: `93c74622427f7e10b7d1cf5726f70c929c49f7e0`
