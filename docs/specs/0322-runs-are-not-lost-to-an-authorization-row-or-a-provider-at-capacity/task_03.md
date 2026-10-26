---
task: task_03
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
status: pending
type: backend
complexity: medium
---

# Task 03: The Baseline clauses name the authorization row's subject

## Overview

Two Baseline clauses tell every adopter that an outside-evidence row the
Spec's authorization keeps offline is exempt only in the exact network
status. After task_01 the policy reads the form `blocked (environment:
<subject> prohibited by Spec authorization)`, so the portable rule must say
the same. This Task rewords both sentences in place, raises their rule
versions, regenerates the Baseline's derived files and this repository's
managed guides through the sanctioned commands, and updates the test that
pins the sentences.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md`](references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md)
   of 2026-10-26 (`_prd.md` → Sources 1) by implementing API Contract 6 of
   `_techspec.md`.
2. MUST replace, in the `guidance` of
   `clause.spec.project-constraints-06-outside-evidence`
   (`internal/baseline/assets/modules/spec-workflow.json`) and of
   `clause.autonomous.loop-01-qa-once`
   (`internal/baseline/assets/modules/autonomous-work.json`), only the
   sentence that names `blocked (environment: network prohibited by Spec
   authorization)` with the sentence API Contract 6 gives for that clause,
   and raise the `version` of `rule.spec.project-constraints` and of
   `rule.autonomous.loop` by one over the value the tree holds. The clauses
   keep their `id`, `enforcement` and `replaces`; nothing else in either
   module is edited by hand.
3. MUST then run, in this order, and hand-edit none of their outputs:
   - `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
     which writes the next version of both modules and records them in
     `internal/baseline/module-versions.json`;
   - `make baseline-digests`, which rewrites the derived files this Task
     declares (measured on a prototype: the two Standard TypeScript formatter
     goldens, the Standard TypeScript profile, the three catalog files and
     the four plan characterization goldens);
   - `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
     the Managed Refresh of this repository, which rewrites
     `docs/agents/spec-routing.md`, `docs/agents/autonomous-work.md` and
     `docs/agents/setup-context.json`, then the same command again, which must
     report no file change. Run both with `HOME` set to a temporary
     directory, so the refresh reads and writes no real Roundfix Home.
   A file the commands change that this Task does not declare is a
   stop-and-report condition.
4. MUST set the two sentence constants of
   `internal/baseline/network_prohibited_row_clause_test.go` to API
   Contract 6's sentences and change nothing else in that test.
5. MUST search the repository outside `docs/history/`, `docs/adr/`,
   `CHANGELOG.md` and this Spec for "authorization forbids network access",
   and report every match left; the skills and guides task_04 changes are the
   only expected ones.
6. MUST NOT change the Source Baseline corpus, manifest, accounting or index,
   `internal/baseline/assets/retention/**`, any other module, any skill, any
   user guide page, `.roundfixrc.yml`, the Makefile, the CI workflows or any
   Go production file.

## Subtasks

- [ ] Reword both clause sentences and raise both rule versions.
- [ ] Record the module versions, run `make baseline-digests` and the Managed Refresh twice.
- [ ] Update the clause test and sweep for the old wording.

## Acceptance Criteria

- [ ] The embedded clauses, the Standard TypeScript goldens and this repository's `spec-routing.md` and `autonomous-work.md` carry the subject form, and none carries the old sentence.
- [ ] A second Managed Refresh changes nothing, and the Setup Manifest test passes.
- [ ] The Baseline package's catalog, characterization and module-version tests pass.

## Consumers

- the two clause sentences and the catalog digests — consumers: `internal/baseline`, `internal/docscontract`

## Context

- instruction: `docs/adr/0335-a-stop-that-is-not-the-work-s-fault-does-not-cost-the-run.md`
- instruction: `internal/baseline/module_versions_test.go`
- instruction: `internal/baseline/catalog_test.go`
- instruction: `internal/baseline/plan_characterization_test.go`
- instruction: `internal/baseline/plan_test.go`
- instruction: `internal/baseline/compatibility_corpus_test.go`
- instruction: `internal/baseline/clause_characterization_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/network_prohibited_row_clause_test.go`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/autonomous-work.md`
- interface: `docs/agents/setup-context.json`

## Verification

- `for file in docs/agents/spec-routing.md docs/agents/autonomous-work.md internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'authorization forbids that lookup, recorded as' || { printf '%s lacks the subject form\n' "$file" >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'authorization forbids network access, recorded as'; then printf '%s keeps the old sentence\n' "$file" >&2; exit 1; fi; done; go test -timeout 15m -count=1 -run '^(TestTheGuidesExemptANetworkProhibitedOutsideEvidenceRow|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestCatalogDiagnosticCharacterization|TestBaselinePlanCharacterization|TestFormatterComposition|TestBaselineCompatibilityCorpus|TestNoTwoBaselineClausesShareText|TestBaselineClauseForceIsCharacterized)$' ./internal/baseline` — expected: exit 0. Before this Task the guides and goldens carry the old sentence, so the first check fails. Afterwards the four rendered copies carry the subject form only, and the clause test, the module record, the catalog and the derived-file characterizations pass.
- `tr -s '[:space:]' ' ' < internal/baseline/assets/modules/autonomous-work.json | grep -q -- 'recorded as .blocked (environment: <subject> prohibited by Spec authorization)., where .<subject>. names what the authorization forbids' || { echo 'autonomous-work module lacks the subject form' >&2; exit 1; }; go test -timeout 15m -count=1 -tags docscontract -run '^TestTheSetupManifestIsCurrent$' ./internal/docscontract` — expected: exit 0. Before this Task the module carries the network sentence, so the first check fails. Afterwards the module carries the subject form and this repository's Setup Manifest matches the regenerated guides.

## References

- `_prd.md` → Sources 1; Goal 5; Core Feature 3; Success Metric 7
- `_techspec.md` → Current behavior 8; API Contract 6; Testing Approach; Build Order 3; Risks & Considerations
- ADR-0335; ADR-0276
