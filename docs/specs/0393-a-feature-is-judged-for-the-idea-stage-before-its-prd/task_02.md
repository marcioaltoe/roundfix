---
task: task_02
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
status: pending
type: docs
complexity: medium
---

# Task 02: The route clauses judge a feature for the Idea stage before its PRD

## Overview

The three route clauses of the Baseline's Spec routing rule keep their ids and
force and carry the rule: before the PRD of a feature, three criteria are
judged, by the command for a recorded item or a request of the session and by
the agent when the judge does not answer; an indicated stage starts and may be
dismissed, an undecided one is the user's choice, and nothing is refused. A
new test pins each new sentence and each guarantee the clauses already gave
with its own phrase in the rendered guide. Every adopter's managed refresh
writes the new words into its Spec routing guide, and this repository's own
guide and setup manifest are refreshed the same way.

## Requirements

1. MUST change, in `internal/baseline/assets/modules/spec-workflow.json`, only
   the `guidance` of `clause.spec.routing-01-large-initiative`,
   `clause.spec.routing-02-standard-feature` and
   `clause.spec.routing-05-task-graph` to the texts of `_techspec.md` → API
   Contract 5, word for word, and raise the `version` of `rule.spec.routing`
   by one from the value the Task finds. MUST keep each clause's id,
   `enforcement: mandatory` and absence of `replaces`, the Surface Transcript
   sentence of the third clause byte for byte, the file's two-space JSON
   formatting and every other clause, the grouping clauses of Spec 0391
   included; MUST NOT add or remove a clause (Invariant 6).
2. MUST add `internal/baseline/idea_stage_route_clauses_test.go` with
   `TestTheRouteClausesJudgeAFeatureForTheIdeaStage`, calling `t.Parallel()`
   first and built like `authoring_evidence_clauses_test.go`: one named
   subtest for every sentence of API Contract 5 and for every kept phrase of
   Invariant 6, each asserting that the phrase occurs once in its clause, that
   the clause is `mandatory` and carries no `replaces`, and that the phrase
   occurs once, whitespace-normalized, in the Standard TypeScript golden and
   in this repository's Spec routing guide. The sentences "The judgment
   informs and never refuses: no command fails for it, and a dismissal needs
   only its recorded reason." and "A fix or a refactor is never judged for the
   Idea stage." MUST each have their own subtest.
3. MUST regenerate, in this order and never by hand: the module record with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   (it writes the next version of `spec-workflow`; the Spec names no number);
   the derived Baseline files with `make baseline-digests`; and, as the Task's
   last step, this repository's managed guide and setup manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   whose second run reports no change. The derived files are measured on the
   tree the Task runs on: the files the prototype rewrote are declared below,
   and any other file those commands rewrite is recorded by the Daemon.
4. MUST NOT change any Source Baseline file (the corpus, manifest and index of
   `baseline.standard-typescript-monorepo-0.0.1` stay byte-identical: a Source
   Baseline is immutable and retains the clauses by id and enforcement), any
   other module, template or profile source by hand, any existing Baseline
   test, or `.roundfixrc.yml`. Of the Governed Paths the Task may change only
   `internal/baseline/assets/modules/spec-workflow.json`,
   `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
   `docs/agents/spec-routing.md` and `docs/agents/setup-context.json`, the
   last four through the commands above.

## Subtasks

- [ ] Reword the three clauses and raise the rule version.
- [ ] Add the phrase test.
- [ ] Record the module version and regenerate the derived files.
- [ ] Refresh this repository's guide and setup manifest, and confirm a second refresh changes nothing.

## Acceptance Criteria

- [ ] A resolved profile's Spec routing guide tells the agent to judge three criteria before a feature's PRD, to run the command for a recorded item or a request of the session, to answer the criteria itself when the judge does not, to start an indicated Idea stage, to ask the user for an undecided one, to record the judgment, and never to refuse.
- [ ] Each new sentence and each kept guarantee is pinned by its own phrase; the three clauses keep their ids, force and Source Baseline rows, and the grouping clauses are unchanged.
- [ ] The derived Baseline files and this repository's guide and setup manifest are current, and a second managed refresh reports no change.

## Context

- instruction: `docs/adr/0403-a-feature-is-judged-for-the-idea-stage-before-its-prd-and-the-judgment-never-refuses.md`
- instruction: `internal/baseline/authoring_evidence_clauses_test.go`
- instruction: `internal/baseline/plan.go`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- creates: `internal/baseline/idea_stage_route_clauses_test.go`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`

## Consumers

- The reworded clauses and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's managed guide and setup manifest — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheRouteClausesJudgeAFeatureForTheIdeaStage|TestTheAuthoringEvidenceClausesCarryTheirText|TestTheAuthoringEvidenceClausesRenderInTheGuide|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheRouteClausesJudgeAFeatureForTheIdeaStage ' || { echo 'missing pass: TestTheRouteClausesJudgeAFeatureForTheIdeaStage' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 30m -run 'Baseline|Guidance|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the phrase test does not exist, so its pass line is missing; afterwards every phrase renders once, the clauses keep their force and the Surface Transcript sentence, the module version is recorded, and the Baseline suite and the CLI's Baseline tests pass.

- `tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- 'The Idea stage is indicated when any criterion is yes, not indicated when every criterion is no, and undecided otherwise.' || { echo 'the spec routing guide lacks the Idea-stage rule' >&2; exit 1; }; out="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"state":"current"' || { printf '%s\n' "$out"; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide holds no Idea-stage rule; afterwards this repository's guide carries it, a managed refresh reports `current`, and the documentation contracts, the setup manifest included, pass.

## References

- `_prd.md` → Goals; User Story 2; User Story 3; User Story 6; User Story 7; Core Feature 5; Core Feature 6; Core Feature 9; Success Metric 5
- `_techspec.md` → Current behavior 5; API Contract 5; Invariant 6; Declared breaks; Testing Approach; Build Order 2
- ADR-0403
