---
task: task_03
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
status: completed
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


## Result

Implemented the C5 slice for the 2026-10-02 Backlog Entry, "The review flags
an authorized QA override archive". Task status and declared Verification
remain Daemon-owned; no declared Verification command was run.

### Implementation and acceptance evidence

1. Archived override eligibility: `conventionRegions` adds C5 for any file
   under an archived Spec only when its `_prd.md` at the immutable review
   head records `qa_override: true` with non-blank approval and reason.
   `TestOverrideArchiveIsEligibleForConventionC5` covers the built-in history
   archive and a custom `_archived` root, PRD, failed Task and partial QA
   Report anchors, preservation of C1/C2/C3, and checkout edits that remove
   authority without changing the reviewed head.
2. Ineligible records: `TestArchiveWithoutOverrideIsNotEligibleForC5` covers
   missing PRDs, absent or false override markers, missing or whitespace-only
   approval/reason, body-only metadata, unclosed front matter, active Specs,
   and uncommitted grants. `TestConventionC5MalformedRecordFailsClosed`
   covers malformed YAML and duplicate override keys.
3. Prompt and record contract: `deliveryConventionsVersion` is now
   `roundfix/delivery-conventions/v2`, and C5 carries the TechSpec's exact
   text. `TestReviewPromptCarriesConventionC5` checks the prompt's C5 text
   and the shared version in both prompt and record validation.
   `TestConventionC5ValidatorHonorsEligibility` checks an eligible sealed
   dismissal and fail-closed rejection of an ineligible C5 rule or blank
   reason. Every new test uses a disposable Git repository; no provider is
   called. Existing validation tests and validator implementation are unedited.
4. Documentation and Skill: the review guide and canonical review reference
   describe v2, C5 and its eligibility, and extend the convention range
   through C5. Both Roundfix Skill version fields rose from the starting
   tree's `0.1.20` to `0.1.21`; `make skills-sync` regenerated the mirrors,
   and the version recorder added its digest. A Python assertion check
   confirmed guide/reference v2 and C5 text, both canonical/mirror byte
   comparisons and the recorded version. Baseline digest regeneration
   reported no derived changes.

### Focused checks

- Before implementation, `rtk proxy go test ./internal/cli -run
  'Test(OverrideArchive|ArchiveWithoutOverride|ReviewPromptCarriesConvention)'
  -count=1` exited 1 on missing C5 eligibility and prompt text, after
  correcting an initially mistyped outcome constant in the new test.
- After implementation, `rtk proxy go test ./internal/cli -run
  'Test(OverrideArchive|ArchiveWithoutOverride|ReviewPrompt|ReviewConvention|ReviewValidation|ReviewValidator|Convention)'
  -count=1` exited 0 (`internal/cli`, 10.720s), including the existing
  convention and fail-closed validator tests.
- `rtk make skills-sync` exited 0.
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` exited 0 and recorded `0.1.21`.
- `rtk make baseline-digests` exited 0 with `changed: false`.
- `rtk make verify-incremental` initially exited 2 because the sandbox
  denied process-table access in two existing force-stop tests and a local
  listener in `TestJevRouterGateChecksTheReportedKeyLimit`. The same command
  rerun with those permissions available exited 0: formatting, vet, the
  repository test suite, skill checks and build passed. The CLI suite took
  173.268s and the Daemon suite took 45.245s on that rerun.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The starting worktree already had the Daemon's `status: in_progress` change
in this Task file. All implementation changes stay within Task 03's declared
paths; no other Task or Task Graph was edited, and no workspace commit, push
or Pull Request was made. No follow-up implementation was added.

## Carry-forward provenance

- Source Run: `run_20261003T213431Z_8d253be688342055`
- Source commit: `0cb99efcb2e58c85e53067649fd3510493690468`
