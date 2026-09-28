---
task: task_04
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: pending
type: backend
complexity: high
---

# Task 04: A blocking review after archive parks publication

## Overview

An archived Spec stays byte-identical, so a finding on a candidate that
archives a Spec cannot be corrected by a Task inside that Spec. Today nothing
says so. `roundfix review` records no archived Spec, the Delivery Queue parks
such a verdict as `review-findings`, the blocker that invites a corrective
Task, and `roundfix deliver retry` answers a moved archived head with a
generic refusal. This Task records the Specs a candidate archives. A findings
verdict on such a candidate parks its item as `corrective-spec-required:
<slug>`, and the retry refuses the item once its head moved, because the
correction is a new Spec with its own authorization and QA gate. The blocker is
written by the Delivery Queue owner and read by the operator through `deliver
status` and `deliver retry`. No queue grant, Run budget or corrective-Task
ceiling may turn it into permission to author that Spec.

## Requirements

1. MUST make `reviewCandidateSpecContexts` in `internal/cli/review.go` also
   return, in a new `reviewSpecContextResult` field, the sorted slugs of
   changed Spec folders under any root other than the first root that
   `reviewCandidateSpecRoots` returns. That root is the active Spec Root, and
   each slug counts whether its context is carried or skipped. MUST keep every
   function signature unchanged.
2. MUST add the record field `archivedSpecs` (JSON `archivedSpecs`, always
   present, `[]` on paths that compute no Spec context), set from that result
   on a fresh review and kept from the record on a reused one.
3. MUST make `finishReviewCommand` write, for a `findings` record whose
   `archivedSpecs` is not empty, one stderr line that names those slugs. The
   line says that an archived Spec is never corrected in place and contains
   `author a corrective Spec with its own authorization and QA gate`. Exit
   codes stay as they are.
4. MUST add `ArchivedSpecs []string` to `delivery.ReviewResult` in
   `internal/delivery/engine.go`, filled by `deliveryReviewResult` in
   `internal/cli/deliver_workflow.go`.
5. MUST add `BlockerCorrectiveSpecRequired = "corrective-spec-required"`. MUST
   make `reviewCandidate` park a `ReviewOutcomeFindings` result with archived
   Specs as `corrective-spec-required: <slug>, <slug>` (sorted, joined by `, `),
   and keep `review-findings` for findings with none.
6. MUST make `Engine.Retry` handle a blocker with the
   `corrective-spec-required` prefix before its archived and carry-forward
   branches. After `UseItemBranch` and `InspectItem`:
   - When the item head differs from the parked candidate head, return an error
     naming the archived Specs, both heads and the phrase `author a corrective
     Spec with its own authorization and QA gate`, leaving the stored item
     unchanged, so `roundfix deliver retry` exits `2` and starts no owner.
   - When the head is unchanged, return the item to `reviewing` without
     carry-forward.
7. MUST describe `archivedSpecs` and the stderr line in the review section, and
   the blocker and the retry row in the delivery section, of
   `docs/user-guide/commands.md` and `.agents/skills/roundfix/SKILL.md`, using
   the phrase `corrective-spec-required` in both. Regenerate
   `skills/roundfix/SKILL.md` with `make skills-sync`. MUST add a Corrective
   Spec entry to `CONTEXT.md` naming it the new Spec, with its own
   authorization and QA gate, that corrects a finding on an archived Spec.
8. MUST add `internal/cli/review_archived_spec_test.go` with:
   - `TestReviewRecordsTheSpecsACandidateArchives`
   - `TestReviewRecordsNoArchivedSpecForAnActiveSpec`
   - `TestReviewNamesTheCorrectiveSpecForFindingsAfterArchive`
   - `TestReviewAfterArchiveWithoutFindingsNamesNoCorrectiveSpec`
   - `TestDeliveryReviewResultCarriesArchivedSpecs`
9. MUST add `internal/delivery/corrective_spec_test.go` with:
   - `TestDeliveryEngineParksFindingsAfterArchiveAsCorrectiveSpecRequired`
   - `TestDeliveryEngineParksFindingsWithoutArchiveAsReviewFindings`
   - `TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved`, which asserts the
     stored item is unchanged
   - `TestRetryReturnsAnUnchangedCorrectiveSpecItemToReviewing`, which asserts
     no carry-forward call
10. MUST keep `TestDeliveryEngineParksDeclaredBlockersAndContinues`,
    `TestRetryRefusesAnArchivedItemWhoseHeadMoved` and
    `TestReviewReadsAnArchivedSpecFromTheCandidate` green and unchanged.

## Subtasks

- [ ] Record the archived Specs and name the corrective Spec on findings.
- [ ] Park and refuse in the Delivery Queue.
- [ ] Document the blocker, the retry row and the glossary term.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A review record names every Spec its candidate archives.
- [ ] A findings verdict on such a candidate parks as
      `corrective-spec-required: <slug>`, and one without archived Specs still
      parks as `review-findings`.
- [ ] A retry refuses a moved head and returns an unchanged one to
      `reviewing`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/delivery/engine.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`
- creates: `internal/cli/review_archived_spec_test.go`
- creates: `internal/delivery/corrective_spec_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewRecordsTheSpecsACandidateArchives|TestReviewRecordsNoArchivedSpecForAnActiveSpec|TestReviewNamesTheCorrectiveSpecForFindingsAfterArchive|TestReviewAfterArchiveWithoutFindingsNamesNoCorrectiveSpec|TestDeliveryReviewResultCarriesArchivedSpecs|TestDeliveryEngineParksFindingsAfterArchiveAsCorrectiveSpecRequired|TestDeliveryEngineParksFindingsWithoutArchiveAsReviewFindings|TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved|TestRetryReturnsAnUnchangedCorrectiveSpecItemToReviewing|TestDeliveryEngineParksDeclaredBlockersAndContinues|TestRetryRefusesAnArchivedItemWhoseHeadMoved|TestReviewReadsAnArchivedSpecFromTheCandidate)$" ./internal/cli ./internal/delivery 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewRecordsTheSpecsACandidateArchives TestReviewRecordsNoArchivedSpecForAnActiveSpec TestReviewNamesTheCorrectiveSpecForFindingsAfterArchive TestReviewAfterArchiveWithoutFindingsNamesNoCorrectiveSpec TestDeliveryReviewResultCarriesArchivedSpecs TestDeliveryEngineParksFindingsAfterArchiveAsCorrectiveSpecRequired TestDeliveryEngineParksFindingsWithoutArchiveAsReviewFindings TestRetryRefusesACorrectiveSpecItemWhoseHeadMoved TestRetryReturnsAnUnchangedCorrectiveSpecItemToReviewing TestDeliveryEngineParksDeclaredBlockersAndContinues TestRetryRefusesAnArchivedItemWhoseHeadMoved TestReviewReadsAnArchivedSpecFromTheCandidate; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "corrective-spec-required" && tr -s '[:space:]' ' ' < .agents/skills/roundfix/SKILL.md | grep -qF -- "corrective-spec-required" && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Corrective Spec**" && diff -r .agents/skills/roundfix skills/roundfix >/dev/null` — expected: exit 0; before this Task none of the nine new named tests exists and no guide names `corrective-spec-required`, so the command fails.

## References

- [_techspec.md](_techspec.md) — A blocking review after archive
- `_prd.md` → Goal 4; Core Feature 4; Success Metric 4
- `_techspec.md` → API Contracts 1 and 5; Testing Approach 4
