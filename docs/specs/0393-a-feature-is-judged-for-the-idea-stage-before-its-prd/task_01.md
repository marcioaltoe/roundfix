---
task: task_01
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
status: pending
type: backend
complexity: high
---

# Task 01: The judge answers the Idea-stage criteria for one item or one request of the session

## Overview

`roundfix spec judge` gains the stage `idea`. It judges one item that is not
yet a Spec, found in the Follow-up Command's inventory by its reference, or one
request of the session read from a file, asks the judge three measured
criteria in one request, and prints the verdict with each probability and the
lines of the record. When the judge does not answer it lists the criteria for
the agent to answer and exits zero. The `spec` command guide and the Roundfix
skill's `spec` reference describe the stage. The slice is verifiable alone:
its tests run the command against a temporary repository with a fake judge.

## Requirements

1. MUST add the `idea-stage` judgment to the judge's questions file with the
   text of `_techspec.md` → Data Models, word for word, and load it into
   `Questions` with the types of `_techspec.md` → Interfaces. The questions,
   their summaries, both thresholds, the item length and the types that are
   not judged MUST live only in that file (Invariant 7).
2. MUST implement `JudgeIdeaStage` and `IdeaStageVerdict` as `_techspec.md` →
   API Contract 1 says, with its checks in that order and its reasons word
   for word. Hazard: the item's text is written by another project or typed by
   the user and may hold a secret, so the secret scan runs before any request
   and a match sends nothing; inside a Run, without a key, at the ceiling and
   with an unreadable Judge Log nothing is sent either (Invariant 2). An
   answer from a model outside the pin MUST NOT become a verdict (Invariant 3).
   The Judge Log row holds the state's hash and never the item's text.
3. MUST implement the stage `idea` of `roundfix spec judge` as API Contract 2
   says: `--item` resolved only through `ideaStageItem` over
   `followUpInventory`, `--request` read only through `ideaStageRequestItem`,
   exactly one of the two, the Waiting Item guard before the judge, exit 0
   for every report, and exit 2 with each error of API Contract 2, word for
   word. Hazard: a request file holds what the user typed; of it only the
   first line, cut to 200 characters, and the first 4,000 characters may be
   sent, with the project's name, and the command never writes, moves or
   deletes the file (Invariant 1).
4. MUST print the text of API Contract 3 and Surface Transcripts 1 to 5, and
   the JSON document of Data Models with its keys in that order. The output
   MUST always show each answer with its probability and MUST NOT print
   `judged-by: judge` unless the judge answered (Invariant 5). The criteria
   lines are read from the questions file, never from a second copy.
5. MUST change the three help texts as API Contract 4 says, keeping every
   existing line, and move the unknown-stage message to the one that names
   `idea` in the code, in the test that pins it and in every other test that
   pins it on the default branch when the Task runs. MUST NOT edit
   `internal/cli/cli_test.go`, a Governed Path this Spec does not grant: the
   new assertions go into the new test file.
6. MUST add the tests `_techspec.md` → Testing Approach names for
   `internal/judge` and `internal/cli`, each calling `t.Parallel()` first, with
   temporary repositories and homes, a fake judge transport, a fixed clock and
   fake key values built at run time. Each Surface Transcript is asserted
   whole by its one named test.
7. MUST change the canonical Roundfix `spec` reference and the `spec` command
   guide as API Contract 6 says for those two files, word for word, copy the
   reference to its mirror with `make skills-sync`, and add
   `TestTheSpecReferenceDocumentsTheIdeaStage` to a new skills test file,
   with the helper that reads both copies of a skill file and requires each
   phrase once, which task_03 reuses. MUST
   raise the Roundfix skill version only with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   and never write a version number. MUST NOT change the `### QA settlement`
   section of any skill or the grouping sentence of the reference.
8. MUST re-record the Behavior Surface Record as the Task's last edit, after
   the last change to a help text or the guide, only with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change the Follow-up Command's output, any other stage of the
   judge beyond its usage text, the judge's keys, recipients, ceiling or log
   format, or any other command's output or exit code (Invariant 4). Of the
   Governed Paths the Task may change only
   `.agents/skills/roundfix/references/spec.md`,
   `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Add the judgment to the questions file and its loader.
- [ ] Implement the judgment with its guards, its request and its Judge Log row.
- [ ] Add the stage, its two inputs, its text and JSON output and its errors to the command.
- [ ] Write the judge and command tests, the five transcript tests included.
- [ ] Update the help texts, the command guide and the skill reference, and record the skill version.
- [ ] Re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] With a judge that answers, the command prints the verdict, each answer with its probability, the rule with both thresholds, the record lines and the judge's cost, in text and in JSON, and appends one Judge Log row; the three verdicts follow the thresholds at their boundaries.
- [ ] The one request holds the project's name, one item's kind, type, origin, title and text cut at 4,000 characters, and the three questions; of a request file it holds only the first line cut to 200 characters and the first 4,000 characters, and never the file's path.
- [ ] Nothing is sent for an item that matches a secret pattern, a Waiting Item, an item typed fix, perf, refactor or finding, inside a Run, without a key, at the ceiling or with an unreadable Judge Log, and each exits 0 naming its reason; a reference outside the inventory, an active Spec, an unreadable or empty request file and a usage error exit 2.
- [ ] Every existing test of the judge, of `roundfix spec judge` and of `roundfix follow-up` passes, with one edit: the unknown-stage message.
- [ ] The command guide and both copies of the skill reference describe the stage, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0403-a-feature-is-judged-for-the-idea-stage-before-its-prd-and-the-judgment-never-refuses.md`
- instruction: `docs/adr/0404-the-idea-stage-judge-reads-one-item-that-is-not-yet-a-spec.md`
- instruction: `internal/judge/follow_up.go`
- instruction: `internal/judge/follow_up_test.go`
- instruction: `internal/cli/follow_up_inventory.go`
- interface: `internal/judge/questions.json`
- interface: `internal/judge/questions.go`
- interface: `internal/judge/client.go`
- creates: `internal/judge/idea_stage.go`
- creates: `internal/judge/idea_stage_test.go`
- interface: `internal/cli/spec_judge.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/spec_judge_test.go`
- creates: `internal/cli/spec_judge_idea.go`
- creates: `internal/cli/spec_judge_idea_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/idea_stage_judgment_test.go`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The judge's questions file, client and answer evaluation — consumers: `./internal/judge`
- The `roundfix spec judge` command, its help and the follow-up inventory it reads — consumers: `./internal/cli`
- The Roundfix `spec` reference and the owned-skill version record — consumers: `./skills`
- The `spec` command guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestIdeaStageAsksThreeCriteriaInOneRequest|TestIdeaStageVerdictUsesTheThresholds|TestIdeaStageSendsNothingItMustNot|TestIdeaStageFailsOpen|TestIdeaStageCutsTextOnARuneBoundary|TestIdeaStageLogsItsJudgment|TestTheIdeaStageQuestionsAreTheMeasuredOnes)$' ./internal/judge 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestIdeaStageAsksThreeCriteriaInOneRequest TestIdeaStageVerdictUsesTheThresholds TestIdeaStageSendsNothingItMustNot TestIdeaStageFailsOpen TestIdeaStageCutsTextOnARuneBoundary TestIdeaStageLogsItsJudgment TestTheIdeaStageQuestionsAreTheMeasuredOnes; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/judge` — expected: exit 0. Before this Task none of the seven tests exists, so the first pass line is missing; afterwards the judgment asks three criteria in one request, reads the thresholds at their boundaries, sends nothing it must not, fails open, logs its row, holds the measured questions, and the judge package passes.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestSpecJudgeIdeaStageRequiresAnItem|TestSpecJudgeIdeaStageRefusesAnUnknownReference|TestSpecJudgeIdeaStageDoesNotJudgeAFix|TestSpecJudgeIdeaStageListsTheCriteriaWhenTheJudgeDoesNotAnswer|TestSpecJudgeIdeaStageListsTheCriteriaForARequest|TestSpecJudgeIdeaStagePrintsTheVerdictAndTheRecord|TestIdeaStageItemComesOnlyFromTheInventory|TestSpecJudgeIdeaStageReadsAMessageEntry|TestSpecJudgeIdeaStageDoesNotJudgeAWaitingItem|TestSpecJudgeIdeaStageSendsOnlyTheStartOfARequest)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSpecJudgeIdeaStageRequiresAnItem TestSpecJudgeIdeaStageRefusesAnUnknownReference TestSpecJudgeIdeaStageDoesNotJudgeAFix TestSpecJudgeIdeaStageListsTheCriteriaWhenTheJudgeDoesNotAnswer TestSpecJudgeIdeaStageListsTheCriteriaForARequest TestSpecJudgeIdeaStagePrintsTheVerdictAndTheRecord TestIdeaStageItemComesOnlyFromTheInventory TestSpecJudgeIdeaStageReadsAMessageEntry TestSpecJudgeIdeaStageDoesNotJudgeAWaitingItem TestSpecJudgeIdeaStageSendsOnlyTheStartOfARequest; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/cli ./internal/judge && go test -count=1 -timeout 30m -run 'SpecJudge|FollowUp|Help|Usage' ./internal/cli` — expected: exit 0. Before this Task none of the ten tests exists; afterwards each Surface Transcript is asserted whole, an item comes only from the inventory, only the start of a request file is sent, the code builds for Windows, and the existing judge, follow-up and help tests of the command package pass.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSpecReferenceDocumentsTheIdeaStage|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheSpecReferenceDocumentsTheIdeaStage ' || { echo 'missing pass: TestTheSpecReferenceDocumentsTheIdeaStage' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/spec.md | grep -qF -- 'roundfix spec judge --stage idea (--item <reference> [--messages <path>] | --request <path>) [--format <text|json>]' || { echo 'the spec command guide lacks the idea stage synopsis' >&2; exit 1; }; go test -count=1 -timeout 30m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the reference test does not exist; afterwards both copies of the skill reference and the command guide describe the stage, the skill version is recorded, and the skills package and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → User Story 1; User Story 3; User Story 5; User Story 7; Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 4; Core Feature 10; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 8; Success Metric 9
- `_techspec.md` → Current behavior 1 to 4 and 6; Interfaces; Data Models; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 6; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Invariants 1 to 5 and 7; Declared breaks; Testing Approach; Build Order 1
- ADR-0403; ADR-0404; ADR-0200; ADR-0201; ADR-0319
