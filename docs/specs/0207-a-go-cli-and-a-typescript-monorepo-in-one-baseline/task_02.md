---
task: task_02
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: pending
type: backend
complexity: high
---

# Task 02: A setup snapshot composed from upstream setups by name

## Overview

No upstream skill list covers a Go CLI beside a TypeScript monorepo, and a profile takes one Setup Snapshot. This Task lets a Setup Snapshot be composed from named component snapshots, has the asset sync write it in the same run as its components, has the catalog refuse a composition that is invalid, conflicting or stale, and creates the `go-cli-typescript-bun` snapshot. No profile selects it yet; task_03 does.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST create `internal/baseline/setup_composition.go` with `composeSetupSnapshot` as the TechSpec's Interfaces give: skills and activation bundles in component order, a later entry equal to an earlier one of the same name or bundle identifier dropped, an unequal one an error.
2. MUST extend `validateSetups` in `internal/baseline/catalog_validate.go` for a snapshot whose `source` is `{"type": "composed", "setups": [...]}`, reporting `catalog.setup.composition.invalid`, `catalog.setup.composition.conflict` and `catalog.setup.composition.drift` as the TechSpec's Data Models define. The per-skill checks and the digest rule stay as they are.
3. MUST extend the asset sync so that it skips a composed snapshot when it builds snapshots from upstream files, then composes each composed snapshot from the built or unchanged components and plans, checks or writes it like any other snapshot. `assetsSyncSource` gains `Setups`. A composition error MUST fail the run with the invalid-assets finding and write nothing.
4. MUST update the existing tests in `internal/baseline/assets_sync_test.go` whose counts include every setup, and MUST make the parity comparison skip composed snapshots with a comment that names the designed delta ADR-0072 requires. `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` MUST stay byte-identical.
5. MUST create `internal/baseline/setup_composition_test.go` with the six tests the TechSpec's Testing Approach 2 names. The sync tests MUST use temporary asset roots and temporary Git sources only.
6. MUST create `internal/baseline/assets/setups/go-cli-typescript-bun.json` by the TechSpec's "The composed setup procedure": the seed, then the asset sync against the commit the `go` snapshot pins. The run MUST change no other setup file. MUST add the new file to the sanctioned outputs `internal/baseline/derived_ownership_test.go` enumerates, and change no other line of that test.
7. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden or the composed snapshot's skills.
8. MUST NOT change any profile, module, decision or other setup snapshot, and MUST NOT rename or remove a top-level test or an exported function.
9. MUST add `Setups []string` with the JSON tag `setups,omitempty` to `baselineSetupSource` in `skills/baseline_skill_contract_test.go`, so the skill sync step of `make baseline-digests` keeps a composed setup's `source.setups` and the regeneration no longer fails with `catalog.setup.composition.invalid`.
10. MUST make the `replaceSetupSkillDigest` helper in `internal/baseline/catalog_test.go` recompute the composed setup when it edits a component setup, so `TestCatalogDigestExcludesOwnedSkillContent` keeps its owned and external digest assertions unchanged.

## Subtasks

- [ ] Write the composition and its validation.
- [ ] Teach the asset sync to write composed snapshots and adjust the counted tests and the parity skip.
- [ ] Write the six tests.
- [ ] Seed and sync the composed snapshot, register it as a sanctioned output, regenerate, and refresh this repository twice.

## Acceptance Criteria

- [ ] The embedded `go-cli-typescript-bun` snapshot equals the union of `go` and `typescript` in that order, and the catalog loads.
- [ ] A drifted composed snapshot, one with an unknown or composed component, and components that disagree on one skill are each refused with their own code.
- [ ] A sync whose source changes a component rewrites the composed snapshot in the same run; a sync with no source change leaves it byte-identical; `--check` reports its drift.
- [ ] The parity fixture is byte-identical and the parity comparison still passes.

## Context

- instruction: `docs/adr/0204-a-composed-profile-takes-a-setup-composed-from-upstream-setups-by-name.md`
- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/assets_sync_test.go`
- interface: `internal/baseline/derived_ownership_test.go`
- creates: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `skills/baseline_skill_contract_test.go`
- interface: `internal/baseline/catalog_test.go`
- creates: `internal/baseline/setup_composition.go`
- creates: `internal/baseline/setup_composition_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheComposedSetupIsTheUnionOfItsComponents|TestAComposedSetupThatDriftsFromItsComponentsIsRefused|TestAComposedSetupWithAnUnknownOrComposedComponentIsRefused|TestComponentsThatDisagreeOnASkillAreAConflict|TestAnAssetSyncRewritesTheComposedSetupWithItsComponents|TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup|TestCatalogCompatibility|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheComposedSetupIsTheUnionOfItsComponents TestAComposedSetupThatDriftsFromItsComponentsIsRefused TestAComposedSetupWithAnUnknownOrComposedComponentIsRefused TestComponentsThatDisagreeOnASkillAreAConflict TestAnAssetSyncRewritesTheComposedSetupWithItsComponents TestAnAssetSyncWithNoSourceChangeLeavesTheComposedSetup TestCatalogCompatibility TestBaselineCompatibilityCorpus; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && grep -q '"type": "composed"' internal/baseline/assets/setups/go-cli-typescript-bun.json && git diff --quiet HEAD -- internal/baseline/testdata/parity-corpus && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the six new tests and the composed snapshot do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Measured facts; Interfaces; Data Models; The composed setup procedure; Testing Approach 2; Testing Approach 5; Build Order 2
- `_prd.md` → Goal 2; Story 2; Core Feature 1; Success Metric 1; Success Metric 2; Success Metric 7
- `_techspec.md` → API Contract 4; API Contract 5
- ADR-0072, ADR-0081, ADR-0149, ADR-0191, ADR-0192, ADR-0204
