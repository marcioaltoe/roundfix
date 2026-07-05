---
task: task_03
spec: 0008-worktree-isolation
status: completed
type: backend
complexity: high
---

# Task 03: Run implement on the Run Worktree

## Overview

The implement path moves into isolation: preflight keeps its user-checkout
checks (with the dirty-tree rejection demoted to a note), the Run creates
its worktree and executes everything there, outcomes integrate through the
worktree package — introducing the Integration Pending terminal state — and
the Settle Command retargets kept worktrees. Verifiable through end-to-end
CLI tests with concurrent user activity.

## Requirements

1. MUST create the Run Worktree after Run creation (recording `work_dir`)
   and point `TaskPlan.WorkDir`, the Agent Session, prompts, verification,
   snapshots, commits, task-file reads/writes, and the QA step at it — the
   engines themselves unchanged.
2. MUST demote implement's dirty-tree rejection to one stderr note stating
   that overlapping local changes end the Run Integration Pending; the
   default-branch veto, Spec validation, Active Run locks (still keyed on
   the user checkout), and the debris sweep (PruneTerminal) run in
   preflight as before.
3. MUST wire outcomes: Clean → Integrate → success removes worktree and Run
   Branch (then auto-push when configured); integration refusal → the Run
   ends `IntegrationPending` (new terminal state, exit 1, outcome line
   naming the exact `git merge --ff-only roundfix/run-<id>` command);
   Unresolved/Failed/Stopped keep the worktree and print its path.
4. MUST retarget `roundfix settle`: resolve the Run's kept worktree via
   `work_dir`, verify and commit there, then run the same integration
   protocol; the stage-everything contract now scopes to the Run Worktree.
5. MUST keep the final stdout task lines reading from where execution
   happened, and every existing contract not named here byte-stable.

## Subtasks

- [x] Worktree creation and work_dir recording in the implement flow
- [x] Engine and agent wiring to the worktree path
- [x] Preflight demotion and debris sweep
- [x] Outcome integration matrix incl. IntegrationPending state
- [x] Settle retarget with integration
- [x] End-to-end tests with concurrent user edits and commits

## Acceptance Criteria

- [x] Multi-writer proof: a Run executing while the test mutates and
      commits in the user checkout produces task commits containing zero
      user files, and the user's commits survive untouched.
- [x] Outcome matrix: Clean+ff (worktree gone, branch advanced, user dirt
      preserved when non-overlapping); Clean+overlap → IntegrationPending
      exit 1 with the branch unmoved and the printed command working when
      executed by the test; Unresolved keeps the worktree and prints the
      path.
- [x] A dirty user tree no longer blocks implement (note asserted on
      stderr).
- [x] Settle over a kept worktree verifies, commits, integrates, and
      cleans up on success.
- [x] Full suite passes with only the documented deliberate updates.

## Verification

- `rtk go test ./internal/cli/ ./internal/daemon/` — expected: all tests
  pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1–5; Core Features 2–4. `_techspec.md` → Run
flows, API Contracts, Build Order 3. ADR-0023, ADR-0024.

## Result

Status: completed.

Acceptance evidence:

- Multi-writer proof: `TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit` runs implement in a real Run Worktree while committing `user.txt` in the user checkout. The test asserts the Run Branch task commit contains `agent.txt` and excludes `user.txt`, while the user branch still has the user commit and no unintegrated agent file.
- Outcome matrix: `TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt` proves Clean fast-forward cleanup and non-overlapping user dirt preservation; `TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks` proves overlap becomes `IntegrationPending`, leaves the branch unmoved with unstaged user dirt intact, prints `git merge --ff-only roundfix/run-<id>`, and that command fast-forwards after the test clears the overlap; `TestRunImplementUnresolvedKeepsRealRunWorktreeAndPrintsPath` proves Unresolved keeps and prints the Run Worktree.
- Dirty user tree demotion: `TestRunImplementDirtyWorkingTreePrintsNoteAndRuns` asserts the stderr note and a successful implement run instead of preflight rejection.
- Settle retarget: `TestRunSettleRetargetsKeptRunWorktreeAndCleansUpAfterIntegration` creates an unresolved kept Run Worktree, then `roundfix settle` verifies and commits there, integrates into the user checkout, deletes the Run Worktree and Run Branch, and marks the Run Clean.
- Deliberate test contract updates were limited to the new isolation behavior: dirty-tree preflight is now a note, unresolved runs are kept worktrees rather than user-checkout leftovers, and settle help names the Run Worktree scope.

Verification:

- `rtk go test ./internal/cli/ ./internal/daemon/` passed: 289 tests in 2 packages.
- `rtk go test -race ./internal/daemon/` passed: 53 tests in 1 package.
- `rtk go test ./...` passed: 675 tests in 17 packages.
