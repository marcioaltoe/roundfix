---
task: task_01
spec: 0391-the-grouping-bound-becomes-a-criterion
status: pending
type: docs
complexity: medium
---

# Task 01: The Baseline sizes a Spec by one context and by Tasks that earn their verification cycle

## Overview

The spec-workflow module's grouping clause keeps its id and force and
replaces the fixed bound of four implementation Tasks with two criteria, one
context and Tasks that each earn their verification cycle, states the two
reasons, keeps four Tasks as a reference, not a cap, and keeps every guarantee
the old text gave. A new test pins each guarantee and each criterion with its
own phrase in the rendered guide. Every adopter's managed refresh writes the
new words into its spec routing guide, and this repository's own guide and
setup manifest are refreshed the same way.

## Requirements

1. MUST change, in `internal/baseline/assets/modules/spec-workflow.json`, only
   the `guidance` of `clause.spec.sources-02-bound-the-group` to the text of
   `_techspec.md` → API Contract 1, word for word, and raise the `version` of
   `rule.spec.routing` by one from the value the Task finds (10 at
   `251ed81e`; Spec 0393 may raise it first). MUST keep the clause's id,
   `enforcement: mandatory` and absence of `replaces`, the file's two-space
   JSON formatting, and every other clause, rule and module unchanged; MUST NOT
   add or remove a clause.
2. MUST change the `clause.spec.sources-02-bound-the-group` literal in
   `groupingClauses` of `internal/baseline/grouped_sources_clauses_test.go` to
   the same text, and every other test that pins the old clause text on the
   default branch when the Task runs (at `251ed81e` that literal is the only
   one); the four tests that read the literal keep their names and assertions.
   MUST add to the same file the `groupingBoundPhrases` list and
   `TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee`, calling
   `t.Parallel()` first, with the eleven named phrases, the absent phrase and
   the subtest names `_techspec.md` → Testing Approach lists for
   `internal/baseline`.
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
   Baseline is immutable and retains the clause by id and enforcement), any
   other module, template or profile source by hand, or `.roundfixrc.yml`. Of
   the Governed Paths the Task may change only
   `internal/baseline/assets/modules/spec-workflow.json`,
   `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
   `docs/agents/spec-routing.md` and `docs/agents/setup-context.json`, the
   last four through the commands above.

## Subtasks

- [ ] Reword the clause and raise the rule version.
- [ ] Update the pinned literal and add the phrase test.
- [ ] Record the module version and regenerate the derived files.
- [ ] Refresh this repository's guide and setup manifest, and confirm a second refresh changes nothing.

## Acceptance Criteria

- [ ] A resolved profile's spec routing guide sizes a Spec by one context and by Tasks that earn their verification cycle, states both reasons, keeps four Tasks as a reference, and no longer bounds grouping by a count.
- [ ] Each kept guarantee and each new criterion is pinned by its own phrase, the clause keeps its id, force and Source Baseline row, and an adopter's refresh records it retained.
- [ ] The derived Baseline files and this repository's guide and setup manifest are current, and a second managed refresh reports no change.

## Context

- instruction: `docs/adr/0400-a-spec-is-sized-by-one-context-and-by-tasks-that-earn-their-verification-cycle.md`
- instruction: `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`
- instruction: `internal/baseline/plan.go`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/grouped_sources_clauses_test.go`
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

- The reworded clause and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`
- This repository's managed guide and setup manifest — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee|TestTheGroupingClausesCarryTheirForceAndText|TestTheGroupingClausesRenderInTheGuides|TestTheGroupingClausesHaveSourceBaselineRows|TestAnAdopterRetainsTheGroupingClauses|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee ' || { echo 'missing pass: TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 30m -run 'Baseline|Guidance|SetupManifest' ./internal/cli` — expected: exit 0. Before this Task the phrase test does not exist, so its pass line is missing; afterwards every phrase renders once, the clause keeps its force, row and retention, the module version is recorded, and the Baseline suite and the CLI's Baseline tests pass.

- `tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- 'Size a Spec by two criteria, never by a count of its Tasks.' || { echo 'spec routing guide lacks the criterion' >&2; exit 1; }; if tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- 'only while the Spec fits four implementation Tasks'; then echo 'spec routing guide still bounds grouping by a count' >&2; exit 1; fi; f="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json > "$f"; status=$?; grep -q '"state":"current"' "$f" || { cat "$f"; rm -f "$f"; exit 1; }; rm -f "$f"; test "$status" -eq 0 && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide holds the count and not the criterion; afterwards this repository's guide carries the criterion, a managed refresh reports `current`, and the documentation contracts, the setup manifest included, pass.

## References

- `_prd.md` → Goals 1 to 5; Core Feature 1; Success Metric 1; Success Metric 4
- `_techspec.md` → Current behavior 1 to 3; API Contract 1; Invariants 1 to 4; Declared breaks; Testing Approach; Build Order 1
- ADR-0400
