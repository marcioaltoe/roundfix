---
task: task_02
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
status: pending
type: backend
complexity: medium
---

# Task 02: External triage states its rules as clauses with force

## Overview

The removed `triage` skill carried the functions external triage still needs, and the `external-triage` module states them in one unlabelled rule. This Task gives the module eight clauses with force, adds their Source Baseline rows, and declares that the first clause replaces the old rule's Source Baseline entry. A Standard TypeScript Monorepo adopter with external triage enabled is refused on refresh today ("1 unaccounted clause(s): rule.external-triage"); after this Task the same refresh reports the entry `replaced`.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST give `rule.external-triage` the eight clauses of the TechSpec's Fixed texts (task_02), with their exact identifiers, order, force and guidance, and MUST move the module to `setup-context-driven/module-v3` with `"repositoryExtensions": []`. The first clause MUST declare `"replaces": ["rule.external-triage"]`. MUST raise the module, guide and rule versions by one from the Task's base.
2. MUST add the eight Source Baseline rows exactly as the TechSpec's Source Baseline table and layout state, and MUST keep the `rule.external-triage` entry. The retention disposition of `rule.external-triage` MUST be `replaced` by `clause.external-triage.classify-before-labelling`; no other entry's disposition changes.
3. MUST add the eight clauses to the force record of `internal/baseline/clause_characterization_test.go`, and raise `maintainedSourceBaselineEntries` in `internal/baseline/preservation_test.go` by eight from its value on the Task's base. No other line of either file changes.
4. MUST create `internal/baseline/external_triage_clauses_test.go` with the four tests of Testing Approach 2, built on `externalTriageClauseFindings` and `newExternalTriageAdopter` as the TechSpec's Interfaces describe.
5. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`. MUST NOT hand-edit an offset, a digest, a golden or a pin.
6. MUST NOT cite a Spec or ADR number in any clause text, and MUST NOT change `decisions.json`: `triage.external` stays a boolean.

## Subtasks

- [ ] Write the eight clauses and the declared replacement.
- [ ] Add the Source Baseline rows.
- [ ] Update the force record and the entry count.
- [ ] Write the four tests.
- [ ] Regenerate, then refresh this repository twice.

## Acceptance Criteria

- [ ] The Standard TypeScript Monorepo `external-triage.md` golden renders each of the eight clauses as `- **<force>**: <guidance>` and no unlabelled bullet.
- [ ] A missing clause, a changed force and an unlabelled bullet are each reported by the helper.
- [ ] A Source Baseline adopter with external triage enabled refreshes with `rule.external-triage` reported `replaced` and the successor named; without the declaration the refresh is refused naming the rule.
- [ ] The force record, the duplicate-text check and the Source Baseline checks pass.

## Context

- instruction: `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- instruction: `docs/adr/0058-baseline-upgrades-fail-closed-on-unaccounted-rule-removal.md`
- instruction: `docs/adr/0060-source-baselines-are-exhaustive-and-project-agnostic.md`
- instruction: `docs/adr/0202-a-stack-rule-carries-force-and-a-preference-names-what-licenses-the-exception.md`
- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- instruction: `internal/baseline/plan.go`
- instruction: `internal/baseline/clause_replacement_test.go`
- instruction: `internal/baseline/assets/modules/bun.json`
- interface: `internal/baseline/assets/modules/external-triage.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/external-triage.md`
- interface: `internal/baseline/assets/source-baselines/index.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/external-triage.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/preservation_test.go`
- creates: `internal/baseline/external_triage_clauses_test.go`
- interface: `docs/agents/setup-context.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheExternalTriageGuideStatesItsClausesWithForce|TestAMissingOrUnlabelledExternalTriageClauseIsReported|TestTheExternalTriageRuleIsReplacedForSourceBaselineAdopters|TestAnUndeclaredExternalTriageReplacementIsRefused|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestReadoptionCompatibilityMaintainedFixture|TestSourceBaselineGuidanceComposition|TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestCatalogDiagnosticCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheExternalTriageGuideStatesItsClausesWithForce TestAMissingOrUnlabelledExternalTriageClauseIsReported TestTheExternalTriageRuleIsReplacedForSourceBaselineAdopters TestAnUndeclaredExternalTriageReplacementIsRefused TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestReadoptionCompatibilityMaintainedFixture TestSourceBaselineGuidanceComposition TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus TestCatalogDiagnosticCharacterization; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && grep -q 'clause.external-triage.classify-before-labelling' internal/baseline/assets/source-baselines/index.json && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the four new tests do not exist and no Source Baseline row names the new clauses, so the command fails.

## References

- `_prd.md` → Goal 3; User Story 3; Core Feature 3; Success Metric 3; Acceptance evidence
- `_techspec.md` → Measured facts; Data Models; Fixed texts (task_02); Existing tests that change; API Contract 3; Testing Approach 2 and 5; Build Order 2
- ADR-0206, ADR-0058, ADR-0060, ADR-0099, ADR-0202, ADR-0186
