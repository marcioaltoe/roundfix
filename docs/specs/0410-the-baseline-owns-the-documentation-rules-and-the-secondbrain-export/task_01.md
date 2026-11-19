---
task: task_01
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: pending
type: docs
complexity: high
---

# Task 01: The Baseline states the documentation rules and the export rule

## Overview

The Baseline's docs-layout rule gains the Amendment List, the exception to the lifecycle rule, the folder variation and the retirement of a Lesson Learned under the ids of four existing clauses; the cross-project rule stops sending a closing Finding's lesson as a message; and two new rules carry the **Lesson Learned** and the **Standard Export List** with the **Export Variation** and the **Knowledge Mark**. A contract test pins every guarantee of the seven texts by its own phrase and their rendering in the three guides, the module records and derived Baseline files are regenerated, this repository's three guides and Setup Manifest are refreshed, and `CONTEXT.md` defines the five new terms and revises **Managed Refresh**.

## Requirements

1. MUST change `internal/baseline/assets/modules/context-workflow.json` as `_techspec.md` → API Contracts 1 to 5 prescribe, word for word, each `guidance` one JSON string with no line break, keeping the file's two-space JSON formatting:
   - `clause.context.adr-01-template` keeps its text and its code block and gains, after the block, the paragraph of API Contract 1;
   - `clause.context.adr-03-legacy-compatibility` takes the text of API Contract 2, which keeps the sentence "Do not rewrite existing ADRs solely to adopt lifecycle metadata." whole, with its period;
   - in `clause.context.docs-one-job-per-directory` everything before "Preserve repository-authored bytes outside setup markers" is replaced by the text of API Contract 3, and the rest of the clause is unchanged;
   - in `clause.context.findings-09-archive` the phrase of API Contract 4 is replaced; the clause still opens with "A retirement deletes the file." and keeps its `replaces` list unchanged;
   - the module gains the rule `rule.context.lessons`, `version` 1, `coverage` `["coverage.context-workflow"]`, with the two clauses of API Contract 5 in this order, `clause.context.lesson-01-record` and `clause.context.lesson-02-test`, each `enforcement` `mandatory` and neither declaring `replaces`; the lesson form of the first is a `markdown` block inside its `guidance`, with the frontmatter `created: YYYY-MM-DD`, the title line and the eight labelled lines in the order API Contract 5 gives;
   - the supporting guide `guide.docs-layout` lists `rule.context.lessons` right after `rule.context.docs-layout`, and its guide `version` rises by one from the value the Task finds.
2. MUST change `clause.core.message-the-owning-project` in `internal/baseline/assets/modules/core.json`: the two sentences from "When a Finding closes, state whether" through "pinned as `<repository>@<commit>:<path>`." become the sentence of API Contract 6, and every other sentence of the clause stays.
3. MUST add to `internal/baseline/assets/modules/secondbrain.json` the rule `rule.secondbrain.knowledge-export`, `version` 1, `coverage` `["coverage.secondbrain"]`, with the three clauses of API Contract 7 in this order, `clause.secondbrain.export-01-standard-list`, `clause.secondbrain.export-02-variation` and `clause.secondbrain.export-03-knowledge-mark`, each `mandatory` with no `replaces`; the supporting guide `guide.secondbrain` lists it last, after `rule.secondbrain.hermes-escalation`, and its guide `version` rises by one.
4. MUST NOT list `rule.context.lessons` or `rule.secondbrain.knowledge-export` under `requiredRules` in any profile, and MUST NOT change any other rule, clause, guide or field of the three modules (`references/2026-11-19-prototype-and-governed-paths.md`, section 2: a new clause in a required rule would need a hand-written Source Baseline row).
5. MUST add the five clause ids of API Contract 9 to the closed list of clause forces in `internal/baseline/clause_characterization_test.go`, each `mandatory`, in the list's order, and change no other line of that file.
6. MUST update the three existing tests that pin text this Task changes, and no other existing test (`_techspec.md` → Testing Approach):
   - `TestTheCaptureTriggerAsksForAConditionalLesson` in `internal/baseline/retirement_rule_clauses_test.go` takes the phrases "send no message about what it taught" and "Lesson Learned" in place of "ADR, Archive Record or reference", "Only then send one message" and the pinned-path form, and keeps "When a Finding closes";
   - the constant `docsOneJobPerDirectoryGuidance` in `internal/baseline/repository_rule_placement_test.go` takes the text of API Contract 3 with the unchanged rest of the clause;
   - `TestTheDocsLayoutRegionStaysLean` in `internal/baseline/retirement_rule_clauses_test.go` holds the managed region `guide.docs-layout` of the Standard TypeScript golden to at most 13,000 bytes (API Contract 8).
   `TestADRLifecycleContract` in `internal/baseline/plan_test.go`, a Governed Path this Spec does not grant, MUST stay byte-identical and pass.
7. MUST create `internal/baseline/documentation_rules_clauses_test.go` with two tests, each calling `t.Parallel()` first, that read the module files as `internal/baseline/retirement_rule_clauses_test.go` does. Every phrase is compared on whitespace-normalized text and is required exactly once in the `guidance` named.
   - `TestTheBaselineModulesCarryTheDocumentationRules`, with these subtests:
     - `the_lesson_rule`: `rule.context.lessons` exists once with version 1, the coverage `coverage.context-workflow` and exactly the two clauses of requirement 1, with those ids, forces and order and no `replaces`; `guide.docs-layout` lists it right after `rule.context.docs-layout`.
     - `the_export_rule`: `rule.secondbrain.knowledge-export` exists once with version 1, the coverage `coverage.secondbrain` and exactly the three clauses of requirement 3, with those ids, forces and order and no `replaces`; `guide.secondbrain` lists it last.
     - `required_by_no_profile`: no profile under `internal/baseline/assets/profiles/` lists either new rule under `requiredRules`.
     - `names_no_record`: none of the seven new or reworded clause texts holds `ADR-` followed by a digit or `Spec ` followed by four digits (ADR-0186).
     - `adr/amended_stays_accepted` (`clause.context.adr-01-template`): "An accepted ADR that a later accepted ADR changes in part, without replacing it, stays `accepted`."; the phrase occurs after the clause's last code fence.
     - `adr/the_later_adr_writes_it`: "In the change that adds the later ADR, its header gains `amended_by: [ADR-NNNN]`"
     - `adr/dated_notice`: "a dated notice that begins `Amended by ADR-NNNN (YYYY-MM-DD)` and says what changed."
     - `adr/the_update_fills`: "`roundfix baseline update` fills the list from the notices an ADR already carries."
     - `legacy/sentence_kept` (`clause.context.adr-03-legacy-compatibility`): "Do not rewrite existing ADRs solely to adopt lifecycle metadata."
     - `legacy/the_exception`: "The one exception is a header field that a guide of this repository defines for existing ADRs, such as `amended_by`: it is written into them and changes no other line"
     - `legacy/stays_legacy`: "a header that holds no lifecycle field leaves its ADR a legacy ADR."
     - `layout/design` (`clause.context.docs-one-job-per-directory`): "`docs/design/` for design artifacts of the product, such as an interface, API or data design, and never a decision or a plan"
     - `layout/lessons`: "`docs/lessons/` for Lessons Learned"
     - `layout/root_reference`: "including one that would otherwise sit at the repository root, the `DESIGN.md` design contract excepted"
     - `layout/variation`: "A folder under `docs/` or a reference document at the root that this clause does not name is a variation: name its path and its job in `docs/agents/specific-repository.md`, or move its content to the directory whose job it is."
     - `layout/doctor_reports`: "`roundfix doctor` reports one that is neither."
     - `layout/rest_kept`: "Preserve repository-authored bytes outside setup markers, and add a repository rule only to `docs/agents/specific-repository.md`."
     - `retirement/lesson_deleted` (`clause.context.findings-09-archive`): "and a Lesson Learned once an accepted ADR, a guide or a skill states what it teaches are deleted"; the `guidance` still begins with "A retirement deletes the file." and `replaces` is not empty.
     - `message/no_message` (`clause.core.message-the-owning-project`): "When a Finding closes, send no message about what it taught: the docs-layout guide says when that becomes a Lesson Learned in this repository."
     - `lesson/record_path` (`clause.context.lesson-01-record`): "It is one file, `docs/lessons/YYYY-MM-DD-<kebab-slug>.md`, in this complete copyable form"
     - `lesson/labels_as_written`: "keep each label as written and write the rest in the language of the repository's documentation"
     - `lesson/eight_lines`: the `guidance` holds `created: YYYY-MM-DD` and the eight lines `- Context:`, `- Tried:`, `- Went wrong:`, `- Cause:`, `- Solved by:`, `- Applies when:`, `- Does not apply when:` and `- Links:`, once each, in that order.
     - `lesson/the_test` (`clause.context.lesson-02-test`): "Write a Lesson Learned only when every part of this test holds: something was tried, it went wrong, the cause is proven, something solved it, and no accepted ADR, reference document or glossary entry already states it."
     - `lesson/otherwise_nothing`: "Otherwise write nothing."
     - `lesson/two_moments`: "when a Finding closes, in the change that answers it or retires it, and when a Spec that needed a correction, a failed QA or a retry is archived."
     - `lesson/in_that_change`: "Write the file in that change, so the Pull Request that carries it shows the lesson to a person, who may delete it."
     - `export/two_parts` (`clause.secondbrain.export-01-standard-list`): "This repository publishes to the Secondbrain only what `.secondbrain-export` lists, and that file has two parts."
     - `export/markers`: "between the lines `# roundfix:standard-export:begin` and `# roundfix:standard-export:end`"
     - `export/standard_items`: "the README, the glossary files, `docs/user-guide/`, `docs/references/`, `docs/design/`, `docs/lessons/` and each ADR whose header says `knowledge: true`, one path per line."
     - `export/baseline_owns`: "The Baseline owns those lines: `roundfix baseline update` rewrites them, and the only line added by hand is the path of an ADR, in order, when the ADR is written with that mark."
     - `export/never_process`: "Never publish process: `AGENTS.md`, `CLAUDE.md`, `docs/agents/`, the Spec folders, `docs/backlog/`, `docs/findings/`, `docs/handoffs/` and `docs/history/` stay out of both parts."
     - `variation/own_lines` (`clause.secondbrain.export-02-variation`): "The second part is the repository's own lines, everything outside those two markers, which no Baseline operation changes."
     - `variation/exclusion`: "or drop a standard item with an exclusion line: `!` followed by the item as the list writes it, or `!docs/adr/*` for every ADR."
     - `variation/reason`: "Give each such line its reason in a comment line directly above it; one comment stands for the lines that follow it without a blank line."
     - `variation/doctor`: "`roundfix doctor` reports a line without a reason, a process path, and a standard list that differs from what the update would write."
     - `mark/add_it` (`clause.secondbrain.export-03-knowledge-mark`): "add `knowledge: true` or `knowledge: false` to its header, on the line after `status`."
     - `mark/the_test`: "It is `true` when the decision states a rule of the business or of the product that someone outside this repository's development can use, and `false` when it is about how this repository is built, tested, organized or delivered, or about a tool it uses."
     - `mark/unmarked_unpublished`: "An ADR without the mark is not published."
     - `mark/existing_once`: "The ADRs a repository already has are marked once, by the same test, by the Agent of the session that applies the Baseline update that brings this rule"
     - `mark/one_reviewed_commit`: "in one commit that a person reviews in a Pull Request, followed by a second update so the list names them."
   - `TestTheGuidesRenderTheDocumentationRulesOnce`, with three subtests. Each asserts that the whole `guidance` of each of the seven new or reworded clauses, whitespace-normalized and read from the module file, occurs exactly once in the guide that renders it (`docs/agents/docs-layout.md` for the five clauses of `context-workflow`, `docs/agents/agent-instructions.md` for the one of `core`, `docs/agents/secondbrain.md` for the three of the export rule), the two lesson clauses after the last clause of `rule.context.docs-layout`, and the three export clauses after the last clause of `rule.secondbrain.hermes-escalation`:
     - `plan_postimage`: in the postimages of a plan built with `newProjectDecisionPlanRepository`, `buildProjectDecisionPlan`, `standardTypeScriptDecisions` and `planPostimage`;
     - `formatter_golden`: in the three guides under `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/`;
     - `repository_guide`: in this repository's three guides.
8. MUST add to `CONTEXT.md`, after the entry of **Managed Refresh** and each followed by a blank line, the five entries of `_techspec.md` → API Contract 21 word for word, **Standard Export List**, **Export Variation**, **Knowledge Mark**, **Amendment List** and **Lesson Learned**, through the `domain-modeling` skill; and add to the entry of **Managed Refresh**, after the sentence that ends "converts Pending History on approval.", the sentence "It also keeps the Standard Export List and fills each Amendment List, in a documentation part of its plan that the Plan Digest covers." No other entry of `CONTEXT.md` changes.
9. MUST regenerate, in this order and never by hand: the module versions with `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`, which writes the next version of `context-workflow`, `core` and `secondbrain` (this Task names no number); the derived Baseline files with `make baseline-digests`; and, as the Task's last step, this repository's three guides and Setup Manifest with `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, run again until a run prints `File changes: 0`. This Task edits no skill, so `make skills-sync` is not needed before `make baseline-digests`. The files the prototype saw those commands rewrite (`references/2026-11-19-prototype-and-governed-paths.md`, section 4) are declared below; any other file they rewrite on the tree the Task runs on is recorded by the Daemon.
10. MUST keep every clause that exists at the Delivery Base with its id and its force, and the text of every clause other than the five this Task rewords; every test of `./internal/baseline` that exists at the Delivery Base passes, the three of requirement 6 updated (`_prd.md` → Success Metric 9).
11. MUST NOT change any skill, any file under `docs/user-guide/`, `.secondbrain-export`, any ADR, any Go file outside `_test.go` files, any Source Baseline or Setup Snapshot under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/`, any other module, template or profile source, any guide under `docs/agents/` other than `agent-instructions.md`, `docs-layout.md`, `secondbrain.md` and `setup-context.json`, `AGENTS.md`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration. Of the Governed Paths it changes only the three modules by hand, and the Standard TypeScript profile, its three formatter goldens, this repository's three guides and `docs/agents/setup-context.json` only through the commands of requirement 9; all eleven are listed in `_authorization.md` `paths`.

## Subtasks

- [ ] Reword the four docs-layout clauses and the message clause, and add the lesson rule and the export rule with their guide entries.
- [ ] Add the five ids to the closed list of clause forces and update the three tests that pin the changed text.
- [ ] Add the clause and rendering tests.
- [ ] Add the five glossary entries and the sentence of Managed Refresh.
- [ ] Record the module versions, regenerate the derived files and refresh this repository's guides and Setup Manifest until a refresh changes nothing.

## Acceptance Criteria

- [ ] A plan for the Standard TypeScript profile renders the seven new or reworded clauses exactly once in their guides, the lesson clauses after the docs-layout clauses and the export clauses last in the Secondbrain guide; the formatter goldens and this repository's guides hold them the same way.
- [ ] Each guarantee of API Contracts 1 to 7 is pinned by its own named subtest, and removing its phrase from the module fails that subtest.
- [ ] Every clause of the Delivery Base keeps its id and force, the docs-layout region is at most 13,000 bytes, and a further Managed Refresh reports this repository current.
- [ ] `CONTEXT.md` defines the five new terms and the revised Managed Refresh, and the glossary contracts pass.

## Context

- instruction: `docs/adr/0438-a-repository-publishes-a-standard-export-list-the-baseline-keeps-and-varies-it-only-in-its-own-lines.md`
- instruction: `docs/adr/0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md`
- instruction: `docs/adr/0440-a-lesson-learned-is-written-only-when-a-test-holds-and-leaves-when-a-rule-states-it.md`
- instruction: `internal/baseline/plan_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/secondbrain.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/retirement_rule_clauses_test.go`
- interface: `internal/baseline/repository_rule_placement_test.go`
- creates: `internal/baseline/documentation_rules_clauses_test.go`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/secondbrain.md`
- interface: `docs/agents/setup-context.json`
- interface: `CONTEXT.md`

## Consumers

- The modules `context-workflow`, `core` and `secondbrain`, their new rules and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's three guides, its Setup Manifest and the glossary — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheBaselineModulesCarryTheDocumentationRules|TestTheGuidesRenderTheDocumentationRulesOnce|TestBaselineClauseForceIsCharacterized|TestAnUnlistedBaselineClauseIsReported|TestEveryBaselineModuleVersionIsRecorded|TestTheCaptureTriggerAsksForAConditionalLesson|TestTheBaselineSendsRepositoryRulesToTheSpecificRepositoryGuide|TestTheDocsLayoutRegionStaysLean|TestADRLifecycleContract)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheBaselineModulesCarryTheDocumentationRules/the_lesson_rule TestTheBaselineModulesCarryTheDocumentationRules/the_export_rule TestTheBaselineModulesCarryTheDocumentationRules/required_by_no_profile TestTheBaselineModulesCarryTheDocumentationRules/names_no_record TestTheBaselineModulesCarryTheDocumentationRules/adr/amended_stays_accepted TestTheBaselineModulesCarryTheDocumentationRules/adr/the_later_adr_writes_it TestTheBaselineModulesCarryTheDocumentationRules/adr/dated_notice TestTheBaselineModulesCarryTheDocumentationRules/adr/the_update_fills TestTheBaselineModulesCarryTheDocumentationRules/legacy/sentence_kept TestTheBaselineModulesCarryTheDocumentationRules/legacy/the_exception TestTheBaselineModulesCarryTheDocumentationRules/legacy/stays_legacy TestTheBaselineModulesCarryTheDocumentationRules/layout/design TestTheBaselineModulesCarryTheDocumentationRules/layout/lessons TestTheBaselineModulesCarryTheDocumentationRules/layout/root_reference TestTheBaselineModulesCarryTheDocumentationRules/layout/variation TestTheBaselineModulesCarryTheDocumentationRules/layout/doctor_reports TestTheBaselineModulesCarryTheDocumentationRules/layout/rest_kept TestTheBaselineModulesCarryTheDocumentationRules/retirement/lesson_deleted TestTheBaselineModulesCarryTheDocumentationRules/message/no_message TestTheBaselineModulesCarryTheDocumentationRules/lesson/record_path TestTheBaselineModulesCarryTheDocumentationRules/lesson/labels_as_written TestTheBaselineModulesCarryTheDocumentationRules/lesson/eight_lines TestTheBaselineModulesCarryTheDocumentationRules/lesson/the_test TestTheBaselineModulesCarryTheDocumentationRules/lesson/otherwise_nothing TestTheBaselineModulesCarryTheDocumentationRules/lesson/two_moments TestTheBaselineModulesCarryTheDocumentationRules/lesson/in_that_change TestTheBaselineModulesCarryTheDocumentationRules/export/two_parts TestTheBaselineModulesCarryTheDocumentationRules/export/markers TestTheBaselineModulesCarryTheDocumentationRules/export/standard_items TestTheBaselineModulesCarryTheDocumentationRules/export/baseline_owns TestTheBaselineModulesCarryTheDocumentationRules/export/never_process TestTheBaselineModulesCarryTheDocumentationRules/variation/own_lines TestTheBaselineModulesCarryTheDocumentationRules/variation/exclusion TestTheBaselineModulesCarryTheDocumentationRules/variation/reason TestTheBaselineModulesCarryTheDocumentationRules/variation/doctor TestTheBaselineModulesCarryTheDocumentationRules/mark/add_it TestTheBaselineModulesCarryTheDocumentationRules/mark/the_test TestTheBaselineModulesCarryTheDocumentationRules/mark/unmarked_unpublished TestTheBaselineModulesCarryTheDocumentationRules/mark/existing_once TestTheBaselineModulesCarryTheDocumentationRules/mark/one_reviewed_commit TestTheGuidesRenderTheDocumentationRulesOnce/plan_postimage TestTheGuidesRenderTheDocumentationRulesOnce/formatter_golden TestTheGuidesRenderTheDocumentationRulesOnce/repository_guide TestBaselineClauseForceIsCharacterized TestAnUnlistedBaselineClauseIsReported TestEveryBaselineModuleVersionIsRecorded TestTheCaptureTriggerAsksForAConditionalLesson TestTheBaselineSendsRepositoryRulesToTheSpecificRepositoryGuide TestTheDocsLayoutRegionStaysLean TestADRLifecycleContract; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 40m -run 'Baseline|GuidanceComposition|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the two new tests do not exist, so their pass lines are missing; afterwards every guarantee and the rendering pass, the closed list holds the five forces, the module versions are recorded, the three updated tests pass, and the Baseline suite and the CLI's Baseline tests pass.
- `for pair in 'docs/agents/docs-layout.md|and a Lesson Learned once an accepted ADR, a guide or a skill states what it teaches are deleted' 'docs/agents/docs-layout.md|Write a Lesson Learned only when every part of this test holds' 'docs/agents/agent-instructions.md|send no message about what it taught' 'docs/agents/secondbrain.md|An ADR without the mark is not published.'; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; bin_dir="$(mktemp -d)" || exit 1; home_dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$bin_dir/roundfix" ./cmd/roundfix || { rm -rf "$bin_dir" "$home_dir"; exit 1; }; out="$(HOME="$home_dir" "$bin_dir/roundfix" baseline update --repo . --no-skills --no-history --format json </dev/null 2>&1)"; status=$?; rm -rf "$bin_dir" "$home_dir"; test "$status" -eq 0 || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"state":"current"' || { printf '%s\n' "$out"; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task this repository's guides hold none of the new text; afterwards they hold it, a Managed Refresh planned under a fresh home reports the repository `current`, and the documentation contracts, the Setup Manifest and the glossary included, pass.
- `grep -F -- '**Standard Export List**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Standard Export List\*\*(:| —( |$))' && grep -F -- '**Export Variation**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Export Variation\*\*(:| —( |$))' && grep -F -- '**Knowledge Mark**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Knowledge Mark\*\*(:| —( |$))' && grep -F -- '**Amendment List**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Amendment List\*\*(:| —( |$))' && grep -F -- '**Lesson Learned**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Lesson Learned\*\*(:| —( |$))' && grep -F -- '**Managed Refresh**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Managed Refresh\*\*(:| —( |$))' && { tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'It also keeps the Standard Export List and fills each Amendment List, in a documentation part of its plan that the Plan Digest covers.' || { printf 'missing phrase in CONTEXT.md: the Managed Refresh sentence\n' >&2; exit 1; }; }` — expected: exit 0. Before this Task none of the five terms is defined and Managed Refresh lacks the sentence; afterwards each term has its Glossary Definition line and Managed Refresh carries the new sentence.

## References

- `_prd.md` → Goals; User Stories 3, 5, 7, 8, 9 and 10; Core Features 2, 3, 6, 9, 10, 11, 12, 13, 14, 15 and 20; Success Metric 1; Success Metric 8; Success Metric 9
- `_techspec.md` → Baseline text; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; API Contract 21; Testing Approach; Build Order 1
- [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), sections 2, 3 and 4
- ADR-0438, ADR-0439, ADR-0440, ADR-0186, ADR-0222, ADR-0257
