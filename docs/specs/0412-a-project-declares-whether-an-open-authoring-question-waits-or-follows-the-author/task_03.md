---
task: task_03
spec: 0412-a-project-declares-whether-an-open-authoring-question-waits-or-follows-the-author
status: pending
type: docs
complexity: medium
---

# Task 03: The Spec workflow says what a session does with an Open Authoring Question

## Overview

The Baseline's Spec workflow gains the rule `rule.spec.open-authoring-question`, with two Normative Clauses: the first says what an **Open Authoring Question** is, that a present user is asked, that with nobody present the Spec waits unless following is in force, and which questions are a **Reserved Question** and always wait; the second says how a **Delegated Decision** is recorded and what follows the user's later answer or a wait. A contract test pins every guarantee of both clauses by its own phrase and their rendering in the Spec routing guide. Every adopter's Managed Refresh writes both clauses into its Spec routing guide, and this repository's own guide and Setup Manifest are refreshed the same way.

## Requirements

1. MUST add to `internal/baseline/assets/modules/spec-workflow.json` the rule object of `_techspec.md` → Interfaces, immediately before the rule `rule.spec.artifacts`: id `rule.spec.open-authoring-question`, `version` 1, `coverage` `["coverage.spec-workflow"]`, and two clauses in this order, `clause.spec.open-authoring-question-01-wait-or-follow` with `enforcement` `stop-and-ask` and the `guidance` of API Contract 4, and `clause.spec.open-authoring-question-02-record` with `enforcement` `mandatory` and the `guidance` of API Contract 5. Each `guidance` is one JSON string with no line break, word for word, written in raw UTF-8 as the rest of the file is, and no clause declares `replaces`. The file keeps its two-space JSON formatting.
2. MUST add `rule.spec.open-authoring-question` to the `rules` list of the guide `guide.spec-routing` immediately before `rule.spec.project-constraints`, which on a tree that holds Spec 0409 puts it right after `rule.spec.rule-challenge`, and raise that guide's `version` by one from the value the Task finds. MUST NOT change any other rule, clause, guide or field of the module, and MUST NOT list the new rule under `requiredRules` in any profile (`_techspec.md` → Invariants 7, 9 and 10).
3. MUST add the two clause ids to the closed list `characterizedBaselineForce` of `internal/baseline/clause_characterization_test.go`, `clause.spec.open-authoring-question-01-wait-or-follow` as `stop-and-ask` and `clause.spec.open-authoring-question-02-record` as `mandatory`, and change no other entry of that list.
4. MUST create `internal/baseline/open_authoring_question_clauses_test.go` with two tests, each calling `t.Parallel()` first. Every phrase below is compared on whitespace-normalized text and is required exactly once in its clause's `guidance`.
   - `TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses` reads the module file, as the clause tests of `internal/baseline/authoring_evidence_clauses_test.go` do, and runs these subtests:
     - `the_rule`: the rule exists once with version 1, the module's coverage id and exactly the two clauses of requirement 1, with those ids, forces and order and no `replaces`; it is the rule right before `rule.spec.artifacts` in the module's `rules`; the guide `guide.spec-routing` lists it right before `rule.spec.project-constraints`; and no profile under `internal/baseline/assets/profiles/` lists it under `requiredRules`.
     - `names_no_record`: neither `guidance` holds `ADR-` followed by a digit or `Spec ` followed by four digits (`_techspec.md` → Invariant 9).
     - `ask/what_the_question_is`: "An Open Authoring Question is a question the author of a PRD, a TechSpec or a Task Graph needs answered to finish its phase and that no source settles"
     - `ask/the_question_contract`: "with two or three options, the recommended one first and its reason"
     - `ask/the_reserved_statement`: "and say whether it is a Reserved Question and why."
     - `ask/the_present_user`: "When the user is present, ask: a question the user can see waits for the user"
     - `ask/adopts_nothing`: "a recommendation, silence, a dismissal or a timeout adopts nothing."
     - `ask/three_ways_told`: "Treat the user as absent only when the brief that started the session says so, the user said so, or the session runs as a continuation that an unattended flow armed."
     - `ask/the_wait`: "With the user absent the Spec waits, unless following is in force"
     - `ask/the_preference_grant`: "reports `follow` as the project's Authoring Question Preference"
     - `ask/the_answer_grant`: "quotes the user's answer, given before the flow, that the session follows the author."
     - `ask/never_writes_the_preference`: "Never write that preference into Project Config yourself."
     - `ask/reserved_idea_stage`: "the Idea stage and the choice an Idea-stage judgment leaves to the user;"
     - `ask/reserved_rule_challenge`: "a Rule Challenge;"
     - `ask/reserved_authorization`: "every authorization, which is a commit, a push, a Pull Request, a merge, a release, a Governed Path, a destructive Git operation, a database statement"
     - `ask/reserved_scope`: "removes one of their non-goals or contradicts a recorded answer of the user"
     - `ask/reserved_spend`: "that needs a paid service, a paid plan or a quota the repository does not already use"
     - `ask/reserved_data`: "that sends repository or client content to a service or a person the repository does not already send it to;"
     - `ask/reserved_shape`: "a question without two or three options and a recommended one;"
     - `ask/reserved_unclassified`: "and a question you cannot classify."
     - `ask/two_readers`: "An author that a session started returns the question and never adopts its own recommendation; the session checks the author's statement against this list, and when either reads the question as reserved or cannot tell, the Spec waits."
     - `record/adopts_the_recommended_option`: "When following is in force and an Open Authoring Question is not a Reserved Question, adopt the option the author recommended and go on with the phase."
     - `record/the_entry_head`: "the line `- delegated: YYYY-MM-DD, <phase>`, where the phase is `prd`, `techspec` or `tasks`"
     - `record/the_five_fields`: "followed by five lines indented two spaces, `- question:`, `- options:`, `- adopted:` with the option, an em dash and the author's reason"
     - `record/the_grant_and_answer_fields`: "`- under:` with `Project Config` or `the user's answer of YYYY-MM-DD`, and `- answer:` with `none yet`."
     - `record/never_a_dated_answer`: "Never write it as a dated answer of a person or under the autonomous-route line."
     - `record/later_phases_in_the_prd`: "A question of the TechSpec or the Tasks phase is recorded in the PRD, and that phase's artifact names the entry."
     - `record/report_and_status`: "List every Delegated Decision of the flow, and their count, in the session's report;"
     - `record/inside_the_authorization`: "The Spec goes on through authoring and delivery inside the authorization it already holds"
     - `record/the_later_answer`: "write it into the entry's `answer` line as the date, `kept` or `changed`, an em dash and the user's words."
     - `record/after_delivery_starts`: "leave the plan on the default branch as it is and make the change a correction of the Spec in delivery or a Backlog Entry"
     - `record/the_waiting_entry`: "only that Spec waits, as a Backlog Entry that names the Spec folder when one exists"
     - `record/the_resume`: "When the user answers, record the answer as a dated line of the PRD's `## Decisions`, let the Spec adopt the entry and resume the phase."
     - `record/the_other_waits`: "The waits of the Idea stage and of a Rule Challenge keep their own rules."
   - `TestTheSpecRoutingGuideRendersTheOpenAuthoringQuestionClausesOnce` asserts `_techspec.md` → Invariant 8 in three subtests: each clause's whole `guidance`, whitespace-normalized, occurs exactly once, after the `guidance` of the last clause of the rule the guide lists right before `rule.spec.open-authoring-question` and before the `guidance` of the first clause of `rule.spec.project-constraints`, all read from the module file:
     - `plan_postimage`: in the postimage of `docs/agents/spec-routing.md` of a plan built with `newProjectDecisionPlanRepository`, `buildProjectDecisionPlan`, `standardTypeScriptDecisions` and `planPostimage`;
     - `formatter_golden`: in `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`;
     - `repository_guide`: in this repository's `docs/agents/spec-routing.md`.
5. MUST regenerate, in this order and never by hand (`_techspec.md` → Invariants 13 and 14): the module version with `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`, which writes the next version of `spec-workflow` (this Task names no number); the derived Baseline files with `make baseline-digests`; and, as the Task's last step, this repository's Spec routing guide and Setup Manifest with `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, run again until a run prints `File changes: 0`. The files the prototype of `references/2026-11-19-prototype-and-governed-paths.md`, section 3, saw those commands rewrite are declared below; any other file they rewrite on the tree the Task runs on is recorded by the Daemon.
6. MUST keep, byte for byte, every clause that exists at the Delivery Base, with its id, force and text, the Rule Challenge clauses included, and every test of `./internal/baseline` that exists there, other than the two entries requirement 3 adds (`_techspec.md` → Invariant 10).
7. MUST NOT change any skill, `CONTEXT.md`, any ADR, any file under `docs/user-guide/`, any Source Baseline or Setup Snapshot under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/`, any other module, template or profile source, any guide under `docs/agents/` other than `spec-routing.md` and `setup-context.json`, `AGENTS.md`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration. Of the Governed Paths it changes only `internal/baseline/assets/modules/spec-workflow.json` by hand, and `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `docs/agents/spec-routing.md` and `docs/agents/setup-context.json` only through the commands of requirement 5, all five listed in `_authorization.md`.

## Subtasks

- [ ] Add the rule with its two clauses and list it in the Spec routing guide.
- [ ] Add the two ids to the closed list of clause forces.
- [ ] Add the clause and rendering tests.
- [ ] Record the module version, regenerate the derived files and refresh this repository's guide and Setup Manifest until a refresh changes nothing.

## Acceptance Criteria

- [ ] A plan for a profile that includes the Spec workflow renders both clauses exactly once in its Spec routing guide, before the Project Constraints clauses; the Standard TypeScript golden and this repository's guide hold them the same way.
- [ ] Each guarantee of API Contracts 4 and 5 is pinned by its own named subtest, and removing its phrase from the module fails that subtest.
- [ ] Every clause of the Delivery Base keeps its id, force and text; the derived Baseline files and this repository's guide and Setup Manifest are current, and a further Managed Refresh reports the repository current.

## Context

- instruction: `docs/adr/0442-an-open-authoring-question-waits-unless-the-project-declares-that-the-session-follows-the-author.md`
- instruction: `internal/baseline/authoring_evidence_clauses_test.go`
- instruction: `internal/baseline/plan_test.go`
- instruction: `internal/baseline/profile_alignment_test.go`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/clause_characterization_test.go`
- creates: `internal/baseline/open_authoring_question_clauses_test.go`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`

## Consumers

- The `spec-workflow` module, its two new clauses and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's Spec routing guide and Setup Manifest — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses|TestTheSpecRoutingGuideRendersTheOpenAuthoringQuestionClausesOnce|TestBaselineClauseForceIsCharacterized|TestAnUnlistedBaselineClauseIsReported|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/the_rule TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/names_no_record TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/what_the_question_is TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/the_question_contract TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/the_reserved_statement TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/the_present_user TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/adopts_nothing TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/three_ways_told TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/the_wait TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/the_preference_grant TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/the_answer_grant TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/never_writes_the_preference TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_idea_stage TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_rule_challenge TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_authorization TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_scope TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_spend TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_data TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_shape TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/reserved_unclassified TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/ask/two_readers TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/adopts_the_recommended_option TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/the_entry_head TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/the_five_fields TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/the_grant_and_answer_fields TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/never_a_dated_answer TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/later_phases_in_the_prd TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/report_and_status TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/inside_the_authorization TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/the_later_answer TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/after_delivery_starts TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/the_waiting_entry TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/the_resume TestTheSpecWorkflowModuleCarriesTheOpenAuthoringQuestionClauses/record/the_other_waits TestTheSpecRoutingGuideRendersTheOpenAuthoringQuestionClausesOnce/plan_postimage TestTheSpecRoutingGuideRendersTheOpenAuthoringQuestionClausesOnce/formatter_golden TestTheSpecRoutingGuideRendersTheOpenAuthoringQuestionClausesOnce/repository_guide TestBaselineClauseForceIsCharacterized TestAnUnlistedBaselineClauseIsReported TestEveryBaselineModuleVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 40m -run 'Baseline|GuidanceComposition|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the two tests do not exist, so their pass lines are missing; afterwards every guarantee and the rendering pass, the closed list holds the two forces, the module version is recorded, and the Baseline suite and the CLI's Baseline tests pass.
- `for phrase in 'Never write that preference into Project Config yourself.' 'The waits of the Idea stage and of a Rule Challenge keep their own rules.'; do tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/spec-routing.md: %s\n' "$phrase" >&2; exit 1; }; done; bin_dir="$(mktemp -d)" || exit 1; home_dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$bin_dir/roundfix" ./cmd/roundfix || { rm -rf "$bin_dir" "$home_dir"; exit 1; }; out="$(HOME="$home_dir" "$bin_dir/roundfix" baseline update --repo . --no-skills --no-history --format json </dev/null 2>&1)"; status=$?; rm -rf "$bin_dir" "$home_dir"; test "$status" -eq 0 || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"state":"current"' || { printf '%s\n' "$out"; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task this repository's Spec routing guide holds neither clause; afterwards it holds both, a Managed Refresh planned under a fresh home reports the repository `current`, and the documentation contracts, the Setup Manifest included, pass.

## References

- `_prd.md` → Goals; User Stories 3, 5, 8 and 9; Core Features 1, 2, 3, 7, 8, 9, 10, 14 and 15; Success Metric 5; Success Metric 8
- `_techspec.md` → Open Authoring Question rule; Interfaces; Invariants 7-10, 13 and 14; API Contract 4; API Contract 5; Testing Approach; Build Order 3
- `_prd.md` → Decisions, the Delegated Decision of 2026-11-19 of the `tasks` phase: how is step 3 of the Build Order cut into Tasks?
- [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), sections 2 and 3
- ADR-0442, ADR-0257, ADR-0186, ADR-0222, ADR-0436
