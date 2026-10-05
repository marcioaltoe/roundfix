---
task: task_03
spec: 0228-queue-items-that-stay-current-with-main
status: completed
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

## Result

Implemented the bounded review-correction retry path. `EngineDependencies`
now accepts a `ReviewCorrectionProver`; the command workflow provides it. A
proved moved head is appended to the candidate commits and resumes at
`reviewing`. A refused proof or proof error names its reason in the existing
corrective-Spec refusal and does not persist any item change. Proof errors
are also logged. With no prover, the refusal text remains unchanged; an
unchanged head bypasses the proof.

The workflow reads the item checkout's persisted review record and the
artifact directory's disposition ledger. It requires a findings record at
the parked candidate, descendant ancestry, one valid disposition for every
standing finding, and changed paths confined to the named archived Specs.
Fixing commits must lie between the candidate and current head. The shared
`reviewDispositionValidAtHead` helper preserves the ceiling's existing
head-containment rule while adding the candidate bound for correction proof.
No queue schema, blocker name, Park Class or other retry branch changed.

### Acceptance evidence

- Accepted correction: `TestAReviewOnlyCorrectionIsProvedFromTheReviewRecord`
  uses a disposable real repository with one archive correction commit, one
  finding fixed by that commit and one dismissed with evidence. It also
  asserts that the persisted record and new head select review round 2.
  `TestRetryReturnsAReviewOnlyCorrectionToReview` uses the real queue store
  and a fake prover to assert `reviewing`, the appended candidate and cleared
  blocker.
- Refused correction: `TestACorrectionOutsideTheArchivedSpecIsRefused`
  covers an unrelated file, a sibling archive directory sharing the slug
  prefix and another Spec. `TestACorrectionWithAnUndisposedFindingIsRefused`
  covers missing/duplicate dispositions, empty dismissal evidence, mismatched
  finding text, dispositions at another head, a fix preceding the candidate
  and a fix absent from the proposed head.
  `TestACorrectionThatDoesNotDescendIsRefused` covers non-descendant ancestry;
  `TestReviewCorrectionRequiresTheParkedReviewRecord` covers another record
  head, another repository and an outcome other than findings.
  `TestRetryRefusesACorrectionTheProofRejects` asserts the exact reason-bearing
  refusal and unchanged persisted item for both refusal and proof error,
  including an error returned alongside an accepted result.
- Existing behavior: `corrective_spec_test.go`, `archived_head_retry_test.go`
  and `review_lineage_test.go` remain unedited. Their tests passed in the
  incremental repository check. All new repositories, records, ledgers and
  queue stores use disposable directories, not the real `~/.roundfix`.

### Focused checks

- Initial red check:
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test ./internal/delivery -run '^TestRetryReturnsAReviewOnlyCorrectionToReview$' -count=1`
  exited 1 because `ReviewCorrection` and the prover dependency were absent.
- Focused check after implementation:
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -count=1 ./internal/delivery ./internal/cli -run 'Test(Retry.*(Correction|Corrective)|AReviewOnlyCorrection|ACorrection|ReviewCorrection|ReviewCeiling|ReviewLineage)'`
  exited 0 for both packages. Earlier fixture runs exposed missing parent
  directories and a missing Specs Root; the fixture now creates those paths.
- `GOCACHE=/tmp/roundfix-task03-gocache rtk make verify-incremental`
  exited 2 in the sandbox: existing process-owner tests could not read the
  process table, and existing HTTP tests could not bind localhost listeners.
  The same incremental target rerun as
  `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy make verify-incremental`
  with host permissions exited 0: formatting, vet, repository tests, skill
  checks and build passed.

The Daemon's pre-existing `status: in_progress` is preserved. The declared
Verification commands were not run; settlement remains Daemon-owned. No
commit, push or Pull Request was made. Documentation and skill changes remain
in task_04's slice.

## Carry-forward provenance

- Source Run: `run_20261005T101255Z_f599a5cc8ea7f4c2`
- Source commit: `7285be89765cc528894f49d08aae02ae3f04fd2b`
