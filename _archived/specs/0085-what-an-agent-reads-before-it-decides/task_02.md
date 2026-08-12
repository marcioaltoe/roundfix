---
task: task_02
spec: 0085-what-an-agent-reads-before-it-decides
status: completed
type: backend
complexity: medium
---

# Task 02: Give the archive layout one owner

## Overview

Five packages compose the archive path themselves, so the layout is five copies
of one fact. This Task adds the resolver that answers where a retired artifact of
each kind lives, with a closed set of kinds so no caller can invent a directory
the checker never reads. No consumer moves yet; that is Task 03.

## Requirements

1. MUST expose a resolver returning the repository-relative directory for a
   retired artifact kind.
2. MUST close the set of kinds to Specs, findings, ADRs, and backlog entries;
   an open set would let a caller invent a directory nothing reads.
3. MUST return the paths Task 01 recorded, so this Task changes no behaviour.
4. MUST NOT change any consumer; the resolver has no callers when this Task
   settles.

## Subtasks

- [ ] Add the kind type and its closed set.
- [ ] Add the resolver.
- [ ] Prove it answers today's paths.

## Acceptance Criteria

- [ ] The resolver returns the recorded directory for each of the four kinds.
- [ ] The kind set is closed and an unknown kind is rejected.
- [ ] No consumer calls the resolver yet.

## Bounded scope

This Task may create or modify only:

- `internal/spec/archive.go`
- `internal/spec/archive_test.go`
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_02.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestArchiveDir' -count=1 -v 2>&1 | grep -q '^--- PASS: TestArchiveDirAnswersEveryRetiredKind'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/spec -run '^TestArchiveDir' -count=1 -v 2>&1 | grep -q '^--- PASS: TestArchiveDirRejectsAnUnknownKind'` — expected: exits 0.
- `grep -q 'func ArchiveDir' internal/spec/archive.go` — expected: exits 0. This function does not exist before this Task.

## References

- `_prd.md` → the archive read path.
- `_techspec.md` → Build Order 2; Interfaces.

## Result

Implemented the closed retired-artifact kind set and its repository-relative
directory resolver. The resolver preserves the four locations characterized by
Task 01 and returns an empty directory for an unknown kind. No consumer moved
onto the resolver in this slice.

- Criterion 1: `TestArchiveDirAnswersEveryRetiredKind` exercises Specs,
  findings, ADRs, and backlog entries and pins their recorded directories.
- Criterion 2: the resolver accepts only the four `ArchiveKind` constants in
  its switch; `TestArchiveDirRejectsAnUnknownKind` proves a fabricated kind
  resolves to no directory.
- Criterion 3: a focused `ArchiveDir(` call-site sweep across `internal/spec`,
  `internal/speccheck`, `internal/specaudit`, `internal/worktree`, and
  `internal/cli` reported only the resolver and its two tests, with no consumer
  call.

Focused checks:

- Before implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-spec0085-task02-gocache go test ./internal/spec -run 'ArchiveDir' -count=1`
  failed because the new type, constants, and resolver were undefined.
- After implementation, the same focused command passed
  (`ok roundfix/internal/spec`).
- `rtk grep -n "ArchiveDir(" internal/spec/*.go internal/speccheck/*.go internal/specaudit/*.go internal/worktree/*.go internal/cli/*.go`
  reported only `internal/spec/archive.go` and `internal/spec/archive_test.go`.

The Daemon-owned `## Verification` commands were not run.
