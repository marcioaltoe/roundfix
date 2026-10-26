---
task: task_02
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
status: pending
type: backend
complexity: medium
---

# Task 02: Doctor warns when a worktree or cache filesystem is low on disk

## Overview

`roundfix doctor` reads the Run Database's free pages but never a
filesystem's free space, so it said nothing while the disk filled twice.
This Task exports the free-space reader `gc compact` already uses, and adds a
`disk-space:` check that warns, after `storage:`, when the filesystem holding
the worktree location or a declared Shared Build Cache has less than 10 GiB
free or cannot be read. A healthy machine sees Doctor's output unchanged, and
no existing test reads the host's disk.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md`](references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md)
   of 2026-10-26 (`_prd.md` → Sources 1) by implementing API Contract 6 of
   `_techspec.md` with the line and rules of Interfaces.
2. MUST add `store.AvailableBytes(path)` in `internal/store/available_bytes.go`
   by calling the existing `availableTemporaryCapacity`; no second
   platform reader.
3. MUST add the `availableBytes` command dependency (default
   `store.AvailableBytes`) in `internal/cli/cli.go`, `HealthCheckDiskSpace`
   in `internal/cli/health.go`, the `diskSpace` doctor dependency with its
   default in `internal/cli/doctor.go`, called after `storage` and before
   `linux-lane` only when not nil, and the check in
   `internal/cli/doctor_disk_space.go`: the worktree location resolved with
   `roundconfig.ResolveWorktreeLocation`, then each declared cache resolved
   as GC resolves it; each read at its nearest existing ancestor through the
   `availableBytes` dependency; `DR-DISK-LOW` below 10,737,418,240 bytes
   with the free MiB rounded down; `DR-DISK-UNREAD` with the error;
   `next: roundfix gc` once when any location is low; the line printed only
   with a finding, always `warn`.
4. MUST make `commandEnvironmentForTest` and
   `updateCommandDependenciesForTest` in `internal/cli/cli_test.go` install an
   `availableBytes` fake that reports 1 TiB, so no existing test reads the
   host's free space.
5. MUST add `internal/store/available_bytes_test.go` with
   `TestAvailableBytesReadsTheFilesystemOfAPath` and
   `internal/cli/doctor_disk_space_test.go` with
   `TestDoctorWarnsWhenALocationIsLowOnDisk` (subtests `below the threshold`,
   `at the threshold`, `unreadable`, `missing location`), with the inputs
   and assertions of `_techspec.md` → Testing Approach; each top-level test
   calls `t.Parallel()` first or carries a `// Sequential: <reason>` comment.
6. MUST keep every existing doctor test passing unchanged.
7. MUST NOT change any document, `internal/store/compaction_capacity_unix.go`
   or `internal/store/compaction_capacity_windows.go`, or `gc compact`.

## Subtasks

- [ ] Export the free-space reader.
- [ ] Add the command and doctor dependencies and the hermetic test default.
- [ ] Add the `disk-space:` check and its tests.

## Acceptance Criteria

- [ ] Doctor prints the exact `disk-space:` line below 10,240 MiB free, between `storage:` and the next line, and keeps its exit code.
- [ ] At exactly 10,240 MiB, and on a healthy machine, Doctor's output is unchanged; an unreadable location reports `DR-DISK-UNREAD`.
- [ ] The free space is read at the nearest existing ancestor, and no existing test reads the host's disk.

## Consumers

- `internal/store.AvailableBytes` — consumers: `./internal/cli`
- the `doctorDependencies.diskSpace` and `commandDependencies.availableBytes` seams — consumers: `./internal/cli`

## Context

- instruction: `internal/store/compaction_capacity_unix.go`
- instruction: `internal/store/compaction_capacity_windows.go`
- instruction: `internal/cli/doctor_storage.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `internal/cli/doctor_linux_lane_test.go`
- instruction: `internal/cli/readiness.go`
- creates: `internal/store/available_bytes.go`
- creates: `internal/store/available_bytes_test.go`
- creates: `internal/cli/doctor_disk_space.go`
- creates: `internal/cli/doctor_disk_space_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAvailableBytesReadsTheFilesystemOfAPath|TestDoctorWarnsWhenALocationIsLowOnDisk|TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn|TestDoctor.*Storage.*)$' ./internal/store ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAvailableBytesReadsTheFilesystemOfAPath TestDoctorWarnsWhenALocationIsLowOnDisk/below_the_threshold TestDoctorWarnsWhenALocationIsLowOnDisk/at_the_threshold TestDoctorWarnsWhenALocationIsLowOnDisk/unreadable TestDoctorWarnsWhenALocationIsLowOnDisk/missing_location TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go build -buildvcs=false ./internal/store ./internal/cli && GOOS=linux go build -buildvcs=false ./internal/store ./internal/cli` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check; with the threshold compared with less-than-or-equal the at-the-threshold subtest fails, with the next action per finding the below-the-threshold subtest fails, and with the configured path read instead of its ancestor the missing-location subtest fails. Afterwards the new and existing doctor tests pass and both packages build for Windows and Linux.

## References

- `_prd.md` → Goals; Core Feature 3; User Experience; Success Metric 5
- `_techspec.md` → Current behavior 5-6; Interfaces; Data Models; API Contract 6; Invariant 5; Testing Approach; Build Order 3
