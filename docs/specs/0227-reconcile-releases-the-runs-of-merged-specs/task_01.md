---
task: task_01
spec: 0227-reconcile-releases-the-runs-of-merged-specs
status: pending
type: backend
complexity: high
---

# Task 01: Merge evidence releases a terminal Run of a merged Spec whose files diverged later

## Overview

`roundfix reconcile` kept six Runs of Specs archived and merged on main
because the merged-head proof compares the files of every non-Task commit with
a default branch that kept moving (backlog entry of 2026-10-04,
[2026-10-04-reconcile-keeps-runs-of-merged-specs.md](references/2026-10-04-reconcile-keeps-runs-of-merged-specs.md)).
This Task adds merge evidence as the fallback after that proof fails, so the
Reconcile Command and the delivery owner release such a Run, while every Spec
without merge evidence keeps today's proof. It is verifiable on its own
through the worktree package and the public reconcile command.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-04, "Reconcile keeps the Runs of
   Specs that are already merged", for terminal Runs.
2. MUST implement merge evidence as `_techspec.md` → Merge evidence steps 1
   to 3 define it: the default-branch head holds the Spec's archived
   `_prd.md`; the delivery commit is the commit that added that file, found
   by its path; the delivery commit is not reachable from the Run Branch head;
   every Task commit names the Spec and its Task is completed in the archived
   Spec at the default head; every other commit is superseded.
3. MUST classify a Run proven by merge evidence `superseded` with the reason
   of Merge evidence step 3, and MUST record the default-branch head and name
   as the proof head and ref, so `--apply` re-proves it and refuses when the
   default branch moved after inspection.
4. MUST consult merge evidence only where today's proof would leave the Run
   `unintegrated` (Merge evidence step 4); a `safe` or `superseded` result
   from today's proof MUST keep its reason byte for byte.
5. MUST keep a dirty Run Worktree `dirty` unless every dirty path passes
   ADR-0212's scope rule at the head that holds the archived Spec (Merge
   evidence step 5); a dirty path outside that scope MUST never be released.
6. MUST keep a Run `unintegrated` when its Spec has no archived `_prd.md` on
   the default branch, when its Run Branch holds the delivery commit, when a
   Task commit belongs to another Spec, or when a Task commit's Task is not
   completed in the archived Spec, naming the commit as today.
7. MUST cover the PRD's measured fixture through the public reconcile command:
   a Run that inherited an operator commit from its item branch, with no
   Delivery Queue merge record, after a later default-branch commit edited
   one of that commit's files and renamed another, reports `superseded` in a
   dry-run and `released` after `--apply`.
8. MUST let the delivery owner's release after a merge release an earlier Run
   of the merged Spec whose head is not contained in the recorded candidate
   and whose inherited files differ from the candidate's.
9. MUST rename `TestArchivedSpecDirectoryProofStillComparesOtherCommitPaths`
   to `TestArchivedSpecWithMergeEvidenceSupersedesOtherCommitPaths` and make
   it expect `superseded`, the one declared break; every other existing test
   of `internal/worktree` and `internal/cli` MUST pass unedited.
10. MUST NOT read GitHub or the network, change the Run Database schema, or
    change the content proof, QA supersession or Task status proof for a Spec
    without merge evidence.

## Subtasks

- [ ] Add the merge-evidence proof and its reason beside the merged-head proof.
- [ ] Fall back to it where today's proof leaves a Run unintegrated, dirty branch included.
- [ ] Add the worktree tests, including the stale-apply and refusal cases.
- [ ] Rename and update the declared-break test.
- [ ] Add the reconcile and delivery-release tests through the public entry points.

## Acceptance Criteria

- [ ] The new tests of `internal/worktree/merge_evidence_test.go` pass: a
      diverged Run is `superseded` with and without a merge record, and a
      pending Task, a missing archived `_prd.md`, a held delivery commit, an
      undeclared dirty path and a moved default branch at apply each keep the
      Run.
- [ ] The measured fixture reports `superseded`, then `released`, through
      `roundfix reconcile`.
- [ ] The existing refusal tests named in Verification pass unedited.

## Context

- interface: `internal/worktree/merged_head.go`
- interface: `internal/worktree/worktree.go`
- creates: `internal/worktree/merge_evidence_test.go`
- interface: `internal/worktree/merged_spec_leftovers_test.go`
- creates: `internal/cli/reconcile_merge_evidence_test.go`
- instruction: `internal/worktree/merged_head_test.go`
- instruction: `internal/cli/reconcile_merged_test.go`
- instruction: `internal/cli/deliver_release.go`
- instruction: `docs/user-guide/commands/reconcile.md`
- instruction: `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`
- instruction: `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestMergeEvidenceSupersedesARunWhoseFilesDivergedAfterTheMerge|TestMergeEvidenceSupersedesADivergedRunWithAMergeRecord|TestMergeEvidenceKeepsARunWithAPendingTaskCommit|TestMergeEvidenceRequiresTheArchivedPRD|TestMergeEvidenceKeepsARunThatHoldsTheDeliveryCommit|TestMergeEvidenceKeepsAnUndeclaredDirtyPath|TestMergeEvidenceApplyRefusesAMovedDefaultBranch|TestArchivedSpecWithMergeEvidenceSupersedesOtherCommitPaths|TestMergedHeadRefusesAnUnrepresentedCommit|TestMergedHeadRefusesATaskNotCompletedAtTheMergedHead|TestMergedHeadRefusesAQAReportTheMergedHeadDoesNotSupersede|TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty|TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated)$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestMergeEvidenceSupersedesARunWhoseFilesDivergedAfterTheMerge TestMergeEvidenceSupersedesADivergedRunWithAMergeRecord TestMergeEvidenceKeepsARunWithAPendingTaskCommit TestMergeEvidenceRequiresTheArchivedPRD TestMergeEvidenceKeepsARunThatHoldsTheDeliveryCommit TestMergeEvidenceKeepsAnUndeclaredDirtyPath TestMergeEvidenceApplyRefusesAMovedDefaultBranch TestArchivedSpecWithMergeEvidenceSupersedesOtherCommitPaths TestMergedHeadRefusesAnUnrepresentedCommit TestMergedHeadRefusesATaskNotCompletedAtTheMergedHead TestMergedHeadRefusesAQAReportTheMergedHeadDoesNotSupersede TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; if grep -q "func TestArchivedSpecDirectoryProofStillComparesOtherCommitPaths" internal/worktree/merged_spec_leftovers_test.go; then printf 'stale test name remains\n' >&2; exit 1; fi` — expected: exit 0; before this Task the new tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestReconcileReleasesAMergedSpecRunWhoseFilesDivergedLater|TestDeliverReleaseReleasesAnEarlierRunOfTheMergedSpec|TestReconcilePreservesARunTheMergeRecordDoesNotRepresent|TestReconcileReleasesAnArchivedSpecRunWithoutARecord|TestDeliverReleaseKeepsAnUnrepresentedRun)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReconcileReleasesAMergedSpecRunWhoseFilesDivergedLater TestDeliverReleaseReleasesAnEarlierRunOfTheMergedSpec TestReconcilePreservesARunTheMergeRecordDoesNotRepresent TestReconcileReleasesAnArchivedSpecRunWithoutARecord TestDeliverReleaseKeepsAnUnrepresentedRun; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, and with the tests but without the change the diverged Run stays unintegrated, so the command fails.

## References

- `_prd.md` → Goals 1 and 4; Core Feature 1; Success Metrics 1, 2 and 5; Acceptance evidence
- `_techspec.md` → Merge evidence; API Contract 1; Testing Approach 1 and 2; Build Order 1
- ADR-0232; ADR-0161; ADR-0212; ADR-0053
