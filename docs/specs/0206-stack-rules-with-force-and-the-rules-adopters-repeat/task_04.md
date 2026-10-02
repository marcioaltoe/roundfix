---
task: task_04
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
status: pending
type: backend
complexity: medium
---

# Task 04: Recurring Spec-workflow and TypeScript rules become clauses

## Overview

Four more recurring adopter rules become clauses: every Task Verification command fails before the change, a Task names the finding it answers, no test reads a Spec artifact, and a TypeScript row fixture is typed from the schema. Each gains its Source Baseline row. The slice is verified through the rendered Spec-routing, docs-layout and TypeScript guides and a check that every promoted clause of a Standard TypeScript Monorepo module has its Source Baseline row.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST add, in `internal/baseline/assets/modules/spec-workflow.json` and `typescript.json`, the four clauses the TechSpec's "Exact texts" gives for task_04, and MUST keep every existing clause byte-identical.
2. MUST add the four task_04 Source Baseline rows the TechSpec's table gives, each after its named anchor, and MUST NOT change or remove any existing row.
3. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_04.
4. MUST add the four clauses to the force record and raise the maintained Source Baseline entry count by four.
5. MUST create `internal/baseline/promoted_spec_and_typescript_clauses_test.go` with `sourceBaselineRowFindings` and the four tests the TechSpec's Testing Approach 4 names. The Source Baseline test MUST cover all twelve promoted clauses of task_01 and task_04 with their force.
6. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`;
   a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a snapshot or a generated guide.
7. MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Add the three Spec-workflow clauses and the TypeScript clause.
- [ ] Add the four Source Baseline rows.
- [ ] Raise versions, update the force record and the maintained count, create the new test file.
- [ ] Regenerate and refresh twice.

## Acceptance Criteria

- [ ] The rendered Spec-routing, docs-layout and TypeScript guides state the four clauses with their force in the Standard TypeScript Monorepo plan, and the Spec-routing and docs-layout clauses in the Go CLI/TUI plan. The composed profile Spec 0207 shipped selects the same modules; its repeated-clause and convergence tests pass in the repository Verification.
- [ ] Every promoted clause of a Standard TypeScript Monorepo module has a Source Baseline row with its force, and a baseline missing one is reported.
- [ ] The maintained Source Baseline count, the force record and the catalog validation pass, and a second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0203-a-repository-rule-becomes-a-baseline-clause-only-when-it-recurs.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/modules/typescript.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/source-baselines/index.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/typescript-bun.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md`
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
- creates: `internal/baseline/promoted_spec_and_typescript_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheSpecAndTypeScriptGuidesStateThePromotedRules|TestThePromotedSpecAndTypeScriptClausesCarryTheirForce|TestEveryPromotedClauseHasItsSourceBaselineRow|TestAMissingSourceBaselineRowIsReported|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestReadoptionCompatibilityMaintainedFixture|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestStandardTypeScriptStructuralClauseRetention)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheSpecAndTypeScriptGuidesStateThePromotedRules TestThePromotedSpecAndTypeScriptClausesCarryTheirForce TestEveryPromotedClauseHasItsSourceBaselineRow TestAMissingSourceBaselineRowIsReported TestBaselineClauseForceIsCharacterized TestReadoptionCompatibilityMaintainedFixture TestCatalogCompatibility; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/spec-routing.md|When a Task answers a recorded finding, name that finding and its date" "docs/agents/spec-routing.md|so every command fails on the tree before the Task's change" "docs/agents/docs-layout.md|Do not make a test, fixture, or build step read a file under the Spec root" "internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md|from the schema's inferred row type"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the four new tests do not exist and the guides lack the clauses, so the command fails.

## References

- `_prd.md` → Goal 4; Story 4; Core Feature 4; Success Metric 6; Success Metric 7
- `_techspec.md` → Candidate rules; Exact texts (task_04); Source Baseline rows; Version changes; Existing tests that change; API Contract 4; Testing Approach 4; Testing Approach 5; Build Order 4
- ADR-0058, ADR-0059, ADR-0060, ADR-0081, ADR-0099, ADR-0149, ADR-0186, ADR-0203
