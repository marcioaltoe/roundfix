---
task: task_04
spec: 0235-one-qa-partial-policy
status: completed
type: docs
complexity: medium
---

# Task 04: The Spec routing and autonomous-work guides exempt a network-denied outside-evidence row

## Overview

Two Baseline clauses state the old rule to every adopter. The outside-evidence
clause says the QA gate "holds Pull Request preparation until the row is
satisfied", and the delivery-order clause names only the Pull Request row as
exempt. Adopters read those guides before the skills. This Task appends one
sentence to each clause, as the TechSpec's "Exact texts" gives it. It then
regenerates the derived Baseline files and this repository's guides with
the repository's own commands.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-06
   ([settlement refuses a partial that archive accepts](references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md))
   by appending the TechSpec's "Exact texts" sentence to the guidance of
   `clause.spec.project-constraints-06-outside-evidence` in
   `internal/baseline/assets/modules/spec-workflow.json`, after its existing
   last sentence, and the second sentence to the delivery-order clause of
   `internal/baseline/assets/modules/autonomous-work.json`, after
   "never decides a qualifying `partial`.". Every existing sentence stays
   byte-identical, so `TestAuthorizationClausesStateTodaysRefusals` keeps
   passing unchanged and adopters keep their current text. No clause is
   removed or renamed, so no Source Baseline retention disposition applies.
2. MUST add `internal/baseline/network_denied_row_clause_test.go` with
   `TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow`. It asserts both
   added sentences in the embedded clauses, in the Standard TypeScript
   Monorepo formatter goldens of `spec-routing.md` and `autonomous-work.md`,
   and in this repository's `docs/agents/spec-routing.md` and
   `docs/agents/autonomous-work.md`.
3. MUST run `make baseline-digests` twice (the second reports
   `"changed":false`), then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
   twice; the second refresh MUST report `File changes: 0`. A probe in a
   disposable clone at `bde616d1` showed the regenerated set: the two
   formatter goldens, the profile, `catalog.diagnostics.golden.json`,
   `catalog.digest`, `catalog.normalized.json`, the four plan
   characterization goldens, `docs/agents/spec-routing.md`,
   `docs/agents/autonomous-work.md` and `docs/agents/setup-context.json`.
   MUST NOT hand-edit a snapshot, golden, digest or generated guide.
4. MUST NOT change any other clause, rule or module version, any skill,
   `CONTEXT.md` or `CHANGELOG.md`.

## Subtasks

- [ ] Append the two sentences to their clauses.
- [ ] Regenerate the derived Baseline files and run the Managed Refresh.
- [ ] Add the clause test.

## Acceptance Criteria

- [ ] Both generated guides and both formatter goldens state that a
      network-denied outside-evidence row never decides a qualifying
      `partial`, and keep their existing sentences.
- [ ] A second `make baseline-digests` and a second Managed Refresh change
      nothing.

## Context

- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- instruction: `docs/adr/0104-a-spec-accepts-on-evidence-it-did-not-author.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
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
- creates: `internal/baseline/network_denied_row_clause_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow|TestAuthorizationClausesStateTodaysRefusals|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestNoTwoBaselineClausesShareText)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow TestAuthorizationClausesStateTodaysRefusals TestCatalogCompatibility TestBaselinePlanCharacterization; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for file in docs/agents/spec-routing.md docs/agents/autonomous-work.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'blocked only because the Run sandbox denied network access' || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done` — expected: exit 0; before this Task the clause test does not exist and neither guide names the network-denied row, so the command fails; after it the clause test, the catalog and plan characterizations and the pinned clause test pass, and both guides carry the sentence.

## References

- `_prd.md` → Goals; Core Features 6; Success Metric 5
- `_techspec.md` → Exact texts; System Architecture; Build Order 4
- ADR-0240; ADR-0104

## Result

Implemented the two Baseline clause additions from the TechSpec and added
`TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow`, which checks the
embedded clauses, both Standard TypeScript Monorepo formatter goldens, and
both repository guides.

Focused evidence from this turn:

- `rtk make baseline-digests`: regenerated the expected derived catalog,
  profile, plan characterizations, and formatter goldens.
- Second `rtk make baseline-digests`: passed with `changed:false`.
- `GOCACHE=/private/tmp/roundfix-0235-gocache go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`:
  verified and applied the three managed guide/setup changes.
- Second managed refresh with the same command: verified with `File changes: 0`
  and `Idempotence: verified`.
- `gofmt -w internal/baseline/network_denied_row_clause_test.go` followed by
  `GOCACHE=/private/tmp/roundfix-0235-gocache go test -count=1 -run
  '^TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow$' ./internal/baseline`:
  passed.

Acceptance evidence:

- The added clause test covers both generated guides and both formatter
  goldens, including the network-denied outside-evidence wording, while the
  source clauses retain their existing sentences.
- The second digest regeneration reported `changed:false`, and the second
  managed refresh reported `File changes: 0` with idempotence verified.

## Carry-forward provenance

- Source Run: `run_20261006T120012Z_3cae1c6ab9ca3edb`
- Source commit: `9d55539ccfad9ff420f2184b7e481f81ed2b17dd`
