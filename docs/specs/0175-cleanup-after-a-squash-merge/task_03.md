---
task: task_03
spec: 0175-cleanup-after-a-squash-merge
status: pending
type: backend
complexity: medium
---

# Task 03: Reconcile reads the Delivery Queue merge record

## Overview

`roundfix reconcile` calls `runworktree.InspectTerminalRun` from
`inspectReconcileRuns` in `internal/cli/reconcile.go`, so it never sees the
merge the delivery owner recorded. `applyReconcileBranchDisposition` can also
replace a Run's classification with the Branch Disposition's, so a Run the
merged-head proof released would read "would discard with
--discard-superseded" instead of "would release with --apply". This Task feeds
the Delivery Queue's merge records into the merged-head proof of task_01. The
records come from the Run Database: `mergeCandidate` in
`internal/delivery/engine.go` wrote them after GitHub reported the pull request
merged at the recorded head. The reconcile report reads them, and `--apply`
acts on them.

## Requirements

1. MUST add `internal/cli/merged_release.go` with
   `loadDeliveryMergedHeads(ctx, reader *store.Store, repository string)
   ([]runworktree.MergedHead, error)`. It reads `reader.DeliveryQueue` for the
   resolved repository path and, when it differs, for the checkout's Git root
   as loaded, and maps each item at stage `merged` with a non-empty merge commit
   and candidate head to a `MergedHead`: Spec slug, item branch, last candidate
   commit, merge commit, pull request number. A missing queue MUST yield no
   records and no error. A queue read failure MUST fail the command as an
   operational failure.
2. MUST make `inspectReconcileRuns` call `runworktree.InspectTerminalRunMerged`
   with those records for every selected Run, in dry-run and in every mutation
   mode.
3. MUST keep a Run's `safe` or `superseded` classification and its apply action
   when the merged-head proof produced it. `applyReconcileBranchDisposition`
   may still change it under `--discard-superseded` only, whose behavior is
   unchanged.
4. MUST keep the `roundfix-reconcile/v1` JSON and text shapes unchanged; the
   `evidence` field carries the proof reason.
5. MUST describe in `docs/user-guide/commands.md` (reconcile) and in the Run
   Worktree reconciliation section of `.agents/skills/roundfix/SKILL.md` that a
   Run of a merged Spec is proven against the merged head. The description MUST
   name the Delivery Queue merge record and the default branch carrying the
   archived Spec as its two sources, use the phrase `merged head`, and update
   the state tables for `safe` and `superseded`. Then regenerate
   `skills/roundfix/SKILL.md` with `make skills-sync`.
6. MUST put new tests in `internal/cli/reconcile_merged_test.go`, driving the
   command through its public runner against a real Git repository and a seeded
   Run Database, and MUST NOT remove or rename an existing test.

## Subtasks

- [ ] Implement the record loading and the merged-head inspection.
- [ ] Update the reconcile guide and the Roundfix skill.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] With a merge record for the Spec, the Spec 0172 shape is reported `superseded` with evidence naming the pull request, and `--apply` removes its Run Worktree and Run Branch.
- [ ] A dry-run removes nothing, a record for another Spec does not release the Run, and an unrepresented commit keeps the Run with a reason naming it.
- [ ] A merged-head-proven Run whose Branch Disposition is superseded keeps the action `would release with --apply`.
- [ ] Without any record, a Run of a Spec archived on the default branch is released by `--apply`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/reconcile.go`
- creates: `internal/cli/merged_release.go`
- creates: `internal/cli/reconcile_merged_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReconcileReleasesARunProvenByTheDeliveryMergeRecord|TestReconcileDryRunLeavesAMergedSpecRunInPlace|TestReconcileIgnoresAMergeRecordForAnotherSpec|TestReconcilePreservesARunTheMergeRecordDoesNotRepresent|TestReconcileKeepsTheApplyActionForAMergedRunWithASupersededDisposition|TestReconcileReleasesAnArchivedSpecRunWithoutARecord|TestReconcileDiscardWritesTheRecordBeforeRemoving|TestReconcileDiscardRefusesAnUnreachableCommit|TestReconcileWithoutTheFlagRemovesNothing|TestRunReconcileSupersededJSONAndApply|TestRunReconcileJSONMatchesTextFields|TestRunReconcileApplyMixedResults)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileReleasesARunProvenByTheDeliveryMergeRecord TestReconcileDryRunLeavesAMergedSpecRunInPlace TestReconcileIgnoresAMergeRecordForAnotherSpec TestReconcilePreservesARunTheMergeRecordDoesNotRepresent TestReconcileKeepsTheApplyActionForAMergedRunWithASupersededDisposition TestReconcileReleasesAnArchivedSpecRunWithoutARecord TestReconcileDiscardWritesTheRecordBeforeRemoving TestReconcileDiscardRefusesAnUnreachableCommit TestReconcileWithoutTheFlagRemovesNothing TestRunReconcileSupersededJSONAndApply TestRunReconcileJSONMatchesTextFields TestRunReconcileApplyMixedResults; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "merged head" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "merged head" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and neither guide says `merged head`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Reconcile reads the merge record
- `_prd.md` → Core Feature 3; Success Metric 1
- `_techspec.md` → API Contracts 1-2
