---
task: task_03
spec: 0227-reconcile-releases-the-runs-of-merged-specs
status: pending
type: backend
complexity: high
---

# Task 03: Reconcile lists and releases the item branches of merged Specs

## Overview

`roundfix reconcile` never reported the item branches of 0205 and 0217 that
the Delivery Queue had recreated, and the operator deleted them by hand
(backlog entry of 2026-10-04,
[2026-10-04-reconcile-keeps-runs-of-merged-specs.md](references/2026-10-04-reconcile-keeps-runs-of-merged-specs.md)).
This Task adds item branch reconciliation: a full scan lists every
`roundfix/deliver-<slug>-<16 hex>` branch that no live item records, proves
it with task_01's merge evidence, and `--apply` deletes it with its clean
worktree. It is verifiable on its own through the worktree package and the
public reconcile command.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-04, "Reconcile keeps the Runs of
   Specs that are already merged", for recreated item branches.
2. MUST implement `_techspec.md` → Item branch reconciliation steps 1 to 5:
   only a full scan inspects item branches; a branch a Delivery Queue item of
   this repository records while its stage is not `merged` is preserved; any
   other branch is releasable only with merge evidence for its slug and a
   proven supersession of its commits by the delivery; a registered worktree
   must be at the derived item path and clean.
3. MUST report item branches as API Contract 3 defines: an
   `itemBranchCandidates` list present, possibly empty, on a full scan and
   absent when a Run ID selects one Run; `debrisSummary` counts
   `itemBranchCandidates` and `itemBranchesApplied`; preserved item branches
   go to `preservedCandidates` with kind `itemBranch` and their reason; the
   other debris entries carry no new field.
4. MUST print the text form of Surface Transcript 1: an
   `Item branch candidate: <branch> (Spec <slug>)` block with `head`,
   `worktree`, `proof`, `action` and `refusal-reason` lines, a proof starting
   `item branch is superseded by the delivery of Spec`, and the debris
   summary line ending with
   `item-branch-candidates=<n> item-branches-applied=<n>`.
5. MUST make `--apply` re-inspect each candidate and refuse it, leaving the
   branch, when its head, delivery commit, default-branch head or worktree
   state changed; otherwise remove a clean worktree without force and delete
   the branch.
6. MUST preserve the branch of a live item, an item branch whose worktree is
   dirty or not at the derived path, and an item branch of a Spec without
   merge evidence.
7. MUST leave `TestRunReconcileJSONMatchesTextFields` in
   `internal/cli/cli_test.go`, a Governed Path this Spec does not touch, and
   every other existing reconcile test passing unedited.
8. MUST NOT touch Run Worktrees or Run Branches through this path, write a
   Branch Disposition record, read GitHub, or delete a branch on a dry-run.

## Subtasks

- [ ] Add item branch inspection and apply beside the item worktree code.
- [ ] Read the live item branches from the Delivery Queue during a full scan.
- [ ] Add the report list, the summary counts and the text block.
- [ ] Add the worktree tests for release, live, unmerged, dirty and moved-head cases.
- [ ] Add the reconcile tests, including the Surface Transcript lines and the Run-ID omission.

## Acceptance Criteria

- [ ] An item branch of a merged Spec is listed by the dry-run with the
      Surface Transcript 1 lines and is gone after `--apply`, its clean
      worktree with it.
- [ ] A live item's branch, a dirty item worktree, an unmerged Spec's branch
      and a moved head each stay.
- [ ] A Run-ID report carries no `itemBranchCandidates` key.

## Context

- creates: `internal/worktree/item_branch_reconcile.go`
- creates: `internal/worktree/item_branch_reconcile_test.go`
- interface: `internal/cli/reconcile.go`
- creates: `internal/cli/reconcile_item_branch_test.go`
- instruction: `internal/worktree/merged_head.go`
- instruction: `internal/cli/merged_release.go`
- instruction: `internal/cli/deliver_workflow.go`
- instruction: `internal/cli/cli_test.go`
- instruction: `docs/user-guide/commands/reconcile.md`
- instruction: `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestItemBranchOfAMergedSpecIsReleasable|TestItemBranchOfALiveItemIsPreserved|TestItemBranchOfAnUnmergedSpecIsPreserved|TestItemBranchWithADirtyWorktreeIsPreserved|TestApplyItemBranchRemovesTheCleanWorktreeAndBranch|TestApplyItemBranchRefusesAMovedHead)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestItemBranchOfAMergedSpecIsReleasable TestItemBranchOfALiveItemIsPreserved TestItemBranchOfAnUnmergedSpecIsPreserved TestItemBranchWithADirtyWorktreeIsPreserved TestApplyItemBranchRemovesTheCleanWorktreeAndBranch TestApplyItemBranchRefusesAMovedHead; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec|TestReconcileKeepsTheItemBranchOfALiveQueueItem|TestReconcileWithARunIDOmitsItemBranchCandidates|TestRunReconcileJSONMatchesTextFields|TestReconcileDryRunReportsStaleStagingWithoutRemovingIt)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec TestReconcileKeepsTheItemBranchOfALiveQueueItem TestReconcileWithARunIDOmitsItemBranchCandidates TestRunReconcileJSONMatchesTextFields TestReconcileDryRunReportsStaleStagingWithoutRemovingIt; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three new tests do not exist, so the command fails.

## References

- `_prd.md` → Goals 3 and 4; Core Feature 3; Success Metric 4
- `_techspec.md` → Item branch reconciliation; API Contract 3; Surface Transcript 1; Testing Approach 4; Build Order 3
- ADR-0232; ADR-0053; ADR-0115
