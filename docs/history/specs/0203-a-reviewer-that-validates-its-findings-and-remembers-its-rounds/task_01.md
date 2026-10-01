---
task: task_01
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: completed
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

## Result

Implemented this Task's anchor-validation slice. The prompt retains its
existing verdict grammar and carries the exact finding-grammar sentence and
C1–C4 Delivery Conventions. Each finding records its parsed anchor and
validation. The validator consumes the exact candidate diff used in the
prompt, without fetching a second diff. Anchored findings stand;
missing/out-of-diff anchors are dismissed as `unanchored`. The outcome,
reuse, missing-disposition diagnostics and dispose refusal now respect
validation dismissals. No convention judgment, lineage or session changes
belong to this diff; those remain with task_02–task_04.

The three new Go files are absent from committed HEAD (confirmed with
`rtk proxy git -c core.fsmonitor=false ls-tree --name-only HEAD --
internal/cli/review_conventions.go internal/cli/review_validation.go
internal/cli/review_validation_test.go`, which returned no paths). The
Daemon-provided `status: in_progress` was preserved. No Task Graph or other
Task file was edited, and no commit, push or Pull Request was made.

Focused implementation evidence:

- `GOCACHE=/tmp/roundfix-task01-cache rtk proxy go test ./internal/cli -run
  '(Review|SplitReview|DeliveryReview)' -count=1`: exit 0, final run
  `ok roundfix/internal/cli 15.511s`. This includes all existing review tests,
  the eight new acceptance tests and the additional parser/reuse negatives.
- `GOCACHE=/tmp/roundfix-task01-cache rtk make verify-incremental`: final
  frozen-tree run with required process-table access exited 0. Formatting,
  vet, repository tests, skill sync/checks and CLI build passed. The prior
  sandboxed attempt hit stop-command process-table restrictions and detected
  edits made while its suite guard was running. The first frozen rerun
  passed the CLI package but hit the unrelated Daemon test
  `TestTaskBudgetReasonNamesTheSettlementThatRenewedIt` (its 200 ms budget
  expired before the second Task began). That unchanged test passed in an
  isolated run, and the subsequent incremental run passed with the successful
  package caches. No unrelated test, budget or verification configuration
  was changed.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.
- A Python read/assert inspection of the review command guide confirmed
  `roundfix/delivery-conventions/v1`, `dismissed-by-validation`, `Failure:`
  and `was dismissed by validation` are present.

Acceptance evidence:

| Criterion | Implementation and focused evidence |
| --- | --- |
| Anchor forms and refusal cases | `TestReviewFindingAnchorParsesEachForm` covers bare/backticked, ranged, comma-listed, end/colon/parenthesis boundaries and Unicode whitespace; missing lines, trailing letters, a non-leading anchor, reversed ranges, whitespace inside paths and numeric overflow are refused. |
| Hunk lines and whole files | `TestReviewDiffIndexHoldsHunkLinesAndWholeFiles` covers context/added lines, multiple hunks, out-of-hunk lines, overlapping ranges, pure deletions, deleted/binary/rename-only/mode-only files and absent files. Separate negative tests prove hunk content cannot overwrite file headers, and Git-quoted paths, including a trailing backslash, retain their identity. |
| Prompt grammar and conventions | `TestReviewPromptCarriesTheDeliveryConventionsAndTheFindingGrammar` checks the added sentence, the existing sentence byte for byte, the versioned block, the Cn lines and placement after the head. `TestReviewPromptAsksForOneListItemPerFinding` and `TestReviewPromptWithoutSpecIsUnchanged` remain unedited and pass. |
| Mixed dismissed/standing findings | `TestReviewDismissesAnUnanchoredFindingAndKeepsTheAnchoredOne` reproduces Transcript 2 without lineage: exit 1, findings outcome, exact stderr, both anchors/validation reasons and persisted validation metadata. A separate missing-anchor case checks its distinct dismissal reason. |
| All findings outside the diff | `TestReviewFindingsAllDismissedByValidationExitZero` proves exit 0, findings-dismissed, two dismissal diagnostics, no ledger and no dispositions; it also proves a findings record with no standing finding is refused. |
| No readable anchors | `TestReviewBlocksFindingsThatNameNoFileAndLine` reproduces Transcript 3 without lineage: exit 2, exact reason/stderr, no finding items/text, retained answer path and validator not-needed. |
| Same-head reuse | `TestReviewReuseCountsValidationDismissalsAsDismissed` proves a validation dismissal plus one evidence-backed operator dismissal reuses as findings-dismissed without Agent activity. `TestReviewReuseListsOnlyStandingFindingsAsMissing` separately proves only the standing F2 gets a missing-disposition diagnostic. |
| Dispose refusal | `TestReviewDisposeRefusesAFindingDismissedByValidation` proves exit 2, exact refusal text, empty stdout and byte-identical existing ledger. |
| Existing contracts | The focused review selection passes all existing review tests, including `TestReviewRecordRoundTripsEachOutcome`, `TestReviewRecordRefusesInvalidFindingsDismissed` and `TestReviewClassifiesVerdictVariants`. Legacy records retain their field names/order, optional new fields and standing findings. |
| Command guide | The existing guide describes C1–C4 and their version, anchor grammar/indexing, outcomes, validation status, reuse and dispose refusal; the phrase inspection above confirms the required public terms. |

Existing tests updated only to supply a finding anchored in their fixture diff,
with every top-level name and existing behavioral assertion retained:

- `TestReviewCommandExitsOneAndRecordsFindings`
- `TestReviewClassifiesVerdictVariants`
- `TestReviewRecognisesAColonlessFindingsHeader`
- `TestReviewRecognisesEmphasizedFindingsHeader`
- `TestReviewReadsFindingsNoneAsNoFindings`
- `TestReviewIgnoresAQuotedVerdictInsideFindings`
- `TestReviewKeepsTheRawAnswer`
- `TestReviewReportsFindingsDismissedWithoutAskingTheReviewer`
- `TestReviewKeepsStandingFindingsWithoutAskingTheReviewer`
- `TestReviewIgnoresADismissalOfDifferentText`
- `TestReviewIgnoresAFixWhenClearingAHead`
- `TestReviewAsksAgainAfterTheHeadMoves`
- `TestReviewAsksAgainForADifferentBase`
- `TestReviewAsksAgainForADifferentProvider`
- `TestReviewRecordsTheSpecsACandidateArchives`
- `TestReviewNamesTheCorrectiveSpecForFindingsAfterArchive`
- `TestReviewInOneCheckoutLeavesAnotherCheckoutsRecord`
- `TestReviewDisposeReadsOnlyItsCheckoutsRecord`
- `TestReviewClassifiesTheFinalMessageAfterProgressText`

The quoted-verdict fixture keeps the quoted clean verdict within an anchored
finding; the inline-none fixture marks its following anchored finding as a
separate item. Both retain their previous verdict and text-preservation
assertions. Disposition-only and merge-base fixtures that this change did not
invalidate were left unedited.

Declared Verification commands were not run. Task status and settlement remain
Daemon-owned. The Roundfix skill reference/version change remains task_04.

## Carry-forward provenance

- Source Run: `run_20261001T145917Z_95f0b38476c6ebcc`
- Source commit: `b645d2999344e5decbf55f98a2d105e6e7e99a26`
