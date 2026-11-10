---
task: task_01
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
status: pending
type: backend
complexity: medium
---

# Task 01: The Follow-up Command ranks a message that declares blocked work first, without the judge

## Overview

The `roundfix/follow-up-messages/v1` file gains the optional `blocked_work`
key, which the session writing the file fills from what a message says. An
entry that carries it ranks first in `unblocking` with basis `declared`, is
never sent to the judge, and its text line names the blocked work; a
malformed value is refused naming the message. The command guide and both
`follow-up` reference copies document it.

## Requirements

1. MUST add `BlockedWork string` with JSON tag `blocked_work,omitempty` to
   `followUpItem` in `internal/cli/follow_up.go`, declared right after
   `Blocks`, and `BlockedWork string` with JSON tag `blocked_work` to the
   entry struct `followUpMessages` decodes in
   `internal/cli/follow_up_inventory.go`.
2. MUST validate and carry the value as `_techspec.md` → API Contracts 1 and
   2 say: trim it; refuse a value with a line break with
   `message "<id>": blocked_work must be one line` and a non-empty value on an
   entry with a blank `from` with
   `message "<id>": blocked_work needs from, the project whose work is blocked`,
   both through the reader's existing failure (exit `2`, empty standard
   output); keep every existing refusal and its text.
3. MUST order, judge and render as API Contracts 3, 4 and 5 say: add
   `followUpDeclaredOrder` to `internal/cli/follow_up.go`; in
   `runFollowUpCommand` keep declared items out of the judge request and of
   `followUpOrder`, list them first (oldest `created_at`, then `ref`) with
   priority `1`, `unblocking`, basis `declared`, `blocks` equal to
   `requested_by` and null `answers` and `reason`, and renumber ranks; set the
   judge reason `only declared or waiting items` when no item goes to the
   judge and a declared item is listed; count in `Judge: answered <n> of <m>`
   only items whose basis is `judge` or `type`; append `: <blocked_work>` to a
   declared item's text line after ` — blocks <blocks>`. MUST NOT read a
   message body to decide a declaration (Invariant 1).
4. MUST create `internal/cli/follow_up_blocked_work_test.go` with the four
   tests `_techspec.md` → Testing Approach names for `internal/cli`, each
   calling `t.Parallel()` first:
   `TestFollowUpRanksADeclaredBlockingMessageFirstWithoutTheJudge` (its input
   lists the newer declared message before the older one, gives one value
   with surrounding spaces, and the fake judge answers only the two judged
   titles), `TestFollowUpRefusesAMalformedBlockedWork` (rows `line break` and
   `without from`), `TestFollowUpOnlyDeclaredItemsAskNoJudge` and
   `TestFollowUpDeclaredBlockedWorkTranscripts`, whose rows `declared first`
   and `without from` assert Surface Transcripts 1 and 2 whole: every line of
   standard output and standard error, in order, and the exit code, with only
   the transcript's placeholders replaced. No existing test changes.
5. MUST change `docs/user-guide/commands/follow-up.md` and the canonical
   `.agents/skills/roundfix/references/follow-up.md` as API Contract 7 says,
   word for word, then run `make skills-sync` to copy the reference to
   `skills/roundfix/references/follow-up.md`, and raise the Roundfix skill
   version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   never hand-edit a mirror, a version field or the record.
6. MUST, as the Task's last step, after its last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/follow-up.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Add the field to the reader and the item, with its two refusals.
- [ ] Order declared items first, keep them from the judge, and render their line.
- [ ] Add the four tests, the transcripts included.
- [ ] Document the field, sync and record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A message entry with `blocked_work` lists first in `unblocking` with basis `declared`, oldest first, and never reaches the judge.
- [ ] A malformed `blocked_work` exits `2` naming the message; a file without the field keeps its report.
- [ ] The guide, both reference copies, the skill version and the Behavior Surface Record agree with the command.

## Context

- instruction: `docs/adr/0383-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks.md`
- instruction: `docs/adr/0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md`
- instruction: `internal/cli/follow_up_test.go`
- instruction: `internal/cli/follow_up_waiting_test.go`
- interface: `internal/cli/follow_up.go`
- interface: `internal/cli/follow_up_inventory.go`
- creates: `internal/cli/follow_up_blocked_work_test.go`
- interface: `internal/cli/live_documents_docscontract_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The messages file's `blocked_work` key and the item's `blocked_work` and `declared` basis — consumers: `./internal/cli`, `./skills`
- The command guide and the `follow-up` reference — consumers: `./internal/docscontract`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestFollowUpRanksADeclaredBlockingMessageFirstWithoutTheJudge|TestFollowUpRefusesAMalformedBlockedWork|TestFollowUpOnlyDeclaredItemsAskNoJudge|TestFollowUpDeclaredBlockedWorkTranscripts)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFollowUpRanksADeclaredBlockingMessageFirstWithoutTheJudge TestFollowUpRefusesAMalformedBlockedWork TestFollowUpOnlyDeclaredItemsAskNoJudge TestFollowUpDeclaredBlockedWorkTranscripts; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'FollowUp' ./internal/cli && go test -count=1 -timeout 30m -tags docscontract -run 'TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes' ./internal/cli` — expected: exit 0. Before this Task the four tests do not exist, so their pass lines are missing; afterwards they pass, every existing Follow-up Command test passes unchanged, and the guides list the item's fields in struct order.

- `for f in docs/user-guide/commands/follow-up.md .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md; do grep -q '^## Declared blocked work$' "$f" || { echo "missing section: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'not_asked. with reason .only declared or waiting items' || { echo "missing judge reason: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide has the section; afterwards the guide and both reference copies document the field, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record pass.

## References

- `_prd.md` → Goal 4; User Story 4; Core Feature 5; Core Feature 7; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 6
- `_techspec.md` → Current behavior 1-4; API Contracts 1-7; Surface Transcript 1; Surface Transcript 2; Invariants 1-3; Declared breaks; Testing Approach; Build Order 1
- ADR-0383
