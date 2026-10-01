---
task: task_01
spec: 0200-a-skill-snapshot-that-matches-its-upstream
status: completed
type: backend
complexity: medium
---

# Task 01: The asset sync validates the catalog it produces

## Overview

`roundfix baseline assets sync` validates the current catalog before it reads the source, so a refresh that follows an upstream rename is refused whichever side is edited first. This Task makes the sync validate only the catalog it would write, and makes the parity fixture's synthetic digests a regenerated value, so task_02 can add rows without computing digests by hand. No asset changes.

## Requirements

1. MUST remove the validation of the current catalog that `syncAssets` performs before it inspects the source, and MUST validate the overlay catalog (current assets plus generated snapshot overrides) on every run, including a run with no overrides, as the TechSpec's Interfaces describe.
2. MUST keep both refusal messages, their codes, categories and actions: with no overrides an invalid catalog reports "Go-owned canonical Baseline assets are invalid"; with overrides it reports "Generated setup snapshots are incompatible with the Baseline catalog". A refused run writes nothing, in check mode and in refresh mode.
3. MUST add `internal/baseline/assets_sync_synthetic_test.go` with `assetsSyncSyntheticSkillFile`, `assetsSyncSyntheticTreeDigest` and `parityFixtureDigestFindings` as the TechSpec sketches them, make `buildAssetsSyncSource` write `assetsSyncSyntheticSkillFile`, and make `regenerateBaselineCompatibilitySetups` set every `github` skill's `treeDigest` from `assetsSyncSyntheticTreeDigest` before it computes the setup digest.
4. MUST add `internal/baseline/assets_sync_upstream_rename_test.go` with the four sync tests of the TechSpec's Testing Approach 1. The rename renames `handoff` to `handoff-next` in the target's `core` module (`requiredSkills` and `skillDispatch`) and in the source's skill directory and setup lists. The no-drift case first refreshes, then breaks the target module, then syncs again.
5. MUST add `TestTheParityFixtureDigestsFollowTheSyntheticSource` (every `github` `treeDigest` in the committed fixture equals the synthetic digest of its path) and `TestASyntheticDigestThatDiffersIsReported` (a literal fixture with one altered digest is reported).
6. MUST run `make baseline-digests` and confirm it reports `"changed":false`: the derivation reproduces the committed fixture.
7. MUST add, after "Without `--check`, Roundfix validates the generated catalog in memory" in the "Canonical asset synchronization" section of `docs/user-guide/context-driven-development.md`, the sentence: "Both modes validate the catalog as the refresh would leave it, so a module edit that names a skill only the refreshed snapshots list lands together with the refresh."
8. MUST NOT change any asset, any existing assertion, or the fixture's bytes.

## Subtasks

- [ ] Move the sync's catalog validation to the produced catalog.
- [ ] Share the synthetic skill content and derive the fixture digests from it.
- [ ] Add the rename, refusal and digest tests.
- [ ] Document the validation in the user guide.

## Acceptance Criteria

- [ ] A sync whose target names a skill only the refreshed source provides refreshes, and its check mode reports drift.
- [ ] A sync that leaves the catalog invalid, and an invalid catalog with no drift, are refused and write nothing.
- [ ] The committed fixture digests equal the synthetic derivation, and an altered digest is reported.
- [ ] `make baseline-digests` changes nothing.

## Context

- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/assets_sync_test.go`
- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `docs/user-guide/context-driven-development.md`
- instruction: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- creates: `internal/baseline/assets_sync_synthetic_test.go`
- creates: `internal/baseline/assets_sync_upstream_rename_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestAssetSyncFollowsASkillRenamedUpstream|TestAssetSyncCheckReportsDriftForARenameTheRefreshRepairs|TestAssetSyncRefusesARenameNoSetupProvides|TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift|TestTheParityFixtureDigestsFollowTheSyntheticSource|TestASyntheticDigestThatDiffersIsReported|TestAssetsSyncProvenanceAndPreMutationRefusals|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract|TestBaselineCompatibilityCorpus)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAssetSyncFollowsASkillRenamedUpstream TestAssetSyncCheckReportsDriftForARenameTheRefreshRepairs TestAssetSyncRefusesARenameNoSetupProvides TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift TestTheParityFixtureDigestsFollowTheSyntheticSource TestASyntheticDigestThatDiffersIsReported TestAssetsSyncProvenanceAndPreMutationRefusals TestAssetsSyncCompatibilityMatchesMaintainedPythonContract TestBaselineCompatibilityCorpus; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "Both modes validate the catalog as the refresh would leave it" || { printf 'missing phrase in %s: %s\n' docs/user-guide/context-driven-development.md "Both modes validate the catalog as the refresh would leave it" >&2; exit 1; }` — expected: exit 0; before this Task the six new tests do not exist and the sentence is absent, so the command fails.

## References

- `_prd.md` → Goal 1; User Story 1; Core Feature 1; Success Metric 1; Declared breaks (asset sync)
- `_techspec.md` → Interfaces; Measured facts; API Contract 1; Testing Approach 1; Build Order 1
- ADR-0072, ADR-0191

## Result

Implemented the Task's slice for Daemon Verification. `syncAssets` now validates
its produced overlay after building snapshots, including an empty overlay. It
preserves the two refusal messages, finding codes, categories and actions,
and validates before either check-mode drift reporting or refresh writes.

The synthetic source builder and fixture regeneration share the skill-file
content and portable tree digest helpers. Regeneration derives every GitHub
skill's tree digest before computing its setup digest. The user guide now
contains the required sentence about validating the catalog as the refresh
would leave it. No existing assertion was changed.

Focused evidence by acceptance criterion:

- Refresh and check after an upstream rename: the new
  `TestAssetSyncFollowsASkillRenamedUpstream` and
  `TestAssetSyncCheckReportsDriftForARenameTheRefreshRepairs` pass. Before the
  production edit, both failed with the current-catalog `handoff-next`
  outside-setup refusal. The tests rename the target core requirements and
  dispatch entry, and the temporary source directory and all setup lists.
- Invalid produced catalog and invalid catalog without drift:
  `TestAssetSyncRefusesARenameNoSetupProvides` and
  `TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift` pass in both refresh
  and check subtests. They assert the refusal category, code, severity,
  message prefix and action, and compare the complete asset tree preimage.
  The no-drift case first refreshes successfully, then breaks the target core.
- Synthetic derivation and altered digest:
  `TestTheParityFixtureDigestsFollowTheSyntheticSource` and
  `TestASyntheticDigestThatDiffersIsReported` pass. The latter uses a literal
  fixture with one altered GitHub digest and requires a finding naming it.
- No regeneration changes: `GOCACHE=/tmp/roundfix-task01-gocache rtk proxy make
  baseline-digests` exited 0 and reported
  `{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}`.

Focused check command (exit 0):

```sh
GOCACHE=/tmp/roundfix-task01-gocache rtk proxy go test -count=1 ./internal/baseline -run 'TestAssetSync(Follows|CheckReports|Refuses|StillRefuses)|TestTheParityFixtureDigests|TestASyntheticDigest|TestAssetsSync(CheckIsReadOnly|Provenance|CompatibilityMatches)'
```

This also exercises existing read-only, pre-mutation refusal and maintained
Python compatibility assertions. `git diff --check` exited 0. A direct byte
comparison against `git show HEAD:<fixture>` confirmed the parity fixture is
unchanged; `git diff --name-only HEAD -- internal/baseline/assets <fixture>`
was empty. Whitespace-normalized guide inspection confirmed the entire
required sentence.

Repository incremental check: `GOCACHE=/tmp/roundfix-task01-gocache rtk make
verify-incremental` exited 0 on a stable-tree rerun with the process-table
access required by CLI force-stop integration tests. It covered formatting,
vet, package tests, skill synchronization checks, skill checks and the build.
The first sandboxed attempt exited 2: the CLI process-table probe was denied,
and suiteguard detected the final helper/guide/Result edits made during that
attempt. No edits were made during the successful rerun.

The declared `## Verification` command was not run; the Daemon owns it and
Task status. No commit, push or pull request was made. No follow-up scope was
implemented.

## Carry-forward provenance

- Source Run: `run_20261001T013631Z_c59b100cbb9f7af1`
- Source commit: `cdae752d7c5b0dc8b65ad0ddd50f16cdeb1e2090`
