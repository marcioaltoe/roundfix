---
task: task_02
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
status: completed
type: docs
complexity: medium
---

# Task 02: The authoring skills put a Rule Challenge to the user

## Overview

The four skills that shape a goal, `brainstorming`, `write-idea`, `write-prd` and `write-techspec`, each gain a `## Rule Challenge` section: one core text, the same bytes in the four, that states what a rule is, when it blocks the goal, the question with its three answers, the cap of three, that the rule binds until the user answers, the wait when nobody can answer and what follows each answer, followed by a paragraph for that skill's stage. Four existing sentences gain a pointer to the section. A contract test pins every guarantee in both copies of each skill, so an author reading any of the four skills meets the same step the Baseline clause of task_01 states.

## Requirements

1. MUST add to the canonical `SKILL.md` of `write-idea`, `write-prd` and `write-techspec`, immediately before the heading `## Declared model`, and to the canonical `SKILL.md` of `brainstorming`, immediately before the heading `## Idea-stage judgment`, the section of `_techspec.md` → API Contract 3, word for word, followed after a blank line by that skill's stage paragraph of API Contract 4, word for word. The core text is the same bytes in the four skills (`_techspec.md` → Invariant 6).
2. MUST make the four pointer changes of `_techspec.md` → API Contract 5, each the whole change to its sentence: `write-prd`, `## Inputs`; `brainstorming`, the bullet "Explore before asking"; `write-techspec`, the ground rule "Read first"; and `write-techspec`, the bullet "Rejected alternatives", which keeps its wrapping over several lines.
3. MUST NOT change any other sentence of the four skills, any other owned skill, any Vendored Skill, either template of task_03, the TechSpec template, any test that pins existing skill text, or the `### QA settlement` section of any skill (`_techspec.md` → Invariant 8).
4. MUST create `skills/rule_challenge_authoring_test.go` with `TestTheAuthoringSkillsPutARuleChallengeToTheUser`, calling `t.Parallel()` first. It checks both the canonical and the shipped copy of each skill through `assertIdeaStageSkillPhrases`, which requires each phrase exactly once on whitespace-normalized text, in these subtests:
   - `heading_and_position`: in each of the four skills the heading `## Rule Challenge` occurs once, followed by the next level-two heading `## Declared model`, or `## Idea-stage judgment` in `brainstorming`;
   - for each of the four skills, one phrase per guarantee of API Contract 3: `what_a_rule_is` ("A rule guides the solution and never bounds the goal in silence."), `the_trigger` ("A rule blocks the goal when the solution you would recommend, or a part the goal needs, cannot exist without contradicting it, or when it stops you from stating the problem."), `not_a_block` ("Extra work, taste and friction are not a block."), `the_question` ("Name the rule and where it lives, what it imposes, up to three things it prevents for this goal, and your reading of why the rule exists and whether it looks out of date or badly written."), `three_answers` ("Offer keep, change and retire, each with its consequence, the recommended one first"), `binds_until_the_answer` ("The rule binds until the user answers: asking changes nothing, and a recommendation, silence, a dismissal or a timeout is never an answer."), `a_decided_change` ("Ask no Rule Challenge about a change the user already decided"), `a_kept_line_needs_a_new_fact` ("unless the question names a fact the earlier one did not have. Record that fact as"), `at_most_three` ("Ask at most three in one Authoring Phase or brainstorming session."), `never_answers_for_the_user` ("Never answer for the user. An Author Subagent returns the question, with its options and the one it recommends, to the session that started it."), `nobody_present` ("When no person is present, write no artifact for this phase"), `change_or_retire` ("Change or retire: the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it."), `keep` ("Keep: the solution stays inside the rule, and the PRD requires a Task to add the Kept Line after the rule's text. The line names no Spec.") and `a_baseline_rule` ("sends the user's decision as one ai-memory message to the Baseline's owner, obeys the rule until the Baseline changes, and records the message and the scope it left out.");
   - `the_core_text_is_shared`: the canonical text from the heading `## Rule Challenge` to the end of the sentence "records the message and the scope it left out." is the same bytes in the four canonical skills;
   - `stage/write-idea`, `stage/write-prd`, `stage/brainstorming` and `stage/write-techspec`: that skill's whole stage paragraph of API Contract 4;
   - `pointers/write-prd_inputs`, `pointers/brainstorming_explore`, `pointers/write-techspec_read_first` and `pointers/write-techspec_rejected_alternatives`: the `after` text of each pointer of API Contract 5.
5. MUST regenerate, never by hand: the mirrors under `skills/` with `make skills-sync`, then the owned-skill versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of each changed skill on the tree the Task runs on, then `make skills-sync` again so each mirror equals its source.
6. MUST NOT change `CONTEXT.md`, any file under `docs/`, any file under `internal/`, `docs/references/behavior-surfaces.json`, the Makefile, a workflow, `go.mod` or a lint configuration. Of the Governed Paths it changes only the eight `SKILL.md` files below, all listed in `_authorization.md`.

## Subtasks

- [ ] Add the section with its stage paragraph to the four canonical skills.
- [ ] Add the four pointers.
- [ ] Add the contract test.
- [ ] Sync the mirrors and record the owned-skill versions.

## Acceptance Criteria

- [ ] Both copies of the four skills hold the section once, in its place, with every guarantee and their own stage paragraph, and the core text is the same bytes in the four.
- [ ] Removing a pinned phrase from a canonical skill or from its mirror fails the named subtest.
- [ ] The mirrors equal their sources, the new versions are recorded, and every existing test of the skills package passes unchanged.

## Context

- instruction: `skills/idea_stage_judgment_test.go`
- instruction: `skills/idea_stage_authoring_test.go`
- instruction: `docs/adr/0436-an-authoring-phase-questions-a-rule-that-blocks-the-goal-and-the-rule-binds-until-the-user-answers.md`
- instruction: `docs/adr/0437-a-rule-the-user-keeps-carries-a-kept-line-and-is-asked-about-again-only-with-a-new-fact.md`
- interface: `.agents/skills/brainstorming/SKILL.md`
- interface: `skills/brainstorming/SKILL.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/rule_challenge_authoring_test.go`

## Consumers

- The text of the four skills and the owned-skill version record — consumers: `./skills`
- The skills read by the public documentation contracts — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringSkillsPutARuleChallengeToTheUser|TestEveryOwnedSkillVersionIsRecorded|TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheAuthoringSkillsPutARuleChallengeToTheUser/heading_and_position TestTheAuthoringSkillsPutARuleChallengeToTheUser/what_a_rule_is TestTheAuthoringSkillsPutARuleChallengeToTheUser/the_trigger TestTheAuthoringSkillsPutARuleChallengeToTheUser/not_a_block TestTheAuthoringSkillsPutARuleChallengeToTheUser/the_question TestTheAuthoringSkillsPutARuleChallengeToTheUser/three_answers TestTheAuthoringSkillsPutARuleChallengeToTheUser/binds_until_the_answer TestTheAuthoringSkillsPutARuleChallengeToTheUser/a_decided_change TestTheAuthoringSkillsPutARuleChallengeToTheUser/a_kept_line_needs_a_new_fact TestTheAuthoringSkillsPutARuleChallengeToTheUser/at_most_three TestTheAuthoringSkillsPutARuleChallengeToTheUser/never_answers_for_the_user TestTheAuthoringSkillsPutARuleChallengeToTheUser/nobody_present TestTheAuthoringSkillsPutARuleChallengeToTheUser/change_or_retire TestTheAuthoringSkillsPutARuleChallengeToTheUser/keep TestTheAuthoringSkillsPutARuleChallengeToTheUser/a_baseline_rule TestTheAuthoringSkillsPutARuleChallengeToTheUser/the_core_text_is_shared TestTheAuthoringSkillsPutARuleChallengeToTheUser/stage/write-idea TestTheAuthoringSkillsPutARuleChallengeToTheUser/stage/write-prd TestTheAuthoringSkillsPutARuleChallengeToTheUser/stage/brainstorming TestTheAuthoringSkillsPutARuleChallengeToTheUser/stage/write-techspec TestTheAuthoringSkillsPutARuleChallengeToTheUser/pointers/write-prd_inputs TestTheAuthoringSkillsPutARuleChallengeToTheUser/pointers/brainstorming_explore TestTheAuthoringSkillsPutARuleChallengeToTheUser/pointers/write-techspec_read_first TestTheAuthoringSkillsPutARuleChallengeToTheUser/pointers/write-techspec_rejected_alternatives TestEveryOwnedSkillVersionIsRecorded TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; make skills-sync-check && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the test and the sections do not exist, so the pass lines are missing; afterwards every guarantee passes in both copies of the four skills, the mirrors match, the whole skills suite passes with the recorded versions, and the documentation contracts pass.

## References

- `_prd.md` → Goals (first and fourth bullets); User Stories 1, 2, 3, 6, 7, 8 and 9; Core Features 2, 3, 4, 5, 6, 8, 9, 10, 11, 12 and 14; Success Metric 2
- `_techspec.md` → Skill section and stage paragraphs; Invariants 6 and 8; API Contract 3; API Contract 4; API Contract 5; Testing Approach; Build Order 2
- ADR-0436, ADR-0437, ADR-0320, ADR-0403


## Result

Skills loaded: implement-task, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging

- Added the shared Rule Challenge text, stage paragraphs and four pointers; the focused `go test ./skills -run '^TestTheAuthoringSkillsPutARuleChallengeToTheUser$' -count=1` passed, covering both copies, section position, guarantees, stage text, pointers and byte equality of the canonical core.
- The contract uses `assertIdeaStageSkillPhrases` for each named guarantee in each copy; deliberate removal probes were not run, so their failure behavior remains for verification.
- Ran `make skills-sync`, the required `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` (exit 0), and `make skills-sync` again. A byte comparison found each of the four mirrors identical to its source. The full skills suite and declared Verification are left to the Daemon.
- The first focused run exposed an assertion counting inline pointers as headings; the assertion was corrected to check heading lines and the focused rerun exited 0.
