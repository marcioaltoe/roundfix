---
task: task_01
spec: 0278-a-fixed-round-2-finding-gets-one-correction-review-and-the-plan-names-only-what-it-delivers
status: pending
type: backend
complexity: high
---

# Task 01: A fixed round-2 finding gets one Correction Review, and the retry admits its fix

## Overview

On 2026-10-17 round 2 of Spec 0276's pre-PR review found two real defects in
the operator's CI fix. The fix that replaced them could reach `main` only by a
hand review and a hand merge (`_prd.md` → Sources 1). Today the next review
after round 2 closes `ceiling-closed` without reading the fix, and the
Delivery Retry refuses a fix outside the archived Spec's declared paths.

This Task adds the Correction Review: one review of the delta from the round-2
head, opened only when every standing round-2 finding was fixed. It also lets
the retry of such a park return the item to review whatever paths the fix
touched. It is verifiable on its own through review and retry tests that use a
fake reviewer.

## Requirements

1. MUST answer the operator observation of 2026-10-17 (Spec 0276, pull
   request #568) cited in `_prd.md` → Sources 1 by implementing API Contracts
   1 to 4 of `_techspec.md` (ADR-0286).
2. MUST add `Correction bool` with the JSON name `correction` (omitted when
   false) to `reviewLineage`. In the ceiling branch of the review command,
   before `closeAtCeiling`, it MUST plan a Correction Review exactly when API
   Contract 1 holds:
   - the prior record has outcome `findings` and a round-2 lineage without
     `correction`;
   - the head moved and descends from the prior head;
   - every standing finding has exactly one valid disposition at the prior
     head, and that disposition is `fixed` by a commit between the prior head
     and the current head.

   Any other case MUST reach `closeAtCeiling` as today.
3. MUST run a Correction Review through the round-2 path:
   - the delta from the prior head, never segmented, blocked above
     `reviewDiffBound`, in a new Agent Session;
   - the prompt of `buildCorrectionReviewPrompt` with the exact sentences and
     delimiters of API Contract 2;
   - a record whose lineage is `{round: 2, correction: true, previousHead,
     previousFindings, previousDispositions}`.

   The data that reaches the reviewer is the delta and the round-2 findings
   with their fixing commits, sent to the configured review provider exactly
   as round 2 sends its delta.
4. MUST follow API Contract 3:
   - A clean or findings-dismissed Correction Review advances the Delivery
     Queue as a clean round 2 does.
   - Findings park as round-2 findings do.
   - A blocked Correction Review repeats as a Correction Review.
   - A review after a record with `correction: true` reaches the ceiling and
     never plans a second Correction Review.
5. MUST make `ProveReviewCorrection` skip the declared-path bound only when
   the parked record's lineage is round 2 without `correction` and every
   standing finding's one valid disposition is `fixed` (API Contract 4). It
   MUST keep descent, the record check, the one-valid-disposition rule and
   every existing refusal text for every other park.
6. MUST add these tests, each calling `t.Parallel()` first, compiling for
   Windows and using the review and retry fixtures already used beside them,
   with a fake reviewer that counts its calls. A test MUST NOT write a literal
   executable: a fake runtime binary comes from `testfixture.FixtureBinary` or
   the review fixtures already in the package.

   The tests are:
   - in the new `internal/cli/review_correction_test.go`:
     `TestCorrectionReviewReviewsTheFixOnce`,
     `TestCorrectionReviewFindingsReachTheCeiling` and
     `TestADismissedRoundTwoFindingKeepsTheCeiling`;
   - in `internal/cli/deliver_review_correction_test.go`:
     `TestRetryAdmitsAFixedRoundTwoParkOutsideDeclaredPaths` and
     `TestRetryKeepsTheBoundForADismissedRoundTwoFinding`.

   Every existing review lineage, ceiling and correction test MUST keep its
   name and pass. A test MUST NOT rename or remove an existing test.
7. MUST document the Correction Review in these files, each naming
   `Correction Review`:
   - `docs/user-guide/commands/review.md`, where its lineage section and the
     paragraph that says a third review never calls the reviewer gain the one
     exception;
   - `docs/user-guide/commands/deliver.md`, where the moved-head correction
     paragraph gains the round-2 exemption;
   - `.agents/skills/roundfix/references/review.md` and
     `.agents/skills/roundfix/references/deliver.md`.

   MUST search the guides and skill copies for text saying that the change
   after round 2 is never reviewed, or that a retry always keeps the
   declared-path bound, and qualify each. No text inside any
   `### QA settlement` section may change.
8. MUST add **Correction Review** to `CONTEXT.md` and revise **Reviewer
   Lineage** through domain-modeling, with ADR-0286's link. The revised
   **Reviewer Lineage** keeps "third round" under _Avoid_.
9. MUST run these commands, because this Task edits user guide pages and the
   Roundfix references. No record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the lineage field, the applicability check and the Correction Review path.
- [ ] Add the retry exemption.
- [ ] Add the five tests.
- [ ] Document the Correction Review in the guides, the references and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] After round-2 findings all fixed, exactly one reviewer call reviews the delta from the round-2 head, and its record carries `correction: true`.
- [ ] A Correction Review never repeats, and a dismissal keeps today's ceiling.
- [ ] The retry admits a fixed round-2 park outside the declared paths and refuses the dismissed case with today's text.

## Consumers

- The Pre-PR Review Record's lineage and the review correction proof — consumers: `./internal/cli`, `./internal/delivery`

## Context

- instruction: `docs/adr/0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md`
- instruction: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- instruction: `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- interface: `internal/cli/review_lineage.go`
- interface: `internal/cli/review.go`
- interface: `internal/cli/deliver_review_correction.go`
- interface: `internal/cli/deliver_review_correction_test.go`
- creates: `internal/cli/review_correction_test.go`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestCorrectionReviewReviewsTheFixOnce|TestCorrectionReviewFindingsReachTheCeiling|TestADismissedRoundTwoFindingKeepsTheCeiling|TestRetryAdmitsAFixedRoundTwoParkOutsideDeclaredPaths|TestRetryKeepsTheBoundForADismissedRoundTwoFinding)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCorrectionReviewReviewsTheFixOnce TestCorrectionReviewFindingsReachTheCeiling TestADismissedRoundTwoFindingKeepsTheCeiling TestRetryAdmitsAFixedRoundTwoParkOutsideDeclaredPaths TestRetryKeepsTheBoundForADismissedRoundTwoFinding; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/review.md docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/review.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Correction Review' || { printf 'missing Correction Review in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Correction Review**:' || { printf 'missing glossary term Correction Review\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Reviewer Lineage**:' || { printf 'missing glossary term Reviewer Lineage\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'Review|Retry|Correction|Ceiling' ./internal/cli || exit 1; go test -count=1 ./internal/delivery || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the five tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the review and retry tests, the documentation, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 1; Core Feature 1; Core Feature 2; Core Feature 5; Success Metric 1; Success Metric 2; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 7; Invariants 1 and 2; Current behavior; Testing Approach; Build Order 1
- ADR-0286; ADR-0197; ADR-0272; ADR-0280; ADR-0189; ADR-0256; ADR-0259
