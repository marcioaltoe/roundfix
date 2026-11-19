---
task: task_01
spec: 0409-the-idea-and-prd-question-rules-that-block-the-goal
status: completed
type: docs
complexity: medium
---

# Task 01: The Spec workflow puts a Rule Challenge to the user

## Overview

The Baseline's Spec workflow gains the rule `rule.spec.rule-challenge`, with two Normative Clauses: the first tells an author of brainstorming, an Idea or a PRD to put a **Rule Challenge** to the user instead of dropping a solution or contradicting a rule in silence, and the second says how the question is recorded, what follows each answer and what happens when nobody can answer. A contract test pins every guarantee of both clauses by its own phrase and their rendering in the Spec routing guide. Every adopter's Managed Refresh writes both clauses into its Spec routing guide, and this repository's own guide and Setup Manifest are refreshed the same way.

## Requirements

1. MUST add to `internal/baseline/assets/modules/spec-workflow.json` the rule object of `_techspec.md` → Interfaces, immediately before the rule `rule.spec.artifacts`: id `rule.spec.rule-challenge`, `version` 1, `coverage` `["coverage.spec-workflow"]`, and two clauses in this order, `clause.spec.rule-challenge-01-ask` with `enforcement` `stop-and-ask` and the `guidance` of API Contract 1, and `clause.spec.rule-challenge-02-answer` with `enforcement` `mandatory` and the `guidance` of API Contract 2. Each `guidance` is one JSON string with no line break, word for word, and no clause declares `replaces`. The file keeps its two-space JSON formatting.
2. MUST add `rule.spec.rule-challenge` to the `rules` list of the guide `guide.spec-routing`, between `rule.spec.routing` and `rule.spec.project-constraints`, and raise that guide's `version` by one from the value the Task finds. MUST NOT change any other rule, clause, guide or field of the module, and MUST NOT list the new rule under `requiredRules` in any profile (`_techspec.md` → Invariants 1, 2 and 5).
3. MUST add the two clause ids to the closed list `characterizedBaselineForce` of `internal/baseline/clause_characterization_test.go`, `clause.spec.rule-challenge-01-ask` as `stop-and-ask` and `clause.spec.rule-challenge-02-answer` as `mandatory`, in the list's order, and change no other line of that file.
4. MUST create `internal/baseline/rule_challenge_clauses_test.go` with two tests, each calling `t.Parallel()` first. Every phrase below is compared on whitespace-normalized text and is required exactly once in its clause's `guidance`.
   - `TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses` reads the module file, as the tests of `internal/baseline/authoring_evidence_clauses_test.go` do, and runs these subtests:
     - `the_rule`: the rule exists once with version 1, the module's coverage id and exactly the two clauses of requirement 1, with those ids, forces and order and no `replaces`; the guide `guide.spec-routing` lists it right after `rule.spec.routing` and right before `rule.spec.project-constraints`; and no profile under `internal/baseline/assets/profiles/` lists it under `requiredRules`.
     - `names_no_record`: neither `guidance` holds `ADR-` followed by a digit or `Spec ` followed by four digits (`_techspec.md` → Invariant 4).
     - `ask/both_silences`: "never drop or narrow a solution in silence because a rule forbids it, and never contradict a rule in silence: put a Rule Challenge to the user."
     - `ask/what_a_rule_is`: "A rule is a decision the repository recorded: an accepted ADR, a Normative Clause of a guide under"
     - `ask/the_trigger`: "A rule blocks the goal when the solution you would recommend, or a part the goal needs, cannot exist without contradicting it, or when it stops you from stating the problem"
     - `ask/not_a_block`: "extra work, taste and friction are not a block."
     - `ask/the_question`: "name the rule and where it lives, what it imposes, up to three things it prevents for this goal, and your reading of why it exists and whether it looks out of date or badly written"
     - `ask/three_answers`: "then offer keep, change and retire, each with its consequence and the recommended one first."
     - `ask/binds_until_the_answer`: "The rule binds until the user answers: asking changes nothing"
     - `ask/never_an_answer`: "a recommendation, silence, a dismissal or a timeout is never an answer."
     - `ask/a_decided_change`: "Ask no Rule Challenge about a change the user already decided, because the request asks for it or it follows from an answer of the user that the artifact records"
     - `ask/a_kept_line_needs_a_new_fact`: "Ask none about a rule that carries a Kept Line unless the question names a fact the earlier one did not have."
     - `ask/at_most_three`: "Ask at most three in one Authoring Phase or brainstorming session"
     - `ask/the_others_keep_binding`: "list the others, which keep binding, in the artifact's open questions."
     - `ask/the_techspec_asks`: "The TechSpec author asks the same question before a design takes what an accepted ADR rejected, or otherwise contradicts its decision, unless the user already decided it."
     - `ask/implementation_guard_rail`: "In a Run, in a Task and in the QA gate a rule is a guard-rail: flag a conflict there and ask no Rule Challenge."
     - `answer/record_fields`: "one fielded record per question:" and "plus `new fact` when the rule already carried a Kept Line."
     - `answer/none_when_no_rule_blocked`: "An Idea and a feature PRD always carry the section, which reads `None.` when no rule blocked the goal"
     - `answer/the_prd_carries_the_records`: "the PRD carries the Idea's records, and a question the TechSpec author asks is recorded in the PRD."
     - `answer/change_or_retire`: "On change or retire, the PRD states the change of the rule as a requirement, and a Task of the same Spec changes the rule together with the work that depended on it"
     - `answer/adr_lifecycle`: "an ADR through the lifecycle of `docs/agents/docs-layout.md`, any other rule in its own file."
     - `answer/kept_line_form`: "On keep, a Task of the Spec adds a Kept Line after the rule's text, `Kept (YYYY-MM-DD): <what was asked>; <why it stays>.`, which names no Spec"
     - `answer/newest_last`: "a rule keeps every Kept Line, the newest last"
     - `answer/no_spec_follows`: "when no Spec follows, the session that asked adds the line."
     - `answer/a_baseline_rule`: "A rule inside `setup-context-driven` markers changes only in the Baseline, and its Kept Line goes in `docs/agents/specific-repository.md`, under a heading of its own, naming the rule"
     - `answer/a_proposal_to_the_owner`: "a repository that does not own the Baseline sends the user's decision as one message to the Baseline's owner, obeys the rule until the Baseline changes, and records the message and the scope it left out."
     - `answer/nobody_present`: "When no person is present to answer, write no artifact for that phase and never answer for the user"
     - `answer/only_that_spec_waits`: "only that Spec waits, as a Backlog Entry that holds the whole question, names the Spec folder when one exists and carries `waiting_on` and `since` in its frontmatter, while other work goes on."
     - `answer/the_author_subagent`: "An Author Subagent returns the question to the session that started it."
   - `TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce` asserts `_techspec.md` → Invariant 3 in three subtests: each clause's whole `guidance`, whitespace-normalized, occurs exactly once, after the `guidance` of the last clause of `rule.spec.routing` and before the `guidance` of the first clause of `rule.spec.project-constraints`, all three read from the module file:
     - `plan_postimage`: in the postimage of `docs/agents/spec-routing.md` of a plan built with `newProjectDecisionPlanRepository`, `buildProjectDecisionPlan`, `standardTypeScriptDecisions` and `planPostimage`;
     - `formatter_golden`: in `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`;
     - `repository_guide`: in this repository's `docs/agents/spec-routing.md`.
5. MUST regenerate, in this order and never by hand (`_techspec.md` → Invariant 10): the module version with `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`, which writes the next version of `spec-workflow` (this Task names no number); the derived Baseline files with `make baseline-digests`; and, as the Task's last step, this repository's Spec routing guide and Setup Manifest with `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, run again until a run prints `File changes: 0`. The files the prototype of `references/2026-11-19-prototype-and-governed-paths.md`, section 3, saw those commands rewrite are declared below; any other file they rewrite on the tree the Task runs on is recorded by the Daemon.
6. MUST keep, byte for byte, every clause that exists at the Delivery Base, with its id, force and text, and every test of `./internal/baseline` that exists there, `internal/baseline/plan_test.go` included, other than the two lines requirement 3 adds (`_techspec.md` → Invariant 5).
7. MUST NOT change any skill, `CONTEXT.md`, any file under `docs/user-guide/`, any Source Baseline or Setup Snapshot under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/`, any other module, template or profile source, any guide under `docs/agents/` other than `spec-routing.md` and `setup-context.json`, `AGENTS.md`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration. Of the Governed Paths it changes only `internal/baseline/assets/modules/spec-workflow.json` by hand, and `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `docs/agents/spec-routing.md` and `docs/agents/setup-context.json` only through the commands of requirement 5, all five listed in `_authorization.md`.

## Subtasks

- [ ] Add the rule with its two clauses and list it in the Spec routing guide.
- [ ] Add the two ids to the closed list of clause forces.
- [ ] Add the clause and rendering tests.
- [ ] Record the module version, regenerate the derived files and refresh this repository's guide and Setup Manifest until a refresh changes nothing.

## Acceptance Criteria

- [ ] A plan for a profile that includes the Spec workflow renders both clauses exactly once in its Spec routing guide, after the route clauses and before the Project Constraints clauses; the Standard TypeScript golden and this repository's guide hold them the same way.
- [ ] Each guarantee of API Contracts 1 and 2 is pinned by its own named subtest, and removing its phrase from the module fails that subtest.
- [ ] Every clause of the Delivery Base keeps its id, force and text; the derived Baseline files and this repository's guide and Setup Manifest are current, and a further Managed Refresh reports the repository current.

## Context

- instruction: `docs/adr/0436-an-authoring-phase-questions-a-rule-that-blocks-the-goal-and-the-rule-binds-until-the-user-answers.md`
- instruction: `docs/adr/0437-a-rule-the-user-keeps-carries-a-kept-line-and-is-asked-about-again-only-with-a-new-fact.md`
- instruction: `internal/baseline/authoring_evidence_clauses_test.go`
- instruction: `internal/baseline/plan_test.go`
- instruction: `internal/baseline/profile_alignment_test.go`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/clause_characterization_test.go`
- creates: `internal/baseline/rule_challenge_clauses_test.go`
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

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses|TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce|TestBaselineClauseForceIsCharacterized|TestAnUnlistedBaselineClauseIsReported|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/the_rule TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/names_no_record TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/both_silences TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/what_a_rule_is TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/the_trigger TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/not_a_block TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/the_question TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/three_answers TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/binds_until_the_answer TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/never_an_answer TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/a_decided_change TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/a_kept_line_needs_a_new_fact TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/at_most_three TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/the_others_keep_binding TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/the_techspec_asks TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/ask/implementation_guard_rail TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/record_fields TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/none_when_no_rule_blocked TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/the_prd_carries_the_records TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/change_or_retire TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/adr_lifecycle TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/kept_line_form TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/newest_last TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/no_spec_follows TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/a_baseline_rule TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/a_proposal_to_the_owner TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/nobody_present TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/only_that_spec_waits TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses/answer/the_author_subagent TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce/plan_postimage TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce/formatter_golden TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce/repository_guide TestBaselineClauseForceIsCharacterized TestAnUnlistedBaselineClauseIsReported TestEveryBaselineModuleVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 40m -run 'Baseline|GuidanceComposition|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the two tests do not exist, so their pass lines are missing; afterwards every guarantee and the rendering pass, the closed list holds the two forces, the module version is recorded, and the Baseline suite and the CLI's Baseline tests pass.
- `for phrase in 'In a Run, in a Task and in the QA gate a rule is a guard-rail: flag a conflict there and ask no Rule Challenge.' 'An Author Subagent returns the question to the session that started it.'; do tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/spec-routing.md: %s\n' "$phrase" >&2; exit 1; }; done; bin_dir="$(mktemp -d)" || exit 1; home_dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$bin_dir/roundfix" ./cmd/roundfix || { rm -rf "$bin_dir" "$home_dir"; exit 1; }; out="$(HOME="$home_dir" "$bin_dir/roundfix" baseline update --repo . --no-skills --no-history --format json </dev/null 2>&1)"; status=$?; rm -rf "$bin_dir" "$home_dir"; test "$status" -eq 0 || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"state":"current"' || { printf '%s\n' "$out"; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task this repository's Spec routing guide holds neither clause; afterwards it holds both, a Managed Refresh planned under a fresh home reports the repository `current`, and the documentation contracts, the Setup Manifest included, pass.

## References

- `_prd.md` → Goals; User Stories 1, 2, 3, 8, 9 and 10; Core Features 1, 2, 3, 4, 5, 6, 8, 9, 10, 11, 12 and 13; Success Metric 1; Success Metric 4
- `_techspec.md` → Rule Challenge rule; Derived Baseline files and Managed Refresh; Interfaces; Invariants 1-5, 9 and 10; API Contract 1; API Contract 2; Testing Approach; Build Order 1
- [`references/2026-11-19-prototype-and-governed-paths.md`](references/2026-11-19-prototype-and-governed-paths.md), sections 2 and 3
- ADR-0436, ADR-0437, ADR-0257, ADR-0186, ADR-0222

## Result

Skills loaded: implement-task, golang-testing, testing-boss, tech-writer, setup-context-driven, roundfix, evidence-gate

- Added the Rule Challenge rule, its two clauses and the Spec routing guide ordering; the closed force list includes both clause IDs. Existing clause text and forces are unchanged in the module diff.
- Rendering criterion: `go test ./internal/baseline -run '^TestTheSpecRoutingGuideRendersTheRuleChallengeClausesOnce$' -count=1` passed, covering the plan postimage, formatter golden and repository guide.
- Phrase criterion: `go test ./internal/baseline -run '^TestTheSpecWorkflowModuleCarriesTheRuleChallengeClauses$' -count=1` passed, including the named phrase checks, rule structure, profile exclusions and record-name check.
- Refresh criterion: the sanctioned module version recorder and `make baseline-digests` exited successfully; the Managed Refresh updated two files and its second run reported `File changes: 0`. The module and characterization diffs preserve the existing clauses and add only the two force entries.
- The Daemon owns the declared Verification and Task status; neither was changed here.
