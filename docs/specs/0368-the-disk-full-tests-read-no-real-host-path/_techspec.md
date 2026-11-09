---
spec: 0368-the-disk-full-tests-read-no-real-host-path
prd: _prd.md
created: 2026-11-09
---

# The disk-full tests read no real host path — Technical Spec

## Executive Summary

`doctorDiskSpaceAncestor` in `internal/cli/doctor_disk_space.go` gains one
condition: an `os.Stat` error that is `syscall.ENOTDIR` walks to the parent,
as `os.ErrNotExist` already does. The disk-full tests in
`internal/cli/deliver_disk_full_test.go` build every path that reaches the
filesystem under their own `t.TempDir()`, and two new tests prove the walk
with a file created inside a temporary root. The trade-off accepted is that
a `worktree.location` misconfigured below a file is now measured at its
nearest existing directory by `doctor`'s `disk-space:` line instead of
warning `DR-DISK-UNREAD`; the guide already describes the nearest existing
directory, and the misconfiguration still fails where a worktree is
created.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local filesystem reads only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0370 defines the Disk-Full Park,
  ADR-0370: "`<N>` is its free space at the park", which this Spec makes
  hold when a file stands in the failing path. ADR-0327 and ADR-0334 carry
  ADR-0370's amendment notices on the Shared Build Cache trim; not
  applicable, the trim is unchanged. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: not applicable — no protected tooling configuration,
  script, ignore file, plugin declaration or version pin changes, and no
  Governed Path (`paths: []`). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0368-the-disk-full-tests-read-no-real-host-path/_authorization.md`.

## System Architecture

- `internal/cli/doctor_disk_space.go` (task_01): `doctorDiskSpaceAncestor`
  walks past `os.ErrNotExist` and `syscall.ENOTDIR`; its other branches are
  unchanged. Its three callers, `inspectFreeSpace` (the `doctor`
  `disk-space:` line and the queue's low-disk warning in
  `deliver_free_space.go`) and `diskFullFailure` (`deliver_disk_full.go`,
  the Disk-Full Park), change no code.
- `internal/cli/deliver_disk_full_test.go` (task_01): the three tests that
  reach the filesystem take their paths from `t.TempDir()`; the
  parsing-only `TestADiskFullPathMayHoldSpaces` keeps literals outside
  `/tmp/`; `TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory` is new.
- `internal/cli/doctor_disk_space_test.go` (task_01):
  `TestTheDiskSpaceAncestorWalksPastAFile` is new.

## Implementation Design

### Interfaces

```go
// internal/cli — signature unchanged
// doctorDiskSpaceAncestor returns the nearest existing directory of path.
// A missing component (os.ErrNotExist) or a file where a directory is
// expected (syscall.ENOTDIR) walks to the parent; any other stat error is
// returned as "stat disk-space location %q: %w".
func doctorDiskSpaceAncestor(path string) (string, error)
```

### Data Models

No changes.

### API Contracts

None. No command output, flag, finding code or record changes; the
Disk-Full Park and `doctor` print the same text with a measured value.

### Documentation

None. `docs/user-guide/commands/doctor.md` already says each location is
"read at its nearest existing directory", and
`docs/user-guide/commands/deliver.md` says `<N>` "is its free space when the
item parked"; the code now does both when a file stands in the path.

### Surface Transcripts

None. No command output, flag or finding text changes; the Disk-Full Park
and `doctor` print their existing lines with a measured value.

## Coverage Map

- Goal 1, Story 1, Core Feature 1, Success Metric 1 → Testing Approach
  (the hermetic test paths); Build Order 1.
- Goal 2, Story 2, Core Features 2 and 3, Success Metrics 2 and 3 →
  Interfaces (`doctorDiskSpaceAncestor`); Testing Approach; Build Order 1.

Deviations from the PRD: None.

## Absence Evidence

1. The `ENOTDIR` walk in the free-space lookup — command: `grep -rn 'ENOTDIR' internal --include='*.go'`; result: no match (exit 1) at `cfc08683`.
2. A disk-full test that measures a path below a file — command: `grep -rn 'func Test' internal/cli/deliver_disk_full_test.go internal/cli/doctor_disk_space_test.go`; result: `TestTheItemBuildOnAFullDiskIsADiskFullFailure`, `TestTheRepositoryGateOnAFullDiskParksDiskFull`, `TestDiskFullFailureUsesOnlyTheLogTailAndPreservesTheError`, `TestADiskFullPathMayHoldSpaces` and `TestDoctorWarnsWhenALocationIsLowOnDisk`; none creates a file where a directory is expected.

## ADR Alternatives

- ADR-0370: kept; the design takes none of its rejected options (clearing
  or emptying the cache, `-trimpath`, a shell trim, refusing a Run below the
  floor, moving the cache) and only makes the measured `<N>` hold.
- ADR-0327: not applicable; the Shared Build Cache trim rules are unchanged.
- ADR-0334: not applicable; when Runs and the queue gate trim the cache is
  unchanged.

## Integration Points

None outside the repository. The Delivery Queue's Disk-Full Park, its
low-disk warning and `roundfix doctor` call the changed function without a
change of their own.

## Testing Approach

- `TestADiskFullPathBelowAFileIsReadAtItsNearestDirectory`
  (`deliver_disk_full_test.go`): a file `<root>/x` in a temporary root and a
  retained log naming `<root>/x/_pkg_.a`; `diskFullFailure` returns that
  path with the stubbed bytes, and the stub saw exactly one read, of
  `<root>`. On the authoring base it returns 0 bytes and no read.
- `TestTheDiskSpaceAncestorWalksPastAFile` (`doctor_disk_space_test.go`):
  an existing directory, an absent child, a file and a path below a file
  each resolve to the temporary root; the subtest
  `an unreadable directory still fails` locks a directory (mode 0, restored
  in `t.Cleanup`) and expects an error that is not `os.ErrNotExist`,
  skipped on Windows and for root, where directory modes do not deny a
  lookup.
- The three existing tests that reach the filesystem replace `/tmp/...`
  literals with `{root}` placeholders resolved to their own `t.TempDir()`;
  their assertions are otherwise unchanged.
- No test creates a file in `/tmp`, `/private/tmp` or `$HOME`, starts a
  process beyond the existing fixture's shell commands, or reaches the
  network. New top-level tests call `t.Parallel()` first. The test files
  compile for Windows (`GOOS=windows go vet ./internal/cli`).
- A prototype on `cfc08683` passed these tests and every
  `DiskFull|DiskSpace|LowOnDisk` test of `./internal/cli`; removing the
  `ENOTDIR` condition failed both new tests, and walking past every error
  failed the unreadable-directory subtest.

## Build Order

1. The `ENOTDIR` walk, the hermetic test paths and the two new tests.
2. Final QA gate (depends on: 1).

No step needs another Spec's code.

## Risks & Considerations

- A `worktree.location` below a file no longer warns `DR-DISK-UNREAD` in
  `doctor`; creating a worktree there still fails, and the guide's text
  already matches the new behavior.
- On Windows `os.Stat` reports a file in the path prefix as a missing path,
  which already walked; the change is inert there.

## Glossary

None.

## Decisions

- Walk past `ENOTDIR` and nothing else; no ADR, since ADR-0370 and the
  guides already state the behavior.
- One implementation Task and the QA gate.
