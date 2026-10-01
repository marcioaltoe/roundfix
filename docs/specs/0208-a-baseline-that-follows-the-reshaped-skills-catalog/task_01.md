---
task: task_01
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
status: pending
type: backend
complexity: high
---

# Task 01: The setup snapshots take their upstream names, and the removed skills leave the catalog

## Overview

Upstream renamed `go-tui`, `rust-cli` and `typescript-bun` to `go`, `rust` and `typescript`, retired `go-cli`, and removed `review` and `triage`. This Task seeds the three new snapshots, points the three built-in profiles at them, drops `review` from `core` and `triage` from `typescript`, and refreshes the snapshots to `b3c45a4` with the asset sync in one run, because the sync refuses a catalog that still requires a skill the refreshed snapshots omit. It then updates the parity fixture, the existing tests and the derived files, and this repository's guides.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST seed `go.json`, `rust.json` and `typescript.json` as the TechSpec's Data Models describe, delete `go-cli.json`, `go-tui.json`, `rust-cli.json` and `typescript-bun.json`, and change each built-in profile's `setup` in place: `go-cli-tui` to `go`, `rust-cli` to `rust`, `standard-typescript-monorepo` to `typescript`. Profile identifiers MUST NOT change.
2. MUST apply the task_01 Fixed texts to `core.json` and `typescript.json`, raising each module version by one from the Task's base. `review` MUST NOT be renamed to `code-review` or any other skill.
3. MUST run the TechSpec's refresh procedure against `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`, cloning only the local checkout `~/dev/skills`. MUST stop and report, without cloning from the network, when that checkout lacks the commit. MUST NOT write `~/dev/skills`. A rerun of the sync with `--check` MUST exit `0`, and every snapshot MUST keep the Roundfix-owned `roundfix` entry and the `minimumVersion` values of its predecessor.
4. MUST apply the TechSpec's fixture transform to `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and MUST leave every other parity fixture byte-identical (ADR-0072; the corpus is frozen).
5. MUST make exactly the existing-test edits the TechSpec's "Existing tests that change" lists for task_01, and MUST create `internal/baseline/upstream_removed_skills_test.go` with the four tests of Testing Approach 1, built on `removedSkillFindings` and `setupNameFindings` as the TechSpec's Interfaces describe. The setup check MUST name only the three listed profiles and the four retired setup names, so a later composed snapshot passes it.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` twice; the second MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a snapshot's skills or a digest.
7. MUST NOT change a clause, a Source Baseline file, a retention transition, `lock-hash-compatibility-v1.json`, any `.agents/skills/` file or `skills-lock.json`.

## Subtasks

- [ ] Seed the three snapshots, delete the four retired ones, and point the profiles at the new names.
- [ ] Drop `review` and `triage` from their modules.
- [ ] Refresh from `b3c45a4` and transform the parity fixture.
- [ ] Update the existing tests and write the removed-skill and setup-name tests.
- [ ] Regenerate, then refresh this repository twice.

## Acceptance Criteria

- [ ] The embedded catalog holds the setups `go`, `rust` and `typescript` and no retired setup, and each built-in profile takes its renamed setup.
- [ ] No module, trigger, bundle or snapshot names any of the seventeen removed skills, and a planted name in each field is reported.
- [ ] A retired setup name, a wrong `source.path`, and a setup without the Roundfix-owned entry are each reported.
- [ ] The asset sync tests, the parity comparison and the frozen-enumeration test pass with three setups.
- [ ] This repository's dispatch guide has no `trigger.core.review`, and a second refresh changes nothing.

## Context

- instruction: `docs/adr/0206-the-baseline-takes-upstream-setup-names-and-drops-a-removed-skill.md`
- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- instruction: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/assets/setups/go-cli.json`
- interface: `internal/baseline/assets/setups/go-tui.json`
- interface: `internal/baseline/assets/setups/rust-cli.json`
- interface: `internal/baseline/assets/setups/typescript-bun.json`
- creates: `internal/baseline/assets/setups/go.json`
- creates: `internal/baseline/assets/setups/rust.json`
- creates: `internal/baseline/assets/setups/typescript.json`
- interface: `internal/baseline/assets/profiles/go-cli-tui.json`
- interface: `internal/baseline/assets/profiles/rust-cli.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/typescript.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/catalog_test.go`
- interface: `internal/baseline/assets_sync_test.go`
- interface: `internal/baseline/upstream_skill_names_test.go`
- interface: `internal/baseline/derived_ownership_test.go`
- interface: `internal/baseline/catalog_dispatch_outside_setup_test.go`
- interface: `internal/cli/baseline_assets_sync_test.go`
- creates: `internal/baseline/upstream_removed_skills_test.go`
- interface: `docs/agents/skill-dispatch.md`
- interface: `docs/agents/setup-context.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestNoCatalogEntryNamesASkillRemovedUpstream|TestARemovedSkillNameInTheCatalogIsReported|TestEveryBuiltInProfileTakesItsRenamedUpstreamSetup|TestARetiredOrUnownedSetupIsReported|TestTheGoCLITUIProfileTakesTheGoSetup|TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract|TestAssetsSyncCheckIsReadOnlyAndReportsDrift|TestAssetsSyncProvenanceAndPreMutationRefusals|TestOutputsForCommand|TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestCatalogDiagnosticCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNoCatalogEntryNamesASkillRemovedUpstream TestARemovedSkillNameInTheCatalogIsReported TestEveryBuiltInProfileTakesItsRenamedUpstreamSetup TestARetiredOrUnownedSetupIsReported TestTheGoCLITUIProfileTakesTheGoSetup TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent TestAssetsSyncCompatibilityMatchesMaintainedPythonContract TestAssetsSyncCheckIsReadOnlyAndReportsDrift TestAssetsSyncProvenanceAndPreMutationRefusals TestOutputsForCommand TestEveryBuiltInProfileSetupListsEverySkillItsGuidesName TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus TestCatalogDiagnosticCharacterization; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for retired in go-cli go-tui rust-cli typescript-bun; do test ! -e "internal/baseline/assets/setups/$retired.json" || { printf 'retired snapshot remains: %s\n' "$retired" >&2; exit 1; }; done && ! grep -q 'trigger.core.review' docs/agents/skill-dispatch.md && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the four new tests do not exist and `docs/agents/skill-dispatch.md` names `trigger.core.review`, so the command fails.

## References

- `_prd.md` → Goals 1-2; User Stories 1-2; Core Features 1-2; Success Metrics 1-2; Declared breaks
- `_techspec.md` → Measured facts; Data Models; Fixed texts (task_01); The refresh procedure; Existing tests that change; API Contracts 1, 2 and 4; Testing Approach 1 and 5; Build Order 1
- ADR-0206, ADR-0191, ADR-0072, ADR-0081, ADR-0149
