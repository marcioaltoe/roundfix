---
task: task_03
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
status: completed
type: docs
complexity: medium
---

# Task 03: The Idea and the PRD record each Rule Challenge, and a Spec waits for the answer

## Overview

The Idea template and the PRD template gain the `## Rule challenges` section, with its fielded record and the single entry `None.` when no rule blocked the goal. The Roundfix skill's `modes` reference says that `--unattended` and `--follow-up` ask a **Rule Challenge** and never answer it, and that a Spec without an answer waits. The glossary defines **Rule Challenge** and **Kept Line**, and the user guide describes the question, the three answers and the record. Contract tests pin the record section and the `modes` paragraphs in both copies, and the documentation contracts hold the glossary and the guide.

## Requirements

1. MUST add the record section of `_techspec.md` → API Contract 6, word for word, to the canonical `write-idea/references/idea-template.md` immediately before its `## Decisions` heading, with the example line `- asked: YYYY-MM-DD, idea`, and to the canonical `write-prd/references/prd-template.md` immediately before its `## Glossary` heading, with `- asked: YYYY-MM-DD, prd`. The two sections are the same text except that phase word (`_techspec.md` → Invariant 7). The TechSpec template does not change.
2. MUST make the two changes of `_techspec.md` → API Contract 7 to the canonical Roundfix `references/modes.md`, word for word: under `## --unattended`, the new paragraph right after the paragraph that begins "The Idea stage is the user's"; under `## --follow-up`, the sentence at the end of step 5, after "its Waiting Declaration naming the Idea stage.". The `--follow-up` section keeps eight numbered steps, and every other sentence of the reference, those Spec 0411 changes included, stays as the Task finds it.
3. MUST create `skills/rule_challenge_record_test.go` with two tests, each calling `t.Parallel()` first and checking both the canonical and the shipped copy through `assertIdeaStageSkillPhrases`, each phrase exactly once on whitespace-normalized text:
   - `TestTheTemplatesCarryTheRuleChallengeRecord`, with the subtests `idea/heading_and_position`, `idea/fields`, `idea/none`, `idea/phase` and `prd/heading_and_position`, `prd/fields`, `prd/none`, `prd/phase`: the heading `## Rule challenges` occurs once and the next level-two heading is `## Decisions` in the Idea template and `## Glossary` in the PRD template; one phrase for each of the seven fields `rule`, `imposes`, `prevents`, `asked`, `answer`, `effect` and `new fact`, each taken from the section so that it occurs once; the entry `None.` of the comment; and `asked: YYYY-MM-DD, idea` or `asked: YYYY-MM-DD, prd`.
   - `TestTheSkillModesNeverAnswerARuleChallenge`, with the subtests `unattended`, holding "A Rule Challenge is the user's too: this mode asks it and never answers it." and "A Spec whose Rule Challenge gets no answer writes no artifact for that phase", and `follow_up_step_five`, holding "A Spec whose author returns a Rule Challenge goes on only after the user answered it; without an answer it waits as" and requiring it after "its Waiting Declaration naming the Idea stage." and before the line that starts step 6.
4. MUST regenerate, never by hand: the mirrors under `skills/` with `make skills-sync`, then the owned-skill versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of `write-idea`, `write-prd` and `roundfix` on the tree the Task runs on, then `make skills-sync` again so each mirror equals its source.
5. MUST, through the `domain-modeling` skill, add to `CONTEXT.md` the two entries of `_techspec.md` → API Contract 8, **Rule Challenge** and **Kept Line**, word for word, right after the entry of Idea Stage Judgment and its `_Avoid_` line. No other glossary entry changes.
6. MUST add to `docs/user-guide/context-driven-development.md`, immediately before the heading `### Spec Event Announcements`, the subsection `### Rule Challenges` of `_techspec.md` → API Contract 9, word for word.
7. MUST, as its last edit, after the last change to the user guide, re-record the Behavior Surface Record with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, never by hand. The guide page's covering skill `write-prd` changes in task_02 and in this Task, so `docs/references/skill-coverage.json` does not change.
8. MUST NOT change any other owned skill or skill file, any Vendored Skill, the `### QA settlement` section of any skill, any ADR, any file under `docs/agents/` or `internal/`, `docs/references/skill-coverage.json`, `docs/references/coverage-record.json`, the Makefile, a workflow, `go.mod` or a lint configuration. Of the Governed Paths it changes only the eight files below that `_authorization.md` lists: both copies of the `SKILL.md` of `write-idea`, `write-prd` and `roundfix`, the canonical `idea-template.md` and `modes.md`, and both copies of `prd-template.md`.

## Subtasks

- [ ] Add the record section to the Idea and PRD templates.
- [ ] Add the paragraph and the sentence of the `modes` reference.
- [ ] Add the two contract tests, sync the mirrors and record the owned-skill versions.
- [ ] Add the two glossary entries and the user guide subsection, and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of both templates carry the `## Rule challenges` section in its place, with the seven fields, `None.` and their own phase word.
- [ ] Both copies of the `modes` reference say that each mode asks a Rule Challenge and never answers it, and step 5 of `--follow-up` waits for the answer; the section keeps eight steps.
- [ ] `CONTEXT.md` defines both terms, the user guide holds the subsection, the Behavior Surface Record is current, and the documentation contracts pass.

## Context

- instruction: `skills/idea_stage_judgment_test.go`
- instruction: `skills/roundfix_modes_test.go`
- instruction: `docs/adr/0436-an-authoring-phase-questions-a-rule-that-blocks-the-goal-and-the-rule-binds-until-the-user-answers.md`
- instruction: `docs/adr/0437-a-rule-the-user-keeps-carries-a-kept-line-and-is-asked-about-again-only-with-a-new-fact.md`
- interface: `.agents/skills/write-idea/references/idea-template.md`
- interface: `skills/write-idea/references/idea-template.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/rule_challenge_record_test.go`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The two templates, the `modes` reference and the owned-skill version record — consumers: `./skills`
- The glossary, the user guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheTemplatesCarryTheRuleChallengeRecord|TestTheSkillModesNeverAnswerARuleChallenge|TestTheFollowUpModeRunsItsStepsInOrder|TestTheUnattendedModeAsksWithinTheBaselineLimit|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheTemplatesCarryTheRuleChallengeRecord/idea/heading_and_position TestTheTemplatesCarryTheRuleChallengeRecord/idea/fields TestTheTemplatesCarryTheRuleChallengeRecord/idea/none TestTheTemplatesCarryTheRuleChallengeRecord/idea/phase TestTheTemplatesCarryTheRuleChallengeRecord/prd/heading_and_position TestTheTemplatesCarryTheRuleChallengeRecord/prd/fields TestTheTemplatesCarryTheRuleChallengeRecord/prd/none TestTheTemplatesCarryTheRuleChallengeRecord/prd/phase TestTheSkillModesNeverAnswerARuleChallenge/unattended TestTheSkillModesNeverAnswerARuleChallenge/follow_up_step_five TestTheFollowUpModeRunsItsStepsInOrder TestTheUnattendedModeAsksWithinTheBaselineLimit TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; make skills-sync-check && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the two tests do not exist, so their pass lines are missing; afterwards both copies of the templates and of the `modes` reference hold the pinned text, the existing mode tests pass, the mirrors match and the whole skills suite passes with the recorded versions.
- `grep -F -- '**Rule Challenge**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Rule Challenge\*\*(:| —( |$))' || { echo 'CONTEXT.md lacks the definition of Rule Challenge' >&2; exit 1; }; grep -F -- '**Kept Line**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Kept Line\*\*(:| —( |$))' || { echo 'CONTEXT.md lacks the definition of Kept Line' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'The dated line a rule gains when the user answers a Rule Challenge with keep' || { echo 'the Kept Line entry lacks its definition' >&2; exit 1; }; for phrase in 'In brainstorming, the Idea and the PRD, a recorded decision guides the solution and never bounds the goal in silence.' 'In a Run, in a Task and in the QA gate a rule is a guard-rail, and no Rule Challenge is asked there'; do tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/user-guide/context-driven-development.md: %s\n' "$phrase" >&2; exit 1; }; done; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither term nor the subsection exists; afterwards the glossary defines both terms, the guide holds the subsection, and the documentation contracts, the glossary contracts and the current Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals (second and fourth bullets); User Stories 4, 6 and 7; Core Features 7, 9, 10, 14 and 15; Success Metric 2; Success Metric 3; Glossary
- `_techspec.md` → Record section; Modes reference; Glossary, user guide and decision records; Data Models; Invariants 7 and 8; API Contract 6; API Contract 7; API Contract 8; API Contract 9; Testing Approach; Build Order 2; Build Order 3
- ADR-0436, ADR-0437, ADR-0403, ADR-0353

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, golang-testing, testing-boss, evidence-gate.

- Both template copies have the prescribed record, fields, empty case, phase and placement. The focused `go test ./skills -run '^(TestTheTemplatesCarryTheRuleChallengeRecord|TestTheSkillModesNeverAnswerARuleChallenge)$' -count=1` passed (12 tests and subtests).
- Both modes copies have the prescribed waiting text; the focused tests above check the phrases and placement in step 5. The existing eight-step contract is left for Daemon Verification.
- The glossary and guide contain the prescribed text. `make skills-sync`, the required owned-skill version recorder, and a second `make skills-sync` succeeded. The Behavior Surface Record recorder is invoked after this entry; documentation contracts are left for Daemon Verification.
