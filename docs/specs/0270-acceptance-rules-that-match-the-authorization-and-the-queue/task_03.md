---
task: task_03
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
status: pending
type: docs
complexity: medium
---

# Task 03: The Baseline guides and the glossary name both acceptance rules

## Overview

Two Baseline clauses list the outside-evidence rows that never decide a
qualifying partial, and adopters read those guides before the skills. This
Task appends the two sentences of `_techspec.md` → Exact texts to them,
regenerates the derived Baseline files and this repository's guides with the
repository's own commands, and writes the glossary entries.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-13,
   [outside evidence blocked by the Spec's own authorization does not qualify](references/2026-10-13-outside-evidence-blocked-by-the-spec-authorization-does-not-qualify.md),
   by appending the `_techspec.md` → Exact texts sentence to the guidance of
   `clause.spec.project-constraints-06-outside-evidence` in
   `internal/baseline/assets/modules/spec-workflow.json`, after its last
   sentence, and the second sentence to `clause.autonomous.loop-01-qa-once` in
   `internal/baseline/assets/modules/autonomous-work.json`, after the sentence
   that names the network-denied row. Every existing sentence stays
   byte-identical, so `TestAuthorizationClausesStateTodaysRefusals` and
   `TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow` keep passing
   unchanged. No clause is removed or renamed, so no Source Baseline retention
   disposition applies.
2. MUST raise each changed module to the next version with
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then run `make baseline-digests` twice (the second reports
   `"changed":false`), then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   twice with `HOME` set to a temporary directory; the second refresh MUST
   report `File changes: 0`. The authoring prototype measured the regenerated
   set: the two module version lines, `internal/baseline/module-versions.json`,
   the profile, the two formatter goldens, the three catalog files, the four
   plan characterization goldens, `docs/agents/spec-routing.md`,
   `docs/agents/autonomous-work.md` and `docs/agents/setup-context.json`. MUST
   NOT hand-edit a snapshot, golden, digest, record or generated guide.
3. MUST add `internal/baseline/network_prohibited_row_clause_test.go` with
   `TestTheGuidesExemptANetworkProhibitedOutsideEvidenceRow`, which asserts
   both new sentences in the embedded clauses, in the two Standard TypeScript
   Monorepo formatter goldens and in this repository's two guides, on
   whitespace-normalized text. It calls `t.Parallel()` first.
4. MUST write, through `domain-modeling`, the **Network-Prohibited Row** and
   **Pinned Byte Invariant** entries of `_techspec.md` → Exact texts into
   `CONTEXT.md` after the Network-Denied Row entry, and append the Delivery
   Base sentence to the existing **Delivery Base** entry.
5. MUST NOT change any other clause, rule or module, any skill,
   `CHANGELOG.md` or the Source Baseline corpus.

## Subtasks

- [ ] Append the two sentences and record the module versions.
- [ ] Regenerate the derived Baseline files and run the Managed Refresh twice.
- [ ] Add the clause test.
- [ ] Write the glossary entries.

## Acceptance Criteria

- [ ] Both generated guides and both formatter goldens name the
      network-prohibited outside-evidence row and keep every existing sentence.
- [ ] A second `make baseline-digests` and a second Managed Refresh change
      nothing.
- [ ] `CONTEXT.md` defines both new terms and the changed Delivery Base.

## Consumers

- The Baseline catalog and the guides it renders — consumers: `./internal/baseline`

## Context

- instruction: `docs/adr/0276-a-row-the-authorization-keeps-offline-qualifies-and-a-byte-invariant-compares-with-the-delivery-base.md`
- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/module-versions.json`
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
- interface: `CONTEXT.md`
- creates: `internal/baseline/network_prohibited_row_clause_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheGuidesExemptANetworkProhibitedOutsideEvidenceRow|TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow|TestAuthorizationClausesStateTodaysRefusals|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestNoTwoBaselineClausesShareText|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheGuidesExemptANetworkProhibitedOutsideEvidenceRow TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow TestAuthorizationClausesStateTodaysRefusals TestCatalogCompatibility TestBaselinePlanCharacterization TestEveryBaselineModuleVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for file in docs/agents/spec-routing.md docs/agents/autonomous-work.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'blocked (environment: network prohibited by Spec authorization)' || { printf 'missing sentence in %s\n' "$file" >&2; exit 1; }; done; glossary="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; for phrase in '**Network-Prohibited Row**:' '**Pinned Byte Invariant**:' '**Delivery Base**:' 'compares it with the Delivery Base, never with a named commit.'; do printf '%s\n' "$glossary" | grep -qF -- "$phrase" || { printf 'CONTEXT.md lacks %s\n' "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task the clause test does not exist and neither guide nor the glossary names the new row, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 2; User Stories 4; Core Features 4; Success Metric 5
- [_techspec.md](_techspec.md) — Exact texts; API Contract 4; Invariants 6; Build Order 3
- ADR-0276; ADR-0240; ADR-0104; ADR-0244; ADR-0259; ADR-0271
