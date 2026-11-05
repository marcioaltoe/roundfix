---
task: task_02
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
status: pending
type: backend
complexity: medium
---

# Task 02: A conflict-free merge of the default branch keeps the Reviewer Lineage

## Overview

`decideReviewLineage` starts round 1 whenever the merge base moved, so an item
that merged the default branch after `review-findings` got a full review
again and new findings; Fiscus saw it on Specs 0066 and 0074 on 2026-11-05.
This Task keeps the lineage when the move is only conflict-free merges of the
default branch, starts the round-2 and Correction Review delta from the
reviewed head carried onto the new merge base, and makes every restart name
its reason, as ADR-0366 records. It answers point 2 of the Fiscus report
adopted as
[references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md](references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md).

## Requirements

1. MUST change `decideReviewLineage` (`internal/cli/review_lineage.go`) to
   decide in the order of `_techspec.md` → API Contract 3 and return round 1
   with `reviewLineagePlan.Restart` set to the reason of API Contract 4, word
   for word, at the first failure; when every condition holds after a moved
   merge base, the plan continues exactly as with an unmoved one. A Git
   failure other than exit status 1 of `git merge-base --is-ancestor` stays an
   error, as today.
2. MUST add `DeltaBase string` with JSON name `deltaBase,omitempty` to
   `reviewLineage`, and make `runReviewCommand` (`internal/cli/review.go`),
   for every round-2 plan with a `PreviousHead`, set it before the reuse
   check to `reviewMeasuredFrom(<merge base>, PreviousHead, HEAD)` when that
   differs from `PreviousHead`, else to `""`.
3. MUST make `runConfiguredReviewSession` start the round-2 and Correction
   Review delta from `DeltaBase` when it is set, with its `measuredBase` (the
   start of the delta's size refusal) set to that same object, and make
   `buildRoundTwoPrompt` and `buildCorrectionReviewPrompt` write, when it is
   set, at the places API Contract 3 names, this line:

   ```text
   The candidate merged the default branch since then; the diff starts from tree <deltaBase>, the reviewed head with the default branch's changes applied, so it holds only the candidate's own change.
   ```

4. MUST keep unchanged: every row of `TestReviewLineageDecidesEachRound`, the
   round-2 bound and its refusal, anchors against the full candidate diff,
   session reuse, the ceiling, the Correction Review conditions, and the
   prompts of a lineage whose merge base did not move (`_techspec.md` →
   Invariants 2 to 4).
5. MUST NOT change `planReviewSegments`, `reviewMeasuredFrom`,
   `internal/cli/deliver.go`, any document, or any existing test.
6. MUST add `internal/cli/review_main_merge_lineage_test.go`, opening with the
   `Suite`, `Invariant`, `Boundary IN` and `Boundary OUT` comment lines, with
   the three tests of `_techspec.md` → Testing Approach for task_02:
   - `TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch`;
   - `TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix`;
   - `TestReviewLineageRestartNamesItsReason`, with the rows `clean merge`,
     `provider`, `rewritten head`, `same head`, `base not descending`, `merge
     of another branch` and `conflict`, each asserting `Restart` word for word
     and the round.
   The file defines the helper `mergeAdvancedDefaultBranch(t, fixture, name,
   content)`, which commits one file on `main`, merges `main` into
   `feature/review` with `--no-edit` and returns the new `main` tip; task_03
   reuses it. The two command tests run `runDefaultBranchReview` from task_01
   with `newLineageRunner`, dispose with `review dispose F1 --fixed-by <commit>`,
   and assert that the prompt holds the fix and the line of Requirement 3 and
   holds neither the default branch's change `main moved on` nor
   `BEGIN CANDIDATE DIFF`. Each top-level test calls `t.Parallel()` first and
   uses only `t.TempDir()` repositories.

## Subtasks

- [ ] Decide the lineage across a moved merge base and name each restart.
- [ ] Record `deltaBase`, start the delta from it and say so in both prompts.
- [ ] Add the lineage tests.

## Acceptance Criteria

- [ ] After round-1 findings, a fix and a conflict-free merge of the default branch, the next review is round 2 and reads only the fix.
- [ ] A round-2 fix merged the same way gets its Correction Review on the fix alone.
- [ ] A changed provider, a rewritten head, a moved base at the same head, a base that does not descend, a merge of another branch and a conflicting merge each restart at round 1 with their own reason; the existing lineage table is unchanged.

## Context

- instruction: `docs/adr/0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md`
- interface: `internal/cli/review_lineage.go`
- interface: `internal/cli/review.go`
- creates: `internal/cli/review_main_merge_lineage_test.go`
- instruction: `internal/cli/review_lineage_test.go`

## Consumers

- `decideReviewLineage`, its `Restart`, the `lineage.deltaBase` record field and the round-2 and Correction Review prompts — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch|TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix|TestReviewLineageRestartNamesItsReason|TestReviewLineageDecidesEachRound)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix TestReviewLineageRestartNamesItsReason TestReviewLineageDecidesEachRound; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/cli` — expected: exit 0. Before this Task the three new tests do not exist, so their pass lines are missing; afterwards they and the existing lineage table pass, and the whole CLI suite passes.

## References

- `_prd.md` → Goals; Core Feature 3; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 1 and 2; API Contract 3; API Contract 4; Invariants 1 to 4; Testing Approach; Build Order 2
