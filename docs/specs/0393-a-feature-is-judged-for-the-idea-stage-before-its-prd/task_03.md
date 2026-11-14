---
task: task_03
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
status: completed
type: docs
complexity: medium
---

# Task 03: The authoring skills and the follow-up mode run and record the Idea-stage judgment

## Overview

`brainstorming` and `write-prd` judge a feature for the Idea stage before its
PRD, `write-idea` shows the judgment and offers the dismissal when the rule
starts it, both templates carry the fielded record, and the Roundfix skill's
follow-up mode runs the judgment for each planned feature and leaves a feature
that needs the user waiting. The glossary gains **Idea Stage Judgment** and
the changed **Advisory Judgment**, and ADR-0201 gains its amendment notice. A
new skills test pins each rule with its own phrase in both copies of every
changed file.

## Requirements

1. MUST change the canonical `brainstorming`, `write-idea` and `write-prd`
   skills and the Roundfix `modes` reference as `_techspec.md` → API Contract
   6 says for those files, and the two templates as API Contract 7 says, word
   for word and at the places named, and copy each to its mirror with `make
   skills-sync`. MUST NOT change any other byte of those files: the grouping
   sentences of Spec 0391, the research-subagent text of Spec 0392 and the
   `### QA settlement` section of any skill stay as the Task finds them.
2. The skill texts state two hazards and MUST keep them word for word: of a
   request of the session only the project's name, its first line and its
   first 4,000 characters reach the judge, after the secret scan, and the
   session removes the file it wrote; and the agent's own answers are never
   presented as the judge's, nor does any agent or delegated subagent answer
   the Idea stage's questions or an undecided judgment for the user.
3. MUST add `TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage` and
   `TestTheFollowUpModeRunsTheIdeaStageJudgment` in a new
   `skills/idea_stage_authoring_test.go`, each calling `t.Parallel()` first and
   reading both copies of each file through the helper task_01 added in
   `skills/idea_stage_judgment_test.go`, with one named subtest for each rule
   `_techspec.md` → Testing Approach lists for `skills`: the three criteria,
   the three verdict rules, the request file with what reaches the judge,
   "Never present your own answers as the judge's", the dismissal option, the
   waiting rule, every key of the record in both templates, and the four
   sentences of the `modes` reference.
4. MUST raise the versions of the four skills only with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   and never write a version number.
5. MUST write the glossary through `domain-modeling` as API Contract 8 says,
   word for word: the changed **Advisory Judgment** and the new **Idea Stage
   Judgment** entry right after it, each ADR written as the glossary's
   Markdown link.
6. MUST add to ADR-0201 the amendment notice of API Contract 9, word for
   word, at the end of the file, and set its `updated_at` as that contract
   says; nothing else in the ADR changes.
7. MUST NOT change any Go source outside `skills/idea_stage_authoring_test.go`,
   any Baseline module, any guide under `docs/agents/` or `docs/user-guide/`,
   or the Roundfix `spec` reference. Of the Governed Paths the Task may change
   only the canonical files and the `SKILL.md` mirrors of the four skills and
   the mirror of the PRD template, as `_authorization.md` lists them.

## Subtasks

- [ ] Add the judgment sections to brainstorming, write-idea and write-prd, and the record to both templates.
- [ ] Add the four sentences to the follow-up and unattended modes.
- [ ] Sync the mirrors, add the two tests and record the skill versions.
- [ ] Write the two glossary entries and the amendment notice in ADR-0201.

## Acceptance Criteria

- [ ] Both copies of brainstorming, write-idea and write-prd tell the agent to judge three criteria before a feature's PRD, to hand a request of the session to the judge as a file whose start alone is sent, to answer the criteria itself and say so when the judge does not, to start an indicated stage with a dismissal option, to ask the user for an undecided one, and to record the judgment in the fielded form.
- [ ] Both copies of the modes reference run the judgment in the plan, ask the Idea stage's questions while the user is present, and leave an unanswered feature waiting as a Backlog Entry with its Waiting Declaration.
- [ ] The glossary defines Idea Stage Judgment and the changed Advisory Judgment, ADR-0201 carries its notice, and the glossary, skill and documentation contracts pass.

## Context

- instruction: `docs/adr/0403-a-feature-is-judged-for-the-idea-stage-before-its-prd-and-the-judgment-never-refuses.md`
- instruction: `docs/adr/0404-the-idea-stage-judge-reads-one-item-that-is-not-yet-a-spec.md`
- interface: `.agents/skills/brainstorming/SKILL.md`
- interface: `skills/brainstorming/SKILL.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-idea/references/idea-template.md`
- interface: `skills/write-idea/references/idea-template.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `skills/roundfix_layout_test.go`
- creates: `skills/idea_stage_authoring_test.go`
- interface: `CONTEXT.md`
- interface: `docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md`

## Consumers

- The four owned skills, their templates and the owned-skill version record — consumers: `./skills`
- The glossary and the amended decision record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage|TestTheFollowUpModeRunsTheIdeaStageJudgment|TestTheSpecReferenceDocumentsTheIdeaStage|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage TestTheFollowUpModeRunsTheIdeaStageJudgment; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./skills` — expected: exit 0. Before this Task the two tests do not exist, so the first pass line is missing; afterwards both copies of every changed skill file carry each rule once, the mirrors match, the versions are recorded and the skills package passes.

- `grep -F -- '**Idea Stage Judgment**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Idea Stage Judgment\*\*(:| —( |$))' && grep -F -- '**Advisory Judgment**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Advisory Judgment\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'or about an item that is not yet a Spec, such as whether a feature needs the Idea stage, that the authoring model must answer' && tr -s '[:space:]' ' ' < docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md | grep -qF -- 'Amended by [ADR-0404](0404-the-idea-stage-judge-reads-one-item-that-is-not-yet-a-spec.md) (2026-11-14).' && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the glossary has no Idea Stage Judgment entry, so the first check fails; afterwards both terms are defined, ADR-0201 carries its notice, and the documentation contracts, the glossary's own rules included, pass.

## References

- `_prd.md` → User Story 2; User Story 3; User Story 4; User Story 6; User Story 7; Core Feature 6; Core Feature 7; Core Feature 8; Core Feature 9; Core Feature 10; Success Metric 6; Success Metric 7; Glossary
- `_techspec.md` → Current behavior 6 and 7; Data Models (the record); API Contract 6; API Contract 7; API Contract 8; API Contract 9; Declared breaks; Testing Approach; Build Order 3
- ADR-0403; ADR-0404; ADR-0201

## Result

Skills loaded: implement-task, coding-guidelines, testing-boss, golang-testing, tech-writer, domain-modeling, evidence-gate

- Added the Idea-stage judgment contract to the canonical `brainstorming`, `write-idea` and `write-prd` skills, including the request-file boundary, fail-open agent answers, verdict routing, dismissal, waiting, and fielded-record rules.
- Added the fielded `## Idea-stage judgment` record to both templates and the four follow-up/unattended mode sentences. `make skills-sync` completed, and every edited canonical/mirror pair compares equal.
- Added `TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage` and `TestTheFollowUpModeRunsTheIdeaStageJudgment`; both call `t.Parallel()` first and use the task_01 helper through named rule subtests.
- Updated `CONTEXT.md` with the changed Advisory Judgment and the new Idea Stage Judgment, and appended ADR-0404's amendment notice to ADR-0201 with the contracted timestamp.
- Verification feedback repair: traced the failure to two doubled-backtick test phrases and one assertion wording mismatch, then corrected only `skills/idea_stage_authoring_test.go`.
- Focused post-repair check: `gofmt -w skills/idea_stage_authoring_test.go && go test -count=1 ./skills -run '^(TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage|TestTheFollowUpModeRunsTheIdeaStageJudgment)$'` passed (`Go test: 9 passed in 1 packages`). Declared Verification remains for the Daemon.

## Carry-forward provenance

- Source Run: `run_20261114T210733Z_b8e5ccd1f4543b30`
- Source commit: `799a903fcdcd351fc992184cc7baba9c890f4072`
