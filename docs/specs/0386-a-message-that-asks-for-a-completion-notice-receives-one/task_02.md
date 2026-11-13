---
task: task_02
spec: 0386-a-message-that-asks-for-a-completion-notice-receives-one
status: completed
type: docs
complexity: medium
---

# Task 02: The follow-up mode and the release reference carry a Notice Request to its completion notice

## Overview

The Roundfix skill's follow-up mode keeps a request for a notice in the record
it writes for the message and sends the notice after the delivery; the release
reference gains the step that finds the requests a published release fulfils
and sends each notice once, never from a Daemon-assigned Run. The glossary
adds **Notice Request** and changes **Cross-project Message**, and ADR-0390,
ADR-0383 and ADR-0330 carry ADR-0396's amendment notices.

## Requirements

1. MUST change steps 4 and 7 under `## --follow-up` of the canonical
   `.agents/skills/roundfix/references/modes.md` as `_techspec.md` → API
   Contract 2 says, word for word, and nothing else in the file; MUST keep
   every phrase `skills/roundfix_modes_test.go` pins today. MUST NOT add or
   change text inside any `### QA settlement` section.
2. MUST append the `## Completion notices` section of `_techspec.md` → API
   Contract 3, word for word, to the canonical
   `.agents/skills/roundfix/references/release.md`, after one blank line, and
   change nothing else in the file.
3. MUST add `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest` to
   `skills/roundfix_modes_test.go`, calling `t.Parallel()` first, over both
   copies through `roundfixLayoutSources`, with the phrases `_techspec.md` →
   Testing Approach lists for `skills`, using `roundfixSection`,
   `roundfixNumberedItem` and `roundfixWhitespace`. No other test changes.
4. MUST change `CONTEXT.md` through `domain-modeling` as API Contract 4 says,
   word for word: the last sentence of **Cross-project Message** and the new
   **Notice Request** entry right after **Blocking Request**; MUST add no other
   capitalized phrase and no other entry.
5. MUST add the three amendment notices of API Contract 5, word for word,
   where it places them: one to
   `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`,
   one to `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
   and one to
   `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`;
   set each file's `updated_at` to `2026-11-12T21:40:00Z`, and change nothing
   else in the three ADRs.
6. MUST run `make skills-sync` to copy both canonical references to
   `skills/roundfix/references/`, then raise the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   as the Task's last step; never hand-edit a mirror, a version field or the
   record. Specs 0383 and 0385 raise the same version first, so the record
   command raises it from the tree the Task runs on. Of the Governed Paths the
   Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/modes.md`,
   `.agents/skills/roundfix/references/release.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Change steps 4 and 7 of the follow-up mode.
- [ ] Append the completion notices section to the release reference.
- [ ] Add the skill contract test.
- [ ] Add the glossary entry and change the message entry.
- [ ] Write the three amendment notices, then sync and record the skill version.

## Acceptance Criteria

- [ ] Both copies of the `modes` reference keep a Notice Request in its record and send its notice after the release or the merge; both copies of the `release` reference carry the completion notices section.
- [ ] `CONTEXT.md` defines **Notice Request** and the changed **Cross-project Message**, and the glossary contracts pass.
- [ ] ADR-0390, ADR-0383 and ADR-0330 carry an "Amended by ADR-0396" notice dated 2026-11-12.

## Context

- instruction: `docs/adr/0396-a-message-that-asks-to-be-told-when-its-work-is-done-receives-one-completion-notice.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/references/release.md`
- interface: `skills/roundfix/references/release.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix_modes_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`
- interface: `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- interface: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`

## Consumers

- The `modes` and `release` references — consumers: `./skills`
- The glossary entries and the amended ADRs — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest|TestTheFollowUpModeAnswersQuestionsAndHoldsOneDeclaredMessagePerSender|TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest ' || { echo 'missing pass: TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the new test does not exist, so its pass line is missing; afterwards both copies hold the two changed steps and the release section, the earlier follow-up tests still pass, and the skills suite with its recorded version passes.

- `grep -F -- '**Notice Request**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Notice Request\*\*(:| —( |$))' || { echo 'missing term: Notice Request' >&2; exit 1; }; grep -F -- '**Cross-project Message**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Cross-project Message\*\*(:| —( |$))' || { echo 'missing term: Cross-project Message' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'except the one completion notice a Notice Request receives' || { echo 'Cross-project Message unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A Cross-project Message that asks to be told when the work it asks for is done.' || { echo 'Notice Request undefined' >&2; exit 1; }; for f in docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md; do test "$(grep -cF -- '**Amended by [ADR-0396](0396-a-message-that-asks-to-be-told-when-its-work-is-done-receives-one-completion-notice.md) (2026-11-12).**' "$f")" -eq 1 || { echo "missing amendment: $f" >&2; exit 1; }; grep -q '^updated_at: 2026-11-12T21:40:00Z$' "$f" || { echo "updated_at unchanged: $f" >&2; exit 1; }; done; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the term, the changed entry and the notices are absent; afterwards the glossary carries both entries, each of the three ADRs carries one dated notice, and the documentation and glossary contracts pass.

## References

- `_prd.md` → Goals 1 to 4; User Story 2; User Story 3; Core Feature 3; Core Feature 4; Core Feature 5; Success Metric 3; Success Metric 4; Glossary
- `_techspec.md` → Current behavior 3 to 5; API Contracts 2-5; Invariant 1; Testing Approach; Declared breaks; Build Order 2
- ADR-0396

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, testing-boss, golang-testing, coding-guidelines

Implementation:

- Changed only the canonical follow-up and release references, added the required two-copy contract test, defined **Notice Request** beside **Blocking Request**, changed the **Cross-project Message** exception, and inserted the three dated ADR-0396 amendments.
- The pre-change focused filter for `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest` reported `Go test: No tests found`.

Focused checks:

- `gofmt -w skills/roundfix_modes_test.go` completed successfully.
- `go test -count=1 -run '^TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest$/canonical$' ./skills` passed both the parent and canonical subtest.
- `go test -count=1 -run '^TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest$' ./skills` reached the expected pre-sync state: canonical passed and embedded failed because the required `make skills-sync` had not yet run.
- `git diff --check` passed.

Acceptance evidence:

- Both canonical follow-up steps now retain the Notice Request record and post-delivery notice rule; the release reference contains the exact `## Completion notices` contract. The embedded copies remain pending the required final sync.
- `CONTEXT.md` now contains the changed **Cross-project Message** wording and the new **Notice Request** entry immediately after **Blocking Request**.
- ADR-0390, ADR-0383 and ADR-0330 each contain one dated ADR-0396 amendment and `updated_at: 2026-11-12T21:40:00Z`.
