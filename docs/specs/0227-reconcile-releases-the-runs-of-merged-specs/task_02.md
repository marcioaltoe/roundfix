---
task: task_02
spec: 0227-reconcile-releases-the-runs-of-merged-specs
status: pending
type: backend
complexity: medium
---

# Task 02: The post-merge cleanup proves a squash merge onto a moved default branch

## Overview

Every Delivery Queue merge on 2026-10-04 whose default branch moved during
the item ended with `warning: cleanup failed: release merged Spec "<slug>"
Runs: candidate head "<sha>" is not represented by merge commit "<sha>"`
(backlog entry of 2026-10-04,
[2026-10-04-reconcile-keeps-runs-of-merged-specs.md](references/2026-10-04-reconcile-keeps-runs-of-merged-specs.md)).
This Task adds the `git merge-tree --write-tree` check to the cleanup's merge
evidence, so a merge commit holding the candidate's changes plus later
default-branch commits is proven, and every other difference still refuses.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-04, "Reconcile keeps the Runs of
   Specs that are already merged", for the post-merge cleanup warning.
2. MUST keep the ancestry and equal-tree checks of the cleanup's merge
   evidence and add, after both miss, the check of `_techspec.md` → The
   cleanup proof for a moved default branch: the merge commit's tree equals
   the single tree line `git merge-tree --write-tree <merge>^1 <candidate>`
   prints with a zero exit.
3. MUST refuse with today's text,
   `candidate head "<head>" is not represented by merge commit "<merge>"`,
   on a conflict, a non-zero exit, a merge commit without a parent, or a
   different tree, so a Git without `--write-tree` fails closed.
4. MUST release the Spec's Runs through `ReleaseMergedRuns` when the merge
   commit is the squash merge of the candidate onto a default branch that
   gained another commit during the item.
5. MUST keep refusing a merge commit whose content differs from that merge,
   such as a squash whose shared file the merge resolved differently from the
   candidate.
6. MUST leave the tests of `internal/cli/deliver_release_evidence_test.go`,
   `internal/cli/deliver_merged_release_test.go` and
   `internal/cli/deliver_release_refresh_test.go` passing unedited.
7. MUST NOT fetch anything beyond the default-branch refresh the cleanup
   already performs, read GitHub, or change the cleanup warning's text.

## Subtasks

- [ ] Add the merge-tree comparison after the ancestry and equal-tree checks.
- [ ] Add the moved-main acceptance test with a disposable repository.
- [ ] Add the altered-merge refusal test.

## Acceptance Criteria

- [ ] A squash merge onto a moved default branch releases the Spec's Run, and
      an altered merge still refuses with "not represented by merge commit".
- [ ] The existing release evidence tests named in Verification pass unedited.

## Context

- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_release_moved_main_test.go`
- instruction: `internal/cli/deliver_release_evidence_test.go`
- instruction: `internal/cli/reconcile_merged_test.go`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `docs/adr/0232-merge-evidence-releases-the-runs-and-item-branches-of-a-merged-spec.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReleaseMergedRunsAcceptsASquashMergeOntoAMovedDefaultBranch|TestReleaseMergedRunsRefusesAMergeCommitThatAltersTheCandidate|TestReleaseMergedRunsRefusesAnUnrelatedCandidateHead|TestReleaseMergedRunsStillReleasesWithValidEvidence|TestReleaseMergedRunsRefusesAMergeCommitOffTheDefaultBranch)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReleaseMergedRunsAcceptsASquashMergeOntoAMovedDefaultBranch TestReleaseMergedRunsRefusesAMergeCommitThatAltersTheCandidate TestReleaseMergedRunsRefusesAnUnrelatedCandidateHead TestReleaseMergedRunsStillReleasesWithValidEvidence TestReleaseMergedRunsRefusesAMergeCommitOffTheDefaultBranch; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, and with the tests but without the change the moved-main merge refuses with not represented by merge commit, so the command fails.

## References

- `_prd.md` → Goal 2; Core Feature 2; Success Metric 3; Acceptance evidence
- `_techspec.md` → The cleanup proof for a moved default branch; API Contract 2; Testing Approach 3; Build Order 2
- ADR-0232; ADR-0161
