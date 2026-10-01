---
task: task_01
spec: 0209-sources-that-share-a-context-share-a-spec
status: pending
type: backend
complexity: medium
---

# Task 01: Sources that share a context share a Spec, as Baseline clauses with force

## Overview

The Baseline does not tell an author to group sources that share a context,
to stop grouping at a size, or to extend an open record instead of minting a
new one. This Task adds the three mandatory clauses of `_techspec.md` → Exact
clause texts to the Spec workflow and CONTEXT workflow modules, gives each its
Source Baseline row, regenerates the derived files and refreshes this
repository's guides (ADR-0208). It answers the maintainer's requests of
2026-10-01 quoted in `_prd.md`. It is verifiable on its own: the rendered
guides state the clauses with force, and a Source Baseline adopter's plan
records each one `retained`.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST append to `rule.spec.routing` in `internal/baseline/assets/modules/spec-workflow.json`, after `clause.spec.verification-two-tiers`, the clauses `clause.spec.sources-01-group-by-shared-context` and `clause.spec.sources-02-bound-the-group`, and to `rule.context.docs-layout` in `internal/baseline/assets/modules/context-workflow.json`, after `clause.context.inbox-02-fleet-flow`, the clause `clause.context.inbox-03-extend-before-minting`, each `mandatory` with the exact guidance of `_techspec.md` → Exact clause texts. Every existing clause MUST stay byte-identical, and no clause gets a `replaces` list.
2. MUST add the three Source Baseline rows of `_techspec.md` → Source Baseline rows: corpus entry, `manifest.json` row and `index.json` `entryIds` identifier, each after its named anchor. No existing row changes.
3. MUST raise by one, from the value on this Task's starting commit, the six versions `_techspec.md` → Version changes lists, keeping each module's formatting.
4. MUST update the existing tests `_techspec.md` → Existing tests that change names for this Task: the force record gains the three clauses as `mandatory`, and the maintained Source Baseline entry count rises by three. No other existing test line changes.
5. MUST create `internal/baseline/grouped_sources_clauses_test.go` with `TestTheGroupingClausesCarryTheirForceAndText`, `TestTheGroupingClausesRenderInTheGuides`, `TestTheGroupingClausesHaveSourceBaselineRows` and `TestAnAdopterRetainsTheGroupingClauses`, following `_techspec.md` → Testing Approach 1. The texts are literals in the test, so a reworded clause fails it; the retention test MUST assert a ready plan and the `retained` disposition for each clause, and no `unaccounted` clause.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`; a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a snapshot or a generated guide.
7. MUST prove each new gate can fail. The Result MUST record one sabotage of a clause text (for example dropping its last sentence) and one of a Source Baseline row (for example removing its manifest row), each with the test or check that failed, and that the source was restored.

## Subtasks

- [ ] Add the three clauses and raise the versions.
- [ ] Add the three Source Baseline rows.
- [ ] Update the two existing tests and create the new test file.
- [ ] Regenerate, refresh twice, and record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/spec-routing.md` states both Spec-workflow clauses and `docs/agents/docs-layout.md` the extend clause, each labelled `mandatory`.
- [ ] The Standard TypeScript Monorepo goldens state the same clauses.
- [ ] A Source Baseline adopter's Managed Refresh plan is ready and records the three clauses `retained`.
- [ ] The force record, the duplicate-text check, the record-citation check and the catalog validation pass, and a second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`
- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/source-baselines/index.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/preservation_test.go`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/grouped_sources_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheGroupingClausesCarryTheirForceAndText|TestTheGroupingClausesRenderInTheGuides|TestTheGroupingClausesHaveSourceBaselineRows|TestAnAdopterRetainsTheGroupingClauses|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestReadoptionCompatibilityMaintainedFixture|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheGroupingClausesCarryTheirForceAndText TestTheGroupingClausesRenderInTheGuides TestTheGroupingClausesHaveSourceBaselineRows TestAnAdopterRetainsTheGroupingClauses TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestShippedGuidanceCitesNoRepositoryRecord TestReadoptionCompatibilityMaintainedFixture TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for pair in "docs/agents/spec-routing.md|- **mandatory**: One Spec may adopt several Inbox Entries, Backlog Entries, and Findings whose context is similar or complementary" "docs/agents/spec-routing.md|is advisory: adopt the suggested source or state why it stays apart." "docs/agents/spec-routing.md|- **mandatory**: Group sources into one Spec only while the Spec fits four implementation Tasks plus its QA gate." "docs/agents/docs-layout.md|- **mandatory**: Before minting a Finding or Backlog Entry, including in Triage" "docs/agents/docs-layout.md|extend a fitting Finding with a dated addendum, because a Finding is immutable history."; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || exit 1; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the four new tests do not exist and the guides lack the clauses, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1, 2, 3 and 4; User Stories 1, 2, 3 and 6; Core Features 1, 2, 3 and 4; Success Metrics 1 and 2; Declared breaks
- [_techspec.md](_techspec.md) — Exact clause texts; Source Baseline rows; Version changes; Existing tests that change; API Contract 6; API Contract 7; Testing Approach 1; Build Order 1
- ADR-0208; ADR-0186; ADR-0058; ADR-0060; ADR-0081; ADR-0149
