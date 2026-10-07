---
task: task_02
spec: 0245-generated-records-that-hold-across-specs-and-platforms
status: completed
type: test
complexity: high
---

# Task 02: The Coverage Record lists every release platform and is the same bytes on any host

## Overview

`TestCoverageEquivalence` builds the Coverage Record from `go test -list` on
the host, so a record written on macOS fails Linux CI and the reverse. This
Task collects the record statically for every `GOOS` in
`dist/npm/platforms.json`, lists each platform-limited test with its
platforms, compares the full record on any host, and anchors the collection
to `go test -list` on the host. It answers problem 2 of the adopted Backlog
Entry
[generated records break when Specs are authored in parallel or tested on another platform](references/2026-10-07-generated-records-that-parallel-work-or-another-platform-breaks.md)
of 2026-10-07.

## Requirements

1. MUST change `internal/spec/coverage_test.go` to the `CoverageRecord` shape
   and functions of `_techspec.md` → Interfaces (Coverage Record) and to
   Invariants 6-12. `collectToolchainCoverage` keeps today's `go test -list`
   parse unchanged, and `TestCoverageEquivalence` and
   `-update-coverage-record` use the static collection.
2. MUST read the platforms only from `dist/npm/platforms.json`, run `go list`
   with `GOOS`, `CGO_ENABLED=0` and `GOWORK=off` set per platform, and run no
   test binary built for another system. Data comes from the repository's
   own package list and test files, and is written only into the record.
3. MUST keep a record without `platforms` comparing as one implicit platform,
   so `TestCompareCoverageRecordsReportsMissingTest`,
   `TestCompareCoverageRecordsReportsAddedTestWithoutRegression`,
   `TestMarshalCoverageRecordIsDeterministic` and
   `TestCoverageRecordCountsNoPackageUnderDocs` pass without an edit.
4. MUST add the tests of `_techspec.md` → Testing Approach 2 to the new file
   `internal/spec/coverage_platform_test.go`.
   `TestCoverageCollectionListsEachPlatformOnlyTest` builds a module under
   `t.TempDir()` with its own `go.mod`, `dist/npm/platforms.json` and test
   files built only on `darwin`, only on `linux`, only on `windows`, and on
   `unix`, and asserts the same record whatever the host is.
5. MUST re-record `docs/references/coverage-record.json` with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   and never edit it by hand. A second run reproduces the same bytes.

## Subtasks

- [ ] Add the static collection, the platform fields and their validation.
- [ ] Compare per platform with today's messages kept for all-platform tests.
- [ ] Add the fixture-module, anchor, regression, validation and matrix tests.
- [ ] Re-record the Coverage Record.

## Acceptance Criteria

- [ ] The record lists the macOS-only detach test under `darwin`, the
      Linux-only store tests under `linux` and the Windows-only store tests
      under `windows`, and re-recording reproduces the same bytes.
- [ ] A test lost on one platform is a regression on every host, named with
      that platform.
- [ ] On the host, the static collection equals `go test -list` with no
      difference.

## Context

- interface: `internal/spec/coverage_test.go`
- creates: `internal/spec/coverage_platform_test.go`
- interface: `docs/references/coverage-record.json`
- instruction: `dist/npm/platforms.json`
- instruction: `internal/cli/detach_fixture_group_darwin_test.go`
- instruction: `internal/store/process_linux_test.go`
- instruction: `internal/store/process_windows_test.go`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`

## Verification

- `out="$(go test -count=1 -v ./internal/spec -run '^(TestCoverageEquivalence|TestCoverageCollectionMatchesTheToolchainOnThisPlatform|TestCoverageCollectionListsEachPlatformOnlyTest|TestCompareCoverageRecordsReportsAPlatformRegression|TestCoverageRecordRefusesAMalformedPlatformEntry|TestCoveragePlatformsComeFromTheReleaseMatrix|TestCompareCoverageRecordsReportsMissingTest|TestCompareCoverageRecordsReportsAddedTestWithoutRegression|TestMarshalCoverageRecordIsDeterministic|TestCoverageRecordCountsNoPackageUnderDocs)$' 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCoverageEquivalence TestCoverageCollectionMatchesTheToolchainOnThisPlatform TestCoverageCollectionListsEachPlatformOnlyTest TestCompareCoverageRecordsReportsAPlatformRegression TestCoverageRecordRefusesAMalformedPlatformEntry TestCoveragePlatformsComeFromTheReleaseMatrix TestCompareCoverageRecordsReportsMissingTest TestCompareCoverageRecordsReportsAddedTestWithoutRegression TestMarshalCoverageRecordIsDeterministic TestCoverageRecordCountsNoPackageUnderDocs; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four new tests and the anchor test do not exist, so the command fails.
- `t="$(mktemp -d)" || exit 1; tr -s '[:space:]' ' ' < docs/references/coverage-record.json > "$t/normalized.txt" || exit 1; for phrase in '"platforms": [ "darwin", "linux", "windows" ]' '"TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded": [ "darwin" ]' '"TestParseProcStatStartTimeRejectsMissingField": [ "linux" ]' '"TestWindowsProcessStartParsesRecordedIdentity": [ "windows" ]'; do grep -qF -- "$phrase" "$t/normalized.txt" || { printf 'missing in the Coverage Record: %s\n' "$phrase" >&2; exit 1; }; done; cp docs/references/coverage-record.json "$t/record.json" || exit 1; go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1 >/dev/null || exit 1; cmp docs/references/coverage-record.json "$t/record.json"` — expected: exit 0; before this Task the record has no `platforms` field and lists no Linux- or Windows-only test, so the command fails.

## References

- `_prd.md` → Goal 3; Goal 4; Core Feature 2; Success Metric 4; Success Metric 5
- `_techspec.md` → Measured starting point; Interfaces; Invariants 6-12; Data Models; API Contract 3; Testing Approach 2; Build Order 2
- ADR-0250

## Result

Implemented the Task 02 slice: the Coverage Record now collects the sorted,
unique release `goos` values from `dist/npm/platforms.json`, lists package test
files with `GOOS`, `CGO_ENABLED=0` and `GOWORK=off`, and parses each file once
across platforms. Common tests remain in `packages`; limited tests carry their
platform sets in `platformTests`. No foreign test binary runs. The existing
`go test -list` parser remains in `collectToolchainCoverage` for the host anchor.
Comparison reports lost and gained platforms and preserves the existing
all-platform and package messages. Legacy records retain one implicit platform.
Validation and marshaling cover platform tokens, ordering, proper subsets,
unknown packages, duplicate placement and deterministic copies.

Starting evidence: the checked-in record had no `platforms` field or required
Darwin, Linux and Windows entries; `coverage_platform_test.go` did not exist.
The pre-existing Task file change was the Daemon's `pending` to `in_progress`
transition. Task status, checkboxes, authored Verification and the Task Graph
were left unchanged.

### Acceptance evidence

1. Release-platform entries and repeatable bytes: the required generator
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   generated `docs/references/coverage-record.json`; it was not hand-edited.
   JSON inspection confirmed `platforms: [darwin, linux, windows]`,
   `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded: [darwin]`,
   `TestParseProcStatStartTimeRejectsMissingField: [linux]` and
   `TestWindowsProcessStartParsesRecordedIdentity: [windows]`, with 48
   platform-limited tests total. A second generator run exited 0. Comparing
   saved first-run bytes with the second output found them identical, SHA-256
   `9b8b3b0a343c0a281048f0446a52e6d0709b6531a1e535b8b67c8502f30b30da`.
2. Platform-specific regression: `TestCompareCoverageRecordsReportsAPlatformRegression`
   passed cases for losing only Linux from a common test, losing a Unix-only
   test, losing all platforms, additions and release-platform/package removal.
   The Linux-only loss reports exactly
   `coverage regression: package "roundfix/example" no longer executes "TestKept" on linux`.
   These comparisons contain no host-dependent branch. The fixture-module test
   passed the same explicit expected Darwin, Linux, Windows and Unix record,
   including external tests, aliased imports, invalid-signature exclusions,
   a Linux-only package and docs exclusion.
3. Host anchor: `TestCoverageCollectionMatchesTheToolchainOnThisPlatform`
   passed on `darwin` with no additions or regressions against `go test -list`.

### Focused checks

Checks used `GOCACHE=/private/tmp/roundfix-0245-task02-gocache GOPROXY=off`.
The first host-anchor attempt using the shared cache failed with
`operation not permitted` reading a cache entry; the task-scoped cache resolved
that environment restriction.

- `go test ./internal/spec -run '^(TestCoverageCollection|TestCompareCoverageRecords|TestCoverageRecordRefuses|TestCoveragePlatforms|TestMarshalCoverageRecord|TestCoverageRecordCounts)' -count=1 -v`
  — exit 0; nine top-level checks passed, including the unchanged legacy
  missing-test, added-test, marshal and docs-package tests. Log:
  `/private/tmp/roundfix-0245-task02-focused.log`.
- Required generator, second invocation — exit 0; byte equality checked in
  Python against `/private/tmp/roundfix-0245-task02-record.json`.
- `GOCACHE=/private/tmp/roundfix-0245-task02-gocache GOPROXY=off rtk make verify-incremental`
  — first run exited 2 at `make test`; only
  `TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` and
  `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`
  failed because the sandbox denied reading the process table. The coverage
  package and every other package passed. Rerun with host process-table
  permission exited 0, including formatting, vet, tests, skill checks and
  CLI build. Logs: `/private/tmp/roundfix-0245-task02-incremental.log` and
  `/private/tmp/roundfix-0245-task02-incremental-permitted.log`.
- `git -c core.fsmonitor=false diff --check` — exit 0 after the Result append.
  Postflight paths are exactly the collector, new platform test file,
  generated record and assigned Task file.

### Follow-up outside this Task

The first generator invocation passed `TestCoverageEquivalence` and wrote the
record, then exited 1 because the existing `internal/spec/main_test.go` suite
guard reported `modified: docs/references/coverage-record.json`. The second
invocation passed because the bytes were unchanged. No guard was disabled.
For a future regeneration that changes the record, the sanctioned-regeneration
integration needs a command declaration, explicit output authority for this
record, and registration by the recording test. The current Spec authorization's
`Sanctioned regeneration` section declares only `make baseline-digests`, and
`TestCoverageEquivalence` does not call `suiteguard.DeclareSanctionedRegeneration`.
Those authorization and guard-integration changes require a separately scoped
follow-up; this Task does not alter authorization records or guard wiring.

Daemon Verification remains pending; neither authored Verification command was
executed. No commit, push or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261007T115541Z_03741c3fcc656349`
- Source commit: `55555df76f71c617dd46c6df6f1c873fa07755b7`
