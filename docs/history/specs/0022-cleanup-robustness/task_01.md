---
task: task_01
spec: 0022-cleanup-robustness
status: completed
type: backend
complexity: medium
---

# Task 01: Forced worktree removal and warn-and-continue clean path

## Overview

Make Roundfix-owned worktree removal survive bootstrap debris, and stop a
post-integration cleanup failure from failing a Clean Run. Demoable with a
worktree full of untracked files.

## Requirements

1. MUST use `git worktree remove --force` at every Roundfix-owned removal
   site: Clean Run cleanup, Task Worktree cleanup, force-stop reap, and the
   preflight sweep.
2. MUST downgrade a cleanup failure occurring after successful integration
   to: one stderr warning shaped like
   `roundfix: Run Worktree cleanup failed; kept <path>: <reason>`, one
   Daemon-source Run Event, and an unchanged outcome, report, and exit code.
3. MUST keep pre-integration failure handling and non-integrated outcomes
   (kept worktrees) byte-for-byte unchanged.

## Subtasks

- [x] `--force` on all removal call sites
- [x] Warn-and-continue on the post-integration clean paths (implement,
      resolve/watch)
- [x] Journal event for the kept worktree
- [x] Tests: removal with untracked debris; failing cleanup after Clean
      leaves report and exit code identical and emits warning + event

## Acceptance Criteria

- [x] A worktree containing untracked files and a nested directory is
      removed successfully on Clean cleanup.
- [x] With a cleanup failure injected after integration, the command's
      stdout report and exit code are byte-identical to the success case,
      stderr carries the warning naming the kept path, and the journal has
      the event.
- [x] Existing kept-worktree behavior for non-integrated outcomes is
      unchanged.

## Verification

- `rtk go test ./internal/worktree/ ./internal/cli/` — expected: all tests
  pass, including the new debris and warn-and-continue tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1-2. `_techspec.md` → System Architecture; API
Contracts; Build Order 1; Risks.

## Result

Implemented forced removal for Roundfix-owned Run Worktree, Task Worktree,
and terminal Worktree reap paths by adding `--force` to each
`git worktree remove` call. Clean post-integration cleanup failures in
Implement, Resolve, and Watch now warn, journal a Daemon status event, and
continue to the existing Clean completion/report path.

Evidence:

- `rtk go test ./internal/worktree/ ./internal/cli/`: passed, 348 tests in
  2 packages. Before the production fix, the new debris and cleanup-failure
  tests failed against the old behavior.
- `rtk make verify`: passed. It ran `rtk go test ./...` with 887 tests in
  19 packages, `rtk go run -buildvcs=false ./cmd/roundfix skills check`
  with all bundled skills passing, and
  `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`.

Acceptance criteria:

- Debris removal: `TestCleanupCleanRemovesRunWorktreeWithUntrackedDebris`
  passed with an untracked file plus nested `node_modules/cache` directory.
  `TestCleanupTaskRemovesTaskWorktreeWithUntrackedDebris` and
  `TestPruneTerminalReapsOnlyEmptyTerminalRunAndTaskBranches` cover the
  Task Worktree and terminal reap removal sites with the same debris shape.
- Post-integration cleanup failure: `TestRunImplementCleanCleanupFailureWarnsAndJournalsWithoutChangingReportOrExit`,
  `TestRunCleanCleanupFailureWarnsAndJournalsWithoutChangingReportOrExit`,
  and `TestRunWatchCleanCleanupFailureWarnsAndJournalsWithoutChangingReportOrExit`
  passed. Each compares stdout and exit code against a Clean success case,
  checks exactly one warning naming the kept path, verifies the Run remains
  Clean, and finds the Daemon-source journal event with the path and reason.
- Kept-worktree behavior for non-integrated outcomes stayed covered by the
  unchanged tests in the same focused gate, including
  `TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks`,
  `TestRunImplementUnresolvedKeepsRealRunWorktreeAndPrintsPath`,
  `TestRunImplementFailedTaskEndsUnresolvedAndKeepsWorktree`,
  `TestRunResolveWorktreeIsolationExcludesConcurrentUserCommit`, and
  `TestIntegrateTaskReturnsConflictAndLeavesRunBranchUnmoved`.
