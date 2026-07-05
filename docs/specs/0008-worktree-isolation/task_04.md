---
task: task_04
spec: 0008-worktree-isolation
status: completed
type: backend
complexity: high
---

# Task 04: Run resolve and watch on the Run Worktree

## Overview

The review path gets the same isolation: resolve and watch Runs execute in
their Run Worktree, the batch commits integrate before any Final Push, and
watch reuses one worktree across its Rounds. Verifiable through the review
engines' existing suites over worktree WorkDirs plus new integration-order
tests.

## Requirements

1. MUST create the Run Worktree for resolve and watch Runs (recording
   `work_dir`) and point `CyclePlan.GitRoot` — and with it agent cwd,
   verification, snapshots, batch commits — at it; watch reuses the single
   worktree across all its Rounds.
2. MUST integrate before pushing: Final Push (and the merge-readiness
   confirm phase that follows it) runs only after a successful integration
   onto the user's branch; integration refusal ends the Run
   `IntegrationPending` with the push skipped and the command printed —
   the pushed branch must never be ahead of an unintegrated local one.
3. MUST keep fetch Runs untouched (no Agent, no worktree) and Compatible
   Artifacts reads/writes on the shared Artifact Directory exactly as
   today.
4. MUST apply the same outcome lifecycle as implement: integrated Clean
   cleans up; every other outcome keeps the worktree and prints the path.
5. MUST leave every review contract byte-stable except the documented new
   report lines and outcome.

## Subtasks

- [x] Worktree creation and wiring in resolve and watch flows
- [x] Integration-before-push ordering incl. the confirm phase
- [x] Outcome lifecycle parity with implement
- [x] Review-suite regression over worktree WorkDirs

## Acceptance Criteria

- [x] A resolve Run with concurrent user commits produces batch commits
      free of user files; the review suite passes over the worktree
      WorkDir without behavioral edits.
- [x] Ordering test: push is invoked only after integration succeeds; the
      refusal fixture ends IntegrationPending with zero push calls.
- [x] Watch executes multiple Rounds in one worktree and cleans up only on
      integrated Clean.
- [x] Fetch behavior byte-identical; full suite passes.

## Verification

- `rtk go test ./internal/cli/ ./internal/daemon/ ./internal/watch/` —
  expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 3, 4; Core Features 2, 3. `_techspec.md` → Run
flows, Risks (watch rounds), Build Order 4. ADR-0023, ADR-0024, ADR-0019
(confirm phase ordering).

## Result

Status: completed.

Acceptance evidence:

- Resolve isolation: `TestRunResolveWorktreeIsolationExcludesConcurrentUserCommit` runs resolve in a real Run Worktree while committing `user.txt` in the user checkout. The test asserts the Run Branch commit contains `agent.txt`, excludes `user.txt`, the user branch keeps its commit, the Agent ran with `GitRoot == work_dir`, and the refusal ended `IntegrationPending`.
- Integration-before-push: `TestRunResolvePushRunsAfterSuccessfulIntegrationAndCleansWorktree` uses a pusher that fails unless the user branch already contains the integrated agent commit. It also asserts the push runs from the Run Worktree and that Clean removes the worktree and Run Branch. The concurrent-user refusal fixture asserts zero push calls and prints `git merge --ff-only roundfix/run-<id>`.
- Watch lifecycle: `TestRunWatchReusesOneRealWorktreeAcrossRoundsAndCleansOnIntegratedClean` drives two watched Rounds through the same recorded worktree, verifies Final Push runs after each integrated clean Round, and asserts terminal Clean removes the worktree and Run Branch.
- Fetch unchanged: existing fetch assertions in `TestRunOperationalCommandAcceptsMVPFlags/fetch` still pass with the same no-Agent/no-commit/no-push contract; no fetch worktree path is created or printed.

Verification:

- `rtk go test ./internal/cli/ ./internal/daemon/ ./internal/watch/` passed: 307 tests in 3 packages.
- `rtk go test ./...` passed: 678 tests in 17 packages.
