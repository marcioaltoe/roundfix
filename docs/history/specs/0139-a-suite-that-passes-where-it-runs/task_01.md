---
task: task_01
spec: 0139-a-suite-that-passes-where-it-runs
status: completed
type: test
complexity: medium
---

# Task 01: Stop ACPX fixtures churning links to the test binary

## Overview

The ACPX adapter fixtures stand in for acpx and runtime adapters by re-executing
the test binary. They create a hard link to that binary for every test, and the
test's temporary-directory cleanup deletes the link afterwards. Under the
parallel suite, macOS sometimes kills an exec that lands during that link churn.
The child then dies before its first write, and tests fail on missing fixture
files.

This slice switches the per-test adapter fixture to a symlink. For the
package-directory adapter, which must resolve to a real path, it creates the
links once, before the package's tests run.

## Requirements

1. MUST make the per-test fake adapter fixture (`provisionFakeAdapter`) create a
   symlink to the running test binary instead of a hard link. The fixture's
   sidecar lookup MUST keep working.
2. MUST stop creating a per-test hard link for the package-directory adapter
   (`installSymlinkedPackageAdapter`). Its package-directory links MUST be
   created once, by a helper named `provisionPackageAdapterLinks`, called from
   the package's `TestMain` before the tests run. Those links MUST be removed
   after the tests return.
3. MUST NOT create a hard link to the test binary anywhere outside `TestMain`
   and `provisionPackageAdapterLinks`.
4. MUST keep `TestFixtureBinarySurvivesConcurrentExec` and every existing
   assertion passing.
5. MUST NOT add a timeout, retry or skip, remove `t.Parallel`, or change
   production code.

## Subtasks

- [ ] Switch the per-test adapter fixture to a symlink.
- [ ] Move package-directory adapter links to one-time setup in `TestMain`.
- [ ] Confirm the concurrent-exec fixture test and a stressed package run pass.

## Acceptance Criteria

- [ ] The per-test adapter fixture creates a symlink to the test binary; today
      it creates a hard link.
- [ ] No hard link to the test binary is created outside `TestMain` and
      `provisionPackageAdapterLinks`; today two per-test sites create one.
- [ ] `go test -count=5 -parallel 16 ./internal/agent` passes on this machine,
      where the same stressed run failed before this change.

## Context

- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `awk 'index($0, "func provisionFakeAdapter(") == 1 {inside=1} inside {print} inside && $0 == "}" {inside=0}' internal/agent/acpx_runner_test.go | grep -q 'os.Symlink(os.Args\[0\]'` — the per-test fixture uses a symlink; this fails today.
- `test -f internal/agent/acpx_runner_test.go && ! awk 'index($0, "func TestMain(") == 1 || index($0, "func provisionPackageAdapterLinks(") == 1 {skip=1} skip == 0 {print} skip == 1 && $0 == "}" {skip=0}' internal/agent/acpx_runner_test.go | grep -q 'os.Link(os.Args\[0\]'` — no per-test hard link to the test binary remains; this fails today.
- `grep -q 'func provisionPackageAdapterLinks(' internal/agent/acpx_runner_test.go || exit 1; go test -count=5 -parallel 16 ./internal/agent` — the stressed package run passes with the concurrent-exec fixture test included.

## References

- `_prd.md` → Goal 1; User Story 1; Core Feature 1; Regression locks.
- `_techspec.md` → Implementation Design: ACPX fixtures; Testing Approach 1;
  Build Order 1.

## Result

Implementation:

- `provisionFakeAdapter` now symlinks each per-test command path to the compiled
  test binary. The existing command-path sidecar remains beside that symlink.
- `TestMain` now calls `provisionPackageAdapterLinks` before the package suite,
  reuses its package-directory targets from parallel tests, and removes their
  temporary root after the suite returns. `installSymlinkedPackageAdapter`
  creates only the per-test command symlink and sidecar.

Focused-check evidence:

- Acceptance criterion 1: `rtk rg -n 'os\.(Link|Symlink)\(os\.Args\[0\]' internal/agent/acpx_runner_test.go`
  found the per-test `os.Symlink` call in
  `provisionFakeAdapter`. The focused Go run below exercised its sidecar output
  through `TestFixtureBinarySurvivesConcurrentExec`.
- Acceptance criterion 2: the same source inspection found one `os.Link` call,
  inside `provisionPackageAdapterLinks`. `rtk rg -n 'provisionPackageAdapterLinks|removePackageAdapterLinks|packageAdapterLinks' internal/agent/acpx_runner_test.go`
  confirmed setup and cleanup are wired
  through `TestMain` and the package adapter installer reads the shared targets.
- Acceptance criterion 3: `rtk go test -run '^(TestFixtureBinarySurvivesConcurrentExec|TestCheckAdapterProvesOfficialClaudePackageAndVersion|TestCheckAdapterClassifiesUnreadyClaudeAdapters)$' -count=3 -parallel 16 ./internal/agent`
  passed 36 tests. The first sandboxed
  invocation could not read the Go build cache; the identical focused command
  passed after cache access was granted. The declared five-count whole-package
  stress command remains for Daemon Verification.
- `rtk git diff --check` exited 0.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  those checks and settlement.
