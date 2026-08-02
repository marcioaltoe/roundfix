---
task: task_01
spec: 0071-verification-cost
status: completed
type: test
complexity: medium
---

# Task 01: Record which tests the suite executes

## Overview

Every later slice changes how tests run, and the Spec's one non-negotiable is
that coverage does not move. This Task records which test functions the suite
executes, per package, so that promise becomes an assertion rather than a
claim. It changes no behavior; it is the rail the rest of the Spec is measured
against, and it only works if it lands first.

## Requirements

1. MUST record, per package, the sorted set of test function names the suite
   executes, in a deterministic form.
2. MUST cover every package the repository builds, not only the two heavy ones.
3. MUST fail when a previously recorded test function disappears or is renamed,
   naming the package and the missing name.
4. MUST report a newly added test function without failing, since adding
   coverage is not a regression.
5. MUST be regenerable through an explicit flag, so an intended change is
   re-recorded deliberately.
6. MUST NOT change any production behavior, test behavior, or exported API.

## Subtasks

- [ ] Enumerate executed test functions per package.
- [ ] Record them deterministically as a golden.
- [ ] Fail on disappearance, report on addition.
- [ ] Add the explicit regeneration flag.

## Acceptance Criteria

- [ ] The record covers every package returned by the repository's package list.
- [ ] Comparing on the unmodified tree passes.
- [ ] Deliberately removing one recorded test name makes the comparison fail and
      name that package and test.
- [ ] Adding a new test function is reported and does not fail.
- [ ] Two consecutive comparisons produce the same result and rewrite no golden.
- [ ] `git status --porcelain` shows no path outside this Spec's folder,
      `internal/`, and this task file.

## Context

- interface: `internal/baseline/catalog_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0; the recorded set matches the tree.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0 on a second consecutive run, proving the comparison is
  stable and self-recording is gated.
- `go test ./internal/spec -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 6; Non-Goals (no test deleted, skipped, or
  weakened).
- `_techspec.md` → Implementation Design: CoverageRecord; Build Order 1.

## Result

Implemented a test-only coverage-equivalence harness in `internal/spec`. It
uses `go list ./...` as the package authority, records each package's sorted
top-level `Test...` names in `coverage-record.json`, fails with package/test
diagnostics for removals, logs additions without failing, and rewrites the
record only when `-update-coverage-record` is explicit. No production file,
existing test body, or exported API changed.

Focused checks:

- `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false go test ./internal/spec -run '^Test(CompareCoverageRecords|MarshalCoverageRecord)' -count=1`
  exited 0. The focused cases exercise missing-test diagnostics, non-failing
  addition reporting, and deterministic marshaling.
- `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -update-coverage-record`
  exited 0 and deliberately created the record through its explicit update
  flag.
- Two consecutive non-update runs of
  `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1`
  exited 0. The record's SHA-256 was
  `286c812bf3b144027e2d7cc0dd156b79e4e0fe90812c02bd21c1800fa5c1a2d1`
  before and after both runs.
- Final focused run
  `GOCACHE="$PWD/.gocache" GOFLAGS=-buildvcs=false go test ./internal/spec -run 'CoverageRecord|CoverageEquivalence' -count=1`
  exited 0 after the last Go edit.

Acceptance evidence:

1. The record contains all 24 packages returned by the repository package
   list, including the package with no test functions, and 1,493 test names.
2. The real-tree comparison exited 0 in each non-update focused run.
3. `TestCompareCoverageRecordsReportsMissingTest` removes `TestRemoved` from
   the actual set and requires a regression naming both
   `roundfix/internal/spec` and `TestRemoved`.
4. `TestCompareCoverageRecordsReportsAddedTestWithoutRegression` requires the
   addition report for `TestAdded` while requiring zero regressions.
5. The two comparison runs produced the same result and preserved the exact
   golden hash; deterministic marshaling also has focused coverage.
6. `git -c core.fsmonitor=false status --short` showed only this task file,
   `docs/specs/0071-verification-cost/coverage-record.json`, and
   `internal/spec/coverage_test.go`.

The Daemon-owned commands under `## Verification` were not run. No follow-up
outside this Task's slice was discovered.
