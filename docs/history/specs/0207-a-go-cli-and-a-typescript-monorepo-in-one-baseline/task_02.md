---
task: task_02
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: completed
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

## Result

Implemented the composed Setup Snapshot slice for Daemon Verification. Task
status remains Daemon-owned; the authored Verification command was not run,
and no commit, push or Pull Request was made.

### Implementation

- Composition preserves component order, drops equal entries sharing a skill
  name or activation bundle identifier, and refuses unequal entries.
- Catalog validation reports composition invalidity, conflict and drift
  separately, retaining the existing skill checks and digest rule.
- Asset sync builds upstream snapshots first, composes from those results,
  validates the resulting catalog, then writes through its existing atomic
  transaction. Plans are sorted by file path: the new composed filename sorts
  before `go.json`, even though its identifier sorts after `go`.
- Both sync and skill-regeneration source types preserve `source.setups`.
  The catalog digest fixture helper keeps shared component entries equal and
  recomposes after changing a component; its owned/external digest assertions
  are unchanged.
- Added the six named composition tests, using private temporary asset roots
  and Git sources for sync. Existing sync counts include four snapshots.
  The parity comparison names the ADR-0072 designed delta.
- Seeded the new snapshot, then generated its skills and digest only through
  asset sync from a temporary local clone at the Go snapshot's pinned commit,
  `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`. Sync reported exactly one update:
  `go-cli-typescript-bun`. The snapshot has 112 skills and 10 bundles. Added its
  sanctioned output by changing only the enumeration line in
  `derived_ownership_test.go`.

### Acceptance evidence

| Criterion | Implementation and focused-check evidence |
| --- | --- |
| Ordered union and loadable catalog | `TestTheComposedSetupIsTheUnionOfItsComponents` checks the embedded union independently, including full entries and their order, plus explicit equal-entry deduplication. It passed in the focused run; the Baseline package also passed incremental verification. |
| Distinct refusal codes | The drift, unknown/composed component and conflicting component tests passed. Cases cover skill and bundle drift, unknown and composed components, insufficient/duplicate/malformed component lists, and unequal shared skill and bundle entries. They assert `.drift`, `.invalid` and `.conflict` respectively. |
| Same-run refresh, unchanged bytes and check drift | Both new sync tests passed. A committed source skill change rewrites its component and composition in one run. Repeated unchanged sync leaves the entire asset tree byte-identical. Check mode reports composed drift without writing, including composed-only drift. A composition conflict returns `AssetsSyncInvalid` with the existing invalid-assets finding and preserves the complete asset preimage. |
| Frozen parity | `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` passed in the focused run. Byte comparison against `HEAD` confirms the parity corpus and `asset-sync.json` are unchanged; the fixture SHA-256 is `5e89fa3377d4a41bdc45a5a4a0e92de10857289ecf6ed4ef3e7e81f5b7d2e233`. |

### Commands and outcomes

All Go commands used `GOCACHE=/tmp/roundfix-task02-gocache`.

- `go test -count=1 ./internal/baseline -run 'ComposedSetup|ComponentsThat|AssetSync|CatalogDigestExcludesOwnedSkillContent'`
  — exit 0 after repairing the path-ordering defect exposed by the first
  focused run and correcting the conflict fixtures to mutate shared entries.
- `go run -buildvcs=false ./cmd/roundfix baseline assets sync --source-dir /var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task02-upstream-vmezhgaz/skills/setups --format text`
  — exit 0, one composed-snapshot update. The initial sandbox attempt could
  not create the Git-private transaction directory; the approved retry
  succeeded. No network source was used and no existing setup changed.
- `rtk make baseline-digests` — exit 0; generated the catalog digest,
  normalized catalog and four plan-characterization goldens. No pin, golden
  or composed skill list was hand-edited.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  — executed twice with Git transaction access, both exit 0. The first
  updated only `docs/agents/setup-context.json`; the second reported
  `File changes: 0` and verified idempotence.
- `rtk make verify-incremental` — exit 0 on the approved rerun. The initial
  sandbox run passed the Baseline package but failed two existing CLI
  process-owner tests with `read process table ... operation not permitted`.
  The rerun passed those tests, analyzer, package tests, skill checks and build.
- Changed-file inspection and byte comparisons against `HEAD` confirm no
  profile, module, decision, existing setup, parity corpus, Source Baseline
  corpus or transition changed. Every new changed path belongs to this Task's
  Context or its assigned Task file.
- `git -c core.fsmonitor=false diff --check` — exit 0.

No follow-up implementation was added. Daemon Verification and settlement
remain pending.

## Carry-forward provenance

- Source Run: `run_20261002T024113Z_8885f5993a89b1e3`
- Source commit: `eefcf28a88e47c369c76397d492e0b57457ffecc`
