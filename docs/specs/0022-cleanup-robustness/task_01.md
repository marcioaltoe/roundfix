---
task: task_01
spec: 0022-cleanup-robustness
status: pending
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

- [ ] `--force` on all removal call sites
- [ ] Warn-and-continue on the post-integration clean paths (implement,
      resolve/watch)
- [ ] Journal event for the kept worktree
- [ ] Tests: removal with untracked debris; failing cleanup after Clean
      leaves report and exit code identical and emits warning + event

## Acceptance Criteria

- [ ] A worktree containing untracked files and a nested directory is
      removed successfully on Clean cleanup.
- [ ] With a cleanup failure injected after integration, the command's
      stdout report and exit code are byte-identical to the success case,
      stderr carries the warning naming the kept path, and the journal has
      the event.
- [ ] Existing kept-worktree behavior for non-integrated outcomes is
      unchanged.

## Verification

- `rtk go test ./internal/worktree/ ./internal/cli/` — expected: all tests
  pass, including the new debris and warn-and-continue tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1-2. `_techspec.md` → System Architecture; API
Contracts; Build Order 1; Risks.
