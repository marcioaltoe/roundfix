---
task: task_03
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
status: completed
type: docs
complexity: medium
---

# Task 03: The follow-up mode replies only to declared blocked work, and the vocabulary and decisions say so

## Overview

The Roundfix skill's follow-up mode fills `blocked_work` from what a message
says, plans a declared message first, replies only to it, and sends its own
findings to other projects expecting no reply, without a wake-up or a Waiting
Declaration to watch them. The glossary adds **Blocking Request** and changes
**Cross-project Message** and **Follow-up Priority**, and ADR-0319 and
ADR-0330 carry ADR-0383's amendment notices.

## Requirements

1. MUST change steps 1, 2, 4 and 8 under `## --follow-up` of the canonical
   `.agents/skills/roundfix/references/modes.md` as `_techspec.md` → API
   Contract 9 says, word for word, and nothing else in the file; MUST NOT add
   or change text inside any `### QA settlement` section. MUST keep every
   phrase `skills/roundfix_modes_test.go` already pins in steps 1, 4 and 8.
2. MUST add `TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork` to
   `skills/roundfix_modes_test.go`, calling `t.Parallel()` first, over both
   copies (`roundfixLayoutSources`, `roundfixSection`,
   `roundfixNumberedItem`, `roundfixWhitespace`), requiring step 1 to hold
   "Give a message entry `blocked_work` only when" and "never from a fixed word
   or phrase", step 2 "is the plan's first item", step 4 "only when its
   message declared blocked work" and not "when a reply is due", and step 8
   "that expects no reply" and "write no Waiting Declaration".
3. MUST change `CONTEXT.md` through `domain-modeling` as API Contract 10 says,
   word for word: the second sentence of **Cross-project Message** and the
   new sentence after it, the new **Blocking Request** entry right after
   **Cross-project Message** with its `_Avoid_` line, and the first paragraph
   of **Follow-up Priority**; MUST add no other capitalized phrase.
4. MUST add to `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
   and `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
   the amendment notices of API Contract 11, word for word, where it places
   them, set each file's `updated_at` to `2026-11-10T19:00:00Z`, and change
   nothing else in either ADR.
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

- [ ] Change steps 1, 2, 4 and 8 of the follow-up mode.
- [ ] Add the modes contract test.
- [ ] Change the glossary entries and add the new term.
- [ ] Write the two amendment notices, then sync and record the skill version.

## Acceptance Criteria

- [ ] Both copies of the `modes` reference fill `blocked_work` from the message, plan a declared message first, reply only to it, and send findings expecting no reply.
- [ ] `CONTEXT.md` defines **Blocking Request** and the changed entries, and the glossary contracts pass.
- [ ] ADR-0319 and ADR-0330 each carry an "Amended by ADR-0383" notice dated 2026-11-10.

## Context

- instruction: `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- instruction: `docs/adr/0338-the-follow-up-mode-claims-the-handoff-and-leaves-no-handed-in-message-pending.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix_modes_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- interface: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`

## Consumers

- The `modes` reference's follow-up steps — consumers: `./skills`
- The glossary entries and the amended ADRs — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork|TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork ' || { echo 'missing pass: TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the new test does not exist, so its pass line is missing; afterwards both copies hold the four changed steps, the earlier pinned phrases still pass, and the skills suite with its recorded version passes.

- `grep -F -- '**Blocking Request**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Blocking Request\*\*(:| —( |$))' || { echo 'missing term: Blocking Request' >&2; exit 1; }; grep -F -- '**Cross-project Message**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Cross-project Message\*\*(:| —( |$))' || { echo 'missing term: Cross-project Message' >&2; exit 1; }; grep -F -- '**Follow-up Priority**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Follow-up Priority\*\*(:| —( |$))' || { echo 'missing term: Follow-up Priority' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'triages each Blocking Request first and the other messages oldest first' || { echo 'Cross-project Message unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A Blocking Request the caller marks leads' || { echo 'Follow-up Priority unchanged' >&2; exit 1; }; for f in docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md; do grep -qF -- '**Amended by [ADR-0383](0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md) (2026-11-10).**' "$f" || { echo "missing amendment: $f" >&2; exit 1; }; grep -q '^updated_at: 2026-11-10T19:00:00Z$' "$f" || { echo "updated_at unchanged: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the term, the changed entries and the notices are absent; afterwards the glossary defines **Blocking Request** and the changed entries, both ADRs carry the dated notice, and the documentation and glossary contracts pass.

## References

- `_prd.md` → User Story 2; User Story 3; Core Feature 6; Core Feature 7; Success Metric 5; Success Metric 6; Glossary
- `_techspec.md` → Current behavior 6; API Contracts 9-11; Testing Approach; Build Order 3
- ADR-0383

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, golang-testing, testing-boss, evidence-gate

Implementation and focused checks:

- Updated steps 1, 2, 4 and 8 under `## --follow-up` in the canonical modes reference, preserving the existing handoff and settlement wording; `make skills-sync` completed successfully and `cmp` confirmed the canonical and mirror references match.
- Added `TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork`, with `t.Parallel()` first and checks over all four `roundfixLayoutSources` using the existing section, numbered-item and whitespace helpers.
- Updated the `Cross-project Message` and `Follow-up Priority` definitions, added `Blocking Request` with its `_Avoid_` line, and added the dated ADR-0383 amendment notices at the specified locations. Only the pre-existing Daemon status change and this Result changed in the Task file.
- `gofmt -d skills/roundfix_modes_test.go`: passed with no output.
- `GOCACHE=/private/tmp/roundfix-task0373-gocache rtk proxy go test ./skills -run '^TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork$' -count=1`: passed (`ok roundfix/skills 0.263s`); the initial cache-only attempt was blocked by the sandbox's global Go cache, then passed with the temporary cache.
- `git diff --check`: passed.

Acceptance evidence:

- Both modes copies contain the four changed steps, and the focused contract test passed across all four layouts.
- `CONTEXT.md` contains the required glossary entries and exact contracted phrases; the full docscontract command remains for Daemon Verification.
- Both ADRs contain the linked `Amended by ADR-0383 (2026-11-10)` notice and `updated_at: 2026-11-10T19:00:00Z`; the full docscontract command remains for Daemon Verification.

The Task status and declared Verification commands were left for the Daemon. The owned-skill version-record command is the final required Task step and is next.
