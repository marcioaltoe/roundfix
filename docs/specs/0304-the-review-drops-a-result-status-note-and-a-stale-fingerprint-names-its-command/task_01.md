---
task: task_01
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
status: pending
type: backend
complexity: high
---

# Task 01: The pre-PR review dismisses a Result Status Note of a completed Task without the validator

## Overview

The pre-PR review sends the finding "the Result says `in_progress`, the front
matter says `completed`" to the sealed validator, which kept it standing on
Spec 0290 and on Spec 0297 (2026-10-22) by reading Delivery Convention C2 as a
requirement; the operator dismissed both with evidence. This Task rewords C2
as version 4, dismisses that pair by validation in code before the validator
is asked, keeps every other finding on today's path, and carries the change
into the review guide, both copies of the Roundfix review reference, the
skill's version, the Behavior Surface Record and the glossary (`_prd.md` →
Core Features 1-3).

## Requirements

1. MUST answer the adopted Backlog Entry
   [`references/2026-10-22-the-pre-pr-review-parks-on-a-result-status-note-that-convention-c2-allows.md`](references/2026-10-22-the-pre-pr-review-parks-on-a-result-status-note-that-convention-c2-allows.md)
   of 2026-10-22, whose findings come from the pre-PR reviews of Spec 0290
   and of Spec 0297 on 2026-10-22.
2. MUST set `deliveryConventionsVersion` in
   `internal/cli/review_conventions.go` to `roundfix/delivery-conventions/v4`
   and C2's text to exactly the text of `_techspec.md` → API Contract 1,
   leaving C1, C3, C4 and C5 byte-identical.
3. MUST add `reviewResultStatusNote` and the constant
   `reviewResultStatusNoteReason` to
   `internal/cli/review_convention_validator.go` with the signature and the
   reason text of `_techspec.md` → Interfaces and the conditions of API
   Contract 2, and call it in `validateReviewConventions` as API Contract 3
   says: a finding it accepts is recorded `dismissed-by-validation` with rule
   `convention:C2` and that reason and is not asked; when no finding is left
   to ask, the record settles exactly as after a validator run, so a review
   whose findings are all dismissed exits 0 with outcome
   `findings-dismissed` and `validator` `not-needed`. Nothing else in the
   validator, its prompt, its grammar or its fail-closed handling changes.
4. MUST create `internal/cli/review_result_status_note_test.go` with
   `TestReviewDismissesAResultStatusNoteOfACompletedTaskWithoutTheValidator`
   and `TestReviewKeepsAContradictionThatIsNotAResultStatusNote` and every
   subtest `_techspec.md` → Testing Approach names, each through `roundfix
   review` with the existing fake reviewer and fake sealed validator of
   `review_convention_validator_test.go`. Each top-level test and subtest
   calls `t.Parallel()` first. The `0297 shape` fixture's Result ends with the
   paragraph "The declared Task Verification and repository Verification were
   not run;" / "they remain Daemon-owned. Task status remains
   `in_progress`, and no commit," / "push or pull request was made." and its
   anchor is that paragraph's first line; the `0290 shape` fixture's Result
   opens with "Implemented the documentation slice while leaving the
   Daemon-owned status" / "`in_progress` unchanged:" and its anchor is the
   blank line right after `## Result`. Each dismissal subtest asserts the
   whole stderr line `roundfix: review finding F1 dismissed by validation
   (convention:C2): <reason>`.
5. MUST make these declared breaks and no other change to existing tests:
   `TestDeliveryConventionC2ExemptsTheQATask`
   (`internal/cli/review_convention_qa_task_test.go`) expects the v4 version
   and the new C2 text; in `TestReviewValidatorDismissesADaemonSettlementAsC2`
   (`internal/cli/review_convention_validator_test.go`) the finding text
   becomes "The Task is `completed` while its Result gives no evidence.
   Failure: an invalid Task state transition." because the old text states
   the pair this Task dismisses without a question; every other assertion of
   that test keeps its bytes.
6. MUST change, in `docs/user-guide/commands/review.md` and in the canonical
   `.agents/skills/roundfix/references/review.md`:
   - the version `roundfix/delivery-conventions/v3` to
     `roundfix/delivery-conventions/v4`, and the C2 bullet to the new C2 text,
     wrapped as the other bullets are;
   - after the paragraph that ends "with the dismissal's rule and reason.",
     add this paragraph: "A Result Status Note is a line of a non-QA Task's
     `## Result` that contains `in_progress` or `Daemon-owned` in any letter
     case; the Agent writes it while the Daemon holds the Task `in_progress`.
     Before the validator is asked, Roundfix dismisses a finding with rule
     `convention:C2` when, at the review head, the whole anchor lies in a Task
     file's C2 region, the Task's front matter has one `status: completed`
     line and no `type: qa`, its `## Result` holds a Result Status Note, and
     the finding's text contains `completed` and one of `in_progress`, `in
     progress` or `Daemon-owned` in any letter case. Its reason is
     `<reasonText>`. A review whose findings are all dismissed this way reads
     `findings-dismissed`, and its `validator` stays `not-needed`.", with
     `<reasonText>` replaced by the value of `reviewResultStatusNoteReason`;
   - replace "Only anchored findings inside a convention region or lacking a
     failure clause are sent to the validator." with "Only anchored findings
     inside a convention region or lacking a failure clause, and not
     dismissed as a Result Status Note, are sent to the validator."
7. MUST add **Result Status Note** to `CONTEXT.md` through `domain-modeling`,
   right after the **Delivery Convention** entry, as "**Result Status
   Note**:" followed on the next line by "A line of a non-QA Task's `##
   Result` that reports the Task's status as `in_progress` or Daemon-owned,
   which the Agent writes while the Daemon holds the Task `in_progress`; it
   does not contradict the `completed` status the Daemon writes after the
   Task's Verification passes, and the Pre-PR Review Command dismisses a
   finding that names that pair without the validator (ADR-0317)." and the
   line "_Avoid_: Status contradiction, stale status, Result status claim".
8. MUST revise **Delivery Convention** in `CONTEXT.md` through
   `domain-modeling` by appending to its definition line, after the period
   that ends it, the sentence "A finding that names a Result Status Note
   against a `completed` Task is dismissed by validation without the sealed
   validator (ADR-0317)." Its `_Avoid_` line keeps its bytes. In both
   entries the ADR reference takes the glossary's link form, a link to
   `docs/adr/0317-a-result-status-note-of-a-completed-task-is-dismissed-without-the-validator.md`
   as the other entries write theirs.
9. MUST run these commands, because this Task edits an owned skill and a
   guide. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `SKILL.md` of the Roundfix skill by
   one patch over the value the tree holds when this Task starts, because the
   record command refuses changed content under an unchanged version. The
   Agent runs them itself before handing back, in this order, and runs the
   last one again after any later edit of the guide or the reference
   ("re-record after the last guide edit"), because the Daemon never
   regenerates a record; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     review guide's fingerprint changes.
10. MUST search the guides, the agent guides and the skill copies for wording
    that says the Daemon writes a Task's `## Result`, or that only the sealed
    validator can dismiss an anchored finding, and report what it found.
11. MUST NOT change the Agent prompt (`internal/agent/spec_prompt.go`), the
    `implement-task` skill, any ADR, `docs/references/skill-coverage.json`,
    `.roundfixrc.yml`, any other guide page, or any text inside a `### QA
    settlement` section.

## Subtasks

- [ ] Reword C2 as version 4 and update its test.
- [ ] Add the Result Status Note dismissal and its two tests; adapt the declared break.
- [ ] Update the review guide and the canonical reference, add the glossary entries.
- [ ] Raise the skill version, sync the mirrors, record the version and re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] A finding of the 0297 shape, the 0290 shape or a front-matter anchor naming the pair is dismissed with rule `convention:C2` and the fixed reason, without a sealed request, and the review exits 0 with `findings-dismissed`.
- [ ] A failed Task, a QA Task, a Result without a note, a note outside the Result and a finding naming another failure each reach the validator once and stand on its verdict; an anchor outside the C2 region stands without a question.
- [ ] The prompt, the validator input, the guide and both reference copies carry version 4 and the same C2 text.
- [ ] `CONTEXT.md` defines **Result Status Note** and revises **Delivery Convention**; the skill version is recorded, the mirrors are in sync and the Behavior Surface Record is current.

## Consumers

- Delivery Conventions version 4 and C2 text — consumers: `./internal/cli`, `./internal/docscontract`
- review validation outcome of a Result Status Note finding — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0317-a-result-status-note-of-a-completed-task-is-dismissed-without-the-validator.md`
- instruction: `internal/docscontract/command_documentation_test.go`
- interface: `internal/cli/review_conventions.go`
- interface: `internal/cli/review_convention_validator.go`
- creates: `internal/cli/review_result_status_note_test.go`
- interface: `internal/cli/review_convention_validator_test.go`
- interface: `internal/cli/review_convention_qa_task_test.go`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `grep -qF 'roundfix/delivery-conventions/v4' internal/cli/review_conventions.go || { echo 'conventions are not version 4' >&2; exit 1; }; for name in TestReviewDismissesAResultStatusNoteOfACompletedTaskWithoutTheValidator TestReviewKeepsAContradictionThatIsNotAResultStatusNote; do go test -list "^${name}\$" ./internal/cli | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^(TestReviewDismissesAResultStatusNoteOfACompletedTaskWithoutTheValidator|TestReviewKeepsAContradictionThatIsNotAResultStatusNote|TestReviewValidator.*|TestReviewConvention.*|TestDeliveryConventionC2ExemptsTheQATask|TestReviewDisposeRefusesAConventionDismissal)$' ./internal/cli` — expected: exit 0. Before this Task the conventions are version 3 and neither new test exists, so the command fails at its first check. Afterwards the conventions are version 4, both tests exist, and they pass with the validator, convention and dispose tests of the package.
- `for f in docs/user-guide/commands/review.md .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'roundfix/delivery-conventions/v4' || { printf 'missing version 4 in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A Result Status Note is a line of a non-QA Task' || { printf 'missing dismissal in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'and not dismissed as a Result Status Note, are sent to the validator' || { printf 'missing validator sentence in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Result Status Note**:' || { echo 'missing glossary term' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A finding that names a Result Status Note against a' || { echo 'missing Delivery Convention revision' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract -run '^(TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem|TestTheDeliveryConventionContractReportsAStaleConvention|TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract` — expected: exit 0. Before this Task no document names version 4 or a Result Status Note, so the command fails at its first phrase check. Afterwards the guide and both reference copies carry version 4, the dismissal and the validator sentence, the glossary defines **Result Status Note** and revises **Delivery Convention**, the mirrors are in sync, the skill version is recorded, and the documentation contracts that hold the conventions and the Skill Coverage Map pass.

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metrics 1-4; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4 and 7; Invariants 1-5; Testing Approach; Build Order 1
- ADR-0317
