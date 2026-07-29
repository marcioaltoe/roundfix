---
task: task_03
spec: 0054-tooling-task-and-verification-hygiene
status: completed
type: backend
complexity: medium
---

# Task 03: Refuse executable files in Daemon commits

## Overview

A build artifact is never legitimate Work Item output, but the Daemon's
snapshot-diff staging has no file-mode filter, so a stray binary in the
worktree reached a Task commit once already. Drop executable files at the
same seam that already drops repository-external and symlink-crossing paths,
so every commit surface — Batch, Task, and QA — is covered by one guard.

## Requirements

1. MUST drop a changed path from staging when it is a regular file carrying
   any execute permission, reporting the refusal with the path and its mode.
2. MUST apply to every commit the Daemon creates, by attaching to the shared
   stageable-path filter rather than to one call site.
3. MUST keep the existing drop classes and their reporting behavior
   unchanged, and MUST keep committing the remaining paths exactly as today.
4. MUST preserve the contract that an explicitly selected tracked path is
   still staged even when an ignore rule matches it.
5. SHOULD word the refusal so an operator can tell a build artifact from a
   deliberately executable repository file.

## Subtasks

- [ ] Add the executable-mode drop to the shared stageable-path filter.
- [ ] Report the refused path and mode through the existing dropped-path
      channel.
- [ ] Cover the guard for Task, Batch, and QA commit paths.

## Acceptance Criteria

- [ ] A Task whose worktree gained an executable file commits its other
      changes and reports the executable as refused with its path and mode.
- [ ] The repository-external and symlink-crossing drops keep their current
      behavior and wording.
- [ ] An explicitly selected tracked path matched by an ignore rule is still
      staged.
- [ ] A Task with no executable change produces a byte-identical commit to
      today.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/daemon_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/daemon/` — expected: pass, including the new executable-drop coverage and the unchanged ignore-override contract.

## References

`_prd.md` → User Story 6, Core Feature 6; `_techspec.md` → Build Order 4,
Interfaces: executable staging drop; ADR-0045.

## Result

### Implementation

- `FilterStageablePaths` now drops a worktree path only when `Lstat` proves it
  is a regular file and any owner, group, or other execute bit is set. The
  external-path and symlink-crossing checks still run first with their existing
  reasons and reporting.
- Executable drops reuse `DroppedStagePath` and the existing dropped-commit
  event channel. Their payload includes `path`, reason `executable file`, and
  the numeric permission `mode`; progress and event summaries state that both
  build artifacts and deliberately executable repository files are invalid
  Work Item output.
- Task and QA commits retain their existing shared filter calls. Batch commits
  now pass each snapshot-diff path through the same filter while preserving
  the original path representation and order for every retained path.
- Commits still proceed with all retained paths. `GitCommitter` was unchanged,
  so an explicitly selected tracked path continues to use `git add -f` even
  when an ignore rule matches it.

### Focused checks

- Red signal: the new focused Task, Batch, QA, and execute-bit tests initially
  failed to compile because `DroppedStagePath` had no `Mode` field. The first
  attempt reached the host Go cache and was denied by the sandbox; rerunning
  with the Spec's repository-local `.gocache` produced the intended product
  failure.
- The focused executable selection passed 7 tests:
  `TestFilterStageablePathsDropsRegularFileWithAnyExecutePermission`,
  `TestTaskCommitDropsExecutableFileAndCommitsRemainingPaths`,
  `TestResolveCycleDropsExecutableFileAndCommitsRemainingBatchPaths`, and
  `TestQACommitDropsExecutableFileAndCommitsRemainingPaths`.
- A combined focused regression selection passed 19 tests. It included the
  executable selection plus
  `TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths`,
  `TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged`,
  `TestTaskCycleQAReportExternalProceedsWithoutStaging`,
  `TestGitCommitterStagesSelectedTrackedPathMatchedByGlobalIgnore`,
  `TestTaskCycleRealRepoCommitsPerTaskExcludingPreexistingDirt`,
  `TestResolveCycleStagesOnlyAgentTouchedPaths`, and
  `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport`.
- `rtk git diff --check` passed. Final changed-path inspection reported only
  `internal/daemon/engine.go`, `internal/daemon/engine_test.go`,
  `internal/daemon/task_engine.go`, `internal/daemon/task_engine_test.go`, and
  this Task file.

### Acceptance evidence

- **Task executable refusal:** the focused Task commit test observed
  `bin/roundfix` omitted with reason `executable file` and mode `0755`, while
  the Task file and ordinary source change remained in the commit request.
  Owner-, group-, and other-only execute-bit cases were each rejected.
- **Existing drop behavior:** the focused symlink, external Task file, and
  external QA Report regressions passed with their existing reason and warning
  assertions unchanged.
- **Ignore override:** the real-Git ignored tracked-path regression passed;
  the selected tracked file was committed and the unselected ignored file
  remained outside the commit.
- **No-executable commit stability:** the existing real-repository Task test
  passed with its exact commit messages, trailers, file sets, ordering, and
  pre-existing-dirt exclusion unchanged. The ordinary Batch and QA regression
  selections also retained their exact staged paths.
- **All Daemon commit surfaces:** focused Task, Batch, and QA tests each
  observed the executable drop through the shared filter and a commit of the
  remaining paths.

### Follow-up

- The Task's declared `## Verification` commands were not run; the Daemon owns
  them and terminal settlement.
