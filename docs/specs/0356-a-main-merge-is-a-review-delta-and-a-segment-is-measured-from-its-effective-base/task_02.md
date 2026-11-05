---
task: task_02
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

Implemented the ordered lineage checks and all six exact restart reasons.
A moved merge base preserves the existing round planning only after head
ancestry, base ancestry, every new merge's secondary parents and the
conflict-free carry have been checked. Git failures propagate as errors.
Round-2 plans compute optional `lineage.deltaBase` before reuse; both delta
prompts and their size measurements start from that object and include the
specified sentence when it is present.

Acceptance evidence from focused checks:

- Round 1 to round 2: `TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch`
  exercises findings, an actual fix, `dispose F1 --fixed-by` and a merge of
  `main`. It asserts round 2, the original reviewed head and findings, the
  fixing disposition, a persisted tree `deltaBase`, its JSON name, the exact
  prompt sentence and its placement, delta markers, and exclusion of
  `main moved on` and the full-candidate marker. Its `oversized carried delta`
  subtest asserts the exact size refusal and stderr, no upstream file in the
  refusal, no segments and no reviewer readiness or prompt call.
- Correction Review: `TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix`
  exercises clean round 1, round-2 findings on `delta.txt`, its fix and
  disposition, and the merge of `main`. It asserts the Correction Review's
  carried tree, persisted record, fix-only prompt, exact sentence placement
  and dispositions, and reuse at the same head without another reviewer call.
- Restarts: `TestReviewLineageRestartNamesItsReason` asserts all seven named
  rows, each round and exact reason. The provider row also rewrites the head
  to prove provider precedence; two unreviewed branch merges prove that the
  earliest offending merge is named. Additional subtests inject failures at
  head/base/parent ancestry, merge listing and carrying, and exercise a real
  Git failure other than exit 1. `TestReviewLineageDecidesEachRound` passes
  every existing row without changes to that test.

Focused commands and outcomes:

- Before the behavioral change, after adding only the new record/plan fields,
  `rtk proxy go test ./internal/cli -run '^(TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch|TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix|TestReviewLineageRestartNamesItsReason)$' -count=1 -timeout=5m`
  exited 1: both command flows restarted at round 1, and every restart reason
  was empty.
- Final regression check:
  `rtk proxy go test ./internal/cli -run '^(TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch|TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix|TestReviewLineageRestartNamesItsReason)$' -count=1 -timeout=5m -v`
  exited 0, including the oversized carried delta and Git-failure subtests.
- Compatibility check:
  `rtk proxy go test ./internal/cli -run '^(TestReviewLineage.*|TestCorrectionReview.*|TestReviewRound.*|TestReviewCeiling.*|TestBlockedCorrectionReviewRepeatsItsDelta|TestADismissedRoundTwoFindingKeepsTheCeiling|TestReviewScopesRoundTwoDelta|TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented|TestReviewBlocksAboveTheBoundWithoutAProviderCall)$' -count=1 -timeout=5m -v`
  exited 0. This covers the existing lineage table, full-candidate anchors,
  session continuation and fallback, ceiling and Correction Review conditions,
  scoped deltas and both providers' existing size refusals.
- `rtk proxy gofmt -l internal/cli/review.go internal/cli/review_lineage.go internal/cli/review_main_merge_lineage_test.go`
  exited 0 with no output; `rtk proxy git diff --check` exited 0.

Scope: production edits are confined to `review.go` and `review_lineage.go`;
the new tests and reusable `mergeAdvancedDefaultBranch` helper are in
`review_main_merge_lineage_test.go`. The pre-existing Daemon status change
in this Task file is preserved. No existing test, other Task, Task Graph,
measurement helper, segmentation code, delivery code or other document changed.
No follow-up work was identified within this slice. Declared Verification,
the full CLI suite, repository gates and Task settlement remain with the Daemon;
none was run or claimed here.
