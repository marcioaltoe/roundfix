---
task: task_02
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
status: pending
type: backend
complexity: high
---

# Task 02: The follow-up holds Waiting Items back, lists them with --all and reports a malformed declaration

## Overview

An open Backlog Entry or an unresolved Finding may carry a Waiting
Declaration, `waiting_on` and `since`. With a valid one it is a Waiting Item
in the class `waiting`: `roundfix follow-up` leaves it out of the list, counts
it on its source line, never sends it to the judge, and lists it after the
ranked items with `--all`. A malformed declaration is reported on its item,
which stays in its ordinary class. A repository without the fields keeps its
report bytes. This answers the first request of the adopted Backlog Entry of
2026-11-02, "The follow-up names what waits on someone else and where each
item went"
([references](references/2026-11-02-the-follow-up-names-what-waits-and-where-each-item-went.md)),
reported by fluxus, with the maintainer's `--all` answer of 2026-11-03.

## Requirements

1. MUST first characterize the starting tree: before any code change, create
   `internal/cli/follow_up_waiting_test.go` with
   `TestFollowUpKeepsItsBytesWithoutAWaitingDeclaration` over a fixture
   holding a handed-in message, an open Backlog Entry, a pending Finding and a
   visible note, none declaring a wait, and record its text and
   `--format json` output, exactly as the starting tree prints them, as
   constants in which only the repository root is substituted. The test must
   pass unchanged after the work (`_techspec.md` → Invariant 1).
2. MUST read the declaration as `_techspec.md` → Interfaces and API Contract 2
   state: `followUpMarkdown` gains `WaitingOn` and `Since` as `yaml.Node`
   fields, and `followUpWaiting` returns the trimmed `waiting_on`, the
   `since` text or exactly one of the four `waiting_error` texts
   `waiting_on must be one line of non-empty text`, `since is missing`,
   `since must be a YYYY-MM-DD date`, `since needs waiting_on`, in the
   contract's order. `followUpInventory` applies it to `backlog` and
   `finding` items only, and sets `Waiting` on the `backlog` and `findings`
   sources; a non-scalar value never fails the read.
3. MUST split, order and render as API Contracts 3, 4 and 6 state: only
   ranked items reach `judge.PrioritizeFollowUp` and `followUpOrder`; without
   `--all` no Waiting Item is listed; with `--all` they follow every ranked
   item sorted by `since` and then reference, with priority 5,
   `priority_name` and `basis` `waiting`, null `answers`, `blocks` and
   `reason`; the source, item and judge lines print the new suffixes and
   every other line keeps its bytes.
4. MUST add `--all` to `parseFollowUpCommand`, refuse it with `--status` with
   `--all applies only to the follow-up list; drop it with --status` after
   the existing `--messages` check, and change `followUpUsage` and the root
   usage line in `internal/cli/cli.go` exactly as API Contract 5 and Surface
   Transcript 1 state. `runFollowUpStatus` stays unchanged (API Contract 7).
5. MUST add the optional JSON keys of `_techspec.md` → Interfaces and Data
   Models with `omitempty`: `waiting` on a source, `waiting_on`, `since` and
   `waiting_error` on an item.
6. MUST add, in `internal/cli/follow_up_waiting_test.go`,
   `TestFollowUpHoldsBackWaitingItems`, `TestFollowUpAllListsWaitingItemsLast`
   and `TestFollowUpReportsAnInvalidWaitingDeclaration` with the cases
   `_techspec.md` → Testing Approach lists, using the fake `gh` and judge of
   `follow_up_test.go`; the judge cases set a fixture key and assert the
   titles the fake judge received. Every new top-level test and subtest calls
   `t.Parallel()` first.
7. MUST add `TestFollowUpHoldsBackWaitingItemsTranscript`,
   `TestFollowUpAllListsWaitingItemsTranscript` and
   `TestFollowUpStatusRefusesAllTranscript`, each asserting the whole Surface
   Transcript 2, 3 or 4 in one test: every line of standard output and
   standard error, in order, and the exit code, over the committed fixture
   the Surface Transcripts section describes, with only `<namespace>`,
   `<root>`, `<reason>` and `<pull request source>` replaced by the fixture's
   values.
8. MUST update `TestFollowUpListHelpMatchesTranscript`
   (`internal/cli/follow_up_test.go`) and `TestFollowUpHelpTranscript`
   (`internal/cli/follow_up_status_test.go`) to Surface Transcript 1 and the
   new root synopsis, keeping both test names (a declared break,
   `_techspec.md` → Declared breaks).
9. MUST document the contract in `docs/user-guide/commands/follow-up.md` and
   the canonical `.agents/skills/roundfix/references/follow-up.md` (API
   Contract 9): the synopsis with `[--all]`; the Waiting Declaration and its
   two fields; the class `waiting`; the sentence "A Waiting Item is never sent
   to the judge, with or without --all."; the held-back count on the source
   line; the four `waiting_error` texts; `--all` refused with `--status`; and
   the new keys in the field sentences "Each source carries these JSON
   fields:" and "Each item carries these JSON fields:", in struct order.
10. MUST, through `domain-modeling`, add **Waiting Declaration** and
    **Waiting Item** to `CONTEXT.md`, each with a link to ADR-0353, and revise
    **Follow-up Priority** (a Waiting Item is in the fifth class `waiting`,
    outside the four) and **Follow-up Command** (it holds Waiting Items back
    unless `--all` is given).
11. MUST create
    `docs/adr/0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md`
    from `_techspec.md` → ADR-0353, with the lifecycle overlay it names and
    the sentence "A Waiting Item is never sent to the judge".
12. MUST copy the canonical skill with `make skills-sync`, then raise both
    version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
    `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
    -record-skill-versions -count=1`. Never hand-edit a mirror or a version
    field.
13. MUST re-record the Behavior Surface Record as the LAST step, after the
    last edit of any guide, with `go test -count=1 -tags docscontract
    ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
    -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
14. MUST NOT touch `internal/judge`, `internal/speccheck`,
    `internal/cli/follow_up_status.go`, the `modes` and `deliver` references,
    any Baseline asset, `docs/agents/`, any other skill, any `### QA
    settlement` section, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the current text and JSON bytes in the new test file.
- [ ] Read the declaration and count Waiting Items per source.
- [ ] Split the ranked and Waiting Items, add `--all`, the JSON keys and the text suffixes.
- [ ] Add the behavior and transcript tests and update the two help tests.
- [ ] Update the guide, reference, `CONTEXT.md` and ADR-0353; sync, record the skill version, re-record the fingerprints last.

## Acceptance Criteria

- [ ] A repository without a declaration prints the characterized bytes, and only the help gains `--all`.
- [ ] A valid declaration holds the item back with a count on its source line, `--all` lists it last with its reason and date, and the judge never receives it.
- [ ] Each malformed declaration lists its item in its ordinary class with its reason, and no source becomes unavailable.
- [ ] The four Surface Transcripts are asserted in full, and the guide, reference, glossary and ADR-0353 state the contract.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0319-the-follow-up-order-is-code-over-four-jev-yes-no-answers.md`
- interface: `internal/cli/follow_up.go`
- interface: `internal/cli/follow_up_inventory.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/follow_up_test.go`
- interface: `internal/cli/follow_up_status_test.go`
- creates: `internal/cli/follow_up_waiting_test.go`
- instruction: `internal/cli/live_documents_docscontract_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- creates: `docs/adr/0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md`

## Consumers

- The follow-up report (`roundfix/follow-up/v1`, its text, flags and help) — consumers: `./internal/cli`
- The follow-up guide's and reference's JSON field sentences — consumers: `./internal/cli`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestFollowUpKeepsItsBytesWithoutAWaitingDeclaration|TestFollowUpHoldsBackWaitingItems|TestFollowUpAllListsWaitingItemsLast|TestFollowUpReportsAnInvalidWaitingDeclaration|TestFollowUpHoldsBackWaitingItemsTranscript|TestFollowUpAllListsWaitingItemsTranscript|TestFollowUpStatusRefusesAllTranscript|TestFollowUpListHelpMatchesTranscript|TestFollowUpHelpTranscript|TestFollowUpListsEveryPendingItemOnce|TestFollowUpOrdersItemsByTheFourAnswers|TestFollowUpReportsSkippedItemsAndAStoppedJudge|TestFollowUpSaysNothingToDo|TestFollowUpStatusRefusesMessages|TestFollowUpIgnoresHiddenInboxEntries)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFollowUpKeepsItsBytesWithoutAWaitingDeclaration TestFollowUpHoldsBackWaitingItems TestFollowUpAllListsWaitingItemsLast TestFollowUpReportsAnInvalidWaitingDeclaration TestFollowUpHoldsBackWaitingItemsTranscript TestFollowUpAllListsWaitingItemsTranscript TestFollowUpStatusRefusesAllTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the seven new tests exists, so the first pass line is missing; afterwards the characterized bytes hold, Waiting Items are held back, listed with `--all` and never judged, malformed declarations are reported, the four transcripts match, and the existing follow-up tests pass.
- `for f in docs/user-guide/commands/follow-up.md .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md docs/adr/0353-a-waiting-item-is-held-back-from-the-follow-up-list-unless-all-is-asked.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'A Waiting Item is never sent to the judge' || { printf 'missing judge rule in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/follow-up.md .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'waiting_on must be one line of non-empty text' || { printf 'missing validation reasons in %s\n' "$f" >&2; exit 1; }; done; grep -F -- '**Waiting Declaration**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Waiting Declaration\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md: Waiting Declaration' >&2; exit 1; }; grep -F -- '**Waiting Item**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Waiting Item\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md: Waiting Item' >&2; exit 1; }; grep -F -- '**Follow-up Priority**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Follow-up Priority\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md: Follow-up Priority' >&2; exit 1; }; grep -F -- '**Follow-up Command**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Follow-up Command\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md: Follow-up Command' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'the fifth class' || { echo 'missing Follow-up Priority revision' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'holds Waiting Items back' || { echo 'missing Follow-up Command revision' >&2; exit 1; }; cmp .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m -tags docscontract -run '^TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes$' ./internal/cli && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide, reference, glossary entry or ADR names the Waiting Item, so the first check fails; afterwards they do, the field sentences name the new JSON keys, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; Core Features 2 to 5 and 7; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Glossary; Decisions
- `_techspec.md` → Current behavior 2 to 5; Interfaces; Data Models; API Contracts 2 to 7 and 9; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Invariants 1 to 4 and 6; Declared breaks; ADR-0353; Testing Approach; Build Order 2
