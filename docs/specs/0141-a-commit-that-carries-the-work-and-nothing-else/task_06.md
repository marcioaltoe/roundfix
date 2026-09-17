---
task: task_06
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: completed
type: backend
complexity: medium
---

# Task 06: Close the two gaps the review found

## Overview

Independent review found two ways the new rules read a path more coarsely than
they promise. The index question matches a directory prefix, so an untracked
executable that replaced a tracked directory is treated as tracked. And the
worktree snapshot collapses a new untracked directory to the directory itself,
so excluding the Verification's window can discard a file the QA Agent wrote
into the same new directory. This slice closes both, and it is the second and
final corrective Task the contract allows.

## Requirements

1. MUST treat a path as tracked only when the index holds that exact path, never
   when the path merely names a directory whose contents are tracked.
2. MUST keep refusing an untracked executable that replaced a tracked directory,
   with the reason and mode the filter reports today.
3. MUST exclude from the QA Report commit only what the repository Verification
   actually wrote, so a file the QA Agent writes into a directory the
   Verification created still reaches the commit.
4. MUST keep the report the mechanical stage seeded and its evidence in the
   commit, and leave a Spec whose Verification writes nothing unchanged.
5. MUST NOT change the shared snapshot contract for other callers unless the
   change is proved by their own tests.

## Subtasks

- [x] Make the index question exact.
- [x] Cover the directory-replacement case at the filter's unit seam.
- [x] Make the Verification window compare files rather than collapsed
      directories.
- [x] Cover an Agent file written into a directory the Verification created.

## Acceptance Criteria

- [x] An untracked executable that replaced a tracked directory is refused with
      its mode.
- [x] A tracked executable is still kept.
- [x] A QA Report commit keeps an Agent file written into a directory the
      repository Verification created, and still excludes the Verification's own
      file in that directory.
- [x] A Spec whose Verification writes nothing produces the same commit as
      before this Task.

## Context

- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/daemon.go`

## Verification

- `out="$(go test -count=1 -run "^TestFilterStageablePathsRefusesExecutableReplacingATrackedDirectory$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestQAReportCommitKeepsAgentEvidenceInANewDirectory$" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `grep -q "TestFilterStageablePathsRefusesExecutableReplacingATrackedDirectory" internal/daemon/task_engine_test.go && go test -count=1 -run "^TestFilterStageablePaths" ./internal/daemon` — expected: exit 0; every filter case, including the tracked one, holds after the question becomes exact. Before this Task the directory-replacement case does not exist, so the command fails.

## References

`_prd.md` → Core Features 1 and 3; Goals 1 and 3; Success Metrics 1 and 3;
Regression locks;
`_techspec.md` → Implementation Design: A tracked path is stageable, The QA
Report commit excludes the Verification's writes; API Contracts 2-3; ADR-0157.

## Result

Implementation:

- Made the index lookup exact by reading NUL-delimited `git ls-files` results
  and accepting only the requested path, never an indexed descendant beneath a
  directory prefix. The same exact answer now governs the absent-path check.
- Added a QA-only file-granularity snapshot path using
  `--untracked-files=all`. The shared `WorktreeSnapshotter.Snapshot` contract
  and its default porcelain invocation remain unchanged for other callers.
- Added regressions for an executable that replaces a tracked directory and
  for a QA Agent file written beside a Verification-created file in the same
  new directory.

Focused checks and acceptance evidence:

- Red signal before the production change: the combined new-test run failed
  because `tracked-directory` was kept and `qa-output/agent.txt` was absent
  from the QA Report commit.
- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk proxy go test -count=1 -run '^(TestFilterStageablePathsRefusesExecutableReplacingATrackedDirectory|TestQAReportCommitKeepsAgentEvidenceInANewDirectory)$' ./internal/daemon` — exit 0 after the change. The filter refuses the replacement with mode `0755`; the QA commit contains the Agent file, excludes the Verification file and leaves that file uncommitted.
- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk proxy go test -count=1 -run '^(TestFilterStageablePaths|TestQAReportCommit)' ./internal/daemon` — exit 0. This includes the existing tracked-executable case, the existing QA report/evidence cases and the no-Verification-write case.
- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk go test -count=1 ./internal/daemon` — exit 0 with 344 tests, covering the unchanged shared snapshot behavior beside the QA-specific path.
- `GOCACHE=/private/tmp/roundfix-task06-gocache rtk make verify-incremental` — the sandboxed attempt reached all packages but two existing force-stop integration tests could not read the process table. The permitted rerun exited 0, including all Go packages, skill checks and the build.
- The authored `## Verification` commands were intentionally left for the
  Daemon.

Follow-up:

- The PRD, TechSpec and this Task cite ADR-0157, but this worktree contains no
  `docs/adr/0157-*.md`. Their embedded tracked-path decision is consistent and
  was implemented; restoring or correcting the missing ADR artifact is outside
  this Task's slice.
