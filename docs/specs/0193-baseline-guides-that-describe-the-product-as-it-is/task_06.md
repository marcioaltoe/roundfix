---
task: task_06
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: pending
type: backend
complexity: medium
---

# Task 06: A removed Source Baseline clause is replaced, not lost

## Overview

task_01 removed the duplicate entry `rule.backend.boundary-contracts` from `backend.json`. The Standard TypeScript Source Baseline records that entry, and `classifySourceClauseTransition` in `internal/baseline/plan.go` knows only `retained` and `unaccounted`. An adopter whose Setup Manifest declares that Source Baseline now gets the removed entry reported `unaccounted`, and its next Baseline update is refused ("retention transition has 1 unaccounted clause(s)"). The first delivery Run's QA gate failed on `TestStandardTypeScriptStructuralClauseRetention`, the test that guards this. The same gate failed on `TestADRLifecycleContract`, which still pins the ADR sentence task_04 scoped on purpose. This Task gives the removal a declared successor and moves both contracts to the facts this Spec establishes. The TechSpec section "Clause replacement for adopters (task_06)" fixes the design.

## Requirements

1. MUST accept on a module clause an optional `replaces` field: a list of clause IDs that clause takes over.
2. MUST make `classifySourceClauseTransition` report a Source Baseline normative clause that is absent from the selected catalog as `replaced` (the existing `ClauseReplaced`), with the replacing clause's ID as its only target and a reason naming it, when exactly one selected clause lists it in `replaces` and that clause's enforcement equals the Source Baseline entry's. Every other absence MUST stay `unaccounted`, with today's reason. Retained clauses MUST keep today's classification and reason.
3. MUST make catalog validation refuse, each with its own diagnostic code under `catalog.clause.replaces.`: a `replaces` value that is not a list of unique non-empty strings; an ID that is still a clause ID anywhere in the catalog; an ID claimed by two clauses.
4. MUST add `"replaces": ["rule.backend.boundary-contracts"]` to `clause.backend.boundary-contracts` in `internal/baseline/assets/modules/backend.json`, change nothing else in that clause, and render no guide differently.
5. MUST change `TestStandardTypeScriptStructuralClauseRetention` in `internal/baseline/plan_test.go` so that `rule.backend.boundary-contracts` is expected `replaced`, targeting `clause.backend.boundary-contracts`, whose enforcement and guidance equal the Source Baseline entry's. Every other clause the test lists MUST stay expected `retained`, with its current checks.
6. MUST change the fragment `Only \`accepted\` is active.` in `TestADRLifecycleContract` to `For an ADR that carries lifecycle frontmatter, only \`accepted\` is active.`, the sentence task_04 wrote. It MUST change no other fragment or fixture of that test.
7. MUST add `internal/baseline/clause_replacement_test.go` with:
   - `TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor`: a Standard TypeScript adopter manifest whose managed artifacts drifted produces a plan that is not refused, and whose delta records `rule.backend.boundary-contracts` as `replaced`;
   - `TestAnUndeclaredRemovalStaysUnaccounted`: the same catalog with the `replaces` declaration removed reports `unaccounted`, and the plan is refused;
   - `TestAReplacementWithOtherForceStaysUnaccounted`: a declaration whose clause has a different enforcement reports `unaccounted`;
   - `TestClauseReplacementDeclarationsAreValidated`: the three refusals of Requirement 3, each by its code.
8. MUST regenerate the catalog snapshots, digest pins and plan goldens only through their sanctioned commands (`make baseline-digests` and the package's golden update path), never by hand, and MUST leave a second refresh reporting no change.
9. MUST NOT edit the Source Baseline corpus or index, any other governed test file, or any rendered guide.

## Subtasks

- [ ] Characterize: run the two failing tests and record the failure.
- [ ] Add the `replaces` producer and its validation.
- [ ] Declare the successor in `backend.json` and regenerate derived files.
- [ ] Move the two plan contracts and add the four new tests.

## Acceptance Criteria

- [ ] The whole `./internal/baseline` package passes, including the two tests the first Run failed.
- [ ] Removing the `replaces` declaration makes `TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor` fail.
- [ ] `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json` exits 0.

## Context

- instruction: `docs/specs/0193-baseline-guides-that-describe-the-product-as-it-is/_techspec.md`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/clause_replacement_test.go`

## Verification

- `out="$(go test -count=1 -v ./internal/baseline 2>&1)" || { printf "%s\\n" "$out" | grep -E -- "--- FAIL|_test.go:" ; exit 1; }; for name in TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor TestAnUndeclaredRemovalStaysUnaccounted TestAReplacementWithOtherForceStaysUnaccounted TestClauseReplacementDeclarationsAreValidated TestStandardTypeScriptStructuralClauseRetention TestADRLifecycleContract TestBaselinePlanCharacterization TestCatalogCompatibility; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the package fails on `TestStandardTypeScriptStructuralClauseRetention` and `TestADRLifecycleContract`, and the four new tests do not exist.

## References

- task_01, task_04
- `_techspec.md` → Clause texts; Build Order
- QA Report of Run `run_20260930T153457Z_b5736250c995b99a` (verification log `batch-005-attempt-1.log`)
