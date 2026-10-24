---
task: task_01
spec: 0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff
status: pending
type: backend
complexity: medium
---

# Task 01: A Common Setup Snapshot keeps only the skills its components share

## Overview

A Setup Snapshot today mirrors one upstream list or is the union of several
(the composed kind). A toolchain without an upstream list, such as Python,
needs the skills its sibling lists share. This Task adds the common kind: the
catalog validates it beside the composed kind and the asset sync writes it
after its components. No embedded snapshot uses it yet; task_02 adds the
`python` one (`_prd.md` → Core Feature 1; ADR-0322).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-22-a-baseline-profile-for-python-repositories.md`](references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected item 5 (the profile works with the skills that
   exist today).
2. MUST add to `internal/baseline/setup_composition.go` the functions
   `derivedSetupSource`, `deriveSetupSnapshot` and `commonSetupSnapshot` of
   `_techspec.md` → Interfaces, with the semantics of API Contract 1 and the
   document shape of Data Model 1; `setupCompositionComponents` MUST refuse a
   component of either derived kind and keep its other refusals, and
   `composeSetupSnapshot` MUST keep its behavior.
3. MUST make `validateSetups` in `internal/baseline/catalog_validate.go` check
   every derived snapshot through `derivedSetupSource` and
   `deriveSetupSnapshot`, reporting the existing codes; the drift detail for a
   common snapshot MUST be exactly "skills or activation bundles differ from
   the components' common part", and the composed detail MUST keep its bytes.
4. MUST make `syncAssets` in `internal/baseline/assets_sync.go` skip both
   derived kinds in its upstream refresh loop and derive each one after it
   through `deriveSetupSnapshot`, keeping its findings and transaction.
5. MUST make the three helpers of `internal/baseline/assets_sync_test.go` that
   skip a `composed` snapshot (the template source check, the maintained
   parity row and `buildAssetsSyncSource`) skip a `common` one too, and change
   nothing else in that file.
6. MUST create `internal/baseline/common_setup_test.go` with the three tests
   of `_techspec.md` → Testing Approach (Build Order 1), using the existing
   helpers `cloneEmbeddedAssets`, `writeCompositionFixture`,
   `requireCatalogDiagnostic`, `catalogDiagnosticByCode`,
   `newAssetsSyncTarget`, `newAssetsSyncSource`, `readAssetsSyncJSON` and
   `writeAssetsSyncJSON`; a fixture written through `writeCompositionFixture`
   first creates its map entry. Each top-level test and subtest calls
   `t.Parallel()` first.
7. MUST keep every existing test of `internal/baseline` passing without other
   edits, including the composed setup tests of `setup_composition_test.go`.
8. MUST NOT add any embedded Setup Snapshot, module, profile or template, nor
   change any asset under `internal/baseline/assets/`.

## Subtasks

- [ ] Add the derivation functions and the derived-component refusal.
- [ ] Validate and sync both derived kinds.
- [ ] Extend the helper skips and add the three tests.

## Acceptance Criteria

- [ ] The common derivation keeps only equal shared entries in the first component's order, with a `common` source and the canonical digest.
- [ ] A drifted common snapshot is reported with the common detail; one component or a derived component is invalid.
- [ ] The asset sync rewrites a stale common snapshot from its refreshed components.

## Consumers

- `commonSetupSnapshot`, `deriveSetupSnapshot` — consumers: `./internal/baseline`

## Context

- instruction: `docs/adr/0322-a-python-profile-takes-the-skills-the-go-and-rust-setups-share.md`
- instruction: `docs/adr/0204-a-composed-profile-takes-a-setup-composed-from-upstream-setups-by-name.md`
- instruction: `internal/baseline/setup_composition_test.go`
- interface: `internal/baseline/setup_composition.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/assets_sync_test.go`
- creates: `internal/baseline/common_setup_test.go`

## Verification

- `for name in TestACommonSetupKeepsOnlyTheSkillsEveryComponentCarries TestACommonSetupThatDiffersFromItsComponentsIsReported TestTheAssetSyncWritesACommonSetupAfterItsComponents; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^(TestACommonSetupKeepsOnlyTheSkillsEveryComponentCarries|TestACommonSetupThatDiffersFromItsComponentsIsReported|TestTheAssetSyncWritesACommonSetupAfterItsComponents|TestTheComposedSetupIsTheUnionOfItsComponents|TestAComposedSetupThatDriftsFromItsComponentsIsRefused|TestAComposedSetupWithAnUnknownOrComposedComponentIsRefused|TestComponentsThatDisagreeOnASkillAreAConflict|TestAnAssetSyncRewritesTheComposedSetupWithItsComponents|TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract|TestCatalogCompatibility)$' ./internal/baseline && GOOS=windows go vet ./internal/baseline` — expected: exit 0. Before this Task none of the three tests exists, so the command fails at its first check. Afterwards they exist and pass with the composed setup tests, the asset-sync tests and the catalog, and the package's tests compile for Windows.

## References

- `_prd.md` → Goals; Core Feature 1; Success Metric 3; User Story 4
- `_techspec.md` → Current behavior 1-3; Interfaces; Data Model 1; API Contract 1; Invariants 1-2; Testing Approach; Build Order 1
- ADR-0322; ADR-0204; ADR-0191
