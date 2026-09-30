---
task: task_02
spec: 0200-a-skill-snapshot-that-matches-its-upstream
status: pending
type: backend
complexity: high
---

# Task 02: The setup snapshots follow upstream by name

## Overview

The four setup snapshots move from `14fdf46` to `a4e18e4`, the catalog follows the three upstream renames, the Go CLI/TUI profile takes the new `go-tui` setup, and the `go` module requires the three skills it already dispatches. The asset sync writes the snapshots; the sanctioned regeneration and the managed refresh write everything derived.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST make the module edits the TechSpec's "Fixed texts" give for task_02 in `core.json`, `rust.json`, `typescript.json` and `go.json`, raising each module's version by one from its value on this Task's base. Replace strings and numbers in place; do not pass a module through a JSON encoder. No clause, rule or guide changes.
2. MUST seed `internal/baseline/assets/setups/go-tui.json` as the TechSpec's Data Models describe, and change only the `setup` value of `internal/baseline/assets/profiles/go-cli-tui.json` to `go-tui`.
3. MUST run the TechSpec's refresh procedure: a clean detached clone of `marcioaltoe/skills` at `a4e18e4fa223196b51d0fd8224e5a33b84f97717` with the GitHub origin, then `roundfix baseline assets sync --source-dir <clone>/setups`. MUST NOT hand-edit a snapshot after the sync, and MUST NOT write `~/dev/skills`.
4. MUST rebuild the parity fixture's `manifest.setups`, `managedEntryLedger` and `normalizedOutput.result` with the TechSpec's fixture transform, and leave its other fields as they are.
5. MUST update exactly these existing test literals: the three setup counts `3` in `internal/baseline/assets_sync_test.go` (`Errors`, `Info`, `SetupIDs()`) and the `Info: 3` of the parity comparison become `4`; `TestOutputsForCommand` appends `internal/baseline/assets/setups/go-tui.json` to the frozen enumeration it reads, in sorted position, with a comment naming ADR-0191; `goTUISkills` in `internal/cli/doctor_test.go` gains the three Go skills in alphabetical position.
6. MUST run `make baseline-digests`, then `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`, and a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a catalog snapshot, a fixture digest, the parity manifest or a rendered guide.
7. MUST add `internal/baseline/upstream_skill_names_test.go` with the five tests of the TechSpec's Testing Approach 2. The rename table holds the three pairs; the rename check reads every module's `requiredSkills` and `skillDispatch`, every activation bundle and every setup snapshot, and requires each renamed trigger to be `trigger.<module>.<new name>`. The negative tests run the same functions on literal documents.
8. MUST keep every setup's recorded `minimumVersion` values and the Roundfix-owned entry Spec 0195 keeps.

## Subtasks

- [ ] Rename the skills in the modules, require the Go skills, seed `go-tui` and point the profile at it.
- [ ] Run the asset sync from a clean clone at `a4e18e4`.
- [ ] Rebuild the parity fixture rows and update the four named test literals.
- [ ] Regenerate, refresh this repository's guides, and confirm the second refresh is a no-op.
- [ ] Add the rename, profile-setup and dispatch-requirement tests.

## Acceptance Criteria

- [ ] No module, bundle or setup names `context7`, `feature-systems-pattern` or `rust`, and a catalog that does is reported.
- [ ] Every snapshot pins `a4e18e4`, and `go-cli-tui` takes `go-tui`, which lists `bubbletea` and `tui-design`.
- [ ] Every module requires every skill it dispatches, and a module that does not is reported.
- [ ] The asset-sync, ownership, corpus and doctor tests named in Verification pass.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- interface: `internal/baseline/assets/setups/go-cli.json`
- interface: `internal/baseline/assets/setups/rust-cli.json`
- interface: `internal/baseline/assets/setups/typescript-bun.json`
- creates: `internal/baseline/assets/setups/go-tui.json`
- interface: `internal/baseline/assets/profiles/go-cli-tui.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/go.json`
- interface: `internal/baseline/assets/modules/rust.json`
- interface: `internal/baseline/assets/modules/typescript.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/assets_sync_test.go`
- interface: `internal/baseline/derived_ownership_test.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/skill-dispatch.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/upstream_skill_names_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestNoCatalogEntryNamesASkillRenamedUpstream|TestARenamedSkillNameInTheCatalogIsReported|TestTheGoCLITUIProfileTakesTheGoTUISetup|TestEveryModuleRequiresEverySkillItDispatches|TestADispatchedSkillNoModuleRequiresIsReported|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract|TestAssetsSyncCheckIsReadOnlyAndReportsDrift|TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent|TestOutputsForCommand|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNoCatalogEntryNamesASkillRenamedUpstream TestARenamedSkillNameInTheCatalogIsReported TestTheGoCLITUIProfileTakesTheGoTUISetup TestEveryModuleRequiresEverySkillItDispatches TestADispatchedSkillNoModuleRequiresIsReported TestAssetsSyncCompatibilityMatchesMaintainedPythonContract TestAssetsSyncCheckIsReadOnlyAndReportsDrift TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent TestOutputsForCommand TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestAuthorialSkillSync)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestAuthorialSkillSync" || { printf 'missing pass: %s\n' TestAuthorialSkillSync >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestRunDoctorDerivesExternalSkillRequirementFromSetupManifest)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestRunDoctorDerivesExternalSkillRequirementFromSetupManifest" || { printf 'missing pass: %s\n' TestRunDoctorDerivesExternalSkillRequirementFromSetupManifest >&2; exit 1; }` — expected: exit 0; before this Task the five new tests do not exist and the asset-sync tests count three setups, so the command fails; the doctor test fails until `go.json` requires the three Go skills and its literal names them.
- `for pair in "docs/agents/skill-dispatch.md|trigger.core.context7-cli" "internal/baseline/assets/profiles/go-cli-tui.json|\"setup\": \"go-tui\"" "internal/baseline/assets/setups/go-tui.json|skills/05-implementation-loop/bubbletea" "internal/baseline/assets/setups/go-tui.json|a4e18e4fa223196b51d0fd8224e5a33b84f97717" "internal/baseline/assets/setups/go-cli.json|a4e18e4fa223196b51d0fd8224e5a33b84f97717" "internal/baseline/assets/setups/rust-cli.json|skills/05-implementation-loop/rust-expert" "internal/baseline/assets/setups/typescript-bun.json|skills/03-engineering-design/app-renderer-systems"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task `go-tui.json` does not exist and the guide names `context7`, so the command fails.

## References

- `_prd.md` → Goals 1-2; User Story 2; Core Feature 2; Success Metrics 2, 3 and 6; Declared breaks; Prerequisites
- `_techspec.md` → Measured facts; Data Models; Fixed texts; The refresh procedure; API Contract 3; Testing Approach 2 and 5; Build Order 2
- ADR-0072, ADR-0081, ADR-0103, ADR-0149, ADR-0191

## Result
