---
task: task_03
spec: 0228-queue-items-that-stay-current-with-main
status: pending
type: backend
complexity: medium
---

# Task 03: A review-only correction after archive returns the item to review round 2

## Overview

On 2026-10-04 the pre-PR review of Spec 0225's archived candidate raised two
findings on the archived Spec's own records; the operator fixed one in a
docs-only commit, dismissed the other with evidence, and `roundfix deliver
retry` refused with `corrective-spec-required`, so the operator ran review
round 2 and opened Pull Request #382 by hand (Backlog Entry
[a docs fix after archive forces a manual Pull Request](references/2026-10-04-a-docs-fix-after-archive-forces-a-manual-pull-request.md),
recorded 2026-10-04; operator log entries 137, 138 and 158). This Task lets
the retry return such an item to `reviewing` when the workflow proves the
correction answers only that review (ADR-0233).

## Requirements

1. MUST add `ReviewCorrection`, `ReviewCorrectionProver` and the
   `Corrections` engine dependency of the TechSpec's Interfaces to
   `internal/delivery/engine.go`, and implement steps 1, 2 and 4 of "A
   review-only correction returns to round 2" in the
   `corrective-spec-required` branch of `Engine.Retry`: on an accepted proof
   append the head to the candidate commits and set `reviewing`; on a refused
   proof or a proof error refuse with API Contract 3's text, the reason in
   parentheses after the head, leaving the item unchanged.
2. MUST keep the refusal byte-identical to today's when the engine has no
   prover, and keep an unchanged head returning to `reviewing` without a
   proof; the other retry branches, blocker names and Park Classes do not
   change.
3. MUST implement step 3 of the same section as
   `commandDeliveryWorkflow.ProveReviewCorrection` in the new file
   `internal/cli/deliver_review_correction.go`, reading the review record and
   disposition ledger the Pre-PR Review Command persisted for the item
   worktree, and wire it as the engine's `Corrections` in
   `newCommandDeliveryEngine`. The disposition test MUST be one helper in
   `internal/cli/review_lineage.go` shared with the review's ceiling, which
   keeps its behavior.
4. MUST add `internal/delivery/review_correction_retry_test.go` with
   `TestRetryReturnsAReviewOnlyCorrectionToReview` and
   `TestRetryRefusesACorrectionTheProofRejects` over a real queue store and a
   fake prover, and `internal/cli/deliver_review_correction_test.go` with
   `TestAReviewOnlyCorrectionIsProvedFromTheReviewRecord`,
   `TestACorrectionOutsideTheArchivedSpecIsRefused`,
   `TestACorrectionWithAnUndisposedFindingIsRefused` and
   `TestACorrectionThatDoesNotDescendIsRefused` over a real repository and a
   disposable artifact directory, as Testing Approach 3 describes. The tests
   of `internal/delivery/corrective_spec_test.go`,
   `internal/delivery/archived_head_retry_test.go` and
   `internal/cli/review_lineage_test.go` pass unedited, and no test writes
   under the real `~/.roundfix`.

## Subtasks

- [ ] Add the prover dependency and the accepted branch of the corrective retry.
- [ ] Prove ancestry, dispositions and the changed paths from the persisted review.
- [ ] Share the disposition test with the review ceiling.
- [ ] Add the engine and workflow tests.

## Acceptance Criteria

- [ ] An item whose head adds one commit under its archived Spec, with one
      finding fixed by that commit and one dismissed with evidence, resumes at
      `reviewing` with the head as its newest candidate.
- [ ] A change outside the archived Spec's directory, an undisposed finding, a
      fix the head does not contain, a record at another head or a
      non-descendant head refuses, naming the reason, and leaves the item
      unchanged.
- [ ] Every existing corrective-Spec retry, archived-head retry and review
      lineage test passes unedited.

## Context

- interface: `internal/delivery/engine.go`
- creates: `internal/delivery/review_correction_retry_test.go`
- creates: `internal/cli/deliver_review_correction.go`
- creates: `internal/cli/deliver_review_correction_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/review_lineage.go`
- instruction: `internal/cli/review.go`
- instruction: `internal/delivery/corrective_spec_test.go`
- instruction: `internal/delivery/archived_head_retry_test.go`
- instruction: `internal/cli/review_lineage_test.go`
- instruction: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRetryReturnsAReviewOnlyCorrectionToReview|TestRetryRefusesACorrectionTheProofRejects|TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved|TestRetryReturnsAnUnchangedCorrectiveSpecItemToReviewing|TestRetryOfACorrectiveSpecParkStillRefusesAMovedHead|TestRetryOfAnArchivedItemWithACorrectionReturnsToReview)$" ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRetryReturnsAReviewOnlyCorrectionToReview TestRetryRefusesACorrectionTheProofRejects TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved TestRetryReturnsAnUnchangedCorrectiveSpecItemToReviewing TestRetryOfACorrectiveSpecParkStillRefusesAMovedHead TestRetryOfAnArchivedItemWithACorrectionReturnsToReview; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two retry tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestAReviewOnlyCorrectionIsProvedFromTheReviewRecord|TestACorrectionOutsideTheArchivedSpecIsRefused|TestACorrectionWithAnUndisposedFindingIsRefused|TestACorrectionThatDoesNotDescendIsRefused|TestReviewCeilingClosesOnDispositions|TestReviewCeilingBlocksWithoutCallingTheReviewer|TestReviewCeilingClosedRecordRequiresDispositionEvidence)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAReviewOnlyCorrectionIsProvedFromTheReviewRecord TestACorrectionOutsideTheArchivedSpecIsRefused TestACorrectionWithAnUndisposedFindingIsRefused TestACorrectionThatDoesNotDescendIsRefused TestReviewCeilingClosesOnDispositions TestReviewCeilingBlocksWithoutCallingTheReviewer TestReviewCeilingClosedRecordRequiresDispositionEvidence; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four proof tests do not exist, so the command fails.

## References

- `_prd.md` → Goals 3 and 4; Core Feature 3; Success Metric 3; Success Metric 4
- `_techspec.md` → A review-only correction returns to round 2; Interfaces; API Contract 3; Testing Approach 3; Build Order 3
- ADR-0233; ADR-0197; ADR-0223; ADR-0165
