---
task: task_01
spec: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
status: completed
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


## Result

Implemented task_01's catalog slice. The three profiles retain their identifiers
and select `go`, `rust` and `typescript`; the four retired snapshots are gone.
`core` version 15 drops `review`, and `typescript` version 6 drops `triage`,
including their dispatch entries. All other module fields and clauses are
unchanged. The asset sync refreshed the seeded snapshots from a disposable
local clone of `~/dev/skills` at
`b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`; the source checkout was not written
and no network clone was used.

The prescribed existing-test edits and four new tests are in place.
`removedSkillFindings` sweeps the seventeen removed names;
`setupNameFindings` checks only the three named profiles and four retired
names, source-path consistency and the Roundfix-owned entry. Negative cases
assert exact findings for all seventeen names in all five fields, the retired
setups, a wrong profile setup, a wrong source path, a missing owned entry and
an externally sourced `roundfix`. A later composed setup is accepted.

Focused evidence (all Go commands used
`GOCACHE=/tmp/roundfix-0208-task01-gocache`):

| Acceptance criterion | Implementation and current-turn evidence |
| --- | --- |
| Three renamed setups and profile selection | `go test ./internal/baseline -count=1 -run 'Test(NoCatalogEntryNamesASkillRemovedUpstream\|ARemovedSkillNameInTheCatalogIsReported\|EveryBuiltInProfileTakesItsRenamedUpstreamSetup\|ARetiredOrUnownedSetupIsReported)$'` exited 0. Before the edits, the two positive tests failed on the old profiles, retired setups and removed names. |
| Removed names absent and planted names reported | The same focused run passed both removed-skill tests, including 85 planted-name cases. No replacement review skill was added. |
| Retired names, wrong source paths and missing ownership reported | The same focused run passed both setup tests and their negative cases. A separate JSON comparison against each predecessor at HEAD confirmed every owned `minimumVersion` was preserved; all three retain the Roundfix-owned `roundfix` entry. |
| Asset sync, parity and frozen enumeration with three setups | Asset sync wrote three snapshots (43, 36 and 101 skills); its rerun with `--check --format text` exited 0 and printed `setup-context-driven audit: ok`. `make baseline-digests` exited 0, including its strict catalog check and sanctioned parity regeneration. The TechSpec's exact jq transform changed only `asset-sync.json`; a byte comparison confirmed all other parity fixtures unchanged and its `plannedByteSequence` and `fileIdentities` preserved. The stable-tree `make verify-incremental` rerun exited 0, including the Baseline and CLI suites; details are recorded below. |
| Repository dispatch and refresh convergence | The required `baseline update --repo . --no-skills --yes --format text` ran twice, both exiting 0. The first changed only the dispatch guide and Setup Manifest; the second reported `File changes: 0` and verified idempotence. Inspection confirms the dispatch guide omits `trigger.core.review`. |

The sync initially encountered sandbox denial while creating its Git-private
transaction directory. The authorized rerun with the required access exited 0;
the first managed refresh used the same access. The refresh reports the existing
nested-carrier warnings for the formatter golden and Source Baseline corpus;
those carriers were preserved.

Scope postflight compared tracked and untracked changes with this Task's original Context:
all 32 changed or new paths are allowed. Source Baselines, retention transitions,
`.agents/skills/`, `skills-lock.json` and the frozen lock compatibility fixture
were not changed. No Task Graph or other Task was edited, and no commit, push
or Pull Request was created. Task status and the authored Verification command
remain Daemon-owned; the declared Verification command was not run.

Repository-required incremental check: the first
`rtk proxy make verify-incremental` exited 2. Baseline, CLI, Daemon and Spec-check
tests reported `PASS`, but their repository guards correctly rejected my
concurrent edit of this Result section. The diagnostic named only this Task
file. The stable-tree rerun exited 0: formatting, vet, all package tests (including
Baseline sync, parity, frozen enumeration and CLI tests), skill checks and the
build passed. This Result evidence was updated only after the rerun exited. Logs:
`/tmp/roundfix-0208-task01-incremental.log` and
`/tmp/roundfix-0208-task01-incremental-rerun.log`.


### Verification Feedback repair

The settlement consistency diagnostic identified four Context interfaces that
still pointed to snapshots this Task intentionally deleted. Removed those stale
interface entries; the three replacement snapshot declarations remain, and the
Requirements still document the four deletions. No catalog, test, Task status,
Task Graph or other Task changed in this repair.

Focused repair check:
`GOCACHE=/tmp/roundfix-0208-task01-gocache go run -buildvcs=false ./cmd/roundfix spec check 0208-a-baseline-that-follows-the-reshaped-skills-catalog --format json`
exited 0 with no findings or repair inputs and `verification.executed: false`.
The JSON retains skips for optional absent artifacts; no unresolved Context
finding remains. `git diff --check` also exited 0. This repair changed only this
Task file's Context and Result; the retired files remain absent. The Daemon's
authored Verification was not rerun by the Agent.
