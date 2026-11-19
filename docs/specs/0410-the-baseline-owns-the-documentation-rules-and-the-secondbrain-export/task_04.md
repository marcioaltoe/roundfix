---
task: task_04
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: pending
type: docs
complexity: medium
---

# Task 04: The skills ask the Lesson Learned question and the update marks the existing ADRs

## Overview

The two moments of a lesson reach the skills that run them: `archive-spec` reads five signs before the archive and applies the Lesson Learned test, and the Roundfix skill's `corrections` reference applies it after a declared Correction or a retry that follows a repair. The `update` Skill Mode shows the documentation part, asks before a plan that stops publishing a path, and has the session's Agent mark the existing ADRs in one reviewed commit. The user guide explains what a repository publishes, and one contract test pins every sentence the skills gained in this Spec, in both copies.

## Requirements

1. MUST add to the canonical `archive-spec` skill, before `## Steps`, the section `## Lesson Learned` of `_techspec.md` → API Contract 14, word for word, under its own heading; no text inside a `### QA settlement` section of any skill changes.
2. MUST add to the end of the canonical Roundfix `corrections` reference the section `## Lesson Learned` of API Contract 15, word for word.
3. MUST change the section `## --update` of the canonical Roundfix `modes` reference as API Contract 16 says, inserting each text by its anchor as the section reads on the default branch when the Task runs: step 3 gains its sentence after the sentence that ends "and its `planDigest`.", step 4 gains its sentence after the sentence that ends "ask the same single question before applying either.", and the paragraph on marking the existing ADRs follows step 5. The section keeps five numbered steps and every sentence it holds at the Delivery Base.
4. MUST add to `docs/user-guide/context-driven-development.md`, before `### Spec Event Announcements`, the subsection `### Knowledge a repository publishes` of API Contract 20: it describes the Standard Export List, the Export Variation, the Knowledge Mark, the Amendment List and the Lesson Learned, shows the export file of the PRD's first illustration (`_prd.md` → User Experience) and the lesson form, and holds the sentences "A Lesson Learned is written only when there is a lesson, and it leaves when a rule takes it in.", "An ADR without the Knowledge Mark is not published." and "The marks of the ADRs a repository already has are one change a person reviews."
5. MUST create `skills/documentation_rules_steps_test.go` with `TestTheSkillsCarryTheDocumentationSteps`, which calls `t.Parallel()` first and checks each phrase on whitespace-normalized text, exactly once, in both copies of the file named (`.agents/skills/...` and `skills/...`); removing the phrase from either copy fails its subtest (`_prd.md` → Success Metric 7; `_techspec.md` → API Contracts 14 to 19):
   - `archive/before_steps`: `## Lesson Learned` comes before `## Steps` in `archive-spec`'s `SKILL.md`;
   - `archive/five_signs`: the five bullets of API Contract 14, each by its own phrase: "adopts a Finding;", "holds a QA Report whose verdict is not", "the Task Graph holds a Corrective Task;", "lists a declaration other than" and "a Task of the Run failed its Verification and ran again, or the Delivery Queue item was retried.";
   - `archive/no_sign_nothing`: "With no sign, write nothing.";
   - `archive/commit_with_archive`: "commit it with the archive, so the same Pull Request shows it to a person, who may delete it.";
   - `archive/report_the_failed_part`: "say in the report that the test was applied and which part failed.";
   - `archive/never_fill`: "Never write a lesson to fill the folder.";
   - `archive/queue_defers`: "When the Delivery Queue writes the archive, no Agent runs this step there";
   - `archive/skip_without_rule` and `corrections/skip_without_rule`: "A repository whose docs-layout guide has no Lesson Learned rule skips this section.";
   - `corrections/cheapest_moment`: "A Correction is where a delivery's lesson is cheapest to state, because its declarer has just named the cause.";
   - `corrections/in_the_correcting_commit`: "in the commit that makes the correction, on the Item Branch, so the delivery Pull Request shows it.";
   - `corrections/after_merge`: "When the item has already merged, or the queue made the correction with no session, write the file in the session's next change to the repository.";
   - `modes/step3_documentation`: "member too: the lines the Standard Export List gains and loses";
   - `modes/step4_stops_exporting`: "Ask that same question, naming the paths, when";
   - `modes/marking`: "mark each ADR that lacks the mark by that rule's test, changing no other line.";
   - `modes/review_in_pr`: "say that the marks are a judgment for a person to review in the Pull Request, in a commit that holds only the marks and the list.";
   - `modes/auto_marks_nothing`: "mark nothing and report the count.";
   - `modes/five_steps`: the `## --update` section holds the numbered steps 1 to 5 and no step 6;
   - `baseline/documentation_part`: "part of its plan that the Plan Digest covers and the same confirmation applies.";
   - `baseline/marks_none`: "It reports the count of ADRs without a Knowledge Mark and marks none.";
   - `baseline/marker_refused`: "A missing or repeated marker is refused: nothing is written to that file, and the rest of the plan applies.";
   - `baseline/no_history`: "does not skip this part.";
   - `setup/two_checks`: "reports two documentation checks as warnings that never change its exit.";
   - `setup/repairs`: "repairs the list, an Amendment List and a retired ADR, and the rest is an edit by hand.";
   - `release/sentence`: "are planning evidence too, and the export list".
6. MUST regenerate, never by hand: `make skills-sync`; the owned-skill versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which raises both version fields of both copies of the `SKILL.md` of `archive-spec` and of the Roundfix skill on this Task's tree; and, after the guide edit, the Behavior Surface Record with the `-record-skill-coverage` step, `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST NOT change any Go file outside `skills/documentation_rules_steps_test.go`, any Baseline module, template, profile, golden or Source Baseline, any guide under `docs/agents/`, any ADR, `.secondbrain-export`, `CONTEXT.md`, any skill other than `archive-spec` and the Roundfix skill's `corrections` and `modes` references and version fields, any Vendored Skill, `skills/baseline_skill_contract_test.go`, `docs/references/skill-coverage.json`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration. Of the Governed Paths it changes only `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/references/corrections.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, all listed in `_authorization.md` `paths`.

## Subtasks

- [ ] Add the Lesson Learned section to `archive-spec` and to the `corrections` reference.
- [ ] Add the documentation part, the question and the marking paragraph to the `update` Skill Mode.
- [ ] Add the user-guide subsection.
- [ ] Add the skill contract test; sync the skills, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of `archive-spec` carry the lesson step with its five signs, and both copies of the `corrections`, `modes`, `baseline`, `setup` and `release` references carry their sentences, each pinned by its own subtest that fails when the phrase is removed from either copy.
- [ ] The `--update` section keeps five steps and every sentence it held at the Delivery Base.
- [ ] The user guide describes what a repository publishes, before the Spec Event Announcements.

## Context

- instruction: `docs/agents/docs-layout.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/references/corrections.md`
- interface: `skills/roundfix/references/corrections.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/documentation_rules_steps_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `archive-spec` skill and the Roundfix skill's `corrections` and `modes` references, read by the skill contract tests — consumers: `./skills`
- The user guide page and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheSkillsCarryTheDocumentationSteps$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in archive/before_steps archive/five_signs archive/no_sign_nothing archive/commit_with_archive archive/report_the_failed_part archive/never_fill archive/queue_defers archive/skip_without_rule corrections/skip_without_rule corrections/cheapest_moment corrections/in_the_correcting_commit corrections/after_merge modes/step3_documentation modes/step4_stops_exporting modes/marking modes/review_in_pr modes/auto_marks_nothing modes/five_steps baseline/documentation_part baseline/marks_none baseline/marker_refused baseline/no_history setup/two_checks setup/repairs release/sentence; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheSkillsCarryTheDocumentationSteps/$name " || { echo "missing pass: TestTheSkillsCarryTheDocumentationSteps/$name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && make skills-sync-check` — expected: exit 0. Before this Task the test does not exist; afterwards every sentence the skills gained in this Spec passes in both copies, the skills package passes whole, the owned-skill versions are recorded and the mirrors match.
- `file=docs/user-guide/context-driven-development.md; for phrase in 'A Lesson Learned is written only when there is a lesson, and it leaves when a rule takes it in.' 'An ADR without the Knowledge Mark is not published.' 'The marks of the ADRs a repository already has are one change a person reviews.'; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; knowledge="$(grep -n '^### Knowledge a repository publishes$' "$file" | cut -d: -f1)"; events="$(grep -n '^### Spec Event Announcements$' "$file" | cut -d: -f1)"; test -n "$knowledge" && test -n "$events" && test "$knowledge" -lt "$events" || { echo 'the subsection is missing or after Spec Event Announcements' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide has no such subsection; afterwards it does, before the Spec Event Announcements, and the documentation contracts, the Behavior Surface Record and the glossary included, pass.

## References

- `_prd.md` → Goals; User Stories 6, 8, 9 and 13; Core Features 7, 13, 14, 16 and 20; Success Metric 7; Success Metric 8
- `_techspec.md` → Skills, guides and glossary; API Contract 14; API Contract 15; API Contract 16; API Contract 17; API Contract 18; API Contract 19; API Contract 20; Testing Approach; Build Order 4
- ADR-0440, ADR-0417, ADR-0420, ADR-0256
