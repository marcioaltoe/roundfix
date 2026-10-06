---
task: task_01
spec: 0241-retire-the-fool-autoresearch-and-council
status: pending
type: backend
complexity: high
---

# Task 01: The Baseline retires council and the-fool and the asset sync keeps them out

## Overview

The `context-workflow` module requires and dispatches `council` and
`the-fool`, and all four Setup Snapshots list them because the upstream setup
lists still do. This Task adds the Retired Skill set, makes the asset sync
drop a Retired Skill from every snapshot it writes, removes both skills from
the module, refreshes the snapshots from upstream `b3c45a4`, and regenerates
the parity fixture, the derived pins and this repository's managed guides.

This is an authorized tooling Task. Its Governed Paths are the module, the
four snapshots, the derived profile and formatter golden, and the two managed
guides, all bounded in `_authorization.md`.

## Requirements

1. MUST create `internal/baseline/retired_skills.go` with `retiredSkills` and
   `RetiredSkills()` as the TechSpec's Interfaces describe: exactly
   `council` and `the-fool`.
2. MUST make `buildAssetsSyncSnapshot` in `internal/baseline/assets_sync.go`
   follow Invariant 1: skip a Retired Skill from the upstream list before the
   duplicate checks, and skip a recorded `repo` entry named as a Retired Skill
   when it re-appends omitted repository-owned entries. MUST remove `council`
   from `assetsSyncRepoOwnedSkills`. No finding is added and no other entry's
   handling changes.
3. MUST apply the task_01 Fixed text to
   `internal/baseline/assets/modules/context-workflow.json` (version 22, the
   two skills and their two dispatch objects removed) and change no other
   line, so a concurrent Spec's edits to the module's clauses merge cleanly.
4. MUST run the TechSpec's refresh procedure against
   `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`, cloning only the local
   checkout `~/dev/skills`, and MUST stop and report, without cloning from
   the network, when that checkout lacks the commit. MUST NOT write
   `~/dev/skills`. The `--check` rerun MUST exit `0`, and every snapshot MUST
   keep its other entries, its `roundfix` entry and its `minimumVersion`
   values.
5. MUST apply the TechSpec's fixture transform to
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
   and leave every other parity fixture byte-identical (ADR-0072).
6. MUST create `internal/baseline/retired_skills_test.go` with
   `TestNoCatalogEntryNamesARetiredSkill` and
   `TestARetiredSkillNameInTheCatalogIsReported` (Testing Approach 1), built
   on the shared `namedSkillFindings` that
   `internal/baseline/upstream_removed_skills_test.go` gains as "Existing
   tests that change" describes; the removed-skill findings keep their text.
7. MUST create `internal/baseline/assets_sync_retired_test.go` with
   `TestAssetSyncDropsARetiredSkillTheUpstreamListNames` (Testing Approach
   2), including a recorded `repo` `council` entry and a non-retired control
   that the sync keeps.
8. MUST run `make baseline-digests`, then the managed refresh twice as the
   procedure says; the second MUST report `File changes: 0`. MUST NOT
   hand-edit a pin, a golden, a snapshot's skills or a digest.
9. MUST NOT change a clause, a Source Baseline file, a retention transition,
   any `.agents/skills/` or `skills/` file, `skills-lock.json` or
   `CONTEXT.md`.

## Subtasks

- [ ] Add the Retired Skill set and the sync filter.
- [ ] Edit the module and refresh the snapshots from the local upstream clone.
- [ ] Transform the parity fixture and regenerate the derived pins.
- [ ] Add the catalog and sync tests.
- [ ] Refresh this repository's managed guides twice.

## Acceptance Criteria

- [ ] No module, trigger, bundle or snapshot names `council` or `the-fool`,
      and a planted name in any of them is reported.
- [ ] A sync over an upstream list naming both, with a recorded
      repository-owned `council`, writes neither; a non-retired skill is kept.
- [ ] The sync `--check` against `b3c45a4` exits `0`.
- [ ] `docs/agents/skill-dispatch.md` names neither trigger, and a second
      managed refresh changes nothing.

## Context

- instruction: `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- instruction: `internal/baseline/assets_sync_owned_membership_test.go`
- creates: `internal/baseline/retired_skills.go`
- creates: `internal/baseline/retired_skills_test.go`
- creates: `internal/baseline/assets_sync_retired_test.go`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/upstream_removed_skills_test.go`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/setups/go.json`
- interface: `internal/baseline/assets/setups/rust.json`
- interface: `internal/baseline/assets/setups/typescript.json`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/skill-dispatch.md`
- interface: `docs/agents/setup-context.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestNoCatalogEntryNamesARetiredSkill|TestARetiredSkillNameInTheCatalogIsReported|TestAssetSyncDropsARetiredSkillTheUpstreamListNames|TestNoCatalogEntryNamesASkillRemovedUpstream|TestARemovedSkillNameInTheCatalogIsReported|TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits|TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract|TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestCatalogDiagnosticCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNoCatalogEntryNamesARetiredSkill TestARetiredSkillNameInTheCatalogIsReported TestAssetSyncDropsARetiredSkillTheUpstreamListNames TestNoCatalogEntryNamesASkillRemovedUpstream TestARemovedSkillNameInTheCatalogIsReported TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent TestAssetsSyncCompatibilityMatchesMaintainedPythonContract TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus TestCatalogDiagnosticCharacterization; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && ! grep -Eq 'trigger.context-workflow.(council|the-fool)' docs/agents/skill-dispatch.md && ! grep -Eq '"(council|the-fool)"' internal/baseline/assets/modules/context-workflow.json internal/baseline/assets/setups/go.json internal/baseline/assets/setups/rust.json internal/baseline/assets/setups/typescript.json internal/baseline/assets/setups/go-cli-typescript-bun.json && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the three new tests do not exist and the module, snapshots and dispatch guide name both skills, so the command fails.

## References

- `_prd.md` → Goal 1; User Stories 1-2; Core Features 1-2; Success Metrics 1-2; Declared breaks
- `_techspec.md` → Measured facts; Interfaces; Fixed texts (task_01); The refresh procedure; Existing tests that change; API Contract 1; Testing Approach 1-2; Build Order 1
- ADR-0246, ADR-0191, ADR-0072, ADR-0081, ADR-0149
