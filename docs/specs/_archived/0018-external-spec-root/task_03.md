---
task: task_03
spec: 0018-external-spec-root
status: completed
type: backend
complexity: medium
---

# Task 03: Commit boundary: staging guard and settle without a commit

## Overview

Make the Daemon's commit boundary degrade instead of dying: paths that are
external to the repository or cross a symbolic link are dropped from staging
with a journaled warning, and a Task whose only changes are external
artifacts settles completed without a commit. This is the fix for the field
failure where the daemon died on `git add` after verification had already
passed and the Run stayed dead unattended.

## Requirements

1. MUST filter commit staging paths, dropping any path that resolves outside
   the repository working tree or crosses a symbolic link, and staging the
   rest.
2. MUST journal one Daemon event per dropped path with its reason and print a
   progress warning shaped like
   `roundfix: task file <path> kept outside the repository; committed without it`.
3. MUST settle the Task `completed` without creating a commit when nothing
   stageable changed, still publishing the normal settled event — never fail
   the Task or the Run because artifacts live elsewhere.
4. MUST apply the symbolic-link guard unconditionally — an un-configured
   symlinked layout degrades to "artifact not committed", never to a failed
   commit.
5. MUST leave QA Report commits under the same rule: an external QA Report is
   not staged and the QA step proceeds.
6. MUST keep internal-root staging behavior unchanged, including the settled
   task file riding in its Task commit.

## Subtasks

- [x] Staging filter with external and symlink-crossing detection
- [x] Journal events and progress warnings for dropped paths
- [x] Settle-without-commit path when nothing stageable changed
- [x] QA Report commit under the same rule
- [x] Daemon tests: mixed paths, only-external changes, symlinked task file,
      unchanged internal behavior

## Acceptance Criteria

- [x] A Task commit in a repository whose task file path crosses a symlink
      succeeds, contains only repository paths, and the journal names the
      dropped path with its reason.
- [x] A Task whose only change is its external task file settles `completed`
      with no commit and the standard settled event.
- [x] With an internal root, the task file still rides in its Task commit and
      existing daemon tests pass unchanged.
- [x] The QA step proceeds when the QA Report is external, without staging
      it.

## Verification

- `rtk go test ./internal/daemon/` — expected: all tests pass, including the
  new staging guard and settle-without-commit tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1, 3; Core Features 4-5. `_techspec.md` →
Interfaces: filterStageablePaths; API Contracts (journal and warning shapes);
Build Order 3. ADR-0035.

## Result

- Added a Daemon staging guard at the Task and QA commit boundaries. It keeps
  repository-local paths, drops external paths and paths crossing a symbolic
  link, journals one `daemon.commit` event per dropped path with
  `decision: dropped` and a reason, and prints the required progress warning.
- A completed Task whose only changed artifact is external now settles without
  calling the committer; the normal settled Task event is still published.
- QA Report commits use the same guard, so an external QA Report records a
  verdict and proceeds without staging the report.
- Internal-root behavior remains covered by the existing daemon commit tests;
  the settled task file still rides in the Task commit.
- Evidence:
  - `TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths`
    covers the symlink-crossing task file with a repository path staged.
  - `TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged`
    covers external-only Task settlement with no commit.
  - `TestTaskCycleQAReportExternalProceedsWithoutStaging` covers external QA
    Report handling.
  - `TestTaskCycleCommitStagesSnapshotDiffPlusTaskFile` and
    `TestTaskCycleRealRepoCommitsPerTaskExcludingPreexistingDirt` continue to
    cover unchanged internal-root staging.
- Verification passed:
  - `rtk go test ./internal/daemon/` — 62 tests passed in 1 package.
  - `rtk make verify` — `rtk go test ./...` passed 859 tests in 18 packages;
    `roundfix skills check` passed; build completed.
