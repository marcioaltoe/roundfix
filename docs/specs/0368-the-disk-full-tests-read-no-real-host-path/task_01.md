---
task: task_01
spec: 0368-the-disk-full-tests-read-no-real-host-path
status: pending
type: backend
complexity: medium
---

# Task 01: The disk-full tests use their own temporary paths and the lookup walks past a file

## Overview

Two disk-full tests wrote the literal path `/tmp/x/_pkg_.a` into their fake
build output, and an empty file `/tmp/x` on the host made the free-space
lookup fail with `ENOTDIR`, so Spec 0365's first Run parked for a defect it
did not touch. This Task answers the adopted Backlog Entry recorded by the
operator on 2026-11-09
([`references/2026-11-09-the-disk-full-tests-read-a-real-host-path.md`](references/2026-11-09-the-disk-full-tests-read-a-real-host-path.md)):
every disk-full test that reaches the filesystem builds its paths under its
own `t.TempDir()`, and `doctorDiskSpaceAncestor` walks past a file that
stands where a directory is expected.

## Requirements

1. MUST, in `internal/cli/doctor_disk_space.go` `doctorDiskSpaceAncestor`,
   walk to the parent when `os.Stat` fails with `syscall.ENOTDIR`, exactly as
   it does for `os.ErrNotExist`; the stat condition becomes
   `!errors.Is(err, os.ErrNotExist) && !errors.Is(err, syscall.ENOTDIR)`.
   Every other branch, the function's signature and both error texts
   (`stat disk-space location %q: %w` and
   `find existing ancestor of %q: %w`) stay as they are. Add a doc comment
   saying that a missing component or a file where a directory is expected
   walks to the parent and any other stat error is returned.
2. MUST, in `internal/cli/deliver_disk_full_test.go`, make every path that
   reaches `diskFullFailure` or `doctorDiskSpaceAncestor` lie under the
   test's own `t.TempDir()`:
   - `TestTheItemBuildOnAFullDiskIsADiskFullFailure`, case
     `path in the log`: output `write <root>/x/_pkg_.a: no space left on device`
     and expected path `<root>/x/_pkg_.a`, `<root>` a `t.TempDir()` of the
     subtest;
   - `TestTheRepositoryGateOnAFullDiskParksDiskFull`, case `full disk`:
     the same output, and the blocker `disk-full: <root>/x/_pkg_.a: 2 MiB free`;
   - `TestDiskFullFailureUsesOnlyTheLogTailAndPreservesTheError`: the paths
     `/tmp/last` and `/tmp/error` become `<dir>/last` and `<dir>/error`
     under the subtest's existing `dir`, including the wrapped
     `syscall.ENOSPC` cause, and `/first` becomes `<dir>/first` in the
     `last path wins` case; the `outside tail` case keeps `/first`, which
     lies outside the read tail and never reaches the filesystem;
   - `TestADiskFullPathMayHoldSpaces`, which only parses text and never
     touches the filesystem, keeps literal paths but names none under
     `/tmp/`: `/tmp/x/_pkg_.a` becomes `/build/x/_pkg_.a` and `/tmp/last`
     becomes `/build/last`.

   Every assertion of those four tests stays as it is apart from the paths.
   After this Task the file holds no `/tmp/` text, comments included.
3. MUST add to `internal/cli/deliver_disk_full_test.go`
   `TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory`: in a
   `t.TempDir()` root, create an empty regular file `<root>/x` and a log
   holding `write <root>/x/_pkg_.a: no space left on device`; call
   `diskFullFailure` with a context whose `availableBytes` records each path
   it reads and returns `5*1024*1024 + 7`; assert `Path` is
   `<root>/x/_pkg_.a`, `AvailableBytes` is `5*1024*1024 + 7`, the error
   wraps the cause, and the recorded reads are exactly `[<root>]`.
4. MUST add to `internal/cli/doctor_disk_space_test.go`
   `TestTheDiskSpaceAncestorWalksPastAFile`: in a `t.TempDir()` root holding
   an empty file `<root>/file`, `doctorDiskSpaceAncestor` returns `<root>`
   with no error for `<root>`, `<root>/absent/child`, `<root>/file` and
   `<root>/file/x/_pkg_.a`; and a subtest named
   `an unreadable directory still fails` creates `<root>/locked`, sets its
   mode to 0 (restoring 0700 in `t.Cleanup`), and asserts that
   `doctorDiskSpaceAncestor(<root>/locked/child)` returns an error that is
   not `os.ErrNotExist`. The subtest calls `t.Skip` when `runtime.GOOS` is
   `windows` or `os.Geteuid()` is 0, where a directory mode does not deny a
   lookup.
5. Both new top-level tests MUST call `t.Parallel()` first. No test may
   create, read or stat a path in `/tmp`, `/private/tmp` or `$HOME` other
   than through `t.TempDir()`; the test files MUST compile for Windows.
6. MUST NOT change `internal/cli/deliver_disk_full.go`,
   `internal/cli/deliver_free_space.go`, any other file, guide, skill, ADR
   or record. `docs/agents/cli.md`, the CLI guide this Task follows, and the
   `doctor` and `deliver` guides are read only: the latter two already
   say each location is read at its nearest existing directory and that the
   Disk-Full Park names the free space when the item parked, which this Task
   makes true when a file stands in the path.

## Subtasks

- [ ] Walk past `ENOTDIR` in `doctorDiskSpaceAncestor`.
- [ ] Move the disk-full test paths under their temporary directories.
- [ ] Add the two tests.

## Acceptance Criteria

- [ ] A failing path below a file is measured at the nearest existing directory, and an unreadable directory still fails the lookup.
- [ ] No disk-full test names a `/tmp/` path, and the four existing tests pass with their assertions unchanged.

## Context

- interface: `internal/cli/doctor_disk_space.go`
- interface: `internal/cli/deliver_disk_full_test.go`
- interface: `internal/cli/doctor_disk_space_test.go`
- interface: `docs/agents/cli.md`
- instruction: `docs/user-guide/commands/doctor.md`
- instruction: `docs/user-guide/commands/deliver.md`

## Consumers

- `doctorDiskSpaceAncestor`, read by `inspectFreeSpace` (`doctor` and the queue's low-disk warning) and `diskFullFailure` (the Disk-Full Park) — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory|TestTheDiskSpaceAncestorWalksPastAFile)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory TestTheDiskSpaceAncestorWalksPastAFile; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run 'DiskFull|DiskSpace|LowOnDisk' ./internal/cli` — expected: exit 0. Before this Task neither test exists, so their pass lines are missing; afterwards both pass with every disk-full, disk-space and low-disk test of the package.
- `if grep -nE -- 'tm[p]/' internal/cli/deliver_disk_full_test.go; then echo 'deliver_disk_full_test.go still names a shared temporary path' >&2; exit 1; fi; go test -count=1 -timeout 20m -run '^(TestTheItemBuildOnAFullDiskIsADiskFullFailure|TestTheRepositoryGateOnAFullDiskParksDiskFull|TestDiskFullFailureUsesOnlyTheLogTailAndPreservesTheError|TestADiskFullPathMayHoldSpaces|TestDoctorWarnsWhenALocationIsLowOnDisk)$' ./internal/cli && GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the file names `/tmp/x/_pkg_.a`, so the first check fails; afterwards no `/tmp/` path remains, the existing tests pass and the package vets for Windows.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Interfaces; System Architecture; Testing Approach; Build Order 1
- ADR-0370
