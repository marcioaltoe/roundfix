---
task: task_01
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
status: pending
type: backend
complexity: high
---

# Task 01: The asset sync imports the upstream python preset at the commit every setup names

## Overview

Upstream published the curated `python` preset and, at
`65d5f365cace27c4b468db1723213723b10c870b`, replaced its `handoff` skill with
six ai-memory skills. Spec 0317, delivered first, has moved every other
upstream Setup Snapshot to that commit. This Task imports `python` at the same
commit through the asset sync, the only writer of a Setup Snapshot
(ADR-0191), seeding the new snapshot as Spec 0208 seeded its renamed ones;
the sync leaves the other snapshots unchanged, so the embedded snapshots keep
agreeing on every external skill (ADR-0241). The parity row, the derived
files and this repository's Setup Manifest follow through their sanctioned
commands; no skill is restored (`_prd.md` → Core Feature 1; ADR-0322).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-22-a-baseline-profile-for-python-repositories.md`](references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected item 5 (the profile takes the curated Python
   skills).
2. MUST first check that the local checkout `~/dev/skills` holds commit
   `65d5f365cace27c4b468db1723213723b10c870b` and that
   `internal/baseline/assets/setups/go.json` already names it (Spec 0317 has
   merged); when either fails, MUST stop and report before writing anything.
   MUST NOT fetch from or otherwise reach a network remote.
3. MUST create `internal/baseline/assets/setups/python.json` as the bytes of
   `internal/baseline/assets/setups/rust.json` with exactly the two
   replacements of `_techspec.md` → Data Model 1, and then run, verbatim and
   in order, the commands of `_techspec.md` → API Contract 1: the fetch into
   an empty repository, the binary build, the sync, the `--check` (exit
   `0`), the parity transform of
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
   `make baseline-digests`, and the managed refresh twice (the second
   reports `File changes: 0`). Every command that runs the binary gets a
   fresh temporary `HOME`; the binary is built under the Run's own
   environment and rebuilt after the sync so the refresh reads the new
   snapshot.
4. MUST hand-edit no snapshot beyond the seed's two values, no derived file
   and no parity digest, and MUST restore no skill. The sync must leave
   `go.json`, `rust.json`, `typescript.json` and `go-cli-typescript-bun.json`
   byte-identical; a change to them, or any file these commands change that
   this Task does not declare, is a stop-and-report condition.
5. MUST make the declared breaks of `_techspec.md` → Testing Approach (Build
   Order 1) and no other change to existing tests:
   `internal/baseline/assets_sync_test.go` (the four counts of four become
   five) and `internal/baseline/derived_ownership_test.go`
   (`TestOutputsForCommand` adds
   `internal/baseline/assets/setups/python.json`).
6. MUST create `internal/baseline/python_setup_test.go` with
   `TestThePythonSetupIsThePublishedUpstreamPreset` and
   `TestEveryUpstreamSetupFollowsThePythonPresetRevision` as `_techspec.md` →
   Testing Approach describes them, with API Contract 2's 52 names as a
   literal list. Each top-level test calls `t.Parallel()` first.
7. MUST NOT change any Baseline module, profile, template, coverage entry,
   other Setup Snapshot, Go production file, skill, `skills-lock.json`, user
   guide page, `.roundfixrc.yml`, the Makefile or `go.mod`/`go.sum`.

## Subtasks

- [ ] Check the commit and Spec 0317's snapshots, seed `python.json`, sync and check.
- [ ] Transform the parity row, regenerate the digests and refresh twice.
- [ ] Make the declared test changes and add the two tests.

## Acceptance Criteria

- [ ] The `python` snapshot is the 52-skill upstream preset at `65d5f36`, and every upstream snapshot names that commit.
- [ ] The sync's `--check` exits `0`, the other snapshots keep their bytes, and the embedded snapshots agree on every external skill.
- [ ] This repository's skills and lock are unchanged, it still holds every required skill at the `go` snapshot, and a second managed refresh reports no file change.

## Consumers

- the `python` Setup Snapshot — consumers: `./internal/baseline`, `./skills`

## Context

- instruction: `docs/adr/0322-a-python-profile-imports-the-published-upstream-python-preset.md`
- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `docs/adr/0241-a-repository-profile-takes-its-skill-contracts-from-the-embedded-setup-snapshots.md`
- instruction: `internal/baseline/assets_sync.go`
- instruction: `internal/baseline/assets/setups/go.json`
- creates: `internal/baseline/assets/setups/python.json`
- creates: `internal/baseline/python_setup_test.go`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/assets_sync_test.go`
- interface: `internal/baseline/derived_ownership_test.go`

## Verification

- `for name in TestThePythonSetupIsThePublishedUpstreamPreset TestEveryUpstreamSetupFollowsThePythonPresetRevision; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/baseline && go test -count=1 ./skills` — expected: exit 0. Before this Task neither test exists, so the command fails at its first check. Afterwards both exist and the whole package passes: the snapshots agree on every external skill, this repository holds its required skills at the `go` snapshot, the asset-sync parity row matches, and the derived-file characterizations are current; the skills package passes unchanged.
- `test -f internal/baseline/assets/setups/python.json || { echo 'missing python setup' >&2; exit 1; }; for f in go rust typescript python; do grep -qF -- '"ref": "65d5f365cace27c4b468db1723213723b10c870b"' "internal/baseline/assets/setups/$f.json" || { printf '%s not at the preset commit\n' "$f" >&2; exit 1; }; done; grep -qF -- '"name": "ai-memory-handoff"' internal/baseline/assets/setups/python.json || { echo 'python setup lacks the ai-memory skills' >&2; exit 1; }; go test -count=1 -run '^(TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill|TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot|TestAssetsSyncCompatibilityMatchesMaintainedPythonContract)$' ./internal/baseline` — expected: exit 0. Before this Task the `python` snapshot does not exist, so the command fails at its first check. Afterwards every upstream snapshot names the preset commit, the `python` snapshot holds the ai-memory skills, and the agreement, required-skill and parity tests pass.

## References

- `_prd.md` → Goals 1-2; Core Feature 1; Success Metrics 1, 2 and 5; User Story 4
- `_techspec.md` → Current behavior 1-5; Data Models 1-2; API Contracts 1-2; Invariants 1-3; Testing Approach; Build Order 1
- ADR-0322; ADR-0191; ADR-0241; ADR-0206; ADR-0072
