---
task: task_01
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: pending
type: backend
complexity: high
---

# Task 01: The prompt carries the Delivery Conventions and a finding parks only when its anchor is in the diff

## Overview

`roundfix review` parks a delivery on every finding the reviewer lists. The
reviewer is never told what a delivery writes by design, and nothing checks
that a finding names a line of the candidate. This Task adds the Delivery
Conventions and the finding grammar to the prompt. It reads each finding's
anchor and dismisses, with a recorded reason, a finding whose anchor the
candidate diff does not hold. It extends the record and the outcome rules so
that only a standing finding parks.

The data comes from the reviewer's final message and from the candidate diff
Roundfix already computes. Both are untrusted. The record is read by the
Delivery Queue from standard output, and by `review dispose` and reuse from
the checkout's record directory. The convention validator is task_02; this
Task leaves every anchored finding standing.

## Requirements

1. MUST add `internal/cli/review_conventions.go` with the four conventions
   `C1` to `C4`, the version `roundfix/delivery-conventions/v1` and the exact
   sentences in `_techspec.md` → The finding grammar and the prompt.
   `buildReviewPrompt` MUST keep its signature and every current line, and
   MUST add the grammar sentence and the conventions block as that section
   states. `TestReviewPromptAsksForOneListItemPerFinding` and
   `TestReviewPromptWithoutSpecIsUnchanged` MUST stay green without edits.
2. MUST add `internal/cli/review_validation.go` with
   `parseReviewFindingAnchor`, `indexReviewDiff`, `reviewDiffIndex.holds` and
   `reviewFindingHasFailureClause`, following `_techspec.md` → The finding
   grammar and the prompt, and Validation.
3. MUST extend the record as `_techspec.md` → Data Models states: `anchor` and
   `validation` on each finding item, and a top-level `validation`. Every
   existing field MUST keep its name, order and meaning. A record without the
   new fields MUST still read, with its findings counted as standing.
4. MUST validate every finding after classification in `runReviewCommand`.
   - A finding with no readable anchor, or with one the diff does not hold, is
     `dismissed-by-validation` with rule `unanchored` and the reason the
     TechSpec names.
   - Every other finding stands with reason `anchored in the candidate diff`.
   - The top-level `validation.validator` is `not-needed`.

   The anchor check MUST read the same diff text `reviewCandidateDiff` returns
   for the prompt.
5. MUST decide the outcome by `_techspec.md` → Validation (outcome table):
   `findings` with exit `1` when a finding stands; `findings-dismissed` with
   exit `0` when none stands; `blocked` with reason `findings name no file and
   line` and exit `2` when no finding carries a readable anchor. It MUST print
   one standard-error line per dismissed finding, in the form of Surface
   Transcript 2.
6. MUST make reuse at the same head count `dismissed-by-validation` as
   dismissed and list only standing findings as missing. `validateReviewRecord`
   MUST accept `findings-dismissed` when every finding is dismissed by
   validation or by exactly one evidence-backed `dismissed` disposition. It
   MUST refuse a `findings` record with no standing finding.
7. MUST make `roundfix review dispose` refuse a finding dismissed by
   validation, with `finding "<ID>" was dismissed by validation (<rule>)` and
   exit `2`, writing nothing to the ledger.
8. MUST update only the existing tests this change invalidates, keeping their
   names. These are the tests whose fake reviewer answers cite a path or line
   outside their fixture diff, in `internal/cli/review_test.go`,
   `internal/cli/review_disposition_test.go`,
   `internal/cli/review_head_bound_test.go`,
   `internal/cli/review_merge_base_test.go`,
   `internal/cli/review_record_checkout_test.go`,
   `internal/cli/review_archived_spec_test.go` and
   `internal/cli/review_final_message_test.go`. Each MUST cite a line its
   fixture diff holds, such as `review.txt:2`, and keep asserting what it
   asserted before. The Result MUST name each updated test. A top-level test
   MUST NOT be renamed or removed.
9. MUST put the new tests in `internal/cli/review_validation_test.go`, using
   `newReviewCommandFixture`, temporary repositories and fake runners, never a
   real reviewer.
10. MUST describe, in the `roundfix review` command guide
    `docs/user-guide/commands/review.md` that Spec 0194 creates, the Delivery
    Conventions with the version string `roundfix/delivery-conventions/v1`,
    the finding grammar, anchor validation, the status
    `dismissed-by-validation` and the `dispose` refusal. The Roundfix skill's
    review reference and version change in task_04.

## Subtasks

- [ ] Add the conventions and the grammar sentence to the prompt.
- [ ] Parse anchors and index the candidate diff.
- [ ] Validate each finding and record the result.
- [ ] Decide the outcome, standard error, reuse and `dispose` from standing findings.
- [ ] Update the existing tests whose answers cite lines outside their fixture.
- [ ] Describe the conventions and validation in the review command guide.
- [ ] Add a test for each acceptance criterion, with each negative case separate.

## Acceptance Criteria

- [ ] `TestReviewFindingAnchorParsesEachForm`: backticked, bare, ranged and
      comma-listed anchors parse to the first range. A missing line, a line
      followed by a letter, and an anchor that is not at the start are
      refused.
- [ ] `TestReviewDiffIndexHoldsHunkLinesAndWholeFiles`: context and added
      lines of a hunk are held, and a line outside every hunk is not. A pure
      deletion holds its `+c` line, and a deleted, binary or rename-only file
      holds every line.
- [ ] `TestReviewPromptCarriesTheDeliveryConventionsAndTheFindingGrammar`: the
      prompt holds the grammar sentence and the four `Cn.` lines with the
      version, and still holds today's grammar sentence byte for byte.
- [ ] `TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne`
      reproduces Surface Transcript 2 without its `lineage` field, which task_03 and task_04 add: outcome `findings`, exit `1`, and the
      standard-error line for F1.
- [ ] `TestReviewFindingsAllDismissedByValidationExitZero`: two findings
      anchored outside the diff record `findings-dismissed`, with exit `0` and
      no disposition in the ledger.
- [ ] `TestReviewBlocksFindingsThatNameNoFileAndLine` reproduces Surface
      Transcript 3 without its `lineage` field.
- [ ] `TestReviewReuseCountsValidationDismissalsAsDismissed`: a record with
      one finding dismissed by validation and one standing finding carrying an
      evidence-backed dismissal is reused as `findings-dismissed`, with exit
      `0` and no runner call.
- [ ] `TestReviewDisposeRefusesAFindingDismissedByValidation`: exit `2`, the
      refusal line, and ledger bytes unchanged.
- [ ] The existing tests named in Requirement 8 keep their names and pass, as
      do `TestReviewPromptAsksForOneListItemPerFinding`,
      `TestReviewPromptWithoutSpecIsUnchanged`,
      `TestReviewRecordRoundTripsEachOutcome`,
      `TestReviewRecordRefusesInvalidFindingsDismissed` and
      `TestReviewClassifiesVerdictVariants`.
- [ ] The review command guide carries `roundfix/delivery-conventions/v1`
      and `dismissed-by-validation`.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_conventions.go`
- creates: `internal/cli/review_validation.go`
- creates: `internal/cli/review_validation_test.go`
- interface: `internal/cli/review_test.go`
- interface: `internal/cli/review_disposition_test.go`
- interface: `internal/cli/review_head_bound_test.go`
- interface: `internal/cli/review_merge_base_test.go`
- interface: `internal/cli/review_record_checkout_test.go`
- interface: `internal/cli/review_archived_spec_test.go`
- interface: `internal/cli/review_final_message_test.go`
- creates: `docs/user-guide/commands/review.md`
- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewFindingAnchorParsesEachForm|TestReviewDiffIndexHoldsHunkLinesAndWholeFiles|TestReviewPromptCarriesTheDeliveryConventionsAndTheFindingGrammar|TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne|TestReviewFindingsAllDismissedByValidationExitZero|TestReviewBlocksFindingsThatNameNoFileAndLine|TestReviewReuseCountsValidationDismissalsAsDismissed|TestReviewDisposeRefusesAFindingDismissedByValidation|TestReviewPromptAsksForOneListItemPerFinding|TestReviewPromptWithoutSpecIsUnchanged|TestReviewRecordRoundTripsEachOutcome|TestReviewRecordRefusesInvalidFindingsDismissed|TestReviewClassifiesVerdictVariants)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewFindingAnchorParsesEachForm TestReviewDiffIndexHoldsHunkLinesAndWholeFiles TestReviewPromptCarriesTheDeliveryConventionsAndTheFindingGrammar TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne TestReviewFindingsAllDismissedByValidationExitZero TestReviewBlocksFindingsThatNameNoFileAndLine TestReviewReuseCountsValidationDismissalsAsDismissed TestReviewDisposeRefusesAFindingDismissedByValidation TestReviewPromptAsksForOneListItemPerFinding TestReviewPromptWithoutSpecIsUnchanged TestReviewRecordRoundTripsEachOutcome TestReviewRecordRefusesInvalidFindingsDismissed TestReviewClassifiesVerdictVariants; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the eight new tests do not exist, so the command fails.
- `test -f internal/cli/review_validation.go || { printf 'missing internal/cli/review_validation.go\n' >&2; exit 1; }; go test -count=1 -run '^TestReview' ./internal/cli` — expected: exit 0. Every existing `TestReview*` test passes with the updated fixtures. Before this Task, `internal/cli/review_validation.go` is absent, so the command fails.
- `file=docs/user-guide/commands/review.md; test -f "$file" || { printf 'missing file: %s\n' "$file" >&2; exit 1; }; for phrase in "roundfix/delivery-conventions/v1" "dismissed-by-validation"; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task the phrases are absent.

## References

- [_prd.md](_prd.md) — Goals 1–3; User Stories 1–3; Core Features 1, 2, 3 and 5; Success Metric 2; Declared breaks
- [_techspec.md](_techspec.md) — The finding grammar and the prompt; Validation; Data Models; Invariants 1, 3, 6 and 9; API Contracts 1–4; Surface Transcripts 2 and 3; Testing Approach 1 and 5; Build Order 1
- ADR-0196; ADR-0153; ADR-0169; ADR-0174
