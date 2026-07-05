---
task: task_03
spec: 0008-worktree-isolation
status: pending
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

- [ ] Worktree creation and work_dir recording in the implement flow
- [ ] Engine and agent wiring to the worktree path
- [ ] Preflight demotion and debris sweep
- [ ] Outcome integration matrix incl. IntegrationPending state
- [ ] Settle retarget with integration
- [ ] End-to-end tests with concurrent user edits and commits

## Acceptance Criteria

- [ ] Multi-writer proof: a Run executing while the test mutates and
      commits in the user checkout produces task commits containing zero
      user files, and the user's commits survive untouched.
- [ ] Outcome matrix: Clean+ff (worktree gone, branch advanced, user dirt
      preserved when non-overlapping); Clean+overlap → IntegrationPending
      exit 1 with the branch unmoved and the printed command working when
      executed by the test; Unresolved keeps the worktree and prints the
      path.
- [ ] A dirty user tree no longer blocks implement (note asserted on
      stderr).
- [ ] Settle over a kept worktree verifies, commits, integrates, and
      cleans up on success.
- [ ] Full suite passes with only the documented deliberate updates.

## Verification

- `rtk go test ./internal/cli/ ./internal/daemon/` — expected: all tests
  pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1–5; Core Features 2–4. `_techspec.md` → Run
flows, API Contracts, Build Order 3. ADR-0023, ADR-0024.
