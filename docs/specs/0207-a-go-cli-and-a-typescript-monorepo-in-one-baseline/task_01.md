---
task: task_01
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: pending
type: backend
complexity: high
---

# Task 01: A repository records its frontend layout

## Overview

The frontend guide makes the systems layout mandatory. This Task adds the optional Frontend Layout Decision `frontend.layout`, lets a clause name the decision value it applies under, renders the layout sentence for each state, and accounts for a clause a recorded value turns off as a reasoned rejection. The two systems clauses keep their identifiers, enforcement and bytes, so a repository that records nothing, including every existing adopter and every repository-owned profile, keeps its clauses and gains one sentence.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST characterize first: before any edit, record in the Result which existing tests in `./internal/baseline` and `./internal/cli` pass, so a later failure is attributable.
2. MUST add the decision field `optional` and the clause field `appliesWhen` with the semantics and diagnostics the TechSpec's Interfaces and Data Models give: `normalizePlanDecisions` never reports an optional decision missing; `clauseApplies` selects a gated clause under the recorded value, or under the decision's `default` when an optional decision is unrecorded; `selectedClauseEnforcement` and `artifactRenderValues` skip a clause that does not apply; an optional decision's `renderBindings` do not make their artifact decision-controlled, and its token renders `renderUnrecordedProjectDecision` while unrecorded.
3. MUST extend `classifySourceClauseTransition` so a prior managed clause that is absent only because a recorded decision value turns it off gets the `reasoned-rejection` disposition, with the decision as its target and the reason the TechSpec gives. A prior clause absent for any other reason MUST stay `unaccounted`.
4. MUST apply the TechSpec's "Fixed texts" for task_01 to `decisions.json`, `frontend.json`, the frontend guide template and `templates/index.json`, and the Data Models change to the Standard TypeScript Monorepo Profile (`entryDecisions` and `architecture.frontend`). The guidance of `clause.frontend.organize-by-system` and `clause.frontend.public-system-boundary` MUST stay byte-identical.
5. MUST add the three layout sentences of Fixed texts to `project_decision_render.go`, taking the suggested value from the decision's `default`.
6. MUST create `internal/baseline/frontend_layout_decision_test.go` with the nine tests the TechSpec's Testing Approach 1 names. The refusal tests MUST load an overlay of the embedded assets and assert the diagnostic code; the reasoned-rejection test MUST plan against a Setup Manifest that held both systems clauses and now records `repository-defined`; the unaccounted test MUST remove a clause with no recorded decision and see `unaccounted`.
7. MUST add `clause.frontend.follow-recorded-layout` as `mandatory` to the table in `internal/baseline/clause_characterization_test.go`, and the row `{name: "frontend layout", id: "frontend.layout", want: "systems"}` to `TestHumanBaselineDecisionDefaults`. An existing scripted interview in `internal/cli/baseline_human_test.go` that adopts a profile with the frontend module gains one answer for the new question and nothing else.
8. MUST update `docs/user-guide/context-driven-development.md` as Fixed texts gives for task_01.
9. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`. MUST NOT hand-edit a pin, golden or generated guide, and MUST keep each asset's formatting by replacing strings and numbers in place.
10. MUST NOT rename or remove a clause, a top-level test or an exported function, and MUST NOT change the Source Baseline corpus, the retention transitions or the parity corpus. If an existing test outside the Context fails, the Task stops and reports it, unless the same test also fails on this Task's starting commit in the same environment (for example the process-table tests in `internal/cli/orphan_unix_test.go`, which a sandboxed Agent cannot run): such a test is recorded in the Result with that proof and does not stop the Task.

## Subtasks

- [ ] Characterize the passing tests.
- [ ] Add optional decisions, clause gates, the layout render and the retention disposition.
- [ ] Add the decision, gate the systems clauses, add the repository-layout clause and the template token.
- [ ] Write the nine tests and update the two characterization tables.
- [ ] Update the public guide, regenerate, and refresh this repository twice.

## Acceptance Criteria

- [ ] With no recorded layout, the Standard TypeScript Monorepo frontend guide renders the suggestion sentence and the six clauses it renders today; with `systems`, the recorded sentence and the same clauses; with `repository-defined`, the recorded sentence, the repository-layout clause and neither systems clause.
- [ ] An unrecorded `frontend.layout` is not a missing decision, and `baseline update` of a manifest without it asks nothing.
- [ ] An optional decision without a default and a clause gate on a required decision are refused with their codes.
- [ ] A clause turned off by a recorded value is a reasoned rejection naming `frontend.layout`; a clause missing for no recorded reason is still unaccounted.
- [ ] Both systems clauses are still retained by `TestStandardTypeScriptStructuralClauseRetention`, and a second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0205-a-repository-records-its-frontend-layout-and-an-unrecorded-layout-follows-the-suggestion.md`
- instruction: `docs/adr/0058-baseline-upgrades-fail-closed-on-unaccounted-rule-removal.md`
- interface: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/assets/modules/frontend.json`
- interface: `internal/baseline/assets/templates/guides/frontend.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md`
- interface: `internal/baseline/catalog_load.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/project_decision_render.go`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/custom_profile_test.go`
- interface: `internal/cli/baseline_human_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/frontend_layout_decision_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAnUnrecordedFrontendLayoutStatesTheSuggestionAndKeepsTheSystemsClauses|TestARecordedSystemsLayoutKeepsTheSystemsClauses|TestARecordedRepositoryDefinedLayoutRendersOnlyItsOwnClause|TestAnUnrecordedOptionalDecisionIsNotMissing|TestAnOptionalDecisionWithoutADefaultIsRefused|TestAClauseGateOnARequiredDecisionIsRefused|TestAClauseARecordedLayoutTurnsOffIsAReasonedRejection|TestAClauseMissingWithoutARecordedDecisionIsStillUnaccounted|TestTheProfileStatesTheLayoutSuggestionTheCatalogDefaults|TestStandardTypeScriptStructuralClauseRetention|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAnUnrecordedFrontendLayoutStatesTheSuggestionAndKeepsTheSystemsClauses TestARecordedSystemsLayoutKeepsTheSystemsClauses TestARecordedRepositoryDefinedLayoutRendersOnlyItsOwnClause TestAnUnrecordedOptionalDecisionIsNotMissing TestAnOptionalDecisionWithoutADefaultIsRefused TestAClauseGateOnARequiredDecisionIsRefused TestAClauseARecordedLayoutTurnsOffIsAReasonedRejection TestAClauseMissingWithoutARecordedDecisionIsStillUnaccounted TestTheProfileStatesTheLayoutSuggestionTheCatalogDefaults TestStandardTypeScriptStructuralClauseRetention TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && cli="$(go test -count=1 -v -run '^TestHumanBaselineDecisionDefaults$' ./internal/cli 2>&1)" || { printf "%s\\n" "$cli"; exit 1; }; printf "%s\\n" "$cli" | grep -q -- "--- PASS: TestHumanBaselineDecisionDefaults/frontend_layout" || { printf 'missing pass: frontend layout default\n' >&2; exit 1; }; golden=internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md && for phrase in "No frontend layout is recorded." "Organize frontend feature code by domain system." "Import another system through that system's public boundary"; do tr -s '[:space:]' ' ' < "$golden" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$golden" "$phrase" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the nine new tests and the `frontend_layout` subtest do not exist and the golden has no layout sentence, so the command fails.

## References

- [_techspec.md](_techspec.md) — Interfaces; Data Models; Fixed texts; Testing Approach 1; Testing Approach 5; Build Order 1
- `_prd.md` → Goal 4; Goal 5; Story 5; Story 6; Core Feature 5; Core Feature 6; Success Metric 5; Success Metric 6; Success Metric 7
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3
- ADR-0058, ADR-0060, ADR-0063, ADR-0067, ADR-0081, ADR-0099, ADR-0149, ADR-0186, ADR-0205
