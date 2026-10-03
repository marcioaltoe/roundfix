---
task: task_03
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
status: pending
type: backend
complexity: medium
---

# Task 03: The review knows an authorized QA Archive Override

## Overview

On 2026-10-02 the pre-PR review of Specs 0215 and 0218 parked each item as
`corrective-spec-required` on the finding that the Spec was archived with a
failed QA Task and a partial QA verdict, although each archive was made with
`roundfix archive --qa-override` under the maintainer's standing authorization
and its `_prd.md` recorded the override. This Task adds the fifth Delivery
Convention, so the review is told about an authorized override archive and the
validator may dismiss a finding that only restates it.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-02, "The review flags an
   authorized QA override archive", as `_techspec.md` → Convention C5 states.
2. MUST raise `deliveryConventionsVersion` to
   `roundfix/delivery-conventions/v2` and add convention C5 with the text of
   `_techspec.md` → Convention C5 (API Contract 5).
3. MUST make `conventionRegions` add `C5` for an anchor under an archive root
   only when that Spec's `_prd.md` at the review head has front matter with
   `qa_override: true` and a non-blank `qa_override_approval` and
   `qa_override_reason`; an anchor in an archived Spec without that record
   MUST NOT be eligible for C5, and C1 to C4 MUST keep their regions.
4. MUST keep the validator's fail-closed rules and the existing review
   validation tests passing unedited.
5. MUST add the tests named in Verification to the new file
   `internal/cli/review_override_convention_test.go`, each over a temporary
   repository and with no provider call.
6. MUST describe C5 and the new version in `docs/user-guide/commands/review.md`
   and the Roundfix Skill's `review` reference, changing the rule range from
   C1–C4 to C1–C5; MUST raise the Roundfix Skill's version by one patch level
   above the version on this Task's starting tree in both front-matter fields,
   run `make skills-sync`, and re-record the version.

## Subtasks

- [ ] Add C5 and raise the conventions version.
- [ ] Compute the C5 region from the archived override record.
- [ ] Add the tests.
- [ ] Describe it in the guide and the Skill, and raise its version.

## Acceptance Criteria

- [ ] A finding anchored on a file of an archived Spec whose `_prd.md` records
      the override lists `C5` among its eligible rules.
- [ ] A finding anchored in an archived Spec without the override record, or
      with a blank approval or reason, does not.
- [ ] The review prompt names `roundfix/delivery-conventions/v2` and the C5
      text, and the record's validation names the same version.
- [ ] The guide and the Skill describe C5, the mirrors equal their canonical
      files, and the raised version is recorded.

## Context

- creates: `internal/cli/review_override_convention_test.go`
- interface: `internal/cli/review_conventions.go`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/cli/review.go`
- instruction: `internal/cli/review_convention_validator.go`
- instruction: `internal/cli/review_validation_test.go`
- instruction: `docs/specs/0219-a-delivery-that-survives-archive-requeue-and-review/references/2026-10-02-the-review-flags-an-authorized-qa-override-archive.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestOverrideArchiveIsEligibleForConventionC5|TestArchiveWithoutOverrideIsNotEligibleForC5|TestReviewPromptCarriesConventionC5)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestOverrideArchiveIsEligibleForConventionC5 TestArchiveWithoutOverrideIsNotEligibleForC5 TestReviewPromptCarriesConventionC5; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `grep -qF -- "roundfix/delivery-conventions/v2" internal/cli/review_conventions.go || { printf 'conventions version not raised\n' >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestReviewPromptCarriesConventionC5|TestReviewValidation.*|TestConvention.*)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestReviewPromptCarriesConventionC5" || { printf 'missing pass\n' >&2; exit 1; }` — expected: exit 0; the existing validation and convention tests run beside the new one; before this Task the version is v1, so the command fails.
- `for file in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "roundfix/delivery-conventions/v2" || { printf 'missing version in %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "C5." || { printf 'missing C5 in %s\n' "$file" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && out="$(go test -count=1 -v -run "^TestEveryOwnedSkillVersionIsRecorded$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass\n' >&2; exit 1; }` — expected: exit 0; before this Task the guide and the reference name version v1 and no C5, so the command fails.

## References

- `_prd.md` → Goals; User Story 5; Core Features 6, 9; Success Metric 5; Acceptance evidence
- `_techspec.md` → Convention C5; API Contract 5; Testing Approach 3; Build Order 3
- ADR-0223; ADR-0196; ADR-0197; ADR-0154
