---
task: task_02
spec: 0263-archive-and-cancel-keep-what-git-ignores
status: pending
type: backend
complexity: medium
---

# Task 02: A failed folder removal keeps what was written and prints how to finish

## Overview

When the final `os.RemoveAll` of the Spec folder fails, both commands have
already written the record and every other output, yet archive reports "No
side effects" with exit 2 and neither command names the step that finishes
(`_techspec.md` → Current behavior). This Task makes `Archive` and `Cancel`
return a typed folder removal error with their filled result, and makes both
commands keep every write, exit 1 and print the record and the exact
finishing command. It is verifiable on its own with a read-only subdirectory
in a temporary repository.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   [archive and cancel remove ignored files in the Spec folder](references/2026-10-12-archive-and-cancel-remove-ignored-files-in-the-spec-folder.md),
   for its failed-removal half.
2. MUST add `FolderRemovalError` in the new `internal/spec/folder_removal.go`
   following `_techspec.md` → Interfaces, and MUST make `Archive` and `Cancel`
   return it, wrapping the `os.RemoveAll` error, together with their filled
   result, following API Contract 2. Nothing written before the removal is
   rolled back.
3. MUST handle that error in both commands following API Contract 3 and
   Surface Transcripts 3 and 4: the Archive Command first records
   `--pull-request` when given; both print nothing on standard output, exactly
   the two standard-error lines, and exit 1. `cancelWriteError` MUST stop
   matching the old `remove cancelled Spec;` text, and write failures before
   the removal MUST keep their current messages and exit codes.
4. MUST add `folderRemovalFinishCommand` and the recovery report in the new
   `internal/cli/folder_removal_recovery.go`, returning
   `chmod -R u+w '<folder>' && rm -rf '<folder>'` for the absolute folder,
   single-quoted with each embedded `'` written as `'\''`.
5. MUST change the help exit-code lines of API Contract 4: archive gains
   `  1  archive could not finish after writing; stderr says how to finish`
   between `0` and `2`, and cancel's `1` reads
   `cancel write or removal failed`.
6. MUST add the tests named in Verification, each calling `t.Parallel()`
   first. `TestArchiveFolderRemovalFailureRecoversForward` asserts every line
   of Surface Transcript 3 and `TestCancelFolderRemovalFailureRecoversForward`
   every line of Surface Transcript 4, each with only the declared variations;
   both then run the printed finishing command with `sh -c` and assert the tree
   a successful run leaves. The failure comes from a committed subdirectory
   made read-only with `chmod 0555` and restored in `t.Cleanup`.
   `TestCancelTranscriptPrintsHelp`, `TestRunArchiveHelp` and
   `TestCancelWriteFailureUsesExitOne` MUST pass with the new help and the
   unchanged write-failure prefix.
7. MUST describe the recovery in `docs/user-guide/commands/archive.md` and
   `docs/user-guide/commands/cancel.md`, each containing `chmod -R u+w`, add
   `internal/cli/folder_removal_recovery.go` to the sources of the
   `command archive` and `command cancel` entries of
   `docs/references/skill-coverage.json`, and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add the typed error and return it with the filled result from both writers.
- [ ] Print the recovery lines and record the pull request in archive.
- [ ] Add the finishing command formatter and the help exit lines.
- [ ] Add the recovery tests in both packages.
- [ ] Update both user guides and the Behavior Surface Record.

## Acceptance Criteria

- [ ] With a read-only subdirectory, archive and cancel exit 1, keep the record and every other write, and print the finishing command.
- [ ] Running the printed command leaves the tree a successful archive or cancellation leaves.
- [ ] Archive with `--pull-request` records the number even when the removal fails.
- [ ] A folder whose name holds a space and a quote yields a command a POSIX shell runs correctly.

## Consumers

- `spec.Archive` and `spec.Cancel` results and errors — consumers: `./internal/spec`, `./internal/cli`
- Archive Command exit code, read by the Delivery Queue archive step — consumers: `./internal/cli`
- Behavior Surface Record of the `archive` and `cancel` help and guides — consumers: `./internal/docscontract`

## Context

- instruction: `docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`
- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- creates: `internal/spec/folder_removal.go`
- creates: `internal/spec/folder_removal_test.go`
- creates: `internal/cli/folder_removal_recovery_test.go`
- interface: `internal/spec/archive.go`
- interface: `internal/spec/cancel.go`
- creates: `internal/cli/folder_removal_recovery.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/cancel.go`
- interface: `internal/cli/archive_test.go`
- interface: `internal/cli/cancel_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/cancel.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `for f in docs/user-guide/commands/archive.md docs/user-guide/commands/cancel.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'chmod -R u+w' || { printf 'missing finishing command in %s\n' "$f" >&2; exit 1; }; done; grep -qF -- 'internal/cli/folder_removal_recovery.go' docs/references/skill-coverage.json || { printf 'coverage map lacks the recovery source\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task neither guide names the finishing command and the map lacks the recovery source, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestArchiveReturnsAFolderRemovalError|TestCancelReturnsAFolderRemovalError)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveReturnsAFolderRemovalError TestCancelReturnsAFolderRemovalError; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Archive|Cancel' ./internal/spec || exit 1` — expected: exit 0. Before this Task the two tests do not exist, so the PASS check fails; afterwards they and every archive and cancel test of the package pass.
- `out="$(go test -count=1 -v -run '^(TestArchiveFolderRemovalFailureRecoversForward|TestCancelFolderRemovalFailureRecoversForward|TestArchiveFolderRemovalFailureRecordsThePullRequest|TestFolderRemovalFinishCommandQuotesTheFolder)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveFolderRemovalFailureRecoversForward TestCancelFolderRemovalFailureRecoversForward TestArchiveFolderRemovalFailureRecordsThePullRequest TestFolderRemovalFinishCommandQuotesTheFolder; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Archive|Cancel' ./internal/cli || exit 1` — expected: exit 0. Before this Task the four tests do not exist, so the PASS check fails; afterwards they and every archive, cancel and Delivery Queue archive test pass.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 2; Core Feature 3; Success Metric 2; Success Metric 3
- `_techspec.md` → API Contract 2; API Contract 3; API Contract 4; Surface Transcript 3; Surface Transcript 4; Interfaces; Invariants 3 and 4; Testing Approach; Build Order 2
- ADR-0269; ADR-0261; ADR-0260
