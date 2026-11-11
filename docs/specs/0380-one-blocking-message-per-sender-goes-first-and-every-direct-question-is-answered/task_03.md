---
task: task_03
spec: 0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
status: pending
type: docs
complexity: medium
---

# Task 03: The follow-up mode, the vocabulary and the decisions follow the bounded precedence and reply rule

## Overview

The Roundfix skill's follow-up mode fills a message entry's `from` from the
sender ai-memory returns, gives `blocked_work` to every declaring message,
keeps the place the command gives a held-back one and the maintainer's
judgment over the order, and answers each direct question. The glossary
changes **Cross-project Message**, **Blocking Request** and **Follow-up
Priority**, and ADR-0383, ADR-0319 and ADR-0330 carry ADR-0390's amendment
notices.

## Requirements

1. MUST change steps 1, 2 and 4 under `## --follow-up` of the canonical
   `.agents/skills/roundfix/references/modes.md` as `_techspec.md` → API
   Contract 7 says, word for word, and nothing else in the file; MUST NOT add
   or change text inside any `### QA settlement` section. MUST keep every
   phrase `skills/roundfix_modes_test.go` pins in steps 1, 4 and 8 other than
   the step 4 phrase API Contract 7 replaces.
2. MUST rename `TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork` in
   `skills/roundfix_modes_test.go` to
   `TestTheFollowUpModeAnswersQuestionsAndHoldsOneDeclaredMessagePerSender`
   and change its assertions as `_techspec.md` → Testing Approach says for
   `skills`, keeping `t.Parallel()` first and its use of
   `roundfixLayoutSources`, `roundfixSection`, `roundfixNumberedItem` and
   `roundfixWhitespace`. No other test changes.
3. MUST change `CONTEXT.md` through `domain-modeling` as API Contract 8 says,
   word for word: the definitions of **Cross-project Message** and
   **Blocking Request** and the first paragraph of **Follow-up Priority**;
   MUST add no other capitalized phrase and no new entry.
4. MUST add the four amendment notices of API Contract 9, word for word,
   where it places them: two to
   `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`,
   one to `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
   and one to
   `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`;
   set each file's `updated_at` to `2026-11-11T18:00:00Z`, and change nothing
   else in the three ADRs.
5. MUST run `make skills-sync` to copy the canonical reference to
   `skills/roundfix/references/modes.md`, then raise the Roundfix skill
   version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   as the Task's last step; never hand-edit a mirror, a version field or the
   record. Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/modes.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Change steps 1, 2 and 4 of the follow-up mode.
- [ ] Rename and update the modes contract test.
- [ ] Change the three glossary entries.
- [ ] Write the four amendment notices, then sync and record the skill version.

## Acceptance Criteria

- [ ] Both copies of the `modes` reference fill `from` from ai-memory's sender, give every declaring message `blocked_work`, keep the command's place and the maintainer's judgment, and answer every direct question.
- [ ] `CONTEXT.md` carries the changed entries, and the glossary contracts pass.
- [ ] ADR-0383, ADR-0319 and ADR-0330 carry an "Amended by ADR-0390" notice dated 2026-11-11.

## Context

- instruction: `docs/adr/0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md`
- instruction: `docs/adr/0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix_modes_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- interface: `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- interface: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`

## Consumers

- The `modes` reference's follow-up steps — consumers: `./skills`
- The glossary entries and the amended ADRs — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheFollowUpModeAnswersQuestionsAndHoldsOneDeclaredMessagePerSender|TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheFollowUpModeAnswersQuestionsAndHoldsOneDeclaredMessagePerSender ' || { echo 'missing pass: TestTheFollowUpModeAnswersQuestionsAndHoldsOneDeclaredMessagePerSender' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the renamed test does not exist, so its pass line is missing; afterwards both copies hold the three changed steps, the earlier pinned phrases still pass, and the skills suite with its recorded version passes.

- `grep -F -- '**Blocking Request**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Blocking Request\*\*(:| —( |$))' || { echo 'missing term: Blocking Request' >&2; exit 1; }; grep -F -- '**Cross-project Message**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Cross-project Message\*\*(:| —( |$))' || { echo 'missing term: Cross-project Message' >&2; exit 1; }; grep -F -- '**Follow-up Priority**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Follow-up Priority\*\*(:| —( |$))' || { echo 'missing term: Follow-up Priority' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'triages first the oldest Blocking Request of each sending project' || { echo 'Cross-project Message unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'so it takes precedence once per sender' || { echo 'Blocking Request unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'The oldest Blocking Request of each sender that the caller marks leads' || { echo 'Follow-up Priority unchanged' >&2; exit 1; }; for f in docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md; do grep -qF -- '**Amended by [ADR-0390](0390-one-blocking-message-per-sender-takes-precedence-and-every-direct-question-is-answered.md) (2026-11-11).**' "$f" || { echo "missing amendment: $f" >&2; exit 1; }; grep -q '^updated_at: 2026-11-11T18:00:00Z$' "$f" || { echo "updated_at unchanged: $f" >&2; exit 1; }; done; test "$(grep -cF -- '**Amended by [ADR-0390]' docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md)" -eq 2 || { echo 'ADR-0383 needs two notices' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the changed entries and the notices are absent; afterwards the glossary carries the changed entries, the three ADRs carry the dated notices (two in ADR-0383), and the documentation and glossary contracts pass.

## References

- `_prd.md` → Goals 1 to 3; User Story 1; User Story 2; Core Feature 4; Core Feature 5; Success Metric 4; Success Metric 5; Glossary
- `_techspec.md` → Current behavior 6; API Contracts 7-9; Testing Approach; Declared breaks; Build Order 3
- ADR-0390
