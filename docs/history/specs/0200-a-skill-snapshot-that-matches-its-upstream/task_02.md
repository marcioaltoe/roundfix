---
task: task_02
spec: 0200-a-skill-snapshot-that-matches-its-upstream
status: completed
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

The four setup snapshots now follow upstream commit
`a4e18e4fa223196b51d0fd8224e5a33b84f97717`. The catalog follows all three
renames, the Go CLI/TUI profile selects `go-tui`, and the Go module requires
its three previously unrequired dispatched skills. The same-session
Verification Feedback authorized repairing the Doctor count assertion; it now
derives the external count from the asserted skill list.
No Task status, Task Graph, other Task file, commit, push or Pull Request was
changed by this Agent.

### Implementation and acceptance evidence

| Acceptance criterion | Implementation and current evidence |
| --- | --- |
| No module, bundle or setup uses an upstream-renamed skill; invalid names are reported | Added the five authored tests in `internal/baseline/upstream_skill_names_test.go`. The detector reads every embedded module's requirements and dispatch, every activation bundle, and every setup; renamed dispatch triggers must use `trigger.<module>.<new name>`. Literal negative documents cover each old name in requirements, dispatch, bundles and setups, plus stale triggers. Focused tests passed. |
| Every snapshot pins the upstream commit; Go CLI/TUI selects the TUI setup | Seeded `go-tui` by byte-copying `go-cli` and replacing only its id, then ran the public asset sync. It reported four updated snapshots at exit 0. The profile test passed and confirms `bubbletea` and `tui-design`. Read-only JSON checks confirmed every GitHub entry and snapshot source uses the full pinned commit. Membership counts are 47, 49, 36 and 110 for `go-cli`, `go-tui`, `rust-cli` and `typescript-bun`, including the retained Roundfix entry. |
| Every module requires every dispatched skill; missing requirements are reported | Added the three Go requirements in alphabetical position. The sweep passed across all modules; the negative test uses literal documents with both the `skill` and fallback `id` forms and asserts the exact missing-requirement finding. |
| Asset-sync, ownership, corpus and Doctor checks | Focused asset-sync, output-ownership and synthetic parity checks passed. Sanctioned regeneration and the incremental check's full Baseline suite passed. Doctor's required skill-set comparison and output assertion now agree on 12 external skills. The focused Go CLI/TUI Doctor subtest passed after the Verification Feedback repair. The earlier incremental check failed on the stale assertion; it was not rerun after repair. Declared Verification remains Daemon-owned and was not executed by this Agent. |
| Second managed refresh is a no-op | First authorized update reported two changes, solely `docs/agents/skill-dispatch.md` and `docs/agents/setup-context.json`. The second reported `File changes: 0`, `Baseline update: verified`, and `Idempotence: verified`, exit 0. |

The four module versions increased once from the Task base: core 12→13,
go 2→3, rust 2→3 and typescript 4→5. A structural comparison against HEAD
confirmed that every other module field outside version, requirements and
dispatch is unchanged, preserving clauses, rules and guides. The profile's
only source change is its setup value. A separate comparison confirmed every
setup's recorded owned minimum versions are unchanged and every setup retains
one repository-owned `roundfix` entry.

The parity fixture was rebuilt with the TechSpec's exact jq transform. A
read-only comparison confirmed all fields outside the three specified
projections stayed identical, including `plannedByteSequence` and
`fileIdentities`. Digests, the parity manifest, pins, catalog snapshots,
formatter golden and plan goldens were written only by `make baseline-digests`.
Rendered repository guidance was written only by the public managed refresh.

### Commands and outcomes

CLI and subsequent checks used
`GOCACHE=/tmp/roundfix-0200-task02-gocache` after the default shared cache
refused access. The sync and managed refresh ran with sandbox escalation only
because their recoverable journals live in this Run worktree's Git-private
directory outside the writable workspace. No automatic approval rejection
occurred.

- Cloned `/Users/marcio/dev/skills` with `git clone --quiet --no-local` into
  `/tmp/roundfix-0200-task02-skills`, checked out the full commit detached,
  and set origin to `https://github.com/marcioaltoe/skills.git`. The final
  `rev-parse HEAD` returned the pinned commit and porcelain status was empty.
  The source checkout at `~/dev/skills` was never written.
- `go test -count=1 ./internal/baseline -run 'Test(NoCatalogEntry|ARenamedSkillName|TheGoCLITUIProfile|EveryModuleRequires|ADispatchedSkillNoModule)'`:
  before asset edits, exposed the old names, old profile setup and three Go
  requirements; literal negative cases passed.
- `go run -buildvcs=false ./cmd/roundfix baseline assets sync --source-dir /tmp/roundfix-0200-task02-skills/setups --format text`:
  exit 0, four snapshot updates, no errors.
- TechSpec jq fixture transform, followed by `make baseline-digests`:
  exit 0; regeneration reported `ok:true, changed:true` and rewrote its eleven
  derived artifacts. No generated digest or rendered guide was hand-edited.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`, twice:
  exit 0 both times; first two file changes, second zero. Existing nested
  carrier warnings were preserved and skills were skipped.
- `go test -count=1 ./internal/baseline -run 'Test(NoCatalogEntry|ARenamedSkillName|TheGoCLITUIProfile|EveryModuleRequires|ADispatchedSkillNoModule|AssetsSync|BaselineAssetsSync|OutputsForCommand|TheParityFixture|ASyntheticDigest)'`:
  exit 0. This sequential run supersedes a check that overlapped regeneration
  and tripped the suite guard on a concurrently regenerated plan golden.
- `go test -count=1 ./internal/cli -run 'TestRunDoctorDerivesExternalSkillRequirement'`:
  exit 1; expected text says `26 required: 14 Roundfix-owned, 9 external`,
  actual text correctly says `26 required: 14 Roundfix-owned, 12 external`.
- `make verify-incremental`: exit 2 at its test target, solely the same Doctor
  output assertion. Formatter and vet checks passed; every other package,
  including the full Baseline suite and owned-skills tests, passed. Later
  prerequisites after the failed test target were not reached.
- Read-only Python structural/preservation checks and `git diff --check`:
  exit 0. Changed-path inspection found only this Task's declared paths and
  sanctioned generated outputs, plus the Task's pre-existing Daemon status
  change. Neither `_tasks.md` nor another Task file changed.

### Verification Feedback repair

Inspected the Daemon's attempt-1 diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261001T013631Z_c59b100cbb9f7af1/verification/batch-002-attempt-1.log`
and the related Doctor test. The diagnostic confirms the only reported failure
is the separate output assertion's hardcoded external count. The three Go
requirements intentionally increase the external skill list from 9 to 12;
production Doctor correctly reports that list's count.

The same-session request to repair this Task's slice resolves the earlier
scope clarification. Replaced the assertion's `9 external` with
`%d external` and added `len(goTUISkills)` as its third format argument. The
existing exact skill-list assertion remains intact. No production behavior,
Verification command, status or additional path changed.

- `GOCACHE=/tmp/roundfix-0200-task02-gocache go test -count=1 ./internal/cli -run 'TestRunDoctorDerivesExternalSkillRequirement/Go_CLI_and_TUI_modules_retain_their_skills'`:
  exit 0 after the repair.
- `gofmt -w internal/cli/doctor_test.go` and `git diff --check`: exit 0.
- Changed-path postflight: all 27 changed paths remain in the Task's declared
  Context, including sanctioned generated outputs. The Task's authored
  projection outside Daemon status and Result remains identical to HEAD.

The Daemon owns the next full configured Verification sequence. This Agent
has not rerun either declared Verification command or claimed a terminal
Task verdict.

## Carry-forward provenance

- Source Run: `run_20261001T013631Z_c59b100cbb9f7af1`
- Source commit: `1dca31a74e7c5511bd06be5f145299b17f9083cd`
