---
task: task_02
spec: 0074-git-spawn-economy
status: completed
type: backend
complexity: high
---

# Task 02: Batch object reads in skills restore

## Overview

`internal/baseline/skills_restore_git.go` reads acquired skill files by
spawning `git cat-file blob <sha>` once per tree entry. This Task replaces
the loop's per-file spawns with one `git cat-file --batch` process per read
scope: object names stream in on stdin, contents stream back in request
order, and the subprocess count becomes proportional to restore operations
instead of file counts.

The reader is private to this package (the techspec's decision: the two
loops own different error vocabularies). Its correctness rail is the batch
protocol's framing — `<sha> SP <type> SP <size> LF <content> LF`, with
`missing` for unknown objects — and the existing restore digest tests,
which fail on any content drift.

## Requirements

1. MUST add a batch object reader that feeds one `cat-file --batch` process
   and returns contents in request order, reading exactly `size` bytes plus
   the trailing newline per object.
2. MUST route the per-file loop in skills restore through the reader; the
   `git cat-file` spawn disappears from the loop body.
3. MUST map a `missing` reply and a mid-stream process death to the same
   error surface the per-file loop produced, preserving the restore error
   vocabulary (`SkillsRestoreExecution`, `source.read-failed`).
4. MUST keep every existing skills-restore test passing unmodified; the
   restore digests are the characterization.
5. MUST cover in unit tests: multi-object reads in order, a missing object,
   a zero-byte blob, content containing the framing delimiters, and process
   death mid-stream.

## Subtasks

- [ ] Implement the batch reader with framing-exact reads.
- [ ] Route the restore loop through it.
- [ ] Map both failure modes onto the existing error surface.
- [ ] Add the five unit cases.

## Acceptance Criteria

- [ ] Restoring N files spawns one `cat-file` process, not N — proven by a
      test that counts spawns through the package's git runner seam.
- [ ] All existing skills-restore tests pass unmodified.
- [ ] The five framing unit cases pass.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`
      and this task file.

## Verification

- `go test ./internal/baseline -count=1 -run 'SkillsRestore|BatchObject' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0; nothing
  regressed.
- `grep -n "cat-file" internal/baseline/skills_restore_git.go | grep -qv -- "--batch" && exit 1 || exit 0`
  — expected: exit 0; no per-file cat-file spawn remains outside the batch
  path.

## References

- `_prd.md` → Core Feature 1; Goals 1, 4.
- `_techspec.md` → Implementation Design (Interfaces); Integration Points
  (the batch framing); Risks (framing bugs corrupt content silently).
- ADR-0090.

## Result

Implemented one private `git cat-file --batch` reader for each skills-restore
tree read. The reader keeps one child process open, sends object names in tree
order, parses the size-prefixed response, consumes the trailing newline, and
waits for the child on every exit path. Skills restore now maps batch start,
read, missing-object, truncated-stream, and close failures through
`SkillsRestoreExecution` with `source.read-failed`.

Focused implementation evidence by acceptance criterion:

1. One batch process for N files: the package runner-seam test
   `TestSkillsRestoreReadsManyFilesThroughOneBatchProcess` restored three
   objects, observed one `OpenBatch` call, one `ls-tree` call, and one reader
   close. The focused command
   `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/baseline -count=1 -run '^TestSkillsRestore(Read|Batch)'`
   passed 4 cases.
2. Existing skills-restore characterization: the five pre-existing top-level
   tests ran without edits via
   `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/baseline -count=1 -run '^TestSkillsRestore(OfflinePreviewApplyAndIdempotence|ProvenanceAndPreMutationRefusals|StalePlanDoesNotMutate|RollbackRestoresSkillAndLockPreimage|CompatibilityMatchesMaintainedPythonShape)$'`
   and passed 14 cases including subtests.
3. Required framing cases:
   `GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/baseline -count=1 -run '^TestBatchObjectReader(ReturnsMultipleObjectsInRequestOrder|ReportsMissingObject|ReturnsZeroByteBlob|PreservesFramingDelimitersInContent|ReportsProcessDeathMidStream)$'`
   passed all 5 cases. The restore failure table separately confirmed that a
   missing object and mid-stream death both retain the
   `SkillsRestoreExecution` / `source.read-failed` surface.
4. Changed-path scope: `git -c core.fsmonitor=false status --porcelain`
   listed only this task file plus `internal/baseline/skills_restore_git.go`
   and `internal/baseline/skills_restore_git_test.go`; `git diff --check`
   exited 0.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
