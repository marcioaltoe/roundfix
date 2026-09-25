---
task: task_04
spec: 0171-a-deterministic-suite-under-load
status: completed
type: backend
complexity: medium
---

# Task 04: No live release lookup in tests

## Overview

`maybeReportVersionFreshness` in `internal/cli/upgrade.go` skips only for a
development version or an empty HOME. Tests run the checked-in release version,
so every test of `fetch`, `resolve`, `watch` or `implement` whose version cache
is stale calls `app.LatestRelease`, which runs `gh api` against GitHub. Only
`internal/cli/upgrade_test.go` injects a fake. Spec 0165's QA gate hit
`api.github.com` again. The data comes from GitHub and lands in
`<HOME>/.roundfix/version-check.json`, which the production freshness check reads
back on every later command in that HOME.

## Requirements

1. MUST add `internal/cli/version_freshness_isolation_test.go` whose `init`
   assigns `versionFreshnessDeps.latestRelease` to an offline lookup that
   returns the tag `v0.0.0-offline-test-lookup`, returns no assets and no error,
   and starts no process. Nothing else may reassign the package default; per-test
   fakes keep using `withVersionFreshnessFakeDeps`.
2. MUST keep the offline tag below every release, so a cache holding it can never
   produce an upgrade warning, and MUST NOT change `upgrade.go`,
   `internal/app/release.go` or any other production file.
3. MUST define `seedFreshVersionCache(t, homeDir)` in the new file, writing
   `<HOME>/.roundfix/version-check.json` with a current `checked_at` and no
   `latest_version`, and MUST make `runRoundfixBinaryMacro` in
   `internal/cli/implement_test.go` call it before it runs the built binary.
4. MUST make
   `TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch` in
   `internal/daemon/run_disposition_characterization_test.go` seed the same
   fresh cache in its HOME before running the built binary, and assert that the
   cache is byte-identical afterwards.
5. MUST prove with separate named tests, each in a fresh test-owned HOME:
   - `TestSuiteDefaultReleaseLookupStaysOffline`: the default dependencies'
     release lookup returns the offline tag;
   - `TestOperationalCommandsRecordTheOfflineLookup`: `fetch`, `resolve`,
     `watch` and `implement` run in process through the suite's default
     dependencies write `latest_version: 0.0.0-offline-test-lookup` into their
     HOME's cache;
   - `TestHelperProcessRecordsTheOfflineLookup`: `implement` run through
     `runCLIHelper` writes the same value;
   - `TestBuiltBinaryLeavesTheSeededVersionCacheUntouched`: a built binary run
     for `implement` leaves the seeded cache byte-identical and prints no
     upgrade warning;
   - `TestPerTestReleaseLookupOverridesTheSuiteDefault`: a per-test fake that
     returns a newer release still produces the upgrade warning, the mirror of
     the offline default.
6. MUST NOT edit `internal/cli/cli_test.go` or
   `docs/references/coverage-record.json`, both governed, and MUST keep every
   recorded top-level test name.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Each of the four operational commands records the offline tag in its
      test HOME, in process and through the helper subprocess.
- [ ] A built binary leaves a seeded fresh cache byte-identical, in
      `internal/cli` and in the Daemon characterization test.
- [ ] A per-test fake still overrides the suite default.

## Context

- creates: `internal/cli/version_freshness_isolation_test.go`
- interface: `internal/cli/upgrade_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/daemon/run_disposition_characterization_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSuiteDefaultReleaseLookupStaysOffline|TestOperationalCommandsRecordTheOfflineLookup|TestHelperProcessRecordsTheOfflineLookup|TestBuiltBinaryLeavesTheSeededVersionCacheUntouched|TestPerTestReleaseLookupOverridesTheSuiteDefault|TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch)$" ./internal/cli ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSuiteDefaultReleaseLookupStaysOffline TestOperationalCommandsRecordTheOfflineLookup TestHelperProcessRecordsTheOfflineLookup TestBuiltBinaryLeavesTheSeededVersionCacheUntouched TestPerTestReleaseLookupOverridesTheSuiteDefault TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done; grep -q "version-check.json" internal/daemon/run_disposition_characterization_test.go && grep -q "seedFreshVersionCache(" internal/cli/implement_test.go` — expected: exit 0; before this Task none of the five new tests exists and neither file seeds the version cache, so the command fails.

## References

- [_prd.md](_prd.md) — Core Feature 4; Success Metric 4
- [_techspec.md](_techspec.md) — The release lookup; API Contract 3

## Result

Implementation:

- The `internal/cli` test package now installs one offline suite default during
  package initialization. It returns `v0.0.0-offline-test-lookup`, no assets
  and no error; the tag sorts below `0.0.0`, and the existing per-test
  dependency override remains the only other test seam.
- `seedFreshVersionCache` writes a current `checked_at` with no
  `latest_version`. The built-binary macro calls it before every invocation,
  and the Daemon characterization seeds the equivalent cache before its built
  `resolve` command and compares the bytes afterwards.
- Five named tests cover the suite default, all four in-process operational
  commands, the helper subprocess, the built binary and the per-test override.
  No production file, governed CLI test, recorded top-level test name or
  coverage record changed.

Focused checks:

- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1
  -run '^(TestSuiteDefaultReleaseLookupStaysOffline|TestOperationalCommandsRecordTheOfflineLookup|TestHelperProcessRecordsTheOfflineLookup|TestBuiltBinaryLeavesTheSeededVersionCacheUntouched|TestPerTestReleaseLookupOverridesTheSuiteDefault)$'
  ./internal/cli` passed (exit 0, 2.164 s).
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1
  -run '^TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch$'
  ./internal/daemon` passed (exit 0, 2.343 s).
- `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache go test -count=1
  -run '^TestAgentSelectionProfilesMacro/invalid_task_type_blocks_every_proof_and_run_side_effect$'
  ./internal/cli` passed (exit 0, 1.569 s), exercising an existing built-binary
  macro path after it gained cache seeding.
- The first sandboxed `rtk env GOCACHE=/private/tmp/roundfix-task04-gocache
  make verify-incremental` attempt reached the package tests but failed because
  two force-stop integration tests could not read the host process table
  (`operation not permitted`). The single failing test reproduced the same
  environment restriction. Re-running the unchanged incremental gate with
  process-table access passed vet, all package tests, skill checks and the build
  (exit 0).
- `rtk git diff --check` passed. The authored `## Verification` command was
  not run; the Daemon owns it.

Acceptance evidence:

- Offline recording: `TestOperationalCommandsRecordTheOfflineLookup` passed
  for `fetch`, `resolve`, `watch` and `implement`, with a new HOME per subtest;
  `TestHelperProcessRecordsTheOfflineLookup` passed through `runCLIHelper` in
  its own HOME. Every cache decoded to
  `latest_version: 0.0.0-offline-test-lookup`.
- Built binaries: `TestBuiltBinaryLeavesTheSeededVersionCacheUntouched` and
  `TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch`
  passed after comparing the seeded cache before and after the command byte for
  byte. The CLI test also observed no upgrade warning.
- Override: `TestPerTestReleaseLookupOverridesTheSuiteDefault` passed after a
  per-test `v1.1.0` fake replaced the suite default, wrote `1.1.0` to the cache
  and printed the expected upgrade warning.

Follow-ups: none.
