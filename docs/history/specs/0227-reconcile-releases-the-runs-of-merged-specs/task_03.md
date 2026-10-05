---
task: task_03
spec: 0227-reconcile-releases-the-runs-of-merged-specs
status: completed
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

## Result

Implemented the item branch slice answering the 2026-10-04 Backlog Entry,
"Reconcile keeps the Runs of Specs that are already merged".

A full scan now inspects matching local item branches and preserves every
branch recorded by a non-merged Delivery Queue item. Other branches require
task_01's local delivery and Task supersession evidence. Every registered
worktree on the branch must be clean and at the configured derived item path.
Apply reloads Queue ownership and re-inspects the Git evidence, refuses a
changed branch head, default head, delivery or worktree, removes a clean
worktree without force, then deletes only that item branch.

The report includes the item candidate text block, its captured head, proof
and action, and the two debris summary counts. Full-scan JSON carries an
`itemBranchCandidates` array, including `[]`; Run-ID JSON omits it.
Preserved item branches carry their branch, head and refusal reason in
`preservedCandidates`. Other debris entries omit the item-only fields.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Dry-run lists a merged Spec's item branch and apply removes its clean worktree and branch | `TestItemBranchOfAMergedSpecIsReleasable`, `TestApplyItemBranchRemovesTheCleanWorktreeAndBranch` and `TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec` exercise real Git repositories, the Surface Transcript lines, text/JSON counts, dry-run preservation, explicit apply and the subsequent empty full scan. The direct item apply check also asserts the unrelated Run Worktree and Run Branch survive. |
| Live, dirty, unmerged and moved-head branches stay | The named live, unmerged, dirty and moved-head worktree tests, plus `TestReconcileKeepsTheItemBranchOfALiveQueueItem` and `TestReconcileItemBranchSafety`, cover these refusals and assert the branch and remaining worktree survive. Queue tests cover queued, running and parked items, a merged item, and ownership acquired after proof. `TestItemBranchSafetyBoundaries` also covers wrong/duplicate worktree registrations, moved default heads, dirtiness or worktree removal after proof, post-delivery branches, pending/other-Spec Task commits, branches without worktrees, unrelated refs and ambiguous short refs. |
| A Run-ID report has no `itemBranchCandidates` key | `TestReconcileWithARunIDOmitsItemBranchCandidates` checks the raw JSON key and proves Run-ID apply leaves the item branch/worktree intact. `TestReconcileEmptyFullScanIncludesItemBranches` proves an empty full scan emits `[]`. The unchanged `TestRunReconcileJSONMatchesTextFields` also passed in the focused reconcile selection. |

### Focused checks

- The pre-implementation check of `TestItemBranchSafetyBoundaries`, with a
  task-scoped cache, failed compilation because `ItemBranchReconciliation`,
  `InspectItemBranches` and `ApplyItemBranch` did not exist.
- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy go test ./internal/worktree ./internal/cli -run 'ItemBranch|ItemBranches|Reconcile' -count=1`
  exited 0: worktree 6.400s, CLI 28.646s. This exercised the new item tests and
  existing reconcile tests without editing any existing test.
- The initial `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk make verify-incremental`
  exited 2 after the test target failed on sandbox restrictions: two existing
  CLI force-stop tests could not read the process table, and existing daemon
  and router tests could not bind localhost HTTP listeners. Formatting and
  vet passed, and the worktree package passed. The same incremental check
  was rerun with the required process-table and localhost access.
- The retry of `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk make verify-incremental`
  with the required test access exited 0: formatting, vet, the full Go suite,
  skill sync/checks and build passed. CLI reran in 165.038s, daemon in
  41.872s and router in 10.256s; other packages reused the incremental cache.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

### Handoff boundary

No declared Task Verification command was run. The Daemon-owned
`status: in_progress` remains unchanged, and Task settlement is left to the
Daemon. No existing test, other Task, Task Graph, skill, dependency or tooling
configuration was edited. The user guide and skill descriptions already
arrived with this behavior documented by task_04. No commit, push or Pull
Request was created. No follow-up outside this slice was required.

## Carry-forward provenance

- Source Run: `run_20261005T123704Z_9101d15c662d2abe`
- Source commit: `ed71d3e254693ac2efcbfdc33865290b16623406`
