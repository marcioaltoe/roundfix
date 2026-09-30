---
task: task_03
spec: 0199-stack-rules-that-say-what-they-mean
status: pending
type: backend
complexity: medium
---

# Task 03: Backend and frontend guides name their workspace, and the stated HTTP suggestion is REST

## Overview

The backend and frontend guides do not say which part of a repository they govern, and their clauses cannot be reworded: the Source Baseline retains them by their text. This Task adds one scope sentence to each guide's template and leaves every clause byte-identical. It also aligns three stale statements with the decision catalog. The catalog suggests REST for the HTTP Contract Decision, and the prompt and the update already read it, while the profile's declared default, the Baseline contract's sentence and the public guide's table still say Post-only. Two preservation tests prove that a repository with a recorded decision keeps it and that a repository with none is only offered REST.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST insert, in `internal/baseline/assets/templates/guides/backend.md` and `internal/baseline/assets/templates/guides/frontend.md`, the scope paragraph the TechSpec's "Exact texts" gives for each, immediately before the rules token.
2. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_03. In `internal/baseline/assets/modules/backend.json` and `internal/baseline/assets/modules/frontend.json` only those version numbers change: every clause stays byte-identical.
3. MUST change the profile's `httpContract.default` to `REST`, the sentence in `internal/baseline/assets/contract-v1.json` and the HTTP contract row of the suggested-values table in `docs/user-guide/context-driven-development.md`, each as the TechSpec gives it. MUST NOT change `internal/baseline/assets/decisions.json`, the prompt, or how a Setup Manifest decision is resolved.
4. MUST create `internal/baseline/stack_scope_and_http_default_test.go` with `profileHTTPDefaultFindings` and the five tests the TechSpec's Testing Approach 3 and 4 name for that file. The recorded-decision test MUST cover both modes and MUST compare the recorded exceptions.
5. MUST create `internal/cli/baseline_http_default_statement_test.go` with `httpDefaultStatementFindings` and the two tests the TechSpec's Testing Approach 4 names. The statement test MUST read the mode from the embedded decision catalog, not from a literal.
6. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep each JSON file's existing
   formatting by replacing strings and version numbers in place.
7. MUST NOT edit `internal/baseline/plan_test.go`, `internal/cli/baseline_documentation_contract_test.go`, `internal/cli/baseline_human_test.go` or any other existing test file. The new tests call their helpers from the new files.
8. MUST NOT change the `auth.provider` decision, its suggested route scope, or any other row of the public guide's table.

## Subtasks

- [ ] Add the two scope paragraphs and raise the template, guide and module versions.
- [ ] Align the profile default, the contract sentence and the public guide's row.
- [ ] Create the scope test, the profile default check and the two preservation tests.
- [ ] Create the statement check and its negative test.
- [ ] Regenerate the pins, the two goldens and the Setup Manifest.

## Acceptance Criteria

- [ ] The rendered backend and frontend guides each carry their scope sentence before their rules, and the structural-clause retention test passes unchanged.
- [ ] The profile default, the contract sentence and the public guide's row name the decision catalog's default mode, and a statement that names the other mode is reported.
- [ ] A Setup Manifest that records `Post-only` resolves to `Post-only` with its exceptions unchanged and no new decision; the same holds for `REST`.
- [ ] A Setup Manifest with no HTTP Contract Decision reports it as a new decision suggested as `REST` and adopts nothing.
- [ ] A second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0190-a-stack-rule-names-the-language-or-workspace-it-governs.md`
- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- interface: `internal/baseline/assets/templates/guides/backend.md`
- interface: `internal/baseline/assets/templates/guides/frontend.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/modules/frontend.json`
- interface: `internal/baseline/assets/contract-v1.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/user-guide/context-driven-development.md`
- creates: `internal/baseline/stack_scope_and_http_default_test.go`
- creates: `internal/cli/baseline_http_default_statement_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheBackendAndFrontendGuidesSayWhatTheyGovern|TestTheProfileHTTPDefaultMatchesTheDecisionCatalog|TestAProfileHTTPDefaultThatDisagreesIsReported|TestARecordedHTTPContractDecisionSurvivesAnUpdate|TestAnAbsentHTTPContractDecisionIsSuggestedAsREST|TestStandardTypeScriptStructuralClauseRetention|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheBackendAndFrontendGuidesSayWhatTheyGovern TestTheProfileHTTPDefaultMatchesTheDecisionCatalog TestAProfileHTTPDefaultThatDisagreesIsReported TestARecordedHTTPContractDecisionSurvivesAnUpdate TestAnAbsentHTTPContractDecisionIsSuggestedAsREST TestStandardTypeScriptStructuralClauseRetention TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && golden=internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents && for pair in "$golden/backend.md|These rules govern the repository's TypeScript backend workspace." "$golden/frontend.md|These rules govern the repository's web frontend workspace." "internal/baseline/assets/contract-v1.json|the interactive workflow suggests REST until"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && if tr -s '[:space:]' ' ' < internal/baseline/assets/contract-v1.json | grep -qF -- "the interactive workflow suggests Post-only until"; then printf 'stale HTTP suggestion in contract-v1.json\n' >&2; exit 1; fi && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the five new named tests exists and the contract still suggests Post-only, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault|TestAStatementThatNamesAnotherHTTPDefaultIsReported)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault TestAStatementThatNamesAnotherHTTPDefaultIsReported; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither named test exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Exact texts; Version changes; Left as they are, deliberately; Testing Approach 3; Testing Approach 4; Testing Approach 5; Build Order 3
- `_prd.md` → Goal 2; Goal 3; Goal 5; Core Feature 3; Core Feature 4; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8
- `_techspec.md` → API Contract 3; API Contract 4
- ADR-0058, ADR-0059, ADR-0060, ADR-0061, ADR-0063, ADR-0186, ADR-0190
