---
task: task_02
spec: 0245-generated-records-that-hold-across-specs-and-platforms
status: pending
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
