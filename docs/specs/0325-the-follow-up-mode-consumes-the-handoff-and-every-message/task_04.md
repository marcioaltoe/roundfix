---
task: task_04
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
status: pending
type: test
complexity: low
---

# Task 04: The modes test catches the restored "without consuming" sentence and names the status refusal

## Overview

Corrective Task 1 of 2. It answers findings F-01 and F-02 of the QA Report of
2026-10-26 (Run `run_20261026T232417Z_2824de576977e0a5`, rows Q05f and Q05c,
task_03 Requirement 5). The delivered skill text is correct; only
`TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage` in
`skills/roundfix_modes_test.go` is too weak or names the wrong diagnostic:

- F-01: replacing "List this project's messages with `memory_message_list`"
  with "List this project's messages without consuming them" in
  `skills/roundfix/references/modes.md` keeps the test green. The negative
  check searches the lowercase phrase "list this project's messages without
  consuming them" case-sensitively, the mutation starts with an uppercase
  "List", and no positive assertion requires the removed
  `memory_message_list` sentence.
- F-02: removing "accepts no handoff, " from `## --status` fails the test, but
  with "--status omits its handoff refusal" instead of the diagnostic task_03
  Requirement 5 declares, "--status section does not say it accepts no
  handoff".

## Requirements

1. MUST make the `Optional sources` check in
   `TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage` require the
   sentence "List this project's messages with `memory_message_list`"
   (whitespace-normalized, as the other phrases are) and fail with the exact
   message `Optional sources still lists messages only to leave them pending`
   when that sentence is missing or when the section contains "without
   consuming them" in any letter case. Keep every other `Optional sources`
   phrase check unchanged.
2. MUST change the `--status` diagnostic to the exact message
   `--status section does not say it accepts no handoff`, keeping the phrase
   it checks unchanged.
3. MUST keep both layouts (`canonical` and `embedded`) running the same
   checks, and every other assertion and diagnostic of the test unchanged.
4. MUST prove both negative controls before handing back, each in its own
   temporary `git archive` copy of the tree outside the checkout, running
   `go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills`
   there:
   - replace "List this project's messages with `memory_message_list`" with
     "List this project's messages without consuming them" in
     `skills/roundfix/references/modes.md`: the `embedded` subtest fails with
     `Optional sources still lists messages only to leave them pending`;
   - remove "accepts no handoff, " from `## --status` in
     `skills/roundfix/references/modes.md`: the `embedded` subtest fails with
     `--status section does not say it accepts no handoff`.
   Write nothing under the checkout for these proofs, and leave no copy behind.
5. MUST NOT change the skill text, its mirrors, its version record, any other
   test, guide, record or planning file. The findings are test gaps, not
   product defects.

## Subtasks

- [ ] Require the `memory_message_list` sentence and reject "without consuming them" in any case, with the declared diagnostic.
- [ ] Rename the `--status` diagnostic to the declared one.
- [ ] Prove both sabotages fail with the declared messages in temporary copies.

## Acceptance Criteria

- [ ] The test passes on the delivered skill text in both layouts.
- [ ] The F-01 mutation fails `embedded` with `Optional sources still lists messages only to leave them pending`.
- [ ] The F-02 mutation fails `embedded` with `--status section does not say it accepts no handoff`.

## Consumers

None. The Task changes one test file and no contract another package reads.

## Context

- interface: `skills/roundfix_modes_test.go`
- instruction: `skills/roundfix/references/modes.md`
- instruction: `docs/specs/0325-the-follow-up-mode-consumes-the-handoff-and-every-message/task_03.md`

## Verification

- `grep -qF 'Optional sources still lists messages only to leave them pending' skills/roundfix_modes_test.go && grep -qF -- '--status section does not say it accepts no handoff' skills/roundfix_modes_test.go && go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills && go test -timeout 15m -count=1 ./skills` — expected: exit 0. Before this Task neither declared diagnostic exists in the test file, so the first check fails; afterwards both diagnostics exist, the named test passes in both layouts and the whole `skills` package passes.

## References

- `_techspec.md` → Testing Approach
- `task_03.md` → Requirement 5
