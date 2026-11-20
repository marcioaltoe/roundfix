---
task: task_06
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
status: pending
type: backend
complexity: low
---

# Task 06: The status view prints the repository text of this Spec as one terminal-safe line

## Overview

Corrective Task 1 of 2. The pre-PR review of 2026-11-19 (provider codex, artifact `043751504adb2ac5`, reviewed head `b38daf3d`, finding F1, anchored at `internal/cli/follow_up_status_delegated.go:72`) found that the question and the adopted option of a Delegated Decision, which a PRD controls, are written to the terminal as they are, so a crafted entry can inject ANSI control sequences or forged status lines into the text of `roundfix follow-up --status`. The finding is valid: `_techspec.md` and the user guide call the text of an entry repository text, and repository text is data.

The status view already has the rule for repository text: `followUpQueueText` in `internal/cli/follow_up_queue.go` folds `\r\n`, `\r` and `\n` to one space and turns every control character (C0, DEL, C1) and every line or paragraph separator into a space, then cuts the text at 120 characters for the Spec queue table. This Task reuses that rule without the cut, so a long question is still printed whole, as API Contract 2 shows, and applies it to every value of repository or Project Config text that this Spec added to the text of the status view: the Delegated Decisions block, the `Open authoring questions` line and the detail of the `specs` source. The JSON output keeps the raw strings, which JSON escapes.

## Requirements

1. MUST answer finding F1 of the pre-PR review of 2026-11-19: no value this Spec adds to the text of `roundfix follow-up --status` reaches standard output with a control character, a line break or a line or paragraph separator.
2. MUST add to `internal/cli/follow_up_queue.go` the function `followUpTerminalText(value string) string`, which does exactly the folding and the mapping `followUpQueueText` does now, with no cut, and make `followUpQueueText` return `followUpTerminalText(value)` cut at 120 characters as now, so every output of `followUpQueueText` stays byte for byte the same.
3. MUST make `renderFollowUpDelegatedDecisions` in `internal/cli/follow_up_status_delegated.go` pass the Spec, the date, the phase, the question, the option, the answer and the problem through `followUpTerminalText`, keeping the two line forms of `_techspec.md` → API Contract 2, the option and the answer still taken before the first ` — `.
4. MUST make `renderFollowUpStatusText` in `internal/cli/follow_up_status.go` pass the preference, the checkout value and the detail of the `Open authoring questions` line, and the detail of each line under `Sources:`, through `followUpTerminalText`. A value with no control character, line break or separator prints unchanged.
5. MUST leave `followUpDelegatedDecision`, `followUpAuthoringQuestions`, `followUpSource` and the JSON document unchanged: `--format json` keeps the raw question, option and detail.
6. MUST add `TestTheStatusViewPrintsRepositoryTextAsOneTerminalSafeLine` to `internal/cli/follow_up_status_delegated_test.go`, with the subtests:
   - `an_ansi_escape_in_the_question`: through the CLI, with the fixture and helpers of `TestTheStatusViewListsTheDelegatedDecisionsOfActiveSpecs`, the Spec `0001-expired-imports` holds `statusDelegatedEntry` with the question `which \x1b[31mfailure\x1b[0m behavior?` and the adopted line `retry\x1b[2J automatically — it preserves the sync retry behavior`; the text block between `Delegated decisions:` and `Sources:` is exactly `Delegated decisions: 1\n  0001-expired-imports 2026-11-20 prd: which  [31mfailure [0m behavior? -> retry [2J automatically (answer: none yet)\n`, standard output holds no byte `0x1b`, and the JSON row keeps the raw question and adopted text;
   - `a_carriage_return_in_the_question`: through the CLI, the question is `which failure\r  0002-forged 2026-11-20 prd: forged -> yes (answer: none yet)`; the block is exactly `Delegated decisions: 1\n  0001-expired-imports 2026-11-20 prd: which failure   0002-forged 2026-11-20 prd: forged -> yes (answer: none yet) -> retry automatically up to 3 times (answer: none yet)\n`, standard output holds no `\r`, and the JSON row keeps the raw question;
   - `a_newline_in_the_question`: a PRD field is read one line at a time, so this case calls `renderFollowUpDelegatedDecisions` with `statusDelegatedRow("0001-expired-imports")` whose question is `first\nDelegated decisions: 9` and whose adopted text is `retry\x1b[2J — it preserves the sync retry behavior`; the result is exactly `Delegated decisions: 1\n  0001-expired-imports 2026-11-20 prd: first Delegated decisions: 9 -> retry [2J (answer: none yet)\n`;
   - `a_normal_entry_is_unchanged`: `renderFollowUpDelegatedDecisions` of `statusDelegatedRow("0001-expired-imports")` is exactly `Delegated decisions: 1\n  0001-expired-imports 2026-11-20 prd: which failure behavior does an expired import have? -> retry automatically up to 3 times (answer: none yet)\n`, and a row whose question is 300 characters `é` prints the 300 characters whole, with no `…`;
   - `the_preference_line_and_the_specs_source_are_one_line`: `renderFollowUpStatusText` of a `followUpStatusDocument` whose `AuthoringQuestions` is `{Preference: "wait", Source: "default", Checkout: "fol\x1b[2Jlow", Detail: "the default branch's Project Config was not read: bad\nOpen authoring questions: follow (Project Config)\x1b[0m"}` and whose `Sources` hold `{Name: "specs", State: "unavailable", Detail: "read PRD for Spec 0001-x\x1b[31m: denied\r\nforged"}` holds exactly one line that begins `Open authoring questions:`, equal to `Open authoring questions: wait (default); this checkout's Project Config gives fol [2Jlow, which counts only once the default branch holds it; the default branch's Project Config was not read: bad Open authoring questions: follow (Project Config) [0m`, holds the line `  specs: unavailable: read PRD for Spec 0001-x [31m: denied forged`, and holds no rune for which `unicode.IsControl` is true other than `\n`.
7. MUST keep every existing test of `internal/cli` unchanged and passing, `TestTheNoteIsOneLineOfRepositoryText`, `TestTheStatusViewListsTheDelegatedDecisionsOfActiveSpecs`, `TestTheStatusViewAuthoringSurfaceTranscripts`, `TestTheStatusViewKeepsEveryUnreadablePRDInTheSpecsSource` and `TestTheStatusViewReportsTheAuthoringQuestionPreference` among them, and change no guide, glossary, record, skill, Baseline module, the Makefile, a workflow, `go.mod` or a lint configuration.

## Subtasks

- [ ] Split `followUpTerminalText` out of `followUpQueueText`.
- [ ] Render the Delegated Decisions block, the preference line and the source details through it.
- [ ] Add the five subtests.

## Acceptance Criteria

- [ ] An ANSI escape, a carriage return or a newline in a Delegated Decision, the preference line or a source detail comes out of the text of the status view as spaces on one line, and the JSON keeps the raw text.
- [ ] A normal entry prints exactly as before, a long question is not cut, and every output of `followUpQueueText` is unchanged.

## Context

- interface: `internal/cli/follow_up_queue.go`
- interface: `internal/cli/follow_up_queue_test.go`
- interface: `internal/cli/follow_up_status.go`
- interface: `internal/cli/follow_up_status_delegated.go`
- interface: `internal/cli/follow_up_status_delegated_test.go`
- interface: `internal/cli/follow_up_status_authoring_test.go`

## Consumers

- `followUpTerminalText` and `followUpQueueText` — consumers: `./internal/cli`
- The text of `roundfix follow-up --status` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheStatusViewPrintsRepositoryTextAsOneTerminalSafeLine|TestTheNoteIsOneLineOfRepositoryText)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheStatusViewPrintsRepositoryTextAsOneTerminalSafeLine/an_ansi_escape_in_the_question TestTheStatusViewPrintsRepositoryTextAsOneTerminalSafeLine/a_carriage_return_in_the_question TestTheStatusViewPrintsRepositoryTextAsOneTerminalSafeLine/a_newline_in_the_question TestTheStatusViewPrintsRepositoryTextAsOneTerminalSafeLine/a_normal_entry_is_unchanged TestTheStatusViewPrintsRepositoryTextAsOneTerminalSafeLine/the_preference_line_and_the_specs_source_are_one_line TestTheNoteIsOneLineOfRepositoryText/line_breaks_become_spaces TestTheNoteIsOneLineOfRepositoryText/control_characters_never_reach_the_terminal TestTheNoteIsOneLineOfRepositoryText/a_long_waiting_declaration_is_cut; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m -run '^(TestTheStatusViewListsTheDelegatedDecisionsOfActiveSpecs|TestTheStatusViewAuthoringSurfaceTranscripts|TestTheStatusViewKeepsEveryUnreadablePRDInTheSpecsSource|TestTheStatusViewReportsTheAuthoringQuestionPreference|TestTheNoteIsOneLineOfRepositoryText|TestAnUnreadableSourceLeavesTheOtherRows|TestTheQueueRowsComeInTableOrder)$' ./internal/cli` — expected: exit 0. Before this Task the test does not exist, so its five subtests print no pass line and the command exits 1; with only the test added, the escape, carriage return, newline and preference subtests fail because the text is printed raw, and the command exits 1; afterwards all five pass, the three `followUpQueueText` cases still pass, and the existing status and queue tests keep their exact text, so a normal entry, the preference line and the source lines print as before.

## References

- Pre-PR review of 2026-11-19 (artifact `043751504adb2ac5`, head `b38daf3d`) → F1
- `_techspec.md` → API Contract 2; Invariant 2
- task_02 → the Delegated Decisions block; task_01 → the `Open authoring questions` line
