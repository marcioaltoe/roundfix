---
task: task_03
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
status: completed
type: backend
complexity: high
---

# Task 03: The Roundfix skill is in every setup and has a dispatch trigger

## Overview

The `typescript-bun` setup lists the Roundfix skill; the `go-cli` and `rust-cli` setups do not, and no module tells an Agent when to load it. This Task adds the entry to the two setups, makes the `autonomous-work` module require the skill and carry one trigger for it, and teaches the asset sync to keep an owned entry the upstream list omits. Without the sync rule the next refresh from upstream would drop the entry, because the upstream lists for these two setups do not name the skill.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST append the setup entry the TechSpec's "Fixed texts" gives as the last entry of `skills` in `internal/baseline/assets/setups/go-cli.json` and `rust-cli.json`. No other entry, no `minimumVersion` and no other field changes by hand; the setup `digest` is regenerated.
2. MUST append the same entry, in the fixture's own key order, to the `go-cli` and `rust-cli` setups of `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`. The regeneration recomputes that file's digests and `v1/manifest.json`.
3. MUST set `requiredSkills` and add the `skillDispatch` entry in `internal/baseline/assets/modules/autonomous-work.json` exactly as the TechSpec gives them, and raise the module version by one. Replace strings and numbers in place; do not pass the file through a JSON encoder. No rule, clause or guide changes.
4. MUST make `buildAssetsSyncSnapshot` in `internal/baseline/assets_sync.go` keep each entry of the current snapshot whose source type is `repo` and whose name the upstream list did not yield, after the upstream entries and in the snapshot's recorded order, with its recorded path and minimum version. An entry of any other source type that the upstream list omits is still dropped.
5. MUST run `go test -buildvcs=false ./skills -run '^TestAuthorialSkillSync$' -update -count=1`, then `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a catalog snapshot, the parity manifest or a generated guide.
6. MUST add `internal/baseline/roundfix_skill_membership_test.go` and `internal/baseline/assets_sync_owned_membership_test.go` with the five tests the TechSpec's Testing Approach 3 names. The profile test MUST read every profile of the embedded catalog and find the trigger among the dispatch entries of its selected modules. The two sync tests build their upstream source in a temporary directory.
7. MUST NOT add, remove or rename any other skill in any setup, change `internal/baseline/assets/setups/typescript-bun.json`, or edit any existing test.

## Subtasks

- [ ] Add the entry to the two setups and the parity fixture.
- [ ] Require and dispatch the skill in `autonomous-work`.
- [ ] Keep recorded owned entries in the asset sync.
- [ ] Regenerate, refresh this repository's guides, and confirm the second refresh is a no-op.
- [ ] Add the membership, dispatch and sync tests, each negative case separate.

## Acceptance Criteria

- [ ] All three setups list the Roundfix skill, and a setup without it is reported by the check.
- [ ] Every profile's selected modules dispatch the skill, and `docs/agents/skill-dispatch.md` carries `trigger.autonomous-work.roundfix`.
- [ ] A sync from an upstream list that omits the Roundfix skill keeps its entry and drops an external skill the list omits.
- [ ] The existing asset-sync compatibility test passes against the fixture.
- [ ] A second managed refresh is a no-op.

## Context

- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `internal/baseline/assets/setups/typescript-bun.json`
- instruction: `internal/baseline/assets_sync_test.go`
- interface: `internal/baseline/assets/setups/go-cli.json`
- interface: `internal/baseline/assets/setups/rust-cli.json`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/skill-dispatch.md`
- creates: `internal/baseline/roundfix_skill_membership_test.go`
- creates: `internal/baseline/assets_sync_owned_membership_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestEverySetupListsTheRoundfixSkill|TestASetupWithoutTheRoundfixSkillIsReported|TestEveryProfileDispatchesTheRoundfixSkill|TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits|TestAssetSyncStillDropsAnExternalSkillTheUpstreamListOmits|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEverySetupListsTheRoundfixSkill TestASetupWithoutTheRoundfixSkillIsReported TestEveryProfileDispatchesTheRoundfixSkill TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits TestAssetSyncStillDropsAnExternalSkillTheUpstreamListOmits TestAssetsSyncCompatibilityMatchesMaintainedPythonContract TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization TestBaselineCompatibilityCorpus; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/skill-dispatch.md|trigger.autonomous-work.roundfix" "internal/baseline/assets/setups/go-cli.json|skills/06-review-repair/roundfix" "internal/baseline/assets/setups/rust-cli.json|skills/06-review-repair/roundfix"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the five new named tests exists and neither setup lists the skill, so the command fails.

## References

- [_techspec.md](_techspec.md) — System Architecture; Interfaces; Fixed texts; Testing Approach 3
- `_prd.md` → Goal 3; Core Feature 3; Success Metrics 4 and 6
- `_techspec.md` → API Contracts 4 and 5; Testing Approach 5
- [references/2026-09-30-an-owned-skill-older-than-the-bundle-passes-readiness.md](references/2026-09-30-an-owned-skill-older-than-the-bundle-passes-readiness.md)
- ADR-0081, ADR-0149, ADR-0189

## Result

Implemented this Task's slice for Daemon Verification. Task status and the
declared Verification command remain Daemon-owned; no repository commit,
push, Pull Request, Task Graph edit or other Task edit was made.

The Go and Rust setups now append the exact Roundfix entry, with the same
entry appended in the parity fixture's key order. `autonomous-work` moves
from module version 11 to 12, requires `roundfix` and carries the fixed
dispatch trigger. Asset sync appends omitted recorded `repo` entries after
upstream entries in their recorded order, preserving the complete entries;
omitted entries with other source types still disappear.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| All three setups list Roundfix; a missing entry is reported | `TestEverySetupListsTheRoundfixSkill` and `TestASetupWithoutTheRoundfixSkillIsReported` passed. Before the source edits, the membership test reported `go-cli` and `rust-cli`. A structural comparison against HEAD confirmed all previous entries and minimum versions remain unchanged and `typescript-bun.json` is unchanged. |
| Every profile dispatches Roundfix; this repository's guide carries the trigger | `TestEveryProfileDispatchesTheRoundfixSkill` passed for `go-cli-tui`, `rust-cli` and `standard-typescript-monorepo`, inspecting every embedded profile's selected module dispatch entries. Before the source edits all three lacked the trigger. The public managed refresh generated `trigger.autonomous-work.roundfix` in `docs/agents/skill-dispatch.md` (line 88). |
| Sync keeps an omitted owned skill and drops an omitted external skill | `TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits` passed against a committed temporary upstream list, asserting upstream-first order and unchanged Roundfix and custom-owned entries. Before the sync edit it returned only the upstream entry. `TestAssetSyncStillDropsAnExternalSkillTheUpstreamListOmits` passed independently for omitted `github` and `local` entries. Both tests create their upstream source in temporary directories. |
| Asset-sync compatibility remains intact | `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` passed against the regenerated fixture. |
| Second managed refresh is a no-op | The second confirmed public refresh exited 0, reported `File changes: 0`, `approved Baseline Plan is already applied` and `Idempotence: verified`. |

### Commands and outcomes

All Go commands used `GOCACHE=/tmp/roundfix-task03-gocache`.

- `go test -buildvcs=false ./skills -run '^TestAuthorialSkillSync$' -update -count=1`: exit 0.
- `make baseline-digests`: exit 0; regenerated setup digests, formatter dispatch golden, Standard TypeScript profile pin, catalog snapshots, plan goldens, parity fixture digests and parity manifest through the sanctioned workflow. No derived value or guide was hand-edited.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`: first sandbox attempt refused Git-private transaction-directory creation. Retried with filesystem permission: exit 0, two managed files updated. Repeated after apply: exit 0, zero file changes.
- `go test -buildvcs=false ./internal/baseline -run '^(TestEverySetupListsTheRoundfixSkill|TestASetupWithoutTheRoundfixSkillIsReported|TestEveryProfileDispatchesTheRoundfixSkill|TestAssetSyncKeepsAnOwnedSkillTheUpstreamListOmits|TestAssetSyncStillDropsAnExternalSkillTheUpstreamListOmits|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract)$' -count=1 -v`: exit 0; all six tests passed.
- `rtk make verify-incremental`: initial attempt failed because the concurrent managed refresh changed repository bytes during suite-guard checks and the sandbox denied process-table access to process-stop tests. Reran with required permission after all source and guide writes stopped: exit 0; formatting, vet, package tests, skill synchronization/checks and CLI build passed.
- Structural scope checks and `git -c core.fsmonitor=false diff --check`: passed. Module rules, clauses and guide versions remain unchanged. Changed-file postflight found only this Task's Context paths, its two new tests and its own Task file.

The complete declared `## Verification` command was not run. Daemon
Verification and settlement remain pending; there are no implementation
follow-ups outside this slice.

## Carry-forward provenance

- Source Run: `run_20260930T162501Z_8ec6b68021df9cef`
- Source commit: `e30d995ab648a0e661047ab9295ccd805621a672`
